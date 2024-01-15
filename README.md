# Used Technologies:
- golang
- GIN
- zap
- postgresql
- go testing
- testify
- docker

# Run postgresql (with docker)
```shell
docker run --name gomsg_db -p 5432:5432 -e POSTGRES_USER=admin -e POSTGRES_PASSWORD=admin -e -d postgres:16
```
You can also add ```-v < your local path >:/var/lib/postgresql/data ``` to store database locally and have access to it after container is stopped (change ```< your local path >``` to your local path)

# Run (with docker)
```shell
docker-compose up -d
```
running in localhost:8080

# Run (local, need go installed)
```shell
go run main.go
```
running in localhost:8080
