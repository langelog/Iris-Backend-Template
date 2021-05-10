package models

import (
    "backend/models/user/UserRole"
    "backend/models/user/UserStatus"
    "backend/utils/misc"
    "errors"
    "fmt"
    "github.com/dgrijalva/jwt-go"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"
    "log"
    "os"
    "time"
)

type User struct {
    gorm.Model
    Username         string                `json:"username" gorm:"unique_index;not null"`
    Password         string                `json:"password" gorm:"not null"`
    Role             UserRole.UserRole     `json:"role" gorm:"not null"`
    ActivationTicket string                `json:"ticket" gorm:"unique_index;column:ticket"`
    ActivationStatus UserStatus.UserStatus `json:"status"`
    Token            string                `json:"token" gorm:"-"`
}

func (tuser *User) Fetch(id uint) error {
    return db.Table("users").Where("id = ?", id).First(tuser).Error
}

func (tuser *User) Create() error {
    if err := tuser.Validate(); err != nil {
        return errors.New(fmt.Sprintf("failed to validate user: %s", err.Error()))
    }

    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(tuser.Password), bcrypt.DefaultCost)
    tuser.Password = string(hashedPassword)
    tuser.ActivationStatus = UserStatus.Transition

    db.Create(tuser)

    if tuser.ID <= 0 {
        log.Println("Failed to create user:", tuser.Username)
        return errors.New("failed to create user")
    }
    return nil
}

func (tuser *User) GenerateNewToken() (Token, bool) {
    tokenDuration := time.Duration(misc.GetEnvInt("JWT_DURATION_MIN", 60)) * time.Minute
    tk := Token{UserId: tuser.ID}
    tk.IssuedAt = time.Now().Unix()
    tk.ExpiresAt = time.Now().Add(tokenDuration).Unix()

    token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
    tokenString, _ := token.SignedString([]byte(os.Getenv("JWT_PASS")))
    tuser.Token = tokenString
    claims, ok := token.Claims.(Token)
    return claims, ok
}

func (tuser *User) Validate() error {
    if len(tuser.Username) < 3 {
        return errors.New("username must be at least 3 chars")
    }
    if len(tuser.Password) < 6 {
        return errors.New("password/ticket must have at least 6 chars")
    }

    temp := &User{}
    err := db.Table("users").Where("username = ?", tuser.Username).First(temp).Error
    if err != nil && err != gorm.ErrRecordNotFound {
        return errors.New("connection error. Please retry")
    }
    if temp.Username != "" {
        return errors.New("username is already in use")
    }

    return nil
}

func (tuser *User) Login(pass string) (Token, error) {
    if err := tuser.ValidatePassword(pass); err != nil {
        return Token{}, errors.New(fmt.Sprintf("Could not validate: %s", err.Error()))
    }

    // with valid access, generate jwt
    token, ok := tuser.GenerateNewToken()
    if !ok {
        return Token{}, errors.New("failed to generate token")
    }
    return token, nil
}

func (tuser *User) ValidatePassword(pass string) error {
    err := db.Table("users").Where("username = ?", tuser.Username).First(tuser).Error
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return errors.New("invalid login credentials. Please try again")
        }
        return errors.New("connection error. Please try again")
    }
    err = bcrypt.CompareHashAndPassword([]byte(tuser.Password), []byte(pass))
    if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
        return errors.New("invalid login credentials. Please try again")
    }

    return nil
}

func (tuser *User) UpdatePassword(password string) error {
    if lengthFlag, number, lower, upper, special := misc.VerifyPassword(password); !(lengthFlag && number && lower && upper && special) {
        return errors.New("too week password. It must include at least one number, lowercase, uppercase and a special character. Length must be bigger than 8 chars")
    }

    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    // When password is updated, activation ticket is changed to match user id (activation ticket no longer available)
    err := db.Model(tuser).Updates(User{
        Password:         string(hashedPassword),
        ActivationStatus: UserStatus.Active,
        ActivationTicket: fmt.Sprintf("%v", tuser.ID),
    }).Error
    if err != nil {
        return errors.New("there was an error while updating user password")
    }
    return nil
}

func (tuser *User) Activate(ticket string, pass string) error {
    // find valid user with ticket
    err := db.Table("users").Where("ticket = ?", ticket).First(tuser).Error
    if err != nil {
        return errors.New("could not recognize ticket")
    }
    // check if is not activated already
    if tuser.ActivationStatus == UserStatus.Active {
        return errors.New("user is already active")
    }

    if err := tuser.UpdatePassword(pass); err != nil {
        return err
    }

    return nil
}
