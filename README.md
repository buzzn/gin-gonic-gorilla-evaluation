# Gin Gonic Api Mock
A minimal mockup-implementation of the [api-defintion]{https://github.com/buzzn/broker-evaluation/blob/master/api.md}.

## How to start
First install the requirements
```
go get github.com/gin-gonic/gin
go get github.com/gorilla/websocket
```

Then start the http server mock:
```
go run main.go
```
## How to use 
This implemention runs a http-server at port 8088. To get the hitlist type:
```
curl http://localhost/8088/hitlist
```

To get the live data: Read from websocket `ws://localhost/8088/live`

