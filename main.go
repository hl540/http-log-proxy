package main

import (
	"flag"
	"github.com/hl540/http-log-proxy/config"
	"github.com/hl540/http-log-proxy/internal/http_log_proxy"
	"github.com/hl540/http-log-proxy/storage"
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

	// 创建proxy Handler
	httpLogProxy := http_log_proxy.NewHttpLogProxy(storageProvider)
	log.Println("starting http server on http://127.0.0.1:" + port)
	if err := http.ListenAndServe(":8080", httpLogProxy); err != nil {
		log.Fatalf("listen failed: %v", err)
	}
}
