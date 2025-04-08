# CoAuth

CoAuth 是一个提供 OAuth 认证服务的项目，支持多种认证流程和功能。本项目使用 Go 语言开发，借助 gRPC 实现高效的服务通信。

## 开始
### 安装依赖工具
在开始使用 CoAuth 之前，你需要安装必要的工具来生成 Protobuf 和 gRPC 代码。运行以下命令进行安装：
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

## 生成API
安装好工具后，你可以使用以下命令生成 API 代码：
```bash
 .\scripts\gen_pb.cmd
```

## 项目结构
```plaintext
.editorconfig
.git/
.gitignore
README.md
api/
  v1/
    coauth.proto
conf/
  conf.yaml
database/
  mysql.sql
go.mod
internal/
  biz/
  bootstrap/
  cmd/
  config/
  data/
  database/
  di/
  entities/
  model/
  server/
  service/
main.go
pkg/
  logger/
scripts/
  build.cmd
  gen_pb.cmd
tests/
  client_test.go
third_party/
  errors/
  google/
  openapi/
  validate/
```
### 主要目录说明
- api ：包含 gRPC 服务的 Protobuf 定义文件。
- conf ：存放项目的配置文件，如 conf.yaml 。
- database ：包含数据库相关的脚本和配置。
- internal ：项目的核心代码，包含业务逻辑、数据访问、服务层等。
- pkg ：存放项目的公共工具包，如日志记录器。
- scripts ：包含用于生成代码和构建项目的脚本。
- tests ：包含项目的单元测试和集成测试代码。
- third_party ：存放第三方依赖和工具。

## 配置文件
项目的配置文件位于 conf/conf.yaml ，以下是一个示例配置：
```yaml
server:
  grpc:
    addr: 0.0.0.0:10101
    timeout: 30
log:
  level: debug
  filename: app.log
db:
  driver: mysql
  source: root:123456@tcp(127.0.0.1:13306)/coauth?charset=utf8mb4&parseTime=True&loc=Local
  log_level: info
  singular_table: true
  prepare_stmt: true
  allow_global_update: false
  max_open_conns: 100
  max_idle_conns: 100
  conn_max_lifetime: 3600
```
### 配置说明
- server ：gRPC 服务器的配置，包括监听地址和超时时间。
- log ：日志配置，包括日志级别和日志文件名称。
- db ：数据库配置，包括数据库驱动、连接源、最大连接数等。

## 运行项目
在完成配置和代码生成后，你可以运行以下命令启动 CoAuth 服务：

### 构建可执行文件
```bash
.\scripts\build.cmd
```

### 运行
```bash
.\bin\coauth.exe -c conf\conf.yaml
```

## 测试
项目提供了单元测试和集成测试代码，你可以使用以下命令运行测试：

```bash
go test -v ./...
 ```

## 贡献
如果你想为 CoAuth 项目做出贡献，请遵循以下步骤：

1. Fork 本仓库。
2. 创建一个新的分支： git checkout -b feature/your-feature 。
3. 提交你的更改： git commit -m 'Add some feature' 。
4. 推送至远程分支： git push origin feature/your-feature 。
5. 提交 Pull Request。
