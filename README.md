# Micro-Order

## Overview
The **Micro-Order** project is a microservice-based application that consists of three main services:
- **Owner-Service**
- **Product-Service**
- **Order-Service**

## Services

### Owner-Service
This service is responsible for managing the owners of the products. It includes functionalities such as:
- Creating a new owner (user)
- Updating owner details
- Deleting an owner
- Retrieving owner information

### Product-Service
This service handles all operations related to products. It includes functionalities such as:
- Adding a new product
- Updating product details
- Deleting a product
- Retrieving product information

### Order-Service
This service manages the orders placed by customers. It includes functionalities such as:
- Creating a new order
- Updating order details
- Deleting an order
- Retrieving order information

Detailed information can be found in each service README file.

## Getting Started

### Prerequisites
- [Go 1.20 or higher](https://go.dev/)
- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)
- [YQ - CLI-based yaml processor](https://github.com/mikefarah/yq)
- [Make](https://www.gnu.org/software/make/)
- [Air](https://github.com/air-verse/air)

### Installation
1. Clone the repository:
    ```sh
    git clone https://github.com/emmrys-jay/micro-order.git
    ```
2. Navigate to the project directory:
    ```sh
    cd micro-order
    ```
3. Start all required containers for databases and cache using Docker Compose:
    ```sh
    make service-up
    ```
4. Navigate to the home directory of each of the services in different terminals and start them in development mode
    ```sh
    cd owner-service
    make air

    cd product-service
    make air

    cd order-service
    make air
    ```
5. All the default configs are being used, and they can be found in `config-sample.yml` in each service 

### Usage
- Access the services via their respective endpoints:
  - Owner-Service: `http://localhost:8081`
  - Product-Service: `http://localhost:8082`
  - Order-Service: `http://localhost:8080`

### TODO
- [ ] Add mechanism to handle consecutive consumer errors when handling events.
- [ ] Security fixes
- [ ] Swagger documentation for each service

## Contributing
Contributions are welcome! Please follow these steps:
1. Fork the repository.
2. Create a new branch (`git checkout -b feature/your-feature`).
3. Commit your changes (`git commit -am 'Add new feature'`).
4. Push to the branch (`git push origin feature/your-feature`).
5. Create a new Pull Request.

## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contact
For any inquiries or issues, please contact:
- Emmrys Jay - [jonathanemma121@gmail.com](jonathanemma121@gmail.com)
