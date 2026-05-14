# Damaged Product Management - Quick Reference Guide

## TL;DR - For Developers

### Add a Product/Variant as Damaged

```bash
curl -X PATCH http://127.0.0.1:8088/api/stock/mark-damaged \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": "prod_8b35c2b9",
    "variant_sku": "SL-001-PLAS",
    "quantity": 10,
    "reason": "defective_batch"
  }'
```

### Get All Damaged Products

```bash
curl -X GET "http://127.0.0.1:8088/api/stock/damaged?limit=50&offset=0" \
  -H "Authorization: Bearer {token}"
```

### Check Stock Summary (with Damaged Items)

```bash
curl -X GET "http://127.0.0.1:8088/api/stock/summary" \
  -H "Authorization: Bearer {token}"
```

---

## Key Concepts

| Term | Description | Example |
|------|-------------|---------|
| **Current Stock** | Always: purchased - sold | 116,105 units |
| **Available Stock** | current - reserved - damaged | 116,095 units |
| **Damaged Stock** | Marked as defective/broken | 10 units |
| **Damage Reason** | Why marked as damaged | defective_batch, broken, expired |

---

## Data Changes

### BEFORE Marking as Damaged:
```
current_stock = 116105
available_stock = 116105
damaged_stock = 0
```

### AFTER Marking 10 Units as Damaged:
```
current_stock = 116105 (UNCHANGED)
available_stock = 116095 (REDUCED by 10)
damaged_stock = 10 (INCREASED by 10)
```

---

## Database Fields Added

### ProductStock Table:
```go
DamagedStock   float64   // Units marked as damaged
DamageReason   string    // Reason
DamagedAt      *time.Time // Timestamp
DamagedBy      string    // User ID who marked
```

### VariantStock Table:
```go
DamagedStock   float64   // Same as above
DamageReason   string
DamagedAt      *time.Time
DamagedBy      string
```

---

## Code Examples

### Mark Product as Damaged (Service Layer):

```go
userID := "user123"
productID := "prod_8b35c2b9"
quantity := 10.0
reason := "defective_batch"

if err := stockService.MarkProductAsDamaged(productID, quantity, reason, userID); err != nil {
    log.Printf("Error: %v", err)
    return
}
log.Println("Product marked as damaged successfully")
```

### Mark Variant as Damaged:

```go
variantSKU := "SL-001-PLAS"
quantity := 5.0
reason := "broken"

if err := variantService.MarkVariantAsDamaged(variantSKU, quantity, reason, userID); err != nil {
    log.Printf("Error: %v", err)
    return
}
log.Println("Variant marked as damaged successfully")
```

### Get Damaged Products:

```go
offset := 0
limit := 50

damaged, total, err := stockService.GetDamagedProducts(offset, limit)
if err != nil {
    log.Printf("Error: %v", err)
    return
}

for _, product := range damaged {
    fmt.Printf("Product: %s, Damaged: %.2f units\n", 
        product.ProductName, product.DamagedStock)
}
```

---

## Valid Damage Reasons

- `defective_batch` - Batch has quality defects
- `broken` - Product physically damaged
- `expired` - Product past expiration date
- `contaminated` - Product contaminated
- `lost` - Product lost in warehouse
- `theft` - Product stolen
- `quality_issue` - General quality issue
- `other` - Other reason (provide details)

---

## Error Handling

```go
if err != nil {
    switch err.Error() {
    case "insufficient available stock":
        // Not enough available units to mark as damaged
        fmt.Println("Cannot mark more units than available")
        
    case "product stock not found":
        // Product doesn't exist
        fmt.Println("Product not found in inventory")
        
    case "damage reason is required":
        // No reason provided
        fmt.Println("Must provide reason for damage")
        
    default:
        // Other error
        fmt.Printf("Error: %v\n", err)
    }
}
```

---

## Stock Ledger Entry

When a product is marked as damaged, a ledger entry is created:

```json
{
  "movement_type": "DAMAGE",
  "quantity": -10,
  "amount": -51.70,
  "reference_type": "damage_record",
  "reference_number": "DMG-1714298245",
  "notes": "Damage reason: defective_batch",
  "created_by": "user123"
}
```

---

## Response Structure

### Success Response:
```json
{
  "success": true,
  "message": "Variant marked as damaged successfully",
  "type": "variant",
  "damaged_stock": 10,
  "available_stock": 116095,
  "damage_reason": "defective_batch",
  "damaged_at": "2026-04-28T10:30:45.123Z"
}
```

### Error Response:
```json
{
  "success": false,
  "error": "insufficient available stock: have 5, trying to mark 10 as damaged"
}
```

---

## Routes

| Method | Endpoint | Purpose |
|--------|----------|---------|
| PATCH | `/api/stock/mark-damaged` | Mark product/variant as damaged |
| GET | `/api/stock/damaged` | Get all damaged items |
| GET | `/api/stock/summary` | Stock summary (includes damaged_products) |

---

## Implementation Files

| File | Purpose |
|------|---------|
| `models/product_stock.go` | ProductStock model with damaged fields |
| `models/variant_stock.go` | VariantStock model with damaged fields |
| `services/stock_management.service.go` | Mark product as damaged logic |
| `services/variant_stock_management.service.go` | Mark variant as damaged logic |
| `handlers/stock_management.handler.go` | API endpoints |
| `routes/routes.go` | Route registration |
| `repo/product_stock.repository.go` | Database queries for damaged products |
| `repo/variant_stock.repository.go` | Database queries for damaged variants |
| `migrations/006_add_damaged_stock_tracking.sql` | Database changes |

---

## Middleware Required

Both endpoints require:
- `middleware.AuthMiddleware()` - User must be authenticated
- `middleware.AdminMiddleware()` - User must have admin role

```go
stockRoutes.Patch("/mark-damaged", middleware.AdminMiddleware(), handler)
stockRoutes.Get("/damaged", middleware.AdminMiddleware(), handler)
```

---

## Useful Queries

### Find Damaged Products by Reason:
```sql
SELECT product_id, product_name, damaged_stock, damage_reason, damaged_at
FROM product_stocks
WHERE damaged_stock > 0
ORDER BY damaged_at DESC;
```

### Calculate Total Damaged Value:
```sql
SELECT 
    SUM(damaged_stock * average_cost) as total_damaged_value
FROM product_stocks
WHERE damaged_stock > 0;
```

### Find Recently Damaged Items:
```sql
SELECT product_id, product_name, damaged_stock, damaged_at
FROM product_stocks
WHERE damaged_stock > 0
AND damaged_at >= DATE_SUB(NOW(), INTERVAL 7 DAY)
ORDER BY damaged_at DESC;
```

---

## Common Tasks

### Task 1: Report All Damaged Items This Week

```bash
curl -X GET "http://127.0.0.1:8088/api/stock/damaged?limit=100" \
  -H "Authorization: Bearer {token}" | \
  jq '.damaged_items[] | select(.damaged_at >= now - 604800)'
```

### Task 2: Calculate Loss from Damaged Sleev Variant

```bash
curl -X GET "http://127.0.0.1:8088/api/stock/damaged" \
  -H "Authorization: Bearer {token}" | \
  jq '.damaged_items[] | select(.variant_sku == "SL-001-PLAS") | .damaged_value'
```

### Task 3: Find Highest Value Damaged Items

```bash
curl -X GET "http://127.0.0.1:8088/api/stock/damaged?limit=100" \
  -H "Authorization: Bearer {token}" | \
  jq '.damaged_items | sort_by(-.damaged_value) | .[0:10]'
```

---

## Performance Considerations

### Indices Created:
- `idx_damaged_stock` - Filter by damaged items
- `idx_product_damaged_status` - Joint index for product + damage status
- `idx_variant_damaged_stock` - Quick variant lookups
- `idx_damaged_at` - Time-based queries

### Query Performance:
- Getting damaged products: O(1) with index
- Filtering by reason: O(n) - Consider adding index if frequently queried
- Time-range queries: O(1) with `damaged_at` index

---

## Troubleshooting

### Issue: "Insufficient available stock" Error

**Cause**: Trying to mark more damaged units than are available

**Solution**:
```go
// Check available stock first
if stock.AvailableStock < damageQuantity {
    fmt.Println("Not enough available stock")
    // Reduce damage quantity
}
```

### Issue: Endpoint Returns 401 Unauthorized

**Cause**: User not authenticated or not admin

**Solution**:
```bash
# Verify token is valid
# Verify user has admin role
# Check middleware configuration
```

### Issue: Damaged Items Not Appearing in Summary

**Cause**: Response structure changed

**Solution**: Check `damaged_products` array in response (separate from `stocks` array)

---

## Logging

### Key Log Messages:

```
[DAMAGE_TRACKING] Marking product as damaged: Product=prod_xxx, Qty=10, Reason=defective_batch
[DAMAGE_TRACKING] Success: Damaged stock updated from 0 to 10 units
[DAMAGE_TRACKING] Marking variant as damaged: SKU=SL-001-PLAS, Qty=5, Reason=broken
```

---

## Security Notes

- ✅ Requires authentication (JWT token)
- ✅ Requires admin role (middleware)
- ✅ User ID tracked for audit trail
- ✅ All changes logged in stock ledger
- ✅ Input validation on quantity and reason

---

## Next Steps

1. Run migration: `006_add_damaged_stock_tracking.sql`
2. Restart service with updated code
3. Test endpoints with curl/Postman
4. Monitor logs for [DAMAGE_TRACKING] entries
5. Integrate into dashboard for visualization
6. Set up alerts for high damage values

---

## Links

- [Full API Documentation](./DAMAGED_PRODUCT_API.md)
- [Implementation Guide](./DAMAGED_PRODUCT_IMPLEMENTATION_GUIDE.md)
- [Flowchart](./DAMAGED_PRODUCT_IMPLEMENTATION_GUIDE.md#1-system-architecture)

---

## Questions?

Check the implementation files or API documentation for more details!
