package mongo

import (
	"context"
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	dbMongo "github.com/foxcool/homework/7-prometheus/app/db/mongo"
	"github.com/foxcool/homework/7-prometheus/app/storage"
)

const UserCollectionName = "users"

type UserStorage struct {
	client *dbMongo.Client

	storage.UserStorage
}

func newUserStorage(client *dbMongo.Client) *UserStorage {
	return &UserStorage{
		client: client,
	}
}

func (s *UserStorage) StoreUser(in storage.User) (storage.User, error) {
	collection := s.client.Collection(UserCollectionName)
	ctx, _ := context.WithTimeout(background, time.Second*1)
	var err error

	_, err = collection.InsertOne(ctx, in)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return storage.User{}, storage.ErrUserAlreadyExists
		}

		return storage.User{}, err
	}

	return storage.User{
		ID:         in.ID,
		FirstName:  in.FirstName,
		LastName:   in.LastName,
		MiddleName: in.MiddleName,
		Mobile:     in.Mobile,
		Email:      in.Email,
	}, nil
}

func (s *UserStorage) UpdateUser(id string, in storage.User) (storage.User, error) {
	collection := s.client.Collection(UserCollectionName)
	ctx, _ := context.WithTimeout(background, time.Second*1)

	var (
		arrayFilters []interface{}
	)

	filter := bson.M{
		"_id": id,
	}

	updateOptions := bson.M{}
	update := bson.M{}
	addToSet := bson.M{}

	if in.FirstName != nil {
		update["firstName"] = *in.FirstName
	}

	if in.LastName != nil {
		update["lastName"] = *in.LastName
	}

	if in.MiddleName != nil {
		update["middleName"] = *in.MiddleName
	}

	if in.Email != nil {
		update["email"] = *in.Email
	}

	if in.Mobile != nil {
		update["mobile"] = *in.Mobile
	}

	if len(update) != 0 {
		updateOptions["$set"] = update
	}
	if len(addToSet) != 0 {
		updateOptions["$addToSet"] = addToSet
	}
	if len(updateOptions) == 0 {
		return storage.User{}, nil
	}
	result := collection.FindOneAndUpdate(ctx, filter, updateOptions,
		options.FindOneAndUpdate().SetReturnDocument(options.After).SetArrayFilters(options.ArrayFilters{
			Filters: arrayFilters,
		}))
	if result.Err() != nil {
		switch {
		case errors.Is(result.Err(), mongo.ErrNoDocuments):
			return storage.User{}, storage.ErrUserNotFound
		case strings.Contains(result.Err().Error(), "duplicate"):
			return storage.User{}, storage.ErrUserAlreadyExists
		}

		return storage.User{}, result.Err()
	}

	var user storage.User
	if err := result.Decode(&user); err != nil {
		return storage.User{}, err
	}

	return user, nil
}

func (s *UserStorage) DeleteUser(id string) error {
	collection := s.client.Collection(UserCollectionName)
	ctx, _ := context.WithTimeout(background, time.Second*1)

	res := collection.FindOneAndDelete(ctx, bson.M{"_id": id})
	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return storage.ErrUserNotFound
		}

		return res.Err()
	}

	return nil
}

func (s *UserStorage) GetUser(in storage.User) (storage.User, error) {
	collection := s.client.Collection(UserCollectionName)
	ctx, _ := context.WithTimeout(background, time.Second*1)

	filter := bson.M{}

	if in.ID != nil {
		filter["_id"] = in.ID
	}

	if in.Email != nil {
		filter["email.value"] = in.Email
	}

	if in.Mobile != nil {
		filter["mobile.value"] = in.Mobile
	}

	if len(filter) == 0 {
		return storage.User{}, storage.ErrBadInput
	}

	result := collection.FindOne(ctx, filter)
	if result.Err() != nil {
		switch result.Err() {
		case mongo.ErrNoDocuments:
			return storage.User{}, storage.ErrUserNotFound
		}

		return storage.User{}, result.Err()
	}

	var user storage.User
	if err := result.Decode(&user); err != nil {
		return storage.User{}, err
	}

	return user, nil
}

func (s *UserStorage) GetUsers() ([]storage.User, error) {
	collection := s.client.Collection(UserCollectionName)
	ctx, _ := context.WithTimeout(background, time.Second*10)

	users := make([]storage.User, 0)
	if cur, err := collection.Find(ctx, bson.M{}, options.Find()); err != nil {
		return nil, err
	} else {
		if cur.Err() != nil {
			return nil, cur.Err()
		}
		defer cur.Close(ctx)

		for cur.Next(ctx) {
			var user storage.User
			if err := cur.Decode(&user); err != nil {
				return nil, err
			}

			users = append(users, user)
		}
	}

	return users, nil
}
