# API Documentation - Zoho Inventory Structure

## Table of Contents

### 1. [Setup & Configuration](#setup--configuration)
   - Organization Setup
   - Company Management
   - Regional Settings
   - Tax Configuration

### 2. [Master Data Management](#master-data-management)
   - Items & Inventory
   - Customers
   - Vendors
   - Manufacturers & Brands

### 3. [Purchase Management](#purchase-management)
   - Purchase Orders (PO)
   - Bills (Vendor Invoices)
   - Purchase Tracking
   - Vendor Management

### 4. [Sales Management](#sales-management)
   - Sales Orders (SO)
   - Invoices (Customer)
   - Sales Tracking
   - Customer Management

### 5. [Manufacturing & Bill of Materials](#manufacturing--bill-of-materials)
   - Item Groups (BOM)
   - Production Orders
   - Manufacturing Tracking
   - Inventory Consumption

### 6. [Inventory Management](#inventory-management)
   - Real-time Inventory Tracking
   - Stock Movements & Journal
   - Inventory Aggregation
   - Supply Chain Summary

### 7. [Fulfillment & Shipping](#fulfillment--shipping)
   - Packages
   - Shipments
   - Order Tracking

### 8. [Financial Management](#financial-management)
   - Payments (Invoices & Bills)
   - Salespersons
   - Tax Management

### 9. [User & Authentication](#user--authentication)
   - Authentication
   - User Management
   - Admin Management

### 10. [Utilities](#utilities)
   - Helper/Lookup Routes
   - Forward Auth
   - Support
   - Health Check

---

# Setup & Configuration

## Organization Setup

### 1. Complete Company Setup (One-time Wizard)
- **Method:** `POST`
- **Endpoint:** `/companies/setup`
- **Authentication:** Bearer Token (Required)
- **Description:** Complete company setup with all details in one request
- **Request Body:**
```json
{
  "company_name": "ABC Corporation",
  "gstin": "18AABCT1234H1Z0",
  "pan": "AAACT1234H",
  "email": "info@abc.com",
  "phone": "9876543210",
  "business_type_id": 1,
  "country": "India",
  "state": "Maharashtra",
  "city": "Mumbai"
}
```
- **Response:** Returns company with all settings configured

## Company Management

### 1. Create Company
- **Method:** `POST`
- **Endpoint:** `/companies`
- **Authentication:** Bearer Token (Required)

### 2. Get All Companies
- **Method:** `GET`
- **Endpoint:** `/companies?limit=10&offset=0`
- **Authentication:** Bearer Token (Required)
- **Pagination:** Supports limit & offset

### 3. Get Company by ID
- **Method:** `GET`
- **Endpoint:** `/companies/:id`
- **Authentication:** Bearer Token (Required)

### 4. Update Company
- **Method:** `PUT`
- **Endpoint:** `/companies/:id`
- **Authentication:** Bearer Token (Required)

### 5. Delete Company
- **Method:** `DELETE`
- **Endpoint:** `/companies/:id`
- **Authentication:** Bearer Token + SuperAdmin Role (Required)

## Company Contact Management

### 1. Upsert Company Contact
- **Method:** `PUT`
- **Endpoint:** `/companies/:id/contact`
- **Authentication:** Bearer Token (Required)
- **Request Body:**
```json
{
  "contact_person": "John Doe",
  "email": "john@abc.com",
  "phone": "9876543210"
}
```

### 2. Get Company Contact
- **Method:** `GET`
- **Endpoint:** `/companies/:id/contact`
- **Authentication:** Bearer Token (Required)

## Company Address Management

### 1. Upsert Company Address
- **Method:** `PUT`
- **Endpoint:** `/companies/:id/address`
- **Authentication:** Bearer Token (Required)
- **Request Body:**
```json
{
  "street": "123 Main Street",
  "city": "Mumbai",
  "state": "Maharashtra",
  "country": "India",
  "postal_code": "400001"
}
```

### 2. Get Company Address
- **Method:** `GET`
- **Endpoint:** `/companies/:id/address`
- **Authentication:** Bearer Token (Required)

## Bank Master Data

### 1. Create Bank
- **Method:** `POST`
- **Endpoint:** `/banks`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Description:** Add bank to master bank list
- **Request Body:**
```json
{
  "bank_name": "HDFC Bank",
  "ifsc_code": "HDFC0001234",
  "branch_name": "Mumbai Main Branch",
  "branch_code": "HDFC123",
  "address": "123 Banking Street",
  "city": "Mumbai",
  "state": "Maharashtra",
  "postal_code": "400001",
  "country": "India",
  "is_active": true
}
```

### 2. Get All Banks
- **Method:** `GET`
- **Endpoint:** `/banks?limit=10&offset=0`
- **Authentication:** None (Public)
- **Pagination:** Supports limit & offset

### 3. Get Bank by ID
- **Method:** `GET`
- **Endpoint:** `/banks/:id`
- **Authentication:** None (Public)

### 4. Update Bank
- **Method:** `PUT`
- **Endpoint:** `/banks/:id`
- **Authentication:** Bearer Token + Admin Role (Required)

### 5. Delete Bank
- **Method:** `DELETE`
- **Endpoint:** `/banks/:id`
- **Authentication:** Bearer Token + SuperAdmin Role (Required)

---

## Bank & Payment Details

### 1. Create Company Bank Detail
- **Method:** `POST`
- **Endpoint:** `/companies/:id/bank-details`
- **Authentication:** Bearer Token (Required)
- **Description:** Link bank account to company
- **Request Body:**
```json
{
  "bank_id": 1,
  "account_holder_name": "ABC Corporation",
  "account_number": "1234567890",
  "is_primary": true
}
```
- **Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "company_id": 1,
    "bank_id": 1,
    "bank": {
      "id": 1,
      "bank_name": "HDFC Bank",
      "ifsc_code": "HDFC0001234",
      "branch_name": "Mumbai Main Branch"
    },
    "account_holder_name": "ABC Corporation",
    "account_number": "1234567890",
    "is_primary": true,
    "is_active": true,
    "created_at": "2026-02-15T10:00:00Z",
    "updated_at": "2026-02-15T10:00:00Z"
  }
}
```

### 2. Get Company Bank Details
- **Method:** `GET`
- **Endpoint:** `/companies/:id/bank-details`
- **Authentication:** Bearer Token (Required)
- **Response:** Returns all bank details linked to the company with complete bank information from banks table

### 3. Update Company Bank Detail
- **Method:** `PUT`
- **Endpoint:** `/companies/bank-details/:id`
- **Authentication:** Bearer Token (Required)

### 4. Delete Company Bank Detail
- **Method:** `DELETE`
- **Endpoint:** `/companies/bank-details/:id`
- **Authentication:** Bearer Token (Required)

### 5. Upsert UPI Detail
- **Method:** `PUT`
- **Endpoint:** `/companies/:id/upi-details`
- **Authentication:** Bearer Token (Required)
- **Request Body:**
```json
{
  "upi_id": "abc@hdfc"
}
```

### 6. Get UPI Detail
- **Method:** `GET`
- **Endpoint:** `/companies/:id/upi-details`
- **Authentication:** Bearer Token (Required)

## Regional Settings

### 1. Upsert Regional Settings
- **Method:** `PUT`
- **Endpoint:** `/companies/:id/regional-settings`
- **Authentication:** Bearer Token (Required)
- **Request Body:**
```json
{
  "default_currency": "INR",
  "timezone": "Asia/Kolkata",
  "date_format": "DD/MM/YYYY"
}
```

### 2. Get Regional Settings
- **Method:** `GET`
- **Endpoint:** `/companies/:id/regional-settings`
- **Authentication:** Bearer Token (Required)

## Invoice Settings

### 1. Upsert Invoice Settings
- **Method:** `PUT`
- **Endpoint:** `/companies/:id/invoice-settings`
- **Authentication:** Bearer Token (Required)
- **Request Body:**
```json
{
  "invoice_prefix": "INV",
  "invoice_start_number": 1,
  "show_logo": true,
  "show_signature": false,
  "round_off_total": true
}
```

### 2. Get Invoice Settings
- **Method:** `GET`
- **Endpoint:** `/companies/:id/invoice-settings`
- **Authentication:** Bearer Token (Required)

## Tax Configuration

### 1. Upsert Tax Settings
- **Method:** `PUT`
- **Endpoint:** `/companies/:id/tax-settings`
- **Authentication:** Bearer Token (Required)
- **Request Body:**
```json
{
  "gst_enabled": true,
  "tax_type_id": 1
}
```

### 2. Get Tax Settings
- **Method:** `GET`
- **Endpoint:** `/companies/:id/tax-settings`
- **Authentication:** Bearer Token (Required)

### 3. Create Tax
- **Method:** `POST`
- **Endpoint:** `/taxes`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Request Body:**
```json
{
  "name": "CGST 9%",
  "tax_type": "CGST",
  "rate": 9
}
```

### 4. Get All Taxes
- **Method:** `GET`
- **Endpoint:** `/taxes?limit=10&offset=0`
- **Authentication:** Bearer Token (Required)
- **Pagination:** Supports limit & offset

### 5. Get Tax by ID
- **Method:** `GET`
- **Endpoint:** `/taxes/:id`
- **Authentication:** Bearer Token (Required)

### 6. Update Tax
- **Method:** `PUT`
- **Endpoint:** `/taxes/:id`
- **Authentication:** Bearer Token + Admin Role (Required)

### 7. Delete Tax
- **Method:** `DELETE`
- **Endpoint:** `/taxes/:id`
- **Authentication:** Bearer Token + SuperAdmin Role (Required)

---

# Master Data Management

## Items & Inventory

### 1. Create Item
- **Method:** `POST`
- **Endpoint:** `/items`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Description:** Create single or variant item with inventory tracking
- **Request Body:**
```json
{
  "name": "Water Bottle 500ml",
  "type": "goods",
  "item_details": {
    "structure": "variants",
    "unit": "Unit",
    "sku": "WB-500-001",
    "upc": "UPC-123",
    "ean": "EAN-123",
    "mpn": "MPN-123",
    "isbn": "",
    "description": "Plastic water bottle",
    "attributes": [
      {
        "key": "color",
        "options": ["blue", "red", "green"]
      }
    ],
    "variants": [
      {
        "sku": "WB-500-BLU",
        "attribute_map": {"color": "blue"},
        "selling_price": 50,
        "cost_price": 30,
        "stock_quantity": 100
      }
    ]
  },
  "sales_info": {
    "account": "Sales Revenue",
    "selling_price": 50,
    "currency": "INR",
    "description": "Sales description"
  },
  "purchase_info": {
    "account": "Purchase Account",
    "cost_price": 30,
    "currency": "INR",
    "preferred_vendor_id": 1,
    "description": "Purchase description"
  },
  "inventory": {
    "track_inventory": true,
    "inventory_account": "Inventory",
    "inventory_valuation_method": "FIFO",
    "reorder_point": 20
  },
  "return_policy": {
    "returnable": true
  }
}
```
- **Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "id": "item_bottle_001",
    "name": "Water Bottle 500ml",
    "type": "goods",
    "item_details": {
      "structure": "variants",
      "unit": "Unit",
      "sku": "WB-500-001",
      "description": "Plastic water bottle",
      "attribute_definitions": [
        {
          "key": "color",
          "options": ["blue", "red", "green"]
        }
      ],
      "variants": [
        {
          "sku": "WB-500-BLU",
          "attribute_map": {"color": "blue"},
          "selling_price": 50,
          "cost_price": 30,
          "stock_quantity": 100
        }
      ]
    },
    "sales_info": {
      "account": "Sales Revenue",
      "selling_price": 50,
      "currency": "INR"
    },
    "purchase_info": {
      "account": "Purchase Account",
      "cost_price": 30,
      "currency": "INR",
      "preferred_vendor_id": 1
    },
    "inventory": {
      "track_inventory": true,
      "inventory_account": "Inventory",
      "inventory_valuation_method": "FIFO",
      "reorder_point": 20
    },
    "return_policy": {
      "returnable": true
    },
    "created_at": "2026-02-15T10:00:00Z",
    "updated_at": "2026-02-15T10:00:00Z"
  },
  "message": "Item created successfully"
}
```

### 2. Get All Items
- **Method:** `GET`
- **Endpoint:** `/items?limit=10&offset=0`
- **Authentication:** Bearer Token (Required)
- **Pagination:** Supports limit & offset

### 3. Get Item by ID
- **Method:** `GET`
- **Endpoint:** `/items/:id`
- **Authentication:** Bearer Token (Required)

### 4. Update Item
- **Method:** `PUT`
- **Endpoint:** `/items/:id`
- **Authentication:** Bearer Token + Admin Role (Required)

### 5. Delete Item
- **Method:** `DELETE`
- **Endpoint:** `/items/:id`
- **Authentication:** Bearer Token + Admin Role (Required)

## Opening Stock

### 1. Create Opening Stock
- **Method:** `POST`
- **Endpoint:** `/items/opening-stock`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Description:** Set opening stock for items
- **Request Body:**
```json
{
  "opening_stock": 100,
  "opening_stock_rate_per_unit": 25.50
}
```
- **Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "opening_stock": 100,
    "opening_stock_rate_per_unit": 25.50,
    "updated_at": "2026-02-15T10:00:00Z"
  }
}
```

### 2. Get Opening Stock by Item
- **Method:** `GET`
- **Endpoint:** `/items/:id/opening-stock`
- **Authentication:** Bearer Token (Required)

### 3. Update Opening Stock
- **Method:** `PUT`
- **Endpoint:** `/items/opening-stock/:id`
- **Authentication:** Bearer Token + Admin Role (Required)

### 4. Get Variant Opening Stock
- **Method:** `GET`
- **Endpoint:** `/items/:itemId/variants/:variantId/opening-stock`
- **Authentication:** Bearer Token (Required)
- **Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "variant_id": 1,
    "variant_sku": "WB-500-BLU",
    "opening_stock": 100,
    "opening_stock_rate_per_unit": 25.50,
    "updated_at": "2026-02-15T10:00:00Z"
  }
}
```

### 5. Upsert Variant Opening Stock
- **Method:** `PUT`
- **Endpoint:** `/items/:itemId/variants/:variantId/opening-stock`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Request Body:**
```json
{
  "variants": [
    {
      "variant_id": 1,
      "opening_stock": 100,
      "opening_stock_rate_per_unit": 25.50
    },
    {
      "variant_id": 2,
      "opening_stock": 150,
      "opening_stock_rate_per_unit": 25.50
    }
  ]
}
```

## Stock Summary

### 1. Get Stock Summary
- **Method:** `GET`
- **Endpoint:** `/items/:items/stock-summary`
- **Authentication:** Bearer Token (Required)
- **Description:** Get current stock summary for items

---

## Vendors

### 1. Create Vendor
- **Method:** `POST`
- **Endpoint:** `/vendors`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Request Body:**
```json
{
  "vendor_name": "ABC Supplies",
  "email": "info@abcsupplies.com",
  "phone": "9876543210",
  "contact_person": "Mr. Sharma",
  "gst_in": "18AABCT1234H1Z0",
  "country": "India",
  "state": "Maharashtra",
  "city": "Mumbai"
}
```

### 2. Get All Vendors
- **Method:** `GET`
- **Endpoint:** `/vendors?limit=10&offset=0`
- **Authentication:** Bearer Token (Required)
- **Pagination:** Supports limit & offset

### 3. Get Vendor by ID
- **Method:** `GET`
- **Endpoint:** `/vendors/:id`
- **Authentication:** Bearer Token (Required)

### 4. Update Vendor
- **Method:** `PUT`
- **Endpoint:** `/vendors/:id`
- **Authentication:** Bearer Token + Admin Role (Required)

### 5. Delete Vendor
- **Method:** `DELETE`
- **Endpoint:** `/vendors/:id`
- **Authentication:** Bearer Token + Admin Role (Required)

---

## Customers

### 1. Create Customer
- **Method:** `POST`
- **Endpoint:** `/customers`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Request Body:**
```json
{
  "customer_name": "John's Enterprise",
  "email": "john@enterprise.com",
  "phone": "9876543210",
  "contact_person": "John Doe",
  "gst_in": "18AABCT1234H1Z0",
  "country": "India",
  "state": "Maharashtra",
  "city": "Mumbai",
  "credit_limit": 100000
}
```

### 2. Get All Customers
- **Method:** `GET`
- **Endpoint:** `/customers?limit=10&offset=0`
- **Authentication:** Bearer Token (Required)
- **Pagination:** Supports limit & offset

### 3. Get Customer by ID
- **Method:** `GET`
- **Endpoint:** `/customers/:id`
- **Authentication:** Bearer Token (Required)

### 4. Update Customer
- **Method:** `PUT`
- **Endpoint:** `/customers/:id`
- **Authentication:** Bearer Token + Admin Role (Required)

### 5. Delete Customer
- **Method:** `DELETE`
- **Endpoint:** `/customers/:id`
- **Authentication:** Bearer Token + Admin Role (Required)

---

## Manufacturers & Brands

### 1. Create Manufacturer
- **Method:** `POST`
- **Endpoint:** `/manufacturers`
- **Authentication:** Bearer Token + SuperAdmin Role (Required)
- **Request Body:**
```json
{
  "name": "Dell Inc.",
  "description": "Computer manufacturer",
  "country": "USA"
}
```

### 2. Get All Manufacturers
- **Method:** `GET`
- **Endpoint:** `/manufacturers?limit=10&offset=0`
- **Authentication:** None (Public)
- **Pagination:** Supports limit & offset

### 3. Get Manufacturer by ID
- **Method:** `GET`
- **Endpoint:** `/manufacturers/:id`
- **Authentication:** None (Public)

### 4. Update Manufacturer
- **Method:** `PUT`
- **Endpoint:** `/manufacturers/:id`
- **Authentication:** Bearer Token + SuperAdmin Role (Required)

### 5. Delete Manufacturer
- **Method:** `DELETE`
- **Endpoint:** `/manufacturers/:id`
- **Authentication:** Bearer Token + SuperAdmin Role (Required)

### 6. Create Brand
- **Method:** `POST`
- **Endpoint:** `/brands`
- **Authentication:** Bearer Token + SuperAdmin Role (Required)
- **Request Body:**
```json
{
  "name": "Dell",
  "description": "Dell brand products",
  "logo_url": "https://example.com/dell-logo.png"
}
```

### 7. Get All Brands
- **Method:** `GET`
- **Endpoint:** `/brands?limit=10&offset=0`
- **Authentication:** None (Public)
- **Pagination:** Supports limit & offset

### 8. Get Brand by ID
- **Method:** `GET`
- **Endpoint:** `/brands/:id`
- **Authentication:** None (Public)

### 9. Update Brand
- **Method:** `PUT`
- **Endpoint:** `/brands/:id`
- **Authentication:** Bearer Token + SuperAdmin Role (Required)

### 10. Delete Brand
- **Method:** `DELETE`
- **Endpoint:** `/brands/:id`
- **Authentication:** Bearer Token + SuperAdmin Role (Required)

---

# Purchase Management

## Purchase Orders (Core)

### 1. Create Purchase Order
- **Method:** `POST`
- **Endpoint:** `/purchase-orders`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Description:** Create PO with line items and automatic inventory tracking
- **Request Body:**
```json
{
  "vendor_id": 1,
  "reference_no": "REF-001",
  "date": "2026-02-07",
  "delivery_date": "2026-02-14",
  "payment_terms": "net_30",
  "delivery_address_type": "organization",
  "organization_name": "Our Warehouse",
  "organization_address": "123 Business Street",
  "shipment_preference": "Air",
  "line_items": [
    {
      "item_id": "item_001",
      "quantity": 100,
      "rate": 150,
      "account": "Purchases"
    }
  ],
  "discount": 1000,
  "discount_type": "amount",
  "tax_id": 1,
  "adjustment": 100,
  "notes": "Notes for vendor",
  "terms_and_conditions": "Payment due within 30 days"
}
```
- **Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "id": "po_001",
    "purchase_order_no": "PO-001",
    "vendor_id": 1,
    "vendor": {
      "id": 1,
      "display_name": "ABC Supplies"
    },
    "date": "2026-02-07T00:00:00Z",
    "delivery_date": "2026-02-14T00:00:00Z",
    "delivery_address_type": "organization",
    "organization_name": "Our Warehouse",
    "payment_terms": "net_30",
    "reference_no": "REF-001",
    "line_items": [
      {
        "id": 1,
        "purchase_order_id": "po_001",
        "item_id": "item_001",
        "quantity": 100,
        "received_quantity": 0,
        "rate": 150,
        "amount": 15000,
        "account": "Purchases"
      }
    ],
    "sub_total": 15000,
    "discount": 1000,
    "tax_amount": 2240,
    "adjustment": 100,
    "total": 16340,
    "status": "draft",
    "inventory_synced": false,
    "created_at": "2026-02-07T10:00:00Z",
    "updated_at": "2026-02-07T10:00:00Z"
  },
  "message": "Purchase Order created successfully"
}
```

### 2. Get All Purchase Orders
- **Method:** `GET`
- **Endpoint:** `/purchase-orders?limit=10&offset=0`
- **Authentication:** Bearer Token (Required)
- **Pagination:** Supports limit & offset

### 3. Get Purchase Order by ID
- **Method:** `GET`
- **Endpoint:** `/purchase-orders/:id`
- **Authentication:** Bearer Token (Required)

### 4. Update Purchase Order
- **Method:** `PUT`
- **Endpoint:** `/purchase-orders/:id`
- **Authentication:** Bearer Token + Admin Role (Required)

### 5. Delete Purchase Order
- **Method:** `DELETE`
- **Endpoint:** `/purchase-orders/:id`
- **Authentication:** Bearer Token + Admin Role (Required)

## Purchase Order Status Management

### 1. Update PO Status
- **Method:** `PATCH`
- **Endpoint:** `/purchase-orders/:id/status`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Description:** Update PO status (triggers inventory sync on receipt)
- **Request Body:**
```json
{
  "status": "received"
}
```
- **Status Values:** draft, sent, partially_received, received, cancelled

### 2. Get POs by Status
- **Method:** `GET`
- **Endpoint:** `/purchase-orders/status/:status`
- **Authentication:** Bearer Token (Required)

## Purchase Order Filtering

### 1. Get POs by Vendor
- **Method:** `GET`
- **Endpoint:** `/purchase-orders/vendor/:vendorId`
- **Authentication:** Bearer Token (Required)

### 2. Get POs by Customer
- **Method:** `GET`
- **Endpoint:** `/purchase-orders/customer/:customerId`
- **Authentication:** Bearer Token (Required)

---

## Bills (Vendor Invoices)

### 1. Create Bill
- **Method:** `POST`
- **Endpoint:** `/bills`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Description:** Create bill from vendor (links to PO)
- **Request Body:**
```json
{
  "vendor_id": 1,
  "billing_address": "123 Vendor Street",
  "order_number": "PO-12345",
  "purchase_order_id": "po_001",
  "bill_date": "2026-02-07",
  "due_date": "2026-03-07",
  "payment_terms": "net_30",
  "subject": "Bill for materials",
  "line_items": [
    {
      "item_id": "item_001",
      "quantity": 100,
      "rate": 150,
      "description": "Raw materials",
      "account": "Purchases"
    }
  ],
  "discount": 1000,
  "tax_id": 1,
  "adjustment": 500,
  "notes": "Thank you for your supply"
}
```
- **Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "id": "bill_001",
    "bill_number": "BILL-001",
    "vendor_id": 1,
    "purchase_order_id": "po_001",
    "bill_date": "2026-02-07T00:00:00Z",
    "due_date": "2026-03-07T00:00:00Z",
    "payment_terms": "net_30",
    "subject": "Bill for materials",
    "line_items": [
      {
        "id": 1,
        "bill_id": "bill_001",
        "item_id": "item_001",
        "quantity": 100,
        "rate": 150,
        "amount": 15000,
        "account": "Purchases",
        "description": "Raw materials"
      }
    ],
    "sub_total": 15000,
    "discount": 1000,
    "tax_amount": 2240,
    "adjustment": 500,
    "total": 16740,
    "status": "draft",
    "inventory_synced": false,
    "created_at": "2026-02-07T10:00:00Z",
    "updated_at": "2026-02-07T10:00:00Z"
  },
  "message": "Bill created successfully"
}
```

### 2. Get All Bills
- **Method:** `GET`
- **Endpoint:** `/bills?limit=10&offset=0`
- **Authentication:** Bearer Token (Required)
- **Pagination:** Supports limit & offset

### 3. Get Bill by ID
- **Method:** `GET`
- **Endpoint:** `/bills/:id`
- **Authentication:** Bearer Token (Required)

### 4. Update Bill
- **Method:** `PUT`
- **Endpoint:** `/bills/:id`
- **Authentication:** Bearer Token + Admin Role (Required)

### 5. Delete Bill
- **Method:** `DELETE`
- **Endpoint:** `/bills/:id`
- **Authentication:** Bearer Token + Admin Role (Required)

## Bill Status Management

### 1. Update Bill Status
- **Method:** `PATCH`
- **Endpoint:** `/bills/:id/status`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Request Body:**
```json
{
  "status": "sent"
}
```
- **Status Values:** draft, sent, partial, paid, overdue, void

### 2. Get Bills by Status
- **Method:** `GET`
- **Endpoint:** `/bills/status/:status`
- **Authentication:** Bearer Token (Required)

### 3. Get Bills by Vendor
- **Method:** `GET`
- **Endpoint:** `/bills/vendor/:vendorId`
- **Authentication:** Bearer Token (Required)

---

# Sales Management

## Sales Orders (Core)

### 1. Create Sales Order
- **Method:** `POST`
- **Endpoint:** `/sales-orders`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Description:** Create SO (automatically reserves inventory)
- **Request Body:**
```json
{
  "customer_id": 1,
  "reference_no": "REF-001",
  "sales_order_date": "2026-02-07",
  "expected_shipment_date": "2026-02-14",
  "payment_terms": "net_30",
  "delivery_method": "courier",
  "salesperson_id": 1,
  "line_items": [
    {
      "item_id": "item_001",
      "quantity": 10,
      "rate": 5000
    }
  ],
  "shipping_charges": 500,
  "tax_id": 1,
  "adjustment": 100,
  "customer_notes": "Handle with care",
  "terms_and_conditions": "Standard terms apply"
}
```
- **Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "id": "so_001",
    "sales_order_no": "SO-001",
    "customer_id": 1,
    "customer": {
      "id": 1,
      "display_name": "John's Enterprise"
    },
    "salesperson_id": 1,
    "reference_no": "REF-001",
    "sales_order_date": "2026-02-07T00:00:00Z",
    "expected_shipment_date": "2026-02-14T00:00:00Z",
    "payment_terms": "net_30",
    "delivery_method": "courier",
    "line_items": [
      {
        "id": 1,
        "sales_order_id": "so_001",
        "item_id": "item_001",
        "quantity": 10,
        "invoiced_quantity": 0,
        "rate": 5000,
        "amount": 50000
      }
    ],
    "sub_total": 50000,
    "shipping_charges": 500,
    "tax_amount": 9000,
    "adjustment": 100,
    "total": 59600,
    "status": "draft",
    "inventory_reserved": false,
    "inventory_deducted": false,
    "created_at": "2026-02-07T10:00:00Z",
    "updated_at": "2026-02-07T10:00:00Z"
  },
  "message": "Sales Order created successfully"
}
```

### 2. Get All Sales Orders
- **Method:** `GET`
- **Endpoint:** `/sales-orders?limit=10&offset=0`
- **Authentication:** Bearer Token (Required)
- **Pagination:** Supports limit & offset

### 3. Get Sales Order by ID
- **Method:** `GET`
- **Endpoint:** `/sales-orders/:id`
- **Authentication:** Bearer Token (Required)

### 4. Update Sales Order
- **Method:** `PUT`
- **Endpoint:** `/sales-orders/:id`
- **Authentication:** Bearer Token + Admin Role (Required)

### 5. Delete Sales Order
- **Method:** `DELETE`
- **Endpoint:** `/sales-orders/:id`
- **Authentication:** Bearer Token + Admin Role (Required)

## Sales Order Status Management

### 1. Update SO Status
- **Method:** `PATCH`
- **Endpoint:** `/sales-orders/:id/status`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Description:** Update SO status (triggers inventory reservation)
- **Request Body:**
```json
{
  "status": "confirmed"
}
```
- **Status Values:** draft, sent, confirmed, partial_shipped, shipped, delivered, cancelled

### 2. Get SOs by Status
- **Method:** `GET`
- **Endpoint:** `/sales-orders/status/:status`
- **Authentication:** Bearer Token (Required)

## Sales Order Filtering

### 1. Get SOs by Customer
- **Method:** `GET`
- **Endpoint:** `/sales-orders/customer/:customerId`
- **Authentication:** Bearer Token (Required)

---

## Invoices (Customer Billing)

### 1. Create Invoice
- **Method:** `POST`
- **Endpoint:** `/invoices`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Description:** Create invoice (deducts inventory on creation)
- **Request Body:**
```json
{
  "sales_order_id": "so_001",
  "customer_id": 1,
  "invoice_date": "2026-02-07",
  "due_date": "2026-03-07",
  "terms": "net_30",
  "subject": "Invoice for products",
  "salesperson_id": 1,
  "order_number": "SO-001",
  "line_items": [
    {
      "item_id": "item_001",
      "quantity": 10,
      "rate": 5000
    }
  ],
  "shipping_charges": 500,
  "tax_id": 1,
  "adjustment": 100,
  "customer_notes": "Thank you for your business",
  "terms_and_conditions": "Payment due on date specified"
}
```
- **Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "id": "inv_001",
    "invoice_number": "INV-001",
    "customer_id": 1,
    "customer": {
      "id": 1,
      "display_name": "John's Enterprise"
    },
    "sales_order_id": "so_001",
    "order_number": "SO-001",
    "invoice_date": "2026-02-07T00:00:00Z",
    "due_date": "2026-03-07T00:00:00Z",
    "terms": "net_30",
    "subject": "Invoice for products",
    "salesperson_id": 1,
    "line_items": [
      {
        "id": 1,
        "invoice_id": "inv_001",
        "item_id": "item_001",
        "quantity": 10,
        "rate": 5000,
        "amount": 50000,
        "inventory_synced": false
      }
    ],
    "sub_total": 50000,
    "shipping_charges": 500,
    "tax_amount": 9000,
    "adjustment": 100,
    "total": 59600,
    "status": "draft",
    "inventory_synced": false,
    "payment_received": false,
    "created_at": "2026-02-07T10:00:00Z",
    "updated_at": "2026-02-07T10:00:00Z"
  },
  "message": "Invoice created successfully"
}
```

### 2. Get All Invoices
- **Method:** `GET`
- **Endpoint:** `/invoices?limit=10&offset=0`
- **Authentication:** Bearer Token (Required)
- **Pagination:** Supports limit & offset

### 3. Get Invoice by ID
- **Method:** `GET`
- **Endpoint:** `/invoices/:id`
- **Authentication:** Bearer Token (Required)

### 4. Update Invoice
- **Method:** `PUT`
- **Endpoint:** `/invoices/:id`
- **Authentication:** Bearer Token + Admin Role (Required)

### 5. Delete Invoice
- **Method:** `DELETE`
- **Endpoint:** `/invoices/:id`
- **Authentication:** Bearer Token + Admin Role (Required)

## Invoice Status Management

### 1. Update Invoice Status
- **Method:** `PATCH`
- **Endpoint:** `/invoices/:id/status`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Description:** Update invoice status (triggers inventory sync/deduction)
- **Request Body:**
```json
{
  "status": "sent"
}
```
- **Status Values:** draft, sent, partial, paid, overdue, void

### 2. Get Invoices by Status
- **Method:** `GET`
- **Endpoint:** `/invoices/status/:status`
- **Authentication:** Bearer Token (Required)

## Invoice Filtering

### 1. Get Invoices by Customer
- **Method:** `GET`
- **Endpoint:** `/customers/:customerId/invoices`
- **Authentication:** Bearer Token (Required)

---

# Manufacturing & Bill of Materials

## Item Groups (BOM)

### 1. Create Item Group
- **Method:** `POST`
- **Endpoint:** `/item-groups`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Description:** Create BOM (Bill of Materials) with component items
- **Request Body:**
```json
{
  "name": "300ml Water Bottle Set",
  "description": "Complete 300ml water bottle with packaging",
  "is_active": true,
  "components": [
    {
      "item_id": "item_bottle_001",
      "variant_id": 1,
      "quantity": 100
    },
    {
      "item_id": "item_cap_001",
      "variant_id": 1,
      "quantity": 100
    }
  ]
}
```
- **Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "id": "ig_001",
    "name": "300ml Water Bottle Set",
    "description": "Complete 300ml water bottle with packaging",
    "is_active": true,
    "components": [
      {
        "id": 1,
        "item_group_id": "ig_001",
        "item_id": "item_bottle_001",
        "variant_id": 1,
        "quantity": 100
      },
      {
        "id": 2,
        "item_group_id": "ig_001",
        "item_id": "item_cap_001",
        "variant_id": 1,
        "quantity": 100
      }
    ],
    "created_at": "2026-02-15T10:00:00Z",
    "updated_at": "2026-02-15T10:00:00Z"
  },
  "message": "Item Group created successfully"
}
```

### 2. Get All Item Groups
- **Method:** `GET`
- **Endpoint:** `/item-groups?limit=10&offset=0`
- **Authentication:** Bearer Token (Required)
- **Pagination:** Supports limit & offset

### 3. Get Item Group by ID
- **Method:** `GET`
- **Endpoint:** `/item-groups/:id`
- **Authentication:** Bearer Token (Required)

### 4. Update Item Group
- **Method:** `PUT`
- **Endpoint:** `/item-groups/:id`
- **Authentication:** Bearer Token + Admin Role (Required)

### 5. Delete Item Group
- **Method:** `DELETE`
- **Endpoint:** `/item-groups/:id`
- **Authentication:** Bearer Token + Admin Role (Required)

---

## Production Orders (Manufacturing)

### 1. Create Production Order
- **Method:** `POST`
- **Endpoint:** `/production-orders`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Description:** Create manufacturing order from BOM
- **Request Body:**
```json
{
  "item_group_id": "ig_001",
  "quantity_to_manufacture": 100,
  "planned_start_date": "2026-02-15",
  "planned_end_date": "2026-02-20",
  "notes": "Batch 01 - Regular production"
}
```
- **Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "id": "po_mfg_001",
    "production_order_no": "PO-MFG-001",
    "item_group_id": "ig_001",
    "quantity_to_manufacture": 100,
    "quantity_manufactured": 0,
    "status": "planned",
    "planned_start_date": "2026-02-15T00:00:00Z",
    "planned_end_date": "2026-02-20T00:00:00Z",
    "actual_start_date": null,
    "actual_end_date": null,
    "manufactured_date": null,
    "inventory_synced": false,
    "production_order_items": [
      {
        "id": 1,
        "production_order_id": "po_mfg_001",
        "item_group_component_id": 1,
        "quantity_required": 10000,
        "quantity_consumed": 0,
        "inventory_synced": false
      },
      {
        "id": 2,
        "production_order_id": "po_mfg_001",
        "item_group_component_id": 2,
        "quantity_required": 10000,
        "quantity_consumed": 0,
        "inventory_synced": false
      }
    ],
    "notes": "Batch 01 - Regular production",
    "created_at": "2026-02-15T10:00:00Z",
    "updated_at": "2026-02-15T10:00:00Z"
  },
  "message": "Production Order created successfully"
}
```

### 2. Get All Production Orders
- **Method:** `GET`
- **Endpoint:** `/production-orders?limit=10&offset=0`
- **Authentication:** Bearer Token (Required)
- **Pagination:** Supports limit & offset

### 3. Get Production Order by ID
- **Method:** `GET`
- **Endpoint:** `/production-orders/:id`
- **Authentication:** Bearer Token (Required)

### 4. Update Production Order
- **Method:** `PUT`
- **Endpoint:** `/production-orders/:id`
- **Authentication:** Bearer Token + Admin Role (Required)

### 5. Delete Production Order
- **Method:** `DELETE`
- **Endpoint:** `/production-orders/:id`
- **Authentication:** Bearer Token + Admin Role (Required)

## Production Order Status Management

### 1. Update Production Order Status
- **Method:** `PATCH`
- **Endpoint:** `/production-orders/:id/status`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Description:** Update order status (consumes components, creates finished goods)
- **Request Body:**
```json
{
  "status": "completed"
}
```
- **Status Values:** planned, in_progress, completed, cancelled

### 2. Get Production Orders by Status
- **Method:** `GET`
- **Endpoint:** `/production-orders/status/:status`
- **Authentication:** Bearer Token (Required)

---

# Inventory Management

## Real-time Inventory Tracking

### 1. Get Inventory Balance
- **Method:** `GET`
- **Endpoint:** `/api/inventory/balance/:item_id?variant_id=1`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Description:** Get current stock with reservations and in-transit
- **Response:**
```json
{
  "success": true,
  "data": {
    "item_id": "item_bottle_001",
    "variant_id": 1,
    "current_quantity": 500,
    "reserved_quantity": 100,
    "in_transit_quantity": 50,
    "available_quantity": 350,
    "reorder_level": 50,
    "average_rate": 30.5,
    "last_inventory_sync_at": "2026-02-16T10:00:00Z"
  }
}
```

### 2. Get All Inventory Balances
- **Method:** `GET`
- **Endpoint:** `/api/inventory/balance?limit=10&offset=0`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Pagination:** Supports limit & offset

## Stock Movements & Journal

### 1. Get Stock Movements
- **Method:** `GET`
- **Endpoint:** `/api/inventory/movements?item_id=item_001&date_from=2026-02-01&date_to=2026-02-16`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Description:** Get audit trail of all inventory movements
- **Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "item_id": "item_bottle_001",
      "variant_id": 1,
      "movement_type": "purchase_received",
      "quantity": 500,
      "reference_no": "PO-2026-001",
      "status": "completed",
      "notes": "Received from vendor",
      "created_at": "2026-02-15T10:00:00Z",
      "created_by": "user_123"
    }
  ]
}
```

### 2. Get Inventory Journal
- **Method:** `GET`
- **Endpoint:** `/api/inventory/journal?item_id=item_001&limit=10&offset=0`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Description:** Get detailed ledger entries

### 3. Get Inventory Aggregation
- **Method:** `GET`
- **Endpoint:** `/api/inventory/aggregation/:item_id?variant_id=1`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Description:** Get summary of purchases, manufacturing, and sales
- **Response:**
```json
{
  "success": true,
  "data": {
    "item_id": "item_bottle_001",
    "variant_id": 1,
    "total_purchased": 500,
    "total_manufactured": 0,
    "total_consumed_in_mfg": 100,
    "total_sold": 250,
    "average_rate": 1.5
  }
}
```

### 4. Get Supply Chain Summary
- **Method:** `GET`
- **Endpoint:** `/api/inventory/supply-chain/:item_id?variant_id=1`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Description:** Complete supply chain view (purchase → manufacturing → sales)
- **Response:**
```json
{
  "success": true,
  "data": {
    "item_id": "item_bottle_001",
    "variant_id": 1,
    "opening_stock": 50,
    "total_po_quantity": 500,
    "total_mfg_qty": 0,
    "total_consumed_in_mfg": 100,
    "total_so_quantity": 250,
    "total_invoiced_qty": 200,
    "current_qty": 150
  }
}
```

---

# Fulfillment & Shipping

## Packages

### 1. Create Package
- **Method:** `POST`
- **Endpoint:** `/packages`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Description:** Create package for shipment
- **Request Body:**
```json
{
  "sales_order_id": "SO-123456",
  "customer_id": 1,
  "package_items": [
    {
      "item_id": "ITEM-001",
      "quantity": 10,
      "item_description": "Laptops"
    }
  ],
  "package_weight": 25,
  "package_dimensions": "50x40x30"
}
```

### 2. Get All Packages
- **Method:** `GET`
- **Endpoint:** `/packages?limit=10&offset=0`
- **Authentication:** Bearer Token (Required)
- **Pagination:** Supports limit & offset

### 3. Get Package by ID
- **Method:** `GET`
- **Endpoint:** `/packages/:id`
- **Authentication:** Bearer Token (Required)

### 4. Update Package
- **Method:** `PUT`
- **Endpoint:** `/packages/:id`
- **Authentication:** Bearer Token + Admin Role (Required)

### 5. Delete Package
- **Method:** `DELETE`
- **Endpoint:** `/packages/:id`
- **Authentication:** Bearer Token + Admin Role (Required)

## Package Status Management

### 1. Update Package Status
- **Method:** `PATCH`
- **Endpoint:** `/packages/:id/status`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Request Body:**
```json
{
  "status": "packed"
}
```
- **Status Values:** created, packed, shipped, delivered, cancelled

### 2. Get Packages by Status
- **Method:** `GET`
- **Endpoint:** `/packages/status/:status`
- **Authentication:** Bearer Token (Required)

## Package Filtering

### 1. Get Packages by Customer
- **Method:** `GET`
- **Endpoint:** `/packages/customer/:customer_id`
- **Authentication:** Bearer Token (Required)

### 2. Get Packages by Sales Order
- **Method:** `GET`
- **Endpoint:** `/packages/sales-order/:sales_order_id`
- **Authentication:** Bearer Token (Required)

---

## Shipments

### 1. Create Shipment
- **Method:** `POST`
- **Endpoint:** `/shipments`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Description:** Create shipment for delivery
- **Request Body:**
```json
{
  "package_id": "PKG-123456",
  "sales_order_id": "SO-123456",
  "customer_id": 1,
  "shipment_date": "2026-02-07",
  "expected_delivery": "2026-02-14",
  "carrier": "FedEx",
  "tracking_number": "FEDEX123456",
  "shipping_address": "123 Main Street, Mumbai"
}
```

### 2. Get All Shipments
- **Method:** `GET`
- **Endpoint:** `/shipments?limit=10&offset=0`
- **Authentication:** Bearer Token (Required)
- **Pagination:** Supports limit & offset

### 3. Get Shipment by ID
- **Method:** `GET`
- **Endpoint:** `/shipments/:id`
- **Authentication:** Bearer Token (Required)

### 4. Update Shipment
- **Method:** `PUT`
- **Endpoint:** `/shipments/:id`
- **Authentication:** Bearer Token + Admin Role (Required)

### 5. Delete Shipment
- **Method:** `DELETE`
- **Endpoint:** `/shipments/:id`
- **Authentication:** Bearer Token + Admin Role (Required)

## Shipment Status Management

### 1. Update Shipment Status
- **Method:** `PATCH`
- **Endpoint:** `/shipments/:id/status`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Request Body:**
```json
{
  "status": "shipped"
}
```
- **Status Values:** created, shipped, in_transit, delivered, cancelled

### 2. Get Shipments by Status
- **Method:** `GET`
- **Endpoint:** `/shipments/status/:status`
- **Authentication:** Bearer Token (Required)

## Shipment Filtering

### 1. Get Shipments by Customer
- **Method:** `GET`
- **Endpoint:** `/shipments/customer/:customer_id`
- **Authentication:** Bearer Token (Required)

### 2. Get Shipments by Package
- **Method:** `GET`
- **Endpoint:** `/shipments/package/:package_id`
- **Authentication:** Bearer Token (Required)

### 3. Get Shipments by Sales Order
- **Method:** `GET`
- **Endpoint:** `/shipments/sales-order/:sales_order_id`
- **Authentication:** Bearer Token (Required)

---

# Financial Management

## Payments (Invoice & Bill Reconciliation)

### 1. Create Payment (Invoice)
- **Method:** `POST`
- **Endpoint:** `/payments`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Description:** Record payment against invoice
- **Request Body:**
```json
{
  "invoice_id": "INV-123456",
  "amount": 5000,
  "payment_date": "2026-02-07",
  "payment_method": "bank_transfer",
  "reference_number": "TXN-12345",
  "notes": "Payment received"
}
```
- **Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "invoice_id": "INV-123456",
    "amount": 5000,
    "payment_date": "2026-02-07",
    "status": "completed",
    "created_at": "2026-02-07T10:00:00Z"
  }
}
```

### 2. Get Payment by ID
- **Method:** `GET`
- **Endpoint:** `/payments/:id`
- **Authentication:** Bearer Token (Required)

### 3. Delete Payment
- **Method:** `DELETE`
- **Endpoint:** `/payments/:id`
- **Authentication:** Bearer Token + Admin Role (Required)

### 4. Get Payments by Invoice
- **Method:** `GET`
- **Endpoint:** `/invoices/:invoiceId/payments`
- **Authentication:** Bearer Token (Required)

---

## Salespersons

### 1. Create Salesperson
- **Method:** `POST`
- **Endpoint:** `/salespersons`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Request Body:**
```json
{
  "name": "Raj Kumar",
  "email": "raj@example.com",
  "phone": "9876543210",
  "department": "Sales",
  "target_revenue": 1000000
}
```

### 2. Get All Salespersons
- **Method:** `GET`
- **Endpoint:** `/salespersons?limit=10&offset=0`
- **Authentication:** Bearer Token (Required)
- **Pagination:** Supports limit & offset

### 3. Get Salesperson by ID
- **Method:** `GET`
- **Endpoint:** `/salespersons/:id`
- **Authentication:** Bearer Token (Required)

### 4. Update Salesperson
- **Method:** `PUT`
- **Endpoint:** `/salespersons/:id`
- **Authentication:** Bearer Token + Admin Role (Required)

### 5. Delete Salesperson
- **Method:** `DELETE`
- **Endpoint:** `/salespersons/:id`
- **Authentication:** Bearer Token + SuperAdmin Role (Required)

---

# User & Authentication

## Authentication

### 1. Register with Email
- **Method:** `POST`
- **Endpoint:** `/auth/register/email`
- **Authentication:** None (Public)
- **Request Body:**
```json
{
  "email": "user@example.com",
  "password": "SecurePassword123!",
  "name": "John Doe"
}
```
- **Response:**
```json
{
  "success": true,
  "data": {
    "id": "user_123",
    "email": "user@example.com",
    "name": "John Doe",
    "created_at": "2026-02-07T10:00:00Z"
  }
}
```

### 2. Register with Phone
- **Method:** `POST`
- **Endpoint:** `/auth/register/phone`
- **Authentication:** None (Public)
- **Request Body:**
```json
{
  "phone": "+919876543210",
  "password": "SecurePassword123!",
  "name": "John Doe"
}
```

### 3. Register with Google
- **Method:** `POST`
- **Endpoint:** `/auth/register/google`
- **Authentication:** None (Public)
- **Request Body:**
```json
{
  "google_token": "eyJhbGciOiJSUzI1NiIs..."
}
```

### 4. Login with Email
- **Method:** `POST`
- **Endpoint:** `/auth/login/email`
- **Authentication:** None (Public)
- **Request Body:**
```json
{
  "email": "user@example.com",
  "password": "SecurePassword123!"
}
```
- **Response:**
```json
{
  "success": true,
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5c...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5c...",
    "user": {
      "id": "user_123",
      "email": "user@example.com"
    }
  }
}
```

### 5. Login with Phone
- **Method:** `POST`
- **Endpoint:** `/auth/login/phone`
- **Authentication:** None (Public)
- **Request Body:**
```json
{
  "phone": "+919876543210",
  "password": "SecurePassword123!"
}
```

### 6. Login with Google
- **Method:** `POST`
- **Endpoint:** `/auth/login/google`
- **Authentication:** None (Public)
- **Request Body:**
```json
{
  "google_token": "eyJhbGciOiJSUzI1NiIs..."
}
```

### 7. Login with Apple
- **Method:** `POST`
- **Endpoint:** `/auth/login/apple`
- **Authentication:** None (Public)
- **Request Body:**
```json
{
  "apple_token": "eyJhbGciOiJSUzI1NiIs..."
}
```

### 8. Login with Password
- **Method:** `POST`
- **Endpoint:** `/auth/login/password`
- **Authentication:** None (Public)
- **Request Body:**
```json
{
  "email": "user@example.com",
  "password": "SecurePassword123!"
}
```

### 9. Verify OTP
- **Method:** `POST`
- **Endpoint:** `/auth/verify-otp`
- **Authentication:** None (Public)
- **Request Body:**
```json
{
  "email": "user@example.com",
  "otp": "123456"
}
```

### 10. Change Password
- **Method:** `POST`
- **Endpoint:** `/auth/change-password`
- **Authentication:** Bearer Token (Required)
- **Request Body:**
```json
{
  "old_password": "OldPassword123!",
  "new_password": "NewPassword123!"
}
```

### 11. Reset Password (Admin)
- **Method:** `POST`
- **Endpoint:** `/auth/reset-password`
- **Authentication:** Bearer Token + SuperAdmin Role (Required)
- **Request Body:**
```json
{
  "user_id": "user_123",
  "new_password": "NewPassword123!"
}
```

### 12. Refresh Token
- **Method:** `POST`
- **Endpoint:** `/auth/refresh-token`
- **Authentication:** None (Public)
- **Request Body:**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5c..."
}
```

### 13. Validate Token
- **Method:** `POST`
- **Endpoint:** `/auth/validate-token`
- **Authentication:** None (Public)
- **Request Body:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5c..."
}
```

---

## User Management

### 1. Create User (SuperAdmin)
- **Method:** `POST`
- **Endpoint:** `/users`
- **Authentication:** Bearer Token + SuperAdmin Role (Required)
- **Request Body:**
```json
{
  "email": "user@example.com",
  "phone": "+919876543210",
  "name": "John Doe",
  "password": "SecurePassword123!",
  "user_type": "admin",
  "role": "admin"
}
```

### 2. Get All Users
- **Method:** `GET`
- **Endpoint:** `/users?limit=10&offset=0`
- **Authentication:** Bearer Token + SuperAdmin Role (Required)

### 3. Get User by ID
- **Method:** `GET`
- **Endpoint:** `/users/:id`
- **Authentication:** Bearer Token (Required)

### 4. Update User
- **Method:** `PUT`
- **Endpoint:** `/users/:id`
- **Authentication:** Bearer Token (Required)

### 5. Delete User
- **Method:** `DELETE`
- **Endpoint:** `/users/:id`
- **Authentication:** Bearer Token + SuperAdmin Role (Required)

---

## Admin Management

### 1. Create Admin
- **Method:** `POST`
- **Endpoint:** `/admin/create`
- **Authentication:** Bearer Token + SuperAdmin Role (Required)
- **Request Body:**
```json
{
  "email": "admin@example.com",
  "name": "Admin User",
  "password": "SecurePassword123!"
}
```

### 2. Get All Admins
- **Method:** `GET`
- **Endpoint:** `/admin/list?limit=10&offset=0`
- **Authentication:** Bearer Token + SuperAdmin Role (Required)

### 3. Get Admin by ID
- **Method:** `GET`
- **Endpoint:** `/admin/:id`
- **Authentication:** Bearer Token + SuperAdmin Role (Required)

### 4. Update Admin
- **Method:** `PUT`
- **Endpoint:** `/admin/:id`
- **Authentication:** Bearer Token + SuperAdmin Role (Required)

### 5. Delete Admin
- **Method:** `DELETE`
- **Endpoint:** `/admin/:id`
- **Authentication:** Bearer Token + SuperAdmin Role (Required)

---

# Utilities

## Helper/Lookup Routes

### 1. Get Business Types
- **Method:** `GET`
- **Endpoint:** `/helpers/business-types`
- **Authentication:** None (Public)
- **Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "Manufacturing",
      "description": "Manufacturing business"
    }
  ]
}
```

### 2. Get Countries
- **Method:** `GET`
- **Endpoint:** `/helpers/countries`
- **Authentication:** None (Public)

### 3. Get States by Country
- **Method:** `GET`
- **Endpoint:** `/helpers/countries/:country_id/states`
- **Authentication:** None (Public)

### 4. Get Tax Types
- **Method:** `GET`
- **Endpoint:** `/helpers/tax-types`
- **Authentication:** None (Public)

---

## Forward Auth Routes (Traefik Middleware)

### 1. Forward Auth
- **Method:** `GET`
- **Endpoint:** `/forward-auth/`
- **Authentication:** Query parameter `token` (Required)
- **Response:** Returns 200 if authenticated, 401 if not

### 2. Product Auth
- **Method:** `GET`
- **Endpoint:** `/forward-auth/product`
- **Authentication:** Query parameter `token` (Required)

### 3. Customer Auth
- **Method:** `GET`
- **Endpoint:** `/forward-auth/customer`
- **Authentication:** Query parameter `token` (Required)

---

## Support

### 1. Create Support Ticket
- **Method:** `POST`
- **Endpoint:** `/public/support`
- **Authentication:** None (Public)
- **Request Body:**
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "subject": "Technical Issue",
  "description": "I'm facing an issue with...",
  "priority": "high"
}
```
- **Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "ticket_number": "TKT-000001",
    "status": "open",
    "created_at": "2026-02-07T10:00:00Z"
  }
}
```

---

## Health Check

### 1. Health Check
- **Method:** `GET`
- **Endpoint:** `/health`
- **Authentication:** None (Public)
- **Response:**
```json
{
  "status": "ok",
  "service": "github.com/bbapp-org/auth-service",
  "version": "1.0.0"
}
```

---

# Common Response Formats

## Success Response
```json
{
  "success": true,
  "data": {
    // response data
  },
  "message": "Operation successful"
}
```

## Error Response
```json
{
  "success": false,
  "error": "Error message",
  "status_code": 400
}
```

## Paginated Response
```json
{
  "success": true,
  "data": [
    // array of items
  ],
  "total": 100,
  "limit": 10,
  "offset": 0
}
```

---

# Authentication Header

All protected endpoints require:
```
Authorization: Bearer <access_token>
```

---

# HTTP Status Codes

| Code | Meaning |
|------|---------|
| 200  | OK - Successful GET/PUT/PATCH request |
| 201  | Created - Successful POST request |
| 204  | No Content - Successful DELETE request |
| 400  | Bad Request - Invalid parameters or body |
| 401  | Unauthorized - Missing/invalid token |
| 403  | Forbidden - Insufficient permissions |
| 404  | Not Found - Resource not found |
| 409  | Conflict - Duplicate/conflicting resource |
| 500  | Internal Server Error - Server error |

---

# API Features

## Real-time Inventory Synchronization
- PurchaseOrder → Automatic inventory update on receipt (DeliveryDateActual)
- SalesOrder → Automatic inventory reservation on confirmation (InventoryReserved)
- Invoice → Automatic inventory deduction on creation (InventorySynced)
- ProductionOrder → Automatic component consumption & finished goods creation (ManufacturedDate)

## Audit Trail
- All transactions tracked with timestamps (CreatedAt, UpdatedAt)
- User tracking (CreatedBy, UpdatedBy)
- Complete movement history via StockMovement and InventoryJournal
- Reference tracking (ReferenceNo, ReferenceID, ReferenceType)

## Pagination
Most GET endpoints support:
- `limit` - Items per page (default: 10, max: 100)
- `offset` - Items to skip (default: 0)

## Filtering & Sorting
All list endpoints support filtering by:
- Status
- Date range
- Related entity (Customer, Vendor, etc.)

## Swagger Documentation

Interactive API documentation available at:
```
GET /docs/
```

---

**Last Updated:** February 16, 2026
**API Version:** 1.0.0
