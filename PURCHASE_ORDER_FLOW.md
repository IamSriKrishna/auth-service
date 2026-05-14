# Purchase Order Flow Documentation

## Overview
This document outlines the complete workflow for creating and managing purchase orders, from product selection through vendor payment.

---

## Flow Diagram

```
┌─────────────────────────────────────────────────────────────────────┐
│                     PURCHASE ORDER WORKFLOW                         │
└─────────────────────────────────────────────────────────────────────┘

    Step 1: Select Products                    Step 2: Create PO
    ┌──────────────────┐                       ┌──────────────────┐
    │ GET /products    │────────────────────→  │ POST /purchase-  │
    │                  │                       │      orders      │
    │ Response: List   │                       │                  │
    │ of Products with │                       │ Input: Vendor,   │
    │ variants & SKU   │                       │ Line Items,      │
    │                  │                       │ Delivery Info    │
    └──────────────────┘                       └──────────────────┘
                                                        │
                                                        ↓
                                                 ┌──────────────────┐
                                                 │ PO Created with  │
                                                 │ Status: draft    │
                                                 │                  │
                                                 │ Response: PO ID, │
                                                 │ Line Items,      │
                                                 │ Total Amount     │
                                                 └──────────────────┘
                                                        │
                                                        ↓
                                      ┌─────────────────────────────────┐
                                      │ Update PO Status (Optional)     │
                                      │ Draft → Sent → Received         │
                                      └─────────────────────────────────┘
                                                        │
                                                        ↓
    Step 3: Create Payment           ┌──────────────────────────────────┐
    ┌──────────────────────────────→ │ POST /vendor-payments            │
    │                                 │                                  │
    │ Input:                           │ Input: PO ID, Vendor ID,       │
    │ - PurchaseOrder ID              │        Amount, Payment Mode     │
    │ - Vendor ID                      │                                  │
    │ - Payment Amount                 │ Response: Payment ID, Status    │
    │ - Payment Mode                   │          pending/partial/       │
    │ - Reference Number              │          completed              │
    │                                 └──────────────────────────────────┘
    │                                                   │
    │                                                   ↓
    └─────────────────────────────────────────────────────
                                           Record Payment
                                    POST /:id/record-payment
                                           
                                    Update paid_amount &
                                    payment_status


```

---

## Step 1: Get Available Products

### Endpoint
```
GET /products
```

### Request Parameters
```
?limit=10&offset=0
```

### Response
```json
{
  "products": [
    {
      "id": "prod_8b35c2b9",
      "name": "Sleev",
      "product_details": {
        "unit": "pieces",
        "base_sku": "SL-001",
        "description": "PREMIUM Sleev",
        "manufacturer_id": 1,
        "manufacturer": {
          "id": 1,
          "name": "Manu"
        },
        "attribute_definitions": [
          {
            "key": "material",
            "options": ["plastic"]
          }
        ],
        "variants": [
          {
            "sku": "SL-001-PLAS",
            "variant_name": "plastic",
            "attribute_map": {
              "material": "plastic"
            },
            "selling_price": 10,
            "cost_price": 5,
            "stock_quantity": 0,
            "reorder_level": 0,
            "is_active": true
          }
        ]
      },
      "sales_info": {
        "account": "SALES_REVENUE",
        "selling_price": 10
      },
      "purchase_info": {
        "account": "PURCHASE_EXPENSE",
        "cost_price": 5
      },
      "created_at": "2026-04-11T15:36:17.742+05:30",
      "updated_at": "2026-04-11T15:36:17.742+05:30",
      "user_id": "2",
      "user_name": "krishna@gmail.com",
      "company_id": 1,
      "company_name": "TATA"
    },
    {
      "id": "prod_b8274277",
      "name": "Cap",
      "product_details": {
        "unit": "pieces",
        "base_sku": "CP-001",
        "description": "Premium Cap",
        "manufacturer_id": 1,
        "manufacturer": {
          "id": 1,
          "name": "Manu"
        },
        "attribute_definitions": [
          {
            "key": "color",
            "options": ["red"]
          }
        ],
        "variants": [
          {
            "sku": "CP-001-RED",
            "variant_name": "red",
            "attribute_map": {
              "color": "red"
            },
            "selling_price": 2,
            "cost_price": 0.5,
            "stock_quantity": 0,
            "reorder_level": 0,
            "is_active": true
          }
        ]
      },
      "purchase_info": {
        "account": "PURCHASE_EXPENSE",
        "cost_price": 0.5
      },
      "created_at": "2026-04-11T15:35:16.687+05:30",
      "updated_at": "2026-04-11T15:35:16.687+05:30",
      "user_id": "2",
      "user_name": "krishna@gmail.com",
      "company_id": 1,
      "company_name": "TATA"
    },
    {
      "id": "prod_eed38bb7",
      "name": "Water Bottle 300ml",
      "product_details": {
        "unit": "pieces",
        "base_sku": "WB-300",
        "description": "Premium Water Bottle",
        "manufacturer_id": 1,
        "attribute_definitions": [
          {
            "key": "size",
            "options": ["300ml"]
          }
        ],
        "variants": [
          {
            "sku": "WB-300-300ML",
            "variant_name": "300ml",
            "attribute_map": {
              "size": "300ml"
            },
            "selling_price": 10,
            "cost_price": 5,
            "stock_quantity": 0,
            "reorder_level": 0,
            "is_active": true
          }
        ]
      },
      "purchase_info": {
        "account": "PURCHASE_EXPENSE",
        "cost_price": 5
      },
      "created_at": "2026-04-11T15:23:55.75+05:30",
      "updated_at": "2026-04-11T15:23:55.75+05:30",
      "user_id": "2",
      "user_name": "krishna@gmail.com",
      "company_id": 1,
      "company_name": "TATA"
    }
  ],
  "total": 3
}
```

### Key Product Details to Use for PO
- **Product ID**: `prod_8b35c2b9` (use for creating line items)
- **SKU**: `SL-001-PLAS` (variant SKU for tracking)
- **Cost Price**: `5` (use as rate in PO line item)
- **Variant Name**: For reference in PO

---

## Step 2: Create Purchase Order

### Endpoint
```
POST /purchase-orders
Header: Authorization: Bearer {token}
```

### Request Body (Full DTO)

```json
{
  "vendor_id": 1,
  "delivery_address_type": "organization",
  "organization_name": "Our Warehouse",
  "organization_address": "123 Storage Lane, City, State 12345",
  "reference_no": "REF-2026-001",
  "date": "2026-04-24T10:30:00Z",
  "delivery_date": "2026-05-05T10:30:00Z",
  "payment_terms": "Net 30",
  "shipment_preference": "Standard",
  "line_items": [
    {
      "product_id": "prod_8b35c2b9",
      "product_name": "Sleev",
      "sku": "SL-001-PLAS",
      "account": "PURCHASE_EXPENSE",
      "quantity": 100,
      "rate": 5.00
    },
    {
      "product_id": "prod_b8274277",
      "product_name": "Cap",
      "sku": "CP-001-RED",
      "account": "PURCHASE_EXPENSE",
      "quantity": 500,
      "rate": 0.50
    },
    {
      "product_id": "prod_eed38bb7",
      "product_name": "Water Bottle 300ml",
      "sku": "WB-300-300ML",
      "account": "PURCHASE_EXPENSE",
      "quantity": 200,
      "rate": 5.00
    }
  ],
  "discount": 50,
  "discount_type": "amount",
  "tax_type": "GST",
  "tax_id": 1,
  "adjustment": 0,
  "notes": "Priority delivery requested",
  "terms_and_conditions": "Please ensure quality check before dispatch",
  "attachments": []
}
```

### Request DTO Structure

```typescript
CreatePurchaseOrderInput {
  vendor_id: uint (required)
  delivery_address_type: string (required, oneof: "organization", "customer")
  delivery_address_id?: uint
  organization_name: string
  organization_address: string
  customer_id?: uint
  reference_no: string
  date: time.Time (required, ISO 8601)
  delivery_date: time.Time (required, ISO 8601)
  payment_terms: string (required)
  shipment_preference: string
  line_items: [
    {
      product_id: string (required, UUID format)
      product_name: string (required)
      sku: string
      account: string (required, accounting code)
      quantity: float64 (required, > 0)
      rate: float64 (required, > 0)
    }
  ] (required, min 1 item)
  discount: float64 (>= 0)
  discount_type: string (oneof: "percentage", "amount")
  tax_type?: string
  tax_id?: uint
  adjustment: float64 (>= 0)
  notes: string
  terms_and_conditions: string
  attachments: []string
}
```

### Response (201 Created)

```json
{
  "success": true,
  "message": "Purchase Order created successfully",
  "data": {
    "id": "po_f8c3e9d2",
    "purchase_order_no": "PO-2026-00001",
    "vendor_id": 1,
    "vendor": {
      "id": 1,
      "display_name": "ABC Supplies",
      "company_name": "ABC Trading Co",
      "email_address": "supplier@abc.com",
      "work_phone": "+91-9876543210"
    },
    "delivery_address_type": "organization",
    "organization_name": "Our Warehouse",
    "organization_address": "123 Storage Lane, City, State 12345",
    "reference_no": "REF-2026-001",
    "date": "2026-04-24T10:30:00Z",
    "delivery_date": "2026-05-05T10:30:00Z",
    "payment_terms": "Net 30",
    "shipment_preference": "Standard",
    "line_items": [
      {
        "id": 1,
        "product_id": "prod_8b35c2b9",
        "product_name": "Sleev",
        "sku": "SL-001-PLAS",
        "account": "PURCHASE_EXPENSE",
        "quantity": 100,
        "received_quantity": 0,
        "rate": 5.00,
        "amount": 500.00
      },
      {
        "id": 2,
        "product_id": "prod_b8274277",
        "product_name": "Cap",
        "sku": "CP-001-RED",
        "account": "PURCHASE_EXPENSE",
        "quantity": 500,
        "received_quantity": 0,
        "rate": 0.50,
        "amount": 250.00
      },
      {
        "id": 3,
        "product_id": "prod_eed38bb7",
        "product_name": "Water Bottle 300ml",
        "sku": "WB-300-300ML",
        "account": "PURCHASE_EXPENSE",
        "quantity": 200,
        "received_quantity": 0,
        "rate": 5.00,
        "amount": 1000.00
      }
    ],
    "sub_total": 1750.00,
    "discount": 50.00,
    "discount_type": "amount",
    "tax_type": "GST",
    "tax_id": 1,
    "tax": {
      "id": 1,
      "name": "GST 18%",
      "tax_type": "GST",
      "rate": 18.0
    },
    "tax_amount": 306.00,
    "adjustment": 0.00,
    "total": 2006.00,
    "notes": "Priority delivery requested",
    "terms_and_conditions": "Please ensure quality check before dispatch",
    "status": "draft",
    "created_at": "2026-04-24T10:35:22.123+05:30",
    "updated_at": "2026-04-24T10:35:22.123+05:30",
    "user_id": "user_123",
    "user_name": "krishna@gmail.com",
    "company_id": 1,
    "company_name": "TATA"
  }
}
```

### Response DTO Structure

```typescript
PurchaseOrderOutput {
  id: string (UUID)
  purchase_order_no: string (auto-generated PO number)
  vendor_id: uint
  vendor?: VendorInfo {
    id: uint
    display_name: string
    company_name: string
    email_address: string
    work_phone: string
  }
  delivery_address_type: string
  delivery_address_id?: uint
  organization_name: string
  organization_address: string
  customer_id?: uint
  customer?: CustomerInfo {
    id: uint
    display_name: string
    company_name: string
    email: string
    phone: string
  }
  reference_no: string
  date: time.Time
  delivery_date: time.Time
  payment_terms: string
  shipment_preference: string
  line_items: [
    {
      id: uint
      product_id?: string
      product_name: string
      sku: string
      account: string
      quantity: float64
      received_quantity: float64
      rate: float64
      amount: float64
    }
  ]
  sub_total: float64
  discount: float64
  discount_type: string
  tax_type?: string
  tax_id?: uint
  tax?: TaxInfo {
    id: uint
    name: string
    tax_type: string
    rate: float64
  }
  tax_amount: float64
  adjustment: float64
  total: float64
  notes: string
  terms_and_conditions: string
  status: string (draft, sent, partially_received, received, cancelled)
  attachments: []string
  created_at: time.Time
  updated_at: time.Time
  user_id: string
  user_name: string
  company_id: uint
  company_name: string
}
```

### Calculations Made by System
```
Line Item Amount = Quantity × Rate
                 = 100 × 5.00 = 500.00

Sub Total = Sum of all line item amounts
          = 500.00 + 250.00 + 1000.00 = 1750.00

Discount = 50.00 (amount type)

Discounted Subtotal = Sub Total - Discount
                    = 1750.00 - 50.00 = 1700.00

Tax Amount = Discounted Subtotal × (Tax Rate / 100)
           = 1700.00 × (18 / 100) = 306.00

Total = Discounted Subtotal + Tax Amount + Adjustment
      = 1700.00 + 306.00 + 0.00 = 2006.00
```

---

## Step 3: Create Vendor Payment

### Endpoint
```
POST /vendor-payments
Header: Authorization: Bearer {token}
```

### Request Body (Full DTO)

```json
{
  "purchase_order_id": 1,
  "vendor_id": 1,
  "payment_mode": "online",
  "amount": 2006.00,
  "payment_date": "2026-04-25T14:00:00Z",
  "reference_number": "TXN-2026-00001",
  "notes": "Payment for PO-2026-00001 - Advance payment"
}
```

### Request DTO Structure

```typescript
CreateVendorPaymentInput {
  purchase_order_id: uint (required)
  vendor_id: uint (required)
  payment_mode: string (required, oneof: "cash", "online")
  amount: float64 (required, > 0)
  payment_date: time.Time (required, ISO 8601)
  reference_number: string (optional)
  notes: string (optional)
}
```

### Response (201 Created)

```json
{
  "success": true,
  "message": "Vendor payment created successfully",
  "data": {
    "id": 1,
    "payment_number": "VP-2026-00001",
    "purchase_order_id": "po_f8c3e9d2",
    "purchase_order": {
      "id": "po_f8c3e9d2",
      "purchase_order_no": "PO-2026-00001",
      "total": 2006.00,
      "status": "draft"
    },
    "vendor_id": 1,
    "vendor": {
      "id": 1,
      "display_name": "ABC Supplies",
      "company_name": "ABC Trading Co",
      "email_address": "supplier@abc.com"
    },
    "payment_mode": "online",
    "amount": 2006.00,
    "paid_amount": 0.00,
    "remaining_amount": 2006.00,
    "payment_status": "pending",
    "payment_date": "2026-04-25T14:00:00Z",
    "reference_number": "TXN-2026-00001",
    "notes": "Payment for PO-2026-00001 - Advance payment",
    "created_at": "2026-04-24T10:40:00.123+05:30",
    "updated_at": "2026-04-24T10:40:00.123+05:30",
    "created_by_user_name": "krishna@gmail.com",
    "created_by_company_name": "TATA"
  }
}
```

### Response DTO Structure

```typescript
VendorPaymentOutput {
  id: uint
  payment_number: string (auto-generated)
  purchase_order_id: string
  purchase_order?: PurchaseOrder
  vendor_id: uint
  vendor?: VendorInfo {
    id: uint
    display_name: string
    company_name: string
    email_address: string
  }
  payment_mode: string (cash, online)
  amount: float64 (total amount to be paid)
  paid_amount: float64 (amount already paid)
  remaining_amount: float64 (still pending)
  payment_status: string (pending, partial, completed)
  payment_date: time.Time
  reference_number: string (transaction ID, check number, etc.)
  notes: string
  created_at: time.Time
  updated_at: time.Time
  created_by_user_name: string
  created_by_company_name: string
}
```

---

## Step 4: Record Payment (Update Payment Received)

### Endpoint
```
POST /vendor-payments/:id/record-payment
Header: Authorization: Bearer {token}
```

### Request Body (Full DTO)

```json
{
  "paid_amount": 1500.00,
  "payment_mode": "online",
  "reference_number": "TXN-2026-00001-PARTIAL",
  "notes": "Partial payment received"
}
```

### Request DTO Structure

```typescript
RecordVendorPaymentInput {
  paid_amount: float64 (required, > 0)
  payment_mode: string (required, oneof: "cash", "online")
  reference_number: string
  notes: string
}
```

### Response (200 OK)

```json
{
  "success": true,
  "message": "Payment recorded successfully",
  "data": {
    "id": 1,
    "payment_number": "VP-2026-00001",
    "purchase_order_id": "po_f8c3e9d2",
    "vendor_id": 1,
    "vendor": {
      "id": 1,
      "display_name": "ABC Supplies",
      "company_name": "ABC Trading Co"
    },
    "payment_mode": "online",
    "amount": 2006.00,
    "paid_amount": 1500.00,
    "remaining_amount": 506.00,
    "payment_status": "partial",
    "payment_date": "2026-04-25T14:00:00Z",
    "reference_number": "TXN-2026-00001-PARTIAL",
    "notes": "Partial payment received",
    "created_at": "2026-04-24T10:40:00.123+05:30",
    "updated_at": "2026-04-24T10:45:30.456+05:30",
    "created_by_user_name": "krishna@gmail.com",
    "created_by_company_name": "TATA"
  }
}
```

### Payment Status Update Logic
```
If paid_amount == 0:
  payment_status = "pending"
  remaining_amount = amount

If 0 < paid_amount < amount:
  payment_status = "partial"
  remaining_amount = amount - paid_amount

If paid_amount >= amount:
  payment_status = "completed"
  remaining_amount = 0
```

---

## Additional Operations

### 5. Update Purchase Order Status

**Endpoint**
```
PUT /purchase-orders/:id/status
Header: Authorization: Bearer {token}
```

**Request**
```json
{
  "status": "sent"
}
```

**Valid Status Transitions**
```
draft → sent → partially_received → received
                               ↓
                           cancelled
```

---

### 6. Get Purchase Order Details

**Endpoint**
```
GET /purchase-orders/:id
Header: Authorization: Bearer {token}
```

**Response**
```json
{
  "success": true,
  "data": {
    "id": "po_f8c3e9d2",
    "purchase_order_no": "PO-2026-00001",
    "total": 2006.00,
    "status": "sent",
    ...
  }
}
```

---

### 7. Get Payments for Purchase Order

**Endpoint**
```
GET /vendor-payments/purchase-order/:poId?limit=10&offset=0
Header: Authorization: Bearer {token}
```

**Response**
```json
{
  "success": true,
  "data": {
    "vendor_payments": [
      {
        "id": 1,
        "payment_number": "VP-2026-00001",
        "paid_amount": 1500.00,
        "remaining_amount": 506.00,
        "payment_status": "partial"
      }
    ],
    "total": 1
  }
}
```

---

## Error Responses

### Invalid Request (400)
```json
{
  "success": false,
  "error": "vendor_id is required"
}
```

### Vendor Not Found (404)
```json
{
  "success": false,
  "error": "Vendor not found"
}
```

### Purchase Order Not Found (404)
```json
{
  "success": false,
  "error": "Purchase Order not found"
}
```

### Insufficient Permission (403)
```json
{
  "success": false,
  "error": "Insufficient permissions to perform this action"
}
```

### Server Error (500)
```json
{
  "success": false,
  "error": "Internal server error"
}
```

---

## Summary

| Step | Endpoint | Method | Status | Key Fields |
|------|----------|--------|--------|-----------|
| 1 | `/products` | GET | - | product_id, sku, cost_price |
| 2 | `/purchase-orders` | POST | 201 | vendor_id, line_items, total |
| 2b | `/purchase-orders/:id/status` | PUT | 200 | status (draft→sent→received) |
| 3 | `/vendor-payments` | POST | 201 | po_id, vendor_id, amount |
| 4 | `/vendor-payments/:id/record-payment` | POST | 200 | paid_amount, reference_number |

---

## Important Notes

1. **Product Selection**: Use product IDs and SKUs from the `/products` endpoint
2. **Line Item Calculation**: System automatically calculates `amount = quantity × rate`
3. **Tax Calculation**: Tax is applied on discounted subtotal
4. **Payment Status**: Automatically updated based on paid_amount vs amount
5. **Remaining Amount**: `remaining_amount = total_amount - paid_amount`
6. **Reference Numbers**: Required for tracking transactions (check number, online transaction ID, etc.)
7. **Vendor ID**: Must be a valid vendor in the system
8. **Line Items**: Minimum 1 item required per purchase order
