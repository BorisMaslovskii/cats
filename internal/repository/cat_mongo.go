package repository

import (
	"context"
	"errors"

	"github.com/BorisMaslovskii/cats/internal/model"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	log "github.com/sirupsen/logrus"
)

// catRepositoryMongo struct
type catRepositoryMongo struct {
	conn *mongo.Collection
}

// NewRepoMongo func creates new catRepositoryMongo
func NewRepoMongo(conn *mongo.Collection) CatRepository {
	return &catRepositoryMongo{conn: conn}
}

// GetAll gets all cats from catRepositoryMongo
func (r *catRepositoryMongo) GetAll(ctx context.Context) ([]*model.Cat, error) {
	cats := make([]*model.Cat, 0)
	cursor, err := r.conn.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer func() {
		err := cursor.Close(ctx)
		if err != nil {
			log.Errorf("cursor Close error %v", err)
		}
	}()

	for cursor.Next(ctx) {
		cat := &model.Cat{}
		err := cursor.Decode(cat)
		if err != nil {
			return nil, err
		}
		cats = append(cats, cat)
	}
	return cats, nil
}

// GetByID func gets a cat by id from catRepositoryMongo
func (r *catRepositoryMongo) GetByID(ctx context.Context, id uuid.UUID) (*model.Cat, error) {
	cat := &model.Cat{}
	err := r.conn.FindOne(ctx, bson.M{"_id": id}).Decode(cat)
	if err != nil {
		return nil, err
	}
	return cat, nil
}

// Create func creates a new cat into catRepositoryMongo
func (r *catRepositoryMongo) Create(ctx context.Context, cat *model.Cat) (uuid.UUID, error) {
	uid := uuid.New()

	_, err := r.conn.InsertOne(ctx, bson.M{"_id": uid, "name": cat.Name, "color": cat.Color})
	if err != nil {
		return uuid.Nil, err
	}

	return uid, nil
}

// Delete func deletes a cat from catRepositoryMongo
func (r *catRepositoryMongo) Delete(ctx context.Context, id uuid.UUID) error {
	res, err := r.conn.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("cat was not found")
	}
	return nil
}

// Update func updates a cat in catRepositoryMongo
func (r *catRepositoryMongo) Update(ctx context.Context, id uuid.UUID, cat *model.Cat) error {
	update := bson.D{primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: "name", Value: cat.Name}, primitive.E{Key: "color", Value: cat.Color}}}}
	res, err := r.conn.UpdateByID(ctx, id, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("cat was not found")
	}
	return nil
}
