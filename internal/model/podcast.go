package model

import "time"

type Language string

const (
	LanguageEnglish  Language = "en"
	LanguageSpanish  Language = "es"
	LanguageBoth     Language = "both"
)

type Level string

const (
	LevelBeginner     Level = "beginner"
	LevelIntermediate Level = "intermediate"
	LevelAdvanced     Level = "advanced"
)

type Podcast struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ImageURL    string    `json:"imageUrl"`
	Language    Language  `json:"language"`
	Level       Level     `json:"level"`
	Country     string    `json:"country"`
	Topic       string    `json:"topic"`
	URL         string    `json:"url"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type CreatePodcastInput struct {
	Title       string   `json:"title" validate:"required,min=1,max=255"`
	Description string   `json:"description" validate:"required,min=1,max=1000"`
	ImageURL    string   `json:"imageUrl" validate:"required,url"`
	Language    Language `json:"language" validate:"required,oneof=en es both"`
	Level       Level    `json:"level" validate:"required,oneof=beginner intermediate advanced"`
	Country     string   `json:"country" validate:"required,min=1,max=100"`
	Topic       string   `json:"topic" validate:"required,min=1,max=100"`
	URL         string   `json:"url" validate:"required,url"`
}

type UpdatePodcastInput struct {
	Title       *string   `json:"title,omitempty" validate:"omitempty,min=1,max=255"`
	Description *string   `json:"description,omitempty" validate:"omitempty,min=1,max=1000"`
	ImageURL    *string   `json:"imageUrl,omitempty" validate:"omitempty,url"`
	Language    *Language `json:"language,omitempty" validate:"omitempty,oneof=en es both"`
	Level       *Level    `json:"level,omitempty" validate:"omitempty,oneof=beginner intermediate advanced"`
	Country     *string   `json:"country,omitempty" validate:"omitempty,min=1,max=100"`
	Topic       *string   `json:"topic,omitempty" validate:"omitempty,min=1,max=100"`
	URL         *string   `json:"url,omitempty" validate:"omitempty,url"`
}

type PodcastFilters struct {
	Language *Language `json:"language,omitempty"`
	Level    *Level    `json:"level,omitempty"`
	Country  *string   `json:"country,omitempty"`
	Topic    *string   `json:"topic,omitempty"`
}
