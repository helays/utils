package template_engine

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
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
	mu             sync.RWMutex
	templateLayout *template.Template            // 基础模板
	templateSets   map[string]*template.Template // 按模板名存储独立模板集
	funcMap        template.FuncMap
	devMode        bool
	fsys           fs.FS  // 保存传入的文件系统实例
	fsysPath       string // 生产模式的虚拟路径（如"tpl"）
	localPath      string // 开发模式下的本地路径
	layoutDir      string
}

// New 创建模板引擎实例
// fsys: 生产模式使用的 embed.FS
// fsysPath: 生产模式下的虚拟路径
// localPath: 开发模式下的本地模板目录路径
// devMode: 是否为开发模式
func New(fsys fs.FS, fsysPath, localPath string, devMode bool) *Engine {
	e := &Engine{
		funcMap:   builtinFuncMap(),
		devMode:   devMode,
		fsys:      fsys,
		fsysPath:  fsysPath,
		localPath: localPath,
	}

	return e
}

// SetLayoutDir 设置layout目录名称（如"layouts"）
func (e *Engine) SetLayoutDir(dir string) {
	e.layoutDir = filepath.ToSlash(dir)
}

// AddFunc 添加自定义函数
func (e *Engine) AddFunc(name string, fn any) {
	// 1. 更新函数映射表
	e.funcMap[name] = fn
}

func (e *Engine) Load() error {
	// 这里需要上锁，开发模式下可能会存在多个线程同时访问
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.layoutDir != "" {
		if err := e.loadLayout(); err != nil {
			return err
		}
	}

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

// 加载layout模板
func (e *Engine) loadLayout() error {
	e.templateLayout = template.New("").Funcs(e.funcMap)
	var (
		subFs fs.FS
		err   error
	)

	if e.devMode || e.fsys == nil {
		subFs = os.DirFS(e.localPath)
	} else {
		// 生产模式：从 embed.FS 递归加载
		subFs, err = fs.Sub(e.fsys, e.fsysPath)
		if err != nil {
			return err
		}
	}
	// 先打开layout目录
	subFs, err = fs.Sub(subFs, e.layoutDir)
	if err != nil {
		return err
	}
	return e.walkDir(e.layoutDir, subFs, true, 1)
}

func (e *Engine) loadTemplates() error {
	e.templateSets = make(map[string]*template.Template)
	var (
		subFs fs.FS
		err   error
	)
	if e.devMode || e.fsys == nil {
		// 开发模式：从本地文件系统递归加载；未启用embed.FS时，使用本地文件系统
		subFs = os.DirFS(e.localPath)
	} else {
		// 生产模式：从 embed.FS 递归加载
		subFs, err = fs.Sub(e.fsys, e.fsysPath)
		if err != nil {
			return err
		}
	}
	return e.walkDir("", subFs, false, 1)
}

// 统一使用 fs.FS 接口处理，通过 relPath 维护相对路径
func (e *Engine) walkDir(relPath string, fsys fs.FS, isLayout bool, level int) error {
	entries, err := fs.ReadDir(fsys, ".")
	if err != nil {
		return err
	}

	for _, entry := range entries {
		fileName := entry.Name()
		entryPath := filepath.Join(relPath, fileName)
		if !isLayout && level == 1 && fileName == e.layoutDir {
			continue
		}
		if entry.IsDir() {
			subfs, _err := fs.Sub(fsys, entry.Name())
			if _err != nil {
				return _err
			}
			if err = e.walkDir(entryPath, subfs, isLayout, level+1); err != nil {
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
		if isLayout {
			if _, err = e.templateLayout.New(entryPath).Parse(string(content)); err != nil {
				return err
			}
		} else {
			// 创建全新的模板实例并设置名称
			var tpl *template.Template
			if e.templateLayout != nil {
				if tpl, err = e.clone(entryPath); err != nil {
					return err
				}
				//// 或者使用tree
				//if err = e.tree(tpl); err != nil {
				//	return err
				//}
			} else {
				tpl = template.New(entryPath).Funcs(e.funcMap)
			}
			if _, err = tpl.Parse(string(content)); err != nil {
				return err
			}
			e.templateSets[entryPath] = tpl
		}

	}
	return nil
}

// 更简单,性能稍微有影响
func (e *Engine) clone(entryPath string) (*template.Template, error) {
	tpl, err := e.templateLayout.Clone()
	if err != nil {
		return nil, err
	}
	// 设置当前模板名称
	return tpl.New(entryPath), nil
}

// 实现较复杂，性能更佳
func (e *Engine) tree(tpl *template.Template) error {
	for _, ltpl := range e.templateLayout.Templates() {
		if _, err := tpl.AddParseTree(ltpl.Name(), ltpl.Tree); err != nil {
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
		if err := e.Load(); err != nil { // 使用保存的 e.fsys
			return err
		}
	}
	// 统一使用 / 作为路径分隔符
	name = filepath.ToSlash(name)
	tpl, ok := e.templateSets[name]
	if !ok {
		return fmt.Errorf("template %s not found", name)
	}
	return tpl.ExecuteTemplate(w, name, data)
}
