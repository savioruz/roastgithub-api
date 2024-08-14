package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/savioruz/roastgithub-api/app/models"
	"github.com/savioruz/roastgithub-api/pkg/repository"
	"github.com/savioruz/roastgithub-api/pkg/utils"
	"strings"
)

// GetRoast func get content from AI models
// @Description Get roast by username and data
// @Summary get roast by username and data
// @Tags Roast
// @Accept json
// @Produce json
// @Param data body models.ContentRequest true "Prompt"
// @Success 200 {object} models.ContentResponse
// @Router /roast [post]
func GetRoast(c *fiber.Ctx) error {
	var req models.ContentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	if req.Lang != models.LangAuto && req.Lang != models.LangID && req.Lang != models.LangEN {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid language request"})
	}

	ctx := context.Background()
	githubService := utils.NewGithubService()

	userProfile, err := githubService.GetUserProfile(ctx, req.Username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get user profile"})
	}

	if userProfile == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
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
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get user repositories"})
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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to convert GitHub data to JSON"})
	}

	var prompt string
	if language == models.LangID {
		prompt = fmt.Sprintf("%s %s. Detailnya seperti ini: %s\nProfile markdown: ```%s```",
			repository.BasePromptID,
			req.Username,
			data,
			readmeResponse,
		)
	} else {
		prompt = fmt.Sprintf("%s %s. Here are the details: %s\nProfile markdown: ```%s```",
			repository.BasePromptEN,
			req.Username,
			data,
			readmeResponse,
		)
	}

	geminiService := utils.NewGeminiService()

	resp, err := geminiService.GenerateContent(ctx, prompt)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": models.ContentResponse{
			GeneratedContent: resp,
		},
	})
}
