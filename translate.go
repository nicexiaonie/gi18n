package gi18n

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// ========== 语言设置 ==========

// SetLang 设置当前语言
func (b *Bundle) SetLang(lang string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.currentLang = normalizeLanguageTag(lang)
}

// GetLang 获取当前语言
func (b *Bundle) GetLang() string {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.currentLang
}

// Languages 获取支持的语言列表
func (b *Bundle) Languages() []string {
	b.mu.RLock()
	defer b.mu.RUnlock()
	result := make([]string, len(b.supported))
	copy(result, b.supported)
	return result
}

// SetDefaultLang 设置默认语言
func (b *Bundle) SetDefaultLang(lang string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.defaultLang = normalizeLanguageTag(lang)
}

// SetFallbackLang 设置回退语言
func (b *Bundle) SetFallbackLang(lang string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.fallbackLang = normalizeLanguageTag(lang)
	b.clearLocalizerCache()
}

// ========== 核心翻译方法 ==========

// T 统一翻译入口
//
// 简单翻译:
//
//	bundle.T("confirm")
//
// 指定语言:
//
//	bundle.T("confirm", WithLang("zh-CN"))
//
// 带参数:
//
//	bundle.T("greeting", WithData("Name", "张三"))
//
// 复数:
//
//	bundle.T("items", WithCount(5))
//
// 组合使用:
//
//	bundle.T("items", WithLang("en"), WithCount(5))
//
// 从 Context 获取语言:
//
//	bundle.T("hello", WithContext(ctx))
func (b *Bundle) T(id string, opts ...Option) string {
	tc := &translateConfig{}
	for _, opt := range opts {
		opt(tc)
	}

	// 语言优先级: 显式指定 > Context > 当前语言
	lang := b.GetLang()
	if tc.ctx != nil {
		if ctxLang, ok := tc.ctx.Value(langCtxKey).(string); ok && ctxLang != "" {
			lang = ctxLang
		}
	}
	if tc.lang != "" {
		lang = tc.lang
	}

	// 构建 LocalizeConfig
	lc := &i18n.LocalizeConfig{
		MessageID: id,
	}

	if tc.data != nil {
		lc.TemplateData = tc.data
	}

	if tc.count != nil {
		lc.PluralCount = *tc.count
		if lc.TemplateData == nil {
			lc.TemplateData = map[string]interface{}{"Count": *tc.count}
		} else if m, ok := lc.TemplateData.(map[string]interface{}); ok {
			m["Count"] = *tc.count
		}
	}

	loc := b.getLocalizer(lang)
	msg, err := loc.Localize(lc)
	if err != nil {
		b.handleMiss(lang, id)
		if b.missPolicy == MissReturnEmpty {
			return ""
		}
		return id
	}
	return msg
}

// ========== 全局核心函数 ==========

// SetLang 设置当前语言（全局）
func SetLang(lang string) { Default().SetLang(lang) }

// GetLang 获取当前语言（全局）
func GetLang() string { return Default().GetLang() }

// Languages 获取支持的语言列表（全局）
func Languages() []string { return Default().Languages() }

// SetDefaultLang 设置默认语言（全局）
func SetDefaultLang(lang string) { Default().SetDefaultLang(lang) }

// SetFallbackLang 设置回退语言（全局）
func SetFallbackLang(lang string) { Default().SetFallbackLang(lang) }

// T 统一翻译入口（全局）
//
//	gi18n.T("confirm")
//	gi18n.T("confirm", gi18n.WithLang("zh-CN"))
//	gi18n.T("greeting", gi18n.WithData("Name", "张三"))
//	gi18n.T("items", gi18n.WithCount(5))
//	gi18n.T("hello", gi18n.WithContext(ctx))
func T(id string, opts ...Option) string {
	return Default().T(id, opts...)
}

// ========== 已废弃的实例方法（向后兼容） ==========

// Deprecated: Use SetLang instead.
func (b *Bundle) SetLanguage(lang string) { b.SetLang(lang) }

// Deprecated: Use GetLang instead.
func (b *Bundle) GetLanguage() string { return b.GetLang() }

// Deprecated: Use Languages instead.
func (b *Bundle) Langs() []string { return b.Languages() }

// Deprecated: Use Languages instead.
func (b *Bundle) GetLanguages() []string { return b.Languages() }

// Deprecated: Use T(id) instead.
func (b *Bundle) Translate(id string) string { return b.T(id) }

// Deprecated: Use T(id, WithLang(lang)) instead.
func (b *Bundle) TL(lang, id string) string { return b.T(id, WithLang(lang)) }

// Deprecated: Use T(id, WithLang(lang)) instead.
func (b *Bundle) TranslateLang(lang, id string) string { return b.T(id, WithLang(lang)) }

// Deprecated: Use T(id, WithData(args...)) instead.
func (b *Bundle) Tf(id string, args ...interface{}) string {
	return b.T(id, WithData(args...))
}

// Deprecated: Use T(id, WithData(args...)) instead.
func (b *Bundle) TranslateWith(id string, args ...interface{}) string {
	return b.T(id, WithData(args...))
}

// Deprecated: Use T(id, WithLang(lang), WithData(args...)) instead.
func (b *Bundle) TLf(lang, id string, args ...interface{}) string {
	return b.T(id, WithLang(lang), WithData(args...))
}

// Deprecated: Use T(id, WithLang(lang), WithData(args...)) instead.
func (b *Bundle) TranslateLangWith(lang, id string, args ...interface{}) string {
	return b.T(id, WithLang(lang), WithData(args...))
}

// Deprecated: Use T(id, WithCount(count)) instead.
func (b *Bundle) Tp(id string, count int, args ...interface{}) string {
	opts := []Option{WithCount(count)}
	if len(args) > 0 {
		opts = append(opts, WithData(args...))
	}
	return b.T(id, opts...)
}

// Deprecated: Use T(id, WithCount(count)) instead.
func (b *Bundle) TranslatePlural(id string, count int, args ...interface{}) string {
	return b.Tp(id, count, args...)
}

// Deprecated: Use T(id, WithLang(lang), WithCount(count)) instead.
func (b *Bundle) TLp(lang, id string, count int, args ...interface{}) string {
	opts := []Option{WithLang(lang), WithCount(count)}
	if len(args) > 0 {
		opts = append(opts, WithData(args...))
	}
	return b.T(id, opts...)
}

// Deprecated: Use T(id, WithLang(lang), WithCount(count)) instead.
func (b *Bundle) TranslateLangPlural(lang, id string, count int, args ...interface{}) string {
	return b.TLp(lang, id, count, args...)
}

// Deprecated: Use T(id, WithMap(data)) instead.
func (b *Bundle) TMap(id string, data map[string]interface{}) string {
	return b.T(id, WithMap(data))
}

// Deprecated: Use T(id, WithLang(lang), WithMap(data)) instead.
func (b *Bundle) TLMap(lang, id string, data map[string]interface{}) string {
	return b.T(id, WithLang(lang), WithMap(data))
}

// ========== 已废弃的全局函数（向后兼容） ==========

// Deprecated: Use SetLang instead.
func SetLanguage(lang string) { Default().SetLang(lang) }

// Deprecated: Use GetLang instead.
func GetLanguage() string { return Default().GetLang() }

// Deprecated: Use Languages instead.
func Langs() []string { return Default().Languages() }

// Deprecated: Use Languages instead.
func GetLanguages() []string { return Default().Languages() }

// Deprecated: Use T(id) instead.
func Translate(id string) string { return Default().T(id) }

// Deprecated: Use T(id, WithLang(lang)) instead.
func TL(lang, id string) string { return Default().TL(lang, id) }

// Deprecated: Use T(id, WithLang(lang)) instead.
func TranslateLang(lang, id string) string { return Default().TranslateLang(lang, id) }

// Deprecated: Use T(id, WithData(args...)) instead.
func Tf(id string, args ...interface{}) string { return Default().Tf(id, args...) }

// Deprecated: Use T(id, WithData(args...)) instead.
func TranslateWith(id string, args ...interface{}) string {
	return Default().TranslateWith(id, args...)
}

// Deprecated: Use T(id, WithLang(lang), WithData(args...)) instead.
func TLf(lang, id string, args ...interface{}) string {
	return Default().TLf(lang, id, args...)
}

// Deprecated: Use T(id, WithLang(lang), WithData(args...)) instead.
func TranslateLangWith(lang, id string, args ...interface{}) string {
	return Default().TranslateLangWith(lang, id, args...)
}

// Deprecated: Use T(id, WithCount(count)) instead.
func Tp(id string, count int, args ...interface{}) string {
	return Default().Tp(id, count, args...)
}

// Deprecated: Use T(id, WithCount(count)) instead.
func TranslatePlural(id string, count int, args ...interface{}) string {
	return Default().TranslatePlural(id, count, args...)
}

// Deprecated: Use T(id, WithLang(lang), WithCount(count)) instead.
func TLp(lang, id string, count int, args ...interface{}) string {
	return Default().TLp(lang, id, count, args...)
}

// Deprecated: Use T(id, WithLang(lang), WithCount(count)) instead.
func TranslateLangPlural(lang, id string, count int, args ...interface{}) string {
	return Default().TranslateLangPlural(lang, id, count, args...)
}

// Deprecated: Use T(id, WithMap(data)) instead.
func TMap(id string, data map[string]interface{}) string {
	return Default().TMap(id, data)
}

// Deprecated: Use T(id, WithLang(lang), WithMap(data)) instead.
func TLMap(lang, id string, data map[string]interface{}) string {
	return Default().TLMap(lang, id, data)
}
