
- `grpcgateway`: A basic REST gateway forwarding requests onto services using grpc.
- `products`: A very bare-bones product service
- `trade`: A very bare-bones trade service

### Getting Started

If you want to run the micro services directly, please make sure your docker start.
And then, please run the docker-compose.yml.
After all services run, you can test them using the below curl.

1. Get products directly
```bash
curl -X GET http://localhost:8080/v1/products

➜ {"products":[]}
```

2. Create a product
```bash
curl -X POST http://localhost:8080/v1/products -d '{"name":"New Product 1", "price":"49.99"}'
```

3. Get products again
```bash
curl -X GET http://localhost:8080/v1/products

➜ {"products":[{"id":"1","name":"New Product 1","price":"49.99"}]}%   
```

4. Create Sales No Discount
```bash
curl -X POST http://localhost:8080/v1/sales -d '{"lineItems": [{"productId": "1", "quantity": 2}]}'
   
➜ {"saleId":"","lineItems":[{"productId":"1","quantity":2}],"totalPrice":99.98}
```

5. Create Sales With Discount
```bash
curl -X POST http://localhost:8080/v1/sales -d '{"lineItems": [{"productId": "1", "quantity": 2}], "discountAmount":10}'

➜ {"saleId":"","lineItems":[{"productId":"1","quantity":2}],"totalPrice":89.98}
```

6. Create Sales With Discount
```bash
curl -X POST http://localhost:8080/v1/sales -d '{"lineItems": [{"productId": "1", "quantity": 2}], "discountAmount":100}'

➜ {"saleId":"","lineItems":[{"productId":"1","quantity":2}],"totalPrice":0}
```


Codes structure:
```
sales/
├── grpcgateway/
│   ├── middleware/
│   │   └── ratelimit/
│   │       ├── impl.go
│   │       ├── ratelimit.go
│   │       └── ratelimit_test.go
│   ├── protos/
│   │   ├── google/
│   │   │   ├── api/
│   │   │   │   ├── annotations.proto
│   │   │   │   └── http.proto
│   │   ├── products/
│   │   │   ├── product.pb.go
│   │   │   ├── product.pb.gw.go
│   │   │   ├── product.proto
│   │   │   └── product_grpc.pb.go
│   │   └── trade/
│   │       ├── trade.pb.go
│   │       ├── trade.pb.gw.go
│   │       ├── trade.proto
│   │       └── trade_grpc.pb.go
│   ├── routes/
│   │   └── route.go
│   ├── Dockerfile
│   ├── gen-proto.sh
│   ├── go.mod
│   ├── go.sum
│   └── main.go
│
├── products/
│   ├── internal/
│   │   ├── repos/
│   │   │   └── product.go
│   │   └── services/
│   │       ├── product.go
│   │       └── product_test.go
│   ├── proto/
│   │   ├── product.pb.go
│   │   ├── product.proto
│   │   └── product_grpc.pb.go
│   ├── Dockerfile
│   ├── gen-proto.sh
│   ├── go.mod
│   ├── go.sum
│   └── main.go
│
├── trade/
│   ├── internal/
│   │   ├── repos/
│   │   │   └── trade.go
│   │   └── services/
│   │       ├── trade.go
│   │       └── trade_test.go
│   ├── proto/
│   │   ├── trade.pb.go
│   │   ├── trade.proto
│   │   └── trade_grpc.pb.go
│   ├── Dockerfile
│   ├── gen-proto.sh
│   ├── go.mod
│   ├── go.sum
│   └── main.go
│
├── docker-compose.yml
└── readme.md
```

- `sales/` is the root directory containing all your microservices.
- `grpcgateway/` contains the API gateway service, its middleware, and related files.
- `products/` and trade/ are directories for each respective service with their own internal/ logic, proto/ definitions, and Docker configurations.
- `protos/` inside grpcgateway/ contains the compiled protobuf files and the gateway definitions.
- Each service (grpcgateway, products, and trade) has a Dockerfile and a gen-proto.sh script for building the Docker image and generating protobuf files respectively.
- The docker-compose.yml file is in the root sales/ directory, which orchestrates the containers for all the services.
readme.md is also in the root sales/ directory, providing documentation for the entire project.