CREATE DATABASE IF NOT EXISTS ridehailing
DEFAULT CHARACTER SET utf8mb4
DEFAULT COLLATE utf8mb4_unicode_ci;

USE ridehailing;

CREATE TABLE IF NOT EXISTS users (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    phone VARCHAR(20) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    nickname VARCHAR(50) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'passenger',
    default_role VARCHAR(20) NOT NULL DEFAULT 'passenger',
    real_name VARCHAR(50) DEFAULT NULL,
    id_card VARCHAR(32) DEFAULT NULL,
    real_name_verified TINYINT(1) NOT NULL DEFAULT 0,
    avatar VARCHAR(255) DEFAULT NULL,
    email VARCHAR(100) NOT NULL,
    gender VARCHAR(20) DEFAULT NULL,
    birthday DATETIME NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at DATETIME(3) NULL,
    updated_at DATETIME(3) NULL,
    UNIQUE KEY uk_users_phone (phone),
    UNIQUE KEY uk_users_email (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS trips (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    driver_id BIGINT UNSIGNED NOT NULL,
    vehicle_type VARCHAR(20) NOT NULL DEFAULT 'car',
    start_city VARCHAR(50) NOT NULL,
    end_city VARCHAR(50) NOT NULL,
    departure_time DATETIME NOT NULL,
    arrival_time DATETIME NOT NULL,
    seat_total INT NOT NULL,
    seat_available INT NOT NULL,
    price_cent INT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'published',
    created_at DATETIME(3) NULL,
    updated_at DATETIME(3) NULL,
    KEY idx_driver_status (driver_id, status),
    KEY idx_trip_search (start_city, end_city, departure_time, status),
    CONSTRAINT fk_trips_driver_id
        FOREIGN KEY (driver_id) REFERENCES users(id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS trip_stops (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    trip_id BIGINT UNSIGNED NOT NULL,
    stop_order INT NOT NULL,
    stop_name VARCHAR(50) NOT NULL,
    plan_arrival_time DATETIME NULL,
    plan_departure_time DATETIME NULL,
    created_at DATETIME(3) NULL,
    updated_at DATETIME(3) NULL,
    UNIQUE KEY idx_trip_stop_order (trip_id, stop_order),
    CONSTRAINT fk_trip_stops_trip_id
        FOREIGN KEY (trip_id) REFERENCES trips(id)
        ON UPDATE CASCADE
        ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS orders (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    order_no VARCHAR(32) NOT NULL,
    user_id BIGINT UNSIGNED NOT NULL,
    trip_id BIGINT UNSIGNED NOT NULL,
    ticket_count INT NOT NULL,
    seat_type VARCHAR(30) NOT NULL DEFAULT 'standard',
    amount INT NOT NULL,
    pay_status VARCHAR(20) NOT NULL DEFAULT 'paid',
    order_status VARCHAR(30) NOT NULL DEFAULT 'pending_verification',
    refund_status VARCHAR(20) NOT NULL DEFAULT 'none',
    created_at DATETIME(3) NULL,
    updated_at DATETIME(3) NULL,
    UNIQUE KEY uk_orders_order_no (order_no),
    KEY idx_orders_user_id (user_id),
    KEY idx_orders_trip_id (trip_id),
    KEY idx_orders_order_status (order_status),
    KEY idx_orders_refund_status (refund_status),
    CONSTRAINT fk_orders_user_id
        FOREIGN KEY (user_id) REFERENCES users(id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT,
    CONSTRAINT fk_orders_trip_id
        FOREIGN KEY (trip_id) REFERENCES trips(id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

SET FOREIGN_KEY_CHECKS = 0;
DELETE FROM orders;
DELETE FROM trip_stops;
DELETE FROM trips;
DELETE FROM users;
ALTER TABLE orders AUTO_INCREMENT = 1;
ALTER TABLE trip_stops AUTO_INCREMENT = 1;
ALTER TABLE trips AUTO_INCREMENT = 1;
ALTER TABLE users AUTO_INCREMENT = 1;
SET FOREIGN_KEY_CHECKS = 1;

INSERT INTO users (
    id, phone, password_hash, nickname, role, default_role,
    real_name, id_card, real_name_verified, avatar, email,
    gender, birthday, status, created_at, updated_at
) VALUES
    (
        1, '13800010001', '$2a$10$FztOLW3BKiP48I7Hoc/idOx/69Gxm.7JuCDZzKQ90c/51uc6/rSDi',
        'driver_li', 'driver', 'driver',
        'Li Ming', '330101199001011111', 1, NULL, 'driver1@example.com',
        'male', '1990-01-01 00:00:00', 'active', '2026-05-01 09:00:00.000', '2026-05-01 09:00:00.000'
    ),
    (
        2, '13800010002', '$2a$10$FztOLW3BKiP48I7Hoc/idOx/69Gxm.7JuCDZzKQ90c/51uc6/rSDi',
        'driver_chen', 'driver', 'driver',
        'Chen Yu', '330102199202022222', 1, NULL, 'driver2@example.com',
        'female', '1992-02-02 00:00:00', 'active', '2026-05-01 09:10:00.000', '2026-05-01 09:10:00.000'
    ),
    (
        3, '13800020001', '$2a$10$FztOLW3BKiP48I7Hoc/idOx/69Gxm.7JuCDZzKQ90c/51uc6/rSDi',
        'passenger_wang', 'passenger', 'passenger',
        'Wang Lei', '330103199503033333', 1, NULL, 'passenger1@example.com',
        'male', '1995-03-03 00:00:00', 'active', '2026-05-01 09:20:00.000', '2026-05-01 09:20:00.000'
    ),
    (
        4, '13800020002', '$2a$10$FztOLW3BKiP48I7Hoc/idOx/69Gxm.7JuCDZzKQ90c/51uc6/rSDi',
        'passenger_zhao', 'passenger', 'passenger',
        'Zhao Lin', '330104199604044444', 0, NULL, 'passenger2@example.com',
        'female', '1996-04-04 00:00:00', 'active', '2026-05-01 09:30:00.000', '2026-05-01 09:30:00.000'
    ),
    (
        5, '13800090001', '$2a$10$FztOLW3BKiP48I7Hoc/idOx/69Gxm.7JuCDZzKQ90c/51uc6/rSDi',
        'admin_root', 'admin', 'admin',
        'Admin User', '330105198805055555', 1, NULL, 'admin@example.com',
        'male', '1988-05-05 00:00:00', 'active', '2026-05-01 09:40:00.000', '2026-05-01 09:40:00.000'
    );

INSERT INTO trips (
    id, driver_id, vehicle_type, start_city, end_city,
    departure_time, arrival_time, seat_total, seat_available,
    price_cent, status, created_at, updated_at
) VALUES
    (
        1, 1, 'car', 'Hangzhou', 'Suzhou',
        '2026-05-07 08:30:00', '2026-05-07 11:10:00', 24, 21,
        16800, 'published', '2026-05-05 10:00:00.000', '2026-05-05 10:00:00.000'
    ),
    (
        2, 1, 'car', 'Suzhou', 'Nanjing',
        '2026-05-07 13:20:00', '2026-05-07 15:50:00', 18, 16,
        16000, 'published', '2026-05-05 10:30:00.000', '2026-05-05 10:30:00.000'
    ),
    (
        3, 2, 'car', 'Shanghai', 'Hangzhou',
        '2026-05-08 09:00:00', '2026-05-08 11:30:00', 30, 26,
        18000, 'published', '2026-05-05 11:00:00.000', '2026-05-05 11:00:00.000'
    ),
    (
        4, 2, 'car', 'Nanjing', 'Wuxi',
        '2026-05-04 18:00:00', '2026-05-04 20:00:00', 12, 12,
        12000, 'closed', '2026-05-03 16:00:00.000', '2026-05-03 16:00:00.000'
    );

INSERT INTO trip_stops (
    id, trip_id, stop_order, stop_name,
    plan_arrival_time, plan_departure_time, created_at, updated_at
) VALUES
    (1, 1, 1, 'Jiaxing',  '2026-05-07 09:30:00', '2026-05-07 09:35:00', '2026-05-05 10:00:00.000', '2026-05-05 10:00:00.000'),
    (2, 1, 2, 'Kunshan',  '2026-05-07 10:35:00', '2026-05-07 10:40:00', '2026-05-05 10:00:00.000', '2026-05-05 10:00:00.000'),
    (3, 2, 1, 'Changzhou','2026-05-07 14:25:00', '2026-05-07 14:30:00', '2026-05-05 10:30:00.000', '2026-05-05 10:30:00.000'),
    (4, 3, 1, 'Songjiang', '2026-05-08 09:35:00', '2026-05-08 09:40:00', '2026-05-05 11:00:00.000', '2026-05-05 11:00:00.000'),
    (5, 3, 2, 'Tongxiang', '2026-05-08 10:35:00', '2026-05-08 10:40:00', '2026-05-05 11:00:00.000', '2026-05-05 11:00:00.000'),
    (6, 4, 1, 'Changzhou', '2026-05-04 19:00:00', '2026-05-04 19:05:00', '2026-05-03 16:00:00.000', '2026-05-03 16:00:00.000');

INSERT INTO orders (
    id, order_no, user_id, trip_id, ticket_count, seat_type,
    amount, pay_status, order_status, refund_status, created_at, updated_at
) VALUES
    (
        1, 'ORD202605060001', 3, 1, 2, 'standard',
        33600, 'paid', 'pending_verification', 'none',
        '2026-05-06 08:15:00.000', '2026-05-06 08:15:00.000'
    ),
    (
        2, 'ORD202605060002', 4, 1, 1, 'standard',
        16800, 'paid', 'completed', 'none',
        '2026-05-06 08:40:00.000', '2026-05-06 08:40:00.000'
    ),
    (
        3, 'ORD202605060003', 3, 2, 2, 'standard',
        32000, 'paid', 'pending_verification', 'requested',
        '2026-05-06 11:00:00.000', '2026-05-06 11:20:00.000'
    ),
    (
        4, 'ORD202605060004', 4, 3, 4, 'standard',
        72000, 'paid', 'completed', 'none',
        '2026-05-06 12:00:00.000', '2026-05-06 12:30:00.000'
    );
