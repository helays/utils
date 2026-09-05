package dataType

import (
	"encoding/json"
	"testing"
)

func TestBool_UnmarshalJSON_VariousTypes(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want bool
	}{
		{"true", `true`, true},
		{"false", `false`, false},
		{"digit one", `1`, true},
		{"digit zero", `0`, false},
		{"digit non-zero", `2`, true},
		{"string true", `"true"`, true},
		{"string false", `"false"`, false},
		{"string one", `"1"`, true},
		{"string zero", `"0"`, false},
		{"string empty", `""`, false},
		{"null", `null`, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var b Bool
			if err := json.Unmarshal([]byte(c.in), &b); err != nil {
				t.Fatalf("json.Unmarshal(%q) error: %v", c.in, err)
			}
			if b.Bool() != c.want {
				t.Fatalf("Bool() = %v, want %v (input %q)", b.Bool(), c.want, c.in)
			}
		})
	}
}

func TestBool_UnmarshalJSON_DefaultField(t *testing.T) {
	// 嵌入到结构体，验证常规字段解码路径
	type holder struct {
		Flag Bool `json:"flag"`
	}
	var h holder
	raw := `{"flag":"1"}`
	if err := json.Unmarshal([]byte(raw), &h); err != nil {
		t.Fatalf("json.Unmarshal struct error: %v", err)
	}
	if !h.Flag.Bool() {
		t.Fatalf("Flag should be true, got %v", h.Flag.Bool())
	}
}

func TestBool_MarshalJSON_RoundTrip(t *testing.T) {
	for _, want := range []bool{true, false} {
		b := NewBool(want)
		raw, err := json.Marshal(b)
		if err != nil {
			t.Fatalf("Marshal error: %v", err)
		}
		var back Bool
		if err := json.Unmarshal(raw, &back); err != nil {
			t.Fatalf("Unmarshal(%s) error: %v", raw, err)
		}
		if back.Bool() != want {
			t.Fatalf("round trip: got %v, want %v (raw %s)", back.Bool(), want, raw)
		}
	}
}

func TestBool_UnmarshalJSON_Invalid(t *testing.T) {
	cases := []string{
		`[]`,         // 数组
		`{}`,         // 对象
		`"yes"`,      // ParseBool 无法识别
		`"abc"`,      // ParseBool 无法识别
		`[`,          // 非法 JSON
		`"unterminated`, // 非法 JSON
	}
	for _, in := range cases {
		var b Bool
		if err := json.Unmarshal([]byte(in), &b); err == nil {
			t.Errorf("json.Unmarshal(%q) expected error, got nil", in)
		}
	}
}

func TestBool_UnmarshalJSON_Empty(t *testing.T) {
	var b Bool
	if err := json.Unmarshal(nil, &b); err == nil && b.Bool() {
		t.Fatalf("empty input should give false without error, got true")
	}
}