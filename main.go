package main

import (
	"flag"
	"github.com/hl540/http-log-proxy/config"
	"github.com/hl540/http-log-proxy/internal/dashboard"
	"github.com/hl540/http-log-proxy/internal/http_log_proxy"
	"github.com/hl540/http-log-proxy/storage"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

var configFile string
var port string

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	flag.StringVar(&configFile, "f", "config.yaml", "the config file")
	flag.StringVar(&port, "p", "8080", "the http server port")
	flag.Parse()
}

func main() {
	// 加载配置
	conf, err := config.Load(configFile)
	if err != nil {
		log.Fatal(err)
	}

	// 加载存储
	storageProvider, err := storage.Load(conf)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	// 创建proxy Handler
	httpLogProxy := http_log_proxy.NewHttpLogProxy(storageProvider)
	mux.Handle("/", httpLogProxy)

	// 创建dashboard路由
	router := httprouter.New()
	dashboard.NewHandler(storageProvider).Register(router)
	mux.Handle("/dashboard/", router)

	log.Println("starting http server on http://127.0.0.1:" + port)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("listen failed: %v", err)
	}
}
