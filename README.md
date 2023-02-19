# TikTok
a simple demo for TikTok backend

# Quick Start
## Setup Basic Dependence
docker-compose up -d

## Run Base RPC Server
cd script/base\
sh bootstrap.sh

## Run Social RPC Server
cd script/social\
sh bootstrap.sh

## Run api RPC Server
cd script/api\
sh bootstrap.sh

## API Requests
### Register
curl --location --request POST 'localhost:8888/douyin/user/register/' \
--header 'Content-Type: application/json' \
--data-raw '{
    "username":"lorain",
    "password":"123456"
}''
### response
// successful
{
    "status_code": 0,
    "status_msg": "操作成功",
    "user_id": 3,
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjozLCJ1c2VybmFtZSI6ImNjIiwiaXNzIjoiZG91eWluLXNlcnZpY2UiLCJleHAiOjE2NzY4MTgwNDJ9.Y2BEW4Pcdxz6AjXE63RzWjX438vreVygfpd0ku6hu2I"

}

### Login
#### will return jwt token
curl --location --request POST 'localhost:8080/douyin/douyin/user/login/' \
--header 'Content-Type: application/json' \
--data-raw '{
    "username":"lorain",
    "password":"123456"
}'

#### response
// successful
{
    "status_code": 0,
    "status_msg": "操作成功",
    "user_id": 1,
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6ImFhIiwiaXNzIjoiZG91eWluLXNlcnZpY2UiLCJleHAiOjE2NzY4MTgwNTZ9.aSoyiZ71IOmF8TnOi8jgXo2IoEh4vKM1irDMhG1pG4k"
}
