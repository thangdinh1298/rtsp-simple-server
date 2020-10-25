package serverhttp

import (
	"net/http"
	"github.com/aler9/rtsp-simple-server/publisherman"
	"encoding/json"
	"fmt"
)

type Logger interface {
	Log(string, ...interface{})
}

type Server struct {
	httpServer *http.Server
	logger Logger
}

func (s *Server) getRoomInfo(w http.ResponseWriter, req *http.Request) {
	room := req.URL.Query().Get("room")
	s.logger.Log("[HTTP server] get room %s", room)

	if room == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	info, err := publisherman.GetInstance().GetRoomInfo(room)
	if err != nil {
		s.logger.Log("[HTTP server] ERR: %+v", err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%+v", err)
		return
	}

	bytes, err :=  json.Marshal(info)
	if err != nil {
		s.logger.Log("[HTTP server] ERR: %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Printf("%+v\n", info)
	w.Write(bytes)
}

func (s *Server) getRoomsInfo(w http.ResponseWriter, req *http.Request) {
	s.logger.Log("[HTTP server] get all rooms")
	infos := publisherman.GetInstance().GetRoomsInfo()

	bytes, err :=  json.Marshal(infos)
	if err != nil {
		s.logger.Log("[HTTP server] ERR: %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(bytes)
}

func New(programLogger Logger) {
	s := &Server {
		logger: programLogger,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/room", s.getRoomInfo)
	mux.HandleFunc("/rooms", s.getRoomsInfo)

	httpServer := &http.Server {
		Handler: mux,
		Addr: ":8080",
	}
	s.httpServer = httpServer

	go func() {
		s.logger.Log("[HTTP server] opened on port 8080")
		err := s.httpServer.ListenAndServe()

		if err != nil {
			s.logger.Log("[HTTP server] failed to start: %+v\n", err)
		}
	}()
}
