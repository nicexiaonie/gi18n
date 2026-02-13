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

// SetLanguage SetLang 的别名
func (b *Bundle) SetLanguage(lang string) {
	b.SetLang(lang)
}

// GetLang 获取当前语言
func (b *Bundle) GetLang() string {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.currentLang
}

// GetLanguage GetLang 的别名
func (b *Bundle) GetLanguage() string {
	return b.GetLang()
}

// Langs 获取支持的语言列表
func (b *Bundle) Langs() []string {
	b.mu.RLock()
	defer b.mu.RUnlock()
	result := make([]string, len(b.supported))
	copy(result, b.supported)
	return result
}

// GetLanguages Langs 的别名
func (b *Bundle) GetLanguages() []string {
	return b.Langs()
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

// ========== 翻译方法（实例） ==========

// T 简单翻译
func (b *Bundle) T(id string) string {
	return b.TL(b.GetLang(), id)
}

// Translate T 的别名
func (b *Bundle) Translate(id string) string {
	return b.T(id)
}

// TL 指定语言的简单翻译
func (b *Bundle) TL(lang, id string) string {
	loc := b.getLocalizer(lang)
	msg, err := loc.Localize(&i18n.LocalizeConfig{
		MessageID: id,
	})
	if err != nil {
		return id // 返回原始 ID 作为后备
	}
	return msg
}

// TranslateLang TL 的别名
func (b *Bundle) TranslateLang(lang, id string) string {
	return b.TL(lang, id)
}

// Tf 带参数的翻译
// 参数格式: key1, value1, key2, value2, ...
func (b *Bundle) Tf(id string, args ...interface{}) string {
	return b.TLf(b.GetLang(), id, args...)
}

// TranslateWith Tf 的别名
func (b *Bundle) TranslateWith(id string, args ...interface{}) string {
	return b.Tf(id, args...)
}

// TLf 指定语言的带参数翻译
func (b *Bundle) TLf(lang, id string, args ...interface{}) string {
	data := argsToMap(args...)
	loc := b.getLocalizer(lang)
	msg, err := loc.Localize(&i18n.LocalizeConfig{
		MessageID:    id,
		TemplateData: data,
	})
	if err != nil {
		return id
	}
	return msg
}

// TranslateLangWith TLf 的别名
func (b *Bundle) TranslateLangWith(lang, id string, args ...interface{}) string {
	return b.TLf(lang, id, args...)
}

// Tp 带复数的翻译
func (b *Bundle) Tp(id string, count int, args ...interface{}) string {
	return b.TLp(b.GetLang(), id, count, args...)
}

// TranslatePlural Tp 的别名
func (b *Bundle) TranslatePlural(id string, count int, args ...interface{}) string {
	return b.Tp(id, count, args...)
}

// TLp 指定语言的带复数翻译
func (b *Bundle) TLp(lang, id string, count int, args ...interface{}) string {
	data := argsToMap(args...)
	data["Count"] = count

	loc := b.getLocalizer(lang)
	msg, err := loc.Localize(&i18n.LocalizeConfig{
		MessageID:    id,
		TemplateData: data,
		PluralCount:  count,
	})
	if err != nil {
		return id
	}
	return msg
}

// TranslateLangPlural TLp 的别名
func (b *Bundle) TranslateLangPlural(lang, id string, count int, args ...interface{}) string {
	return b.TLp(lang, id, count, args...)
}

// TMap 使用 map 参数的翻译
func (b *Bundle) TMap(id string, data map[string]interface{}) string {
	return b.TLMap(b.GetLang(), id, data)
}

// TLMap 指定语言使用 map 参数的翻译
func (b *Bundle) TLMap(lang, id string, data map[string]interface{}) string {
	loc := b.getLocalizer(lang)
	msg, err := loc.Localize(&i18n.LocalizeConfig{
		MessageID:    id,
		TemplateData: data,
	})
	if err != nil {
		return id
	}
	return msg
}

// ========== 全局函数 ==========

// SetLang 设置当前语言（全局）
func SetLang(lang string) {
	Default().SetLang(lang)
}

// SetLanguage SetLang 的别名（全局）
func SetLanguage(lang string) {
	Default().SetLanguage(lang)
}

// GetLang 获取当前语言（全局）
func GetLang() string {
	return Default().GetLang()
}

// GetLanguage GetLang 的别名（全局）
func GetLanguage() string {
	return Default().GetLanguage()
}

// Langs 获取支持的语言列表（全局）
func Langs() []string {
	return Default().Langs()
}

// GetLanguages Langs 的别名（全局）
func GetLanguages() []string {
	return Default().GetLanguages()
}

// SetDefaultLang 设置默认语言（全局）
func SetDefaultLang(lang string) {
	Default().SetDefaultLang(lang)
}

// SetFallbackLang 设置回退语言（全局）
func SetFallbackLang(lang string) {
	Default().SetFallbackLang(lang)
}

// T 简单翻译（全局）
func T(id string) string {
	return Default().T(id)
}

// Translate T 的别名（全局）
func Translate(id string) string {
	return Default().Translate(id)
}

// TL 指定语言的简单翻译（全局）
func TL(lang, id string) string {
	return Default().TL(lang, id)
}

// TranslateLang TL 的别名（全局）
func TranslateLang(lang, id string) string {
	return Default().TranslateLang(lang, id)
}

// Tf 带参数的翻译（全局）
func Tf(id string, args ...interface{}) string {
	return Default().Tf(id, args...)
}

// TranslateWith Tf 的别名（全局）
func TranslateWith(id string, args ...interface{}) string {
	return Default().TranslateWith(id, args...)
}

// TLf 指定语言的带参数翻译（全局）
func TLf(lang, id string, args ...interface{}) string {
	return Default().TLf(lang, id, args...)
}

// TranslateLangWith TLf 的别名（全局）
func TranslateLangWith(lang, id string, args ...interface{}) string {
	return Default().TranslateLangWith(lang, id, args...)
}

// Tp 带复数的翻译（全局）
func Tp(id string, count int, args ...interface{}) string {
	return Default().Tp(id, count, args...)
}

// TranslatePlural Tp 的别名（全局）
func TranslatePlural(id string, count int, args ...interface{}) string {
	return Default().TranslatePlural(id, count, args...)
}

// TLp 指定语言的带复数翻译（全局）
func TLp(lang, id string, count int, args ...interface{}) string {
	return Default().TLp(lang, id, count, args...)
}

// TranslateLangPlural TLp 的别名（全局）
func TranslateLangPlural(lang, id string, count int, args ...interface{}) string {
	return Default().TranslateLangPlural(lang, id, count, args...)
}

// TMap 使用 map 参数的翻译（全局）
func TMap(id string, data map[string]interface{}) string {
	return Default().TMap(id, data)
}

// TLMap 指定语言使用 map 参数的翻译（全局）
func TLMap(lang, id string, data map[string]interface{}) string {
	return Default().TLMap(lang, id, data)
}

// ========== 辅助函数 ==========

// argsToMap 将可变参数转换为 map
// 格式: key1, value1, key2, value2, ...
func argsToMap(args ...interface{}) map[string]interface{} {
	data := make(map[string]interface{})
	for i := 0; i+1 < len(args); i += 2 {
		if key, ok := args[i].(string); ok {
			data[key] = args[i+1]
		}
	}
	return data
}
