package cors

import (
	"gin-center/configs/config"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORSMiddleware(cfg *config.GlobalConfig) gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     cfg.Server.CORS.AllowOrigins,
		AllowMethods:     cfg.Server.CORS.AllowMethods,
		AllowHeaders:     cfg.Server.CORS.AllowHeaders,
		ExposeHeaders:    cfg.Server.CORS.ExposeHeaders,
		AllowCredentials: cfg.Server.CORS.AllowCredentials,
		MaxAge:           time.Duration(cfg.Server.CORS.MaxAge) * time.Second,
	})
}
