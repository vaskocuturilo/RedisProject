package com.example.api.service;

import com.example.api.dto.EventDto;
import com.example.api.entity.EventJpaEntity;
import com.example.api.entity.EventRedisEntity;
import com.example.api.mapper.EventMapper;
import com.example.api.repository.EventJpaRepository;
import com.example.api.repository.EventRedisRepository;
import jakarta.persistence.EntityNotFoundException;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import org.springframework.util.Assert;

import java.util.Collections;
import java.util.List;
import java.util.UUID;

@Service
@Slf4j
@RequiredArgsConstructor
public class EventService {

    private final EventRedisRepository eventRedisRepository;
    private final EventMapper eventMapper;
    private final EventJpaRepository eventJpaRepository;

    public EventDto create(final EventDto dto) {
        final EventDto eventDto = new EventDto(UUID.randomUUID().toString(), dto.title(), dto.description());

        final EventJpaEntity savedEntity = eventJpaRepository.save(eventMapper.toJpaEntity(eventDto));

        eventRedisRepository.save(eventMapper.toRedisEntity(eventDto));

        log.info("IN create - saved user: {}", savedEntity);

        return eventMapper.toDto(savedEntity);
    }

    public EventDto get(final String id) {
        Assert.hasText(id, "Event id must not be null or blank");

        return eventRedisRepository
                .findById(id)
                .map(cached -> {
                    log.info("Cache hit for Event id= {}", id);

                    return eventMapper.toDto(cached);
                }).orElseGet(() ->
                        eventJpaRepository.findById(id)
                                .map(entity -> {
                                    final EventDto dto = eventMapper.toDto(entity);
                                    eventRedisRepository.save(eventMapper.toRedisEntity(dto));
                                    log.info("Event id = {} cached after loading from Postgres", id);
                                    return dto;
                                }).orElseThrow(() -> new EntityNotFoundException("Event is not found with id = %s".formatted(id))));
    }

    public List<EventDto> getAll() {
        List<EventRedisEntity> cachedEvents = (List<EventRedisEntity>) eventRedisRepository.findAll();

        if (!cachedEvents.isEmpty()) {
            log.info("Cache hit for all Events, count={}", cachedEvents.size());
            return cachedEvents.stream().map(eventMapper::toDto).toList();
        }

        log.info("Cache miss for all Events, loading from Postgres");

        List<EventJpaEntity> entities = (List<EventJpaEntity>) eventJpaRepository.findAll();

        if (entities.isEmpty()) {
            return Collections.emptyList();
        }

        List<EventDto> dtos = entities.stream()
                .map(eventMapper::toDto)
                .toList();

        List<EventRedisEntity> redisEntities = dtos.stream()
                .map(eventMapper::toRedisEntity)
                .toList();

        eventRedisRepository.saveAll(redisEntities);

        log.info("All Events cached after loading from Postgres, count={}", redisEntities.size());

        return dtos;

    }

    public EventDto update(final String id, final EventDto dto) {
        log.info("Updating  Event id = {} in Postgres and Redis", id);

        final EventDto eventDto = new EventDto(id, dto.title(), dto.description());

        final EventJpaEntity updatedEntity = eventJpaRepository.save(eventMapper.toJpaEntity(eventDto));

        eventRedisRepository.save(eventMapper.toRedisEntity(eventDto));

        return eventMapper.toDto(updatedEntity);
    }

    public void delete(String id) {
        log.info("Deleting  Event id = {} in Postgres and Redis", id);
        eventRedisRepository.deleteById(id);
        eventJpaRepository.deleteById(id);
    }
}
