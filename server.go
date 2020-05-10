package nimkv

import (
  // "fmt"
  "log"
  "strings"
  "encoding/json"
  "io/ioutil"
  "net/http"
  "github.com/julienschmidt/httprouter"
  "github.com/go-yaml/yaml"
)

var cache Cacher

func indexHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
  w.WriteHeader(http.StatusOK)
  w.Write([]byte("nimKV - (nimble Key-Value Store)"))
}

func getItemHandler(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
  searchKey := params.ByName("key")

  item, err := cache.GetItem(searchKey)

  if err != nil {
    w.WriteHeader(http.StatusNotFound)
    w.Write([]byte(err.Error()))
    return
  }

  itemJson, _ := json.Marshal(item)

  w.Header().Set("Content-Type", "application/json; charset=utf-8")
  w.WriteHeader(http.StatusOK)
  w.Write([]byte(itemJson))
}

func getAllItemsHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
  items := cache.GetAllItems()
  itemsJson, _ := json.Marshal(items)

  w.Header().Set("Content-Type", "application/json; charset=utf-8")
  w.WriteHeader(http.StatusOK)
  w.Write([]byte(itemsJson))
}

func deleteItemHandler(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
  searchKey := params.ByName("key")

  err := cache.DeleteItem(searchKey)

  if err != nil {
    w.WriteHeader(http.StatusNotFound)
    w.Write([]byte(err.Error()))
    return
  }

  w.WriteHeader(http.StatusOK)
}

func purgeCacheHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
  cache.Purge()

  w.WriteHeader(http.StatusOK)
}

func setItemHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
  item := cacheItem{}

  err := json.NewDecoder(req.Body).Decode(&item)

  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    w.Write([]byte(err.Error()))
    return
  }

  cache.SetItemWithExpiry(item.Key, item.Value, item.TTL)
  w.WriteHeader(http.StatusOK)
}

func initRouter() *httprouter.Router {
  router := httprouter.New()

  router.GET("/", indexHandler)
  router.GET("/cache/items/:key", getItemHandler)
  router.GET("/cache/items/", getAllItemsHandler)

  router.POST("/cache/items", setItemHandler)

  router.DELETE("/cache/items/:key", deleteItemHandler)
  router.DELETE("/cache/items/", purgeCacheHandler)

  return router
}

func setupCache() {
  cacheConfig, err := ioutil.ReadFile("config.yaml")

  if err != nil {
    log.Fatal("Error trying to read config.yaml.")
  }

  cacheBase := cacheBase{}
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
