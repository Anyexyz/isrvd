# Isrvd Code Wiki

## 项目概述

Isrvd 是一个轻量级 Web 服务器管理工具，基于 **Go + Vue 3** 构建，提供文件管理、在线编辑、Docker/Swarm 管理、APISIX 管理和实时终端等功能。

### 核心特性

- **系统概览** - CPU、内存、硬盘、网络实时监控
- **文件管理** - 浏览、上传、下载、在线编辑、压缩解压、权限修改
- **Docker 服务** - 容器、镜像、网络、卷的完整管理，支持终端和实时统计
- **Docker Swarm** - 服务、节点、任务的完整管理
- **APISIX 管理** - 路由、Consumer、Upstream、IP 白名单管理
- **Web 终端** - xterm.js 实时 Shell 交互
- **成员管理** - 多用户支持，用户隔离、独立家目录、权限控制
- **AI 助手** - 内置 Page-Agent，支持自然语言操作与页面感知
- **移动端适配** - 响应式布局

---

## 项目架构

### 整体架构图

```
┌─────────────────────────────────────────────────────────────┐
│                         前端 (Vue 3)                        │
│  ┌──────────────┐  ┌──────────────┐  ┌───────────────────┐ │
│  │  系统概览    │  │  文件管理    │  │  Docker/Swarm    │ │
│  │  APISIX管理  │  │  Web终端    │  │  成员管理        │ │
│  └──────────────┘  └──────────────┘  └───────────────────┘ │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                   后端 API 层 (Go + Gin)                     │
│  ┌───────────────────────────────────────────────────────┐ │
│  │              路由与认证中间件                          │ │
│  └───────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                   业务服务层                                 │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌────────────┐ │
│  │  Docker  │  │  Swarm   │  │  APISIX  │  │  System    │ │
│  │  Compose │  │  Filer   │  │  Agent   │  │  Shell     │ │
│  └──────────┘  └──────────┘  └──────────┘  └────────────┘ │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                   底层封装层 (pkgs)                          │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌────────────┐ │
│  │docker    │  │swarm     │  │apisix    │  │compose     │ │
│  │archive   │  │          │  │          │  │            │ │
│  └──────────┘  └──────────┘  └──────────┘  └────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

### 目录结构

```
/workspace/
├── cmd/server/              # 程序入口
│   └── main.go
├── config/                  # 配置模块
│   ├── types.go             # 配置类型定义
│   ├── load.go              # 配置加载
│   └── save.go              # 配置保存
├── internal/                # 内部模块
│   ├── helper/              # 辅助工具
│   │   ├── response.go      # API响应处理
│   │   └── websocket.go     # WebSocket工具
│   ├── registry/            # 服务注册
│   │   ├── apisix.go
│   │   ├── compose.go
│   │   ├── docker.go
│   │   └── registry.go
│   ├── server/              # Web服务器
│   │   ├── app.go           # 应用初始化与路由
│   │   ├── auth.go          # 认证中间件
│   │   └── handler_*.go     # 各模块处理器
│   └── service/             # 业务服务层
│       ├── apisix/
│       ├── compose/
│       ├── docker/
│       ├── swarm/
│       └── system/
├── pkgs/                    # 公共包（可复用）
│   ├── apisix/              # APISIX 客户端封装
│   ├── archive/             # 压缩/解压工具
│   ├── compose/             # Docker Compose 封装
│   ├── docker/              # Docker 客户端封装
│   └── swarm/               # Docker Swarm 封装
├── webview/                 # 前端 (Vue 3)
│   ├── src/
│   │   ├── component/       # 组件库
│   │   ├── views/           # 页面视图
│   │   │   ├── apisix/
│   │   │   ├── docker/
│   │   │   ├── filer/
│   │   │   ├── swarm/
│   │   │   └── system/
│   │   ├── router/
│   │   ├── service/
│   │   ├── helper/
│   │   ├── store/
│   │   └── app.vue
│   └── package.json
├── build/                   # 构建相关
│   ├── docker/              # Docker 构建
│   └── marketplace/         # 应用市场构建
├── public/                  # 公共资源（嵌入）
├── go.mod                   # Go 依赖管理
└── README.md
```

---

## 主要模块说明

### 1. 入口模块 (`cmd/server/`)

**文件**: [main.go](file:///workspace/cmd/server/main.go)

**职责**:
- 初始化配置加载
- 初始化服务注册
- 启动 Web 服务器

### 2. 配置模块 (`config/`)

**核心文件**: [types.go](file:///workspace/config/types.go), [load.go](file:///workspace/config/load.go), [save.go](file:///workspace/config/save.go)

**主要配置结构**:
- `Config` - 整体配置
- `Server` - 服务器配置（监听地址、JWT密钥等）
- `AgentConfig` - LLM代理配置
- `ApisixConfig` - APISIX 配置
- `DockerConfig` - Docker 配置
- `MemberConfig` - 成员配置

### 3. 服务注册模块 (`internal/registry/`)

**文件**: [registry.go](file:///workspace/internal/registry/registry.go), [apisix.go](file:///workspace/internal/registry/apisix.go), [docker.go](file:///workspace/internal/registry/docker.go), [compose.go](file:///workspace/internal/registry/compose.go)

**职责**:
- 初始化并注册底层服务实例
- 解耦业务层与底层封装
- 处理服务初始化失败的情况

### 4. Web服务器模块 (`internal/server/`)

**核心文件**: [app.go](file:///workspace/internal/server/app.go), [auth.go](file:///workspace/internal/server/auth.go)

**App 结构**:
```go
type App struct {
    *gin.Engine
    dockerSvc   *svcDocker.Service
    swarmSvc    *svcSwarm.Service
    apisixSvc   *svcApisix.Service
    composeSvc  *svcCompose.DeployService
    systemSvc   *svcSystem.Service
    settingsSvc *svcSystem.SettingsService
    memberSvc   *svcSystem.MemberService
}
```

**API 路由分组**:
- `/api/login` - 登录
- `/api/filer/*` - 文件管理
- `/api/docker/*` - Docker 管理
- `/api/swarm/*` - Swarm 管理
- `/api/apisix/*` - APISIX 管理
- `/api/compose/*` - Compose 部署
- `/api/system/*` - 系统管理
- `/api/agent/*` - AI 代理
- `/ws/*` - WebSocket (Shell/Exec)

### 5. 业务服务层 (`internal/service/`)

#### Docker 服务 (`internal/service/docker/`)
- 容器管理
- 镜像管理
- 网络管理
- 卷管理
- 镜像仓库管理

#### Swarm 服务 (`internal/service/swarm/`)
- 节点管理
- 服务管理
- 任务管理

#### APISIX 服务 (`internal/service/apisix/`)
- 路由管理
- Consumer 管理
- 插件管理
- IP 白名单管理

#### Compose 服务 (`internal/service/compose/`)
- Docker Compose 部署
- Swarm Stack 部署

#### System 服务 (`internal/service/system/`)
- 系统监控
- 成员管理
- 配置管理

### 6. 公共包层 (`pkgs/`)

#### Docker 封装 (`pkgs/docker/`)
**核心文件**: [client.go](file:///workspace/pkgs/docker/client.go)

```go
type DockerService struct {
    client *client.Client
    config *DockerConfig
}
```

**功能**:
- 封装 Docker SDK 客户端
- 提供容器、镜像、网络、卷等操作
- 统一错误处理

#### APISIX 封装 (`pkgs/apisix/`)
- 封装 APISIX Admin API
- 提供路由、Consumer、Upstream 等操作

#### Compose 封装 (`pkgs/compose/`)
- 使用 compose-spec 库解析 Compose 文件
- 提供部署功能

#### Archive 封装 (`pkgs/archive/`)
- 提供 zip 压缩/解压功能

### 7. 辅助工具 (`internal/helper/`)

**响应处理**: [response.go](file:///workspace/internal/helper/response.go)
```go
type APIResponse struct {
    Success bool   `json:"success"`
    Message string `json:"message,omitempty"`
    Payload any    `json:"payload,omitempty"`
}
```

**WebSocket 工具**: [websocket.go](file:///workspace/internal/helper/websocket.go)

### 8. 前端模块 (`webview/`)

**框架**: Vue 3 + Vite + TypeScript + Tailwind CSS

**核心依赖**:
- Vue Router - 路由管理
- Axios - HTTP 客户端
- Xterm.js - Web 终端
- CodeMirror - 代码编辑器
- Chart.js - 图表展示
- Page-Agent - AI 助手

**页面路由**: 见 [webview/src/router/index.ts](file:///workspace/webview/src/router/index.ts)

---

## 关键类与函数说明

### 后端核心

#### 1. `NewApp()` 函数

**位置**: [app.go](file:///workspace/internal/server/app.go#L30-L70)

**职责**: 初始化整个应用
- 创建 Gin 引擎
- 初始化所有业务服务
- 注册路由
- 启动 HTTP 服务器

#### 2. `AuthMiddleware()` 中间件

**位置**: [auth.go](file:///workspace/internal/server/auth.go#L15-L20)

**职责**: 认证中间件工厂
- 支持 JWT 认证
- 支持代理 Header 认证（用于内网环境）

#### 3. `RespondSuccess()` / `RespondError()`

**位置**: [response.go](file:///workspace/internal/helper/response.go#L17-L31)

**职责**: 统一 API 响应格式

#### 4. `DockerService` 类

**位置**: [pkgs/docker/client.go](file:///workspace/pkgs/docker/client.go#L13-L50)

**主要方法**:
- `NewDockerService(cfg *DockerConfig)` - 创建 Docker 服务
- `GetInfo(ctx)` - 获取 Docker 概览
- `ContainerList()`, `ContainerCreate()`, 等 - 容器操作
- `ImageList()`, `ImagePull()`, 等 - 镜像操作
- 网络、卷、仓库等操作

### 前端核心

#### 1. `createApp()` 入口

**位置**: [webview/src/main.ts](file:///workspace/webview/src/main.ts#L13)

**职责**: 创建并挂载 Vue 应用

#### 2. API 服务

**位置**: [webview/src/service/](file:///workspace/webview/src/service/)

**职责**: 封装后端 API 调用

---

## 技术栈与依赖

### 后端技术栈

| 技术/依赖 | 版本 | 用途 |
|----------|------|------|
| Go | 1.25.0 | 开发语言 |
| Gin | 1.12.0 | Web 框架 |
| Docker SDK | 24.0.9+ | Docker 客户端 |
| compose-spec | 2.10.2 | Compose 文件解析 |
| Gorilla WebSocket | 1.5.3 | WebSocket 支持 |
| JWT | 5.3.1 | 认证令牌 |
| gopsutil | 3.24.5 | 系统监控 |
| pango | 0.12.0 | 通用工具库 |

### 前端技术栈

| 技术/依赖 | 版本 | 用途 |
|----------|------|------|
| Vue | 3.5.18 | 前端框架 |
| Vite | 7.1.2 | 构建工具 |
| Vue Router | 4.6.4 | 路由管理 |
| Tailwind CSS | 3.4.19 | UI 框架 |
| Axios | 1.11.0 | HTTP 客户端 |
| Xterm.js | 5.5.0 | Web 终端 |
| CodeMirror | 6.0.2 | 代码编辑器 |
| Chart.js | 4.5.1 | 图表 |
| Page-Agent | 1.8.0 | AI 助手 |

---

## 运行方式

### 1. Docker 部署

**镜像**: `rehiy/isrvd:latest`

```bash
docker run -d \
  --name isrvd \
  -p 8080:8080 \
  -p 9080:9080 \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v /mnt/isrvd:/data \
  rehiy/isrvd:latest
```

**Docker Compose**:
```yaml
services:
  isrvd:
    image: rehiy/isrvd:latest
    container_name: isrvd
    restart: unless-stopped
    ports:
      - "8080:8080"
      - "9080:9080"
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
      - /mnt/isrvd:/data
```

**端口说明**:
- `8080` - Isrvd Web 管理界面
- `9080` - APISIX HTTP 代理
- `9180` - APISIX Admin API（不对外暴露）

**数据目录结构**:
```
/data/
├── conf/       # 配置文件
├── etcd/       # etcd 数据
├── storage/    # 文件管理根目录
└── container/  # 容器数据
```

### 2. 二进制部署

从 [Releases](https://github.com/rehiy/isrvd/releases) 下载对应平台的二进制文件。

**配置示例** (`config.yml`):
```yaml
server:
  listenAddr: ":8080"
  jwtSecret: "your-secret-key"
  rootDirectory: "."

members:
  - username: admin
    password: admin123
    homeDirectory: public
    allowTerminal: true
```

**运行**:
```bash
./isrvd
```

### 3. 开发环境

**前置要求**:
- Go 1.21+
- Node.js 16+

**编译脚本**:
```bash
./build.sh
```

**前端开发**:
```bash
cd webview
npm install
npm run dev
```

---

## 工作流程示例

### 1. 容器创建流程

```
前端 (Vue)
    │
    ▼
Docker 组件调用 API
    │
    ▼
POST /api/docker/container/create
    │
    ▼
handler_docker.go: dockerCreateContainer()
    │
    ▼
internal/service/docker/*: 业务逻辑
    │
    ▼
pkgs/docker/*: 底层封装
    │
    ▼
Docker SDK → Docker Daemon
    │
    ▼
返回响应 (APIResponse)
```

### 2. WebSocket 终端流程

```
前端连接
    │
    ▼
GET /ws/shell?token=...
    │
    ▼
AuthMiddleware() 验证
    │
    ▼
shellWebSocket() 处理
    │
    ▼
启动 PTY 进程
    │
    ▼
双向数据流
    │
    └── Xterm.js ←→ Gorilla WebSocket ←→ PTY
```

---

## 认证与安全

### 认证方式

1. **JWT 认证** (默认)
   - 登录获取 Token
   - Header: `Authorization: Bearer <token>`
   - WebSocket 可通过 query 参数 `?token=<token>`

2. **代理 Header 认证** (可选)
   - 配置 `server.proxyHeaderName`
   - 从 HTTP Header 获取用户名
   - 适用于内网环境与其他认证系统集成

### 安全建议

- 生产环境必须修改 `jwtSecret`
- 首次运行后通过 Docker logs 获取随机密码并修改
- 使用 HTTPS 反向代理
- 限制 Docker socket 的访问权限

---

## 扩展与开发

### 添加新 API 端点

1. 在 `internal/server/handler_*.go` 添加处理器函数
2. 在 `app.go` 的 `setupRouter()` 中注册路由
3. 如需要，在 `internal/service/` 添加业务逻辑
4. 前端在 `webview/src/service/api.ts` 添加 API 调用

### 添加新页面

1. 在 `webview/src/views/` 创建页面组件
2. 在 `webview/src/router/index.ts` 添加路由
3. 在导航组件中添加菜单入口

---

## 许可证

GPL-3.0
