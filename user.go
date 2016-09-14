package main

import (
  "fmt"
  "os"
  "log"
  "github.com/gin-gonic/gin"
  "github.com/jackc/pgx"
)

var pool *pgx.ConnPool

func getUserHandler(c *gin.Context) {
  id := c.Params.ByName("id")

  var name string
  err := pool.QueryRow("select name from users where id=$1", id).Scan(&name)
  if err != nil {

    switch err {
    default:
      log.Panic(err)
      c.JSON(500, gin.H{"message": "internal server error"})
    case pgx.ErrNoRows:
      c.JSON(404, gin.H{
        "id": id,
        "message": "user is not found"})
    }
    return
  }
  c.JSON(200, gin.H{"id": id, "name": name})
}

func createUserHandler(c *gin.Context) {
  var userJson struct {
    Name string `json:"name" binding:"required`
  }
  if c.BindJSON(&userJson) == nil {
    c.JSON(200, gin.H{"id": "999", "name": userJson.Name})
  }
}

func main() {
  var err error
  pool, err = pgx.NewConnPool(extractConfig())
  if err != nil {
    fmt.Fprintln(os.Stderr, "Unable to connect to database:", err)
    os.Exit(1)
  }

  r := gin.Default()

  r.GET("/user/:id", getUserHandler)
  r.POST("/user", createUserHandler)
  r.Run(":8080")
}

func extractConfig() pgx.ConnPoolConfig {
  var config pgx.ConnPoolConfig

  config.Host = os.Getenv("DB_HOST")
  if config.Host == "" {
    config.Host = "localhost"
  }

  config.User = os.Getenv("DB_USER")
  if config.User == "" {
    config.User = os.Getenv("USER")
  }

  config.Password = os.Getenv("DB_PASSWORD")

  config.Database = os.Getenv("DB_DATABASE")
  if config.Database == "" {
    config.Database = "postgres"
  }

  return config
}