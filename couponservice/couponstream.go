package couponservice

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/jasonlvhit/gocron"
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

// function to flush to the database
func flushdb(rdbi *redis.Client) {
	rdbi.FlushDB(ctx)
}

// This function is called by the main server to get the redis instance.
// The instance is again returned to the controller for this package. Now all the functions can access this like zap logger.
func RedisInstanceGenerator(logger *zap.Logger) *redis.Client {

	// declare the essentials and make redis connections
	var host = "localhost"
	var port = "6379"

	// get credentials from the Environment Variables
	if os.Getenv("REDIS_HOST") != "" {
		host = os.Getenv("REDIS_HOST")
	}
	if string(os.Getenv("REDIS_PORT")) != "" {
		port = string(os.Getenv("REDIS_PORT"))
	}

	// Declare the redis client
	// We can keep a password as an environment variable. However, I wont use this for simplicity.
	client := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: "",
		DB:       0,
	})

	// check status of connection.
	// It returns "PONG"
	_, err := client.Ping(ctx).Result()

	if err != nil {
		logger.Error("Redis connection failed")
		return nil
	}

	// Lets call zap to get the instance.
	logger.Info("Redis Instance Started", zap.Any("server details", map[string]interface{}{
		"Host":    client.Options().Addr,
		"Network": client.Options().Network,
	}))

	// this will make sure that every day at 12:00 AM, the database is flushed.
	gocron.Every(1).Day().At("00:00").Do(flushdb, client)

	// return the client to be used to router controller.
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
	if _, ok := r.Header["Couponregion"]; !ok {
		ctrl.logger.Warn("Coupon Region was not found in the header")
		handleNotInHeader(rw, r, "Region")
		return
	}

	coupon := store.Coupon{
		Name:        r.Header["Couponname"][0],
		VendorName:  r.Header["Couponvendor"][0],
		Code:        r.Header["Couponcode"][0],
		Description: r.Header["Coupondescription"][0],
		Region:      r.Header["Couponregion"][0],
	}

	couponJson, err := json.Marshal(coupon)
	if err != nil {
		ctrl.logger.Error("Cannot Marshal coupon", zap.Error(err))
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("An Internal server error ocurred"))
		return
	}

	// this automatically handles cases for new users.
	// if we want an already existing list and not create a new list, we use lpushx/rpushx
	_, err = ctrl.rdbi.RPush(ctx, coupon.VendorName, []interface{}{couponJson}).Result()
	if err != nil {
		ctrl.logger.Error("Fatal Redis Error", zap.Error(err))
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("An Internal server error ocurred"))
		return
	}

	// We have 4 regions. I want 4 streams to be there for each region. This will be given by Couponregion.
	// Any number of consumers can poll from this stream coming from that region.

	// it will be wise to put a string check over here. We dont need to add new streams if region is corrupt in Http request
	// Also, we need pretty flexible region names. We need to minimize Hardcoding in our application.
	region := r.Header["Couponregion"][0]
	if region != "APAC" && region != "NA" && region != "EU" && region == "SA" {
		// region does not exist.
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Provided Region does not exist"))
		return
	}

	err = ctrl.rdbi.XAdd(ctx, &redis.XAddArgs{
		Stream: "coupon-" + r.Header["Couponregion"][0],
		Values: map[string]interface{}{
			"coupon": couponJson,
		},
	}).Err()

	if err != nil {
		ctrl.logger.Error("Fatal Redis Error", zap.Error(err))
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("An Internal server error ocurred"))
		return
	}

	ctrl.logger.Info("Coupon added to the stream")

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Coupon added to Region Stream"))
}

// GetCouponForInternalValidation returns all the coupons of the certain vendor.
// we might need this for internal validation if the need arises. It must be remembered
// that the Redis database will be flushed every 24 hours. We dont want to fill the Database with excess data.
func (ctrl *StreamController) GetCouponForInternalValidation(rw http.ResponseWriter, r *http.Request) {
	// we ger vendor name in the header.
	if _, ok := r.Header["Vendorname"]; !ok {
		ctrl.logger.Warn("Vendor Name was not found in the header")
		handleNotInHeader(rw, r, "vendor")
		return
	}

	res, err := ctrl.rdbi.LRange(ctx, r.Header["Vendorname"][0], 0, 0).Result()

	if err != nil {
		ctrl.logger.Error("Fatal Redis Error", zap.Error(err))
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("An Internal server error ocurred"))
		return
	}

	ctrl.logger.Info("Redis Info", zap.Any("listdata", res))

	userResponse, err := json.Marshal(res)

	if err != nil {
		ctrl.logger.Error("Cannot Marshal []string", zap.Error(err))
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("An Internal server error ocurred"))
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(userResponse))
}

// Lets assume that at any moment, You want to purge the stream. We need a method for that too.
// We should use role authentication here.
func (ctrl *StreamController) PurgeStream(rw http.ResponseWriter, r *http.Request) {
	//Lets take in the region
	if _, ok := r.Header["Region"]; !ok {
		ctrl.logger.Warn("Region was not found in the header")
		handleNotInHeader(rw, r, "region")
		return
	}

	streamName := "coupon-" + r.Header["Region"][0]

	// this deletes the redis stream for the given region
	status, err := ctrl.rdbi.Del(ctx, streamName).Result()

	if err != nil {
		ctrl.logger.Error("Fatal Redis Error", zap.Error(err))
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("An Internal server error ocurred"))
		return
	}

	if status >= 1 {
		ctrl.logger.Info("Stream Deleted", zap.Any("streamName", streamName))
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte("Successfully deleted the stream"))
		return
	}

	ctrl.logger.Info("Stream Does not Exist", zap.Any("streamName", streamName))
	rw.WriteHeader(http.StatusBadRequest)
	rw.Write([]byte("Stream does not exist."))
}
