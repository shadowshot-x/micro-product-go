# Authentication Microservice. \[REST\]
## JWT
I will be using JWT for authentication. This requires a key. H256 algo will be used to generate the signature.

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

## Running the Hashed Command

`curl http://localhost:9090/signin --header 'Email:abc@gmail.com' --header 'Passwordhash:hashedme1'`

`curl http://localhost:9090/signup --header 'Email:newuser@example.com' --header 'Passwordhash:hashedme1' --header 'Username:user77' --header 'Fullname:test user'`