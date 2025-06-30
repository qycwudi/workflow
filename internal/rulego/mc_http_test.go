package rulego

import (
	"encoding/base64"
	"os"
	"strings"
	"testing"
)

func TestHttpCallNode_prepareFormDataRequestBody(t *testing.T) {
	path := "_base64_content_filename"

	// 原始方法（已优化）
	t.Run("Original Method", func(t *testing.T) {
		if !strings.HasPrefix(path, "_base64_") {
			t.Errorf("path should have prefix _base64_")
		}
		// 读取 content ,读取 filename
		content := strings.Split(strings.TrimPrefix(path, "_base64_"), "_")[0]
		filename := strings.TrimPrefix(path, "_base64_"+content+"_")
		t.Logf("content: %s", content)
		t.Logf("filename: %s", filename)
	})

	// 测试所有优化方案
	testCases := []struct {
		name         string
		input        string
		wantContent  string
		wantFilename string
		wantErr      bool
	}{
		{"Valid Path", "_base64_content_filename", "content", "filename", false},
		{"Valid Path with Multiple Underscores", "_base64_my_content_my_filename_ext", "my", "content_my_filename_ext", false},
		{"Invalid Prefix", "base64_content_filename", "", "", true},
		{"Missing Content", "_base64__filename", "", "", true},
		{"Missing Filename", "_base64_content_", "", "", true},
		{"No Underscore After Content", "_base64_contentfilename", "", "", true},
		{"Empty String", "", "", "", true},
	}

	methods := []struct {
		name string
		fn   func(string) (string, string, error)
	}{
		{"Index", parseBase64PathWithIndex},
	}

	for _, method := range methods {
		t.Run(method.name, func(t *testing.T) {
			for _, tc := range testCases {
				t.Run(tc.name, func(t *testing.T) {
					content, filename, err := method.fn(tc.input)

					if tc.wantErr {
						if err == nil {
							t.Errorf("expected error but got none")
						}
						return
					}

					if err != nil {
						t.Errorf("unexpected error: %v", err)
						return
					}

					if content != tc.wantContent {
						t.Errorf("content = %q, want %q", content, tc.wantContent)
					}

					if filename != tc.wantFilename {
						t.Errorf("filename = %q, want %q", filename, tc.wantFilename)
					}

					t.Logf("content: %s, filename: %s", content, filename)
				})
			}
		})
	}
}

// 性能基准测试
func BenchmarkParseBase64Path(b *testing.B) {
	path := "_base64_content_filename"

	b.Run("Original", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			content := strings.Split(strings.TrimPrefix(path, "_base64_"), "_")[0]
			filename := strings.TrimPrefix(path, "_base64_"+content+"_")
			_, _ = content, filename
		}
	})

	b.Run("Index", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			parseBase64PathWithIndex(path)
		}
	})

}

func TestToBase64(t *testing.T) {
	// 读取文件，把文件转成 base64字符串
	content, err := os.ReadFile("./chain/base.json")
	if err != nil {
		t.Errorf("read file error: %v", err)
	}
	t.Logf("content: %s", base64.StdEncoding.EncodeToString(content))
}
