package todo

import (
	"fmt"
	"strings"
	"github.com/jmoiron/sqlx"
	"os"
)

var db *DB

type DB struct {
	engine *sqlx.DB
	source string
}

func (d *DB) connect() {
	var err error = nil
	d.engine, err = sqlx.Connect("sqlite3", d.source)
	if err != nil {
		panic(err)
	}
}

func (d *DB) InitDB() *sqlx.DB {
	if d.engine == nil {
		d.connect()
	}
	return d.engine
}

func (d *DB) Close() {
	if d.engine != nil {
		d.engine.Close()
	}
}

func (d *DB) CreateTable(table string, columns map[string]string ) {
	template := "CREATE TABLE IF NOT EXISTS %s (%s)"
	colTemplate := "%s %s"
	cols := []string{}
	for name, dataType := range columns {
		cols = append(cols, fmt.Sprintf(colTemplate, name, dataType))
	}
	query := fmt.Sprintf(template, table, strings.Join(cols, ","))
	d.engine.MustExec(query)
}


func NewDB(source string) *DB {
	db = &DB{source: source}
	db.InitDB()
	return db
}

func DeleteDB(force bool) {
	// ask if user is sure
	if !force {
		fmt.Println("Are you sure you want to delete the database? (y/n)")
		var response string
		fmt.Scanln(&response)
		if response != "y" {
			fmt.Println("Aborting")
			return
		}
	}
	source := db.source
	db.Close()
	db = nil
	err := os.Remove(source)
	if err != nil {
		panic(err)
	}
}
