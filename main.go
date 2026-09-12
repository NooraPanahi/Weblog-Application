package main

import (
	"log"

	"github.com/NooraPanahi/Weblog-Application.git/internal/config"
	"github.com/NooraPanahi/Weblog-Application.git/internal/database"
	"github.com/NooraPanahi/Weblog-Application.git/internal/handler"
	"github.com/NooraPanahi/Weblog-Application.git/internal/middleware"
	"github.com/NooraPanahi/Weblog-Application.git/internal/repo"
	"github.com/NooraPanahi/Weblog-Application.git/internal/service"
	"github.com/NooraPanahi/Weblog-Application.git/internal/session"
	"github.com/labstack/echo/v4"
)

func main () {
	cfg := config.Load()

	db,err := database.NewPostgres(cfg.DatabaseURL)

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	e := echo.New()

	render :=handler.NewTemplateRenderer()
	e.Renderer = render


	userRepo := repo.NewUserRepo(db)
	authService := service.NewAuthService(userRepo)
	sessionManager := session.NewManager(cfg.SessionSecret)

	authHandler := handler.NewAuthHandler(authService, sessionManager)

	weblogRepo := repo.NewWeblogRepo(db)
	weblogService := service.NewWeblogService(weblogRepo)
	homeHandler := handler.NewHomeHandler(weblogService)

	commentRepo := repo.NewCommentRepo(db)
	commentService := service.NewCommentService(commentRepo, weblogService)
	commentHandler := handler.NewCommentHandler(commentService)

	shareRepo := repo.NewWeblogShareRepo(db)
	shareService := service.NewWeblogShareService(shareRepo, userRepo, weblogRepo)

	shareHandler := handler.NewWeblogShareHandler(shareService, weblogService, commentService)

	weblogHandler := handler.NewWeblogHandler(weblogService, commentService)


	e.GET("/register", authHandler.ShowRegister, middleware.RequireGuest(sessionManager))
	e.POST("/register", authHandler.Register, middleware.RequireGuest(sessionManager))

	e.GET("/login", authHandler.ShowLogin, middleware.RequireGuest(sessionManager))
	e.POST("/login", authHandler.Login,  middleware.RequireGuest(sessionManager))

	e.POST("/logout", authHandler.Logout, middleware.RequireAuth(sessionManager, userRepo))

	e.GET("/weblog/create", weblogHandler.ShowCreate, middleware.RequireAuth(sessionManager, userRepo))
	e.POST("/weblog/create", weblogHandler.Create, middleware.RequireAuth(sessionManager, userRepo))

	e.GET("/weblog/:id", weblogHandler.Detail, middleware.RequireAuth(sessionManager, userRepo))

	e.POST("/weblog/:id/share", shareHandler.Share, middleware.RequireAuth(sessionManager, userRepo))
	e.POST("/weblog/:id/delete", weblogHandler.Delete, middleware.RequireAuth(sessionManager, userRepo))
	e.GET("/", homeHandler.Home, middleware.RequireAuth(sessionManager, userRepo))

	e.POST("/weblog/:id/comments", commentHandler.Create, middleware.RequireAuth(sessionManager, userRepo))
	e.POST("/comments/:id/delete", commentHandler.Delete, middleware.RequireAuth(sessionManager, userRepo))

	e.Static("/static", "statics")
	e.Logger.Fatal(e.Start(":"+ cfg.Port))
}