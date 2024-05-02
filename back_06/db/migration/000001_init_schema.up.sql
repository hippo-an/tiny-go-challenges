CREATE TABLE users
(
    id           SERIAL PRIMARY KEY,
    first_name   VARCHAR(255)        NOT NULL,
    last_name    VARCHAR(255)        NOT NULL,
    email        VARCHAR(255) UNIQUE NOT NULL,
    password     VARCHAR(255)        NOT NULL,
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL ,
    updated_at   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL ,
    access_level INTEGER             NOT NULL
);

CREATE TABLE rooms
(
    id         SERIAL PRIMARY KEY,
    room_name  VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL ,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE restrictions
(
    id               SERIAL PRIMARY KEY,
    restriction_name VARCHAR(255) NOT NULL ,
    created_at       TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL ,
    updated_at       TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);


CREATE TABLE reservations
(
    id         SERIAL PRIMARY KEY,
    first_name VARCHAR(255)        NOT NULL,
    last_name  VARCHAR(255)        NOT NULL,
    email      VARCHAR(255) UNIQUE NOT NULL,
    phone      VARCHAR(255)        NOT NULL,
    start_date DATE                NOT NULL,
    end_date   DATE                NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    room_id    INTEGER NOT NULL,
    CONSTRAINT fk_reservation_room FOREIGN KEY (room_id)
        REFERENCES rooms (id) ON DELETE CASCADE ON UPDATE CASCADE
);



CREATE TABLE room_restrictions
(
    id             SERIAL PRIMARY KEY,
    start_date     DATE NOT NULL,
    end_date       DATE NOT NULL,
    room_id        INTEGER NOT NULL,
    reservation_id INTEGER NULL ,
    restriction_id INTEGER NOT NULL,
    created_at     TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at     TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT fk_restriction_room FOREIGN KEY (room_id)
        REFERENCES rooms (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT fk_restriction_reservation FOREIGN KEY (reservation_id)
        REFERENCES reservations (id)
        ON DELETE SET NULL
        ON UPDATE CASCADE,
    CONSTRAINT fk_restriction_restriction FOREIGN KEY (restriction_id)
        REFERENCES restrictions (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

