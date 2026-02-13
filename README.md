# gi18n

基于 [go-i18n](https://github.com/nicksnyder/go-i18n) 封装的简单易用的 Go 国际化库。

## 特性

- **零配置启动**: 内置 6 种语言的通用词条，开箱即用
- **多格式支持**: JSON、YAML、TOML 全部支持
- **灵活加载**: 支持路径加载、embed.FS 加载、内容加载
- **嵌套展平**: 支持嵌套结构，自动展平为扁平 key
- **语言标签兼容**: 兼容 `zh-CN`、`zh_CN`、`zh-Hans` 等格式
- **回退机制**: 找不到翻译时自动回退到默认语言
- **HTTP 中间件**: 内置 HTTP 中间件，支持 Gin 等框架
- **两套 API**: 提供简短版和完整版方法名

## 安装

```bash
go get github.com/nicexiaonie/gi18n
```

## 快速开始

```go
package main

import (
    "fmt"
    "github.com/nicexiaonie/gi18n"
)

func main() {
    // 内置语言包已自动加载，直接使用
    gi18n.SetLang("zh-CN")
    fmt.Println(gi18n.T("confirm"))  // 输出: 确定
    fmt.Println(gi18n.T("cancel"))   // 输出: 取消

    // 切换语言
    gi18n.SetLang("en")
    fmt.Println(gi18n.T("confirm"))  // 输出: OK

    // 带参数
    gi18n.SetLang("zh-CN")
    fmt.Println(gi18n.Tf("greeting", "Name", "张三"))  // 输出: 你好，张三！

    // 复数
    fmt.Println(gi18n.Tp("items", 5))  // 输出: 5 个项目
}
```

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
zhData := []byte(`{"hello": "你好", "world": "世界"}`)
gi18n.LoadContent("zh-CN", "json", zhData)
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

### go-i18n 完整格式（支持复数）

```json
{
  "items": {
    "id": "items",
    "description": "项目数量",
    "one": "{{.Count}} 个项目",
    "other": "{{.Count}} 个项目"
  }
}
```

## API 参考

### 翻译方法

| 简短版 | 完整版 | 说明 |
|-------|--------|-----|
| `T(id)` | `Translate(id)` | 简单翻译 |
| `TL(lang, id)` | `TranslateLang(lang, id)` | 指定语言翻译 |
| `Tf(id, args...)` | `TranslateWith(id, args...)` | 带参数翻译 |
| `TLf(lang, id, args...)` | `TranslateLangWith(...)` | 指定语言带参数 |
| `Tp(id, count, args...)` | `TranslatePlural(...)` | 复数翻译 |
| `TLp(lang, id, count, args...)` | `TranslateLangPlural(...)` | 指定语言复数 |
| `TC(ctx, id)` | `TranslateContext(ctx, id)` | 从 context 获取语言 |

### 语言设置

| 简短版 | 完整版 | 说明 |
|-------|--------|-----|
| `SetLang(lang)` | `SetLanguage(lang)` | 设置当前语言 |
| `GetLang()` | `GetLanguage()` | 获取当前语言 |
| `Langs()` | `GetLanguages()` | 获取支持的语言列表 |

### 加载方法

| 方法 | 说明 |
|-----|------|
| `Load(dir)` | 从目录加载 |
| `LoadFS(fs, root)` | 从 embed.FS 加载 |
| `LoadContent(lang, format, data)` | 从内容加载 |
| `LoadMessages(lang, messages)` | 加载消息映射 |

## HTTP 中间件

### 标准库

```go
mux := http.NewServeMux()
mux.Handle("/", gi18n.Middleware(nil)(yourHandler))
```

### Gin

```go
r := gin.Default()
r.Use(func(c *gin.Context) {
    middleware := gi18n.Middleware(nil)
    middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        c.Request = r
        c.Next()
    })).ServeHTTP(c.Writer, c.Request)
})

// 在 handler 中使用
func handler(c *gin.Context) {
    msg := gi18n.TC(c.Request.Context(), "hello")
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
middleware := gi18n.Middleware(cfg)
```

语言检测优先级（可配置）：
1. URL 参数 `?lang=zh-CN`
2. Cookie `lang=zh-CN`
3. Accept-Language 头

## 内置语言包

内置 6 种语言的通用词条：

- `en` - 英语
- `zh-CN` - 简体中文
- `zh-TW` - 繁体中文
- `ja` - 日语
- `ko` - 韩语
- `ru` - 俄语

内置词条包括：`confirm`、`cancel`、`save`、`delete`、`edit`、`submit`、`reset`、`search`、`close`、`back`、`next`、`prev`、`yes`、`no`、`success`、`failed`、`error`、`warning`、`info`、`loading`、`required`、`optional`、`invalid`、`username`、`password`、`email`、`phone`、`login`、`logout`、`register`、`welcome`、`greeting`、`items`

### 禁用内置语言包

编译时添加 tag：

```bash
go build -tags=gi18n_no_builtin
```

## 高级用法

### 多实例

```go
// 创建独立实例
bundle := gi18n.New(&gi18n.Config{
    DefaultLang:  "zh-CN",
    FallbackLang: "en",
})
bundle.Load("./locales")
bundle.T("hello")
```

### 获取底层 go-i18n Bundle

```go
bundle := gi18n.Default().GetBundle()
// 使用 go-i18n 原生 API
```

## License

MIT
