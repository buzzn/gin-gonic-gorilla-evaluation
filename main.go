package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Error to be reported to the user
type Error struct {
	Name       string `json:"errorName"`
	Desciption string `json:"errorDescription"`
}

// Hint to be shown to the user.
type Hint struct {
	Icon        string `json:"icon"`
	Headline    string `json:"headline"`
	Description string `json:"description"`
}

// Challange can be masterd by the user to prove he is capable of saving energy.
type Challange struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Total       int    `json:"total"`
	Succeeded   int    `json:"succeeded"`
}

// UserProfile with basic informationa about the current user.
type UserProfile struct {
	Nick           string  `form:"nick" json:"nick"`
	FlatSize       float64 `form:"flatsize" json:"faltSize"`
	FlatPopulation int     `form:"flatpopulation" json:"flatPopulation"`
}

// HitlistEntry represents one entry int the hitlist.
type HitlistEntry struct {
	Nick     string `json:"nick"`
	Position int    `json:"position"`
	Value    string `json:"value"`
}

// LiveChartEntry represents an entry in the production/consumption chart.
type LiveChartEntry struct {
	Date            time.Time `json:"date"`
	UserConsumption int       `json:"userConsumption"`
	GroupConsumtion int       `json:"groupConsumption"`
	GroupProduction int       `json:"groupProduction"`
	SelfSufficiency int       `json:"selfSufficiency"`
}

// MaxValuesHistory tells how many history values can be requested at once
const MaxValuesHistory = 10000

// PasswordMinLength specifies the minimum length of a user's password
const PasswordMinLength = 8

func main() {
	var challangeRunning Challange
	currentUser := UserProfile{"DarkNight", 120.0, 2}

	hints := [3]Hint{
		Hint{"HAIR_DRYER", "Don't dry your hair hot!", "Heating the air consumes a lot of energy. Besides, it is way better for your hair to dry it cold, belive me ;)"},
		Hint{"PLUG", "Turn of your plugs.", "Devices in stand-by still consume energy. Turning the plug off will save you a log of energy."},
		Hint{"CLOTHES_DRYER", "Dont't use a clothes dryer.", "These devices consume the most energy in your household. Besides your trousers will life longer when they dry on fresh air."}}

	challanges := [3]Challange{
		Challange{0, "Quiet thorught the nights", "Can you consume less than 1kwh 5 nights in a row?", 0, 0},
		Challange{1, "Take it easy", "Can you consume less than 20 kwh for one week?", 0, 0},
		Challange{2, "Beat youre self", "Consume less than in tha last week.", 0, 0}}

	hitlist := [6]HitlistEntry{
		HitlistEntry{"Lisl", 1, "1000kW"},
		HitlistEntry{"Dieter", 2, "1300kW"},
		HitlistEntry{"Heinz", 3, "1400kW"},
		HitlistEntry{"Gunigunde", 4, "1700kW"},
		HitlistEntry{"Sylvia Maria Eva", 5, "2000kW"},
		HitlistEntry{"Y", 6, "2400kW"}}

	route := gin.Default()

	route.GET("/hints", func(c *gin.Context) {
		c.JSON(200, hints)
	})

	route.GET("/challanges", func(c *gin.Context) {
		c.JSON(200, challanges)
	})

	route.GET("/challanges/start/:id", func(c *gin.Context) {
		idString := c.Param("id")

		id, error := strconv.Atoi(idString)
		if error != nil {
			c.JSON(200, Error{"Unknown id", idString + " is not a known challange."})
			return
		}

		if id >= len(challanges) {
			c.JSON(200, Error{"Unknown id", idString + " is not a known challange."})
			return
		}

		challangeRunning = challanges[id]
		c.JSON(200, challangeRunning)
	})

	route.GET("/challanges/status", func(c *gin.Context) {
		c.JSON(200, challangeRunning)
	})

	route.GET("/individual-consumption-history/begin/:begin/end/:end/tics/:tics", func(c *gin.Context) {
		beginDateString := c.Param("begin")
		endDateString := c.Param("end")
		ticsString := c.Param("tics")

		beginDate, error := time.Parse(time.RFC3339, beginDateString)
		if error != nil {
			c.JSON(200, Error{"Invalid begin date", beginDateString + " is not a valid RFC3339 date "})
			return
		}

		endDate, error := time.Parse(time.RFC3339, endDateString)
		if error != nil {
			c.JSON(200, Error{"Invalid end date", endDateString + " is not a valid RFC3339 date "})
			return
		}

		tics, error := strconv.Atoi(ticsString)
		if error != nil {
			c.JSON(200, Error{"Invalid values for tics ", ticsString + " is not a valid integer."})
			return
		}

		if beginDate.After(endDate) {
			c.JSON(200, Error{"Invalid values for dates ", "End date must be after begin date."})
			return
		}

		totalValues := endDate.Sub(beginDate).Seconds() / float64(tics)
		if totalValues > MaxValuesHistory {
			c.JSON(200, Error{"Too many values requested ", "Max tics allowed: " + string(MaxValuesHistory)})
			return
		}

		results := make([]float64, int(totalValues))
		last := float64(0)
		for i := 0; i < int(totalValues); i++ {
			last += rand.Float64() * 10
			results[i] = last
		}

		c.JSON(200, results)
	})

	route.GET("/group-consumption-history/begin/:begin/end/:end/tics/:tics", func(c *gin.Context) {
		beginDateString := c.Param("begin")
		endDateString := c.Param("end")
		ticsString := c.Param("tics")

		beginDate, error := time.Parse(time.RFC3339, beginDateString)
		if error != nil {
			c.JSON(200, Error{"Invalid begin date", beginDateString + " is not a valid RFC3339 date "})
			return
		}

		endDate, error := time.Parse(time.RFC3339, endDateString)
		if error != nil {
			c.JSON(200, Error{"Invalid end date", endDateString + " is not a valid RFC3339 date "})
			return
		}

		tics, error := strconv.Atoi(ticsString)
		if error != nil {
			c.JSON(200, Error{"Invalid values for tics ", ticsString + " is not a valid integer."})
			return
		}

		if beginDate.After(endDate) {
			c.JSON(200, Error{"Invalid values for dates ", "End date must be after begin date."})
			return
		}

		totalValues := endDate.Sub(beginDate).Seconds() / float64(tics)
		if totalValues > MaxValuesHistory {
			c.JSON(200, Error{"Too many values requested ", "Max tics allowed: " + string(MaxValuesHistory)})
			return
		}

		results := make([]float64, int(totalValues))
		last := float64(0)
		for i := 0; i < int(totalValues); i++ {
			last += rand.Float64() * 100
			results[i] = last
		}

		c.JSON(200, results)
	})

	route.GET("/profile", func(c *gin.Context) {
		c.JSON(200, currentUser)
	})

	route.POST("/profile", func(c *gin.Context) {
		if err := c.ShouldBindJSON(&currentUser); err != nil {
			c.JSON(400, Error{"Can not parse profile.", "Maybe there are missing values?"})
			return
		}

		c.JSON(200, currentUser)
	})

	route.GET("/hitlist", func(c *gin.Context) {
		c.JSON(200, hitlist)
	})

	route.POST("/password", func(c *gin.Context) {
		newPassword := c.PostForm("password")
		token := c.PostForm("token")

		if token != "expected" {
			c.JSON(400, Error{"Unknown token", "Try again with a valid token."})
			return
		}

		if len(newPassword) < PasswordMinLength {
			c.JSON(400, Error{"Password too short", "Try again with a password at least  " + strconv.Itoa(PasswordMinLength) + " chars long."})
			return
		}

		c.JSON(200, gin.H{})
	})

	route.POST("/update-password", func(c *gin.Context) {
		newPassword := c.PostForm("password")

		if len(newPassword) < PasswordMinLength {
			c.JSON(400, Error{"Password too short", "Try again with a password at least  " + strconv.Itoa(PasswordMinLength) + " chars long."})
			return
		}

		c.JSON(200, gin.H{})
	})

	route.GET("/live", func(c *gin.Context) {
		var wsupgrader = websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}

		conn, err := wsupgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			fmt.Println("Failed to create websocket")
			return
		}

		e := LiveChartEntry{time.Now(), rand.Int() % 10, rand.Int() % 100, rand.Int() % 10, rand.Int() % 100}
		for {
			e = LiveChartEntry{time.Now(), e.UserConsumption + rand.Int()%10, e.GroupConsumtion + rand.Int()%100, e.GroupProduction + rand.Int()%100, rand.Int() % 100}
			conn.WriteJSON(e)
			time.Sleep(1000000000)
		}
	})

	route.Run(":8088")
}
