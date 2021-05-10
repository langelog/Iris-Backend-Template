package context

import (
    "backend/models"
    "backend/utils/misc"
    "encoding/json"
    "errors"
    "github.com/kataras/iris/v12"
    "log"
)

type Context struct {
    iris.Context
}

func Ctx(f func(Context)) iris.Handler {
    return func(ctx iris.Context) {
        f(Context{Context: ctx})
    }
}

func (ctx Context) GetUser() (*models.User, error) {
    user, ok := ctx.Values().Get("user").(*models.User)
    if !ok {
        return nil, errors.New("could not identify user. Please renew your token")
    }
    return user, nil
}

func (ctx Context) ParseBody(targetStructure interface{}) error {
    if err := json.NewDecoder(ctx.Request().Body).Decode(targetStructure); err != nil {
        return errors.New("could not parse input")
    }
    return nil
}

func (ctx Context) Reply(status int, msg misc.Msg) {
    ctx.StatusCode(status)
    if _, err := ctx.JSON(msg); err != nil {
        log.Println("Error while replying to client", err.Error())
    }
}

func (ctx Context) ReplyForbidden(msg string) {
    ctx.Reply(iris.StatusForbidden, misc.Msg{
        "msg": msg,
    })
}

func (ctx Context) ReplyBadRequest(msg string) {
    ctx.Reply(iris.StatusBadRequest, misc.Msg{
        "msg": msg,
    })
}

func (ctx Context) ReplyInternalError(msg string) {
    ctx.Reply(iris.StatusInternalServerError, misc.Msg{
        "msg": msg,
    })
}