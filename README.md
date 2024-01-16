# Used Technologies:
- golang
- GIN
- zap
- postgresql
- go testing
- testify
- docker

# Run postgresql (with docker)
_[you need to install docker before](https://docs.docker.com/engine/install/)_
```shell
docker run --name gomsg_db -p 5432:5432 -e POSTGRES_USER=admin -e POSTGRES_PASSWORD=admin -e -d postgres:16
```
You can also add ```-v < your local path >:/var/lib/postgresql/data ``` to store database locally and have access to it after container is stopped (change ```< your local path >``` to your local path)

# Run (with docker)
_[you need to install docker before](https://docs.docker.com/engine/install/)_
```shell
git clone https://github.com/Anton-beep/gomsg
cd gomsg
docker build -t gomsg .
docker run -p 8080:8080 --name gomsg --link gomsg_db:gomsg_db gomsg
```
running in localhost:8080

# Run (local)
_[you need to install go before](https://go.dev/doc/install)_
```shell
go run main.go
```
running in localhost:8080
