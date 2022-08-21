package store

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"main/src/pkg/types"
	"net/http"
	"os"
	"time"
)

func GetBuildingPermits(db *sql.DB) {
	fmt.Println("GetBuildingPermits: Collecting Building Permits Data")

	// This function is NOT complete
	// It provides code-snippets for the data source: https://data.cityofchicago.org/Buildings/Building-Permits/ydr8-5enu/data

	// Data Collection needed from data source:
	// https://data.cityofchicago.org/Buildings/Building-Permits/ydr8-5enu/data

	drop_table := `drop table if exists building_permits`
	_, err := db.Exec(drop_table)
	if err != nil {
		panic(err)
	}

	create_table := `CREATE TABLE IF NOT EXISTS "building_permits" (
						"id"   SERIAL , 
						"permit_id" VARCHAR(255) UNIQUE, 
						"permit_code" VARCHAR(255), 
						"permit_type" VARCHAR(255),  
						"review_type"      VARCHAR(255), 
						"application_start_date"      VARCHAR(255), 
						"issue_date"      VARCHAR(255), 
						"processing_time"      VARCHAR(255), 
						"street_number"      VARCHAR(255), 
						"street_direction"      VARCHAR(255), 
						"street_name"      VARCHAR(255), 
						"suffix"      VARCHAR(255), 
						"work_description"      TEXT, 
						"building_fee_paid"      VARCHAR(255), 
						"zoning_fee_paid"      VARCHAR(255), 
						"other_fee_paid"      VARCHAR(255), 
						"subtotal_paid"      VARCHAR(255), 
						"building_fee_unpaid"      VARCHAR(255), 
						"zoning_fee_unpaid"      VARCHAR(255), 
						"other_fee_unpaid"      VARCHAR(255), 
						"subtotal_unpaid"      VARCHAR(255), 
						"building_fee_waived"      VARCHAR(255), 
						"zoning_fee_waived"      VARCHAR(255), 
						"other_fee_waived"      VARCHAR(255), 
						"subtotal_waived"      VARCHAR(255), 
						"total_fee"      VARCHAR(255), 
						"contact_1_type"      VARCHAR(255), 
						"contact_1_name"      VARCHAR(255), 
						"contact_1_city"      VARCHAR(255), 
						"contact_1_state"      VARCHAR(255), 
						"contact_1_zipcode"      VARCHAR(255), 
						"reported_cost"      VARCHAR(255), 
						"pin1"      VARCHAR(255), 
						"pin2"      VARCHAR(255), 
						"community_area"      VARCHAR(255), 
						"census_tract"      VARCHAR(255), 
						"ward"      VARCHAR(255), 
						"xcoordinate"      DOUBLE PRECISION ,
						"ycoordinate"      DOUBLE PRECISION ,
						"latitude"      DOUBLE PRECISION ,
						"longitude"      DOUBLE PRECISION,
						PRIMARY KEY ("id") 
					);`

	_, _err := db.Exec(create_table)
	if _err != nil {
		panic(_err)
	}

	fmt.Println("Created Table for Building Permits")

	// While doing unit-testing keep the limit value to 500
	// later you could change it to 1000, 2000, 10,000, etc.
	var url = "https://data.cityofchicago.org/resource/building-permits.json?$limit=500"

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

	fmt.Println("Received data from SODA REST API for Building Permits")

	body, _ := ioutil.ReadAll(res.Body)
	var building_data_list types.BuildingPermitsJsonRecords
	json.Unmarshal(body, &building_data_list)

	s := fmt.Sprintf("\n\n Building Permits: number of SODA records received = %d\n\n", len(building_data_list))
	io.WriteString(os.Stdout, s)

	for i := 0; i < len(building_data_list); i++ {

		// We will execute defensive coding to check for messy/dirty/missing data values
		// There are different methods to deal with messy/dirty/missing data.
		// We will use the simplest method: drop records that have messy/dirty/missing data
		// Any record that has messy/dirty/missing data we don't enter it in the data lake/table

		permit_id := building_data_list[i].Id
		if permit_id == "" {
			continue
		}

		permit_code := building_data_list[i].Permit_Code
		if permit_code == "" {
			continue
		}

		permit_type := building_data_list[i].Permit_type
		if permit_type == "" {
			continue
		}

		review_type := building_data_list[i].Review_type
		if review_type == "" {
			continue
		}

		application_start_date := building_data_list[i].Application_start_date
		if application_start_date == "" {
			continue
		}
		issue_date := building_data_list[i].Issue_date
		if issue_date == "" {
			continue
		}
		processing_time := building_data_list[i].Processing_time
		if processing_time == "" {
			continue
		}

		street_number := building_data_list[i].Street_number
		if street_number == "" {
			continue
		}
		street_direction := building_data_list[i].Street_direction
		if street_direction == "" {
			continue
		}
		street_name := building_data_list[i].Street_name
		if street_name == "" {
			continue
		}
		suffix := building_data_list[i].Suffix
		if suffix == "" {
			continue
		}
		work_description := building_data_list[i].Work_description
		if work_description == "" {
			continue
		}
		building_fee_paid := building_data_list[i].Building_fee_paid
		if building_fee_paid == "" {
			continue
		}
		zoning_fee_paid := building_data_list[i].Zoning_fee_paid
		if zoning_fee_paid == "" {
			continue
		}
		other_fee_paid := building_data_list[i].Other_fee_paid
		if other_fee_paid == "" {
			continue
		}
		subtotal_paid := building_data_list[i].Subtotal_paid
		if subtotal_paid == "" {
			continue
		}
		building_fee_unpaid := building_data_list[i].Building_fee_unpaid
		if building_fee_unpaid == "" {
			continue
		}
		zoning_fee_unpaid := building_data_list[i].Zoning_fee_unpaid
		if zoning_fee_unpaid == "" {
			continue
		}
		other_fee_unpaid := building_data_list[i].Other_fee_unpaid
		if other_fee_unpaid == "" {
			continue
		}
		subtotal_unpaid := building_data_list[i].Subtotal_unpaid
		if subtotal_unpaid == "" {
			continue
		}
		building_fee_waived := building_data_list[i].Building_fee_waived
		if building_fee_waived == "" {
			continue
		}
		zoning_fee_waived := building_data_list[i].Zoning_fee_waived
		if zoning_fee_waived == "" {
			continue
		}
		other_fee_waived := building_data_list[i].Other_fee_waived
		if other_fee_waived == "" {
			continue
		}

		subtotal_waived := building_data_list[i].Subtotal_waived
		if subtotal_waived == "" {
			continue
		}
		total_fee := building_data_list[i].Total_fee
		if total_fee == "" {
			continue
		}

		contact_1_type := building_data_list[i].Contact_1_type
		if contact_1_type == "" {
			continue
		}

		contact_1_name := building_data_list[i].Contact_1_name
		if contact_1_name == "" {
			continue
		}

		contact_1_city := building_data_list[i].Contact_1_city
		if contact_1_city == "" {
			continue
		}
		contact_1_state := building_data_list[i].Contact_1_state
		if contact_1_state == "" {
			continue
		}

		contact_1_zipcode := building_data_list[i].Contact_1_zipcode
		if contact_1_zipcode == "" {
			continue
		}

		reported_cost := building_data_list[i].Reported_cost
		if reported_cost == "" {
			continue
		}

		pin1 := building_data_list[i].Pin1
		if pin1 == "" {
			continue
		}

		pin2 := building_data_list[i].Pin2

		community_area := building_data_list[i].Community_area

		census_tract := building_data_list[i].Census_tract
		if census_tract == "" {
			continue
		}

		ward := building_data_list[i].Ward
		if ward == "" {
			continue
		}

		xcoordinate := building_data_list[i].Xcoordinate

		ycoordinate := building_data_list[i].Ycoordinate

		latitude := building_data_list[i].Latitude
		if latitude == "" {
			continue
		}

		longitude := building_data_list[i].Longitude
		if longitude == "" {
			continue
		}

		sql := `INSERT INTO building_permits ("permit_id", "permit_code", "permit_type","review_type",
				"application_start_date",
				"issue_date",
				"processing_time",
				"street_number",
				"street_direction",
				"street_name",
				"suffix",
				"work_description",
				"building_fee_paid",
				"zoning_fee_paid",
				"other_fee_paid",
				"subtotal_paid",
				"building_fee_unpaid",
				"zoning_fee_unpaid",
				"other_fee_unpaid",
				"subtotal_unpaid",
				"building_fee_waived",
				"zoning_fee_waived",
				"other_fee_waived",
				"subtotal_waived",
				"total_fee",	
				"contact_1_type",
				"contact_1_name",
				"contact_1_city",
				"contact_1_state",
				"contact_1_zipcode",
				"reported_cost",
				"pin1",
				"pin2",
				"community_area",
				"census_tract",
				"ward",
				"xcoordinate",
				"ycoordinate",
				"latitude",
				"longitude") values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10,$11, $12, $13, $14, $15,$16, $17, $18, $19, $20,$21, $22, $23, $24, $25,$26, $27, $28, $29,$30,$31, $32, $33, $34, $35,$36, $37, $38, $39, $40)`

		_, err = db.Exec(
			sql,
			permit_id,
			permit_code,
			permit_type,
			review_type,
			application_start_date,
			issue_date,
			processing_time,
			street_number,
			street_direction,
			street_name,
			suffix,
			work_description,
			building_fee_paid,
			zoning_fee_paid,
			other_fee_paid,
			subtotal_paid,
			building_fee_unpaid,
			zoning_fee_unpaid,
			other_fee_unpaid,
			subtotal_unpaid,
			building_fee_waived,
			zoning_fee_waived,
			other_fee_waived,
			subtotal_waived,
			total_fee,
			contact_1_type,
			contact_1_name,
			contact_1_city,
			contact_1_state,
			contact_1_zipcode,
			reported_cost,
			pin1,
			pin2,
			community_area,
			census_tract,
			ward,
			xcoordinate,
			ycoordinate,
			latitude,
			longitude)

		if err != nil {
			log.Print("Failed to insert building permit data: ", err)
			panic(err)
		}

	}

	fmt.Println("Completed Inserting Rows into the Building Permits Table")
}
