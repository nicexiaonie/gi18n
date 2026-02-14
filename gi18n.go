// Package gi18n 提供简单易用的国际化封装
// 基于 go-i18n 库，提供零配置、开箱即用的国际化能力
//
// 核心 API 只有一个翻译入口 T()，通过 Option 组合实现所有场景：
//
//	gi18n.T("confirm")                                    // 简单翻译
//	gi18n.T("confirm", gi18n.WithLang("zh-CN"))           // 指定语言
//	gi18n.T("greeting", gi18n.WithData("Name", "张三"))    // 带参数
//	gi18n.T("items", gi18n.WithCount(5))                  // 复数
//	gi18n.T("hello", gi18n.WithContext(ctx))               // 从 Context
package gi18n

import (
	"strings"
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
	missHandler  func(lang, id string)
	missPolicy   MissPolicy
	logger       Logger
}

// Config 初始化配置
type Config struct {
	DefaultLang  string // 默认语言，默认 "en"
	FallbackLang string // 回退语言，默认 "en"

	// MissHandler 翻译缺失回调（可选）
	// 当翻译 key 不存在时触发，可用于日志记录或监控
	MissHandler func(lang, id string)

	// MissPolicy 翻译缺失策略，默认 MissReturnID
	MissPolicy MissPolicy

	// Logger 日志接口（可选），兼容 slog/zap/logrus
	Logger Logger
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
	var missHandler func(lang, id string)
	var missPolicy MissPolicy
	var logger Logger

	if cfg != nil {
		if cfg.DefaultLang != "" {
			defaultLang = cfg.DefaultLang
		}
		if cfg.FallbackLang != "" {
			fallbackLang = cfg.FallbackLang
		}
		missHandler = cfg.MissHandler
		missPolicy = cfg.MissPolicy
		logger = cfg.Logger
	}

	tag := parseLanguageTag(defaultLang)

	b := &Bundle{
		bundle:       i18n.NewBundle(tag),
		currentLang:  defaultLang,
		defaultLang:  defaultLang,
		fallbackLang: fallbackLang,
		supported:    make([]string, 0),
		missHandler:  missHandler,
		missPolicy:   missPolicy,
		logger:       logger,
	}

	b.registerUnmarshalers()
	return b
}

// Init 初始化全局实例（替换默认实例）
func Init(cfg *Config) {
	once.Do(func() {}) // 确保 Default() 不会覆盖
	globalBundle = New(cfg)
}

// parseLanguageTag 解析语言标签，兼容多种格式
func parseLanguageTag(lang string) language.Tag {
	normalized := normalizeLanguageTag(lang)
	tag, err := language.Parse(normalized)
	if err != nil {
		return language.English
	}
	return tag
}

// normalizeLanguageTag 标准化语言标签: zh_CN -> zh-CN
func normalizeLanguageTag(lang string) string {
	return strings.ReplaceAll(lang, "_", "-")
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

	if loc, ok := b.localizers.Load(normalized); ok {
		return loc.(*i18n.Localizer)
	}

	loc := i18n.NewLocalizer(b.bundle, normalized, b.fallbackLang)
	b.localizers.Store(normalized, loc)
	return loc
}

// clearLocalizerCache 清空 Localizer 缓存
func (b *Bundle) clearLocalizerCache() {
	b.localizers = sync.Map{}
}

// handleMiss 处理翻译缺失
func (b *Bundle) handleMiss(lang, id string) {
	if b.missHandler != nil {
		b.missHandler(lang, id)
	}
	if b.logger != nil {
		b.logger.Warn("gi18n: missing translation", "lang", lang, "id", id)
	}
}

// GetBundle 获取底层的 go-i18n Bundle（高级用法）
func (b *Bundle) GetBundle() *i18n.Bundle {
	return b.bundle
}
