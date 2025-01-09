goctl api plugin -plugin goctl-swagger="swagger -filename workflow.json -host 127.0.0.1 -basepath /api -schemes https,wss" -api ../workflow.api -dir .

docker run --rm -p 8083:8080 -e SWAGGER_JSON=workflow.json -v $PWD:/foo swaggerapi/swagger-ui

https://app.swaggerhub.com/apis/1227694865/your-api/1.0.0