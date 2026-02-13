package gi18n

import (
	"context"
	"net/http"
	"strings"
)

type ctxKey struct{}

var langCtxKey = ctxKey{}

// WithLang 将语言设置注入到 context 中
func WithLang(ctx context.Context, lang string) context.Context {
	return context.WithValue(ctx, langCtxKey, normalizeLanguageTag(lang))
}

// LangFromContext 从 context 获取语言
func LangFromContext(ctx context.Context) string {
	if lang, ok := ctx.Value(langCtxKey).(string); ok {
		return lang
	}
	return Default().GetLang()
}

// TC 从 context 获取语言并翻译
func TC(ctx context.Context, id string) string {
	return TL(LangFromContext(ctx), id)
}

// TranslateContext TC 的别名
func TranslateContext(ctx context.Context, id string) string {
	return TC(ctx, id)
}

// TCf 从 context 获取语言并带参数翻译
func TCf(ctx context.Context, id string, args ...interface{}) string {
	return TLf(LangFromContext(ctx), id, args...)
}

// TranslateContextWith TCf 的别名
func TranslateContextWith(ctx context.Context, id string, args ...interface{}) string {
	return TCf(ctx, id, args...)
}

// TCp 从 context 获取语言并带复数翻译
func TCp(ctx context.Context, id string, count int, args ...interface{}) string {
	return TLp(LangFromContext(ctx), id, count, args...)
}

// TranslateContextPlural TCp 的别名
func TranslateContextPlural(ctx context.Context, id string, count int, args ...interface{}) string {
	return TCp(ctx, id, count, args...)
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
func Middleware(cfg *MiddlewareConfig) func(http.Handler) http.Handler {
	if cfg == nil {
		cfg = DefaultMiddlewareConfig()
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			lang := detectLanguage(r, cfg)
			ctx := WithLang(r.Context(), lang)
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

	// 使用默认语言
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

	// 简单解析，取第一个语言
	// Accept-Language: zh-CN,zh;q=0.9,en;q=0.8
	parts := strings.Split(header, ",")
	if len(parts) == 0 {
		return ""
	}

	// 取第一个，去掉权重
	lang := strings.TrimSpace(parts[0])
	if idx := strings.Index(lang, ";"); idx > 0 {
		lang = lang[:idx]
	}

	return lang
}

// ========== Gin 适配 ==========

// GinMiddleware 返回 Gin 风格的中间件函数
// 使用方式: router.Use(gi18n.GinMiddleware(nil))
func GinMiddleware(cfg *MiddlewareConfig) func(c interface{ Next(); Request() *http.Request; Set(key string, value interface{}) }) {
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

// LangFromGin 从 Gin context 获取语言
func LangFromGin(c interface{ GetString(key string) string }) string {
	if lang := c.GetString("gi18n_lang"); lang != "" {
		return lang
	}
	return Default().GetLang()
}
