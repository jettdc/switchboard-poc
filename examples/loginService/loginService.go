package main
 
import (
	"database/sql"
	"fmt"
	"net/http"
	_ "github.com/lib/pq"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"os"
	"strconv"
  )


func InitializeEnv() error {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return err
	}
	return nil
}


func PostgresQuery(username string, userpassword string, host string, port string, user string, password string, dbname string) (string, error) {
	intport, err := strconv.Atoi(port)
	postgresqlDbInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, intport, user, password, dbname)

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

	var token string
	row := db.QueryRow("SELECT token FROM person where username=$1 and password=$2", username, userpassword)
	if err := row.Scan(&token); err != nil{
		fmt.Println("error", err)
		return "", err
	}else{
		return token, err
	}
	
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}


func main(){
	
	InitializeEnv()
	router := gin.Default()
	router.POST("/loginJSON", func(c *gin.Context) {
		var json User
		if err := c.ShouldBindJSON(&json); err == nil {
			fmt.Println("error - %+v", err)
		}
		// else{
		// 	fmt.Println("error - %+v", err)
		// }

		fmt.Println("get username as: ", json.Username)

		token, err := PostgresQuery(json.Username, json.Password, os.Getenv("host"), os.Getenv("port"), os.Getenv("user"), os.Getenv("password"), os.Getenv("dbname"))
		if err != nil{
			c.JSON(http.StatusUnauthorized, "")
		} else{
			c.JSON(http.StatusOK, gin.H{
				"status": "ok",
				"token": token,
			})
		}
	})
	router.Run("localhost:"+os.Getenv("localhostport")) //question: modify this as os.getenv too?
}