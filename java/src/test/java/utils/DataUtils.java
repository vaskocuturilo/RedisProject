package utils;

import com.example.api.dto.EventDto;
import com.example.api.entity.EventJpaEntity;
import com.example.api.entity.EventRedisEntity;

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

    public static EventJpaEntity getEventEntityPersisted() {
        return EventJpaEntity
                .builder()
                .id(UUID.randomUUID().toString())
                .title("Saved Title")
                .description("Saved Description")
                .build();
    }

    public static EventDto getEventDtoTransient() {
        return EventDto
                .builder()
                .id(UUID.randomUUID().toString())
                .title("Title Dto")
                .description("Description Dto")
                .build();
    }

    public static EventRedisEntity getEventRedisEntityTransient() {
        return EventRedisEntity
                .builder()
                .id(UUID.randomUUID().toString())
                .title("Title Dto")
                .description("Description Dto")
                .build();
    }
}
