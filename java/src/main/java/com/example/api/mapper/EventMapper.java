package com.example.api.mapper;

import com.example.api.dto.EventDto;
import com.example.api.entity.EventJpaEntity;
import com.example.api.entity.EventRedisEntity;
import org.mapstruct.Mapper;

@Mapper(componentModel = "spring")
public interface EventMapper {

    //Redis
    EventDto toDto(EventRedisEntity eventRedisEntity);

    EventRedisEntity toRedisEntity(EventDto eventDto);


    //JPA
    EventDto toDto(EventJpaEntity eventJpaEntity);

    EventJpaEntity toJpaEntity(EventDto eventDto);
}

