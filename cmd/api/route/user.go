package route

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/linqcod/user-service/internal/handlers"
	"github.com/linqcod/user-service/internal/repository"
	"github.com/linqcod/user-service/internal/usecase"
	"go.uber.org/zap"
	"time"
)

func NewUserRouter(logger *zap.SugaredLogger, timeout time.Duration, db *sql.DB, validate *validator.Validate, group *gin.RouterGroup) {
	ageApiUsecase := usecase.NewAgeApiUsecase()
	genderApiUsecase := usecase.NewGenderApiUsecase()
	nationalityApiUsecase := usecase.NewNationalityApiUsecase()

	userRepo := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo, ageApiUsecase, genderApiUsecase, nationalityApiUsecase, timeout)
	userHandler := handlers.NewUserHandler(logger, validate, userUsecase)

	group.POST("/users", userHandler.Create)
	group.DELETE("/users/:id", userHandler.Delete)
	group.PATCH("/users/:id", userHandler.Update)
	group.GET("/users/:id", userHandler.GetById)
	group.GET("/users", userHandler.GetFilteredUsers)
}
