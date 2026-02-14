package gi18n

import "context"

// Option 翻译选项，用于配置 T() 方法的行为
type Option func(*translateConfig)

// translateConfig 翻译内部配置
type translateConfig struct {
	lang  string
	data  map[string]interface{}
	count *int
	ctx   context.Context
}

// WithLang 指定翻译目标语言
//
//	gi18n.T("confirm", gi18n.WithLang("zh-CN"))
func WithLang(lang string) Option {
	return func(c *translateConfig) {
		c.lang = lang
	}
}

// WithData 设置模板参数（key-value 对）
//
//	gi18n.T("greeting", gi18n.WithData("Name", "张三"))
//	gi18n.T("info", gi18n.WithData("Name", "张三", "Age", 18))
func WithData(kv ...interface{}) Option {
	return func(c *translateConfig) {
		if c.data == nil {
			c.data = make(map[string]interface{})
		}
		for i := 0; i+1 < len(kv); i += 2 {
			if key, ok := kv[i].(string); ok {
				c.data[key] = kv[i+1]
			}
		}
	}
}

// WithMap 使用 map 设置模板参数
//
//	gi18n.T("greeting", gi18n.WithMap(map[string]interface{}{"Name": "张三"}))
func WithMap(data map[string]interface{}) Option {
	return func(c *translateConfig) {
		c.data = data
	}
}

// WithCount 设置复数计数
//
//	gi18n.T("items", gi18n.WithCount(5))
func WithCount(n int) Option {
	return func(c *translateConfig) {
		c.count = &n
	}
}

// WithContext 从 context.Context 获取语言设置
//
//	ctx := gi18n.ContextWithLang(ctx, "zh-CN")
//	gi18n.T("hello", gi18n.WithContext(ctx))
func WithContext(ctx context.Context) Option {
	return func(c *translateConfig) {
		c.ctx = ctx
	}
}
