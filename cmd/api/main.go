package main
import (
	"fmt"
	"os"
	"context"
	"github.com/adi-ch1234/campusslots/internal/config"
	"github.com/adi-ch1234/campusslots/internal/telemetry"
	"github.com/adi-ch1234/campusslots/internal/database"
)

func main(){
	cfg,err:= config.Load()
	if err!=nil{
		fmt.Fprintf(os.Stderr,"Failed to load configuration %v\n",err)
		os.Exit(1)
	}
	logger:=telemetry.NewLogger(cfg.AppEnv, cfg.LogLevel)
	ctx:=context.Background();	
	pgPool,err:=database.NewPostgresPool(ctx,cfg.DatabaseURL)
	if err!=nil{
		logger.Error("failed to connect to postgres", "error", err)
		os.Exit(1)
	}
	defer pgPool.Close()
	logger.Info("Postgres is connected")
	redisClient,err:=database.NewRedisClient(ctx,cfg.RedisAddr)
	if err!=nil{
		logger.Warn("redis unavailable at startup; continuing with degraded caching",
			"redis_addr", cfg.RedisAddr,
			"error", err,
		)
	}else{
		logger.Info("connected to redis", "redis_addr", cfg.RedisAddr)
	}
	defer redisClient.Close();
	logger.Info("startup checks passed",
		"env", cfg.AppEnv,
		"http_addr", cfg.HTTPAddr,
	)
}