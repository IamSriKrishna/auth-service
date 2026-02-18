# Complete Production to Invoice Workflow
## 500ml Water Bottle: From Items to Final Invoice

This guide demonstrates the complete workflow for manufacturing and selling 100 units of 500ml Complete Water Bottles.

---

## Troubleshooting Common Errors

### Error: "variant items must define attributes"

**Backend Validation:** The API now validates variant items automatically using the `ValidateVariantAttributes()` method.

**Requirements:**
1. `attribute_definitions` MUST be defined (array of objects with `key` and `options`)
2. Each variant MUST have an `attribute_map` that maps ALL defined attribute keys
3. Attribute values in `attribute_map` MUST match exactly with one of the defined options (case-sensitive)
4. Variants cannot have extra attributes not defined in `attribute_definitions`

**Validation Flow:**
```
Request arrives → ValidateVariantAttributes() checks:
  ✓ attribute_definitions exists
  ✓ variants exist
  ✓ Each variant has attribute_map
  ✓ All defined attributes present in each variant's attribute_map
  ✓ All attribute values match defined options (case-sensitive)
  ✓ No extra undefined attributes in variants
  → ✅ PASS (Item created) or ❌ FAIL (Detailed error message)
```

**Solution - Correct JSON Structure:**

```json
{
  "item_details": {
    "structure": "variants",
    "attribute_definitions": [
      {
        "key": "Cap Type",
        "options": ["Regular Cap", "Sports Cap"]
      }
    ],
    "variants": [
      {
        "sku": "WTR-BOT-500-REG",
        "attribute_map": {
          "Cap Type": "Regular Cap"
        },
        "selling_price": 15,
        "cost_price": 8,
        "stock_quantity": 5000
      },
      {
        "sku": "WTR-BOT-500-SPORT",
        "attribute_map": {
          "Cap Type": "Sports Cap"
        },
        "selling_price": 18,
        "cost_price": 9.5,
        "stock_quantity": 3000
      }
    ]
  }
}
```

**Common Mistakes to Avoid:**

1. ❌ Missing `attribute_definitions`
   ```json
   "variants": [
     { "sku": "ABC", "selling_price": 10, ... }
   ]
   ```
   **Error:** "variant items must define attributes"

2. ❌ Variant missing `attribute_map`
   ```json
   {
     "attribute_definitions": [{ "key": "Color", "options": ["Red"] }],
     "variants": [
       { "sku": "ABC", "selling_price": 10 }
     ]
   }
   ```
   **Error:** "variant 1 (ABC) must have attribute_map defined"

3. ❌ Missing required attribute in variant
   ```json
   {
     "attribute_definitions": [
       { "key": "Color", "options": ["Red", "Blue"] },
       { "key": "Size", "options": ["S", "M", "L"] }
     ],
     "variants": [
       { "sku": "ABC", "attribute_map": { "Color": "Red" } }
     ]
   }
   ```
   **Error:** "variant 1 (ABC) missing required attribute 'Size'"

4. ❌ Invalid attribute value (not in options)
   ```json
   {
     "attribute_definitions": [
       { "key": "Color", "options": ["Red", "Blue"] }
     ],
     "variants": [
       { "sku": "ABC", "attribute_map": { "Color": "Green" } }
     ]
   }
   ```
   **Error:** "variant 1 (ABC) has invalid value 'Green' for attribute 'Color'. Valid options: [Red Blue]"

5. ❌ Case sensitivity (Red ≠ red)
   ```json
   {
     "attribute_definitions": [
       { "key": "Color", "options": ["Red", "Blue"] }
     ],
     "variants": [
       { "sku": "ABC", "attribute_map": { "Color": "red" } }
     ]
   }
   ```
   **Error:** "variant 1 (ABC) has invalid value 'red' for attribute 'Color'. Valid options: [Red Blue]"

6. ❌ Extra undefined attributes
   ```json
   {
     "attribute_definitions": [
       { "key": "Color", "options": ["Red"] }
     ],
     "variants": [
       { "sku": "ABC", "attribute_map": { "Color": "Red", "Brand": "Nike" } }
     ]
   }
   ```
   **Error:** "variant 1 (ABC) has undefined attribute 'Brand'"

### Error: "Duplicate entry 'WTR-BOT-500-REG' for key 'variants.idx_variants_sku'"

**Cause:** The variant SKU already exists in the database. Variant SKUs must be globally unique across all items.

**Problem:** When you attempt to create the same item twice (or the same variant SKU twice), the database rejects it because of the unique constraint on the `variants.sku` field.

**Solutions:**

**Option 1: Use Different SKUs (Recommended)**
Use different variant SKUs in your request. Modify the SKUs to include a suffix or version:
```json
{
  "variants": [
    {
      "sku": "WTR-BOT-500-REG-V2",  // Changed from WTR-BOT-500-REG
      "attribute_map": { "Cap Type": "Regular Cap" },
      "selling_price": 15,
      "cost_price": 8,
      "stock_quantity": 5000
    }
  ]
}
```

**Option 2: Delete and Recreate (Development Only)**
If this is a test environment and you need to recreate the item:
1. Delete the existing item with that variant
2. Create a new item with the same or different SKU

**Option 3: Update Instead of Create**
If you want to update an existing item's variant (not recommended):
- Use the PUT `/items/{item_id}` endpoint instead of POST
- Update the existing variant's price/quantity rather than creating new ones

**Option 4: Query Existing Items**
Check if the item already exists:
```bash
GET /items
```
Look for items with SKU `WTR-BOT-500-BASE`. If it exists with the variant `WTR-BOT-500-REG`, skip Step 1.1 and use the existing item ID in Step 2 (Item Group creation).

**For Testing the Complete Workflow:**
If you've already created the items, skip Step 1 and proceed directly to Step 2 (Create Item Group) using the existing item IDs.

### Error: "variant not found" or "variant_sku not found"

**Cause:** The variant SKU being referenced doesn't exist in the item.

**Solution:** 
- Check that the item was created with variants showing the correct SKUs in the list
- Use the exact SKU as shown in the item details (case-sensitive)
- Confirm all components in the Item Group reference existing variant_skus
- Verify the SKU matches what was used during item creation
- If you get this error in Item Group creation, ensure the variant was created in Step 1

---

## Backend Implementation Details

### Variant Attribute Validation System

The backend now includes comprehensive validation for variant items through the `ValidateVariantAttributes()` method:

**What Changed:**
- JSON field renamed from `attributes` → `attribute_definitions` for clarity
- Validation moved to DTO layer (input validation) before service processing
- Validation occurs automatically on all item creation/update requests
- Detailed error messages guide users to the correct structure

**How It Works:**
1. When you POST an item with `structure: "variants"`, the validation runs automatically
2. It checks that attribute_definitions exist with proper keys and options
3. It verifies each variant has a matching attribute_map with all required attributes
4. It validates that attribute values exist in the defined options
5. It prevents extra undefined attributes in variants
6. Any validation failure returns a specific error message before the item is created

**Why This Matters:**
- **Prevents Data Corruption:** Ensures variants always have complete attribute mappings
- **Better Error Messages:** Users see exactly what's wrong and how to fix it
- **Single Source of Truth:** Attribute definitions defined once, used by all variants
- **Consistency Across Items:** All variant items follow the same validation rules

### JSON Field Name Change

**Old Field (No longer valid):**
```json
"attributes": [ { "key": "...", "options": [...] } ]
```

**New Field (Use this):**
```json
"attribute_definitions": [ { "key": "...", "options": [...] } ]
```

The field was renamed to clarify that these are the *definitions* of available attributes, not the *instances* of attributes on specific variants. Each variant's actual attribute assignments go in `attribute_map`.

---

## Table of Contents

0. [Backend Implementation Details](#backend-implementation-details)
1. [Troubleshooting Common Errors](#troubleshooting-common-errors)
2. [Step 1: Create Base Items](#step-1-create-base-items)
2. [Step 2: Create Item Group (BOM)](#step-2-create-item-group-bom)
3. [Step 3: Set Opening Stock](#step-3-set-opening-stock)
4. [Step 4: Create Production Order](#step-4-create-production-order)
5. [Step 5: Create Sales Order](#step-5-create-sales-order)
6. [Step 6: Create Invoice](#step-6-create-invoice)
7. [Step 7: Record Payment](#step-7-record-payment)

---

## Step 1: Create Base Items

### ⚠️ Important: Check Before Creating Items

**Before starting Step 1, verify if these items already exist in your system:**

```bash
# Query existing items
GET /items
```

Look for items with these base SKUs:
- `WTR-BOT-500-BASE` 
- `CAP-20MM-BASE`
- `LBL-1000-CUSTOM`

**If they already exist:** Skip to Step 2 and use the existing item IDs.

**If they don't exist:** Follow Step 1 to create them.

---

### 1.1 Create 500ml Water Bottle Item

**IMPORTANT:** For items with `structure: "variants"`, you MUST define attributes FIRST before adding variants.

**Endpoint:** `POST /items`

#### Step-by-Step (UI Form):
1. **Basic Information**
   - Item Name: `500ml Premium Drinking Water Bottle`
   - Item Type: Choose `Goods`
   - Structure: Choose `Variants`
   - Unit: `piece`
   - SKU (Base): `WTR-BOT-500-BASE`
   - UPC: `8904220500001`
   - EAN: `8904220500001`

2. **Define Variant Attributes FIRST**
   - Attribute Name: `Cap Type`
   - Options: `Regular Cap, Sports Cap` (comma-separated)
   - Click **"+ Add Attribute"** button
   - ✅ Wait for attribute to appear in the table before proceeding

3. **Then Add Variants**
   - Variant SKU: `WTR-BOT-500-REG`
   - Selling Price: `15`
   - Cost Price: `8`
   - Stock Quantity: `5000`
   - Cap Type: Select `Regular Cap` from dropdown
   - Click **"+ Add Variant"**
   
   - Variant SKU: `WTR-BOT-500-SPORT`
   - Selling Price: `18`
   - Cost Price: `9.5`
   - Stock Quantity: `3000`
   - Cap Type: Select `Sports Cap` from dropdown
   - Click **"+ Add Variant"**

4. **Sales & Purchase Information** - Fill in as shown in request below
5. **Inventory Management** - Enable "Track Inventory"
6. **Return Policy** - Enable "Item is Returnable"

**API Request:**
```json
{
  "name": "500ml Premium Drinking Water Bottle",
  "type": "goods",
  "item_details": {
    "structure": "variants",
    "unit": "piece",
    "sku": "WTR-BOT-500-BASE",
    "upc": "8904220500001",
    "ean": "8904220500001",
    "description": "Premium 500ml PET drinking water bottle with tamper-proof cap. BPA-free, food-grade plastic. Available with regular or sports cap options.",
    "attribute_definitions": [
      {
        "key": "Cap Type",
        "options": [
          "Regular Cap",
          "Sports Cap"
        ]
      }
    ],
    "variants": [
      {
        "sku": "WTR-BOT-500-REG",
        "attribute_map": {
          "Cap Type": "Regular Cap"
        },
        "selling_price": 15,
        "cost_price": 8,
        "stock_quantity": 5000
      },
      {
        "sku": "WTR-BOT-500-SPORT",
        "attribute_map": {
          "Cap Type": "Sports Cap"
        },
        "selling_price": 18,
        "cost_price": 9.5,
        "stock_quantity": 3000
      }
    ]
  },
  "sales_info": {
    "account": "Sales Revenue - Water Bottles",
    "selling_price": 16.5,
    "currency": "INR",
    "description": "500ml drinking water bottles retail sales"
  },
  "purchase_info": {
    "account": "Cost of Goods Purchased",
    "cost_price": 8.75,
    "currency": "INR",
    "preferred_vendor_id": 1,
    "description": "Purchase from bottle manufacturer"
  },
  "inventory": {
    "track_inventory": true,
    "inventory_account": "Inventory - Water Bottles",
    "inventory_valuation_method": "FIFO",
    "reorder_point": 1000
  },
  "return_policy": {
    "returnable": true
  }
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "item_26673b1a",
    "name": "500ml Premium Drinking Water Bottle",
    "type": "goods",
    "item_details": {
      "structure": "variants",
      "sku": "WTR-BOT-500-BASE",
      "variants": [
        {
          "sku": "WTR-BOT-500-REG",
          "attribute_map": { "Cap Type": "Regular Cap" }
        },
        {
          "sku": "WTR-BOT-500-SPORT",
          "attribute_map": { "Cap Type": "Sports Cap" }
        }
      ]
    },
    "created_at": "2026-02-17T15:36:43.898+05:30"
  },
  "message": "Item created successfully"
}
```

### 1.2 Create 20mm Tamper-proof Caps Item

**IMPORTANT:** Same rule applies - define attributes BEFORE variants.

**Endpoint:** `POST /items`

#### Step-by-Step (UI Form):
1. **Basic Information**
   - Item Name: `20mm Tamper-proof Bottle Cap - Pack of 100`
   - Item Type: `Goods`
   - Structure: `Variants`
   - Unit: `pack`
   - SKU (Base): `CAP-20MM-BASE`

2. **Define Variant Attributes FIRST**
   - Attribute Name: `Color`
   - Options: `White, Blue` (comma-separated)
   - Click **"+ Add Attribute"**
   - ✅ Wait for attribute to appear before adding variants

3. **Then Add Variants**
   - Variant SKU: `CAP-20MM-WHITE`
   - Selling Price: `45`
   - Cost Price: `20`
   - Stock Quantity: `500`
   - Color: Select `White`
   - Click **"+ Add Variant"**
   
   - Variant SKU: `CAP-20MM-BLUE`
   - Selling Price: `50`
   - Cost Price: `22`
   - Stock Quantity: `300`
   - Color: Select `Blue`
   - Click **"+ Add Variant"**

**API Request:**
```json
{
  "name": "20mm Tamper-proof Bottle Cap - Pack of 100",
  "type": "goods",
  "item_details": {
    "structure": "variants",
    "unit": "pack",
    "sku": "CAP-20MM-BASE",
    "description": "20mm tamper-proof food-grade plastic bottle caps. Perfect for 500ml water bottles. Sold per pack of 100 pieces.",
    "attribute_definitions": [
      {
        "key": "Color",
        "options": ["White", "Blue"]
      }
    ],
    "variants": [
      {
        "sku": "CAP-20MM-WHITE",
        "attribute_map": { "Color": "White" },
        "selling_price": 45,
        "cost_price": 20,
        "stock_quantity": 500
      },
      {
        "sku": "CAP-20MM-BLUE",
        "attribute_map": { "Color": "Blue" },
        "selling_price": 50,
        "cost_price": 22,
        "stock_quantity": 300
      }
    ]
  },
  "sales_info": {
    "account": "Sales Revenue - Packaging",
    "selling_price": 47.5,
    "currency": "INR",
    "description": "20mm tamper caps for 500ml bottles"
  },
  "purchase_info": {
    "account": "Cost of Caps",
    "cost_price": 21,
    "currency": "INR",
    "preferred_vendor_id": 1,
    "description": "Tamper-proof caps from supplier"
  },
  "inventory": {
    "track_inventory": true,
    "inventory_account": "Inventory - Packaging",
    "reorder_point": 100
  },
  "return_policy": {
    "returnable": false
  }
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "item_386700c6",
    "name": "20mm Tamper-proof Bottle Cap - Pack of 100",
    "type": "goods",
    "created_at": "2026-02-17T15:38:20.193+05:30"
  },
  "message": "Item created successfully"
}
```

### 1.3 Create Water Labels Item

**Endpoint:** `POST /items`

**Request:**
```json
{
  "name": "Custom Water Bottle Labels - 1000 pieces",
  "type": "goods",
  "item_details": {
    "structure": "single",
    "unit": "pack",
    "sku": "LBL-1000-CUSTOM",
    "description": "Custom printed waterproof water bottle labels. 1000 pieces per pack. High-quality vinyl, resistant to moisture and UV."
  },
  "sales_info": {
    "account": "Sales Revenue - Packaging",
    "selling_price": 500,
    "currency": "INR",
    "description": "Custom water labels for branding"
  },
  "purchase_info": {
    "account": "Cost of Labels",
    "cost_price": 300,
    "currency": "INR",
    "preferred_vendor_id": 1,
    "description": "Custom label printing from supplier"
  },
  "inventory": {
    "track_inventory": true,
    "inventory_account": "Inventory - Packaging Labels",
    "reorder_point": 50
  },
  "return_policy": {
    "returnable": false
  }
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "item_8195d60a",
    "name": "Custom Water Bottle Labels - 1000 pieces",
    "type": "goods",
    "created_at": "2026-02-17T17:02:25.895+05:30"
  },
  "message": "Item created successfully"
}
```

---

## Step 2: Create Item Group (BOM)

### 2.1 Create 500ml Complete Water Bottle Package

This combines the bottle, cap, and label into one finished product.

**Endpoint:** `POST /item-groups`

**Request:**
```json
{
  "name": "500ml Complete Water Bottle Packaged",
  "description": "Complete packaged 500ml water bottle including cap and label. Ready for retail sale.",
  "is_active": true,
  "components": [
    {
      "item_id": "item_26673b1a",
      "variant_sku": "WTR-BOT-500-REG",
      "quantity": 1,
      "variant_details": {
        "type": "bottle",
        "capacity": "500ml",
        "cap_type": "Regular Cap"
      }
    },
    {
      "item_id": "item_386700c6",
      "variant_sku": "CAP-20MM-WHITE",
      "quantity": 0.01,
      "variant_details": {
        "type": "cap",
        "color": "White"
      }
    },
    {
      "item_id": "item_8195d60a",
      "quantity": 0.001,
      "variant_details": {
        "type": "label"
      }
    }
  ]
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "ig_a1b2c3d4",
    "name": "500ml Complete Water Bottle Packaged",
    "description": "Complete packaged 500ml water bottle including cap and label. Ready for retail sale.",
    "is_active": true,
    "components": [
      {
        "id": 1,
        "item_group_id": "ig_a1b2c3d4",
        "item_id": "item_26673b1a",
        "item": {
          "id": "item_26673b1a",
          "name": "500ml Premium Drinking Water Bottle",
          "sku": "WTR-BOT-500-BASE"
        },
        "variant_sku": "WTR-BOT-500-REG",
        "quantity": 1,
        "variant_details": {
          "type": "bottle",
          "capacity": "500ml",
          "cap_type": "Regular Cap"
        }
      },
      {
        "id": 2,
        "item_group_id": "ig_a1b2c3d4",
        "item_id": "item_386700c6",
        "item": {
          "id": "item_386700c6",
          "name": "20mm Tamper-proof Bottle Cap - Pack of 100",
          "sku": "CAP-20MM-BASE"
        },
        "variant_sku": "CAP-20MM-WHITE",
        "quantity": 0.01,
        "variant_details": {
          "type": "cap",
          "color": "White"
        }
      },
      {
        "id": 3,
        "item_group_id": "ig_a1b2c3d4",
        "item_id": "item_8195d60a",
        "item": {
          "id": "item_8195d60a",
          "name": "Custom Water Bottle Labels - 1000 pieces",
          "sku": "LBL-1000-CUSTOM"
        },
        "quantity": 0.001,
        "variant_details": {
          "type": "label"
        }
      }
    ],
    "created_at": "2026-02-17T11:30:00Z"
  },
  "message": "Item Group created successfully"
}
```

---

## Step 3: Set Opening Stock

### 3.1 Set Initial Stock for the Complete Bottle Package

**Endpoint:** `PUT /items/item_26673b1a/opening-stock`

**Note:** For the item group, we're setting stock for the base bottle variant. The complete packaged bottle would typically be manufactured from components.

**Request:**
```json
{
  "opening_stock": 100,
  "opening_stock_rate_per_unit": 23.75
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "item_26673b1a",
    "name": "500ml Premium Drinking Water Bottle",
    "opening_stock": 100,
    "opening_stock_value": 2375,
    "opening_stock_rate_per_unit": 23.75
  },
  "message": "Opening stock updated successfully"
}
```

---

## Step 4: Create Production Order

### 4.1 Create Production Order for 100 Units

**Endpoint:** `POST /production-orders`

**Request:**
```json
{
  "reference_no": "PROD-2026-500ML-001",
  "date": "2026-02-18",
  "scheduled_completion_date": "2026-02-25",
  "item_group_id": "ig_a1b2c3d4",
  "quantity": 100,
  "warehouse_id": 1,
  "notes": "Production of 100 units of 500ml Complete Water Bottle Packaged. Use white caps and standard labels.",
  "status": "planned"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "prod_xyz789",
    "reference_no": "PROD-2026-500ML-001",
    "item_group_id": "ig_a1b2c3d4",
    "item_group_name": "500ml Complete Water Bottle Packaged",
    "quantity": 100,
    "warehouse_id": 1,
    "status": "planned",
    "scheduled_completion_date": "2026-02-25",
    "components_required": [
      {
        "item_id": "item_26673b1a",
        "item_name": "500ml Premium Drinking Water Bottle",
        "variant_sku": "WTR-BOT-500-REG",
        "quantity_required": 100
      },
      {
        "item_id": "item_386700c6",
        "item_name": "20mm Tamper-proof Bottle Cap - Pack of 100",
        "variant_sku": "CAP-20MM-WHITE",
        "quantity_required": 1
      },
      {
        "item_id": "item_8195d60a",
        "item_name": "Custom Water Bottle Labels - 1000 pieces",
        "quantity_required": 0.1
      }
    ],
    "created_at": "2026-02-18T10:00:00Z"
  },
  "message": "Production Order created successfully"
}
```

### 4.2 Update Production Order Status to Completed

**Endpoint:** `PUT /production-orders/prod_xyz789`

**Request:**
```json
{
  "status": "completed",
  "completed_date": "2026-02-24",
  "notes": "Production completed. 100 units assembled and quality checked."
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "prod_xyz789",
    "reference_no": "PROD-2026-500ML-001",
    "status": "completed",
    "completed_date": "2026-02-24",
    "quantity": 100,
    "created_at": "2026-02-18T10:00:00Z",
    "updated_at": "2026-02-24T14:30:00Z"
  },
  "message": "Production Order updated successfully"
}
```

---

## Step 5: Create Sales Order

### 5.1 Create Sales Order for 100 Units

**Endpoint:** `POST /sales-orders`

**Request:**
```json
{
  "customer_id": 12,
  "reference_no": "SO-2026-500ML-RETAIL-001",
  "order_date": "2026-02-24",
  "delivery_date": "2026-03-03",
  "delivery_address_type": "customer",
  "payment_terms": "net_30",
  "discount": 0,
  "discount_type": "amount",
  "tax_type": "SGST",
  "tax_id": 1,
  "notes": "Order for 100 units of 500ml Complete Water Bottle Packaged for retail distribution.",
  "line_items": [
    {
      "item_id": "item_26673b1a",
      "variant_sku": "WTR-BOT-500-REG",
      "quantity": 100,
      "unit_price": 16.50,
      "description": "500ml water bottles with regular caps"
    }
  ]
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "SO-2026-500ML-RETAIL-001",
    "reference_no": "SO-2026-500ML-RETAIL-001",
    "customer_id": 12,
    "customer": {
      "id": 12,
      "display_name": "Fresh Water Retail"
    },
    "order_date": "2026-02-24",
    "delivery_date": "2026-03-03",
    "status": "draft",
    "line_items": [
      {
        "id": 1,
        "item_id": "item_26673b1a",
        "item_name": "500ml Premium Drinking Water Bottle",
        "variant_sku": "WTR-BOT-500-REG",
        "quantity": 100,
        "unit_price": 16.50,
        "line_total": 1650.00,
        "tax_amount": 297.00,
        "total_with_tax": 1947.00
      }
    ],
    "subtotal": 1650.00,
    "tax_total": 297.00,
    "total": 1947.00,
    "created_at": "2026-02-24T15:00:00Z"
  },
  "message": "Sales Order created successfully"
}
```

### 5.2 Update Sales Order Status to Confirmed

**Endpoint:** `PUT /sales-orders/SO-2026-500ML-RETAIL-001`

**Request:**
```json
{
  "status": "confirmed"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "SO-2026-500ML-RETAIL-001",
    "reference_no": "SO-2026-500ML-RETAIL-001",
    "status": "confirmed",
    "updated_at": "2026-02-24T15:30:00Z"
  },
  "message": "Sales Order confirmed successfully"
}
```

---

## Step 6: Create Invoice

### 6.1 Create Invoice from Sales Order

**Endpoint:** `POST /invoices`

**Request:**
```json
{
  "customer_id": 12,
  "reference_no": "INV-2026-500ML-001",
  "invoice_date": "2026-02-25",
  "due_date": "2026-03-27",
  "sales_order_id": "SO-2026-500ML-RETAIL-001",
  "invoice_status": "draft",
  "payment_method": "bank_transfer",
  "notes": "Invoice for 100 units of 500ml Complete Water Bottle Packaged. Payment due within 30 days.",
  "line_items": [
    {
      "item_id": "item_26673b1a",
      "variant_sku": "WTR-BOT-500-REG",
      "quantity": 100,
      "unit_price": 16.50,
      "description": "500ml water bottles with regular caps"
    }
  ],
  "shipping_address": {
    "attention": "Store Manager",
    "address_line_1": "Fresh Water Retail Store",
    "city": "Bangalore",
    "state": "Karnataka",
    "postal_code": "560034",
    "country": "India"
  }
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "inv_abc123",
    "reference_no": "INV-2026-500ML-001",
    "customer_id": 12,
    "customer": {
      "id": 12,
      "display_name": "Fresh Water Retail"
    },
    "invoice_date": "2026-02-25",
    "due_date": "2026-03-27",
    "invoice_status": "draft",
    "line_items": [
      {
        "id": 1,
        "item_id": "item_26673b1a",
        "item_name": "500ml Premium Drinking Water Bottle",
        "variant_sku": "WTR-BOT-500-REG",
        "quantity": 100,
        "unit_price": 16.50,
        "line_total": 1650.00,
        "tax_amount": 297.00,
        "total_with_tax": 1947.00
      }
    ],
    "subtotal": 1650.00,
    "tax_total": 297.00,
    "total": 1947.00,
    "created_at": "2026-02-25T10:00:00Z"
  },
  "message": "Invoice created successfully"
}
```

### 6.2 Update Invoice Status to Sent

**Endpoint:** `PUT /invoices/inv_abc123`

**Request:**
```json
{
  "invoice_status": "sent"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "inv_abc123",
    "reference_no": "INV-2026-500ML-001",
    "invoice_status": "sent",
    "updated_at": "2026-02-25T11:00:00Z"
  },
  "message": "Invoice sent successfully"
}
```

---

## Step 7: Record Payment

### 7.1 Record Payment Received

**Endpoint:** `POST /invoices/inv_abc123/payments`

**Request:**
```json
{
  "payment_date": "2026-02-28",
  "amount": 1947.00,
  "payment_method": "bank_transfer",
  "reference_no": "BANK-TXN-2026-001",
  "notes": "Payment received from Fresh Water Retail for Invoice INV-2026-500ML-001"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "pay_xyz456",
    "invoice_id": "inv_abc123",
    "payment_date": "2026-02-28",
    "amount": 1947.00,
    "payment_method": "bank_transfer",
    "reference_no": "BANK-TXN-2026-001",
    "status": "completed",
    "created_at": "2026-02-28T14:00:00Z"
  },
  "message": "Payment recorded successfully"
}
```

### 7.2 Update Invoice Status to Paid

**Endpoint:** `PUT /invoices/inv_abc123`

**Request:**
```json
{
  "invoice_status": "paid"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "inv_abc123",
    "reference_no": "INV-2026-500ML-001",
    "invoice_status": "paid",
    "amount_paid": 1947.00,
    "amount_due": 0.00,
    "updated_at": "2026-02-28T14:30:00Z"
  },
  "message": "Invoice marked as paid"
}
```

---

## Workflow Summary

```
┌─────────────────┐
│  Create Items   │
│ (3 base items)  │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Create BOM      │
│  (Item Group)   │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Set Opening     │
│ Stock (100 qty) │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Production      │
│ Order           │
│ (100 units)     │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Sales Order     │
│ (100 units)     │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Create Invoice  │
│ (INV-2026-...)  │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Record Payment  │
│ (Amount Paid)   │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Invoice Status  │
│ = PAID          │
└─────────────────┘
```

---

## Key Metrics

| Metric | Value |
|--------|-------|
| Total Units | 100 |
| Unit Price | ₹16.50 |
| Subtotal | ₹1,650.00 |
| Tax (18% SGST) | ₹297.00 |
| **Total Amount** | **₹1,947.00** |
| Components per Unit | 3 (bottle + cap + label) |
| Production Days | 6 days |
| Delivery Days | 7 days |
| Payment Terms | Net 30 days |

---

## Notes

- All IDs shown are examples from your system
- Actual IDs will be generated by the API endpoints
- Timestamps follow ISO 8601 format with IST timezone (+05:30)
- Status progression: `draft` → `confirmed` → `sent` → `paid`
- Inventory is automatically adjusted after production and sales
- All calculations include applicable taxes
