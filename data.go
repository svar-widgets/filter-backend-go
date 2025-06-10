package main

// demo works with PersonData table only

import (
	"fmt"
	"io"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/xbsoftware/querysql/v2"
)

type PersonData struct {
	ID        int    `json:"id" db:"id"`
	LastName  string `json:"last_name" db:"last_name"`
	FirstName string `json:"first_name" db:"first_name"`
	Birthdate string `json:"birthdate" db:"birthdate"`
	Country   string `json:"country" db:"country"`
	City      string `json:"city" db:"city"`
	Email     string `json:"email" db:"email"`
	Phone     string `json:"phone" db:"phone"`
	Job       string `json:"job" db:"job"`
	Address   string `json:"address" db:"address"`
	CompanyID int    `json:"company_id" db:"company_id"`
	Notify    int    `json:"notify" db:"notify"`
	Age       int    `json:"age" db:"age"`
}

func getDataFromDB(w io.Writer, db *sqlx.DB, name string, body []byte) error {
	var filter = querysql.Filter{}
	var err error

	if len(body) > 0 {
		filter, err = querysql.FromJSON(body)
		if err != nil {
			return err
		}
	}

	query, data, err := querysql.GetSQL(filter, nil)
	if err != nil {
		return err
	}

	t := make([]PersonData, 0)
	sql := "select * from " + name
	if query != "" {
		sql += " where " + query
	}

	err = db.Select(&t, sql, data...)
	if err != nil {
		return err
	}

	format.JSON(w, 200, t)
	return nil
}

func getSuggestFromDB(w io.Writer, db *sqlx.DB, table, field string) error {
	sql := fmt.Sprintf("select distinct %s from %s ORDER BY %s ASC", field, table, field)

	// string fields
	if field == "last_name" || field == "first_name" || field == "birthdate" || field == "country" || field == "city" || field == "address" || field == "job" || field == "phone" || field == "email" {
		out := make([]string, 0)
		err := db.Select(&out, sql)
		if err != nil {
			fmt.Printf("%+v", err)
		}
		format.JSON(w, 200, out)
		return nil
	}

	// numeric fields
	if field == "company_id" || field == "notify" || field == "age" {
		out := make([]float64, 0)
		err := db.Select(&out, sql)
		if err != nil {
			log.Printf("%+v", err)
		}
		format.JSON(w, 200, out)
		return nil
	}

	return nil
}
