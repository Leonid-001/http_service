
# My First Go Project

## Endpoints Description

### Hash Password

POST request to http://localhost:8080/hash  
Immediately returns the request id. The hashed password is not available for 5 seconds.

```shell
curl --data “password=angryMonkey” http://localhost:8080/hash
42
```   

### Get Hashed Password

GET request to http://localhost:8080/hash/###  
Returns the hashed password if available.

```shell
curl --data “password=angryMonkey” http://localhost:8080/hash/42  
ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q==
``` 

### Get statistics

HTTP Get request to curl http://localhost:8080/stats  
Returns the statistics.

```shell
curl http://localhost:8080/stats
{“total”: 5, “average”: 943}
```  

### Shutdown service

HTTP Get request to curl http://localhost:8080/shutdown  
Initiates a graceful shutdown.

```shell
curl http://localhost:8080/shutdown
```  

## Run and Test Instructions

Run instructions

```shell
go run .
```

Test instructions

```shell
go test
```

## Obvious Improvements
1. Real configuration options
2. Better error handling
3. Real logging solution
4. Designed as a single process. This solution doesn't support scaling out.

## Runtime Limitations
1. The total number of POST requests is limited to 18446744073709551615
2. Traffic spikes will probably cause issues
3. Designed as a single process. This solution doesn't support scaling out.

