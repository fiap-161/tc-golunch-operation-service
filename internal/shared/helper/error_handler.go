package helper

import (
	"net/http"

	"github.com/gin-gonic/gin"

	apperror "github.com/fiap-161/tc-golunch-operation-service/internal/shared/errors"
)

func HandleError(c *gin.Context, err error) {
	status := http.StatusInternalServerError
	message := "Internal Server Error"

	switch err.(type) {
	case *apperror.ValidationError:
		status = http.StatusBadRequest
		message = "Validation failed"
	case *apperror.UnauthorizedError:
		status = http.StatusUnauthorized
		message = "Unauthorized"
	case *apperror.NotFoundError:
		status = http.StatusBadRequest
		message = "Invalid resource"
	}

	c.JSON(status, apperror.ErrorDTO{
		Message:      message,
		MessageError: err.Error(),
	})

	c.Abort()
}
