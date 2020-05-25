package nimkv


import (
  "net/http"
  "net/http/httptest"
  "testing"
  "github.com/julienschmidt/httprouter"
)

func TestIndexHandler(t *testing.T) {
  router := httprouter.New()
  router.GET("/", indexHandler)

  req, err := http.NewRequest("GET", "/", nil)

  if err != nil {
    t.Error(err)
  }

  recorder := httptest.NewRecorder()

  router.ServeHTTP(recorder, req)

  if status := recorder.Code; status != http.StatusOK {
    t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
  }

  if recorder.Body.String() != "nimKV - (nimble Key-Value Store)" {
    t.Error("Unexpected response.")
  }
}
