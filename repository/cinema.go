package repository

import (
	"cinema-management/structs"
	"database/sql"
)

func GetAllCinema(db *sql.DB) (result []structs.Cinema, err error) {
	sql := "SELECT * FROM cinema"

	rows, err := db.Query(sql)
	if err != nil {
		return
	}

	defer rows.Close()
	for rows.Next() {
		var cinema structs.Cinema

		err = rows.Scan(&cinema.ID, &cinema.Nama, &cinema.Lokasi, &cinema.Rating)
		if err != nil {
			return
		}

		result = append(result, cinema)
	}

	return
}

func InsertCinema(db *sql.DB, cinema structs.Cinema) (err error) {
	sql := "INSERT INTO cinema(id, name, location, rating) VALUES ($1, $2, $3, $4)"

	errs := db.QueryRow(sql, cinema.ID, cinema.Nama, cinema.Lokasi, cinema.Rating)

	return errs.Err()
}

func UpdateCinema(db *sql.DB, cinema structs.Cinema) (err error) {
	sql := "UPDATE cinema SET name = $1, location = $2, rating = $3 WHERE id = $4"

	errs := db.QueryRow(sql, cinema.Nama, cinema.Lokasi, cinema.Rating, cinema.ID)

	return errs.Err()
}

func DeleteCinema(db *sql.DB, cinema structs.Cinema) (err error) {
	sql := "DELETE FROM cinema WHERE id = $1"

	errs := db.QueryRow(sql, cinema.ID)
	return errs.Err()
}
