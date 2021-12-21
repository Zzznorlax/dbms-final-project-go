package controller

import (
	"dbmsbackend/domain"
	"dbmsbackend/service"
	"dbmsbackend/util"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type User struct {
	userService *service.User
}

func (controller *User) Initialize(config *util.Config) (err error) {
	controller.userService = new(service.User)
	err = controller.userService.Initialize(config)

	if err != nil {
		err = fmt.Errorf("initializing user controller: %w", err)
	}

	return err
}

func (controller *User) New(c *gin.Context) {

	type newUserDTO struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		Password string `json:"password"`
		Kind     string `json:"type"`
	}

	dto := new(newUserDTO)
	err := c.BindJSON(dto)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, util.GeneralAPIResponse{
			Status:  http.StatusUnprocessableEntity,
			Message: err.Error(),
		})
		return
	}

	entity, err := controller.userService.New(dto.Name, dto.Email, dto.Phone, domain.UserType(dto.Kind), dto.Password)

	if err != nil {
		c.JSON(http.StatusInternalServerError, util.GeneralAPIResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, entity.ToDTO())
}

func (controller *User) Update(c *gin.Context) {

	rawID := c.Param("id")

	id, err := strconv.Atoi(rawID)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, util.GeneralAPIResponse{
			Message: err.Error(),
		})
		return
	}

	type updateUserDTO struct {
		Name  string `json:"name"`
		Phone string `json:"phone"`
	}

	dto := new(updateUserDTO)

	err = c.BindJSON(dto)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}

	entity, err := controller.userService.Update(id, dto.Name, dto.Phone)

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, util.GeneralAPIResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, *entity.ToDTO())
}

func (controller *User) Login(config *util.Config) gin.HandlerFunc {

	return func(c *gin.Context) {
		type loginDTO struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		type tokenDTO struct {
			ID    int    `json:"id"`
			Token string `json:"token"`
		}

		type tokenRespDTO struct {
			Data tokenDTO `json:"data"`
		}

		dto := new(loginDTO)

		err := c.BindJSON(dto)

		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, err)
			return
		}

		user, err := controller.userService.Login(dto.Email, dto.Password)

		if err != nil {
			c.JSON(http.StatusUnauthorized, err)
			return
		}

		token, err := util.NewToken(config, user.ID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, util.GeneralAPIResponse{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, tokenRespDTO{
			Data: tokenDTO{
				ID:    user.ID,
				Token: token,
			},
		})
	}
}

func (controller *User) GetCurrentUser(c *gin.Context) {

	id := c.GetInt("userID")

	entity, err := controller.userService.GetByID(id)

	if err != nil {
		if err, ok := err.(*util.ErrNotFound); ok {
			c.JSON(http.StatusNotFound, util.GeneralAPIResponse{
				Message: err.Error(),
			})
			return

		} else {
			c.JSON(http.StatusInternalServerError, util.GeneralAPIResponse{
				Message: err.Error(),
			})
			return
		}

	}

	c.JSON(http.StatusOK, *entity.ToDTO())
}
