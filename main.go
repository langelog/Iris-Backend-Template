package main

import (
    "fmt"
    "github.com/joho/godotenv"
    "github.com/kataras/iris/v12"
    "github.com/kataras/iris/v12/middleware/logger"
    "github.com/kataras/iris/v12/middleware/recover"
    "log"
    "os"
)

type Context struct {
    iris.Context
}

func main() {
    if e := godotenv.Load(); e != nil {
        log.Println("Could not load .env file", e.Error())
    }

    app := iris.New()

    app.Logger().SetLevel("debug")

    app.Use(recover.New())
    app.Use(logger.New())

    //booksAPI := app.Party("/books")
    //{
    //    m := mvc.New(booksAPI)
    //    m.Handle(new(BookController))
    //}

    Routes(app)
    // todo: incorporate socket-io

    if err := app.Run(
        iris.Addr(fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT"))),
        iris.WithoutServerError(iris.ErrServerClosed)); err != nil {
        log.Fatalln(err)
    }
}


