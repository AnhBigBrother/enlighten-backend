package models

import (
	"time"

	"github.com/AnhBigBrother/enlighten-backend/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	Password     string    `json:"password"`
	Image        string    `json:"image"`
	Bio          string    `json:"bio"`
	RefreshToken string    `json:"refresh_token"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Author struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	Image        string    `json:"image"`
	TotalPosts   int64     `json:"total_posts"`
	TotalUpvoted int64     `json:"total_upvoted"`
}

func FormatDatabaseUser(dbUser database.User) User {
	return User{
		ID:           dbUser.ID,
		Email:        dbUser.Email,
		Name:         dbUser.Name,
		Password:     dbUser.Password,
		Image:        dbUser.Image.String,
		Bio:          dbUser.Bio.String,
		RefreshToken: dbUser.RefreshToken.String,
		CreatedAt:    dbUser.CreatedAt,
		UpdatedAt:    dbUser.UpdatedAt,
	}
}
