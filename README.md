# TikTok
a simple demo for TikTok backend

# Quick Start
## Setup Basic Dependence
    docker-compose up
## Run User RPC Server
cd cmd/user
sh build.sh
sh output/bootstrap.sh

## Run API Server
cd cmd/api
go run .

## API Requests
### Register
curl --location --request POST 'localhost:8080/v1/user/register' \
--header 'Content-Type: application/json' \
--data-raw '{
    "username":"lorain",
    "password":"123456"
}''
### response
// successful
{
    "code": 0,
    "message": "Success",
    "data": null
}

// failed
{
    "code": 10003,
    "message": "User already exists",
    "data": null
}

### Login
#### will return jwt token
curl --location --request POST 'localhost:8080/v1/user/login' \
--header 'Content-Type: application/json' \
--data-raw '{
    "username":"lorain",
    "password":"123456"
}'

#### response
// successful
{
    "code": 0,
    "expire": "2022-12-3T01:56:46+08:00",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NDI1Mjg2MDYsImlkIjoxLCJvcmlnX2lhdCI6MTY0MjUyNTAwNn0.k7Ah9G4Enap9YiDP_rKr5HSzF-fc3cIxwMZAGeOySqU"
}

// failed
{
    "code": 10004,
    "message": "Authorization failed",
    "data": null
}
