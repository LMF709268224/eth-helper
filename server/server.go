package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/julienschmidt/httprouter"
)

// GenericHTTPRsp 通用的http response回复
type GenericHTTPRsp struct {
	// 错误码，0表示成功
	ErrCode int `json:"errorCode"`
	// 字符串消息，用于客户端显示具体的错误信息等
	ErrMessage string `json:"errorMessage"`
	//
	Result interface{} `json:"result,omitempty"`
}

func replyGeneric(w http.ResponseWriter, msg string, result interface{}) {
	gr := &GenericHTTPRsp{
		// ErrCode:    errCode,
		ErrMessage: msg,
		Result:     result,
	}

	b, err := json.Marshal(gr)
	if err != nil {
		fmt.Printf("replyGeneric failed: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

// StartHTTPServer 开启http服务
func StartHTTPServer(port string) {
	router := httprouter.New()
	router.POST("/CreateAddresss", createAddresss)

	log.Fatal(http.ListenAndServe(port, router))
}
