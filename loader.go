package gi18n

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"gopkg.in/yaml.v3"
)

// registerUnmarshalers 注册各格式的解析器
func (b *Bundle) registerUnmarshalers() {
	b.bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	b.bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)
	b.bundle.RegisterUnmarshalFunc("yml", yaml.Unmarshal)
	b.bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
}

// Load 从目录加载语言文件
// 支持 .json, .yaml, .yml, .toml 格式
// 文件名格式: {语言标记}.{扩展名}，如 en.json, zh-CN.yaml
func (b *Bundle) Load(dir string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("gi18n: failed to read directory %s: %w", dir, err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if err := b.loadFile(dir, entry.Name()); err != nil {
			return err
		}
	}

	b.clearLocalizerCache()
	return nil
}

// loadFile 加载单个文件
func (b *Bundle) loadFile(dir, filename string) error {
	ext := strings.ToLower(filepath.Ext(filename))
	if !isSupportedExt(ext) {
		return nil
	}

	filePath := filepath.Join(dir, filename)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("gi18n: failed to read file %s: %w", filePath, err)
	}

	lang := extractLangFromFilename(filename)
	return b.loadData(lang, ext, data)
}

// LoadFS 从 embed.FS 加载语言文件
func (b *Bundle) LoadFS(fsys embed.FS, root string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	err := fs.WalkDir(fsys, root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(d.Name()))
		if !isSupportedExt(ext) {
			return nil
		}

		data, err := fsys.ReadFile(path)
		if err != nil {
			return fmt.Errorf("gi18n: failed to read file %s: %w", path, err)
		}

		lang := extractLangFromFilename(d.Name())
		return b.loadData(lang, ext, data)
	})

	if err != nil {
		return err
	}

	b.clearLocalizerCache()
	return nil
}

// LoadContent 从字节内容加载语言包
// lang: 语言标记，如 "en", "zh-CN"
// format: 格式，如 "json", "yaml", "toml"
// data: 文件内容
func (b *Bundle) LoadContent(lang, format string, data []byte) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	ext := "." + strings.TrimPrefix(format, ".")
	if err := b.loadData(lang, ext, data); err != nil {
		return err
	}

	b.clearLocalizerCache()
	return nil
}

// LoadMessages 直接加载消息映射（简化格式）
func (b *Bundle) LoadMessages(lang string, messages map[string]string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	tag := parseLanguageTag(lang)
	for id, text := range messages {
		if err := b.bundle.AddMessages(tag, &i18n.Message{
			ID:    id,
			Other: text,
		}); err != nil {
			return fmt.Errorf("gi18n: failed to add message %s: %w", id, err)
		}
	}

	b.addSupported(lang)
	b.clearLocalizerCache()
	return nil
}

// loadData 加载数据到 bundle
func (b *Bundle) loadData(lang, ext string, data []byte) error {
	// 先尝试解析为通用格式，处理嵌套和简化写法
	processed, err := b.preprocessData(data, ext)
	if err != nil {
		return err
	}

	filename := fmt.Sprintf("%s%s", normalizeLanguageTag(lang), ext)
	if _, err := b.bundle.ParseMessageFileBytes(processed, filename); err != nil {
		return fmt.Errorf("gi18n: failed to parse message file %s: %w", filename, err)
	}

	b.addSupported(lang)
	return nil
}

// preprocessData 预处理数据，处理嵌套和简化写法
func (b *Bundle) preprocessData(data []byte, ext string) ([]byte, error) {
	// 解析为通用 map
	var raw map[string]interface{}
	var err error

	switch ext {
	case ".json":
		err = json.Unmarshal(data, &raw)
	case ".yaml", ".yml":
		err = yaml.Unmarshal(data, &raw)
	case ".toml":
		err = toml.Unmarshal(data, &raw)
	default:
		return data, nil
	}

	if err != nil {
		// 解析失败，返回原始数据让 go-i18n 处理
		return data, nil
	}

	// 展平嵌套结构并转换简化写法
	flattened := b.flattenMessages("", raw)

	// 转换为 go-i18n 格式
	result := make(map[string]interface{})
	for id, value := range flattened {
		result[id] = value
	}

	// 重新编码为 JSON（go-i18n 内部统一处理）
	return json.Marshal(result)
}

// flattenMessages 展平嵌套结构
func (b *Bundle) flattenMessages(prefix string, data map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for key, value := range data {
		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}

		switch v := value.(type) {
		case string:
			// 简化写法: "hello": "你好" -> {"id": "hello", "other": "你好"}
			result[fullKey] = map[string]interface{}{
				"id":    fullKey,
				"other": v,
			}
		case map[string]interface{}:
			if isMessageObject(v) {
				// go-i18n 消息对象，添加 id
				v["id"] = fullKey
				result[fullKey] = v
			} else {
				// 嵌套命名空间，继续展平
				nested := b.flattenMessages(fullKey, v)
				for k, nv := range nested {
					result[k] = nv
				}
			}
		default:
			// 其他类型，尝试转为字符串
			result[fullKey] = map[string]interface{}{
				"id":    fullKey,
				"other": fmt.Sprintf("%v", v),
			}
		}
	}

	return result
}

// messageObjectKeys go-i18n 消息对象的特征字段（用于快速查找）
var messageObjectKeys = map[string]struct{}{
	"one": {}, "other": {}, "zero": {},
	"two": {}, "few": {}, "many": {},
	"description": {}, "hash": {},
}

// isMessageObject 判断是否为 go-i18n 消息对象
func isMessageObject(obj map[string]interface{}) bool {
	for key := range obj {
		if _, ok := messageObjectKeys[key]; ok {
			return true
		}
	}
	return false
}

// extractLangFromFilename 从文件名提取语言标记
func extractLangFromFilename(filename string) string {
	name := strings.TrimSuffix(filename, filepath.Ext(filename))
	return normalizeLanguageTag(name)
}

// isSupportedExt 判断是否支持的扩展名
func isSupportedExt(ext string) bool {
	switch ext {
	case ".json", ".yaml", ".yml", ".toml":
		return true
	}
	return false
}

// ========== 全局函数 ==========

// Load 从目录加载语言文件（全局）
func Load(dir string) error {
	return Default().Load(dir)
}

// LoadFS 从 embed.FS 加载语言文件（全局）
func LoadFS(fsys embed.FS, root string) error {
	return Default().LoadFS(fsys, root)
}

// LoadContent 从内容加载语言包（全局）
func LoadContent(lang, format string, data []byte) error {
	return Default().LoadContent(lang, format, data)
}

// LoadMessages 加载消息映射（全局）
func LoadMessages(lang string, messages map[string]string) error {
	return Default().LoadMessages(lang, messages)
}
