package utils

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
	"os"
	"strings"
)

type GeminiService struct {
	client *genai.Client
	model  *genai.GenerativeModel
}

func NewGeminiService(key string) *GeminiService {
	if key == "" {
		key = os.Getenv("GEMINI_API_KEY")
		if key == "" {
			log.Fatal("Consider setting GEMINI_API_KEY environment variable")
		}
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(key))
	if err != nil {
		log.Fatal(err)
	}

	model := client.GenerativeModel("gemini-1.5-flash")
	model.SetTemperature(1.0)
	model.SetMaxOutputTokens(180)
	model.SafetySettings = []*genai.SafetySetting{
		{
			Category:  genai.HarmCategoryDangerousContent,
			Threshold: genai.HarmBlockLowAndAbove,
		},
		{
			Category:  genai.HarmCategoryHarassment,
			Threshold: genai.HarmBlockMediumAndAbove,
		},
	}

	return &GeminiService{
		client: client,
		model:  model,
	}
}

// GenerateContent func generate content from AI models
func (s *GeminiService) GenerateContent(ctx context.Context, prompt string) (string, error) {
	resp, err := s.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", err
	}

	var result strings.Builder
	if resp != nil && resp.Candidates != nil {
		for _, c := range resp.Candidates {
			if c.Content != nil {
				for _, part := range c.Content.Parts {
					result.WriteString(fmt.Sprintf("%v", part))
				}
			}
		}
	}

	generatedContent := result.String()

	defer func() {
		if err := s.client.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	return generatedContent, nil
}
