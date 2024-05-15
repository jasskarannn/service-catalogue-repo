# Service Catalogue Service

## Tech Stack - 
Go, PostreSQL
## Dependencies - 
Goose

## Run Migration - 
### goose -dir database/migrations postgres "postgres://postgres:postgres@db:5432/catalogue_service?sslmode=disable" up
### postgres://postgres:postgres@db:5432/catalogue_service?sslmode=disable
^ this a connection string and need to be updated according to configuration of postgres connection. 

## Build the service - 
go build -o service-catalogue main.go  

## Run the service - 
./service-catalogue

## API Endpoints of the service serve the following purposes -
1. Check DB Health
2. Add a new 'service' card
3. Add a new 'version'
4. Search for a service based on input text with pagination.
5. Retrieve all services with pagination.
6. Retrieve a particular service.

## Video Walkthrough of the service API Endpoints


https://github.com/jasskarannn/service-catalogue-repo/assets/59541154/d5189640-cad2-48b0-9853-efcd5d8ba0c6

## Docker 
docker-compose up --build --force-recreate
![image](https://github.com/jasskarannn/service-catalogue-repo/assets/59541154/cb6a3d64-e119-461d-9fe0-5850baaa578f)

