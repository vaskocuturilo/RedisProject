package com.example.api.exception;

import jakarta.persistence.EntityNotFoundException;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.ControllerAdvice;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.bind.annotation.ResponseStatus;
import org.springframework.web.bind.annotation.RestControllerAdvice;

import java.time.LocalDateTime;
import java.util.Map;
import java.util.MissingFormatArgumentException;

@RestControllerAdvice
@ControllerAdvice
public class GlobalExceptionHandler {

    private static final String KEY = "message";

    @ExceptionHandler(EntityNotFoundException.class)
    public ResponseEntity<ErrorResponseDto> handleRegionNotFound(final EntityNotFoundException exception) {
        final ErrorResponseDto regionNotFoundException = new ErrorResponseDto(
                "Not Found",
                HttpStatus.NOT_FOUND.value(),
                exception.getMessage(),
                LocalDateTime.now());

        return ResponseEntity
                .status(HttpStatus.NOT_FOUND)
                .contentType(MediaType.APPLICATION_JSON)
                .body(regionNotFoundException);
    }

    @ExceptionHandler(IllegalArgumentException.class)
    public ResponseEntity<Map<String, String>> handleIllegalArgument(IllegalArgumentException ex) {
        return ResponseEntity
                .status(HttpStatus.BAD_REQUEST)
                .body(Map.of(KEY, ex.getMessage()));
    }

    @ExceptionHandler(MissingFormatArgumentException.class)
    public ResponseEntity<Map<String, String>> handleMissingFormat(MissingFormatArgumentException ex) {
        return ResponseEntity
                .status(HttpStatus.BAD_REQUEST)
                .body(Map.of(KEY, "Missing format argument: " + ex.getMessage()));
    }

    @ExceptionHandler(RateLimitExceededException.class)
    @ResponseStatus(HttpStatus.TOO_MANY_REQUESTS)  // 429
    public Map<String, String> handleRateLimit(RateLimitExceededException ex) {
        return Map.of(
                "error", "Too Many Requests",
                KEY, ex.getMessage()
        );
    }

    @ExceptionHandler(Exception.class)
    public ResponseEntity<Map<String, String>> handleGeneric(Exception ex) {
        return ResponseEntity
                .status(HttpStatus.INTERNAL_SERVER_ERROR)
                .body(Map.of(KEY, "Unexpected error: " + ex.getMessage()));
    }
}
