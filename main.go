package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Good struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
}

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:Zxcvbnm123@localhost:5432/gonurlan?sslmode=disable"
	}

	db, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := gin.Default()

	r.Use(cors.Default())

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World")
	})

	r.GET("/api/goods", func(c *gin.Context) {
		rows, err := db.Query(context.Background(), `
			SELECT id, name, price, COALESCE(description, '')
			FROM goods
			ORDER BY id DESC
		`)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		goods := []Good{}

		for rows.Next() {
			var g Good

			err := rows.Scan(
				&g.ID,
				&g.Name,
				&g.Price,
				&g.Description,
			)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			goods = append(goods, g)
		}

		c.JSON(http.StatusOK, goods)
	})

	r.Run(":8080")
}