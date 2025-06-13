package controllers

import (
	"cinema-management/database"
	"net/http"
	"strconv"

	"cinema-management/repository"

	"cinema-management/structs"

	"github.com/gin-gonic/gin"
)

func GetaLLCinema(c *gin.Context) {
	var (
		result gin.H
	)

	cinema, err := repository.GetAllCinema(database.DbConnection)

	if err != nil {
		result = gin.H{
			"result": err.Error(),
		}
	} else {
		result = gin.H{
			"result": cinema,
		}
	}

	c.JSON(http.StatusOK, result)
}

func InsertCinema(c *gin.Context) {
	var cinema structs.Cinema

	err := c.BindJSON(&cinema)
	if err != nil {
		panic(err)
	}

	err = repository.InsertCinema(database.DbConnection, cinema)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, cinema)
}

func UpdateCinema(c *gin.Context) {
	var cinema structs.Cinema
	id, _ := strconv.Atoi(c.Param("id"))

	err := c.BindJSON(&cinema)
	if err != nil {
		panic(err)
	}

	cinema.ID = id

	err = repository.UpdateCinema(database.DbConnection, cinema)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, cinema)
}

func DeleteCinema(c *gin.Context) {
	var cinema structs.Cinema
	id, _ := strconv.Atoi(c.Param("id"))

	cinema.ID = id
	err := repository.DeleteCinema(database.DbConnection, cinema)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, cinema)
}
