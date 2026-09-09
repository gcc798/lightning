package iam

import (
	"strings"
	"testing"
)

func TestWechatUsername(t *testing.T) {
	t.Parallel()
	first := wechatUsername("openid-a")
	second := wechatUsername("openid-b")
	if first == second {
		t.Fatal("different OpenIDs produced the same username")
	}
	if !strings.HasPrefix(first, "wx_") || len(first) != 35 {
		t.Fatalf("wechatUsername() = %q", first)
	}
	if got := wechatUsername("openid-a"); got != first {
		t.Fatalf("wechatUsername() is not deterministic: %q != %q", got, first)
	}
}

func TestParseUserAgent(t *testing.T) {
	t.Parallel()
	browser, operatingSystem := parseUserAgent("Mozilla/5.0 (Mac OS X) AppleWebKit Chrome/125 Safari/537.36")
	if browser != "Chrome" || operatingSystem != "macOS" {
		t.Fatalf("parseUserAgent() = %q, %q", browser, operatingSystem)
	}
}
