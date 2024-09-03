package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
)

type NumbersRequest struct {
    Numbers []int32 `json:"numbers"`
}

type ResultResponse struct {
    Result int32 `json:"result"`
}

func calculateHandler(w http.ResponseWriter, r *http.Request) {
    var numbersReq NumbersRequest

    // Decode JSON body into NumbersRequest struct
    err := json.NewDecoder(r.Body).Decode(&numbersReq)
    if err != nil {
        http.Error(w, "Invalid input. Please provide a JSON array of int32 numbers.", http.StatusBadRequest)
        return
    }

    // Calculate the sum of the numbers
    var sum int32
    for _, num := range numbersReq.Numbers {
        sum += num
    }

    // Return the result as JSON
    result := ResultResponse{Result: sum}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(result)
}

func main() {
    http.HandleFunc("/calculate", calculateHandler)

    fmt.Println("Server is running on http://localhost:8080...")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatalf("could not start server: %v\n", err)
    }
}
