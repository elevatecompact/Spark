package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/elevatecompact/spark/services/commerce-service/internal/domain"
	"github.com/elevatecompact/spark/services/commerce-service/internal/repository"
)

type CommerceService struct {
	repo *repository.CommerceRepository
	evt  domain.EventProducer
}

func NewCommerceService(repo *repository.CommerceRepository, evt domain.EventProducer) *CommerceService {
	return &CommerceService{repo: repo, evt: evt}
}

func (s *CommerceService) CreateProduct(ctx context.Context, p *domain.Product) (*domain.Product, error) {
	p.ID = uuid.New()
	p.CreatedAt = time.Now()
	p.IsFeatured = false
	if p.Inventory == 0 {
		p.Inventory = -1
	}
	if err := s.repo.CreateProduct(ctx, p); err != nil {
		return nil, fmt.Errorf("create product: %w", err)
	}
	s.evt.Publish(ctx, "commerce.product.created", map[string]interface{}{
		"productId": p.ID, "creatorId": p.CreatorID, "type": p.Type,
	})
	return p, nil
}

func (s *CommerceService) GetProduct(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	p, err := s.repo.GetProduct(ctx, id)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (s *CommerceService) UpdateProduct(ctx context.Context, p *domain.Product) error {
	if err := s.repo.UpdateProduct(ctx, p); err != nil {
		return err
	}
	s.evt.Publish(ctx, "commerce.product.updated", map[string]interface{}{
		"productId": p.ID, "creatorId": p.CreatorID,
	})
	return nil
}

func (s *CommerceService) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.DeleteProduct(ctx, id); err != nil {
		return err
	}
	s.evt.Publish(ctx, "commerce.product.deleted", map[string]interface{}{
		"productId": id,
	})
	return nil
}

func (s *CommerceService) ListProducts(ctx context.Context, creatorID *uuid.UUID, category string, featured bool, limit, offset int) ([]domain.Product, error) {
	return s.repo.ListProducts(ctx, creatorID, category, featured, limit, offset)
}

func (s *CommerceService) CreateVariant(ctx context.Context, v *domain.ProductVariant) error {
	v.ID = uuid.New()
	return s.repo.CreateVariant(ctx, v)
}

func (s *CommerceService) GetVariants(ctx context.Context, productID uuid.UUID) ([]domain.ProductVariant, error) {
	return s.repo.GetVariants(ctx, productID)
}

func (s *CommerceService) AddCartItem(ctx context.Context, userID uuid.UUID, productID uuid.UUID, variantID *uuid.UUID, quantity int) (*domain.CartItem, error) {
	if quantity < 1 {
		return nil, errors.New("quantity must be at least 1")
	}
	item := &domain.CartItem{
		ID:        uuid.New(),
		UserID:    userID,
		ProductID: productID,
		VariantID: variantID,
		Quantity:  quantity,
		AddedAt:   time.Now(),
	}
	if err := s.repo.AddCartItem(ctx, item); err != nil {
		return nil, err
	}
	s.evt.Publish(ctx, "commerce.cart.updated", map[string]interface{}{
		"userId": userID, "itemCount": quantity,
	})
	return item, nil
}

func (s *CommerceService) GetCart(ctx context.Context, userID uuid.UUID) ([]domain.CartItem, error) {
	return s.repo.GetCart(ctx, userID)
}

func (s *CommerceService) UpdateCartItemQuantity(ctx context.Context, itemID uuid.UUID, quantity int) error {
	return s.repo.UpdateCartItemQuantity(ctx, itemID, quantity)
}

func (s *CommerceService) RemoveCartItem(ctx context.Context, itemID uuid.UUID) error {
	productID, err := s.repo.GetCartItemProduct(ctx, itemID)
	if err != nil {
		return err
	}
	if err := s.repo.RemoveCartItem(ctx, itemID); err != nil {
		return err
	}
	_ = productID
	return nil
}

func (s *CommerceService) ClearCart(ctx context.Context, userID uuid.UUID) error {
	return s.repo.ClearCart(ctx, userID)
}

func (s *CommerceService) Checkout(ctx context.Context, input domain.CreateOrderInput) (*domain.Order, error) {
	if len(input.Items) == 0 {
		return nil, errors.New("order must have at least one item")
	}
	order := &domain.Order{
		ID:              uuid.New(),
		BuyerID:         input.BuyerID,
		MerchantID:      input.MerchantID,
		Status:          domain.OrderStatusPending,
		Currency:        input.Currency,
		ShippingAddress: input.ShippingAddress,
		CreatedAt:       time.Now(),
	}
	var total int64
	for _, item := range input.Items {
		total += item.UnitPriceCents * int64(item.Quantity)
		if item.Quantity < 1 {
			return nil, errors.New("item quantity must be at least 1")
		}
		if err := s.repo.DeductInventory(ctx, item.ProductID, item.Quantity); err != nil {
			return nil, fmt.Errorf("inventory deduction: %w", err)
		}
	}
	order.TotalCents = total
	if err := s.repo.CreateOrder(ctx, order); err != nil {
		return nil, err
	}
	var orderItems []domain.OrderItem
	for _, in := range input.Items {
		oi := domain.OrderItem{
			ID:                uuid.New(),
			OrderID:           order.ID,
			ProductID:         in.ProductID,
			VariantID:         in.VariantID,
			Quantity:          in.Quantity,
			UnitPriceCents:    in.UnitPriceCents,
			FulfillmentStatus: domain.FulfillmentPending,
		}
		orderItems = append(orderItems, oi)
	}
	if err := s.repo.CreateOrderItems(ctx, orderItems); err != nil {
		return nil, err
	}
	order.Items = orderItems
	s.evt.Publish(ctx, "commerce.order.placed", map[string]interface{}{
		"orderId": order.ID, "buyerId": order.BuyerID, "merchantId": order.MerchantID,
		"totalCents": order.TotalCents, "currency": order.Currency, "status": order.Status,
		"placedAt": order.CreatedAt,
	})
	return order, nil
}

func (s *CommerceService) GetOrder(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	o, err := s.repo.GetOrder(ctx, id)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (s *CommerceService) ListOrders(ctx context.Context, buyerID *uuid.UUID, merchantID *uuid.UUID, limit, offset int) ([]domain.Order, error) {
	return s.repo.ListOrders(ctx, buyerID, merchantID, limit, offset)
}

func (s *CommerceService) CancelOrder(ctx context.Context, id uuid.UUID) error {
	o, err := s.repo.GetOrder(ctx, id)
	if err != nil {
		return err
	}
	if o.Status != domain.OrderStatusPending && o.Status != domain.OrderStatusPaid {
		return errors.New("order cannot be cancelled")
	}
	if err := s.repo.UpdateOrderStatus(ctx, id, domain.OrderStatusCancelled); err != nil {
		return err
	}
	s.evt.Publish(ctx, "commerce.order.cancelled", map[string]interface{}{
		"orderId": id, "buyerId": o.BuyerID,
	})
	return nil
}

func (s *CommerceService) FulfillOrder(ctx context.Context, orderID uuid.UUID) error {
	o, err := s.repo.GetOrder(ctx, orderID)
	if err != nil {
		return err
	}
	if o.Status != domain.OrderStatusPaid {
		return errors.New("order must be paid before fulfillment")
	}
	if err := s.repo.UpdateOrderStatus(ctx, orderID, domain.OrderStatusFulfilled); err != nil {
		return err
	}
	for _, item := range o.Items {
		url := fmt.Sprintf("https://dl.spark.dev/orders/%s/items/%s", orderID, item.ID)
		if err := s.repo.FulfillOrderItem(ctx, item.ID, url); err != nil {
			return err
		}
	}
	s.evt.Publish(ctx, "commerce.order.fulfilled", map[string]interface{}{
		"orderId": orderID, "merchantId": o.MerchantID,
	})
	return nil
}

func (s *CommerceService) RefundOrder(ctx context.Context, orderID uuid.UUID) error {
	o, err := s.repo.GetOrder(ctx, orderID)
	if err != nil {
		return err
	}
	if o.Status == domain.OrderStatusRefunded {
		return errors.New("order already refunded")
	}
	if err := s.repo.UpdateOrderStatus(ctx, orderID, domain.OrderStatusRefunded); err != nil {
		return err
	}
	s.evt.Publish(ctx, "commerce.order.refunded", map[string]interface{}{
		"orderId": orderID, "buyerId": o.BuyerID,
	})
	return nil
}

func (s *CommerceService) GetDownloads(ctx context.Context, orderID uuid.UUID) ([]domain.OrderItem, error) {
	return s.repo.GetDownloads(ctx, orderID)
}

func (s *CommerceService) RetryFulfillment(ctx context.Context, itemID uuid.UUID) error {
	oi, err := s.repo.GetOrderItem(ctx, itemID)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("https://dl.spark.dev/orders/%s/items/%s", oi.OrderID, oi.ID)
	return s.repo.FulfillOrderItem(ctx, itemID, url)
}

func (s *CommerceService) CreateReview(ctx context.Context, productID, userID uuid.UUID, rating int, title, body string) (*domain.Review, error) {
	if rating < 1 || rating > 5 {
		return nil, errors.New("rating must be 1-5")
	}
	verified, _ := s.repo.IsVerifiedPurchase(ctx, productID, userID)
	rev := &domain.Review{
		ID:                uuid.New(),
		ProductID:         productID,
		UserID:            userID,
		Rating:            rating,
		Title:             title,
		Body:              body,
		IsVerifiedPurchase: verified,
		CreatedAt:         time.Now(),
	}
	if err := s.repo.CreateReview(ctx, rev); err != nil {
		return nil, err
	}
	s.evt.Publish(ctx, "commerce.review.submitted", map[string]interface{}{
		"reviewId": rev.ID, "productId": productID, "userId": userID, "rating": rating,
	})
	return rev, nil
}

func (s *CommerceService) ListReviews(ctx context.Context, productID uuid.UUID, limit, offset int) ([]domain.Review, error) {
	return s.repo.ListReviews(ctx, productID, limit, offset)
}

func (s *CommerceService) GetMerchantDashboard(ctx context.Context, merchantID uuid.UUID) (*domain.MerchantDashboard, error) {
	return s.repo.GetMerchantDashboard(ctx, merchantID)
}

func (s *CommerceService) GetPayouts(ctx context.Context, merchantID uuid.UUID) ([]domain.Payout, error) {
	return s.repo.GetPayouts(ctx, merchantID)
}

func (s *CommerceService) GetMerchantProducts(ctx context.Context, merchantID uuid.UUID) ([]domain.Product, error) {
	return s.repo.GetMerchantProducts(ctx, merchantID)
}

func (s *CommerceService) FeatureProduct(ctx context.Context, id uuid.UUID, featured bool) error {
	return s.repo.FeatureProduct(ctx, id, featured)
}

func (s *CommerceService) GetAdminRevenue(ctx context.Context) (*domain.AdminRevenue, error) {
	return s.repo.GetAdminRevenue(ctx)
}

func (s *CommerceService) ConfigureStorefront(ctx context.Context, merchantID uuid.UUID, config map[string]interface{}) error {
	log.Info().Interface("config", config).Str("merchantId", merchantID.String()).Msg("storefront configured")
	return nil
}
