package controller

import (
	"fmt"
	"math"
	model "model"
	"strings"
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
	Players    []model.Player
	Forms      []FormMounth
	Progress   []Progress
	Surface    []Surface
	HeadToHead []HeadToHead
	Technic    Technic
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
		result.HeadToHead = headToHead(c, player1, player2)

		dateEnd = time.Now()
		dateStart = dateEnd.AddDate(0, -1, 0)
		result.Technic = technicPlayers(c, player1.ID, player2.ID, dateStart, dateEnd)

	}

	c.JSON(200, result)
}

type Score struct {
	score1 int
	score2 int
}
type HeadToHead struct {
	Date         string
	Tournament   string
	Stage        string
	Score        []string
	Winner       string
	TournamentID int
}

func headToHead(c *gin.Context, player1 model.Player, player2 model.Player) (results []HeadToHead) {

	tournamenentIDs := []int{}
	var games1 []model.Game

	model.Connect.
		Where("player1 = ?", player1.ID).
		Where("player2 = ?", player2.ID).
		Find(&games1)

	for _, game := range games1 {

		result := HeadToHead{
			Date:         game.DateEvent,
			Stage:        game.Stage,
			Score:        strings.Split(game.Scores, ";"),
			Winner:       player1.Name,
			TournamentID: game.Tournir,
		}

		tournamenentIDs = append(tournamenentIDs, game.Tournir)
		results = append(results, result)
	}

	var games2 []model.Game

	model.Connect.
		Where("player2 = ?", player1.ID).
		Where("player1 = ?", player2.ID).
		Find(&games2)

	for _, game := range games2 {

		result := HeadToHead{
			Date:         game.DateEvent,
			Stage:        game.Stage,
			Score:        strings.Split(game.Scores, ";"),
			Winner:       player2.Name,
			TournamentID: game.Tournir,
		}

		tournamenentIDs = append(tournamenentIDs, game.Tournir)
		results = append(results, result)
	}

	var tournirs []model.Tournament
	model.Connect.
		Where("ID IN (?)", tournamenentIDs).
		Find(&tournirs)

	for i, result := range results {
		for _, tournir := range tournirs {
			if tournir.ID == result.TournamentID {
				results[i].Tournament = fmt.Sprintf("%s %s %s", tournir.Name, tournir.Type, tournir.Surface)
			}
		}
	}
	return
}

type Technic struct {
	QualityServe1 int /*разница между ейсами и двойными ошибками*/
	QualityServe2 int
	Serve1        float64 /*процент первой подачи*/
	Serve2        float64
	AvgServeWon1  float64 /*процент выигранных мячей на подаче*/
	AvgServeWon2  float64
	AvgReturnWon1 float64 /*процент выигранных мячей на приёме */
	AvgReturnWon2 float64
}

func technicPlayers(c *gin.Context, player1 int, player2 int, date1 time.Time, date2 time.Time) (result Technic) {

	result1 := getTechnicPlayer(player1, date1, date2)
	result2 := getTechnicPlayer(player2, date1, date2)

	result = Technic{
		Serve1:        result1.Serve1,
		Serve2:        result2.Serve1,
		QualityServe1: result1.QualityServe1,
		QualityServe2: result2.QualityServe1,
		AvgServeWon1:  result1.AvgReturnWon1,
		AvgServeWon2:  result2.AvgReturnWon1,
		AvgReturnWon1: result1.AvgReturnWon1,
		AvgReturnWon2: result2.AvgReturnWon1,
	}

	return

}

func getTechnicPlayer(player int, date1 time.Time, date2 time.Time) (result Technic) {

	var games []model.Game
	model.Connect.
		Where("player1 = ? or player1 = ?", player, player).
		Where("dateEvent < ?", date1).
		Where("dateEvent > ?", date2).
		Find(&games)

	serve := []int{}
	avgServe := []int{}
	avgReturn := []int{}
	for _, game := range games {

		if game.Player1 == player {
			result.QualityServe1 += game.Aces1 - game.DoubleFaults1
			serve = append(serve, game.Serve1)
			avgServe = append(avgServe, game.Serve1PointsWon1)
			avgReturn = append(avgReturn, game.Serve1ReturnPointsWon1)
		}

		if game.Player2 == player {
			result.QualityServe1 += game.Aces2 - game.DoubleFaults2
			serve = append(serve, game.Serve2)
			avgServe = append(avgServe, game.Serve1PointsWon2)
			avgReturn = append(avgReturn, game.Serve1ReturnPointsWon2)
		}

	}

	for _, _serve := range serve {
		result.Serve1 = result.Serve1 + float64(_serve)
	}
	if len(serve) > 0 {
		result.Serve1 = math.Round(result.Serve1/float64(len(serve))) * 100
	}

	for _, _avgServe := range avgServe {
		result.Serve1 = result.Serve1 + float64(_avgServe)
	}
	if len(avgServe) > 0 {
		result.AvgServeWon1 = math.Round(result.AvgServeWon1/float64(len(avgServe))) * 100
	}

	for _, _avgReturn := range avgReturn {
		result.AvgReturnWon1 = result.AvgReturnWon1 + float64(_avgReturn)
	}
	if len(avgReturn) > 0 {
		result.AvgReturnWon1 = math.Round(result.AvgReturnWon1/float64(len(avgReturn))) * 100
	}

	return
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
