package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
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

func addPercentage(c *gin.Context) {
	body, eror := ioutil.ReadAll(c.Request.Body)
	if eror != nil {
		fmt.Println("Error in reading from request", eror)
		c.JSON(500, gin.H{"Error": "Failed to read"})
		return
	}
	fmt.Printf("The recieved request is of type %T", body)
	fmt.Println("The recieved request is ", string(body))

	// insert into db
	if _, err := db.Exec("INSERT INTO percentage(id,percent) VALUES (DEFAULT,$1)", string(body)); err != nil {
		fmt.Println("Error in inserting into db", err)
		c.JSON(500, gin.H{"Error": "Failed to insert"})
		return
	}

	c.JSON(200, gin.H{"Push": "Success"})
	return

}
func getPercentage(c *gin.Context) {
	var percent int
	row := db.QueryRow("select percent from percentage order by id desc limit 1")
	errr := row.Scan(&percent)
	if errr == sql.ErrNoRows {
		fmt.Println("No rows were returned!", errr)
		c.JSON(500, gin.H{"Error": "Failed to read"})
		return
	} else if errr != nil {
		fmt.Println("Error in selecting from db", errr)
		c.JSON(500, gin.H{"Error": "Failed to read"})
		return
	}
	fmt.Println("The percentage value is ", percent)
	c.JSON(200, gin.H{"Percentage": percent})
	return
}
func getPercentvalue() (string, error) {
	fmt.Println("---Running in getPercentvalue---")
	var percent int
	row := db.QueryRow("select percent from percentage order by id desc limit 1")
	errr := row.Scan(&percent)
	if errr == sql.ErrNoRows {
		fmt.Println("No rows were returned!", errr)
		return "", errr
	} else if errr != nil {
		fmt.Println("Error in selecting from db", errr)
		return "", errr
	}
	fmt.Println("The percentage value is ", percent)
	return string(percent), nil
}

/*
// For getting realtime distance value
func getPercentage(c *gin.Context) {
	fmt.Println("--Running in getPercentage---")
	// in Mac the path to py3 is /Library/Frameworks/Python.framework/Versions/3.6/bin/python3
	// in linux /usr/bin/python3
	cmd := exec.Command("/Library/Frameworks/Python.framework/Versions/3.6/bin/python3", "/Users/souvikhaldar/Development/go/src/github.com/heroku/go-getting-started/mrbin.py")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error in Output", err.Error())
		c.JSON(500, gin.H{"Error": string(err.Error())})
		return
	}
	fmt.Println("The output of the script is", string(out))
	fmt.Printf("The type is %T", out)
	c.JSON(200, gin.H{"Distance": string(out)})
}
*/

// this function is for putting the data from arduino to the server

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
	router.LoadHTMLGlob("templates/*.html")
	router.Static("/static", "static")
	//router.LoadHTMLFiles("templates/aboutus.html")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	var percent = string(53)
	percent, err = getPercentvalue()
	if err != nil {
		fmt.Println("Some issue is getting value from db", err)
	}

	router.GET("/rts.html", func(c *gin.Context) {
		fmt.Println("Final value of percent is", percent)
		c.HTML(http.StatusOK, "rts.html", gin.H{"percent": 34})
	})
	router.GET("/aboutus.html", func(c *gin.Context) {
		c.HTML(http.StatusOK, "aboutus.html", nil)
	})
	router.GET("/route.html", func(c *gin.Context) {
		c.HTML(http.StatusOK, "route.html", nil)
	})

	router.GET("/repeat", repeatFunc)
	router.GET("/db", dbFunc)
	router.GET("/getPercent", getPercentage)
	router.POST("/addPercent", addPercentage)

	router.Run(":" + port)
}
