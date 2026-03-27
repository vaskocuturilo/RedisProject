package com.example.api.repository;

import com.example.api.entity.UserJpaEntity;
import org.springframework.data.repository.CrudRepository;
import org.springframework.stereotype.Repository;

@Repository
public interface UserJpaRepository extends CrudRepository<UserJpaEntity, String> {
}
