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
