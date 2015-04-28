package main

import(
	"os"
	"bufio"
	"database/sql"
	"gopkg.in/gorp.v1"
	_ "github.com/lib/pq"
	"log"
	"time"
	"strings"
	"fmt"
)


func main() {
	// initialize the DbMap
	dbmap := initDb()
	defer dbmap.Db.Close()

	// delete any existing rows
	err := dbmap.TruncateTables()
	checkErr(err, "TruncateTables failed")

	readLines("all_products.txt", dbmap)
}

type Audible struct {
	// db tag lets you specify the column name if it differs from the struct field
	Id      				int64 `db:"post_id"`
	Created 				int64
	LastUpdated 			string
	Name	   				string
	Category				string
	Keywords				string
	ShortDescription 		string
	LongDescription			string
	Sku						string
	Asin					string
	Isbn					string
	OurPrice				string
	RetailPrice				string
	BuyUrl					string
	ThumbNailUrl			string
	LargeImageUrl			string
	AverageCustomerRating	string
	Author 					string
	Publisher				string
	AudioLength				string
	SampleUrl				string
	ReleaseDate				string
	Narrator				string
	Faithfulness			string
	NumberOfCredits			string
}

func readLines(path string, dbmap *gorp.DbMap) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		arr := strings.Split(line, "\t")
		fmt.Printf(arr[0])
		a := newAudible(arr)

		// insert rows - auto increment PKs will be set properly after the insert
		err = dbmap.Insert(&a)
		checkErr(err, "Insert failed")
	}
}

func newAudible(line []string) Audible {
	return Audible{
		Created: time.Now().UnixNano(),
		LastUpdated: 			line[0],
		Name: 					line[1],
		Category: 				line[2],
		Keywords: 				line[3],
		ShortDescription: 		line[4],
		LongDescription: 		line[5],
		Sku: 					line[6],
		Asin: 					line[7],
		Isbn: 					line[8],
		OurPrice: 				line[9],
		RetailPrice: 			line[10],
		BuyUrl: 				line[11],
		ThumbNailUrl: 			line[12],
		LargeImageUrl: 			line[13],
		AverageCustomerRating: 	line[14],
		Author: 				line[15],
		Publisher: 				line[16],
		AudioLength: 			line[17],
		SampleUrl: 				line[18],
		ReleaseDate: 			line[19],
		Narrator: 				line[20],
		Faithfulness: 			line[21],
		NumberOfCredits: 		line[22],
	}
}

func initDb() *gorp.DbMap {
	// connect to db using standard Go database/sql API
	// use whatever database/sql driver you wish
	db, err := sql.Open("postgres", "postgres://audiotime:audiotime@localhost/audiotime?sslmode=disable")
	checkErr(err, "postgres.Open failed")

	// construct a gorp DbMap
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}

	// add a table, setting the table name to 'audibles' and
	// specifying that the Id property is an auto incrementing PK
	dbmap.AddTableWithName(Audible{}, "audibles").SetKeys(true, "Id")

	// create the table. in a production system you'd generally
	// use a migration tool, or create the tables via scripts
	err = dbmap.CreateTablesIfNotExists()
	checkErr(err, "Create tables failed")

	return dbmap
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}