package main

import (
	"context"
	"fmt"
	"task/config"
	"task/flood"
	"task/mongo"
	"time"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg)

	mongodb, err := mongo.NewMongoConnection(mongo.MongoConnectionInfo{
		Host:     cfg.Mongo.Host,
		Port:     cfg.Mongo.Port,
		Username: cfg.Mongo.Username,
		Password: cfg.Mongo.Password,
		DBName:   cfg.Mongo.DBName,
		Timeout:  cfg.Mongo.Timeout,
	})
	if err != nil {
		panic(err)
	}

	floodControl := flood.NewMongoFloodControl(mongodb, cfg.Window, cfg.MaxCalls)

	// Моделирование вызовов
	сases := []struct {
		userID int64
		wait   time.Duration
		cnt    int
	}{
		{userID: 1, wait: 1, cnt: 7},
		{userID: 1, wait: 6, cnt: 2},
	}

	// Проходим по тестовым случаям
	for _, c := range сases {
		fmt.Printf("Testing UserID: %d\n", c.userID)
		for i := 0; i < c.cnt; i++ {
			time.Sleep(c.wait * time.Second)
			ok, err := floodControl.Check(context.Background(), c.userID)
			if err != nil {
				fmt.Printf("Error checking flood control: %v\n", err)
				continue
			}
			fmt.Println(ok)
		}
		fmt.Println()
	}
}

// FloodControl интерфейс, который нужно реализовать.
// Рекомендуем создать директорию-пакет, в которой будет находиться реализация.
type FloodControl interface {
	// Check возвращает false если достигнут лимит максимально разрешенного
	// кол-ва запросов согласно заданным правилам флуд контроля.
	Check(ctx context.Context, userID int64) (bool, error)
}
