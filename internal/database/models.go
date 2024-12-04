// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package database

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Voted string

const (
	VotedUp   Voted = "up"
	VotedDown Voted = "down"
)

func (e *Voted) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = Voted(s)
	case string:
		*e = Voted(s)
	default:
		return fmt.Errorf("unsupported scan type for Voted: %T", src)
	}
	return nil
}

type NullVoted struct {
	Voted Voted
	Valid bool // Valid is true if Voted is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullVoted) Scan(value interface{}) error {
	if value == nil {
		ns.Voted, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.Voted.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullVoted) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.Voted), nil
}

type Comment struct {
	ID              uuid.UUID
	Comment         string
	AuthorID        uuid.UUID
	PostID          uuid.UUID
	ParentCommentID uuid.NullUUID
	UpVoted         int32
	DownVoted       int32
	CreatedAt       time.Time
}

type CommentVote struct {
	ID        uuid.UUID
	Voted     Voted
	VoterID   uuid.UUID
	CommentID uuid.UUID
	CreatedAt time.Time
}

type Post struct {
	ID            uuid.UUID
	Title         string
	Content       string
	AuthorID      uuid.UUID
	UpVoted       int32
	DownVoted     int32
	CommentsCount int32
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type PostVote struct {
	ID        uuid.UUID
	Voted     Voted
	VoterID   uuid.UUID
	PostID    uuid.UUID
	CreatedAt time.Time
}

type User struct {
	ID           uuid.UUID
	Email        string
	Name         string
	Password     string
	Image        sql.NullString
	RefreshToken sql.NullString
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
