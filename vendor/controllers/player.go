package controller

import (
	"fmt"
	model "model"

	"github.com/gin-gonic/gin"
)

func GetPlayers(c *gin.Context) {

	model.InitDB()
	var players []model.Player
	model.Connect.Find(&players)
	fmt.Println(players)

	c.JSON(200, players)
}
