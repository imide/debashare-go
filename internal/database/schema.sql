CREATE TABLE rooms (
                       id VARCHAR(36) PRIMARY KEY,
                       created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE files (
                       id VARCHAR(36) PRIMARY KEY,
                       room_id VARCHAR(36) NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
                       name VARCHAR(255) NOT NULL,
                       size BIGINT NOT NULL,
                       uploaded_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
