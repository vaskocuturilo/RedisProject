package com.example.api.repository;

import com.example.api.entity.EventRedisEntity;
import org.springframework.data.repository.CrudRepository;
import org.springframework.stereotype.Repository;

@Repository
public interface EventRedisRepository extends CrudRepository<EventRedisEntity, String> {
}
