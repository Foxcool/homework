package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	dbMongo "github.com/foxcool/homework/5-k8s-crud/db/mongo"
	"github.com/foxcool/homework/5-k8s-crud/storage"
)

var (
	background = context.Background()
)

func New(client *dbMongo.Client) (*storage.Storage, error) {
	if err := ensure(client); err != nil {
		return nil, err
	}

	return &storage.Storage{
		UserStorage: newUserStorage(client),
	}, nil
}

func ensure(client *dbMongo.Client) error {
	_, err := client.Collection(UserCollectionName).Indexes().CreateMany(background, []mongo.IndexModel{
		{
			Keys: bson.M{
				"mobile": 1,
			},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.M{
				"email": 1,
			},
			Options: options.Index().SetUnique(true),
		},
	})
	if err != nil {
		return err
	}

	return nil
}
