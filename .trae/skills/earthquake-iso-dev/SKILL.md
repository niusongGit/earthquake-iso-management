---
name: "earthquake-iso-dev"
description: "开发地震国际标准管理平台（Go+Gin+GORM后端，Vue3+Element-Plus前端）。当用户在此项目上进行编码、调试、添加功能或修改任何部分时触发调用。"
---

# 地震国际标准管理平台 - 开发技能

## 项目概述

地震国际标准管理平台，用于管理地震相关的国际ISO标准文档。提供面向公众的文档浏览/搜索界面，以及后台管理员的文档增删改查管理功能。

## 技术栈

- **后端**: Go (Gin框架)、MySQL (GORM ORM)、JWT认证
- **前端**: Vue3 + Element Plus，兼容手机端和PC端
- **日志**: go.uber.org/zap，按天分割日志文件，保留最新20天
- **部署**: Go embed 将 `staticweb/` 嵌入二进制文件，实现单文件部署

## 目录结构

```
earthquake-iso-management/
├── config.json                  # 项目配置文件（数据库、默认管理员账号、JWT密钥等）
├── main.go                      # 程序入口：初始化配置、数据库、路由，启动服务，embed静态文件
├── go.mod
├── go.sum
├── internal/
│   ├── config/
│   │   └── config.go            # 从config.json加载配置
│   ├── logger/
│   │   └── logger.go            # zap日志初始化 + lumberjack按天分割
│   ├── database/
│   │   └── database.go          # MySQL初始化、自动迁移、默认管理员创建
│   ├── model/
│   │   ├── admin.go             # 管理员模型（GORM）
│   │   └── document.go          # 国际标准文档模型 + 请求/响应结构体
│   ├── handler/
│   │   ├── admin_auth.go        # 管理员登录处理器
│   │   ├── admin_document.go    # 后台文档增删改查处理器
│   │   └── front_document.go    # 前台文档查询/详情/预览/下载处理器
│   ├── middleware/
│   │   └── auth.go              # JWT生成与认证中间件
│   ├── response/
│   │   └── response.go          # 统一JSON响应格式
│   ├── router/
│   │   └── router.go            # 路由注册 + SPA静态文件服务
│   └── service/
│       ├── admin.go             # 管理员业务逻辑
│       └── document.go          # 文档业务逻辑
├── uploads/                     # PDF附件上传目录
├── web/                         # Vue3前端源码
│   ├── package.json
│   ├── vite.config.js
│   ├── index.html
│   └── src/
│       ├── main.js
│       ├── App.vue
│       ├── router/
│       │   └── index.js         # 前端路由（前台+后台）
│       ├── api/
│       │   ├── request.js       # axios封装（拦截器、token、401处理）
│       │   ├── admin.js         # 后台API
│       │   └── front.js         # 前台API
│       ├── views/
│       │   ├── front/
│       │   │   ├── Home.vue     # 主页（文档列表+搜索+筛选+排序+管理入口）
│       │   │   └── Detail.vue   # 文档详情页（字段独立展示+PDF预览/下载）
│       │   └── admin/
│       │       ├── Login.vue    # 管理员登录（已登录自动跳转+错误提示）
│       │       ├── Layout.vue   # 后台布局（侧边栏+标题）
│       │       └── Document.vue # 文档管理（PC表格+手机卡片+表单弹窗）
│       └── assets/
│           └── logo.png
├── staticweb/                   # Vue编译后的静态文件（Go embed嵌入目标）
├── logs/                        # 日志文件目录（按天分割）
└── html-dome/                   # HTML原型（参考用，不参与编译）
    ├── index.html
    └── logo.png
```

## 数据库配置（config.json）

```json
{
  "mysql": {
    "address": "localhost",
    "port": 3306,
    "db-name": "iso_standard_document",
    "username": "root",
    "password": "123456",
    "config": "charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai",
    "max-idle-conns": 10,
    "max-open-conns": 100
  },
  "admin": {
    "username": "admin",
    "password": "admin123"
  },
  "jwt": {
    "secret": "earthquake-iso-secret-key",
    "expire-hours": 24
  },
  "server": {
    "port": 8080
  }
}
```

默认管理员账号在config.json中配置，首次启动时自动插入到admin数据表中。

## 数据模型

### admin 管理员表
| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint (主键) | 自增 |
| username | string (唯一) | 管理员用户名 |
| password | string | Bcrypt加密后的密码 |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |

### document 国际标准文档表
| 字段 | 数据库列名 | 类型 | 说明 |
|------|-----------|------|------|
| id | id | uint (主键) | 自增 |
| iso_code | iso_code | string(50) | 国际标准编号，如ISOxxxx（去除前后空格） |
| name | name | string(255) | 国际标准名称 |
| type | type | string(20) | 标准类型：IS、TS、PAS、TR、IWA、Guides |
| standard_belongs_to | standard_belongs_to | string(100) | 标准所属：ISO/TC（输入框默认值） |
| belongs_to | belongs_to | string(50) | 所属：SC、WG（单选枚举） |
| summary | summary | text | 摘要 |
| scope | scope | text | 范围 |
| publish_date | publish_date | date | 发布日期（YYYY-MM-DD） |
| first_publish_code | first_publish_code | string(100) | 首次发布编号 |
| current_stage | current_stage | string(20) | 当前阶段：PWI、NP、WD、CD、DIS、FDIS、IS |
| earthquake_relevance | earthquake_relevance | int(默认1) | 地震相关度：1-5星级 |
| created_at | created_at | datetime | 系统记录创建时间 |
| updated_at | updated_at | datetime | 更新时间 |
| attachment | attachment | string(500) | PDF附件路径 |

## API设计

### 前台API（无需认证）
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/front/documents | 文档列表（分页、搜索、筛选、排序） |
| GET | /api/front/documents/:id | 文档详情 |
| GET | /api/front/documents/:id/preview | 在线预览PDF（Content-Disposition: inline） |
| GET | /api/front/documents/:id/download | 下载PDF附件（Content-Disposition: attachment） |

列表查询参数：`keyword`（关键词）、`stage`（当前阶段）、`sort`（createTime/publishDate）、`page`、`pageSize`

### 后台管理API
| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | /api/admin/login | 管理员登录，返回JWT令牌 | 否 |
| POST | /api/admin/documents | 创建文档（multipart/form-data，含PDF上传） | 是 |
| GET | /api/admin/documents | 文档列表（分页、搜索、筛选） | 是 |
| GET | /api/admin/documents/:id | 文档详情 | 是 |
| PUT | /api/admin/documents/:id | 更新文档（multipart/form-data） | 是 |
| DELETE | /api/admin/documents/:id | 删除文档（含附件文件） | 是 |

## 前端需求

### 前台页面（无需登录）
1. **主页**：文档卡片列表，做好分页功能
2. **搜索栏**：页面顶部，关键词搜索 + 当前阶段下拉筛选（PWI/NP/WD/CD/DIS/FDIS/IS）
3. **排序选项卡**：列表上方，"创建时间"（降序，默认选中）和"发布日期"（降序）
4. **文档详情页**：所有字段独立展示，PDF附件提供在线预览（iframe）和下载功能
5. **后台入口**：header右侧"管理"按钮，已登录直接跳后台，未登录跳登录页

### 后台管理页面
1. **登录页**：用户名/密码表单，已登录自动跳转后台，登录失败显示红色错误提示
2. **文档管理**：
   - PC端：表格展示，支持增删改查
   - 手机端：卡片列表展示（与前台风格一致）
   - 新增/编辑弹窗：label-position="top"，手机端宽度适配屏幕

### 设计参考
- HTML原型位于 `html-dome/index.html`
- Vue3 + Element Plus 实现必须与HTML原型的布局、色调和响应式行为保持一致
- 原型色彩变量：
  - `--primary-color: #1e3a8a`（深蓝色，代表权威）
  - `--secondary-color: #2563eb`（科技蓝）
  - `--bg-color: #f8fafc`
  - `--star-color: #f59e0b`
  - `--border-color: #e2e8f0`
  - `--text-muted: #64748b`

### 响应式设计
- 手机端（<768px）：单列布局，卡片式文档列表，弹窗宽度 `calc(100vw - 20px)`，表单label和输入控件分行显示
- PC端（>=768px）：宽屏布局，详情页左右分栏（信息+PDF预览），搜索栏水平排列

## 前端关键实现细节

1. **el-dialog样式覆盖**：el-dialog默认传送到body，scoped样式无效，需用非scoped的style块通过`.doc-dialog`类名选择器覆盖
2. **401拦截器**：request.js中401拦截需排除登录接口（`/admin/login`），否则登录失败会触发页面刷新跳转
3. **登录按钮**：需设置 `native-type="button"` 防止触发表单默认提交行为
4. **PDF预览**：使用独立preview接口（Content-Disposition: inline），iframe加载；下载使用download接口（Content-Disposition: attachment），点击按钮时window.open
5. **文档列表元数据布局**：标准所属+发布日期一行，首发编号+所属一行，地震相关度单独一行

## 日志要求

- 使用 `go.uber.org/zap` 进行结构化日志记录
- 使用 `gopkg.in/natefinch/lumberjack.v2` 实现按天分割日志
- 日志文件存放在 `logs/` 目录
- 仅保留最新20天的日志文件
- 日志格式：生产环境JSON格式，开发环境控制台友好格式
- 配置示例：
  ```go
  &lumberjack.Logger{
      Filename:   "logs/app.log",
      MaxSize:    100, // MB
      MaxBackups: 3,
      MaxAge:     20,  // 天 - 保留最新20天
      Compress:   true,
  }
  ```

## Go Embed部署

```go
//go:embed staticweb
var staticFiles embed.FS
```

`staticweb/` 目录存放编译后的Vue3前端文件。Go服务器提供静态文件服务，并处理SPA路由（未匹配路由回退到index.html）。

## 开发规范

1. **ISO编号输入**：保存前必须去除前后空格
2. **标准所属字段**：输入框默认值为"ISO/TC"
3. **所属字段**：单选枚举SC/WG，下拉选择
4. **地震相关度**：1-5星级评分，显示为★/☆
5. **PDF处理**：通过multipart表单上传，存储在 `uploads/` 目录，提供预览（inline）和下载（attachment）
6. **分页**：默认每页10条
7. **错误处理**：统一的JSON错误响应格式 `{"code": int, "message": string}`
8. **密码安全**：管理员密码使用bcrypt加密存储
9. **JWT认证**：基于令牌的认证，过期时间在config.json中配置
10. **前端构建**：在 `web/` 目录执行 `npm run build`，输出到 `staticweb/`
11. **GORM comment**：所有模型字段需添加GORM comment标签，使用中文注释
12. **form标签**：请求结构体需同时添加 `json` 和 `form` 标签，确保Gin能正确绑定multipart/form-data请求
