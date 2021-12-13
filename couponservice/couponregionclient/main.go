package main

import (
	"fmt"
	"os"

	"github.com/go-redis/redis"
)

// A very simple consumer application
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

	// lets create a consumer group
	// status, err := client.XGroupCreate(area, "client-EU", "0").Result()
	// fmt.Println(status)

	// if err != nil {
	// 	fmt.Println("Could not create group", err)
	// 	return
	// }

	// XReadGroup is used by groups to read from stream.
	// Here we are using a special consumer id ">". This ensures that the message we are getting has
	// never been sent to another consumer.
	// However, if you dont specify ">", you will get a pending messages which have not been acknowledged.
	// XAck command removes the message from the history.

	// if we set NoAck true, this means our message is added to the message history of pending messages.
	// We can call XAck when the coupon code has been used by the client.
	streamData, err := client.XReadGroup(&redis.XReadGroupArgs{
		Group:    "client-EU",
		Consumer: "consumer-1",
		Streams:  []string{area, ">"},
		Count:    1,
		NoAck:    true,
	}).Result()

	// Here XReadGroup will wait indefinitely for messages if the stream is empty. As soon as
	// there is a message, it will poll the COUNT number of messages and then exit.

	if err != nil {
		fmt.Println("Could not Read from stream", err)
		return
	}

	fmt.Println(streamData)
}
