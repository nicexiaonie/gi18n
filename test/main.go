package main

import (
	"context"
	"fmt"

	"github.com/nicexiaonie/gi18n"
)

// printSeparator 打印分隔线
func printSeparator(title string) {
	fmt.Printf("\n========== %s ==========\n", title)
}

// printResult 打印非翻译类测试结果
func printResult(name, result string) {
	fmt.Printf("  %-38s => %s\n", name, result)
}

// printTrans 以直观格式打印翻译测试结果
// 格式: [语言] "key" {参数}  →  结果
func printTrans(lang, key, params, result string) {
	tag := fmt.Sprintf("[%s]", lang)
	keyStr := fmt.Sprintf("%q", key)
	var desc string
	if params != "" {
		desc = fmt.Sprintf("%-8s %-20s %s", tag, keyStr, params)
	} else {
		desc = fmt.Sprintf("%-8s %s", tag, keyStr)
	}
	fmt.Printf("  %-50s →  %s\n", desc, result)
}

func main() {
	fmt.Println("=================================================")
	fmt.Println("           gi18n 国际化工具测试")
	fmt.Println("           (新 API + 向后兼容)")
	fmt.Println("=================================================")

	// ========== 1. 初始化和配置测试 ==========
	printSeparator("1. 初始化和配置")

	bundle := gi18n.Default()
	printResult("Default() 创建全局实例", "成功")

	customBundle := gi18n.New(&gi18n.Config{
		DefaultLang:  "zh-CN",
		FallbackLang: "en",
	})
	printResult("New() 创建自定义实例", "成功")

	// 测试 MissHandler 和 Logger
	missBundle := gi18n.New(&gi18n.Config{
		DefaultLang: "en",
		MissHandler: func(lang, id string) {
			fmt.Printf("  [MissHandler] lang=%s, id=%s\n", lang, id)
		},
		MissPolicy: gi18n.MissReturnID,
	})
	printResult("New() with MissHandler", "成功")

	gi18n.Init(&gi18n.Config{
		DefaultLang:  "en",
		FallbackLang: "en",
	})
	printResult("Init() 全局初始化", "成功")

	// ========== 2. 加载语言文件 ==========
	printSeparator("2. 加载语言文件")

	testMessagesEN := map[string]string{
		"confirm":      "OK",
		"cancel":       "Cancel",
		"greeting":     "Hello, {{.Name}}!",
		"items":        "{{.Count}} items",
		"search":       "Search",
		"login":        "Login",
		"test.hello":   "Hello Test",
		"test.goodbye": "Goodbye Test",
	}
	testMessagesZH := map[string]string{
		"confirm":  "确定",
		"cancel":   "取消",
		"greeting": "你好，{{.Name}}！",
		"items":    "{{.Count}} 个项目",
		"search":   "搜索",
		"login":    "登录",
	}

	if err := gi18n.LoadMessages("en", testMessagesEN); err != nil {
		fmt.Printf("  LoadMessages(en) 失败: %v\n", err)
	} else {
		printResult("LoadMessages(en) 加载英文", "成功")
	}

	if err := gi18n.LoadMessages("zh-CN", testMessagesZH); err != nil {
		fmt.Printf("  LoadMessages(zh-CN) 失败: %v\n", err)
	} else {
		printResult("LoadMessages(zh-CN) 加载中文", "成功")
	}

	jsonContent := []byte(`{
		"test.content": {
			"id": "test.content",
			"other": "Content Test"
		}
	}`)
	if err := gi18n.LoadContent("en", "json", jsonContent); err != nil {
		fmt.Printf("  LoadContent() 失败: %v\n", err)
	} else {
		printResult("LoadContent() 从字节加载", "成功")
	}

	if err := customBundle.LoadMessages("zh-CN", map[string]string{
		"custom.test": "自定义测试",
	}); err != nil {
		fmt.Printf("  自定义实例加载失败: %v\n", err)
	} else {
		printResult("自定义实例 LoadMessages()", "成功")
	}

	// ========== 3. 新 API: T() + Option ==========
	printSeparator("3. 新 API: T() + Option 组合")

	gi18n.SetLang("en")
	fmt.Println("  格式: [语言] \"key\" {参数}  →  翻译结果")
	fmt.Println()

	// 简单翻译
	printTrans("en", "confirm", "", gi18n.T("confirm"))

	// 指定语言
	printTrans("zh-CN", "confirm", "(WithLang)", gi18n.T("confirm", gi18n.WithLang("zh-CN")))

	// 带参数
	printTrans("en", "greeting", "{Name=Alice}", gi18n.T("greeting", gi18n.WithData("Name", "Alice")))

	// 指定语言 + 参数
	printTrans("zh-CN", "greeting", "{Name=张三} (WithLang)", gi18n.T("greeting", gi18n.WithLang("zh-CN"), gi18n.WithData("Name", "张三")))

	// 复数
	printTrans("en", "items", "{count=1}", gi18n.T("items", gi18n.WithCount(1)))
	printTrans("en", "items", "{count=5}", gi18n.T("items", gi18n.WithCount(5)))

	// 指定语言 + 复数
	printTrans("zh-CN", "items", "{count=99} (WithLang)", gi18n.T("items", gi18n.WithLang("zh-CN"), gi18n.WithCount(99)))

	// Map 参数
	printTrans("en", "greeting", "{Name=Bob} (WithMap)", gi18n.T("greeting", gi18n.WithMap(map[string]interface{}{"Name": "Bob"})))

	// Context
	ctx := gi18n.ContextWithLang(context.Background(), "zh-CN")
	printTrans("zh-CN", "confirm", "(WithContext)", gi18n.T("confirm", gi18n.WithContext(ctx)))

	// Context + 参数
	printTrans("zh-CN", "greeting", "{Name=李华} (WithContext)", gi18n.T("greeting", gi18n.WithContext(ctx), gi18n.WithData("Name", "李华")))

	// ========== 4. MissHandler 测试 ==========
	printSeparator("4. MissHandler 测试")

	missBundle.LoadMessages("en", map[string]string{"hello": "Hello"})
	printTrans("en", "hello", "(已注册)", missBundle.T("hello"))

	fmt.Print("  触发 MissHandler: ")
	result := missBundle.T("nonexistent.key")
	printTrans("en", "nonexistent.key", "(未注册, MissReturnID)", result)

	// MissReturnEmpty 策略
	emptyBundle := gi18n.New(&gi18n.Config{
		DefaultLang: "en",
		MissPolicy:  gi18n.MissReturnEmpty,
	})
	emptyResult := emptyBundle.T("nonexistent")
	printTrans("en", "nonexistent", "(MissReturnEmpty)", fmt.Sprintf("%q  空=%v", emptyResult, emptyResult == ""))

	// ========== 5. 语言管理 ==========
	printSeparator("5. 语言管理")

	gi18n.SetLang("en")
	printResult("SetLang(\"en\")  → GetLang()", gi18n.GetLang())

	gi18n.SetLang("zh-CN")
	printResult("SetLang(\"zh-CN\")  → GetLang()", gi18n.GetLang())

	langs := gi18n.Languages()
	printResult("Languages()", fmt.Sprintf("%v", langs))

	gi18n.SetLang("zh_CN")
	printResult("SetLang(\"zh_CN\") 自动标准化  → GetLang()", gi18n.GetLang())

	// ========== 6. 多实例隔离 ==========
	printSeparator("6. 多实例隔离")

	instance1 := gi18n.New(&gi18n.Config{DefaultLang: "en"})
	instance2 := gi18n.New(&gi18n.Config{DefaultLang: "zh-CN"})

	instance1.LoadMessages("en", map[string]string{"app.name": "Application"})
	instance2.LoadMessages("zh-CN", map[string]string{"app.name": "应用程序"})

	printTrans("en", "app.name", "(instance1)", instance1.T("app.name"))
	printTrans("zh-CN", "app.name", "(instance2)", instance2.T("app.name"))
	printResult("instance1.GetLang()", instance1.GetLang())
	printResult("instance2.GetLang()", instance2.GetLang())

	// ========== 7. 向后兼容（已废弃方法） ==========
	printSeparator("7. 向后兼容（已废弃方法仍可用）")

	gi18n.SetLang("en")
	fmt.Println("  格式: [语言] \"key\" {参数}  →  翻译结果")
	fmt.Println()

	printTrans("en", "confirm", "[Deprecated] Translate()", gi18n.Translate("confirm"))
	printTrans("zh-CN", "confirm", "[Deprecated] TL(lang, key)", gi18n.TL("zh-CN", "confirm"))
	printTrans("en", "greeting", "{Name=Test} [Deprecated] Tf()", gi18n.Tf("greeting", "Name", "Test"))
	printTrans("zh-CN", "greeting", "{Name=测试} [Deprecated] TLf()", gi18n.TLf("zh-CN", "greeting", "Name", "测试"))
	printTrans("en", "items", "{count=5} [Deprecated] Tp()", gi18n.Tp("items", 5))
	printTrans("zh-CN", "items", "{count=10} [Deprecated] TLp()", gi18n.TLp("zh-CN", "items", 10))

	printTrans("en", "search", "(bundle) [Deprecated] Translate()", bundle.Translate("search"))
	printTrans("zh-CN", "search", "(bundle) [Deprecated] TL()", bundle.TL("zh-CN", "search"))

	ctxEN := gi18n.ContextWithLang(context.Background(), "en")
	ctxZH := gi18n.ContextWithLang(context.Background(), "zh-CN")

	printTrans("en", "login", "(ctxEN) [Deprecated] TC()", gi18n.TC(ctxEN, "login"))
	printTrans("zh-CN", "login", "(ctxZH) [Deprecated] TC()", gi18n.TC(ctxZH, "login"))
	printTrans("en", "greeting", "{Name=Frank} (ctxEN) [Deprecated] TCf()", gi18n.TCf(ctxEN, "greeting", "Name", "Frank"))
	printTrans("en", "items", "{count=7} (ctxEN) [Deprecated] TCp()", gi18n.TCp(ctxEN, "items", 7))

	// ========== 8. 边界情况 ==========
	printSeparator("8. 边界情况")

	gi18n.SetLang("zh-CN") // 当前语言是 zh-CN，but zh-CN 没有 not.exist，会 fallback 到 en
	printTrans("zh-CN", "not.exist", "(key 不存在)", gi18n.T("not.exist"))
	printTrans("en", "greeting", "(WithData 无参数)", gi18n.T("greeting", gi18n.WithData()))
	printTrans("en", "items", "{count=0}", gi18n.T("items", gi18n.WithCount(0)))
	printTrans("en", "items", "{count=-1}", gi18n.T("items", gi18n.WithCount(-1)))

	// ========== 9. 新 API 总结 ==========
	printSeparator("9. 新 API 总结")

	apis := []string{
		"核心翻译: T(id, opts...)  — 唯一翻译入口",
		"选项: WithLang(lang)     — 指定语言",
		"选项: WithData(kv...)    — 模板参数",
		"选项: WithMap(m)         — Map 参数",
		"选项: WithCount(n)       — 复数",
		"选项: WithContext(ctx)   — 从 Context 获取语言",
		"语言: SetLang / GetLang / Languages / SetFallbackLang",
		"加载: Load / LoadFS / LoadContent / LoadMessages",
		"中间件: Middleware(cfg)   — 标准 HTTP 中间件",
		"上下文: ContextWithLang / LangFromContext",
		"配置: Config { MissHandler, MissPolicy, Logger }",
	}

	for i, api := range apis {
		fmt.Printf("  %2d. %s\n", i+1, api)
	}

	// ========== 完成 ==========
	fmt.Println("\n=================================================")
	fmt.Println("           所有测试完成！")
	fmt.Println("=================================================")
}
