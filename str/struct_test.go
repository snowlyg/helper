package str

import "testing"

func TestJoin(t *testing.T) {
	t.Run("字符串拼接", func(t *testing.T) {
		if got := Join("abc", " ", "def"); got != "abc def" {
			t.Errorf("Join() = %v, want %v", got, "abc def")
		}
	})
	t.Run("字符串拼接单个字符", func(t *testing.T) {
		if got := Join("abc"); got != "abc" {
			t.Errorf("Join() = %v, want %v", got, "abc")
		}
	})
	t.Run("中文字符串拼接", func(t *testing.T) {
		if got := Join("中文字符串拼接", " ", "你好"); got != "中文字符串拼接 你好" {
			t.Errorf("Join() = %v, want %v", got, "中文字符串拼接 你好")
		}
	})
	t.Run("中文字符串拼接单个字符", func(t *testing.T) {
		if got := Join("中文字符串拼接"); got != "中文字符串拼接" {
			t.Errorf("Join() = %v, want %v", got, "中文字符串拼接")
		}
	})
}
