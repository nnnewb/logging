package logging

import (
	"net/url"
	"testing"
)

// 测试 lumberjack scheme 的 URI 解析是否得到预期的结果。
// 测试包括对 opaque URI、带 authority 的 URI 和不带 authority 的 URI 。
// 三类 URI 的测试都包括带 query 和 不带 query 的测试用例。
//
// 用例仅测试 url.Parse 取得的 url.URL 结构体字段是否符合预期，并不测试 LumberjackSink 的功能。
// TestLumberjackSinkURI 测试通过URI配置LumberjackSink的正确性。
// 该函数通过解析不同配置的URI来验证LumberjackSink配置的正确性。
func TestLumberjackSinkURI(t *testing.T) {
	// 定义测试用例结构体，包含URI和预期的配置参数。
	var testcases = []struct {
		uri        string
		path       string
		maxBackups string
		maxSize    string
		maxAge     string
		localTime  bool
		compress   bool
	}{
		// 测试用例1: 最小配置，仅指定日志文件路径。
		{
			uri:        "lumberjack:foo.log",
			path:       "foo.log",
			maxBackups: "",
			maxSize:    "",
			maxAge:     "",
			localTime:  false,
			compress:   false,
		},
		// 不主动解码 url 编码的路径
		{
			uri:        "lumberjack:logs%2ffoo.log",
			path:       "logs%2ffoo.log",
			maxBackups: "",
			maxSize:    "",
			maxAge:     "",
			localTime:  false,
			compress:   false,
		},
		{
			uri:        "lumberjack:%2fvar%2flog%2ffoo.log",
			path:       "%2fvar%2flog%2ffoo.log",
			maxBackups: "",
			maxSize:    "",
			maxAge:     "",
			localTime:  false,
			compress:   false,
		},
		// 测试用例2: 完整配置，指定所有可配置参数。
		{
			uri:        "lumberjack:foo.log?max_backups=5&max_size=10&max_age=7&local_time=true&compress=true",
			path:       "foo.log",
			maxBackups: "5",
			maxSize:    "10",
			maxAge:     "7",
			localTime:  true,
			compress:   true,
		},
		// 测试用例3: 使用绝对路径，仅指定日志文件路径。
		{
			uri:        "lumberjack:/tmp/foo.log",
			path:       "/tmp/foo.log",
			maxBackups: "",
			maxSize:    "",
			maxAge:     "",
			localTime:  false,
			compress:   false,
		},
		// 测试用例4: 使用绝对路径和部分配置参数。
		{
			uri:        "lumberjack:/tmp/foo.log?max_backups=3&max_size=5&max_age=2&local_time=false&compress=false",
			path:       "/tmp/foo.log",
			maxBackups: "3",
			maxSize:    "5",
			maxAge:     "2",
			localTime:  false,
			compress:   false,
		},
		// 测试用例5: 包含主机名的URI，仅指定日志文件路径。
		{
			uri:        "lumberjack://localhost/tmp/foo.log",
			path:       "/tmp/foo.log",
			maxBackups: "",
			maxSize:    "",
			maxAge:     "",
			localTime:  false,
			compress:   false,
		},
		// 测试用例6: 包含主机名的URI和完整配置参数。
		{
			uri:        "lumberjack://localhost/tmp/foo.log?max_backups=10&max_size=20&max_age=14&local_time=true&compress=true",
			path:       "/tmp/foo.log",
			maxBackups: "10",
			maxSize:    "20",
			maxAge:     "14",
			localTime:  true,
			compress:   true,
		},
	}

	// 遍历所有测试用例，解析URI并验证解析结果是否与预期一致。
	for _, tc := range testcases {
		// 使用tc.uri作为测试用例的标识符。
		t.Run(tc.uri, func(t *testing.T) {
			// 解析URI。
			u, err := url.Parse(tc.uri)
			if err != nil {
				// 如果解析失败，输出错误信息并跳过后续测试。
				t.Fatalf("failed to parse URI: %v", err)
			}

			// 验证URI的scheme是否为lumberjack。
			if u.Scheme != "lumberjack" {
				t.Errorf("expected scheme 'lumberjack', got '%s'", u.Scheme)
			}

			// 验证URI的路径是否与预期一致。
			if u.Opaque != "" && u.Opaque != tc.path {
				t.Errorf("expected path '%s', got '%s'", tc.path, u.Opaque)
			}
			if u.Path != "" && u.Path != tc.path {
				t.Errorf("expected path '%s', got '%s'", tc.path, u.Path)
			}

			// 解析并验证URI的查询参数是否与预期一致。
			query := u.Query()
			if tc.maxBackups != "" && query.Get("max_backups") != tc.maxBackups {
				t.Errorf("expected max_backups '%s', got '%s'", tc.maxBackups, query.Get("max_backups"))
			}
			if tc.maxSize != "" && query.Get("max_size") != tc.maxSize {
				t.Errorf("expected max_size '%s', got '%s'", tc.maxSize, query.Get("max_size"))
			}
			if tc.maxAge != "" && query.Get("max_age") != tc.maxAge {
				t.Errorf("expected max_age '%s', got '%s'", tc.maxAge, query.Get("max_age"))
			}
			if tc.localTime && query.Get("local_time") != "true" {
				t.Errorf("expected local_time 'true', got '%s'", query.Get("local_time"))
			}
			if tc.compress && query.Get("compress") != "true" {
				t.Errorf("expected compress 'true', got '%s'", query.Get("compress"))
			}
		})
	}
}
