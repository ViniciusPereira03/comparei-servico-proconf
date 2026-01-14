package repository

import (
	"comparei-servico-proconf/internal/domain/proconf"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProconfRepository struct {
	collection *mongo.Collection
}

func NewProconfRepository(client *mongo.Client, dbName, collectionName string) *ProconfRepository {
	coll := client.Database(dbName).Collection(collectionName)
	return &ProconfRepository{collection: coll}
}

func (r *MongoRepository) Create(u *proconf.Proconf) (*proconf.Proconf, error) {
	return u, nil
}
