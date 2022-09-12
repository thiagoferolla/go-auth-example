package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type Router struct {
	Engine  *gin.Engine
	Database *sqlx.DB
}

func NewRouter(engine *gin.Engine, db *sqlx.DB) *Router {
	r := &Router{Engine: engine, Database: db}

	r.RegisterRoutes(engine)

	return r
}

func (r *Router) RegisterRoutes(server *gin.Engine) {
	RegisterAuthRoutes(server, r.Database)
}