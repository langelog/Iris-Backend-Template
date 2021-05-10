package models

import (
    "backend/models/user/UserRole"
    "github.com/joho/godotenv"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "log"
    "os"
)

var db *gorm.DB

func init() {
    var err error = nil
    db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
    if err != nil {
        panic("failed to connect database")
    }

    // Migrate the schema
    err = db.AutoMigrate(&User{})
    if err != nil {
        return
    }

    initDefaultUsers()
}

func initDefaultUsers() {
    if e := godotenv.Load(); e != nil {
        log.Println("Could not load .env file", e.Error())
    }

    adminUser := User{
        Username: "langelog",
        Password: "testing",
        ActivationTicket: os.Getenv("ADMIN_TICKET"),
        Role: UserRole.Administrator,
    }
    log.Println("TICKET!!!!!!", os.Getenv("ADMIN_TICKET"))
    log.Println("DURATION", os.Getenv("JWT_DURATION_MIN"))

    _ = adminUser.Create()
}
