package com.example.api.rest;

import com.example.api.dto.EventDto;
import com.example.api.dto.PageResponse;
import com.example.api.service.EventService;
import com.fasterxml.jackson.databind.ObjectMapper;
import jakarta.persistence.EntityNotFoundException;
import org.hamcrest.CoreMatchers;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.mockito.BDDMockito;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.web.servlet.WebMvcTest;
import org.springframework.data.domain.Pageable;
import org.springframework.http.MediaType;
import org.springframework.test.context.ActiveProfiles;
import org.springframework.test.context.bean.override.mockito.MockitoBean;
import org.springframework.test.web.servlet.MockMvc;
import org.springframework.test.web.servlet.ResultActions;
import org.springframework.test.web.servlet.result.MockMvcResultHandlers;
import org.springframework.test.web.servlet.result.MockMvcResultMatchers;
import utils.DataUtils;

import java.util.List;

import static org.hamcrest.Matchers.hasSize;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.ArgumentMatchers.anyString;
import static org.mockito.BDDMockito.then;
import static org.mockito.Mockito.times;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.*;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.jsonPath;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.status;

@WebMvcTest(EventControllerV1.class)
@ActiveProfiles("test")
class EventControllerV1Test {

    @MockitoBean
    private EventService eventService;

    @Autowired
    private MockMvc mockMvc;

    @Autowired
    private ObjectMapper objectMapper;

    private static final String ENDPOINT_PATH = "/api/v1/events";

    @Test
    @DisplayName("Test create event functionality")
    void givenEventObject_whenCreate_thenSuccessResponse() throws Exception {
        //given
        final EventDto eventDto = DataUtils.getEvent1DtoTransient();

        BDDMockito.given(eventService.create(any(EventDto.class))).willReturn(DataUtils.getEvent1DtoTransient());

        //when
        final ResultActions result = mockMvc
                .perform(post(ENDPOINT_PATH)
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsBytes(eventDto)));

        //then
        result.andDo(MockMvcResultHandlers.print())
                .andExpect(status().isCreated())
                .andExpect(jsonPath("$.id").isNotEmpty())
                .andExpect(jsonPath("$.title").isNotEmpty())
                .andExpect(jsonPath("$.description").isNotEmpty())
                .andExpect(jsonPath("$.title", CoreMatchers.notNullValue()))
                .andExpect(jsonPath("$.title", CoreMatchers.is("Title Dto 1")))
                .andExpect(jsonPath("$.description", CoreMatchers.notNullValue()))
                .andExpect(jsonPath("$.description", CoreMatchers.is("Description Dto 1")));
    }

    @Test
    @DisplayName("Test update the developer functionality")
    void givenEventDto_whenUpdateEvent_thenSuccessResponse() throws Exception {
        //given
        final EventDto eventDto = DataUtils.getEvent1DtoTransient();

        BDDMockito.given(eventService.update(anyString(), any(EventDto.class))).willReturn(DataUtils.getEvent1DtoTransient());

        //when
        final ResultActions result = mockMvc
                .perform(put(ENDPOINT_PATH + "/" + eventDto.id())
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsBytes(eventDto)));

        //then
        result.andDo(MockMvcResultHandlers.print())
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andExpect(jsonPath("$.id").isNotEmpty())
                .andExpect(jsonPath("$.title").isNotEmpty())
                .andExpect(jsonPath("$.description").isNotEmpty())
                .andExpect(jsonPath("$.title", CoreMatchers.notNullValue()))
                .andExpect(jsonPath("$.title", CoreMatchers.is("Title Dto 1")))
                .andExpect(jsonPath("$.description", CoreMatchers.notNullValue()))
                .andExpect(jsonPath("$.description", CoreMatchers.is("Description Dto 1")));
    }

    @Test
    @DisplayName("Test update the event with incorrect id functionality")
    void givenEventDtoWithIncorrectId_whenUpdateEvent_thenErrorResponse() throws Exception {
        //given
        final EventDto developerDto = DataUtils.getEvent1DtoTransient();

        BDDMockito.given(eventService.update(anyString(), any(EventDto.class)))
                .willThrow(new EntityNotFoundException("The event not found"));

        //when
        final ResultActions result = mockMvc.perform(put(ENDPOINT_PATH + "/" + "1")
                .contentType(MediaType.APPLICATION_JSON)
                .content(objectMapper.writeValueAsBytes(developerDto)));
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
        final EventDto eventDto = DataUtils.getEvent1DtoTransient();

        BDDMockito.given(eventService.get(anyString())).willReturn(DataUtils.getEvent1DtoTransient());

        //when
        final ResultActions result = mockMvc.perform(get(ENDPOINT_PATH + "/" + eventDto.id())
                .contentType(MediaType.APPLICATION_JSON));

        //then
        result.andDo(MockMvcResultHandlers.print())
                .andExpect(MockMvcResultMatchers.status().isOk())
                .andExpect(jsonPath("$.id").isNotEmpty())
                .andExpect(jsonPath("$.title").isNotEmpty())
                .andExpect(jsonPath("$.description").isNotEmpty())
                .andExpect(jsonPath("$.title", CoreMatchers.notNullValue()))
                .andExpect(jsonPath("$.title", CoreMatchers.is("Title Dto 1")))
                .andExpect(jsonPath("$.description", CoreMatchers.notNullValue()))
                .andExpect(jsonPath("$.description", CoreMatchers.is("Description Dto 1")));
    }

    @Test
    @DisplayName("Test get by id with incorrect id functionality")
    void givenIncorrectId_whenGetById_thenErrorResponse() throws Exception {
        //given
        String id = "1";

        BDDMockito.given(eventService.get(anyString()))
                .willThrow(new EntityNotFoundException("Event is not found with id = %s".formatted(id)));

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
        PageResponse<EventDto> response = new PageResponse<>(
                List.of(
                        DataUtils.getEvent1DtoTransient(),
                        DataUtils.getEvent2DtoTransient(),
                        DataUtils.getEvent3DtoTransient()),
                0,
                10,
                3,
                1,
                "title",
                "ASC");

        BDDMockito.given(eventService.getAll(any(Pageable.class)))
                .willReturn(response);

        //when
        final ResultActions result = mockMvc.perform(get(ENDPOINT_PATH)
                .contentType(MediaType.APPLICATION_JSON));

        //then
        result.andDo(MockMvcResultHandlers.print())
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.content", hasSize(3)))
                .andExpect(jsonPath("$.content[0].id").isNotEmpty())
                .andExpect(jsonPath("$.content[0].title").isNotEmpty())
                .andExpect(jsonPath("$.content[0].description").isNotEmpty())

                .andExpect(jsonPath("$.page").value(0))
                .andExpect(jsonPath("$.size").value(10))
                .andExpect(jsonPath("$.totalElements").value(3))
                .andExpect(jsonPath("$.totalPages").value(1))
                .andExpect(jsonPath("$.sortBy").value("title"))
                .andExpect(jsonPath("$.direction").value("ASC"));
    }

    @Test
    @DisplayName("Test delete event by id functionality")
    void givenId_whenDelete_thenSuccessResponse() throws Exception {
        //given
        BDDMockito.doNothing().when(eventService).delete(anyString());

        //when
        final ResultActions result = mockMvc.perform(delete(ENDPOINT_PATH + "/" + 1)
                .contentType(MediaType.APPLICATION_JSON));

        //then
        then(eventService).should(times(1)).delete("1");
        result.andDo(MockMvcResultHandlers.print())
                .andExpect(MockMvcResultMatchers.status().isNoContent());
    }


    @Test
    @DisplayName("Test delete event with incorrect id functionality")
    void givenIncorrectId_whenDelete_thenErrorResponse() throws Exception {
        //given
        BDDMockito.doThrow(new EntityNotFoundException("The event not found")).when(eventService).delete(anyString());

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