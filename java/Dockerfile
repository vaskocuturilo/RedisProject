FROM gradle:jdk-21-and-24 AS builder

WORKDIR /app

COPY --chown=gradle:gradle . .

RUN gradle clean build -x test --no-daemon

FROM eclipse-temurin:24-jdk-alpine

WORKDIR /app

COPY --from=builder /app/build/libs/api-1.0.0.jar app.jar

ENTRYPOINT ["java", "-jar","-Dspring.profiles.active=docker", "app.jar"]