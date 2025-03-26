package main

import (
	"database/sql"
	"fmt"
	"os"

	"sample/internal/db"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var queries *db.Queries
var sqlDB *sql.DB

func initDB() error {
	// .envファイルの読み込み（ローカル環境の場合のみ）
	if err := godotenv.Load(); err != nil {
		// .envファイルが存在しない場合は無視（Cloud Run環境）
		fmt.Println("Warning: .env file not found")
	}

	// Cloud SQL接続情報
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	instance := os.Getenv("DB_INSTANCE")
	database := os.Getenv("DB_NAME")

	// Cloud Run環境ではUnixドメインソケットを使用
	dsn := fmt.Sprintf("%s:%s@unix(/cloudsql/%s)/%s?parseTime=true",
		user, password, instance, database)

	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	// 接続テスト
	err = sqlDB.Ping()
	if err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	queries = db.New(sqlDB)
	return nil
}

func main() {
	// データベース接続の初期化
	if err := initDB(); err != nil {
		fmt.Printf("Failed to initialize database: %v\n", err)
		os.Exit(1)
	}
	defer sqlDB.Close()

	r := gin.Default()

	// サンプルデータを取得するエンドポイント
	r.GET("/data", func(c *gin.Context) {
		items, err := queries.GetItems(c.Request.Context())
		if err != nil {
			c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to query database: %v", err)})
			return
		}

		c.JSON(200, gin.H{
			"data": items,
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(fmt.Sprintf(":%s", port))
}
