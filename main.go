package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
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

    // Calculate the sum of the numbers
    var sum int32
    for _, num := range numbers {
        sum += num
    }

    // Create the result object
    result := struct {
        Result int32 `json:"result"`
    }{
        Result: sum,
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
