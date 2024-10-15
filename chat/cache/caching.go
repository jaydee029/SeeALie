package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jaydee029/SeeALie/chat/internal/database"
	"github.com/redis/go-redis/v9"
)

func SetCacheFriends(ctx context.Context, cache *redis.Client, key string, value []database.Find_friendsRow) error {

	expirytime := 24 * time.Hour
	json, err := json.Marshal(value)
	if err != nil {
		return err
	}

	_, err = cache.Set(ctx, key, json, expirytime).Result()
	if err != nil {
		return err
	}

	return nil
}

func GetCacheFriends(ctx context.Context, cache *redis.Client, key string) ([]database.Find_friendsRow, error) {

	data, err := cache.Get(ctx, key).Result()
	if err != nil {
		return []database.Find_friendsRow{}, err
	}

	var friendsTable []database.Find_friendsRow

	err = json.Unmarshal([]byte(data), &friendsTable)
	if err != nil {
		return []database.Find_friendsRow{}, err
	}

	return friendsTable, nil
}
