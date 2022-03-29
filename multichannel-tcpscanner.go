package main

import (
	"fmt"
	"net"
	"sort"
)

func worker(ports, results chan int) {
	// func declared to accept two channels
	for p:= range ports {
	address := fmt.Sprintf("scanme.nmap.org:%d",p)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		results <- 0 // return/send 0 if the port is closed
		continue
	}
	conn.Close()
	results <- p // send the port if it is open
	}
}

func main() {
	ports:= make(chan int, 100)
	results:= make(chan int) // separate chan 
	// to communicate results from worker with the main thread
	var  openports []int // using slice to store the result; later can be sorted
	
	for i:= 0; i< cap(ports); i++ {
		go worker(ports,results)
	}
	
	go func() {
	// result-gathering-loop has to start before >100 items of work continue
		for i:= 1; i<= 1024; i++ {
			ports <- i
		}
	}()
	
	for i:= 0; i< 1024; i++ { // result-gathering-loop recieves 
	// on results chan 1024; if port != 0, appended to slice
		port := <-results
		if port != 0 {
			openports = append(openports,port)
		}
	}
	
	close(ports)
	close(results)
	sort.Ints(openports) // sorting that slice full of results
	for _,port := range openports {
		fmt.Printf("%d open\n",port)
	}
}
