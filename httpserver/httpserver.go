package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	versionKey     = "VERSION"
	defaultVersion = "0.0.0"
)

type Config struct {
	Version string `json:"VERSION,omitempty"`
}

func init() {
	data, err := os.ReadFile("./config.json")
	if err != nil {
		log.Printf("read config file failed, error: %s\n", err.Error())
		return
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		log.Printf("marshal config failed, error: %s\n", err.Error())
		return
	}

	if cfg.Version == "" {
		cfg.Version = defaultVersion
	}
	os.Setenv(versionKey, cfg.Version)

}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", HandlerFunc)
	//4.访问 localhost/healthz 时，用healthz函数处理，返回 200→http.StatusOK
	mux.HandleFunc("/healthz", healthz)

	if err := http.ListenAndServe("localhost:8000", mux); err != nil {
		log.Fatalf("start http server failed, error: %s\n", err.Error())
	}
}

func HandlerFunc(w http.ResponseWriter, r *http.Request) {
	status, rc := "accepted", http.StatusAccepted
	innerHandler(status, rc, w, r)
}

func healthz(w http.ResponseWriter, r *http.Request) {
	status, rc := "healthy", http.StatusOK
	innerHandler(status, rc, w, r)
}

func innerHandler(status string, rc int, w http.ResponseWriter, r *http.Request) {
	// 1.接收客户端 request，并将 request 中带的 header 写入 response header
	rsp := http.Response{Status: status, StatusCode: rc, Header: r.Header}

	// 2.读取当前系统的环境变量中的 VERSION 配置，并写入 response header
	rsp.Header[versionKey] = []string{os.Getenv(versionKey)}

	data, _ := json.Marshal(rsp)
	w.Write(data)
	remoteIP := strings.SplitN(r.RemoteAddr, ":", 2)[0]

	// 3.Server 端记录访问日志包括客户端 IP，HTTP 返回码，输出到 server 端的标准输出
	log.Printf("client ip: %s, http return code: %d, response: %+v\n", remoteIP, rc, rsp)
}
