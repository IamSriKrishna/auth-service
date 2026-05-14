# Stock Management API - DTOs & Request/Response Examples

## Table of Contents
1. [Mark Damaged (PATCH)](#1-mark-damaged-patch)
2. [Get Damaged Products (GET)](#2-get-damaged-products-get)
3. [Get Stock Summary (GET)](#3-get-stock-summary-get)

---

## 1. Mark Damaged (PATCH)

### Endpoint
```
PATCH /api/stock/mark-damaged
Authorization: Bearer {admin_token}
Content-Type: application/json
```

### Input DTO

```go
type MarkDamagedInput struct {
	ProductID  string  `json:"product_id" validate:"required"`
	VariantSKU *string `json:"variant_sku"`
	Quantity   float64 `json:"quantity" validate:"required,gt=0"`
	Reason     string  `json:"reason" validate:"required"`
}
```

### Request Examples

#### Example 1: Mark Variant as Damaged
```json
{
  "product_id": "prod_8b35c2b9",
  "variant_sku": "SL-001-PLAS",
  "quantity": 10,
  "reason": "defective_batch"
}
```

#### Example 2: Mark Product as Damaged (Without Variant)
```json
{
  "product_id": "prod_12345",
  "quantity": 50,
  "reason": "broken"
}
```

#### Example 3: Mark as Expired
```json
{
  "product_id": "prod_test_15_units",
  "quantity": 5,
  "reason": "expired"
}
```

### Response DTO (Fiber Map)

```go
{
  "success": bool,
  "message": string,
  "type": "variant|product",
  "variant_sku": string,          // if variant
  "product_id": string,            // if product
  "damaged_stock": float64,
  "available_stock": float64,
  "damage_reason": string,
  "damaged_at": time.Time
}
```

### Response Examples

#### Success Response - Variant Damaged
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

#### Success Response - Product Damaged
```json
{
  "success": true,
  "message": "Product marked as damaged successfully",
  "type": "product",
  "product_id": "prod_12345",
  "damaged_stock": 50,
  "available_stock": 5000,
  "damage_reason": "broken",
  "damaged_at": "2026-04-28T11:15:30.456Z"
}
```

### Error Responses

#### Error: Insufficient Stock
```json
{
  "success": false,
  "error": "insufficient available stock: have 5, trying to mark 10 as damaged"
}
```

#### Error: Invalid Quantity
```json
{
  "success": false,
  "error": "damage quantity must be positive"
}
```

#### Error: Product Not Found
```json
{
  "success": false,
  "error": "product stock not found: invalid_product_id"
}
```

#### Error: Missing Required Field
```json
{
  "success": false,
  "error": "Invalid request body"
}
```

### Validation Rules
- **product_id**: Required, must be valid UUID or product ID
- **variant_sku**: Optional, if provided routes to variant service
- **quantity**: Required, must be > 0
- **reason**: Required, non-empty string

### Valid Damage Reasons
```
defective_batch    - Batch has defects
broken             - Item physically broken
expired            - Item past expiration date
contaminated       - Item contaminated or unsafe
lost               - Item lost in warehouse
theft              - Item stolen/missing
quality_issue      - Quality doesn't meet standards
other              - Other reasons
```

---

## 2. Get Damaged Products (GET)

### Endpoint
```
GET /api/stock/damaged?limit=50&offset=0
Authorization: Bearer {admin_token}
Content-Type: application/json
```

### Query Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| limit | integer | 50 | Number of records to fetch |
| offset | integer | 0 | Number of records to skip |

### Output DTO

```go
type DamagedItemOutput struct {
	Type             string    `json:"type"`             // "product" or "variant"
	ProductID        string    `json:"product_id"`
	ProductName      string    `json:"product_name"`
	SKU              string    `json:"sku,omitempty"`             // for products
	VariantSKU       string    `json:"variant_sku,omitempty"`     // for variants
	VariantName      string    `json:"variant_name,omitempty"`
	DamagedStock     float64   `json:"damaged_stock"`
	DamageReason     string    `json:"damage_reason"`
	DamagedAt        *time.Time `json:"damaged_at"`
	DamagedBy        string    `json:"damaged_by"`
	AverageCost      float64   `json:"average_cost"`
	DamagedValue     float64   `json:"damaged_value"`
}

type GetDamagedProductsResponse struct {
	Success           bool                   `json:"success"`
	DamagedItems      []DamagedItemOutput    `json:"damaged_items"`
	DamagedProducts   int64                  `json:"damaged_products"`
	DamagedVariants   int64                  `json:"damaged_variants"`
	TotalDamaged      int64                  `json:"total_damaged"`
	TotalDamagedValue float64                `json:"total_damaged_value"`
	Limit             int                    `json:"limit"`
	Offset            int                    `json:"offset"`
}
```

### Response Example - Multiple Damaged Items

```json
{
  "success": true,
  "damaged_items": [
    {
      "type": "variant",
      "product_id": "prod_8b35c2b9",
      "product_name": "sleev",
      "variant_sku": "SL-001-PLAS",
      "variant_name": "Plastic Sleeve",
      "damaged_stock": 10,
      "damage_reason": "defective_batch",
      "damaged_at": "2026-04-28T10:30:45.123Z",
      "damaged_by": "user123",
      "average_cost": 5.17,
      "damaged_value": 51.70
    },
    {
      "type": "variant",
      "product_id": "prod_8b35c2b9",
      "product_name": "sleev",
      "variant_sku": "SL-002-KRAFT",
      "variant_name": "Kraft Sleeve",
      "damaged_stock": 5,
      "damage_reason": "broken",
      "damaged_at": "2026-04-27T14:22:10.000Z",
      "damaged_by": "admin",
      "average_cost": 3.50,
      "damaged_value": 17.50
    },
    {
      "type": "product",
      "product_id": "pg_test_15_units",
      "product_name": "test",
      "sku": "pg_test_15_units",
      "damaged_stock": 3,
      "damage_reason": "expired",
      "damaged_at": "2026-04-26T09:00:00.000Z",
      "damaged_by": "warehouse_user",
      "average_cost": 82.50,
      "damaged_value": 247.50
    }
  ],
  "damaged_products": 1,
  "damaged_variants": 2,
  "total_damaged": 3,
  "total_damaged_value": 316.70,
  "limit": 50,
  "offset": 0
}
```

### Response Example - No Damaged Items

```json
{
  "success": true,
  "damaged_items": [],
  "damaged_products": 0,
  "damaged_variants": 0,
  "total_damaged": 0,
  "total_damaged_value": 0,
  "limit": 50,
  "offset": 0
}
```

### Response Example - Paginated Results

```json
{
  "success": true,
  "damaged_items": [
    {
      "type": "variant",
      "product_id": "prod_item_1",
      "product_name": "Product 1",
      "variant_sku": "VAR-001",
      "damaged_stock": 100,
      "damage_reason": "defective_batch",
      "damaged_at": "2026-04-28T10:30:45.123Z",
      "damaged_by": "user1",
      "average_cost": 10.00,
      "damaged_value": 1000.00
    }
  ],
  "damaged_products": 0,
  "damaged_variants": 50,
  "total_damaged": 50,
  "total_damaged_value": 15000.00,
  "limit": 1,
  "offset": 0
}
```

---

## 3. Get Stock Summary (GET)

### Endpoint
```
GET /api/stock/summary?limit=100&offset=0
Authorization: Bearer {token}
Content-Type: application/json
```

### Query Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| limit | integer | 100 | Number of records to fetch |
| offset | integer | 0 | Number of records to skip |
| view_user_id | string | current_user | User ID to view (admin only) |

### Output DTOs

```go
type StockSummaryItemOutput struct {
	ProductID      string    `json:"product_id"`
	ProductName    string    `json:"product_name"`
	SKU            string    `json:"sku"`
	VariantName    string    `json:"variant_name,omitempty"`
	Type           string    `json:"type"`              // "product" or "variant"
	CurrentStock   float64   `json:"current_stock"`
	PurchasedTotal float64   `json:"purchased_total"`
	SoldTotal      float64   `json:"sold_total"`
	ReservedStock  float64   `json:"reserved_stock"`
	AvailableStock float64   `json:"available_stock"`
	DamagedStock   float64   `json:"damaged_stock"`
	AverageCost    float64   `json:"average_cost"`
	StockValue     float64   `json:"stock_value"`
	LastPurchased  *time.Time `json:"last_purchased,omitempty"`
	LastSold       *time.Time `json:"last_sold,omitempty"`
}

type DamagedProductSummaryOutput struct {
	ProductID    string    `json:"product_id"`
	ProductName  string    `json:"product_name"`
	SKU          string    `json:"sku,omitempty"`
	VariantSKU   string    `json:"variant_sku,omitempty"`
	VariantName  string    `json:"variant_name,omitempty"`
	Type         string    `json:"type"`
	DamagedStock float64   `json:"damaged_stock"`
	AverageCost  float64   `json:"average_cost"`
	DamagedValue float64   `json:"damaged_value"`
	DamageReason string    `json:"damage_reason"`
	DamagedAt    *time.Time `json:"damaged_at"`
	DamagedBy    string    `json:"damaged_by"`
}

type GetStockSummaryResponse struct {
	Stocks                 []StockSummaryItemOutput        `json:"stocks"`
	DamagedProducts        []DamagedProductSummaryOutput   `json:"damaged_products"`
	DamagedCount           int                             `json:"damaged_count"`
	Total                  int64                           `json:"total"`
	TotalStockValue        float64                         `json:"total_stock_value"`
	TotalSoldProductValue  float64                         `json:"total_sold_product_value"`
	TotalDamagedValue      float64                         `json:"total_damaged_value"`
}
```

### Response Example - Full Summary with Damaged Items

```json
{
  "stocks": [
    {
      "product_id": "prod_8b35c2b9",
      "product_name": "sleev",
      "sku": "SL-001-PLAS",
      "variant_name": "Plastic Sleeve",
      "type": "variant",
      "current_stock": 116105,
      "purchased_total": 116500,
      "sold_total": 395,
      "reserved_stock": 0,
      "available_stock": 116095,
      "damaged_stock": 10,
      "average_cost": 5.17,
      "stock_value": 600451.85,
      "last_purchased": "2026-04-20T10:00:00Z",
      "last_sold": "2026-04-28T15:30:00Z"
    },
    {
      "product_id": "prod_8b35c2b9",
      "product_name": "sleev",
      "sku": "SL-002-KRAFT",
      "variant_name": "Kraft Sleeve",
      "type": "variant",
      "current_stock": 85000,
      "purchased_total": 90000,
      "sold_total": 5000,
      "reserved_stock": 0,
      "available_stock": 85000,
      "damaged_stock": 0,
      "average_cost": 3.50,
      "stock_value": 297500.00,
      "last_purchased": "2026-04-15T08:00:00Z",
      "last_sold": "2026-04-28T14:00:00Z"
    },
    {
      "product_id": "pg_test_15_units",
      "product_name": "test",
      "sku": "pg_test_15_units",
      "type": "product",
      "current_stock": 15,
      "purchased_total": 20,
      "sold_total": 5,
      "reserved_stock": 0,
      "available_stock": 12,
      "damaged_stock": 3,
      "average_cost": 82.50,
      "stock_value": 1237.50,
      "last_purchased": "2026-04-10T12:00:00Z",
      "last_sold": "2026-04-26T09:30:00Z"
    }
  ],
  "damaged_products": [
    {
      "product_id": "prod_8b35c2b9",
      "product_name": "sleev",
      "variant_sku": "SL-001-PLAS",
      "variant_name": "Plastic Sleeve",
      "type": "variant",
      "damaged_stock": 10,
      "average_cost": 5.17,
      "damaged_value": 51.70,
      "damage_reason": "defective_batch",
      "damaged_at": "2026-04-28T10:30:45.123Z",
      "damaged_by": "user123"
    },
    {
      "product_id": "pg_test_15_units",
      "product_name": "test",
      "sku": "pg_test_15_units",
      "type": "product",
      "damaged_stock": 3,
      "average_cost": 82.50,
      "damaged_value": 247.50,
      "damage_reason": "expired",
      "damaged_at": "2026-04-26T09:00:00Z",
      "damaged_by": "warehouse_user"
    }
  ],
  "damaged_count": 2,
  "total": 3,
  "total_stock_value": 899189.35,
  "total_sold_product_value": 18925.00,
  "total_damaged_value": 299.20
}
```

### Response Example - No Damaged Items

```json
{
  "stocks": [
    {
      "product_id": "prod_8b35c2b9",
      "product_name": "sleev",
      "sku": "SL-001-PLAS",
      "type": "variant",
      "current_stock": 116105,
      "purchased_total": 116500,
      "sold_total": 395,
      "reserved_stock": 0,
      "available_stock": 116105,
      "damaged_stock": 0,
      "average_cost": 5.17,
      "stock_value": 600451.85,
      "last_purchased": "2026-04-20T10:00:00Z",
      "last_sold": "2026-04-28T15:30:00Z"
    }
  ],
  "damaged_products": [],
  "damaged_count": 0,
  "total": 1,
  "total_stock_value": 600451.85,
  "total_sold_product_value": 2042.15,
  "total_damaged_value": 0
}
```

### Response Example - Only Damaged Products

```json
{
  "stocks": [],
  "damaged_products": [
    {
      "product_id": "prod_test_001",
      "product_name": "Defective Items",
      "sku": "DEF-001",
      "type": "product",
      "damaged_stock": 50,
      "average_cost": 10.00,
      "damaged_value": 500.00,
      "damage_reason": "quality_issue",
      "damaged_at": "2026-04-28T08:00:00Z",
      "damaged_by": "qa_user"
    }
  ],
  "damaged_count": 1,
  "total": 0,
  "total_stock_value": 0,
  "total_sold_product_value": 0,
  "total_damaged_value": 500.00
}
```

---

## Stock Calculation Formulas

### Available Stock Formula
```
available_stock = current_stock - reserved_stock - damaged_stock
```

### Stock Value Formula
```
stock_value = current_stock × average_cost
```

### Damaged Value Formula
```
damaged_value = damaged_stock × average_cost
```

### Example Calculation
```
Product: sleev (SL-001-PLAS)
├─ current_stock: 116,105 units
├─ reserved_stock: 0 units
├─ damaged_stock: 10 units
├─ average_cost: $5.17 per unit
│
├─ Formulas:
│  ├─ available_stock = 116,105 - 0 - 10 = 116,095 ✓
│  ├─ stock_value = 116,105 × $5.17 = $600,451.85 ✓
│  └─ damaged_value = 10 × $5.17 = $51.70 ✓
│
└─ Result:
   ├─ Available for sale: 116,095 units
   ├─ Total value: $600,451.85
   └─ Damaged value: $51.70
```

---

## Error Response Examples

### 401 Unauthorized
```json
{
  "success": false,
  "error": "unauthorized"
}
```

### 400 Bad Request
```json
{
  "success": false,
  "error": "Invalid request body"
}
```

### 403 Forbidden (Not Admin)
```json
{
  "success": false,
  "error": "forbidden: admin access required"
}
```

### 500 Internal Server Error
```json
{
  "success": false,
  "error": "failed to update product stock with damage"
}
```

---

## Integration Flow Example

### Step 1: Mark Items as Damaged
```bash
curl -X PATCH http://127.0.0.1:8088/api/stock/mark-damaged \
  -H "Authorization: Bearer {token}" \
  -d '{
    "product_id": "prod_8b35c2b9",
    "variant_sku": "SL-001-PLAS",
    "quantity": 10,
    "reason": "defective_batch"
  }'
```
Response: Damaged count increases to 10

### Step 2: Get Damaged Products
```bash
curl -X GET http://127.0.0.1:8088/api/stock/damaged \
  -H "Authorization: Bearer {token}"
```
Response: Shows sleev with 10 damaged units

### Step 3: Check Stock Summary
```bash
curl -X GET http://127.0.0.1:8088/api/stock/summary \
  -H "Authorization: Bearer {token}"
```
Response: Shows sleev in both `stocks` (with damaged_stock: 10) and `damaged_products` arrays

---

## Data Types Reference

| Field | Type | Format | Example |
|-------|------|--------|---------|
| product_id | string | UUID or custom ID | "prod_8b35c2b9" |
| variant_sku | string | SKU format | "SL-001-PLAS" |
| quantity | float64 | Positive number | 10, 10.5 |
| reason | string | Predefined values | "defective_batch" |
| damaged_stock | float64 | Positive number | 10.0 |
| damaged_value | float64 | Currency | 51.70 |
| damaged_at | time.Time | ISO 8601 | "2026-04-28T10:30:45.123Z" |
| damaged_by | string | User ID | "user123" |
| type | string | "product" or "variant" | "variant" |

---

**API Status**: Production Ready
**Last Updated**: May 1, 2026
**Version**: 1.0
