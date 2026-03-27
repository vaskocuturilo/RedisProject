package com.example.api.repository;

import com.example.api.entity.EventJpaEntity;
import org.springframework.data.repository.CrudRepository;
import org.springframework.stereotype.Repository;

@Repository
public interface EventJpaRepository extends CrudRepository<EventJpaEntity, String> {
}
