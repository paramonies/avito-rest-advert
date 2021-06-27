package handler

import "github.com/gin-gonic/gin"

type errorMessage struct {
	Message string `json:"error"`
}

type statusMessage struct {
	Status string `json:"status"`
}

func SendErrorResponse(ctx *gin.Context, statusCode int, message string) {
	ctx.AbortWithStatusJSON(statusCode, errorMessage{message})
}
