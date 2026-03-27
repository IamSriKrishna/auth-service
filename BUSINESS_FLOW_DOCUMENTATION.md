# Business Flow Documentation: Purchase to Sale to Shipment

## Overview
This document illustrates the complete business workflow from item creation through purchase orders, bills, sales orders, packages, shipments, and invoices using real data.

---

## 1. Item Created

### Request Input (Item Creation)
```json
{
  "items": [
    {
      "id": "item_b17b66ff",
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
        "preferred_vendor_id": 7,
        "preferred_vendor": {
          "id": 7,
          "display_name": "AquaPlast Industries",
          "email_address": "rajesh.kumar@aquaplast.com",
          "work_phone": "08041234567"
        },
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
      },
      "created_at": "2026-03-09T14:52:31.284+05:30",
      "updated_at": "2026-03-09T14:52:31.284+05:30",
      "user_id": "19",
      "user_name": "srik@company.com",
      "company_id": 1,
      "company_name": "Tech Innovations Pvt Ltd"
    }
  ],
  "total": 2
}
```

### Item Response Output
```json
{
  "id": "item_b17b66ff",
  "name": "500ml Premium Drinking Water Bottle",
  "type": "goods",
  "sku": "WTR-BOT-500-BASE",
  "upc": "8904220500001",
  "ean": "8904220500001",
  "description": "Premium 500ml PET drinking water bottle with tamper-proof cap",
  "unit": "piece",
  "selling_price": 16.5,
  "cost_price": 8.75,
  "currency": "INR",
  "inventory_tracking": true,
  "inventory_valuation_method": "FIFO",
  "reorder_point": 1000,
  "returnable": true,
  "created_by": "19",
  "created_by_user_name": "srik@company.com",
  "created_by_company_id": 1,
  "created_by_company_name": "Tech Innovations Pvt Ltd",
  "created_at": "2026-03-09T14:52:31.284+05:30",
  "updated_at": "2026-03-09T14:52:31.284+05:30"
}
```

### Item Variants
| Variant SKU | Cap Type | Selling Price | Cost Price | Stock |
|-------------|----------|---------------|-----------|-------|
| WTR-BOT-500-REG | Regular Cap | ₹15 | ₹8 | 5000 |
| WTR-BOT-500-SPORT | Sports Cap | ₹18 | ₹9.5 | 3000 |

---

## 2. Purchase Order (Buying from Vendor)

### Purchase Order Request Input
```json
{
  "purchase_order_no": "PO-2026-000001",
  "vendor_id": "7",
  "vendor_name": "AquaPlast Industries",
  "delivery_address_type": "vendor_warehouse",
  "delivery_address_id": "15",
  "organization_name": "Tech Innovations Pvt Ltd",
  "organization_address": "123 Business Park, Mumbai, Maharashtra, India",
  "reference_no": "REF-2026-001",
  "date": "2026-03-09T00:00:00+05:30",
  "delivery_date": "2026-03-19T00:00:00+05:30",
  "payment_terms": "net_30",
  "shipment_preference": "standard_shipping",
  "line_items": [
    {
      "item_id": "item_b17b66ff",
      "variant_sku": "WTR-BOT-500-REG",
      "account": "Cost of Goods Purchased",
      "description": "500ml Water Bottle - Regular Cap",
      "quantity": 500,
      "rate": 8.0,
      "amount": 4000.0,
      "variant_details": {
        "Cap Type": "Regular Cap"
      }
    },
    {
      "item_id": "item_b17b66ff",
      "variant_sku": "WTR-BOT-500-SPORT",
      "account": "Cost of Goods Purchased",
      "description": "500ml Water Bottle - Sports Cap",
      "quantity": 300,
      "rate": 9.5,
      "amount": 2850.0,
      "variant_details": {
        "Cap Type": "Sports Cap"
      }
    }
  ],
  "sub_total": 6850.0,
  "discount": 0,
  "discount_type": "amount",
  "tax_type": "CGST",
  "tax_id": "4",
  "tax_amount": 685.0,
  "adjustment": 0,
  "total": 7535.0,
  "notes": "First bulk order from AquaPlast Industries",
  "terms_and_conditions": "Payment due within 30 days. Goods inspected upon delivery.",
  "status": "draft"
}
```

**Note:** `created_by`, `created_by_user_name`, `created_by_company_id`, and `created_by_company_name` are **automatically populated from JWT token** and do not need to be provided in the request. They will be extracted from the authenticated user's session.
```

### Purchase Order Response Output
```json
{
  "id": "po_8f2a1d4c",
  "purchase_order_no": "PO-2026-000001",
  "vendor_id": "7",
  "vendor": {
    "id": "7",
    "name": "AquaPlast Industries",
    "email": "rajesh.kumar@aquaplast.com",
    "phone": "08041234567"
  },
  "ref_number": "REF-2026-001",
  "po_date": "2026-03-09T00:00:00+05:30",
  "delivery_date": "2026-03-19T00:00:00+05:30",
  "delivery_date_actual": null,
  "payment_terms": "net_30",
  "shipment_preference": "standard_shipping",
  "delivery_address": {
    "address_type": "vendor_warehouse",
    "organization_name": "Tech Innovations Pvt Ltd",
    "address": "123 Business Park, Mumbai, Maharashtra, India"
  },
  "line_items": [
    {
      "id": "1",
      "item_id": "item_b17b66ff",
      "item_name": "500ml Premium Drinking Water Bottle",
      "variant_sku": "WTR-BOT-500-REG",
      "variant_details": {
        "Cap Type": "Regular Cap"
      },
      "quantity": 500,
      "rate": 8.0,
      "amount": 4000.0,
      "received_quantity": 0
    },
    {
      "id": "2",
      "item_id": "item_b17b66ff",
      "item_name": "500ml Premium Drinking Water Bottle",
      "variant_sku": "WTR-BOT-500-SPORT",
      "variant_details": {
        "Cap Type": "Sports Cap"
      },
      "quantity": 300,
      "rate": 9.5,
      "amount": 2850.0,
      "received_quantity": 0
    }
  ],
  "sub_total": 6850.0,
  "discount": 0.0,
  "discount_type": "amount",
  "tax_type": "CGST",
  "tax_id": "4",
  "tax_amount": 685.0,
  "adjustment": 0.0,
  "total": 7535.0,
  "notes": "First bulk order from AquaPlast Industries",
  "terms_and_conditions": "Payment due within 30 days. Goods inspected upon delivery.",
  "status": "draft",
  "inventory_synced": false,
  "created_at": "2026-03-09T14:52:31.284+05:30",
  "created_by": "19"
}
```

### Purchase Order Summary
| Field | Value |
|-------|-------|
| **PO Number** | PO-2026-000001 |
| **Vendor** | AquaPlast Industries |
| **Total Items** | 2 Line Items (800 units total) |
| **Sub Total** | ₹6,850.00 |
| **Tax (CGST)** | ₹685.00 |
| **Total** | ₹7,535.00 |
| **Payment Terms** | Net 30 |
| **Status** | Draft |

---

## 3. Bill (Vendor Invoice)

### Bill Request Input
```json
{
  "bill_number": "INV-AP-2026-001",
  "vendor_id": "7",
  "vendor_name": "AquaPlast Industries",
  "billing_address": "AquaPlast Industries, 45 Industrial Area, Pune, Maharashtra, India",
  "order_number": "PO-2026-000001",
  "purchase_order_id": "po_8f2a1d4c",
  "bill_date": "2026-03-09T00:00:00+05:30",
  "due_date": "2026-04-08T00:00:00+05:30",
  "payment_terms": "net_30",
  "subject": "Invoice for Premium Drinking Water Bottles",
  "line_items": [
    {
      "item_id": "item_b17b66ff",
      "item_name": "500ml Premium Drinking Water Bottle",
      "variant_sku": "WTR-BOT-500-REG",
      "account": "Cost of Goods Purchased",
      "description": "500ml Water Bottle - Regular Cap",
      "quantity": 500,
      "rate": 8.0,
      "amount": 4000.0,
      "variant_details": {
        "Cap Type": "Regular Cap"
      }
    },
    {
      "item_id": "item_b17b66ff",
      "item_name": "500ml Premium Drinking Water Bottle",
      "variant_sku": "WTR-BOT-500-SPORT",
      "account": "Cost of Goods Purchased",
      "description": "500ml Water Bottle - Sports Cap",
      "quantity": 300,
      "rate": 9.5,
      "amount": 2850.0,
      "variant_details": {
        "Cap Type": "Sports Cap"
      }
    }
  ],
  "sub_total": 6850.0,
  "discount": 0,
  "tax_type": "CGST",
  "tax_id": "4",
  "tax_amount": 685.0,
  "adjustment": 0,
  "total": 7535.0,
  "notes": "Payment invoice from AquaPlast Industries",
  "status": "draft",
  "created_by": "19"
}
```

### Bill Response Output
```json
{
  "id": "bill_c3f8b6a1",
  "bill_number": "INV-AP-2026-001",
  "vendor_id": "7",
  "vendor": {
    "id": "7",
    "name": "AquaPlast Industries",
    "email": "rajesh.kumar@aquaplast.com",
    "phone": "08041234567"
  },
  "billing_address": "AquaPlast Industries, 45 Industrial Area, Pune, Maharashtra, India",
  "order_number": "PO-2026-000001",
  "purchase_order_id": "po_8f2a1d4c",
  "bill_date": "2026-03-09T00:00:00+05:30",
  "due_date": "2026-04-08T00:00:00+05:30",
  "payment_terms": "net_30",
  "subject": "Invoice for Premium Drinking Water Bottles",
  "line_items": [
    {
      "id": "1",
      "item_id": "item_b17b66ff",
      "item_name": "500ml Premium Drinking Water Bottle",
      "variant_sku": "WTR-BOT-500-REG",
      "account": "Cost of Goods Purchased",
      "description": "500ml Water Bottle - Regular Cap",
      "quantity": 500,
      "rate": 8.0,
      "amount": 4000.0,
      "inventory_synced": false,
      "variant_details": {
        "Cap Type": "Regular Cap"
      }
    },
    {
      "id": "2",
      "item_id": "item_b17b66ff",
      "item_name": "500ml Premium Drinking Water Bottle",
      "variant_sku": "WTR-BOT-500-SPORT",
      "account": "Cost of Goods Purchased",
      "description": "500ml Water Bottle - Sports Cap",
      "quantity": 300,
      "rate": 9.5,
      "amount": 2850.0,
      "inventory_synced": false,
      "variant_details": {
        "Cap Type": "Sports Cap"
      }
    }
  ],
  "sub_total": 6850.0,
  "discount": 0.0,
  "tax_type": "CGST",
  "tax_id": "4",
  "tax_amount": 685.0,
  "adjustment": 0.0,
  "total": 7535.0,
  "notes": "Payment invoice from AquaPlast Industries",
  "status": "draft",
  "inventory_synced": false,
  "created_at": "2026-03-09T14:52:31.284+05:30",
  "created_by": "19",
  "created_by_user_name": "srik@company.com"
}
```

### Bill Summary
| Field | Value |
|-------|-------|
| **Bill Number** | INV-AP-2026-001 |
| **Vendor** | AquaPlast Industries |
| **Bill Date** | 2026-03-09 |
| **Due Date** | 2026-04-08 |
| **Total Amount** | ₹7,535.00 |
| **Status** | Draft |

---

## 4. Sales Order (Selling to Customer)

### Sales Order Request Input
```json
{
  "sales_order_no": "SO-2026-000001",
  "customer_id": "15",
  "customer_name": "Premium Retail Store",
  "salesperson_id": "8",
  "salesperson_name": "Raj Kumar",
  "reference_no": "RETAIL-SO-001",
  "sales_order_date": "2026-03-10T00:00:00+05:30",
  "expected_shipment_date": "2026-03-25T00:00:00+05:30",
  "payment_terms": "net_15",
  "delivery_method": "Standard Shipping",
  "line_items": [
    {
      "item_id": "item_b17b66ff",
      "item_name": "500ml Premium Drinking Water Bottle",
      "variant_sku": "WTR-BOT-500-REG",
      "description": "500ml Water Bottle with Regular Cap - Premium Quality",
      "quantity": 250,
      "rate": 15.0,
      "amount": 3750.0,
      "variant_details": {
        "Cap Type": "Regular Cap"
      }
    },
    {
      "item_id": "item_b17b66ff",
      "item_name": "500ml Premium Drinking Water Bottle",
      "variant_sku": "WTR-BOT-500-SPORT",
      "description": "500ml Water Bottle with Sports Cap - Premium Quality",
      "quantity": 150,
      "rate": 18.0,
      "amount": 2700.0,
      "variant_details": {
        "Cap Type": "Sports Cap"
      }
    }
  ],
  "sub_total": 6450.0,
  "shipping_charges": 200.0,
  "tax_type": "IGST",
  "tax_id": "5",
  "tax_amount": 665.0,
  "adjustment": 0,
  "total": 7315.0,
  "customer_notes": "Please ensure tamper-proof packaging. Deliver to main warehouse.",
  "terms_and_conditions": "Goods are returnable within 7 days if defective. Payment due within 15 days.",
  "status": "draft",
  "created_by": "19"
}
```

### Sales Order Response Output
```json
{
  "id": "so_5e9c3b2f",
  "sales_order_no": "SO-2026-000001",
  "customer_id": "15",
  "customer": {
    "id": "15",
    "name": "Premium Retail Store",
    "email": "contact@premiumretail.com",
    "phone": "02243456789"
  },
  "salesperson_id": "8",
  "salesperson": {
    "id": "8",
    "name": "Raj Kumar",
    "email": "raj@company.com"
  },
  "reference_no": "RETAIL-SO-001",
  "sales_order_date": "2026-03-10T00:00:00+05:30",
  "expected_shipment_date": "2026-03-25T00:00:00+05:30",
  "payment_terms": "net_15",
  "delivery_method": "Standard Shipping",
  "line_items": [
    {
      "id": "1",
      "item_id": "item_b17b66ff",
      "item_name": "500ml Premium Drinking Water Bottle",
      "variant_sku": "WTR-BOT-500-REG",
      "description": "500ml Water Bottle with Regular Cap - Premium Quality",
      "quantity": 250,
      "invoiced_quantity": 0,
      "rate": 15.0,
      "amount": 3750.0,
      "variant_details": {
        "Cap Type": "Regular Cap"
      }
    },
    {
      "id": "2",
      "item_id": "item_b17b66ff",
      "item_name": "500ml Premium Drinking Water Bottle",
      "variant_sku": "WTR-BOT-500-SPORT",
      "description": "500ml Water Bottle with Sports Cap - Premium Quality",
      "quantity": 150,
      "invoiced_quantity": 0,
      "rate": 18.0,
      "amount": 2700.0,
      "variant_details": {
        "Cap Type": "Sports Cap"
      }
    }
  ],
  "sub_total": 6450.0,
  "shipping_charges": 200.0,
  "tax_type": "IGST",
  "tax_id": "5",
  "tax_amount": 665.0,
  "adjustment": 0.0,
  "total": 7315.0,
  "customer_notes": "Please ensure tamper-proof packaging. Deliver to main warehouse.",
  "terms_and_conditions": "Goods are returnable within 7 days if defective. Payment due within 15 days.",
  "status": "draft",
  "inventory_reserved": false,
  "inventory_deducted": false,
  "created_at": "2026-03-10T14:52:31.284+05:30",
  "created_by": "19"
}
```

### Sales Order Summary
| Field | Value |
|-------|-------|
| **SO Number** | SO-2026-000001 |
| **Customer** | Premium Retail Store |
| **Total Items** | 2 Line Items (400 units) |
| **Sub Total** | ₹6,450.00 |
| **Shipping** | ₹200.00 |
| **Tax (IGST)** | ₹665.00 |
| **Total** | ₹7,315.00 |
| **Status** | Draft |

---

## 5. Package (Packing for Shipment)

### Package Request Input
```json
{
  "package_slip_no": "PKG-2026-000001",
  "sales_order_id": "so_5e9c3b2f",
  "customer_id": "15",
  "customer_name": "Premium Retail Store",
  "package_date": "2026-03-20T00:00:00+05:30",
  "items": [
    {
      "sales_order_item_id": "1",
      "item_id": "item_b17b66ff",
      "item_name": "500ml Premium Drinking Water Bottle",
      "variant_sku": "WTR-BOT-500-REG",
      "ordered_qty": 250,
      "packed_qty": 250,
      "variant_details": {
        "Cap Type": "Regular Cap"
      }
    },
    {
      "sales_order_item_id": "2",
      "item_id": "item_b17b66ff",
      "item_name": "500ml Premium Drinking Water Bottle",
      "variant_sku": "WTR-BOT-500-SPORT",
      "ordered_qty": 150,
      "packed_qty": 150,
      "variant_details": {
        "Cap Type": "Sports Cap"
      }
    }
  ],
  "status": "created",
  "internal_notes": "All items packed and quality checked. Ready for shipment.",
  "created_by": "19"
}
```

### Package Response Output
```json
{
  "id": "pkg_7a4d2e8c",
  "package_slip_no": "PKG-2026-000001",
  "sales_order_id": "so_5e9c3b2f",
  "sales_order": {
    "id": "so_5e9c3b2f",
    "sales_order_no": "SO-2026-000001"
  },
  "customer_id": "15",
  "customer": {
    "id": "15",
    "name": "Premium Retail Store"
  },
  "package_date": "2026-03-20T00:00:00+05:30",
  "items": [
    {
      "id": "1",
      "package_id": "pkg_7a4d2e8c",
      "sales_order_item_id": "1",
      "item_id": "item_b17b66ff",
      "item_name": "500ml Premium Drinking Water Bottle",
      "variant_sku": "WTR-BOT-500-REG",
      "ordered_qty": 250,
      "packed_qty": 250,
      "variant_details": {
        "Cap Type": "Regular Cap"
      }
    },
    {
      "id": "2",
      "package_id": "pkg_7a4d2e8c",
      "sales_order_item_id": "2",
      "item_id": "item_b17b66ff",
      "item_name": "500ml Premium Drinking Water Bottle",
      "variant_sku": "WTR-BOT-500-SPORT",
      "ordered_qty": 150,
      "packed_qty": 150,
      "variant_details": {
        "Cap Type": "Sports Cap"
      }
    }
  ],
  "status": "created",
  "internal_notes": "All items packed and quality checked. Ready for shipment.",
  "created_at": "2026-03-20T14:52:31.284+05:30",
  "created_by": "19"
}
```

### Package Summary
| Field | Value |
|-------|-------|
| **Package Slip No** | PKG-2026-000001 |
| **Sales Order** | SO-2026-000001 |
| **Customer** | Premium Retail Store |
| **Items Packed** | 400 units total |
| **Status** | Created |

---

## 6. Shipment (Shipping the Package)

### Shipment Request Input
```json
{
  "shipment_no": "SHIP-2026-000001",
  "package_id": "pkg_7a4d2e8c",
  "sales_order_id": "so_5e9c3b2f",
  "customer_id": "15",
  "customer_name": "Premium Retail Store",
  "ship_date": "2026-03-21T00:00:00+05:30",
  "carrier": "Blue Dart Express",
  "tracking_no": "BD123456789IN",
  "tracking_url": "https://www.bluedartexpress.com/tracking?awb=BD123456789IN",
  "shipping_charges": 200.0,
  "status": "created",
  "notes": "Fragile - Handle with care. Tamper-proof packaging applied.",
  "created_by": "19"
}
```

### Shipment Response Output
```json
{
  "id": "ship_9b6f4e2a",
  "shipment_no": "SHIP-2026-000001",
  "package_id": "pkg_7a4d2e8c",
  "package": {
    "id": "pkg_7a4d2e8c",
    "package_slip_no": "PKG-2026-000001"
  },
  "sales_order_id": "so_5e9c3b2f",
  "sales_order": {
    "id": "so_5e9c3b2f",
    "sales_order_no": "SO-2026-000001"
  },
  "customer_id": "15",
  "customer": {
    "id": "15",
    "name": "Premium Retail Store",
    "email": "contact@premiumretail.com"
  },
  "ship_date": "2026-03-21T00:00:00+05:30",
  "carrier": "Blue Dart Express",
  "tracking_no": "BD123456789IN",
  "tracking_url": "https://www.bluedartexpress.com/tracking?awb=BD123456789IN",
  "shipping_charges": 200.0,
  "status": "created",
  "notes": "Fragile - Handle with care. Tamper-proof packaging applied.",
  "created_at": "2026-03-21T14:52:31.284+05:30",
  "updated_at": "2026-03-21T14:52:31.284+05:30",
  "created_by": "19"
}
```

### Shipment Summary
| Field | Value |
|-------|-------|
| **Shipment No** | SHIP-2026-000001 |
| **Carrier** | Blue Dart Express |
| **Tracking No** | BD123456789IN |
| **Ship Date** | 2026-03-21 |
| **Shipping Charges** | ₹200.00 |
| **Status** | Created |

---

## 7. Invoice (Customer Invoice)

### Invoice Request Input
```json
{
  "invoice_number": "INV-2026-000001",
  "customer_id": "15",
  "customer_name": "Premium Retail Store",
  "order_number": "SO-2026-000001",
  "sales_order_id": "so_5e9c3b2f",
  "invoice_date": "2026-03-21T00:00:00+05:30",
  "terms": "net_15",
  "due_date": "2026-04-05T00:00:00+05:30",
  "salesperson_id": "8",
  "salesperson_name": "Raj Kumar",
  "subject": "Invoice for Premium Drinking Water Bottles - Quality Goods",
  "line_items": [
    {
      "item_id": "item_b17b66ff",
      "item_name": "500ml Premium Drinking Water Bottle",
      "variant_sku": "WTR-BOT-500-REG",
      "description": "500ml Water Bottle with Regular Cap - Premium Quality",
      "quantity": 250,
      "rate": 15.0,
      "amount": 3750.0,
      "variant_details": {
        "Cap Type": "Regular Cap"
      }
    },
    {
      "item_id": "item_b17b66ff",
      "item_name": "500ml Premium Drinking Water Bottle",
      "variant_sku": "WTR-BOT-500-SPORT",
      "description": "500ml Water Bottle with Sports Cap - Premium Quality",
      "quantity": 150,
      "rate": 18.0,
      "amount": 2700.0,
      "variant_details": {
        "Cap Type": "Sports Cap"
      }
    }
  ],
  "sub_total": 6450.0,
  "shipping_charges": 200.0,
  "tax_type": "IGST",
  "tax_id": "5",
  "tax_amount": 665.0,
  "adjustment": 0,
  "total": 7315.0,
  "customer_notes": "Thank you for your business. Please ensure payment by due date.",
  "terms_and_conditions": "Payment must be received within 15 days from invoice date. Returns accepted only for defective goods within 7 days.",
  "status": "draft",
  "created_by": "19"
}
```

### Invoice Response Output
```json
{
  "id": "inv_2c7a9f6e",
  "invoice_number": "INV-2026-000001",
  "customer_id": "15",
  "customer": {
    "id": "15",
    "name": "Premium Retail Store",
    "email": "contact@premiumretail.com",
    "phone": "02243456789"
  },
  "order_number": "SO-2026-000001",
  "sales_order_id": "so_5e9c3b2f",
  "invoice_date": "2026-03-21T00:00:00+05:30",
  "terms": "net_15",
  "due_date": "2026-04-05T00:00:00+05:30",
  "salesperson_id": "8",
  "salesperson": {
    "id": "8",
    "name": "Raj Kumar",
    "email": "raj@company.com"
  },
  "subject": "Invoice for Premium Drinking Water Bottles - Quality Goods",
  "line_items": [
    {
      "id": "1",
      "invoice_id": "inv_2c7a9f6e",
      "item_id": "item_b17b66ff",
      "item_name": "500ml Premium Drinking Water Bottle",
      "variant_sku": "WTR-BOT-500-REG",
      "variant": {
        "sku": "WTR-BOT-500-REG",
        "selling_price": 15.0
      },
      "description": "500ml Water Bottle with Regular Cap - Premium Quality",
      "quantity": 250,
      "rate": 15.0,
      "amount": 3750.0,
      "inventory_synced": false,
      "variant_details": {
        "Cap Type": "Regular Cap"
      }
    },
    {
      "id": "2",
      "invoice_id": "inv_2c7a9f6e",
      "item_id": "item_b17b66ff",
      "item_name": "500ml Premium Drinking Water Bottle",
      "variant_sku": "WTR-BOT-500-SPORT",
      "variant": {
        "sku": "WTR-BOT-500-SPORT",
        "selling_price": 18.0
      },
      "description": "500ml Water Bottle with Sports Cap - Premium Quality",
      "quantity": 150,
      "rate": 18.0,
      "amount": 2700.0,
      "inventory_synced": false,
      "variant_details": {
        "Cap Type": "Sports Cap"
      }
    }
  ],
  "sub_total": 6450.0,
  "shipping_charges": 200.0,
  "tax_type": "IGST",
  "tax_id": "5",
  "tax": {
    "id": "5",
    "name": "Integrated Goods and Services Tax (IGST)",
    "percentage": 10.3
  },
  "tax_amount": 665.0,
  "adjustment": 0.0,
  "total": 7315.0,
  "customer_notes": "Thank you for your business. Please ensure payment by due date.",
  "terms_and_conditions": "Payment must be received within 15 days from invoice date. Returns accepted only for defective goods within 7 days.",
  "payment_received": false,
  "payments": [],
  "payment_splits": [],
  "status": "draft",
  "inventory_synced": false,
  "created_at": "2026-03-21T14:52:31.284+05:30",
  "updated_at": "2026-03-21T14:52:31.284+05:30",
  "created_by": "19",
  "created_by_user_name": "srik@company.com",
  "created_by_company_id": 1,
  "created_by_company_name": "Tech Innovations Pvt Ltd"
}
```

### Invoice Summary
| Field | Value |
|-------|-------|
| **Invoice Number** | INV-2026-000001 |
| **Customer** | Premium Retail Store |
| **Invoice Date** | 2026-03-21 |
| **Due Date** | 2026-04-05 |
| **Total Items** | 2 Line Items (400 units) |
| **Sub Total** | ₹6,450.00 |
| **Shipping** | ₹200.00 |
| **Tax (IGST @ 10.3%)** | ₹665.00 |
| **Total Amount Due** | ₹7,315.00 |
| **Payment Received** | No |
| **Status** | Draft |

---

## Complete Business Flow Summary

### Timeline
```
2026-03-09: Item Created (WTR-BOT-500)
            ↓
2026-03-09: Purchase Order Created (PO-2026-000001)
            ↓
2026-03-09: Bill Created (INV-AP-2026-001)
            ↓
2026-03-10: Sales Order Created (SO-2026-000001)
            ↓
2026-03-20: Package Created (PKG-2026-000001)
            ↓
2026-03-21: Shipment Created (SHIP-2026-000001)
            ↓
2026-03-21: Invoice Created (INV-2026-000001)
```

### Financial Summary

| Transaction | Quantity | Unit Price | Total Amount | Description |
|-------------|----------|------------|--------------|-------------|
| **Purchase Order** | 800 units | ₹8.75 average | ₹7,535.00 | Buying from AquaPlast Industries |
| **Bill** | 800 units | ₹8.75 average | ₹7,535.00 | Vendor invoice to pay |
| **Sales Order** | 400 units | ₹16.50 average | ₹7,315.00 | Selling to Premium Retail Store |
| **Invoice** | 400 units | ₹16.50 average | ₹7,315.00 | Customer invoice to collect |
| **Gross Profit** | 400 units | ₹7.75 per unit | ₹3,100.00 | Revenue minus Cost (before overhead) |

### Key Inventory Movements

1. **Purchase**: 800 units purchased (500 Regular + 300 Sports)
   - 500 Regular Cap @ ₹8/unit = ₹4,000
   - 300 Sports Cap @ ₹9.5/unit = ₹2,850

2. **Sales**: 400 units sold (250 Regular + 150 Sports)
   - 250 Regular Cap @ ₹15/unit = ₹3,750
   - 150 Sports Cap @ ₹18/unit = ₹2,700

3. **Remaining Stock**: 400 units
   - 250 Regular Cap (owned value: ₹2,000)
   - 150 Sports Cap (owned value: ₹1,425)

### All Model Fields Included

#### Item Fields ✓
- ID, Name, Type, SKU, UPC, EAN, Description, Unit, Selling/Cost Price, Currency, Inventory Tracking, Variants, Created By, Created At

#### Purchase Order Fields ✓
- ID, PO Number, Vendor, Delivery Address, Reference No, Date, Delivery Date, Payment Terms, Line Items (with variants), Sub Total, Discount, Tax, Adjustment, Total, Notes, Terms & Conditions, Status, Inventory Sync, Created By

#### Bill Fields ✓
- ID, Bill Number, Vendor, Billing Address, Order Number, Bill Date, Due Date, Payment Terms, Line Items (with variants), Sub Total, Discount, Tax, Adjustment, Total, Notes, Status, Attachments, Created By

#### Sales Order Fields ✓
- ID, SO Number, Customer, Salesperson, Reference No, SO Date, Expected Shipment Date, Payment Terms, Delivery Method, Line Items (with variants), Sub Total, Shipping, Tax, Adjustment, Total, Customer Notes, Terms & Conditions, Status, Inventory Reserved/Deducted, Created By

#### Package Fields ✓
- ID, Package Slip No, Sales Order, Customer, Package Date, Items (with quantities and variants), Status, Internal Notes, Created By

#### Shipment Fields ✓
- ID, Shipment No, Package, Sales Order, Customer, Ship Date, Carrier, Tracking No, Tracking URL, Shipping Charges, Status, Notes, Created By

#### Invoice Fields ✓
- ID, Invoice Number, Customer, Sales Order, Invoice Date, Terms, Due Date, Salesperson, Subject, Line Items (with variants), Sub Total, Shipping, Tax, Adjustment, Total, Customer Notes, Terms & Conditions, Payment Received, Payments, Status, Inventory Sync, Created By

---

## API Endpoints Used (Routes Configuration)

```
ITEM CREATION:
POST /auth/manage/items

PURCHASE ORDER:
POST /purchase-orders
GET /purchase-orders/:id
PATCH /purchase-orders/:id/status

BILL:
POST /bills
GET /bills/:id
PATCH /bills/:id/status

SALES ORDER:
POST /sales-orders
GET /sales-orders/:id
PATCH /sales-orders/:id/status

PACKAGE:
POST /packages
GET /packages/:id
PATCH /packages/:id/status

SHIPMENT:
POST /shipments
GET /shipments/:id
PATCH /shipments/:id/status

INVOICE:
POST /invoices
GET /invoices/:id
PATCH /invoices/:id/status

PAYMENT:
POST /payments
GET /invoices/:invoiceId/payments
```

---

## Status Transitions

### Purchase Order Status Flow
```
draft → open → received → cancelled
```

### Bill Status Flow
```
draft → open → overdue → paid → partially_paid → cancelled
```

### Sales Order Status Flow
```
draft → confirmed → packed → shipped → delivered → cancelled
```

### Package Status Flow
```
created → picking → packed → shipped → cancelled
```

### Shipment Status Flow
```
created → in_transit → delivered → cancelled
```

### Invoice Status Flow
```
draft → open → overdue → paid → partially_paid → cancelled → written_off
```

---

## Audit Trail Fields (All Entities)

- `created_at`: Timestamp when record was created
- `created_by`: User ID who created the record
- `created_by_user_name`: Username for audit trail
- `created_by_company_id`: Company ID for multi-tenant support
- `created_by_company_name`: Company name for audit trail
- `updated_at`: Last modification timestamp
- `updated_by`: User ID who last updated the record
- `updated_by_user_name`: Username for audit trail

---

## Notes

- All monetary values are in INR
- Variant tracking is enabled for products with multiple options
- Inventory tracking uses FIFO valuation method
- GST/Tax configurations differ between purchase (CGST) and sales (IGST)
- All entities support soft deletes and audit trails
- Relationships maintain referential integrity with ON_UPDATE:CASCADE and ON_DELETE:RESTRICT/SET NULL
