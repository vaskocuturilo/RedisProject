package com.example.api.rest;

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
    public UserDto create(@RequestBody UserDto userDto) {
        return userService.create(userDto);
    }

    @PutMapping("/{id}")
    public UserDto update(@PathVariable String id, @RequestBody UserDto userDto) {
        return userService.update(id, userDto);
    }

    @GetMapping("/{id}")
    public UserDto get(@PathVariable String id) {
        return userService.get(id);
    }

    @DeleteMapping("/{id}")
    public void delete(@PathVariable String id) {
        userService.delete(id);
    }
}
