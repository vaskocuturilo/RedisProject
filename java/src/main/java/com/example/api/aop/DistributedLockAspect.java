package com.example.api.aop;

import com.example.api.annotation.DistributedLock;
import com.example.api.component.LockKeyResolver;
import com.example.api.exception.LockAcquisitionException;
import com.example.api.service.DistributedLockService;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.aspectj.lang.ProceedingJoinPoint;
import org.aspectj.lang.annotation.Around;
import org.aspectj.lang.annotation.Aspect;
import org.aspectj.lang.reflect.MethodSignature;
import org.springframework.stereotype.Component;

import java.util.UUID;

@Slf4j
@Aspect
@Component
@RequiredArgsConstructor
public class DistributedLockAspect {
    private final DistributedLockService lockService;
    private final LockKeyResolver keyResolver;

    @Around("@annotation(distributedLock)")
    public Object around(ProceedingJoinPoint joinPoint, DistributedLock distributedLock) throws Throwable {

        MethodSignature signature = (MethodSignature) joinPoint.getSignature();
        String resolvedKey = "lock:" + keyResolver.resolve(
                distributedLock.key(),
                signature,
                joinPoint.getArgs()
        );

        long leaseTimeMillis = distributedLock.timeUnit()
                .toMillis(distributedLock.leaseTime());
        long timeoutMillis = distributedLock.timeUnit()
                .toMillis(distributedLock.timeout());

        String lockValue = UUID.randomUUID().toString();

        if (!tryAcquireWithRetry(resolvedKey, lockValue, leaseTimeMillis, timeoutMillis)) {
            throw new LockAcquisitionException(
                    "Could not acquire lock for key: " + resolvedKey
            );
        }

        try {
            return joinPoint.proceed();
        } finally {
            boolean released = lockService.release(resolvedKey, lockValue);
            if (!released) {
                log.warn("Lock was not released for key '{}' — may have expired or been overridden", resolvedKey);
            }
        }
    }

    private boolean tryAcquireWithRetry(
            String key, String lockValue,
            long leaseTimeMillis, long timeoutMillis) {

        long deadline = System.currentTimeMillis() + timeoutMillis;
        long retryIntervalMillis = 100L;

        while (System.currentTimeMillis() < deadline) {
            if (lockService.acquire(key, lockValue, leaseTimeMillis)) {
                return true;
            }
            try {
                Thread.sleep(retryIntervalMillis);
            } catch (InterruptedException _) {
                Thread.currentThread().interrupt();
                return false;
            }
        }

        return false;
    }
}
