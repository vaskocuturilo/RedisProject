package com.example.api.it;

import org.springframework.test.context.DynamicPropertyRegistry;
import org.springframework.test.context.DynamicPropertySource;
import org.testcontainers.containers.GenericContainer;
import org.testcontainers.containers.PostgreSQLContainer;
import org.testcontainers.junit.jupiter.Container;

public abstract class AbstractRestControllerBaseTest {

    @Container
    static final PostgreSQLContainer POSTGRES_SQL_CONTAINER;

    @Container
    static final GenericContainer<?> REDIS_CONTAINER;


    static {
        POSTGRES_SQL_CONTAINER = new PostgreSQLContainer("postgres:17")
                .withUsername("postgres")
                .withPassword("password")
                .withDatabaseName("events_testcontainers");


        REDIS_CONTAINER = new GenericContainer<>("redis:latest")
                .withExposedPorts(6379);


        POSTGRES_SQL_CONTAINER.start();
        REDIS_CONTAINER.start();
    }

    @DynamicPropertySource
    public static void dynamicPropertySource(final DynamicPropertyRegistry registry) {
        registry.add("spring.datasource.url", POSTGRES_SQL_CONTAINER::getJdbcUrl);
        registry.add("spring.datasource.username", POSTGRES_SQL_CONTAINER::getUsername);
        registry.add("spring.datasource.password", POSTGRES_SQL_CONTAINER::getPassword);

        registry.add("spring.data.redis.host", REDIS_CONTAINER::getHost);
        registry.add("spring.data.redis.port", REDIS_CONTAINER::getFirstMappedPort);
    }
}
