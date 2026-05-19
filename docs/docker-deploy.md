# Docker 打包部署

本项目使用 Docker Compose 一次启动 4 个服务：

- `frontend`：Nginx 静态站点，对外暴露 `http://localhost:8080`
- `backend`：Go API 服务，容器内监听 `:8081`
- `mysql`：MySQL 8.4，首次启动自动执行 `backend/sql/ridehailing_demo_seed.sql`
- `redis`：验证码和 AI 限流缓存

## 1. 准备环境变量

复制示例环境变量文件：

```powershell
Copy-Item .env.example .env
```

如果要使用 AI、Embedding、Rerank 能力，编辑 `.env` 填入对应 API Key。只体验基础流程时可以先留空。

## 2. 构建并启动

在项目根目录执行：

```powershell
docker compose up -d --build
```

Linux 服务器使用 Podman Compose 时建议执行：

```bash
export GOPROXY=https://goproxy.cn,direct
podman-compose down
podman-compose up -d --build
```

如果上一次构建在 `go mod download` 阶段超时，后端镜像不会生成，Podman 可能会继续尝试从远程拉取 `ridehailingsystem_backend`。本项目已在 `docker-compose.yml` 中显式使用 `localhost/ridehailingsystem_backend:latest` 和 `localhost/ridehailingsystem_frontend:latest`，重新构建前先 `podman-compose down` 可以清理失败状态。

Podman Compose 1.0.x 对嵌套 `depends_on` 会生成 `--requires` 依赖图，某些版本在重建时可能报 `depends on container ... not found in input list`。前端只是 Nginx 静态服务，请求后端时由 Nginx 通过服务名代理，因此前端不需要容器级 `depends_on`。

启动完成后访问：

```text
http://localhost:8080
```

前端请求 `/api` 会由 Nginx 自动转发到后端容器，不需要在浏览器里访问后端端口。

## 3. 查看日志

```powershell
docker compose logs -f backend
docker compose logs -f frontend
docker compose logs -f mysql
```

## 4. 停止服务

```powershell
docker compose down
```

如果想连数据库数据一起清空，执行：

```powershell
docker compose down -v
```

注意：`-v` 会删除 MySQL、Redis、知识库 memory 等 Docker 卷数据。

## 5. 重新初始化数据库

MySQL 初始化脚本只会在 `mysql-data` 卷第一次创建时执行。如果改了 `backend/sql/ridehailing_demo_seed.sql`，需要清空卷后重启：

```powershell
docker compose down -v
docker compose up -d --build
```

## 6. 默认账号

种子 SQL 中的默认密码为：

```text
123456
```

管理员账号：

```text
phone: 18800000000
email: admin_root@tripverse.local
```

## 7. 生产部署建议

上线前至少修改：

- `docker-compose.yml` 中的 `JWT_SECRET`
- MySQL root 密码和 `MYSQL_DSN`
- 是否开启 `ENABLE_MOCK_PAYMENT`
- `.env` 中的真实模型服务 API Key

生产环境可以把 `frontend` 的 `8080:80` 改成由服务器 Nginx、Caddy 或云负载均衡统一代理。
