package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strconv"
)

func calculateHandler(w http.ResponseWriter, r *http.Request) {
    // Declare a slice of int32 to hold the incoming numbers
    var numbers []int32

    // Decode JSON array directly into the numbers slice
    err := json.NewDecoder(r.Body).Decode(&numbers)
    if err != nil {
        http.Error(w, "Invalid input. Please provide a JSON array of int32 numbers.", http.StatusBadRequest)
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
    w.Write([]byte(result))
}

func main() {
    http.HandleFunc("/calculate", calculateHandler)

    fmt.Println("Server is running on http://localhost:8080...")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatalf("could not start server: %v\n", err)
    }
}
