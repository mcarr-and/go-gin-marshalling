Test project to understand the 

[Go Gin framework](https://github.com/gin-gonic/gin#gin-web-framework)

## Project includes:

* Unit testing sending and receiving JSON to the server representing an in-memory music store.
  * Get all albums
  * Get album by ID
  * Get album by ID that is not found
  * Post to create new album
  * Post album without required JSON values to be valid 
* Benchmark tests for throughput for all unit tests

## TODO
* Adding CI server integration
* Fuzz testing
* adding a database 
* Docker container
* K3D cluster to run Docker container in K8s.
* Helm chart to add Gin Server and Database
* Skaffold to setup the K8s server for this project.

## Source Site for creating this project

Golang tutorial for Gin music store: https://go.dev/doc/tutorial/web-service-gin. 

Go does JSON marshalling and binding in Gin: https://blog.logrocket.com/gin-binding-in-go-a-tutorial-with-examples/

Go Gin testing: https://semaphoreci.com/community/tutorials/test-driven-development-of-go-web-applications-with-gin

Test benchmarking: https://blog.logrocket.com/benchmarking-golang-improve-function-performance/

Gin Examples: https://github.com/gin-gonic/examples