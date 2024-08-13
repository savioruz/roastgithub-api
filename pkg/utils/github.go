package utils

import (
	"context"
	"fmt"
	"github.com/google/go-github/v63/github"
	"os"
	"roastgithub-api/app/models"
)

type GithubService struct {
	client *github.Client
}

func NewGithubService() *GithubService {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		panic("consider setting GITHUB_TOKEN environment variable")
	}

	client := github.NewClient(nil).WithAuthToken(token)

	return &GithubService{
		client: client,
	}
}

// GetUserProfile func get user profile by username
func (s *GithubService) GetUserProfile(ctx context.Context, username string) (*models.ProfileResponse, error) {
	user, _, err := s.client.Users.Get(ctx, username)
	if err != nil {
		return nil, err
	}

	return &models.ProfileResponse{
		Name:      user.Name,
		Bio:       user.Bio,
		Company:   user.Company,
		Location:  user.Location,
		Followers: user.Followers,
		Following: user.Following,
		Repos:     user.PublicRepos,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

// GetUserRepositories func get user repositories by username
func (s *GithubService) GetUserRepositories(ctx context.Context, username string) ([]*models.Repository, error) {
	opt := &github.RepositoryListByUserOptions{
		Type:      "owner",
		Sort:      "updated",
		Direction: "desc",
	}

	repos, _, err := s.client.Repositories.ListByUser(ctx, username, opt)
	if err != nil {
		return nil, err
	}

	var userRepos []*models.Repository
	for _, repo := range repos {
		userRepos = append(userRepos, &models.Repository{
			Name:        repo.GetName(),
			Description: repo.GetDescription(),
			Lang:        repo.GetLanguage(),
			Stars:       repo.GetStargazersCount(),
			Forks:       repo.GetForksCount(),
			Watchers:    repo.GetWatchersCount(),
			OpenIssues:  repo.GetOpenIssuesCount(),
			CreatedAt:   repo.GetCreatedAt().Format("2006-01-02"),
			UpdatedAt:   repo.GetUpdatedAt().Format("2006-01-02"),
		})
	}

	return userRepos, nil
}

// GetReadme func get README.md content by username
func (s *GithubService) GetReadme(ctx context.Context, username string) (string, error) {
	var readmeContent string
	var err error

	readmeContent, err = s.getReadmeFromBranch(ctx, username, "main")
	if err != nil {
		readmeContent, err = s.getReadmeFromBranch(ctx, username, "master")
		if err != nil {
			return "", fmt.Errorf("failed to get README.md: %v", err)
		}
	}

	return readmeContent, nil
}

// getReadmeFromBranch func get README.md content by branch
func (s *GithubService) getReadmeFromBranch(ctx context.Context, username, branch string) (string, error) {
	opt := &github.RepositoryContentGetOptions{Ref: branch}

	fileContent, _, _, err := s.client.Repositories.GetContents(ctx, username, username, "README.md", opt)
	if err != nil {
		return "", err
	}

	content, err := fileContent.GetContent()
	if err != nil {
		return "", err
	}

	return content, nil
}
