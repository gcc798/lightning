package middleware

import (
	"net/url"
	"testing"
)

func TestSanitizeJSON(t *testing.T) {
	t.Parallel()

	input := []byte(`{"username":"admin","password":"secret","nested":{"refresh_token":"refresh","safe":1},"items":[{"accessKeySecret":"key"}]}`)
	want := `{"items":[{"accessKeySecret":"***"}],"nested":{"refresh_token":"***","safe":1},"password":"***","username":"admin"}`
	if got := sanitizeJSON(input); got != want {
		t.Fatalf("sanitizeJSON() = %s, want %s", got, want)
	}
}

func TestSanitizeJSONRejectsUnstructuredBody(t *testing.T) {
	t.Parallel()

	if got := sanitizeJSON([]byte("password=secret")); got != "[请求体不是有效 JSON]" {
		t.Fatalf("sanitizeJSON() = %q", got)
	}
}

func TestSanitizeQuery(t *testing.T) {
	t.Parallel()

	values := url.Values{
		"code":  {"123456"},
		"page":  {"1"},
		"token": {"bearer"},
	}
	got := sanitizeQuery(values)
	if got.Get("code") != redactedValue || got.Get("token") != redactedValue {
		t.Fatalf("sensitive query values were not redacted: %v", got)
	}
	if got.Get("page") != "1" {
		t.Fatalf("safe query value changed: %v", got)
	}
}

func TestIsSensitiveField(t *testing.T) {
	t.Parallel()

	for _, field := range []string{"password", "newPassword", "refresh_token", "Access-Key-Secret", "Authorization"} {
		if !isSensitiveField(field) {
			t.Errorf("isSensitiveField(%q) = false", field)
		}
	}
	for _, field := range []string{"username", "pageNum", "configCode"} {
		if isSensitiveField(field) {
			t.Errorf("isSensitiveField(%q) = true", field)
		}
	}
}
