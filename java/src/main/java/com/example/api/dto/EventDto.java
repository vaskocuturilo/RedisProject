package com.example.api.dto;

import lombok.Builder;

import java.io.Serializable;

@Builder
public record EventDto(String id, String title, String description) implements Serializable {
}
