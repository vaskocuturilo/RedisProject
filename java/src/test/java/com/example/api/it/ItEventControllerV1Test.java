package com.example.api.it;

import com.example.api.dto.EventDto;
import com.example.api.entity.EventJpaEntity;
import com.example.api.repository.EventJpaRepository;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.hamcrest.CoreMatchers;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.web.servlet.AutoConfigureMockMvc;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.data.redis.connection.RedisConnectionFactory;
import org.springframework.http.MediaType;
import org.springframework.test.context.ActiveProfiles;
import org.springframework.test.web.servlet.MockMvc;
import org.springframework.test.web.servlet.ResultActions;
import org.springframework.test.web.servlet.result.MockMvcResultHandlers;
import org.springframework.test.web.servlet.result.MockMvcResultMatchers;
import org.testcontainers.junit.jupiter.Testcontainers;
import utils.DataUtils;

import static org.assertj.core.api.Assertions.assertThat;
import static org.hamcrest.Matchers.hasSize;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.*;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.jsonPath;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.status;

@ActiveProfiles("test")
@AutoConfigureMockMvc
@SpringBootTest(webEnvironment = SpringBootTest.WebEnvironment.RANDOM_PORT)
@Testcontainers
class ItEventControllerV1Test extends AbstractRestControllerBaseTest {

    @Autowired
    private EventJpaRepository eventJpaRepository;

    @Autowired
    private MockMvc mockMvc;

    @Autowired
    private ObjectMapper objectMapper;

    @Autowired
    private RedisConnectionFactory redisConnectionFactory;

    private static final String ENDPOINT_PATH = "/api/v1/events";

    @BeforeEach
    void setUp() {
        eventJpaRepository.deleteAll();

        redisConnectionFactory.getConnection()
                .serverCommands()
                .flushAll();
    }

    @Test
    @DisplayName("Test create event functionality")
    void givenEventObject_whenCreate_thenSuccessResponse() throws Exception {
        //given
        final EventJpaEntity eventEntity = DataUtils.getEvent1EntityPersisted();

        //when
        final ResultActions result = mockMvc
                .perform(post(ENDPOINT_PATH)
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsBytes(eventEntity)));

        //then
        result.andDo(MockMvcResultHandlers.print())
                .andExpect(status().isCreated())
                .andExpect(jsonPath("$.id").isNotEmpty())
                .andExpect(jsonPath("$.title").isNotEmpty())
                .andExpect(jsonPath("$.description").isNotEmpty())
                .andExpect(jsonPath("$.title", CoreMatchers.notNullValue()))
                .andExpect(jsonPath("$.title", CoreMatchers.is("Title 1")))
                .andExpect(jsonPath("$.description", CoreMatchers.notNullValue()))
                .andExpect(jsonPath("$.description", CoreMatchers.is("Description 1")));

        assertThat(eventJpaRepository.count()).isEqualTo(1);

    }

    @Test
    @DisplayName("Test update the developer functionality")
    void givenEventDto_whenUpdateEvent_thenSuccessResponse() throws Exception {
        //given
        final EventJpaEntity existEvent = DataUtils.getEvent1EntityPersisted();

        eventJpaRepository.save(existEvent);

        final EventDto updateDto = new EventDto(null, "New Title", "New Description");

        //when
        final ResultActions result = mockMvc
                .perform(put(ENDPOINT_PATH + "/" + existEvent.getId())
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsBytes(updateDto)));

        //then
        result.andDo(MockMvcResultHandlers.print())
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andExpect(jsonPath("$.id").isNotEmpty())
                .andExpect(jsonPath("$.title").isNotEmpty())
                .andExpect(jsonPath("$.description").isNotEmpty())
                .andExpect(jsonPath("$.title", CoreMatchers.is("New Title")))
                .andExpect(jsonPath("$.description", CoreMatchers.is("New Description")));

        EventJpaEntity updated = eventJpaRepository.findById(existEvent.getId()).orElseThrow();
        assertThat(updated.getTitle()).isEqualTo("New Title");
    }

    @Test
    @DisplayName("Test update the event with incorrect id functionality")
    void givenEventDtoWithIncorrectId_whenUpdateEvent_thenErrorResponse() throws Exception {
        //given
        final EventJpaEntity eventJpaEntity = DataUtils.getEvent1EntityPersisted();

        eventJpaRepository.save(eventJpaEntity);

        //when
        final ResultActions result = mockMvc.perform(put(ENDPOINT_PATH + "/" + "1")
                .contentType(MediaType.APPLICATION_JSON)
                .content(objectMapper.writeValueAsBytes(eventJpaEntity)));
        //then
        result.andDo(MockMvcResultHandlers.print())
                .andExpect(MockMvcResultMatchers.status().isNotFound())
                .andExpect(jsonPath("$.status", CoreMatchers.is(404)))
                .andExpect(jsonPath("$.message", CoreMatchers.is("The event not found")));

    }

    @Test
    @DisplayName("Test get by id event functionality")
    void givenId_whenGetById_thenSuccessResponse() throws Exception {
        //given
        final EventJpaEntity eventJpaEntity = DataUtils.getEvent1EntityPersisted();

        eventJpaRepository.save(eventJpaEntity);

        //when
        final ResultActions result = mockMvc.perform(get(ENDPOINT_PATH + "/" + eventJpaEntity.getId())
                .contentType(MediaType.APPLICATION_JSON));

        //then
        result.andDo(MockMvcResultHandlers.print())
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andExpect(jsonPath("$.id").isNotEmpty())
                .andExpect(jsonPath("$.title").isNotEmpty())
                .andExpect(jsonPath("$.description").isNotEmpty())
                .andExpect(jsonPath("$.title", CoreMatchers.is("Title 1")))
                .andExpect(jsonPath("$.description", CoreMatchers.is("Description 1")));
    }

    @Test
    @DisplayName("Test get by id with incorrect id functionality")
    void givenIncorrectId_whenGetById_thenErrorResponse() throws Exception {
        //given
        String id = "1";

        //when
        final ResultActions result = mockMvc.perform(get(ENDPOINT_PATH + "/" + id)
                .contentType(MediaType.APPLICATION_JSON));
        //then
        result.andDo(MockMvcResultHandlers.print())
                .andExpect(MockMvcResultMatchers.status().isNotFound())
                .andExpect(jsonPath("$.status", CoreMatchers.is(404)))
                .andExpect(jsonPath("$.message", CoreMatchers.is("Event is not found with id = %s".formatted(id))));

    }

    @Test
    @DisplayName("Test get all events functionality")
    void givenThreeEvents_whenGetByAll_thenSuccessResponse() throws Exception {
        //given
        final EventJpaEntity eventJpaEntity1 = DataUtils.getEvent1EntityPersisted();
        final EventJpaEntity eventJpaEntity2 = DataUtils.getEvent2EntityPersisted();
        final EventJpaEntity eventJpaEntity3 = DataUtils.getEvent3EntityPersisted();

        eventJpaRepository.save(eventJpaEntity1);
        eventJpaRepository.save(eventJpaEntity2);
        eventJpaRepository.save(eventJpaEntity3);

        //when
        final ResultActions result = mockMvc.perform(get(ENDPOINT_PATH)
                .contentType(MediaType.APPLICATION_JSON));

        //then
        result.andDo(MockMvcResultHandlers.print())
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andExpect(jsonPath("$[*]", hasSize(7)))
                .andExpect(jsonPath("$.content[0].id").isNotEmpty())
                .andExpect(jsonPath("$.content[0].title").isNotEmpty())
                .andExpect(jsonPath("$.content[0].description").isNotEmpty())

                .andExpect(jsonPath("$.content[0].title", CoreMatchers.is("Title 1")))
                .andExpect(jsonPath("$.content[1].title", CoreMatchers.is("Title 2")))
                .andExpect(jsonPath("$.content[2].title", CoreMatchers.is("Title 3")));
    }

    @Test
    @DisplayName("Test delete event by id functionality")
    void givenId_whenDelete_thenSuccessResponse() throws Exception {
        //given
        EventJpaEntity eventJpaEntity = DataUtils.getEvent1EntityPersisted();

        eventJpaRepository.save(eventJpaEntity);

        //when
        final ResultActions result = mockMvc.perform(delete(ENDPOINT_PATH + "/" + eventJpaEntity.getId())
                .contentType(MediaType.APPLICATION_JSON));

        //then
        result.andDo(MockMvcResultHandlers.print())
                .andExpect(MockMvcResultMatchers.status().isNoContent());
        assertThat(eventJpaRepository.existsById(eventJpaEntity.getId())).isFalse();
    }

    @Test
    @DisplayName("Test delete event with incorrect id functionality")
    void givenIncorrectId_whenDelete_thenErrorResponse() throws Exception {
        //given

        //when
        final ResultActions result = mockMvc.perform(delete(ENDPOINT_PATH + "/" + 1)
                .contentType(MediaType.APPLICATION_JSON));

        //then
        result.andDo(MockMvcResultHandlers.print())
                .andExpect(MockMvcResultMatchers.status().isNotFound())
                .andExpect(jsonPath("$.status", CoreMatchers.is(404)))
                .andExpect(jsonPath("$.message", CoreMatchers.is("The event not found")));
    }
}