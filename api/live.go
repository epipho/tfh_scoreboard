package api

import (
	"sync"

	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
)

type Event struct {
	ID string `json:"id"`
}

type StartedEvent struct {
	Event
	Name    string    `json:"name"`
	Class   string    `json:"class"`
	CurRank int       `json:"cur_rank"`
	Ranks   []float32 `json:"ranks"`
}

type UpdatedEvent struct {
	Event
	Score float32 `json:"score"`
}

type FinalizedEvent struct {
	Event
}

type Broadcaster struct {
	sync.Mutex
	conns map[*websocket.Conn]bool
}

func (b *Broadcaster) AddConn(ws *websocket.Conn) {
	b.Lock()
	defer b.Unlock()
	b.conns[ws] = true
}

func (b *Broadcaster) RemoveConn(ws *websocket.Conn) {
	b.Lock()
	defer b.Unlock()
	delete(b.conns, ws)
}

func (b *Broadcaster) Broadcast(data interface{}) {
	b.Lock()
	defer b.Unlock()
	for ws := range b.conns {
		websocket.JSON.Send(ws, data)
	}
}

func (b *Broadcaster) Started(name string, class string, cur_rank int, ranks []float32) {
	evt := &StartedEvent{
		Event:   Event{ID: "started"},
		Name:    name,
		Class:   class,
		CurRank: cur_rank,
		Ranks:   ranks,
	}
	b.Broadcast(evt)
}

func (b *Broadcaster) Updated(score float32) {
	evt := &UpdatedEvent{
		Event: Event{ID: "updated"},
		Score: score,
	}
	b.Broadcast(evt)
}

func (b *Broadcaster) Finalized() {
	evt := &FinalizedEvent{
		Event: Event{ID: "finalized"},
	}
	b.Broadcast(evt)
}

func NewBroadcaster() *Broadcaster {
	return &Broadcaster{
		conns: make(map[*websocket.Conn]bool),
	}
}

func Live(b *Broadcaster) func(c echo.Context) error {
	return func(c echo.Context) error {
		websocket.Handler(func(ws *websocket.Conn) {
			defer ws.Close()
			b.AddConn(ws)
			// block until receiving any message or EOF
			msg := ""
			websocket.Message.Receive(ws, &msg)

			b.RemoveConn(ws)
		}).ServeHTTP(c.Response(), c.Request())
		return nil
	}
}
