# Docker / Podman 部署说明

本项目包含 4 个服务：

- `frontend`：Nginx 静态前端，对外端口 `8080`
- `backend`：Go API 服务，容器内监听 `8081`
- `mysql`：MySQL 8.4，业务库 `ridehailing_demo`
- `redis`：验证码、限流和临时缓存

## 1. 端口与前端接口

前端浏览器只需要访问：

```text
http://服务器公网IP:8080
```

前端请求统一走 `/api`，例如：

```text
http://服务器公网IP:8080/api/tickets/search
```

`frontend/nginx.conf` 已把 `/api` 转发到容器网络里的后端：

```nginx
location /api/ {
    proxy_pass http://backend:8081/api/;
}
```

所以服务器安全组通常只需要放行 `8080`。`backend:8081` 是容器内部端口，不需要浏览器直接访问。

## 2. 构建与启动

Linux 服务器使用 Podman Compose：

```bash
cd /root/railway/RideHailingSystem
export GOPROXY=https://goproxy.cn,direct
podman-compose up -d --build
```

如果之前 Podman 留下旧依赖容器，先清理容器再启动，数据卷不会被删除：

```bash
podman rm -f ridehailingsystem_frontend_1 ridehailingsystem_backend_1 ridehailingsystem_redis_1 ridehailingsystem_mysql_1
podman-compose up -d --build
```

查看运行状态：

```bash
podman ps
podman logs -f ridehailingsystem_backend_1
```

## 3. 重新创建数据库并导入数据

`backend/sql/ridehailing_demo_seed.sql` 已包含完整建表和大量演示数据，会 `DROP DATABASE` 后重新创建 `ridehailing_demo`。

在服务器项目目录执行：

```bash
podman exec -i ridehailingsystem_mysql_1 mysql -uroot -phyh --default-character-set=utf8mb4 < backend/sql/ridehailing_demo_seed.sql
```

如果 MySQL 容器还没有启动：

```bash
podman-compose up -d mysql
podman exec -i ridehailingsystem_mysql_1 mysql -uroot -phyh --default-character-set=utf8mb4 < backend/sql/ridehailing_demo_seed.sql
```

如果想让 MySQL 初始化脚本从零自动执行，也可以删除数据卷后重建：

```bash
podman-compose down
podman volume rm ridehailingsystem_mysql-data
podman-compose up -d --build
```

注意：删除 `ridehailingsystem_mysql-data` 会清空现有数据库。

## 4. 已添加的演示数据

种子数据包含：

- 管理员：1 个
- 乘客：120 个
- 司机：24 个
- 车辆、司机资质、乘客实名信息
- 12 条城际线路，按多天生成班次
- 220 笔订单，覆盖待支付、待核销、已完成、取消、退款中、已退款、退款驳回
- 支付记录、电子票、通知、价格提醒
- Token 使用量、风险事件、审计日志

默认密码都是：

```text
123456
```

常用测试账号：

```text
管理员：18800000000 / 123456
乘客：13000000001 / 123456
司机：15000000001 / 123456
```

登录页可以选择乘客、司机、管理员登录；注册页只允许乘客和司机注册。

## 5. 如何查看数据

### 方式一：页面查看

访问：

```text
http://服务器公网IP:8080
```

管理员登录后可以查看：

- 用户管理
- 订单与退款
- Token 使用量
- 风险事件
- 知识库

乘客登录后可以查看：

- 搜索车票
- 下单和支付
- 订单详情
- 电子票与通知

司机登录后可以查看：

- 司机工作台
- 班次列表
- 班次订单核销
- 收入统计

### 方式二：进入 MySQL 查看

```bash
podman exec -it ridehailingsystem_mysql_1 mysql -uroot -phyh --default-character-set=utf8mb4 ridehailing_demo
```

常用 SQL：

```sql
SELECT COUNT(*) FROM users;
SELECT id, phone, nickname, role, status FROM users LIMIT 10;
SELECT id, start_city, end_city, departure_time, seat_available FROM trips LIMIT 10;
SELECT id, order_no, pay_status, order_status, refund_status FROM orders LIMIT 10;
SELECT id, title, status, created_at FROM risk_events ORDER BY id DESC LIMIT 10;
```

### 方式三：通过接口查看

无需登录的票务搜索接口：

```bash
curl "http://127.0.0.1:8080/api/tickets/search?startCity=杭州&endCity=上海"
```

需要登录的接口要先在页面登录，或使用登录接口拿到 token 后再访问。

## 6. 常见问题

如果页面还是旧内容，通常是浏览器缓存：

```bash
podman-compose up -d --build
```

然后浏览器强制刷新，或打开无痕窗口访问。

如果 Go 依赖下载超时：

```bash
export GOPROXY=https://goproxy.cn,direct
podman-compose up -d --build
```

如果数据库中文显示问号，确认导入命令带了：

```bash
--default-character-set=utf8mb4
```

本项目的建库 SQL 已使用：

```sql
CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci
```
