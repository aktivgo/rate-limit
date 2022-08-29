package app

import (
	"fmt"
	"time"
)

const (
	RequestsCount         = 5
	WindowBetweenRequests = time.Second
	BurstyRequestsCount   = 3
)

func Run() {
	limiter()

	time.Sleep(time.Second)
	fmt.Println()

	burstyLimiter()
}

func limiter() {
	requests := make(chan int, RequestsCount)

	for i := 0; i < RequestsCount; i++ {
		requests <- i + 1
	}
	close(requests)

	limiter := time.NewTicker(WindowBetweenRequests).C

	for req := range requests {
		<-limiter
		fmt.Println("request", req, time.Now())
	}
}

func burstyLimiter() {
	burstyLimiter := make(chan time.Time, BurstyRequestsCount)

	for i := 0; i < BurstyRequestsCount; i++ {
		burstyLimiter <- time.Now()
	}

	go func() {
		for t := range time.NewTicker(WindowBetweenRequests).C {
			burstyLimiter <- t
		}
	}()

	burstyRequests := make(chan int, RequestsCount)

	for i := 0; i < RequestsCount; i++ {
		burstyRequests <- i + 1
	}
	close(burstyRequests)

	for req := range burstyRequests {
		<-burstyLimiter
		fmt.Println("request", req, time.Now())
	}
}
