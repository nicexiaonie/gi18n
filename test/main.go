package main

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/nicexiaonie/gi18n"
)

// 辅助函数：打印分隔线
func printSeparator(title string) {
	fmt.Printf("\n========== %s ==========\n", title)
}

// 辅助函数：打印测试结果
func printResult(name, result string) {
	fmt.Printf("  %-40s => %s\n", name, result)
}

// 获取builtin目录的绝对路径
func getBuiltinDir() string {
	_, filename, _, _ := runtime.Caller(0)
	projectRoot := filepath.Dir(filepath.Dir(filename))
	return filepath.Join(projectRoot, "builtin")
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

	builtinDir := getBuiltinDir()
	fmt.Printf("  Builtin 目录: %s\n", builtinDir)

	if err := gi18n.Load(builtinDir); err != nil {
		fmt.Printf("  Load() 失败: %v\n", err)
	} else {
		printResult("Load() 从目录加载", "成功")
	}

	testMessages := map[string]string{
		"test.hello":   "Hello Test",
		"test.goodbye": "Goodbye Test",
	}
	if err := gi18n.LoadMessages("en", testMessages); err != nil {
		fmt.Printf("  LoadMessages() 失败: %v\n", err)
	} else {
		printResult("LoadMessages() 直接加载消息", "成功")
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

	// 简单翻译
	printResult("T(\"confirm\") [en]",
		gi18n.T("confirm"))

	// 指定语言
	printResult("T(\"confirm\", WithLang(\"zh-CN\"))",
		gi18n.T("confirm", gi18n.WithLang("zh-CN")))

	// 带参数
	printResult("T(\"greeting\", WithData(\"Name\", \"Alice\"))",
		gi18n.T("greeting", gi18n.WithData("Name", "Alice")))

	// 指定语言 + 参数
	printResult("T(\"greeting\", WithLang, WithData)",
		gi18n.T("greeting", gi18n.WithLang("zh-CN"), gi18n.WithData("Name", "张三")))

	// 复数
	printResult("T(\"items\", WithCount(1)) [en]",
		gi18n.T("items", gi18n.WithCount(1)))
	printResult("T(\"items\", WithCount(5)) [en]",
		gi18n.T("items", gi18n.WithCount(5)))

	// 指定语言 + 复数
	printResult("T(\"items\", WithLang, WithCount)",
		gi18n.T("items", gi18n.WithLang("zh-CN"), gi18n.WithCount(99)))

	// Map 参数
	printResult("T(\"greeting\", WithMap(...))",
		gi18n.T("greeting", gi18n.WithMap(map[string]interface{}{"Name": "Bob"})))

	// Context
	ctx := gi18n.ContextWithLang(context.Background(), "zh-CN")
	printResult("T(\"confirm\", WithContext(ctx))",
		gi18n.T("confirm", gi18n.WithContext(ctx)))

	// Context + 参数
	printResult("T(\"greeting\", WithContext, WithData)",
		gi18n.T("greeting", gi18n.WithContext(ctx), gi18n.WithData("Name", "李华")))

	// ========== 4. MissHandler 测试 ==========
	printSeparator("4. MissHandler 测试")

	missBundle.LoadMessages("en", map[string]string{"hello": "Hello"})
	printResult("missBundle.T(\"hello\")", missBundle.T("hello"))
	fmt.Print("  调用不存在的 key: ")
	result := missBundle.T("nonexistent.key")
	printResult("missBundle.T(\"nonexistent.key\")", result)

	// MissReturnEmpty 策略
	emptyBundle := gi18n.New(&gi18n.Config{
		DefaultLang: "en",
		MissPolicy:  gi18n.MissReturnEmpty,
	})
	emptyResult := emptyBundle.T("nonexistent")
	printResult("MissReturnEmpty 策略",
		fmt.Sprintf("'%s' (空字符串=%v)", emptyResult, emptyResult == ""))

	// ========== 5. 语言管理 ==========
	printSeparator("5. 语言管理")

	gi18n.SetLang("en")
	printResult("SetLang(\"en\")", gi18n.GetLang())

	gi18n.SetLang("zh-CN")
	printResult("SetLang(\"zh-CN\")", gi18n.GetLang())

	langs := gi18n.Languages()
	printResult("Languages()", fmt.Sprintf("%v", langs))

	// 测试下划线标准化
	gi18n.SetLang("zh_CN")
	printResult("SetLang(\"zh_CN\") 自动标准化", gi18n.GetLang())

	// ========== 6. 多实例隔离 ==========
	printSeparator("6. 多实例隔离")

	instance1 := gi18n.New(&gi18n.Config{DefaultLang: "en"})
	instance2 := gi18n.New(&gi18n.Config{DefaultLang: "zh-CN"})

	instance1.LoadMessages("en", map[string]string{"app.name": "Application"})
	instance2.LoadMessages("zh-CN", map[string]string{"app.name": "应用程序"})

	printResult("instance1.T(\"app.name\") [en]", instance1.T("app.name"))
	printResult("instance2.T(\"app.name\") [zh-CN]", instance2.T("app.name"))
	printResult("instance1.GetLang()", instance1.GetLang())
	printResult("instance2.GetLang()", instance2.GetLang())

	// ========== 7. 向后兼容（已废弃方法） ==========
	printSeparator("7. 向后兼容（已废弃方法仍可用）")

	gi18n.SetLang("en")
	printResult("[Deprecated] Translate(\"confirm\")", gi18n.Translate("confirm"))
	printResult("[Deprecated] TL(\"zh-CN\", \"confirm\")", gi18n.TL("zh-CN", "confirm"))
	printResult("[Deprecated] Tf(\"greeting\", ...)", gi18n.Tf("greeting", "Name", "Test"))
	printResult("[Deprecated] TLf(\"zh-CN\", \"greeting\", ...)",
		gi18n.TLf("zh-CN", "greeting", "Name", "测试"))
	printResult("[Deprecated] Tp(\"items\", 5)", gi18n.Tp("items", 5))
	printResult("[Deprecated] TLp(\"zh-CN\", \"items\", 10)",
		gi18n.TLp("zh-CN", "items", 10))

	// 实例方法向后兼容
	printResult("[Deprecated] bundle.Translate(\"search\")", bundle.Translate("search"))
	printResult("[Deprecated] bundle.TL(\"zh-CN\", \"search\")", bundle.TL("zh-CN", "search"))

	// Context 方法向后兼容
	ctxEN := gi18n.ContextWithLang(context.Background(), "en")
	ctxZH := gi18n.ContextWithLang(context.Background(), "zh-CN")

	printResult("[Deprecated] TC(ctxEN, \"login\")", gi18n.TC(ctxEN, "login"))
	printResult("[Deprecated] TC(ctxZH, \"login\")", gi18n.TC(ctxZH, "login"))
	printResult("[Deprecated] TCf(ctxEN, \"greeting\", ...)",
		gi18n.TCf(ctxEN, "greeting", "Name", "Frank"))
	printResult("[Deprecated] TCp(ctxEN, \"items\", 7)",
		gi18n.TCp(ctxEN, "items", 7))

	// ========== 8. 边界情况 ==========
	printSeparator("8. 边界情况")

	printResult("T(\"not.exist\") 不存在的 key", gi18n.T("not.exist"))
	printResult("T() 空参数 WithData", gi18n.T("greeting", gi18n.WithData()))
	printResult("T() WithCount(0)", gi18n.T("items", gi18n.WithCount(0)))
	printResult("T() WithCount(-1)", gi18n.T("items", gi18n.WithCount(-1)))

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
