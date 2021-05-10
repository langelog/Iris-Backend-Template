package controllers

import (
    m "backend/models"
    "github.com/kataras/iris/v12"
)

type UserController struct {

}

func (c *UserController) Get() []m.User {
    var response []m.User

    return response
}

func (c *UserController) Post(user m.User) int {
    return iris.StatusNotFound
}
