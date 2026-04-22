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
public class RateLimiterService {
    private final StringRedisTemplate redisTemplate;

    private static final RedisScript<Long> RATE_LIMIT_SCRIPT = RedisScript.of(
            """
                    local count = redis.call('INCR', KEYS[1])
                    if count == 1 then
                        redis.call('EXPIRE', KEYS[1], ARGV[1])
                    else
                        -- Ensure TTL is always set even if first call missed it
                        if redis.call('TTL', KEYS[1]) == -1 then
                            redis.call('EXPIRE', KEYS[1], ARGV[1])
                        end
                    end
                    return count
                    """,
            Long.class
    );


    public boolean isAllowed(final String key, int limit, int durationSeconds) {
        String redisKey = "rate_limit:" + key;

        try {
            long count = Objects.requireNonNullElse(
                    redisTemplate.execute(
                            RATE_LIMIT_SCRIPT,
                            Collections.singletonList(redisKey),
                            String.valueOf(durationSeconds)
                    ),
                    0L
            );
            return count <= limit;

        } catch (RedisException ex) {
            log.warn("Redis unavailable during rate limiting, failing open: {}", ex.getMessage());
            return true;
        }
    }
}
