package security_policy

import (
	"errors"
	"fmt"
	"strings"
)

// PasswordPolicy 密码策略
type PasswordPolicy struct {
	PasswordMinLength      int  `json:"password_min_length"`      // 密码最小长度
	PasswordRequireUpper   bool `json:"password_require_upper"`   // 需要大写字母
	PasswordRequireLower   bool `json:"password_require_lower"`   //需要小写字母
	PasswordRequireDigit   bool `json:"password_require_digit"`   //需要数字
	PasswordRequireSpecial bool `json:"password_require_special"` // 需要特殊字符
	PasswordExpireDays     int  `json:"password_expire_days"`     // 密码过期天数
	PasswordHistorySize    int  `json:"password_history_size"`    // 密码历史记录数量
	PasswordReuseLimit     int  `json:"password_reuse_limit"`     // 密码重复使用限制
	AllowWeakPassword      bool `json:"allow_weak_password"`      // 允许弱密码
	// 新增：防止连续键盘按键配置
	PreventSequentialChars bool `json:"prevent_sequential_chars"` // 防止连续字符(如: abc, 123)
	PreventKeyboardPattern bool `json:"prevent_keyboard_pattern"` // 防止键盘模式(如: qwert, asdfg)
	MaxSequentialLen       int  `json:"max_sequential_len"`       // 最大连续字符长度
	MaxRepeatedChars       int  `json:"max_repeated_chars"`       // 最大重复字符数，比如连续aaaa，ccccc，ddddd
}

// ValidatePassword 验证密码是否符合策略
func (p *PasswordPolicy) ValidatePassword(password string) error {
	if len(password) < p.PasswordMinLength {
		return fmt.Errorf("密码长度不能少于%d位", p.PasswordMinLength)
	}

	// 检查大写字母
	if p.PasswordRequireUpper && !containsUpper(password) {
		return errors.New("密码必须包含大写字母")
	}

	// 检查小写字母
	if p.PasswordRequireLower && !containsLower(password) {
		return errors.New("密码必须包含小写字母")
	}

	// 检查数字
	if p.PasswordRequireDigit && !containsDigit(password) {
		return errors.New("密码必须包含数字")
	}

	// 检查特殊字符
	if p.PasswordRequireSpecial && !containsSpecial(password) {
		return errors.New("密码必须包含特殊字符")
	}

	// 检查连续字符
	if p.PreventSequentialChars && p.hasSequentialChars(password) {
		return fmt.Errorf("密码包含%d个以上连续字符", p.MaxSequentialLen)
	}

	// 检查键盘模式
	if p.PreventKeyboardPattern && p.hasKeyboardPattern(password) {
		return errors.New("密码包含键盘模式")
	}

	// 检查重复字符
	if p.hasRepeatedChars(password) {
		return fmt.Errorf("密码包含%d个以上重复字符", p.MaxRepeatedChars)
	}

	// 检查弱密码（如果配置了不允许弱密码）
	if !p.AllowWeakPassword && p.isWeakPassword(password) {
		return errors.New("密码强度太弱，请使用更复杂的密码")
	}

	return nil
}

// 辅助验证函数
func containsUpper(s string) bool {
	for _, r := range s {
		if r >= 'A' && r <= 'Z' {
			return true
		}
	}
	return false
}

func containsLower(s string) bool {
	for _, r := range s {
		if r >= 'a' && r <= 'z' {
			return true
		}
	}
	return false
}

func containsDigit(s string) bool {
	for _, r := range s {
		if r >= '0' && r <= '9' {
			return true
		}
	}
	return false
}

func containsSpecial(s string) bool {
	specialChars := "!@#$%^&*()_+-=[]{}|;:,.<>?/~`"
	for _, r := range s {
		if strings.ContainsRune(specialChars, r) {
			return true
		}
	}
	return false
}

// hasSequentialChars 检查连续字符 (abc, 123, 987等)
func (p *PasswordPolicy) hasSequentialChars(password string) bool {
	if p.MaxSequentialLen < 2 {
		return false
	}

	runes := []rune(password)
	for i := 0; i <= len(runes)-p.MaxSequentialLen; i++ {
		// 检查正向连续
		isSequential := true
		for j := 1; j < p.MaxSequentialLen; j++ {
			if runes[i+j] != runes[i+j-1]+1 {
				isSequential = false
				break
			}
		}
		if isSequential {
			return true
		}

		// 检查反向连续
		isReverseSequential := true
		for j := 1; j < p.MaxSequentialLen; j++ {
			if runes[i+j] != runes[i+j-1]-1 {
				isReverseSequential = false
				break
			}
		}
		if isReverseSequential {
			return true
		}
	}
	return false
}

// hasKeyboardPattern 检查键盘模式
func (p *PasswordPolicy) hasKeyboardPattern(password string) bool {
	if p.MaxSequentialLen < 2 {
		return false
	}

	// 常见的键盘行
	keyboardRows := []string{
		"qwertyuiop",
		"asdfghjkl",
		"zxcvbnm",
		"1234567890",
	}

	runes := []rune(strings.ToLower(password))
	for i := 0; i <= len(runes)-p.MaxSequentialLen; i++ {
		segment := string(runes[i : i+p.MaxSequentialLen])

		for _, row := range keyboardRows {
			// 检查正向
			if strings.Contains(row, segment) {
				return true
			}
			// 检查反向
			reversedSegment := reverseString(segment)
			if strings.Contains(row, reversedSegment) {
				return true
			}
		}
	}
	return false
}

// hasRepeatedChars 检查重复字符
func (p *PasswordPolicy) hasRepeatedChars(password string) bool {
	if p.MaxRepeatedChars < 1 {
		return false
	}

	runes := []rune(password)
	if len(runes) == 0 {
		return false
	}

	currentChar := runes[0]
	count := 1

	for i := 1; i < len(runes); i++ {
		if runes[i] == currentChar {
			count++
			if count > p.MaxRepeatedChars {
				return true
			}
		} else {
			currentChar = runes[i]
			count = 1
		}
	}
	return false
}

// isWeakPassword 检查弱密码
func (p *PasswordPolicy) isWeakPassword(password string) bool {
	weakPasswords := []string{
		"password", "123456", "qwerty", "admin", "welcome",
		"abc123", "password1", "12345678", "123456789",
		"111111", "123123", "admin123", "000000",
	}

	lowerPassword := strings.ToLower(password)
	for _, weak := range weakPasswords {
		if lowerPassword == weak {
			return true
		}
	}

	// 检查是否全是相同字符
	if isAllSameChars(password) {
		return true
	}

	return false
}

// 辅助函数
func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func isAllSameChars(s string) bool {
	if len(s) == 0 {
		return false
	}
	first := rune(s[0])
	for _, char := range s {
		if char != first {
			return false
		}
	}
	return true
}

// GetPasswordStrength 获取密码强度评分 (0-100)
func (p *PasswordPolicy) GetPasswordStrength(password string) int {
	score := 0

	// 长度评分
	length := len(password)
	if length >= 8 {
		score += 25
	} else if length >= 6 {
		score += 15
	} else {
		score += 5
	}

	// 字符类型评分
	typeCount := 0
	if containsUpper(password) {
		typeCount++
	}
	if containsLower(password) {
		typeCount++
	}
	if containsDigit(password) {
		typeCount++
	}
	if containsSpecial(password) {
		typeCount++
	}

	switch typeCount {
	case 4:
		score += 40
	case 3:
		score += 30
	case 2:
		score += 20
	case 1:
		score += 10
	}

	// 模式惩罚
	if p.hasSequentialChars(password) {
		score -= 15
	}
	if p.hasKeyboardPattern(password) {
		score -= 15
	}
	if p.hasRepeatedChars(password) {
		score -= 10
	}
	if p.isWeakPassword(password) {
		score -= 30
	}

	// 确保分数在 0-100 范围内
	if score < 0 {
		return 0
	}
	if score > 100 {
		return 100
	}
	return score
}
