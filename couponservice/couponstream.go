package couponservice

import (
	"context"
	"os"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

var ctx = context.Background()

// StreamController is the Upload route handler
type StreamController struct {
	logger *zap.Logger
}

// NewCouponStreamController returns a frsh Stream controller
func NewCouponStreamController(logger *zap.Logger) *StreamController {
	return &StreamController{
		logger: logger,
	}
}

func (ctrl *StreamController) CouponStreamGenerator() {

	var host = "localhost"
	var port = "6379"
	if os.Getenv("REDIS_HOST") != "" {
		host = os.Getenv("REDUS_HOST")
	}
	if string(os.Getenv("REDIS_PORT")) != "" {
		port = string(os.Getenv("REDIS_PORT"))
	}
	client := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: "",
		DB:       0,
	})

	// congratulations, we set up a connection and pushed a key-value pair to redis
	status := client.Set(ctx, "couponCodeSetCheck", "running", 0)
	ctrl.logger.Info("Check out the status", zap.Any("redis status", status))
}
