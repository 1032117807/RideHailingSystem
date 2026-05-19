# 前端通用 Skill

## 1. Skill 定位

这是一个用于生成、优化和统一前端页面的通用 Skill。  
适用于 Web 管理后台、官网、仪表盘、数据看板、表单系统、低代码页面、SaaS 产品界面等前端场景。

整体风格以 **浅蓝色、干净、现代、柔和、专业** 为主，强调清晰的信息层级、舒适的视觉间距和良好的交互体验。

---

## 2. 设计目标

- 页面整体清爽、简洁、易读。
- 色调以浅蓝色为主，避免高饱和刺眼颜色。
- 组件风格统一，适合企业级系统和通用前端项目。
- 布局具备响应式能力，兼容桌面端、平板和移动端。
- 代码结构清晰，方便维护、扩展和复用。
- 交互反馈明确，包括 hover、active、focus、loading、empty、error 等状态。

---

## 3. 视觉风格

### 3.1 主色调

推荐使用浅蓝色系作为主视觉：

```css
:root {
  --color-primary: #3b82f6;
  --color-primary-light: #dbeafe;
  --color-primary-soft: #eff6ff;
  --color-primary-hover: #2563eb;

  --color-bg: #f8fbff;
  --color-surface: #ffffff;
  --color-surface-soft: #f1f7ff;

  --color-border: #dbeafe;
  --color-border-strong: #bfdbfe;

  --color-text-main: #0f172a;
  --color-text-secondary: #475569;
  --color-text-muted: #94a3b8;

  --color-success: #22c55e;
  --color-warning: #f59e0b;
  --color-error: #ef4444;
}
```

### 3.2 推荐色彩使用

| 类型 | 颜色 | 用途 |
|---|---|---|
| 页面背景 | `#f8fbff` | 页面整体背景 |
| 卡片背景 | `#ffffff` | 内容卡片、弹窗、表单容器 |
| 浅蓝背景 | `#eff6ff` | 区块背景、提示区域 |
| 主按钮 | `#3b82f6` | 主要操作按钮 |
| Hover 状态 | `#2563eb` | 按钮悬停 |
| 边框 | `#dbeafe` | 卡片、输入框、分割线 |
| 主文本 | `#0f172a` | 标题、正文重点 |
| 次级文本 | `#475569` | 描述文字、说明信息 |

---

## 4. 字体与排版

### 4.1 字体

优先使用系统字体：

```css
font-family:
  Inter,
  -apple-system,
  BlinkMacSystemFont,
  "Segoe UI",
  "PingFang SC",
  "Microsoft YaHei",
  Arial,
  sans-serif;
```

### 4.2 字号建议

| 场景 | 字号 | 字重 |
|---|---:|---:|
| 页面主标题 | 28px - 36px | 700 |
| 区块标题 | 20px - 24px | 600 |
| 卡片标题 | 16px - 18px | 600 |
| 正文 | 14px - 16px | 400 |
| 辅助说明 | 12px - 14px | 400 |

---

## 5. 布局规范

### 5.1 页面结构

推荐页面结构：

```text
App
├── Header / TopBar
├── Sidebar / Navigation
├── Main Content
│   ├── Page Header
│   ├── Filter / Search Area
│   ├── Data Cards / Charts / Table
│   └── Pagination / Footer Actions
└── Footer
```

### 5.2 间距规范

```css
--space-xs: 4px;
--space-sm: 8px;
--space-md: 16px;
--space-lg: 24px;
--space-xl: 32px;
--space-2xl: 48px;
```

使用建议：

- 卡片内边距：`20px - 24px`
- 页面左右边距：`24px - 40px`
- 区块之间间距：`24px - 32px`
- 表单项间距：`16px - 20px`

---

## 6. 组件规范

### 6.1 Button 按钮

按钮应具备清晰状态：

```css
.button-primary {
  background: #3b82f6;
  color: #ffffff;
  border: 1px solid #3b82f6;
  border-radius: 10px;
  padding: 10px 18px;
  font-weight: 600;
  transition: all 0.2s ease;
}

.button-primary:hover {
  background: #2563eb;
  border-color: #2563eb;
  box-shadow: 0 8px 20px rgba(59, 130, 246, 0.22);
}

.button-secondary {
  background: #eff6ff;
  color: #2563eb;
  border: 1px solid #bfdbfe;
}
```

按钮类型：

- Primary：页面主操作，如提交、保存、创建。
- Secondary：次级操作，如取消、返回、筛选。
- Ghost：低强调操作，如更多、查看详情。
- Danger：删除、移除等风险操作。

---

### 6.2 Card 卡片

```css
.card {
  background: #ffffff;
  border: 1px solid #dbeafe;
  border-radius: 18px;
  padding: 24px;
  box-shadow: 0 12px 30px rgba(37, 99, 235, 0.08);
}
```

卡片设计要求：

- 圆角建议 `16px - 20px`。
- 阴影要轻，不要过重。
- 卡片之间保持足够留白。
- 适合展示数据、表单、图表、说明区块。

---

### 6.3 Input 输入框

```css
.input {
  width: 100%;
  height: 42px;
  border: 1px solid #bfdbfe;
  border-radius: 10px;
  padding: 0 12px;
  background: #ffffff;
  color: #0f172a;
  outline: none;
  transition: all 0.2s ease;
}

.input:focus {
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.15);
}
```

---

### 6.4 Table 表格

表格适合管理后台和数据系统：

```css
.table {
  width: 100%;
  border-collapse: separate;
  border-spacing: 0;
  background: #ffffff;
  border: 1px solid #dbeafe;
  border-radius: 16px;
  overflow: hidden;
}

.table th {
  background: #eff6ff;
  color: #1e3a8a;
  font-weight: 600;
  padding: 14px 16px;
}

.table td {
  padding: 14px 16px;
  border-top: 1px solid #e0f2fe;
  color: #334155;
}
```

表格状态：

- Loading：显示骨架屏。
- Empty：显示空状态插画或提示。
- Error：显示错误提示和重试按钮。
- Hover：行悬停使用浅蓝背景。

---

### 6.5 Modal 弹窗

弹窗设计：

- 背景遮罩使用半透明深色。
- 弹窗主体使用白色卡片。
- 圆角建议 `20px`。
- 标题、内容、操作按钮区域分明。
- 关闭按钮位置清晰。

---

## 7. 页面类型模板

### 7.1 管理后台首页

适合结构：

```text
页面标题 + 欢迎语
数据概览卡片
趋势图表
快捷入口
最近操作 / 最新数据表格
```

视觉建议：

- 背景使用 `#f8fbff`
- 数据卡片使用白色
- 图表主色使用浅蓝渐变
- 重要数字使用深蓝色

---

### 7.2 表单页面

适合结构：

```text
页面标题
说明文字
表单卡片
分组字段
底部操作按钮
```

设计要求：

- 表单字段分组明确。
- 必填项标注清晰。
- 错误提示靠近对应字段。
- 提交按钮放在明显位置。

---

### 7.3 数据列表页

适合结构：

```text
页面标题
搜索与筛选区域
操作按钮区
数据表格
分页器
```

设计要求：

- 搜索区可折叠。
- 表格支持加载、空状态、错误状态。
- 批量操作要有二次确认。

---

### 7.4 官网 Landing Page

适合结构：

```text
Hero 首屏
核心卖点
产品功能
使用流程
客户案例
价格方案
FAQ
CTA 行动按钮
Footer
```

设计要求：

- 首屏使用浅蓝渐变背景。
- CTA 按钮突出但不过分刺眼。
- 内容区块间距充足。
- 图文结合，避免大段文字堆积。

---

## 8. 响应式规范

### 8.1 断点建议

```css
--breakpoint-sm: 640px;
--breakpoint-md: 768px;
--breakpoint-lg: 1024px;
--breakpoint-xl: 1280px;
```

### 8.2 响应式原则

- 移动端优先保证可读性。
- 卡片在移动端纵向排列。
- 表格在移动端可横向滚动。
- 导航栏在小屏幕下转为抽屉菜单。
- 按钮在移动端适当增大点击区域。

---

## 9. 交互状态

每个组件都应考虑以下状态：

- Default：默认状态。
- Hover：鼠标悬停。
- Active：点击中。
- Focus：键盘聚焦。
- Disabled：不可用。
- Loading：加载中。
- Empty：无数据。
- Error：错误。
- Success：成功反馈。

示例：

```css
.interactive {
  transition:
    background-color 0.2s ease,
    border-color 0.2s ease,
    box-shadow 0.2s ease,
    transform 0.2s ease;
}

.interactive:hover {
  transform: translateY(-1px);
}
```

---

## 10. 动效规范

动效应轻量、自然，不影响性能。

推荐：

```css
.fade-in {
  animation: fadeIn 0.24s ease-out;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(6px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
```

动效原则：

- 不使用夸张弹跳。
- 不使用过多闪烁。
- 页面进入、卡片 hover、弹窗打开可使用轻动效。
- 动画时长建议 `160ms - 280ms`。

---

## 11. 可访问性要求

- 文本与背景保持足够对比度。
- 所有交互元素支持键盘操作。
- 表单控件要有 label。
- 图标按钮需要 aria-label。
- 错误提示不能只依赖颜色。
- 禁用状态要清晰但仍可读。

---

## 12. 推荐技术栈

可根据项目选择：

- React / Vue / Next.js / Nuxt
- TypeScript
- Tailwind CSS
- Ant Design / Element Plus / shadcn/ui
- ECharts / Recharts
- Framer Motion
- Zustand / Pinia
- React Query / TanStack Query
- Vite

---

## 13. Tailwind 浅蓝主题建议

```js
// tailwind.config.js
export default {
  theme: {
    extend: {
      colors: {
        brand: {
          50: "#eff6ff",
          100: "#dbeafe",
          200: "#bfdbfe",
          300: "#93c5fd",
          400: "#60a5fa",
          500: "#3b82f6",
          600: "#2563eb",
          700: "#1d4ed8",
          800: "#1e40af",
          900: "#1e3a8a"
        }
      },
      borderRadius: {
        card: "18px",
        button: "10px"
      },
      boxShadow: {
        soft: "0 12px 30px rgba(37, 99, 235, 0.08)"
      }
    }
  }
}
```

---

## 14. 前端生成提示词模板

当需要生成页面时，可以使用以下提示词：

```text
请生成一个现代化前端页面，整体采用浅蓝色调，风格干净、专业、柔和。

要求：
1. 使用响应式布局。
2. 页面背景使用极浅蓝或蓝白渐变。
3. 卡片使用白色背景、浅蓝边框、柔和阴影。
4. 主按钮使用蓝色，hover 时加深。
5. 字体层级清晰，留白充足。
6. 所有组件包含 hover、focus、loading、empty、error 等状态设计。
7. 代码结构清晰，组件可复用。
8. 优先使用 TypeScript 和现代前端写法。
9. 如果使用 Tailwind CSS，请使用浅蓝色系。
10. 不要使用过度复杂或高饱和的视觉效果。
```

---

## 15. 代码质量要求

- 使用语义化 HTML。
- 组件命名清晰。
- 样式变量统一管理。
- 避免重复代码。
- 表单、表格、弹窗等组件应可复用。
- 对异步请求做好 loading、error、empty 状态处理。
- 对复杂页面拆分为多个组件。
- 关键交互应有明确反馈。
- 保持代码可读性优先。

---

## 16. 默认页面风格描述

默认生成页面时，使用以下风格：

```text
浅蓝色科技感，蓝白配色，圆角卡片，柔和阴影，清晰留白，现代 SaaS 风格，适合企业后台、数据看板、产品官网和通用 Web 应用。
```

---

## 17. 禁止事项

- 不使用大面积深色背景，除非用户明确要求。
- 不使用刺眼的纯蓝高饱和背景。
- 不堆叠过多阴影和渐变。
- 不生成混乱的布局。
- 不忽略移动端适配。
- 不只生成静态 UI，应考虑真实交互状态。
- 不使用难以维护的内联样式堆砌。

---

## 18. 输出要求

生成前端代码时，应尽量包含：

- 页面结构说明。
- 关键组件。
- 样式变量。
- 响应式布局。
- 交互状态。
- 空状态、加载状态和错误状态。
- 简洁清晰的代码注释。
- 必要时提供可运行示例。

---

## 19. 适用场景

该 Skill 可用于：

- 生成后台管理系统页面。
- 生成 SaaS 官网。
- 生成数据看板。
- 生成表单系统。
- 优化已有前端页面。
- 统一项目 UI 风格。
- 生成组件库规范。
- 生成 Tailwind / CSS 主题配置。
- 生成 React / Vue 页面代码。
- 生成产品原型页面。

---

## 20. 总结

该前端通用 Skill 的核心原则是：

```text
浅蓝主色、蓝白背景、现代圆角、柔和阴影、清晰层级、响应式布局、组件可复用、交互状态完整。
```

在任何前端页面生成任务中，都应优先保证页面干净、专业、舒适、易维护。
