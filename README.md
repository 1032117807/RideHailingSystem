# RideHailingSystem

本项目已整理为适合前后端同仓开发的结构，前端原型、后端代码和业务文档分目录维护，避免后续开发时文件混杂在根目录。

## 目录结构

```text
RideHailingSystem/
  backend/
    src/
  docs/
    买车票.md
    前后端分离规范.md
  frontend/
    assets/
    *.html
  .editorconfig
  README.md
```

## 目录约定

- `frontend/`：前端页面原型、静态资源、运行时配置
- `backend/`：后端服务代码
- `docs/`：需求、设计、接口、规范文档

## 当前入口

- 前端预览入口：`frontend/index.html`
- 前后端分离规范：`docs/前后端分离规范.md`
- 系统设计文档：`docs/买车票.md`

## 后续建议

### 前端

- 继续在 `frontend/` 里维护原型
- 后续迁移到 Vue / React 时，建议直接把 `frontend/` 升级为正式前端工程

### 后端

- 在 `backend/src/` 下开始创建后端项目
- 后续按 `controller / service / repository / model` 或你使用的框架约定拆分

### 文档

- 需求、接口、数据库设计统一放到 `docs/`
- 不再把业务文档散落在项目根目录
