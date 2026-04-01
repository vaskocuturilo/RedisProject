package com.example.api.rest;

import com.example.api.annotation.DistributedLock;
import com.example.api.annotation.RateLimiter;
import com.example.api.dto.EventDto;
import com.example.api.service.EventService;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping("/api/v1/events")
public class EventControllerV1 {

    private final EventService eventService;

    public EventControllerV1(EventService eventService) {
        this.eventService = eventService;
    }


    @PostMapping
    @RateLimiter(limit = 3, duration = 60, key = "create-event")
    @DistributedLock(leaseTime = 10, key = "create-event")
    public EventDto create(@RequestBody EventDto eventDto) {
        return eventService.create(eventDto);
    }

    @PutMapping("/{id}")
    @RateLimiter(limit = 3, duration = 60, key = "put-event")
    @DistributedLock(leaseTime = 10, key = "put-event")
    public EventDto update(@PathVariable String id, @RequestBody EventDto eventDto) {
        return eventService.update(id, eventDto);
    }

    @GetMapping
    @RateLimiter(limit = 3, duration = 60, key = "get-all-events")
    @DistributedLock(leaseTime = 10, key = "get-all-events")
    public List<EventDto> getAll() {
        return eventService.getAll();
    }

    @GetMapping("/{id}")
    @RateLimiter(limit = 3, duration = 60, key = "get-event")
    @DistributedLock(leaseTime = 10, key = "get-event")
    public EventDto get(@PathVariable String id) {
        return eventService.get(id);
    }

    @DeleteMapping("/{id}")
    @RateLimiter(limit = 3, duration = 60, key = "delete-event")
    @DistributedLock(leaseTime = 10, key = "delete-event")
    public void delete(@PathVariable String id) {
        eventService.delete(id);
    }
}
