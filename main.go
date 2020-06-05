package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/awa/go-iap/appstore"
	_ "github.com/heroku/x/hmetrics/onload"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/appstore", verifyAppstoreReceipt)

	port, ok := os.LookupEnv("PORT")

	if !ok {
		port = "8080"
	}

	log.Printf("Starting verify receipt server on port %s\n", port)
	if err := http.ListenAndServe(":" + port, nil); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}

type ReceiptData struct {
	Receipt string
}


func verifyAppstoreReceipt(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/appstore" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "POST":
		var receipt ReceiptData
		err := json.NewDecoder(r.Body).Decode(&receipt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		client := appstore.New()
		req := appstore.IAPRequest{
			ReceiptData: receipt.Receipt,
		}
		resp := &appstore.IAPResponse{}
		ctx := context.Background()
		error := client.Verify(ctx, req, resp)
		if error != nil {
			fmt.Fprintf(w, "Failed = %v\n", error)
		} else {
			fmt.Fprintf(w, "Success = %v\n", resp)
		}

	default:
		fmt.Fprintf(w, "Sorry, only POST methods are supported.")
	}
}
