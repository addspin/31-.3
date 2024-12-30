package mongo

import (
	"GoNews/pkg/storage"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Store represents the MongoDB storage.
type Store struct {
	db *mongo.Database
}

// New creates a new MongoDB store.
func New(uri string) (*Store, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	// Ping the database to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	db := client.Database("gonews")
	return &Store{db: db}, nil
}

// Posts returns all posts from the database.
func (s *Store) Posts() ([]storage.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := s.db.Collection("posts").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var posts []storage.Post
	if err = cursor.All(ctx, &posts); err != nil {
		return nil, err
	}
	return posts, nil
}

// AddPost adds a new post to the database.
func (s *Store) AddPost(p storage.Post) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := s.db.Collection("posts").InsertOne(ctx, p)
	return err
}

// UpdatePost updates an existing post in the database.
func (s *Store) UpdatePost(p storage.Post) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := s.db.Collection("posts").UpdateOne(
		ctx,
		bson.M{"id": p.ID},
		bson.M{"$set": p},
	)
	return err
}

// DeletePost removes a post from the database.
func (s *Store) DeletePost(p storage.Post) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := s.db.Collection("posts").DeleteOne(ctx, bson.M{"id": p.ID})
	return err
}
