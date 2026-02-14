package gi18n

import (
	"context"
	"net/http"
	"strings"
)

type ctxKey struct{}

var langCtxKey = ctxKey{}

// ========== Context 集成 ==========

// ContextWithLang 将语言设置注入到 context 中
//
//	ctx := gi18n.ContextWithLang(ctx, "zh-CN")
//	gi18n.T("hello", gi18n.WithContext(ctx))
func ContextWithLang(ctx context.Context, lang string) context.Context {
	return context.WithValue(ctx, langCtxKey, normalizeLanguageTag(lang))
}

// LangFromContext 从 context 获取语言
// 如果 context 中没有语言信息，返回全局默认语言
func LangFromContext(ctx context.Context) string {
	if lang, ok := ctx.Value(langCtxKey).(string); ok {
		return lang
	}
	return Default().GetLang()
}

// ========== HTTP 中间件 ==========

// LangSource 语言来源类型
type LangSource int

const (
	SourceQuery  LangSource = iota // URL 参数 ?lang=zh-CN
	SourceHeader                   // Accept-Language 头
	SourceCookie                   // Cookie
)

// MiddlewareConfig 中间件配置
type MiddlewareConfig struct {
	// 语言来源优先级，默认: Query > Cookie > Header
	Sources []LangSource
	// URL 参数名，默认 "lang"
	QueryParam string
	// Cookie 名，默认 "lang"
	CookieName string
	// 默认语言，默认使用全局设置
	DefaultLang string
}

// DefaultMiddlewareConfig 默认中间件配置
func DefaultMiddlewareConfig() *MiddlewareConfig {
	return &MiddlewareConfig{
		Sources:     []LangSource{SourceQuery, SourceCookie, SourceHeader},
		QueryParam:  "lang",
		CookieName:  "lang",
		DefaultLang: "",
	}
}

// Middleware 创建 HTTP 中间件
//
// 标准库:
//
//	mux.Handle("/", gi18n.Middleware(nil)(handler))
//
// Gin:
//
//	r.Use(func(c *gin.Context) {
//	    handler := gi18n.Middleware(nil)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//	        c.Request = r
//	        c.Next()
//	    }))
//	    handler.ServeHTTP(c.Writer, c.Request)
//	})
func Middleware(cfg *MiddlewareConfig) func(http.Handler) http.Handler {
	if cfg == nil {
		cfg = DefaultMiddlewareConfig()
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			lang := detectLanguage(r, cfg)
			ctx := ContextWithLang(r.Context(), lang)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// detectLanguage 检测请求的语言
func detectLanguage(r *http.Request, cfg *MiddlewareConfig) string {
	for _, source := range cfg.Sources {
		var lang string
		switch source {
		case SourceQuery:
			lang = r.URL.Query().Get(cfg.QueryParam)
		case SourceCookie:
			if cookie, err := r.Cookie(cfg.CookieName); err == nil {
				lang = cookie.Value
			}
		case SourceHeader:
			lang = parseAcceptLanguage(r.Header.Get("Accept-Language"))
		}

		if lang != "" {
			return normalizeLanguageTag(lang)
		}
	}

	if cfg.DefaultLang != "" {
		return normalizeLanguageTag(cfg.DefaultLang)
	}
	return Default().GetLang()
}

// parseAcceptLanguage 解析 Accept-Language 头
func parseAcceptLanguage(header string) string {
	if header == "" {
		return ""
	}

	// 取第一个语言，去掉权重
	// Accept-Language: zh-CN,zh;q=0.9,en;q=0.8
	parts := strings.Split(header, ",")
	if len(parts) == 0 {
		return ""
	}

	lang := strings.TrimSpace(parts[0])
	if idx := strings.Index(lang, ";"); idx > 0 {
		lang = lang[:idx]
	}

	return lang
}

// ========== 已废弃方法（向后兼容） ==========

// Deprecated: Use ContextWithLang instead.
// 注意: 此函数在新版本中签名已变更为 Option 类型的 WithLang(lang string)。
// 原 WithLang(ctx, lang) 请改用 ContextWithLang(ctx, lang)。

// Deprecated: Use T(id, WithContext(ctx)) instead.
func TC(ctx context.Context, id string) string {
	return T(id, WithContext(ctx))
}

// Deprecated: Use T(id, WithContext(ctx)) instead.
func TranslateContext(ctx context.Context, id string) string {
	return TC(ctx, id)
}

// Deprecated: Use T(id, WithContext(ctx), WithData(args...)) instead.
func TCf(ctx context.Context, id string, args ...interface{}) string {
	return T(id, WithContext(ctx), WithData(args...))
}

// Deprecated: Use T(id, WithContext(ctx), WithData(args...)) instead.
func TranslateContextWith(ctx context.Context, id string, args ...interface{}) string {
	return TCf(ctx, id, args...)
}

// Deprecated: Use T(id, WithContext(ctx), WithCount(count)) instead.
func TCp(ctx context.Context, id string, count int, args ...interface{}) string {
	opts := []Option{WithContext(ctx), WithCount(count)}
	if len(args) > 0 {
		opts = append(opts, WithData(args...))
	}
	return T(id, opts...)
}

// Deprecated: Use T(id, WithContext(ctx), WithCount(count)) instead.
func TranslateContextPlural(ctx context.Context, id string, count int, args ...interface{}) string {
	return TCp(ctx, id, count, args...)
}

// Deprecated: Use Middleware instead. Gin 可直接使用标准 Middleware 适配。
func GinMiddleware(cfg *MiddlewareConfig) func(c interface {
	Next()
	Request() *http.Request
	Set(key string, value interface{})
}) {
	if cfg == nil {
		cfg = DefaultMiddlewareConfig()
	}

	return func(c interface {
		Next()
		Request() *http.Request
		Set(key string, value interface{})
	}) {
		lang := detectLanguage(c.Request(), cfg)
		c.Set("gi18n_lang", lang)
		c.Next()
	}
}

// Deprecated: Use T(id, WithContext(ctx)) with standard Middleware instead.
func LangFromGin(c interface{ GetString(key string) string }) string {
	if lang := c.GetString("gi18n_lang"); lang != "" {
		return lang
	}
	return Default().GetLang()
}
