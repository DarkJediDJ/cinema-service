# Cinema-service is a golang API
  This Web API designed to manage cinema activity.
  
  This API provides function to register and authorize user.
  
  Standart premissions are:
  * Buy ticket
  * Download ticket
  
  Superadmin can add different privileges to users to make them admins.
  
  This privileges provides CRUD operations on:
  * Halls
  * Movies
  * Sessions
  * Tickets
  
## Project Layout

The cinetickets uses the following project layout:
```
.
├── api                  
│   └── halls
│   └── movies 
│   └── sessions 
│   └── tickets
│   └── users
│   └── user_privileges
│   └── server.go 
├── api
│   └── cinetickets
│       └── main.go
├── db
│   └── migrations
├── internal
│   └── repository
│       └── halls
│       └── movies 
│       └── sessions 
│       └── tickets
│       └── users
│       └── user_privileges
│   └── service
│       └── halls
│       └── movies 
│       └── sessions 
│       └── tickets
│       └── users
│       └── user_privileges             
│   ├── errors.go           
│   ├── interfaces.go                      
├── package
│   └── aws
│   └── generator
│   └── grpc
│   └── jwt
│   └── encryption.go
├── test
│   └── mockService.go
└──
            
```

## Setup localy

* Clone this repository
* Run `go mod tidy`

### Configure `.env` file:

* `DB_HOST = host`
* `DB_NAME = name`
* `DB_PORT = port`
* `DB_USER = user`
* `DB_PASSWORD = password`
* `ACCESS_SECRET = key`

### Configure AWS
* https://aws.amazon.com/cli/?nc1=h_ls

### Setup ticketgenerator service
* clone https://github.com/DarkJediDJ/ticketgenerator
* set environment variables
* `go run cmd/ticketgenerator/main.go`

### Run 
* `go run cmd/cinetickets/main.go`

## Docker containers
* https://hub.docker.com/repository/docker/darkjedidj/ticketgenerator
* https://hub.docker.com/repository/docker/darkjedidj/cinetickets

## Hosting API
* http://cinema-alb-dev-o81jt53c-906642332.us-east-1.elb.amazonaws.com:8085/swagger/index.html#/

## Testing 
* `go run test ./...`
