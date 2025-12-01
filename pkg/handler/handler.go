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
	h.initSearchRoutes(router)
	return router
}

func (h *Handler) initAuthRoutes(router *gin.Engine) {
	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}
}

func (h *Handler) initGameRoutes(router *gin.Engine) {
	games := router.Group("/games")
	{
		games.GET("/", h.getAllGames)
		games.GET("/:id", h.getGamesById)
		games.GET(":id/reviews", h.getGamesReviews)
		games.GET("/leaderboard", h.getLeaderboard)
		games.GET(":id/rating-history", h.getRatingHistory)
	}
}

func (h *Handler) initReviewsRoutes(router *gin.Engine) {
	reviews := router.Group("/reviews")
	{
		reviews.GET("/:id", h.getReviewById)
	}
}

func (h *Handler) initAPIRoutes(router *gin.Engine) {
	api := router.Group("/api", h.userIdentity)
	{
		gamesAuth := api.Group("/games")
		{
			gamesAuth.DELETE("/:id", h.deleteGame)
			gamesAuth.POST("/", h.createGame)
			gamesAuth.PUT("/:id", h.changeGame)

			reviewsAuth := gamesAuth.Group(":id/reviews")
			{
				reviewsAuth.DELETE("/:review_id", h.deleteReview)
				reviewsAuth.POST("/", h.createReview)
				reviewsAuth.PUT("/:review_id", h.changeReview)
			}
		}

	}
}

func (h *Handler) initSearchRoutes(router *gin.Engine) {
	search := router.Group("/search")
	{
		search.GET("/", h.searchGame)
	}
}
