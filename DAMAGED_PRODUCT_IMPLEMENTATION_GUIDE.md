# Damaged Product Management - Complete Implementation Guide

## Architecture Overview

This document outlines the complete implementation of the damaged product tracking system, including data flow, API endpoints, and integration points.

---

## 1. System Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                      API Request Layer                          │
│  PATCH /api/stock/mark-damaged  |  GET /api/stock/damaged      │
└──────────────┬────────────────────────────────────────────────┬─┘
               │                                                  │
               ▼                                                  ▼
    ┌──────────────────────────┐              ┌──────────────────────────┐
    │  StockManagementHandler  │              │  StockManagementHandler  │
    │  MarkProductAsDamaged()  │              │  GetDamagedProducts()    │
    └──────────────┬───────────┘              └──────────────┬───────────┘
                   │                                         │
                   ▼                                         ▼
    ┌──────────────────────────────────┐    ┌──────────────────────────────┐
    │    StockManagementService        │    │  StockManagementService      │
    │    MarkProductAsDamaged()        │    │  GetDamagedProducts()        │
    └──────────────┬────────────────────┘    └──────────────┬──────────────┘
                   │                                         │
                   ▼                                         ▼
    ┌──────────────────────────────────────────────────────────────┐
    │            ProductStockRepository                            │
    │    GetDamagedProducts()  |  GetDamagedProductsByUser()      │
    └──────────────┬───────────────────────────────┬──────────────┘
                   │                               │
                   ▼                               ▼
    ┌──────────────────────────────────────────────────────────────┐
    │                    Database                                  │
    │  product_stocks table  |  stock_ledgers table               │
    │  damage_logs table     |  damage_recovery table             │
    └──────────────────────────────────────────────────────────────┘
```

---

## 2. Data Flow for Marking Product as Damaged

### Step-by-Step Process:

```
REQUEST: PATCH /api/stock/mark-damaged
├─ Body: { product_id, variant_sku?, quantity, reason }
│
├─ 1. Handler validates request
│    └─ Check: product_id (required)
│    └─ Check: quantity > 0
│    └─ Check: reason not empty
│
├─ 2. Extract user context (user_id from JWT)
│
├─ 3. Route to appropriate service
│    ├─ IF variant_sku provided
│    │  └─ Call: variantStockMgmt.MarkVariantAsDamaged()
│    │     ├─ Get variant by SKU
│    │     ├─ Validate available_stock >= quantity
│    │     ├─ UPDATE variant_stocks:
│    │     │  ├─ damaged_stock += quantity
│    │     │  ├─ available_stock -= quantity
│    │     │  ├─ damage_reason = reason
│    │     │  ├─ damaged_at = NOW()
│    │     │  └─ damaged_by = user_id
│    │     ├─ CREATE VariantStockMovement:
│    │     │  ├─ movement_type = "DAMAGE"
│    │     │  ├─ quantity = -quantity
│    │     │  ├─ reference_type = "damage_record"
│    │     │  ├─ notes = damage reason
│    │     │  └─ created_by = user_id
│    │     └─ RETURN updated stock
│    │
│    └─ ELSE (product-level)
│       └─ Call: service.MarkProductAsDamaged()
│          ├─ Get product stock by product_id
│          ├─ Validate available_stock >= quantity
│          ├─ UPDATE product_stocks:
│          │  ├─ damaged_stock += quantity
│          │  ├─ available_stock -= quantity
│          │  ├─ damage_reason = reason
│          │  ├─ damaged_at = NOW()
│          │  └─ damaged_by = user_id
│          ├─ CREATE StockLedger:
│          │  ├─ movement_type = "DAMAGE"
│          │  ├─ quantity = -quantity
│          │  ├─ reference_type = "damage_record"
│          │  ├─ notes = damage reason
│          │  └─ created_by = user_id
│          └─ RETURN updated stock
│
└─ RESPONSE: 200 OK with updated stock values
```

---

## 3. Stock Calculation Formula

### Current Stock Balance Calculation:

```go
// BEFORE marking as damaged
available_stock = current_stock - reserved_stock - damaged_stock
                = current_stock - 0 - 0
                = current_stock

// Example: 116,105 units available

// AFTER marking 10 units as damaged
damaged_stock += 10              // 0 + 10 = 10
available_stock -= 10            // 116,105 - 10 = 116,095

// New calculation
available_stock = current_stock - reserved_stock - damaged_stock
                = 116,105 - 0 - 10
                = 116,095 ✓
```

### Key Formula:
```
available_stock = current_stock - reserved_stock - damaged_stock

Where:
- current_stock = purchased_stock - sold_stock (ALWAYS)
- reserved_stock = units reserved for sales orders
- damaged_stock = units marked as damaged/defective
```

---

## 4. Database Schema Updates

### ProductStock Table Changes:
```sql
-- Added columns
ALTER TABLE product_stocks
ADD COLUMN damaged_stock DOUBLE NOT NULL DEFAULT 0,
ADD COLUMN damage_reason TEXT,
ADD COLUMN damaged_at TIMESTAMP NULL,
ADD COLUMN damaged_by VARCHAR(255);

-- Added indices for performance
CREATE INDEX idx_damaged_stock ON product_stocks(damaged_stock);
CREATE INDEX idx_product_damaged_status ON product_stocks(product_id, damaged_stock);
```

### VariantStock Table Changes:
```sql
-- Added columns
ALTER TABLE variant_stocks
ADD COLUMN damaged_stock DOUBLE NOT NULL DEFAULT 0,
ADD COLUMN damage_reason TEXT,
ADD COLUMN damaged_at TIMESTAMP NULL,
ADD COLUMN damaged_by VARCHAR(255);

-- Added indices for performance
CREATE INDEX idx_variant_damaged_stock ON variant_stocks(damaged_stock);
CREATE INDEX idx_variant_damaged_status ON variant_stocks(variant_sku, damaged_stock);
```

### New Tables for Audit Trail:
```sql
-- Detailed damage logging
CREATE TABLE damage_logs (
    id INT PRIMARY KEY AUTO_INCREMENT,
    product_id VARCHAR(255),
    variant_sku VARCHAR(100),
    damaged_quantity DOUBLE,
    damage_reason TEXT,
    damage_category VARCHAR(50),  -- defective, broken, expired, etc.
    stock_ledger_id INT UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255),
    notes TEXT
);

-- Track damage recovery/returns
CREATE TABLE damage_recovery (
    id INT PRIMARY KEY AUTO_INCREMENT,
    damage_log_id INT UNIQUE,
    recovery_type VARCHAR(50),  -- returned, fixed, refund, written_off
    recovered_quantity DOUBLE,
    recovery_notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255)
);
```

---

## 5. API Request/Response Examples

### Example 1: Mark Variant as Damaged

**REQUEST:**
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

**RESPONSE (200):**
```json
{
  "success": true,
  "message": "Variant marked as damaged successfully",
  "type": "variant",
  "variant_sku": "SL-001-PLAS",
  "damaged_stock": 10,
  "available_stock": 116095,
  "damage_reason": "defective_batch",
  "damaged_at": "2026-04-28T10:30:45.123Z"
}
```

### Example 2: Get All Damaged Products

**REQUEST:**
```bash
curl -X GET "http://127.0.0.1:8088/api/stock/damaged?limit=10&offset=0" \
  -H "Authorization: Bearer {token}"
```

**RESPONSE (200):**
```json
{
  "success": true,
  "damaged_items": [
    {
      "type": "variant",
      "product_id": "prod_8b35c2b9",
      "product_name": "sleev",
      "variant_sku": "SL-001-PLAS",
      "variant_name": "SL-001-PLAS",
      "damaged_stock": 10,
      "damage_reason": "defective_batch",
      "damaged_at": "2026-04-28T10:30:45.123Z",
      "damaged_by": "user123",
      "average_cost": 5.17,
      "damaged_value": 51.70
    },
    {
      "type": "product",
      "product_id": "pg_test_15_units",
      "product_name": "test",
      "sku": "pg_test_15_units",
      "damaged_stock": 5,
      "damage_reason": "broken",
      "damaged_at": "2026-04-27T09:15:00.000Z",
      "damaged_by": "admin",
      "average_cost": 82.50,
      "damaged_value": 412.50
    }
  ],
  "damaged_products": 1,
  "damaged_variants": 1,
  "total_damaged": 2,
  "total_damaged_value": 464.20,
  "limit": 10,
  "offset": 0
}
```

### Example 3: Stock Summary with Damaged Items

**REQUEST:**
```bash
curl -X GET "http://127.0.0.1:8088/api/stock/summary" \
  -H "Authorization: Bearer {token}"
```

**RESPONSE (200):**
```json
{
  "stocks": [
    {
      "product_id": "prod_8b35c2b9",
      "product_name": "sleev",
      "sku": "SL-001-PLAS",
      "current_stock": 116105,
      "purchased_total": 116500,
      "sold_total": 395,
      "reserved_stock": 0,
      "available_stock": 116095,
      "damaged_stock": 10,
      "average_cost": 5.17,
      "stock_value": 600451.85,
      "type": "variant"
    }
  ],
  "damaged_products": [
    {
      "type": "variant",
      "product_id": "prod_8b35c2b9",
      "product_name": "sleev",
      "variant_sku": "SL-001-PLAS",
      "damaged_stock": 10,
      "average_cost": 5.17,
      "damaged_value": 51.70,
      "damage_reason": "defective_batch",
      "damaged_at": "2026-04-28T10:30:45.123Z"
    }
  ],
  "damaged_count": 1,
  "total": 5,
  "total_stock_value": 619638.47,
  "total_sold_product_value": 3036.53,
  "total_damaged_value": 51.70
}
```

---

## 6. Service Layer Implementation Details

### MarkProductAsDamaged Flow:

```go
func (s *stockManagementService) MarkProductAsDamaged(
    productID string,
    quantity float64,
    reason string,
    userID string,
) error {
    // 1. Validate inputs
    if quantity <= 0 {
        return errors.New("damage quantity must be positive")
    }
    
    // 2. Fetch current stock
    stock, err := s.productStockRepo.GetByProductID(productID)
    if err != nil {
        return fmt.Errorf("product stock not found: %s", productID)
    }
    
    // 3. Validate sufficient available stock
    if stock.AvailableStock < quantity {
        return fmt.Errorf("insufficient available stock")
    }
    
    // 4. Update stock object
    stock.DamagedStock += quantity
    stock.AvailableStock -= quantity
    stock.DamageReason = reason
    stock.DamagedAt = time.Now()
    stock.DamagedBy = userID
    stock.UpdatedAt = time.Now()
    
    // 5. Save to database
    if err := s.productStockRepo.Update(stock); err != nil {
        return fmt.Errorf("failed to update stock: %w", err)
    }
    
    // 6. Create ledger entry for audit trail
    amount := quantity * stock.AverageCost
    ledger := &models.StockLedger{
        ProductID:       productID,
        MovementType:    "DAMAGE",
        Quantity:        -quantity,           // Negative
        Rate:            stock.AverageCost,
        Amount:          -amount,             // Negative
        ReferenceType:   "damage_record",
        ReferenceID:     uuid.New().String(),
        ReferenceNumber: fmt.Sprintf("DMG-%d", time.Now().Unix()),
        Notes:           fmt.Sprintf("Damage reason: %s", reason),
        CreatedAt:       time.Now(),
        CreatedBy:       userID,
    }
    
    if err := s.stockLedgerRepo.Create(ledger); err != nil {
        log.Printf("Warning: Failed to create ledger entry: %v", err)
        // Continue - ledger entry is not critical
    }
    
    return nil
}
```

---

## 7. Repository Query Examples

### Get All Damaged Products:

```go
func (r *productStockRepository) GetDamagedProducts(
    offset, limit int,
) ([]models.ProductStock, int64, error) {
    var stocks []models.ProductStock
    var total int64
    
    // Query all products where damaged_stock > 0
    query := r.db.Model(&models.ProductStock{}).
        Where("damaged_stock > ?", 0)
    
    err := query.Count(&total).
        Preload("Product").
        Offset(offset).
        Limit(limit).
        Order("damaged_at DESC").
        Find(&stocks).Error
    
    return stocks, total, err
}
```

### Get Damaged Variants by User:

```go
func (r *variantStockRepository) GetDamagedVariantsByUser(
    userID uint,
    offset, limit int,
) ([]models.VariantStock, int64, error) {
    var stocks []models.VariantStock
    var total int64
    
    // Query variants belonging to user with damaged_stock > 0
    err := r.db.Model(&models.VariantStock{}).
        Joins("LEFT JOIN products ON variant_stocks.product_id = products.id").
        Where("(products.created_by = ? OR variant_stocks.product_id LIKE 'pg_%') "+
              "AND variant_stocks.damaged_stock > ?", 
              userID, 0).
        Count(&total).
        Offset(offset).
        Limit(limit).
        Order("variant_stocks.damaged_at DESC").
        Find(&stocks).Error
    
    return stocks, total, err
}
```

---

## 8. Error Handling

### Common Error Scenarios:

```json
// Error 1: Insufficient stock
{
  "success": false,
  "error": "insufficient available stock: have 5, trying to mark 10 as damaged"
}

// Error 2: Invalid reason
{
  "success": false,
  "error": "damage reason is required"
}

// Error 3: Product not found
{
  "success": false,
  "error": "product stock not found: invalid_product_id"
}

// Error 4: Authentication failed
{
  "success": false,
  "error": "unauthorized"
}

// Error 5: Server error
{
  "success": false,
  "error": "failed to update product stock with damage"
}
```

---

## 9. Integration Checklist

- [x] Models updated (ProductStock, VariantStock)
- [x] Migration file created (006_add_damaged_stock_tracking.sql)
- [x] Service interfaces updated
- [x] Service implementations added
- [x] Repository interfaces updated
- [x] Repository implementations added
- [x] Handler methods implemented
- [x] Routes registered
- [x] Stock summary endpoint updated
- [x] API documentation created
- [ ] Unit tests (TODO)
- [ ] Integration tests (TODO)
- [ ] Load testing (TODO)
- [ ] Documentation review (TODO)

---

## 10. Testing Guide

### Manual Test: Mark Product as Damaged

```bash
# 1. Get current stock
curl -X GET http://127.0.0.1:8088/api/stock/summary \
  -H "Authorization: Bearer {token}"

# 2. Mark 10 sleev units as damaged
curl -X PATCH http://127.0.0.1:8088/api/stock/mark-damaged \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": "prod_8b35c2b9",
    "variant_sku": "SL-001-PLAS",
    "quantity": 10,
    "reason": "defective_batch"
  }'

# 3. Verify damaged stock appears in summary
curl -X GET http://127.0.0.1:8088/api/stock/summary \
  -H "Authorization: Bearer {token}"

# 4. Get damaged products report
curl -X GET http://127.0.0.1:8088/api/stock/damaged \
  -H "Authorization: Bearer {token}"

# Expected: damaged_products array should contain the damaged sleev variant
```

---

## 11. Deployment Steps

1. **Backup Database**
   ```bash
   mysqldump -u user -p database > backup.sql
   ```

2. **Run Migration**
   ```bash
   mysql -u user -p database < migrations/006_add_damaged_stock_tracking.sql
   ```

3. **Deploy Code**
   ```bash
   git pull origin main
   go build
   ./auth-service
   ```

4. **Verify**
   ```bash
   # Test PATCH /api/stock/mark-damaged endpoint
   # Test GET /api/stock/damaged endpoint
   # Test GET /api/stock/summary endpoint
   ```

---

## 12. Monitoring & Logging

### Key Metrics to Monitor:

```
- Total damaged stock value
- Damaged items count
- Damage reasons distribution
- Damaged stock trend over time
- User who marked items as damaged
```

### Log Patterns:

```
[DAMAGE_TRACKING] Marking product as damaged: Product=prod_xxx, Qty=10, Reason=defective_batch
[DAMAGE_TRACKING] Success: Damaged stock updated from 0 to 10 units
[DAMAGE_TRACKING] Marking variant as damaged: SKU=SL-001-PLAS, Qty=5, Reason=broken
```

---

## Support & Troubleshooting

For issues with the damaged product feature, check:

1. **Database migration applied**: `SELECT * FROM product_stocks LIMIT 1;` should show `damaged_stock` column
2. **Routes registered**: Check `/api/stock/mark-damaged` is accessible
3. **Permissions**: User must have admin role for these endpoints
4. **Stock availability**: Ensure available_stock >= damage quantity

---

## References

- [API Documentation](./DAMAGED_PRODUCT_API.md)
- [Flowchart Diagram](#architecture-overview)
- Models: [ProductStock](./models/product_stock.go), [VariantStock](./models/variant_stock.go)
- Services: [StockManagementService](./services/stock_management.service.go)
- Handlers: [StockManagementHandler](./handlers/stock_management.handler.go)
- Migration: [006_add_damaged_stock_tracking.sql](./migrations/006_add_damaged_stock_tracking.sql)
