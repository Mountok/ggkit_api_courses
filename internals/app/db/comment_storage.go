package db

import (
	"context"
	"ggkit_learn_service/internals/app/models"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

type CommentStorage struct {
	db *pgxpool.Pool
}

func NewCommentStorage(db *pgxpool.Pool) *CommentStorage {
	return &CommentStorage{db: db}
}

func (s *CommentStorage) CreateComment(comment models.Comment) error {

	query := `
		INSERT INTO comments (user_id, theme_id, content)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	var id int
	row := s.db.QueryRow(context.Background(), query, comment.UserID, comment.ThemeID, comment.Content)
	err := row.Scan(&id)
	if err != nil {
		logrus.Printf("Error creating comment: %v", err)
		return err
	}
	return nil
}

func (s *CommentStorage) GetCommentsByThemeID(themeID int, user_id string, user_role string) ([]models.Comment, error) {
	query := `
	SELECT 
		c.id,
		c.user_id,
		c.theme_id,
		c.content,
		c.created_at,
		u.email
	FROM 
		comments c
	JOIN 
		users u ON u.id = c.user_id
	WHERE 
		c.theme_id = $1 AND c.user_id = $2;
	`

	queryForAll := `
	SELECT 
		c.id,
		c.user_id,
		c.theme_id,
		c.content,
		c.created_at,
		u.email
	FROM 
		comments c
	JOIN 
		users u ON u.id = c.user_id
	WHERE 
		c.theme_id = $1
	ORDER BY c.id ASC;
	`
	if user_role == "admin" {
		rows, err := s.db.Query(context.Background(), queryForAll, themeID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		var comments []models.Comment = []models.Comment{}
		for rows.Next() {
			var comment models.Comment
			err := rows.Scan(&comment.ID, &comment.UserID, &comment.ThemeID, &comment.Content, &comment.CreatedAt, &comment.Email)
			if err != nil {
				logrus.Printf("Error getting comments: %v", err)
				return nil, err
			}
			comments = append(comments, comment)
		}
		return comments, nil
	} else {
		rows, err := s.db.Query(context.Background(), query, themeID, user_id)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		var comments []models.Comment = []models.Comment{}
		for rows.Next() {
			var comment models.Comment
			err := rows.Scan(&comment.ID, &comment.UserID, &comment.ThemeID, &comment.Content, &comment.CreatedAt, &comment.Email)
			if err != nil {
				logrus.Printf("Error getting comments: %v", err)
				return nil, err
			}
			comments = append(comments, comment)
		}
		return comments, nil
	}

}

func (s *CommentStorage) DeleteComment(commentID int, user_id string, user_role string) error {
	query := `
		DELETE FROM comments
		WHERE id = $1 AND user_id = $2
	`

	queryDeleteReplyes := "DELETE FROM admin_replies WHERE comment_id = $1"

	if user_role == "admin" {
		queryForAdmin := `
		DELETE FROM comments
		WHERE id = $1
		`
		_, err := s.db.Exec(context.Background(), queryDeleteReplyes, commentID)
		if err != nil {
			logrus.Printf("Error deleting comment: %v", err)
			return err
		}

		_, err = s.db.Exec(context.Background(), queryForAdmin, commentID)
		if err != nil {
			logrus.Printf("Error deleting comment: %v", err)
			return err
		}
		return nil
	}

	tx, err := s.db.Begin(context.Background())
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback(context.Background())
			logrus.Printf("Transaction rolback with error: %s \n", err)
		} else {
			err := tx.Commit(context.Background())
			if err != nil {
				logrus.Fatalf("Unable to commit transaction: %v", err)
			}
		}
	}()

	_, err = tx.Exec(context.Background(),queryDeleteReplyes, commentID)
	if err != nil {
		return err
	}
	_, err = tx.Exec(context.Background(),query,commentID,user_id)
	if err != nil {
		return err
	}
	return nil
}

func (s *CommentStorage) CreateReply(reply models.AdminReply) error {
	query := `
		INSERT INTO admin_replies (comment_id, content, user_id)
		VALUES ($1, $2, $3)
		RETURNING id	
	`
	var id int
	row := s.db.QueryRow(context.Background(), query, reply.CommentID, reply.Content, reply.UserID)
	err := row.Scan(&id)
	if err != nil {
		logrus.Printf("Error creating reply: %v", err)
		return err
	}
	return nil
}

func (s *CommentStorage) GetRepliesByCommentID(commentID int) ([]models.CommentAdminReply, error) {
	query := `
		SELECT 
			ar.id,
			ar.comment_id,
			ar.content,
			ar.user_id,
			ar.created_at,
			u.email
		FROM 
    		admin_replies ar
		JOIN 
			users u ON u.id = ar.user_id
		WHERE ar.comment_id = $1
	`
	rows, err := s.db.Query(context.Background(), query, commentID)
	if err != nil {
		logrus.Printf("Error getting replies: %v", err)
		return nil, err
	}
	defer rows.Close()

	var replies []models.CommentAdminReply
	for rows.Next() {
		var reply models.CommentAdminReply
		err := rows.Scan(&reply.ID, &reply.CommentID, &reply.Content, &reply.UserID, &reply.CreatedAt, &reply.Email)
		if err != nil {
			logrus.Printf("Error getting replies: %v", err)
			return nil, err
		}
		replies = append(replies, reply)
	}
	return replies, nil
}
