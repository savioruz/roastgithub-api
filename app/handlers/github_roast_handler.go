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
	"strings"
	"time"
)

// GetGithubRoast func get content from AI models
// @Description Get roast by username and data
// @Summary get roast by username and data
// @Tags Roast
// @Accept json
// @Produce json
// @Param data body models.GithubRequest true "Prompt"
// @Success 200 {object} models.GithubContentResponseSuccess
// @Failure 400 {object} models.ContentResponseFailure
// @Failure 500 {object} models.ContentResponseFailure
// @Router /roast/github [post]
func GetGithubRoast(c *fiber.Ctx) error {
	var req models.GithubRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ContentResponseFailure{
			Error: "Failed to parse request body",
		})
	}

	validate := utils.NewValidator()
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": utils.ValidatorErrors(err),
		})
	}
	key := req.Key

	ctx := context.Background()

	redisClient, err := cache.NewRedisConnection()
	if err != nil {
		log.Errorf("Couldn't connect to Redis: %v", err)
	}
	defer func(redisClient *cache.RedisClient) {
		err := redisClient.Close()
		if err != nil {
			log.Errorf("Failed to close Redis connection: %v", err)
		}
	}(redisClient)

	cacheKey := fmt.Sprintf("githubroast:content:%s:%s", req.Username, req.Lang)

	cachedContent, err := redisClient.Get(ctx, cacheKey)
	if err == nil && cachedContent != "" {
		var cachedData models.GithubContentResponse
		if err := json.Unmarshal([]byte(cachedContent), &cachedData); err == nil {
			return c.Status(fiber.StatusOK).JSON(models.GithubContentResponseSuccess{
				Data: cachedData,
			})
		}
	}

	githubService := utils.NewGithubService()

	userProfile, err := githubService.GetUserProfile(ctx, req.Username)
	if err != nil {
		log.Errorf("Failed to get user profile: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.ContentResponseFailure{
			Error: "Failed to get user profile",
		})
	}

	var language models.Language
	if strings.EqualFold(string(req.Lang), string(models.LangAuto)) {
		if userProfile.Location != nil && strings.Contains(*userProfile.Location, "Indonesia") {
			language = models.LangID
		} else {
			language = models.LangEN
		}
	} else {
		language = req.Lang
	}

	var githubData models.GithubData
	githubData.ProfileResponse = userProfile
	var readmeResponse string

	if *userProfile.Repos > 0 {
		userRepos, err := githubService.GetUserRepositories(ctx, req.Username)
		if err != nil {
			log.Errorf("Failed to get user repositories: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(models.ContentResponseFailure{
				Error: "Failed to get user repositories",
			})
		}

		var repos []models.Repository
		for _, repo := range userRepos {
			repos = append(repos, *repo)
		}

		githubData.Repositories = &repos

		if len(repos) > 0 {
			readmeResponse, _ = githubService.GetReadme(ctx, req.Username)
			if readmeResponse == "" {
				readmeResponse = ""
			}
		}
	}

	data, err := json.Marshal(githubData)
	if err != nil {
		log.Errorf("Failed to convert GitHub data to JSON: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.ContentResponseFailure{
			Error: "Failed to convert GitHub data to JSON",
		})
	}

	var prompt string
	if language == models.LangID {
		prompt = fmt.Sprintf("%s profile github berikut: %s. Detailnya seperti ini: %s\nProfile markdown: ```%s```",
			repository.BasePromptID,
			req.Username,
			data,
			readmeResponse,
		)
	} else {
		prompt = fmt.Sprintf("%s github profiles: %s. Here are the details: %s\nProfile markdown: ```%s```",
			repository.BasePromptEN,
			req.Username,
			data,
			readmeResponse,
		)
	}

	geminiService := utils.NewGeminiService(*key)

	resp, err := geminiService.GenerateContent(ctx, prompt)
	if err != nil {
		log.Errorf("Failed to generate content: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.ContentResponseFailure{
			Error: "Failed to generate content",
		})
	}

	cacheData := models.GithubContentResponse{
		Username:         *userProfile.Username,
		AvatarURL:        *userProfile.AvatarURL,
		GeneratedContent: resp,
	}

	cachedData, err := json.Marshal(cacheData)
	if err != nil {
		log.Errorf("Failed to marshal cache data: %v", err)
	}

	err = redisClient.Set(ctx, cacheKey, cachedData, 4*time.Hour)
	if err != nil {
		log.Errorf("Failed to cache content: %v", err)
	}

	return c.Status(fiber.StatusOK).JSON(models.GithubContentResponseSuccess{
		Data: models.GithubContentResponse{
			Username:         *userProfile.Username,
			AvatarURL:        *userProfile.AvatarURL,
			GeneratedContent: resp,
		},
	})
}
