package template_engine

import (
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
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
// localPath: 开发模式下的本地模板目录路径
// devMode: 是否为开发模式
func New(fsys fs.FS, fsysPath, localPath string, devMode bool) (*Engine, error) {
	e := &Engine{
		funcMap: template.FuncMap{
			"safeHTML": func(s string) template.HTML {
				return template.HTML(s)
			},
		},
		devMode:   devMode,
		fsys:      fsys,
		fsysPath:  fsysPath,
		localPath: localPath,
	}

	if err := e.loadTemplates(); err != nil {
		return nil, err
	}

	return e, nil
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
			subfs, err := fs.Sub(fsys, entry.Name())
			if err != nil {
				return err
			}
			if err := e.walkDir(entryPath, subfs); err != nil {
				return err
			}
			continue
		}

		if filepath.Ext(entry.Name()) != ".html" {
			continue
		}

		content, err := fs.ReadFile(fsys, entry.Name())
		if err != nil {
			return err
		}

		name := strings.TrimSuffix(entryPath, ".html")
		fmt.Println(name, "---")
		if _, err := e.templates.New(name).Parse(string(content)); err != nil {
			return err
		}
	}
	return nil
}

func (e *Engine) Render(w io.Writer, name string, data interface{}) error {
	if e.devMode {
		if err := e.loadTemplates(); err != nil { // 使用保存的 e.fsys
			return err
		}
	}
	// 统一使用 / 作为路径分隔符
	name = filepath.ToSlash(name)
	return e.templates.ExecuteTemplate(w, name, data)
}
