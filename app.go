package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// App ...
type App struct {
	Router *gin.Engine
	DB     *sql.DB
}

func (a *App) getEndpoints(c *gin.Context) {
	endpoints, err := getEndpoints(a.DB, 0, 10)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, endpoints)
}

func (a *App) updateEndpoint(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid endpoint ID")
		return
	}

	var e endpoint
	if c.ShouldBind(&e) != nil {
		c.String(http.StatusBadRequest, "Invalid request body")
	}
	e.ID = id

	err = e.updateEndpoint(a.DB)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, e)
}

func (a *App) getEndpoint(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid endpoint ID")
		return
	}

	e := endpoint{ID: id}

	if err := e.getEndpoint(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			c.JSON(http.StatusNotFound, gin.H{"error": "Endpoint not found"})
		default:
			c.JSON(http.StatusInternalServerError, err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, e)
}

func (a *App) createEndpoint(c *gin.Context) {
	var e endpoint

	if c.ShouldBind(&e) != nil {
		c.String(http.StatusBadRequest, "Invalid request body")
		return
	}

	if !e.isValidURL() {
		c.String(http.StatusBadRequest, "Invalid URL")
		return
	}

	err := e.createEndpoint(a.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, e)
}

func (a *App) deleteEndpoint(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid endpoint ID")
		return
	}

	e := endpoint{ID: id}

	if err := e.deleteEndpoint(a.DB); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": "success",
	})
}

func (a *App) getClientContent(c *gin.Context) {
	url := c.Param("url")

	e := endpoint{URL: url}

	if err := e.getEndpointByURL(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			c.JSON(http.StatusNotFound, gin.H{"error": "Endpoint not found"})
		default:
			c.JSON(http.StatusInternalServerError, err.Error())
		}
		return
	}

	var result map[string]interface{}
	json.Unmarshal([]byte(e.Content), &result)

	c.JSON(http.StatusOK, result)
}

func (a *App) initializeRoutes() {
	a.Router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{os.Getenv("CLIENT_URL")},
		AllowCredentials: true,
		AllowMethods:     []string{"PUT", "PATCH", "GET", "DELETE"},
		AllowHeaders:     []string{"Origin", "content-type"},
	}))

	a.Router.GET("/endpoints", a.getEndpoints)
	a.Router.POST("/endpoint", a.createEndpoint)
	a.Router.GET("/endpoint/:id", a.getEndpoint)
	a.Router.PUT("/endpoint/:id", a.updateEndpoint)
	a.Router.DELETE("/endpoint/:id", a.deleteEndpoint)
	a.Router.GET("/temp/:url", a.getClientContent)
}

// Initialize ...
func (a *App) Initialize(connectionString string) {
	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = gin.Default()

	a.initializeRoutes()
}

// Run ...
func (a *App) Run() {
	a.Router.Run()
}
