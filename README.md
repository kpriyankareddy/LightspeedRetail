# LightspeedRetail

Running the program: we can simple run the main.go file to start the server on host 127.0.0.1:8080

## Testing stage1
command: curl.exe -X GET http://127.0.0.1:8080/products
output: [{"id":"58a69d6b-d831-422a-ac55-f5705d186be2","name":"Chrome Toaster","price":100},{"id":"7f39bb69-3d52-49e3-9b9a-9acf5d2be3aa","name":"Copper Kettle","price":49.99},{"id":"5447d643-9761-471a-9b27-21d5b43d2629","name":"Mixing Bowl","price":20}]

## Testing stage2
command: curl.exe -X POST http://127.0.0.1:8080/products -H "Content-Type: application/json" -d "{\"name\":\"Electric Kettle\",\"price\":35.50}"
output: [{"id":"58a69d6b-d831-422a-ac55-f5705d186be2","name":"Chrome Toaster","price":100},{"id":"7f39bb69-3d52-49e3-9b9a-9acf5d2be3aa","name":"Copper Kettle","price":49.99},{"id":"5447d643-9761-471a-9b27-21d5b43d2629","name":"Mixing Bowl","price":20},
{"id":"35c36797-4fa2-4434-9c3b-443a58c7547e","name":"Electric Kettle","price":35.5}]

command2: curl.exe -X POST http://127.0.0.1:8080/products -H "Content-Type: application/json" -d "{}"
output: {"error":"Invalid product details: name is required"}

command3: curl.exe -X POST http://127.0.0.1:8080/products -H "Content-Type: application/json" -d "{\"name\":\"\",\"price\":20.00}"
output: {"error":"Invalid product details: name is required"}

command4: curl.exe -X POST http://127.0.0.1:8080/products -H "Content-Type: application/json" -d "{\"name\":\"Faulty Item\",\"price\":-20}"
output: {"error":"Invalid product details: price must be greater than zero"}

command5: curl.exe -X POST http://127.0.0.1:8080/products -H "Content-Type: application/json" -d "{\"name\":\"Gift Item\",\"price\":0}"
output: {"error":"Invalid product details: price must be greater than zero"}

command6: curl.exe -X POST http://127.0.0.1:8080/products -H "Content-Type: application/json" -d "{\"name\":\"Wireless Charger\"}"
output: {"error":"Invalid product details: price must be greater than zero"}

command7: curl.exe -X POST http://127.0.0.1:8080/products -H "Content-Type: application/json" -d "{invalid-json}"
output: {"error":"Invalid request body"}

## Testing stage3
command1: curl.exe -X POST http://127.0.0.1:8080/sales -H "Content-Type: application/json" -d "{\"items\":[{\"product_id\":\"eea449a2-5cc7-4fe1-829a-be8dbf3c7091\",\"quantity\":2}]}"
output: {"items":[{"product_id":"eea449a2-5cc7-4fe1-829a-be8dbf3c7091","quantity":2,"total":99.98}],"total":99.98}

command2: curl.exe -X POST http://127.0.0.1:8080/sales -H "Content-Type: application/json" -d "{\"items\":[]}"
output: {"error":"Sale must contain at least one item"}

command3: curl.exe -X POST http://127.0.0.1:8080/sales -H "Content-Type: application/json" -d "{\"items\":[{\"product_id\":\"00000000-0000-0000-0000-000000000000\",\"quantity\":1}]}"
output: {"error":"Product not found"}

command4: curl.exe -X POST http://127.0.0.1:8080/sales -H "Content-Type: application/json" -d "{\"items\":[{\"product_id\":\"550e8400-e29b-41d4-a716-446655440000\",\"quantity\":-1}]}"
output: {"error":"Quantity must be a positive integer"}

command5: curl.exe -X POST http://127.0.0.1:8080/sales -H "Content-Type: application/json" -d "{\"items\":[{\"product_id\":\"550e8400-e29b-41d4-a716-446655440000\",\"quantity\":1.5}]}"
output: {"error":"Invalid request body"}

command6: curl.exe -X POST http://127.0.0.1:8080/sales -H "Content-Type: application/json" -d "{\"items\":[{\"product_id\":\"550e8400-e29b-41d4-a716-446655440000\"}]}"
output: {"error":"Quantity must be a positive integer"}

command7: curl.exe -X POST http://127.0.0.1:8080/sales -H "Content-Type: application/json" -d "{\"items\":[{\"product_id\":\"a767983f-9d79-4d19-9fa3-582ed7ba3d83\",\"quantity\":2},{\"product_id\":\"c790aab5-3e7e-4028-b92e-c6fe6711c39b\",\"quantity\":3}]}"
output: {"items":[{"product_id":"a767983f-9d79-4d19-9fa3-582ed7ba3d83","quantity":2,"total":200},{"product_id":"c790aab5-3e7e-4028-b92e-c6fe6711c39b","quantity":3,"total":60}],"total":260}

## Testing stage4
command: {"items":[{"product_id":"7a4cfaa7-b0f2-4fca-b1df-86efedbd9e14","quantity":2,"total":99.98,"discount":4.99},{"product_id":"e0544a89-668f-4da0-861d-4ae72c40ef6c","quantity":3,"total":300,"discount":15.01}],"total":379.98,"discount":20}
output: {"items":[{"product_id":"7a4cfaa7-b0f2-4fca-b1df-86efedbd9e14","quantity":2,"total":99.98,"discount":4.99},{"product_id":"e0544a89-668f-4da0-861d-4ae72c40ef6c","quantity":3,"total":300,"discount":15.01}],"total":379.98,"discount":20}

### Discount Application Logic
1. Total Discount is Distributed Proportionally  
   - Each item's discount share is calculated as:  
     \[
     \text{Item Discount} = \left( \frac{\text{Item Total}}{\text{Sale Total}} \right) \times \text{Discount Amount}
     \]
   - Example:  
     ```
     Total Sale Amount: $100  
     Discount: $20  
     Item A Total: $40  
     Item B Total: $60  
     ```
     - Item A Discount = (40/100) × 20 = **$8**
     - Item B Discount = (60/100) × 20 = **$12**

2. Avoiding Rounding Errors  
   - To ensure that the total discount **exactly matches** the requested amount:
     - The discount is **rounded down** for all items except the last one.
     - The last item is assigned any remaining discount to **fix rounding differences**.

