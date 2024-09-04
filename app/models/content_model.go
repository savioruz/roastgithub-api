package models

import "mime/multipart"

type Language string

const (
	LangAuto Language = "auto"
	LangID   Language = "id"
	LangEN   Language = "en"
)

type GithubRequest struct {
	Username string   `json:"username" validate:"required,min=6,max=32"`
	Key      *string  `json:"key,omitempty" validate:"omitempty"`
	Lang     Language `json:"lang" validate:"required,oneof=auto id en"`
}

type ResumeRequest struct {
	File *multipart.FileHeader `form:"file" validate:"required"`
	Key  *string               `form:"key" validate:"omitempty"`
	Lang *string               `form:"lang" validate:"required,oneof=id en"`
}

type ContentResponse struct {
	GeneratedContent string `json:"generated_content"`
}

type GithubContentResponse struct {
	Username         string `json:"username"`
	AvatarURL        string `json:"avatar_url"`
	GeneratedContent string `json:"generated_content"`
}

type GithubContentResponseSuccess struct {
	Data GithubContentResponse `json:"data"`
}

type ContentResponseSuccess struct {
	Data ContentResponse `json:"data"`
}

type ContentResponseFailure struct {
	Error string `json:"error"`
}
