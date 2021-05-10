package middlewares

import (
    "backend/models"
    "backend/models/user/UserRole"
    "backend/utils/context"
    "github.com/kataras/iris/v12"
)

func RoleChecker(roles UserRole.RoleList) func(iris.Context) {
    return context.Ctx(func(ctx context.Context) {
        var user *models.User
        var err error

        if user, err = ctx.GetUser(); err != nil {
            ctx.ReplyBadRequest(err.Error())
            return
        }
        for _, v := range roles {
            if user.Role == v {
                ctx.Next()
                return
            }
        }
        ctx.ReplyForbidden("You're not allowed")
        return
    })
}