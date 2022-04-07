package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type jsonData struct {
	Number int
	String string
	Bool   bool
}
type twoNumber struct {
	Right int
	Left  int
}
type City struct {
	ID          int    `json:"id,omitempty"  db:"ID"`
	Name        string `json:"name,omitempty"  db:"Name"`
	CountryCode string `json:"countryCode,omitempty"  db:"CountryCode"`
	District    string `json:"district,omitempty"  db:"District"`
	Population  int    `json:"population,omitempty"  db:"Population"`
}

func main() {
	e := echo.New()

	e.GET("/hello/:username", helloHandler)
	e.GET("/greeting", func(c echo.Context) error {
		log.Print("work")
		return c.String(http.StatusOK, "Hello, database.\n")
	})
	e.GET("/json", jsonHandler)
	e.GET("/database", database)
	e.POST("/post", postHnadler)
	e.POST("/add", addHnadler)
	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong\n")
	})
	e.GET("/fizzbuzz", fizzBuzz)
	e.Logger.Fatal(e.Start(os.Getenv("PORT")))
}
func database(c echo.Context) error {
	db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOSTNAME"), os.Getenv("DB_PORT"), os.Getenv("DB_DATABASE")))
	if err != nil {
		log.Fatalf("Cannot Connect to Database: %s", err)
	}

	fmt.Println("Connected!")
	var city City
	if err := db.Get(&city, "SELECT * FROM city WHERE Name='Tokyo'"); errors.Is(err, sql.ErrNoRows) {
		log.Printf("no such city Name = %s", "Tokyo")
	} else if err != nil {
		log.Fatalf("DB Error: %s", err)
	}

	fmt.Printf("Tokyoの人口は%d人です\n", city.Population)
	return c.String(http.StatusOK, strconv.Itoa(city.Population))
}

func addHnadler(c echo.Context) error {
	data := new(twoNumber)
	err := c.Bind(data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
	}
	ans := data.Left + data.Right
	//return c.JSON(http.StatusOK, data)
	return c.JSON(http.StatusOK, gin.H{"answer": ans})
}
func fizzBuzz(c echo.Context) error {
	countstr := c.QueryParam("count")

	count, err := strconv.Atoi(countstr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "It must a be integer")
	}
	if count%5 == 0 && count%3 == 0 {
		return c.String(http.StatusOK, "FizzBuzz")
	}
	if count%5 == 0 {
		return c.String(http.StatusOK, "Buzz")
	}
	if count%3 == 0 {
		return c.String(http.StatusOK, "Fizz")
	}
	return c.String(http.StatusOK, countstr)
}
func helloHandler(c echo.Context) error {
	userID := c.Param("username")
	return c.String(http.StatusOK, "Hello, "+userID+".\n")
}
func postHnadler(c echo.Context) error {
	data := new(jsonData)
	err := c.Bind(data)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("%+v", data))
	}
	return c.JSON(http.StatusOK, data)
}
func jsonHandler(c echo.Context) error {
	res := jsonData{
		Number: 10,
		String: "hoge",
		Bool:   false,
	}
	return c.JSON(http.StatusOK, &res)
}
