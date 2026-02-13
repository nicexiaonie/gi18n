// Package gi18n 提供简单易用的国际化封装
// 基于 go-i18n 库，提供零配置、开箱即用的国际化能力
package gi18n

import (
	"sync"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

// 全局实例
var (
	globalBundle *Bundle
	once         sync.Once
)

// Bundle 国际化包装器
type Bundle struct {
	bundle       *i18n.Bundle
	mu           sync.RWMutex
	localizers   sync.Map // map[string]*i18n.Localizer
	currentLang  string
	defaultLang  string
	fallbackLang string
	supported    []string
}

// Config 初始化配置
type Config struct {
	DefaultLang  string // 默认语言，默认 "en"
	FallbackLang string // 回退语言，默认 "en"
}

// Default 获取全局默认实例
func Default() *Bundle {
	once.Do(func() {
		globalBundle = New(nil)
	})
	return globalBundle
}

// New 创建新的 Bundle 实例
func New(cfg *Config) *Bundle {
	defaultLang := "en"
	fallbackLang := "en"

	if cfg != nil {
		if cfg.DefaultLang != "" {
			defaultLang = cfg.DefaultLang
		}
		if cfg.FallbackLang != "" {
			fallbackLang = cfg.FallbackLang
		}
	}

	// 解析语言标签
	tag := parseLanguageTag(defaultLang)

	b := &Bundle{
		bundle:       i18n.NewBundle(tag),
		currentLang:  defaultLang,
		defaultLang:  defaultLang,
		fallbackLang: fallbackLang,
		supported:    make([]string, 0),
	}

	// 注册解析器
	b.registerUnmarshalers()

	return b
}

// Init 初始化全局实例
func Init(cfg *Config) {
	globalBundle = New(cfg)
}

// parseLanguageTag 解析语言标签，兼容多种格式
// 支持: zh-CN, zh_CN, zh-Hans 等
func parseLanguageTag(lang string) language.Tag {
	// 统一转换下划线为连字符
	normalized := normalizeLanguageTag(lang)
	tag, err := language.Parse(normalized)
	if err != nil {
		return language.English
	}
	return tag
}

// normalizeLanguageTag 标准化语言标签
func normalizeLanguageTag(lang string) string {
	// 将下划线替换为连字符: zh_CN -> zh-CN
	result := make([]byte, len(lang))
	for i := 0; i < len(lang); i++ {
		if lang[i] == '_' {
			result[i] = '-'
		} else {
			result[i] = lang[i]
		}
	}
	return string(result)
}

// addSupported 添加支持的语言
func (b *Bundle) addSupported(lang string) {
	normalized := normalizeLanguageTag(lang)
	for _, l := range b.supported {
		if l == normalized {
			return
		}
	}
	b.supported = append(b.supported, normalized)
}

// getLocalizer 获取指定语言的 Localizer（带缓存）
func (b *Bundle) getLocalizer(lang string) *i18n.Localizer {
	normalized := normalizeLanguageTag(lang)

	// 尝试从缓存获取
	if loc, ok := b.localizers.Load(normalized); ok {
		return loc.(*i18n.Localizer)
	}

	// 创建新的 Localizer，包含回退语言
	loc := i18n.NewLocalizer(b.bundle, normalized, b.fallbackLang)
	b.localizers.Store(normalized, loc)
	return loc
}

// clearLocalizerCache 清空 Localizer 缓存
func (b *Bundle) clearLocalizerCache() {
	b.localizers = sync.Map{}
}

// GetBundle 获取底层的 go-i18n Bundle（高级用法）
func (b *Bundle) GetBundle() *i18n.Bundle {
	return b.bundle
}
