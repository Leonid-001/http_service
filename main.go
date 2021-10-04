package main

import (
	"context"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const PROCESSING_DELAY_IN_SECONDS = 5

// counts all hash password requests
var RequestCounter uint64

// total processing time for all requests (in microseconds)
var TotalProcessingTime int64

// average processing time for all requests (in microseconds)
var AverageProcessingTime int64

// number of processed requests
var NumberOfProcessedRequests uint64

// my HTTP server instance
var myServer http.Server

// Synchronized structure where keys are request ids and values are the hashed passwords
var hashes sync.Map

func main() {

	log.Println("Starting the http service")

	mux := http.NewServeMux()
	mux.HandleFunc("/hash", HashPassword)
	mux.HandleFunc("/hash/", GetHashedPassword)
	mux.HandleFunc("/stats", Stats)
	mux.HandleFunc("/shutdown", Shutdown)

	myServer = http.Server{Addr: ":8080", Handler: mux}

	if err := myServer.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting listener:
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
}

func HashPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		requestId := atomic.AddUint64(&RequestCounter, 1)
		log.Println("processing POST request", requestId)
		go func() {
			time.Sleep(PROCESSING_DELAY_IN_SECONDS * time.Second)
			log.Println("i am awaken and processing request", requestId)
			startTime := time.Now()
			EncodedValue := HashAndEncodeString(r.FormValue("password"))
			hashes.Store(requestId, EncodedValue)
			elapsedTime := time.Since(startTime)
			AllProcessedRequests := atomic.AddUint64(&NumberOfProcessedRequests, 1)
			UpdateAverageProcessingTime(AllProcessedRequests, elapsedTime.Microseconds())
		}()
		fmt.Fprint(w, requestId)
	} else {
		log.Println("HashPassword undefined behavior for the HTTP method", r.Method)
	}
}

// returns hash passwords, if available
func GetHashedPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		input := strings.TrimPrefix(r.URL.Path, "/hash/")
		log.Println("processing GET request", input)
		// convert the input string to the unsigned int64
		requestId, err := strconv.ParseUint(input, 10, 64)
		if err == nil {
			// return the hashed password only if the request id exists in hashes
			if hashedPassword, ok := hashes.Load(requestId); ok {
				fmt.Fprint(w, hashedPassword)
			}
		} else {
			log.Println("Error processing GET request", input, "is not an integer.")
		}
	} else {
		log.Println("GetHashedPassword undefined behavior for the HTTP method", r.Method)
	}
}

// returns current execution statistics in the JSON format
func Stats(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		log.Println("processing stats request")
		currentTotal := strconv.FormatUint(atomic.LoadUint64(&NumberOfProcessedRequests), 10)
		currentAverage := strconv.FormatInt(atomic.LoadInt64(&AverageProcessingTime), 10)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		jsonString := fmt.Sprintf("{“total”: %s, “average”: %s}", currentTotal, currentAverage)
		w.Write([]byte(jsonString))
	} else {
		log.Println("Stats undefined behavior for the HTTP method", r.Method)
	}
}

// waits PROCESSING_DELAY_IN_SECONDS and triggers graceful service shutdown
func Shutdown(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		log.Println("processing shutdown request")
		go func() {
			time.Sleep(PROCESSING_DELAY_IN_SECONDS * time.Second)
			log.Println("i am awaken and shuting down")
			myServer.Shutdown(context.Background())
		}()
	} else {
		log.Println("Shutdown undefined behavior for the HTTP method", r.Method)
	}
}

// Computes base64 encoding of the SHA512 hash of the provided input string
func HashAndEncodeString(input string) string {
	sha512 := sha512.New()
	sha512.Write([]byte(input))
	return base64.StdEncoding.EncodeToString([]byte(sha512.Sum(nil)))
}

// Computes and updates the new average execution time
func UpdateAverageProcessingTime(CurrentRequest uint64, CurrentExecutionTime int64) {
	newTotalTime := atomic.AddInt64(&TotalProcessingTime, CurrentExecutionTime)
	newAverageValue := int64(uint64(newTotalTime) / CurrentRequest)
	log.Println("newAverageValue", newAverageValue)
	atomic.StoreInt64(&AverageProcessingTime, newAverageValue)
}
