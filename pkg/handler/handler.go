package handler

import (
	"github.com/Neokrid/game-review/pkg/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	h.initAuthRoutes(router)
	h.initGameRoutes(router)
	h.initReviewsRoutes(router)
	h.initAPIRoutes(router)
	return router
}

func (h *Handler) initAuthRoutes(router *gin.Engine) {
	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp) //complete
		auth.POST("/sign-in", h.signIn) //complete
	}
}

func (h *Handler) initGameRoutes(router *gin.Engine) {
	games := router.Group("/games")
	{
		games.GET("/", h.getAllGames)               //complete
		games.GET("/:id", h.getGamesById)           //complete
		games.GET(":id/reviews", h.getGamesReviews) //complete
	}
}

func (h *Handler) initReviewsRoutes(router *gin.Engine) {
	reviews := router.Group("/reviews")
	{
		reviews.GET("/:id", h.getReviewById) //complete
	}
}

func (h *Handler) initAPIRoutes(router *gin.Engine) {
	api := router.Group("/api", h.userIdentity) //complete
	{
		gamesAuth := api.Group("/games")
		{
			gamesAuth.DELETE("/:id", h.deleteGame) //complete
			gamesAuth.POST("/", h.createGame)      //complete
			gamesAuth.PUT("/:id", h.changeGame)    //complete

			reviewsAuth := gamesAuth.Group(":id/reviews")
			{
				reviewsAuth.DELETE("/:review_id", h.deleteReview) //complete
				reviewsAuth.POST("/", h.createReview)             //complete
				reviewsAuth.PUT("/:review_id", h.changeReview)    //complete
			}
		}

	}
}
