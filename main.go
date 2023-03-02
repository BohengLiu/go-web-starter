package main

import (
	"context"
	"database/sql"
	"github.com/redis/go-redis/v9"
	"github.com/uptrace/bun"
	"net/http"

	"github.com/BohengLiu/go-web-starter/model"
	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func main() {
	dsn := "postgres://postgres:@localhost:5432/exampledb?sslmode=disable"
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	db := bun.NewDB(sqldb, pgdialect.New())

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.GET("/ping", func(c *gin.Context) {
		var data []model.User
		db.NewSelect().Model(&data).Column("id", "name", "signup").Scan(context.Background())
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
			"data":    data,
		})
	})
	r.GET("/set", func(c *gin.Context) {
		key := c.DefaultQuery("key", "")
		val := c.DefaultQuery("val", "")
		if key == "" || val == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "not key or val",
			})
			return
		}
		err := rdb.Set(context.Background(), key, val, 0).Err()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "internal error",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "OK",
		})
	})
	r.GET("/get", func(c *gin.Context) {
		key := c.DefaultQuery("key", "")
		if key == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "not key or val",
			})
			return
		}
		val, err := rdb.Get(context.Background(), key).Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "internal error",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "OK",
			"data":    val,
		})
	})
	r.GET("/json-slice", func(c *gin.Context) {
		var val []string
		val = []string{}
		c.JSON(http.StatusOK, gin.H{
			"message": "OK",
			"data":    val,
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
