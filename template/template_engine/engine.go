package template_engine

import (
	"bytes"
	"html/template"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

//
// ━━━━━━神兽出没━━━━━━
// 　　 ┏┓     ┏┓
// 　　┏┛┻━━━━━┛┻┓
// 　　┃　　　　　 ┃
// 　　┃　　━　　　┃
// 　　┃　┳┛　┗┳  ┃
// 　　┃　　　　　 ┃
// 　　┃　　┻　　　┃
// 　　┃　　　　　 ┃
// 　　┗━┓　　　┏━┛　Code is far away from bug with the animal protecting
// 　　　 ┃　　　┃    神兽保佑,代码无bug
// 　　　　┃　　　┃
// 　　　　┃　　　┗━━━┓
// 　　　　┃　　　　　　┣┓
// 　　　　┃　　　　　　┏┛
// 　　　　┗┓┓┏━┳┓┏┛
// 　　　　 ┃┫┫ ┃┫┫
// 　　　　 ┗┻┛ ┗┻┛
//
// ━━━━━━感觉萌萌哒━━━━━━
//
//
// User helay
// Date: 2025/6/15 13:16
//

type TplExt string

var tplExts = map[TplExt]bool{
	".html": true,
}

type Engine struct {
	templates *template.Template
	funcMap   template.FuncMap
	devMode   bool
	fsys      fs.FS  // 保存传入的文件系统实例
	fsysPath  string // 生产模式的虚拟路径（如"tpl"）
	localPath string // 开发模式下的本地路径
}

// New 创建模板引擎实例
// fsys: 生产模式使用的 embed.FS
// fsysPath: 生产模式下的虚拟路径
// localPath: 开发模式下的本地模板目录路径
// devMode: 是否为开发模式
func New(fsys fs.FS, fsysPath, localPath string, devMode bool) *Engine {
	e := &Engine{
		funcMap: template.FuncMap{
			// 时间处理
			"now":     time.Now,
			"date":    formatDate,
			"timeAgo": timeSince, // 实现相对时间显示

			// 字符串处理
			"truncate": truncateString,
			"safeHTML": func(s string) template.HTML { return template.HTML(s) },

			// 数学计算
			"add":    func(a, b int) int { return a + b },
			"divide": func(a, b int) float64 { return float64(a) / float64(b) },

			// 链接处理
			"a":          A,
			"aSafe":      ASafe,
			"aWithQuery": AWithQuery,
			"dict":       Dict,
		},
		devMode:   devMode,
		fsys:      fsys,
		fsysPath:  fsysPath,
		localPath: localPath,
	}

	return e
}

// AddFunc 添加自定义函数
func (e *Engine) AddFunc(name string, fn any) {
	// 1. 更新函数映射表
	e.funcMap[name] = fn
}

func (e *Engine) Load() error {
	return e.loadTemplates()
}

func (e *Engine) SetExts(exts ...TplExt) {
	for _, ext := range exts {
		tplExts[ext] = true
	}
}

func (e *Engine) checkExt(name string) bool {
	return tplExts[TplExt(filepath.Ext(name))]
}

func (e *Engine) loadTemplates() error {
	e.templates = template.New("").Funcs(e.funcMap)
	if e.devMode || e.fsys == nil {
		// 开发模式：从本地文件系统递归加载；未启用embed.FS时，使用本地文件系统
		return e.walkDir("", os.DirFS(e.localPath))
	}

	// 生产模式：从 embed.FS 递归加载
	subFs, err := fs.Sub(e.fsys, e.fsysPath)
	if err != nil {
		return err
	}
	return e.walkDir("", subFs)
}

// 统一使用 fs.FS 接口处理，通过 relPath 维护相对路径
func (e *Engine) walkDir(relPath string, fsys fs.FS) error {
	entries, err := fs.ReadDir(fsys, ".")
	if err != nil {
		return err
	}

	for _, entry := range entries {
		entryPath := filepath.Join(relPath, entry.Name())

		if entry.IsDir() {
			subfs, _err := fs.Sub(fsys, entry.Name())
			if _err != nil {
				return _err
			}
			if err = e.walkDir(entryPath, subfs); err != nil {
				return err
			}
			continue
		}

		if !e.checkExt(filepath.Ext(entry.Name())) {
			continue
		}

		content, _err := fs.ReadFile(fsys, entry.Name())
		if _err != nil {
			return _err
		}
		if _, err = e.templates.New(entryPath).Parse(string(content)); err != nil {
			return err
		}
	}
	return nil
}

// RenderString 渲染模板为字符串
func (e *Engine) RenderString(name string, data any) (string, error) {
	var buf bytes.Buffer
	if err := e.Render(&buf, name, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// Render 渲染模板到 io.Writer
func (e *Engine) Render(w io.Writer, name string, data any) error {
	if e.devMode {
		if err := e.loadTemplates(); err != nil { // 使用保存的 e.fsys
			return err
		}
	}
	// 统一使用 / 作为路径分隔符
	name = filepath.ToSlash(name)
	return e.templates.ExecuteTemplate(w, name, data)
}
