package com.example.api;

import com.example.api.it.AbstractRestControllerBaseTest;
import com.example.api.rest.EventControllerV1;
import com.example.api.service.EventService;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.test.context.ActiveProfiles;

import static org.assertj.core.api.AssertionsForClassTypes.assertThat;

@ActiveProfiles("test")
@SpringBootTest
class ApiApplicationTests extends AbstractRestControllerBaseTest {

    @Autowired
    private EventControllerV1 eventControllerV1;

    @Autowired
    private EventService eventService;

    @Test
    void contextLoads() {
        assertThat(eventControllerV1).isNotNull();
        assertThat(eventService).isNotNull();
    }
}
