package store

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"main/src/pkg/types"
	"net/http"
	"os"
	"time"
)

func GetCCVIDetails(db *sql.DB) {

	fmt.Println("GetCCVIDetails: Collecting CCVI Data")

	drop_table := `drop table if exists ccvi`
	_, err := db.Exec(drop_table)
	if err != nil {
		panic(err)
	}

	create_table := `CREATE TABLE IF NOT EXISTS "ccvi" (
						"id"   SERIAL , 
						"geography_type" VARCHAR(255), 
						"community_area_or_zip" VARCHAR(255), 
						"community_area_name" VARCHAR(255), 
						"ccvi_score" VARCHAR(255), 
						"ccvi_category" VARCHAR(255),												
						PRIMARY KEY ("id") 
					);`

	_, _err := db.Exec(create_table)
	if _err != nil {
		panic(_err)
	}

	fmt.Println("Created Table for ccvi")

	// set limit to 200.
	var url = "https://data.cityofchicago.org/resource/xhc6-88s9.json?$limit=200"

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

	fmt.Println("CCVI: Received data from SODA REST API for CCVI")

	body, _ := ioutil.ReadAll(res.Body)
	var ccvi_data_list types.CCVIJsonRecords
	json.Unmarshal(body, &ccvi_data_list)

	s := fmt.Sprintf("\n\n CCVI number of SODA records received = %d\n\n", len(ccvi_data_list))
	io.WriteString(os.Stdout, s)

	for i := 0; i < len(ccvi_data_list); i++ {

		// We will execute defensive coding to check for messy/dirty/missing data values
		// Any record that has messy/dirty/missing data we don't enter it in the data lake/table

		community_area_or_zip := ccvi_data_list[i].Community_area_or_ZIP_code
		if community_area_or_zip == "" {
			continue
		}

		community_area_name := ccvi_data_list[i].Community_name
		if community_area_name == "" {
			continue
		}

		geography_type := ccvi_data_list[i].Geography_type
		if geography_type == "" {
			continue
		}

		ccvi_category := ccvi_data_list[i].CCVI_category
		ccvi_score := ccvi_data_list[i].CCVI_score

		sql := `INSERT INTO ccvi ("geography_type" , 
			"community_area_name" , 
			"community_area_or_zip" , 
			"ccvi_score" , 
			"ccvi_category" )
			values($1, $2, $3, $4, $5)`

		_, err = db.Exec(
			sql,
			geography_type,
			community_area_name,
			community_area_or_zip,
			ccvi_score,
			ccvi_category)

		if err != nil {
			panic(err)
		}

	}

	fmt.Println("Completed Inserting Rows into the ccvi Table")

}
