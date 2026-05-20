-- TripVerse demo database bootstrap
-- MySQL 8.0+ recommended
-- Default password for seeded users: 123456
-- Admin account:
--   phone: 18800000000
--   email: admin_root@tripverse.local

DROP DATABASE IF EXISTS ridehailing_demo;
CREATE DATABASE ridehailing_demo DEFAULT CHARACTER SET utf8mb4 DEFAULT COLLATE utf8mb4_general_ci;
USE ridehailing_demo;

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS ticket_verifications;
DROP TABLE IF EXISTS electronic_tickets;
DROP TABLE IF EXISTS refund_audit_logs;
DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS driver_settlements;
DROP TABLE IF EXISTS vehicles;
DROP TABLE IF EXISTS driver_profiles;
DROP TABLE IF EXISTS price_alerts;
DROP TABLE IF EXISTS passengers;
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
  UNIQUE KEY idx_trip_stop_order (trip_id, stop_order),
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
  KEY idx_risk_events_status (status),
  KEY idx_risk_events_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE passengers (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  user_id BIGINT UNSIGNED NOT NULL,
  name VARCHAR(50) NOT NULL,
  id_card VARCHAR(32) NOT NULL,
  phone VARCHAR(20) NOT NULL,
  is_default TINYINT(1) NOT NULL DEFAULT 0,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id), KEY idx_passengers_user_id (user_id),
  CONSTRAINT fk_passengers_user FOREIGN KEY (user_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE price_alerts (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  user_id BIGINT UNSIGNED NOT NULL,
  start_city VARCHAR(50) NOT NULL,
  end_city VARCHAR(50) NOT NULL,
  target_price_cent INT NOT NULL,
  start_date DATE NOT NULL,
  end_date DATE NOT NULL,
  status VARCHAR(20) NOT NULL DEFAULT 'active',
  triggered_trip_id BIGINT UNSIGNED NULL,
  triggered_at DATETIME NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id), KEY idx_price_alerts_user_id (user_id), KEY idx_price_alerts_route (start_city, end_city, status),
  CONSTRAINT fk_price_alerts_user FOREIGN KEY (user_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE driver_profiles (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  user_id BIGINT UNSIGNED NOT NULL,
  real_name VARCHAR(50) NOT NULL,
  id_card VARCHAR(32) NOT NULL,
  license_no VARCHAR(64) NOT NULL,
  id_card_image_url VARCHAR(255) NULL,
  license_image_url VARCHAR(255) NULL,
  status VARCHAR(20) NOT NULL DEFAULT 'approved',
  review_note VARCHAR(255) NULL,
  reviewed_at DATETIME NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  UNIQUE KEY idx_driver_profiles_user_id (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE vehicles (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  driver_id BIGINT UNSIGNED NOT NULL,
  plate_no VARCHAR(32) NOT NULL,
  brand VARCHAR(64) NULL,
  model_name VARCHAR(64) NULL,
  seat_count INT NOT NULL,
  status VARCHAR(20) NOT NULL DEFAULT 'active',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id), KEY idx_vehicles_driver_id (driver_id),
  CONSTRAINT fk_vehicles_driver FOREIGN KEY (driver_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE driver_settlements (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  driver_id BIGINT UNSIGNED NOT NULL,
  settlement_date DATE NOT NULL,
  gross_amount_cent INT NOT NULL,
  refund_amount_cent INT NOT NULL,
  service_fee_cent INT NOT NULL,
  net_amount_cent INT NOT NULL,
  status VARCHAR(20) NOT NULL DEFAULT 'pending',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id), KEY idx_driver_settlements_driver_date (driver_id, settlement_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE electronic_tickets (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  order_id BIGINT UNSIGNED NOT NULL,
  user_id BIGINT UNSIGNED NOT NULL,
  trip_id BIGINT UNSIGNED NOT NULL,
  token TEXT NOT NULL,
  token_hash VARCHAR(64) NOT NULL,
  status VARCHAR(20) NOT NULL DEFAULT 'issued',
  expires_at DATETIME NOT NULL,
  verified_at DATETIME NULL,
  verified_by_driver_id BIGINT UNSIGNED NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id), UNIQUE KEY uk_electronic_tickets_order_id (order_id), UNIQUE KEY uk_electronic_tickets_token_hash (token_hash)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE ticket_verifications (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  ticket_id BIGINT UNSIGNED NOT NULL,
  order_id BIGINT UNSIGNED NOT NULL,
  driver_id BIGINT UNSIGNED NOT NULL,
  result VARCHAR(20) NOT NULL,
  message VARCHAR(255) NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id), KEY idx_ticket_verifications_ticket_id (ticket_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE refund_audit_logs (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  order_id BIGINT UNSIGNED NOT NULL,
  refund_status VARCHAR(20) NOT NULL,
  review_note VARCHAR(255) NULL,
  reviewer_id BIGINT UNSIGNED NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id), KEY idx_refund_audit_logs_order_id (order_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE audit_logs (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  actor_user_id BIGINT UNSIGNED NOT NULL,
  actor_role VARCHAR(20) NOT NULL,
  action VARCHAR(80) NOT NULL,
  resource_type VARCHAR(80) NOT NULL,
  resource_id VARCHAR(80) NOT NULL,
  detail TEXT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id), KEY idx_audit_logs_action (action), KEY idx_audit_logs_resource (resource_type, resource_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

SET @demo_now = TIMESTAMP('2026-05-20 12:00:00');
SET @password_hash = '$2a$10$MiDz18JPMzRMdnRVcAkj5O1IR2rMgzo2mG8oWdMugByyzuAFdpyXm';

CREATE TEMPORARY TABLE seed_numbers AS
SELECT ones.d + tens.d * 10 + hundreds.d * 100 + 1 AS n
FROM (SELECT 0 d UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4 UNION ALL SELECT 5 UNION ALL SELECT 6 UNION ALL SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9) ones
CROSS JOIN (SELECT 0 d UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4 UNION ALL SELECT 5 UNION ALL SELECT 6 UNION ALL SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9) tens
CROSS JOIN (SELECT 0 d UNION ALL SELECT 1 UNION ALL SELECT 2) hundreds;
ALTER TABLE seed_numbers ADD PRIMARY KEY (n);

INSERT INTO users (id, phone, password_hash, nickname, role, default_role, real_name, id_card, real_name_verified, avatar, email, gender, birthday, status, created_at, updated_at)
VALUES (1, '18800000000', @password_hash, 'admin_root', 'admin', 'admin', '系统管理员', '110101198001010001', 1, '', 'admin_root@tripverse.local', 'unknown', '1980-01-01 00:00:00', 'active', @demo_now, @demo_now);

INSERT INTO users (phone, password_hash, nickname, role, default_role, real_name, id_card, real_name_verified, avatar, email, gender, birthday, status, created_at, updated_at)
SELECT CONCAT('13', LPAD(n, 9, '0')), @password_hash, CONCAT('乘客', LPAD(n, 3, '0')), 'passenger', 'passenger', CONCAT('旅客', LPAD(n, 3, '0')), CONCAT('1101011990', LPAD(n, 8, '0')), IF(MOD(n,4)=0,0,1), '', CONCAT('passenger', LPAD(n,3,'0'), '@tripverse.local'), IF(MOD(n,2)=0,'female','male'), DATE_ADD('1990-01-01', INTERVAL n DAY), CASE WHEN MOD(n,53)=0 THEN 'frozen' WHEN MOD(n,37)=0 THEN 'disabled' ELSE 'active' END, DATE_SUB(@demo_now, INTERVAL MOD(n,180) DAY), DATE_SUB(@demo_now, INTERVAL MOD(n,30) DAY)
FROM seed_numbers WHERE n <= 120;

INSERT INTO users (phone, password_hash, nickname, role, default_role, real_name, id_card, real_name_verified, avatar, email, gender, birthday, status, created_at, updated_at)
SELECT CONCAT('15', LPAD(n, 9, '0')), @password_hash, CONCAT('司机', LPAD(n, 3, '0')), 'driver', 'driver', CONCAT('师傅', LPAD(n, 3, '0')), CONCAT('3201011988', LPAD(n, 8, '0')), 1, '', CONCAT('driver', LPAD(n,3,'0'), '@tripverse.local'), IF(MOD(n,2)=0,'female','male'), DATE_ADD('1988-01-01', INTERVAL n DAY), CASE WHEN MOD(n,17)=0 THEN 'frozen' WHEN MOD(n,29)=0 THEN 'disabled' ELSE 'active' END, DATE_SUB(@demo_now, INTERVAL MOD(n,160) DAY), DATE_SUB(@demo_now, INTERVAL MOD(n,20) DAY)
FROM seed_numbers WHERE n <= 24;

INSERT INTO passengers (user_id, name, id_card, phone, is_default, created_at, updated_at)
SELECT id, real_name, id_card, phone, 1, created_at, updated_at FROM users WHERE role='passenger' AND id <= 80;

INSERT INTO driver_profiles (user_id, real_name, id_card, license_no, status, review_note, reviewed_at, created_at, updated_at)
SELECT id, real_name, id_card, CONCAT('DL', LPAD(id, 8, '0')), 'approved', '资料审核通过', @demo_now, created_at, updated_at FROM users WHERE role='driver';

INSERT INTO vehicles (driver_id, plate_no, brand, model_name, seat_count, status, created_at, updated_at)
SELECT id, CONCAT('浙A', LPAD(id, 5, '0')), CASE MOD(id,4) WHEN 0 THEN '宇通' WHEN 1 THEN '比亚迪' WHEN 2 THEN '金龙' ELSE '奔驰' END, CASE MOD(id,3) WHEN 0 THEN '商务车' WHEN 1 THEN '城际大巴' ELSE '新能源巴士' END, 28 + MOD(id,12), 'active', created_at, updated_at FROM users WHERE role='driver' AND status='active';

CREATE TEMPORARY TABLE seed_routes (route_no INT PRIMARY KEY, start_city VARCHAR(50), end_city VARCHAR(50), vehicle_type VARCHAR(20), departure_hour INT, departure_minute INT, duration_minute INT, seat_total INT, price_cent INT);
INSERT INTO seed_routes VALUES
(1,'杭州','上海','car',7,30,120,32,9800),(2,'杭州','北京','car',8,0,720,42,39800),(3,'苏州','杭州','car',8,20,150,30,10800),(4,'上海','南京','car',9,15,210,34,12800),(5,'宁波','杭州','car',9,40,120,28,8900),(6,'南京','苏州','car',7,20,110,26,9800),(7,'杭州','武汉','car',8,45,420,40,24600),(8,'上海','合肥','car',7,5,220,30,12600),(9,'杭州','宁波','car',13,15,210,28,11800),(10,'苏州','武汉','car',9,45,330,34,19600),(11,'南京','杭州','car',9,30,300,34,17600),(12,'上海','杭州','car',9,35,150,24,9200);

CREATE TEMPORARY TABLE seed_days (day_offset INT PRIMARY KEY);
INSERT INTO seed_days VALUES (-1),(0),(1),(2),(3),(4);

SET @driver_rn := 0;
CREATE TEMPORARY TABLE seed_drivers AS SELECT (@driver_rn := @driver_rn + 1) rn, id FROM users WHERE role='driver' AND status='active' ORDER BY id;
SET @driver_count := (SELECT COUNT(*) FROM seed_drivers);

INSERT INTO trips (driver_id, vehicle_type, start_city, end_city, departure_time, arrival_time, seat_total, seat_available, price_cent, status, created_at, updated_at)
SELECT d.id, r.vehicle_type, r.start_city, r.end_city,
  TIMESTAMP(DATE_ADD(DATE(@demo_now), INTERVAL dy.day_offset DAY), MAKETIME(r.departure_hour, r.departure_minute, 0)),
  DATE_ADD(TIMESTAMP(DATE_ADD(DATE(@demo_now), INTERVAL dy.day_offset DAY), MAKETIME(r.departure_hour, r.departure_minute, 0)), INTERVAL r.duration_minute MINUTE),
  r.seat_total, r.seat_total, r.price_cent, IF(dy.day_offset < 0, 'closed', 'published'), DATE_SUB(@demo_now, INTERVAL 7 DAY), @demo_now
FROM seed_routes r JOIN seed_days dy JOIN seed_drivers d ON d.rn = MOD(r.route_no + dy.day_offset + @driver_count * 10, @driver_count) + 1;

INSERT INTO trip_stops (trip_id, stop_order, stop_name, plan_arrival_time, plan_departure_time, created_at, updated_at)
SELECT id, 1, CONCAT(start_city, '中途站'), DATE_ADD(departure_time, INTERVAL 55 MINUTE), DATE_ADD(departure_time, INTERVAL 65 MINUTE), created_at, updated_at FROM trips;

SET @passenger_rn := 0;
CREATE TEMPORARY TABLE seed_passengers AS SELECT (@passenger_rn := @passenger_rn + 1) rn, id FROM users WHERE role='passenger' AND status='active' ORDER BY id;
SET @trip_rn := 0;
CREATE TEMPORARY TABLE seed_trip_ids AS SELECT (@trip_rn := @trip_rn + 1) rn, id, price_cent, departure_time FROM trips ORDER BY id;
SET @passenger_count := (SELECT COUNT(*) FROM seed_passengers);
SET @trip_count := (SELECT COUNT(*) FROM seed_trip_ids);

INSERT INTO orders (order_no, user_id, trip_id, ticket_count, seat_type, amount, pay_status, order_status, refund_status, refund_review_note, refund_reviewed_at, payment_expire_at, created_at, updated_at)
SELECT CONCAT('ORD202605', LPAD(s.n, 6, '0')), p.id, t.id, IF(MOD(s.n,11)=0,2,1), 'standard', t.price_cent * IF(MOD(s.n,11)=0,2,1),
  IF(MOD(s.n,10) IN (1,2,3,4,5,6,7), 'paid', 'unpaid'),
  CASE WHEN MOD(s.n,17)=0 THEN 'cancelled' WHEN t.departure_time < @demo_now AND MOD(s.n,5) IN (0,1) THEN 'completed' WHEN MOD(s.n,10) IN (8,9) THEN 'pending_payment' ELSE 'pending_verification' END,
  CASE WHEN MOD(s.n,23)=0 THEN 'requested' WHEN MOD(s.n,29)=0 THEN 'rejected' WHEN MOD(s.n,31)=0 THEN 'refunded' ELSE 'none' END,
  CASE WHEN MOD(s.n,23)=0 THEN '乘客申请退款，等待人工审核' WHEN MOD(s.n,29)=0 THEN '票已使用，不符合退款规则' WHEN MOD(s.n,31)=0 THEN '审核通过，已原路退款' ELSE '' END,
  CASE WHEN MOD(s.n,29)=0 OR MOD(s.n,31)=0 THEN DATE_SUB(@demo_now, INTERVAL MOD(s.n,72) HOUR) ELSE NULL END,
  CASE WHEN MOD(s.n,10) IN (8,9) THEN DATE_ADD(DATE_SUB(@demo_now, INTERVAL MOD(s.n,36) HOUR), INTERVAL 15 MINUTE) ELSE NULL END,
  DATE_SUB(@demo_now, INTERVAL MOD(s.n,360) HOUR), DATE_SUB(@demo_now, INTERVAL MOD(s.n,180) HOUR)
FROM seed_numbers s JOIN seed_passengers p ON p.rn = MOD(s.n - 1, @passenger_count) + 1 JOIN seed_trip_ids t ON t.rn = MOD(s.n * 7 - 1, @trip_count) + 1 WHERE s.n <= 220;

UPDATE trips t LEFT JOIN (SELECT trip_id, SUM(IF(order_status <> 'cancelled', ticket_count, 0)) sold_count FROM orders GROUP BY trip_id) s ON s.trip_id=t.id SET t.seat_available=GREATEST(t.seat_total-IFNULL(s.sold_count,0),0), t.updated_at=@demo_now;

INSERT INTO payments (payment_no, order_id, user_id, amount, channel, status, paid_at, created_at, updated_at)
SELECT CONCAT('PAY202605', LPAD(id,6,'0')), id, user_id, amount, 'mock', CASE WHEN pay_status='paid' THEN 'paid' WHEN order_status='cancelled' THEN 'closed' ELSE 'pending' END, IF(pay_status='paid', DATE_ADD(created_at, INTERVAL 5 MINUTE), NULL), created_at, updated_at FROM orders;

INSERT INTO electronic_tickets (order_id, user_id, trip_id, token, token_hash, status, expires_at, created_at, updated_at)
SELECT id, user_id, trip_id, CONCAT('demo-ticket-', id), SHA2(CONCAT('demo-ticket-', id), 256), 'issued', DATE_ADD(@demo_now, INTERVAL 30 DAY), created_at, updated_at FROM orders WHERE pay_status='paid' AND order_status <> 'cancelled';

INSERT INTO notifications (user_id, type, title, content, related_order_id, is_read, read_at, created_at, updated_at)
SELECT user_id, CASE WHEN refund_status='refunded' THEN 'refund_approved' WHEN refund_status='rejected' THEN 'refund_rejected' ELSE 'order_expired' END,
  CASE WHEN refund_status='refunded' THEN '退款成功' WHEN refund_status='rejected' THEN '退款驳回' ELSE '订单提醒' END,
  CASE WHEN refund_status='refunded' THEN CONCAT('订单 ', order_no, ' 已完成退款，请留意账户到账。') WHEN refund_status='rejected' THEN CONCAT('订单 ', order_no, ' 退款申请未通过，请查看审核备注。') ELSE CONCAT('订单 ', order_no, ' 已取消或支付超时。') END,
  id, IF(MOD(id,3)=0,1,0), IF(MOD(id,3)=0,DATE_ADD(updated_at, INTERVAL 2 HOUR),NULL), DATE_ADD(updated_at, INTERVAL 1 HOUR), DATE_ADD(updated_at, INTERVAL 1 HOUR)
FROM orders WHERE refund_status IN ('refunded','rejected') OR (order_status='cancelled' AND pay_status='unpaid');

INSERT INTO price_alerts (user_id, start_city, end_city, target_price_cent, start_date, end_date, status, created_at, updated_at)
SELECT id, '杭州', '上海', 9000, DATE(@demo_now), DATE_ADD(DATE(@demo_now), INTERVAL 7 DAY), IF(MOD(id,3)=0,'triggered','active'), created_at, updated_at FROM users WHERE role='passenger' AND id <= 40;

INSERT INTO driver_settlements (driver_id, settlement_date, gross_amount_cent, refund_amount_cent, service_fee_cent, net_amount_cent, status, created_at, updated_at)
SELECT id, DATE_SUB(DATE(@demo_now), INTERVAL MOD(id,7) DAY), 200000 + id * 1000, MOD(id,5)*3000, 12000, 188000 + id * 800, IF(MOD(id,3)=0,'settled','pending'), @demo_now, @demo_now FROM users WHERE role='driver' AND status='active';

INSERT INTO token_usages (user_id, role, feature, request_kind, provider, model, prompt_tokens, completion_tokens, total_tokens, request_count, created_at)
SELECT u.id, u.role, CASE MOD(s.n,4) WHEN 0 THEN 'passenger_ai' WHEN 1 THEN 'driver_ai_draft' WHEN 2 THEN 'knowledge_search' ELSE 'knowledge_ingest' END, CASE MOD(s.n,3) WHEN 0 THEN 'chat' WHEN 1 THEN 'embedding' ELSE 'rerank' END, 'openai-compatible', CASE MOD(s.n,3) WHEN 0 THEN 'qwen-plus' WHEN 1 THEN 'text-embedding-v3' ELSE 'rerank-v1' END, 200+MOD(s.n*13,1200), IF(MOD(s.n,3)=0,80+MOD(s.n*7,900),0), 200+MOD(s.n*13,1200)+IF(MOD(s.n,3)=0,80+MOD(s.n*7,900),0), 1, DATE_SUB(@demo_now, INTERVAL MOD(s.n,720) HOUR)
FROM seed_numbers s JOIN users u ON u.id = MOD(s.n * 5, (SELECT COUNT(*) FROM users)) + 1 WHERE s.n <= 300;

INSERT INTO risk_events (severity, event_type, subject_type, subject_id, fingerprint, title, detail, status, metrics_json, created_at, updated_at) VALUES
('high','ai_rate_limit','user','12','ai_rate_limit:passenger-chat:12','AI 调用频率过高','用户 12 在 passenger-chat 场景下触发 AI 限流阈值。','open','{"scope":"passenger-chat","limit":12}',DATE_SUB(@demo_now, INTERVAL 2 HOUR),DATE_SUB(@demo_now, INTERVAL 2 HOUR)),
('medium','token_spike','user','45','token_spike:45:202605201130','Token 使用量异常','用户 45 近 15 分钟 token 使用量明显升高。','open','{"currentTotalTokens":15680,"requestCount":18}',DATE_SUB(@demo_now, INTERVAL 40 MINUTE),DATE_SUB(@demo_now, INTERVAL 40 MINUTE)),
('low','token_spike','user','1','token_spike:1:202605201000','管理端导入提醒','管理员导入种子数据后触发低风险审计提醒。','acknowledged','{"currentTotalTokens":4200}',DATE_SUB(@demo_now, INTERVAL 3 HOUR),DATE_SUB(@demo_now, INTERVAL 2 HOUR));

INSERT INTO refund_audit_logs (order_id, refund_status, review_note, reviewer_id, created_at)
SELECT id, refund_status, refund_review_note, 1, refund_reviewed_at FROM orders WHERE refund_status IN ('refunded','rejected') AND refund_reviewed_at IS NOT NULL;

INSERT INTO audit_logs (actor_user_id, actor_role, action, resource_type, resource_id, detail, created_at)
VALUES (1,'admin','seed_database','database','ridehailing_demo','重建数据库并导入演示数据',@demo_now);

SET FOREIGN_KEY_CHECKS = 1;

SELECT COUNT(*) AS user_count FROM users;
SELECT COUNT(*) AS trip_count FROM trips;
SELECT COUNT(*) AS order_count FROM orders;
SELECT COUNT(*) AS notification_count FROM notifications;
SELECT COUNT(*) AS token_usage_count FROM token_usages;
SELECT COUNT(*) AS risk_event_count FROM risk_events;
