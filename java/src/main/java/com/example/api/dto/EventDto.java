package com.example.api.dto;

import java.io.Serializable;

public record EventDto(String id, String title, String description) implements Serializable {
}
