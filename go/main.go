package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var db *sql.DB

type Product struct {
	Id    int64
	Name  string
	Price int
}

func main() {
	// Load in the `.env` file
	err := godotenv.Load()
	if err != nil {
		log.Print("failed to load env", err)
	}

	// Open a connection to the database
	db, err = sql.Open("mysql", os.Getenv("DSN"))
	if err != nil {
		log.Print("failed to open db connection", err)
	}

	// Build router & define routes
	router := gin.Default()
	router.Use(cors.Default()) // Add the cors middleware
	router.GET("/products", GetProducts)
	router.GET("/products/:productId", GetSingleProduct)
	router.POST("/products", CreateProduct)
	router.PUT("/products/:productId", UpdateProduct)
	router.DELETE("/products/:productId", DeleteProduct)

	// Run the router
	router.Run()
}

func GetProducts(c *gin.Context) {
	query := "SELECT * FROM products"
	res, err := db.Query(query)
	defer res.Close()
	if err != nil {
		log.Print("(GetProducts) db.Query", err)
	}

	products := []Product{}
	for res.Next() {
		var product Product
		err := res.Scan(&product.Id, &product.Name, &product.Price)
		if err != nil {
			log.Print("(GetProducts) res.Scan", err)
		}
		products = append(products, product)
	}

	c.JSON(http.StatusOK, products)
}

func GetSingleProduct(c *gin.Context) {
	productId := c.Param("productId")
	productId = strings.ReplaceAll(productId, "/", "")
	productIdInt, err := strconv.Atoi(productId)
	if err != nil {
		log.Print("(GetSingleProduct) strconv.Atoi", err)
	}

	var product Product
	query := `SELECT * FROM products WHERE id = ?`
	err = db.QueryRow(query, productIdInt).Scan(&product.Id, &product.Name, &product.Price)
	if err != nil {
		log.Print("(GetSingleProduct) db.Exec", err)
	}

	c.JSON(http.StatusOK, product)
}

func CreateProduct(c *gin.Context) {
	var newProduct Product
	err := c.BindJSON(&newProduct)
	if err != nil {
		log.Print("(CreateProduct) c.BindJSON", err)
	}

	query := `INSERT INTO products (name, price) VALUES (?, ?)`
	res, err := db.Exec(query, newProduct.Name, newProduct.Price)
	if err != nil {
		log.Print("(CreateProduct) db.Exec", err)
	}
	newProduct.Id, err = res.LastInsertId()
	if err != nil {
		log.Print("(CreateProduct) res.LastInsertId", err)
	}

	c.JSON(http.StatusOK, newProduct)
}

func UpdateProduct(c *gin.Context) {
	var updates Product
	err := c.BindJSON(&updates)
	if err != nil {
		log.Print("(UpdateProduct) c.BindJSON", err)
	}

	productId := c.Param("productId")
	productId = strings.ReplaceAll(productId, "/", "")
	productIdInt, err := strconv.Atoi(productId)
	if err != nil {
		log.Print("(UpdateProduct) strconv.Atoi", err)
	}

	query := `UPDATE products SET name = ?, price = ? WHERE id = ?`
	_, err = db.Exec(query, updates.Name, updates.Price, productIdInt)
	if err != nil {
		log.Print("(UpdateProduct) db.Exec", err)
	}

	c.Status(http.StatusOK)
}

func DeleteProduct(c *gin.Context) {
	productId := c.Param("productId")

	productId = strings.ReplaceAll(productId, "/", "")
	productIdInt, err := strconv.Atoi(productId)
	if err != nil {
		log.Print("(DeleteProduct) strconv.Atoi", err)
	}
	query := `DELETE FROM products WHERE id = ?`
	_, err = db.Exec(query, productIdInt)
	if err != nil {
		log.Print("(DeleteProduct) db.Exec", err)
	}

	c.Status(http.StatusOK)
}
