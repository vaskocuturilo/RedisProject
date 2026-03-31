package com.example.api.rest;

import com.example.api.annotation.RateLimiter;
import com.example.api.dto.UserDto;
import com.example.api.service.UserService;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/api/v1/users")
public class UserControllerV1 {

    private final UserService userService;

    public UserControllerV1(UserService userService) {
        this.userService = userService;
    }

    @PostMapping
    @RateLimiter(limit = 3, duration = 60, key = "create-user")
    public UserDto create(@RequestBody UserDto userDto) {
        return userService.create(userDto);
    }

    @PutMapping("/{id}")
    @RateLimiter(limit = 3, duration = 60, key = "put-user")
    public UserDto update(@PathVariable String id, @RequestBody UserDto userDto) {
        return userService.update(id, userDto);
    }

    @GetMapping("/{id}")
    @RateLimiter(limit = 3, duration = 60, key = "get-user")
    public UserDto get(@PathVariable String id) {
        return userService.get(id);
    }

    @DeleteMapping("/{id}")
    @RateLimiter(limit = 3, duration = 60, key = "delete-user")
    public void delete(@PathVariable String id) {
        userService.delete(id);
    }
}
