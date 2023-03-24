package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var transact *TransactionLogger

func main() {
	err := initializeTransactionLog()
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()

	r.Use(loggingMiddleware)

	r.HandleFunc("/v1/{key}", keyValuePutHandler).Methods("PUT")
	r.HandleFunc("/v1/{key}", keyValueGeHandler).Methods("GET")
	r.HandleFunc("/v1/{key}", keyValueDeleteHandler).Methods("DELETE")

	r.HandleFunc("/v1", notAllowedHeader)
	r.HandleFunc("/v1/{key}", notAllowedHeader)

	log.Fatal(http.ListenAndServe(":8080", r))
}

func initializeTransactionLog() error {
	var err error

	transact, err = NewTransactionLogger("/tmp/transactions.log")
	if err != nil {
		return fmt.Errorf("failed to create transaction logger: %w", err)
	}

	events, errors := transact.ReadEvents()
	count, ok, e := 0, true, Event{}

	for ok && err == nil {
		select {
		case err, ok = <-errors:

		case e, ok = <-events:
			switch e.EventType {
			case EventDelete: // Got a DELETE event!
				err = Delete(e.Key)
				count++
			case EventPut: // Got a PUT event!
				err = Put(e.Key, e.Value)
				count++
			}
		}
	}

	log.Printf("%d events replayed\n", count)

	transact.Run()

	return err
}
