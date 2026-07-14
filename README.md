# XINFRA - 统一运维入口平台

XINFRA 是一个统一运维入口平台，整合多个运维子系统，提供统一的登录认证、子系统导航、审计日志和 Ansible 运维调度功能。

## 项目特性

- 🔐 **统一认证**：支持 LDAP + SAML + OAuth 2.0 认证
- 🎯 **子系统导航**：卡片式展示 6 个运维子系统，支持 SSO 跳转
- 📊 **审计面板**：登录审计 + 运维操作审计，完整的操作追溯
- 🚀 **Ansible 调度**：Playbook 执行、实时日志推送
- 🔒 **安全可靠**：JWT 认证、HTTPS 传输、敏感配置 K8s Secret 管理

## 技术栈

| 层级 | 技术 |
|------|------|
| 后端 | Go 1.21+ / Gin |
| 前端 | Vue 3.3+ / Element Plus / Vite |
| 数据库 | MySQL 8.0 |
| 缓存 | Redis |
| 部署 | Kubernetes + Nginx |

## 项目结构

```
xinfra/
├── frontend/              # 前端 Vue 3 项目
│   ├── src/
│   │   ├── api/           # API 请求封装
│   │   ├── components/    # 可复用组件
│   │   ├── composables/   # 组合式函数
│   │   ├── layouts/       # 页面布局
│   │   ├── router/        # 路由配置
│   │   ├── stores/        # Pinia 状态管理
│   │   ├── views/         # 页面视图
│   │   └── utils/         # 工具函数
│   └── ...
├── server/                # 后端 Go 项目
│   ├── cmd/               # 程序入口
│   ├── internal/          # 内部包（分层架构）
│   ├── pkg/               # 可复用公共包
│   └── migrations/        # 数据库迁移
├── deploy/                # 部署配置
│   ├── docker/            # Docker 构建配置
│   └── k8s/               # Kubernetes 配置
├── docs/                  # 项目文档
└── ...
```

## 快速开始

### 环境要求

- Go 1.21+
- Node.js 18+
- MySQL 8.0
- Redis 7.0+

### 本地开发

1. 克隆项目
   ```bash
   git clone <repository-url>
   cd xinfra
   ```

2. 启动开发环境
   ```bash
   make dev-up
   ```

3. 启动前端开发服务器
   ```bash
   make run-frontend
   ```

4. 启动后端开发服务器
   ```bash
   make run-server
   ```

### 构建部署

1. 构建项目
   ```bash
   make build
   ```

2. 构建 Docker 镜像
   ```bash
   make docker
   ```

3. 部署到 Kubernetes
   ```bash
   kubectl apply -f deploy/k8s/
   ```

## 开发规范

### Git 提交规范

使用语义化提交信息：
- `feat:` 新功能
- `fix:` 修复 bug
- `docs:` 文档更新
- `style:` 代码格式调整
- `refactor:` 重构
- `test:` 测试相关
- `chore:` 构建/工具相关

### 代码规范

- Go 代码遵循 `gofmt` 格式
- Vue/TypeScript 代码遵循 ESLint 规范
- 提交前请确保代码通过所有测试

## 相关文档

- [MVP 方案文档](docs/mvp-skeleton-plan.md)
- [API 接口文档](docs/api.md)（待补充）
- [部署文档](docs/deployment.md)（待补充）

## 许可证

[待添加]
