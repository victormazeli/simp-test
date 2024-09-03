package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "strconv"
    "syscall"
    "time"
)

const (
    maxPayloadSize = 1024 * 1024 // 1MB limit
)

func calculateHandler(w http.ResponseWriter, r *http.Request) {
    // Check if the Content-Type is application/json
    if r.Header.Get("Content-Type") != "application/json" {
        http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
        return
    }

    // Limit the size of the request body
    r.Body = http.MaxBytesReader(w, r.Body, maxPayloadSize)

    // Declare a slice of int32 to hold the incoming numbers
    var numbers []int32

    // Decode JSON array directly into the numbers slice
    err := json.NewDecoder(r.Body).Decode(&numbers)
    if err != nil {
        http.Error(w, "Invalid input or payload too large. Please provide a JSON array of int32 numbers.", http.StatusBadRequest)
        return
    }

    // Calculate the sum of the numbers using int64 to handle large sums
    var sum int64
    for _, num := range numbers {
        sum += int64(num)
    }

    // Convert sum to string to send as raw response
    result := strconv.FormatInt(sum, 10)

    // Set the response header to plain text and return the result as a raw string
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    if _, err := w.Write([]byte(result)); err != nil {
        log.Printf("Failed to write response: %v", err)
    }
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/calculate", calculateHandler)

    srv := &http.Server{
        Addr:    ":8080",
        Handler: mux,
    }

    go func() {
        fmt.Println("Server is running on http://localhost:8080...")
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Could not start server: %v\n", err)
        }
    }()

    // Graceful shutdown on signal interrupt or terminate
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    fmt.Println("Shutting down server...")

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    if err := srv.Shutdown(ctx); err != nil {
        log.Fatalf("Server forced to shutdown: %v", err)
    }

    fmt.Println("Server exiting")
}
