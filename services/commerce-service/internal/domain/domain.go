package domain

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type EventProducer interface {
	Publish(ctx context.Context, eventType string, data interface{})
	Close()
}

type ProductType string

const (
	ProductTypeDigital  ProductType = "digital"
	ProductTypePhysical ProductType = "physical"
	ProductTypeDownload ProductType = "download"
	ProductTypeBundle   ProductType = "bundle"
)

type Product struct {
	ID            uuid.UUID       `json:"id"`
	CreatorID     uuid.UUID       `json:"creatorId"`
	Name          string          `json:"name"`
	Description   string          `json:"description"`
	Type          ProductType     `json:"type"`
	PriceCents    int64           `json:"priceCents"`
	Currency      string          `json:"currency"`
	Category      string          `json:"category"`
	Tags          []string        `json:"tags"`
	MediaURLs     []string        `json:"mediaUrls"`
	Inventory     int64           `json:"inventoryCount"`
	IsActive      bool            `json:"isActive"`
	IsFeatured    bool            `json:"isFeatured"`
	CreatedAt     time.Time       `json:"createdAt"`
}

type ProductVariant struct {
	ID             uuid.UUID `json:"id"`
	ProductID      uuid.UUID `json:"productId"`
	Name           string    `json:"name"`
	PriceCents     *int64    `json:"priceCents,omitempty"`
	InventoryCount int64     `json:"inventoryCount"`
	SortOrder      int       `json:"sortOrder"`
}

type CartItem struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"userId"`
	ProductID uuid.UUID `json:"productId"`
	VariantID *uuid.UUID `json:"variantId,omitempty"`
	Quantity  int       `json:"quantity"`
	AddedAt   time.Time `json:"addedAt"`
	Product   *Product  `json:"product,omitempty"`
	Variant   *ProductVariant `json:"variant,omitempty"`
}

type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusPaid       OrderStatus = "paid"
	OrderStatusFulfilled  OrderStatus = "fulfilled"
	OrderStatusCancelled  OrderStatus = "cancelled"
	OrderStatusRefunded   OrderStatus = "refunded"
)

type Order struct {
	ID               uuid.UUID       `json:"id"`
	BuyerID          uuid.UUID       `json:"buyerId"`
	MerchantID       uuid.UUID       `json:"merchantId"`
	Status           OrderStatus     `json:"status"`
	TotalCents       int64           `json:"totalCents"`
	Currency         string          `json:"currency"`
	PaymentIntentID  *uuid.UUID      `json:"paymentIntentId,omitempty"`
	ShippingAddress  json.RawMessage `json:"shippingAddress,omitempty"`
	PlacedAt         *time.Time      `json:"placedAt,omitempty"`
	FulfilledAt      *time.Time      `json:"fulfilledAt,omitempty"`
	Items            []OrderItem     `json:"items,omitempty"`
	CreatedAt        time.Time       `json:"createdAt"`
}

type FulfillmentStatus string

const (
	FulfillmentPending    FulfillmentStatus = "pending"
	FulfillmentFulfilled  FulfillmentStatus = "fulfilled"
)

type OrderItem struct {
	ID                uuid.UUID         `json:"id"`
	OrderID           uuid.UUID         `json:"orderId"`
	ProductID         uuid.UUID         `json:"productId"`
	VariantID         *uuid.UUID        `json:"variantId,omitempty"`
	Quantity          int               `json:"quantity"`
	UnitPriceCents    int64             `json:"unitPriceCents"`
	FulfillmentStatus FulfillmentStatus `json:"fulfillmentStatus"`
	DownloadURL       string            `json:"downloadUrl,omitempty"`
	FulfilledAt       *time.Time        `json:"fulfilledAt,omitempty"`
}

type Review struct {
	ID                uuid.UUID `json:"id"`
	ProductID         uuid.UUID `json:"productId"`
	UserID            uuid.UUID `json:"userId"`
	Rating            int       `json:"rating"`
	Title             string    `json:"title"`
	Body              string    `json:"body"`
	IsVerifiedPurchase bool     `json:"isVerifiedPurchase"`
	CreatedAt         time.Time `json:"createdAt"`
}

type Payout struct {
	ID          uuid.UUID `json:"id"`
	MerchantID  uuid.UUID `json:"merchantId"`
	AmountCents int64     `json:"amountCents"`
	Currency    string    `json:"currency"`
	Status      string    `json:"status"`
	PeriodStart time.Time `json:"periodStart"`
	PeriodEnd   time.Time `json:"periodEnd"`
	PaidAt      *time.Time `json:"paidAt,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
}

type MerchantDashboard struct {
	TotalSales    int64 `json:"totalSales"`
	TotalRevenue  int64 `json:"totalRevenueCents"`
	PendingPayout int64 `json:"pendingPayoutCents"`
	ProductCount  int64 `json:"productCount"`
	OrderCount    int64 `json:"orderCount"`
}

type AdminRevenue struct {
	TotalRevenueCents   int64 `json:"totalRevenueCents"`
	PlatformFeeCents    int64 `json:"platformFeeCents"`
	MerchantPayoutTotal int64 `json:"merchantPayoutTotalCents"`
	OrderCount          int64 `json:"orderCount"`
	RefundCount         int64 `json:"refundCount"`
}

type CreateOrderInput struct {
	BuyerID         uuid.UUID       `json:"buyerId"`
	MerchantID      uuid.UUID       `json:"merchantId"`
	Items           []OrderItemInput `json:"items"`
	Currency        string          `json:"currency"`
	ShippingAddress json.RawMessage `json:"shippingAddress,omitempty"`
}

type OrderItemInput struct {
	ProductID      uuid.UUID  `json:"productId"`
	VariantID      *uuid.UUID `json:"variantId,omitempty"`
	Quantity       int        `json:"quantity"`
	UnitPriceCents int64      `json:"unitPriceCents"`
}
