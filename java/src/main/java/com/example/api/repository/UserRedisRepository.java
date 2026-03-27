package com.example.api.repository;

import com.example.api.entity.UserRedisEntity;
import org.springframework.data.repository.CrudRepository;
import org.springframework.stereotype.Repository;

@Repository
public interface UserRedisRepository extends CrudRepository<UserRedisEntity, String> {
}
