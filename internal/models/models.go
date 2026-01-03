package models

import (
	"time"
)

type User struct {
	UserID         int       `json:"user_id"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"password"`
	CreatedAt      time.Time `json:"created_at"`
}

type Post struct {
	PostID    int        `json:"post_id"`
	TopicID   int        `json:"topic_id"`
	UserID    int        `json:"user_id"`
	TextTitle string     `json:"text_title"`
	TextBody  string     `json:"text_body"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	IsDeleted bool       `json:"is_deleted"`
	Username  string     `json:"username"`
}

type Comment struct {
	CommentID int        `json:"comment_id"`
	PostID    int        `json:"post_id"`
	ParentID  *int       `json:"parent_id"` //null if is a reply to post
	UserID    int        `json:"user_id"`
	TextBody  string     `json:"text_body"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	IsDeleted bool       `json:"is_deleted"`
	Username  string     `json:"username"`
}

type Topic struct {
	TopicID     int       `json:"topic_id"`
	TopicName   string    `json:"name"`
	Description *string   `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	IsDeleted   bool      `json:"is_deleted"`
}

type AddPostForm struct {
	TopicID   int     `json:"topic_id"`
	UserID    int     `json:"user_id"`
	TextTitle *string `json:"text_title"`
	TextBody  string  `json:"text_body"`
}

type AddCommentForm struct {
	PostID   int    `json:"post_id"`
	ParentID *int   `json:"parent_id"`
	UserID   int    `json:"user_id"`
	TextBody string `json:"text_body"`
}
