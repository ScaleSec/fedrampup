package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

type Journaler struct {
	wg       sync.WaitGroup
	dataCh   chan []string
	quitCh   chan bool
	fd       *os.File
	journal  *csv.Writer
	Filename string
}

func NewJournaler() (*Journaler, error) {
	var err error
	journaler := &Journaler{}
	journaler.dataCh = make(chan []string)
	journaler.quitCh = make(chan bool)

	// Create temporary file to journal to
	journaler.fd, err = ioutil.TempFile("", "fedrampup.csv")
	if err != nil {
		return nil, err
	}
	journaler.Filename = journaler.fd.Name()
	journaler.journal = csv.NewWriter(journaler.fd)
	journaler.wg.Add(1)
	go journaler.listen()
	return journaler, nil
}

func (this *Journaler) Write(data []string) {
	this.dataCh <- data
}

func (this *Journaler) Close() {
	this.quitCh <- true
	this.wg.Wait()
}

func (this *Journaler) Purge() {
	os.Remove(this.Filename)
}

func (this *Journaler) listen() {
	for {
		select {
		case entry := <-this.dataCh:
			fmt.Println(entry)
			if err := this.journal.Write(entry); err != nil {
				log.Fatalln("error writing journal record to csv:", err)
			}
			this.journal.Flush()
		case <-this.quitCh:
			this.fd.Close()
			this.wg.Done()
			return
		}
	}
}
