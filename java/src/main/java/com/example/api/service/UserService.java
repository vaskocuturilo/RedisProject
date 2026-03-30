package com.example.api.service;

import com.example.api.dto.UserDto;
import com.example.api.entity.UserJpaEntity;
import com.example.api.entity.UserRedisEntity;
import com.example.api.mapper.UserMapper;
import com.example.api.repository.UserJpaRepository;
import com.example.api.repository.UserRedisRepository;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;

import java.util.UUID;

@Service
@Slf4j
public class UserService {

    private final UserRedisRepository userRedisRepository;
    private final UserMapper userMapper;
    private final UserJpaRepository userJpaRepository;

    public UserService(UserRedisRepository userRedisRepository, UserMapper userMapper, UserJpaRepository userJpaRepository) {
        this.userRedisRepository = userRedisRepository;
        this.userMapper = userMapper;
        this.userJpaRepository = userJpaRepository;
    }

    public UserDto create(final UserDto dto) {
        final UserDto userDto = new UserDto(UUID.randomUUID().toString(), dto.name(), dto.age(), dto.events());
        final UserJpaEntity savedUser = userJpaRepository.save(userMapper.toJpaEntity(userDto));

        userRedisRepository.save(userMapper.toRedisEntity(userDto));

        log.info("IN create - saved user: {}", savedUser);

        return userMapper.toDto(savedUser);
    }

    public UserDto get(final String id) {
        final UserDto existUser = userRedisRepository.findById(id).map(userMapper::toDto).orElse(null);

        log.info("IN get - get user: {}", existUser);

        return existUser;
    }

    public UserDto update(final String id, final UserDto dto) {
        final UserRedisEntity entity = userMapper.toRedisEntity(dto);

        entity.setId(id);

        log.info("IN update - updated user: {}", entity);

        return userMapper.toDto(userRedisRepository.save(entity));
    }

    public void delete(String id) {
        userRedisRepository.deleteById(id);
    }
}
