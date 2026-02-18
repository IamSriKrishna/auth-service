# Water Bottle Production Workflow - Version 3
## 1L Premium Water Bottle: From Items to Final Invoice

This guide demonstrates the complete workflow for manufacturing and selling 50 units of 1L Premium Water Bottles with unique V3 SKUs to avoid database conflicts.

⚠️ **Note:** All SKUs in this version end with `-V3` to prevent duplicate entry errors. If you've already created V1 or V2 items, those IDs will be different from what's shown here.

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
        "key": "Bottle Size",
        "options": ["1 Liter", "1.5 Liter"]
      }
    ],
    "variants": [
      {
        "sku": "WTR-PLT-1L-STD",
        "attribute_map": {
          "Bottle Size": "1 Liter"
        },
        "selling_price": 25,
        "cost_price": 12,
        "stock_quantity": 2000
      }
    ]
  }
}
```

### Error: "Duplicate entry" for variant SKU

**Cause:** The variant SKU already exists in the database. Variant SKUs must be globally unique across all items.

**Solution:**
- Use unique variant SKUs each time you create new items
- Check `/items` endpoint first to see existing items
- Use versioned SKUs: `-V2`, `-V3`, etc. for different iterations
- Or use timestamps: `-20260218`, `-20260219`, etc.

### Error: "variant not found" or "variant_sku not found"

**Cause:** The variant SKU being referenced doesn't exist in the item or was misspelled.

**Solution:** 
- Verify the exact SKU from item creation response
- Ensure case-sensitive matching
- Confirm the variant was successfully created before using it in Item Groups

---

## Backend Implementation Details

### Variant Attribute Validation System

The backend now includes comprehensive validation for variant items through the `ValidateVariantAttributes()` method:

**What Changed:**
- JSON field: `attribute_definitions` (defines available variant attributes)
- Each variant's actual attributes go in: `attribute_map`
- Validation occurs automatically on item creation/update
- Detailed error messages guide users to correct structure

**How It Works:**
1. Validation runs automatically when you POST/PUT an item with `structure: "variants"`
2. Checks that `attribute_definitions` exist with proper keys and options
3. Verifies each variant has matching `attribute_map` with all required attributes
4. Validates attribute values match defined options exactly (case-sensitive)
5. Prevents extra undefined attributes in variants
6. Returns specific error messages before item creation

---

## Table of Contents

0. [Troubleshooting Common Errors](#troubleshooting-common-errors)
1. [Step 1: Create Base Items](#step-1-create-base-items)
2. [Step 2: Create Item Group (BOM)](#step-2-create-item-group-bom)
3. [Step 3: Set Opening Stock](#step-3-set-opening-stock)
4. [Step 4: Create Production Order](#step-4-create-production-order)
5. [Step 5: Create Sales Order](#step-5-create-sales-order)
6. [Step 6: Create Invoice](#step-6-create-invoice)
7. [Step 7: Record Payment](#step-7-record-payment)

---

## Step 1: Create Base Items

### 1.1 Create 1L Water Bottle Item with Variants

**Endpoint:** `POST /items`

**API Request:**
```json
{
  "name": "1L Premium Drinking Water Bottle - V3",
  "type": "goods",
  "item_details": {
    "structure": "variants",
    "unit": "piece",
    "sku": "WTR-PLT-1L-BASE-V3",
    "upc": "8904220100001",
    "ean": "8904220100001",
    "description": "Premium 1 liter PET drinking water bottle. BPA-free, food-grade plastic. Available in standard or eco-friendly variants.",
    "attribute_definitions": [
      {
        "key": "Bottle Variant",
        "options": [
          "Standard",
          "Eco-Friendly"
        ]
      }
    ],
    "variants": [
      {
        "sku": "WTR-PLT-1L-STD-V3",
        "attribute_map": {
          "Bottle Variant": "Standard"
        },
        "selling_price": 25,
        "cost_price": 12,
        "stock_quantity": 2000
      },
      {
        "sku": "WTR-PLT-1L-ECO-V3",
        "attribute_map": {
          "Bottle Variant": "Eco-Friendly"
        },
        "selling_price": 30,
        "cost_price": 14,
        "stock_quantity": 1500
      }
    ]
  },
  "sales_info": {
    "account": "Sales Revenue - Water Bottles",
    "selling_price": 27.5,
    "currency": "INR",
    "description": "1L drinking water bottles retail sales"
  },
  "purchase_info": {
    "account": "Cost of Goods Purchased",
    "cost_price": 13,
    "currency": "INR",
    "preferred_vendor_id": 1,
    "description": "Purchase from bottle manufacturer"
  },
  "inventory": {
    "track_inventory": true,
    "inventory_account": "Inventory - Water Bottles",
    "inventory_valuation_method": "FIFO",
    "reorder_point": 500
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
    "id": "item_plt1l_v3",
    "name": "1L Premium Drinking Water Bottle - V3",
    "type": "goods",
    "item_details": {
      "structure": "variants",
      "sku": "WTR-PLT-1L-BASE-V3",
      "variants": [
        {
          "sku": "WTR-PLT-1L-STD-V3",
          "attribute_map": { "Bottle Variant": "Standard" }
        },
        {
          "sku": "WTR-PLT-1L-ECO",
          "attribute_map": { "Bottle Variant": "Eco-Friendly" }
        }
      ]
    },
    "created_at": "2026-02-18T10:00:00.000+05:30"
  },
  "message": "Item created successfully"
}
```

### 1.2 Create Bottle Caps Item with Variants

**Endpoint:** `POST /items`

**API Request:**
```json
{
  "name": "28mm Twist-Lock Bottle Caps - V3",
  "type": "goods",
  "item_details": {
    "structure": "variants",
    "unit": "pack",
    "sku": "CAP-28MM-BASE-V3",
    "description": "28mm twist-lock food-grade plastic bottle caps. Sold per pack of 200 pieces.",
    "attribute_definitions": [
      {
        "key": "Material",
        "options": ["Standard Plastic", "Recycled Plastic"]
      }
    ],
    "variants": [
      {
        "sku": "CAP-28MM-STD-V3",
        "attribute_map": { "Material": "Standard Plastic" },
        "selling_price": 60,
        "cost_price": 25,
        "stock_quantity": 400
      },
      {
        "sku": "CAP-28MM-REC-V3",
        "attribute_map": { "Material": "Recycled Plastic" },
        "selling_price": 65,
        "cost_price": 28,
        "stock_quantity": 300
      }
    ]
  },
  "sales_info": {
    "account": "Sales Revenue - Packaging",
    "selling_price": 62.5,
    "currency": "INR",
    "description": "28mm twist caps for 1L bottles"
  },
  "purchase_info": {
    "account": "Cost of Caps",
    "cost_price": 26.5,
    "currency": "INR",
    "preferred_vendor_id": 1,
    "description": "Twist-lock caps from supplier"
  },
  "inventory": {
    "track_inventory": true,
    "inventory_account": "Inventory - Packaging",
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
    "id": "item_cap28_v3",
    "name": "28mm Twist-Lock Bottle Caps - V3",
    "type": "goods",
    "created_at": "2026-02-18T10:15:00.000+05:30"
  },
  "message": "Item created successfully"
}
```

### 1.3 Create Box Packaging Item

**Endpoint:** `POST /items`

**API Request:**
```json
{
  "name": "1L Bottle Shipping Box - V3",
  "type": "goods",
  "item_details": {
    "structure": "single",
    "unit": "box",
    "sku": "BOX-1L-SHIP-V3",
    "description": "Corrugated shipping box for 1L water bottles. Holds 6 bottles per box."
  },
  "sales_info": {
    "account": "Sales Revenue - Packaging",
    "selling_price": 30,
    "currency": "INR",
    "description": "Shipping boxes for 1L bottles"
  },
  "purchase_info": {
    "account": "Cost of Packaging",
    "cost_price": 15,
    "currency": "INR",
    "preferred_vendor_id": 1,
    "description": "Corrugated boxes from supplier"
  },
  "inventory": {
    "track_inventory": true,
    "inventory_account": "Inventory - Packaging Boxes",
    "reorder_point": 20
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
    "id": "item_box1l_v3",
    "name": "1L Bottle Shipping Box - V3",
    "type": "goods",
    "created_at": "2026-02-18T10:30:00.000+05:30"
  },
  "message": "Item created successfully"
}
```

---

## Step 2: Create Item Group (BOM)

### 2.1 Create 1L Complete Water Bottle Package

This combines the bottle, cap, and box into one finished product ready for sale.

**Endpoint:** `POST /item-groups`

**API Request:**
```json
{
  "name": "1L Water Bottle - Complete Package",
  "description": "Complete packaged 1L water bottle including twist-lock cap and shipping box. Ready for retail distribution.",
  "is_active": true,
  "components": [
    {
      "item_id": "item_plt1l_v3",
      "variant_sku": "WTR-PLT-1L-STD-V3",
      "quantity": 1,
      "variant_details": {
        "type": "bottle",
        "capacity": "1 Liter",
        "variant_type": "Standard"
      }
    },
    {
      "item_id": "item_cap28_v3",
      "variant_sku": "CAP-28MM-STD-V3",
      "quantity": 0.005,
      "variant_details": {
        "type": "cap",
        "material": "Standard Plastic"
      }
    },
    {
      "item_id": "item_box1l_v3",
      "quantity": 0.1667,
      "variant_details": {
        "type": "box",
        "note": "1/6th of a box per bottle (6 bottles per box)"
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
    "id": "ig_plt1l_v3_pkg",
    "name": "1L Water Bottle - Complete Package",
    "description": "Complete packaged 1L water bottle including twist-lock cap and shipping box. Ready for retail distribution.",
    "is_active": true,
    "components": [
      {
        "id": 1,
        "item_group_id": "ig_plt1l_v3_pkg",
        "item_id": "item_plt1l_v3",
        "item": {
          "id": "item_plt1l_v3",
          "name": "1L Premium Drinking Water Bottle - V3",
          "sku": "WTR-PLT-1L-BASE-V3"
        },
        "variant_sku": "WTR-PLT-1L-STD-V3",
        "quantity": 1,
        "variant_details": {
          "type": "bottle",
          "capacity": "1 Liter",
          "variant_type": "Standard"
        }
      },
      {
        "id": 2,
        "item_group_id": "ig_plt1l_v3_pkg",
        "item_id": "item_cap28_v3",
        "item": {
          "id": "item_cap28_v3",
          "name": "28mm Twist-Lock Bottle Caps - V3",
          "sku": "CAP-28MM-BASE-V3"
        },
        "variant_sku": "CAP-28MM-STD-V3",
        "quantity": 0.005,
        "variant_details": {
          "type": "cap",
          "material": "Standard Plastic"
        }
      },
      {
        "id": 3,
        "item_group_id": "ig_plt1l_v3_pkg",
        "item_id": "item_box1l_v3",
        "item": {
          "id": "item_box1l_v3",
          "name": "1L Bottle Shipping Box - V3",
          "sku": "BOX-1L-SHIP-V3"
        },
        "quantity": 0.1667,
        "variant_details": {
          "type": "box",
          "note": "1/6th of a box per bottle (6 bottles per box)"
        }
      }
    ],
    "created_at": "2026-02-18T11:00:00Z"
  },
  "message": "Item Group created successfully"
}
```

---

## Step 3: Set Opening Stock

### 3.1 Set Initial Stock for the Bottle

**Endpoint:** `PUT /items/item_plt1l_v3/opening-stock`

**Request:**
```json
{
  "opening_stock": 50,
  "opening_stock_rate_per_unit": 12.5
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "item_plt1l_v3",
    "name": "1L Premium Drinking Water Bottle - V3",
    "opening_stock": 50,
    "opening_stock_value": 625,
    "opening_stock_rate_per_unit": 12.5
  },
  "message": "Opening stock updated successfully"
}
```

---

## Step 4: Create Production Order

### 4.1 Create Production Order for 50 Units

**Endpoint:** `POST /production-orders`

**Request:**
```json
{
  "reference_no": "PROD-2026-1L-V3-001",
  "date": "2026-02-18",
  "scheduled_completion_date": "2026-02-23",
  "item_group_id": "ig_plt1l_v3_pkg",
  "quantity": 50,
  "warehouse_id": 1,
  "notes": "Production of 50 units of 1L Complete Water Bottle Package. Use standard bottles and plastic caps.",
  "status": "planned"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "prod_1lv3_001",
    "reference_no": "PROD-2026-1L-V3-001",
    "item_group_id": "ig_plt1l_v3_pkg",
    "item_group_name": "1L Water Bottle - Complete Package",
    "quantity": 50,
    "warehouse_id": 1,
    "status": "planned",
    "scheduled_completion_date": "2026-02-23",
    "components_required": [
      {
        "item_id": "item_plt1l_v3",
        "item_name": "1L Premium Drinking Water Bottle - V3",
        "variant_sku": "WTR-PLT-1L-STD-V3",
        "quantity_required": 50
      },
      {
        "item_id": "item_cap28_v3",
        "item_name": "28mm Twist-Lock Bottle Caps - V3",
        "variant_sku": "CAP-28MM-STD-V3",
        "quantity_required": 0.25
      },
      {
        "item_id": "item_box1l_v3",
        "item_name": "1L Bottle Shipping Box - V3",
        "quantity_required": 8.335
      }
    ],
    "created_at": "2026-02-18T11:30:00Z"
  },
  "message": "Production Order created successfully"
}
```

### 4.2 Mark Production as Completed

**Endpoint:** `PUT /production-orders/prod_1lv3_001`

**Request:**
```json
{
  "status": "completed",
  "completed_date": "2026-02-22",
  "notes": "Production completed successfully. 50 units assembled and quality inspected."
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "prod_1lv3_001",
    "reference_no": "PROD-2026-1L-V3-001",
    "status": "completed",
    "completed_date": "2026-02-22",
    "quantity": 50,
    "created_at": "2026-02-18T11:30:00Z",
    "updated_at": "2026-02-22T16:00:00Z"
  },
  "message": "Production Order updated successfully"
}
```

---

## Step 5: Create Sales Order

### 5.1 Create Sales Order for 50 Units

**Endpoint:** `POST /sales-orders`

**Request:**
```json
{
  "customer_id": 5,
  "reference_no": "SO-2026-1L-V3-001",
  "order_date": "2026-02-22",
  "delivery_date": "2026-03-01",
  "delivery_address_type": "customer",
  "payment_terms": "net_30",
  "discount": 0,
  "discount_type": "amount",
  "tax_type": "SGST",
  "tax_id": 1,
  "notes": "Order for 50 units of 1L Complete Water Bottle Package for distribution.",
  "line_items": [
    {
      "item_id": "item_plt1l_v3",
      "variant_sku": "WTR-PLT-1L-STD-V3",
      "quantity": 50,
      "unit_price": 30,
      "description": "1L water bottles with standard plastic caps"
    }
  ]
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "SO-2026-1L-V3-001",
    "reference_no": "SO-2026-1L-V3-001",
    "customer_id": 5,
    "customer": {
      "id": 5,
      "display_name": "Water Distribution Partners"
    },
    "order_date": "2026-02-22",
    "delivery_date": "2026-03-01",
    "status": "draft",
    "line_items": [
      {
        "id": 1,
        "item_id": "item_plt1l_v3",
        "item_name": "1L Premium Drinking Water Bottle - V3",
        "variant_sku": "WTR-PLT-1L-STD-V3",
        "quantity": 50,
        "unit_price": 30,
        "line_total": 1500,
        "tax_amount": 270,
        "total_with_tax": 1770
      }
    ],
    "subtotal": 1500,
    "tax_total": 270,
    "total": 1770,
    "created_at": "2026-02-22T10:00:00Z"
  },
  "message": "Sales Order created successfully"
}
```

### 5.2 Confirm Sales Order

**Endpoint:** `PUT /sales-orders/SO-2026-1L-V3-001`

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
    "id": "SO-2026-1L-V3-001",
    "reference_no": "SO-2026-1L-V3-001",
    "status": "confirmed",
    "updated_at": "2026-02-22T10:30:00Z"
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
  "customer_id": 5,
  "reference_no": "INV-2026-1L-V3-001",
  "invoice_date": "2026-02-22",
  "due_date": "2026-03-24",
  "sales_order_id": "SO-2026-1L-V3-001",
  "invoice_status": "draft",
  "payment_method": "bank_transfer",
  "notes": "Invoice for 50 units of 1L Complete Water Bottle Package. Payment due within 30 days.",
  "line_items": [
    {
      "item_id": "item_plt1l_v3",
      "variant_sku": "WTR-PLT-1L-STD-V3",
      "quantity": 50,
      "unit_price": 30,
      "description": "1L water bottles with standard plastic caps"
    }
  ],
  "shipping_address": {
    "attention": "Distribution Manager",
    "address_line_1": "Water Distribution Warehouse",
    "city": "Mumbai",
    "state": "Maharashtra",
    "postal_code": "400001",
    "country": "India"
  }
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "inv_1lv3_001",
    "reference_no": "INV-2026-1L-V3-001",
    "customer_id": 5,
    "customer": {
      "id": 5,
      "display_name": "Water Distribution Partners"
    },
    "invoice_date": "2026-02-22",
    "due_date": "2026-03-24",
    "invoice_status": "draft",
    "line_items": [
      {
        "id": 1,
        "item_id": "item_plt1l_v3",
        "item_name": "1L Premium Drinking Water Bottle - V3",
        "variant_sku": "WTR-PLT-1L-STD-V3",
        "quantity": 50,
        "unit_price": 30,
        "line_total": 1500,
        "tax_amount": 270,
        "total_with_tax": 1770
      }
    ],
    "subtotal": 1500,
    "tax_total": 270,
    "total": 1770,
    "created_at": "2026-02-22T11:00:00Z"
  },
  "message": "Invoice created successfully"
}
```

### 6.2 Send Invoice

**Endpoint:** `PUT /invoices/inv_1lv3_001`

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
    "id": "inv_1lv3_001",
    "reference_no": "INV-2026-1L-V3-001",
    "invoice_status": "sent",
    "updated_at": "2026-02-22T11:30:00Z"
  },
  "message": "Invoice sent successfully"
}
```

---

## Step 7: Record Payment

### 7.1 Record Payment Received

**Endpoint:** `POST /invoices/inv_1lv3_001/payments`

**Request:**
```json
{
  "payment_date": "2026-02-28",
  "amount": 1770,
  "payment_method": "bank_transfer",
  "reference_no": "BANK-TXN-2026-1LV3",
  "notes": "Payment received from Water Distribution Partners for Invoice INV-2026-1L-V3-001"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "pay_1lv3_001",
    "invoice_id": "inv_1lv3_001",
    "payment_date": "2026-02-28",
    "amount": 1770,
    "payment_method": "bank_transfer",
    "reference_no": "BANK-TXN-2026-1LV2",
    "status": "completed",
    "created_at": "2026-02-28T14:00:00Z"
  },
  "message": "Payment recorded successfully"
}
```

### 7.2 Mark Invoice as Paid

**Endpoint:** `PUT /invoices/inv_1lv3_001`

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
    "id": "inv_1lv3_001",
    "reference_no": "INV-2026-1L-V3-001",
    "invoice_status": "paid",
    "amount_paid": 1770,
    "amount_due": 0,
    "updated_at": "2026-02-28T14:30:00Z"
  },
  "message": "Invoice marked as paid"
}
```

---

## Workflow Summary

```
┌─────────────────────┐
│ Create 3 Base Items │
│  (Fresh SKUs)       │
└────────┬────────────┘
         │
         ▼
┌─────────────────────┐
│ Create Item Group   │
│  (BOM/Package)      │
└────────┬────────────┘
         │
         ▼
┌─────────────────────┐
│ Set Opening Stock   │
│  (50 units)         │
└────────┬────────────┘
         │
         ▼
┌─────────────────────┐
│ Production Order    │
│ (Plan → Complete)   │
└────────┬────────────┘
         │
         ▼
┌─────────────────────┐
│ Sales Order         │
│ (Draft → Confirm)   │
└────────┬────────────┘
         │
         ▼
┌─────────────────────┐
│ Create Invoice      │
│ (Draft → Send)      │
└────────┬────────────┘
         │
         ▼
┌─────────────────────┐
│ Record Payment      │
│ & Mark Paid         │
└─────────────────────┘
```

---

## Quick Reference: Items & SKUs

| Item Name | Base SKU | Variant SKU 1 | Variant SKU 2 |
|-----------|----------|---------------|---------------|
| 1L Water Bottle | WTR-PLT-1L-BASE-V3 | WTR-PLT-1L-STD-V3 | WTR-PLT-1L-ECO-V3 |
| 28mm Caps | CAP-28MM-BASE-V3 | CAP-28MM-STD-V3 | CAP-28MM-REC-V3 |
| Shipping Box | BOX-1L-SHIP-V3 | (single item - no variants) | |

---

## Key Metrics

| Metric | Value |
|--------|-------|
| **Order Quantity** | 50 units |
| **Unit Price** | ₹30.00 |
| **Subtotal** | ₹1,500.00 |
| **Tax (18% SGST)** | ₹270.00 |
| **Total Amount** | ₹1,770.00 |
| **Components per Unit** | 3 (bottle + cap + box portion) |
| **Production Timeline** | 4 days |
| **Delivery Timeline** | 7 days |
| **Payment Terms** | Net 30 days |

---

## Notes

- All IDs shown are examples from your system
- Actual IDs will be generated by API endpoints
- Timestamps use ISO 8601 format with IST (+05:30)
- Status flow: `draft` → `confirmed`/`planned` → `completed` → `sent` → `paid`
- All prices and taxes automatically calculated
- Inventory tracked automatically through production and sales
