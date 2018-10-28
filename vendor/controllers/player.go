package controller

import (
	model "model"
	"time"

	"github.com/gin-gonic/gin"
)

func ComparePlayers(c *gin.Context) {

	result := Compare{}
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
		result.Forms = formPlayers(c, player1.ID, player2.ID, dateStart, dateEnd)
		result.Progress = progressPlayers(c, player1.ID, player2.ID, dateStart, dateEnd)

	}

	c.JSON(200, result)
}

// данные по прогрессу игроков за последний год
func progressPlayers(c *gin.Context, player1 int, player2 int, date1 time.Time, date2 time.Time) (results []Progress) {

	var ratings []model.Rating
	model.Connect.
		Where("player = ? or player = ?", player1, player2).
		Where("dateUpdate > ?", date1).
		Where("dateUpdate < ?", date2).
		Find(&ratings)

	for _, rating := range ratings {

		_result := Progress{
			Date: rating.DateUpdate,
		}
		if len(_result.Data) == 0 {
			_result.Data = map[int]int{}
		}
		if rating.Player == player1 {
			_result.Data[0] = rating.Position
		}
		if rating.Player == player2 {
			_result.Data[1] = rating.Position
		}

		results = append(results, _result)
	}

	return
}

type Form struct {
	All  int
	Win  int
	Lose int
}

type FormMounth struct {
	Month     string
	Statistic map[int]Form
}

type Progress struct {
	Date string
	Data map[int]int
}

type Forms map[string]map[int]Form

type Compare struct {
	Forms    []FormMounth
	Progress []Progress
}

// данные по форме игроков за определённый период
func formPlayers(c *gin.Context, player1 int, player2 int, date1 time.Time, date2 time.Time) (results []FormMounth) {

	win := 0
	lose := 0
	all := 0

	for date1.Unix() < date2.Unix() {

		var games1 []model.Game
		var games2 []model.Game
		var month string
		startDate := date1
		date1 = date1.AddDate(0, 1, 0)

		month = date1.Month().String()
		resByMonth := map[int]Form{}

		win = 0
		lose = 0
		all = 0

		model.Connect.
			Where("player1 = ? or player2 = ?", player1, player1).
			Where("dateEvent < ?", date1).
			Where("dateEvent > ?", startDate).
			Find(&games1)
		for _, game := range games1 {

			if game.Player1 == player1 {
				win++
			} else {
				lose++
			}
			all++

			resByMonth[player1] = Form{
				All:  all,
				Win:  win,
				Lose: lose,
			}
		}

		win = 0
		lose = 0
		all = 0

		model.Connect.
			Where("player1 = ? or player2 = ?", player2, player2).
			Where("dateEvent < ?", date1).
			Where("dateEvent > ?", startDate).
			Find(&games2)

		for _, game := range games2 {

			if game.Player1 == player2 {
				win++
			} else {
				lose++
			}
			all++

			resByMonth[player2] = Form{
				All:  all,
				Win:  win,
				Lose: lose,
			}
		}
		result := FormMounth{
			Month:     month,
			Statistic: resByMonth,
		}

		results = append(results, result)

	}
	//str, _ := json.Marshal(results)
	//c.JSON(200, c.BindJSON(results))
	//fmt.Printf("%+v", c.BindJSON(results))
	//c.BindJSON(results)

	return results
	/*fmt.Println(date2.Month().String())
	fmt.Println(date2.Sub(date1).Hours())*/

}
