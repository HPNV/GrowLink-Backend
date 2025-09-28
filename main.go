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
	businessRepo "github.com/HPNV/growlink-backend/repository/business"
	fileRepo "github.com/HPNV/growlink-backend/repository/file"
	projectRepo "github.com/HPNV/growlink-backend/repository/project"
	skillRepo "github.com/HPNV/growlink-backend/repository/skill"
	studentRepo "github.com/HPNV/growlink-backend/repository/student"
	userRepo "github.com/HPNV/growlink-backend/repository/user"

	//service imports
	businessService "github.com/HPNV/growlink-backend/service/business"
	fileService "github.com/HPNV/growlink-backend/service/file"
	projectService "github.com/HPNV/growlink-backend/service/project"
	skillService "github.com/HPNV/growlink-backend/service/skill"
	studentService "github.com/HPNV/growlink-backend/service/student"
	userService "github.com/HPNV/growlink-backend/service/user"

	//delivery imports
	businessDelivery "github.com/HPNV/growlink-backend/delivery/business"
	fileDelivery "github.com/HPNV/growlink-backend/delivery/file"
	projectDelivery "github.com/HPNV/growlink-backend/delivery/project"
	skillDelivery "github.com/HPNV/growlink-backend/delivery/skill"
	studentDelivery "github.com/HPNV/growlink-backend/delivery/student"
	userDelivery "github.com/HPNV/growlink-backend/delivery/user"
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
	user := userRepo.NewUser(db)
	skill := skillRepo.NewSkill(db)
	business := businessRepo.NewBusiness(db)
	student := studentRepo.NewStudent(db)
	project := projectRepo.NewProject(db)
	file := fileRepo.NewFile(db)

	repo := repository.NewRegistry(
		db,
		user,
		skill,
		business,
		student,
		project,
		file,
	)

	return repo
}

func initService(repo repository.IRegistry) *service.Registry {
	user := userService.NewUser(repo)
	skill := skillService.NewSkill(repo)
	business := businessService.NewBusiness(repo)
	student := studentService.NewStudent(repo)
	project := projectService.NewProject(repo)
	file := fileService.NewFile(repo)

	serviceRegistry := service.NewRegistry(
		user,
		skill,
		business,
		student,
		project,
		file,
	)

	return serviceRegistry
}

func initDelivery(service service.IRegistry) delivery.IDelivery {
	user := userDelivery.NewUser(service)
	business := businessDelivery.NewBusiness(service)
	student := studentDelivery.NewStudent(service)
	skill := skillDelivery.NewSkill(service)
	project := projectDelivery.NewProject(service)
	file := fileDelivery.NewFile(service)

	delivery := delivery.NewDelivery(
		user,
		business,
		student,
		skill,
		project,
		file,
	)

	return delivery
}
