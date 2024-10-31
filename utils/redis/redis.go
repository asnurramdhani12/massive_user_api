package redis_utils

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

func DeleteWithPattern(ctx context.Context, rdb *redis.Client, pattern string) error {
	var cursor uint64

	for {
		// Scan for keys matching the pattern
		keys, newCursor, err := rdb.Scan(ctx, cursor, pattern, 0).Result()
		if err != nil {
			return fmt.Errorf("failed to scan keys: %v", err)
		}

		// Delete keys if found
		if len(keys) > 0 {
			if _, err := rdb.Del(ctx, keys...).Result(); err != nil {
				return fmt.Errorf("failed to delete keys: %v", err)
			}
		}

		// Update the cursor for the next iteration
		cursor = newCursor
		if cursor == 0 {
			break // Exit the loop if no more keys to scan
		}
	}

	return nil
}
