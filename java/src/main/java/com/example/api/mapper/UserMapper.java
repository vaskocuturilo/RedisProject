package com.example.api.mapper;

import com.example.api.dto.UserDto;
import com.example.api.entity.UserJpaEntity;
import com.example.api.entity.UserRedisEntity;
import org.mapstruct.Mapper;
import org.mapstruct.MappingConstants;

@Mapper(componentModel = MappingConstants.ComponentModel.SPRING)
public interface UserMapper {

    //Redis
    UserDto toDto(UserRedisEntity userRedisEntity);

    UserRedisEntity toRedisEntity(UserDto userDto);


    //JPA
    UserDto toDto(UserJpaEntity userJpaEntity);

    UserJpaEntity toJpaEntity(UserDto eventDto);
}
