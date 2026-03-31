package com.example.api.aop;

import com.example.api.annotation.RateLimiter;
import com.example.api.exception.RateLimitExceededException;
import com.example.api.service.RateLimiterService;
import jakarta.servlet.http.HttpServletRequest;
import lombok.RequiredArgsConstructor;
import org.aspectj.lang.ProceedingJoinPoint;
import org.aspectj.lang.annotation.Around;
import org.aspectj.lang.annotation.Aspect;
import org.springframework.stereotype.Component;

@Aspect
@Component
@RequiredArgsConstructor
public class RateLimitAspect {

    private final RateLimiterService rateLimiterService;
    private final HttpServletRequest request;

    @Around("@annotation(rateLimit)")
    public Object around(ProceedingJoinPoint joinPoint, RateLimiter rateLimit) throws Throwable {

        String ip = request.getRemoteAddr();
        String methodName = joinPoint.getSignature().toShortString();
        String customKey = rateLimit.key().isEmpty() ? methodName : rateLimit.key();
        String finalKey = ip + ":" + customKey;

        boolean allowed = rateLimiterService.isAllowed(
                finalKey,
                rateLimit.limit(),
                rateLimit.duration()
        );

        if (!allowed) {
            throw new RateLimitExceededException(
                    "Rate limit exceeded. Max " + rateLimit.limit() +
                            " requests per " + rateLimit.duration() + "s."
            );
        }

        return joinPoint.proceed();
    }
}
