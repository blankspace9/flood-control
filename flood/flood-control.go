package flood

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoFloodControl struct {
	mongo    *mongo.Database
	window   time.Duration
	maxCalls int
}

func NewMongoFloodControl(mongo *mongo.Database, window time.Duration, maxCalls int) *MongoFloodControl {
	return &MongoFloodControl{
		mongo:    mongo,
		window:   window,
		maxCalls: maxCalls,
	}
}

func (rfc *MongoFloodControl) Check(ctx context.Context, userID int64) (bool, error) {
	now := time.Now().Unix()

	filter := bson.M{
		"userID":    userID,
		"timestamp": bson.M{"$lt": now - int64(rfc.window.Seconds())},
	}

	_, err := rfc.mongo.Collection("calls").DeleteMany(ctx, filter)
	if err != nil {
		return false, err
	}

	count, err := rfc.mongo.Collection("calls").CountDocuments(ctx, bson.M{})
	if err != nil {
		return false, err
	}

	rfc.mongo.Collection("calls").InsertOne(ctx, map[string]interface{}{
		"userID":    userID,
		"timestamp": now,
	})

	if count+1 > int64(rfc.maxCalls) {
		return false, nil
	}

	return true, nil
}
