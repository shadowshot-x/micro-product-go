# Authentication Microservice. \[REST\]
## JWT
I will creating my own JWT for authentication. This requires a key. H256 algo will be used to generate the signature.

JWT has the following components
1. Header \[base64 encoded\]
2. Payload \[base64 encoded\]
3. Signature \[combination of header and payload hashed using a private key\]

## Swagger
This will be used for Documenting the API

## Schema for the User
1. Email
2. Full Name
3. Password hash
4. Username
5. CreateDate

## Things Learnt\[Covered in Article\] :-
1. Using Gorilla MUX for Routing and Subroutes
2. Implementing my own JWT Logic
3. Creation of Modules and handling module specific data
4. Writing Handlers for Sign In and Sign Up
5. Creating Middleware
6. Containerize the Application using Docker

## Running the Hashed Command

`curl http://localhost:9090/auth/signin --header 'Email:abc@gmail.com' --header 'Passwordhash:hashedme1'`

`curl http://localhost:9090/auth/signup --header 'Email:newuser@example.com' --header 'Passwordhash:hashedme1' --header 'Username:user77' --header 'Fullname:test user'`