package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func NewHttpServer(addr string) *http.Server {
	httpserver := newHTTPServer()

	r := mux.NewRouter()

	r.HandleFunc("/", httpserver.handleProduce).Methods("POST")
	r.HandleFunc("/", httpserver.handleConsume).Methods("GET")
	return &http.Server{
		Addr:    addr,
		Handler: r,
	}
}

type httpServer struct {
	Log *Log
}

func newHTTPServer() *httpServer {
	return &httpServer{
		Log: NewLog(),
	}
}

type ProduceRequest struct {
	Record Record `json:"record`
}

type ProduceResponse struct {
	Offset uint64 `json:"offset"`
}

type ConsumeRequest struct {
	Offset uint64 `json:"offset"`
}

type ConsumeResponse struct {
	Record Record `json:"record"`
}

func (s *httpServer) handleProduce(rw http.ResponseWriter, r *http.Request) {
	// create empty ProduceRequest struct
	var req ProduceRequest

	// decode request body into the new struct
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	// append the decoded record to the log and get the new offset value
	off, err := s.Log.Append(req.Record)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	// create the new ProduceResponse
	res := ProduceResponse{Offset: off}

	// encode the ProduceResponse and send it through the response writer
	err = json.NewEncoder(rw).Encode(res)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (s *httpServer) handleConsume(rw http.ResponseWriter, r *http.Request) {
	var req ConsumeRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	rec, err := s.Log.Read(req.Offset)
	if err != nil {
		if err == ErrOffsetNotFound {
			http.Error(rw, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	res := ConsumeResponse{Record: rec}

	err = json.NewEncoder(rw).Encode(res)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}
