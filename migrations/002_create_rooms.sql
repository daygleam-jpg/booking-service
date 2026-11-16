CREATE TABLE rooms(
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    hotel_id BIGINT NOT NULL REFERENCES hotels(id) ON DELETE CASCADE,
    number VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'available'
       CHECK(status IN('available', 'maintenance')),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (hotel_id, number)
    );
