package openapi

import (
	"encoding/json"
	"os"
	"testing"
)

func TestCriticalOperations(t *testing.T) {
	data, err := os.ReadFile("swagger.json")
	if err != nil {
		t.Fatal(err)
	}
	var doc struct {
		Paths map[string]map[string]json.RawMessage `json:"paths"`
	}
	if err := json.Unmarshal(data, &doc); err != nil {
		t.Fatal(err)
	}
	for path, method := range map[string]string{"/login": "post", "/api/v1/user": "post", "/api/v1/config": "post", "/api/v1/attachment/page": "post"} {
		if _, ok := doc.Paths[path][method]; !ok {
			t.Errorf("missing %s %s", method, path)
		}
	}
}
