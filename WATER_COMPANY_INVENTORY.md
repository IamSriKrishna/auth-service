# Water Company Inventory Setup Guide
## Complete Product Catalog with DTO Inputs & Outputs

This guide demonstrates the complete water company inventory setup following Zoho inventory patterns, with exact DTO request/response examples.

---

## Table of Contents (Zoho Inventory Sequence)

### Master Setup Phase
1. [Step 1: Create Vendor (Suppliers)](#step-1-create-vendor-suppliers)
2. [Step 2: Create Customer (Buyers)](#step-2-create-customer-buyers)

### Items Management Phase
3. [Step 3: Create 500ml Water Bottle](#step-3-create-500ml-water-bottle)
4. [Step 4: Create 20L Water Cooler Bottle](#step-4-create-20l-water-cooler-bottle)
5. [Step 5: Create Plastic Caps (by Size)](#step-5-create-plastic-caps-by-size)
6. [Step 6: Create Water Labels](#step-6-create-water-labels)

### Inventory Management Phase
7. [Step 7: Set Opening Stock](#step-7-set-opening-stock)
8. [Step 8: Create Item Groups (BOM)](#step-8-create-item-groups-bom)

### Purchase Workflow Phase
9. [Step 9: Create Purchase Order](#step-9-create-purchase-order)
10. [Step 10: Create Bill (Vendor Invoice)](#step-10-create-bill-vendor-invoice)
11. [Step 11: Record Payment Made](#step-11-record-payment-made)

### Sales Workflow Phase
12. [Step 12: Create Sales Order](#step-12-create-sales-order)
13. [Step 13: Create Invoice](#step-13-create-invoice)
14. [Step 14: Create Shipment](#step-14-create-shipment)
15. [Step 15: Record Payment Received](#step-15-record-payment-received)

### Reference
16. [Complete API Reference](#complete-api-reference)

---

## Product Overview

### Water Company Product Hierarchy

```
WATER PRODUCTS
├── Bottles
│   ├── 500ml PET Bottles (with variants: Regular Cap, Sports Cap)
│   ├── 1L PET Bottles
│   ├── 2L PET Bottles
│   └── 20L Polycarbonate Cooler Bottles
├── Caps & Closures
│   ├── 20mm Tamper-proof Caps (for 500ml)
│   ├── 28mm Caps (for 1L and 2L)
│   └── 90mm Large Caps (for 20L bottles)
├── Labels & Packaging
│   ├── Custom Water Labels
│   ├── Packaging Boxes
│   └── Shipping Cartons
└── Accessories
    ├── Water Cooler Stands
    └── Bottle Crates
```

---

## Step 1: Create Vendor (Suppliers)

### Zoho Menu Path: Purchases → Vendors

Vendors are suppliers of water bottles, caps, labels, and packaging materials.

### 1.1 Vendor Details

**Vendor Name:** AquaPlast Industries  
**Business Type:** Bottle Manufacturer  
**Location:** Bangalore, Karnataka  
**Contact:** supplier@aquaplast.com

### 1.2 DTO Input Request (Create Vendor)

**Endpoint:** `POST /vendors`  
**Authentication:** Bearer Token + SuperAdmin Role

```json
{
  "salutation": "Mr.",
  "first_name": "Rajesh",
  "last_name": "Kumar",
  "company_name": "AquaPlast Industries Pvt Ltd",
  "display_name": "AquaPlast Industries",
  "email_address": "rajesh.kumar@aquaplast.com",
  "work_phone": "08041234567",
  "work_phone_code": "+91",
  "mobile": "9876543210",
  "mobile_code": "+91",
  "vendor_language": "English",
  "other_details": {
    "pan": "AABCT5678H",
    "is_msme_registered": true,
    "currency": "INR",
    "payment_terms": "Net 45",
    "tds": "2%",
    "enable_portal": true,
    "website_url": "https://www.aquaplast.com",
    "department": "Sales",
    "designation": "Regional Manager"
  },
  "billing_address": {
    "attention": "Accounts Department",
    "street": "123 Industrial Estate",
    "address_line2": "Block A",
    "city": "Bangalore",
    "state": "Karnataka",
    "country": "India",
    "postal_code": "560001",
    "phone": "08041234567",
    "phone_code": "+91"
  },
  "shipping_address": {
    "attention": "Warehouse Manager",
    "street": "123 Industrial Estate",
    "address_line2": "Block A, Warehouse 1",
    "city": "Bangalore",
    "state": "Karnataka",
    "country": "India",
    "postal_code": "560001",
    "phone": "08041234567"
  },
  "contact_persons": [
    {
      "salutation": "Mr.",
      "first_name": "Suresh",
      "last_name": "Singh",
      "email_address": "suresh.singh@aquaplast.com",
      "mobile": "9876543211"
    }
  ],
  "bank_details": [
    {
      "bank_id": 1,
      "account_holder_name": "AquaPlast Industries Pvt Ltd",
      "account_number": "1234567890123456",
      "reenter_account_number": "1234567890123456"
    }
  ]
}
```

### 1.3 DTO Output Response

**Status:** 201 Created

```json
{
  "success": true,
  "data": {
    "id": 5,
    "first_name": "Rajesh",
    "last_name": "Kumar",
    "display_name": "AquaPlast Industries",
    "company_name": "AquaPlast Industries Pvt Ltd",
    "email_address": "rajesh.kumar@aquaplast.com",
    "work_phone": "08041234567",
    "mobile": "9876543210",
    "other_details": {
      "pan": "AABCT5678H",
      "is_msme_registered": true,
      "currency": "INR",
      "payment_terms": "Net 45",
      "tds": "2%",
      "enable_portal": true,
      "website_url": "https://www.aquaplast.com"
    },
    "billing_address": {
      "id": 1,
      "address_line1": "123 Industrial Estate",
      "city": "Bangalore",
      "state": "Karnataka",
      "country_region": "India",
      "pin_code": "560001"
    },
    "contact_persons": [
      {
        "id": 1,
        "first_name": "Suresh",
        "email_address": "suresh.singh@aquaplast.com"
      }
    ],
    "bank_details": [
      {
        "id": 1,
        "bank_id": 1,
        "account_holder_name": "AquaPlast Industries Pvt Ltd",
        "account_number": "****7890"
      }
    ],
    "created_at": "2026-02-17T09:00:00Z"
  },
  "message": "Vendor created successfully"
}
```

### 1.4 Additional Vendors to Create

You may need multiple vendors for different items:

| Vendor | Products | Email |
|--------|----------|-------|
| AquaPlast Industries | PET Bottles | rajesh@aquaplast.com |
| CapMaster Solutions | Plastic Caps | info@capmaster.com |
| LabelTech Pvt Ltd | Custom Labels | sales@labeltech.com |

---

## Step 2: Create Customer (Buyers)

### Zoho Menu Path: Sales → Customers

Customers are buyers of water bottles and related products.

### 2.1 Customer Details

**Customer Name:** Fresh Water Retail Pvt Ltd  
**Business Type:** Retail Store Chain  
**Location:** Bangalore  
**Contact:** orders@freshwaterretail.com

### 2.2 DTO Input Request (Create Customer)

**Endpoint:** `POST /customers`  
**Authentication:** Bearer Token + SuperAdmin Role

```json
{
  "salutation": "Mr.",
  "first_name": "Amit",
  "last_name": "Singh",
  "company_name": "Fresh Water Retail Pvt Ltd",
  "display_name": "Fresh Water Retail",
  "email_address": "amit.singh@freshwaterretail.com",
  "work_phone": "08041234500",
  "work_phone_code": "+91",
  "mobile": "9876543200",
  "mobile_code": "+91",
  "customer_language": "English",
  "other_details": {
    "pan": "BBCDE1234H",
    "currency": "INR",
    "payment_terms": "Net 30",
    "enable_portal": true
  },
  "billing_address": {
    "attention": "Finance Department",
    "street": "456 Market Street",
    "address_line2": "Building B",
    "city": "Bangalore",
    "state": "Karnataka",
    "country": "India",
    "postal_code": "560002",
    "phone": "08041234500",
    "phone_code": "+91"
  },
  "shipping_address": {
    "attention": "Store Manager",
    "street": "789 Retail Plaza",
    "address_line2": "Store 1",
    "city": "Bangalore",
    "state": "Karnataka",
    "country": "India",
    "postal_code": "560003",
    "phone": "08041234501"
  },
  "contact_persons": [
    {
      "salutation": "Mr.",
      "first_name": "Pradeep",
      "last_name": "Sharma",
      "email_address": "pradeep.sharma@freshwaterretail.com",
      "mobile": "9876543201"
    }
  ]
}
```

### 2.3 DTO Output Response

**Status:** 201 Created

```json
{
  "success": true,
  "data": {
    "id": 12,
    "first_name": "Amit",
    "last_name": "Singh",
    "display_name": "Fresh Water Retail",
    "company_name": "Fresh Water Retail Pvt Ltd",
    "email_address": "amit.singh@freshwaterretail.com",
    "work_phone": "08041234500",
    "mobile": "9876543200",
    "other_details": {
      "pan": "BBCDE1234H",
      "currency": "INR",
      "payment_terms": "Net 30",
      "enable_portal": true
    },
    "billing_address": {
      "id": 1,
      "address_line1": "456 Market Street",
      "city": "Bangalore",
      "country_region": "India",
      "pin_code": "560002"
    },
    "shipping_address": {
      "id": 2,
      "address_line1": "789 Retail Plaza",
      "city": "Bangalore",
      "country_region": "India",
      "pin_code": "560003"
    },
    "contact_persons": [
      {
        "id": 1,
        "first_name": "Pradeep",
        "email_address": "pradeep.sharma@freshwaterretail.com"
      }
    ],
    "created_at": "2026-02-17T09:30:00Z"
  },
  "message": "Customer created successfully"
}
```

### 2.4 Additional Customers to Create

| Customer | Type | Email | Location |
|----------|------|-------|----------|
| Fresh Water Retail | Retail Chain | orders@freshwaterretail.com | Bangalore |
| AquaOffice Solutions | Corporate Buyer | procurement@aquaoffice.com | Mumbai |
| Wellness Stores India | Retail Distributor | supply@wellnessstores.com | Delhi |

---

## Product Overview

### Water Company Product Hierarchy

```
WATER PRODUCTS
├── Bottles
│   ├── 500ml PET Bottles (with variants: Regular Cap, Sports Cap)
│   ├── 1L PET Bottles
│   ├── 2L PET Bottles
│   └── 20L Polycarbonate Cooler Bottles
├── Caps & Closures
│   ├── 20mm Tamper-proof Caps (for 500ml)
│   ├── 28mm Caps (for 1L and 2L)
│   └── 90mm Large Caps (for 20L bottles)
├── Labels & Packaging
│   ├── Custom Water Labels
│   ├── Packaging Boxes
│   └── Shipping Cartons
└── Accessories
    ├── Water Cooler Stands
    └── Bottle Crates
```

---

## Step 3: Create 500ml Water Bottle

### 3.1 Product Details
- **Product Name:** 500ml Premium Drinking Water Bottle
- **Type:** goods (physical product)
- **Structure:** variants (same product with different options)
- **Unit:** piece
- **Variants:** Regular Cap, Sports Cap
- **Tracking:** Yes (track inventory)

### 3.2 DTO Input Request

**Endpoint:** `POST /items`  
**Authentication:** Bearer Token + Admin Role

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
    "attributes": [
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
        "selling_price": 15.00,
        "cost_price": 8.00,
        "stock_quantity": 5000
      },
      {
        "sku": "WTR-BOT-500-SPORT",
        "attribute_map": {
          "Cap Type": "Sports Cap"
        },
        "selling_price": 18.00,
        "cost_price": 9.50,
        "stock_quantity": 3000
      }
    ]
  },
  "sales_info": {
    "account": "Sales Revenue - Water Bottles",
    "selling_price": 16.50,
    "currency": "INR",
    "description": "500ml drinking water bottles retail sales"
  },
  "purchase_info": {
    "account": "Cost of Goods Purchased",
    "cost_price": 8.75,
    "currency": "INR",
    "preferred_vendor_id": 5,
    "description": "Purchase from bottle supplier"
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

### 3.3 DTO Output Response

**Status:** 201 Created

```json
{
  "success": true,
  "data": {
    "id": "ITEM-WTR-500-001",
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
          "options": ["Regular Cap", "Sports Cap"]
        }
      ],
      "variants": [
        {
          "sku": "WTR-BOT-500-REG",
          "attribute_map": {
            "Cap Type": "Regular Cap"
          },
          "selling_price": 15.00,
          "cost_price": 8.00,
          "stock_quantity": 5000
        },
        {
          "sku": "WTR-BOT-500-SPORT",
          "attribute_map": {
            "Cap Type": "Sports Cap"
          },
          "selling_price": 18.00,
          "cost_price": 9.50,
          "stock_quantity": 3000
        }
      ]
    },
    "sales_info": {
      "account": "Sales Revenue - Water Bottles",
      "selling_price": 16.50,
      "currency": "INR",
      "description": "500ml drinking water bottles retail sales"
    },
    "purchase_info": {
      "account": "Cost of Goods Purchased",
      "cost_price": 8.75,
      "currency": "INR",
      "preferred_vendor_id": 5,
      "description": "Purchase from bottle supplier"
    },
    "inventory": {
      "track_inventory": true,
      "inventory_account": "Inventory - Water Bottles",
      "inventory_valuation_method": "FIFO",
      "reorder_point": 1000
    },
    "return_policy": {
      "returnable": true
    },
    "created_at": "2026-02-17T10:00:00Z",
    "updated_at": "2026-02-17T10:00:00Z"
  },
  "message": "Item created successfully"
}
```

---

## Step 4: Create 20L Water Cooler Bottle

### 4.1 Product Details
- **Product Name:** 20 Litre Polycarbonate Water Cooler Bottle
- **Type:** goods
- **Structure:** single (no variants needed for bulk bottles)
- **Unit:** piece
- **Tracking:** Yes (track inventory)
- **Premium pricing:** Water cooler bottles

### 2.2 DTO Input Request

**Endpoint:** `POST /items`

```json
{
  "name": "20 Litre Polycarbonate Water Cooler Bottle",
  "type": "goods",
  "item_details": {
    "structure": "single",
    "unit": "piece",
    "sku": "WTR-COOLER-20L",
    "upc": "8904220200001",
    "ean": "8904220200001",
    "description": "20 litre premium polycarbonate water cooler bottle. High capacity, durable, transparent, suitable for office and institutional use. Deposit/return model."
  },
  "sales_info": {
    "account": "Sales Revenue - Water Coolers",
    "selling_price": 150.00,
    "currency": "INR",
    "description": "20L cooler bottles for offices and institutions"
  },
  "purchase_info": {
    "account": "Cost of Goods Purchased",
    "cost_price": 75.00,
    "currency": "INR",
    "preferred_vendor_id": 5,
    "description": "Bulk water bottle from supplier"
  },
  "inventory": {
    "track_inventory": true,
    "inventory_account": "Inventory - Cooler Bottles",
    "inventory_valuation_method": "FIFO",
    "reorder_point": 50
  },
  "return_policy": {
    "returnable": true
  }
}
```

### 2.3 DTO Output Response

**Status:** 201 Created

```json
{
  "success": true,
  "data": {
    "id": "ITEM-WTR-20L-002",
    "name": "20 Litre Polycarbonate Water Cooler Bottle",
    "type": "goods",
    "item_details": {
      "structure": "single",
      "unit": "piece",
      "sku": "WTR-COOLER-20L",
      "upc": "8904220200001",
      "ean": "8904220200001",
      "description": "20 litre premium polycarbonate water cooler bottle. High capacity, durable, transparent, suitable for office and institutional use. Deposit/return model."
    },
    "sales_info": {
      "account": "Sales Revenue - Water Coolers",
      "selling_price": 150.00,
      "currency": "INR",
      "description": "20L cooler bottles for offices and institutions"
    },
    "purchase_info": {
      "account": "Cost of Goods Purchased",
      "cost_price": 75.00,
      "currency": "INR",
      "preferred_vendor_id": 5,
      "description": "Bulk water bottle from supplier"
    },
    "inventory": {
      "track_inventory": true,
      "inventory_account": "Inventory - Cooler Bottles",
      "inventory_valuation_method": "FIFO",
      "reorder_point": 50
    },
    "return_policy": {
      "returnable": true
    },
    "created_at": "2026-02-17T10:15:00Z",
    "updated_at": "2026-02-17T10:15:00Z"
  },
  "message": "Item created successfully"
}
```

---

## Step 5: Create Plastic Caps (by Size)

### 5.1 Caps Overview

Water bottle caps are sold separately by size:

| Cap Size | Product | Fits Bottle | Quantity/Pack |
|----------|---------|-------------|--------------|
| 20mm | Regular/Sports Caps | 500ml | 100 pieces/pack |
| 28mm | Standard Cap | 1L, 2L | 50 pieces/pack |
| 90mm | Large Cap | 20L Cooler | 20 pieces/pack |

### 5.2 Create 20mm Tamper-proof Caps

#### DTO Input Request

**Endpoint:** `POST /items`

```json
{
  "name": "20mm Tamper-proof Bottle Cap - Pack of 100",
  "type": "goods",
  "item_details": {
    "structure": "variants",
    "unit": "pack",
    "sku": "CAP-20MM-BASE",
    "description": "20mm tamper-proof food-grade plastic bottle caps. Perfect for 500ml water bottles. Sold per pack of 100 pieces. High quality, reusable.",
    "attributes": [
      {
        "key": "Color",
        "options": ["White", "Blue"]
      }
    ],
    "variants": [
      {
        "sku": "CAP-20MM-WHITE",
        "attribute_map": {
          "Color": "White"
        },
        "selling_price": 45.00,
        "cost_price": 20.00,
        "stock_quantity": 500
      },
      {
        "sku": "CAP-20MM-BLUE",
        "attribute_map": {
          "Color": "Blue"
        },
        "selling_price": 50.00,
        "cost_price": 22.00,
        "stock_quantity": 300
      }
    ]
  },
  "sales_info": {
    "account": "Sales Revenue - Packaging",
    "selling_price": 47.50,
    "currency": "INR",
    "description": "20mm tamper caps for 500ml bottles"
  },
  "purchase_info": {
    "account": "Cost of Caps",
    "cost_price": 21.00,
    "currency": "INR",
    "preferred_vendor_id": 6,
    "description": "Tamper-proof caps from supplier"
  },
  "inventory": {
    "track_inventory": true,
    "inventory_account": "Inventory - Packaging",
    "inventory_valuation_method": "FIFO",
    "reorder_point": 100
  },
  "return_policy": {
    "returnable": false
  }
}
```

#### DTO Output Response

**Status:** 201 Created

```json
{
  "success": true,
  "data": {
    "id": "ITEM-CAP-20MM-003",
    "name": "20mm Tamper-proof Bottle Cap - Pack of 100",
    "type": "goods",
    "item_details": {
      "structure": "variants",
      "unit": "pack",
      "sku": "CAP-20MM-BASE",
      "description": "20mm tamper-proof food-grade plastic bottle caps. Perfect for 500ml water bottles. Sold per pack of 100 pieces. High quality, reusable.",
      "attribute_definitions": [
        {
          "key": "Color",
          "options": ["White", "Blue"]
        }
      ],
      "variants": [
        {
          "sku": "CAP-20MM-WHITE",
          "attribute_map": {
            "Color": "White"
          },
          "selling_price": 45.00,
          "cost_price": 20.00,
          "stock_quantity": 500
        },
        {
          "sku": "CAP-20MM-BLUE",
          "attribute_map": {
            "Color": "Blue"
          },
          "selling_price": 50.00,
          "cost_price": 22.00,
          "stock_quantity": 300
        }
      ]
    },
    "sales_info": {
      "account": "Sales Revenue - Packaging",
      "selling_price": 47.50,
      "currency": "INR",
      "description": "20mm tamper caps for 500ml bottles"
    },
    "purchase_info": {
      "account": "Cost of Caps",
      "cost_price": 21.00,
      "currency": "INR",
      "preferred_vendor_id": 6,
      "description": "Tamper-proof caps from supplier"
    },
    "inventory": {
      "track_inventory": true,
      "inventory_account": "Inventory - Packaging",
      "inventory_valuation_method": "FIFO",
      "reorder_point": 100
    },
    "return_policy": {
      "returnable": false
    },
    "created_at": "2026-02-17T10:30:00Z",
    "updated_at": "2026-02-17T10:30:00Z"
  },
  "message": "Item created successfully"
}
```

### 3.3 Create 28mm Standard Caps (for 1L/2L)

#### DTO Input Request

**Endpoint:** `POST /items`

```json
{
  "name": "28mm Standard Bottle Cap - Pack of 50",
  "type": "goods",
  "item_details": {
    "structure": "single",
    "unit": "pack",
    "sku": "CAP-28MM-STD",
    "description": "28mm standard food-grade plastic bottle caps. Suitable for 1L and 2L water bottles. Pack of 50 pieces."
  },
  "sales_info": {
    "account": "Sales Revenue - Packaging",
    "selling_price": 60.00,
    "currency": "INR",
    "description": "28mm caps for 1L and 2L bottles"
  },
  "purchase_info": {
    "account": "Cost of Caps",
    "cost_price": 28.00,
    "currency": "INR",
    "preferred_vendor_id": 6
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

#### DTO Output Response

**Status:** 201 Created

```json
{
  "success": true,
  "data": {
    "id": "ITEM-CAP-28MM-004",
    "name": "28mm Standard Bottle Cap - Pack of 50",
    "type": "goods",
    "item_details": {
      "structure": "single",
      "unit": "pack",
      "sku": "CAP-28MM-STD",
      "description": "28mm standard food-grade plastic bottle caps. Suitable for 1L and 2L water bottles. Pack of 50 pieces."
    },
    "sales_info": {
      "account": "Sales Revenue - Packaging",
      "selling_price": 60.00,
      "currency": "INR",
      "description": "28mm caps for 1L and 2L bottles"
    },
    "purchase_info": {
      "account": "Cost of Caps",
      "cost_price": 28.00,
      "currency": "INR",
      "preferred_vendor_id": 6
    },
    "inventory": {
      "track_inventory": true,
      "inventory_account": "Inventory - Packaging",
      "reorder_point": 50
    },
    "return_policy": {
      "returnable": false
    },
    "created_at": "2026-02-17T10:45:00Z",
    "updated_at": "2026-02-17T10:45:00Z"
  },
  "message": "Item created successfully"
}
```

### 3.4 Create 90mm Large Caps (for 20L)

#### DTO Input Request

```json
{
  "name": "90mm Large Water Cooler Cap",
  "type": "goods",
  "item_details": {
    "structure": "single",
    "unit": "piece",
    "sku": "CAP-90MM-LARGE",
    "description": "90mm large polycarbonate cap for 20L water cooler bottles. Universal fit. Durable and long-lasting."
  },
  "sales_info": {
    "account": "Sales Revenue - Packaging",
    "selling_price": 35.00,
    "currency": "INR",
    "description": "Large caps for 20L cooler bottles"
  },
  "purchase_info": {
    "account": "Cost of Caps",
    "cost_price": 15.00,
    "currency": "INR",
    "preferred_vendor_id": 6
  },
  "inventory": {
    "track_inventory": true,
    "inventory_account": "Inventory - Packaging",
    "reorder_point": 100
  },
  "return_policy": {
    "returnable": true
  }
}
```

#### DTO Output Response

```json
{
  "success": true,
  "data": {
    "id": "ITEM-CAP-90MM-005",
    "name": "90mm Large Water Cooler Cap",
    "type": "goods",
    "item_details": {
      "structure": "single",
      "unit": "piece",
      "sku": "CAP-90MM-LARGE",
      "description": "90mm large polycarbonate cap for 20L water cooler bottles. Universal fit. Durable and long-lasting."
    },
    "sales_info": {
      "account": "Sales Revenue - Packaging",
      "selling_price": 35.00,
      "currency": "INR",
      "description": "Large caps for 20L cooler bottles"
    },
    "purchase_info": {
      "account": "Cost of Caps",
      "cost_price": 15.00,
      "currency": "INR",
      "preferred_vendor_id": 6
    },
    "inventory": {
      "track_inventory": true,
      "inventory_account": "Inventory - Packaging",
      "reorder_point": 100
    },
    "return_policy": {
      "returnable": true
    },
    "created_at": "2026-02-17T11:00:00Z",
    "updated_at": "2026-02-17T11:00:00Z"
  },
  "message": "Item created successfully"
}
```

---

## Step 4: Create Water Labels

### 4.1 Product Details

- **Product Name:** Custom Water Bottle Labels - 1000 pieces
- **Type:** goods (packaging material)
- **Structure:** single
- **Unit:** pack
- **SKU:** LBL-1000-CUSTOM

### 4.2 DTO Input Request

**Endpoint:** `POST /items`

```json
{
  "name": "Custom Water Bottle Labels - 1000 pieces",
  "type": "goods",
  "item_details": {
    "structure": "single",
    "unit": "pack",
    "sku": "LBL-1000-CUSTOM",
    "description": "Custom printed waterproof water bottle labels. 1000 pieces per pack. High-quality vinyl, resistant to moisture and UV. Printed with company logo, batch number, date."
  },
  "sales_info": {
    "account": "Sales Revenue - Packaging",
    "selling_price": 500.00,
    "currency": "INR",
    "description": "Custom water labels for branding"
  },
  "purchase_info": {
    "account": "Cost of Labels",
    "cost_price": 300.00,
    "currency": "INR",
    "preferred_vendor_id": 7,
    "description": "Custom label printing from supplier"
  },
  "inventory": {
    "track_inventory": true,
    "inventory_account": "Inventory - Packaging Labels",
    "inventory_valuation_method": "FIFO",
    "reorder_point": 50
  },
  "return_policy": {
    "returnable": false
  }
}
```

### 4.3 DTO Output Response

**Status:** 201 Created

```json
{
  "success": true,
  "data": {
    "id": "ITEM-LBL-1K-006",
    "name": "Custom Water Bottle Labels - 1000 pieces",
    "type": "goods",
    "item_details": {
      "structure": "single",
      "unit": "pack",
      "sku": "LBL-1000-CUSTOM",
      "description": "Custom printed waterproof water bottle labels. 1000 pieces per pack. High-quality vinyl, resistant to moisture and UV. Printed with company logo, batch number, date."
    },
    "sales_info": {
      "account": "Sales Revenue - Packaging",
      "selling_price": 500.00,
      "currency": "INR",
      "description": "Custom water labels for branding"
    },
    "purchase_info": {
      "account": "Cost of Labels",
      "cost_price": 300.00,
      "currency": "INR",
      "preferred_vendor_id": 7,
      "description": "Custom label printing from supplier"
    },
    "inventory": {
      "track_inventory": true,
      "inventory_account": "Inventory - Packaging Labels",
      "inventory_valuation_method": "FIFO",
      "reorder_point": 50
    },
    "return_policy": {
      "returnable": false
    },
    "created_at": "2026-02-17T11:15:00Z",
    "updated_at": "2026-02-17T11:15:00Z"
  },
  "message": "Item created successfully"
}
```

---

## Step 5: Set Opening Stock

### 5.1 Opening Stock for 500ml Bottles

**Endpoint:** `PUT /items/ITEM-WTR-500-001/opening-stock`

#### DTO Input Request

```json
{
  "opening_stock": 10000,
  "opening_stock_rate_per_unit": 8.00
}
```

#### DTO Output Response

**Status:** 200 OK

```json
{
  "success": true,
  "data": {
    "id": "ITEM-WTR-500-001",
    "name": "500ml Premium Drinking Water Bottle",
    "opening_stock": 10000,
    "opening_stock_value": 80000,
    "opening_stock_rate_per_unit": 8.00
  },
  "message": "Opening stock updated successfully"
}
```

### 5.2 Opening Stock for Variants (Per Variant)

**Endpoint:** `PUT /items/ITEM-WTR-500-001/variants/opening-stock`

#### DTO Input Request

```json
{
  "variants": [
    {
      "variant_id": 1,
      "opening_stock": 6000,
      "opening_stock_rate_per_unit": 8.00
    },
    {
      "variant_id": 2,
      "opening_stock": 4000,
      "opening_stock_rate_per_unit": 9.50
    }
  ]
}
```

#### DTO Output Response

**Status:** 200 OK

```json
{
  "success": true,
  "data": {
    "id": "ITEM-WTR-500-001",
    "name": "500ml Premium Drinking Water Bottle",
    "variants": [
      {
        "sku": "WTR-BOT-500-REG",
        "attribute_map": {
          "Cap Type": "Regular Cap"
        },
        "opening_stock": 6000,
        "opening_stock_value": 48000,
        "opening_stock_rate_per_unit": 8.00
      },
      {
        "sku": "WTR-BOT-500-SPORT",
        "attribute_map": {
          "Cap Type": "Sports Cap"
        },
        "opening_stock": 4000,
        "opening_stock_value": 38000,
        "opening_stock_rate_per_unit": 9.50
      }
    ],
    "total_value": 86000
  },
  "message": "Variant opening stock updated successfully"
}
```

---

## Step 8: Create Item Groups (BOM)

Item Groups (Bill of Materials) combine multiple items into a finished product ready for sale. This allows you to define how a complete water bottle is assembled from individual components (bottle + cap + label).

**Authentication:** Bearer Token + Admin Role

### 8.1 500ml Complete Water Bottle (Bottle + Cap + Label)

**Endpoint:** `POST /item-groups`

#### DTO Input Request

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

#### DTO Output Response

**Status:** 201 Created

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
        },
        "created_at": "2026-02-17T11:30:00Z",
        "updated_at": "2026-02-17T11:30:00Z"
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
        },
        "created_at": "2026-02-17T11:30:00Z",
        "updated_at": "2026-02-17T11:30:00Z"
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
        },
        "created_at": "2026-02-17T11:30:00Z",
        "updated_at": "2026-02-17T11:30:00Z"
      }
    ],
    "created_at": "2026-02-17T11:30:00Z",
    "updated_at": "2026-02-17T11:30:00Z"
  },
  "message": "Item Group created successfully"
}
```

### 8.2 20L Water Cooler Complete Set

**Endpoint:** `PUT /item-groups/{id}` (Update) or `POST /item-groups` (Create)

#### DTO Input Request

```json
{
  "name": "20L Water Cooler Complete Set",
  "description": "20L polycarbonate bottle with large cap and labels. Complete for water cooler use.",
  "is_active": true,
  "components": [
    {
      "item_id": "item_e4c9a5c8",
      "quantity": 1,
      "variant_details": {
        "type": "bottle",
        "capacity": "20L"
      }
    },
    {
      "item_id": "item_3382c1ab",
      "quantity": 1,
      "variant_details": {
        "type": "cap",
        "size": "90mm"
      }
    },
    {
      "item_id": "item_8195d60a",
      "quantity": 0.005,
      "variant_details": {
        "type": "large_label"
      }
    }
  ]
}
```

#### DTO Output Response

```json
{
  "success": true,
  "data": {
    "id": "ig_x9y8z7w6",
    "name": "20L Water Cooler Complete Set",
    "description": "20L polycarbonate bottle with large cap and labels. Complete for water cooler use.",
    "is_active": true,
    "components": [
      {
        "id": 1,
        "item_group_id": "ig_x9y8z7w6",
        "item_id": "item_e4c9a5c8",
        "item": {
          "id": "item_e4c9a5c8",
          "name": "20 Litre Polycarbonate Water Cooler Bottle",
          "sku": "WTR-COOLER-20L"
        },
        "quantity": 1,
        "variant_details": {
          "type": "bottle",
          "capacity": "20L"
        },
        "created_at": "2026-02-17T11:45:00Z",
        "updated_at": "2026-02-17T11:45:00Z"
      },
      {
        "id": 2,
        "item_group_id": "ig_x9y8z7w6",
        "item_id": "item_3382c1ab",
        "item": {
          "id": "item_3382c1ab",
          "name": "90mm Large Water Cooler Cap",
          "sku": "CAP-90MM-LARGE"
        },
        "quantity": 1,
        "variant_details": {
          "type": "cap",
          "size": "90mm"
        },
        "created_at": "2026-02-17T11:45:00Z",
        "updated_at": "2026-02-17T11:45:00Z"
      },
      {
        "id": 3,
        "item_group_id": "ig_x9y8z7w6",
        "item_id": "item_8195d60a",
        "item": {
          "id": "item_8195d60a",
          "name": "Custom Water Bottle Labels - 1000 pieces",
          "sku": "LBL-1000-CUSTOM"
        },
        "quantity": 0.005,
        "variant_details": {
          "type": "large_label"
        },
        "created_at": "2026-02-17T11:45:00Z",
        "updated_at": "2026-02-17T11:45:00Z"
      }
    ],
    "created_at": "2026-02-17T11:45:00Z",
    "updated_at": "2026-02-17T11:45:00Z"
  },
  "message": "Item Group created successfully"
}
```

---

## Step 9: Additional ItemGroup Operations

### 9.1 Get All Item Groups with Search & Pagination

**Endpoint:** `GET /item-groups?limit=10&offset=0&search=500ml`

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": "ig_a1b2c3d4",
      "name": "500ml Complete Water Bottle Packaged",
      "description": "Complete packaged 500ml water bottle...",
      "is_active": true,
      "components": [...],
      "created_at": "2026-02-17T11:30:00Z",
      "updated_at": "2026-02-17T11:30:00Z"
    }
  ],
  "total": 2,
  "page": 1,
  "page_size": 10
}
```

### 9.2 Get Item Group by ID

**Endpoint:** `GET /item-groups/{id}`

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "ig_a1b2c3d4",
    "name": "500ml Complete Water Bottle Packaged",
    "description": "Complete packaged 500ml water bottle...",
    "is_active": true,
    "components": [...],
    "created_at": "2026-02-17T11:30:00Z",
    "updated_at": "2026-02-17T11:30:00Z"
  }
}
```

### 9.3 Update Item Group

**Endpoint:** `PUT /item-groups/{id}`

**Request:**
```json
{
  "name": "500ml Complete Water Bottle Updated",
  "description": "Updated description",
  "is_active": true,
  "components": [
    {
      "item_id": "ITEM-WTR-500-001",
      "variant_id": 1,
      "quantity": 1
    }
  ]
}
```

### 9.4 Search Item Group by Name

**Endpoint:** `GET /item-groups/search/by-name?name=500ml`

### 9.5 Delete Item Group

**Endpoint:** `DELETE /item-groups/{id}`

**Response:**
```json
{
  "success": true,
  "message": "Item Group deleted successfully"
}
```

---

## Step 10: Create Purchase Order

### 10.1 Purchase Order for Bulk Water Bottles

**Endpoint:** `POST /purchase-orders`  
**Authentication:** Bearer Token + Admin Role

#### 10.2 DTO Input Request (With Multiple Line Items)

```json
{
  "vendor_id": 5,
  "delivery_address_type": "organization",
  "organization_name": "AquaPlast Industries",
  "organization_address": "Warehouse #1, 123 Industrial Park, Bangalore",
  "reference_no": "BOTTLE-BULK-FEB",
  "date": "2026-02-17",
  "delivery_date": "2026-03-03",
  "payment_terms": "net_45",
  "shipment_preference": "Truck Transport",
  "discount": 0,
  "discount_type": "amount",
  "tax_type": "SGST",
  "tax_id": 1,
  "notes": "Bulk order for retail stock. Regular Cap and Sports Cap variants.",
  "line_items": [
    {
      "item_id": "item_1a2b3c4d",
      "variant_id": 1,
      "account": "Cost of Goods Purchased",
      "quantity": 5000,
      "rate": 8.00,
      "variant_details": {
        "Cap Type": "Regular Cap"
      }
    },
    {
      "item_id": "item_1a2b3c4d",
      "variant_id": 2,
      "account": "Cost of Goods Purchased",
      "quantity": 3000,
      "rate": 9.50,
      "variant_details": {
        "Cap Type": "Sports Cap"
      }
    }
  ]
}
```

#### 9.3 DTO Output Response

**Status:** 201 Created

```json
{
  "success": true,
  "data": {
    "id": "12345678-1234-1234-1234-123456789012",
    "po_number": "PO-20260217-0001",
    "reference_no": "BOTTLE-BULK-FEB",
    "vendor_id": 5,
    "vendor_name": "AquaPlast Industries",
    "delivery_address_type": "organization",
    "organization_name": "AquaPlast Industries",
    "organization_address": "Warehouse #1, 123 Industrial Park, Bangalore",
    "po_date": "2026-02-17",
    "delivery_date": "2026-03-03",
    "payment_terms": "net_45",
    "shipment_preference": "Truck Transport",
    "status": "draft",
    "notes": "Bulk order for retail stock. Regular Cap and Sports Cap variants.",
    "line_items": [
      {
        "id": 1,
        "item_id": "item_1a2b3c4d",
        "item_name": "500ml Premium Drinking Water Bottle",
        "variant_id": 1,
        "variant_name": "Regular Cap",
        "quantity": 5000,
        "rate": 8.00,
        "amount": 40000,
        "account": "Cost of Goods Purchased",
        "variant_details": {
          "Cap Type": "Regular Cap"
        },
        "created_at": "2026-02-17T12:00:00Z"
      },
      {
        "id": 2,
        "item_id": "item_1a2b3c4d",
        "item_name": "500ml Premium Drinking Water Bottle",
        "variant_id": 2,
        "variant_name": "Sports Cap",
        "quantity": 3000,
        "rate": 9.50,
        "amount": 28500,
        "account": "Cost of Goods Purchased",
        "variant_details": {
          "Cap Type": "Sports Cap"
        },
        "created_at": "2026-02-17T12:01:00Z"
      }
    ],
    "sub_total": 68500,
    "discount": 0,
    "tax_amount": 0,
    "adjustment": 0,
    "total": 68500,
    "created_at": "2026-02-17T12:00:00Z",
    "updated_at": "2026-02-17T12:00:00Z"
  },
  "message": "Purchase order created successfully with 2 line items"
}
```

### 9.4 Add Line Items to Purchase Order

**Endpoint:** `POST /purchase-orders/PO-2026-WTR-001/line-items`

#### DTO Input Request

```json
{
  "item_id": "ITEM-WTR-500-001",
  "variant_id": 1,
  "quantity": 5000,
  "rate": 8.00,
  "description": "500ml Regular Cap bottles (5000 units)",
  "tax_id": 1,
  "warehouse_id": 1
}
```

#### DTO Output Response

**Status:** 201 Created

```json
{
  "success": true,
  "data": {
    "id": 1,
    "purchase_order_id": "PO-2026-WTR-001",
    "line_item_number": 1,
    "item_id": "ITEM-WTR-500-001",
    "item_name": "500ml Premium Drinking Water Bottle",
    "variant_id": 1,
    "variant_name": "Regular Cap",
    "quantity": 5000,
    "unit": "piece",
    "rate": 8.00,
    "amount": 40000,
    "tax_id": 1,
    "tax_percentage": 5,
    "tax_amount": 2000,
    "total": 42000,
    "warehouse_id": 1,
    "warehouse_name": "Warehouse #1",
    "description": "500ml Regular Cap bottles (5000 units)",
    "created_at": "2026-02-17T12:00:00Z"
  },
  "message": "Line item added successfully"
}
```

### 9.5 Add Second Line Item (Sports Cap Variant)

**Endpoint:** `POST /purchase-orders/PO-2026-WTR-001/line-items`

#### DTO Input Request

```json
{
  "item_id": "ITEM-WTR-500-001",
  "variant_id": 2,
  "quantity": 3000,
  "rate": 9.50,
  "description": "500ml Sports Cap bottles (3000 units)",
  "tax_id": 1,
  "warehouse_id": 1
}
```

#### DTO Output Response

**Status:** 201 Created

```json
{
  "success": true,
  "data": {
    "id": 2,
    "purchase_order_id": "PO-2026-WTR-001",
    "line_item_number": 2,
    "item_id": "ITEM-WTR-500-001",
    "item_name": "500ml Premium Drinking Water Bottle",
    "variant_id": 2,
    "variant_name": "Sports Cap",
    "quantity": 3000,
    "unit": "piece",
    "rate": 9.50,
    "amount": 28500,
    "tax_id": 1,
    "tax_percentage": 5,
    "tax_amount": 1425,
    "total": 29925,
    "warehouse_id": 1,
    "warehouse_name": "Warehouse #1",
    "description": "500ml Sports Cap bottles (3000 units)",
    "created_at": "2026-02-17T12:05:00Z"
  },
  "message": "Line item added successfully"
}
```

### 9.6 Complete Purchase Order

**Endpoint:** `PUT /purchase-orders/PO-2026-WTR-001/confirm`

#### DTO Input Request

```json
{
  "action": "confirm"
}
```

#### DTO Output Response

**Status:** 200 OK

```json
{
  "success": true,
  "data": {
    "id": 1,
    "po_number": "PO-2026-WTR-001",
    "vendor_id": 5,
    "vendor_name": "AquaPlast Industries",
    "po_date": "2026-02-17",
    "expected_delivery_date": "2026-03-03",
    "status": "Confirmed",
    "line_items": [
      {
        "id": 1,
        "item_id": "ITEM-WTR-500-001",
        "quantity": 5000,
        "rate": 8.00,
        "total": 42000
      },
      {
        "id": 2,
        "item_id": "ITEM-WTR-500-001",
        "quantity": 3000,
        "rate": 9.50,
        "total": 29925
      }
    ],
    "total_amount": 68500,
    "tax_total": 3425,
    "grand_total": 71925,
    "created_at": "2026-02-17T12:00:00Z",
    "confirmed_at": "2026-02-17T12:10:00Z"
  },
  "message": "Purchase order confirmed successfully"
}
```

### 9.7 Inventory Tracking for Purchase Order

When a Purchase Order is **confirmed**, inventory is automatically tracked and updated in the system.

#### Inventory Updates on PO Confirmation

**Endpoint:** `GET /inventory-balance?item_id=item_1a2b3c4d`

```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "item_id": "item_1a2b3c4d",
      "item_name": "500ml Premium Drinking Water Bottle",
      "variant_id": 1,
      "variant_name": "Regular Cap",
      "current_quantity": 5000,
      "reserved_quantity": 0,
      "available_quantity": 5000,
      "in_transit_quantity": 0,
      "average_rate": 8.00,
      "last_received_date": "2026-02-17",
      "last_sold_date": null,
      "updated_at": "2026-02-17T12:00:00Z"
    },
    {
      "id": 2,
      "item_id": "item_1a2b3c4d",
      "item_name": "500ml Premium Drinking Water Bottle",
      "variant_id": 2,
      "variant_name": "Sports Cap",
      "current_quantity": 3000,
      "reserved_quantity": 0,
      "available_quantity": 3000,
      "in_transit_quantity": 0,
      "average_rate": 9.50,
      "last_received_date": "2026-02-17",
      "last_sold_date": null,
      "updated_at": "2026-02-17T12:00:00Z"
    }
  ],
  "message": "Inventory balance retrieved successfully"
}
```

#### Inventory Journal Entry

When PO is received, this entry is created:

```json
{
  "id": 1,
  "item_id": "item_1a2b3c4d",
  "variant_id": 1,
  "transaction_type": "PURCHASE_ORDER_RECEIVED",
  "quantity": 5000,
  "reference_type": "PurchaseOrder",
  "reference_id": "12345678-1234-1234-1234-123456789012",
  "reference_no": "PO-20260217-0001",
  "notes": "500ml Regular Cap - Received from AquaPlast Industries",
  "created_at": "2026-02-17T12:00:00Z",
  "created_by": "admin_user"
}
```

---

## Step 10: Create Bill (Vendor Invoice)

### 10.1 Bill Details

When goods are received from the vendor, a bill (vendor invoice) is created to record the purchase.

**Endpoint:** `POST /bills`  
**Authentication:** Bearer Token + Admin Role

### 10.2 DTO Input Request

```json
{
  "vendor_id": 5,
  "company_id": 1,
  "bill_number": "VENDOR-INV-2026-001",
  "reference_po_number": "PO-2026-WTR-001",
  "bill_date": "2026-02-20",
  "due_date": "2026-03-06",
  "delivery_date": "2026-02-20",
  "billing_address": {
    "attention": "Accounts Department",
    "street": "123 Industrial Estate",
    "city": "Bangalore",
    "state": "Karnataka",
    "country": "India",
    "postal_code": "560001"
  },
  "line_items": [
    {
      "po_line_item_id": 1,
      "item_id": "ITEM-WTR-500-001",
      "variant_id": 1,
      "quantity_billed": 5000,
      "rate": 8.00,
      "description": "500ml Regular Cap bottles - as per PO"
    },
    {
      "po_line_item_id": 2,
      "item_id": "ITEM-WTR-500-001",
      "variant_id": 2,
      "quantity_billed": 3000,
      "rate": 9.50,
      "description": "500ml Sports Cap bottles - as per PO"
    }
  ],
  "notes": "Invoice for goods received on 2026-02-20. All items in good condition.",
  "payment_terms": "net_45"
}
```

### 10.3 DTO Output Response

**Status:** 201 Created

```json
{
  "success": true,
  "data": {
    "id": 1,
    "bill_number": "VENDOR-INV-2026-001",
    "reference_po_number": "PO-2026-WTR-001",
    "vendor_id": 5,
    "vendor_name": "AquaPlast Industries",
    "company_id": 1,
    "bill_date": "2026-02-20",
    "due_date": "2026-03-06",
    "delivery_date": "2026-02-20",
    "status": "Open",
    "line_items": [
      {
        "id": 1,
        "item_id": "ITEM-WTR-500-001",
        "item_name": "500ml Premium Drinking Water Bottle",
        "variant_id": 1,
        "variant_name": "Regular Cap",
        "quantity_billed": 5000,
        "rate": 8.00,
        "amount": 40000,
        "tax_percentage": 5,
        "tax_amount": 2000,
        "total": 42000
      },
      {
        "id": 2,
        "item_id": "ITEM-WTR-500-001",
        "item_name": "500ml Premium Drinking Water Bottle",
        "variant_id": 2,
        "variant_name": "Sports Cap",
        "quantity_billed": 3000,
        "rate": 9.50,
        "amount": 28500,
        "tax_percentage": 5,
        "tax_amount": 1425,
        "total": 29925
      }
    ],
    "subtotal": 68500,
    "tax_total": 3425,
    "grand_total": 71925,
    "amount_due": 71925,
    "notes": "Invoice for goods received on 2026-02-20. All items in good condition.",
    "payment_terms": "net_45",
    "created_at": "2026-02-20T13:00:00Z",
    "updated_at": "2026-02-20T13:00:00Z"
  },
  "message": "Bill created successfully"
}
```

---

## Step 11: Record Payment Made

### 11.1 Payment Details

Record the payment to the vendor for the bill received.

**Endpoint:** `POST /payments`  
**Authentication:** Bearer Token + Admin Role

### 11.2 DTO Input Request

```json
{
  "payment_type": "bill_payment",
  "bill_id": 1,
  "bill_number": "VENDOR-INV-2026-001",
  "vendor_id": 5,
  "company_id": 1,
  "payment_date": "2026-03-02",
  "payment_method": "bank_transfer",
  "bank_account_id": 1,
  "amount_paid": 71925,
  "reference_number": "TXN-2026-003821",
  "notes": "Payment for purchase order PO-2026-WTR-001. Invoice VENDOR-INV-2026-001."
}
```

### 11.3 DTO Output Response

**Status:** 201 Created

```json
{
  "success": true,
  "data": {
    "id": 1,
    "payment_number": "PAY-2026-001",
    "payment_type": "bill_payment",
    "bill_id": 1,
    "bill_number": "VENDOR-INV-2026-001",
    "vendor_id": 5,
    "vendor_name": "AquaPlast Industries",
    "company_id": 1,
    "payment_date": "2026-03-02",
    "payment_method": "bank_transfer",
    "bank_account_id": 1,
    "bank_account_name": "Company Bank Account - SBI",
    "amount_paid": 71925,
    "bill_amount": 71925,
    "pending_amount": 0,
    "reference_number": "TXN-2026-003821",
    "status": "Completed",
    "notes": "Payment for purchase order PO-2026-WTR-001. Invoice VENDOR-INV-2026-001.",
    "created_at": "2026-03-02T10:30:00Z",
    "updated_at": "2026-03-02T10:30:00Z"
  },
  "message": "Payment recorded successfully"
}
```

### 11.4 Updated Bill Status

**Endpoint:** `GET /bills/1`

#### DTO Output Response

**Status:** 200 OK

```json
{
  "success": true,
  "data": {
    "id": 1,
    "bill_number": "VENDOR-INV-2026-001",
    "vendor_id": 5,
    "vendor_name": "AquaPlast Industries",
    "bill_date": "2026-02-20",
    "due_date": "2026-03-06",
    "status": "Paid",
    "grand_total": 71925,
    "amount_paid": 71925,
    "amount_due": 0,
    "payments": [
      {
        "id": 1,
        "payment_number": "PAY-2026-001",
        "payment_date": "2026-03-02",
        "amount": 71925,
        "method": "bank_transfer",
        "reference": "TXN-2026-003821"
      }
    ],
    "created_at": "2026-02-20T13:00:00Z",
    "updated_at": "2026-03-02T10:30:00Z"
  },
  "message": "Bill details retrieved successfully"
}
```

---

## Step 12: Create Sales Order

### 12.1 Sales Order Details

Create a sales order for the retail customer to purchase water bottles.

**Endpoint:** `POST /sales-orders`  
**Authentication:** Bearer Token + Admin Role

### 12.2 DTO Input Request

```json
{
  "customer_id": 12,
  "company_id": 1,
  "so_number": "SO-2026-WTR-001",
  "reference_no": "RETAIL-FEB-001",
  "sales_order_date": "2026-02-17",
  "expected_shipment_date": "2026-02-19",
  "delivery_address": {
    "attention": "Store Manager",
    "street": "789 Retail Plaza",
    "address_line2": "Store 1",
    "city": "Bangalore",
    "state": "Karnataka",
    "country": "India",
    "postal_code": "560003",
    "phone": "08041234501"
  },
  "payment_terms": "net_30",
  "delivery_method": "courier",
  "courier_company": "FedEx",
  "salesperson_id": 3,
  "notes": "Retail order for water bottles. Please ship within 2 days."
}
```

### 12.3 DTO Output Response

**Status:** 201 Created

```json
{
  "success": true,
  "data": {
    "id": 1,
    "so_number": "SO-2026-WTR-001",
    "reference_no": "RETAIL-FEB-001",
    "customer_id": 12,
    "customer_name": "Fresh Water Retail",
    "company_id": 1,
    "sales_order_date": "2026-02-17",
    "expected_shipment_date": "2026-02-19",
    "delivery_address": {
      "id": 1,
      "attention": "Store Manager",
      "street": "789 Retail Plaza",
      "city": "Bangalore",
      "state": "Karnataka",
      "country": "India",
      "postal_code": "560003"
    },
    "status": "Draft",
    "payment_terms": "net_30",
    "delivery_method": "courier",
    "courier_company": "FedEx",
    "salesperson_id": 3,
    "salesperson_name": "Rajesh Sharma",
    "notes": "Retail order for water bottles. Please ship within 2 days.",
    "line_items": [],
    "total_amount": 0,
    "tax_total": 0,
    "grand_total": 0,
    "created_at": "2026-02-17T12:15:00Z",
    "updated_at": "2026-02-17T12:15:00Z"
  },
  "message": "Sales order created successfully"
}
```

### 12.4 Add Multiple Line Items to Sales Order (With Inventory Validation)

**Endpoint:** `POST /sales-orders`  
**Authentication:** Bearer Token + Admin Role

#### DTO Input Request (With Multiple Line Items)

```json
{
  "customer_id": 12,
  "sales_order_date": "2026-02-17",
  "expected_shipment_date": "2026-02-19",
  "payment_terms": "net_30",
  "delivery_method": "courier",
  "courier_company": "FedEx",
  "salesperson_id": 3,
  "notes": "Retail order for water bottles. Please ship within 2 days.",
  "line_items": [
    {
      "item_id": "item_1a2b3c4d",
      "variant_id": 1,
      "quantity": 1000,
      "rate": 15.00,
      "variant_details": {
        "Cap Type": "Regular Cap"
      }
    },
    {
      "item_id": "item_1a2b3c4d",
      "variant_id": 2,
      "quantity": 500,
      "rate": 18.00,
      "variant_details": {
        "Cap Type": "Sports Cap"
      }
    }
  ]
}
```

#### DTO Output Response (Inventory Validation Success)

**Status:** 201 Created

```json
{
  "success": true,
  "data": {
    "id": "87654321-4321-4321-4321-210987654321",
    "so_number": "SO-20260217-0001",
    "customer_id": 12,
    "customer_name": "Fresh Water Retail",
    "sales_order_date": "2026-02-17",
    "expected_shipment_date": "2026-02-19",
    "status": "draft",
    "payment_terms": "net_30",
    "delivery_method": "courier",
    "courier_company": "FedEx",
    "salesperson_id": 3,
    "salesperson_name": "Rajesh Sharma",
    "notes": "Retail order for water bottles. Please ship within 2 days.",
    "line_items": [
      {
        "id": 1,
        "item_id": "item_1a2b3c4d",
        "item_name": "500ml Premium Drinking Water Bottle",
        "variant_id": 1,
        "variant_name": "Regular Cap",
        "quantity": 1000,
        "rate": 15.00,
        "amount": 15000,
        "available_inventory": 5000,
        "variant_details": {
          "Cap Type": "Regular Cap"
        },
        "created_at": "2026-02-17T12:15:00Z"
      },
      {
        "id": 2,
        "item_id": "item_1a2b3c4d",
        "item_name": "500ml Premium Drinking Water Bottle",
        "variant_id": 2,
        "variant_name": "Sports Cap",
        "quantity": 500,
        "rate": 18.00,
        "amount": 9000,
        "available_inventory": 3000,
        "variant_details": {
          "Cap Type": "Sports Cap"
        },
        "created_at": "2026-02-17T12:16:00Z"
      }
    ],
    "sub_total": 24000,
    "tax_amount": 0,
    "adjustment": 0,
    "total": 24000,
    "created_at": "2026-02-17T12:15:00Z",
    "updated_at": "2026-02-17T12:16:00Z"
  },
  "message": "Sales order created successfully with 2 line items. Inventory validated."
}
```

#### Error Response (Insufficient Inventory)

**Status:** 400 Bad Request

```json
{
  "success": false,
  "error": "Insufficient inventory for variant_id 1. Required: 1000 units, Available: 500 units, Item: 500ml Premium Drinking Water Bottle (Regular Cap)"
}
```

### 12.5 Inventory Reservation on Sales Order Confirmation

When a Sales Order is **confirmed**, inventory is reserved for the order.

**Endpoint:** `PUT /sales-orders/{id}/confirm`

```json
{
  "success": true,
  "message": "Sales order confirmed successfully. Inventory reserved for 1000 units (Regular Cap) and 500 units (Sports Cap)"
}
```

#### Inventory Balance After Reservation

```json
{
  "id": 1,
  "item_id": "item_1a2b3c4d",
  "variant_id": 1,
  "variant_name": "Regular Cap",
  "current_quantity": 5000,
  "reserved_quantity": 1000,
  "available_quantity": 4000,
  "in_transit_quantity": 0,
  "updated_at": "2026-02-17T13:00:00Z"
}
```

### 12.6 Inventory Tracking Summary

#### Endpoint: `GET /inventory-journal?item_id=item_1a2b3c4d`

```json
{
  "success": true,
  "data": [
    {
      "id": 2,
      "item_id": "item_1a2b3c4d",
      "variant_id": 1,
      "transaction_type": "SALES_ORDER_RESERVED",
      "quantity": 1000,
      "reference_type": "SalesOrder",
      "reference_id": "87654321-4321-4321-4321-210987654321",
      "reference_no": "SO-20260217-0001",
      "notes": "Regular Cap variant reserved for retail order",
      "created_at": "2026-02-17T13:00:00Z",
      "created_by": "admin_user"
    },
    {
      "id": 1,
      "item_id": "item_1a2b3c4d",
      "variant_id": 1,
      "transaction_type": "PURCHASE_ORDER_RECEIVED",
      "quantity": 5000,
      "reference_type": "PurchaseOrder",
      "reference_id": "12345678-1234-1234-1234-123456789012",
      "reference_no": "PO-20260217-0001",
      "notes": "500ml Regular Cap - Received from AquaPlast Industries",
      "created_at": "2026-02-17T12:00:00Z",
      "created_by": "admin_user"
    }
  ]
}
```

### 12.7 Confirm Sales Order

**Endpoint:** `PUT /sales-orders/SO-2026-WTR-001/confirm`

#### DTO Input Request

```json
{
  "action": "confirm"
}
```

#### DTO Output Response

**Status:** 200 OK

```json
{
  "success": true,
  "data": {
    "id": 1,
    "so_number": "SO-2026-WTR-001",
    "customer_id": 12,
    "customer_name": "Fresh Water Retail",
    "sales_order_date": "2026-02-17",
    "expected_shipment_date": "2026-02-19",
    "status": "Confirmed",
    "line_items": [
      {
        "id": 1,
        "item_id": "ITEM-WTR-500-001",
        "quantity": 1000,
        "rate": 15.00,
        "total": 15750
      },
      {
        "id": 2,
        "item_id": "ITEM-WTR-500-001",
        "quantity": 500,
        "rate": 18.00,
        "total": 9450
      }
    ],
    "total_amount": 24000,
    "tax_total": 1200,
    "grand_total": 25200,
    "created_at": "2026-02-17T12:15:00Z",
    "confirmed_at": "2026-02-17T12:25:00Z"
  },
  "message": "Sales order confirmed successfully"
}
```

---

## Step 13: Create Invoice

### 13.1 Invoice Details

Generate an invoice (customer invoice) from the confirmed sales order for billing purposes.

**Endpoint:** `POST /invoices`  
**Authentication:** Bearer Token + Admin Role

### 13.2 DTO Input Request

```json
{
  "customer_id": 12,
  "company_id": 1,
  "invoice_number": "INV-2026-001",
  "reference_so_number": "SO-2026-WTR-001",
  "invoice_date": "2026-02-18",
  "due_date": "2026-03-20",
  "delivery_date": "2026-02-19",
  "billing_address": {
    "attention": "Finance Department",
    "street": "456 Market Street",
    "address_line2": "Building B",
    "city": "Bangalore",
    "state": "Karnataka",
    "country": "India",
    "postal_code": "560002"
  },
  "shipping_address": {
    "attention": "Store Manager",
    "street": "789 Retail Plaza",
    "address_line2": "Store 1",
    "city": "Bangalore",
    "state": "Karnataka",
    "country": "India",
    "postal_code": "560003"
  },
  "line_items": [
    {
      "so_line_item_id": 1,
      "item_id": "ITEM-WTR-500-001",
      "variant_id": 1,
      "quantity_invoiced": 1000,
      "rate": 15.00,
      "description": "500ml Regular Cap bottles (1000 units)"
    },
    {
      "so_line_item_id": 2,
      "item_id": "ITEM-WTR-500-001",
      "variant_id": 2,
      "quantity_invoiced": 500,
      "rate": 18.00,
      "description": "500ml Sports Cap bottles (500 units)"
    }
  ],
  "notes": "Invoice for water bottle retail sale. Payment due within 30 days.",
  "payment_terms": "net_30"
}
```

### 13.3 DTO Output Response

**Status:** 201 Created

```json
{
  "success": true,
  "data": {
    "id": 1,
    "invoice_number": "INV-2026-001",
    "reference_so_number": "SO-2026-WTR-001",
    "customer_id": 12,
    "customer_name": "Fresh Water Retail",
    "company_id": 1,
    "invoice_date": "2026-02-18",
    "due_date": "2026-03-20",
    "delivery_date": "2026-02-19",
    "status": "Open",
    "billing_address": {
      "id": 1,
      "attention": "Finance Department",
      "street": "456 Market Street",
      "city": "Bangalore",
      "state": "Karnataka",
      "country": "India",
      "postal_code": "560002"
    },
    "shipping_address": {
      "id": 2,
      "attention": "Store Manager",
      "street": "789 Retail Plaza",
      "city": "Bangalore",
      "state": "Karnataka",
      "country": "India",
      "postal_code": "560003"
    },
    "line_items": [
      {
        "id": 1,
        "item_id": "ITEM-WTR-500-001",
        "item_name": "500ml Premium Drinking Water Bottle",
        "variant_id": 1,
        "variant_name": "Regular Cap",
        "quantity_invoiced": 1000,
        "rate": 15.00,
        "amount": 15000,
        "tax_percentage": 5,
        "tax_amount": 750,
        "total": 15750
      },
      {
        "id": 2,
        "item_id": "ITEM-WTR-500-001",
        "item_name": "500ml Premium Drinking Water Bottle",
        "variant_id": 2,
        "variant_name": "Sports Cap",
        "quantity_invoiced": 500,
        "rate": 18.00,
        "amount": 9000,
        "tax_percentage": 5,
        "tax_amount": 450,
        "total": 9450
      }
    ],
    "subtotal": 24000,
    "tax_total": 1200,
    "grand_total": 25200,
    "amount_due": 25200,
    "notes": "Invoice for water bottle retail sale. Payment due within 30 days.",
    "payment_terms": "net_30",
    "created_at": "2026-02-18T10:00:00Z",
    "updated_at": "2026-02-18T10:00:00Z"
  },
  "message": "Invoice created successfully"
}
```

---

## Step 14: Create Shipment

### 14.1 Shipment Details

Record the shipment of goods to the customer based on the sales order.

**Endpoint:** `POST /shipments`  
**Authentication:** Bearer Token + Admin Role

### 14.2 DTO Input Request

```json
{
  "customer_id": 12,
  "sales_order_id": "SO-2026-WTR-001",
  "so_number": "SO-2026-WTR-001",
  "company_id": 1,
  "shipment_date": "2026-02-19",
  "shipment_number": "SHIP-2026-001",
  "delivery_method": "courier",
  "courier_company": "FedEx",
  "tracking_number": "FDX-2026-78901234",
  "from_address": {
    "street": "123 Industrial Park",
    "city": "Bangalore",
    "state": "Karnataka",
    "country": "India",
    "postal_code": "560001"
  },
  "to_address": {
    "attention": "Store Manager",
    "street": "789 Retail Plaza",
    "address_line2": "Store 1",
    "city": "Bangalore",
    "state": "Karnataka",
    "country": "India",
    "postal_code": "560003"
  },
  "line_items": [
    {
      "so_line_item_id": 1,
      "item_id": "ITEM-WTR-500-001",
      "variant_id": 1,
      "quantity_shipped": 1000,
      "description": "500ml Regular Cap bottles (1000 units)"
    },
    {
      "so_line_item_id": 2,
      "item_id": "ITEM-WTR-500-001",
      "variant_id": 2,
      "quantity_shipped": 500,
      "description": "500ml Sports Cap bottles (500 units)"
    }
  ],
  "expected_delivery_date": "2026-02-21",
  "notes": "Shipment of water bottles for Fresh Water Retail store order."
}
```

### 14.3 DTO Output Response

**Status:** 201 Created

```json
{
  "success": true,
  "data": {
    "id": 1,
    "shipment_number": "SHIP-2026-001",
    "sales_order_id": "SO-2026-WTR-001",
    "so_number": "SO-2026-WTR-001",
    "customer_id": 12,
    "customer_name": "Fresh Water Retail",
    "company_id": 1,
    "shipment_date": "2026-02-19",
    "delivery_method": "courier",
    "courier_company": "FedEx",
    "tracking_number": "FDX-2026-78901234",
    "status": "Shipped",
    "from_address": {
      "id": 1,
      "street": "123 Industrial Park",
      "city": "Bangalore",
      "state": "Karnataka",
      "country": "India",
      "postal_code": "560001"
    },
    "to_address": {
      "id": 2,
      "attention": "Store Manager",
      "street": "789 Retail Plaza",
      "city": "Bangalore",
      "state": "Karnataka",
      "country": "India",
      "postal_code": "560003"
    },
    "line_items": [
      {
        "id": 1,
        "item_id": "ITEM-WTR-500-001",
        "item_name": "500ml Premium Drinking Water Bottle",
        "variant_id": 1,
        "variant_name": "Regular Cap",
        "quantity_shipped": 1000
      },
      {
        "id": 2,
        "item_id": "ITEM-WTR-500-001",
        "item_name": "500ml Premium Drinking Water Bottle",
        "variant_id": 2,
        "variant_name": "Sports Cap",
        "quantity_shipped": 500
      }
    ],
    "expected_delivery_date": "2026-02-21",
    "actual_delivery_date": null,
    "notes": "Shipment of water bottles for Fresh Water Retail store order.",
    "created_at": "2026-02-19T09:00:00Z",
    "updated_at": "2026-02-19T09:00:00Z"
  },
  "message": "Shipment created successfully"
}
```

### 14.4 Update Shipment Status to Delivered

**Endpoint:** `PUT /shipments/SHIP-2026-001/delivery-confirmation`

#### DTO Input Request

```json
{
  "actual_delivery_date": "2026-02-21",
  "delivered_by": "FedEx Courier",
  "notes": "Order received by customer in good condition.",
  "signature_required": true,
  "signed_by": "Store Manager - Rajesh"
}
```

#### DTO Output Response

**Status:** 200 OK

```json
{
  "success": true,
  "data": {
    "id": 1,
    "shipment_number": "SHIP-2026-001",
    "sales_order_id": "SO-2026-WTR-001",
    "status": "Delivered",
    "shipment_date": "2026-02-19",
    "expected_delivery_date": "2026-02-21",
    "actual_delivery_date": "2026-02-21",
    "delivered_by": "FedEx Courier",
    "signed_by": "Store Manager - Rajesh",
    "tracking_number": "FDX-2026-78901234",
    "line_items": [
      {
        "item_id": "ITEM-WTR-500-001",
        "quantity_shipped": 1500
      }
    ],
    "notes": "Order received by customer in good condition.",
    "created_at": "2026-02-19T09:00:00Z",
    "updated_at": "2026-02-21T14:30:00Z"
  },
  "message": "Shipment marked as delivered successfully"
}
```

---

## Step 15: Record Payment Received

### 15.1 Payment Receipt Details

Record the payment received from the customer for the invoice.

**Endpoint:** `POST /payments`  
**Authentication:** Bearer Token + Admin Role

### 15.2 DTO Input Request

```json
{
  "payment_type": "invoice_payment",
  "invoice_id": 1,
  "invoice_number": "INV-2026-001",
  "customer_id": 12,
  "company_id": 1,
  "payment_date": "2026-03-10",
  "payment_method": "cheque",
  "cheque_number": "CHQ-2026-5678",
  "bank_account_id": 1,
  "amount_received": 25200,
  "reference_number": "CHQ-2026-5678",
  "notes": "Cheque received from Fresh Water Retail for invoice INV-2026-001."
}
```

### 15.3 DTO Output Response

**Status:** 201 Created

```json
{
  "success": true,
  "data": {
    "id": 1,
    "payment_number": "CUST-PAY-2026-001",
    "payment_type": "invoice_payment",
    "invoice_id": 1,
    "invoice_number": "INV-2026-001",
    "customer_id": 12,
    "customer_name": "Fresh Water Retail",
    "company_id": 1,
    "payment_date": "2026-03-10",
    "payment_method": "cheque",
    "cheque_number": "CHQ-2026-5678",
    "bank_account_id": 1,
    "bank_account_name": "Company Bank Account - SBI",
    "amount_received": 25200,
    "invoice_amount": 25200,
    "pending_amount": 0,
    "reference_number": "CHQ-2026-5678",
    "status": "Completed",
    "notes": "Cheque received from Fresh Water Retail for invoice INV-2026-001.",
    "created_at": "2026-03-10T11:00:00Z",
    "updated_at": "2026-03-10T11:00:00Z"
  },
  "message": "Payment received successfully"
}
```

### 15.4 Updated Invoice Status

**Endpoint:** `GET /invoices/1`

#### DTO Output Response

**Status:** 200 OK

```json
{
  "success": true,
  "data": {
    "id": 1,
    "invoice_number": "INV-2026-001",
    "customer_id": 12,
    "customer_name": "Fresh Water Retail",
    "invoice_date": "2026-02-18",
    "due_date": "2026-03-20",
    "status": "Paid",
    "grand_total": 25200,
    "amount_received": 25200,
    "amount_due": 0,
    "line_items": [
      {
        "item_id": "ITEM-WTR-500-001",
        "quantity_invoiced": 1000,
        "rate": 15.00,
        "total": 15750
      },
      {
        "item_id": "ITEM-WTR-500-001",
        "quantity_invoiced": 500,
        "rate": 18.00,
        "total": 9450
      }
    ],
    "payments": [
      {
        "id": 1,
        "payment_number": "CUST-PAY-2026-001",
        "payment_date": "2026-03-10",
        "amount": 25200,
        "method": "cheque",
        "cheque_number": "CHQ-2026-5678"
      }
    ],
    "created_at": "2026-02-18T10:00:00Z",
    "updated_at": "2026-03-10T11:00:00Z"
  },
  "message": "Invoice details retrieved successfully"
}
```

---

---

## Complete API Reference

### Summary of Water Company Endpoints

| Operation | Method | Endpoint | Purpose |
|-----------|--------|----------|---------|
| Create Item | POST | `/items` | Create bottles, caps, labels |
| Get Item | GET | `/items/{id}` | Retrieve item details |
| Update Item | PUT | `/items/{id}` | Update pricing, description |
| Set Opening Stock | PUT | `/items/{id}/opening-stock` | Initialize inventory |
| Set Variant Stock | PUT | `/items/{id}/variants/opening-stock` | Initialize per-variant stock |
| Create Item Group | POST | `/item-groups` | Create BOM (bottle + cap + label) |
| Get Item Groups | GET | `/item-groups` | List all item groups with search |
| Get Item Group | GET | `/item-groups/{id}` | Retrieve item group details |
| Update Item Group | PUT | `/item-groups/{id}` | Update BOM components |
| Delete Item Group | DELETE | `/item-groups/{id}` | Remove item group |
| Search Item Group | GET | `/item-groups/search/by-name` | Find by name |
| Create Purchase Order | POST | `/purchase-orders` | Order from suppliers |
| Add PO Items | POST | `/purchase-orders/{po_id}/line-items` | Add bottles/caps to order |
| Create Sales Order | POST | `/sales-orders` | Sell to customers |
| Add SO Items | POST | `/sales-orders/{so_id}/line-items` | Add bottles to sale |
| Create Invoice | POST | `/invoices` | Generate customer invoice |
| Create Shipment | POST | `/shipments` | Track delivery |

### Authentication Required

All endpoints except GET require:
- **Bearer Token** in Authorization header
- **Role-based access:** Admin or SuperAdmin

---

## Inventory Valuation Methods

For water company products, use:
- **FIFO (First In, First Out)** - Recommended for perishable water bottles
- **Weighted Average** - For consistent costing across batches

---

## Key Points for Water Company

1. **Bottle Variants:** Store different cap types as variants for 500ml bottles
2. **Cap Categories:** Separate items for different sizes (20mm, 28mm, 90mm)
3. **Returnable Items:** Mark cooler bottles as returnable for deposit/return model
4. **Inventory Tracking:** Enable for all water bottles and packaging
5. **Reorder Points:** Set minimum stock levels (e.g., 1000 units for 500ml)
6. **Item Groups:** Create BOMs for packaged products (bottle + cap + label)
   - Defines which components are needed to create finished product
   - Supports exact ingredient tracking for manufacturing
   - Enables costing and profitability analysis
7. **Tax Configuration:** Water is typically 5% GST in India

---

## Next Steps

1. ✅ Create all water company products (bottles, caps, labels)
2. ✅ Set opening stock in warehouse
3. ✅ Create Item Groups (BOMs) for packaged products
4. 📋 Create Purchase Orders with suppliers
5. 📋 Create Sales Orders with customers
6. 📋 Generate Invoices and track shipments
7. 📋 Monitor inventory levels and reorder as needed

---

**Document Version:** 1.0  
**Last Updated:** February 17, 2026  
**Prepared for:** Water Company Inventory Management System
