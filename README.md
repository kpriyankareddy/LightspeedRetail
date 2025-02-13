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
