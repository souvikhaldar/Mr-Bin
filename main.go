package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type DistanceStruct struct {
	Distance float64
}

var (
	repeat int
	db     *sql.DB
)

func repeatFunc(c *gin.Context) {
	var buffer bytes.Buffer
	for i := 0; i < repeat; i++ {
		buffer.WriteString("Hello from Go!")
	}
	c.String(http.StatusOK, buffer.String())
}

// this function is for putting the data from arduino to the server
func addPercentage(c *gin.Context) {
	height := 45.0
	data, er := ioutil.ReadAll(c.Request.Body)
	if er != nil {
		fmt.Println("Error in reading from request body", er)
		c.JSON(500, gin.H{"Error": "Reading"})
		return
	}
	var dist DistanceStruct
	oops := json.Unmarshal(data, &dist)
	if oops != nil {
		fmt.Println("Error in unmarshalling", oops)
		c.JSON(500, gin.H{"error": "unmarshalling failed"})
		return
	}
	fmt.Println("The recieved distance is ", dist.Distance)
	percentvalue := (dist.Distance / height) * 100
	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS percentage (id bigserial,percent int)"); err != nil {
		fmt.Println("Error in creating percentage table", err)
		c.JSON(500, gin.H{"error": "table creation failed"})
		return
	}
	if _, err := db.Exec("INSERT INTO percentage(percent) VALUES ($1)", int(percentvalue)); err != nil {
		fmt.Println("Error in inserting percentage table", err)
		c.JSON(500, gin.H{"error": "table insertion failed"})
		return
	}
	c.JSON(200, gin.H{"status": "insert success"})
	return
}

func dbFunc(c *gin.Context) {
	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS ticks (tick timestamp)"); err != nil {
		c.String(http.StatusInternalServerError,
			fmt.Sprintf("Error creating database table: %q", err))
		return
	}

	if _, err := db.Exec("INSERT INTO ticks VALUES (now())"); err != nil {
		c.String(http.StatusInternalServerError,
			fmt.Sprintf("Error incrementing tick: %q", err))
		return
	}

	rows, err := db.Query("SELECT tick FROM ticks")
	if err != nil {
		c.String(http.StatusInternalServerError,
			fmt.Sprintf("Error reading ticks: %q", err))
		return
	}

	defer rows.Close()
	for rows.Next() {
		var tick time.Time
		if err := rows.Scan(&tick); err != nil {
			c.String(http.StatusInternalServerError,
				fmt.Sprintf("Error scanning ticks: %q", err))
			return
		}
		c.String(http.StatusOK, fmt.Sprintf("Read from DB: %s\n", tick.String()))
	}
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	var err error
	//tStr := os.Getenv("REPEAT")
	//repeat, err = strconv.Atoi(tStr)
	//if err != nil {
	//	log.Print("Error converting $REPEAT to an int: %q - Using default", err)
	//	repeat = 5
	//}

	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	router.GET("/repeat", repeatFunc)
	router.GET("/db", dbFunc)
	router.POST("/post", addPercentage)

	router.Run(":" + port)
}
