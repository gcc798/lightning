package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"gopkg.in/yaml.v3"
)

const openAPISpecPath = "api/system/v1.openapi.yaml"

func registerSwaggerEndpoints(srv *khttp.Server) {
	if srv == nil {
		return
	}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL == nil:
			http.NotFound(w, r)
		case r.URL.Path == "/swagger" || r.URL.Path == "/swagger/" || strings.HasSuffix(r.URL.Path, "/index.html"):
			serveSwaggerIndex(w)
		case strings.HasSuffix(r.URL.Path, "/doc.json") || strings.HasSuffix(r.URL.Path, "/openapi.json"):
			serveOpenAPIJSON(w, r)
		case strings.HasSuffix(r.URL.Path, "/openapi.yaml"):
			serveOpenAPIYAML(w, r)
		default:
			http.NotFound(w, r)
		}
	})
	srv.Handle("/swagger", handler)
	srv.HandlePrefix("/swagger/", handler)
}

func serveSwaggerIndex(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(`<!doctype html>
<html lang="zh-CN">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>lightning swagger</title>
  <style>
    html { box-sizing: border-box; overflow: -moz-scrollbars-vertical; overflow-y: scroll; }
    *, *:before, *:after { box-sizing: inherit; }
    body { margin:0; background:#fafafa; }
  </style>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css" />
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-standalone-preset.js"></script>
  <script>
    window.onload = function() {
      window.ui = SwaggerUIBundle({
        url: "/swagger/doc.json",
        dom_id: "#swagger-ui",
        deepLinking: true,
        presets: [SwaggerUIBundle.presets.apis, SwaggerUIStandalonePreset],
        layout: "StandaloneLayout",
        docExpansion: "none",
        displayRequestDuration: true
      });
    };
  </script>
</body>
</html>`))
}

func serveOpenAPIYAML(w http.ResponseWriter, r *http.Request) {
	content, modTime, err := readOpenAPIContent()
	if err != nil {
		http.Error(w, "openapi spec not found, run `make proto-all` first", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/yaml; charset=utf-8")
	http.ServeContent(w, r, filepath.Base(cleanOpenAPISpecPath()), modTime, bytes.NewReader(content))
}

func serveOpenAPIJSON(w http.ResponseWriter, r *http.Request) {
	content, modTime, err := readOpenAPIContent()
	if err != nil {
		http.Error(w, "openapi spec not found, run `make proto-all` first", http.StatusNotFound)
		return
	}
	var payload any
	if err = yaml.Unmarshal(content, &payload); err != nil {
		http.Error(w, "failed to parse openapi spec", http.StatusInternalServerError)
		return
	}
	jsonContent, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "failed to encode openapi spec", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	http.ServeContent(w, r, "openapi.json", modTime, bytes.NewReader(jsonContent))
}

func readOpenAPIContent() ([]byte, time.Time, error) {
	specPath := cleanOpenAPISpecPath()
	content, err := os.ReadFile(specPath)
	if err != nil {
		return nil, time.Time{}, err
	}
	modTime := time.Time{}
	if info, statErr := os.Stat(specPath); statErr == nil {
		modTime = info.ModTime()
	}
	return content, modTime, nil
}

func cleanOpenAPISpecPath() string {
	specPath := openAPISpecPath
	if filepath.IsAbs(specPath) {
		return specPath
	}
	return filepath.Clean(specPath)
}
