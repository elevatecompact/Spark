# Subscriptions

Subscriptions provide a recurring revenue model for creators and predictable access for viewers. SPARK's subscription system supports multiple tiers, billing cycles, and promotional pricing.

## Subscription Tiers

Creators can define up to three subscription tiers. The basic tier typically offers ad-free viewing and a badge. The mid tier adds exclusive content and early access to new videos. The premium tier includes all previous benefits plus direct messaging, custom emotes, and priority support. Each tier has a minimum and maximum price set by the platform.

## Billing and Renewal

Subscriptions bill on a recurring monthly or annual cycle. Monthly subscriptions renew automatically on the same day each month. Annual subscriptions offer a discount equivalent to two months free compared to monthly billing. Failed payment retries follow an exponential backoff schedule over 5 days. After 5 days of failed retries, the subscription enters a 7-day grace period before cancellation.

## Subscription Management

Viewers can manage subscriptions through their account settings including upgrading, downgrading, canceling, and reactivating. Prorated credits are applied for tier upgrades. Cancelled subscriptions remain active until the end of the current billing period. Subscription history is available for both viewers and creators.

## Creator Payouts

Subscription revenue accrues to the creator's wallet. Payouts are processed on a monthly basis with a 30-day hold period to account for refunds and chargebacks. Revenue reports show per-subscriber metrics, churn rates, and revenue trends.
