package mongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client struct {
	Config *Config

	*mongo.Database
	Ctx context.Context
}

type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

func Connect(config *Config) (*Client, error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)

	var client *mongo.Client
	var err error
	if len(config.Username) > 0 && len(config.Password) > 0 {
		client, err = mongo.NewClient(
			options.Client().ApplyURI(
				fmt.Sprintf("mongodb://%s:%s@%s:%d/%s", config.Username, config.Password, config.Host, config.Port, config.Database),
			),
		)
	} else {
		client, err = mongo.NewClient(
			options.Client().ApplyURI(
				fmt.Sprintf("mongodb://%s:%d/%s", config.Host, config.Port, config.Database),
			),
		)
	}

	if err != nil {
		return nil, err
	}

	if err := client.Connect(ctx); err != nil {
		return nil, err
	}

	return &Client{
		Config: config,

		Database: client.Database(config.Database),

		Ctx: ctx,
	}, nil
}

func (c *Client) Session(fn func() error) (err error) {
	client := c.Client()

	if err = client.UseSession(c.Ctx, func(sc mongo.SessionContext) (err error) {
		if err = sc.StartTransaction(); err != nil {
			return err
		}
		if err = fn(); err != nil {
			sc.AbortTransaction(sc)
			return err
		}

		if err = sc.CommitTransaction(sc); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}
	return nil
}
