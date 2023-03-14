package data

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/klaxn/klaxn-api/internal/config"
)

type Manager struct {
	db     *gorm.DB
	logger logrus.FieldLogger
}

func New(conf *config.DatabaseConfig) (*Manager, error) {
	var err error
	l := logrus.New()

	gormLogger := logger.New(
		l,
		logger.Config{
			SlowThreshold:             time.Second,
			Colorful:                  true,
			IgnoreRecordNotFoundError: false,
			LogLevel:                  logger.Info,
		},
	)

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/London",
		conf.Hostname,
		conf.User,
		conf.Password,
		conf.DatabaseName,
		conf.Port,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, err
	}

	var items []DatabaseItem
	items = append(items, &Escalation{})
	//items = append(items, &Layer{})
	items = append(items, &Service{})
	items = append(items, &Team{})
	items = append(items, &User{})
	items = append(items, &UserContactMethod{})
	items = append(items, &Link{})
	items = append(items, &Alert{})

	for _, item := range items {
		if err := db.AutoMigrate(item); err != nil {
			return nil, err
		}
	}

	return &Manager{db: db, logger: l}, nil
}
