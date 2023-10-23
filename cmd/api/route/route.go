package route

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"time"
)

func Setup(logger *zap.SugaredLogger, timeout time.Duration, db *sql.DB, validate *validator.Validate, gin *gin.Engine) {
	publicRouter := gin.Group("/api/v1")

	NewUserRouter(logger, timeout, db, validate, publicRouter)
}
