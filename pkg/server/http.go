package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"

	"github.com/lefuturiste/jobatator/pkg/store"
	log "github.com/sirupsen/logrus"
)

type gatewayHTTPInput struct {
	Commands []string `json:"commands"`
}

type httpError struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

type gatewayHTTPOutput struct {
	Success bool        `json:"success"`
	Data    []string    `json:"data"`
	Errors  []httpError `json:"errors"`
}

// StartHTTPServer -
func StartHTTPServer() {
	if store.Options.WebPort < 1 {
		return
	}

	http.HandleFunc("/gateway", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		if r.Method != "POST" {
			http.StatusText(http.StatusMethodNotAllowed)
			resp, _ := json.Marshal(gatewayHTTPOutput{
				Success: false,
				Errors:  []httpError{{Message: "Method not allowed", Code: "method-not-allowed"}},
			})
			fmt.Fprintf(w, string(resp))
			return
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.StatusText(http.StatusInternalServerError)
			resp, _ := json.Marshal(gatewayHTTPOutput{
				Success: false,
				Errors:  []httpError{{Message: "Invalid body", Code: "invalid-body"}},
			})
			fmt.Fprintf(w, string(resp))
			return
		}
		//fmt.Println(string(body))
		var bodyFormated gatewayHTTPInput
		err = json.Unmarshal(body, &bodyFormated)
		if err != nil || len(bodyFormated.Commands) == 0 {
			http.StatusText(http.StatusBadRequest)
			resp, _ := json.Marshal(gatewayHTTPOutput{
				Success: false,
				Errors:  []httpError{{Message: "Invalid body", Code: "invalid-body"}},
			})
			fmt.Fprintf(w, string(resp))
			return
		}

		// connect to this socket
		conn, _ := net.Dial("tcp", store.Options.Host+":"+strconv.FormatInt(int64(store.Options.Port), 10))

		output := gatewayHTTPOutput{
			Success: true,
			Errors:  []httpError{},
		}
		for _, command := range bodyFormated.Commands {
			// send to socket
			fmt.Fprintf(conn, command+"\n")
			// listen for reply
			message, _ := bufio.NewReader(conn).ReadString('\n')
			message = message[:len(message)-1]

			output.Data = append(output.Data, message)
		}
		conn.Close()

		resp, _ := json.Marshal(output)
		fmt.Fprintf(w, string(resp))
	})

	webListeningStr := store.Options.Host + ":" + strconv.FormatInt(int64(store.Options.WebPort), 10)

	log.Info("Web server listening on " + webListeningStr)
	http.ListenAndServe(webListeningStr, nil)
}
