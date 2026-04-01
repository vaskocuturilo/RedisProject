package com.example.api.service;

import io.lettuce.core.RedisException;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.data.redis.core.StringRedisTemplate;
import org.springframework.data.redis.core.script.RedisScript;
import org.springframework.stereotype.Service;

import java.util.Collections;
import java.util.Objects;

@Slf4j
@Service
@RequiredArgsConstructor
public class DistributedLockService {
    private final StringRedisTemplate redisTemplate;

    private static final RedisScript<Long> ACQUIRE_SCRIPT = RedisScript.of(
            """
                    if redis.call('SET', KEYS[1], ARGV[1], 'NX', 'PX', ARGV[2]) then
                        return 1
                    else
                        return 0
                    end
                    """,
            Long.class
    );

    private static final RedisScript<Long> RELEASE_SCRIPT = RedisScript.of(
            """
                    if redis.call('GET', KEYS[1]) == ARGV[1] then
                        return redis.call('DEL', KEYS[1])
                    else
                        return 0
                    end
                    """,
            Long.class
    );

    public boolean acquire(String key, String lockValue, long leaseTimeMillis) {
        try {
            long result = Objects.requireNonNullElse(
                    redisTemplate.execute(
                            ACQUIRE_SCRIPT,
                            Collections.singletonList(key),
                            lockValue,
                            String.valueOf(leaseTimeMillis)
                    ),
                    0L
            );
            return result == 1L;

        } catch (RedisException ex) {
            log.warn("Redis unavailable during lock acquire for key '{}': {}", key, ex.getMessage());
            return false;
        }
    }

    public boolean release(String key, String lockValue) {
        try {
            long result = Objects.requireNonNullElse(
                    redisTemplate.execute(
                            RELEASE_SCRIPT,
                            Collections.singletonList(key),
                            lockValue
                    ),
                    0L
            );
            return result == 1L;

        } catch (RedisException ex) {
            log.warn("Redis unavailable during lock release for key '{}': {}", key, ex.getMessage());
            return false;
        }
    }
}
