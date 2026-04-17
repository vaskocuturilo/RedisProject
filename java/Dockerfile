FROM gradle:jdk24 AS builder

WORKDIR /app

COPY --chown=gradle:gradle . .

RUN gradle clean build -x test --no-daemon

FROM eclipse-temurin:24-jre-alpine

WORKDIR /app

COPY --from=builder /app/build/libs/*.jar app.jar

ENTRYPOINT ["java", "-Dspring.profiles.active=docker", "-jar", "app.jar"]