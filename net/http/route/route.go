package route

import (
	"github.com/helays/utils/v2/tools"
)

type Route struct {
	opt *Config
}

func New(opt *Config) *Route {
	r := &Route{
		opt: opt,
	}
	r.opt.Root = tools.Fileabs(r.opt.Root)
	if len(r.opt.Index) == 0 {
		r.opt.Index = []string{"index.html", "index.htm"}
	}
	return r
}
