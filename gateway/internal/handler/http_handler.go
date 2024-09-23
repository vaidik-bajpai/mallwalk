package handler

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/vaidik-bajpai/mallwalk/common"
	"github.com/vaidik-bajpai/mallwalk/gateway/internal/gateway"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	gateway  gateway.APIGateway
	validate *validator.Validate
}

func NewHandler(gateway gateway.APIGateway, validate *validator.Validate) *Handler {
	return &Handler{gateway: gateway, validate: validate}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /user/signup", h.HandleCreateUser)
	mux.HandleFunc("POST /user/login", h.HandleUserLogin)

	mux.HandleFunc("POST /products/create", h.HandleCreateProduct)
	mux.HandleFunc("GET /products/get/{productID}", h.HandleGetProduct)
	mux.HandleFunc("GET /products/list", h.HandleGetProduct)
	mux.HandleFunc("PUT /products/update/{productID}", h.HandleUpdateProduct)
	mux.HandleFunc("DELETE /products/delete/{productID}", h.HandleDeleteProduct)

	mux.HandleFunc("POST /cart/add/{cartID}", h.HandleAddToCart)
	mux.HandleFunc("DELETE /cart/{cartID}/remove/{productID}", h.HandleRemoveFromCart)
	mux.HandleFunc("GET /cart/{cartID}/get", h.HandleListCartItems)

	mux.HandleFunc("POST /order/create", h.HandleCreateOrder)
	mux.HandleFunc("POST /payment/{orderID}", h.HandlePayments)
}

type UserRegister struct {
	Username string `json:"username" validate:"required,min=3,max=20"`
	Password string `json:"password" validate:"required,min=8,max=30"`
	Email    string `json:"email" validate:"required,email"`
}

func (h *Handler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	var userDetails UserRegister

	err := common.ReadJSON(r, &userDetails)
	if err != nil {
		common.WriteError(w, http.StatusBadRequest, "invalid user credentials")
		return
	}

	err = h.validate.Struct(userDetails)
	if err != nil {
		common.WriteError(w, http.StatusBadRequest, "invalid user credentials")
		return
	}

	ctx := context.Background()
	_, err = h.gateway.CreateUser(ctx, toPBCreateUser(&userDetails))
	rStatus := status.Convert(err)
	if rStatus != nil {
		switch {
		case rStatus.Code() == codes.InvalidArgument:
			common.WriteError(w, http.StatusBadRequest, "invalid arguments")
		case rStatus.Code() == codes.AlreadyExists:
			common.WriteError(w, http.StatusConflict, "users already exists")
		default:
			common.WriteError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	common.WriteJSON(w, http.StatusCreated, "user created successfully")
}

type UserLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `password:"password" validate:"required,min=8,max=30"`
}

func (h *Handler) HandleUserLogin(w http.ResponseWriter, r *http.Request) {
	var userDetails UserLogin
	err := common.ReadJSON(r, &userDetails)
	if err != nil {
		common.WriteError(w, http.StatusBadRequest, "invalid user credentials")
		return
	}

	err = h.validate.Struct(userDetails)
	if err != nil {
		common.WriteError(w, http.StatusBadRequest, "invalid user credentials")
		return
	}

	ctx := context.Background()
	res, err := h.gateway.UserLogin(ctx, toPBLoginUser(&userDetails))
	rStatus := status.Convert(err)
	if rStatus != nil {
		switch {
		case rStatus.Code() == codes.InvalidArgument:
			common.WriteError(w, http.StatusBadRequest, "invalid arguments")
		case rStatus.Code() == codes.NotFound:
			common.WriteError(w, http.StatusNotFound, "user not registered")
		default:
			common.WriteError(w, http.StatusInternalServerError, "something went wrong with our server "+err.Error())
		}
		return
	}

	common.WriteJSON(w, http.StatusOK, map[string]string{"token": res.Token})
}

type Product struct {
	ID          primitive.ObjectID `json:"id" validate:"omitempty"`
	Name        string             `json:"name" validate:"required"`
	Price       uint32             `json:"price" validate:"required,min=10"`
	Stock       uint32             `json:"stock" validate:"required,min=1"`
	Description string             `json:"description" validate:"required,min=20,max=300"`
	Category    string             `json:"category" validate:"required,min=4,max=20"`
	Rating      float32            `json:"rating" validate:"required"`
	Image       string             `json:"image" validate:"required"`
}

func (h *Handler) HandleCreateProduct(w http.ResponseWriter, r *http.Request) {
	var product Product
	err := common.ReadJSON(r, &product)
	if err != nil {
		log.Println(err)
		common.WriteError(w, http.StatusBadRequest, "invalid request credentials")
		return
	}

	err = h.validate.Struct(product)
	if err != nil {
		log.Println(err)
		common.WriteError(w, http.StatusBadRequest, "invalid request credentials")
		return
	}

	ctx := context.Background()
	p, err := h.gateway.CreateProduct(ctx, toPBCreateProductRequest(product))
	rStatus := status.Convert(err)
	if rStatus != nil {
		log.Println(err)
		common.WriteError(w, http.StatusInternalServerError, "something went wrong with our servers")
		return
	}

	common.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "product created successfully",
		"product": p,
	})
}

func (h *Handler) HandleGetProduct(w http.ResponseWriter, r *http.Request) {
	productID := r.PathValue("productID")

	p, err := h.gateway.GetProduct(context.Background(), toPBGetProduct(productID))
	rStatus := status.Convert(err)
	if rStatus != nil {
		log.Println(err)
		common.WriteError(w, http.StatusInternalServerError, "something went wrong with our servers")
		return
	}

	common.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"product": p,
	})
}

func (h *Handler) HandleListProducts(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	pageNumber := query.Get("page")
	pageSize := query.Get("size")
	category := query.Get("category")
	minRating := query.Get("min_rating")

	uintPage, err := strconv.ParseUint(pageNumber, 10, 32)
	if err != nil {
		log.Println(err)
		common.WriteError(w, http.StatusBadRequest, "invalid credentials")
		return
	}

	uintSize, err := strconv.ParseUint(pageSize, 10, 32)
	if err != nil {
		log.Println(err)
		common.WriteError(w, http.StatusBadRequest, "invalid credentials")
		return
	}

	var floatRating float64
	if minRating != "" {
		floatRating, err = strconv.ParseFloat(minRating, 32)
		if err != nil {
			log.Println(err)
			common.WriteError(w, http.StatusBadRequest, "invalid credentials")
			return
		}
	}

	ctx := context.Background()
	res, err := h.gateway.ListProduct(ctx, toPBListProductRequest(uint32(uintPage), uint32(uintSize), category, float32(floatRating)))
	rStatus := status.Convert(err)
	if rStatus != nil {
		log.Println(err)
		common.WriteError(w, http.StatusInternalServerError, "something went wrong with our server")
		return
	}

	common.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"products": res,
	})
}

func (h *Handler) HandleUpdateProduct(w http.ResponseWriter, r *http.Request) {
	strID := r.PathValue("productID")

	id, err := primitive.ObjectIDFromHex(strID)
	if err != nil {
		log.Println(err)
		common.WriteError(w, http.StatusBadRequest, "invalid credentials")
		return
	}

	var product Product
	product.ID = id
	err = common.ReadJSON(r, &product)
	if err != nil {
		log.Println(err)
		common.WriteError(w, http.StatusBadRequest, "invalid credentials")
		return
	}

	ctx := context.Background()

	p, err := h.gateway.UpdateProduct(ctx, toPBUpdateProduct(product))
	rStatus := status.Convert(err)
	if rStatus != nil {
		log.Println(err)
		common.WriteError(w, http.StatusInternalServerError, "something went wrong with our server")
		return
	}

	common.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message":     "product was updated successfully",
		"new_product": p,
	})
}

func (h *Handler) HandleDeleteProduct(w http.ResponseWriter, r *http.Request) {
	productID := r.PathValue("productID")

	ctx := context.Background()
	res, err := h.gateway.DeleteProduct(ctx, toPBDeleteProductRequest(productID))
	rStatus := status.Convert(err)
	if rStatus != nil {
		log.Println(err)
		common.WriteError(w, http.StatusInternalServerError, "something went wrong with our server")
		return
	}

	common.WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{
		"message": "product was deleted successfully",
		"success": res.Success,
	})
}

type CartProduct struct {
	ProductID string `json:"product_id" validate:"required"`
	Name      string `json:"name" validate:"required"`
	Price     uint64 `json:"price" validate:"required"`
	Image     string `json:"image" validate:"required"`
	Quantity  uint32 `json:"quantity" validate:"required"`
}

func (h *Handler) HandleAddToCart(w http.ResponseWriter, r *http.Request) {
	cartID := r.PathValue("cartID")

	_, err := primitive.ObjectIDFromHex(cartID)
	if err != nil {
		log.Println(err)
		common.WriteError(w, http.StatusBadRequest, "invalid credentials")
		return
	}

	var item CartProduct
	err = common.ReadJSON(r, &item)
	if err != nil {
		log.Println(err)
		common.WriteError(w, http.StatusBadRequest, "invalid credentials")
		return
	}

	err = h.validate.Struct(item)
	if err != nil {
		log.Println(err)
		common.WriteError(w, http.StatusBadRequest, "invalid credentials")
		return
	}

	ctx := context.Background()
	_, err = h.gateway.AddToCart(ctx, toPBAddToCart(cartID, &item))
	rStatus := status.Convert(err)
	if rStatus != nil {
		log.Println(err)
		common.WriteError(w, http.StatusBadRequest, "something went wrong with our servers")
		return
	}

	common.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "successfully added the item to the cart",
	})
}

func (h *Handler) HandleRemoveFromCart(w http.ResponseWriter, r *http.Request) {
	productID := r.PathValue("productID")
	if _, err := primitive.ObjectIDFromHex(productID); err != nil {
		common.WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}

	cartID := r.PathValue("cartID")
	if _, err := primitive.ObjectIDFromHex(cartID); err != nil {
		log.Println(err)
		common.WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}

	ctx := context.Background()
	_, err := h.gateway.RemoveFromCart(ctx, toPBRemoveFromCart(productID, cartID))
	rStatus := status.Convert(err)
	if rStatus != nil {
		log.Println(err)
		common.WriteError(w, http.StatusInternalServerError, "something went wrong with our servers")
		return
	}

	common.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message": "successfully removed the item from the cart",
	})
}

func (h *Handler) HandleListCartItems(w http.ResponseWriter, r *http.Request) {
	cartID := r.PathValue("cartID")

	_, err := primitive.ObjectIDFromHex(cartID)
	if err != nil {
		log.Println(err)
		common.WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}

	ctx := context.Background()
	cart, err := h.gateway.ViewCart(ctx, toPBViewCart(cartID))
	rStatus := status.Convert(err)
	if rStatus != nil {
		log.Println(err)
		common.WriteError(w, http.StatusInternalServerError, "something went wrong with our servers")
		return
	}

	common.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"cart": cart,
	})
}

func (h *Handler) HandleCreateOrder(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) HandlePayments(w http.ResponseWriter, r *http.Request) {}
