package store

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"main/src/backend/types"
	"net/http"
	"os"
	"time"
)

func GetCovidDetails(db *sql.DB) {
	fmt.Println("GetCovidDetails: Collecting Covid Data")

	drop_table := `drop table if exists covid_weekly`
	_, err := db.Exec(drop_table)
	if err != nil {
		panic(err)
	}

	create_table := `CREATE TABLE IF NOT EXISTS "covid_weekly" (
						"id"   SERIAL , 
						"zip_code" VARCHAR(255), 
						"week_number" VARCHAR(255), 
						"week_start" TIMESTAMP WITH TIME ZONE, 
						"week_end" TIMESTAMP WITH TIME ZONE, 
						"cases_weekly" VARCHAR(255),	
						"cases_cumulative" VARCHAR(255), 
						"case_rate_weekly" VARCHAR(255), 
						"case_rate_cumulative" VARCHAR(255), 
						"percent_tested_positive_weekly" VARCHAR(255), 
						"percent_tested_positive_cumulative" VARCHAR(255),
						"population" VARCHAR(255),											
						PRIMARY KEY ("id") 
					);`

	_, _err := db.Exec(create_table)
	if _err != nil {
		panic(_err)
	}

	fmt.Println("Created Table for covid_weekly")

	// set limit to 500.
	var url = "https://data.cityofchicago.org/resource/yhhz-zm2v.json"

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

	fmt.Println("covid_weekly: Received data from SODA REST API for covid_weekly")

	body, _ := ioutil.ReadAll(res.Body)
	var covid_weekly_data_list types.CovidJsonRecords
	json.Unmarshal(body, &covid_weekly_data_list)

	s := fmt.Sprintf("\n\n Covid_weekly number of SODA records received = %d\n\n", len(covid_weekly_data_list))
	io.WriteString(os.Stdout, s)

	for i := 0; i < len(covid_weekly_data_list); i++ {

		// We will execute defensive coding to check for messy/dirty/missing data values
		// Any record that has messy/dirty/missing data we don't enter it in the data lake/table

		zip_code := covid_weekly_data_list[i].Zip_code
		if zip_code == "" {
			continue
		}

		week_start := covid_weekly_data_list[i].Week_start
		if week_start == "" {
			continue
		}

		week_number := covid_weekly_data_list[i].Week_number
		if week_number == "" {
			continue
		}

		week_end := covid_weekly_data_list[i].Week_end
		if week_end == "" {
			continue
		}

		population := covid_weekly_data_list[i].Population
		if population == "" {
			continue
		}

		percent_tested_positive_weekly := covid_weekly_data_list[i].Percent_tested_positive_weekly
		if percent_tested_positive_weekly == "" {
			continue
		}

		percent_tested_positive_cumulative := covid_weekly_data_list[i].Percent_tested_positive_cumulative
		if percent_tested_positive_cumulative == "" {
			continue
		}

		cases_weekly := covid_weekly_data_list[i].Cases_weekly
		if cases_weekly == "" {
			continue
		}

		cases_cumulative := covid_weekly_data_list[i].Cases_cumulative
		if cases_cumulative == "" {
			continue
		}

		case_rate_weekly := covid_weekly_data_list[i].Case_rate_weekly
		if case_rate_weekly == "" {
			continue
		}

		case_rate_cumulative := covid_weekly_data_list[i].Case_rate_cumulative
		if case_rate_cumulative == "" {
			continue
		}

		sql := `INSERT INTO covid_weekly ("zip_code", 
			"week_number" , 
			"week_start" , 
			"week_end" , 
			"cases_weekly" ,	
			"cases_cumulative" , 
			"case_rate_weekly" , 
			"case_rate_cumulative" , 
			"percent_tested_positive_weekly" , 
			"percent_tested_positive_cumulative" ,
			"population")
			values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

		_, err = db.Exec(
			sql,
			zip_code,
			week_number,
			week_start,
			week_end,
			cases_weekly,
			cases_cumulative,
			case_rate_weekly,
			case_rate_cumulative,
			percent_tested_positive_weekly,
			percent_tested_positive_cumulative,
			population)

		if err != nil {
			panic(err)
		}

	}

	fmt.Println("Completed Inserting Rows into the covid_weekly Table")

}
