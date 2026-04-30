##     

The demo project for Redis is implemented with Golang and Java.

### This project was created to demonstrate how to work with Redis.

- Create a REST API event service (create/put/get/get all/delete events and save to Redis).
- Add Cache-Aside.
- Add Rate Limiter.
- Add Distributed Lock.
- Add Redis.
- Add Redis UI.
- Add Docker file.
- Add Docker Compose.
- Add unit(repository, service and rest) and integration tests with testcontainers.

You will need the following technologies available to try it out:

* Git
* Spring Boot 3+
* Gradle 9+
* JDK 24+
* Golang 1.25+
* Redis 7.4+
* Redis Insight (Redis UI)
* Docker
* Docker compose
* IDE of your choice

### How to run via Spring Boot.

``` ./gradlew bootRun ```

``` docker compose -f "docker-compose-java.yml" up --detach ```

### How to run via Golang.

``` go run .```

``` docker compose -f "docker-compose-golang.yml" up --detach ```
