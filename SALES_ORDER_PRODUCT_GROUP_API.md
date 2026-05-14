# Sales Order with Product Groups - Complete API Documentation

## Overview

Sales Orders have been updated to use **Product Groups** instead of individual products. This enables bundled product management with automatic stock tracking for all components.

---

## 1. DTOs - Input Structures

### CreateSalesOrderInput
```go
type CreateSalesOrderInput struct {
	CustomerID           uint                      `json:"customer_id" validate:"required"`
	ReferenceNo          string                    `json:"reference_no"`
	SalesOrderDate       time.Time                 `json:"sales_order_date" validate:"required"`
	ExpectedShipmentDate time.Time                 `json:"expected_shipment_date" validate:"required"`
	PaymentTerms         string                    `json:"payment_terms" validate:"required"`
	DeliveryMethod       string                    `json:"delivery_method"`
	LineItems            []SalesOrderLineItemInput `json:"line_items" validate:"required,min=1,dive"`
	ShippingCharges      float64                   `json:"shipping_charges" validate:"gte=0"`
	TaxID                *uint                     `json:"tax_id"`
	Adjustment           float64                   `json:"adjustment"`
	CustomerNotes        string                    `json:"customer_notes"`
	TermsAndConditions   string                    `json:"terms_and_conditions"`
	TaxRate              float64                   `json:"tax_rate" validate:"gte=0,lte=100"`
	SalespersonID        *uint                     `json:"salesperson_id"`
	CreatedBy            string                    `json:"created_by"`
}
```

### SalesOrderLineItemInput
```go
type SalesOrderLineItemInput struct {
	ProductGroupID   string  `json:"product_group_id" validate:"required"`
	ProductGroupName string  `json:"product_group_name" validate:"required"`
	Quantity         float64 `json:"quantity" validate:"required,gt=0"`
	Rate             float64 `json:"rate" validate:"required,gt=0"`
	Account          string  `json:"account" validate:"required"`
}
```

### UpdateSalesOrderInput
```go
type UpdateSalesOrderInput struct {
	CustomerID           *uint                     `json:"customer_id"`
	ReferenceNo          *string                   `json:"reference_no"`
	SalesOrderDate       *time.Time                `json:"sales_order_date"`
	ExpectedShipmentDate *time.Time                `json:"expected_shipment_date"`
	PaymentTerms         *string                   `json:"payment_terms"`
	DeliveryMethod       *string                   `json:"delivery_method"`
	LineItems            []SalesOrderLineItemInput `json:"line_items" validate:"omitempty,dive"`
	ShippingCharges      *float64                  `json:"shipping_charges" validate:"omitempty,gte=0"`
	TaxID                *uint                     `json:"tax_id"`
	Adjustment           *float64                  `json:"adjustment"`
	CustomerNotes        *string                   `json:"customer_notes"`
	TermsAndConditions   *string                   `json:"terms_and_conditions"`
	TaxRate              *float64                  `json:"tax_rate" validate:"omitempty,gte=0,lte=100"`
	SalespersonID        *uint                     `json:"salesperson_id"`
}
```

### UpdateSalesOrderStatusInput
```go
type UpdateSalesOrderStatusInput struct {
	Status string `json:"status" validate:"required,oneof=draft sent confirmed partial_delivered delivered paid cancelled"`
}
```

---

## 2. DTOs - Output Structures

### SalesOrderOutput
```go
type SalesOrderOutput struct {
	ID                   string                     `json:"id"`
	SalesOrderNo         string                     `json:"sales_order_no"`
	CustomerID           uint                       `json:"customer_id"`
	ReferenceNo          string                     `json:"reference_no,omitempty"`
	Status               string                     `json:"status"`
	Date                 time.Time                  `json:"date"`
	ExpectedShipmentDate time.Time                  `json:"expected_shipment_date"`
	DeliveryMethod       string                     `json:"delivery_method,omitempty"`
	PaymentTerms         string                     `json:"payment_terms"`
	LineItems            []SalesOrderLineItemOutput `json:"line_items"`
	SubTotal             float64                    `json:"sub_total"`
	ShippingCharges      float64                    `json:"shipping_charges"`
	Adjustment           float64                    `json:"adjustment"`
	TaxRate              float64                    `json:"tax_rate"`
	TaxTotal             float64                    `json:"tax_total"`
	Total                float64                    `json:"total"`
	CustomerNotes        string                     `json:"customer_notes,omitempty"`
	TermsAndConditions   string                     `json:"terms_and_conditions,omitempty"`
	SalespersonID        *uint                      `json:"salesperson_id,omitempty"`
	CreatedAt            time.Time                  `json:"created_at"`
	UpdatedAt            time.Time                  `json:"updated_at"`
}
```

### SalesOrderLineItemOutput
```go
type SalesOrderLineItemOutput struct {
	ID                uint    `json:"id"`
	ProductGroupID    string  `json:"product_group_id"`
	ProductGroupName  string  `json:"product_group_name"`
	Account           string  `json:"account"`
	Quantity          float64 `json:"quantity"`
	DeliveredQuantity float64 `json:"delivered_quantity"`
	Rate              float64 `json:"rate"`
	Amount            float64 `json:"amount"`
}
```

---

## 3. API Endpoints

### Create Sales Order
**POST** `/api/sales-orders`

**Request:**
```json
{
  "customer_id": 1,
  "reference_no": "REF-2024-001",
  "sales_order_date": "2024-05-02T10:00:00Z",
  "expected_shipment_date": "2024-05-10T10:00:00Z",
  "payment_terms": "net_30",
  "delivery_method": "courier",
  "line_items": [
    {
      "product_group_id": "pg-bottle-kit",
      "product_group_name": "Bottle Kit - Complete Set",
      "quantity": 10,
      "rate": 250.00,
      "account": "SALES"
    },
    {
      "product_group_id": "pg-label-pack",
      "product_group_name": "Label Pack - Standard",
      "quantity": 20,
      "rate": 50.00,
      "account": "SALES"
    }
  ],
  "shipping_charges": 150.00,
  "adjustment": -10.00,
  "tax_rate": 18.0,
  "customer_notes": "Please deliver by Friday",
  "terms_and_conditions": "Payment due within 30 days",
  "salesperson_id": 5
}
```

**Response (Success):**
```json
{
  "success": true,
  "message": "Sales order created successfully",
  "data": {
    "id": "so-uuid-12345",
    "sales_order_no": "SO-2024-0001",
    "customer_id": 1,
    "reference_no": "REF-2024-001",
    "status": "draft",
    "date": "2024-05-02T10:00:00Z",
    "expected_shipment_date": "2024-05-10T10:00:00Z",
    "delivery_method": "courier",
    "payment_terms": "net_30",
    "line_items": [
      {
        "id": 1,
        "product_group_id": "pg-bottle-kit",
        "product_group_name": "Bottle Kit - Complete Set",
        "account": "SALES",
        "quantity": 10,
        "delivered_quantity": 0,
        "rate": 250.00,
        "amount": 2500.00
      },
      {
        "id": 2,
        "product_group_id": "pg-label-pack",
        "product_group_name": "Label Pack - Standard",
        "account": "SALES",
        "quantity": 20,
        "delivered_quantity": 0,
        "rate": 50.00,
        "amount": 1000.00
      }
    ],
    "sub_total": 3500.00,
    "shipping_charges": 150.00,
    "adjustment": -10.00,
    "tax_rate": 18.0,
    "tax_total": 630.00,
    "total": 4270.00,
    "customer_notes": "Please deliver by Friday",
    "terms_and_conditions": "Payment due within 30 days",
    "salesperson_id": 5,
    "created_at": "2024-05-02T10:15:30Z",
    "updated_at": "2024-05-02T10:15:30Z"
  }
}
```

---

### Get Sales Order
**GET** `/api/sales-orders/{id}`

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "so-uuid-12345",
    "sales_order_no": "SO-2024-0001",
    "customer_id": 1,
    "reference_no": "REF-2024-001",
    "status": "draft",
    "date": "2024-05-02T10:00:00Z",
    "expected_shipment_date": "2024-05-10T10:00:00Z",
    "delivery_method": "courier",
    "payment_terms": "net_30",
    "line_items": [
      {
        "id": 1,
        "product_group_id": "pg-bottle-kit",
        "product_group_name": "Bottle Kit - Complete Set",
        "account": "SALES",
        "quantity": 10,
        "delivered_quantity": 0,
        "rate": 250.00,
        "amount": 2500.00
      }
    ],
    "sub_total": 3500.00,
    "shipping_charges": 150.00,
    "adjustment": -10.00,
    "tax_rate": 18.0,
    "tax_total": 630.00,
    "total": 4270.00,
    "customer_notes": "Please deliver by Friday",
    "terms_and_conditions": "Payment due within 30 days",
    "salesperson_id": 5,
    "created_at": "2024-05-02T10:15:30Z",
    "updated_at": "2024-05-02T10:15:30Z"
  }
}
```

---

### Get All Sales Orders
**GET** `/api/sales-orders?limit=10&offset=0`

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": "so-uuid-12345",
      "sales_order_no": "SO-2024-0001",
      "customer_id": 1,
      "status": "draft",
      "date": "2024-05-02T10:00:00Z",
      "total": 4270.00,
      "line_items": []
    }
  ],
  "total": 1
}
```

---

### Update Sales Order
**PUT** `/api/sales-orders/{id}`

**Request:**
```json
{
  "reference_no": "REF-2024-001-UPDATED",
  "expected_shipment_date": "2024-05-12T10:00:00Z",
  "line_items": [
    {
      "product_group_id": "pg-bottle-kit",
      "product_group_name": "Bottle Kit - Complete Set",
      "quantity": 15,
      "rate": 260.00,
      "account": "SALES"
    }
  ],
  "tax_rate": 18.0
}
```

**Response:**
```json
{
  "success": true,
  "message": "Sales order updated successfully",
  "data": {
    "id": "so-uuid-12345",
    "sales_order_no": "SO-2024-0001",
    "reference_no": "REF-2024-001-UPDATED",
    "expected_shipment_date": "2024-05-12T10:00:00Z",
    "line_items": [
      {
        "id": 1,
        "product_group_id": "pg-bottle-kit",
        "product_group_name": "Bottle Kit - Complete Set",
        "quantity": 15,
        "rate": 260.00,
        "amount": 3900.00
      }
    ],
    "sub_total": 3900.00,
    "tax_total": 702.00,
    "total": 4752.00
  }
}
```

---

### Update Sales Order Status
**PUT** `/api/sales-orders/{id}/status`

**Request:**
```json
{
  "status": "paid"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Sales order status updated successfully",
  "data": {
    "id": "so-uuid-12345",
    "sales_order_no": "SO-2024-0001",
    "status": "paid",
    "line_items": [
      {
        "id": 1,
        "product_group_id": "pg-bottle-kit",
        "product_group_name": "Bottle Kit - Complete Set",
        "quantity": 15,
        "delivered_quantity": 0,
        "rate": 260.00,
        "amount": 3900.00
      }
    ],
    "total": 4752.00,
    "updated_at": "2024-05-02T14:30:00Z"
  }
}
```

**Status Flow:**
- `draft` → Initial state
- `sent` → Quote sent to customer
- `confirmed` → Customer accepted the quote
- `partial_delivered` → Partial shipment received
- `delivered` → Complete shipment received
- `paid` → Payment received (triggers stock deduction)
- `cancelled` → Order cancelled

---

### Get Sales Orders by Customer
**GET** `/api/sales-orders/customer/{customerId}?limit=10&offset=0`

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": "so-uuid-12345",
      "sales_order_no": "SO-2024-0001",
      "customer_id": 1,
      "status": "paid",
      "total": 4752.00,
      "line_items": []
    }
  ],
  "total": 5
}
```

---

### Get Sales Orders by Status
**GET** `/api/sales-orders/status/{status}?limit=10&offset=0`

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": "so-uuid-12345",
      "sales_order_no": "SO-2024-0001",
      "status": "paid",
      "total": 4752.00,
      "line_items": []
    }
  ],
  "total": 3
}
```

---

### Delete Sales Order
**DELETE** `/api/sales-orders/{id}`

**Response:**
```json
{
  "success": true,
  "message": "Sales order deleted successfully"
}
```

---

## 4. Database Schema

### sales_orders Table
```sql
CREATE TABLE sales_orders (
    id VARCHAR(255) PRIMARY KEY,
    sales_order_no VARCHAR(100) NOT NULL UNIQUE,
    customer_id INT NOT NULL,
    salesperson_id INT,
    reference_no VARCHAR(100),
    date DATETIME NOT NULL,
    expected_shipment_date DATETIME NOT NULL,
    delivery_method VARCHAR(100),
    payment_terms VARCHAR(50) NOT NULL,
    sub_total DOUBLE NOT NULL DEFAULT 0,
    shipping_charges DOUBLE DEFAULT 0,
    adjustment DOUBLE DEFAULT 0,
    tax_id INT,
    tax_rate DOUBLE DEFAULT 0,
    tax_total DOUBLE DEFAULT 0,
    total DOUBLE NOT NULL DEFAULT 0,
    customer_notes TEXT,
    terms_and_conditions TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'draft',
    inventory_reserved BOOLEAN DEFAULT FALSE,
    inventory_deducted BOOLEAN DEFAULT FALSE,
    reserved_date TIMESTAMP NULL,
    deducted_date TIMESTAMP NULL,
    attachments JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    created_by VARCHAR(255),
    updated_by VARCHAR(255),
    
    INDEX idx_customer (customer_id),
    INDEX idx_status (status),
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE RESTRICT
);
```

### sales_order_line_items Table
```sql
CREATE TABLE sales_order_line_items (
    id INT AUTO_INCREMENT PRIMARY KEY,
    sales_order_id VARCHAR(255) NOT NULL,
    product_group_id VARCHAR(255) NOT NULL,
    product_group_name VARCHAR(255) NOT NULL,
    quantity DOUBLE NOT NULL,
    delivered_quantity DOUBLE DEFAULT 0,
    rate DOUBLE NOT NULL,
    amount DOUBLE NOT NULL,
    account VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_sales_order_id (sales_order_id),
    INDEX idx_product_group_id (product_group_id),
    INDEX idx_so_pg_composite (sales_order_id, product_group_id),
    FOREIGN KEY (sales_order_id) REFERENCES sales_orders(id) ON DELETE CASCADE,
    FOREIGN KEY (product_group_id) REFERENCES product_groups(id) ON DELETE RESTRICT
);
```

---

## 5. Stock Management

### Stock Deduction Workflow

When a sales order status changes to **"paid"** or **"delivered"**:

1. **System retrieves** all line items from the sales order
2. **For each line item:**
   - Calls `ProductGroupInventoryService.DeductStock()`
   - Passes: ProductGroupID, Quantity, Reason (Sales Order reference)
3. **Product group inventory** is reduced:
   - CurrentStock decreases by quantity
   - AvailableStock recalculated (CurrentStock - AllocatedStock)
   - TotalSold increases by quantity
   - Transaction logged in ProductGroupTransaction table

### Stock Allocation (Optional)

To reserve stock without deducting:
```bash
POST /api/product-groups/{product_group_id}/allocate-stock

{
  "quantity": 10,
  "sales_order_id": "so-uuid-12345"
}
```

---

## 6. Error Handling

### Common Errors

**Missing Required Fields:**
```json
{
  "success": false,
  "error": "product_group_id is required for each line item"
}
```

**Invalid Customer:**
```json
{
  "success": false,
  "error": "sales order not found"
}
```

**Insufficient Stock:**
```json
{
  "success": false,
  "error": "failed to deduct stock for product group pg-001: insufficient current stock to deduct"
}
```

---

## 7. Example Workflow

### Step 1: Create Sales Order
```bash
curl -X POST http://localhost:8088/api/sales-orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <TOKEN>" \
  -d '{
    "customer_id": 1,
    "sales_order_date": "2024-05-02T10:00:00Z",
    "expected_shipment_date": "2024-05-10T10:00:00Z",
    "payment_terms": "net_30",
    "line_items": [
      {
        "product_group_id": "pg-001",
        "product_group_name": "Bottle Kit",
        "quantity": 10,
        "rate": 250,
        "account": "SALES"
      }
    ],
    "tax_rate": 18
  }'
```

### Step 2: Confirm Order
```bash
curl -X PUT http://localhost:8088/api/sales-orders/{id}/status \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <TOKEN>" \
  -d '{"status": "confirmed"}'
```

### Step 3: Mark as Paid (Triggers Stock Deduction)
```bash
curl -X PUT http://localhost:8088/api/sales-orders/{id}/status \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <TOKEN>" \
  -d '{"status": "paid"}'
```

### Step 4: Verify Stock Updated
```bash
curl -X GET "http://localhost:8088/api/product-groups/pg-001/inventory/status" \
  -H "Authorization: Bearer <TOKEN>"
```

---

## 8. Validation Rules

| Field | Rule | Example |
|-------|------|---------|
| product_group_id | Required, must exist | "pg-bottle-kit" |
| product_group_name | Required, non-empty | "Bottle Kit - Complete Set" |
| quantity | Required, > 0 | 10 |
| rate | Required, > 0 | 250.00 |
| account | Required, non-empty | "SALES" |
| tax_rate | Optional, 0-100 | 18 |
| payment_terms | Required | "net_30", "immediate" |

---

## 9. Migration

**To update existing database:**
```bash
mysql -u root auth < migrations/007_update_sales_order_to_product_groups.sql
```

This will:
- Drop old product-related columns
- Add new product_group columns
- Create proper indexes and foreign keys

