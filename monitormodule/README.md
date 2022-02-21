# Monitor Module
Add a VERY SIMPLE module to monitor the existing microservices. Some of the main tasks for this issue :-
1. Modify each microservice to have a function to tell the health of the application. This we can see from the Prometheus points. We can get number of successful metrics and errors.
2. Implement a new module that will run in parallel inside your main application. It will monitor each microservice application based on a JSON file provided in the config by sending a ping request to the health function after a time interval of some seconds configured in the same JSON.
3. Module will run purely with GoRoutines, Channels and Sync package. The Module will actually be a data pipeline with a defined architecture taking care of fan-in and fan-out. 
4. The data should be written to a client every 10 seconds. The client will be coded and can have very simple stdout listener using REST API.

Data Pipeline Steps :- 

- Init function which initiates the GoRoutines based on a JSON file. It syncs the below 3 steps for Data Pipeline.
- Data Aggregation. Have dedicated GoRoutines to monitor application health and send the results of application in an unbuffered channel.
- Data Transformation.  Output returned by Health functions is in string format. Have multiple GoRoutines read from this 1 unbuffered channel and transform the data to a specific Output struct and Stringify this into a JSON string. Each GoRoutine posts to its own channel.
- Data Transportation. Have another GoRoutine to aggregate the data using Fan-in mechanism copying data of multiple channels into one channel and post a REST request to the Monitoring Data Ingester.