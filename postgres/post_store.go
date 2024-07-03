package postgres

import (
	"fmt"
	"github.com/alekseibykov/goreddit"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PostStore struct {
	*sqlx.DB
}

func (s *PostStore) Post(id uuid.UUID) (goreddit.Post, error) {
	var p goreddit.Post
	if err := s.Get(&p, `SELECT * FROM posts WHERE thread_id = $1`, id); err != nil {
		return goreddit.Post{}, fmt.Errorf("cannot get thread: %w", err)
	}
	return p, nil
}

func (s *PostStore) PostsByThreadID(id uuid.UUID) ([]goreddit.Post, error) {
	var pp []goreddit.Post
	if err := s.Select(&pp, `SELECT * FROM posts WHERE thread_id = $1`, id); err != nil {
		return []goreddit.Post{}, fmt.Errorf("cannot get posts: %w", err)
	}
	return pp, nil
}

func (s *PostStore) CreatePost(p *goreddit.Post) error {
	if err := s.Get(p, `INSERT INTO posts VALUES ($1, $2, $3, $4, $5) RETURNING *`,
		p.ID,
		p.ThreadID,
		p.Title,
		p.Content,
		p.Votes); err != nil {
		return fmt.Errorf("cannot insert post: %w", err)
	}
	return nil
}

func (s *PostStore) UpdatePost(p *goreddit.Post) error {
	if err := s.Get(p, `UPDATE posts SET thread_id = $1, title = $2, content = $3, votes = $4 WHERE id = $5) RETURNING *`,
		p.ThreadID,
		p.Title,
		p.Content,
		p.Votes,
		p.ID); err != nil {
		return fmt.Errorf("cannot update post: %w", err)
	}
	return nil
}

func (s *PostStore) DeletePost(id uuid.UUID) error {
	if _, err := s.Exec(`DELETE FROM posts WHERE id = $1`, id); err != nil {
		return fmt.Errorf("cannot delete post: %w", err)
	}
	return nil
}
