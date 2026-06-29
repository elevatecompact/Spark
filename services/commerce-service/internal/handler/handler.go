package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/elevatecompact/spark/services/commerce-service/internal/domain"
	"github.com/elevatecompact/spark/services/commerce-service/internal/service"
)

type Handler struct {
	svc *service.CommerceService
}

func New(svc *service.CommerceService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(r chi.Router) {
	r.Route("/v1/products", func(r chi.Router) {
		r.Post("/", h.createProduct)
		r.Get("/", h.listProducts)
		r.Get("/{id}", h.getProduct)
		r.Patch("/{id}", h.updateProduct)
		r.Delete("/{id}", h.deleteProduct)
		r.Post("/{id}/variants", h.createVariant)
		r.Get("/{id}/reviews", h.listReviews)
		r.Post("/{id}/reviews", h.createReview)
	})
	r.Route("/v1/cart", func(r chi.Router) {
		r.Get("/", h.getCart)
		r.Post("/items", h.addCartItem)
		r.Patch("/items/{id}", h.updateCartItem)
		r.Delete("/items/{id}", h.removeCartItem)
		r.Delete("/", h.clearCart)
	})
	r.Post("/v1/checkout", h.checkout)
	r.Get("/v1/orders/{id}", h.getOrder)
	r.Get("/v1/orders", h.listOrders)
	r.Post("/v1/orders/{id}/cancel", h.cancelOrder)
	r.Post("/v1/orders/{id}/fulfill", h.fulfillOrder)
	r.Get("/v1/orders/{id}/downloads", h.getDownloads)
	r.Post("/v1/fulfillment/retry", h.retryFulfillment)
	r.Route("/v1/merchant", func(r chi.Router) {
		r.Get("/dashboard", h.merchantDashboard)
		r.Get("/payouts", h.merchantPayouts)
		r.Get("/products", h.merchantProducts)
		r.Post("/storefront", h.configureStorefront)
	})
	r.Route("/v1/admin", func(r chi.Router) {
		r.Post("/products/{id}/feature", h.featureProduct)
		r.Post("/orders/{id}/refund", h.refundOrder)
		r.Get("/revenue", h.adminRevenue)
	})
}

func (h *Handler) createProduct(w http.ResponseWriter, r *http.Request) {
	var p domain.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	result, err := h.svc.CreateProduct(r.Context(), &p)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, result)
}

func (h *Handler) getProduct(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))
	p, err := h.svc.GetProduct(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "product not found")
		return
	}
	writeJSON(w, http.StatusOK, p)
}

func (h *Handler) updateProduct(w http.ResponseWriter, r *http.Request) {
	var p domain.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	p.ID = uuid.MustParse(chi.URLParam(r, "id"))
	if err := h.svc.UpdateProduct(r.Context(), &p); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, p)
}

func (h *Handler) deleteProduct(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))
	if err := h.svc.DeleteProduct(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusNoContent, nil)
}

func (h *Handler) listProducts(w http.ResponseWriter, r *http.Request) {
	var creatorID *uuid.UUID
	if cid := r.URL.Query().Get("creatorId"); cid != "" {
		id := uuid.MustParse(cid)
		creatorID = &id
	}
	category := r.URL.Query().Get("category")
	featured := r.URL.Query().Get("featured") == "true"
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	products, err := h.svc.ListProducts(r.Context(), creatorID, category, featured, limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, products)
}

func (h *Handler) createVariant(w http.ResponseWriter, r *http.Request) {
	productID := uuid.MustParse(chi.URLParam(r, "id"))
	var v domain.ProductVariant
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	v.ProductID = productID
	if err := h.svc.CreateVariant(r.Context(), &v); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, v)
}

func (h *Handler) getCart(w http.ResponseWriter, r *http.Request) {
	userID := extractUserID(r)
	cart, err := h.svc.GetCart(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, cart)
}

func (h *Handler) addCartItem(w http.ResponseWriter, r *http.Request) {
	userID := extractUserID(r)
	var req struct {
		ProductID uuid.UUID  `json:"productId"`
		VariantID *uuid.UUID `json:"variantId,omitempty"`
		Quantity  int        `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	item, err := h.svc.AddCartItem(r.Context(), userID, req.ProductID, req.VariantID, req.Quantity)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (h *Handler) updateCartItem(w http.ResponseWriter, r *http.Request) {
	itemID := uuid.MustParse(chi.URLParam(r, "id"))
	var req struct{ Quantity int `json:"quantity"` }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.svc.UpdateCartItemQuantity(r.Context(), itemID, req.Quantity); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

func (h *Handler) removeCartItem(w http.ResponseWriter, r *http.Request) {
	itemID := uuid.MustParse(chi.URLParam(r, "id"))
	if err := h.svc.RemoveCartItem(r.Context(), itemID); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusNoContent, nil)
}

func (h *Handler) clearCart(w http.ResponseWriter, r *http.Request) {
	userID := extractUserID(r)
	if err := h.svc.ClearCart(r.Context(), userID); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusNoContent, nil)
}

func (h *Handler) checkout(w http.ResponseWriter, r *http.Request) {
	var input domain.CreateOrderInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	order, err := h.svc.Checkout(r.Context(), input)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, order)
}

func (h *Handler) getOrder(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))
	o, err := h.svc.GetOrder(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "order not found")
		return
	}
	writeJSON(w, http.StatusOK, o)
}

func (h *Handler) listOrders(w http.ResponseWriter, r *http.Request) {
	var buyerID, merchantID *uuid.UUID
	if bid := r.URL.Query().Get("buyerId"); bid != "" {
		id := uuid.MustParse(bid)
		buyerID = &id
	}
	if mid := r.URL.Query().Get("merchantId"); mid != "" {
		id := uuid.MustParse(mid)
		merchantID = &id
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	orders, err := h.svc.ListOrders(r.Context(), buyerID, merchantID, limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, orders)
}

func (h *Handler) cancelOrder(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))
	if err := h.svc.CancelOrder(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "cancelled"})
}

func (h *Handler) fulfillOrder(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))
	if err := h.svc.FulfillOrder(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "fulfilled"})
}

func (h *Handler) getDownloads(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))
	items, err := h.svc.GetDownloads(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h *Handler) retryFulfillment(w http.ResponseWriter, r *http.Request) {
	var req struct{ OrderItemID uuid.UUID `json:"orderItemId"` }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.svc.RetryFulfillment(r.Context(), req.OrderItemID); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "retried"})
}

func (h *Handler) createReview(w http.ResponseWriter, r *http.Request) {
	productID := uuid.MustParse(chi.URLParam(r, "id"))
	userID := extractUserID(r)
	var req struct {
		Rating int    `json:"rating"`
		Title  string `json:"title"`
		Body   string `json:"body"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	rev, err := h.svc.CreateReview(r.Context(), productID, userID, req.Rating, req.Title, req.Body)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, rev)
}

func (h *Handler) listReviews(w http.ResponseWriter, r *http.Request) {
	productID := uuid.MustParse(chi.URLParam(r, "id"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	reviews, err := h.svc.ListReviews(r.Context(), productID, limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, reviews)
}

func (h *Handler) merchantDashboard(w http.ResponseWriter, r *http.Request) {
	merchantID := uuid.MustParse(r.URL.Query().Get("merchantId"))
	d, err := h.svc.GetMerchantDashboard(r.Context(), merchantID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, d)
}

func (h *Handler) merchantPayouts(w http.ResponseWriter, r *http.Request) {
	merchantID := uuid.MustParse(r.URL.Query().Get("merchantId"))
	payouts, err := h.svc.GetPayouts(r.Context(), merchantID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, payouts)
}

func (h *Handler) merchantProducts(w http.ResponseWriter, r *http.Request) {
	merchantID := uuid.MustParse(r.URL.Query().Get("merchantId"))
	products, err := h.svc.GetMerchantProducts(r.Context(), merchantID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, products)
}

func (h *Handler) configureStorefront(w http.ResponseWriter, r *http.Request) {
	merchantID := uuid.MustParse(r.URL.Query().Get("merchantId"))
	var config map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.svc.ConfigureStorefront(r.Context(), merchantID, config); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "configured"})
}

func (h *Handler) featureProduct(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))
	var req struct{ Featured bool `json:"featured"` }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.svc.FeatureProduct(r.Context(), id, req.Featured); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "updated"})
}

func (h *Handler) refundOrder(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(chi.URLParam(r, "id"))
	if err := h.svc.RefundOrder(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "refunded"})
}

func (h *Handler) adminRevenue(w http.ResponseWriter, r *http.Request) {
	rev, err := h.svc.GetAdminRevenue(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, rev)
}

func extractUserID(r *http.Request) uuid.UUID {
	uid := r.Header.Get("X-User-ID")
	if uid != "" {
		return uuid.MustParse(uid)
	}
	return uuid.New()
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if v != nil {
		json.NewEncoder(w).Encode(v)
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	log.Error().Int("status", status).Str("msg", msg).Msg("handler error")
	writeJSON(w, status, map[string]string{"error": msg})
}
