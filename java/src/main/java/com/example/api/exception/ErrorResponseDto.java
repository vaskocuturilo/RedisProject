package com.example.api.exception;

import lombok.Builder;

import java.time.LocalDateTime;
@Builder
public record ErrorResponseDto
        (
                String error,
                int errorCode,
                String message,
                LocalDateTime timestamp
        ) {
}
