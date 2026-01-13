package repository

import (
	"context"
	"log"
	"time"

	"comparei-servico-promer/internal/domain/users"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(client *mongo.Client, dbName, collectionName string) *MongoRepository {
	coll := client.Database(dbName).Collection(collectionName)
	return &MongoRepository{collection: coll}
}

// CreateUser insere um novo usuário no MongoDB
func (r *MongoRepository) CreateUser(u *users.User) (*users.User, error) {
	log.Println("EXEC: repo.CreateUser")
	// Insere o documento de usuário
	u.CreatedAt = time.Now()
	u.ModifiedAt = time.Now()
	res, err := r.collection.InsertOne(context.Background(), u)
	if err != nil {
		return nil, err
	}
	// Converter ObjectID para string
	oid := res.InsertedID.(primitive.ObjectID)
	u.ID = oid.Hex()
	return u, nil
}

// GetUser busca o usuário pelo ID do serviço de usuários
func (r *MongoRepository) GetUser(id string) (*users.User, error) {
	var u users.User
	err := r.collection.FindOne(context.Background(), bson.M{"id_usuario": id}).Decode(&u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// UpdateLevelUser atualiza o nível do usuário
func (r *MongoRepository) UpdateLevelUser(u *users.User) error {
	oid, err := primitive.ObjectIDFromHex(u.ID)
	if err != nil {
		return err
	}
	update := bson.M{"$set": bson.M{
		"level":       u.Level,
		"modified_at": time.Now(),
	}}
	_, err = r.collection.UpdateByID(context.Background(), oid, update)
	return err
}

// DeleteUser marca o usuário como inativo
func (r *MongoRepository) DeleteUser(u *users.User) error {
	oid, err := primitive.ObjectIDFromHex(u.ID)
	if err != nil {
		return err
	}
	update := bson.M{"$set": bson.M{"deleted_at": time.Now()}}
	_, err = r.collection.UpdateByID(context.Background(), oid, update)
	return err
}
