package models

import (
	"time"

	"github.com/AnhBigBrother/enlighten-backend/internal/database"
	"github.com/google/uuid"
)

type Post struct {
	ID            uuid.UUID `json:"id"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	AuthorID      uuid.UUID `json:"author_id"`
	AuthorName    string    `json:"author_name"`
	AuthorEmail   string    `json:"author_email"`
	AuthorImage   string    `json:"author_image"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	UpVoted       int32     `json:"up_voted"`
	DownVoted     int32     `json:"down_voted"`
	CommentsCount int32     `json:"comments_count"`
}

func FormatDatabasePost(dbPost database.Post) Post {
	return Post{
		ID:            dbPost.ID,
		Title:         dbPost.Title,
		Content:       dbPost.Content,
		AuthorID:      dbPost.AuthorID,
		CreatedAt:     dbPost.CreatedAt,
		UpdatedAt:     dbPost.UpdatedAt,
		UpVoted:       dbPost.UpVoted,
		DownVoted:     dbPost.DownVoted,
		CommentsCount: dbPost.CommentsCount,
	}
}

func FormatDatabaseGetNewPostsRow(dbPost database.GetNewPostsRow) Post {
	return Post{
		ID:            dbPost.ID,
		Title:         dbPost.Title,
		Content:       dbPost.Content,
		AuthorID:      dbPost.AuthorID,
		AuthorName:    dbPost.AuthorName.String,
		AuthorEmail:   dbPost.AuthorEmail.String,
		AuthorImage:   dbPost.AuthorImage.String,
		UpdatedAt:     dbPost.UpdatedAt,
		UpVoted:       dbPost.UpVoted,
		DownVoted:     dbPost.DownVoted,
		CommentsCount: dbPost.CommentsCount,
		CreatedAt:     dbPost.CreatedAt,
	}
}

func FormatDatabaseGetHotPostsRow(dbPost database.GetHotPostsRow) Post {
	return Post{
		ID:            dbPost.ID,
		Title:         dbPost.Title,
		Content:       dbPost.Content,
		AuthorID:      dbPost.AuthorID,
		AuthorName:    dbPost.AuthorName.String,
		AuthorEmail:   dbPost.AuthorEmail.String,
		AuthorImage:   dbPost.AuthorImage.String,
		UpdatedAt:     dbPost.UpdatedAt,
		UpVoted:       dbPost.UpVoted,
		DownVoted:     dbPost.DownVoted,
		CommentsCount: dbPost.CommentsCount,
		CreatedAt:     dbPost.CreatedAt,
	}
}

func FormatDatabaseGetTopPostsRow(dbPost database.GetTopPostsRow) Post {
	return Post{
		ID:            dbPost.ID,
		Title:         dbPost.Title,
		Content:       dbPost.Content,
		AuthorID:      dbPost.AuthorID,
		AuthorName:    dbPost.AuthorName.String,
		AuthorEmail:   dbPost.AuthorEmail.String,
		AuthorImage:   dbPost.AuthorImage.String,
		UpdatedAt:     dbPost.UpdatedAt,
		UpVoted:       dbPost.UpVoted,
		DownVoted:     dbPost.DownVoted,
		CommentsCount: dbPost.CommentsCount,
		CreatedAt:     dbPost.CreatedAt,
	}
}

func FormatDatabaseGetPostByIdRow(dbPost database.GetPostByIdRow) Post {
	return Post{
		ID:            dbPost.ID,
		Title:         dbPost.Title,
		Content:       dbPost.Content,
		AuthorID:      dbPost.AuthorID,
		AuthorName:    dbPost.AuthorName.String,
		AuthorEmail:   dbPost.AuthorEmail.String,
		AuthorImage:   dbPost.AuthorImage.String,
		UpdatedAt:     dbPost.UpdatedAt,
		UpVoted:       dbPost.UpVoted,
		DownVoted:     dbPost.DownVoted,
		CommentsCount: dbPost.CommentsCount,
		CreatedAt:     dbPost.CreatedAt,
	}
}

func FormatDatabaseGetNewFollowedPostsRow(dbPost database.GetNewFollowedPostsRow) Post {
	return Post{
		ID:            dbPost.ID,
		Title:         dbPost.Title,
		Content:       dbPost.Content,
		AuthorID:      dbPost.AuthorID,
		AuthorName:    dbPost.AuthorName.String,
		AuthorEmail:   dbPost.AuthorEmail.String,
		AuthorImage:   dbPost.AuthorImage.String,
		UpdatedAt:     dbPost.UpdatedAt,
		UpVoted:       dbPost.UpVoted,
		DownVoted:     dbPost.DownVoted,
		CommentsCount: dbPost.CommentsCount,
		CreatedAt:     dbPost.CreatedAt,
	}
}

func FormatDatabaseGetTopFollowedPostsRow(dbPost database.GetTopFollowedPostsRow) Post {
	return Post{
		ID:            dbPost.ID,
		Title:         dbPost.Title,
		Content:       dbPost.Content,
		AuthorID:      dbPost.AuthorID,
		AuthorName:    dbPost.AuthorName.String,
		AuthorEmail:   dbPost.AuthorEmail.String,
		AuthorImage:   dbPost.AuthorImage.String,
		UpdatedAt:     dbPost.UpdatedAt,
		UpVoted:       dbPost.UpVoted,
		DownVoted:     dbPost.DownVoted,
		CommentsCount: dbPost.CommentsCount,
		CreatedAt:     dbPost.CreatedAt,
	}
}

func FormatDatabaseGetHotFollowedPostsRow(dbPost database.GetHotFollowedPostsRow) Post {
	return Post{
		ID:            dbPost.ID,
		Title:         dbPost.Title,
		Content:       dbPost.Content,
		AuthorID:      dbPost.AuthorID,
		AuthorName:    dbPost.AuthorName.String,
		AuthorEmail:   dbPost.AuthorEmail.String,
		AuthorImage:   dbPost.AuthorImage.String,
		UpdatedAt:     dbPost.UpdatedAt,
		UpVoted:       dbPost.UpVoted,
		DownVoted:     dbPost.DownVoted,
		CommentsCount: dbPost.CommentsCount,
		CreatedAt:     dbPost.CreatedAt,
	}
}
