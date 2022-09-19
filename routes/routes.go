package routes

import (
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

	jwtProvider := jwt.NewBaseProvider()
	// emailProvider := email.NewSendgridEmailProvider(os.Getenv("SENDGRID_API_KEY"))
	emailProvider := email.NewMockEmailProvider()
	cacheProvider := cache.NewRedisProvider()

	RegisterAuthRoutes(server, r.Database, *jwtProvider, emailProvider, cacheProvider)
}
