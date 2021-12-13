package main

import (
	"fmt"
	"os"

	"github.com/go-redis/redis"
)

func main() {
	// instances that connect to the stream to read data are consumers.
	// we use consumer groups to get read scalability.

	// declare the essentials and make redis connections
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

	area := "coupon-EU"

	status, err := client.XGroupCreate(area, "client-EU", "0").Result()

	if err != nil {
		fmt.Println("Could not create group")
	}
	fmt.Println(status)

}
