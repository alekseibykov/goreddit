package postgres

import (
	"fmt"
	"github.com/alekseibykov/goreddit"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type CommentStore struct {
	*sqlx.DB
}

func (s *CommentStore) Comment(id uuid.UUID) (goreddit.Comment, error) {
	var c goreddit.Comment
	err := s.Get(&c, "SELECT * FROM comments WHERE post_is = $1", id)
	if err != nil {
		return goreddit.Comment{}, fmt.Errorf("cannot get comment: %w", err)
	}
	return c, nil
}

func (s *CommentStore) CommentsByPost(postID uuid.UUID) ([]goreddit.Comment, error) {
	var cc []goreddit.Comment
	err := s.Select(&cc, "SELECT * FROM comments WHERE post_id = $1", postID)
	if err != nil {
		return []goreddit.Comment{}, fmt.Errorf("cannot get comments: %w", err)
	}
	return cc, nil
}

func (s *CommentStore) CreateComment(c *goreddit.Comment) error {
	err := s.Get(c, `INSERT INTO comments VALUES ($1, $2, $3, $4) RETURNING *`,
		c.ID,
		c.PostID,
		c.Content,
		c.Votes)
	if err != nil {
		return fmt.Errorf("cannot insert comment: %w", err)
	}
	return nil
}

func (s *CommentStore) UpdateComment(c *goreddit.Comment) error {
	err := s.Get(c, `UPDATE comments SET post_id = $1, content = $2, votes = $3 WHERE id = $4 RETURNING *`,
		c.PostID,
		c.Content,
		c.Votes,
		c.ID)
	if err != nil {
		return fmt.Errorf("cannot update comment: %w", err)
	}
	return nil
}

func (s *CommentStore) DeleteComment(id uuid.UUID) error {
	_, err := s.Exec(`DELETE FROM comments WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("cannot delete comment: %w", err)
	}
	return nil
}
