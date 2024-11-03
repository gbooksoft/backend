package main

import (
	"github.com/gbooksoft/backend/internal/routing"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func main() {
	e := gin.Default()

	routing.Add(e)

	log.Debug().Msg("Server started: http://localhost:8080")
	if err := e.Run("localhost:8080"); err != nil {
		log.Fatal().Msgf("server shutdown %v", err)
	}
}
