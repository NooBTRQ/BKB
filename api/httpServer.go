// 对外提供kv存储的服务
package apiServer

import (
	"net/http"
)

func StartHttp() error {

	http.HandleFunc("/Set", setHandle)
	http.HandleFunc("/Get", getHandle)
	return http.ListenAndServe("127.0.0.1:8081", nil)
}

func setHandle(w http.ResponseWriter, r *http.Request) {

}

func getHandle(w http.ResponseWriter, r *http.Request) {

}
