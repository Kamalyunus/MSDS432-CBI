package main

import (
	"fmt"
	"log"
	"main/src/pkg/api"
	"main/src/pkg/store"
	"net/http"
	"os"
	"time"

	"database/sql"

	_ "github.com/lib/pq"
)

func main() {

	var err error

	fmt.Println("Initializing the DB connection")

	// Establish connection to Postgres Database
	db_connection := "user=postgres dbname=chicago_business_intelligence password=root host=/cloudsql/cbi-yunus:us-central1:cbipostgres sslmode=disable port = 5432"

	db, err := sql.Open("postgres", db_connection)
	if err != nil {
		panic(err)
	}

	// Test the database connection
	err = db.Ping()
	if err != nil {
		log.Fatal(fmt.Println("Couldn't Open Connection to database"))
		panic(err)
	}

	// Spin in a loop and pull data from the city of chicago data portal
	// Once every hour, day, week, etc.
	// Though, please note that Not all datasets need to be pulled on daily basis
	// fine-tune the following code-snippet as you see necessary

	dayTicker := time.NewTicker(24 * time.Hour)
	weekTicker := time.NewTicker(24 * 7 * time.Hour)

	for {

		log.Print("starting CBI Microservices ...")
		// pull Taxi Trips and COVID data on daily basis

		//These datasets are one time pull as they are historical
		go store.GetCCVIDetails(db)
		go store.GetCommunityAreaUnemployment(db)

		//one time pull to get history and then refreshed on regular cadence as schedule below
		go store.GetBuildingPermits(db)
		go store.GetTaxiTrips(db)
		go store.GetCovidDetails(db)

		http.HandleFunc("/", api.Handler)
		// Determine port for HTTP service.
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
			log.Printf("defaulting to port %s", port)
		}

		// Start HTTP server.
		log.Printf("listening on port %s", port)
		log.Print("Navigate to Cloud Run services and find the URL of your service")
		log.Print("Use the browser and navigate to your service URL to to check your service has started")

		if err := http.ListenAndServe(":"+port, nil); err != nil {
			log.Fatal(err)
		}

		select {
		case <-dayTicker.C:
			go store.GetBuildingPermits(db)
		case <-weekTicker.C:
			go store.GetTaxiTrips(db)
			go store.GetCovidDetails(db)
		}
	}
}
