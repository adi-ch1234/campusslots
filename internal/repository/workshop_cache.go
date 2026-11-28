package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/adi-ch1234/campusslots/internal/domain"
)

const (
	workshopListCacheKey    = "campusslots:workshops:list"
	workshopDetailKeyFmt    = "campusslots:workshops:%d"
	defaultWorkshopCacheTTL = 30 * time.Second
)

type WorkshopCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewWorkshopCache(client *redis.Client) *WorkshopCache {
	return &WorkshopCache{
		client: client,
		ttl:    defaultWorkshopCacheTTL,
	}
}
func (c *WorkshopCache) GetList(ctx context.Context) ([]domain.Workshop, bool, error) {
	raw, err := c.client.Get(ctx, workshopListCacheKey).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("redis get workshop list: %w", err)
	}
	var workshops []domain.Workshop
	if err := json.Unmarshal(raw, &workshops); err != nil {
		return nil, false, fmt.Errorf("unmarshal cached workshop list: %w", err)
	}
	return workshops, true, nil
}
func (c *WorkshopCache) SetList(ctx context.Context, workshops []domain.Workshop) error {
	raw, err := json.Marshal(workshops)
	if err != nil {
		return fmt.Errorf("marshal workshop list: %w", err)
	}
	if err := c.client.Set(ctx, workshopListCacheKey, raw, c.ttl).Err(); err != nil {
		return fmt.Errorf("redis set workshop list: %w", err)
	}
	return nil
}
func (c *WorkshopCache) GetByID(ctx context.Context, id int64) (domain.Workshop, bool, error) {
	key := fmt.Sprintf(workshopDetailKeyFmt, id)
	raw, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return domain.Workshop{}, false, nil
		}
		return domain.Workshop{}, false, fmt.Errorf("redis get workshop %d: %w", id, err)
	}
	var w domain.Workshop
	if err := json.Unmarshal(raw, &w); err != nil {
		return domain.Workshop{}, false, fmt.Errorf("unmarshal cached workshop %d: %w", id, err)
	}

	return w, true, nil
}
func (c *WorkshopCache) SetByID(ctx context.Context, w domain.Workshop) error {
	raw, err := json.Marshal(w)
	if err != nil {
		return fmt.Errorf("marshal workshop %d: %w", w.ID, err)
	}

	key := fmt.Sprintf(workshopDetailKeyFmt, w.ID)
	if err := c.client.Set(ctx, key, raw, c.ttl).Err(); err != nil {
		return fmt.Errorf("redis set workshop %d: %w", w.ID, err)
	}

	return nil
}
func (c *WorkshopCache) Invalidate(ctx context.Context, workshopID int64) error {
	keys := []string{workshopListCacheKey}
	if workshopID > 0 {
		keys = append(keys, fmt.Sprintf(workshopDetailKeyFmt, workshopID))
	}
	if err := c.client.Del(ctx, keys...).Err(); err != nil {
		return fmt.Errorf("redis delete workshop keys: %w", err)
	}
	return nil
}
