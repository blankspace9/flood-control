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
	session, err := rfc.mongo.Client().StartSession()
	if err != nil {
		return false, err
	}
	defer session.EndSession(ctx)

	now := time.Now().Unix()
	var count int64

	err = mongo.WithSession(ctx, session, func(sessionCtx mongo.SessionContext) error {
		filter := bson.M{
			"userID":    userID,
			"timestamp": bson.M{"$lt": now - int64(rfc.window.Seconds())},
		}

		// Удаляем записи о старых вызовах за пределами текущего интервала
		_, err := rfc.mongo.Collection("calls").DeleteMany(ctx, filter)
		if err != nil {
			return err
		}

		// Считаем количество оставшихся документов (количество вызовов за последние N секунд)
		count, err = rfc.mongo.Collection("calls").CountDocuments(ctx, bson.M{"userID": userID})
		if err != nil {
			return err
		}

		// Вставляем информацию о текущем вызове
		rfc.mongo.Collection("calls").InsertOne(ctx, map[string]interface{}{
			"userID":    userID,
			"timestamp": now,
		})

		return nil
	})
	if err != nil {
		return false, err
	}

	if count+1 > int64(rfc.maxCalls) {
		return false, nil
	}

	return true, nil
}
