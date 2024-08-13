package models

type Language string

const (
	LangAuto Language = "auto"
	LangID   Language = "id"
	LangEN   Language = "en"
)

type ContentRequest struct {
	Username string   `json:"username"`
	Lang     Language `json:"lang"`
}

type ContentResponse struct {
	GeneratedContent string `json:"generated_content"`
}
