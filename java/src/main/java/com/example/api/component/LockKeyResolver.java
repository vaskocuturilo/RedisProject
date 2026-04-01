package com.example.api.component;

import org.aspectj.lang.reflect.MethodSignature;
import org.springframework.expression.ExpressionParser;
import org.springframework.expression.spel.standard.SpelExpressionParser;
import org.springframework.expression.spel.support.StandardEvaluationContext;
import org.springframework.stereotype.Component;

import java.util.Objects;

@Component
public class LockKeyResolver {

    private final ExpressionParser parser = new SpelExpressionParser();

    public String resolve(String keyExpression, MethodSignature signature, Object[] args) {
        if (!keyExpression.contains("#")) {
            return keyExpression;
        }

        StandardEvaluationContext context = new StandardEvaluationContext();
        String[] paramNames = signature.getParameterNames();

        for (int i = 0; i < paramNames.length; i++) {
            context.setVariable(paramNames[i], args[i]);
        }

        return Objects.requireNonNullElse(
                parser.parseExpression(keyExpression).getValue(context, String.class),
                keyExpression
        );
    }
}
