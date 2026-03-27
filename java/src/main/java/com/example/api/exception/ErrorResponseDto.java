package com.example.api.exception;

import java.time.LocalDateTime;

public record ErrorResponseDto
        (
                String error,
                int errorCode,
                String message,
                LocalDateTime timestamp
        ) {
}
