package main

import (
    "backend/controllers"
    "backend/middlewares"
    "backend/models/user/UserRole"
    "backend/utils/context"
    "github.com/kataras/iris/v12"
)

func Routes(app *iris.Application) {
    api := app.Party("/api")
    {
        apiJwt := api.Party("/jwt", context.Ctx(middlewares.JwtAuthentication))
        {
            // --
            apiJwt.Post("/login", context.Ctx(controllers.Login))
            // --
            apiJwt.Post("/activate", context.Ctx(controllers.ActivateUser))
            // --
            apiJwt.Post("/create", middlewares.RoleChecker(UserRole.RoleList{
                    UserRole.Administrator,
                },
                ), context.Ctx(controllers.CreateUser))
            // --

        }
    }
}
