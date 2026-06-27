# commerce-service — Testing Guide
## Unit: Cart total calculation with variants, inventory deduction (optimistic lock test), order state machine, digital fulfillment URL generation, review rating enforcement.
## Integration: Full purchase flow (add to cart→checkout→payment→fulfill→download), inventory count accuracy under concurrent purchases, order cancellation and refund, merchant payout calculation.
## Load: 100 concurrent checkouts for limited-inventory product, 5000 products with complex variant trees.
