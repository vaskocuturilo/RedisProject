package com.example.api.service;

import com.example.api.dto.EventDto;
import com.example.api.entity.EventJpaEntity;
import com.example.api.entity.EventRedisEntity;
import com.example.api.mapper.EventMapper;
import com.example.api.repository.EventJpaRepository;
import com.example.api.repository.EventRedisRepository;
import jakarta.persistence.EntityNotFoundException;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import utils.DataUtils;

import java.util.Optional;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.BDDMockito.*;
import static org.mockito.Mockito.never;
import static org.mockito.Mockito.times;

@ExtendWith(MockitoExtension.class)
class EventServiceTests {

    @InjectMocks
    private EventService eventService;

    @Mock
    private EventJpaRepository eventJpaRepository;

    @Mock
    private EventRedisRepository eventRedisRepository;

    @Mock
    private EventMapper eventMapper;

    @Test
    @DisplayName("Test event save functionality")
    void givenEventObject_whenSave_thenEventSaved() {
        // given
        final EventDto inputDto = DataUtils.getEvent1DtoTransient();

        final EventJpaEntity jpaEntity = DataUtils.getEventEntityTransient();
        final EventJpaEntity savedEntity = DataUtils.getEvent1EntityPersisted();
        final EventRedisEntity redisEntity = DataUtils.getEventRedisEntityTransient();
        final EventDto expectedDto = DataUtils.getEvent1DtoTransient();

        given(eventMapper.toJpaEntity(any(EventDto.class))).willReturn(jpaEntity);
        given(eventJpaRepository.save(jpaEntity)).willReturn(savedEntity);
        given(eventMapper.toRedisEntity(any(EventDto.class))).willReturn(redisEntity);
        given(eventMapper.toDto(savedEntity)).willReturn(expectedDto);

        // when
        final EventDto result = eventService.create(inputDto);

        // then
        assertThat(result).isNotNull();
        assertThat(result.title()).isEqualTo(expectedDto.title());
        assertThat(result.description()).isEqualTo(expectedDto.description());
        assertThat(result.id()).isNotNull();

        then(eventJpaRepository).should(times(1)).save(jpaEntity);
        then(eventRedisRepository).should(times(1)).save(redisEntity);
        then(eventMapper).should(times(1)).toDto(savedEntity);
    }

    @Test
    @DisplayName("Test return DTO from Redis cache when cache hit")
    void givenCachedEvent_whenGet_thenReturnFromRedis() {
        // given
        final String id = "test-id";
        final EventRedisEntity redisEntity = DataUtils.getEventRedisEntityTransient();
        final EventDto expectedDto = DataUtils.getEvent1DtoTransient();

        given(eventRedisRepository.findById(id)).willReturn(Optional.of(redisEntity));
        given(eventMapper.toDto(redisEntity)).willReturn(expectedDto);

        // when
        final EventDto result = eventService.get(id);

        // then
        assertThat(result).isNotNull();
        assertThat(result.id()).isEqualTo(expectedDto.id());
        assertThat(result.title()).isEqualTo(expectedDto.title());
        assertThat(result.description()).isEqualTo(expectedDto.description());

        then(eventRedisRepository).should(times(1)).findById(id);
        then(eventMapper).should(times(1)).toDto(redisEntity);
        then(eventJpaRepository).should(never()).findById(any());
    }

    @Test
    @DisplayName("Test return DTO from JPA and cache it in Redis when cache miss")
    void givenNoCachedEvent_whenGet_thenReturnFromJpaAndCache() {
        // given
        final String id = "test-id";
        final EventJpaEntity jpaEntity = DataUtils.getEvent1EntityPersisted();
        final EventDto expectedDto = DataUtils.getEvent1DtoTransient();
        final EventRedisEntity redisEntity = DataUtils.getEventRedisEntityTransient();

        given(eventRedisRepository.findById(id)).willReturn(Optional.empty());
        given(eventJpaRepository.findById(id)).willReturn(Optional.of(jpaEntity));
        given(eventMapper.toDto(jpaEntity)).willReturn(expectedDto);
        given(eventMapper.toRedisEntity(expectedDto)).willReturn(redisEntity);

        // when
        final EventDto result = eventService.get(id);

        // then
        assertThat(result).isNotNull();
        assertThat(result.id()).isEqualTo(expectedDto.id());
        assertThat(result.title()).isEqualTo(expectedDto.title());
        assertThat(result.description()).isEqualTo(expectedDto.description());

        then(eventRedisRepository).should(times(1)).findById(id);
        then(eventJpaRepository).should(times(1)).findById(id);
        then(eventMapper).should(times(1)).toDto(jpaEntity);
        then(eventRedisRepository).should(times(1)).save(redisEntity);
    }

    @Test
    @DisplayName("Test throw EntityNotFoundException when event not found in Redis or JPA")
    void givenNoEvent_whenGet_thenThrowEntityNotFoundException() {
        // given
        final String id = "non-existing-id";

        given(eventRedisRepository.findById(id)).willReturn(Optional.empty());
        given(eventJpaRepository.findById(id)).willReturn(Optional.empty());

        // when / then
        assertThatThrownBy(() -> eventService.get(id))
                .isInstanceOf(EntityNotFoundException.class)
                .hasMessageContaining(id);

        then(eventRedisRepository).should(times(1)).findById(id);
        then(eventJpaRepository).should(times(1)).findById(id);
        then(eventRedisRepository).should(never()).save(any());
        then(eventMapper).should(never()).toDto(any(EventJpaEntity.class));
    }

    @Test
    @DisplayName("Test update event functionality - success")
    void givenExistingEvent_whenUpdate_thenReturnUpdatedDto() {
        // given
        final String id = "test-id";
        final EventDto inputDto = DataUtils.getEvent1DtoTransient();
        final EventJpaEntity updatedEntity = DataUtils.getEvent1EntityPersisted();
        final EventRedisEntity redisEntity = DataUtils.getEventRedisEntityTransient();
        final EventDto expectedDto = DataUtils.getEvent1DtoTransient();

        given(eventJpaRepository.existsById(id)).willReturn(true);
        given(eventMapper.toJpaEntity(any(EventDto.class))).willReturn(updatedEntity);
        given(eventJpaRepository.save(updatedEntity)).willReturn(updatedEntity);
        given(eventMapper.toRedisEntity(any(EventDto.class))).willReturn(redisEntity);
        given(eventMapper.toDto(updatedEntity)).willReturn(expectedDto);

        // when
        final EventDto result = eventService.update(id, inputDto);

        // then
        assertThat(result).isNotNull();
        assertThat(result.id()).isEqualTo(expectedDto.id());
        assertThat(result.title()).isEqualTo(expectedDto.title());
        assertThat(result.description()).isEqualTo(expectedDto.description());

        then(eventJpaRepository).should(times(1)).existsById(id);
        then(eventJpaRepository).should(times(1)).save(updatedEntity);
        then(eventRedisRepository).should(times(1)).save(redisEntity);
        then(eventMapper).should(times(1)).toDto(updatedEntity);
    }

    @Test
    @DisplayName("Test update event functionality - event not found")
    void givenNonExistingEvent_whenUpdate_thenThrowEntityNotFoundException() {
        // given
        final String id = "non-existing-id";
        final EventDto inputDto = DataUtils.getEvent1DtoTransient();

        given(eventJpaRepository.existsById(id)).willReturn(false);

        // when / then
        assertThatThrownBy(() -> eventService.update(id, inputDto))
                .isInstanceOf(EntityNotFoundException.class)
                .hasMessageContaining("The event not found");

        then(eventJpaRepository).should(times(1)).existsById(id);
        then(eventJpaRepository).should(never()).save(any());
        then(eventRedisRepository).should(never()).save(any());
        then(eventMapper).should(never()).toJpaEntity(any());
    }

    @Test
    @DisplayName("Test delete event functionality - success")
    void givenExistingEvent_whenDelete_thenDeleteFromJpaAndRedis() {
        // given
        final String id = "test-id";

        given(eventJpaRepository.existsById(id)).willReturn(true);
        willDoNothing().given(eventRedisRepository).deleteById(id);
        willDoNothing().given(eventJpaRepository).deleteById(id);

        // when
        eventService.delete(id);

        // then
        then(eventJpaRepository).should(times(1)).existsById(id);
        then(eventRedisRepository).should(times(1)).deleteById(id);
        then(eventJpaRepository).should(times(1)).deleteById(id);
    }

    @Test
    @DisplayName("Test delete event functionality - event not found")
    void givenNonExistingEvent_whenDelete_thenThrowEntityNotFoundException() {
        // given
        final String id = "non-existing-id";

        given(eventJpaRepository.existsById(id)).willReturn(false);

        // when
        assertThatThrownBy(() -> eventService.delete(id))
                .isInstanceOf(EntityNotFoundException.class)
                .hasMessageContaining("The event not found");

        // then
        then(eventJpaRepository).should(times(1)).existsById(id);
        then(eventRedisRepository).should(never()).deleteById(any());
        then(eventJpaRepository).should(never()).deleteById(any());
    }

    @Test
    @DisplayName("Test throw exception when id is null")
    void givenNullId_whenGet_thenThrowException() {
        assertThatThrownBy(() -> eventService.get(null))
                .isInstanceOf(IllegalArgumentException.class);

        then(eventRedisRepository).should(never()).findById(any());
        then(eventJpaRepository).should(never()).findById(any());
    }

    @Test
    @DisplayName("Test throw exception when id is blank")
    void givenBlankId_whenGet_thenThrowException() {
        assertThatThrownBy(() -> eventService.get("   "))
                .isInstanceOf(IllegalArgumentException.class);

        then(eventRedisRepository).should(never()).findById(any());
        then(eventJpaRepository).should(never()).findById(any());
    }
}