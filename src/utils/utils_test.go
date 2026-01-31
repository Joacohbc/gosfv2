package utils

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestServeFileSafe(t *testing.T) {
	// Create a temp file to serve
	tmpFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Write([]byte("content"))
	tmpFile.Close()

	e := echo.New()

	tests := []struct {
		name           string
		filename       string
		expectedHeader string
	}{
		{
			name:           "Safe file (jpg)",
			filename:       "image.jpg",
			expectedHeader: "inline",
		},
		{
			name:           "Dangerous file (html)",
			filename:       "index.html",
			expectedHeader: "attachment",
		},
		{
			name:           "Dangerous file (svg)",
			filename:       "image.svg",
			expectedHeader: "attachment",
		},
		{
			name:           "Safe file (txt)",
			filename:       "notes.txt",
			expectedHeader: "inline",
		},
		{
			name:           "Case insensitive check",
			filename:       "INDEX.HTML",
			expectedHeader: "attachment",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := ServeFileSafe(c, tmpFile.Name(), tt.filename)
			if err != nil {
				t.Fatalf("ServeFileSafe failed: %v", err)
			}

			cd := rec.Header().Get("Content-Disposition")
			if !strings.Contains(cd, tt.expectedHeader) {
				t.Errorf("expected Content-Disposition to contain %q, got %q", tt.expectedHeader, cd)
			}

			if !strings.Contains(cd, tt.filename) {
				t.Errorf("expected Content-Disposition to contain filename %q, got %q", tt.filename, cd)
			}
		})
	}
}
