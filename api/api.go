package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"github.com/kaibling/apiforge/ctxkeys"
	"github.com/kaibling/cerodev/api/auth"
	"github.com/kaibling/cerodev/api/container"
	images "github.com/kaibling/cerodev/api/image"
	"github.com/kaibling/cerodev/api/middleware"
	"github.com/kaibling/cerodev/api/template"
	"github.com/kaibling/cerodev/api/user"
	"github.com/kaibling/cerodev/bootstrap"
)

func Route() chi.Router { //nolint: ireturn
	r := chi.NewRouter()
	r.Mount("/users", user.Route())
	r.Mount("/containers", container.Route())
	r.Mount("/templates", template.Route())
	r.Mount("/images", images.Route())
	r.Mount("/auth", auth.Route())
	r.Mount("/ws", WSRoute())

	return r
}

//	var upgrader = websocket.Upgrader{
//		CheckOrigin: func(r *http.Request) bool { return true }, // allow all for demo
//	}
var upgrader = websocket.Upgrader{}

func WSRoute() chi.Router { //nolint: ireturn
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Use(middleware.Authentication)
		r.Get("/", ws)
	})

	return r
}

func ws(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("upgrade failed:", err)
		return
	}

	conn.SetPongHandler(func(appData string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	token, ok := ctxkeys.GetValue(r.Context(), ctxkeys.TokenKey).(string)
	if !ok {
		fmt.Printf("ctxkeys TokenKey  failed: %s", err)
		return
	}
	wss, err := bootstrap.GetWebSocketService(r.Context())
	if err != nil {
		fmt.Printf("GetWebSocketService failed: %s", err)
		return
	}
	wss.Add(token, conn)
}
