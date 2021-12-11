package couponservice

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/shadowshot-x/micro-product-go/couponservice/store"
	"go.uber.org/zap"
)

var ctx = context.Background()

// StreamController is the Upload route handler
type StreamController struct {
	logger *zap.Logger
	rdbi   *redis.Client
}

// NewCouponStreamController returns a frsh Stream controller
func NewCouponStreamController(logger *zap.Logger, instance *redis.Client) *StreamController {
	return &StreamController{
		logger: logger,
		rdbi:   instance,
	}
}

func handleNotInHeader(rw http.ResponseWriter, r *http.Request, param string) {
	rw.WriteHeader(http.StatusBadRequest)
	rw.Write([]byte(fmt.Sprintf("%s Missing", param)))
}

// This function is called by the main server to get the redis instance.
// The instance is again returned to the controller for this package. Now all the functions can access this like zap logger.
func RedisInstanceGenerator(logger *zap.Logger) *redis.Client {
	var host = "localhost"
	var port = "6379"
	if os.Getenv("REDIS_HOST") != "" {
		host = os.Getenv("REDIS_HOST")
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
	logger.Info("Check out the status", zap.Any("redis status", status))

	return client
}

func (ctrl *StreamController) AddCouponList(rw http.ResponseWriter, r *http.Request) {
	if _, ok := r.Header["Couponname"]; !ok {
		ctrl.logger.Warn("Coupon Name was not found in the header")
		handleNotInHeader(rw, r, "name")
		return
	}
	if _, ok := r.Header["Couponvendor"]; !ok {
		ctrl.logger.Warn("Coupon Vendor was not found in the header")
		handleNotInHeader(rw, r, "Vendor")
		return
	}
	if _, ok := r.Header["Couponcode"]; !ok {
		ctrl.logger.Warn("Coupon Inventory was not found in the header")
		handleNotInHeader(rw, r, "Inventory")
		return
	}
	if _, ok := r.Header["Coupondescription"]; !ok {
		ctrl.logger.Warn("Coupon Description was not found in the header")
		handleNotInHeader(rw, r, "Description")
		return
	}

	coupon := store.Coupon{
		Name:        r.Header["Couponname"][0],
		VendorName:  r.Header["Couponvendor"][0],
		Code:        r.Header["Couponcode"][0],
		Description: r.Header["Coupondescription"][0],
	}

	// this automatically handles cases for new users.
	// if we want an already existing list and not create a new list, we use lpushx/rpushx
	status, err := ctrl.rdbi.RPush(ctx, r.Header["Couponvendor"][0], coupon).Result()
	if err != nil {
		ctrl.logger.Error("Fatal Redis Error", zap.Error(err))
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("An Internal server error ocurred"))
		return
	}

	if status == 0 {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Could not append your Request to Redis"))
		return
	}

	// add logic to add to a stream
}
