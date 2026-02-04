package operators

import (
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"github.com/google/uuid"
	"github.com/helays/utils/v2/net/checkIp"
	"github.com/helays/utils/v2/rule-engine/validator/types"
	"github.com/helays/utils/v2/tools"
	"net/mail"
	"net/url"
	"regexp"
	"strings"
)

func ValidateFormat(operator types.Operator, value any, rule []any) (string, bool) {
	str := tools.Any2string(value)
	switch operator {
	case types.FormatEmail:
		if !ValidateEmail(str) {
			return types.FormatChineseMap[operator], false
		}
	case types.FormatURL:
		if !ValidateURL(str) {
			return types.FormatChineseMap[operator], false
		}
	case types.FormatIP:
		if !checkIp.IsIP(str) {
			return types.FormatChineseMap[operator], false
		}
	case types.FormatPhone:
		if !ValidatePhone(str) {
			return types.FormatChineseMap[operator], false
		}
	case types.FormatIDCard:
		if !FormatIDCard(str) {
			return types.FormatChineseMap[operator], false
		}
	case types.FormatCreditCard:
		if !FormatCreditCard(str) {
			return types.FormatChineseMap[operator], false
		}
	case types.FormatHexColor:
		if !FormatHexColor(str) {
			return types.FormatChineseMap[operator], false
		}
	case types.FormatJSON:
		if !FormatJSON(str) {
			return types.FormatChineseMap[operator], false
		}
	case types.FormatXML:
		if !FormatXML(str) {
			return types.FormatChineseMap[operator], false
		}
	case types.FormatBase64:
		if !FormatBase64(str) {
			return types.FormatChineseMap[operator], false
		}
	case types.FormatUUID:
		if !FormatUUID(str) {
			return types.FormatChineseMap[operator], false
		}
	}
	return "", true
}

// ValidateEmail 验证邮箱格式
func ValidateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// ValidateURL 验证URL格式
func ValidateURL(urlStr string) bool {
	u, err := url.ParseRequestURI(urlStr)
	if err != nil {
		return false
	}
	return u.Scheme != "" && u.Host != ""
}

// ValidatePhone 验证手机号格式(中国)
func ValidatePhone(phone string) bool {
	// 中国手机号正则: 1开头，第二位3-9，后面9位数字
	re := regexp.MustCompile(`^1[3-9]\d{9}$`)
	return re.MatchString(phone)
}

// FormatIDCard 校验身份证号
func FormatIDCard(idCard string) bool {
	// 简单校验：15位或18位，最后一位可以是数字或X
	pattern := `^(\d{15}|\d{17}[\dXx])$`
	matched, err := regexp.MatchString(pattern, idCard)
	return err == nil && matched
}

// FormatCreditCard 校验信用卡号
func FormatCreditCard(cardNumber string) bool {
	// 移除所有非数字字符
	cleaned := strings.ReplaceAll(cardNumber, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")

	// 基本校验：13-19位数字
	if len(cleaned) < 13 || len(cleaned) > 19 {
		return false
	}

	// Luhn算法校验
	sum := 0
	alternate := false
	for i := len(cleaned) - 1; i >= 0; i-- {
		digit := int(cleaned[i] - '0')
		if alternate {
			digit *= 2
			if digit > 9 {
				digit = (digit / 10) + (digit % 10)
			}
		}
		sum += digit
		alternate = !alternate
	}
	return sum%10 == 0
}

// FormatHexColor 校验十六进制颜色值
func FormatHexColor(color string) bool {
	pattern := `^#?([0-9a-fA-F]{3}|[0-9a-fA-F]{6}|[0-9a-fA-F]{8})$`
	matched, err := regexp.MatchString(pattern, color)
	return err == nil && matched
}

// FormatJSON 校验JSON格式
func FormatJSON(jsonStr string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(jsonStr), &js) == nil
}

// FormatXML 校验XML格式
func FormatXML(xmlStr string) bool {
	return xml.Unmarshal([]byte(xmlStr), new(interface{})) == nil
}

// FormatBase64 校验Base64编码
func FormatBase64(base64Str string) bool {
	_, err := base64.StdEncoding.DecodeString(base64Str)
	return err == nil
}

// FormatUUID 校验UUID
func FormatUUID(uuidStr string) bool {
	_, err := uuid.Parse(uuidStr)
	return err == nil
}
