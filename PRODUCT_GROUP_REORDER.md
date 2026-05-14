# Product Group Reorder Feature

## Overview

The **Product Group Reorder** feature allows you to update quantities of existing products within a product group without creating new products or modifying the group structure. This is useful when you need to adjust order quantities while keeping the same products.

---

## Key Characteristics

✅ **Update quantities only** - Change quantity of existing products  
✅ **No product addition** - Cannot add new products to the group  
✅ **No product removal** - Cannot remove existing products from the group  
✅ **All products required** - Must include all existing products in reorder request  
✅ **Stock validation** - Checks if sufficient stock available for quantity increases  
✅ **Auto stock adjustment** - Deducts stock when increasing, releases stock when decreasing  

---

## Endpoint

### POST `/product-groups/:id/reorder`

**Authentication:** Required (Admin role)

**Authorization:** Admin middleware required

---

## Request Body

```json
{
  "products": [
    {
      "product_id": "prod_e65ac7bd",
      "variant_sku": "CP-001-RED",
      "quantity": 100
    },
    {
      "product_id": "prod_319bf367",
      "variant_sku": "WB-001-RED",
      "quantity": 100
    }
  ]
}
```

### Field Details

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `products` | Array | Yes | List of products with quantities |
| `product_id` | String | Yes | ID of existing product (must be in current group) |
| `variant_sku` | String | Yes | Variant SKU (must match existing) |
| `quantity` | Integer | Yes | New quantity (must be positive) |

---

## Response

### Success Response (200 OK)

```json
{
  "status": "success",
  "code": 200,
  "message": "Product group reordered successfully",
  "data": {
    "id": "pg_abc123",
    "reorder_summary": {
      "total_products": 2,
      "updates": [
        {
          "variant_sku": "CP-001-RED",
          "old_quantity": 10,
          "new_quantity": 100,
          "stock_adjusted": -90
        },
        {
          "variant_sku": "WB-001-RED",
          "old_quantity": 10,
          "new_quantity": 100,
          "stock_adjusted": -90
        }
      ]
    },
    "updated_at": "2026-05-03T14:22:00Z"
  }
}
```

---

## Validation Rules

### ❌ Product Not Found
Returns **404** if product group doesn't exist

### ❌ Invalid Product in Request
Returns **400** if any product in request:
- Doesn't exist in the current product group
- Has different variant SKU than stored
- Has mismatched product_id

### ❌ Missing Existing Products
Returns **400** if any existing product is not included in reorder request

**Example Error:**
```json
{
  "status": "error",
  "code": 400,
  "message": "Product with variant_sku 'WB-001-RED' not found in existing products",
  "error": "PRODUCT_NOT_FOUND_IN_GROUP"
}
```

### ❌ Insufficient Stock
Returns **400** if available stock is less than quantity increase

**Example Error:**
```json
{
  "status": "error",
  "code": 400,
  "message": "Insufficient stock for variant_sku 'CP-001-RED'. Available: 50, Required: 90",
  "error": "INSUFFICIENT_STOCK"
}
```

---

## Stock Management

### Quantity Increase

**Original quantity:** 10  
**New quantity:** 100  
**Difference:** 90 additional units needed

**Action:**
1. Check if 90 units available in component stock
2. If available, deduct 90 from component stock
3. Update product group quantity to 100

**Stock after reorder:**
- Component stock: reduced by 90
- Product group inventory: 100

### Quantity Decrease

**Original quantity:** 100  
**New quantity:** 50  
**Difference:** 50 units being released

**Action:**
1. Release 50 units back to component stock
2. Update product group quantity to 50

**Stock after reorder:**
- Component stock: increased by 50
- Product group inventory: 50

### No Change

**Original quantity:** 50  
**New quantity:** 50  
**Difference:** 0

**Action:** No stock adjustment, quantities remain same

---

## Example Workflow

### Initial State

**Product Group: "Cleaning Kit"**
```
Cap (CP-001-RED):           10 units
Water Bottle (WB-001-RED): 10 units
```

**Component Stock:**
```
Cap (CP-001-RED):          200 units available
Water Bottle (WB-001-RED): 300 units available
```

### Reorder Request

```json
POST /product-groups/pg_abc123/reorder

{
  "products": [
    {
      "product_id": "prod_e65ac7bd",
      "variant_sku": "CP-001-RED",
      "quantity": 100
    },
    {
      "product_id": "prod_319bf367",
      "variant_sku": "WB-001-RED",
      "quantity": 50
    }
  ]
}
```

### Validation Process

1. ✅ Product group found
2. ✅ Both products exist in group
3. ✅ No new products added
4. ✅ All existing products included
5. ✅ Check stock:
   - Cap: need 90 more (200 available) ✅
   - Water Bottle: need 40 less (no check needed) ✅
6. ✅ All validations pass

### Stock Adjustments

```
Cap stock:
  Before: 200 units
  Deduct: -90 units (increase from 10 to 100)
  After:  110 units

Water Bottle stock:
  Before: 300 units
  Add:    +10 units (decrease from 10 to 50)
  After:  310 units
```

### Final State

**Product Group: "Cleaning Kit"**
```
Cap (CP-001-RED):          100 units (updated)
Water Bottle (WB-001-RED): 50 units (updated)
```

**Component Stock:**
```
Cap (CP-001-RED):          110 units (adjusted)
Water Bottle (WB-001-RED): 310 units (adjusted)
```

---

## Implementation Details

### Service Layer
**File:** `/app/services/product_group.service.go`  
**Method:** `Reorder(id string, input *input.UpdateProductGroupInput) (*output.ProductGroupOutput, error)`

**Process:**
1. Fetches existing product group
2. Builds map of existing products by `productID:variantSku`
3. Validates all request products exist in map
4. Validates all existing products included in request
5. Checks stock availability for increases
6. Records stock adjustments using VariantStockManagementService
7. Updates component quantities in-place
8. Reinitializes product group inventory
9. Returns updated product group

### Handler Layer
**File:** `/app/handlers/product_group.handler.go`  
**Method:** `ReorderProductGroup(c *fiber.Ctx)`

**Process:**
1. Extracts product group ID from URL parameter
2. Parses request body into UpdateProductGroupInput
3. Calls ProductGroupService.Reorder()
4. Returns response with updated product group

### Route Registration
**File:** `/app/routes/routes.go`

```go
// POST /product-groups/:id/reorder
productGroupRoutes.Post("/:id/reorder", middleware.AuthMiddleware(), middleware.AdminMiddleware(), handlers.ReorderProductGroup)
```

**Route Order:** Specific routes (like `/:id/reorder`) are registered BEFORE generic routes (like `/:id`) to ensure proper matching.

---

## Error Handling

| Error Code | HTTP Status | Cause | Resolution |
|-----------|-------------|-------|-----------|
| `PRODUCT_GROUP_NOT_FOUND` | 404 | Product group doesn't exist | Verify product group ID |
| `PRODUCT_NOT_FOUND_IN_GROUP` | 400 | Product not in current group | Use existing product IDs |
| `VARIANT_SKU_MISMATCH` | 400 | Variant SKU doesn't match | Verify variant SKU matches stored |
| `MISSING_EXISTING_PRODUCTS` | 400 | Existing products not in request | Include all existing products |
| `INSUFFICIENT_STOCK` | 400 | Not enough stock to increase quantity | Reduce quantity or add stock |
| `INVALID_QUANTITY` | 400 | Quantity is negative or zero | Use positive quantity values |

---

## Differences from Create

| Aspect | Create | Reorder |
|--------|--------|---------|
| **Add new products** | ✅ Yes | ❌ No |
| **Update quantities** | N/A (new group) | ✅ Yes |
| **Remove products** | N/A (new group) | ❌ No |
| **Stock check** | Optional | ✅ Required for increases |
| **Use case** | Create new kit | Adjust existing kit quantities |

---

## Best Practices

1. **Always include all products** - Even if not changing quantity
2. **Check stock before reorder** - To avoid insufficient stock errors
3. **Use exact variant SKU** - Case-sensitive, must match exactly
4. **Verify product IDs** - Match the product_id in product group
5. **Monitor stock after reorder** - Adjustments are automatic but verify results

---

## Testing

### Test Case 1: Successful Reorder

**Setup:** Product group with 2 products (qty 10 each)  
**Action:** Increase both to 100  
**Expected:** Success, stock deducted, quantities updated  

### Test Case 2: Insufficient Stock

**Setup:** Product group with product qty 10, only 20 stock available  
**Action:** Try to increase to 100  
**Expected:** 400 error, insufficient stock

### Test Case 3: Missing Product

**Setup:** Product group with 2 products  
**Action:** Only include 1 product in request  
**Expected:** 400 error, missing product

### Test Case 4: New Product

**Setup:** Product group with product A  
**Action:** Try to add product B  
**Expected:** 400 error, product not in group

### Test Case 5: Decrease Quantity

**Setup:** Product group with qty 100  
**Action:** Decrease to 50  
**Expected:** Success, stock released, quantity updated to 50

---

## cURL Examples

### Increase Quantities

```bash
curl -X POST http://localhost:3000/product-groups/pg_abc123/reorder \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "products": [
      {
        "product_id": "prod_001",
        "variant_sku": "CP-001-RED",
        "quantity": 100
      },
      {
        "product_id": "prod_002",
        "variant_sku": "WB-001-RED",
        "quantity": 100
      }
    ]
  }'
```

### Decrease Quantities

```bash
curl -X POST http://localhost:3000/product-groups/pg_abc123/reorder \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "products": [
      {
        "product_id": "prod_001",
        "variant_sku": "CP-001-RED",
        "quantity": 50
      },
      {
        "product_id": "prod_002",
        "variant_sku": "WB-001-RED",
        "quantity": 50
      }
    ]
  }'
```

---

## Related Features

- **Create Product Group** - Create new product groups with initial products
- **Update Product Group** - Modify group metadata (name, description, status)
- **Delete Product Group** - Remove entire product group
- **Stock Management** - View and track stock levels
- **Inventory Tracking** - Monitor product group inventory

