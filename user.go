package main

import (
  "fmt"
  "os"
  "log"
  "github.com/etrubenok/ExampleGolangREST/domain"
  "github.com/gin-gonic/gin"
  "github.com/jackc/pgx"
)

var dbConnection domain.DbConnection

func GetUserHandler(c *gin.Context) {
  id := c.Params.ByName("id")

  user, err := domain.GetUserById(dbConnection, id)

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

  c.JSON(200, user)
}

func CreateUserHandler(c *gin.Context) {
  var userJson struct {
    Name string `json:"name" binding:"required`
  }
  if c.BindJSON(&userJson) == nil {
    c.JSON(200, gin.H{"id": "999", "name": userJson.Name})
  }
}

func SetupRouter() *gin.Engine {
  r := gin.Default()
  r.GET("/user/:id", GetUserHandler)
  r.POST("/user", CreateUserHandler)

  return r
}

func main() {
  var err error
  var pool *pgx.ConnPool
  pool, err = pgx.NewConnPool(extractConfig())

  if err != nil {
    fmt.Fprintln(os.Stderr, "Unable to connect to database:", err)
    os.Exit(1)
  }
  dbConnection = &domain.NewDbConnection{pool}


  router := SetupRouter()
  router.Run(":8080")
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