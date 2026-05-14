# Damaged Product Management API Documentation

## Overview
The Damaged Product Management feature allows you to track defective, damaged, or unusable products and variants separately from regular inventory. When a product/variant is marked as damaged, it is:
- Deducted from available_stock
- Added to damaged_stock
- Recorded in stock ledger with movement_type: "DAMAGE"
- Displayed separately in stock reports

## Data Model Changes

### ProductStock Model
```go
// New fields added:
DamagedStock   float64   // Units marked as damaged
DamageReason   string    // Reason for damage
DamagedAt      *time.Time // When marked as damaged
DamagedBy      string    // User who marked as damaged
```

### VariantStock Model
```go
// New fields added (same as ProductStock):
DamagedStock   float64
DamageReason   string
DamagedAt      *time.Time
DamagedBy      string
```

## API Endpoints

### 1. Mark Product/Variant as Damaged
**Endpoint:** `PATCH /api/stock/mark-damaged`

**Request Body:**
```json
{
  "product_id": "prod_8b35c2b9",
  "variant_sku": "SL-001-PLAS",  // Optional - only for variants
  "quantity": 10,                 // Quantity to mark as damaged
  "reason": "defective_batch"     // Required reason
}
```

**Valid Reason Categories:**
- defective_batch
- broken
- expired
- contaminated
- lost
- theft
- quality_issue
- other

**Success Response (200):**
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

**Error Response (400/500):**
```json
{
  "success": false,
  "error": "insufficient available stock: have 5, trying to mark 10 as damaged"
}
```

---

### 2. Get All Damaged Products & Variants
**Endpoint:** `GET /api/stock/damaged?limit=50&offset=0`

**Query Parameters:**
- `limit` (optional, default: 50) - Number of records to return
- `offset` (optional, default: 0) - Pagination offset

**Response (200):**
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
      "product_id": "prod_test_123",
      "product_name": "Test Product",
      "sku": "TEST-SKU",
      "damaged_stock": 5,
      "damage_reason": "broken",
      "damaged_at": "2026-04-27T15:20:30.456Z",
      "damaged_by": "admin",
      "average_cost": 25.00,
      "damaged_value": 125.00
    }
  ],
  "damaged_products": 1,
  "damaged_variants": 1,
  "total_damaged": 2,
  "total_damaged_value": 176.70,
  "limit": 50,
  "offset": 0
}
```

---

### 3. Get Stock Summary (with Damage Tracking)
**Endpoint:** `GET /api/stock/summary`

**Response (200) - Updated to include damaged products:**
```json
{
  "stocks": [
    {
      "product_id": "prod_8b35c2b9",
      "product_name": "sleev",
      "sku": "SL-001-PLAS",
      "variant_name": "SL-001-PLAS",
      "current_stock": 116105,
      "purchased_total": 116500,
      "sold_total": 395,
      "reserved_stock": 0,
      "available_stock": 116095,
      "damaged_stock": 10,
      "average_cost": 5.17,
      "stock_value": 599893.85,
      "last_purchased": "2026-04-26T13:50:10.338+05:30",
      "last_sold": "2026-04-14T04:29:40.298+05:30",
      "type": "variant"
    }
  ],
  "damaged_products": [
    {
      "type": "variant",
      "product_id": "prod_8b35c2b9",
      "product_name": "sleev",
      "variant_sku": "SL-001-PLAS",
      "variant_name": "SL-001-PLAS",
      "damaged_stock": 10,
      "average_cost": 5.17,
      "damaged_value": 51.70,
      "damage_reason": "defective_batch",
      "damaged_at": "2026-04-28T10:30:45.123Z",
      "damaged_by": "user123"
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

## Stock Calculation Examples

### Before Marking as Damaged:
```
current_stock:    116105 units
available_stock:  116105 units
damaged_stock:    0 units
```

### After PATCH marking 10 units as damaged:
```
current_stock:    116105 units (UNCHANGED - current always = purchased - sold)
available_stock:  116095 units (REDUCED by 10)
damaged_stock:    10 units     (INCREASED by 10)
```

### Key Points:
- ✅ **current_stock** remains unchanged (it's always: purchased - sold)
- ✅ **available_stock** decreases (current - reserved - damaged)
- ✅ **damaged_stock** increases
- ✅ Stock ledger entry created with movement_type: "DAMAGE"
- ✅ Both damaged and available calculated from current, not duplicated

---

## Stock Ledger Tracking

When marking products as damaged, a stock ledger entry is created:

```go
StockLedger {
    ID:              auto-increment,
    ProductID:       "prod_8b35c2b9",
    MovementType:    "DAMAGE",
    Quantity:        -10,              // Negative to indicate reduction
    Rate:            5.17,             // Average cost
    Amount:          -51.70,           // Negative value reduction
    ReferenceType:   "damage_record",
    ReferenceID:     "uuid",
    ReferenceNumber: "DMG-1704009000", // DMG-{timestamp}
    BalanceBeforeQty: 116105,
    BalanceAfterQty: 116095,
    Notes:           "Damage reason: defective_batch",
    CreatedAt:       2026-04-28,
    CreatedBy:       "user123"
}
```

---

## Dashboard Integration

The damaged status is also reflected in the dashboard:

**Endpoint:** `GET /dashboard/stock`

Each product in the response includes:
```json
{
  "product_id": "prod_8b35c2b9",
  "product_name": "sleev",
  "current_stock": 116105,
  "available_stock": 116095,
  "reserved_stock": 0,
  "damaged_stock": 10,
  "status": "in_stock",
  "damage_info": {
    "damaged_quantity": 10,
    "damage_reason": "defective_batch",
    "damaged_at": "2026-04-28T10:30:45.123Z"
  }
}
```

---

## Error Handling

### Case 1: Insufficient Available Stock
```json
{
  "success": false,
  "error": "insufficient available stock: have 5, trying to mark 10 as damaged"
}
```

### Case 2: Missing Required Fields
```json
{
  "success": false,
  "error": "damage reason is required"
}
```

### Case 3: Invalid Product ID
```json
{
  "success": false,
  "error": "product stock not found: invalid_id"
}
```

---

## Use Cases

### 1. Warehouse Damage Report
```bash
GET /api/stock/damaged?limit=100&offset=0
```
Generate comprehensive damage report with all damaged items and their values.

### 2. Insurance Claim
Use the `damaged_value` field to calculate total loss:
```
total_damaged_value = SUM(damaged_quantity * average_cost)
```

### 3. Disposal Tracking
Filter by damage_reason to identify items for disposal:
```
- defective_batch: Return to supplier
- expired: Dispose as waste
- contaminated: Destroy immediately
- other: Manual review required
```

### 4. Quality Control
Track damaged items to identify quality issues:
```
- Monitor damage_reason trends
- Identify batch-level defects
- Trigger supplier investigations
```

---

## Database Migration

Run the migration to add damaged stock tracking:

```bash
# Migration file: 006_add_damaged_stock_tracking.sql
# Creates:
# - damaged_stock column in product_stocks
# - damaged_stock column in variant_stocks
# - DamageLog table for detailed tracking
# - DamageRecovery table for returns/fixes
# - Indices for performance optimization
```

---

## API Security

All endpoints require authentication and admin middleware:
- `middleware.AuthMiddleware()` - User must be authenticated
- `middleware.AdminMiddleware()` - User must have admin role

```go
stockRoutes.Patch("/mark-damaged", middleware.AdminMiddleware(), handler)
stockRoutes.Get("/damaged", middleware.AdminMiddleware(), handler)
```

---

## Response Headers

All responses include standard headers:
```
Content-Type: application/json
Authorization: Bearer {token}
```

---

## Audit Trail

Every damage marking is tracked:
- **createdAt**: Timestamp of damage record
- **createdBy**: User ID who marked as damaged
- **damageReason**: Detailed reason
- **damageCategory**: Categorized reason for analysis

Perfect for compliance and audit requirements.

---

## Performance Optimization

Indices are created for fast queries:
- `idx_damaged_stock` - Filter by damaged_stock > 0
- `idx_product_damaged_status` - Joint index for product + damage status
- `idx_variant_damaged_stock` - Quick variant damage lookups
- `idx_damaged_at` - Time-based queries for recent damages

---

## Future Enhancements

Potential features for future versions:
- [ ] Damage recovery/repair workflow
- [ ] Automatic damage alerts
- [ ] Damage root cause analysis
- [ ] Supplier accountability tracking
- [ ] Damage trend analytics
- [ ] RMA (Return Material Authorization) integration
- [ ] Warranty claim automation

---

## Questions?

For issues or questions about the Damaged Product API:
1. Check this documentation
2. Review the migration file: `006_add_damaged_stock_tracking.sql`
3. Check handler implementation: `stock_management.handler.go`
4. Review service implementation: `stock_management.service.go` and `variant_stock_management.service.go`
