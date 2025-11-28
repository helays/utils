package captcha

import (
	"time"

	"github.com/helays/utils/v2/tools"
)

type Captcha struct {
	opt *Config
}

func New(opt *Config) *Captcha {
	c := &Captcha{opt: opt}
	c.setDefault()
	return c
}

func (c *Captcha) setDefault() {
	c.opt.Text.Length = tools.Ternary(c.opt.Text.Length < 1, 4, c.opt.Text.Length)
	c.opt.Text.ExpireTime = tools.AutoTimeDuration(c.opt.Text.ExpireTime, time.Second, 4*time.Minute)
	c.opt.Text.Width = tools.Ternary(c.opt.Text.Width < 1, 106, c.opt.Text.Width)
	c.opt.Text.Height = tools.Ternary(c.opt.Text.Height < 1, 40, c.opt.Text.Height)
}
