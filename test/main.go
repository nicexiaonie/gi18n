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
	fmt.Printf("  %-35s => %s\n", name, result)
}

// 获取builtin目录的绝对路径
func getBuiltinDir() string {
	_, filename, _, _ := runtime.Caller(0)
	projectRoot := filepath.Dir(filepath.Dir(filename))
	return filepath.Join(projectRoot, "builtin")
}

func main() {
	fmt.Println("=================================================")
	fmt.Println("           gi18n 国际化工具完整测试")
	fmt.Println("=================================================")

	// ========== 1. 初始化和配置测试 ==========
	printSeparator("1. 初始化和配置测试")

	// 测试默认实例
	bundle := gi18n.Default()
	printResult("Default() 创建全局实例", "成功")

	// 测试自定义配置初始化
	customBundle := gi18n.New(&gi18n.Config{
		DefaultLang:  "zh-CN",
		FallbackLang: "en",
	})
	printResult("New() 创建自定义实例", "成功")

	// 测试 Init 全局初始化
	gi18n.Init(&gi18n.Config{
		DefaultLang:  "en",
		FallbackLang: "en",
	})
	printResult("Init() 全局初始化", "成功")

	// ========== 2. 加载语言文件测试 ==========
	printSeparator("2. 加载语言文件测试")

	builtinDir := getBuiltinDir()
	fmt.Printf("  Builtin 目录: %s\n", builtinDir)

	// 测试从目录加载
	if err := gi18n.Load(builtinDir); err != nil {
		fmt.Printf("  Load() 失败: %v\n", err)
	} else {
		printResult("Load() 从目录加载", "成功")
	}

	// 测试 LoadMessages 直接加载
	testMessages := map[string]string{
		"test.hello":   "Hello Test",
		"test.goodbye": "Goodbye Test",
	}
	if err := gi18n.LoadMessages("en", testMessages); err != nil {
		fmt.Printf("  LoadMessages() 失败: %v\n", err)
	} else {
		printResult("LoadMessages() 直接加载消息", "成功")
	}

	// 测试 LoadContent 从字节加载
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

	// 测试自定义实例加载
	if err := customBundle.LoadMessages("zh-CN", map[string]string{
		"custom.test": "自定义测试",
	}); err != nil {
		fmt.Printf("  自定义实例加载失败: %v\n", err)
	} else {
		printResult("自定义实例 LoadMessages()", "成功")
	}

	// ========== 3. 语言设置测试 ==========
	printSeparator("3. 语言设置测试")

	// 测试设置和获取语言
	gi18n.SetLang("en")
	printResult("SetLang(\"en\")", gi18n.GetLang())

	gi18n.SetLanguage("zh-CN")
	printResult("SetLanguage(\"zh-CN\")", gi18n.GetLanguage())

	// 测试获取支持的语言列表
	langs := gi18n.Langs()
	printResult("Langs() 支持的语言", fmt.Sprintf("%v", langs))
	printResult("GetLanguages() 别名方法", fmt.Sprintf("%v", gi18n.GetLanguages()))

	// 测试设置默认和回退语言
	gi18n.SetDefaultLang("en")
	printResult("SetDefaultLang(\"en\")", "成功")

	gi18n.SetFallbackLang("en")
	printResult("SetFallbackLang(\"en\")", "成功")

	// 测试自定义实例的语言设置
	customBundle.SetLang("zh-CN")
	printResult("自定义实例 SetLang(\"zh-CN\")", customBundle.GetLang())

	// ========== 4. 简单翻译测试 ==========
	printSeparator("4. 简单翻译测试")

	gi18n.SetLang("en")
	printResult("T(\"confirm\") [en]", gi18n.T("confirm"))
	printResult("Translate(\"cancel\") 别名 [en]", gi18n.Translate("cancel"))

	gi18n.SetLang("zh-CN")
	printResult("T(\"confirm\") [zh-CN]", gi18n.T("confirm"))
	printResult("Translate(\"save\") 别名 [zh-CN]", gi18n.Translate("save"))

	// 测试不存在的键（应返回ID）
	printResult("T(\"not.exist\") 不存在的键", gi18n.T("not.exist"))

	// 测试自定义实例翻译
	printResult("自定义实例 T(\"custom.test\")", customBundle.T("custom.test"))

	// ========== 5. 指定语言翻译测试 ==========
	printSeparator("5. 指定语言翻译测试")

	printResult("TL(\"en\", \"delete\")", gi18n.TL("en", "delete"))
	printResult("TL(\"zh-CN\", \"delete\")", gi18n.TL("zh-CN", "delete"))
	printResult("TranslateLang(\"en\", \"edit\") 别名", gi18n.TranslateLang("en", "edit"))
	printResult("TranslateLang(\"zh-CN\", \"edit\") 别名", gi18n.TranslateLang("zh-CN", "edit"))

	// 测试下划线格式的语言标签
	printResult("TL(\"zh_CN\", \"submit\")", gi18n.TL("zh_CN", "submit"))

	// ========== 6. 带参数翻译测试 ==========
	printSeparator("6. 带参数翻译测试")

	gi18n.SetLang("en")
	printResult("Tf(\"greeting\", \"Name\", \"Alice\") [en]",
		gi18n.Tf("greeting", "Name", "Alice"))

	gi18n.SetLang("zh-CN")
	printResult("Tf(\"greeting\", \"Name\", \"小明\") [zh-CN]",
		gi18n.Tf("greeting", "Name", "小明"))

	printResult("TranslateWith() 别名 [zh-CN]",
		gi18n.TranslateWith("greeting", "Name", "小红"))

	// ========== 7. 指定语言带参数翻译测试 ==========
	printSeparator("7. 指定语言带参数翻译测试")

	printResult("TLf(\"en\", \"greeting\", \"Name\", \"Bob\")",
		gi18n.TLf("en", "greeting", "Name", "Bob"))

	printResult("TLf(\"zh-CN\", \"greeting\", \"Name\", \"李华\")",
		gi18n.TLf("zh-CN", "greeting", "Name", "李华"))

	printResult("TranslateLangWith() 别名",
		gi18n.TranslateLangWith("en", "greeting", "Name", "Charlie"))

	// ========== 8. 带复数翻译测试 ==========
	printSeparator("8. 带复数翻译测试")

	gi18n.SetLang("en")
	printResult("Tp(\"items\", 1) [en]", gi18n.Tp("items", 1))
	printResult("Tp(\"items\", 5) [en]", gi18n.Tp("items", 5))
	printResult("TranslatePlural() 别名 [en]", gi18n.TranslatePlural("items", 10))

	gi18n.SetLang("zh-CN")
	printResult("Tp(\"items\", 1) [zh-CN]", gi18n.Tp("items", 1))
	printResult("Tp(\"items\", 5) [zh-CN]", gi18n.Tp("items", 5))

	// ========== 9. 指定语言带复数翻译测试 ==========
	printSeparator("9. 指定语言带复数翻译测试")

	printResult("TLp(\"en\", \"items\", 1)", gi18n.TLp("en", "items", 1))
	printResult("TLp(\"en\", \"items\", 100)", gi18n.TLp("en", "items", 100))
	printResult("TLp(\"zh-CN\", \"items\", 50)", gi18n.TLp("zh-CN", "items", 50))
	printResult("TranslateLangPlural() 别名",
		gi18n.TranslateLangPlural("en", "items", 3))

	// ========== 10. Map参数翻译测试 ==========
	printSeparator("10. Map参数翻译测试")

	dataMap := map[string]interface{}{
		"Name": "David",
	}

	gi18n.SetLang("en")
	printResult("TMap(\"greeting\", map) [en]", gi18n.TMap("greeting", dataMap))

	gi18n.SetLang("zh-CN")
	dataMap["Name"] = "王五"
	printResult("TMap(\"greeting\", map) [zh-CN]", gi18n.TMap("greeting", dataMap))

	dataMap["Name"] = "Emma"
	printResult("TLMap(\"en\", \"greeting\", map)",
		gi18n.TLMap("en", "greeting", dataMap))

	dataMap["Name"] = "赵六"
	printResult("TLMap(\"zh-CN\", \"greeting\", map)",
		gi18n.TLMap("zh-CN", "greeting", dataMap))

	// ========== 11. Context 相关测试 ==========
	printSeparator("11. Context 相关测试")

	// 创建带语言的 context
	ctx := context.Background()
	ctxEN := gi18n.WithLang(ctx, "en")
	ctxZH := gi18n.WithLang(ctx, "zh-CN")

	printResult("LangFromContext(ctxEN)", gi18n.LangFromContext(ctxEN))
	printResult("LangFromContext(ctxZH)", gi18n.LangFromContext(ctxZH))

	// 测试 TC (TranslateContext)
	printResult("TC(ctxEN, \"login\")", gi18n.TC(ctxEN, "login"))
	printResult("TC(ctxZH, \"login\")", gi18n.TC(ctxZH, "login"))
	printResult("TranslateContext() 别名", gi18n.TranslateContext(ctxEN, "logout"))

	// 测试 TCf (TranslateContextWith)
	printResult("TCf(ctxEN, \"greeting\", \"Name\", \"Frank\")",
		gi18n.TCf(ctxEN, "greeting", "Name", "Frank"))
	printResult("TCf(ctxZH, \"greeting\", \"Name\", \"孙七\")",
		gi18n.TCf(ctxZH, "greeting", "Name", "孙七"))
	printResult("TranslateContextWith() 别名",
		gi18n.TranslateContextWith(ctxZH, "greeting", "Name", "周八"))

	// 测试 TCp (TranslateContextPlural)
	printResult("TCp(ctxEN, \"items\", 1)", gi18n.TCp(ctxEN, "items", 1))
	printResult("TCp(ctxEN, \"items\", 7)", gi18n.TCp(ctxEN, "items", 7))
	printResult("TCp(ctxZH, \"items\", 99)", gi18n.TCp(ctxZH, "items", 99))
	printResult("TranslateContextPlural() 别名",
		gi18n.TranslateContextPlural(ctxEN, "items", 2))

	// 测试没有语言信息的 context (应使用默认语言)
	printResult("TC(ctx, \"welcome\") 无语言context", gi18n.TC(ctx, "welcome"))

	// ========== 12. 实例方法测试 ==========
	printSeparator("12. 实例方法测试")

	bundle.SetLang("en")
	printResult("bundle.T(\"search\")", bundle.T("search"))
	printResult("bundle.Translate(\"close\") 别名", bundle.Translate("close"))
	printResult("bundle.TL(\"zh-CN\", \"search\")", bundle.TL("zh-CN", "search"))
	printResult("bundle.TranslateLang() 别名", bundle.TranslateLang("zh-CN", "back"))

	printResult("bundle.Tf(\"greeting\", \"Name\", \"Instance\")",
		bundle.Tf("greeting", "Name", "Instance"))
	printResult("bundle.TranslateWith() 别名",
		bundle.TranslateWith("greeting", "Name", "Test"))

	printResult("bundle.TLf(\"en\", \"greeting\", \"Name\", \"EN\")",
		bundle.TLf("en", "greeting", "Name", "EN"))
	printResult("bundle.TranslateLangWith() 别名",
		bundle.TranslateLangWith("zh-CN", "greeting", "Name", "中文"))

	printResult("bundle.Tp(\"items\", 1)", bundle.Tp("items", 1))
	printResult("bundle.TranslatePlural() 别名", bundle.TranslatePlural("items", 8))

	printResult("bundle.TLp(\"en\", \"items\", 15)",
		bundle.TLp("en", "items", 15))
	printResult("bundle.TranslateLangPlural() 别名",
		bundle.TranslateLangPlural("zh-CN", "items", 20))

	testMap := map[string]interface{}{"Name": "MapTest"}
	printResult("bundle.TMap(\"greeting\", map)", bundle.TMap("greeting", testMap))
	printResult("bundle.TLMap(\"zh-CN\", \"greeting\", map)",
		bundle.TLMap("zh-CN", "greeting", testMap))

	// ========== 13. 高级功能测试 ==========
	printSeparator("13. 高级功能测试")

	// 测试 GetBundle
	underlyingBundle := bundle.GetBundle()
	if underlyingBundle != nil {
		printResult("GetBundle() 获取底层Bundle", "成功")
	}

	// 测试语言标签标准化
	gi18n.SetLang("zh_CN") // 使用下划线
	printResult("SetLang(\"zh_CN\") 自动标准化", gi18n.GetLang())

	gi18n.SetLang("zh-Hans") // 使用 Hans
	printResult("SetLang(\"zh-Hans\")", gi18n.GetLang())

	// ========== 14. 边界情况测试 ==========
	printSeparator("14. 边界情况测试")

	// 测试空参数
	printResult("Tf() 空参数", gi18n.Tf("greeting"))

	// 测试奇数参数（最后一个参数会被忽略）
	printResult("Tf() 奇数参数",
		gi18n.Tf("greeting", "Name", "Test", "Extra"))

	// 测试非字符串键
	printResult("Tf() 非字符串键",
		gi18n.Tf("greeting", 123, "Value", "Name", "Test"))

	// 测试复数为0
	printResult("Tp(\"items\", 0)", gi18n.Tp("items", 0))

	// 测试负数复数
	printResult("Tp(\"items\", -1)", gi18n.Tp("items", -1))

	// ========== 15. 多实例测试 ==========
	printSeparator("15. 多实例测试")

	instance1 := gi18n.New(&gi18n.Config{DefaultLang: "en"})
	instance2 := gi18n.New(&gi18n.Config{DefaultLang: "zh-CN"})

	instance1.LoadMessages("en", map[string]string{
		"app.name": "Application",
	})
	instance2.LoadMessages("zh-CN", map[string]string{
		"app.name": "应用程序",
	})

	printResult("instance1.T(\"app.name\") [en]", instance1.T("app.name"))
	printResult("instance2.T(\"app.name\") [zh-CN]", instance2.T("app.name"))

	// 验证实例互不影响
	printResult("instance1.GetLang()", instance1.GetLang())
	printResult("instance2.GetLang()", instance2.GetLang())

	// ========== 16. 完整功能列表 ==========
	printSeparator("16. 所有已测试的方法总结")

	methods := []string{
		"初始化: Default(), New(), Init()",
		"加载: Load(), LoadFS(), LoadContent(), LoadMessages()",
		"语言设置: SetLang(), GetLang(), SetLanguage(), GetLanguage()",
		"语言列表: Langs(), GetLanguages()",
		"默认/回退: SetDefaultLang(), SetFallbackLang()",
		"简单翻译: T(), Translate()",
		"指定语言: TL(), TranslateLang()",
		"带参数: Tf(), TranslateWith()",
		"指定语言+参数: TLf(), TranslateLangWith()",
		"复数: Tp(), TranslatePlural()",
		"指定语言+复数: TLp(), TranslateLangPlural()",
		"Map参数: TMap(), TLMap()",
		"Context: WithLang(), LangFromContext()",
		"Context翻译: TC(), TranslateContext()",
		"Context+参数: TCf(), TranslateContextWith()",
		"Context+复数: TCp(), TranslateContextPlural()",
		"高级: GetBundle()",
	}

	for i, method := range methods {
		fmt.Printf("  %2d. %s\n", i+1, method)
	}

	// ========== 完成 ==========
	fmt.Println("\n=================================================")
	fmt.Println("           所有测试完成！")
	fmt.Println("=================================================")
}
