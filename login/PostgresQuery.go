package login

import (
	"database/sql"
	"fmt"
	"net/http"
	_ "github.com/lib/pq"
  )
  
const (
host     = "localhost"
port     = 5431
user     = "switchboard_admin"
password = "12345687"
dbname   = "go_test"
)


func PostgresQuery(username string) string {
	postgresqlDbInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", postgresqlDbInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Established a successful connection with database!")

	var password string
	row := db.QueryRow("SELECT password FROM person where username=$1", username)
	if err := row.Scan(&password); err != nil{
		fmt.Println("Error with querying data: Didn't find username")
	}
	// else{
	// 	fmt.Println("Found username")
	// }
	return password

}