package captcha

import "time"

type Provider string

func (p Provider) String() string {
	return string(p)
}

// noinspection all
const (
	ProviderText  Provider = "text"  // 文字验证码
	ProviderDrag  Provider = "drag"  // 拖拽验证码
	ProviderClick Provider = "click" // 点选验证码
	ProviderSlide Provider = "slide" // 滑动验证码
)

// Config 验证码配置
type Config struct {
	Provider Provider    `json:"provider" yaml:"provider"`
	Text     TextConfig  `json:"text" yaml:"text"`   // 文字验证码
	Drag     DragConfig  `json:"drag" yaml:"drag"`   // 拖拽验证码
	Click    ClickConfig `json:"click" yaml:"click"` // 点选验证码
	Slide    SlideConfig `json:"slide" yaml:"slide"` // 滑动验证码
}

type TextConfig struct {
	Length        int           `json:"length" yaml:"length"`                 // 验证码长度
	ExpireTime    time.Duration `json:"expire_time" yaml:"expire_time"`       // 验证码有效期
	CaseSensitive bool          `json:"case_sensitive" yaml:"case_sensitive"` // 验证码大小写敏感
	Width         int           `json:"width" yaml:"width"`                   // 图片宽度
	Height        int           `json:"height" yaml:"height"`                 // 图片高度
}

// DragConfig 拖拽验证码配置
type DragConfig struct {
	BackgroundWidth  int  `json:"background_width" yaml:"background_width"`   // 背景图宽度
	BackgroundHeight int  `json:"background_height" yaml:"background_height"` // 背景图高度
	TemplateWidth    int  `json:"template_width" yaml:"template_width"`       // 拼图模板宽度
	TemplateHeight   int  `json:"template_height" yaml:"template_height"`     // 拼图模板高度
	Tolerance        int  `json:"tolerance" yaml:"tolerance"`                 // 容错像素范围
	ShowShadow       bool `json:"show_shadow" yaml:"show_shadow"`             // 是否显示阴影
}

// ClickConfig 点选验证码配置
type ClickConfig struct {
	Width          int      `json:"width" yaml:"width"`                     // 图片宽度
	Height         int      `json:"height" yaml:"height"`                   // 图片高度
	WordCount      int      `json:"word_count" yaml:"word_count"`           // 需要点击的文字数量
	TotalWords     int      `json:"total_words" yaml:"total_words"`         // 总文字数量
	FontSize       int      `json:"font_size" yaml:"font_size"`             // 字体大小
	WordList       []string `json:"word_list" yaml:"word_list"`             // 文字列表
	RandomRotation bool     `json:"random_rotation" yaml:"random_rotation"` // 随机旋转文字
}

// SlideConfig 滑动验证码配置
type SlideConfig struct {
	Width           int  `json:"width" yaml:"width"`                       // 背景图宽度
	Height          int  `json:"height" yaml:"height"`                     // 背景图高度
	TemplateWidth   int  `json:"template_width" yaml:"template_width"`     // 滑块宽度
	TemplateHeight  int  `json:"template_height" yaml:"template_height"`   // 滑块高度
	Tolerance       int  `json:"tolerance" yaml:"tolerance"`               // 容错像素范围
	ShowTrajectory  bool `json:"show_trajectory" yaml:"show_trajectory"`   // 是否显示轨迹
	BackgroundNoise bool `json:"background_noise" yaml:"background_noise"` // 背景干扰
}

// noinspection all
const (
	CaptchaTextKey  = "captcha_text"
	CaptchaDragKey  = "captcha_drag"
	CaptchaClickKey = "captcha_click"
	CaptchaSlideKey = "captcha_slide"
)
