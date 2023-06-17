package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

type Earthquake struct {
	Time            string
	Latitude        float64
	Longitude       float64
	Depth           float64
	Magnitude       float64
	MagType         string
	Nst             *int
	Gap             *float64
	Dmin            *float64
	Rms             float64
	Net             string
	ID              string
	Updated         string
	Place           string
	Type            string
	HorizontalError *float64
	DepthError      float64
	MagError        *float64
	MagNst         * int64
	Status          string
	LocationSource  string
	MagSource       string
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/about", aboutHandler)
	http.HandleFunc("/charts", chartsHandler)
	log.Println("Listening on port 80: http://localhost:80")
	log.Fatal(http.ListenAndServe(":80", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	location := r.URL.Query().Get("location")
	magnitude := r.URL.Query().Get("magnitude")
	magnitudegt := r.URL.Query().Get("magnitudegt")
	date := r.URL.Query().Get("date")

	// Connect to the SQLite database
	db, err := sql.Open("sqlite3", "earthquakes-30-days.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	query := "SELECT * FROM earth_quakes WHERE Place NOT NULL"

	if location != "" {
		query += " AND Place LIKE '%" + location + "%'"
	}

	if magnitude != "" {
		magnitudeValue, _ := strconv.ParseFloat(magnitude, 64)
		if magnitudegt == "greater" {
			query += fmt.Sprintf(" AND mag > %.2f", magnitudeValue)
		} else if magnitudegt == "less" {
			query += fmt.Sprintf(" AND mag < %.2f", magnitudeValue)
		} else {
			query += fmt.Sprintf(" AND mag = %.2f", magnitudeValue)
		}
	}

	if date != "" {
		query += " AND date(time) = '" + date + "'"
	}

	query += " ORDER BY time DESC LIMIT 100"

	fmt.Println(query)

	// Fetch the data from the earth_quakes table
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Define a slice to hold Earthquake objects
	var earthquakes []Earthquake

	// Iterate over the rows and populate earthquakes slice
	for rows.Next() {
		var eq Earthquake
		err := rows.Scan(
			&eq.Time, &eq.Latitude, &eq.Longitude, &eq.Depth, &eq.Magnitude, &eq.MagType,
			&eq.Nst, &eq.Gap, &eq.Dmin, &eq.Rms, &eq.Net, &eq.ID, &eq.Updated, &eq.Place,
			&eq.Type, &eq.HorizontalError, &eq.DepthError, &eq.MagError, &eq.MagNst,
			&eq.Status, &eq.LocationSource, &eq.MagSource,
		)
		if err != nil {
			log.Fatal(err)
		}
		earthquakes = append(earthquakes, eq)
	}

	// Render the template with the earthquakes data
	tmpl, err := template.ParseFiles("./templates/index.html")
	if err != nil {
		log.Fatal(err)
	}
	err = tmpl.Execute(w, earthquakes)
	if err != nil {
		log.Fatal(err)
	}
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/about.html")
	if err != nil {
		log.Fatal(err)
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func chartsHandler(w http.ResponseWriter, r *http.Request) {
	query := "SELECT time, mag, depth FROM earth_quakes WHERE time >= datetime('now', '-7 days')"
	db, err := sql.Open("sqlite3", "earthquakes-30-days.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var earthquakes []Earthquake
	for rows.Next() {
		var eq Earthquake
		err := rows.Scan(&eq.Time, &eq.Magnitude, &eq.Depth)
		if err != nil {
			log.Fatal(err)
		}
		earthquakes = append(earthquakes, eq)
	}

	tmpl, err := template.ParseFiles("templates/chart.html")
	if err != nil {
		log.Fatal(err)
	}
	err = tmpl.Execute(w, earthquakes)
	if err != nil {
		log.Fatal(err)
	}
}
