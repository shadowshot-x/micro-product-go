# Redis Database Usage with Minikube

We are building a coupon service which is essentially like a streaming service for some requests. We will use Redis Streams for Queuing. We will then make the image of the service and deploy the redis instance and image to a minikube kubernetes cluster. The essential commands will be in this Readme.

## Working of the Coupon Service

We have 4 major regions (APAC, NA, SA, EU) with running service who want coupons to distribute to their local users. They constantly demand coupon codes from our server. As merchants add the coupons, they become available for redemption.

## What we are Building?
We are making a microservice that has a coupon code distribution instance. This can be demanded from different consumers. The coupons are stored in our database categorized based on merchants. The merchants add their codes to our database using REST API. For each merchant we have a Redis list. After every coupon addition, messages are added to redis stream based on the coupons in the list. These streams are subscibed by consumers who then get code based on requests. We also purge all the older requests every 24 hours if they have not been fulfilled.

## Redis Database

Redis is an open source database that can be treated as a key-value store that we can use for Database purposes(NoSQL), Caching and for Message brokering.

We will run redis in a Docker container. This way we will have the image ready and we can deploy this to minikube.

`$ docker run --name coupon-redis-instance -p 6379:6379 -d redis`

## Redis Structure

For every new merchant, lets add a new list to the Redis database with vendorname as the key. Each time vendor adds a coupon list, we should add it to the redis db list corresponding to that vendor name.

`curl http://localhost:9090/coupon/addcoupon --request POST --header 'Couponname:off_50_flat' --header 'Couponvendor:vendor1' --header 'Coupondescription:Avail flat 50 off on all products' --header 'Couponcode:EU778' --header 'Couponregion:EU'`

`curl http://localhost:9090/coupon/getvendorcoupons --request GET --header 'Vendorname:vendor1'`

`curl http://localhost:9090/coupon/delregionstream --request DELETE --header 'Region:EU'`

## Running the Redis Consumer Code [for testing with Minikube]

We are using consumer groups of redis to read from a stream. It is a good practice as we want each message to used just once and in FIFO order. 

`go run couponservice/couponregionclient/main.go`

## Minikube Instructions

We will push our image to dockerhub so that this can be downloaded by minikube. Do remember to login first.

`docker push shadowshotx/product-go-micro`

Now, we will deploy the .yaml files to create deployment and services.

Lets start the minikube instance.
`minikube start`

Remember to have `minikube` and `kubectl` installed on your local. 

`kubectl apply -f couponservice/deployments/redis-deployment.yaml`

`kubectl apply -f couponservice/deployments/redis-service.yaml`

Now wait for the Redis Instances to get up and running. Then run the following commands.

`kubectl apply -f couponservice/deployments/go-micro-deployment.yaml`

`kubectl apply -f couponservice/deployments/go-micro-service.yaml`

Now, check the status of the services and pods by running

`kubectl get pods`

`kubectl get services`

Remember to stop your minikube instance.

`minikube stop`

## Workflow of the Architecture

![coupon-redis-architecture drawio](https://user-images.githubusercontent.com/43992469/146631429-cd2b8236-c710-41dc-b29b-70ac6b089f76.png)
