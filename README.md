Test project to understand the 

[Go Gin framework](https://github.com/gin-gonic/gin#gin-web-framework)

Project is a music store with an in-memory database.

## Project includes:

* Unit testing sending and receiving JSON
  * Get all albums
  * Get album by ID
  * Get album by ID that is not found
  * Post to create new album
  * Post album without all the required JSON fields to be a valid object to Gin & V10 
* Benchmark tests for throughput for all unit tests

## TODO
* Adding CI server integration
* Fuzz testing
* Use a database as a data store
* Database migration
* Adding project to a Docker container
* K3D cluster to run Docker container in K8s.
* Helm chart to add Gin Server and Database
* Skaffold to set up the K8s cluster for this project.

## Source Sites for creating this project

Golang tutorial for Gin music store: https://go.dev/doc/tutorial/web-service-gin. 

Go does JSON marshalling and binding in Gin: https://blog.logrocket.com/gin-binding-in-go-a-tutorial-with-examples/

Go Gin testing: https://semaphoreci.com/community/tutorials/test-driven-development-of-go-web-applications-with-gin

Test benchmarking: https://blog.logrocket.com/benchmarking-golang-improve-function-performance/

Gin Examples: https://github.com/gin-gonic/examples