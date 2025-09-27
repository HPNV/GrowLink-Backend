package main

import (
	"fmt"
	"log"

	"github.com/HPNV/growlink-backend/config"
	"github.com/HPNV/growlink-backend/delivery"
	"github.com/HPNV/growlink-backend/migration"
	"github.com/HPNV/growlink-backend/repository"
	"github.com/HPNV/growlink-backend/routing"
	"github.com/HPNV/growlink-backend/service"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	//repository imports
	developerRepo "github.com/HPNV/growlink-backend/repository/developer"

	//service imports
	developerService "github.com/HPNV/growlink-backend/service/developer"

	//delivery imports
	developerDelivery "github.com/HPNV/growlink-backend/delivery/developer"
)

func main() {
	config.Init()

	db, err := connect(config.CFG.DB)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if config.CFG.DB.Migration {
		if err := migration.AutoMigrate(db, "./migration"); err != nil {
			log.Fatal("Failed to run migrations:", err)
		}
		log.Println("Database migrations completed successfully")
	}

	repo := initRepository(db)
	service := initService(repo)
	delivery := initDelivery(service)

	fmt.Println("Starting server on port:", config.CFG.Server)

	routing.NewRoute(config.CFG.Server, delivery).SetupRoutes()
}

func connect(
	config config.DatabaseConfig,
) (*sqlx.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.Name)

	fmt.Println("Connecting to database:", psqlInfo)

	db, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	fmt.Println("Successfully connected to database!")

	return db, nil
}

func initRepository(db *sqlx.DB) *repository.Registry {
	developer := developerRepo.NewDeveloperRepo(db)
	repo := repository.NewRegistry(db, developer)

	return repo
}

func initService(repo repository.IRegistry) *service.Registry {
	developer := developerService.NewDeveloperService(repo)

	serviceRegistry := service.NewRegistry(developer)

	return serviceRegistry
}

func initDelivery(service service.IRegistry) delivery.IDelivery {
	developer := developerDelivery.NewDeveloperService(service)
	delivery := delivery.NewDelivery(developer)

	return delivery
}
