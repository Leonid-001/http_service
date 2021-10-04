
# My First Go Project  

This solution is feature complete. However, it should undergo Load and Stability testing before production deployments.  

## Endpoints    

### Hash Password  

HTTP POST request to http://localhost:8080/hash  
Immediately returns the request id. The hashed password is not available for 5 seconds.

```shell
curl --data “password=angryMonkey” http://localhost:8080/hash
42
```   

### Get Hashed Password  

HTTP GET request to http://localhost:8080/hash/###  
Returns the hashed password for the request ### if available.

```shell
curl http://localhost:8080/hash/42  
ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q==
``` 

```shell
curl http://localhost:8080/hash/333  
``` 

### Get statistics  

HTTP Get request to http://localhost:8080/stats  
Returns the statistics.

```shell
curl http://localhost:8080/stats
{“total”: 5, “average”: 943}
```  

### Shutdown service  

HTTP Get request to http://localhost:8080/shutdown  
Initiates a graceful shutdown.

```shell
curl http://localhost:8080/shutdown
```  

## Run and Test Instructions  

Run the HTTP service:

```shell
go run .
```

Test the HTTP service:

```shell
go test
```

## Obvious Project Improvements  

1. Implement configuration options
1. Remove the dependency on the total request processing time in the running average calculations
1. Sanitize the API inputs
1. Implement better error handling
1. Implement better logging solution
1. This solution doesn't support scaling out.

## Runtime Limitations  

1. The total number of hash password requests is limited to 18446744073709551615
1. The total request processing time is limited to 9223372036854775807 
1. Traffic spikes might cause issues
1. This solution doesn't support scaling out.
