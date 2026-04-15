# QQBot

基于 napcat 的 AI 对话机器人。

**核心定位**：接收 QQ 消息，通过 ARK（大模型 API）生成回复，支持私聊/群聊、会话记忆、并发处理。

## 技术栈

| 层级 | 选型 | 说明 |
|------|------|------|
| 语言 | Go 1.25 | 高并发，天然适配消息处理 |
| Web 框架 | Gin | 高性能 HTTP 路由 |
| LLM | ARK（字节火山引擎） | 提供对话能力 |
| 消息源 | napcat | QQ 机器人协议适配 |
| 配置 | Viper | YAML 配置管理 |

## 项目结构

```
qqbot/
├── cmd/main.go          # 入口：加载配置、启动工作池、监听 HTTP
├── config/              # 配置加载（LLM、服务器、napcat）
├── internal/
│   ├── hook/            # 中间件：Token 验证、消息解析、@ 判断
│   ├── routes/          # Gin 路由注册
│   ├── sendMsg/         # 消息发送（异步写入 napcat API）
│   └── utils/           # 工具函数
└── models/
    ├── workPool.go      # 工作池：4 worker 并发消费消息
    ├── sessionProfessor.go  # 会话管理：维护上下文 + 调用 LLM
    ├── sessionRegister.go   # 会话注册表
    └── messagesRingBuffer.go # 消息环形缓冲区
```

## 快速启动

```bash
# 1. 克隆项目
git clone <repo>
cd qqbot

# 2. 配置
cp config.example.yaml config.yaml
# 编辑 config.yaml，填入 LLM API Key、napcat 地址等

# 3. 运行
go run cmd/main.go
```

## 配置项说明

```yaml
openaiProvider:       # LLM 配置
  key: ""             # ARK API Key
  model: "gpt-5.2"   # 模型名称
  baseURL: ""         # API 地址

server:
  port: ":8077"       # HTTP 监听端口

napcat:
  apiBaseURL: "http://localhost:3000"   # napcat API 地址
  expectedToken: ""     # 期望的 Token（用于验证）
  authToken: ""         # 认证 Token

judgeAt:
  enable: true         # 是否仅回复 @ 消息（false = 所有消息都回复）
```

## Docker 部署

```bash
docker-compose up -d
```

## 消息处理流程

```
napcat (QQ)  →  POST /  →  Hooks (验证/解析)  →  WorkPool  →  SessionProfessor
                                                                ↓
                                                        ARK (LLM)
                                                                ↓
                                                        回复写入 napcat
```

## 核心特性

- **并发处理**：4 worker 工作池，消息队列缓冲 20 条
- **会话记忆**：SessionProfessor 维护每个用户/群的对话上下文
- **灵活配置**：支持仅回复 @ 消息，或全部消息
- **安全验证**：Token 校验，防止非法请求
