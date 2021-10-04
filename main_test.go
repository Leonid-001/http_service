package main

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

func TestHashPassword(t *testing.T) {

	atomic.StoreUint64(&RequestCounter, 4)

	req, err := http.NewRequest("POST", "/hash", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HashPassword)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	result := rr.Body.String()

	if result != "5" {
		t.Errorf("HandlePasswordRequest returned %s; expecting 5", result)
	}
}

func TestGetHashedPasswordWithNoResult(t *testing.T) {

	req, err := http.NewRequest("GET", "/hash/42", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetHashedPassword)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	result := rr.Body.String()

	if result != "" {
		t.Errorf("GetHashedPassword returned %s; expecting no result", result)
	}
}

/*
// this test doesn't work
func TestGetHashedPasswordWithResult(t *testing.T) {

	hashes.Store(42, "ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q==")

	hashedPassword, _ := hashes.Load(42)

	t.Log("value is", hashedPassword)

	req, err := http.NewRequest("GET", "/hash/42", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetHashedPassword)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	result := rr.Body.String()

	if result != "ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q==" {
		t.Errorf("GetHashedPassword returned %s; expecting ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q==", result)
	}
}
*/

func TestShutdown(t *testing.T) {

	req, err := http.NewRequest("GET", "/shutdown", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Shutdown)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	result := rr.Body.String()

	if result != "" {
		t.Errorf("Shutdown returned %s; expecting no result", result)
	}
}

func TestStats(t *testing.T) {

	atomic.StoreUint64(&NumberOfProcessedRequests, 10)
	atomic.StoreInt64(&AverageProcessingTime, 25)

	req, err := http.NewRequest("GET", "/stats", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Stats)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	result := rr.Body.String()

	if result != "{“total”: 10, “average”: 25}" {
		t.Errorf("Stats returned %s; expecting {“total”: 10, “average”: 25}", result)
	}
}

func TestHashAndEncodeString(t *testing.T) {

	result := HashAndEncodeString("angryMonkey")

	if result != "ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q==" {
		t.Errorf("HashAndEncodeString returned %s; expecting \"ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q==\"", result)
	}

}

func TestUpdateAverageProcessingTime(t *testing.T) {

	atomic.StoreUint64(&NumberOfProcessedRequests, 4)
	atomic.StoreInt64(&TotalProcessingTime, 99)

	UpdateAverageProcessingTime(5, 35)

	ave1 := atomic.LoadInt64(&AverageProcessingTime)

	if ave1 != 26 {
		t.Errorf("UpdateAverageProcessingTime set average processing time to %d; expecting 26", ave1)
	}

	UpdateAverageProcessingTime(6, 41)

	ave2 := atomic.LoadInt64(&AverageProcessingTime)

	if ave2 != 29 {
		t.Errorf("UpdateAverageProcessingTime set average processing time to %d; expecting 26", ave2)
	}
}
