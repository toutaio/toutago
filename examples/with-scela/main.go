package main

import (
"context"
"fmt"
"log"
"net/http"

"github.com/toutaio/toutago-scela-bus/pkg/scela"
"github.com/toutaio/toutago/pkg/touta/integration"
)

func main() {
bus := integration.NewScelaBus()
defer bus.Close()

_, err := bus.Subscribe("user.*", scela.HandlerFunc(func(ctx context.Context, msg scela.Message) error {
fmt.Printf("[EVENT] %s: %v\n", msg.Topic(), msg.Payload())
return nil
}))
if err != nil {
log.Fatal(err)
}

http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
if r.Method == http.MethodPost {
userData := map[string]interface{}{
"id":    "123",
"email": "user@example.com",
}
if err := bus.Publish(r.Context(), "user.created", userData); err != nil {
http.Error(w, err.Error(), http.StatusInternalServerError)
return
}
w.WriteHeader(http.StatusCreated)
fmt.Fprintf(w, "User created\n")
return
}

if r.Method == http.MethodDelete {
if err := bus.Publish(r.Context(), "user.deleted", map[string]interface{}{
"id": "123",
}); err != nil {
http.Error(w, err.Error(), http.StatusInternalServerError)
return
}
w.WriteHeader(http.StatusOK)
fmt.Fprintf(w, "User deleted\n")
return
}

http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
})

fmt.Println("Server starting on :8080")
fmt.Println("Try: curl -X POST http://localhost:8080/users")

if err := http.ListenAndServe(":8080", nil); err != nil {
log.Fatal(err)
}
}
