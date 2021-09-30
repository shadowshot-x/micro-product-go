# File Handling Microservice

## MultiPart File Uploading
We are using multi part HTTP Requests which means each request can contain multiple key value pairs for file uploading. Currently we are storing the files in our server itself but a good practice is to rent out a cloud storage and write the logic in handler to send files to this. There are many mechanisms for this including Data Pipelines with Firehose in AWS. I have worked with S3 and it works the best for file upload and is a cheap option.

## Things Learnt
1. Handling Files in Golang
2. MultiPart HTTP Form Requests

## Running the File Upload Service

`curl -v -F file=@/home/ujjwal/Downloads/download.jpeg --header 'Token:SFMyNTY=.eyJhdWQiOiJmcm9udGVuZC5rbm93c2VhcmNoLm1sIiwiZXhwIjoiMTYzMjk4NzQwNiIsImlzcyI6Imtub3dzZWFyY2gubWwifQ==./4hTLnjLW1tt5tdHAq6hph1R7IGm5uJWehheZrMu24M='  localhost:9090/claims/upload `