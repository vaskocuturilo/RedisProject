package com.example.api.repository;

import com.example.api.entity.EventJpaEntity;
import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.jdbc.AutoConfigureTestDatabase;
import org.springframework.boot.test.autoconfigure.orm.jpa.DataJpaTest;
import org.springframework.boot.testcontainers.service.connection.ServiceConnection;
import org.springframework.test.context.junit.jupiter.SpringExtension;
import org.springframework.util.CollectionUtils;
import org.testcontainers.containers.PostgreSQLContainer;
import org.testcontainers.junit.jupiter.Container;
import org.testcontainers.junit.jupiter.Testcontainers;
import utils.DataUtils;

import java.util.List;
import java.util.UUID;

import static org.assertj.core.api.Assertions.assertThat;

@ExtendWith(SpringExtension.class)
@DataJpaTest(properties = {
        "spring.jpa.properties.javax.persistence.validation.mode=none",
        "spring.jpa.hibernate.ddl-auto=create-drop"
})
@AutoConfigureTestDatabase(replace = AutoConfigureTestDatabase.Replace.NONE)
@Testcontainers
class EventJpaRepositoryTest {

    @Container
    @ServiceConnection
    static PostgreSQLContainer<?> postgres = new PostgreSQLContainer<>("postgres:17");

    @Autowired
    private EventJpaRepository eventJpaRepository;

    @BeforeEach
    void setUp() {
        eventJpaRepository.deleteAll();
    }

    @Test
    @DisplayName("Test save event functionality")
    void givenEventObject_whenSave_thenEventIsCreated() {
        //given
        final EventJpaEntity eventToCreate = DataUtils.getEventEntityTransient();

        //when
        final EventJpaEntity savedEntity = eventJpaRepository.save(eventToCreate);

        //then
        assertThat(savedEntity).isNotNull();
        assertThat(savedEntity.getId()).isNotNull();
    }

    @Test
    @DisplayName("Test get event by id functionality")
    void givenEventCreated_whenGet_thenEventIsReturned() {
        //given
        final EventJpaEntity eventToCreate = DataUtils.getEventEntityTransient();

        final EventJpaEntity savedEvent = eventJpaRepository.save(eventToCreate);

        //when
        final EventJpaEntity existEvent = eventJpaRepository.findById(savedEvent.getId()).orElse(null);

        //then
        assertThat(existEvent).isNotNull();
        assertThat(existEvent.getId()).isEqualTo(savedEvent.getId());
    }

    @Test
    @DisplayName("Test the event was not found functionality")
    void givenEventIsNotCreated_whenGetById_thenOptionalIsEmpty() {
        //given
        final UUID randomUUID = UUID.randomUUID();

        //when
        final EventJpaEntity existEvent = eventJpaRepository.findById(randomUUID.toString()).orElse(null);

        //then
        assertThat(existEvent).isNull();
    }

    @Test
    @DisplayName("Test find all events functionality")
    void givenEventsAreStored_whenGetAll_thenAllEventsAreReturned() {
        //given
        final EventJpaEntity eventJpaEntity = DataUtils.getEventEntityTransient();

        eventJpaRepository.saveAll(List.of(eventJpaEntity));

        //when
        final List<EventJpaEntity> developers = eventJpaRepository.findAll();

        //then
        assertThat(CollectionUtils.isEmpty(developers)).isFalse();
        assertThat(developers).hasSize(1).contains(eventJpaEntity);
    }

    @Test
    @DisplayName("Test the event update set title functionality")
    void givenEventsAreStored_whenUpdateTitle_thenEventUpdated() {

        final String updatedTitle = "Updated Title";

        //given
        final EventJpaEntity eventJpaEntity = DataUtils.getEventEntityTransient();

        eventJpaRepository.save(eventJpaEntity);

        final EventJpaEntity existEvent = eventJpaRepository.findById(eventJpaEntity.getId()).orElse(null);

        //when
        Assertions.assertNotNull(existEvent);
        existEvent.setTitle(updatedTitle);

        final EventJpaEntity updateEvent = eventJpaRepository.save(existEvent);

        //then
        assertThat(updateEvent).isNotNull();
        assertThat(updateEvent.getTitle()).isEqualTo(updatedTitle);
    }

    @Test
    @DisplayName("Test the event update set description functionality")
    void givenEventsAreStored_whenUpdateDescription_thenEventUpdated() {

        final String updatedDescription = "Updated Description";

        //given
        final EventJpaEntity eventJpaEntity = DataUtils.getEventEntityTransient();

        eventJpaRepository.save(eventJpaEntity);

        final EventJpaEntity existEvent = eventJpaRepository.findById(eventJpaEntity.getId()).orElse(null);

        //when
        Assertions.assertNotNull(existEvent);
        existEvent.setDescription(updatedDescription);

        final EventJpaEntity updateEvent = eventJpaRepository.save(existEvent);

        //then
        assertThat(updateEvent).isNotNull();
        assertThat(updateEvent.getDescription()).isEqualTo(updatedDescription);
    }

    @Test
    @DisplayName("Test the event deleted functionality")
    void givenEventsAreStored_whenDelete_thenEventDeleted() {
        //given
        final EventJpaEntity eventJpaEntity = DataUtils.getEventEntityTransient();

        eventJpaRepository.save(eventJpaEntity);

        //when
        eventJpaRepository.deleteById(eventJpaEntity.getId());

        final EventJpaEntity isNotExistDeveloper = eventJpaRepository.findById(eventJpaEntity.getId()).orElse(null);

        //then
        assertThat(isNotExistDeveloper).isNull();
    }
}