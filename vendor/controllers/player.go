package controller

import (
	"fmt"
	model "model"
	"time"

	"github.com/gin-gonic/gin"
)

func ComparePlayers(c *gin.Context) {

	model.InitDB()

	code1 := c.Params.ByName("code1")
	code2 := c.Params.ByName("code2")

	var player1 model.Player
	var player2 model.Player
	model.Connect.Where("code = ?", code1).Find(&player1)
	model.Connect.Where("code = ?", code2).Find(&player2)

	if player1.ID > 0 && player2.ID > 0 {

		dateEnd := time.Now()
		dateStart := dateEnd.AddDate(0, -6, 0)
		formPlayers(c, player1.ID, player2.ID, dateStart, dateEnd)

	}

	//c.JSON(200, players)
}

// данные по прогрессу игроков за последний год
func progressPlayers(c *gin.Context, code1 string, code2 string) {

}

// данные по форме игроков за определённый период
func formPlayers(c *gin.Context, player1 int, player2 int, date1 time.Time, date2 time.Time) {

	/*fmt.Println(date2.Format(time.UnRFC822ixDate))
	fmt.Println(date1.Format(time.RFC822))*/
	//then, err := time.Parse("2006-01-02 15:04 MST", "2014-05-03 20:57 UTC")

	fmt.Println(date2.Format("2006-01-02 15:04"))
	/*then, err := time.Parse("2006-02-02 15:04 MST", date2.Format("2006-01-02 15:04"))

	if err != nil {
		fmt.Println(err)
		return
	}*/
	for date1.Unix() < date2.Unix() {

		var games []model.Game

		startDate := date1
		date1 = date1.AddDate(0, 1, 0)

		model.Connect.
			Where("player1 = ? or player2 = ?", player1, player2).
			Where("dateEvent < ?", date1).
			Where("dateEvent > ?", startDate).
			Find(&games)

		fmt.Println(date1.Month().String())

	}
	/*fmt.Println(date2.Month().String())
	fmt.Println(date2.Sub(date1).Hours())*/

}
