# Gin-Center 后台管理系统

## 项目简介

Gin-Center 是一个现代化、高性能的企业级后台管理系统，基于 Golang Gin 框架精心构建。旨在提供安全、高效、可扩展的企业级管理解决方案。

## 项目特点

- 安全的用户认证机制
- 灵活的角色和权限管理
- 全面的日志追踪和监控
- 高性能、低延迟的系统架构
- 可配置、可扩展的模块化设计

## 技术栈

| 类型 | 技术 | 版本/描述 |
|------|------|-----------|
| 后端框架 | Gin | v1.7+ |
| 数据库 | MySQL | v5.7+ |
| ORM | GORM | v2.0+ |
| 认证 | JWT | 基于 Token 的安全认证 |
| 文档 | Swagger | OpenAPI 规范 |
| 缓存 | Redis | v6.0+ |
| 开发工具 | Air | 热重载开发 |
| 日志 | Zap | 高性能日志库 |
| 配置管理 | Viper | 配置文件解析 |

## 项目结构

```
gin-center/
├── configs/            # 配置文件目录
│   ├── app.yaml       # 应用主配置
│   ├── config/        # 配置管理
│   ├── database/      # 数据库迁移
│   └── env/          # 环境配置
├── docs/              # 文档目录
├── infrastructure/    # 基础设施层
│   ├── bootstrap/     # 应用引导
│   ├── cache/        # 缓存实现
│   ├── container/    # 依赖容器
│   ├── database/     # 数据库连接
│   ├── errors/       # 错误处理
│   ├── repository/   # 数据仓储实现
│   └── zaplogger/    # 日志实现
├── internal/          # 内部核心代码
│   ├── application/  # 应用服务层
│   ├── domain/       # 领域层
│   └── types/        # 类型定义
├── pkg/               # 公共工具包
│   ├── circuitbreaker/ # 熔断器
│   ├── http/         # HTTP工具
│   ├── security/     # 安全工具
│   ├── time/         # 时间工具
│   ├── tracer/       # 链路追踪
│   └── utils/        # 通用工具
├── scripts/           # 运维脚本
└── web/              # Web层
    ├── controller/   # 控制器
    ├── middleware/   # 中间件
    └── routes/       # 路由定义
```

## 接口列表

### 用户认证接口

| 接口名称 | 路径 | 方法 | 说明 | 权限 |
|---------|------|------|------|------|
| 用户登录 | `/login` | POST | 用户身份认证 | 公开 |
| 用户注册 | `/register` | POST | 创建新用户账号 | 公开 |
| 获取个人信息 | `/profile` | GET | 查询当前用户详情 | 登录 |
| 更新个人信息 | `/profile` | PUT | 修改用户基本信息 | 登录 |
| 修改密码 | `/password` | PUT | 更新用户密码 | 登录 |

### 管理员接口

| 接口名称 | 路径 | 方法 | 说明 | 权限 |
|---------|------|------|------|------|
| 管理员登录 | `/admin/login` | POST | 管理员身份认证 | 公开 |
| 管理员注册 | `/admin/register` | POST | 创建管理员账号 | 管理员 |
| 获取管理员信息 | `/admin/info` | GET | 查询管理员详情 | 管理员 |
| 更新管理员信息 | `/admin` | PUT | 修改管理员信息 | 管理员 |
| 管理员列表 | `/admin/list` | GET | 分页获取管理员列表 | 超级管理员 |

### 系统管理接口

| 接口名称 | 路径 | 方法 | 说明 | 权限 |
|---------|------|------|------|------|
| 系统信息 | `/system/info` | GET | 获取系统基本信息 | 管理员 |
| 系统配置 | `/system/config` | GET | 获取系统配置详情 | 管理员 |
| 更新系统配置 | `/system/config` | PUT | 修改系统配置 | 超级管理员 |

## 环境要求

| 依赖 | 最低版本 | 推荐版本 |
|------|----------|----------|
| Go | 1.16+ | 1.20+ |
| MySQL | 5.7 | 8.0 |
| Redis | 6.0 | 6.2+ |

## 快速开始

### 1. 克隆项目
```bash
git clone https://github.com/your-username/gin-center.git
cd gin-center
```

### 2. 安装依赖
```bash
go mod download
```

### 3. 配置环境
```bash
# 复制配置文件
cp configs/env/dev.yaml configs/app.yaml

# 编辑配置文件
vim configs/app.yaml
```

### 4. 启动项目

#### 开发模式
```bash
# 安装 Air
go install github.com/cosmtrek/air@latest

# 启动热重载
air
```

#### 生产模式
```bash
# 编译
go build -o gin-center

# 运行
./gin-center
```

## 接口文档

- Swagger文档：`http://localhost:8080/swagger/index.html`

## 安全特性

1. JWT token认证
2. RBAC权限控制
3. 请求频率限制
4. 参数校验
5. 敏感信息脱敏
6. 安全日志记录
7. 密码加密存储
8. 熔断保护机制
9. 链路追踪

## 性能优化

- 多级缓存：Redis + 本地缓存
- 数据库连接池管理
- 请求限流与熔断
- 异步任务处理
- 链路追踪优化
- 日志分级处理
- 资源复用机制

## 贡献指南

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交代码 (`git commit -m '新增：某某特性'`)
4. 推送分支 (`git push origin feature/amazing-feature`)
5. 提交 Pull Request

## 许可证

本项目基于 MIT 许可证开源。详情请查看 `LICENSE` 文件。


**感谢使用 Gin-Center！**