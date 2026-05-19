# RideHailingSystem 功能扩展与技术实现路线图

## 项目目标

当前项目已经具备乘客端、司机端、管理端和 AI 助手的页面基础。下一步目标是把它扩展成一个完整的“AI 出行票务平台”：

- 乘客可以搜索班次、下单、支付、查看电子票、申请退款。
- 司机可以发布班次、管理订单、核验乘客、查看收益。
- 管理员可以管理用户、订单、退款、风控、知识库和模型调用。
- AI 可以辅助购票、解释规则、生成司机发车建议和管理日报。

## 总体技术栈建议

### 前端

- 框架：Vue 3 + TypeScript + Vite
- UI 组件库：Element Plus
- 状态管理：Pinia
- 请求层：Axios + 统一 request 封装
- 数据缓存：TanStack Query for Vue
- 图表：ECharts
- 路由：Vue Router
- 表单校验：VeeValidate 或 Element Plus Form Rules
- 二维码：qrcode
- 扫码：html5-qrcode
- E2E 测试：Playwright

### 后端

- 语言：Go
- Web 框架：Gin
- ORM：GORM 或 sqlc
- 数据库：PostgreSQL
- 缓存：Redis
- 鉴权：JWT
- 权限：RBAC
- 异步任务：Asynq + Redis
- 日志：Zap
- 配置：Viper
- API 文档：Swagger / OpenAPI
- 单元测试：Go testing + testify

### AI 与搜索

- 大模型接口：OpenAI API 或兼容 OpenAI 协议的模型服务
- RAG 向量库：pgvector，后续可升级 Milvus
- 文档切分：按标题、段落、规则条款切分
- Embedding：text-embedding-3-small 或本地 embedding 模型
- 搜索增强：PostgreSQL 全文搜索 + 向量检索混合召回

### 部署与基础设施

- 本地编排：Docker Compose
- 反向代理：Nginx
- 数据库：PostgreSQL 容器
- 缓存：Redis 容器
- 文件存储：MinIO
- CI：GitHub Actions
- 监控：Prometheus + Grafana

## P0：乘客搜索、下单、支付闭环

### 功能

- [ ] 乘客搜索班次
- [ ] 查看班次详情
- [ ] 创建订单
- [ ] 创建支付单
- [ ] 模拟支付成功
- [ ] 查看订单列表和订单详情

### 使用技术

- 前端：Vue 3、Element Plus、Pinia、TanStack Query
- 后端：Go + Gin
- 数据库：PostgreSQL
- 缓存：Redis
- 支付：先实现 Mock Payment，后续接入微信支付或支付宝

### 关键数据表

- `trips`：班次表，保存起点、终点、发车时间、到达时间、票价、座位数、状态。
- `orders`：订单表，保存用户、班次、金额、订单状态、支付状态。
- `payments`：支付表，保存订单号、支付金额、支付渠道、支付状态。
- `passengers`：常用乘车人表。

### 实现逻辑

1. 乘客在搜索页输入起点、终点、日期。
2. 前端调用 `GET /api/tickets/search`。
3. 后端查询 `trips` 表，按日期、城市、状态过滤可售班次。
4. 用户进入详情页后，前端调用 `GET /api/tickets/:ticketId`。
5. 用户提交订单时，后端检查余票并创建 `orders` 记录。
6. 创建订单后锁定库存，订单状态为 `pending_payment`。
7. 支付页调用 `POST /api/payments/create` 创建支付单。
8. Mock 支付成功后更新 `payments.status = success`，同时更新 `orders.status = paid`。
9. 订单列表从 `GET /api/orders/my` 获取当前用户订单。

## P0：电子票二维码与司机核验

### 功能

- [ ] 支付成功后生成电子票
- [ ] 订单详情页展示二维码
- [ ] 司机端扫码核验
- [ ] 核验成功后更新订单状态
- [ ] 防止重复核验

### 使用技术

- 前端二维码生成：qrcode
- 前端扫码：html5-qrcode
- 后端签名：HMAC-SHA256
- 缓存：Redis 存储短期核验 token

### 关键数据表

- `tickets`：电子票表，保存订单 ID、二维码 token、核验状态。
- `ticket_verifications`：核验记录表，保存司机、订单、核验时间、核验结果。

### 实现逻辑

1. 订单支付成功后，后端生成电子票 token。
2. token 内容包含 `orderId`、`userId`、`tripId`、过期时间和签名。
3. 前端订单详情页使用 qrcode 将 token 渲染为二维码。
4. 司机端扫码后调用 `POST /api/driver/tickets/verify`。
5. 后端校验签名、过期时间、订单归属和班次归属。
6. 如果未核验，更新订单状态为 `checked_in`。
7. 写入 `ticket_verifications` 核验记录。
8. 如果重复扫码，返回“已核验”状态，不重复更新。

## P0：退款申请与管理端审核

### 功能

- [ ] 乘客申请退款
- [ ] 后端根据规则判断是否可退款
- [ ] 管理端查看待审核退款
- [ ] 管理员通过或拒绝退款
- [ ] 订单详情展示退款进度

### 使用技术

- 后端状态机：订单状态枚举 + 明确状态流转函数
- 异步任务：Asynq 处理退款通知
- 数据库事务：PostgreSQL transaction
- 管理端表格：Element Plus Table

### 关键数据表

- `refunds`：退款申请表，保存订单、金额、原因、审核状态。
- `refund_audit_logs`：退款审核日志表。
- `notifications`：用户通知表。

### 实现逻辑

1. 乘客在订单详情页点击“申请退款”。
2. 后端检查订单是否已支付、是否已核验、是否超过发车时间。
3. 符合规则则创建 `refunds`，状态为 `pending`。
4. 管理端调用 `GET /api/admin/orders?refundStatus=pending` 获取待审核列表。
5. 管理员点击通过，后端开启事务：
   - 更新 `refunds.status = approved`
   - 更新 `orders.status = refunded`
   - 释放或回滚相关库存
   - 写入审核日志
6. 管理员点击拒绝，则更新 `refunds.status = rejected` 并写入拒绝原因。
7. 通过 Asynq 创建通知任务，给乘客发送退款结果通知。

## P1：乘客端体验增强

### 功能

- [ ] 常用乘车人管理
- [ ] 低价提醒订阅
- [ ] 余票紧张提醒
- [ ] 出发前提醒
- [ ] 退改签规则提示

### 使用技术

- 前端：Vue 3、Element Plus Dialog、Element Plus Form
- 后端：Gin + PostgreSQL
- 定时任务：Asynq Scheduler 或 cron
- 通知：站内信，后续可扩展短信、邮件、微信模板消息

### 关键数据表

- `passengers`：常用乘车人。
- `price_alerts`：低价提醒订阅。
- `notifications`：通知。

### 实现逻辑

1. 用户在个人中心维护常用乘车人。
2. 下单页选择乘车人后自动填充姓名、证件号和手机号。
3. 用户订阅低价提醒时，保存路线、日期范围和目标价格。
4. 定时任务扫描可售班次，如果价格低于目标价格则生成通知。
5. 发车前 30 分钟扫描已支付订单，生成出发提醒通知。

## P1：司机端运营增强

### 功能

- [ ] 司机资质认证
- [ ] 车辆信息管理
- [ ] 班次修改与取消
- [ ] 司机收益结算
- [ ] 路线收入排行
- [ ] 低上座率提醒

### 使用技术

- 前端：Element Plus Upload、Table、Statistic、ECharts
- 后端：Gin + PostgreSQL
- 文件存储：MinIO
- 异步任务：Asynq
- 图表：ECharts

### 关键数据表

- `driver_profiles`：司机认证资料。
- `vehicles`：车辆信息。
- `driver_settlements`：司机结算记录。
- `trip_statistics`：班次统计。

### 实现逻辑

1. 司机上传身份证、驾驶证、车辆信息。
2. 文件存入 MinIO，数据库只保存文件 URL 和审核状态。
3. 管理员审核通过后，司机才可以发布班次。
4. 班次结束后，系统统计订单收入、退款金额、平台服务费。
5. 结算任务按天生成 `driver_settlements`。
6. ECharts 展示收入趋势、路线排行和上座率。
7. 上座率低于阈值时，生成运营建议。

## P1：管理端与权限系统

### 功能

- [ ] RBAC 权限控制
- [ ] 用户角色管理
- [ ] 操作审计日志
- [ ] 风控事件详情
- [ ] 管理端数据看板

### 使用技术

- 后端鉴权：JWT
- 权限模型：RBAC
- 中间件：Gin middleware
- 图表：ECharts
- 审计日志：PostgreSQL + Zap

### 关键数据表

- `roles`：角色表。
- `permissions`：权限表。
- `role_permissions`：角色权限关系表。
- `user_roles`：用户角色关系表。
- `audit_logs`：审计日志。
- `risk_events`：风控事件。

### 实现逻辑

1. 登录后后端签发 JWT，包含用户 ID 和当前角色。
2. 请求进入管理端接口时，Gin 中间件解析 JWT。
3. 中间件检查角色是否拥有当前接口权限。
4. 管理员修改用户、订单、退款、权限时，统一写入 `audit_logs`。
5. 数据看板接口聚合订单量、GMV、退款率、活跃用户等指标。
6. 风控事件由登录、支付、退款等业务流程写入。

## P2：AI 知识库客服助手 RAG

### 功能

- [ ] 管理端上传知识库文档
- [ ] 文档自动切分
- [ ] 生成 embedding
- [ ] 用户提问时向量召回相关规则
- [ ] AI 基于召回内容回答
- [ ] 返回引用来源

### 使用技术

- AI：OpenAI API 或兼容 OpenAI 协议服务
- 向量库：PostgreSQL + pgvector
- 文档处理：Go 后端切分 Markdown / TXT
- 后台任务：Asynq
- 前端：AI Chat 页面 + 管理端知识库页面

### 关键数据表

- `knowledge_documents`：知识库文档。
- `knowledge_chunks`：文档切片。
- `ai_chat_sessions`：AI 会话。
- `ai_chat_messages`：AI 消息记录。
- `ai_usage_logs`：AI 调用日志。

### 实现逻辑

1. 管理员上传退改签规则、乘车规则、客服 FAQ。
2. 后端将文档按标题和段落切分成 chunk。
3. 异步任务为每个 chunk 生成 embedding。
4. 用户在 AI 助手提问时，后端先生成 query embedding。
5. 使用 pgvector 查询最相关的 chunks。
6. 将召回内容和用户问题拼成 prompt。
7. 调用大模型生成回答。
8. 保存会话消息和 token 用量。
9. 前端展示回答和引用来源。

## P2：自然语言购票

### 功能

- [ ] 用户输入自然语言购票需求
- [ ] AI 提取起点、终点、日期、时间、预算、偏好
- [ ] 自动跳转搜索结果页
- [ ] 支持多轮补全缺失信息

### 使用技术

- AI 函数调用 / JSON Schema 输出
- 前端：Vue Router + Pinia 保存搜索条件
- 后端：AI intent parser service

### 实现逻辑

1. 用户输入“明天早上杭州到苏州，帮我找便宜点的车”。
2. 后端调用 AI，要求输出结构化 JSON：
   - `startCity`
   - `endCity`
   - `date`
   - `timeRange`
   - `sortBy`
   - `allowTransfer`
3. 如果缺少必要字段，AI 返回追问问题。
4. 字段完整后，前端将条件带到搜索页。
5. 搜索页调用现有班次搜索接口。

## P2：司机 AI 发车建议

### 功能

- [ ] 根据历史订单生成发车建议
- [ ] 推荐路线、时间、票价和座位数
- [ ] 一键填充发布班次表单

### 使用技术

- 数据分析：PostgreSQL 聚合查询
- 图表：ECharts
- AI：大模型生成运营建议
- 缓存：Redis 缓存热门路线统计

### 实现逻辑

1. 后端统计最近 7 天和 30 天的路线订单量、上座率、退款率。
2. 根据时间段聚合热门出发时间。
3. 将统计结果传给 AI，生成可读建议。
4. AI 返回结构化发车草稿。
5. 前端点击“应用到表单”，自动填充发布班次页面。

## P2：数据看板与风控

### 功能

- [ ] 管理端运营数据看板
- [ ] 热门路线排行
- [ ] 退款原因分布
- [ ] 异常登录检测
- [ ] 恶意退款检测
- [ ] 高频下单检测

### 使用技术

- 图表：ECharts
- 聚合查询：PostgreSQL materialized view
- 缓存：Redis
- 风控规则：规则引擎函数，后续可引入机器学习模型

### 实现逻辑

1. 订单、支付、退款、登录等流程写入业务事件。
2. 定时任务聚合每日指标。
3. 管理端看板读取聚合结果，避免实时扫大表。
4. 风控规则检查同一用户、同一设备、同一 IP 的异常行为。
5. 命中规则后写入 `risk_events`。
6. 管理员可以标记风险事件为已处理或误报。

## 推荐开发顺序

1. 完成订单支付闭环
2. 增加电子票二维码和司机核验
3. 增加退款申请与管理端审核
4. 增加 RBAC 权限和审计日志
5. 增加 AI 知识库客服助手
6. 增加司机收益和发车建议
7. 增加管理端数据看板和风控规则
8. 前端迁移到 Vue 3 + TypeScript + Vite

## 最小可演示版本范围

- [ ] 乘客可以搜索班次
- [ ] 乘客可以创建订单
- [ ] 乘客可以模拟支付成功
- [ ] 乘客可以查看电子票二维码
- [ ] 司机可以核验电子票
- [ ] 乘客可以申请退款
- [ ] 管理员可以审核退款
- [ ] AI 可以基于知识库回答退改签问题

## 当前项目改造建议

- 保留当前静态 HTML 作为原型页面。
- 后续新前端建议创建 `frontend-vue` 目录，避免直接破坏现有页面。
- 后端优先补齐订单、支付、退款、电子票相关接口。
- 数据库优先从 PostgreSQL 开始，避免业务状态长期只放内存或 Mock。
- AI 功能先做知识库问答，不急着做复杂预测。
- 管理端先做退款审核和用户管理，再做复杂风控。
