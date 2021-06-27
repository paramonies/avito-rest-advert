package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/paramonies/avito-rest-advert/internal/app/model"
	"github.com/paramonies/avito-rest-advert/internal/app/service"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()

	router.GET("/user/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, "Hello %s", name)
	})

	router.POST("/create", h.createAdvert)
	router.GET("/get/:id", h.getAdvertById)
	router.GET("/list", h.getList)

	return router
}

func (h *Handler) createAdvert(ctx *gin.Context) {
	var input model.Advert

	if err := ctx.BindJSON(&input); err != nil {
		SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.service.CreateAdvert(input)
	if err != nil {
		SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) getAdvertById(ctx *gin.Context) {
	id := ctx.Param("id")
	advertId, err := strconv.Atoi(id)
	if err != nil {
		SendErrorResponse(ctx, http.StatusBadRequest, "advertisement id must be integer")
		return
	}

	fieldsStr := ctx.Query("fields")
	fieldsValid := make([]string, 0)
	fieldsSet := map[string]bool{"description": true, "pictures": true}

	fields := strings.Split(fieldsStr, ",")
	length := len(fields)
	if length == 1 || length == 2 {
		for _, field := range fields {
			field := strings.ToLower(field)
			if _, ok := fieldsSet[field]; ok {
				fieldsValid = append(fieldsValid, field)
			} else {
				fieldsValid = make([]string, 0)
				break
			}
		}
	}

	advert, err := h.service.GetAdvertById(advertId, fieldsValid)
	if err != nil {
		SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, advert)
}

func (h *Handler) getList(ctx *gin.Context) {
	pageStr := ctx.Query("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	orderBy := ctx.Query("order_by")
	orderBySet := map[string]bool{
		"price_desc": true, "price_asc": true,
		"createdat_desc": true, "createdat_asc": true,
	}
	isValid := false
	if _, ok := orderBySet[orderBy]; ok {
		isValid = true
	}

	if orderBy == "" || !isValid {
		orderBy = "createdat_desc"
	}

	adverts, err := h.service.GetAdvertList(page, orderBy)
	if err != nil {
		SendErrorResponse(ctx, http.StatusOK, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, adverts)

}
