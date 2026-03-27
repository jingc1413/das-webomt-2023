package storage

// import (
// 	"sync"

// 	"github.com/glebarez/sqlite"
// 	"github.com/pkg/errors"
// 	"github.com/sirupsen/logrus"
// 	"gorm.io/gorm"
// )

// var defaultDB *gorm.DB
// var defaultDBSetupOnce sync.Once

// func GetDefaultDB() *gorm.DB {
// 	defaultDBSetupOnce.Do(func() {
// 		db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
// 		if err != nil {
// 			logrus.Fatal(errors.Wrap(err, "open database").Error())
// 		}
// 		defaultDB = db
// 	})
// 	return defaultDB
// }
