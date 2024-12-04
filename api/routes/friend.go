package routes

import (
	"github.com/ffaann02/cosplace-server/api/handler"
	"github.com/gofiber/fiber/v2"
)

func FriendRoutes(app fiber.Router) {
	// commision := app.Group("/commission")

	// commision.Get("/", handler.GetCommisions)
	// commision.Get("/:id", handler.GetCommision)
}

func ProtectedFriendRoutes(app fiber.Router) {
	commision := app.Group("/friend")
	// Fetch List
	commision.Get("/list", handler.GetFriendList)
	commision.Get("/requests", handler.GetFriendRequests)
	commision.Get("/waiting-accept", handler.GetFriendWaitingAccept)
	commision.Get("/suggests", handler.GetSuggestions)
	// Manage request
	commision.Post("/send-request", handler.SendFriendRequest)
	commision.Post("/accept-request", handler.AcceptFriendRequest)
	commision.Post("/reject-request", handler.RejectFriendRequest)
	commision.Post("/cancel-request", handler.CancelFriendRequest)
	// Delete friend
	commision.Post("/delete", handler.DeleteFriend)

	//
	commision.Get("/check-status", handler.CheckFriendStatusWithUsername)
}
