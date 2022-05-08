package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

/*接收客户端 request，并将 request 中带的 header 写入 response header
读取当前系统的环境变量中的 VERSION 配置，并写入 response header
Server 端记录访问日志包括客户端 IP，HTTP 返回码，输出到 server 端的标准输出
当访问 {url}/healthz 时，应返回200*/

func index(w http.ResponseWriter, r *http.Request) { //responsewrite只响应一次
	//w.Write([]byte("hello, cloud"))

	//2.读取当前系统的环境变量中的 VERSION 配置，并写入 response header
	os.Setenv("VERSION","0.0.1")
	version := os.Getenv("VERSION")
	w.Header().Set("VERSION", version)
	//fmt.Printf(version)

	//1.接收客户端 request，并将 request 中带的 header 写入 response header
	for k, v := range r.Header {
		for _, vv := range v {
			fmt.Println(k, v)
			w.Header().Set(k, vv)
		}
	}

	clientIP := getCurrentIP(r)
	log.Printf("response code, %d", 200)
	log.Printf("client ip, %s", clientIP)
}

/* linux上也许适用
func getCurrentIP(r *http.Request) string {
	ip := r.Header.Get("X-Real-IP")

	//host:port -> remoteaddr
	if ip == "" {
		fmt.Printf(r.RemoteAddr)
		ip = strings.Split(r.RemoteAddr,":")[0]

	}

	return ip

}*/

func getCurrentIP(r *http.Request) string {
	xForwardedFor := r.Header.Get("X-Forwarded-For")   
	ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])   
	if ip != "" {      
		return ip   
	} 

	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))   
	if ip != "" {      
		return ip   
	}   
	
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}   
   return ""
}

func healthz(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w,"working\n")
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", index)
	mux.HandleFunc("/healthz", healthz)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("starting httpserver failed, %s\n", err.Error())
	}

}

