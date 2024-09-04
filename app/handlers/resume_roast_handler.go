package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/savioruz/roastgithub-api/app/models"
	"github.com/savioruz/roastgithub-api/pkg/repository"
	"github.com/savioruz/roastgithub-api/pkg/utils"
	"github.com/savioruz/roastgithub-api/platform/cache"
	"os"
	"time"
)

// GetResumeRoast func get content from AI models
// @Description Get roast by resume.pdf
// @Summary get roast by resume.pdf
// @Tags Roast
// @Accept json
// @Produce json
// @Param file formData file true "Resume as PDF"
// @Param lang formData string false "Language for the content default is id"
// @Param key formData string false "Gemini API key"
// @Success 200 {object} models.ContentResponseSuccess
// @Failure 400 {object} models.ContentResponseFailure
// @Failure 500 {object} models.ContentResponseFailure
// @Router /roast/resume [post]
func GetResumeRoast(c *fiber.Ctx) error {
	var req models.ResumeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ContentResponseFailure{
			Error: "Failed to parse request body",
		})
	}

	key := c.FormValue("key", "")
	req.Key = &key
	lang := c.FormValue("lang", "id")
	req.Lang = &lang

	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ContentResponseFailure{
			Error: "Failed to get file from request",
		})
	}

	if file.Size >= (5 * 1024 * 1024) {
		return c.Status(fiber.StatusBadRequest).JSON(models.ContentResponseFailure{
			Error: "File size must be less than 5MB",
		})
	}

	if file.Header.Get("Content-Type") != "application/pdf" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ContentResponseFailure{
			Error: "File must be a PDF",
		})
	}

	req.File = file
	validate := utils.NewValidator()
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": utils.ValidatorErrors(err),
		})
	}

	path := fmt.Sprintf("./%s", file.Filename)
	s := c.SaveFile(file, path)
	if s != nil {
		log.Errorf("Failed to save file: %v", s)
		return c.Status(fiber.StatusInternalServerError).JSON(models.ContentResponseFailure{
			Error: "Failed to save file",
		})
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			log.Fatalf("Failed to remove file: %s", name)
		}
	}(path)

	pdfReader, err := utils.GetPdfText(path)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to read PDF file",
		})
	}

	ctx := context.Background()

	redisClient, err := cache.NewRedisConnection()
	if err != nil {
		log.Errorf("Failed to connect to Redis: %v", err)
	}
	defer func(redisClient *cache.RedisClient) {
		err := redisClient.Close()
		if err != nil {
			log.Errorf("Failed to close Redis connection: %v", err)
		}
	}(redisClient)

	cacheKey := fmt.Sprintf("roast-resume:content:%s:%s", *req.Lang, file.Filename)
	cachedContent, err := redisClient.Get(ctx, cacheKey)
	if err == nil && cachedContent != "" {
		var cachedData models.ContentResponse
		if err := json.Unmarshal([]byte(cachedContent), &cachedData); err == nil {
			return c.Status(fiber.StatusOK).JSON(models.ContentResponseSuccess{
				Data: cachedData,
			})
		}
	}

	var prompt string
	switch *req.Lang {
	case "id":
		prompt = fmt.Sprintf("%s resume ini. Berikut detailnya: %s", repository.BasePromptID, pdfReader)
	case "en":
		prompt = fmt.Sprintf("%s resume. Here are the details: %s", repository.BasePromptEN, pdfReader)
	default:
		prompt = fmt.Sprintf("%s resume ini. Berikut detailnya: %s", repository.BasePromptID, pdfReader)
	}

	g := utils.NewGeminiService(*req.Key)
	r, err := g.GenerateContent(ctx, prompt)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ContentResponseFailure{
			Error: "Failed to generate content",
		})
	}
	cacheData := models.ContentResponse{
		GeneratedContent: r,
	}
	cachedData, err := json.Marshal(cacheData)
	if err != nil {
		log.Errorf("Failed to marshal cache data: %v", err)
	}

	err = redisClient.Set(ctx, cacheKey, cachedData, 4*time.Hour)
	if err != nil {
		log.Errorf("Failed to set cache data: %v", err)
	}

	return c.Status(fiber.StatusOK).JSON(models.ContentResponseSuccess{
		Data: models.ContentResponse{
			GeneratedContent: r,
		},
	})
}
