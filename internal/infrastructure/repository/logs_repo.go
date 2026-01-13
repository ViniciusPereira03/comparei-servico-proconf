package repository

import (
	"comparei-servico-promer/internal/domain/logs"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoLogsRepository struct {
	collection *mongo.Collection
}

func NewLogsRepository(client *mongo.Client, dbName, collectionName string) *MongoLogsRepository {
	coll := client.Database(dbName).Collection(collectionName)
	return &MongoLogsRepository{collection: coll}
}

// CreateLogsConfirmacao insere um novo log de confirmação no MongoDB
func (r *MongoLogsRepository) CreateLogsConfirmacao(log *logs.LogsConfirmacao) (*logs.LogsConfirmacao, error) {
	log.CreatedAt = time.Now()
	res, err := r.collection.InsertOne(context.Background(), log)
	if err != nil {
		return nil, err
	}
	oid := res.InsertedID.(primitive.ObjectID)
	log.ID = oid.Hex()
	return log, nil
}
