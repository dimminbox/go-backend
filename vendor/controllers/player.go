package controller

import (
	"fmt"
	model "model"
	"time"

	"github.com/gin-gonic/gin"
)

type Surface struct {
	Surface     string
	Won1        int
	Won2        int
	All1        int
	All2        int
	AllSurface1 int
	AllSurface2 int
	Use1        int
	Use2        int
	Title1      int
	Title2      int
}

func surfacePlayers(c *gin.Context, player1 int, player2 int) (surfaces []Surface) {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()

	var tournaments []model.Tournament
	surfaceToTournir := map[string][]int{}

	model.Connect.Select("id, surface").Find(&tournaments)
	for _, tournir := range tournaments {
		if tournir.Surface != "" {
			if _, ok := surfaceToTournir[tournir.Surface]; ok {
				surfaceToTournir[tournir.Surface] = append(surfaceToTournir[tournir.Surface], tournir.ID)
			} else {
				surfaceToTournir[tournir.Surface] = []int{tournir.ID}
			}
		}
	}

	for surfaceName, tournirs := range surfaceToTournir {

		var all1 int
		var all2 int
		var count1Win int
		var count2Win int
		var all1Surface int
		var all2Surface int
		var title1 int
		var title2 int

		var games []model.Game

		model.Connect.Where("player1 = ? or player2 = ?", player1, player1).
			Find(&games).
			Count(&all1)

		model.Connect.Where("player1 = ? or player2 = ?", player2, player2).
			Find(&games).
			Count(&all2)

		model.Connect.Where("player1 = ?", player1).
			Where("tournir IN (?)", tournirs).
			Find(&games).
			Count(&count1Win)

		model.Connect.Where("player1 = ?", player2).
			Where("tournir IN (?)", tournirs).
			Find(&games).
			Count(&count2Win)

		model.Connect.Where("player1 = ? or player2 = ?", player1, player1).
			Where("tournir IN (?)", tournirs).
			Find(&games).
			Count(&all1Surface)

		model.Connect.Where("player1 = ? or player2 = ?", player2, player2).
			Where("tournir IN (?)", tournirs).
			Find(&games).
			Count(&all2Surface)

		model.Connect.Where("player1 = ?", player1).
			Where("stage = 'Finals'").
			Find(&games).
			Count(&title1)

		model.Connect.Where("player1 = ?", player2).
			Where("stage = 'Finals'").
			Find(&games).
			Count(&title2)

		surface := Surface{
			Surface:     surfaceName,
			All1:        all1,
			All2:        all2,
			AllSurface1: all1Surface,
			AllSurface2: all2Surface,
			Won1:        count1Win,
			Won2:        count2Win,
			Title1:      title1,
			Title2:      title2,
		}
		surfaces = append(surfaces, surface)

	}

	return surfaces

}

type Compare struct {
	Players  []model.Player
	Forms    []FormMounth
	Progress []Progress
	Surface  []Surface
}

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

		result.Players = append(result.Players, player1)
		result.Players = append(result.Players, player2)
		dateEnd := time.Now()
		dateStart := dateEnd.AddDate(0, -6, 0)
		result.Forms = formPlayers(c, player1.ID, player2.ID, dateStart, dateEnd)
		result.Progress = progressPlayers(c, player1.ID, player2.ID, dateStart, dateEnd)
		result.Surface = surfacePlayers(c, player1.ID, player2.ID)

	}

	c.JSON(200, result)
}

type Progress struct {
	Date string
	Data map[int]int
}

// данные по прогрессу игроков за последний год
func progressPlayers(c *gin.Context, player1 int, player2 int, date1 time.Time, date2 time.Time) (results []Progress) {

	var ratings []model.Rating

	model.Connect.
		Where("player = ? or player = ?", player1, player2).
		Where("dateUpdate > ?", date1).
		Where("dateUpdate < ?", date2).
		Order("dateUpdate ASC").
		Find(&ratings)

	var dateUpdate string
	var _result Progress
	for _, rating := range ratings {

		if dateUpdate != rating.DateUpdate {

			if _result.Date != "" {
				results = append(results, _result)
			}
			dateUpdate = rating.DateUpdate
			_result = Progress{
				Date: dateUpdate,
				Data: map[int]int{},
			}

		}

		if _result.Date != "" {
			if rating.Player == player1 {
				_result.Data[0] = rating.Position
			}
			if rating.Player == player2 {
				_result.Data[1] = rating.Position
			}
		}
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
type Forms map[string]map[int]Form

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

	return results
	/*fmt.Println(date2.Month().String())
	fmt.Println(date2.Sub(date1).Hours())*/

}
