package main
import (
	"fmt"
	"os"
	"github.com/adi-ch1234/campusslots/internal/config"
	"github.com/adi-ch1234/campusslots/internal/telemetry"
)

func main(){
	cfg,err:= config.Load()
	if err!=nil{
		fmt.Fprintf(os.Stderr,"Failed to load configuration %v\n",err)
		os.Exit(1)
	}
	logger:=telemetry.NewLogger(cfg.AppEnv, cfg.LogLevel)
	logger.Info("configuration loaded",
		"env", cfg.AppEnv,
		"http_addr", cfg.HTTPAddr,
		"redis_addr", cfg.RedisAddr,
		"request_timeout", cfg.RequestTimeout.String(),
	)
}