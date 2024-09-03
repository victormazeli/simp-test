package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "math"
)

func calculateHandler(w http.ResponseWriter, r *http.Request) {
    var numbers []int32

    // Decode JSON array directly into the numbers slice
    err := json.NewDecoder(r.Body).Decode(&numbers)
    if err != nil {
        http.Error(w, "Invalid input. Please provide a JSON array of int32 numbers.", http.StatusBadRequest)
        return
    }

    // Initialize sum as an int64 to prevent overflow
    var sum int64

    // Calculate the sum of the numbers
    for _, num := range numbers {
        sum += int64(num)
    }

    // Check if sum exceeds int32 limits
    if sum > math.MaxInt32 || sum < math.MinInt32 {
        http.Error(w, "Sum exceeds int32 range.", http.StatusBadRequest)
        return
    }

    // Cast sum back to int32 for output
    result := struct {
        Result int32 `json:"result"`
    }{
        Result: int32(sum),
    }

    // Set the response header to JSON and return the result
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(result); err != nil {
        http.Error(w, "Failed to encode response.", http.StatusInternalServerError)
        return
    }
}

func main() {
    http.HandleFunc("/calculate", calculateHandler)

    fmt.Println("Server is running on http://localhost:8080...")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatalf("could not start server: %v\n", err)
    }
}
