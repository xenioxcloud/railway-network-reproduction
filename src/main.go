package main

import (
	"context"
	"log"
	"net/http"
	"network-bug/oxylabs"
	"os"
	"time"
)

func main() {
	ctx := context.Background()
	client := oxylabs.NewOxylabsHttpClient(os.Getenv("oxylabs_username"), os.Getenv("oxylabs_password"), os.Getenv("oxylabs_entry"))
	idx := 1
	for {
		req, err := http.NewRequest(http.MethodGet, "https://api.kadaster.nl/lvwoz/wozwaardeloket-api/v1/wozwaarde/nummeraanduiding/1641200001127387", nil)
		if err != nil {
			log.Fatalln(err)
		}
		_, err = client.Do(ctx, req)
		if err != nil && err != oxylabs.Err404 {
			log.Fatalln(err)
		}
		time.Sleep(time.Millisecond * 100)
		log.Printf("request %d succeeded\n", idx)
		idx++
	}
}
