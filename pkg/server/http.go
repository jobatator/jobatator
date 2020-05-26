package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/lefuturiste/jobatator/pkg/store"
	"github.com/lefuturiste/jobatator/pkg/utils"
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

		output := gatewayHTTPOutput{
			Success: true,
			Errors:  []httpError{},
		}
		// if there is a referer header sent by the client, use this value as the cors host
		referer := r.Header.Get("Referer")
		if referer != "" {
			// extract the host from the referer URL using a regex
			expr, _ := regexp.Compile(`http(s)?:\/\/[a-zA-Z0-9.]+(:[0-9]{1,9})?`)
			match := expr.FindStringSubmatch(referer)
			if len(match) > 0 {
				w.Header().Set("Access-Control-Allow-Origin", match[0])
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Referer")
			}
		}

		if r.Method != "POST" && r.Method != "OPTIONS" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			resp, _ := json.Marshal(gatewayHTTPOutput{
				Success: false,
				Errors:  []httpError{{Message: "Method not allowed", Code: "method-not-allowed"}},
			})
			fmt.Fprintf(w, string(resp))
			return
		}
		if r.Method == "OPTIONS" {
			resp, _ := json.Marshal(output)
			fmt.Fprintf(w, string(resp))
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
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
			w.WriteHeader(http.StatusBadRequest)
			resp, _ := json.Marshal(gatewayHTTPOutput{
				Success: false,
				Errors:  []httpError{{Message: "Invalid body", Code: "invalid-body"}},
			})
			fmt.Fprintf(w, string(resp))
			return
		}

		// connect to this socket
		conn, _ := net.Dial("tcp", store.Options.Host+":"+strconv.FormatInt(int64(store.Options.Port), 10))
		jsonData := make([]string, 0)
		jsonIndex := 0
		jsonStartSequence := "$!$!__"
		jsonEndSequence := "__!$!$"
		for _, command := range bodyFormated.Commands {
			// send to socket
			fmt.Fprintf(conn, command+"\n")
			// listen for reply
			message, _ := bufio.NewReader(conn).ReadString('\n')
			message = message[:len(message)-1]

			if len(message) >= 3 && message[0:3] == "Err" {
				code := message
				if len(code) >= 5 {
					code = message[5:]
				}
				output.Errors = append(output.Errors, httpError{
					Code:    code,
					Message: "",
				})
			} else {
				if len(message) > 3 &&
					(message == "null" ||
						message[0:3] == `[{"` ||
						message[0:3] == `[["` ||
						message[0:2] == `{"` ||
						message[0:2] == `["`) {
					// consider this message as json do not insert directly into the data array
					jsonData = append(jsonData, message)
					message = jsonStartSequence + utils.IntToStr(jsonIndex) + jsonEndSequence
					jsonIndex++
				}
				output.Data = append(output.Data, message)
			}
		}
		conn.Close()

		output.Success = len(output.Errors) == 0
		if !output.Success {
			w.WriteHeader(http.StatusBadRequest)
		}
		resp, _ := json.Marshal(output)

		// replace all json ocurrences by directly injecting into the string
		// this code can be dangerous because it can trapped into a infinit loop
		// and it is subject to injection of the '$!$!__0__!$!$' sequence
		out := string(resp)
		hasError := false
		cursor := 0
		for strings.Index(out, jsonStartSequence) != -1 && !hasError {
			startIndex := strings.Index(out, jsonStartSequence)
			if startIndex < cursor {
				hasError = true
				break
			}
			endIndex := strings.Index(out, jsonEndSequence)
			if endIndex < cursor {
				hasError = true
				break
			}
			index := utils.StrToInt(out[startIndex+len(jsonStartSequence) : endIndex])
			if index < 0 || index > len(jsonData)-1 {
				hasError = true
				break
			}
			out = out[0:startIndex-1] + jsonData[index] + out[endIndex+len(jsonEndSequence)+1:]
			cursor = endIndex + len(jsonEndSequence) + 1 + len(jsonData[index])
		}
		fmt.Fprintf(w, out)
	})

	webListeningStr := store.Options.Host + ":" + utils.IntToStr(store.Options.WebPort)

	log.Info("Web server listening on " + webListeningStr)
	http.ListenAndServe(webListeningStr, nil)
}
