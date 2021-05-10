package middlewares

import (
    "backend/models"
    "backend/models/user/UserStatus"
    "backend/utils/context"
    "backend/utils/misc"
    "fmt"
    "github.com/dgrijalva/jwt-go"
    "github.com/kataras/iris/v12"
    "os"
    "strings"
)

var notAuthPaths = []string {
    "/api/jwt/login",
    "/api/jwt/activate",
}

func JwtAuthentication(ctx context.Context) {
    requestPath := ctx.Path()

    for _, value := range notAuthPaths {
        if value == requestPath {
            ctx.Next()
            return
        }
    }

    tokenHeader := ctx.GetHeader("Authorization")
    if tokenHeader == "" {
        ctx.ReplyForbidden("Missing authentication token")
        return
    }

    tokenSplit := strings.Split(tokenHeader, " ")
    if len(tokenSplit) != 2 {
        ctx.ReplyForbidden("Invalid/Malformed Authentication Token")
        return
    }

    tokenPart := tokenSplit[1]
    tk := &models.Token{}

    token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
        return []byte(os.Getenv("JWT_PASS")), nil
    })
    if err != nil {
        ctx.ReplyForbidden(fmt.Sprintf("Token is not valid: %s", err.Error()))
        return
    }

    err = token.Claims.Valid()
    if err != nil {
        ctx.ReplyForbidden(fmt.Sprintf("Token is not valid: %s", err.Error()))
        return
    }

    user := models.User{}

    err = user.Fetch(tk.UserId)
    if err != nil {
        ctx.ReplyForbidden("There was a problem with you authentication. Please generate a new one.")
        return
    }

    if user.ActivationStatus == UserStatus.Transition {
        ctx.Reply(iris.StatusBadRequest, misc.Msg{
            "msg": "Please activate your account first.",
        })
        return
    }

    ctx.Values().Set("user", &user)
    ctx.Next()
}
