package gi18n

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// ========== 初始化测试 ==========

func TestNew_Default(t *testing.T) {
	b := New(nil)
	if b == nil {
		t.Fatal("New(nil) should not return nil")
	}
	if got := b.GetLang(); got != "en" {
		t.Errorf("expected default lang 'en', got '%s'", got)
	}
}

func TestNew_WithConfig(t *testing.T) {
	b := New(&Config{
		DefaultLang:  "zh-CN",
		FallbackLang: "ja",
	})
	if got := b.GetLang(); got != "zh-CN" {
		t.Errorf("expected lang 'zh-CN', got '%s'", got)
	}
}

func TestInit(t *testing.T) {
	Init(&Config{DefaultLang: "fr"})
	if got := Default().GetLang(); got != "fr" {
		t.Errorf("expected lang 'fr', got '%s'", got)
	}
	// 重置
	Init(&Config{DefaultLang: "en"})
}

// ========== 语言管理测试 ==========

func TestSetGetLang(t *testing.T) {
	b := New(nil)
	b.SetLang("zh-CN")
	if got := b.GetLang(); got != "zh-CN" {
		t.Errorf("expected 'zh-CN', got '%s'", got)
	}
}

func TestSetLang_NormalizeUnderscore(t *testing.T) {
	b := New(nil)
	b.SetLang("zh_CN")
	if got := b.GetLang(); got != "zh-CN" {
		t.Errorf("expected 'zh-CN', got '%s'", got)
	}
}

func TestLanguages(t *testing.T) {
	b := New(nil)
	_ = b.LoadMessages("en", map[string]string{"hello": "Hello"})
	_ = b.LoadMessages("zh-CN", map[string]string{"hello": "你好"})
	_ = b.LoadMessages("ja", map[string]string{"hello": "こんにちは"})

	langs := b.Languages()
	if len(langs) != 3 {
		t.Errorf("expected 3 languages, got %d: %v", len(langs), langs)
	}
}

func TestSetFallbackLang(t *testing.T) {
	b := New(nil)
	b.SetFallbackLang("zh-CN")
	// 验证不 panic，缓存被清除
}

func TestSetDefaultLang(t *testing.T) {
	b := New(nil)
	b.SetDefaultLang("ja")
	// 验证不 panic
}

// ========== normalizeLanguageTag 测试 ==========

func TestNormalizeLanguageTag(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"zh_CN", "zh-CN"},
		{"zh-CN", "zh-CN"},
		{"en", "en"},
		{"zh_TW", "zh-TW"},
		{"pt_BR", "pt-BR"},
		{"en_US", "en-US"},
	}

	for _, tt := range tests {
		if got := normalizeLanguageTag(tt.input); got != tt.expected {
			t.Errorf("normalizeLanguageTag(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

// ========== 核心翻译测试 (新 API) ==========

func TestT_Simple(t *testing.T) {
	b := New(&Config{DefaultLang: "en"})
	_ = b.LoadMessages("en", map[string]string{"hello": "Hello"})

	if got := b.T("hello"); got != "Hello" {
		t.Errorf("expected 'Hello', got '%s'", got)
	}
}

func TestT_WithLang(t *testing.T) {
	b := New(&Config{DefaultLang: "en"})
	_ = b.LoadMessages("en", map[string]string{"hello": "Hello"})
	_ = b.LoadMessages("zh-CN", map[string]string{"hello": "你好"})

	if got := b.T("hello", WithLang("zh-CN")); got != "你好" {
		t.Errorf("expected '你好', got '%s'", got)
	}
}

func TestT_WithData(t *testing.T) {
	b := New(&Config{DefaultLang: "en"})
	_ = b.LoadMessages("en", map[string]string{"greeting": "Hello, {{.Name}}!"})

	got := b.T("greeting", WithData("Name", "Alice"))
	if got != "Hello, Alice!" {
		t.Errorf("expected 'Hello, Alice!', got '%s'", got)
	}
}

func TestT_WithData_Multiple(t *testing.T) {
	b := New(&Config{DefaultLang: "en"})
	_ = b.LoadMessages("en", map[string]string{"info": "{{.Name}} is {{.Age}}"})

	got := b.T("info", WithData("Name", "Bob", "Age", 25))
	if got != "Bob is 25" {
		t.Errorf("expected 'Bob is 25', got '%s'", got)
	}
}

func TestT_WithMap(t *testing.T) {
	b := New(&Config{DefaultLang: "en"})
	_ = b.LoadMessages("en", map[string]string{"greeting": "Hello, {{.Name}}!"})

	got := b.T("greeting", WithMap(map[string]interface{}{"Name": "Bob"}))
	if got != "Hello, Bob!" {
		t.Errorf("expected 'Hello, Bob!', got '%s'", got)
	}
}

func TestT_WithCount(t *testing.T) {
	b := New(&Config{DefaultLang: "en"})
	content := []byte(`{
		"items": {
			"one": "{{.Count}} item",
			"other": "{{.Count}} items"
		}
	}`)
	_ = b.LoadContent("en", "json", content)

	if got := b.T("items", WithCount(1)); got != "1 item" {
		t.Errorf("expected '1 item', got '%s'", got)
	}
	if got := b.T("items", WithCount(5)); got != "5 items" {
		t.Errorf("expected '5 items', got '%s'", got)
	}
	if got := b.T("items", WithCount(0)); got != "0 items" {
		t.Errorf("expected '0 items', got '%s'", got)
	}
}

func TestT_WithContext(t *testing.T) {
	b := New(&Config{DefaultLang: "en"})
	_ = b.LoadMessages("en", map[string]string{"hello": "Hello"})
	_ = b.LoadMessages("zh-CN", map[string]string{"hello": "你好"})

	ctx := ContextWithLang(context.Background(), "zh-CN")
	if got := b.T("hello", WithContext(ctx)); got != "你好" {
		t.Errorf("expected '你好', got '%s'", got)
	}
}

func TestT_WithContext_NoLang(t *testing.T) {
	b := New(&Config{DefaultLang: "en"})
	_ = b.LoadMessages("en", map[string]string{"hello": "Hello"})

	// Context 中没有语言信息，应使用 Bundle 的当前语言
	if got := b.T("hello", WithContext(context.Background())); got != "Hello" {
		t.Errorf("expected 'Hello', got '%s'", got)
	}
}

func TestT_Combined_LangAndData(t *testing.T) {
	b := New(&Config{DefaultLang: "en"})
	_ = b.LoadMessages("en", map[string]string{"greeting": "Hello, {{.Name}}!"})
	_ = b.LoadMessages("zh-CN", map[string]string{"greeting": "你好，{{.Name}}！"})

	got := b.T("greeting", WithLang("zh-CN"), WithData("Name", "张三"))
	if got != "你好，张三！" {
		t.Errorf("expected '你好，张三！', got '%s'", got)
	}
}

func TestT_Combined_LangAndCount(t *testing.T) {
	b := New(&Config{DefaultLang: "en"})
	content := []byte(`{"items": {"one": "{{.Count}} item", "other": "{{.Count}} items"}}`)
	_ = b.LoadContent("en", "json", content)
	zhContent := []byte(`{"items": "{{.Count}} 个项目"}`)
	_ = b.LoadContent("zh-CN", "json", zhContent)

	if got := b.T("items", WithLang("zh-CN"), WithCount(10)); got != "10 个项目" {
		t.Errorf("expected '10 个项目', got '%s'", got)
	}
}

func TestT_LangOverridesContext(t *testing.T) {
	b := New(&Config{DefaultLang: "en"})
	_ = b.LoadMessages("en", map[string]string{"hello": "Hello"})
	_ = b.LoadMessages("zh-CN", map[string]string{"hello": "你好"})
	_ = b.LoadMessages("ja", map[string]string{"hello": "こんにちは"})

	ctx := ContextWithLang(context.Background(), "zh-CN")
	// WithLang 优先级 > WithContext
	if got := b.T("hello", WithContext(ctx), WithLang("ja")); got != "こんにちは" {
		t.Errorf("expected 'こんにちは', got '%s'", got)
	}
}

func TestT_Missing_ReturnID(t *testing.T) {
	b := New(nil)
	if got := b.T("nonexistent"); got != "nonexistent" {
		t.Errorf("expected 'nonexistent', got '%s'", got)
	}
}

func TestT_Missing_ReturnEmpty(t *testing.T) {
	b := New(&Config{MissPolicy: MissReturnEmpty})
	if got := b.T("nonexistent"); got != "" {
		t.Errorf("expected empty string, got '%s'", got)
	}
}

func TestT_MissHandler(t *testing.T) {
	var missedLang, missedID string
	b := New(&Config{
		MissHandler: func(lang, id string) {
			missedLang = lang
			missedID = id
		},
	})

	b.T("some.key")

	if missedID != "some.key" {
		t.Errorf("expected missedID 'some.key', got '%s'", missedID)
	}
	if missedLang != "en" {
		t.Errorf("expected missedLang 'en', got '%s'", missedLang)
	}
}

// ========== Logger 测试 ==========

type testLogger struct {
	warnings []string
}

func (l *testLogger) Warn(msg string, args ...any) {
	l.warnings = append(l.warnings, msg)
}

func TestLogger(t *testing.T) {
	logger := &testLogger{}
	b := New(&Config{Logger: logger})

	b.T("nonexistent")

	if len(logger.warnings) != 1 {
		t.Errorf("expected 1 warning, got %d", len(logger.warnings))
	}
	if len(logger.warnings) > 0 && !strings.Contains(logger.warnings[0], "missing translation") {
		t.Errorf("expected warning about missing translation, got '%s'", logger.warnings[0])
	}
}

func TestLogger_NotCalledOnSuccess(t *testing.T) {
	logger := &testLogger{}
	b := New(&Config{Logger: logger})
	_ = b.LoadMessages("en", map[string]string{"hello": "Hello"})

	b.T("hello")

	if len(logger.warnings) != 0 {
		t.Errorf("expected 0 warnings on success, got %d", len(logger.warnings))
	}
}

// ========== 加载测试 ==========

func TestLoadMessages(t *testing.T) {
	b := New(nil)
	err := b.LoadMessages("en", map[string]string{
		"hello": "Hello",
		"world": "World",
	})
	if err != nil {
		t.Fatalf("LoadMessages failed: %v", err)
	}

	if got := b.T("hello"); got != "Hello" {
		t.Errorf("expected 'Hello', got '%s'", got)
	}
	if got := b.T("world"); got != "World" {
		t.Errorf("expected 'World', got '%s'", got)
	}
}

func TestLoadContent_JSON(t *testing.T) {
	b := New(nil)
	data := []byte(`{"hello": "你好", "world": "世界"}`)
	err := b.LoadContent("zh-CN", "json", data)
	if err != nil {
		t.Fatalf("LoadContent failed: %v", err)
	}

	if got := b.T("hello", WithLang("zh-CN")); got != "你好" {
		t.Errorf("expected '你好', got '%s'", got)
	}
}

func TestLoadContent_Nested(t *testing.T) {
	b := New(nil)
	data := []byte(`{
		"common": {
			"confirm": "确定",
			"cancel": "取消"
		},
		"user": {
			"profile": {
				"title": "个人资料"
			}
		}
	}`)
	err := b.LoadContent("zh-CN", "json", data)
	if err != nil {
		t.Fatalf("LoadContent failed: %v", err)
	}

	tests := map[string]string{
		"common.confirm":     "确定",
		"common.cancel":      "取消",
		"user.profile.title": "个人资料",
	}

	for key, expected := range tests {
		if got := b.T(key, WithLang("zh-CN")); got != expected {
			t.Errorf("T(%q) = %q, want %q", key, got, expected)
		}
	}
}

func TestLoadContent_YAML(t *testing.T) {
	b := New(nil)
	data := []byte("hello: 你好\nworld: 世界\n")
	err := b.LoadContent("zh-CN", "yaml", data)
	if err != nil {
		t.Fatalf("LoadContent YAML failed: %v", err)
	}

	if got := b.T("hello", WithLang("zh-CN")); got != "你好" {
		t.Errorf("expected '你好', got '%s'", got)
	}
}

func TestLoadContent_Plural(t *testing.T) {
	b := New(nil)
	data := []byte(`{
		"items": {
			"one": "{{.Count}} item",
			"other": "{{.Count}} items"
		}
	}`)
	err := b.LoadContent("en", "json", data)
	if err != nil {
		t.Fatalf("LoadContent failed: %v", err)
	}

	if got := b.T("items", WithCount(1)); got != "1 item" {
		t.Errorf("expected '1 item', got '%s'", got)
	}
	if got := b.T("items", WithCount(3)); got != "3 items" {
		t.Errorf("expected '3 items', got '%s'", got)
	}
}

// ========== Context 测试 ==========

func TestContextWithLang(t *testing.T) {
	ctx := ContextWithLang(context.Background(), "zh-CN")
	lang := LangFromContext(ctx)
	if lang != "zh-CN" {
		t.Errorf("expected 'zh-CN', got '%s'", lang)
	}
}

func TestContextWithLang_Normalize(t *testing.T) {
	ctx := ContextWithLang(context.Background(), "zh_TW")
	lang := LangFromContext(ctx)
	if lang != "zh-TW" {
		t.Errorf("expected 'zh-TW', got '%s'", lang)
	}
}

func TestLangFromContext_NoLang(t *testing.T) {
	// 确保全局实例已初始化
	Init(&Config{DefaultLang: "en"})
	lang := LangFromContext(context.Background())
	if lang != "en" {
		t.Errorf("expected 'en', got '%s'", lang)
	}
}

// ========== Middleware 测试 ==========

func TestMiddleware_QueryParam(t *testing.T) {
	Init(&Config{DefaultLang: "en"})
	_ = LoadMessages("en", map[string]string{"hello": "Hello"})
	_ = LoadMessages("zh-CN", map[string]string{"hello": "你好"})

	handler := Middleware(nil)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		msg := T("hello", WithContext(r.Context()))
		w.Write([]byte(msg))
	}))

	req := httptest.NewRequest("GET", "/?lang=zh-CN", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if got := w.Body.String(); got != "你好" {
		t.Errorf("expected '你好', got '%s'", got)
	}
}

func TestMiddleware_AcceptLanguage(t *testing.T) {
	Init(&Config{DefaultLang: "en"})
	_ = LoadMessages("en", map[string]string{"hello": "Hello"})
	_ = LoadMessages("zh-CN", map[string]string{"hello": "你好"})

	handler := Middleware(nil)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		msg := T("hello", WithContext(r.Context()))
		w.Write([]byte(msg))
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if got := w.Body.String(); got != "你好" {
		t.Errorf("expected '你好', got '%s'", got)
	}
}

func TestMiddleware_Cookie(t *testing.T) {
	Init(&Config{DefaultLang: "en"})
	_ = LoadMessages("en", map[string]string{"hello": "Hello"})
	_ = LoadMessages("zh-CN", map[string]string{"hello": "你好"})

	handler := Middleware(nil)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		msg := T("hello", WithContext(r.Context()))
		w.Write([]byte(msg))
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{Name: "lang", Value: "zh-CN"})
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if got := w.Body.String(); got != "你好" {
		t.Errorf("expected '你好', got '%s'", got)
	}
}

func TestMiddleware_CustomConfig(t *testing.T) {
	Init(&Config{DefaultLang: "en"})
	_ = LoadMessages("en", map[string]string{"hello": "Hello"})
	_ = LoadMessages("zh-CN", map[string]string{"hello": "你好"})

	cfg := &MiddlewareConfig{
		Sources:    []LangSource{SourceQuery},
		QueryParam: "locale",
	}

	handler := Middleware(cfg)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		msg := T("hello", WithContext(r.Context()))
		w.Write([]byte(msg))
	}))

	req := httptest.NewRequest("GET", "/?locale=zh-CN", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if got := w.Body.String(); got != "你好" {
		t.Errorf("expected '你好', got '%s'", got)
	}
}

// ========== 已废弃方法兼容性测试 ==========

func TestDeprecated_TL(t *testing.T) {
	b := New(nil)
	_ = b.LoadMessages("en", map[string]string{"hello": "Hello"})
	_ = b.LoadMessages("zh-CN", map[string]string{"hello": "你好"})

	if got := b.TL("zh-CN", "hello"); got != "你好" {
		t.Errorf("expected '你好', got '%s'", got)
	}
}

func TestDeprecated_Tf(t *testing.T) {
	b := New(nil)
	_ = b.LoadMessages("en", map[string]string{"greeting": "Hello, {{.Name}}!"})

	if got := b.Tf("greeting", "Name", "Test"); got != "Hello, Test!" {
		t.Errorf("expected 'Hello, Test!', got '%s'", got)
	}
}

func TestDeprecated_TLf(t *testing.T) {
	b := New(nil)
	_ = b.LoadMessages("zh-CN", map[string]string{"greeting": "你好，{{.Name}}！"})

	if got := b.TLf("zh-CN", "greeting", "Name", "测试"); got != "你好，测试！" {
		t.Errorf("expected '你好，测试！', got '%s'", got)
	}
}

func TestDeprecated_Tp(t *testing.T) {
	b := New(nil)
	content := []byte(`{"items": {"one": "{{.Count}} item", "other": "{{.Count}} items"}}`)
	_ = b.LoadContent("en", "json", content)

	if got := b.Tp("items", 1); got != "1 item" {
		t.Errorf("expected '1 item', got '%s'", got)
	}
	if got := b.Tp("items", 5); got != "5 items" {
		t.Errorf("expected '5 items', got '%s'", got)
	}
}

func TestDeprecated_TMap(t *testing.T) {
	b := New(nil)
	_ = b.LoadMessages("en", map[string]string{"greeting": "Hello, {{.Name}}!"})

	got := b.TMap("greeting", map[string]interface{}{"Name": "Map"})
	if got != "Hello, Map!" {
		t.Errorf("expected 'Hello, Map!', got '%s'", got)
	}
}

func TestDeprecated_TC(t *testing.T) {
	Init(&Config{DefaultLang: "en"})
	_ = LoadMessages("en", map[string]string{"hello": "Hello"})
	_ = LoadMessages("zh-CN", map[string]string{"hello": "你好"})

	ctx := ContextWithLang(context.Background(), "zh-CN")
	if got := TC(ctx, "hello"); got != "你好" {
		t.Errorf("expected '你好', got '%s'", got)
	}
}

func TestDeprecated_SetLanguage(t *testing.T) {
	b := New(nil)
	b.SetLanguage("ja")
	if got := b.GetLanguage(); got != "ja" {
		t.Errorf("expected 'ja', got '%s'", got)
	}
}

func TestDeprecated_Langs(t *testing.T) {
	b := New(nil)
	_ = b.LoadMessages("en", map[string]string{"hello": "Hello"})
	_ = b.LoadMessages("ja", map[string]string{"hello": "こんにちは"})

	langs := b.Langs()
	getLangs := b.GetLanguages()

	if len(langs) != len(getLangs) {
		t.Errorf("Langs() and GetLanguages() should return same result")
	}
}

// ========== 多实例隔离测试 ==========

func TestMultiInstance_Isolation(t *testing.T) {
	b1 := New(&Config{DefaultLang: "en"})
	b2 := New(&Config{DefaultLang: "zh-CN"})

	_ = b1.LoadMessages("en", map[string]string{"app": "App1"})
	_ = b2.LoadMessages("zh-CN", map[string]string{"app": "应用2"})

	if got := b1.T("app"); got != "App1" {
		t.Errorf("b1.T('app') = '%s', want 'App1'", got)
	}
	if got := b2.T("app"); got != "应用2" {
		t.Errorf("b2.T('app') = '%s', want '应用2'", got)
	}
}

func TestMultiInstance_LangIsolation(t *testing.T) {
	b1 := New(&Config{DefaultLang: "en"})
	b2 := New(&Config{DefaultLang: "zh-CN"})

	b1.SetLang("ja")

	if b1.GetLang() != "ja" {
		t.Errorf("b1 lang should be 'ja'")
	}
	if b2.GetLang() != "zh-CN" {
		t.Errorf("b2 lang should be 'zh-CN', got '%s'", b2.GetLang())
	}
}

// ========== GetBundle 测试 ==========

func TestGetBundle(t *testing.T) {
	b := New(nil)
	if b.GetBundle() == nil {
		t.Error("GetBundle() should not return nil")
	}
}

// ========== isMessageObject 测试 ==========

func TestIsMessageObject(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]interface{}
		expected bool
	}{
		{
			name:     "plural with other",
			input:    map[string]interface{}{"other": "items"},
			expected: true,
		},
		{
			name:     "plural with one and other",
			input:    map[string]interface{}{"one": "item", "other": "items"},
			expected: true,
		},
		{
			name:     "with description",
			input:    map[string]interface{}{"description": "desc"},
			expected: true,
		},
		{
			name:     "with hash",
			input:    map[string]interface{}{"hash": "abc123"},
			expected: true,
		},
		{
			name:     "nested namespace",
			input:    map[string]interface{}{"confirm": "确定", "cancel": "取消"},
			expected: false,
		},
		{
			name:     "empty",
			input:    map[string]interface{}{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isMessageObject(tt.input); got != tt.expected {
				t.Errorf("isMessageObject(%v) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

// ========== parseAcceptLanguage 测试 ==========

func TestParseAcceptLanguage(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"zh-CN,zh;q=0.9,en;q=0.8", "zh-CN"},
		{"en-US,en;q=0.5", "en-US"},
		{"", ""},
		{"ja", "ja"},
		{"fr;q=0.9", "fr"},
	}

	for _, tt := range tests {
		if got := parseAcceptLanguage(tt.input); got != tt.expected {
			t.Errorf("parseAcceptLanguage(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}
