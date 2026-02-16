# Item Group & Manufacturing System Architecture

## Overview

This document explains the new **Item Group** (Bill of Materials/BOM) and **Manufacturing** system that has been added to your auth-service application.

---

## Key Concepts

### 1. **Item Model Changes**
- ✅ **Removed**: `Brand` field
- ✅ **Removed**: `ManufacturerID` relationship
- **Reason**: Simplified model to make items more generic. Variants handle different specifications.

### 2. **Item Variants** (Already Existed)
- Items can have multiple variants (e.g., 300ml, 500ml, 1000ml water bottles)
- Each variant has:
  - SKU (unique identifier)
  - Attributes (capacity, color, size, etc.)
  - Selling price & Cost price
  - Stock quantity

### 3. **ItemGroup (New)** - Bill of Materials (BOM)
A collection of items that work together to create a final product.

**Example: 300ml Water Bottle ItemGroup**
```
ItemGroup: "300ml Water Bottle"
├── ItemGroupComponent 1:
│   ├── Item: "Bottle"
│   ├── Variant: "300ml Bottle" 
│   └── Quantity: 1
└── ItemGroupComponent 2:
    ├── Item: "Cap"
    ├── Variant: "Cap Size 20mm"
    └── Quantity: 1
```

### 4. **ProductionOrder (New)** - Manufacturing Orders
When you want to manufacture a product using an ItemGroup:

**Example: Manufacture 100 units of "300ml Water Bottle"**
```
ProductionOrder:
├── ItemGroup: "300ml Water Bottle"
├── QuantityToManufacture: 100
├── Status: "planned" → "in_progress" → "completed"
├── PlannedStartDate: 2026-02-20
├── PlannedEndDate: 2026-02-25
└── ProductionOrderItems:
    ├── Item: Bottle (needs 100 × 1 = 100 bottles)
    └── Item: Cap (needs 100 × 1 = 100 caps)
```

### 5. **Inventory Tracking Models (New)**

#### **InventoryBalance**
Real-time inventory status:
- `CurrentQuantity`: Total available
- `ReservedQuantity`: Allocated to open orders
- `AvailableQuantity`: current - reserved
- Tracks `LastReceivedDate`, `LastConsumedDate`, `LastSoldDate`

#### **InventoryAggregation**
Summary metrics:
- `TotalPurchased`: From all purchase orders
- `TotalManufactured`: From all production orders
- `TotalConsumedInMfg`: Used as components in manufacturing
- `TotalSold`: From all sales orders

#### **InventoryJournal**
Complete audit trail of every inventory transaction:
- Transaction types: `purchase`, `manufacture`, `consume`, `sale`, `adjustment`
- Links to reference documents (PO, ProductionOrder, SalesOrder)
- Tracks who made the change and when

#### **SupplyChainSummary**
Complete view of item flow:
- Opening Stock
- Purchase metrics (quantity, amount, average rate)
- Manufacturing metrics (produced, consumed)
- Sales metrics (quantity, amount, average rate)
- Current quantity

---

## Step-by-Step Usage Example

### Step 1: Create Items (Base Products)

Create two items with variants:

**Item 1: Bottle**
```json
{
  "name": "Plastic Bottle",
  "type": "goods",
  "item_details": {
    "structure": "variants",
    "unit": "pieces",
    "attribute_definitions": [
      {
        "key": "capacity",
        "options": ["300ml", "500ml", "1000ml"]
      }
    ]
  }
}
```

**Item 2: Cap**
```json
{
  "name": "Bottle Cap",
  "type": "goods",
  "item_details": {
    "structure": "variants",
    "unit": "pieces",
    "attribute_definitions": [
      {
        "key": "size",
        "options": ["20mm", "25mm", "28mm"]
      }
    ]
  }
}
```

### Step 2: Create ItemGroups (Products)

**ItemGroup 1: 300ml Water Bottle**
```json
{
  "name": "300ml Water Bottle",
  "description": "Ready-to-sell 300ml water bottle package",
  "is_active": true,
  "components": [
    {
      "item_id": "bottle_001",
      "variant_id": 1,  // 300ml Bottle variant
      "quantity": 1,
      "variant_details": {"capacity": "300ml"}
    },
    {
      "item_id": "cap_001",
      "variant_id": 2,  // 20mm Cap variant
      "quantity": 1,
      "variant_details": {"size": "20mm"}
    }
  ]
}
```

**ItemGroup 2: 1000ml Water Bottle**
```json
{
  "name": "1000ml Water Bottle",
  "description": "Ready-to-sell 1000ml water bottle package",
  "is_active": true,
  "components": [
    {
      "item_id": "bottle_001",
      "variant_id": 3,  // 1000ml Bottle variant
      "quantity": 1,
      "variant_details": {"capacity": "1000ml"}
    },
    {
      "item_id": "cap_001",
      "variant_id": 3,  // 28mm Cap variant (larger for bigger bottle)
      "quantity": 1,
      "variant_details": {"size": "28mm"}
    }
  ]
}
```

### Step 3: Create Purchase Orders for Components

**PurchaseOrder: Get Bottles**
```json
{
  "vendor_id": 5,
  "line_items": [
    {
      "item_id": "bottle_001",
      "variant_id": 1,  // 300ml
      "quantity": 100,
      "rate": 2.50
    },
    {
      "item_id": "bottle_001",
      "variant_id": 3,  // 1000ml
      "quantity": 100,
      "rate": 3.50
    }
  ]
}
```

**PurchaseOrder: Get Caps**
```json
{
  "vendor_id": 6,
  "line_items": [
    {
      "item_id": "cap_001",
      "variant_id": 2,  // 20mm
      "quantity": 100,
      "rate": 0.50
    },
    {
      "item_id": "cap_001",
      "variant_id": 3,  // 28mm
      "quantity": 100,
      "rate": 0.60
    }
  ]
}
```

### Step 4: Create ProductionOrder to Manufacture

**ProductionOrder: Manufacture 300ml Water Bottles**
```json
{
  "item_group_id": "grp_300ml_bottle",
  "quantity_to_manufacture": 100,
  "planned_start_date": "2026-02-20",
  "planned_end_date": "2026-02-25",
  "notes": "Standard production run"
}
```

**What happens automatically:**
1. ProductionOrder creates ProductionOrderItems:
   - Needs 100 bottles (300ml variant)
   - Needs 100 caps (20mm variant)
2. Inventory System:
   - Reserves 100 bottles from purchased stock
   - Reserves 100 caps from purchased stock
3. When production completes:
   - Inventory marked as "consumed" in manufacturing
   - New inventory created for finished "300ml Water Bottle" ItemGroup product

### Step 5: Create SalesOrder for Finished Product

**SalesOrder: Sell to Customer**
```json
{
  "customer_id": 10,
  "line_items": [
    {
      "item_id": "grp_300ml_bottle",  // Can now be sold as a finished product
      "quantity": 50,
      "rate": 6.00
    }
  ]
}
```

### Step 6: Create Invoice

**Invoice: Bill Customer**
```json
{
  "customer_id": 10,
  "sales_order_id": "so_001",
  "line_items": [
    {
      "item_id": "grp_300ml_bottle",
      "quantity": 50,
      "rate": 6.00
    }
  ]
}
```

---

## Inventory Tracking Example

### After completing all steps above:

**InventoryBalance for Bottle (300ml variant):**
```
ItemID: bottle_001
VariantID: 1
PurchaseQuantity: 100 (from PurchaseOrder)
ReservedQuantity: 100 (for ProductionOrder)
AvailableQuantity: 0
LastReceivedDate: 2026-02-17
LastConsumedDate: 2026-02-25 (when manufactured)
```

**InventoryBalance for 300ml Water Bottle ItemGroup:**
```
ItemID: grp_300ml_bottle
CurrentQuantity: 100 (manufactured)
ReservedQuantity: 50 (for SalesOrder)
AvailableQuantity: 50
LastReceivedDate: 2026-02-25 (manufactured)
LastSoldDate: 2026-02-28 (invoiced)
```

**InventoryAggregation for Bottle (300ml):**
```
TotalPurchased: 100
TotalManufactured: 0 (it's a component, not a finished product)
TotalConsumedInMfg: 100 (consumed in ProductionOrder)
TotalSold: 0
```

**InventoryAggregation for 300ml Water Bottle:**
```
TotalPurchased: 0 (it's manufactured, not purchased)
TotalManufactured: 100 (from ProductionOrder)
TotalConsumedInMfg: 0 (it's not a component)
TotalSold: 50 (from Invoice)
```

---

## Database Schema Summary

### New Tables Created:
1. **item_groups** - ItemGroup models
2. **item_group_components** - ItemGroupComponent models
3. **production_orders** - ProductionOrder models
4. **production_order_items** - ProductionOrderItem models
5. **inventory_balances** - Current inventory status
6. **inventory_aggregations** - Summary metrics
7. **inventory_journals** - Complete audit trail
8. **supply_chain_summary** - Overview of item flow

### Modified Tables:
- **items** - Removed `brand` and `manufacturer_id` columns

---

## Service Layer Implementation (Next Steps)

You'll need to create services for:

### 1. **ItemGroupService**
- CRUD operations for ItemGroups
- Validate components exist
- Calculate component requirements for given quantity

### 2. **ProductionOrderService**
- Create and manage production orders
- Calculate required component quantities
- Update inventory when production starts/completes
- Track production status

### 3. **InventoryService**
- Manage InventoryBalance
- Manage InventoryAggregation (updates after transactions)
- Record transactions in InventoryJournal
- Calculate SupplyChainSummary
- Reserve/unreserve inventory

---

## API Endpoints (Suggested)

### ItemGroup Endpoints
```
POST   /api/item-groups              - Create ItemGroup
GET    /api/item-groups              - List ItemGroups
GET    /api/item-groups/:id          - Get ItemGroup details
PUT    /api/item-groups/:id          - Update ItemGroup
DELETE /api/item-groups/:id          - Delete ItemGroup
```

### ProductionOrder Endpoints
```
POST   /api/production-orders        - Create ProductionOrder
GET    /api/production-orders        - List ProductionOrders
GET    /api/production-orders/:id    - Get ProductionOrder details
PUT    /api/production-orders/:id    - Update ProductionOrder status
```

### Inventory Endpoints
```
GET    /api/inventory/balance/:item_id                    - Get current balance
GET    /api/inventory/aggregation/:item_id                - Get aggregated metrics
GET    /api/inventory/journal?item_id=&date_from=&date_to - Get audit trail
GET    /api/supply-chain/summary/:item_id                 - Get supply chain view
```

---

## Key Relationships

```
Item
├── ItemDetails
│   └── Variant (multiple)
│       └── VariantAttribute
├── SalesInfo
├── PurchaseInfo
├── Inventory
└── ReturnPolicy

ItemGroup
└── ItemGroupComponent (multiple)
    ├── Item
    └── Variant

ProductionOrder
├── ItemGroup
└── ProductionOrderItem (multiple)
    └── ItemGroupComponent

PurchaseOrder
└── PurchaseOrderLineItem
    ├── Item
    └── Variant

SalesOrder
└── SalesOrderLineItem
    ├── Item (can be Item or ItemGroup)
    └── Variant

Invoice
└── InvoiceLineItem
    ├── Item (can be Item or ItemGroup)
    └── Variant

InventoryBalance
├── Item
└── Variant

InventoryJournal
└── (references any document type)

SupplyChainSummary
├── Item
└── Variant
```

---

## State Transitions

### ProductionOrder Status Flow
```
planned → in_progress → completed
            ↓
        cancelled (from any state)
```

---

## Calculation Examples

### Calculate Required Component Quantities
If you manufacture 100 units of "300ml Water Bottle" ItemGroup:
- Bottles needed: 100 × 1 = 100 bottles
- Caps needed: 100 × 1 = 100 caps

### Calculate Available Inventory
```
AvailableQuantity = CurrentQuantity - ReservedQuantity
```

### Track Supply Chain
```
For "300ml Water Bottle":
- Purchased: 0 (it's a finished good)
- Manufactured: 100
- Sold: 50
- Current Stock: 50

For "Bottle (300ml variant)":
- Purchased: 100
- Consumed in Manufacturing: 100
- Current Stock: 0
```

---

## Next Steps

1. Create repository interfaces for new models
2. Implement repositories for ItemGroup and ProductionOrder
3. Create services with business logic
4. Create DTOs (Input/Output) for API requests/responses
5. Create handlers for REST endpoints
6. Update routes to include new endpoints
7. Add comprehensive logging and error handling
8. Create unit tests for services and repositories

---

## Notes

- The system uses GORM AutoMigrate, so no manual SQL migrations needed
- Set `DROP_ITEM_TABLES=true` env variable to reset item tables during development
- All timestamps are tracked (createdAt, updatedAt)
- All changes are tracked (createdBy, updatedBy)
- Inventory movements can be audited through InventoryJournal table
