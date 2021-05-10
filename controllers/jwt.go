package controllers

import (
    "backend/controllers/types"
    "backend/models"
    "backend/models/user/UserStatus"
    "backend/utils/context"
    "backend/utils/misc"
    "fmt"
    "github.com/google/uuid"
    "github.com/kataras/iris/v12"
    "time"
)

func Login(ctx context.Context) {
    loginRequest := types.LoginRequest{}
    if err := ctx.ParseBody(&loginRequest); err != nil {
        ctx.ReplyBadRequest(err.Error())
        return
    }

    user := models.User{
        Username: loginRequest.Username,
    }

    token, err := user.Login(loginRequest.Password)
    if err != nil {
        ctx.ReplyBadRequest(fmt.Sprintf("Could not login: %s", err.Error()))
        return
    }

    if user.ActivationStatus == UserStatus.Transition {
        ctx.ReplyBadRequest("Please activate account first")
        return
    }

    ctx.Reply(iris.StatusOK, misc.Msg{
        "token": user.Token,
        "issuedAt": time.Unix(token.IssuedAt, 0),
        "expiresAt": time.Unix(token.ExpiresAt, 0),
    })
}

func CreateUser(ctx context.Context) {
    creationRequest := types.CreationRequest{}
    if err := ctx.ParseBody(&creationRequest); err != nil {
        ctx.ReplyBadRequest(err.Error())
        return
    }

    user := models.User{
        Username: creationRequest.Username,
        Password: creationRequest.Password, // this password will be replaced by its hash
        ActivationTicket: uuid.NewString(),
        Role: creationRequest.Role,
    }

    if err := user.Create(); err != nil {
        ctx.ReplyBadRequest(err.Error())
        return
    }

    ctx.Reply(iris.StatusOK, misc.Msg{
        "username": user.Username,
        "ticket": user.ActivationTicket,
    })
}

func ActivateUser(ctx context.Context) {
    activationRequest := types.ActivationRequest{}
    if err := ctx.ParseBody(&activationRequest); err != nil {
        ctx.ReplyBadRequest(err.Error())
        return
    }

    user := models.User{}
    err := user.Activate(activationRequest.ActivationTicket, activationRequest.Password)

    if err != nil {
        ctx.ReplyBadRequest(err.Error())
        return
    }

    token, ok := user.GenerateNewToken()
    if !ok {
        ctx.ReplyInternalError("failed to generate token")
        return
    }

    ctx.Reply(iris.StatusOK, misc.Msg{
        "token": user.Token,
        "issuedAt": time.Unix(token.IssuedAt, 0),
        "expiresAt": time.Unix(token.ExpiresAt, 0),
    })
}
