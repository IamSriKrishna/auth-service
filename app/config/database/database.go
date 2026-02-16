package database

import (
	"github.com/bbapp-org/auth-service/app/config"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

var (
	DB          *gorm.DB 
)

func ConnectDatabase(config *config.Config) {
	var err error

	// Connect to primary database (write)
	primaryDSN := config.GetDatabaseDSN()
	log.Printf("Connecting to MySQL primary database: %s@tcp(%s:%d)/%s",
		config.Database.User,
		config.Database.Host,
		config.Database.Port,
		config.Database.DBName)

	gormConfig := &gorm.Config{}

	// Set logger level based on environment
	if config.App.Environment == "development" || config.App.Environment == "dev" {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
		log.Println("Database logging enabled (development mode)")
	} else {
		gormConfig.Logger = logger.Default.LogMode(logger.Silent)
		log.Println("Database logging disabled (production mode)")
	}

	DB, err = gorm.Open(mysql.Open(primaryDSN), gormConfig)

	if err != nil {
		log.Printf("Failed to connect to MySQL primary database: %v", err)
		log.Printf("Database configuration: Host=%s, Port=%d, User=%s, Database=%s",
			config.Database.Host,
			config.Database.Port,
			config.Database.User,
			config.Database.DBName)
		log.Fatal("Database connection failed")
	}

	log.Println("MySQL primary database connected successfully")

	// Configure read replica using dbresolver plugin if configured
	readReplicaDSN := config.GetReadReplicaDSN()
	if readReplicaDSN != "" {
		// Use read replica user if provided, otherwise fall back to primary user
		replicaUser := config.Database.ReadReplicaUser
		if replicaUser == "" {
			replicaUser = config.Database.User
		}
		log.Printf("Configuring read replica with dbresolver: %s@tcp(%s:%d)/%s",
			replicaUser,
			config.Database.ReadReplicaHost,
			config.Database.ReadReplicaPort,
			config.Database.DBName)

		// Register dbresolver plugin for read/write splitting
		err = DB.Use(dbresolver.Register(dbresolver.Config{
			// Sources: primary database (for writes)
			Sources: []gorm.Dialector{mysql.Open(primaryDSN)},
			// Replicas: read replica (for reads)
			Replicas: []gorm.Dialector{mysql.Open(readReplicaDSN)},
			// Policy: use round-robin for load balancing (if multiple replicas)
			Policy: dbresolver.RandomPolicy{},
		}).SetConnMaxIdleTime(3600).
			SetConnMaxLifetime(7200).
			SetMaxIdleConns(10).
			SetMaxOpenConns(100))

		if err != nil {
			log.Printf("Failed to configure read replica with dbresolver: %v", err)
			log.Printf("Read replica configuration: Host=%s, Port=%d, User=%s, Database=%s",
				config.Database.ReadReplicaHost,
				config.Database.ReadReplicaPort,
				replicaUser,
				config.Database.DBName)
			log.Println("Warning: Read replica configuration failed, using primary database for all operations")
		} else {
			log.Println("Read replica configured successfully with dbresolver - reads will be routed to replica, writes to primary")
		}
	} else {
		// No read replica configured, all operations use primary
		log.Println("Read replica not configured, using primary database for all operations")
	}
}



// GetDB returns the database connection with dbresolver enabled
// dbresolver automatically routes:
// - Read operations (SELECT, First, Find, etc.) to read replica if configured
// - Write operations (Create, Update, Delete, Save, etc.) to primary database
func GetDB() *gorm.DB {
	return DB
}

