package com.example.api.service;

import lombok.RequiredArgsConstructor;
import org.springframework.data.redis.core.StringRedisTemplate;
import org.springframework.stereotype.Service;

import java.util.Objects;
import java.util.concurrent.TimeUnit;

@Service
@RequiredArgsConstructor
public class RateLimiterService {
    private final StringRedisTemplate redisTemplate;


    public boolean isAllowed(final String key, int limit, int durationSeconds) {
        final String redisKey = "rate_limit:" + key;

        final Long count = redisTemplate.opsForValue().increment(key);

        if (Objects.isNull(count)) {
            return true;
        }

        if (count == 1) {
            redisTemplate.expire(redisKey, durationSeconds, TimeUnit.SECONDS);
        }
        return count <= limit;
    }
}
