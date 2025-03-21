package apiserver

import (
	"dialogue/internal/store"
	"html/template"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type server struct {
	store         store.Store
	redisClient   *redis.Client
	router        *gin.Engine
	errorLog      *log.Logger
	templateCache map[string]*template.Template
}

func newServer(store store.Store, redisClient *redis.Client) *server {
	s := &server{
		store:       store,
		redisClient: redisClient,
	}

	s.configureRouter()
	s.configureTemplateCash()
	return s
}

func (s *server) configureTemplateCash() error {
	templateCache, err := newTemplateCache()
	if err != nil {
		return err
	}
	s.templateCache = templateCache
	return nil
}

func (s *server) configureRouter() *gin.Engine {
	s.router = gin.Default()

	s.router.HandleMethodNotAllowed = true

	s.router.Use(s.authenticateMiddleware())

	s.router.Static("/static", "./ui/static")

	s.router.GET("/", s.redirectHomePage)
	s.router.GET("/home", s.homePage)
	s.router.GET("/about", s.about)

	s.router.GET("/newfirstblock", s.emptyFBView)
	s.router.POST("/newfirstblock", s.createFB)
	s.router.GET("/firstblock", s.createdFBView)
	s.router.POST("/firstblock", s.deleteFB)
	s.router.GET("/editfirstblock", s.editFBView)
	s.router.POST("/editfirstblock", s.editFB)

	s.router.GET("/block", s.createdBView)
	s.router.POST("/block", s.deleteB)
	s.router.GET("/editblock", s.editBView)
	s.router.POST("/editblock", s.editB)
	s.router.GET("/{digits:[0-9]+}", s.redirectBlock)

	s.router.GET("/user/signup", s.userSignupView)
	s.router.POST("/user/signup", s.userSignup)

	s.router.GET("/user/login", s.userLoginView)
	s.router.POST("/user/login", s.userLogin)
	s.router.POST("/user/logout", s.userLogout)

	s.router.GET("/account/view", s.accountView)

	s.router.GET("/account/password/update", s.passwordUpdateView)
	s.router.POST("/account/password/update", s.passwordUpdate)

	return s.router
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
