grpcgateway: A basic REST gateway forwarding requests onto services using grpc.
products: A very bare-bones product service
trade: A very bare-bones trade service

If you want to run the micro services directly, please make sure your docker start.
And then, please run the docker-compose.yml.
After all services run, you can test them using the below curl.

1. Get products directly
curl -X GET http://localhost:8080/v1/products

{"products":[]}

2. Create a product
curl -X POST http://localhost:8080/v1/products -d '{"name":"New Product 1", "price":"49.99"}'

3. Get products again
curl -X GET http://localhost:8080/v1/products

{"products":[{"id":"1","name":"New Product 1","price":"49.99"}]}%   

4. Create Sales No Discount
curl -X POST http://localhost:8080/v1/sales -d '{"lineItems": [{"productId": "1", "quantity": 2}]}'
   
{"saleId":"","lineItems":[{"productId":"1","quantity":2}],"totalPrice":99.98}

5. Create Sales With Discount
curl -X POST http://localhost:8080/v1/sales -d '{"lineItems": [{"productId": "1", "quantity": 2}], "discountAmount":10}'

{"saleId":"","lineItems":[{"productId":"1","quantity":2}],"totalPrice":89.98}

6. Create Sales With Discount
curl -X POST http://localhost:8080/v1/sales -d '{"lineItems": [{"productId": "1", "quantity": 2}], "discountAmount":100}'

{"saleId":"","lineItems":[{"productId":"1","quantity":2}],"totalPrice":0}