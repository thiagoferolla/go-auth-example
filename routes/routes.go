package routes

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/thiagoferolla/go-auth/providers/cache"
	"github.com/thiagoferolla/go-auth/providers/email"
	"github.com/thiagoferolla/go-auth/providers/jwt"
)

type Router struct {
	Engine   *gin.Engine
	Database *sqlx.DB
}

func NewRouter(engine *gin.Engine, db *sqlx.DB) *Router {
	r := &Router{Engine: engine, Database: db}

	r.RegisterRoutes(engine)

	return r
}

func (r *Router) RegisterRoutes(server *gin.Engine) {

	jwtProvider := jwt.NewProvider()
	emailProvider := email.NewSendgridEmailProvider(os.Getenv("SENDGRID_API_KEY"))
	cacheProvider := cache.NewRedisProvider()

	RegisterAuthRoutes(server, r.Database, *jwtProvider, emailProvider, cacheProvider)
}
