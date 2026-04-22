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

    public static EventJpaEntity getEvent1EntityPersisted() {
        return EventJpaEntity
                .builder()
                .id(UUID.randomUUID().toString())
                .title("Title 1")
                .description("Description 1")
                .build();
    }

    public static EventJpaEntity getEvent2EntityPersisted() {
        return EventJpaEntity
                .builder()
                .id(UUID.randomUUID().toString())
                .title("Title 2")
                .description("Description 2")
                .build();
    }

    public static EventJpaEntity getEvent3EntityPersisted() {
        return EventJpaEntity
                .builder()
                .id(UUID.randomUUID().toString())
                .title("Title 3")
                .description("Description 3")
                .build();
    }

    public static EventDto getEvent1DtoTransient() {
        return EventDto
                .builder()
                .id(UUID.randomUUID().toString())
                .title("Title Dto 1")
                .description("Description Dto 1")
                .build();
    }

    public static EventDto getEvent2DtoTransient() {
        return EventDto
                .builder()
                .id(UUID.randomUUID().toString())
                .title("Title Dto 2")
                .description("Description Dto 2")
                .build();
    }

    public static EventDto getEvent3DtoTransient() {
        return EventDto
                .builder()
                .id(UUID.randomUUID().toString())
                .title("Title Dto 3")
                .description("Description Dto 3")
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
