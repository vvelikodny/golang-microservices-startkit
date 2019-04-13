package news

import (
	"encoding/json"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
	nats "github.com/nats-io/go-nats"
	"github.com/vvelikodny/golang-microservices-test/errors"
	"github.com/vvelikodny/golang-microservices-test/config"
)

func GetNewsByIdHandler(conn *nats.Conn) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		key := vars["id"]

		id, err := strconv.ParseInt(key, 10, 64)
		if err != nil {
			errors.HttpError(w, fmt.Sprintf("request error: %v", err), http.StatusInternalServerError)
			return
		}

		requestMsg, err := proto.Marshal(&GetNewsByIdRequest{Id: id})
		if err != nil {
			errors.HttpError(w, fmt.Sprintf("Entity marshal error: %v", err), http.StatusInternalServerError)
			return
		}

		msg, err := conn.Request(config.GetNewsChannel, requestMsg, 10*time.Second)
		if err != nil {
			errors.HttpError(w, fmt.Sprintf("request error: %v", err), http.StatusInternalServerError)
			return
		}

		var n News
		if err := proto.Unmarshal(msg.Data, &n); err != nil {
			errors.HttpError(w, fmt.Sprintf("Unmarshal error: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(n); err != nil {
			errors.HttpError(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func CreateNewsHandler(nats *nats.Conn) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var n News
		if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
			errors.HttpError(w, fmt.Sprintf("JSON parsing error: %v", err), http.StatusInternalServerError)
			return
		}

		if _, err := govalidator.ValidateStruct(n); err != nil {
			errors.HttpError(w, fmt.Sprintf("New entity validation error: %v", err), http.StatusInternalServerError)
			return
		}

		requestMsg, err := proto.Marshal(&n)
		if err != nil {
			errors.HttpError(w, fmt.Sprintf("Entity marshal error: %v", err), http.StatusInternalServerError)
			return
		}

		msg, err := nats.Request(config.CreateNewsChannel, requestMsg, 500*time.Millisecond)
		if err != nil {
			errors.HttpError(w, fmt.Sprintf("request error: %v", err), http.StatusInternalServerError)
			return
		}

		if err := proto.Unmarshal(msg.Data, &n); err != nil {
			errors.HttpError(w, fmt.Sprintf("Unmarshal error: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(n); err != nil {
			errors.HttpError(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
