package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/paramonies/avito-rest-advert/internal/app/model"
	"github.com/paramonies/avito-rest-advert/internal/app/service"

	_ "github.com/paramonies/avito-rest-advert/docs"
	ginSwagger "github.com/swaggo/gin-swagger"   // gin-swagger middleware
	"github.com/swaggo/gin-swagger/swaggerFiles" // swagger embed files
)

type Handler struct {
	service service.Service
}

func NewHandler(service service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.POST("/create", h.createAdvert)
	router.GET("/get/:id", h.getAdvertById)
	router.GET("/list", h.getList)

	return router
}

// @Summary создать объявление
// @Tags Advert
// @Description Cоздание нового объявления
// @ID create-advert
// @Accept  json
// @Produce  json
// @Param input body InputAdvert true "Advert info"
// @Success 200 {object} CreateMessageOk
// @Failure 400 {object} CreateMessage400
// @Failure 500 {object} CreateMessage500
// @Router /create [post]
func (h *Handler) createAdvert(ctx *gin.Context) {
	var input model.Advert

	if err := ctx.BindJSON(&input); err != nil {
		SendErrorResponse(ctx, http.StatusBadRequest, "invalid input body")
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

// @Summary получить объявление
// @Tags Advert
// @Description Получить объявление по id
// @ID get-advert-id
// @Accept  html
// @Produce  json
// @Param id path int true "Advert ID"
// @Param fields query string false "Additional Advert fields in response" Enums(description, pictures)
// @Success 200 {object} GetMessageOk
// @Failure 400 {object} GetMessage400
// @Failure 500 {object} GetMessage500
// @Router /get/{id} [get]
func (h *Handler) getAdvertById(ctx *gin.Context) {
	//..../get/:id?fields=description,pictures
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
		SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, advert)
}

// @Summary получить список объявлений
// @Tags Advert
// @Description Получить список объявлений по номеру страницы. На одной странице должно присутствовать 10 объявлений
// @ID get-advert
// @Accept  html
// @Produce  json
// @Param page query int false "Page number"
// @Param order_by query string false "Order field and order destination" Enums(price_desc, price_asc, createdat_desc, createdat_asc)
// @Success 200 {object} ListMessageOk1
// @Failure 404 {object} ListMessage404
// @Failure 500 {object} ListMessage500
// @Router /list [get]
func (h *Handler) getList(ctx *gin.Context) {
	//..../list?page=2&order_by=createdat_desc
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
		SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	if len(adverts) == 0 {
		SendErrorResponse(ctx, http.StatusNotFound, "advertisements not found")
		return
	}

	ctx.JSON(http.StatusOK, adverts)

}
