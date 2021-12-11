package app

import (
	"github.com/aftaab60/bookstore_users-api/controllers/ping"
	"github.com/aftaab60/bookstore_users-api/controllers/users"
)

func mapUrl() {
	router.GET("/ping", ping.Ping)

	router.GET("/users/:userId", users.Get)
	router.POST("/users", users.Create)
	router.PUT("/users/:userId", users.Update)
	router.PATCH("/users/:userId", users.Update)
	router.DELETE("/users/:userId", users.Delete)
	router.GET("/internal/users/search", users.Search)
	router.POST("/users/login", users.Login)
}
