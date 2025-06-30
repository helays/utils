package encodinghelper

import (
	"encoding/binary"
	"fmt"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"strings"
	"unicode/utf16"
)

func ToUTF8(src []byte, srcCode string) ([]byte, error) {
	var decoder *encoding.Decoder

	switch strings.ToUpper(srcCode) {
	case "GBK", "GB2312":
		decoder = simplifiedchinese.GBK.NewDecoder()
	case "GB18030":
		decoder = simplifiedchinese.GB18030.NewDecoder()
	case "BIG5":
		decoder = traditionalchinese.Big5.NewDecoder()
	case "SHIFTJIS":
		decoder = japanese.ShiftJIS.NewDecoder()
	case "EUCJP":
		decoder = japanese.EUCJP.NewDecoder()
	case "EUCKR":
		decoder = korean.EUCKR.NewDecoder()
	case "UTF-16", "UTF-16LE", "UTF-16BE":
		// 处理 UTF-16 需要额外参数
		return convertUTF16(src, srcCode)
	default: // 包括 UTF-8
		return src, nil // 已经是 UTF-8 不需要转换
	}

	result, err := decoder.Bytes(src)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// 单独处理 UTF-16 的情况
func convertUTF16(src []byte, enc string) ([]byte, error) {
	var bo binary.ByteOrder = binary.LittleEndian
	if strings.HasSuffix(enc, "BE") {
		bo = binary.BigEndian
	}

	// 检查 BOM (Byte Order Mark)
	if len(src) >= 2 {
		switch {
		case (src)[0] == 0xFE && (src)[1] == 0xFF:
			bo = binary.BigEndian
			src = (src)[2:]
		case (src)[0] == 0xFF && (src)[1] == 0xFE:
			bo = binary.LittleEndian
			src = (src)[2:]
		}
	}

	// 确保字节数是偶数
	if len(src)%2 != 0 {
		return nil, fmt.Errorf("UTF-16 data has odd length")
	}

	// 将字节转换为 uint16
	u16 := make([]uint16, 0, len(src)/2)
	for i := 0; i < len(src); i += 2 {
		u16 = append(u16, bo.Uint16((src)[i:i+2]))
	}

	// 将 UTF-16 转换为 UTF-8
	return []byte(string(utf16.Decode(u16))), nil
}
