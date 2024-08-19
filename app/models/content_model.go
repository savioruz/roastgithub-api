package models

type Language string

const (
	LangAuto Language = "auto"
	LangID   Language = "id"
	LangEN   Language = "en"
)

type ContentRequest struct {
	Username string   `json:"username" validate:"required,min=6,max=32"`
	Lang     Language `json:"lang" validate:"required,oneof=auto id en"`
}

type ContentResponse struct {
	Username         string `json:"username"`
	AvatarURL        string `json:"avatar_url"`
	GeneratedContent string `json:"generated_content"`
}

type ContentResponseSuccess struct {
	Data ContentResponse `json:"data"`
}

type ContentResponseFailure struct {
	Error string `json:"error"`
}
