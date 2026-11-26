package database
import (
	"fmt"
	"context"
	"time"
	"github.com/redis/go-redis/v9"
)
func NewRedisClient(ctx context.Context, addr string )(*redis.Client,error){
	client:=redis.NewClient(&redis.Options{
		Addr:addr,
		PoolSize:20,
		DialTimeout:2*time.Second,
		ReadTimeout:1*time.Second,
		WriteTimeout:1*time.Second,
	})
	pingCtx,cancel:= context.WithTimeout(ctx,2*time.Second)
	defer cancel();
	if err:= client.Ping(pingCtx).Err();err!=nil{
		// Still return the client along with the error so the caller can decide
		// whether to fail fast or degrade gracefully (go-redis auto-reconnects later).
		return client, fmt.Errorf("ping redis: %w", err)
	}
	return client,nil
}
