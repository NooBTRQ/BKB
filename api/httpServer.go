// 对外提供kv存储的服务
package apiServer

import (
	"BlackKingBar/infrastructure"
	"BlackKingBar/raft"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

func StartHttp(ctx context.Context) {
	http.HandleFunc("/Set", setHandle)
	http.HandleFunc("/Get", getHandle)
	http.HandleFunc("/Delete", deleteHandle)
	cfg := infrastructure.CfgInstance
	srv := &http.Server{Addr: cfg.HttpIP + ":" + cfg.HttpPort}
	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			panic("start httpServer err," + err.Error())
		}
	}()

	<-ctx.Done()
	srv.Shutdown(context.TODO())
}

func getHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	query := r.URL.Query()
	key := query.Get("key")

	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	value, has := raft.R.GetValue(key)
	if has {
		w.Write([]byte(value))
	} else {
		w.Write([]byte("not found"))
	}

}

func setHandle(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var data raft.Command
	json.Unmarshal([]byte(body), &data)
	data.Operation = raft.Set
	data.Done = make(chan bool)
	raft.R.EventCh <- &data

	select {
	case <-time.After(10 * time.Second):
		w.Write([]byte{0})
	case success := <-data.Done:
		if success {
			w.Write([]byte("写入成功！"))
		} else {
			w.Write([]byte("写入失败！"))
		}
	}
}

func deleteHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	query := r.URL.Query()
	key := query.Get("key")

	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data := &raft.Command{Operation: raft.Delete, Key: key}
	data.Done = make(chan bool)

	raft.R.EventCh <- &data

	select {
	case <-time.After(10 * time.Second):
		w.Write([]byte{0})
	case success := <-data.Done:
		if success {
			w.Write([]byte("写入成功！"))
		} else {
			w.Write([]byte("写入失败！"))
		}
	}
}
