package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/hl540/http-log-proxy/configs"
	"github.com/hl540/http-log-proxy/internal/dashboard"
	"github.com/hl540/http-log-proxy/internal/http_log_proxy"
	"github.com/hl540/http-log-proxy/storage"
	"github.com/hl540/http-log-proxy/tools/log"
	"net/http"
)

var configFile string
var port string

func init() {
	flag.StringVar(&configFile, "f", "config.yaml", "the config file")
	flag.StringVar(&port, "p", "8080", "the http server port")
	flag.Parse()
}

func main() {
	// 加载配置
	conf, err := configs.Load(configFile)
	if err != nil {
		log.WithContext(nil).Fatal(err)
	}

	// 加载存储
	storageProvider, err := storage.Load(conf)
	if err != nil {
		log.WithContext(nil).Fatal(err)
	}

	mux := http.NewServeMux()

	// 创建proxy Handler
	httpLogProxy := http_log_proxy.NewHttpLogProxy(storageProvider)
	mux.Handle("/", httpLogProxy)

	// 创建dashboard路由
	dashboardApp := gin.Default()
	dashboardApp.Delims("[[", "]]")
	dashboardApp.LoadHTMLGlob("templates/*")

	dashboard.NewHandler(storageProvider).Register(dashboardApp.RouterGroup)
	mux.Handle("/dashboard/", dashboardApp)

	log.WithContext(nil).Infof("starting http server on http://127.0.0.1:" + port)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.WithContext(nil).Fatalf("listen failed: %v", err)
	}
}
