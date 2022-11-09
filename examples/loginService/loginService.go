package main
 
import (
	"database/sql"
	"fmt"
	"net/http"
	_ "github.com/lib/pq"
	"github.com/gin-gonic/gin"
  )
  
const (
host     = "localhost"
port     = 5431
user     = "switchboard_admin"
password = "12345687"
dbname   = "go_test"
)


func PostgresQuery(username string) string { //ctx context.Context
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
	}else{
		fmt.Println("Found username")
	}
	// fmt.Println("found password as: ", password)
	return password

}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}


func main(){
	
	router := gin.Default()
	router.POST("/loginJSON", func(c *gin.Context) {
		var json User
		if err := c.ShouldBindJSON(&json); err == nil {
			fmt.Println("json receive - %+v", json.Username)
			
		}else{
			fmt.Println("error - %+v", err)
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"data": json,
		})

		fmt.Println("get username as: ", json.Username)

		foundPassword := PostgresQuery(json.Username)
		if foundPassword == json.Password{
			fmt.Println("password is the same; login success")
		}else{
			fmt.Println("login failed")
		}
	})
	router.Run("localhost:8081")
}