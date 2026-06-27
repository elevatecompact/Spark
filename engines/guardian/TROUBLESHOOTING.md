# Troubleshooting

## Users cannot log in
1. Check rate limiting: redis-cli GET rate:login:user:userId.
2. Verify PostgreSQL reachable and user record exists.
3. Check if account is locked.
4. Validate password hashing cost not excessive.

## Token validation failures
1. Check iss claim matches jwt.issuer.
2. Verify token not in revocation list.
3. Ensure signing public key available.
4. Check clock skew < 30s.
