package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
)

const (
	versionKey     = "VERSION"
	defaultVersion = "0.0.0"
)

const (
	logLevelFlag    = "v"
	logLevelKey     = "LOGLEVEL"
	// 3. 日志等级（采用 glog 的 V level 功能，生产环境暂设 3，测试环境暂设 4）
	defaultLogLevel = 3
)

type Config struct {
	Version  string `json:"VERSION,omitempty"`
	LogLevel int    `json:"LOGLEVEL,omitempty"`
}

func main() {
	flag.Set(logLevelFlag, os.Getenv(logLevelKey))
	defer glog.Flush()

	router := mux.NewRouter()
	router.HandleFunc("/", rootHandler).Methods("GET")
	// 2. 探活（开放 "/healthz" RESTful API 为健康检查接口）
	router.HandleFunc("/healthz", healthzHandler).Methods("GET")

	srv := http.Server{
		Addr:    ":80",
		Handler: router,
	}

	// 1. 优雅终止（1.1 捕获 SIGTERM 等信号并完成子进程的优雅终止）
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			glog.V(4).Info("server start failed")
			log.Fatalf("start http server failed, error: %s\n", err.Error())
		}
	}()
	glog.V(2).Info("server started")

	sig := <-done
	glog.V(1).Info("server stopped, receive stop signal:", sig)

	// 1. 优雅终止（1.2 给予一定的时间用于优雅终止）
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		// extra handling here
		cancel()
	}()

	// 1. 优雅终止（1.3 在规定的时间内清理退出的子进程以避免僵尸进程）
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown failed, error: %s\n", err.Error())
	}
	glog.V(1).Info("server exited properly")
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	status, rc := "accepted", http.StatusAccepted
	innerHandler(status, rc, w, r)
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	status, rc := "healthy", http.StatusOK
	innerHandler(status, rc, w, r)
}

func innerHandler(status string, rc int, w http.ResponseWriter, r *http.Request) {
	rsp := http.Response{Status: status, StatusCode: rc, Header: r.Header}
	rsp.Header[versionKey] = []string{os.Getenv(versionKey)}
	rsp.Header[logLevelKey] = []string{os.Getenv(logLevelKey)}

	data, _ := json.Marshal(rsp)
	w.Write(data)
	remoteIP := strings.SplitN(r.RemoteAddr, ":", 2)[0]

	glog.V(3).Infof("client ip: %s, http return code: %d, response: %+v\n", remoteIP, rc, rsp)
}

func initEnv(cfg Config) {
	if cfg.Version == "" {
		cfg.Version = defaultVersion
	}
	// 3. 日志等级（仅测试环境打开 Debug 等级的日志，生产环境 Info 级别起步）
	if cfg.LogLevel == 0 {
		cfg.LogLevel = defaultLogLevel
	}
	os.Setenv(versionKey, cfg.Version)
	os.Setenv(logLevelKey, strconv.Itoa(cfg.LogLevel))
}

func init() {
	// 4. 配置和代码分离（从 config.json 中读取配置信息，写入环境变量）
	data, err := os.ReadFile("./config.json")
	if err != nil {
		glog.V(1).Infof("read config file failed, error: %s\n", err.Error())
		return
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		glog.V(1).Infof("marshal config failed, error: %s\n", err.Error())
		return
	}

	initEnv(cfg)
}
