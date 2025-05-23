package initlizers

import (
	"AI_WEB_SCRAPPER/models"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// loading the env file
func LoadEnv() {
	if godotenv.Load() != nil {
		panic(" faild to load the env file here is err " + godotenv.Load().Error())
	}
}

// an var that will be working as DB connect for creeaitng or doing somthing for the db
var DB *gorm.DB

func ConnectDB() {
	var err error
	// Konfiguracja DSN dla PostgreSQL
	dsn := "host=" + os.Getenv("DB_HOST") +
		" port=" + os.Getenv("DB_PORT") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" sslmode=disable"

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
}

func CreateTables() {
	if DB.AutoMigrate(models.User{}, models.Todos{}) != nil {
		panic("Failed to create tables: " + DB.AutoMigrate(models.User{}, models.Todos{}).Error())
	}
}
