# RideHailingSystem 功能模块实现说明

## 1. 项目整体架构

### 实现思想

项目采用“静态前端 + Go 后端 API + MySQL + Redis”的分层结构。前端页面负责展示和交互，后端负责业务规则、数据状态流转和权限校验，数据库负责持久化，Redis 负责验证码、AI 限流等短期状态。

### 技术栈

- 前端：HTML、CSS、原生 JavaScript
- 后端：Go、net/http、GORM
- 数据库：MySQL
- 缓存：Redis
- 鉴权：JWT
- AI：兼容 OpenAI 协议的模型接口
- 部署：Docker、Docker Compose、Nginx

### 核心代码

- 后端入口：`backend/cmd/server/main.go`
- 路由注册：`backend/internal/router/router.go`
- 配置加载：`backend/internal/config/config.go`
- 前端入口脚本：`frontend/assets/app.js`
- 前端全局样式：`frontend/assets/styles.css`

### 模块关联

`main.go` 初始化数据库、Redis、Repository、Service、Handler，再把 Handler 注入 `router.NewRouter`。请求进入后由 Router 分发到 Handler，Handler 调用 Service，Service 处理业务逻辑并调用 Repository 读写数据库。

## 2. 用户认证与角色系统

### 实现步骤

1. 用户通过手机号密码或邮箱验证码登录。
2. 后端校验用户信息。
3. 登录成功后生成 JWT。
4. 前端保存 token 和用户角色。
5. 后续请求通过 `Authorization: Bearer <token>` 访问受保护接口。
6. 后端中间件解析 JWT，把用户信息写入 request context。

### 核心代码

- 用户模型：`backend/internal/model/user.go`
- 认证服务：`backend/internal/service/auth_service.go`
- 认证接口：`backend/internal/handler/auth_handler.go`
- JWT 工具：`backend/internal/pkg/jwtutil/jwt.go`
- 鉴权中间件：`backend/internal/pkg/middleware/auth.go`
- 前端登录逻辑：`frontend/assets/app.js`

### 模块关联

认证模块为乘客端、司机端、管理端提供统一身份基础。订单、支付、司机班次、管理端审核等接口都会先经过 `AuthMiddleware`。

### 实现思想

认证系统使用 JWT 作为无状态登录凭证，后端通过 token 中的 `userId` 和 `role` 判断当前用户是谁、能访问哪些功能。

## 3. 乘客票务搜索模块

### 实现步骤

1. 乘客在搜索页填写起点、终点、日期。
2. 前端调用 `GET /api/tickets/search`。
3. 后端根据城市、日期、班次状态和余票查询可售班次。
4. 返回直达班次和可扩展的中转方案。
5. 前端渲染搜索结果卡片。
6. 用户点击结果进入班次详情。

### 核心代码

- 班次模型：`backend/internal/model/trip.go`
- 停靠点模型：`backend/internal/model/trip_stop.go`
- 票务服务：`backend/internal/service/ticket_service.go`
- 票务接口：`backend/internal/handler/ticket_handler.go`
- 班次仓储：`backend/internal/repository/trip_repository.go`
- 搜索页面：`frontend/search.html`
- 前端搜索逻辑：`frontend/assets/app.js`

### 模块关联

票务搜索依赖司机端发布的 `trips` 数据。搜索结果中的班次 ID 会继续传给订单模块创建订单。

### 实现思想

搜索模块只展示可售班次，即状态为 `published` 且 `seat_available > 0` 的班次。这样可以避免乘客看到无法购买的结果。

## 4. 订单创建与库存锁定模块

### 实现步骤

1. 乘客在班次详情页进入下单确认页。
2. 前端提交班次 ID、票数、座席类型。
3. 后端校验当前用户必须是乘客。
4. 后端检查班次是否存在、是否可售、余票是否足够。
5. 在数据库事务中扣减余票。
6. 创建订单，状态为 `pending_payment`。
7. 设置支付超时时间。

### 核心代码

- 订单模型：`backend/internal/model/order.go`
- 订单服务：`backend/internal/service/order_service.go`
- 订单接口：`backend/internal/handler/order_handler.go`
- 订单仓储：`backend/internal/repository/order_repository.go`
- 下单页面：`frontend/checkout.html`

### 模块关联

订单模块关联票务班次、乘客用户、支付模块和通知模块。订单创建后进入支付流程，支付成功后进入电子票流程。

### 实现思想

创建订单时使用事务和行锁扣减座位，避免多人同时购买导致超卖。订单创建后并不立即完成，而是进入待支付状态。

## 5. 支付与模拟支付模块

### 实现步骤

1. 乘客进入支付页。
2. 前端调用 `POST /api/payments/create` 创建支付单。
3. 后端检查订单归属和订单状态。
4. 创建支付记录，状态为 `pending`。
5. 本地演示环境调用 Mock 支付成功接口。
6. 后端更新支付状态为 `paid`。
7. 同步更新订单状态为 `pending_verification`。
8. 支付成功后自动生成电子票。

### 核心代码

- 支付模型：`backend/internal/model/payment.go`
- 支付服务：`backend/internal/service/payment_service.go`
- 支付接口：`backend/internal/handler/payment_handler.go`
- 支付仓储：`backend/internal/repository/payment_repository.go`
- 支付页面：`frontend/payment.html`

### 模块关联

支付模块依赖订单模块。支付成功后会调用电子票服务生成电子票，并把订单推进到待核验状态。

### 实现思想

支付模块先实现 Mock Payment，保证业务闭环可演示。真实支付平台可以在 `PaymentService` 中扩展不同 channel。

## 6. 电子票与司机核验模块

### 实现步骤

1. 订单支付成功后自动生成电子票。
2. 电子票包含订单 ID、用户 ID、班次 ID、过期时间、随机 nonce 和 HMAC 签名。
3. 乘客在订单详情页点击“查看电子票”。
4. 前端调用 `GET /api/orders/{orderId}/ticket`。
5. 后端校验订单归属，返回电子票 token。
6. 司机在班次详情页粘贴电子票 token。
7. 前端调用 `POST /api/driver/tickets/verify`。
8. 后端校验签名、过期时间、班次归属、支付状态和退款状态。
9. 核验成功后更新电子票状态为 `verified`。
10. 同步更新订单状态为 `completed`。
11. 写入核验记录。

### 核心代码

- 电子票模型：`backend/internal/model/electronic_ticket.go`
- 电子票服务：`backend/internal/service/electronic_ticket_service.go`
- 电子票接口：`backend/internal/handler/electronic_ticket_handler.go`
- 电子票仓储：`backend/internal/repository/electronic_ticket_repository.go`
- 路由注册：`backend/internal/router/router.go`
- 前端订单详情逻辑：`frontend/assets/app.js`
- 司机班次详情页：`frontend/driver-trip-detail.html`

### 模块关联

电子票模块关联订单、支付、司机班次和核验记录。支付成功负责出票，司机端负责核验，退款通过后会作废未核验电子票。

### 实现思想

电子票 token 使用 HMAC 签名，避免用户伪造票据。数据库保存 token hash，便于定位电子票，同时减少直接依赖明文 token 查询。

## 7. 订单取消与超时关闭模块

### 实现步骤

1. 乘客可取消待支付订单。
2. 后端校验订单归属。
3. 仅允许 `pending_payment` 订单取消。
4. 事务中释放座位。
5. 更新订单状态为 `cancelled`。
6. 后台定时任务每分钟扫描超时未支付订单。
7. 自动取消超时订单并释放座位。
8. 创建站内通知。

### 核心代码

- 订单服务：`backend/internal/service/order_service.go`
- 订单仓储：`backend/internal/repository/order_repository.go`
- 定时任务入口：`backend/cmd/server/main.go`
- 通知模型：`backend/internal/model/notification.go`

### 模块关联

订单取消会影响班次余票。超时关闭还会调用支付仓储关闭待支付支付单，并调用通知模块提醒用户。

### 实现思想

座位是核心库存资源，所有释放和扣减都放在事务中处理，保证订单状态和班次余票一致。

## 8. 退款申请与管理端审核模块

### 实现步骤

1. 乘客对已支付订单发起退款申请。
2. 后端校验订单状态和支付状态。
3. 将订单退款状态更新为 `requested`。
4. 管理员在管理端查看待审核退款。
5. 管理员通过退款时，更新退款状态为 `refunded`。
6. 如果订单仍处于待核验状态，释放座位并取消订单。
7. 作废未核验电子票。
8. 写入退款审核日志。
9. 写入通用审计日志。
10. 创建乘客站内通知。
11. 管理员拒绝退款时，更新状态为 `rejected` 并记录原因。

### 核心代码

- 订单模型：`backend/internal/model/order.go`
- 审计模型：`backend/internal/model/audit.go`
- 订单服务：`backend/internal/service/order_service.go`
- 订单接口：`backend/internal/handler/order_handler.go`
- 审计服务：`backend/internal/service/audit_service.go`
- 审计接口：`backend/internal/handler/audit_handler.go`
- 审计仓储：`backend/internal/repository/audit_repository.go`
- 管理端订单页面：`frontend/admin-orders.html`

### 模块关联

退款模块关联订单、电子票、班次库存、通知和审计。退款通过会影响电子票有效性和座位库存。

### 实现思想

退款是典型状态机流程，必须限制可流转状态，避免已退款订单重复退款，或已核验订单被错误释放座位。

## 9. 通知中心模块

### 实现步骤

1. 后端在订单超时、退款通过、退款拒绝等场景创建通知。
2. 乘客在个人中心获取通知列表。
3. 前端展示未读数量。
4. 用户可以标记单条通知已读。
5. 用户可以一键全部已读。

### 核心代码

- 通知模型：`backend/internal/model/notification.go`
- 通知服务：`backend/internal/service/notification_service.go`
- 通知接口：`backend/internal/handler/notification_handler.go`
- 通知仓储：`backend/internal/repository/notification_repository.go`
- 前端个人中心：`frontend/profile.html`

### 模块关联

通知模块被订单超时任务和退款审核流程调用，也可以扩展到低价提醒、余票提醒和发车提醒。

### 实现思想

通知作为独立模块，不直接参与主流程事务。主流程完成后创建通知，失败不影响核心业务成功。

## 10. 常用乘车人模块

### 实现步骤

1. 乘客新增常用乘车人。
2. 后端校验姓名、证件号、手机号。
3. 保存到 `passengers` 表。
4. 乘客可以查看自己的乘车人列表。
5. 乘客可以删除自己的乘车人。
6. 后续下单页可选择乘车人自动填充信息。

### 核心代码

- 扩展模型：`backend/internal/model/extension.go`
- 扩展服务：`backend/internal/service/extension_service.go`
- 扩展接口：`backend/internal/handler/extension_handler.go`
- 扩展仓储：`backend/internal/repository/extension_repository.go`
- 路由注册：`backend/internal/router/router.go`

### 模块关联

乘车人模块和订单模块关联。当前接口已完成，后续可把乘车人 ID 写入订单，形成更完整的实名购票流程。

### 实现思想

乘车人属于乘客个人私有数据，所有查询、删除都必须按当前登录用户过滤，不能跨用户访问。

## 11. 低价提醒模块

### 实现步骤

1. 乘客创建低价提醒。
2. 填写起点、终点、目标价格、日期范围。
3. 后端保存提醒，状态为 `active`。
4. 乘客可以查看自己的提醒列表。
5. 乘客可以禁用提醒。
6. 后续可增加定时任务扫描班次价格，命中后创建通知。

### 核心代码

- 低价提醒模型：`backend/internal/model/extension.go`
- 低价提醒服务：`backend/internal/service/extension_service.go`
- 低价提醒接口：`backend/internal/handler/extension_handler.go`
- 低价提醒仓储：`backend/internal/repository/extension_repository.go`

### 模块关联

低价提醒模块关联票务搜索和通知中心。扫描到低价班次后，可以写入 `notifications` 通知用户。

### 实现思想

低价提醒先保存用户订阅条件，再由异步任务或定时任务批量处理，避免用户请求时做重计算。

## 12. 司机班次发布与管理模块

### 实现步骤

1. 司机填写起点、终点、发车时间、到达时间、座位数、票价和停靠点。
2. 前端提交到司机班次接口。
3. 后端校验当前用户是司机。
4. 创建 `trips` 记录。
5. 创建停靠点记录。
6. 司机可以查看自己的班次列表。
7. 司机可以查看班次详情和订单列表。

### 核心代码

- 班次模型：`backend/internal/model/trip.go`
- 停靠点模型：`backend/internal/model/trip_stop.go`
- 司机班次服务：`backend/internal/service/trip_service.go`
- 司机班次接口：`backend/internal/handler/driver_trip_handler.go`
- 班次仓储：`backend/internal/repository/trip_repository.go`
- 发布页面：`frontend/driver-publish.html`
- 班次详情页：`frontend/driver-trip-detail.html`

### 模块关联

司机发布的班次会进入乘客搜索模块。班次详情关联订单列表和电子票核验。

### 实现思想

司机只能管理自己发布的班次。班次作为票务系统的核心资源，需要关联库存、订单和收益统计。

## 13. 司机资质与车辆管理模块

### 实现步骤

1. 司机提交真实姓名、身份证、驾驶证号和证件图片地址。
2. 后端保存资质信息，状态为 `pending`。
3. 司机可以查看自己的资质状态。
4. 司机新增车辆信息。
5. 后端保存车牌、品牌、型号、座位数。
6. 司机可以查看自己的车辆列表。

### 核心代码

- 司机资质模型：`backend/internal/model/extension.go`
- 扩展服务：`backend/internal/service/extension_service.go`
- 扩展接口：`backend/internal/handler/extension_handler.go`
- 扩展仓储：`backend/internal/repository/extension_repository.go`

### 模块关联

司机资质和车辆模块关联司机班次。后续可以限制只有资质审核通过且至少有一辆可用车辆的司机才能发布班次。

### 实现思想

资质审核采用状态字段控制流程，先实现资料提交和状态保存，后续管理端可增加审核入口。

## 14. 司机收益与运营数据模块

### 实现步骤

1. 后端统计司机班次。
2. 根据订单计算今日收入、已售票数、待核验数量和退款数量。
3. 前端司机工作台展示关键指标。
4. 收益页展示收入、待结算、路线排行和经营建议。

### 核心代码

- 司机班次服务：`backend/internal/service/trip_service.go`
- 司机班次接口：`backend/internal/handler/driver_trip_handler.go`
- 订单仓储：`backend/internal/repository/order_repository.go`
- 司机工作台页面：`frontend/driver-dashboard.html`
- 司机收益页面：`frontend/driver-income.html`

### 模块关联

收益数据来自班次和订单。退款状态、核验状态会影响司机实际收益。

### 实现思想

收益统计优先使用聚合查询动态计算，后续订单量变大后可以落地到 `driver_settlements` 做结算快照。

## 15. 管理端用户管理模块

### 实现步骤

1. 管理员查询用户列表。
2. 后端按角色、状态、关键词过滤用户。
3. 管理员修改用户角色或状态。
4. 后端更新用户记录。
5. 管理端展示用户汇总数据。

### 核心代码

- 用户模型：`backend/internal/model/user.go`
- 用户服务：`backend/internal/service/user_service.go`
- 用户接口：`backend/internal/handler/user_hander.go`
- 用户仓储：`backend/internal/repository/User_repository.go`
- 管理端用户页面：`frontend/admin-users.html`

### 模块关联

用户管理影响权限、角色导航、订单归属和司机功能入口。

### 实现思想

管理端用户操作必须限制为管理员角色，普通乘客和司机不能访问管理接口。

## 16. 管理端数据看板模块

### 实现步骤

1. 管理员访问管理端首页。
2. 前端请求 `/api/admin/dashboard`。
3. 后端统计用户总数、角色数量、活跃用户、退款数量。
4. 前端渲染统计卡片和待办提示。

### 核心代码

- 订单服务：`backend/internal/service/order_service.go`
- 用户仓储：`backend/internal/repository/User_repository.go`
- 管理端首页：`frontend/admin-dashboard.html`

### 模块关联

看板聚合用户、订单、退款数据，是管理端入口页。

### 实现思想

看板接口只返回聚合后的轻量数据，避免前端拉取大量明细再计算。

## 17. 风控模块

### 实现步骤

1. AI 调用、登录、退款等流程产生风险事件。
2. 后端写入 `risk_events`。
3. 管理员查看风控日志。
4. 管理员可以更新风险状态。

### 核心代码

- 风控模型：`backend/internal/model/risk_event.go`
- 风控服务：`backend/internal/service/risk_service.go`
- 风控接口：`backend/internal/handler/risk_handler.go`
- 风控仓储：`backend/internal/repository/risk_event_repository.go`
- 管理端风控页面：`frontend/admin-risk.html`

### 模块关联

风控模块与 AI 调用、退款审核、登录行为相关。当前主要用于记录和管理风险事件，后续可扩展复杂规则。

### 实现思想

风控事件作为独立日志，不阻塞主业务。命中风险后先记录，再由管理员人工处理。

## 18. AI 助手模块

### 实现步骤

1. 用户在 AI 助手页面输入问题。
2. 前端调用 `/api/ai/chat`。
3. 后端检查用户身份和 AI 限流。
4. AI 服务结合票务服务、订单服务、知识库服务生成回答。
5. 记录 Token 用量。
6. 返回 AI 回复给前端。

### 核心代码

- AI 服务：`backend/internal/service/ai_service.go`
- AI 接口：`backend/internal/handler/ai_handler.go`
- Token 用量模型：`backend/internal/model/token_usage.go`
- Token 用量服务：`backend/internal/service/token_usage_service.go`
- AI 助手页面：`frontend/ai-assistant.html`

### 模块关联

AI 模块可以调用票务搜索、订单查询、知识库检索和司机发车草稿生成能力，是系统的智能入口。

### 实现思想

AI 不直接替代业务接口，而是作为编排层：先理解用户意图，再调用已有业务服务，最后组织自然语言回答。

## 19. 知识库 RAG 模块

### 实现步骤

1. 管理员上传知识库文档。
2. 后端保存文档。
3. 文档按块切分。
4. 生成 embedding。
5. 用户提问时召回相关片段。
6. AI 根据片段生成回答。
7. 管理端可以搜索、查看、重建索引、删除文档。

### 核心代码

- 知识库服务：`backend/internal/service/knowledge_service.go`
- 知识库接口：`backend/internal/handler/knowledge_handler.go`
- 知识库文档目录：`backend/ragDocs`
- 内存文档目录：`backend/memory`
- 管理端知识库页面：`frontend/admin-knowledge.html`

### 模块关联

知识库服务为 AI 助手提供规则依据，尤其适合回答退改签、实名、购票限制等问题。

### 实现思想

RAG 的核心思想是“先检索，再生成”。模型回答前先查系统内可信文档，降低幻觉，提高回答可追溯性。

## 20. Token 用量与模型治理模块

### 实现步骤

1. AI 调用完成后记录模型、用户、场景和 Token 用量。
2. 管理员查看用户用量统计。
3. 管理端可以观察 AI 成本和调用趋势。
4. 风控模块可基于调用频率记录异常。

### 核心代码

- Token 用量模型：`backend/internal/model/token_usage.go`
- Token 用量服务：`backend/internal/service/token_usage_service.go`
- Token 用量接口：`backend/internal/handler/token_usage_handldler.go`
- Token 用量仓储：`backend/internal/repository/token_usage_repository.go`
- 管理端 Token 页面：`frontend/admin-tokens.html`

### 模块关联

Token 用量模块关联 AI 助手、知识库、风控和管理端报表。

### 实现思想

AI 调用是有成本的，必须记录用量并做限流。先记录可观测数据，再进一步做额度控制。

## 21. 前端页面与交互模块

### 实现步骤

1. 每个页面保留静态 HTML 入口。
2. 页面通过 `data-page` 标识当前模块。
3. 前端统一加载 `assets/app.js`。
4. `app.js` 根据当前页面初始化对应逻辑。
5. 所有 API 地址统一写在 `API_ENDPOINTS`。
6. 全局样式写在 `assets/styles.css`。
7. 导航根据登录角色动态渲染。

### 核心代码

- 前端逻辑：`frontend/assets/app.js`
- 前端样式：`frontend/assets/styles.css`
- 首页：`frontend/index.html`
- 搜索页：`frontend/search.html`
- 订单页：`frontend/orders.html`
- 司机端页面：`frontend/driver-*.html`
- 管理端页面：`frontend/admin-*.html`

### 模块关联

前端页面通过统一 API 调用后端，不直接操作数据库。所有业务状态以后端返回为准。

### 实现思想

当前前端是轻量静态原型，便于快速演示和联调。后续如果迁移 Vue，可以保留接口和业务结构，只替换前端实现。

## 22. 部署模块

### 实现步骤

1. 使用 Dockerfile 构建后端镜像。
2. 使用 Dockerfile 构建前端 Nginx 镜像。
3. Docker Compose 启动 MySQL、Redis、后端和前端。
4. 后端通过环境变量读取数据库和 Redis 地址。
5. 前端通过 Nginx 提供静态资源。

### 核心代码

- 后端 Dockerfile：`backend/Dockerfile`
- 前端 Dockerfile：`frontend/Dockerfile`
- Compose 编排：`docker-compose.yml`
- 前端 Nginx 配置：`frontend/nginx.conf`

### 模块关联

部署模块把前端、后端、数据库、缓存连接成完整运行环境。

### 实现思想

通过 Docker Compose 固化本地运行环境，减少“本机能跑、别人不能跑”的问题。

## 23. 后续扩展建议

### 前端工程化

- 新建 `frontend-vue`。
- 使用 Vue 3 + TypeScript + Vite。
- 使用 Element Plus 重建当前页面。
- 使用 Pinia 管理登录态、用户信息、搜索条件。
- 使用 TanStack Query 管理接口缓存和刷新。

### 数据库升级

- 当前项目使用 MySQL。
- 如果严格按路线图升级，可迁移 PostgreSQL。
- RAG 向量检索可使用 pgvector。

### 文件存储

- 司机证件图片和车辆图片可接入 MinIO。
- 数据库只保存文件 URL。

### 监控与 CI

- 使用 GitHub Actions 做构建和测试。
- 使用 Prometheus + Grafana 监控接口耗时、错误率、AI 调用量。

## 24. 模块实现顺序建议

1. 订单支付闭环
2. 电子票与司机核验
3. 退款审核与审计日志
4. 常用乘车人与低价提醒
5. 司机资质与车辆管理
6. 管理端数据看板与风控
7. AI 知识库客服助手
8. Vue 前端工程化迁移
9. PostgreSQL + pgvector RAG 升级
10. MinIO、CI、监控体系

## 25. 各模块技术栈与关键代码块

本节按模块补充“使用了哪些技术栈”和“技术栈对应的关键代码块”。代码块只展示核心片段，完整实现以对应文件为准。

### 25.1 后端启动与依赖注入

#### 使用技术栈

- Go：后端主语言。
- net/http：HTTP Server 和路由承载。
- GORM：数据库连接、自动迁移和 ORM。
- MySQL Driver：连接 MySQL。
- Redis Go Client：连接 Redis。
- 分层架构：Repository、Service、Handler、Router。

#### 关键代码块

文件：`backend/cmd/server/main.go`

```go
db, err := gorm.Open(mysql.Open(cfg.MySQLDSN), &gorm.Config{
	Logger: gormLogger,
})
if err != nil {
	log.Fatalf("connect mysql failed: %v", err)
}

redisClient := redis.NewClient(&redis.Options{
	Addr:     cfg.RedisAddr,
	Password: cfg.RedisPassword,
	DB:       cfg.RedisDB,
})
```

```go
userRepo := repository.NewGormUserRepository(db)
tripRepo := repository.NewGormTripRepository(db)
orderRepo := repository.NewGormOrderRepository(db)
paymentRepo := repository.NewGormPaymentRepository(db)

ticketService := service.NewTicketService(tripRepo)
orderService := service.NewOrderService(
	orderRepo,
	tripRepo,
	userRepo,
	paymentRepo,
	notificationRepo,
	auditRepo,
	electronicTicketService,
)
```

#### 技术思想

后端入口只负责组装依赖，不直接写业务逻辑。业务规则放在 Service，数据访问放在 Repository，HTTP 入参出参放在 Handler。

### 25.2 路由系统

#### 使用技术栈

- net/http ServeMux：注册 API 路由。
- REST API：按资源组织接口。
- 中间件：对需要登录的接口包裹 JWT 校验。
- CORS：支持前端跨域访问。

#### 关键代码块

文件：`backend/internal/router/router.go`

```go
mount(mux, "/api/orders", map[string]http.Handler{
	http.MethodPost: dep.AuthMiddleware(http.HandlerFunc(dep.OrderHandler.CreateOrder)),
})

mount(mux, "/api/orders/{orderId}/ticket", map[string]http.Handler{
	http.MethodGet: dep.AuthMiddleware(http.HandlerFunc(dep.ElectronicTicketHandler.GetMyOrderTicket)),
})

mount(mux, "/api/driver/tickets/verify", map[string]http.Handler{
	http.MethodPost: dep.AuthMiddleware(http.HandlerFunc(dep.ElectronicTicketHandler.VerifyByDriver)),
})
```

```go
func mount(mux *http.ServeMux, path string, handlers map[string]http.Handler) {
	mux.Handle(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		h, ok := handlers[r.Method]
		if !ok {
			response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		h.ServeHTTP(w, r)
	}))
}
```

#### 技术思想

路由层只做请求分发。是否需要登录由路由注册时决定，避免 Handler 内重复写鉴权逻辑。

### 25.3 JWT 认证模块

#### 使用技术栈

- JWT：无状态登录凭证。
- Go Context：在请求上下文中保存当前用户。
- HTTP Middleware：统一拦截受保护接口。

#### 关键代码块

文件：`backend/internal/pkg/middleware/auth.go`

```go
claims, err := jwtManager.ParseToken(tokenStr)
if err != nil {
	response.Error(w, http.StatusUnauthorized, "invalid token")
	return
}

currentUser := &CurrentAuthUser{
	ID:    claims.UserID,
	Role:  claims.Role,
	Phone: claims.Phone,
}

ctx := context.WithValue(r.Context(), userContextKey, currentUser)
next.ServeHTTP(w, r.WithContext(ctx))
```

Handler 读取当前用户：

```go
currentUser, ok := middleware.CurrentUser(r)
if !ok {
	response.Error(w, http.StatusUnauthorized, "unauthorized")
	return
}
```

#### 技术思想

JWT 里携带用户 ID 和角色。业务层通过角色判断当前用户是否允许执行某个操作。

### 25.4 用户与角色模块

#### 使用技术栈

- GORM Model：定义用户表结构。
- JWT：保存登录态。
- RBAC 简化模型：当前项目先用 `role` 字段区分乘客、司机、管理员。

#### 关键代码块

文件：`backend/internal/model/user.go`

```go
const (
	RolePassenger = "passenger"
	RoleDriver    = "driver"
	RoleAdmin     = "admin"
)

type User struct {
	ID               uint   `gorm:"primaryKey" json:"id"`
	Phone            string `gorm:"size:20;uniqueIndex;not null" json:"phone"`
	Nickname         string `gorm:"size:50;not null" json:"nickname"`
	Role             string `gorm:"size:20;not null;default:passenger" json:"role"`
	DefaultRole      string `gorm:"column:default_role;size:20;not null;default:passenger" json:"defaultRole"`
	RealNameVerified bool   `gorm:"column:real_name_verified;not null;default:false" json:"realNameVerified"`
	Status           string `gorm:"size:20;not null;default:active" json:"status"`
}
```

角色校验示例：

```go
if currentUserRole != model.RolePassenger {
	return nil, errors.New("only passenger can create order")
}
```

#### 技术思想

先用简单角色字段完成权限控制，后续可以扩展成 `roles`、`permissions`、`role_permissions` 的完整 RBAC。

### 25.5 票务搜索模块

#### 使用技术栈

- GORM 查询：按城市、日期、状态、余票过滤。
- MySQL 索引：通过 `start_city`、`end_city`、`departure_time`、`status` 提升查询效率。
- 前端 Fetch：搜索页调用 API。

#### 关键代码块

文件：`backend/internal/repository/trip_repository.go`

```go
err := r.db.WithContext(ctx).
	Preload("Stops", func(db *gorm.DB) *gorm.DB {
		return db.Order("stop_order asc")
	}).
	Where("start_city = ? AND end_city = ?", startCity, endCity).
	Where("departure_time >= ? AND departure_time < ?", departureStart, departureEnd).
	Where("status = ?", model.TripStatusPublished).
	Where("seat_available > 0").
	Order("departure_time asc").
	Find(&trips).Error
```

文件：`backend/internal/model/trip.go`

```go
type Trip struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	StartCity     string    `gorm:"size:50;not null;index:idx_trip_search,priority:1" json:"startCity"`
	EndCity       string    `gorm:"size:50;not null;index:idx_trip_search,priority:2" json:"endCity"`
	DepartureTime time.Time `gorm:"not null;index:idx_trip_search,priority:3" json:"departureTime"`
	SeatAvailable int       `gorm:"not null" json:"seatAvailable"`
	Status        string    `gorm:"size:20;not null;default:published;index:idx_trip_search,priority:4" json:"status"`
}
```

#### 技术思想

乘客只能搜索到“已发布且有余票”的班次，避免展示不可购买资源。

### 25.6 订单创建与库存锁定模块

#### 使用技术栈

- GORM Transaction：订单创建和座位扣减放在同一个事务。
- 数据库行锁：使用 `FOR UPDATE` 防止并发超卖。
- 状态机：订单状态从 `pending_payment` 开始流转。

#### 关键代码块

文件：`backend/internal/repository/order_repository.go`

```go
return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
	var trip model.Trip
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&trip, order.TripID).Error
	if err != nil {
		return err
	}

	if trip.SeatAvailable < order.TicketCount {
		return errors.New("not enough seats available")
	}

	trip.SeatAvailable -= order.TicketCount
	if err := tx.Save(&trip).Error; err != nil {
		return err
	}

	return tx.Create(order).Error
})
```

文件：`backend/internal/service/order_service.go`

```go
order := &model.Order{
	OrderNo:         generateOrderNo(),
	UserID:          currentUserID,
	TripID:          input.TripID,
	TicketCount:     input.TicketCount,
	Amount:          trip.PriceCent * input.TicketCount,
	PayStatus:       model.PayStatusUnpaid,
	OrderStatus:     model.OrderStatusPendingPayment,
	RefundStatus:    model.RefundStatusNone,
	PaymentExpireAt: &expireAt,
}
```

#### 技术思想

订单和座位库存必须强一致。创建订单成功但扣座失败，或者扣座成功但订单失败，都会造成数据错乱，所以使用事务。

### 25.7 支付模块

#### 使用技术栈

- Mock Payment：本地演示支付成功。
- GORM Transaction：支付状态和订单状态同步更新。
- 可扩展 Channel：后续可接微信、支付宝。

#### 关键代码块

文件：`backend/internal/repository/payment_repository.go`

```go
payment.Status = model.PaymentStatusPaid
payment.PaidAt = &now

order.PayStatus = model.PayStatusPaid
if order.OrderStatus == model.OrderStatusPendingPayment {
	order.OrderStatus = model.OrderStatusPendingVerification
}

if err := tx.Save(&order).Error; err != nil {
	return err
}
return tx.Save(&payment).Error
```

文件：`backend/internal/service/payment_service.go`

```go
if err := s.paymentRepo.MarkPaid(ctx, paymentID); err != nil {
	return nil, err
}
if s.ticketService != nil {
	paidOrder, orderErr := s.orderRepo.GetByID(ctx, payment.OrderID)
	if orderErr != nil {
		return nil, orderErr
	}
	if paidOrder != nil {
		if _, ticketErr := s.ticketService.EnsureForOrder(ctx, paidOrder); ticketErr != nil {
			return nil, ticketErr
		}
	}
}
```

#### 技术思想

支付成功是订单进入后续履约流程的分界点。支付成功后，订单从待支付变为待核验，并自动生成电子票。

### 25.8 电子票与司机核验模块

#### 使用技术栈

- HMAC-SHA256：电子票签名。
- Base64 URL Encoding：生成可复制、可二维码化的 token。
- SHA256 Hash：数据库保存 token hash。
- GORM Transaction：核验电子票和更新订单状态保持一致。

#### 关键代码块

文件：`backend/internal/service/electronic_ticket_service.go`

```go
payload := fmt.Sprintf(
	"%d:%d:%d:%d:%s",
	order.ID,
	order.UserID,
	order.TripID,
	expiresAt.Unix(),
	hex.EncodeToString(nonceBytes),
)
signature := s.sign(payload)
rawToken := payload + "." + signature
token := base64.RawURLEncoding.EncodeToString([]byte(rawToken))
```

```go
func (s *ElectronicTicketService) sign(payload string) string {
	mac := hmac.New(sha256.New, []byte(s.secret))
	mac.Write([]byte(payload))
	return hex.EncodeToString(mac.Sum(nil))
}
```

文件：`backend/internal/repository/electronic_ticket_repository.go`

```go
ticket.Status = model.ElectronicTicketStatusVerified
ticket.VerifiedAt = &verifiedAt
ticket.VerifiedByDriverID = &driverID

if err := tx.Save(&ticket).Error; err != nil {
	return err
}

return tx.Model(&model.Order{}).
	Where("id = ? AND order_status = ?", ticket.OrderID, model.OrderStatusPendingVerification).
	Update("order_status", model.OrderStatusCompleted).Error
```

#### 技术思想

电子票不能只靠订单 ID，因为订单 ID 容易猜测。使用 HMAC 签名 token 可以防止伪造，司机核验时再校验班次归属。

### 25.9 退款审核与审计模块

#### 使用技术栈

- 状态机：退款状态 `none -> requested -> refunded/rejected`。
- GORM Transaction：通过退款时释放库存和更新订单状态。
- 审计日志：记录管理员操作。
- 通知模块：把审核结果通知乘客。

#### 关键代码块

文件：`backend/internal/repository/order_repository.go`

```go
lockedOrder.RefundStatus = model.RefundStatusRefunded
lockedOrder.RefundReviewNote = reviewNote
lockedOrder.RefundReviewedAt = &now

if lockedOrder.OrderStatus == model.OrderStatusPendingVerification {
	var trip model.Trip
	err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&trip, lockedOrder.TripID).Error
	if err != nil {
		return err
	}

	trip.SeatAvailable += lockedOrder.TicketCount
	lockedOrder.OrderStatus = model.OrderStatusCancelled

	if err := tx.Save(&trip).Error; err != nil {
		return err
	}
}
```

文件：`backend/internal/service/order_service.go`

```go
if err := s.orderRepo.ApproveRefund(ctx, order, reviewNote); err != nil {
	return nil, err
}
if s.ticketService != nil {
	_ = s.ticketService.VoidByOrderID(ctx, order.ID)
}
s.createRefundAuditLog(ctx, order.ID, model.RefundStatusRefunded, currentUserID, reviewNote)
s.createAuditLog(ctx, currentUserID, currentUserRole, "refund.approve", "order", fmt.Sprintf("%d", order.ID), reviewNote)
```

文件：`backend/internal/model/audit.go`

```go
type RefundAuditLog struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	OrderID      uint      `gorm:"column:order_id;not null;index" json:"orderId"`
	RefundStatus string    `gorm:"column:refund_status;size:20;not null;index" json:"refundStatus"`
	ReviewNote   string    `gorm:"column:review_note;size:255" json:"reviewNote"`
	ReviewerID   uint      `gorm:"column:reviewer_id;not null;index" json:"reviewerId"`
	CreatedAt    time.Time `json:"createdAt"`
}
```

#### 技术思想

退款是高风险操作，必须有状态限制、库存回滚、电子票作废、通知和审计记录。

### 25.10 通知中心模块

#### 使用技术栈

- GORM Model：通知持久化。
- 站内信：先实现系统内通知。
- 后续可扩展异步队列：短信、邮件、微信模板消息。

#### 关键代码块

文件：`backend/internal/model/notification.go`

```go
type Notification struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	UserID         uint       `gorm:"column:user_id;not null;index:idx_notifications_user_read_created,priority:1" json:"userId"`
	Type           string     `gorm:"size:50;not null" json:"type"`
	Title          string     `gorm:"size:100;not null" json:"title"`
	Content        string     `gorm:"type:text;not null" json:"content"`
	RelatedOrderID *uint      `gorm:"column:related_order_id;index" json:"relatedOrderId,omitempty"`
	IsRead         bool       `gorm:"column:is_read;not null;default:false;index:idx_notifications_user_read_created,priority:2" json:"isRead"`
}
```

文件：`backend/internal/service/order_service.go`

```go
_ = s.notificationRepo.Create(ctx, &model.Notification{
	UserID:         userID,
	Type:           typ,
	Title:          title,
	Content:        content,
	RelatedOrderID: relatedOrderID,
	IsRead:         false,
})
```

#### 技术思想

通知不应该阻塞主业务。即使通知创建失败，订单退款等主流程也不应该回滚。

### 25.11 常用乘车人模块

#### 使用技术栈

- GORM Model：保存乘车人实名信息。
- REST API：提供新增、列表、删除。
- JWT：保证用户只能操作自己的乘车人。

#### 关键代码块

文件：`backend/internal/model/extension.go`

```go
type Passenger struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"column:user_id;not null;index" json:"userId"`
	Name      string    `gorm:"size:50;not null" json:"name"`
	IDCard    string    `gorm:"column:id_card;size:32;not null" json:"idCard"`
	Phone     string    `gorm:"size:20;not null" json:"phone"`
	IsDefault bool      `gorm:"column:is_default;not null;default:false" json:"isDefault"`
	CreatedAt time.Time `json:"createdAt"`
}
```

文件：`backend/internal/service/extension_service.go`

```go
if role != model.RolePassenger {
	return nil, errors.New("only passenger can manage passengers")
}
passenger.UserID = userID
passenger.Name = strings.TrimSpace(passenger.Name)
passenger.IDCard = strings.TrimSpace(passenger.IDCard)
passenger.Phone = strings.TrimSpace(passenger.Phone)
```

文件：`backend/internal/router/router.go`

```go
mount(mux, "/api/passengers", map[string]http.Handler{
	http.MethodGet:  dep.AuthMiddleware(http.HandlerFunc(dep.PassengerHandler.List)),
	http.MethodPost: dep.AuthMiddleware(http.HandlerFunc(dep.PassengerHandler.Create)),
})
```

#### 技术思想

乘车人信息属于用户私有数据，所有操作都必须绑定当前 JWT 用户 ID。

### 25.12 低价提醒模块

#### 使用技术栈

- GORM Model：保存订阅条件。
- REST API：创建、列表、禁用。
- 日期范围：控制提醒有效期。
- 后续可接定时任务：扫描价格并创建通知。

#### 关键代码块

文件：`backend/internal/model/extension.go`

```go
type PriceAlert struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	UserID          uint      `gorm:"column:user_id;not null;index" json:"userId"`
	StartCity       string    `gorm:"column:start_city;size:50;not null;index" json:"startCity"`
	EndCity         string    `gorm:"column:end_city;size:50;not null;index" json:"endCity"`
	TargetPriceCent int      `gorm:"column:target_price_cent;not null" json:"targetPriceCent"`
	StartDate       time.Time `gorm:"column:start_date;type:date;not null" json:"startDate"`
	EndDate         time.Time `gorm:"column:end_date;type:date;not null" json:"endDate"`
	Status          string    `gorm:"size:20;not null;default:active;index" json:"status"`
}
```

文件：`backend/internal/service/extension_service.go`

```go
if alert.StartDate.IsZero() {
	alert.StartDate = time.Now()
}
if alert.EndDate.IsZero() || alert.EndDate.Before(alert.StartDate) {
	alert.EndDate = alert.StartDate.AddDate(0, 0, 7)
}
```

#### 技术思想

低价提醒不在用户请求时实时扫描全部班次，而是保存订阅条件，交给后续定时任务批量处理。

### 25.13 司机班次与车辆资质模块

#### 使用技术栈

- GORM Model：保存司机资质和车辆。
- REST API：司机资料提交、车辆新增、车辆列表。
- 状态字段：资质审核状态。

#### 关键代码块

文件：`backend/internal/model/extension.go`

```go
type DriverProfile struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	UserID         uint       `gorm:"column:user_id;not null;uniqueIndex" json:"userId"`
	RealName       string     `gorm:"column:real_name;size:50;not null" json:"realName"`
	IDCard         string     `gorm:"column:id_card;size:32;not null" json:"idCard"`
	LicenseNo      string     `gorm:"column:license_no;size:64;not null" json:"licenseNo"`
	Status         string     `gorm:"size:20;not null;default:pending;index" json:"status"`
	ReviewedAt     *time.Time `gorm:"column:reviewed_at" json:"reviewedAt,omitempty"`
}
```

```go
type Vehicle struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	DriverID  uint      `gorm:"column:driver_id;not null;index" json:"driverId"`
	PlateNo   string    `gorm:"column:plate_no;size:32;not null" json:"plateNo"`
	Brand     string    `gorm:"size:64" json:"brand"`
	SeatCount int       `gorm:"column:seat_count;not null" json:"seatCount"`
	Status    string    `gorm:"size:20;not null;default:active;index" json:"status"`
}
```

文件：`backend/internal/router/router.go`

```go
mount(mux, "/api/driver/profile", map[string]http.Handler{
	http.MethodGet: dep.AuthMiddleware(http.HandlerFunc(dep.DriverProfileHandler.GetProfile)),
	http.MethodPut: dep.AuthMiddleware(http.HandlerFunc(dep.DriverProfileHandler.UpsertProfile)),
})

mount(mux, "/api/driver/vehicles", map[string]http.Handler{
	http.MethodGet:  dep.AuthMiddleware(http.HandlerFunc(dep.DriverProfileHandler.ListVehicles)),
	http.MethodPost: dep.AuthMiddleware(http.HandlerFunc(dep.DriverProfileHandler.CreateVehicle)),
})
```

#### 技术思想

司机资料和车辆资料先独立建模。后续发布班次时可以增加校验：资质通过且有可用车辆才允许发车。

### 25.14 AI 助手模块

#### 使用技术栈

- OpenAI-Compatible API：兼容 OpenAI 协议的大模型服务。
- HTTP Client：调用模型接口。
- Redis 限流：控制用户 AI 调用频率。
- Token Usage：记录模型调用成本。
- RAG：结合知识库增强回答。

#### 关键代码块

文件：`backend/cmd/server/main.go`

```go
aiService := service.NewAIService(
	cfg.OpenAIAPIKey,
	cfg.OpenAIBaseURL,
	cfg.OpenAIModel,
	cfg.OpenAITimeout,
	ticketService,
	orderService,
	userService,
	knowledgeService,
	tokenUsageService,
)
```

文件：`backend/internal/router/router.go`

```go
mount(mux, "/api/ai/chat", map[string]http.Handler{
	http.MethodPost: dep.AuthMiddleware(http.HandlerFunc(dep.AIHandler.ChatPassenger)),
})

mount(mux, "/api/ai/driver/create-trip", map[string]http.Handler{
	http.MethodPost: dep.AuthMiddleware(http.HandlerFunc(dep.AIHandler.CreateDriverTripDraft)),
})
```

#### 技术思想

AI 模块不直接修改核心业务数据，而是调用已有 Service 获取数据、组织建议或生成草稿，降低 AI 误操作风险。

### 25.15 知识库 RAG 模块

#### 使用技术栈

- 文档切分：按 chunk 处理长文档。
- Embedding API：把文本转向量。
- Rerank API：对召回结果重排序。
- 本地文件存储：保存知识库文档和索引数据。
- Token Usage：记录 embedding、rerank、chat 调用成本。

#### 关键代码块

文件：`backend/cmd/server/main.go`

```go
knowledgeService := service.NewKnowledgeService(
	cfg.KnowledgeMemoryDir,
	cfg.EmbeddingAPIKey,
	cfg.EmbeddingBaseURL,
	cfg.EmbeddingModel,
	cfg.RerankAPIKey,
	cfg.RerankBaseURL,
	cfg.RerankModel,
	tokenUsageService,
	cfg.KnowledgeChunkSize,
	cfg.KnowledgeChunkOverlap,
	cfg.KnowledgeRecallLimit,
	cfg.KnowledgeTopK,
)
```

文件：`backend/internal/router/router.go`

```go
mount(mux, "/api/admin/knowledge/upload", map[string]http.Handler{
	http.MethodPost: dep.AuthMiddleware(http.HandlerFunc(dep.KnowledgeHandler.UploadDocument)),
})

mount(mux, "/api/admin/knowledge/search", map[string]http.Handler{
	http.MethodPost: dep.AuthMiddleware(http.HandlerFunc(dep.KnowledgeHandler.SearchDocuments)),
})
```

#### 技术思想

RAG 的核心是“先从可信知识库检索，再让模型基于检索内容回答”，避免模型凭空编造票务规则。

### 25.16 Token 用量与限流模块

#### 使用技术栈

- Redis：短窗口限流。
- MySQL：Token 用量持久化。
- GORM：记录和聚合用量。
- 管理端 API：查看用户或场景用量。

#### 关键代码块

文件：`backend/cmd/server/main.go`

```go
aiLimiter := repository.NewRedisAIRateLimiter(redisClient, cfg.AIRateLimitPrefix)
tokenUsageRepo := repository.NewTokenUsageRepository(db)
tokenUsageService := service.NewTokenUsageService(tokenUsageRepo)
```

文件：`backend/internal/router/router.go`

```go
mount(mux, "/api/admin/tokens", map[string]http.Handler{
	http.MethodGet: dep.AuthMiddleware(http.HandlerFunc(dep.TokenUsageHandler.ListAdminUsage)),
})
```

#### 技术思想

AI 调用会产生费用，因此必须记录用量并限制频率。限流用 Redis，历史统计用数据库。

### 25.17 前端 API 调用模块

#### 使用技术栈

- 原生 JavaScript。
- Fetch API。
- 统一 API_ENDPOINTS。
- data-* DOM 钩子。
- sessionStorage/localStorage 保存登录态和临时状态。

#### 关键代码块

文件：`frontend/assets/app.js`

```js
const API_ENDPOINTS = Object.freeze({
  passenger: Object.freeze({
    searchTickets: "/tickets/search",
    createOrder: "/orders",
    orderDetail: "/orders/:orderId",
    orderTicket: "/orders/:orderId/ticket",
    createPayment: "/payments/create",
    mockPaymentSuccess: "/payments/:paymentId/mock-success",
  }),
  driver: Object.freeze({
    trips: "/driver/trips",
    verifyTicket: "/driver/tickets/verify",
  }),
});
```

```js
const result = await api.get(API_ENDPOINTS.passenger.orderTicket, undefined, {
  pathParams: { orderId },
});
```

```js
const result = await api.post(API_ENDPOINTS.driver.verifyTicket, { token });
```

#### 技术思想

前端不把接口路径散落到各处，而是集中在 `API_ENDPOINTS`，页面逻辑通过 `data-*` 找 DOM 节点并绑定交互。

### 25.18 前端样式模块

#### 使用技术栈

- CSS Variables：统一主题变量。
- 响应式 Grid：适配桌面和移动端。
- 状态样式：hover、focus、empty、loading。
- 浅蓝企业级 UI 风格。

#### 关键代码块

文件：`frontend/assets/styles.css`

```css
:root {
  --color-primary: #3b82f6;
  --color-primary-soft: #eff6ff;
  --color-bg: #f8fbff;
  --color-surface: #ffffff;
  --color-border: #dbeafe;
  --color-text-main: #0f172a;
  --shadow-soft: 0 12px 30px rgba(37, 99, 235, 0.08);
  --radius-card: 18px;
}
```

```css
.button-primary {
  color: #fff;
  background: var(--color-primary);
  border-color: var(--color-primary);
  box-shadow: 0 8px 20px rgba(59, 130, 246, 0.22);
}
```

#### 技术思想

使用 CSS 变量统一控制颜色、圆角、阴影和间距，避免每个页面各写一套样式。
