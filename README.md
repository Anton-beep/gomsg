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
```

change file pkg/db/DB.go on the line 24: change localhost to your local ip (starts with 192.168 and could be found by ```ipconfig``` in windows or ```ifconfig``` in linux, for example 192.168.56.1)

````shell
docker build -t gomsg .
docker run -p 8080:8080 -d --name gomsg gomsg
````
running in localhost:8080

# Run (local)
_[you need to install go before](https://go.dev/doc/install)_
```shell
go run main.go
```
running in localhost:8080
