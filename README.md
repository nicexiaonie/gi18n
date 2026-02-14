# gi18n

[![Go Reference](https://pkg.go.dev/badge/github.com/nicexiaonie/gi18n.svg)](https://pkg.go.dev/github.com/nicexiaonie/gi18n)
[![Go Report Card](https://goreportcard.com/badge/github.com/nicexiaonie/gi18n)](https://goreportcard.com/report/github.com/nicexiaonie/gi18n)

简单易用的 Go 国际化库，基于 [go-i18n](https://github.com/nicksnyder/go-i18n) 封装。

**一个 `T()` 函数覆盖所有翻译场景。**

## 特性

- **极简 API** — 一个 `T()` + Option 组合，替代记忆十几个方法
- **零配置** — 内置 6 种语言通用词条，开箱即用
- **多格式** — JSON / YAML / TOML 全部支持
- **嵌套展平** — 支持 `common.confirm` 风格的嵌套 key
- **可观测** — MissHandler 回调 + Logger 接口，缺失翻译不再静默
- **HTTP 中间件** — 内置标准库中间件，自动检测语言
- **线程安全** — 全部方法并发安全

## 安装

```bash
go get github.com/nicexiaonie/gi18n
```

## 30 秒快速开始

```go
package main

import (
    "fmt"
    "github.com/nicexiaonie/gi18n"
)

func main() {
    // 内置语言包已自动加载，直接使用
    gi18n.SetLang("zh-CN")
    fmt.Println(gi18n.T("confirm"))  // 确定
    fmt.Println(gi18n.T("cancel"))   // 取消

    // 切换语言
    fmt.Println(gi18n.T("confirm", gi18n.WithLang("en")))  // OK

    // 带参数
    fmt.Println(gi18n.T("greeting", gi18n.WithData("Name", "张三")))  // 你好，张三！

    // 复数
    fmt.Println(gi18n.T("items", gi18n.WithCount(5)))  // 5 个项目
}
```

## 核心 API

只需记住 **1 个翻译函数 + 5 个选项**：

```go
gi18n.T(id string, opts ...Option) string
```

| 选项 | 说明 | 示例 |
|------|------|------|
| `WithLang(lang)` | 指定语言 | `T("hi", WithLang("zh-CN"))` |
| `WithData(kv...)` | 模板参数 (key-value) | `T("hi", WithData("Name", "张三"))` |
| `WithMap(m)` | 模板参数 (map) | `T("hi", WithMap(data))` |
| `WithCount(n)` | 复数 | `T("items", WithCount(5))` |
| `WithContext(ctx)` | 从 Context 获取语言 | `T("hi", WithContext(ctx))` |

选项可自由组合：

```go
// 指定语言 + 参数
gi18n.T("greeting", gi18n.WithLang("en"), gi18n.WithData("Name", "Alice"))

// 指定语言 + 复数
gi18n.T("items", gi18n.WithLang("en"), gi18n.WithCount(3))

// Context + 参数
gi18n.T("greeting", gi18n.WithContext(ctx), gi18n.WithData("Name", "张三"))
```

### 语言管理

| 函数 | 说明 |
|------|------|
| `SetLang(lang)` | 设置当前语言 |
| `GetLang()` | 获取当前语言 |
| `Languages()` | 获取已加载的语言列表 |
| `SetFallbackLang(lang)` | 设置回退语言 |

### 加载方法

| 函数 | 说明 |
|------|------|
| `Load(dir)` | 从目录加载 |
| `LoadFS(fs, root)` | 从 embed.FS 加载 |
| `LoadContent(lang, format, data)` | 从字节内容加载 |
| `LoadMessages(lang, messages)` | 从 map 直接加载 |

## 加载自定义语言包

### 从目录加载

```go
// locales/en.json, locales/zh-CN.json, ...
gi18n.Load("./locales")
```

### 从 embed.FS 加载

```go
//go:embed locales/*
var localesFS embed.FS

gi18n.LoadFS(localesFS, "locales")
```

### 从内容加载

```go
data := []byte(`{"hello": "你好", "world": "世界"}`)
gi18n.LoadContent("zh-CN", "json", data)
```

### 直接注册消息

```go
gi18n.LoadMessages("zh-CN", map[string]string{
    "hello": "你好",
    "world": "世界",
})
```

## 语言包格式

### 简化格式

```json
{
  "hello": "你好",
  "greeting": "你好，{{.Name}}！"
}
```

### 嵌套格式（自动展平）

```json
{
  "common": {
    "confirm": "确定",
    "cancel": "取消"
  },
  "user": {
    "profile": {
      "title": "个人资料"
    }
  }
}
```

使用时：
```go
gi18n.T("common.confirm")       // 确定
gi18n.T("user.profile.title")   // 个人资料
```

### 复数格式

```json
{
  "items": {
    "one": "{{.Count}} item",
    "other": "{{.Count}} items"
  }
}
```

## 配置

### 基础配置

```go
gi18n.Init(&gi18n.Config{
    DefaultLang:  "zh-CN",
    FallbackLang: "en",
})
```

### 翻译缺失处理

```go
gi18n.Init(&gi18n.Config{
    // 缺失时回调（用于日志/监控）
    MissHandler: func(lang, id string) {
        log.Printf("missing translation: lang=%s, id=%s", lang, id)
    },

    // 缺失策略: MissReturnID（默认）或 MissReturnEmpty
    MissPolicy: gi18n.MissReturnID,
})
```

### 日志集成

```go
// 兼容 slog
gi18n.Init(&gi18n.Config{
    Logger: slog.Default(),
})

// 兼容任何实现了 Warn(msg string, args ...any) 的接口
```

### 多实例

```go
bundle := gi18n.New(&gi18n.Config{
    DefaultLang:  "zh-CN",
    FallbackLang: "en",
})
bundle.Load("./locales")
bundle.T("hello")
```

## HTTP 中间件

### 标准库

```go
mux := http.NewServeMux()
mux.Handle("/", gi18n.Middleware(nil)(yourHandler))

// 在 handler 中使用
func handler(w http.ResponseWriter, r *http.Request) {
    msg := gi18n.T("hello", gi18n.WithContext(r.Context()))
    w.Write([]byte(msg))
}
```

### Gin 集成

```go
r := gin.Default()
r.Use(func(c *gin.Context) {
    handler := gi18n.Middleware(nil)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        c.Request = r
        c.Next()
    }))
    handler.ServeHTTP(c.Writer, c.Request)
})

func handler(c *gin.Context) {
    msg := gi18n.T("hello", gi18n.WithContext(c.Request.Context()))
    c.String(200, msg)
}
```

### 中间件配置

```go
cfg := &gi18n.MiddlewareConfig{
    Sources:     []gi18n.LangSource{gi18n.SourceQuery, gi18n.SourceCookie, gi18n.SourceHeader},
    QueryParam:  "lang",      // URL 参数名
    CookieName:  "lang",      // Cookie 名
    DefaultLang: "en",        // 默认语言
}
gi18n.Middleware(cfg)
```

语言检测优先级（可配置）：
1. URL 参数 `?lang=zh-CN`
2. Cookie `lang=zh-CN`
3. `Accept-Language` 头

## 内置语言包

内置 6 种语言 33 个通用词条，开箱即用：

| 语言 | 标识 |
|------|------|
| 英语 | `en` |
| 简体中文 | `zh-CN` |
| 繁体中文 | `zh-TW` |
| 日语 | `ja` |
| 韩语 | `ko` |
| 俄语 | `ru` |

词条涵盖：`confirm`、`cancel`、`save`、`delete`、`edit`、`submit`、`reset`、`search`、`close`、`back`、`next`、`prev`、`yes`、`no`、`success`、`failed`、`error`、`warning`、`info`、`loading`、`required`、`optional`、`invalid`、`username`、`password`、`email`、`phone`、`login`、`logout`、`register`、`welcome`、`greeting`（带参数）、`items`（复数）

### 禁用内置语言包

```bash
go build -tags=gi18n_no_builtin
```

## 从旧版迁移

如果你使用的是旧版 API（`TL`, `Tf`, `TLf`, `Tp` 等），这些方法仍然可用但已标记为 `Deprecated`。建议迁移到新 API：

| 旧 API | 新 API |
|--------|--------|
| `T(id)` | `T(id)` ✅ 无需改动 |
| `TL(lang, id)` | `T(id, WithLang(lang))` |
| `Tf(id, "Name", "张三")` | `T(id, WithData("Name", "张三"))` |
| `TLf(lang, id, args...)` | `T(id, WithLang(lang), WithData(args...))` |
| `Tp(id, count)` | `T(id, WithCount(count))` |
| `TLp(lang, id, count)` | `T(id, WithLang(lang), WithCount(count))` |
| `TC(ctx, id)` | `T(id, WithContext(ctx))` |
| `TCf(ctx, id, args...)` | `T(id, WithContext(ctx), WithData(args...))` |
| `TMap(id, data)` | `T(id, WithMap(data))` |
| `WithLang(ctx, lang)` | `ContextWithLang(ctx, lang)` |
| `Langs()` | `Languages()` |
| `SetLanguage(lang)` | `SetLang(lang)` |
| `GetLanguage()` | `GetLang()` |
| `GetLanguages()` | `Languages()` |

## License

MIT
