package store

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"main/src/backend/types"
	"os"
	"strconv"
	"time"

	"net/http"

	"github.com/kelvins/geocoder"
)

func GetTaxiTrips(db *sql.DB) {

	// This function is NOT complete
	// It provides code-snippets for the data source: https://data.cityofchicago.org/Transportation/Taxi-Trips/wrvz-psew
	// You need to complete the implmentation and add the data source: https://data.cityofchicago.org/Transportation/Transportation-Network-Providers-Trips/m6dm-c72p

	// Data Collection needed from two data sources:
	// 1. https://data.cityofchicago.org/Transportation/Taxi-Trips/wrvz-psew
	// 2. https://data.cityofchicago.org/Transportation/Transportation-Network-Providers-Trips/m6dm-c72p

	fmt.Println("GetTaxiTrips: Collecting Taxi Trips Data")

	geocoder.ApiKey = "AIzaSyDnYxb-zx89dF3aXn_wojYkodKVtsKRUOA"

	drop_table := `drop table if exists taxi_trips`
	_, err := db.Exec(drop_table)
	if err != nil {
		panic(err)
	}

	create_table := `CREATE TABLE IF NOT EXISTS "taxi_trips" (
					"id"   SERIAL , 
					"trip_id" VARCHAR(255), 
					"trip_start_timestamp" TIMESTAMP WITH TIME ZONE, 
					"trip_end_timestamp" TIMESTAMP WITH TIME ZONE, 
					"pickup_centroid_latitude" DOUBLE PRECISION, 
					"pickup_centroid_longitude" DOUBLE PRECISION, 
					"dropoff_centroid_latitude" DOUBLE PRECISION, 
					"dropoff_centroid_longitude" DOUBLE PRECISION, 
					"pickup_zip_code" VARCHAR(255), 
					"dropoff_zip_code" VARCHAR(255), 
					PRIMARY KEY ("id") 
					);`

	_, _err := db.Exec(create_table)
	if _err != nil {
		panic(_err)
	}

	fmt.Println("Created Table for Taxi Trips")

	// While doing unit-testing keep the limit value to 500
	// later you could change it to 1000, 2000, 10,000, etc.
	var url = "https://data.cityofchicago.org/resource/wrvz-psew.json?$limit=500"

	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    300 * time.Second,
		DisableCompression: true,
	}

	client := &http.Client{Transport: tr}

	res, err := client.Get(url)

	if err != nil {
		panic(err)
	}

	fmt.Println("Received data from SODA REST API for Taxi Trips")

	body_1, _ := ioutil.ReadAll(res.Body)
	var taxi_trips_list_1 types.TaxiTripsJsonRecords
	json.Unmarshal(body_1, &taxi_trips_list_1)

	// Get the Taxi Trip list for rideshare companies like Uber/Lyft list
	// Transportation-Network-Providers-Trips:
	var url_2 = "https://data.cityofchicago.org/resource/m6dm-c72p.json?$limit=500"

	res_2, err := http.Get(url_2)
	if err != nil {
		panic(err)
	}

	fmt.Println("Received data from SODA REST API for Transportation-Network-Providers-Trips")

	body_2, _ := ioutil.ReadAll(res_2.Body)
	var taxi_trips_list_2 types.TaxiTripsJsonRecords
	json.Unmarshal(body_2, &taxi_trips_list_2)

	s := fmt.Sprintf("\n\n Transportation-Network-Providers-Trips number of SODA records received = %d\n\n", len(taxi_trips_list_2))
	io.WriteString(os.Stdout, s)

	// Add the Taxi medallions list & rideshare companies like Uber/Lyft list

	taxi_trips_list := append(taxi_trips_list_1, taxi_trips_list_2...)

	for i := 0; i < len(taxi_trips_list); i++ {

		// We will execute defensive coding to check for messy/dirty/missing data values
		// Any record that has messy/dirty/missing data we don't enter it in the data lake/table

		trip_id := taxi_trips_list[i].Trip_id
		if trip_id == "" {
			continue
		}

		// if trip start/end timestamp doesn't have the length of 23 chars in the format "0000-00-00T00:00:00.000"
		// skip this record

		// get Trip_start_timestamp
		trip_start_timestamp := taxi_trips_list[i].Trip_start_timestamp
		if len(trip_start_timestamp) < 23 {
			continue
		}

		// get Trip_end_timestamp
		trip_end_timestamp := taxi_trips_list[i].Trip_end_timestamp
		if len(trip_end_timestamp) < 23 {
			continue
		}

		pickup_centroid_latitude := taxi_trips_list[i].Pickup_centroid_latitude

		if pickup_centroid_latitude == "" {
			continue
		}

		pickup_centroid_longitude := taxi_trips_list[i].Pickup_centroid_longitude

		if pickup_centroid_longitude == "" {
			continue
		}

		dropoff_centroid_latitude := taxi_trips_list[i].Dropoff_centroid_latitude

		if dropoff_centroid_latitude == "" {
			continue
		}

		dropoff_centroid_longitude := taxi_trips_list[i].Dropoff_centroid_longitude

		if dropoff_centroid_longitude == "" {
			continue
		}

		// Using pickup_centroid_latitude and pickup_centroid_longitude in geocoder.GeocodingReverse
		// we could find the pickup zip-code

		pickup_centroid_latitude_float, _ := strconv.ParseFloat(pickup_centroid_latitude, 64)
		pickup_centroid_longitude_float, _ := strconv.ParseFloat(pickup_centroid_longitude, 64)
		pickup_location := geocoder.Location{
			Latitude:  pickup_centroid_latitude_float,
			Longitude: pickup_centroid_longitude_float,
		}

		// Comment the following line while not unit-testing
		//fmt.Println(pickup_location)

		pickup_address_list, _ := geocoder.GeocodingReverse(pickup_location)
		pickup_address := pickup_address_list[0]
		pickup_zip_code := pickup_address.PostalCode
		if pickup_zip_code == "" {
			continue
		}

		// Using dropoff_centroid_latitude and dropoff_centroid_longitude in geocoder.GeocodingReverse
		// we could find the dropoff zip-code

		dropoff_centroid_latitude_float, _ := strconv.ParseFloat(dropoff_centroid_latitude, 64)
		dropoff_centroid_longitude_float, _ := strconv.ParseFloat(dropoff_centroid_longitude, 64)

		dropoff_location := geocoder.Location{
			Latitude:  dropoff_centroid_latitude_float,
			Longitude: dropoff_centroid_longitude_float,
		}

		dropoff_address_list, _ := geocoder.GeocodingReverse(dropoff_location)
		dropoff_address := dropoff_address_list[0]
		dropoff_zip_code := dropoff_address.PostalCode
		if dropoff_zip_code == "" {
			continue
		}

		sql := `INSERT INTO taxi_trips ("trip_id", "trip_start_timestamp", "trip_end_timestamp", "pickup_centroid_latitude", "pickup_centroid_longitude", "dropoff_centroid_latitude", "dropoff_centroid_longitude", "pickup_zip_code", 
			"dropoff_zip_code") values($1, $2, $3, $4, $5, $6, $7, $8, $9)`

		_, err = db.Exec(
			sql,
			trip_id,
			trip_start_timestamp,
			trip_end_timestamp,
			pickup_centroid_latitude,
			pickup_centroid_longitude,
			dropoff_centroid_latitude,
			dropoff_centroid_longitude,
			pickup_zip_code,
			dropoff_zip_code)

		if err != nil {
			panic(err)
		}
	}

	fmt.Println("Completed Inserting Rows into the TaxiTrips Table")

}
