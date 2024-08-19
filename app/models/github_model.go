package models

import "github.com/google/go-github/v63/github"

type GithubData struct {
	ProfileResponse *ProfileResponse `json:"profile,omitempty"`
	Repositories    *[]Repository    `json:"repositories,omitempty"`
}

type ProfileResponse struct {
	Username  *string           `json:"username,omitempty"`
	Name      *string           `json:"name,omitempty"`
	AvatarURL *string           `json:"avatar_url,omitempty"`
	Bio       *string           `json:"bio,omitempty"`
	Company   *string           `json:"company,omitempty"`
	Location  *string           `json:"location,omitempty"`
	Followers *int              `json:"followers,omitempty"`
	Following *int              `json:"following,omitempty"`
	Repos     *int              `json:"repos,omitempty"`
	CreatedAt *github.Timestamp `json:"created_at,omitempty"`
	UpdatedAt *github.Timestamp `json:"updated_at,omitempty"`
}

type Repository struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Lang        string `json:"lang,omitempty"`
	Stars       int    `json:"stars"`
	Forks       int    `json:"forks"`
	Watchers    int    `json:"watchers"`
	OpenIssues  int    `json:"open_issues"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
