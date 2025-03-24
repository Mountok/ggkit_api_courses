package processor

import (
	"ggkit_learn_service/internals/app/db"
	"ggkit_learn_service/internals/app/models"
)

type CommentProcessor struct {
	db *db.CommentStorage
}

func NewCommentProcessor(db *db.CommentStorage) *CommentProcessor {
	return &CommentProcessor{db: db}
}

func (p *CommentProcessor) CreateComment(comment models.Comment) error {
	return p.db.CreateComment(comment)
}

func (p *CommentProcessor) GetCommentsByThemeID(themeID int, user_id string,user_role string) ([]models.Comment, error) {
	return p.db.GetCommentsByThemeID(themeID, user_id,user_role)
}

func (p *CommentProcessor) DeleteComment(commentID int, user_id string, user_role string) error {
	return p.db.DeleteComment(commentID, user_id, user_role)
}

func (p *CommentProcessor) CreateReply(reply models.AdminReply) error {
	return p.db.CreateReply(reply)
}

func (p *CommentProcessor) GetRepliesByCommentID(commentID int) ([]models.CommentAdminReply, error) {
	return p.db.GetRepliesByCommentID(commentID)
}