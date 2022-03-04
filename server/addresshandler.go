package server

import (
	"eth-helper/ethevent"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

// createAddresss 增加地址
func createAddresss(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	query := r.URL.Query()
	num := query.Get("num")

	n, err := strconv.Atoi(num)
	if err != nil {
		n = 1
	}

	err = ethevent.NewAddresss(n)
	if err != nil {
		log.Errorf("createAddresss NewAddresss err: %v", err.Error())
		replyGeneric(w, fmt.Sprintf("%v", err), nil)

		return
	}

	// fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
	w.Write([]byte(fmt.Sprintf("createAddresss %d", 1)))
}
