package com.example.api.mapper;

import com.example.api.dto.EventDto;
import com.example.api.entity.EventJpaEntity;
import com.example.api.entity.EventRedisEntity;

import java.util.UUID;

public class EventConverter {

    public static EventDto toDto(final EventRedisEntity entity) {
        return new EventDto(UUID.randomUUID().toString(), entity.getTitle(), entity.getDescription());
    }

    public static EventRedisEntity toEntity(final EventDto eventDto) {
        return new EventRedisEntity(UUID.randomUUID().toString(), eventDto.title(), eventDto.description());
    }

    //JPA
    public static EventDto toDto(final EventJpaEntity entity) {
        return new EventDto(UUID.randomUUID().toString(), entity.getTitle(), entity.getDescription());
    }

    public static EventJpaEntity toJpaEntity(final EventDto eventDto) {
        return new EventJpaEntity(UUID.randomUUID().toString(), eventDto.title(), eventDto.description());
    }
}
