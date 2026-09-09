package middleware

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gcc798/lightning/internal/httpx"
	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
)

func TestStringIDConverter_PathParam(t *testing.T) {

	tests := []struct {
		name     string
		pathID   string
		expected int64
		hasError bool
	}{
		{
			name:     "Valid string ID",
			pathID:   "1234567890123456789",
			expected: 1234567890123456789,
			hasError: false,
		},
		{
			name:     "Valid numeric ID",
			pathID:   "123",
			expected: 123,
			hasError: false,
		},
		{
			name:     "Invalid ID",
			pathID:   "invalid",
			expected: 0,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httpx.NewRouter(echo.New())
			r.Use(StringIDConverter())
			r.GET("/test/:id", func(c *echo.Context) {
				if parsedID := c.Get("parsed_id"); parsedID != nil {
					c.JSON(200, map[string]any{"id": parsedID})
				} else {
					c.JSON(400, map[string]any{"error": "ID not parsed"})
				}
			})

			req := httptest.NewRequest("GET", "/test/"+tt.pathID, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if tt.hasError {
				assert.Equal(t, 400, w.Code)
			} else {
				assert.Equal(t, 200, w.Code)
				var response map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &response)
				// JSON numbers are float64 by default
				assert.Equal(t, float64(tt.expected), response["id"])
			}
		})
	}
}

func TestStringIDConverter_JSONBody_SingleID(t *testing.T) {

	tests := []struct {
		name     string
		body     map[string]interface{}
		expected map[string]interface{}
	}{
		{
			name: "Mixed string and number IDs",
			body: map[string]interface{}{
				"userId": "1234567890123456789",
				"orgId":  123,
			},
			expected: map[string]interface{}{
				"userId": float64(1234567890123456789),
				"orgId":  float64(123),
			},
		},
		{
			name: "Non-ID fields unchanged",
			body: map[string]interface{}{
				"userId":   "1234567890123456789",
				"userName": "test",
				"age":      "25",
			},
			expected: map[string]interface{}{
				"userId":   float64(1234567890123456789),
				"userName": "test",
				"age":      "25",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httpx.NewRouter(echo.New())
			r.Use(StringIDConverter())
			r.POST("/test", func(c *echo.Context) {
				var body map[string]interface{}
				if err := c.Bind(&body); err != nil {
					c.JSON(400, map[string]any{"error": err.Error()})
					return
				}
				c.JSON(200, body)
			})

			bodyBytes, _ := json.Marshal(tt.body)
			req := httptest.NewRequest("POST", "/test", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, 200, w.Code)
			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)
			assert.Equal(t, tt.expected, response)
		})
	}
}

func TestStringIDConverter_JSONBody_IDArray(t *testing.T) {

	tests := []struct {
		name     string
		body     map[string]interface{}
		expected map[string]interface{}
	}{
		{
			name: "Mixed string and number ID array",
			body: map[string]interface{}{
				"ids": []interface{}{"1234567890123456789", 123},
			},
			expected: map[string]interface{}{
				"ids": []interface{}{float64(1234567890123456789), float64(123)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httpx.NewRouter(echo.New())
			r.Use(StringIDConverter())
			r.POST("/test", func(c *echo.Context) {
				var body map[string]interface{}
				if err := c.Bind(&body); err != nil {
					c.JSON(400, map[string]any{"error": err.Error()})
					return
				}
				c.JSON(200, body)
			})

			bodyBytes, _ := json.Marshal(tt.body)
			req := httptest.NewRequest("POST", "/test", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, 200, w.Code)
			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)
			assert.Equal(t, tt.expected, response)
		})
	}
}

func TestStringIDConverter_JSONBody_NestedObjects(t *testing.T) {

	body := map[string]interface{}{
		"user": map[string]interface{}{
			"id":    "1234567890123456789",
			"orgId": "2234567890123456789",
		},
		"roles": []interface{}{
			map[string]interface{}{"roleId": "111111111111111111"},
			map[string]interface{}{"roleId": "222222222222222222"},
		},
	}

	expected := map[string]interface{}{
		"user": map[string]interface{}{
			"id":    float64(1234567890123456789),
			"orgId": float64(2234567890123456789),
		},
		"roles": []interface{}{
			map[string]interface{}{"roleId": float64(111111111111111111)},
			map[string]interface{}{"roleId": float64(222222222222222222)},
		},
	}

	r := httpx.NewRouter(echo.New())
	r.Use(StringIDConverter())
	r.POST("/test", func(c *echo.Context) {
		var body map[string]interface{}
		if err := c.Bind(&body); err != nil {
			c.JSON(400, map[string]any{"error": err.Error()})
			return
		}
		c.JSON(200, body)
	})

	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/test", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, expected, response)
}

func TestIsIDField(t *testing.T) {
	tests := []struct {
		fieldName string
		expected  bool
	}{
		{"id", true},
		{"ids", true},
		{"userId", true},
		{"userIds", true},
		{"orgId", true},
		{"roleId", true},
		{"userName", false},
		{"email", false},
		{"status", false},
	}

	for _, tt := range tests {
		t.Run(tt.fieldName, func(t *testing.T) {
			result := isIDField(tt.fieldName)
			assert.Equal(t, tt.expected, result)
		})
	}
}
