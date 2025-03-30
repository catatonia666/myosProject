package apiserver

import (
	"dialogue/internal/services"
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

	services services.Service
}

func newServer(store store.Store, redisClient *redis.Client) *server {
	s := &server{
		store:       store,
		redisClient: redisClient,
		services:    services.New(store),
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

	s.router.GET("/stories/startingblocks/new", s.startingBlockForm)
	s.router.POST("/stories/startingblocks/new", s.startingBlockCreate)
	s.router.GET("/stories/startingblocks/:id", s.startingBlockRender)
	s.router.POST("/stories/startingblocks/:id", s.deleteWholeStory) //delete
	s.router.GET("/stories/startingblocks/:id/edit", s.startingBlockEditionForm)
	s.router.POST("/stories/startingblocks/:id/edit", s.startingBlockEdit) //patch

	s.router.GET("/stories/blocks/:id", s.blockRender)
	s.router.POST("/stories/blocks/:id", s.deleteBlock) //delete
	s.router.GET("/stories/blocks/:id/edit", s.blockEditionForm)
	s.router.POST("/stories/blocks/:id/edit", s.blockEdit) //patch

	s.router.GET("/user/signup", s.userSignupForm)
	s.router.POST("/user/signup", s.userSignup)

	s.router.GET("/user/login", s.userLoginForm)
	s.router.POST("/user/login", s.userLogin)
	s.router.POST("/user/logout", s.userLogout)

	s.router.GET("/account/view", s.userAccountRender)

	s.router.GET("/account/password/update", s.passwordEditionForm)
	s.router.POST("/account/password/update", s.passwordUpdate) //patch
	return s.router
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
