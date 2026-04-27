package com.example.api.rest;

import com.example.api.annotation.DistributedLock;
import com.example.api.annotation.RateLimiter;
import com.example.api.dto.EventDto;
import com.example.api.dto.PageResponse;
import com.example.api.service.EventService;
import lombok.extern.slf4j.Slf4j;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Pageable;
import org.springframework.data.domain.Sort;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.Set;

@RestController
@RequestMapping("/api/v1/events")
@Slf4j
public class EventControllerV1 {

    private final EventService eventService;

    private static final Set<String> ALLOWED_SORT_FIELDS = Set.of("id", "title", "description");

    public EventControllerV1(EventService eventService) {
        this.eventService = eventService;
    }


    @PostMapping
    @RateLimiter(limit = 3, duration = 60, key = "create-event")
    @DistributedLock(leaseTime = 10, key = "create-event")
    public ResponseEntity<EventDto> create(@RequestBody EventDto eventDto) {
        log.debug("Create, Received DTO: {}", eventDto);
        return ResponseEntity.status(HttpStatus.CREATED).body(eventService.create(eventDto));
    }

    @PutMapping("/{id}")
    @RateLimiter(limit = 3, duration = 60, key = "put-event")
    @DistributedLock(leaseTime = 10, key = "put-event")
    public ResponseEntity<EventDto> update(@PathVariable String id, @RequestBody EventDto eventDto) {
        log.debug("Update, Received DTO: {}", eventDto);
        return ResponseEntity.status(HttpStatus.OK).body(eventService.update(id, eventDto));
    }

    @GetMapping
    public ResponseEntity<PageResponse<EventDto>> getAll(
            @RequestParam(defaultValue = "0") int page,
            @RequestParam(defaultValue = "10") int size,
            @RequestParam(defaultValue = "title") String sortBy,
            @RequestParam(defaultValue = "asc") String direction
    ) {
        final int safePage = Math.max(page, 0);

        final int safeSize = Math.clamp(size, 1, 100);

        String safeSortBy = ALLOWED_SORT_FIELDS.contains(sortBy) ? sortBy : "title";

        final Sort sort = direction.equalsIgnoreCase("desc")
                ? Sort.by(safeSortBy).descending()
                : Sort.by(safeSortBy).ascending();

        final Pageable pageable = PageRequest.of(safePage, safeSize, sort);

        log.debug("Get All events page={}, size={}, sortBy={}, direction={}",
                page, size, sortBy, direction);

        return ResponseEntity.status(HttpStatus.OK).body(eventService.getAll(pageable));
    }

    @GetMapping("/{id}")
    @RateLimiter(limit = 3, duration = 60, key = "get-event")
    public ResponseEntity<EventDto> get(@PathVariable String id) {
        log.debug("Get by id: {}", id);
        return ResponseEntity.status(HttpStatus.OK).body(eventService.get(id));
    }

    @DeleteMapping("/{id}")
    @RateLimiter(limit = 3, duration = 60, key = "delete-event")
    @DistributedLock(leaseTime = 10, key = "delete-event")
    public ResponseEntity<Void> delete(@PathVariable String id) {
        log.debug("Delete by id: {}", id);
        eventService.delete(id);
        return ResponseEntity.noContent().build();
    }
}
