# AuthServer

最小 Gin + MySQL 权限服务实践。

## 当前能力

- Gin HTTP 服务
- MySQL / GORM 连接
- 核心表 AutoMigrate
- 默认环境和身份源初始化
- 健康检查
- 本地登录接口骨架
- LDAP / 内部 SSO provider 预留
- 权限判断接口骨架
- 管理员用户、用户组、角色、权限、授权绑定 API
- 业务线、namespace、cluster 管理 API
- Rancher 同步预览接口
- 审计日志落库
- Vue 管理后台

## 启动

```bash
cp .env.example .env
go mod tidy
go run ./cmd/server
```

本地 Docker MySQL 8.0.46：

```bash
docker compose -f docker-compose.mysql.yml up -d
```

这个 MySQL 会映射到本机 `3000` 端口，`.env.example` 默认 DSN 已按这个端口配置：

```env
MYSQL_DSN=auth:auth@tcp(127.0.0.1:3000)/authserver?charset=utf8mb4&parseTime=True&loc=Local
```

连接测试：

```bash
mysql -h127.0.0.1 -P3000 -uauth -pauth authserver
```

前端开发服务：

```bash
cd web
npm install
npm run dev
```

默认访问：

```text
http://127.0.0.1:5173
```

Vite 会把 `/api` 代理到后端 `http://127.0.0.1:8080`。

如果 MySQL 还没有创建库和账号，可以用 root 执行：

```bash
mysql -uroot -p < docs/mysql-init.sql
```

如果你使用的是 Docker Desktop、OrbStack、Colima 或容器里的 MySQL，应用访问 MySQL 时来源 IP 可能不是 `localhost`，而是类似 `192.168.65.1`。这时 MySQL 账号需要允许远程来源，`docs/mysql-init.sql` 已经包含 `'auth'@'%'`。

如果要自动创建本地管理员账号，在 `.env` 里设置：

```bash
BOOTSTRAP_ADMIN_USERNAME=admin
BOOTSTRAP_ADMIN_PASSWORD='your-password'
```

也可以直接使用环境变量：

```bash
MYSQL_DSN='auth:auth@tcp(127.0.0.1:3306)/authserver?charset=utf8mb4&parseTime=True&loc=Local' \
JWT_SECRET='change-this-secret' \
BOOTSTRAP_ADMIN_PASSWORD='your-password' \
go run ./cmd/server
```

## 发版

发布脚本会在本地完成：

- 构建后端和前端 Docker 镜像
- 读取本地 `.env.server`，生成并应用线上 `ConfigMap`/`Secret`
- 导出 `authserver-images-<version>.tar`
- 上传到服务器
- 导入 RKE2 使用的 `containerd`
- 更新 `authserver` namespace 下的 deployment 并等待 rollout 完成

默认目标与你当前线上环境一致：

- 主机：`root@218.11.5.223`
- 私钥：`/Users/mac/dev/xengineer-cs2.pem`
- 远端目录：`/authserver`
- containerd socket：`/run/k3s/containerd/containerd.sock`

发布脚本默认把仓库根目录的 `.env.server` 当作服务器配置源；脚本会在发版时重新生成并应用 `authserver-config` 和 `authserver-secret`。本地开发用的 `.env` 不会被默认发布。

脚本默认不会拦截测试环境配置。如果你要对正式环境启用更严格的保护，可以显式打开：

```bash
STRICT_DEPLOY_ENV=true bash scripts/release.sh 1.0.5
```

开启后会拦截几类明显错误：

- `APP_ENV != prod`
- `MYSQL_DSN` 指向 `127.0.0.1` 或 `localhost`
- `PUBLIC_BASE_URL`、`SAML_ENTITY_ID`、`SAML_ACS_URL`、`WAYEN_LOGIN_URL`、`WAYEN_TARGET_URL` 指向 `127.0.0.1` 或 `localhost`

执行方式：

```bash
bash scripts/release.sh 1.0.5
```

如果需要覆盖默认值，可以传环境变量：

```bash
ENV_FILE=.env.server \
DEPLOY_HOST=218.11.5.223 \
DEPLOY_KEY=/path/to/key.pem \
REMOTE_DIR=/authserver \
bash scripts/release.sh 1.0.5
```

如果只是想在服务器上重启已部署的服务，可以执行：

```bash
bash scripts/restart-authserver.sh
```

也支持只重启单个 deployment：

```bash
bash scripts/restart-authserver.sh backend
bash scripts/restart-authserver.sh nginx
```

## 接口

| 方法 | 路径 | 说明 |
|---|---|---|
| `GET` | `/healthz` | 服务健康检查 |
| `GET` | `/readyz` | 数据库连接检查 |
| `POST` | `/api/v1/login` | 本地登录 |
| `POST` | `/api/v1/login/ldap` | LDAP 登录预留 |
| `GET` | `/api/v1/login/:provider` | SSO 跳转预留 |
| `GET` | `/api/v1/login/:provider/callback` | SSO 回调预留 |
| `POST` | `/api/v1/authz/check` | 权限判断 |
| `GET` | `/api/v1/users/me` | 当前用户 |
| `GET` | `/api/v1/wayen/login` | 根据 token 邮箱登录 Wayen 并跳转 |
| `GET` | `/api/v1/saml/metadata` | SAML SP metadata |
| `GET` | `/api/v1/login/internal-sso` | 发起 SAML SSO 登录 |
| `POST` | `/api/v1/saml/acs` | SAML ACS 回调，当前仅调试打印 |
| `GET` | `/auth/.well-known/openid-configuration` | OIDC Discovery 配置 |
| `GET` | `/auth/oauth/authorize` | OAuth2 Authorization Code 授权入口，供 Wayne 使用 |
| `POST` | `/auth/oauth/token` | OAuth2 code 换 access token |
| `GET` | `/auth/oauth/userinfo` | OAuth2 bearer token 查询当前用户 |
| `GET` | `/auth/oauth/jwks` | OIDC JWKS 公钥 |

## SAML Metadata

服务会根据 `.env` 和 `certs/sp.crt` 动态生成 SP metadata，不需要手工维护 `certs/spmeta`。
发起 SSO 登录时会实时请求 `.env` 里的 `SAML_IDP_METADATA_URL`，不再读取本地 `certs/idpmeta`。

本地后端如果运行在 `8083`，建议配置：

```env
HTTP_ADDR=:8083
PUBLIC_BASE_URL=http://localhost:8083
SAML_ENTITY_ID=http://localhost:8083/api/v1/saml/metadata
SAML_ACS_URL=http://localhost:8083/api/v1/saml/acs
SAML_IDP_METADATA_URL=http://sso-internal.dev.qiniu.io/saml2/meta
SAML_SP_CERT_FILE=certs/sp.crt
SAML_SP_KEY_FILE=certs/sp.key
```

给 IdP 配置的 SP metadata 地址：

```text
http://localhost:8083/api/v1/saml/metadata
```

调试 SSO 登录入口：

```text
http://localhost:8083/api/v1/login/internal-sso
```

ACS 会处理 IdP 回调：

```text
POST /api/v1/saml/acs
```

收到 IdP 回调后，会把 SAMLResponse 解码；如果 Assertion 被加密，会用 `SAML_SP_KEY_FILE` 对应私钥解密，并在后端控制台打印解密后的摘要和 Assertion：

- Response ID
- InResponseTo
- Issuer
- Destination
- NameID
- Attributes
- 解密后的 Assertion XML

当前 SAML 用户落库规则：

- `eduPersonPrincipalName` 同时作为 `users.username` 和 `users.email`。
- `NameID` 写入 `users.external_id`。
- 首次登录会自动创建 `source=saml` 的用户。
- 不再创建 `identity_providers` / `user_identities` 表记录。
- 登录成功后签发 JWT，写入 `authserver_token` HttpOnly cookie，并带着 token 跳回前端。

当前阶段还没有做：

- SAML 签名验证

## Wayne OAuth2 对接

方案 1 下 AuthServer 作为 OAuth2 Provider，Wayne 作为 OAuth2 Client。Wayne 登录时不再校验 AuthServer 里的 Wayne 密码，也不需要 AuthServer 存储 Wayne 资源组信息；AuthServer 只负责 SAML 登录、签发 token、返回 `name/email/display`，Wayne 收到 userinfo 后会在 Wayne 自己数据库里创建或更新本地用户，Wayne 的资源组和权限仍由 Wayne 自己维护。

AuthServer 端配置：

```env
OAUTH_WAYNE_CLIENT_ID=wayne
OAUTH_WAYNE_CLIENT_SECRET=change-this-wayne-client-secret
OAUTH_WAYNE_REDIRECT_URI=http://127.0.0.1:8080/login/oauth2/oauth2
WAYEN_OAUTH_REF=/portal/namespace/1/app
OAUTH_CODE_TTL_SECONDS=120
```

AuthServer 暴露给 Wayne 的 OAuth2 地址：

```text
GET  /auth/oauth/authorize
POST /auth/oauth/token
GET  /auth/oauth/userinfo
```

OIDC Discovery 里的 endpoint 默认由 `OIDC_ISSUER` 拼接，也可以按 endpoint 单独覆盖。浏览器需要访问 `OIDC_AUTHORIZATION_ENDPOINT`，后端系统通常访问 `OIDC_TOKEN_ENDPOINT`、`OIDC_USERINFO_ENDPOINT` 和 `OIDC_JWKS_URI`。

```env
OIDC_ISSUER=http://auth.example.com/auth
OIDC_AUTHORIZATION_ENDPOINT=http://auth.example.com/auth/oauth/authorize
OIDC_TOKEN_ENDPOINT=http://auth-internal.example.com/auth/oauth/token
OIDC_USERINFO_ENDPOINT=http://auth-internal.example.com/auth/oauth/userinfo
OIDC_JWKS_URI=http://auth-internal.example.com/auth/oauth/jwks
CLOUDDM_TARGET_URL=http://authserver-nginx/internal/clouddm
```

`CLOUDDM_TARGET_URL` 用于 AuthServer 后端请求 CloudDM `/requestJumpUrl`。在 k8s 内建议指向 AuthServer nginx 的内部代理路径，由 nginx 转发到 CloudDM Service，并把 `Host` 固定成 CloudDM 公网入口，确保 CloudDM 生成浏览器可访问的 callback。

Wayne `app.conf` 示例：

```ini
[auth.oauth2]
enabled = true
redirect_url = http://127.0.0.1:8080
client_id = wayne
client_secret = change-this-wayne-client-secret
auth_url = http://218.11.5.223/auth/oauth/authorize
token_url = http://218.11.5.223/auth/oauth/token
api_url = http://218.11.5.223/auth/oauth/userinfo
api_mapping = name:name,email:email,display:display
scopes = profile,email
```

Wayne 会把回调地址拼成：

```text
{redirect_url}/login/oauth2/oauth2
```

因此 `OAUTH_WAYNE_REDIRECT_URI` 必须和 Wayne 实际回调地址完全一致。浏览器访问 Wayne OAuth 登录入口后，如果 AuthServer 还没有登录态，会先跳内部 SAML；SAML 成功后再回到 OAuth authorize，签发 code 给 Wayne。
`WAYEN_OAUTH_REF` 是 AuthServer 发起 Wayne 登录时写入 Wayne `next` 参数的登录完成页，默认 `/portal/namespace/1/app`，对应 Wayne `DemoNamespaceId = 1` 的默认 namespace。不要配置成 `oauth` 或 `/oauth`，否则 Wayne 回调会把它当成前端路由跳到 `/oauth`。

## Wayne 授权代理接口

AuthServer 的 Wayne 授权代理接口不要求调用方传 Wayne user ID。后端会从当前 `authserver_token` 里取 `email`，把它作为 Wayne username 传给 Wayne internal API。

对外接口：

```text
GET    /auth/api/v1/wayne/namespaces
GET    /auth/api/v1/wayne/groups
GET    /auth/api/v1/wayne/users/me/roles
GET    /auth/api/v1/wayne/namespaces/:namespaceid/operator-permissions
GET    /auth/api/v1/wayne/apps/:appid/operator-permissions
PUT    /auth/api/v1/wayne/namespaces/:namespaceid/roles
DELETE /auth/api/v1/wayne/namespaces/:namespaceid/roles
PUT    /auth/api/v1/wayne/apps/:appid/roles
DELETE /auth/api/v1/wayne/apps/:appid/roles
```

示例：

```http
PUT /auth/api/v1/wayne/namespaces/1/roles
Authorization: Bearer <authserver_token>
Content-Type: application/json

{
  "groupIds": [10, 11],
  "replace": false,
  "requestId": "req-001",
  "reason": "grant namespace access"
}
```

AuthServer 转发到 Wayne internal API 时会使用 token email：

```text
PUT /api/v1/internal/namespaces/1/users/<token-email>/roles
```

并覆盖请求体中的 `operatorName` 为 token email，忽略外部传入的 `operatorUserId`。

相关配置：

```env
WAYNE_INTERNAL_API_BASE_URL=http://wayne-backend.demo.svc.cluster.local:8080
WAYNE_SERVICE_NAME=xinfra
WAYNE_SERVICE_API_SECRET_KEY=<wayne-service-secret>
```

Wayne internal API 签名规则：

```text
bodyHash = SHA256_HEX(rawBody)
payload = METHOD + "\n" + URI + "\n" + timestamp + "\n" + nonce + "\n" + bodyHash
signature = HMAC_SHA256_HEX(secret, payload)
X-Wayne-Signature = "sha256=" + signature
```

管理员可查看当前 SAML metadata 配置：

```text
GET /api/v1/admin/saml/metadata/config
```

## 管理员接口

管理员接口需要 `Authorization: Bearer <token>`，并且 token 里的用户必须是 `is_admin=true`。

| 方法 | 路径 | 说明 |
|---|---|---|
| `GET` | `/api/v1/admin/users` | 用户列表 |
| `POST` | `/api/v1/admin/users` | 创建用户，可带本地密码 |
| `PUT` | `/api/v1/admin/users/:id` | 更新用户 |
| `DELETE` | `/api/v1/admin/users/:id` | 删除用户 |
| `GET` | `/api/v1/admin/groups` | 用户组列表 |
| `POST` | `/api/v1/admin/groups` | 创建用户组 |
| `POST` | `/api/v1/admin/groups/:id/members` | 用户加入组 |
| `DELETE` | `/api/v1/admin/groups/:id/members/:userID` | 用户移出组 |
| `GET` | `/api/v1/admin/roles` | 角色列表 |
| `POST` | `/api/v1/admin/roles` | 创建角色 |
| `PUT` | `/api/v1/admin/roles/:id/permissions` | 设置角色权限 |
| `GET` | `/api/v1/admin/permissions` | 权限点列表 |
| `POST` | `/api/v1/admin/permissions` | 创建权限点 |
| `GET` | `/api/v1/admin/business-lines` | 业务线列表 |
| `POST` | `/api/v1/admin/business-lines` | 创建业务线 |
| `GET` | `/api/v1/admin/namespaces` | namespace 列表 |
| `POST` | `/api/v1/admin/namespaces` | 创建 namespace |
| `GET` | `/api/v1/admin/environments` | 环境列表 |
| `GET` | `/api/v1/admin/clusters` | 集群列表 |
| `POST` | `/api/v1/admin/clusters` | 创建集群 |
| `GET` | `/api/v1/admin/wayen/credentials` | Wayen 凭据列表，不返回密码 |
| `POST` | `/api/v1/admin/wayen/credentials` | 按邮箱保存 Wayen 密码 |
| `DELETE` | `/api/v1/admin/wayen/credentials/:id` | 删除 Wayen 凭据 |
| `GET` | `/api/v1/admin/role-bindings` | 授权绑定列表 |
| `POST` | `/api/v1/admin/role-bindings` | 创建授权绑定 |
| `DELETE` | `/api/v1/admin/role-bindings/:id` | 删除授权绑定 |
| `POST` | `/api/v1/admin/rancher/sync` | Rancher 同步预览 |
| `GET` | `/api/v1/admin/rancher/sync-status` | Rancher 同步状态 |
| `GET` | `/api/v1/admin/saml/metadata/config` | SAML metadata 配置 |

## Rancher 对齐流程

当前 Rancher 同步接口是 dry-run，不会真实调用 Rancher API。推荐先按下面流程把本系统数据配齐：

1. 创建环境和集群：`rke2-dev`、`rke2-test`、`rke2-prod`。
2. 创建业务线：例如 `pay`。
3. 创建 namespace：例如 `pay`，归属支付业务线。
4. 创建权限点：例如 `deployment:deploy`、`deployment:rollback`。
5. 创建角色：例如 `prod-deployer`。
6. 给角色绑定权限点。
7. 创建用户组：例如 `pay-prod-deployer`。
8. 把 SSO/LDAP 用户加入用户组。
9. 创建授权绑定：`group + role + cluster + namespace`。
10. 调用 `/api/v1/admin/rancher/sync` 查看将要同步到 Rancher 的映射。

核心授权绑定应该长这样：

```json
{
  "subject_type": "group",
  "subject_id": 1,
  "role_id": 1,
  "scope_type": "cluster_namespace",
  "cluster_id": 3,
  "namespace_id": 1
}
```

这表示某个用户组在 prod 集群的 `pay` namespace 下拥有指定角色。
