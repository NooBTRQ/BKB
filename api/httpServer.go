// 对外提供kv存储的服务
package apiServer

import (
	"BlackKingBar/infrastructure"
	"net/http"
)

func StartHttp() error {

	http.HandleFunc("/Set", setHandle)
	http.HandleFunc("/Get", getHandle)
	http.HandleFunc("/Delete", deleteHandle)
	cfg := infrastructure.CfgInstance
	return http.ListenAndServe(cfg.HttpIP+":"+cfg.HttpPort, nil)
}

func setHandle(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world!"))
}

func getHandle(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world!"))
}

func deleteHandle(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world!"))
}
