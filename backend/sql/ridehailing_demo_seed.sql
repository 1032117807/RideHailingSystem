-- TripVerse demo database bootstrap
-- MySQL 8.0+ recommended
-- Default password for seeded users: 123456
-- Admin account:
--   phone: 18800000000
--   email: admin_root@tripverse.local

DROP DATABASE IF EXISTS ridehailing_demo;
CREATE DATABASE ridehailing_demo
  DEFAULT CHARACTER SET utf8mb4
  DEFAULT COLLATE utf8mb4_general_ci;

USE ridehailing_demo;

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS risk_events;
DROP TABLE IF EXISTS token_usages;
DROP TABLE IF EXISTS notifications;
DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS trip_stops;
DROP TABLE IF EXISTS trips;
DROP TABLE IF EXISTS users;

CREATE TABLE users (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  phone VARCHAR(20) NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  nickname VARCHAR(50) NOT NULL,
  role VARCHAR(20) NOT NULL DEFAULT 'passenger',
  default_role VARCHAR(20) NOT NULL DEFAULT 'passenger',
  real_name VARCHAR(50) NULL,
  id_card VARCHAR(32) NULL,
  real_name_verified TINYINT(1) NOT NULL DEFAULT 0,
  avatar VARCHAR(255) NULL,
  email VARCHAR(100) NOT NULL,
  gender VARCHAR(20) NULL,
  birthday DATETIME NULL,
  status VARCHAR(20) NOT NULL DEFAULT 'active',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uk_users_phone (phone),
  UNIQUE KEY uk_users_email (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE trips (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
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
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  KEY idx_driver_status (driver_id, status),
  KEY idx_trip_search (start_city, end_city, departure_time, status),
  CONSTRAINT fk_trips_driver FOREIGN KEY (driver_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE trip_stops (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  trip_id BIGINT UNSIGNED NOT NULL,
  stop_order INT NOT NULL,
  stop_name VARCHAR(50) NOT NULL,
  plan_arrival_time DATETIME NULL,
  plan_departure_time DATETIME NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uk_trip_stop_order (trip_id, stop_order),
  CONSTRAINT fk_trip_stops_trip FOREIGN KEY (trip_id) REFERENCES trips(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE orders (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  order_no VARCHAR(32) NOT NULL,
  user_id BIGINT UNSIGNED NOT NULL,
  trip_id BIGINT UNSIGNED NOT NULL,
  ticket_count INT NOT NULL,
  seat_type VARCHAR(30) NOT NULL DEFAULT 'standard',
  amount INT NOT NULL,
  pay_status VARCHAR(20) NOT NULL DEFAULT 'unpaid',
  order_status VARCHAR(30) NOT NULL DEFAULT 'pending_payment',
  refund_status VARCHAR(20) NOT NULL DEFAULT 'none',
  refund_review_note VARCHAR(255) NULL,
  refund_reviewed_at DATETIME NULL,
  payment_expire_at DATETIME NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uk_orders_order_no (order_no),
  KEY idx_orders_user_id (user_id),
  KEY idx_orders_trip_id (trip_id),
  KEY idx_orders_payment_expire_at (payment_expire_at),
  CONSTRAINT fk_orders_user FOREIGN KEY (user_id) REFERENCES users(id),
  CONSTRAINT fk_orders_trip FOREIGN KEY (trip_id) REFERENCES trips(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE payments (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  payment_no VARCHAR(32) NOT NULL,
  order_id BIGINT UNSIGNED NOT NULL,
  user_id BIGINT UNSIGNED NOT NULL,
  amount INT NOT NULL,
  channel VARCHAR(20) NOT NULL DEFAULT 'mock',
  status VARCHAR(20) NOT NULL DEFAULT 'pending',
  paid_at DATETIME NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY uk_payments_payment_no (payment_no),
  KEY idx_payments_order_id (order_id),
  KEY idx_payments_user_id (user_id),
  CONSTRAINT fk_payments_order FOREIGN KEY (order_id) REFERENCES orders(id),
  CONSTRAINT fk_payments_user FOREIGN KEY (user_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE notifications (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  user_id BIGINT UNSIGNED NOT NULL,
  type VARCHAR(50) NOT NULL,
  title VARCHAR(100) NOT NULL,
  content TEXT NOT NULL,
  related_order_id BIGINT UNSIGNED NULL,
  is_read TINYINT(1) NOT NULL DEFAULT 0,
  read_at DATETIME NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  KEY idx_notifications_user_read_created (user_id, is_read, created_at),
  KEY idx_notifications_related_order_id (related_order_id),
  CONSTRAINT fk_notifications_user FOREIGN KEY (user_id) REFERENCES users(id),
  CONSTRAINT fk_notifications_order FOREIGN KEY (related_order_id) REFERENCES orders(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE token_usages (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  user_id BIGINT UNSIGNED NOT NULL,
  role VARCHAR(20) NOT NULL,
  feature VARCHAR(50) NOT NULL,
  request_kind VARCHAR(20) NOT NULL,
  provider VARCHAR(50) NOT NULL,
  model VARCHAR(100) NOT NULL,
  prompt_tokens INT NOT NULL DEFAULT 0,
  completion_tokens INT NOT NULL DEFAULT 0,
  total_tokens INT NOT NULL DEFAULT 0,
  request_count INT NOT NULL DEFAULT 1,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  KEY idx_token_usages_user_id (user_id),
  KEY idx_token_usages_role (role),
  KEY idx_token_usages_feature (feature),
  KEY idx_token_usages_request_kind (request_kind),
  KEY idx_token_usages_model (model),
  KEY idx_token_usages_created_at (created_at),
  CONSTRAINT fk_token_usages_user FOREIGN KEY (user_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE risk_events (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  severity VARCHAR(20) NOT NULL,
  event_type VARCHAR(50) NOT NULL,
  subject_type VARCHAR(30) NOT NULL,
  subject_id VARCHAR(100) NOT NULL,
  fingerprint VARCHAR(150) NOT NULL,
  title VARCHAR(255) NOT NULL,
  detail TEXT NOT NULL,
  status VARCHAR(20) NOT NULL DEFAULT 'open',
  metrics_json LONGTEXT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  KEY idx_risk_events_severity (severity),
  KEY idx_risk_events_event_type (event_type),
  KEY idx_risk_events_subject_type (subject_type),
  KEY idx_risk_events_subject_id (subject_id),
  KEY idx_risk_events_fingerprint (fingerprint),
  KEY idx_risk_events_status (status),
  KEY idx_risk_events_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

SET @demo_now = TIMESTAMP('2026-05-10 12:00:00');
SET @password_hash = '$2a$10$MiDz18JPMzRMdnRVcAkj5O1IR2rMgzo2mG8oWdMugByyzuAFdpyXm';

INSERT INTO users (
  id, phone, password_hash, nickname, role, default_role, real_name, id_card,
  real_name_verified, avatar, email, gender, birthday, status, created_at, updated_at
) VALUES (
  1, '18800000000', @password_hash, 'admin_root', 'admin', 'admin', '系统管理员',
  '110101198001010001', 1, '', 'admin_root@tripverse.local', 'unknown',
  '1980-01-01 00:00:00', 'active', @demo_now, @demo_now
);

CREATE TEMPORARY TABLE seed_numbers
SELECT ones.d + tens.d * 10 + hundreds.d * 100 + 1 AS n
FROM (
  SELECT 0 AS d UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4
  UNION ALL SELECT 5 UNION ALL SELECT 6 UNION ALL SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9
) ones
CROSS JOIN (
  SELECT 0 AS d UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4
  UNION ALL SELECT 5 UNION ALL SELECT 6 UNION ALL SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9
) tens
CROSS JOIN (
  SELECT 0 AS d UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4
  UNION ALL SELECT 5 UNION ALL SELECT 6 UNION ALL SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9
) hundreds;

ALTER TABLE seed_numbers ADD PRIMARY KEY (n);

INSERT INTO users (
  phone, password_hash, nickname, role, default_role, real_name, id_card,
  real_name_verified, avatar, email, gender, birthday, status, created_at, updated_at
)
SELECT
  CONCAT('13', LPAD(n, 9, '0')),
  @password_hash,
  CONCAT('passenger_', LPAD(n, 3, '0')),
  'passenger',
  'passenger',
  CONCAT('乘客', LPAD(n, 3, '0')),
  CONCAT('1101011990', LPAD(n, 8, '0')),
  CASE WHEN MOD(n, 4) = 0 THEN 0 ELSE 1 END,
  '',
  CONCAT('passenger', LPAD(n, 3, '0'), '@tripverse.local'),
  CASE WHEN MOD(n, 2) = 0 THEN 'female' ELSE 'male' END,
  DATE_ADD('1990-01-01 00:00:00', INTERVAL n DAY),
  CASE
    WHEN MOD(n, 53) = 0 THEN 'frozen'
    WHEN MOD(n, 37) = 0 THEN 'disabled'
    ELSE 'active'
  END,
  DATE_SUB(@demo_now, INTERVAL MOD(n, 180) DAY),
  DATE_SUB(@demo_now, INTERVAL MOD(n, 30) DAY)
FROM seed_numbers
WHERE n <= 180
ORDER BY n;

INSERT INTO users (
  phone, password_hash, nickname, role, default_role, real_name, id_card,
  real_name_verified, avatar, email, gender, birthday, status, created_at, updated_at
)
SELECT
  CONCAT('15', LPAD(n, 9, '0')),
  @password_hash,
  CONCAT('driver_', LPAD(n, 3, '0')),
  'driver',
  'driver',
  CONCAT('司机', LPAD(n, 3, '0')),
  CONCAT('3201011988', LPAD(n, 8, '0')),
  1,
  '',
  CONCAT('driver', LPAD(n, 3, '0'), '@tripverse.local'),
  CASE WHEN MOD(n, 2) = 0 THEN 'female' ELSE 'male' END,
  DATE_ADD('1988-01-01 00:00:00', INTERVAL n DAY),
  CASE
    WHEN MOD(n, 17) = 0 THEN 'frozen'
    WHEN MOD(n, 29) = 0 THEN 'disabled'
    ELSE 'active'
  END,
  DATE_SUB(@demo_now, INTERVAL MOD(n, 220) DAY),
  DATE_SUB(@demo_now, INTERVAL MOD(n, 20) DAY)
FROM seed_numbers
WHERE n <= 40
ORDER BY n;

CREATE TEMPORARY TABLE seed_routes (
  route_no INT PRIMARY KEY,
  start_city VARCHAR(50) NOT NULL,
  end_city VARCHAR(50) NOT NULL,
  vehicle_type VARCHAR(20) NOT NULL,
  departure_hour INT NOT NULL,
  departure_minute INT NOT NULL,
  duration_minute INT NOT NULL,
  seat_total INT NOT NULL,
  price_cent INT NOT NULL
);

INSERT INTO seed_routes VALUES
(1, '北京', '天津', '城际快线', 7, 30, 150, 28, 9800),
(2, '北京', '石家庄', '商务大巴', 8, 10, 210, 38, 12800),
(3, '天津', '济南', '商务大巴', 9, 0, 240, 40, 15600),
(4, '石家庄', '太原', '拼车专线', 8, 40, 220, 18, 11800),
(5, '太原', '西安', '商务大巴', 10, 15, 330, 36, 18800),
(6, '呼和浩特', '银川', '商务大巴', 7, 50, 420, 34, 22800),
(7, '沈阳', '大连', '城际快线', 6, 55, 250, 32, 16800),
(8, '沈阳', '长春', '商务大巴', 9, 20, 230, 30, 14800),
(9, '长春', '哈尔滨', '城际快线', 11, 0, 230, 28, 14600),
(10, '哈尔滨', '沈阳', '商务大巴', 7, 45, 420, 36, 23800),
(11, '上海', '南京', '城际快线', 7, 15, 210, 30, 12800),
(12, '上海', '杭州', '城际快线', 8, 5, 180, 30, 10800),
(13, '南京', '合肥', '拼车专线', 9, 35, 150, 16, 9200),
(14, '杭州', '宁波', '城际快线', 10, 20, 150, 24, 8800),
(15, '福州', '厦门', '商务大巴', 9, 10, 210, 30, 12600),
(16, '南昌', '武汉', '商务大巴', 8, 25, 260, 36, 16200),
(17, '济南', '青岛', '城际快线', 10, 50, 220, 30, 13800),
(18, '郑州', '武汉', '商务大巴', 9, 30, 300, 34, 17600),
(19, '郑州', '西安', '商务大巴', 7, 40, 320, 34, 18200),
(20, '武汉', '长沙', '城际快线', 13, 15, 210, 28, 11800),
(21, '广州', '深圳', '城际快线', 7, 20, 150, 26, 9800),
(22, '广州', '南宁', '商务大巴', 8, 45, 420, 40, 24600),
(23, '深圳', '厦门', '商务大巴', 9, 50, 330, 34, 19800),
(24, '南宁', '海口', '商务大巴', 6, 40, 480, 36, 26800),
(25, '成都', '重庆', '城际快线', 7, 5, 220, 30, 12600),
(26, '成都', '贵阳', '商务大巴', 8, 15, 360, 36, 21800),
(27, '贵阳', '昆明', '商务大巴', 9, 5, 300, 32, 18800),
(28, '昆明', '南宁', '商务大巴', 7, 35, 420, 36, 24800),
(29, '西安', '兰州', '商务大巴', 9, 45, 330, 34, 19600),
(30, '兰州', '西宁', '拼车专线', 13, 5, 150, 16, 9800),
(31, '西宁', '乌鲁木齐', '商务大巴', 6, 25, 900, 38, 46800),
(32, '银川', '乌鲁木齐', '商务大巴', 7, 15, 960, 38, 49800),
(33, '重庆', '武汉', '商务大巴', 10, 35, 540, 40, 28600),
(34, '长沙', '南昌', '城际快线', 14, 10, 210, 24, 11600),
(35, '合肥', '上海', '城际快线', 8, 55, 240, 28, 13600),
(36, '广州', '海口', '商务大巴', 9, 25, 600, 40, 29800),
(37, '北京', '上海', '商务大巴', 7, 0, 720, 42, 39800),
(38, '深圳', '长沙', '商务大巴', 10, 0, 480, 36, 25800),
(39, '哈尔滨', '北京', '商务大巴', 6, 30, 780, 40, 42800),
(40, '拉萨', '成都', '商务大巴', 5, 50, 1260, 34, 68800);

CREATE TEMPORARY TABLE seed_route_stops (
  route_no INT NOT NULL,
  stop_order INT NOT NULL,
  stop_name VARCHAR(50) NOT NULL,
  arrival_offset_min INT NOT NULL,
  departure_offset_min INT NOT NULL,
  PRIMARY KEY (route_no, stop_order)
);

INSERT INTO seed_route_stops VALUES
(1, 1, '廊坊', 60, 70),
(2, 1, '保定', 90, 100),
(3, 1, '德州', 120, 130),
(4, 1, '阳泉', 110, 120),
(5, 1, '临汾', 150, 160),
(6, 1, '鄂尔多斯', 200, 215),
(7, 1, '营口', 110, 120),
(8, 1, '四平', 105, 115),
(9, 1, '大庆', 120, 130),
(10, 1, '长春', 210, 225),
(11, 1, '苏州', 95, 105),
(12, 1, '嘉兴', 80, 90),
(13, 1, '滁州', 65, 75),
(14, 1, '绍兴', 65, 75),
(15, 1, '泉州', 105, 115),
(16, 1, '九江', 120, 130),
(17, 1, '潍坊', 115, 125),
(18, 1, '信阳', 145, 155),
(19, 1, '洛阳', 135, 145),
(20, 1, '岳阳', 100, 110),
(21, 1, '东莞', 70, 80),
(22, 1, '梧州', 210, 225),
(23, 1, '汕尾', 150, 160),
(24, 1, '钦州', 220, 235),
(25, 1, '遂宁', 105, 115),
(26, 1, '宜宾', 170, 180),
(27, 1, '曲靖', 140, 150),
(28, 1, '百色', 205, 220),
(29, 1, '天水', 150, 160),
(30, 1, '海东', 70, 80),
(31, 1, '张掖', 360, 380),
(32, 1, '武威', 320, 340),
(33, 1, '恩施', 260, 280),
(34, 1, '萍乡', 95, 105),
(35, 1, '南京', 120, 130),
(36, 1, '湛江', 280, 300),
(37, 1, '济南', 360, 380),
(38, 1, '韶关', 210, 225),
(39, 1, '长春', 330, 350),
(40, 1, '林芝', 520, 540);

CREATE TEMPORARY TABLE seed_day_offsets (day_offset INT PRIMARY KEY);
INSERT INTO seed_day_offsets VALUES (-2),(-1),(0),(1),(2),(3),(4);

SET @driver_rn := 0;
CREATE TEMPORARY TABLE seed_drivers AS
SELECT (@driver_rn := @driver_rn + 1) AS rn, id
FROM users
WHERE role = 'driver' AND status = 'active'
ORDER BY id;

SET @driver_count := (SELECT COUNT(*) FROM seed_drivers);

INSERT INTO trips (
  driver_id, vehicle_type, start_city, end_city, departure_time, arrival_time,
  seat_total, seat_available, price_cent, status, created_at, updated_at
)
SELECT
  d.id,
  r.vehicle_type,
  r.start_city,
  r.end_city,
  TIMESTAMP(DATE_ADD(DATE(@demo_now), INTERVAL dy.day_offset DAY), MAKETIME(r.departure_hour, r.departure_minute, 0)),
  DATE_ADD(
    TIMESTAMP(DATE_ADD(DATE(@demo_now), INTERVAL dy.day_offset DAY), MAKETIME(r.departure_hour, r.departure_minute, 0)),
    INTERVAL r.duration_minute MINUTE
  ),
  r.seat_total,
  r.seat_total,
  r.price_cent,
  CASE WHEN dy.day_offset < 0 THEN 'closed' ELSE 'published' END,
  DATE_SUB(@demo_now, INTERVAL (ABS(dy.day_offset) + 7) DAY),
  DATE_SUB(@demo_now, INTERVAL ABS(dy.day_offset) DAY)
FROM seed_routes r
JOIN seed_day_offsets dy
JOIN seed_drivers d
  ON d.rn = MOD(r.route_no + dy.day_offset + @driver_count * 10, @driver_count) + 1
ORDER BY r.route_no, dy.day_offset;

INSERT INTO trip_stops (
  trip_id, stop_order, stop_name, plan_arrival_time, plan_departure_time, created_at, updated_at
)
SELECT
  t.id,
  s.stop_order,
  s.stop_name,
  DATE_ADD(t.departure_time, INTERVAL s.arrival_offset_min MINUTE),
  DATE_ADD(t.departure_time, INTERVAL s.departure_offset_min MINUTE),
  t.created_at,
  t.updated_at
FROM trips t
JOIN seed_routes r
  ON r.start_city = t.start_city
 AND r.end_city = t.end_city
JOIN seed_route_stops s
  ON s.route_no = r.route_no
ORDER BY t.id, s.stop_order;

SET @passenger_rn := 0;
CREATE TEMPORARY TABLE seed_passengers AS
SELECT (@passenger_rn := @passenger_rn + 1) AS rn, id
FROM users
WHERE role = 'passenger' AND status = 'active'
ORDER BY id;

SET @trip_rn := 0;
CREATE TEMPORARY TABLE seed_trip_ids AS
SELECT (@trip_rn := @trip_rn + 1) AS rn, id, price_cent, departure_time
FROM trips
ORDER BY id;

SET @passenger_count := (SELECT COUNT(*) FROM seed_passengers);
SET @trip_count := (SELECT COUNT(*) FROM seed_trip_ids);

INSERT INTO orders (
  order_no, user_id, trip_id, ticket_count, seat_type, amount,
  pay_status, order_status, refund_status, refund_review_note,
  refund_reviewed_at, payment_expire_at, created_at, updated_at
)
SELECT
  x.order_no,
  x.user_id,
  x.trip_id,
  x.ticket_count,
  x.seat_type,
  x.amount,
  x.pay_status,
  x.order_status,
  x.refund_status,
  x.refund_review_note,
  x.refund_reviewed_at,
  x.payment_expire_at,
  x.created_at,
  x.updated_at
FROM (
  SELECT
    CONCAT('ORD202605', LPAD(n, 6, '0')) AS order_no,
    p.id AS user_id,
    t.id AS trip_id,
    CASE WHEN MOD(n, 11) = 0 THEN 2 ELSE 1 END AS ticket_count,
    'standard' AS seat_type,
    t.price_cent * CASE WHEN MOD(n, 11) = 0 THEN 2 ELSE 1 END AS amount,
    CASE
      WHEN MOD(n, 10) IN (1,2,3,4,5,6,7) THEN 'paid'
      ELSE 'unpaid'
    END AS pay_status,
    CASE
      WHEN MOD(n, 17) = 0 THEN 'cancelled'
      WHEN t.departure_time < @demo_now AND MOD(n, 5) IN (0,1) THEN 'completed'
      WHEN MOD(n, 10) IN (8,9) THEN 'pending_payment'
      ELSE 'pending_verification'
    END AS order_status,
    CASE
      WHEN MOD(n, 23) = 0 THEN 'requested'
      WHEN MOD(n, 29) = 0 THEN 'rejected'
      WHEN MOD(n, 31) = 0 THEN 'refunded'
      ELSE 'none'
    END AS refund_status,
    CASE
      WHEN MOD(n, 23) = 0 THEN '退款申请已提交，等待管理员审核'
      WHEN MOD(n, 29) = 0 THEN '资料不完整，退款申请已驳回'
      WHEN MOD(n, 31) = 0 THEN '退款已原路退回，请注意查收'
      ELSE ''
    END AS refund_review_note,
    CASE
      WHEN MOD(n, 29) = 0 OR MOD(n, 31) = 0 THEN DATE_SUB(@demo_now, INTERVAL MOD(n, 72) HOUR)
      ELSE NULL
    END AS refund_reviewed_at,
    CASE
      WHEN MOD(n, 10) IN (8,9) THEN DATE_ADD(DATE_SUB(@demo_now, INTERVAL MOD(n, 36) HOUR), INTERVAL 15 MINUTE)
      ELSE NULL
    END AS payment_expire_at,
    DATE_SUB(@demo_now, INTERVAL MOD(n, 360) HOUR) AS created_at,
    DATE_SUB(@demo_now, INTERVAL MOD(n, 180) HOUR) AS updated_at
  FROM seed_numbers s
  JOIN seed_passengers p
    ON p.rn = MOD(s.n - 1, @passenger_count) + 1
  JOIN seed_trip_ids t
    ON t.rn = MOD(s.n * 7 - 1, @trip_count) + 1
  WHERE s.n <= 900
) AS x;

UPDATE trips t
LEFT JOIN (
  SELECT trip_id, SUM(CASE WHEN order_status <> 'cancelled' THEN ticket_count ELSE 0 END) AS sold_count
  FROM orders
  GROUP BY trip_id
) s ON s.trip_id = t.id
SET
  t.seat_available = GREATEST(t.seat_total - IFNULL(s.sold_count, 0), 0),
  t.updated_at = @demo_now;

INSERT INTO payments (
  payment_no, order_id, user_id, amount, channel, status, paid_at, created_at, updated_at
)
SELECT
  CONCAT('PAY202605', LPAD(o.id, 6, '0')),
  o.id,
  o.user_id,
  o.amount,
  'mock',
  CASE
    WHEN o.pay_status = 'paid' THEN 'paid'
    WHEN o.order_status = 'cancelled' THEN 'closed'
    ELSE 'pending'
  END,
  CASE WHEN o.pay_status = 'paid' THEN DATE_ADD(o.created_at, INTERVAL 5 MINUTE) ELSE NULL END,
  o.created_at,
  o.updated_at
FROM orders o;

INSERT INTO notifications (
  user_id, type, title, content, related_order_id, is_read, read_at, created_at, updated_at
)
SELECT
  o.user_id,
  CASE
    WHEN o.refund_status = 'refunded' THEN 'refund_approved'
    WHEN o.refund_status = 'rejected' THEN 'refund_rejected'
    ELSE 'order_expired'
  END,
  CASE
    WHEN o.refund_status = 'refunded' THEN '退款审核通过'
    WHEN o.refund_status = 'rejected' THEN '退款审核驳回'
    ELSE '订单已过期'
  END,
  CASE
    WHEN o.refund_status = 'refunded' THEN CONCAT('订单 ', o.order_no, ' 的退款已处理完成，请注意查收。')
    WHEN o.refund_status = 'rejected' THEN CONCAT('订单 ', o.order_no, ' 的退款申请被驳回，请查看原因。')
    ELSE CONCAT('订单 ', o.order_no, ' 因超时未支付已自动取消。')
  END,
  o.id,
  CASE WHEN MOD(o.id, 3) = 0 THEN 1 ELSE 0 END,
  CASE WHEN MOD(o.id, 3) = 0 THEN DATE_ADD(o.updated_at, INTERVAL 2 HOUR) ELSE NULL END,
  DATE_ADD(o.updated_at, INTERVAL 1 HOUR),
  DATE_ADD(o.updated_at, INTERVAL 1 HOUR)
FROM orders o
WHERE o.refund_status IN ('refunded', 'rejected') OR (o.order_status = 'cancelled' AND o.pay_status = 'unpaid');

SET @active_user_rn := 0;
CREATE TEMPORARY TABLE seed_active_users AS
SELECT (@active_user_rn := @active_user_rn + 1) AS rn, id, role
FROM users
WHERE status = 'active'
ORDER BY id;

SET @active_user_count := (SELECT COUNT(*) FROM seed_active_users);

INSERT INTO token_usages (
  user_id, role, feature, request_kind, provider, model,
  prompt_tokens, completion_tokens, total_tokens, request_count, created_at
)
SELECT
  u.id,
  u.role,
  CASE MOD(s.n, 4)
    WHEN 0 THEN 'passenger_ai'
    WHEN 1 THEN 'driver_ai_draft'
    WHEN 2 THEN 'knowledge_search'
    ELSE 'knowledge_ingest'
  END,
  CASE MOD(s.n, 3)
    WHEN 0 THEN 'chat'
    WHEN 1 THEN 'embedding'
    ELSE 'rerank'
  END,
  'openai-compatible',
  CASE MOD(s.n, 3)
    WHEN 0 THEN 'qwen-plus'
    WHEN 1 THEN 'text-embedding-v3'
    ELSE 'rerank-v1'
  END,
  200 + MOD(s.n * 13, 1200),
  CASE WHEN MOD(s.n, 3) = 0 THEN 80 + MOD(s.n * 7, 900) ELSE 0 END,
  (200 + MOD(s.n * 13, 1200)) + CASE WHEN MOD(s.n, 3) = 0 THEN 80 + MOD(s.n * 7, 900) ELSE 0 END,
  1,
  DATE_SUB(@demo_now, INTERVAL MOD(s.n, 720) HOUR)
FROM seed_numbers s
JOIN seed_active_users u
  ON u.rn = MOD(s.n * 5 - 1, @active_user_count) + 1
WHERE s.n <= 1000;

INSERT INTO risk_events (
  severity, event_type, subject_type, subject_id, fingerprint, title, detail,
  status, metrics_json, created_at, updated_at
) VALUES
('high', 'ai_rate_limit', 'user', '12', 'ai_rate_limit:passenger-chat:12', 'AI 接口限流触发', '用户 12 在 passenger-chat 触发了 AI 限流，系统已拒绝继续请求。', 'open', '{"scope":"passenger-chat","subject":"user:12","limit":12,"retryAfterSecond":60}', DATE_SUB(@demo_now, INTERVAL 2 HOUR), DATE_SUB(@demo_now, INTERVAL 2 HOUR)),
('high', 'ai_rate_limit', 'user', '207', 'ai_rate_limit:driver-create-trip:207', 'AI 接口限流触发', '司机 207 在 driver-create-trip 触发了 AI 限流，存在高频生成班次草稿行为。', 'acknowledged', '{"scope":"driver-create-trip","subject":"user:207","limit":6,"retryAfterSecond":60}', DATE_SUB(@demo_now, INTERVAL 5 HOUR), DATE_SUB(@demo_now, INTERVAL 4 HOUR)),
('medium', 'token_spike', 'user', '45', 'token_spike:45:202605101130', 'Token 使用量异常激增', '用户 45 最近 15 分钟的 token 使用量明显高于历史平均值。', 'open', '{"currentTotalTokens":15680,"currentRequestCount":18,"baselineTokensPer15":2980.5,"baselineRequestsPer15":4.2}', DATE_SUB(@demo_now, INTERVAL 40 MINUTE), DATE_SUB(@demo_now, INTERVAL 40 MINUTE)),
('medium', 'token_spike', 'user', '88', 'token_spike:88:202605101115', 'Token 使用量异常激增', '用户 88 最近 15 分钟的请求数和 token 消耗同时偏高。', 'resolved', '{"currentTotalTokens":12120,"currentRequestCount":14,"baselineTokensPer15":2440.3,"baselineRequestsPer15":3.1}', DATE_SUB(@demo_now, INTERVAL 55 MINUTE), DATE_SUB(@demo_now, INTERVAL 20 MINUTE)),
('low', 'token_spike', 'user', '1', 'token_spike:1:202605100930', '后台检索测试量偏高', '管理员最近一小时频繁执行知识库检索测试，已记录为低优先级观察项。', 'acknowledged', '{"currentTotalTokens":4200,"currentRequestCount":11,"baselineTokensPer15":1800.0,"baselineRequestsPer15":2.0}', DATE_SUB(@demo_now, INTERVAL 3 HOUR), DATE_SUB(@demo_now, INTERVAL 2 HOUR));

SET FOREIGN_KEY_CHECKS = 1;

-- Quick checks
SELECT COUNT(*) AS user_count FROM users;
SELECT COUNT(*) AS trip_count FROM trips;
SELECT COUNT(*) AS order_count FROM orders;
SELECT COUNT(*) AS payment_count FROM payments;
SELECT COUNT(*) AS notification_count FROM notifications;
SELECT COUNT(*) AS token_usage_count FROM token_usages;
SELECT COUNT(*) AS risk_event_count FROM risk_events;
