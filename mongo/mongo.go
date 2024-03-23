package mongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConnectionInfo struct {
	Host     string
	Port     int
	Username string
	Password string
	DBName   string
	Timeout  time.Duration
}

func NewMongoConnection(info MongoConnectionInfo) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), info.Timeout)
	defer cancel()
	opts := options.Client()
	opts.SetAuth(options.Credential{
		Username: info.Username,
		Password: info.Password,
	})
	opts.ApplyURI(fmt.Sprintf("mongodb://%s:%d", info.Host, info.Port))

	dbClient, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	if err := dbClient.Ping(context.Background(), nil); err != nil {
		return nil, err
	}

	db := dbClient.Database(info.DBName)

	return db, nil
}
