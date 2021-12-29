## Technologies to use here

1. Prometheus
2. Grafana
3. Melody

## Idea
Lets create a simple Mutex based application. Anyone from Ops team can access a variable string. This can be obtained and freed using an http request. Anyone trying to access the resource should be given a Wait message while someone else is writing to it. We should be able to trace this and count metrics of different types. Also, lets create a simultaneous text log file to record the user activity who edit this and mutex activity. This string will be broadcasted via Websocket to anyone connected after each mutex release with the username who edited last. This can be for the Senior Management team to look at progress for a single project. 

## Workflow
1. Developer checks the status of the request using API route `/check`. If not available, Developer should not make any other requests.
2. If available, Developer should make request to `/avail`. This might assign the developer the mutex. However, if there is a race condition, We will monitor this.
3. If a Developer makes a `/avail` request, we will monitor that too.
4. To release the mutex, Developer makes a `/release` request. 
5. Each time mutex is released, we broadcast the string status.
6. A websocket can be requested at `/ws` route.

We will try to monitor this using Grafana and Prometheus. This is the main aim of this application.

We will try to set up alerts if the string is occupied for a large amount of time and for every 2 such instances.
