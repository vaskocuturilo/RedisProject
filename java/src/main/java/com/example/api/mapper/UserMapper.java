package com.example.api.mapper;

import com.example.api.dto.UserDto;
import com.example.api.entity.UserRedisEntity;
import org.mapstruct.Mapper;
import org.mapstruct.factory.Mappers;

@Mapper(componentModel = "spring")
public interface UserMapper {

    UserMapper INSTANCE = Mappers.getMapper( UserMapper.class );

    UserDto toDto(UserRedisEntity userRedisEntity);

    UserRedisEntity toEntity(UserDto userDto);
}
