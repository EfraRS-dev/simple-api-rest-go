package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Item struct {
	ProductID string  `json:"product_id" binding:"required"`
	Quantity  int     `json:"quantity" binding:"required,min=1"`
	Price     float64 `json:"price" binding:"required,min=0"`
}

type OrderRequest struct {
	CustomerID string `json:"customer_id" binding:"required"`
	Items      []Item `json:"items" binding:"required,min=1"`
}

type OrderResponse struct {
	OrderID        int       `json:"order_id"`
	CustomerID     string    `json:"customer_id"`
	TotalAmount    float64   `json:"total_amount"`
	ItemsCount     int       `json:"items_count"`
	ProcessingDate time.Time `json:"processing_date"`
}

type Order struct {
	ID          int       `json:"id"`
	CustomerID  string    `json:"customer_id"`
	TotalAmount float64   `json:"total_amount"`
	ItemsCount  int       `json:"items_count"`
	CreatedAt   time.Time `json:"created_at"`
}

var db *sql.DB

func initDB() {
	var err error
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Error opening database: ", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}

	fmt.Println("Successfully connected to Postgres database!")
}

// GET /orders - Obtener todas las órdenes
func getOrders(c *gin.Context) {
	rows, err := db.Query("SELECT id, customer_id, total_amount, items_count, created_at FROM orders ORDER BY created_at DESC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var order Order
		err := rows.Scan(&order.ID, &order.CustomerID, &order.TotalAmount, &order.ItemsCount, &order.CreatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		orders = append(orders, order)
	}

	c.JSON(http.StatusOK, orders)
}

// POST /orders - Crear una nueva orden con items
func createOrder(c *gin.Context) {
	var orderReq OrderRequest
	if err := c.ShouldBindJSON(&orderReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Calcular total_amount e items_count
	var totalAmount float64
	var itemsCount int = 0

	for _, item := range orderReq.Items {
		totalAmount += item.Price * float64(item.Quantity)
		itemsCount += item.Quantity
	}

	// Iniciar transacción
	tx, err := db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error starting transaction"})
		return
	}
	defer tx.Rollback()

	// Insertar orden
	var orderID int
	var createdAt time.Time
	err = tx.QueryRow(
		"INSERT INTO orders (customer_id, total_amount, items_count) VALUES ($1, $2, $3) RETURNING id, created_at",
		orderReq.CustomerID, totalAmount, itemsCount,
	).Scan(&orderID, &createdAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating order: " + err.Error()})
		return
	}

	// Insertar items
	for _, item := range orderReq.Items {
		_, err := tx.Exec(
			"INSERT INTO items (order_id, product_id, quantity, price) VALUES ($1, $2, $3, $4)",
			orderID, item.ProductID, item.Quantity, item.Price,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating items: " + err.Error()})
			return
		}
	}

	// Confirmar transacción
	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error committing transaction"})
		return
	}

	// Preparar respuesta
	response := OrderResponse{
		OrderID:        orderID,
		CustomerID:     orderReq.CustomerID,
		TotalAmount:    totalAmount,
		ItemsCount:     itemsCount,
		ProcessingDate: createdAt,
	}

	c.JSON(http.StatusCreated, response)
}

// DELETE /orders/:id - Eliminar una orden
func deleteOrder(c *gin.Context) {
	id := c.Param("id")

	result, err := db.Exec("DELETE FROM orders WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order deleted successfully"})
}

func main() {
	initDB()
	defer db.Close()

	router := gin.Default()

	// Rutas
	router.GET("/api/orders", getOrders)
	router.POST("/api/orders", createOrder)
	router.DELETE("/api/orders/:id", deleteOrder)

	// Health check
	router.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	log.Println("Server starting on port 5001...")
	router.Run(":5001")
}
