package com.example.api.mapper;

import com.example.api.dto.EventDto;
import com.example.api.entity.EventJpaEntity;
import com.example.api.entity.EventRedisEntity;
import org.mapstruct.Mapper;
import org.mapstruct.MappingConstants;

@Mapper(componentModel = MappingConstants.ComponentModel.SPRING)
public interface EventMapper {


    EventDto toDto(EventRedisEntity eventRedisEntity);

    EventRedisEntity toRedisEntity(EventDto eventDto);


    //JPA
    EventDto toDto(EventJpaEntity eventJpaEntity);

    EventJpaEntity toJpaEntity(EventDto eventDto);
}
