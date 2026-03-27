package com.example.api.rest;

import com.example.api.dto.EventDto;
import com.example.api.service.EventService;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/api/v1/events")
public class EventControllerV1 {

    private final EventService eventService;

    public EventControllerV1(EventService eventService) {
        this.eventService = eventService;
    }


    @PostMapping
    public EventDto create(@RequestBody EventDto eventDto) {
        return eventService.create(eventDto);
    }

    @PutMapping("/{id}")
    public EventDto update(@PathVariable String id, @RequestBody EventDto eventDto) {
        return eventService.update(id, eventDto);
    }

    @GetMapping("/{id}")
    public EventDto get(@PathVariable String id) {
        return eventService.get(id);
    }

    @DeleteMapping("/{id}")
    public void delete(@PathVariable String id) {
        eventService.delete(id);
    }
}
