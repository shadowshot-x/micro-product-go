## Technologies to use here

1. Prometheus
2. Grafana
3. Mutex
4. Melody - A very brief implementation

## Idea
Lets create a simple Mutex based application. Anyone from Ops team can access a variable string. This can be obtained and freed using an http request. Anyone trying to access the resource should be given a Wait message while someone else is writing to it. We should be able to trace this and count metrics of different types. This string will be broadcasted via Websocket to anyone connected after each mutex release with the username who edited last. This can be for the Senior Management team to look at progress for a single project. 

## Workflow
1. Developer checks the status of the request using API route `/check`. If not available, Developer should not make any other requests.
2. If available, Developer should make request to `/avail`. This might assign the developer the mutex. However, if there is a race condition, We will monitor this.
3. If a Developer makes a `/avail` request, we will monitor that too.
4. To release the mutex, Developer makes a `/release` request. 
5. Each time mutex is released, we broadcast the string status.
6. A websocket can be requested at `/ws` route.

We will try to monitor this using Grafana and Prometheus. This is the main aim of this application.

We will try to set up alerts if the string is occupied for a large amount of time and for every 2 such instances.

Check the Status of the Retrospective Resource

`curl localhost:9090/retrospective/check`

Avail the Access to Retrospective String

`curl localhost:9090/retrospective/avail --request GET --header 'Username:ujjwal'`

Change the String

`curl localhost:9090/retrospective/change --request POST --header 'Username:ujjwal' --header 'Retrospective:this is changed retrospective'`

Release the Retrospective

`curl localhost:9090/retrospective/release --request POST --header 'Username:ujjwal'`

## Setting up Prometheus
Prometheus actually looks for metrics. We can expose Prometheus metrics in our project via HTTP. Metrics can be seen just as a count of an event you have sent. This can be a successful or a failed event, you can count these occurances for your application. Prometheus can scale the metric monitoring for our application. We can set alerts based on ratio of successful and error events.

`docker run -p 9091:9090 -v /home/ujjwal/Desktop/opensource/personal/product-go-micro/prometheus/prometheus.yml prom/prometheus`

## Setting up Grafana

We will use the OSS Release using a Docker Image of Grafana.

`docker run -p 3000:3000 grafana/grafana-oss`

Wait for the Migrations to complete and then open localhost:3000. The username and password default is `admin`

Then, Create a Dashboard on this Docker Image. 