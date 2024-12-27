package ws

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/vpramatarov/hardware-monitor/internal/hardware"
)

type Ws struct{}

type httpServer struct {
	messageBuffer    int // number of messages that we can send these subscribers at any given time
	mux              http.ServeMux
	subscribersMutex sync.Mutex
	subscribers      map[*subscriber]struct{}
}

type subscriber struct {
	msgs chan []byte
}

func NewWs() *Ws {
	return &Ws{}
}

func (h *Ws) DisplaySystemData() {
	fmt.Println("Starting monitor server on port 9000...")

	s := NewHttpServer()
	hardwareHtml := hardware.NewHardwareHtml()

	go func(srv *httpServer, h *hardware.HardwareHtml) {
		for {
			systemData := h.GetSystemSection()
			diskData := h.GetDiskSection()
			cpuData := h.GetCpuSection()

			timeStamp := time.Now().Format("2006-01-02 15:04:05")
			msg := []byte(`
				<div hx-swap-oob="innerHTML:#update-timestamp">
					<p><i style="color: green" class="fa fa-circle"></i> ` + timeStamp + `</p>
				</div>
				<div hx-swap-oob="innerHTML:#system-data">` + systemData + `</div>
				<div hx-swap-oob="innerHTML:#cpu-data">` + cpuData + `</div>
				<div hx-swap-oob="innerHTML:#disk-data">` + diskData + `</div>
			`)
			srv.publishMsg(msg)
			time.Sleep(time.Duration(hardware.SecondsInterval) * time.Second)
		}
	}(s, hardwareHtml)

	err := http.ListenAndServe(":9000", &s.mux)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func NewHttpServer() *httpServer {
	s := &httpServer{
		messageBuffer: 10,
		subscribers:   make(map[*subscriber]struct{}),
	}

	s.mux.Handle("/", http.FileServer(http.Dir("./htmx")))
	s.mux.HandleFunc("/ws", s.subscribeHandler)
	return s
}

func (s *httpServer) subscribeHandler(w http.ResponseWriter, r *http.Request) {
	err := s.subscribe(r.Context(), w, r)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (s *httpServer) addSubscriber(subscriber *subscriber) {
	s.subscribersMutex.Lock()
	s.subscribers[subscriber] = struct{}{}
	s.subscribersMutex.Unlock()
	fmt.Println("Added subscriber", subscriber)
}

func (s *httpServer) subscribe(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var wsConn *websocket.Conn
	subscriber := &subscriber{
		msgs: make(chan []byte, s.messageBuffer),
	}
	s.addSubscriber(subscriber)

	wsConn, err := websocket.Accept(w, r, nil)
	if err != nil {
		return err
	}
	defer wsConn.CloseNow()

	ctx = wsConn.CloseRead(ctx)
	for {
		select {
		case msg := <-subscriber.msgs:
			ctx, cancel := context.WithTimeout(ctx, time.Second*hardware.SecondsInterval)
			defer cancel()
			err := wsConn.Write(ctx, websocket.MessageText, msg)
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (s *httpServer) publishMsg(msg []byte) {
	s.subscribersMutex.Lock()
	defer s.subscribersMutex.Unlock()

	for s := range s.subscribers {
		s.msgs <- msg
	}
}
