/**
 * Order Service Documentation
 * 
 * 
 *   Prerequisites
 * - [Go 1.20 or higher](https://go.dev/)
 * - [Docker](https://www.docker.com/)
 * - [Docker Compose](https://docs.docker.com/compose/)
 * - [YQ - CLI-based yaml processor](https://github.com/mikefarah/yq)
 * - [Make](https://www.gnu.org/software/make/)
 * - [Air](https://github.com/air-verse/air)
 * 
 * This service is responsible for managing CRUD operations for orders. It consists of several components:
 * 1. MongoDB Database: Used for storing product order details.
 * 2. Redis cache: Used for caching users and products fetched from other services to prevent always making those calls.
 * 3. HTTP Server: Runs on port 8080 to handle HTTP requests.
 * 4. gRPC Client 1: Fetches user information to verify a user exists during creation of an order.
 * 5. gRPC Client 2: Fetches products when creating an order.
 * 6. RabbitMQ Consumer: Receives products updated from the "product-updates" queue.
 * 
 * To start the database, use the command:
 * ```
 * make service-up
 * ```
 * 
 * To bring down the database, use the command:
 * ```
 * make service-down
 * ```
 * 
 * To bring down the database and remove all volumes, use the command:
 * ```
 * make service-down-add
 * ```
 * 
 * Note: The RabbitMQ instance cannot be run from this service directory. It can only be started from the micro-order parent directory.
 * 
 * To start the service, use one of the following commands:
 * ```
 * make air
 * ```
 * or
 * ```
 * go run cmd/http/main.go
 * ```
 * 
 * You can change configurations in the `config-sample.yml` file.
 */