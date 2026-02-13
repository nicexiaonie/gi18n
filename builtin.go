//go:build !gi18n_no_builtin

package gi18n

import (
	"embed"
)

//go:embed builtin/*.json
var builtinFS embed.FS

func init() {
	// 自动加载内置语言包
	if err := Default().LoadFS(builtinFS, "builtin"); err != nil {
		// 内置包加载失败不应该 panic，只是静默失败
		// 用户可以通过 Langs() 检查已加载的语言
	}
}

// UseBuiltin 手动加载内置语言包（如果需要重新加载）
func UseBuiltin() error {
	return Default().LoadFS(builtinFS, "builtin")
}

// BuiltinFS 返回内置语言包的 embed.FS（高级用法）
func BuiltinFS() embed.FS {
	return builtinFS
}
