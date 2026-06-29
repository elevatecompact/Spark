package consumer

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/elevatecompact/spark/services/commerce-service/internal/domain"
	"github.com/elevatecompact/spark/services/commerce-service/internal/repository"
	"github.com/elevatecompact/spark/services/commerce-service/internal/service"
)

type CommerceServiceConsumer struct {
	repo *repository.CommerceRepository
	svc  *service.CommerceService
	evt  domain.EventProducer
}

func NewCommerceServiceConsumer(repo *repository.CommerceRepository, svc *service.CommerceService, evt domain.EventProducer) *CommerceServiceConsumer {
	return &CommerceServiceConsumer{repo: repo, svc: svc, evt: evt}
}

func (c *CommerceServiceConsumer) HandleWalletSettled(ctx context.Context, orderID uuid.UUID) error {
	o, err := c.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		return err
	}
	if o.Status != domain.OrderStatusPending {
		return nil
	}
	if err := c.repo.UpdateOrderStatus(ctx, orderID, domain.OrderStatusPaid); err != nil {
		return err
	}
	c.evt.Publish(ctx, "commerce.order.paid", map[string]interface{}{
		"orderId": orderID,
	})
	return nil
}

func (c *CommerceServiceConsumer) HandlePaymentRefunded(ctx context.Context, orderID uuid.UUID) error {
	return c.svc.RefundOrder(ctx, orderID)
}

func (c *CommerceServiceConsumer) HandleNotificationPushSent(ctx context.Context, orderID uuid.UUID) error {
	log.Info().Str("orderId", orderID.String()).Msg("delivery confirmation notification sent")
	return nil
}

func (c *CommerceServiceConsumer) HandleIdentityUserDeleted(ctx context.Context, userID uuid.UUID) error {
	orders, err := c.repo.ListOrders(ctx, &userID, nil, 100, 0)
	if err != nil {
		return err
	}
	for _, o := range orders {
		if o.Status == domain.OrderStatusPending || o.Status == domain.OrderStatusPaid {
			if err := c.svc.CancelOrder(ctx, o.ID); err != nil {
				log.Error().Err(err).Str("orderId", o.ID.String()).Msg("failed to cancel order on user delete")
			}
		}
	}
	return nil
}
