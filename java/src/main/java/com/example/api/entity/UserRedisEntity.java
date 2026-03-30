package com.example.api.entity;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;
import org.springframework.data.annotation.Id;
import org.springframework.data.redis.core.RedisHash;

import java.io.Serializable;
import java.util.HashSet;
import java.util.Set;

@RedisHash(value = "User", timeToLive = 3600L)
@Data
@NoArgsConstructor
@AllArgsConstructor
@Builder
public class UserRedisEntity implements Serializable {

    @Id
    private String id;
    private String name;
    private int age;

    @Builder.Default
    private Set<String> events = new HashSet<>();
}
