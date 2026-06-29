package repository

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/elevatecompact/spark/services/commerce-service/internal/domain"
)

type CommerceRepository struct {
	pool *pgxpool.Pool
}

func NewCommerceRepository(pool *pgxpool.Pool) *CommerceRepository {
	return &CommerceRepository{pool: pool}
}

func (r *CommerceRepository) CreateProduct(ctx context.Context, p *domain.Product) error {
	tags, _ := json.Marshal(p.Tags)
	media, _ := json.Marshal(p.MediaURLs)
	_, err := r.pool.Exec(ctx, `
		INSERT INTO products (id, creator_id, name, description, type, price_cents, currency, category, tags, media_urls, inventory_count, is_active, is_featured, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)
	`, p.ID, p.CreatorID, p.Name, p.Description, string(p.Type), p.PriceCents, p.Currency, p.Category, tags, media, p.Inventory, p.IsActive, p.IsFeatured, p.CreatedAt)
	return err
}

func (r *CommerceRepository) GetProduct(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, creator_id, name, description, type, price_cents, currency, category, tags, media_urls, inventory_count, is_active, is_featured, created_at
		FROM products WHERE id=$1
	`, id)
	p := &domain.Product{}
	var tags, media []byte
	err := row.Scan(&p.ID, &p.CreatorID, &p.Name, &p.Description, &p.Type, &p.PriceCents, &p.Currency, &p.Category, &tags, &media, &p.Inventory, &p.IsActive, &p.IsFeatured, &p.CreatedAt)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(tags, &p.Tags)
	json.Unmarshal(media, &p.MediaURLs)
	return p, nil
}

func (r *CommerceRepository) UpdateProduct(ctx context.Context, p *domain.Product) error {
	tags, _ := json.Marshal(p.Tags)
	media, _ := json.Marshal(p.MediaURLs)
	_, err := r.pool.Exec(ctx, `
		UPDATE products SET name=$2, description=$3, type=$4, price_cents=$5, currency=$6, category=$7, tags=$8, media_urls=$9, inventory_count=$10, is_active=$11
		WHERE id=$1
	`, p.ID, p.Name, p.Description, string(p.Type), p.PriceCents, p.Currency, p.Category, tags, media, p.Inventory, p.IsActive)
	return err
}

func (r *CommerceRepository) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM products WHERE id=$1`, id)
	return err
}

func (r *CommerceRepository) ListProducts(ctx context.Context, creatorID *uuid.UUID, category string, featured bool, limit, offset int) ([]domain.Product, error) {
	where := "WHERE is_active=true"
	args := []interface{}{}
	idx := 1
	if creatorID != nil {
		where += " AND creator_id=$" + string(rune('0'+idx))
		args = append(args, *creatorID)
		idx++
	}
	if category != "" {
		where += " AND category=$" + string(rune('0'+idx))
		args = append(args, category)
		idx++
	}
	if featured {
		where += " AND is_featured=true"
	}
	args = append(args, limit, offset)
	q := `SELECT id, creator_id, name, description, type, price_cents, currency, category, tags, media_urls, inventory_count, is_active, is_featured, created_at
	      FROM products ` + where + ` ORDER BY created_at DESC LIMIT $` + string(rune('0'+idx)) + ` OFFSET $` + string(rune('0'+idx+1))
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanProducts(rows)
}

func scanProducts(rows pgx.Rows) ([]domain.Product, error) {
	var res []domain.Product
	for rows.Next() {
		var p domain.Product
		var tags, media []byte
		if err := rows.Scan(&p.ID, &p.CreatorID, &p.Name, &p.Description, &p.Type, &p.PriceCents, &p.Currency, &p.Category, &tags, &media, &p.Inventory, &p.IsActive, &p.IsFeatured, &p.CreatedAt); err != nil {
			return nil, err
		}
		json.Unmarshal(tags, &p.Tags)
		json.Unmarshal(media, &p.MediaURLs)
		res = append(res, p)
	}
	return res, nil
}

func (r *CommerceRepository) CreateVariant(ctx context.Context, v *domain.ProductVariant) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO product_variants (id, product_id, name, price_cents, inventory_count, sort_order)
		VALUES ($1,$2,$3,$4,$5,$6)
	`, v.ID, v.ProductID, v.Name, v.PriceCents, v.InventoryCount, v.SortOrder)
	return err
}

func (r *CommerceRepository) GetVariants(ctx context.Context, productID uuid.UUID) ([]domain.ProductVariant, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, product_id, name, price_cents, inventory_count, sort_order FROM product_variants WHERE product_id=$1 ORDER BY sort_order`, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []domain.ProductVariant
	for rows.Next() {
		var v domain.ProductVariant
		if err := rows.Scan(&v.ID, &v.ProductID, &v.Name, &v.PriceCents, &v.InventoryCount, &v.SortOrder); err != nil {
			return nil, err
		}
		res = append(res, v)
	}
	return res, nil
}

func (r *CommerceRepository) AddCartItem(ctx context.Context, item *domain.CartItem) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO cart_items (id, user_id, product_id, variant_id, quantity, added_at)
		VALUES ($1,$2,$3,$4,$5,$6)
		ON CONFLICT (user_id, product_id, COALESCE(variant_id,'00000000-0000-0000-0000-000000000000')) DO UPDATE SET quantity=EXCLUDED.quantity, added_at=NOW()
	`, item.ID, item.UserID, item.ProductID, item.VariantID, item.Quantity, item.AddedAt)
	return err
}

func (r *CommerceRepository) GetCart(ctx context.Context, userID uuid.UUID) ([]domain.CartItem, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT ci.id, ci.user_id, ci.product_id, ci.variant_id, ci.quantity, ci.added_at FROM cart_items ci WHERE ci.user_id=$1 ORDER BY ci.added_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []domain.CartItem
	for rows.Next() {
		var ci domain.CartItem
		if err := rows.Scan(&ci.ID, &ci.UserID, &ci.ProductID, &ci.VariantID, &ci.Quantity, &ci.AddedAt); err != nil {
			return nil, err
		}
		res = append(res, ci)
	}
	return res, nil
}

func (r *CommerceRepository) UpdateCartItemQuantity(ctx context.Context, itemID uuid.UUID, quantity int) error {
	_, err := r.pool.Exec(ctx, `UPDATE cart_items SET quantity=$2 WHERE id=$1`, itemID, quantity)
	return err
}

func (r *CommerceRepository) RemoveCartItem(ctx context.Context, itemID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM cart_items WHERE id=$1`, itemID)
	return err
}

func (r *CommerceRepository) ClearCart(ctx context.Context, userID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM cart_items WHERE user_id=$1`, userID)
	return err
}

func (r *CommerceRepository) CreateOrder(ctx context.Context, o *domain.Order) error {
	addr, _ := json.Marshal(o.ShippingAddress)
	_, err := r.pool.Exec(ctx, `
		INSERT INTO orders (id, buyer_id, merchant_id, status, total_cents, currency, payment_intent_id, shipping_address, placed_at, fulfilled_at, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
	`, o.ID, o.BuyerID, o.MerchantID, string(o.Status), o.TotalCents, o.Currency, o.PaymentIntentID, addr, o.PlacedAt, o.FulfilledAt, o.CreatedAt)
	return err
}

func (r *CommerceRepository) CreateOrderItems(ctx context.Context, items []domain.OrderItem) error {
	for _, item := range items {
		_, err := r.pool.Exec(ctx, `
			INSERT INTO order_items (id, order_id, product_id, variant_id, quantity, unit_price_cents, fulfillment_status, download_url, fulfilled_at)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		`, item.ID, item.OrderID, item.ProductID, item.VariantID, item.Quantity, item.UnitPriceCents, string(item.FulfillmentStatus), item.DownloadURL, item.FulfilledAt)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *CommerceRepository) GetOrder(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, buyer_id, merchant_id, status, total_cents, currency, payment_intent_id, shipping_address, placed_at, fulfilled_at, created_at
		FROM orders WHERE id=$1
	`, id)
	o := &domain.Order{}
	var addr []byte
	err := row.Scan(&o.ID, &o.BuyerID, &o.MerchantID, &o.Status, &o.TotalCents, &o.Currency, &o.PaymentIntentID, &addr, &o.PlacedAt, &o.FulfilledAt, &o.CreatedAt)
	if err != nil {
		return nil, err
	}
	o.ShippingAddress = addr
	items, err := r.GetOrderItems(ctx, id)
	if err != nil {
		return nil, err
	}
	o.Items = items
	return o, nil
}

func (r *CommerceRepository) GetOrderItems(ctx context.Context, orderID uuid.UUID) ([]domain.OrderItem, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, order_id, product_id, variant_id, quantity, unit_price_cents, fulfillment_status, download_url, fulfilled_at
		FROM order_items WHERE order_id=$1
	`, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []domain.OrderItem
	for rows.Next() {
		var oi domain.OrderItem
		if err := rows.Scan(&oi.ID, &oi.OrderID, &oi.ProductID, &oi.VariantID, &oi.Quantity, &oi.UnitPriceCents, &oi.FulfillmentStatus, &oi.DownloadURL, &oi.FulfilledAt); err != nil {
			return nil, err
		}
		res = append(res, oi)
	}
	return res, nil
}

func (r *CommerceRepository) ListOrders(ctx context.Context, buyerID *uuid.UUID, merchantID *uuid.UUID, limit, offset int) ([]domain.Order, error) {
	where := "WHERE 1=1"
	args := []interface{}{}
	idx := 1
	if buyerID != nil {
		where += " AND buyer_id=$" + string(rune('0'+idx))
		args = append(args, *buyerID)
		idx++
	}
	if merchantID != nil {
		where += " AND merchant_id=$" + string(rune('0'+idx))
		args = append(args, *merchantID)
		idx++
	}
	args = append(args, limit, offset)
	q := `SELECT id, buyer_id, merchant_id, status, total_cents, currency, payment_intent_id, shipping_address, placed_at, fulfilled_at, created_at
	      FROM orders ` + where + ` ORDER BY created_at DESC LIMIT $` + string(rune('0'+idx)) + ` OFFSET $` + string(rune('0'+idx+1))
	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []domain.Order
	for rows.Next() {
		var o domain.Order
		var addr []byte
		if err := rows.Scan(&o.ID, &o.BuyerID, &o.MerchantID, &o.Status, &o.TotalCents, &o.Currency, &o.PaymentIntentID, &addr, &o.PlacedAt, &o.FulfilledAt, &o.CreatedAt); err != nil {
			return nil, err
		}
		o.ShippingAddress = addr
		res = append(res, o)
	}
	return res, nil
}

func (r *CommerceRepository) UpdateOrderStatus(ctx context.Context, id uuid.UUID, status domain.OrderStatus) error {
	_, err := r.pool.Exec(ctx, `UPDATE orders SET status=$2, fulfilled_at=CASE WHEN $2='fulfilled' THEN NOW() ELSE fulfilled_at END WHERE id=$1`, id, string(status))
	return err
}

func (r *CommerceRepository) FulfillOrderItem(ctx context.Context, itemID uuid.UUID, downloadURL string) error {
	_, err := r.pool.Exec(ctx, `UPDATE order_items SET fulfillment_status='fulfilled', download_url=$2, fulfilled_at=NOW() WHERE id=$1`, itemID, downloadURL)
	return err
}

func (r *CommerceRepository) CreateReview(ctx context.Context, rev *domain.Review) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO reviews (id, product_id, user_id, rating, title, body, is_verified_purchase, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	`, rev.ID, rev.ProductID, rev.UserID, rev.Rating, rev.Title, rev.Body, rev.IsVerifiedPurchase, rev.CreatedAt)
	return err
}

func (r *CommerceRepository) ListReviews(ctx context.Context, productID uuid.UUID, limit, offset int) ([]domain.Review, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, product_id, user_id, rating, title, body, is_verified_purchase, created_at
		FROM reviews WHERE product_id=$1 ORDER BY created_at DESC LIMIT $2 OFFSET $3
	`, productID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []domain.Review
	for rows.Next() {
		var rv domain.Review
		if err := rows.Scan(&rv.ID, &rv.ProductID, &rv.UserID, &rv.Rating, &rv.Title, &rv.Body, &rv.IsVerifiedPurchase, &rv.CreatedAt); err != nil {
			return nil, err
		}
		res = append(res, rv)
	}
	return res, nil
}

func (r *CommerceRepository) GetMerchantDashboard(ctx context.Context, merchantID uuid.UUID) (*domain.MerchantDashboard, error) {
	var d domain.MerchantDashboard
	err := r.pool.QueryRow(ctx, `
		SELECT COUNT(*) as product_count FROM products WHERE creator_id=$1 AND is_active=true
	`, merchantID).Scan(&d.ProductCount)
	if err != nil {
		return nil, err
	}
	err = r.pool.QueryRow(ctx, `
		SELECT COALESCE(SUM(total_cents),0), COUNT(*) FROM orders WHERE merchant_id=$1 AND status IN ('paid','fulfilled')
	`, merchantID).Scan(&d.TotalRevenue, &d.OrderCount)
	if err != nil {
		return nil, err
	}
	d.TotalSales = d.OrderCount
	return &d, nil
}

func (r *CommerceRepository) GetPayouts(ctx context.Context, merchantID uuid.UUID) ([]domain.Payout, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, merchant_id, amount_cents, currency, status, period_start, period_end, paid_at, created_at
		FROM payouts WHERE merchant_id=$1 ORDER BY created_at DESC
	`, merchantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []domain.Payout
	for rows.Next() {
		var p domain.Payout
		if err := rows.Scan(&p.ID, &p.MerchantID, &p.AmountCents, &p.Currency, &p.Status, &p.PeriodStart, &p.PeriodEnd, &p.PaidAt, &p.CreatedAt); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, nil
}

func (r *CommerceRepository) GetAdminRevenue(ctx context.Context) (*domain.AdminRevenue, error) {
	var d domain.AdminRevenue
	err := r.pool.QueryRow(ctx, `
		SELECT COALESCE(SUM(total_cents),0), COUNT(*) FROM orders WHERE status IN ('paid','fulfilled')
	`).Scan(&d.TotalRevenueCents, &d.OrderCount)
	if err != nil {
		return nil, err
	}
	d.PlatformFeeCents = d.TotalRevenueCents * 10 / 100
	err = r.pool.QueryRow(ctx, `
		SELECT COALESCE(SUM(amount_cents),0) FROM payouts WHERE status='paid'
	`).Scan(&d.MerchantPayoutTotal)
	if err != nil {
		return nil, err
	}
	err = r.pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM orders WHERE status='refunded'
	`).Scan(&d.RefundCount)
	return &d, err
}

func (r *CommerceRepository) FeatureProduct(ctx context.Context, id uuid.UUID, featured bool) error {
	_, err := r.pool.Exec(ctx, `UPDATE products SET is_featured=$2 WHERE id=$1`, id, featured)
	return err
}

func (r *CommerceRepository) ProductExists(ctx context.Context, id uuid.UUID) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM products WHERE id=$1)`, id).Scan(&exists)
	return exists, err
}

func (r *CommerceRepository) IsVerifiedPurchase(ctx context.Context, productID, userID uuid.UUID) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM order_items oi JOIN orders o ON oi.order_id=o.id WHERE oi.product_id=$1 AND o.buyer_id=$2 AND o.status IN ('paid','fulfilled'))`, productID, userID).Scan(&exists)
	return exists, err
}

func (r *CommerceRepository) GetCartItemProduct(ctx context.Context, itemID uuid.UUID) (*uuid.UUID, error) {
	var productID uuid.UUID
	err := r.pool.QueryRow(ctx, `SELECT product_id FROM cart_items WHERE id=$1`, itemID).Scan(&productID)
	if err != nil {
		return nil, err
	}
	return &productID, nil
}

func (r *CommerceRepository) DeductInventory(ctx context.Context, productID uuid.UUID, qty int) error {
	_, err := r.pool.Exec(ctx, `UPDATE products SET inventory_count=inventory_count-$2 WHERE id=$1 AND inventory_count>0`, productID, qty)
	return err
}

func (r *CommerceRepository) GetMerchantProducts(ctx context.Context, merchantID uuid.UUID) ([]domain.Product, error) {
	return r.ListProducts(ctx, &merchantID, "", false, 1000, 0)
}

func (r *CommerceRepository) GetCartCount(ctx context.Context, userID uuid.UUID) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM cart_items WHERE user_id=$1`, userID).Scan(&count)
	return count, err
}

func (r *CommerceRepository) CreatePayout(ctx context.Context, p *domain.Payout) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO payouts (id, merchant_id, amount_cents, currency, status, period_start, period_end, paid_at, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
	`, p.ID, p.MerchantID, p.AmountCents, p.Currency, p.Status, p.PeriodStart, p.PeriodEnd, p.PaidAt, p.CreatedAt)
	return err
}

func (r *CommerceRepository) GetOrderItem(ctx context.Context, itemID uuid.UUID) (*domain.OrderItem, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, order_id, product_id, variant_id, quantity, unit_price_cents, fulfillment_status, download_url, fulfilled_at
		FROM order_items WHERE id=$1
	`, itemID)
	var oi domain.OrderItem
	err := row.Scan(&oi.ID, &oi.OrderID, &oi.ProductID, &oi.VariantID, &oi.Quantity, &oi.UnitPriceCents, &oi.FulfillmentStatus, &oi.DownloadURL, &oi.FulfilledAt)
	if err != nil {
		return nil, err
	}
	return &oi, nil
}

func (r *CommerceRepository) GetDownloads(ctx context.Context, orderID uuid.UUID) ([]domain.OrderItem, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, order_id, product_id, variant_id, quantity, unit_price_cents, fulfillment_status, download_url, fulfilled_at
		FROM order_items WHERE order_id=$1 AND download_url!='' AND fulfillment_status='fulfilled'
	`, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []domain.OrderItem
	for rows.Next() {
		var oi domain.OrderItem
		if err := rows.Scan(&oi.ID, &oi.OrderID, &oi.ProductID, &oi.VariantID, &oi.Quantity, &oi.UnitPriceCents, &oi.FulfillmentStatus, &oi.DownloadURL, &oi.FulfilledAt); err != nil {
			return nil, err
		}
		res = append(res, oi)
	}
	return res, nil
}

func (r *CommerceRepository) GetOrderByID(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	return r.GetOrder(ctx, id)
}
