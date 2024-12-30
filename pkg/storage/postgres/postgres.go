package postgres

import (
	"GoNews/pkg/storage"
	"database/sql"

	_ "github.com/lib/pq"
)

// Store represents the PostgreSQL storage.
type Store struct {
	db *sql.DB
}

// New creates a new PostgreSQL store.
func New(connStr string) (*Store, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &Store{db: db}, nil
}

// Posts returns all posts from the database.
func (s *Store) Posts() ([]storage.Post, error) {
	rows, err := s.db.Query(`
		SELECT id, title, content, author_id, author_name, created_at, published_at 
		FROM posts
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []storage.Post
	for rows.Next() {
		var p storage.Post
		err := rows.Scan(
			&p.ID,
			&p.Title,
			&p.Content,
			&p.AuthorID,
			&p.AuthorName,
			&p.CreatedAt,
			&p.PublishedAt,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}

// AddPost adds a new post to the database.
func (s *Store) AddPost(p storage.Post) error {
	_, err := s.db.Exec(`
		INSERT INTO posts (title, content, author_id, author_name, created_at, published_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`,
		p.Title,
		p.Content,
		p.AuthorID,
		p.AuthorName,
		p.CreatedAt,
		p.PublishedAt,
	)
	return err
}

// UpdatePost updates an existing post in the database.
func (s *Store) UpdatePost(p storage.Post) error {
	_, err := s.db.Exec(`
		UPDATE posts 
		SET title=$2, content=$3, author_id=$4, author_name=$5, created_at=$6, published_at=$7
		WHERE id=$1
	`,
		p.ID,
		p.Title,
		p.Content,
		p.AuthorID,
		p.AuthorName,
		p.CreatedAt,
		p.PublishedAt,
	)
	return err
}

// DeletePost removes a post from the database.
func (s *Store) DeletePost(p storage.Post) error {
	_, err := s.db.Exec("DELETE FROM posts WHERE id=$1", p.ID)
	return err
}
