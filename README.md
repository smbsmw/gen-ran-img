# Generate Random Image

Microservice architecture using REST + gRPC + Kafka

## Quick start
Local development:
```sh
# Run services in docker
$ make compose-up
```

Swagger is available on `http://localhost:8081/swagger/index.html`


## Example

Make POST request
```
curl -X 'POST' \
  'http://localhost:8081/api/v1/generate_image' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{ "image_name": "my_pic"}'
  
# If ok, returns the same value "my_pic"
# Сreates random image with suffix name "my_pic" in Minio
# Minio http://localhost:9000 minioadmin:minioadmin 
```





