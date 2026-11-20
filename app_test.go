package main

import (
 "encoding/json"
 "net/http"
 "net/http/httptest"
 "strings"
 "testing"
)

func TestRoutes(t *testing.T) {
 for _, tc := range []struct{method, path string; code int; content string}{
  {"GET", "/", 200, "StudyOps"},
  {"GET", "/healthz", 200, `"status":"ok"`},
  {"GET", "/missing", 404, "404"},
  {"POST", "/healthz", 405, "Method Not Allowed"},
 } {
  t.Run(tc.method+tc.path, func(t *testing.T) {
   w := httptest.NewRecorder()
   routes().ServeHTTP(w, httptest.NewRequest(tc.method, tc.path, nil))
   if w.Code != tc.code || !strings.Contains(w.Body.String(), tc.content) { t.Fatalf("status=%d body=%s", w.Code, w.Body.String()) }
  })
 }
}

func TestInfo(t *testing.T) {
 t.Setenv("APP_ENV", "staging")
 w := httptest.NewRecorder()
 routes().ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/info", nil))
 var info struct { Version string; Environment string; Uptime int `json:"uptime_seconds"` }
 if err := json.Unmarshal(w.Body.Bytes(), &info); err != nil { t.Fatal(err) }
 if w.Code != 200 || info.Environment != "staging" || info.Version != version || info.Uptime < 0 { t.Fatalf("unexpected info: %+v", info) }
}
