# Digital Wallet

The SPARK digital wallet manages platform currency for all users, enabling purchases, tips, gifts, and payouts within the ecosystem.

## Wallet Types

Viewer wallets hold platform credits purchased with real currency. Credits can be used for subscriptions, gifts, tips, and marketplace purchases. Creator wallets hold revenue earned from subscriptions, gifts, tips, and pay-per-view sales. Creator balances are denominated in real currency and are withdrawable. Platform wallets hold promotional credits granted by the platform for campaigns and rewards.

## Currency System

Platform credits are purchased at a fixed exchange rate of 1 credit per 1 cent. Credits are non-refundable except as required by law. Minimum purchase amounts are .99. Maximum wallet balance is ,000 for viewers and unlimited for creators. Credits expire after 365 days of inactivity.

## Transactions

The wallet service processes all credit transactions including purchases, transfers, gifts, and payouts. Each transaction is recorded with a unique ID, timestamp, amount, type, and reference entity. Transactions are immutable once committed. The wallet maintains a running balance with atomic updates.

## Security

Wallet operations require authentication and authorization verification. Large transactions require additional verification through email or SMS confirmation. Suspicious activity triggers automated account freezing and manual review. The wallet system is PCI-DSS compliant for credit card processing.

## Reporting

Users can view their transaction history, balance, and statements. Creators have access to earnings reports with revenue breakdowns by source. Export functionality supports CSV and PDF formats for accounting purposes.
