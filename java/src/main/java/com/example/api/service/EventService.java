package com.example.api.service;

import com.example.api.dto.EventDto;
import com.example.api.dto.PageResponse;
import com.example.api.entity.EventJpaEntity;
import com.example.api.mapper.EventMapper;
import com.example.api.repository.EventJpaRepository;
import com.example.api.repository.EventRedisRepository;
import jakarta.persistence.EntityNotFoundException;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.domain.Sort;
import org.springframework.stereotype.Service;
import org.springframework.util.Assert;

import java.util.UUID;

@Service
@Slf4j
@RequiredArgsConstructor
public class EventService {

    private final EventRedisRepository eventRedisRepository;
    private final EventMapper eventMapper;
    private final EventJpaRepository eventJpaRepository;

    public EventDto create(final EventDto dto) {
        log.debug("Creating event: {}", dto);

        final EventDto eventDto = new EventDto(UUID.randomUUID().toString(), dto.title(), dto.description());

        final EventJpaEntity savedEntity = eventJpaRepository.save(eventMapper.toJpaEntity(eventDto));

        final EventDto savedDto = eventMapper.toDto(savedEntity);

        eventRedisRepository.save(eventMapper.toRedisEntity(savedDto));

        log.debug("Event created successfully id={}", savedEntity.getId());

        return savedDto;
    }

    public EventDto get(final String id) {
        Assert.hasText(id, "Event id must not be null or blank");

        return eventRedisRepository
                .findById(id)
                .map(cached -> {
                    log.debug("Cache hit for Event id= {}", id);

                    return eventMapper.toDto(cached);
                }).orElseGet(() ->
                        eventJpaRepository.findById(id)
                                .map(entity -> {
                                    final EventDto dto = eventMapper.toDto(entity);
                                    eventRedisRepository.save(eventMapper.toRedisEntity(dto));
                                    log.debug("Event id = {} cached after loading from Postgres", id);
                                    return dto;
                                }).orElseThrow(() -> new EntityNotFoundException("Event is not found with id = %s".formatted(id))));
    }

    public PageResponse<EventDto> getAll(final Pageable pageable) {
        final Page<EventJpaEntity> page = eventJpaRepository.findAll(pageable);

        final Page<EventDto> dtoPage = page.map(eventMapper::toDto);

        return new PageResponse<>
                (dtoPage.getContent(),
                        dtoPage.getNumber(),
                        dtoPage.getSize(),
                        dtoPage.getTotalElements(),
                        dtoPage.getTotalPages(), resolveSortBy(pageable), resolveDirection(pageable));
    }

    public EventDto update(final String id, final EventDto dto) {
        log.debug("Updating  Event id = {} in Postgres and Redis", id);

        if (!eventJpaRepository.existsById(id)) {
            throw new EntityNotFoundException("The event not found");
        }

        final EventDto eventDto = new EventDto(id, dto.title(), dto.description());


        final EventJpaEntity updatedEntity = eventJpaRepository.save(eventMapper.toJpaEntity(eventDto));

        final EventDto updatedDto = eventMapper.toDto(updatedEntity);

        eventRedisRepository.save(eventMapper.toRedisEntity(updatedDto));

        return updatedDto;
    }

    public void delete(String id) {
        log.debug("Deleting  Event id = {} in Postgres and Redis", id);

        if (!eventJpaRepository.existsById(id)) {
            throw new EntityNotFoundException("The event not found");
        }

        eventRedisRepository.deleteById(id);
        eventJpaRepository.deleteById(id);
    }

    private String resolveSortBy(Pageable pageable) {
        return pageable.getSort().stream()
                .findFirst()
                .map(Sort.Order::getProperty)
                .orElse("unsorted");

    }

    private String resolveDirection(Pageable pageable) {
        return pageable.getSort().stream()
                .findFirst()
                .map(order -> order.getDirection().name())
                .orElse("ASC");

    }
}
