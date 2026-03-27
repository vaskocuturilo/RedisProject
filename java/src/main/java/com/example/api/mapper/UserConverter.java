package com.example.api.mapper;

import com.example.api.dto.UserDto;
import com.example.api.entity.UserJpaEntity;
import com.example.api.entity.UserRedisEntity;

import java.util.UUID;

public class UserConverter {

    public static UserDto toDto(final UserRedisEntity entity) {
        return new UserDto(UUID.randomUUID().toString(), entity.getName(), entity.getAge(), entity.getEvents());
    }

    public static UserRedisEntity toEntity(final UserDto userDto) {
        return new UserRedisEntity(UUID.randomUUID().toString(), userDto.name(), userDto.age(), userDto.events());
    }

    public static UserDto toDto(final UserJpaEntity entity) {
        return new UserDto(UUID.randomUUID().toString(), entity.getName(), entity.getAge(), entity.getEvents());
    }

    public static UserJpaEntity toJpaEntity(final UserDto userDto) {
        return new UserJpaEntity(UUID.randomUUID().toString(), userDto.name(), userDto.age(), userDto.events());
    }
}
