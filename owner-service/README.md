 ## Owner Service Documentation
 
 
 ###  Prerequisites
 - [Go 1.20 or higher](https://go.dev/)
 - [Docker](https://www.docker.com/)
 - [Docker Compose](https://docs.docker.com/compose/)
 - [YQ - CLI-based yaml processor](https://github.com/mikefarah/yq)
 - [Make](https://www.gnu.org/software/make/)
 - [Air](https://github.com/air-verse/air)
 
 This service is responsible for managing CRUD operations for users. It consists of several components:
 1. ***MongoDB Database***: Used for storing user data.
 2. ***HTTP Server***: Runs on port 8081 to handle HTTP requests.
 3. ***gRPC Server***: Runs on port 8091 to handle gRPC requests to get a user.
 4. ***RabbitMQ Producer***: Sends owner updates to the "user-updates" queue.
 
 To start the database, use the command:
 ```
 make service-up
 ```
 
 To bring down the database, use the command:
 ```
 make service-down
 ```
 
 To bring down the database and remove all volumes, use the command:
 ```
 make service-down-add
 ```
 
 Note: The RabbitMQ instance cannot be run from this service directory. It can only be started from the micro-order parent directory.
 
 To start the service, use one of the following commands:
 ```
 make air
 ```
 or
 ```
 go run cmd/http/main.go
 ```
 
 You can change configurations in the `config-sample.yml` file.
