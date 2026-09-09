package proto

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"
)

var (
	serviceDeclaration = regexp.MustCompile(`^\s*service\s+[A-Za-z_]\w*\s*\{`)
	rpcDeclaration     = regexp.MustCompile(`^\s*rpc\s+[A-Za-z_]\w*\s*\(`)
	messageDeclaration = regexp.MustCompile(`^\s*message\s+[A-Za-z_]\w*\s*(?:\{|$)`)
	fieldDeclaration   = regexp.MustCompile(`^\s*(?:(?:repeated|optional)\s+)?(?:map<[^>]+>|[.A-Za-z_]\w*(?:\.[A-Za-z_]\w*)*)\s+[A-Za-z_]\w*\s*=\s*\d+\b`)
)

func TestProtoDeclarationsHaveComments(t *testing.T) {
	_, sourceFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("locate test source")
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(sourceFile), "..", ".."))
	files, err := filepath.Glob(filepath.Join(root, "api", "*", "v1", "*.proto"))
	if err != nil {
		t.Fatalf("find proto files: %v", err)
	}
	if len(files) == 0 {
		t.Fatal("no native proto files found")
	}

	for _, file := range files {
		checkProtoComments(t, file)
	}
}

func checkProtoComments(t *testing.T, file string) {
	t.Helper()
	input, err := os.Open(file)
	if err != nil {
		t.Fatalf("open %s: %v", file, err)
	}
	defer input.Close()

	scanner := bufio.NewScanner(input)
	lineNumber := 0
	previousNonEmpty := ""
	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		if serviceDeclaration.MatchString(line) || rpcDeclaration.MatchString(line) || messageDeclaration.MatchString(line) || fieldDeclaration.MatchString(line) {
			if !strings.HasPrefix(strings.TrimSpace(previousNonEmpty), "//") {
				t.Errorf("%s:%d declaration must have a preceding // comment", file, lineNumber)
			}
		}
		if trimmed := strings.TrimSpace(line); trimmed != "" {
			previousNonEmpty = trimmed
		}
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("read %s: %v", file, err)
	}
}
