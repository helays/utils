package route

import (
	"helay.net/go/utils/v3/tools"
)

type Route struct {
	opt   *Config
	embed []*EmbedInfo
}

func New(opt *Config) *Route {
	if opt == nil {
		panic("route.New 参数不能为空")
	}
	r := &Route{
		opt:   opt,
		embed: make([]*EmbedInfo, 0),
	}
	r.opt.Root = tools.Fileabs(r.opt.Root)
	if r.opt.Index == "" {
		r.opt.Index = "index.html"
	}

	return r
}

// AddEmbed 添加静态文件系统
func (ro *Route) AddEmbed(e *EmbedInfo) {
	// 根据Search 字段，判断是否已经添加过
	if !tools.ContainsFunc(ro.embed, e, func(a *EmbedInfo, b *EmbedInfo) bool { return a.Search == b.Search }) {
		ro.embed = append(ro.embed, e)
	}
}
