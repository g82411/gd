package googleDrive

import (
	"context"
	"fmt"
	"net/http"
)

var server *http.Server

func StartCallbackServer(ctx context.Context, port int, codeChan chan string) {
	server = &http.Server{
		Addr: fmt.Sprintf(":%d", port),
	}
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		codeChan <- code
		close(codeChan)
	})
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		fmt.Printf("Error starting server: %v\n", err)
	} else {
		fmt.Println("Server closed")
	}
}

func StopCallbackServer(ctx context.Context) {
	server.Shutdown(ctx)
}
