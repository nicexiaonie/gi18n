//go:build gi18n_no_builtin

package gi18n

// 当使用 -tags=gi18n_no_builtin 编译时，不加载任何内置语言包
// 用户需要自行加载语言包

func init() {
	// 不加载任何内置语言包
}

// UseBuiltin 在禁用模式下返回错误提示
func UseBuiltin() error {
	return nil // 静默忽略，因为内置包被禁用
}
