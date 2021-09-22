# Golang Server Project Best Practices

## Dependency Injection :-
In simple words, we want our functions and packages to receive the objects they depend on ie. We dont want to declare new instances inside packages to have control over them. For Eg :- Using Structs to declare the private methods and save logger as a variables. The methods can access the value of logger using `g.logger` in their domain.

## Handle Timeouts :-
This is to prevent DOS attacks. Dont make requests make infinitely if your server crashes.

## Graceful Shutdown :-
Wait until current requests are handled and then shutdown the server. We use signal interrupts for this with channels

## Using JSON Encoder :-
Sometimes it is better to use Encoder than json.Marshal as we dont have to use additional buffer. This will matter when there are a lot of concurrent go routines being processed. It is also a bit faster.

## Validation of Data is very Important :-
Middleware ensures Connectivity between 2 or more types applications or components. You can write your validation code in middleware to make sure data is validated and then only goes to your handlers

## Running using Docker Image
`docker build . -t product-go-micro`

`docker run --network host -d a3faa264fcc3`
