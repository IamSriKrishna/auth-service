# Quick Reference: Item Group & Manufacturing System

## What's Done ✅

### Database Models Created
```
✅ ItemGroup              - Bill of Materials definition
✅ ItemGroupComponent     - Individual components in a BOM
✅ ProductionOrder        - Manufacturing order
✅ ProductionOrderItem    - Component requirements for production
✅ InventoryBalance       - Current inventory status
✅ InventoryAggregation   - Summary metrics
✅ InventoryJournal       - Audit trail
✅ SupplyChainSummary     - Complete overview
```

### Item Model Updated
```
✅ Removed: Brand field
✅ Removed: ManufacturerID field
✅ Variants now handle all specifications
```

### Domain Types Added
```
✅ ProductionOrderStatus (planned, in_progress, completed, cancelled)
```

### Migrations
```
✅ All models added to AutoMigrate
✅ Foreign key dependencies configured
✅ Drop functions updated
```

---

## Key Models at a Glance

### ItemGroup (Bill of Materials)
```go
ItemGroup {
    ID: "grp_300ml_bottle"
    Name: "300ml Water Bottle"
    Components: [
        {item: Bottle, variant: 300ml, qty: 1},
        {item: Cap, variant: 20mm, qty: 1}
    ]
}
```

### ProductionOrder (Manufacturing)
```go
ProductionOrder {
    ID: "po_mfg_001"
    ItemGroupID: "grp_300ml_bottle"
    QuantityToManufacture: 100
    Status: "planned" → "in_progress" → "completed"
    PlannedDates: 2026-02-20 to 2026-02-25
}
```

### Inventory Tracking
```
InventoryBalance:
├── Current: 150 units
├── Reserved: 50 units (for pending orders)
└── Available: 100 units

InventoryAggregation:
├── TotalPurchased: 500
├── TotalManufactured: 300
├── TotalSold: 250
└── TotalConsumedInMfg: 100

InventoryJournal:
└── [Audit trail of all transactions]

SupplyChainSummary:
└── Complete view of opening stock, purchases, manufacturing, sales
```

---

## How to Use

### 1. Create ItemGroup
```go
ig := &models.ItemGroup{
    ID: "grp_300ml",
    Name: "300ml Water Bottle",
    Components: []models.ItemGroupComponent{
        {ItemID: "bottle", VariantID: 1, Quantity: 1},
        {ItemID: "cap", VariantID: 2, Quantity: 1},
    },
}
db.Create(ig)
```

### 2. Create ProductionOrder
```go
po := &models.ProductionOrder{
    ID: "po_001",
    ProductionOrderNumber: "PO-001",
    ItemGroupID: "grp_300ml",
    QuantityToManufacture: 100,
    Status: domain.ProductionOrderStatusPlanned,
}
db.Create(po)
```

### 3. Track Inventory
```go
// Get current balance
var balance models.InventoryBalance
db.Where("item_id = ? AND variant_id = ?", "bottle", 1).First(&balance)

// Get supply chain view
var summary models.SupplyChainSummary
db.Where("item_id = ?", "bottle").First(&summary)
```

---

## Database Schema Overview

```
ITEM HIERARCHY:
items
├── item_details (variants)
│   └── variant_attributes
├── sales_info
├── purchase_info
├── inventory

PRODUCT ASSEMBLY:
item_groups
└── item_group_components → links to items

MANUFACTURING:
production_orders
└── production_order_items → links to components

INVENTORY TRACKING:
├── inventory_balances (current status)
├── inventory_aggregations (summaries)
├── inventory_journals (audit trail)
└── supply_chain_summary (complete view)
```

---

## Files Structure

### Models (DONE ✅)
```
app/models/
├── item.go                (MODIFIED - removed brand, manufacturer)
├── item_group.go          (NEW - ItemGroup, ItemGroupComponent)
├── production_order.go    (NEW - ProductionOrder, ProductionOrderItem)
└── inventory_tracking.go  (NEW - All inventory tracking models)
```

### Domain (DONE ✅)
```
app/domain/
└── invoice.domain.go      (MODIFIED - added ProductionOrderStatus)
```

### Migrations (DONE ✅)
```
app/helper/
└── migrations.go          (MODIFIED - added all models)
```

### Documentation (DONE ✅)
```
ROOT/
├── ITEM_GROUP_MANUFACTURING_GUIDE.md     (Complete usage guide)
├── IMPLEMENTATION_SUMMARY.md              (What was done)
├── NEXT_IMPLEMENTATION_GUIDE.md           (DTOs, interfaces, templates)
└── QUICK_REFERENCE.md                     (This file)
```

---

## Next Steps

### Phase 1: Repositories (Estimated: 4 hours)
```
□ ItemGroupRepository
□ ProductionOrderRepository
□ InventoryRepository
└── All CRUD operations + custom queries
```

### Phase 2: DTOs (Estimated: 2 hours)
```
□ app/dto/input/item_group.input.go
□ app/dto/input/production_order.input.go
□ app/dto/output/item_group.output.go
□ app/dto/output/production_order.output.go
└── app/dto/output/inventory.output.go
```

### Phase 3: Services (Estimated: 6 hours)
```
□ ItemGroupService
□ ProductionOrderService
└── InventoryService
```

### Phase 4: Handlers & Routes (Estimated: 4 hours)
```
□ ItemGroupHandler + routes
□ ProductionOrderHandler + routes
└── InventoryHandler + routes
```

### Phase 5: Testing (Estimated: 4 hours)
```
□ Unit tests
□ Integration tests
└── API tests
```

**Total Estimated Time**: ~20 hours

---

## API Endpoints (To Implement)

### ItemGroup API
```
POST   /api/item-groups                     - Create
GET    /api/item-groups                     - List
GET    /api/item-groups/:id                 - Get
PUT    /api/item-groups/:id                 - Update
DELETE /api/item-groups/:id                 - Delete
```

### ProductionOrder API
```
POST   /api/production-orders               - Create
GET    /api/production-orders               - List
GET    /api/production-orders/:id           - Get
PUT    /api/production-orders/:id           - Update
PUT    /api/production-orders/:id/status    - Update status
POST   /api/production-orders/:id/start     - Start
POST   /api/production-orders/:id/complete  - Complete
```

### Inventory API
```
GET    /api/inventory/balance/:item_id                    - Balance
GET    /api/inventory/aggregation/:item_id                - Aggregation
GET    /api/inventory/journal/:item_id                    - Journal
PUT    /api/inventory/balance/:item_id/reserve            - Reserve
PUT    /api/inventory/balance/:item_id/release            - Release
GET    /api/supply-chain/summary/:item_id                 - Summary
```

---

## Real-World Workflow Example

### Step 1: Define Product (ItemGroup)
```
ItemGroup: "300ml Water Bottle"
├── 1 × Bottle (300ml variant)
└── 1 × Cap (20mm variant)
```

### Step 2: Purchase Components
```
PurchaseOrder 1:
├── 100 × Bottle (300ml)
└── Rate: $2.50 each

PurchaseOrder 2:
├── 100 × Cap (20mm)
└── Rate: $0.50 each

Result: InventoryBalance updated
├── Bottle: 100 available
└── Cap: 100 available
```

### Step 3: Manufacture Product
```
ProductionOrder:
├── Manufacture: 100 × "300ml Water Bottle"
├── Status: planned → in_progress → completed
├── Consume: 100 × Bottle (300ml)
└── Consume: 100 × Cap (20mm)

Result: 
├── Bottle: 0 available (consumed)
├── Cap: 0 available (consumed)
└── "300ml Water Bottle": 100 available (manufactured)
```

### Step 4: Sell Product
```
SalesOrder:
├── 50 × "300ml Water Bottle"
└── Rate: $6.00 each

Invoice:
└── Same 50 × "300ml Water Bottle"

Result:
└── "300ml Water Bottle": 50 available (50 sold)
```

### Step 5: View Metrics
```
SupplyChainSummary for Bottle (300ml):
├── Opening: 0
├── Purchased: 100
├── Manufactured: 0
├── Consumed in Mfg: 100
├── Sold: 0
└── Current: 0

SupplyChainSummary for "300ml Water Bottle":
├── Opening: 0
├── Purchased: 0
├── Manufactured: 100
├── Consumed in Mfg: 0
├── Sold: 50
└── Current: 50
```

---

## Important Notes

### Variants Are Key
- Items have variants (300ml, 500ml, 1000ml)
- Components in ItemGroups link to specific variants
- Inventory is tracked per variant
- This allows easy management of different specifications

### Flexible Components
- Components can be any Item + Variant combo
- Quantity can be decimal (for fractional usage)
- VariantDetails stores human-readable info (capacity: 300ml)

### Inventory Precision
- InventoryBalance = Real-time current status
- InventoryAggregation = Summary metrics
- InventoryJournal = Complete audit trail
- SupplyChainSummary = Business metrics

### Status Management
- ProductionOrder statuses are explicit
- Clear state transitions (planned → in_progress → completed)
- Can cancel from any state if needed

---

## Testing the System

### Manual Testing Steps
1. Create an Item with variants
2. Create an ItemGroup with components
3. Create a PurchaseOrder for components
4. Create a ProductionOrder for ItemGroup
5. Check InventoryBalance at each step
6. Create a SalesOrder for the ItemGroup product
7. Create an Invoice
8. View SupplyChainSummary for complete overview

### Expected Results
- Components reserved during manufacturing
- Components consumed when production completes
- New product inventory created
- Sales reduce new product inventory
- All tracked in InventoryJournal

---

## File References

| File | Purpose |
|------|---------|
| [ITEM_GROUP_MANUFACTURING_GUIDE.md](ITEM_GROUP_MANUFACTURING_GUIDE.md) | Complete guide with examples |
| [IMPLEMENTATION_SUMMARY.md](IMPLEMENTATION_SUMMARY.md) | What was implemented |
| [NEXT_IMPLEMENTATION_GUIDE.md](NEXT_IMPLEMENTATION_GUIDE.md) | DTOs, interfaces, templates |
| [QUICK_REFERENCE.md](QUICK_REFERENCE.md) | This quick reference |

---

## Support

If you need help implementing the next phases:
1. Refer to NEXT_IMPLEMENTATION_GUIDE.md for templates
2. Check the models in app/models/ for field names
3. Review ITEM_GROUP_MANUFACTURING_GUIDE.md for workflows
4. Use existing services as patterns (e.g., ItemService)

---

## Summary

🎉 **Database Models**: Complete!  
📊 **Inventory Tracking**: Complete!  
🔧 **Manufacturing Order**: Complete!  
📦 **ItemGroup (BOM)**: Complete!  

Ready for: Repository → Service → Handler → API implementation
