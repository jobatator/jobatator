package server

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/lefuturiste/jobatator/pkg/commands"
	"github.com/lefuturiste/jobatator/pkg/store"
)

// PublishHTTPInput -
type PublishHTTPInput struct {
	Group   string
	Queue   string
	Type    string
	Payload string
}

// StartHTTPServer -
func StartHTTPServer() {
	http.HandleFunc("/publish", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.StatusText(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, "Err: method not allowed")
			return
		}
		if r.Header.Get("Authorization") == "" {
			http.StatusText(http.StatusUnauthorized)
			fmt.Fprintf(w, "Err: not logged")
			return
		}
		sDec, _ := base64.StdEncoding.DecodeString(strings.Replace(r.Header.Get("Authorization"), "Basic ", "", 1))
		authComponents := strings.Split(string(sDec), ":")
		if len(authComponents) != 2 {
			http.StatusText(http.StatusUnauthorized)
			fmt.Fprintf(w, "Err: invalid format for authorization header")
			return
		}
		authComponents[1] = strings.Replace(authComponents[1], "\n", "", 1)
		var user store.User
		for _, val := range store.Options.Users {
			if val.Username == authComponents[0] && val.Password == authComponents[1] {
				user = val
			}
		}
		if user.Username == "" {
			http.StatusText(http.StatusUnauthorized)
			fmt.Fprintf(w, "Err: invalid creds")
			return
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.StatusText(http.StatusInternalServerError)
			fmt.Fprintf(w, "Err: invalid body")
			return
		}

		var bodyFormated PublishHTTPInput
		json.Unmarshal(body, &bodyFormated)

		group, err := commands.UseGroupUniversal(bodyFormated.Group, user)
		if err != nil {
			http.StatusText(http.StatusBadRequest)
			fmt.Fprintf(w, "Err: "+err.Error())
			return
		}
		user.CurrentGroup = group

		parts := map[int]string{
			0: "PUBLISH",
			1: bodyFormated.Queue,
			2: bodyFormated.Type,
			3: bodyFormated.Payload,
		}
		commands.PublishUniversal(parts, user)
		fmt.Fprintf(w, "{\"success\": true}")
	})

	webListeningStr := store.Options.Host + ":" + strconv.FormatInt(int64(store.Options.WebPort), 10)
	http.ListenAndServe(webListeningStr, nil)
}
