package utils;

import com.example.api.entity.EventJpaEntity;

import java.util.UUID;


public class DataUtils {

    public static EventJpaEntity getEventEntityTransient() {
        return EventJpaEntity
                .builder()
                .id(UUID.randomUUID().toString())
                .title("Title")
                .description("Description")
                .build();
    }
}
