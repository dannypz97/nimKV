package nimkv

import (
  // "fmt"
  "log"
  "strings"
  "io/ioutil"
  "net/http"
  "github.com/julienschmidt/httprouter"
  "github.com/go-yaml/yaml"
)

var cache Cacher

func index(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
  w.WriteHeader(http.StatusOK)
  w.Write([]byte("nimKV - (nimble Key-Value Store)"))
}

func getItem(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
  searchKey := params.ByName("key")

  item, err := cache.GetItem(searchKey)

  if err != nil {
    w.WriteHeader(http.StatusNotFound)
    w.Write([]byte(err.Error()))
    return
  }

  w.WriteHeader(http.StatusOK)
  w.Write([]byte(item.Value().(string)))
}

func setItem(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

}

func initRouter() *httprouter.Router {
  router := httprouter.New()

  // GET requests
  router.GET("/", index)
  router.GET("/cache/items/:key", getItem)

  // POST requests
  router.POST("/cache/items", setItem)

  return router
}

func setupCache() {
  cacheConfig, err := ioutil.ReadFile("config.yaml")

  if err != nil {
    log.Fatal("Error trying to read config.yaml.")
  }

  cacheBase := CacheBase{}
  err = yaml.Unmarshal(cacheConfig, &cacheBase)

  if err != nil {
    log.Fatal("Error trying to unmarshal cache config.")
  }

  var errs []error
  switch strings.ToUpper(cacheBase.Type) {
  case "LRU":
    // := can't be used for assignment to package-level variables
    cache, errs = NewLRUCache(&cacheBase)
    if errs != nil {
      log.Fatalf("Errors while building LRUCache: %v", errs)
    }

  default:
    log.Fatalf("Cache Type %v not supported.", cacheBase.Type)
  }
}

func StartServer() {
  log.Println("Starting Server...")

  setupCache()

  log.Fatal(http.ListenAndServe(":8080", initRouter()))
}
