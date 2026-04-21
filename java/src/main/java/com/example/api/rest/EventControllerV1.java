package com.example.api.rest;

import com.example.api.annotation.DistributedLock;
import com.example.api.annotation.RateLimiter;
import com.example.api.dto.EventDto;
import com.example.api.exception.ErrorResponseDto;
import com.example.api.service.EventService;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping("/api/v1/events")
@Slf4j
public class EventControllerV1 {

    private final EventService eventService;

    public EventControllerV1(EventService eventService) {
        this.eventService = eventService;
    }


    @PostMapping
    @RateLimiter(limit = 3, duration = 60, key = "create-event")
    @DistributedLock(leaseTime = 10, key = "create-event")
    public ResponseEntity<EventDto> create(@RequestBody EventDto eventDto) {
        log.info("Received DTO: {}", eventDto);
        return ResponseEntity.status(HttpStatus.CREATED).body(eventService.create(eventDto));
    }

    @PutMapping("/{id}")
    @RateLimiter(limit = 3, duration = 60, key = "put-event")
    @DistributedLock(leaseTime = 10, key = "put-event")
    public ResponseEntity<EventDto> update(@PathVariable String id, @RequestBody EventDto eventDto) {
        return ResponseEntity.status(HttpStatus.OK).body(eventService.update(id, eventDto));
    }

    @GetMapping
    @RateLimiter(limit = 3, duration = 60, key = "get-all-events")
    @DistributedLock(leaseTime = 10, key = "get-all-events")
    public ResponseEntity<List<EventDto>> getAll() {
        return ResponseEntity.status(HttpStatus.OK).body(eventService.getAll());
    }

    @GetMapping("/{id}")
    @RateLimiter(limit = 3, duration = 60, key = "get-event")
    @DistributedLock(leaseTime = 10, key = "get-event")
    public ResponseEntity<EventDto> get(@PathVariable String id) {
        return ResponseEntity.status(HttpStatus.OK).body(eventService.get(id));
    }

    @DeleteMapping("/{id}")
    @RateLimiter(limit = 3, duration = 60, key = "delete-event")
    @DistributedLock(leaseTime = 10, key = "delete-event")
    public ResponseEntity<ErrorResponseDto> delete(@PathVariable String id) {
        eventService.delete(id);
        return ResponseEntity.noContent().build();
    }
}
