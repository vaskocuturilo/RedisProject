package com.example.api.entity;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;
import org.springframework.data.annotation.Id;
import org.springframework.data.redis.core.RedisHash;

import java.io.Serializable;

@RedisHash(value = "Event", timeToLive = 3600L)
@Data
@AllArgsConstructor
@NoArgsConstructor
@Builder
public class EventRedisEntity implements Serializable {

    @Id
    private String id;
    private String title;
    private String description;
}
