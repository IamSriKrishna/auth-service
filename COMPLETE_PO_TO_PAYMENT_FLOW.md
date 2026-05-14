# Complete Purchase Order to Vendor Payment Flow

## Overview
This document shows the complete step-by-step flow from creating a Purchase Order to creating a Vendor Payment with full request and response payloads.

---

## STEP 1: Create Purchase Order

### Endpoint
```
POST http://127.0.0.1:8088/purchase-orders
Header: Authorization: Bearer {token}
```

### Request Body (Full Input)
```json
{
  "vendor_id": 1,
  "delivery_address_type": "organization",
  "organization_name": "TATA",
  "organization_address": "117 Storage Lane, City",
  "reference_no": "PO-REF-001",
  "date": "2026-04-26T00:00:00Z",
  "delivery_date": "2026-05-05T00:00:00Z",
  "payment_terms": "net_30",
  "shipment_preference": "standard_shipping",
  "line_items": [
    {
      "product_id": "prod_8b35c2b9",
      "product_name": "Sleev",
      "sku": "SL-001-PLAS",
      "account": "Cost of Goods Sold",
      "quantity": 100,
      "rate": 5.00
    },
    {
      "product_id": "prod_b8274277",
      "product_name": "Cap",
      "sku": "CP-001-RED",
      "account": "Cost of Goods Sold",
      "quantity": 500,
      "rate": 0.50
    }
  ],
  "discount": 0,
  "discount_type": "amount",
  "tax_id": 1,
  "adjustment": 0,
  "notes": "Standard order",
  "terms_and_conditions": "Net 30 payment terms"
}
```

### Request Input DTO Structure
```go
type CreatePurchaseOrderInput struct {
  VendorID              uint                           // 1
  DeliveryAddressType   string                         // "organization"
  OrganizationName      string                         // "TATA"
  OrganizationAddress   string                         // "117 Storage Lane, City"
  ReferenceNo           string                         // "PO-REF-001"
  Date                  time.Time                      // "2026-04-26T00:00:00Z"
  DeliveryDate          time.Time                      // "2026-05-05T00:00:00Z"
  PaymentTerms          string                         // "net_30"
  ShipmentPreference    string                         // "standard_shipping"
  LineItems             []PurchaseOrderLineItemInput   // Array of line items
  Discount              float64                        // 0
  DiscountType          string                         // "amount"
  TaxID                 *uint                          // 1
  Adjustment            float64                        // 0
  Notes                 string                         // "Standard order"
  TermsAndConditions    string                         // "Net 30 payment terms"
}

type PurchaseOrderLineItemInput struct {
  ProductID   *string  // "prod_8b35c2b9"
  ProductName string   // "Sleev"
  SKU         string   // "SL-001-PLAS"
  Account     string   // "Cost of Goods Sold"
  Quantity    float64  // 100
  Rate        float64  // 5.00
}
```

### Response (201 Created)
```json
{
  "success": true,
  "message": "Purchase order created successfully",
  "data": {
    "id": "38f5b444-4892-4798-a4aa-c3cdfde26f40",
    "purchase_order_no": "PO-20260426-0001",
    "vendor_id": 1,
    "vendor": {
      "id": 1,
      "display_name": "Premium Plastics Manufacturing",
      "company_name": "AquaPlast Industries",
      "email_address": "rajesh.kumar@aquaplast.com",
      "work_phone": "9345927994"
    },
    "delivery_address_type": "organization",
    "organization_name": "TATA",
    "organization_address": "117 Storage Lane, City",
    "reference_no": "PO-REF-001",
    "date": "2026-04-26T00:00:00Z",
    "delivery_date": "2026-05-05T00:00:00Z",
    "payment_terms": "net_30",
    "shipment_preference": "standard_shipping",
    "line_items": [
      {
        "id": 1,
        "product_id": "prod_8b35c2b9",
        "product_name": "Sleev",
        "sku": "SL-001-PLAS",
        "account": "Cost of Goods Sold",
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
        "account": "Cost of Goods Sold",
        "quantity": 500,
        "received_quantity": 0,
        "rate": 0.50,
        "amount": 250.00
      }
    ],
    "sub_total": 750.00,
    "discount": 0,
    "discount_type": "amount",
    "tax_type": null,
    "tax_id": 1,
    "tax": {
      "id": 1,
      "name": "GST 5%",
      "tax_type": "GST",
      "rate": 5
    },
    "tax_amount": 37.50,
    "adjustment": 0,
    "total": 787.50,
    "notes": "Standard order",
    "terms_and_conditions": "Net 30 payment terms",
    "status": "draft",
    "created_at": "2026-04-26T01:15:30.123+05:30",
    "updated_at": "2026-04-26T01:15:30.123+05:30",
    "user_id": "2",
    "user_name": "krishna@gmail.com",
    "company_id": 1,
    "company_name": "TATA"
  }
}
```

### Response Output DTO Structure
```go
type PurchaseOrderOutput struct {
  ID                  string                          // "38f5b444-4892-4798-a4aa-c3cdfde26f40"
  PurchaseOrderNo     string                          // "PO-20260426-0001"
  VendorID            uint                            // 1
  Vendor              *VendorInfo                     // Vendor details
  DeliveryAddressType string                          // "organization"
  OrganizationName    string                          // "TATA"
  OrganizationAddress string                          // "117 Storage Lane, City"
  ReferenceNo         string                          // "PO-REF-001"
  Date                time.Time                       // "2026-04-26T00:00:00Z"
  DeliveryDate        time.Time                       // "2026-05-05T00:00:00Z"
  PaymentTerms        string                          // "net_30"
  ShipmentPreference  string                          // "standard_shipping"
  LineItems           []PurchaseOrderLineItemOutput   // Line items
  SubTotal            float64                         // 750.00
  Discount            float64                         // 0
  DiscountType        string                          // "amount"
  Tax                 *TaxInfo                        // Tax details
  TaxAmount           float64                         // 37.50
  Adjustment          float64                         // 0
  Total               float64                         // 787.50
  Notes               string                          // "Standard order"
  TermsAndConditions  string                          // "Net 30 payment terms"
  Status              string                          // "draft"
  CreatedAt           time.Time                       // "2026-04-26T01:15:30.123+05:30"
  UpdatedAt           time.Time                       // "2026-04-26T01:15:30.123+05:30"
  UserID              string                          // "2"
  UserName            string                          // "krishna@gmail.com"
  CompanyID           uint                            // 1
  CompanyName         string                          // "TATA"
}

type PurchaseOrderLineItemOutput struct {
  ID               uint    // 1
  ProductID        *string // "prod_8b35c2b9"
  ProductName      string  // "Sleev"
  SKU              string  // "SL-001-PLAS"
  Account          string  // "Cost of Goods Sold"
  Quantity         float64 // 100
  ReceivedQuantity float64 // 0
  Rate             float64 // 5.00
  Amount           float64 // 500.00
}
```

### Key Information to Extract
- **PO ID**: `38f5b444-4892-4798-a4aa-c3cdfde26f40` ← **Use this in Step 2**
- **PO Number**: `PO-20260426-0001`
- **Vendor ID**: `1`
- **Total Amount**: `787.50`
- **Status**: `draft`

---

## STEP 2: Create Vendor Payment

### Endpoint
```
POST http://127.0.0.1:8088/vendor-payments
Header: Authorization: Bearer {token}
```

### Request Body (Full Input)
```json
{
  "purchase_order_id": "38f5b444-4892-4798-a4aa-c3cdfde26f40",
  "vendor_id": 1,
  "payment_mode": "online",
  "amount": 787.50,
  "payment_date": "2026-04-26T10:00:00Z",
  "reference_number": "TXN-2026-00001",
  "notes": "Payment for PO-20260426-0001 - Initial payment"
}
```

### Request Input DTO Structure
```go
type CreateVendorPaymentInput struct {
  PurchaseOrderID string    // "38f5b444-4892-4798-a4aa-c3cdfde26f40" (UUID STRING - NOT NUMERIC)
  VendorID        uint      // 1
  PaymentMode     string    // "online" (oneof: "cash", "online")
  Amount          float64   // 787.50
  PaymentDate     time.Time // "2026-04-26T10:00:00Z"
  ReferenceNumber string    // "TXN-2026-00001"
  Notes           string    // "Payment for PO-20260426-0001 - Initial payment"
}
```

### Response (201 Created)
```json
{
  "success": true,
  "message": "Vendor payment created successfully",
  "data": {
    "id": 1,
    "payment_number": "VP-20260426-a1b2c3d4",
    "purchase_order_id": "38f5b444-4892-4798-a4aa-c3cdfde26f40",
    "purchase_order": {
      "id": "38f5b444-4892-4798-a4aa-c3cdfde26f40",
      "purchase_order_no": "PO-20260426-0001",
      "total": 787.50,
      "status": "draft"
    },
    "vendor_id": 1,
    "vendor": {
      "id": 1,
      "display_name": "Premium Plastics Manufacturing",
      "company_name": "AquaPlast Industries",
      "email_address": "rajesh.kumar@aquaplast.com"
    },
    "payment_mode": "online",
    "amount": 787.50,
    "paid_amount": 787.50,
    "remaining_amount": 0.00,
    "payment_status": "completed",
    "payment_date": "2026-04-26T10:00:00Z",
    "reference_number": "TXN-2026-00001",
    "notes": "Payment for PO-20260426-0001 - Initial payment",
    "created_at": "2026-04-26T01:20:15.456+05:30",
    "updated_at": "2026-04-26T01:20:15.456+05:30",
    "created_by_user_name": "krishna@gmail.com",
    "created_by_company_name": "TATA"
  }
}
```

### Response Output DTO Structure
```go
type VendorPaymentOutput struct {
  ID                   uint                  // 1
  PaymentNumber        string                // "VP-20260426-a1b2c3d4"
  PurchaseOrderID      string                // "38f5b444-4892-4798-a4aa-c3cdfde26f40"
  PurchaseOrder        *PurchaseOrder        // Full PO details (omitted if null)
  VendorID             uint                  // 1
  Vendor               *VendorInfo           // Vendor details
  PaymentMode          string                // "online"
  Amount               float64               // 787.50 (Total amount to be paid)
  PaidAmount           float64               // 787.50 (Amount already paid)
  RemainingAmount      float64               // 0.00 (Still pending)
  PaymentStatus        string                // "completed" (pending/partial/completed)
  PaymentDate          time.Time             // "2026-04-26T10:00:00Z"
  ReferenceNumber      string                // "TXN-2026-00001"
  Notes                string                // "Payment for PO-20260426-0001 - Initial payment"
  CreatedAt            time.Time             // "2026-04-26T01:20:15.456+05:30"
  UpdatedAt            time.Time             // "2026-04-26T01:20:15.456+05:30"
  CreatedByUserName    string                // "krishna@gmail.com"
  CreatedByCompanyName string                // "TATA"
}
```

### Automatic Payment Status Calculation
The system automatically calculates payment status based on the amount paid:

```
IF paid_amount == 0:
  status = "pending"
  remaining_amount = amount

IF 0 < paid_amount < amount:
  status = "partial"
  remaining_amount = amount - paid_amount

IF paid_amount >= amount:
  status = "completed"
  remaining_amount = 0
```

In this example:
- Amount: 787.50
- Paid Amount: 787.50
- → Status: **completed** (paid_amount >= amount)
- → Remaining: **0.00**

---

## STEP 3: Fetch Vendor Payments for Purchase Order

### Endpoint
```
GET http://127.0.0.1:8088/vendor-payments/purchase-order/38f5b444-4892-4798-a4aa-c3cdfde26f40?limit=10&offset=0
Header: Authorization: Bearer {token}
```

### Response (200 OK)
```json
{
  "success": true,
  "data": {
    "vendor_payments": [
      {
        "id": 1,
        "payment_number": "VP-20260426-a1b2c3d4",
        "purchase_order_id": "38f5b444-4892-4798-a4aa-c3cdfde26f40",
        "vendor_id": 1,
        "vendor": {
          "id": 1,
          "display_name": "Premium Plastics Manufacturing",
          "company_name": "AquaPlast Industries",
          "email_address": "rajesh.kumar@aquaplast.com"
        },
        "payment_mode": "online",
        "amount": 787.50,
        "paid_amount": 787.50,
        "remaining_amount": 0.00,
        "payment_status": "completed",
        "payment_date": "2026-04-26T10:00:00Z",
        "reference_number": "TXN-2026-00001",
        "notes": "Payment for PO-20260426-0001 - Initial payment",
        "created_at": "2026-04-26T01:20:15.456+05:30",
        "updated_at": "2026-04-26T01:20:15.456+05:30",
        "created_by_user_name": "krishna@gmail.com",
        "created_by_company_name": "TATA"
      }
    ],
    "total": 1
  }
}
```

### Response DTO Structure
```go
type VendorPaymentListResponse struct {
  Data  []VendorPaymentOutput `json:"vendor_payments"`
  Total int64                 `json:"total"`
}
```

---

## Summary: Data Flow

### Step 1 → Step 2 Data Mapping
```
PurchaseOrderOutput          →    CreateVendorPaymentInput
├─ id                        →    purchase_order_id  (STRING UUID)
├─ vendor_id                 →    vendor_id
├─ total                     →    (reference for amount)
└─ status: "draft"           →    (unchanged - PO stays draft)
```

### Example Values:
```
Step 1 Response:
  id: "38f5b444-4892-4798-a4aa-c3cdfde26f40"
  vendor_id: 1
  total: 787.50

Step 2 Request:
  purchase_order_id: "38f5b444-4892-4798-a4aa-c3cdfde26f40"  ← UUID STRING
  vendor_id: 1
  amount: 787.50  ← Can be less than total for partial payments
```

---

## Important Notes

### ⚠️ Critical: Purchase Order ID Type
- **Format**: UUID String (UUID v4)
- **Example**: `38f5b444-4892-4798-a4aa-c3cdfde26f40`
- **NOT**: Integer/numeric ID
- **JSON Field**: `purchase_order_id` must be a string

### Payment Scenarios

**Scenario 1: Full Payment at Once**
```json
{
  "purchase_order_id": "38f5b444-4892-4798-a4aa-c3cdfde26f40",
  "vendor_id": 1,
  "amount": 787.50,
  "payment_date": "2026-04-26T10:00:00Z"
}
→ Status: "completed", Remaining: 0.00
```

**Scenario 2: Partial Payment (Advance)**
```json
{
  "purchase_order_id": "38f5b444-4892-4798-a4aa-c3cdfde26f40",
  "vendor_id": 1,
  "amount": 400.00,
  "payment_date": "2026-04-26T10:00:00Z"
}
→ Status: "partial", Remaining: 387.50
```

**Scenario 3: Future Payment (No amount yet)**
```json
{
  "purchase_order_id": "38f5b444-4892-4798-a4aa-c3cdfde26f40",
  "vendor_id": 1,
  "amount": 0,
  "payment_date": "2026-04-30T10:00:00Z"
}
→ Status: "pending", Remaining: 787.50
```

### Additional Endpoints

**Record Additional Payment** (for partial payments)
```
POST /vendor-payments/:id/record-payment
{
  "paid_amount": 387.50,
  "payment_mode": "online",
  "reference_number": "TXN-2026-00002"
}
→ Updates paid_amount and remaining_amount
→ Updates status to "completed" if fully paid
```

**Fetch All Payments for Vendor**
```
GET /vendor-payments/vendor/1?limit=10&offset=0
```

**Get Single Payment**
```
GET /vendor-payments/1
```
