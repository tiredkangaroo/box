package main

import (
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var DEBUG = os.Getenv("DEBUG") == "true"

func server(db *DB) error {
	app := echo.New()

	api := app.Group("/api", middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOriginFunc: func(origin string) (bool, error) {
			return origin == "http://localhost:5173" && DEBUG, nil
		},
	}))

	api.GET("/mailboxes", func(c echo.Context) error {
		mailboxes, err := db.ListMailboxes(c.Request().Context(), nil)
		if err != nil {
			slog.Error("list mailboxes", "error", err)
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to list mailboxes"})
		}
		return c.JSON(http.StatusOK, mailboxes)
	})

	api.POST("/mailboxes/:id/sync", func(c echo.Context) error {
		id := c.Param("id")

		idInt, err := strconv.Atoi(id)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid mailbox ID"})
		}
		mailbox, err := db.GetMailbox(c.Request().Context(), nil, idInt)
		if err != nil {
			// note: convert db error into proper error response (not found vs connection error)
			slog.Error("get mailbox", "error", err)
			return c.JSON(http.StatusNotFound, echo.Map{"error": "failed to get mailbox"})
		}

		err = syncMailbox(c.Request().Context(), db, *mailbox)
		if err != nil {
			slog.Error("sync mailbox", "error", err)
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to sync mailbox"})
		}
		return c.JSON(http.StatusOK, echo.Map{"message": "mailbox synced"})
	})

	api.GET("/mailboxes/:id/messages", func(c echo.Context) error {
		id := c.Param("id")

		idInt, err := strconv.Atoi(id)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid mailbox ID"})
		}

		mailbox, err := db.GetMailbox(c.Request().Context(), nil, idInt)
		if err != nil {
			slog.Error("get mailbox", "error", err)
			return c.JSON(http.StatusNotFound, echo.Map{"error": "failed to get mailbox"})
		}

		emails, err := db.ListEmails(c.Request().Context(), nil, mailbox.ID)
		if err != nil {
			slog.Error("list emails", "error", err)
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to list emails"})
		}

		return c.JSON(http.StatusOK, emails)
	})

	api.GET("/mailboxes/:id/messages/:message_id", func(c echo.Context) error {
		id := c.Param("id")
		messageID := c.Param("message_id")

		idInt, err := strconv.Atoi(id)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid mailbox ID"})
		}

		mailbox, err := db.GetMailbox(c.Request().Context(), nil, idInt)
		if err != nil {
			slog.Error("get mailbox", "error", err)
			return c.JSON(http.StatusNotFound, echo.Map{"error": "failed to get mailbox"})
		}

		email, err := db.GetEmail(c.Request().Context(), nil, mailbox.ID, messageID)
		if err != nil {
			slog.Error("get email", "error", err)
			return c.JSON(http.StatusNotFound, echo.Map{"error": "failed to get email"})
		}

		return c.JSON(http.StatusOK, email)
	})

	return app.Start(":9000")
}
