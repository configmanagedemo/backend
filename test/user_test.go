package test

import (
	"main/internal/pkg/web/router"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogin(t *testing.T) {
	r := router.InitRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/test", nil)
	r.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Error("return fail")
	}
}
