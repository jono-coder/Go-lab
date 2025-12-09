# NOTES

## A lab type project depicting Golang functionality

### The 'best approach' has tried to be utilised. Simplicity is Go's middle name. 

* A database connection [sqLite; connection pooling]
* A database "Entity"
* A database "Repository"
* A Business "Service"
* A REST Client [includes oauth]
* A REST Server [including Cache-Control]
* Shutdown gracefully using a channel
* A "Service" type registry [using for scheduling but can be used for anythign that can be started and stopped]
* A "Worker Pool" [batches of threads to limit Transactions for the connection pool; any error and the pool dies]
* Using Docker
* Build scripts for tests, build and Docker

### Dependencies

* Chi -- a router for a REST server
* Resty -- a REST client
* Ristretto -- a comprehensive caching solution
* SqLite -- a local database
* Testify -- a mocking testing solution
* Unrolled Secure -- a secure handler for the REST server
* X-Oauth2 -- for oauth
