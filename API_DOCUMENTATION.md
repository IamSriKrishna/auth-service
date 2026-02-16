# API Documentation

## Table of Contents
1. [Authentication](#authentication)
2. [Admin Management](#admin-management)
3. [Users & Partners](#users--partners)
4. [Vendors](#vendors)
5. [Customers](#customers)
6. [Items & Inventory](#items--inventory)
7. [Item Groups (BOM)](#item-groups-bom)
8. [Production Orders (Manufacturing)](#production-orders-manufacturing)
9. [Inventory Tracking](#inventory-tracking)
10. [Manufacturers & Brands](#manufacturers--brands)
11. [Companies](#companies)
12. [Invoices](#invoices)
13. [Bills](#bills)
14. [Payments](#payments)
15. [Tax Configuration](#tax-configuration)
16. [Salespersons](#salespersons)
17. [Purchase Orders](#purchase-orders)
18. [Sales Orders](#sales-orders)
19. [Packages](#packages)
20. [Shipments](#shipments)
21. [Helper/Lookup Routes](#helperlookup-routes)
22. [Forward Auth Routes](#forward-auth-routes)
23. [Support](#support)

---

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
  },
  "message": "User registered successfully"
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
- **Response:** Same as email registration

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
- **Response:** Returns user with JWT token and refresh token

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
    "user_id": "user_123",
    "email": "user@example.com",
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "token_expiry": "2026-02-07T11:00:00Z"
  },
  "message": "Login successful"
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
- **Response:** Same as email login

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
  "username": "john.doe",
  "password": "SecurePassword123!"
}
```

### 9. Validate Token
- **Method:** `POST`
- **Endpoint:** `/auth/validate-token`
- **Authentication:** None (Public)
- **Request Body:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```
- **Response:**
```json
{
  "success": true,
  "valid": true,
  "user_id": "user_123",
  "expires_at": "2026-02-07T11:00:00Z"
}
```

### 10. Refresh Token
- **Method:** `POST`
- **Endpoint:** `/auth/refresh-token`
- **Authentication:** Bearer Token (Required)
- **Request Body:**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```
- **Response:**
```json
{
  "success": true,
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_at": "2026-02-07T11:00:00Z"
  }
}
```

### 11. Get User Info
- **Method:** `GET`
- **Endpoint:** `/auth/user-info`
- **Authentication:** Bearer Token (Required)
- **Response:**
```json
{
  "success": true,
  "data": {
    "id": "user_123",
    "email": "user@example.com",
    "name": "John Doe",
    "role": "admin",
    "created_at": "2026-01-01T10:00:00Z"
  }
}
```

### 12. Logout
- **Method:** `POST`
- **Endpoint:** `/auth/logout`
- **Authentication:** Bearer Token (Required)
- **Response:**
```json
{
  "success": true,
  "message": "Logged out successfully"
}
```

---

## Admin Management

### 1. Create User
- **Method:** `POST`
- **Endpoint:** `/auth/admin/create-user`
- **Authentication:** Bearer Token + SuperAdmin Role (Required)
- **Request Body:**
```json
{
  "email": "newuser@example.com",
  "password": "SecurePassword123!",
  "first_name": "Jane",
  "last_name": "Doe",
  "role": "admin"
}
```
- **Response:**
```json
{
  "success": true,
  "data": {
    "id": "user_456",
    "email": "newuser@example.com",
    "first_name": "Jane",
    "last_name": "Doe",
    "role": "admin",
    "created_at": "2026-02-07T10:00:00Z"
  },
  "message": "User created successfully"
}
```

### 2. Get Users
- **Method:** `GET`
- **Endpoint:** `/auth/admin/users?limit=10&offset=0`
- **Authentication:** Bearer Token + SuperAdmin Role (Required)
- **Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": "user_123",
      "email": "user@example.com",
      "first_name": "John",
      "last_name": "Doe",
      "role": "admin",
      "status": "active",
      "created_at": "2026-01-01T10:00:00Z"
    }
  ],
  "total": 1
}
```

### 3. Get User by ID
- **Method:** `GET`
- **Endpoint:** `/auth/admin/users/:id`
- **Authentication:** Bearer Token + SuperAdmin Role (Required)
- **Response:** Single user object (same as above)

### 4. Update User
- **Method:** `PUT`
- **Endpoint:** `/auth/admin/users/:id`
- **Authentication:** Bearer Token + SuperAdmin Role (Required)
- **Request Body:**
```json
{
  "first_name": "Jane",
  "last_name": "Smith",
  "email": "jane.smith@example.com"
}
```

### 5. Delete User
- **Method:** `DELETE`
- **Endpoint:** `/auth/admin/users/:id`
- **Authentication:** Bearer Token + SuperAdmin Role (Required)
- **Response:**
```json
{
  "success": true,
  "message": "User deleted successfully"
}
```

### 6. Update User Status
- **Method:** `PUT`
- **Endpoint:** `/auth/admin/users/:id/status`
- **Authentication:** Bearer Token + SuperAdmin Role (Required)
- **Request Body:**
```json
{
  "status": "active"
}
```

### 7. Update User Role
- **Method:** `PUT`
- **Endpoint:** `/auth/admin/users/:id/role`
- **Authentication:** Bearer Token + SuperAdmin Role (Required)
- **Request Body:**
```json
{
  "role": "admin"
}
```

### 8. Reset Admin Password
- **Method:** `POST`
- **Endpoint:** `/auth/admin/reset-password`
- **Authentication:** Bearer Token + SuperAdmin Role (Required)
- **Request Body:**
```json
{
  "user_id": "user_123",
  "new_password": "NewSecurePassword123!"
}
```

### 9. Get Dashboard Stats
- **Method:** `GET`
- **Endpoint:** `/auth/admin/dashboard/stats`
- **Authentication:** Bearer Token + SuperAdmin Role (Required)
- **Response:**
```json
{
  "success": true,
  "data": {
    "total_users": 50,
    "total_orders": 150,
    "revenue": 500000,
    "pending_tasks": 25
  }
}
```

---

## Users & Partners

### 1. Create Partner
- **Method:** `POST`
- **Endpoint:** `/auth/manage/create-partner`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Request Body:**
```json
{
  "email": "partner@example.com",
  "password": "SecurePassword123!",
  "company_name": "ABC Corp",
  "contact_person": "John Smith"
}
```

### 2. Get Partners
- **Method:** `GET`
- **Endpoint:** `/auth/manage/partners?limit=10&offset=0`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Response:** List of partner objects

### 3. Reset Partner Password
- **Method:** `PATCH`
- **Endpoint:** `/partners/:partner_id/reset-password`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Request Body:**
```json
{
  "new_password": "NewSecurePassword123!"
}
```

---

## Vendors

### 1. Create Vendor
- **Method:** `POST`
- **Endpoint:** `/vendors`
- **Authentication:** Bearer Token + SuperAdmin Role (Required)
- **Request Body:**
```json
{
  "salutation": "Mr.",
  "first_name": "Rajesh",
  "last_name": "Kumar",
  "company_name": "Kumar Supplies",
  "display_name": "Kumar Supplies",
  "email_address": "rajesh@kumarsupplies.com",
  "work_phone": "9876543210",
  "work_phone_code": "+91",
  "mobile": "9876543210",
  "mobile_code": "+91",
  "vendor_language": "English",
  "gstin": "18AABCT1234H1Z0"
}
```
- **Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "display_name": "Kumar Supplies",
    "company_name": "Kumar Supplies",
    "email_address": "rajesh@kumarsupplies.com",
    "work_phone": "9876543210",
    "mobile": "9876543210",
    "created_at": "2026-02-07T10:00:00Z"
  },
  "message": "Vendor created successfully"
}
```

### 2. Get All Vendors
- **Method:** `GET`
- **Endpoint:** `/vendors?limit=10&offset=0`
- **Authentication:** None (Public)
- **Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "display_name": "Kumar Supplies",
      "company_name": "Kumar Supplies",
      "email_address": "rajesh@kumarsupplies.com"
    }
  ],
  "total": 1
}
```

### 3. Get Vendor by ID
- **Method:** `GET`
- **Endpoint:** `/vendors/:id`
- **Authentication:** None (Public)
- **Response:** Single vendor object

### 4. Update Vendor
- **Method:** `PUT`
- **Endpoint:** `/vendors/:id`
- **Authentication:** Bearer Token + SuperAdmin Role (Required)
- **Request Body:** Same fields as create (all optional)

### 5. Delete Vendor
- **Method:** `DELETE`
- **Endpoint:** `/vendors/:id`
- **Authentication:** Bearer Token + SuperAdmin Role (Required)

---

## Customers

### 1. Create Customer
- **Method:** `POST`
- **Endpoint:** `/customers`
- **Authentication:** Bearer Token + SuperAdmin Role (Required)
- **Request Body:**
```json
{
  "salutation": "Mr.",
  "first_name": "Amit",
  "last_name": "Singh",
  "company_name": "Singh Enterprises",
  "display_name": "Singh Enterprises",
  "email_address": "amit@singh.com",
  "work_phone": "9876543210",
  "mobile": "9876543210",
  "customer_language": "English"
}
```

### 2. Get All Customers
- **Method:** `GET`
- **Endpoint:** `/customers?limit=10&offset=0`
- **Authentication:** None (Public)
- **Response:** List of customers

### 3. Get Customer by ID
- **Method:** `GET`
- **Endpoint:** `/customers/:id`
- **Authentication:** None (Public)

### 4. Update Customer
- **Method:** `PUT`
- **Endpoint:** `/customers/:id`
- **Authentication:** Bearer Token + SuperAdmin Role (Required)

### 5. Delete Customer
- **Method:** `DELETE`
- **Endpoint:** `/customers/:id`
- **Authentication:** Bearer Token + SuperAdmin Role (Required)

### 6. Get Customer Invoices
- **Method:** `GET`
- **Endpoint:** `/customers/:customerId/invoices`
- **Authentication:** Bearer Token (Required)

---

## Items & Inventory

### 1. Create Item
- **Method:** `POST`
- **Endpoint:** `/items`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Request Body:**
```json
{
  "name": "Laptop",
  "type": "product",
  "brand": "Dell",
  "manufacturer_id": 1,
  "item_details": {
    "sku": "LAP-001",
    "hsn_code": "8471.30",
    "description": "Dell Laptop 15 inch"
  },
  "sales_info": {
    "sales_account": "Sales",
    "sales_tax": 18
  },
  "purchase_info": {
    "purchase_account": "Purchases",
    "purchase_tax": 5
  }
}
```

### 2. Get All Items
- **Method:** `GET`
- **Endpoint:** `/items?limit=10&offset=0`
- **Authentication:** None (Public)

### 3. Get Item by ID
- **Method:** `GET`
- **Endpoint:** `/items/:id`
- **Authentication:** None (Public)

### 4. Update Item
- **Method:** `PUT`
- **Endpoint:** `/items/:id`
- **Authentication:** Bearer Token + Admin Role (Required)

### 5. Delete Item
- **Method:** `DELETE`
- **Endpoint:** `/items/:id`
- **Authentication:** Bearer Token + SuperAdmin Role (Required)

### 6. Update Opening Stock
- **Method:** `PUT`
- **Endpoint:** `/items/:id/opening-stock`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Request Body:**
```json
{
  "quantity": 100,
  "warehouse_id": 1
}
```

### 7. Get Opening Stock
- **Method:** `GET`
- **Endpoint:** `/items/:id/opening-stock`
- **Authentication:** Bearer Token + Admin Role (Required)

### 8. Update Variants Opening Stock
- **Method:** `PUT`
- **Endpoint:** `/items/:id/variants/opening-stock`
- **Authentication:** Bearer Token + Admin Role (Required)

### 9. Get Variants Opening Stock
- **Method:** `GET`
- **Endpoint:** `/items/:id/variants/opening-stock`
- **Authentication:** Bearer Token + Admin Role (Required)

### 10. Get Stock Summary
- **Method:** `GET`
- **Endpoint:** `/items/:id/stock-summary`
- **Authentication:** Bearer Token + Admin Role (Required)

---

## Item Groups (BOM)

Item Groups represent Bill of Materials (BOM) - combinations of items that form finished products.
**Example:** "300ml Water Bottle" = 1 × Bottle (300ml) + 1 × Cap (20mm)

### 1. Create Item Group
- **Method:** `POST`
- **Endpoint:** `/api/item-groups`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Request Body:**
```json
{
  "name": "300ml Water Bottle",
  "description": "300ml plastic water bottle with cap",
  "is_active": true,
  "components": [
    {
      "item_id": "item_bottle_001",
      "variant_id": 1,
      "quantity": 1,
      "variant_details": {
        "capacity": "300ml",
        "material": "plastic"
      }
    },
    {
      "item_id": "item_cap_001",
      "variant_id": 2,
      "quantity": 1,
      "variant_details": {
        "size": "20mm"
      }
    }
  ]
}
```
- **Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "id": "grp_300ml_bottle",
    "name": "300ml Water Bottle",
    "description": "300ml plastic water bottle with cap",
    "is_active": true,
    "components": [
      {
        "id": 1,
        "item_group_id": "grp_300ml_bottle",
        "item_id": "item_bottle_001",
        "variant_id": 1,
        "quantity": 1,
        "variant_details": {
          "capacity": "300ml",
          "material": "plastic"
        },
        "created_at": "2026-02-16T10:00:00Z",
        "updated_at": "2026-02-16T10:00:00Z"
      }
    ],
    "created_at": "2026-02-16T10:00:00Z",
    "updated_at": "2026-02-16T10:00:00Z"
  },
  "message": "Item Group created successfully"
}
```

### 2. List Item Groups
- **Method:** `GET`
- **Endpoint:** `/api/item-groups?page=1&page_size=10&search=water&is_active=true`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Query Parameters:**
  - `page` - Page number (default: 1)
  - `page_size` - Items per page (default: 10, max: 100)
  - `search` - Search by name
  - `is_active` - Filter by active status (true/false)
  - `sort` - Sort field (name, created_at)

- **Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "data": [
      {
        "id": "grp_300ml_bottle",
        "name": "300ml Water Bottle",
        "description": "300ml plastic water bottle with cap",
        "is_active": true,
        "components": [...],
        "created_at": "2026-02-16T10:00:00Z",
        "updated_at": "2026-02-16T10:00:00Z"
      }
    ],
    "total": 5,
    "page": 1,
    "page_size": 10,
    "total_page": 1
  }
}
```

### 3. Get Item Group by ID
- **Method:** `GET`
- **Endpoint:** `/api/item-groups/:id`
- **Authentication:** Bearer Token + Admin Role (Required)

- **Response (200 OK):** Same as Create response

### 4. Update Item Group
- **Method:** `PUT`
- **Endpoint:** `/api/item-groups/:id`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Request Body:**
```json
{
  "name": "300ml Water Bottle Premium",
  "description": "Updated description",
  "is_active": true,
  "components": [...]
}
```

### 5. Delete Item Group
- **Method:** `DELETE`
- **Endpoint:** `/api/item-groups/:id`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Response (204 No Content)**

---

## Production Orders (Manufacturing)

Production Orders track the manufacturing of Item Groups from component inventory.

### 1. Create Production Order
- **Method:** `POST`
- **Endpoint:** `/api/production-orders`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Request Body:**
```json
{
  "item_group_id": "grp_300ml_bottle",
  "quantity_to_manufacture": 100,
  "planned_start_date": "2026-02-20T00:00:00Z",
  "planned_end_date": "2026-02-25T00:00:00Z",
  "notes": "Standard production run"
}
```
- **Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "id": "po_mfg_001",
    "production_order_no": "PO-MFG-001",
    "item_group_id": "grp_300ml_bottle",
    "item_group": {
      "id": "grp_300ml_bottle",
      "name": "300ml Water Bottle",
      "components": [...]
    },
    "quantity_to_manufacture": 100,
    "quantity_manufactured": 0,
    "status": "planned",
    "planned_start_date": "2026-02-20T00:00:00Z",
    "planned_end_date": "2026-02-25T00:00:00Z",
    "actual_start_date": null,
    "actual_end_date": null,
    "notes": "Standard production run",
    "production_order_items": [
      {
        "id": 1,
        "production_order_id": "po_mfg_001",
        "item_group_component_id": 1,
        "quantity_required": 100,
        "quantity_consumed": 0,
        "created_at": "2026-02-16T10:00:00Z",
        "updated_at": "2026-02-16T10:00:00Z"
      }
    ],
    "created_at": "2026-02-16T10:00:00Z",
    "updated_at": "2026-02-16T10:00:00Z",
    "created_by": "user_123",
    "updated_by": "user_123"
  },
  "message": "Production Order created successfully"
}
```

### 2. List Production Orders
- **Method:** `GET`
- **Endpoint:** `/api/production-orders?page=1&page_size=10&status=in_progress&item_group_id=grp_300ml_bottle`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Query Parameters:**
  - `page` - Page number (default: 1)
  - `page_size` - Items per page (default: 10)
  - `status` - Filter by status (planned, in_progress, completed, cancelled)
  - `item_group_id` - Filter by ItemGroup
  - `start_date` - Filter from date
  - `end_date` - Filter to date
  - `sort` - Sort field

- **Response (200 OK):** List of Production Orders

### 3. Get Production Order by ID
- **Method:** `GET`
- **Endpoint:** `/api/production-orders/:id`
- **Authentication:** Bearer Token + Admin Role (Required)

### 4. Update Production Order Status
- **Method:** `PUT`
- **Endpoint:** `/api/production-orders/:id/status`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Request Body:**
```json
{
  "status": "in_progress",
  "quantity_manufactured": 0,
  "actual_end_date": null
}
```
- **Status Values:** `planned`, `in_progress`, `completed`, `cancelled`

### 5. Mark Production Order In Progress
- **Method:** `POST`
- **Endpoint:** `/api/production-orders/:id/start`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Response:** Returns updated Production Order with status = "in_progress"

### 6. Complete Production Order
- **Method:** `POST`
- **Endpoint:** `/api/production-orders/:id/complete`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Request Body:**
```json
{
  "quantity_manufactured": 100
}
```
- **Response:** Returns updated Production Order with:
  - `status` = "completed"
  - `quantity_manufactured` = 100
  - Components marked as "consumed" in inventory
  - New product inventory created

### 7. Get Required Components
- **Method:** `GET`
- **Endpoint:** `/api/production-orders/:id/components`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Response:**
```json
{
  "success": true,
  "data": {
    "production_order_id": "po_mfg_001",
    "components": [
      {
        "item_id": "item_bottle_001",
        "variant_id": 1,
        "quantity_required": 100,
        "quantity_available": 150,
        "name": "Bottle (300ml)"
      },
      {
        "item_id": "item_cap_001",
        "variant_id": 2,
        "quantity_required": 100,
        "quantity_available": 200,
        "name": "Cap (20mm)"
      }
    ]
  }
}
```

---

## Inventory Tracking

Comprehensive inventory tracking across purchases, manufacturing, and sales.

### 1. Get Inventory Balance
- **Method:** `GET`
- **Endpoint:** `/api/inventory/balance/:item_id?variant_id=1`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "item_id": "item_bottle_001",
    "variant_id": 1,
    "item": {
      "id": "item_bottle_001",
      "name": "Plastic Bottle",
      "type": "goods"
    },
    "variant": {
      "id": 1,
      "sku": "BOTTLE-300ML",
      "selling_price": 2.5,
      "cost_price": 1.5,
      "stock_quantity": 150
    },
    "current_quantity": 150,
    "reserved_quantity": 50,
    "available_quantity": 100,
    "last_received_date": "2026-02-15T10:00:00Z",
    "last_consumed_date": "2026-02-14T14:30:00Z",
    "last_sold_date": "2026-02-16T09:00:00Z",
    "updated_at": "2026-02-16T10:00:00Z"
  }
}
```

### 2. Get Inventory Aggregation
- **Method:** `GET`
- **Endpoint:** `/api/inventory/aggregation/:item_id?variant_id=1`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Response (200 OK):**
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
    "calculated_at": "2026-02-16T10:00:00Z",
    "updated_at": "2026-02-16T10:00:00Z"
  }
}
```

### 3. Reserve Inventory
- **Method:** `PUT`
- **Endpoint:** `/api/inventory/balance/:item_id/reserve`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Request Body:**
```json
{
  "variant_id": 1,
  "quantity": 50
}
```
- **Response:** Updated InventoryBalance with reserved_quantity increased

### 4. Release Reservation
- **Method:** `PUT`
- **Endpoint:** `/api/inventory/balance/:item_id/release`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Request Body:**
```json
{
  "variant_id": 1,
  "quantity": 50
}
```

### 5. Get Inventory Journal (Audit Trail)
- **Method:** `GET`
- **Endpoint:** `/api/inventory/journal/:item_id?variant_id=1&start_date=2026-02-01&end_date=2026-02-16&page=1&page_size=20`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Query Parameters:**
  - `variant_id` - Filter by variant
  - `start_date` - Filter from date (ISO 8601)
  - `end_date` - Filter to date (ISO 8601)
  - `page` - Page number
  - `page_size` - Items per page

- **Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "data": [
      {
        "id": 1,
        "item_id": "item_bottle_001",
        "variant_id": 1,
        "transaction_type": "purchase",
        "quantity": 500,
        "reference_type": "PurchaseOrder",
        "reference_id": "po_001",
        "reference_no": "PO-2026-001",
        "notes": "Received from vendor",
        "created_at": "2026-02-15T10:00:00Z",
        "created_by": "user_123"
      },
      {
        "id": 2,
        "item_id": "item_bottle_001",
        "variant_id": 1,
        "transaction_type": "consume",
        "quantity": -100,
        "reference_type": "ProductionOrder",
        "reference_id": "po_mfg_001",
        "reference_no": "PO-MFG-001",
        "notes": "Consumed in manufacturing",
        "created_at": "2026-02-16T10:00:00Z",
        "created_by": "user_456"
      }
    ],
    "total": 10,
    "start_date": "2026-02-01T00:00:00Z",
    "end_date": "2026-02-16T23:59:59Z"
  }
}
```

### 6. Get Supply Chain Summary
- **Method:** `GET`
- **Endpoint:** `/api/inventory/supply-chain/:item_id?variant_id=1`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "item_id": "item_bottle_001",
    "variant_id": 1,
    "opening_stock": 50,
    "total_po_quantity": 500,
    "total_po_amount": 750.0,
    "avg_purchase_rate": 1.5,
    "total_prod_qty": 0,
    "total_mfg_qty": 0,
    "total_consumed_in_mfg": 100,
    "total_so_quantity": 250,
    "total_so_amount": 625.0,
    "avg_sales_rate": 2.5,
    "total_invoiced_qty": 200,
    "current_qty": 150,
    "updated_at": "2026-02-16T10:00:00Z"
  }
}
```

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

## Companies

### 1. Complete Company Setup
- **Method:** `POST`
- **Endpoint:** `/companies/setup`
- **Authentication:** Bearer Token (Required)
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

### 2. Create Company
- **Method:** `POST`
- **Endpoint:** `/companies`
- **Authentication:** Bearer Token (Required)

### 3. Get All Companies
- **Method:** `GET`
- **Endpoint:** `/companies?limit=10&offset=0`
- **Authentication:** Bearer Token (Required)

### 4. Get Company by ID
- **Method:** `GET`
- **Endpoint:** `/companies/:id`
- **Authentication:** Bearer Token (Required)

### 5. Update Company
- **Method:** `PUT`
- **Endpoint:** `/companies/:id`
- **Authentication:** Bearer Token (Required)

### 6. Delete Company
- **Method:** `DELETE`
- **Endpoint:** `/companies/:id`
- **Authentication:** Bearer Token + SuperAdmin Role (Required)

### 7. Upsert Company Contact
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

### 8. Get Company Contact
- **Method:** `GET`
- **Endpoint:** `/companies/:id/contact`
- **Authentication:** Bearer Token (Required)

### 9. Upsert Company Address
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

### 10. Get Company Address
- **Method:** `GET`
- **Endpoint:** `/companies/:id/address`
- **Authentication:** Bearer Token (Required)

### 11. Create Bank Detail
- **Method:** `POST`
- **Endpoint:** `/companies/:id/bank-details`
- **Authentication:** Bearer Token (Required)
- **Request Body:**
```json
{
  "bank_name": "HDFC Bank",
  "account_holder": "ABC Corporation",
  "account_number": "1234567890",
  "ifsc_code": "HDFC0001234"
}
```

### 12. Get Bank Details
- **Method:** `GET`
- **Endpoint:** `/companies/:id/bank-details`
- **Authentication:** Bearer Token (Required)

### 13. Update Bank Detail
- **Method:** `PUT`
- **Endpoint:** `/companies/bank-details/:id`
- **Authentication:** Bearer Token (Required)

### 14. Delete Bank Detail
- **Method:** `DELETE`
- **Endpoint:** `/companies/bank-details/:id`
- **Authentication:** Bearer Token (Required)

### 15. Upsert UPI Detail
- **Method:** `PUT`
- **Endpoint:** `/companies/:id/upi-details`
- **Request Body:**
```json
{
  "upi_id": "abc@hdfc"
}
```

### 16. Get UPI Detail
- **Method:** `GET`
- **Endpoint:** `/companies/:id/upi-details`

### 17. Upsert Invoice Settings
- **Method:** `PUT`
- **Endpoint:** `/companies/:id/invoice-settings`
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

### 18. Get Invoice Settings
- **Method:** `GET`
- **Endpoint:** `/companies/:id/invoice-settings`

### 19. Upsert Tax Settings
- **Method:** `PUT`
- **Endpoint:** `/companies/:id/tax-settings`
- **Request Body:**
```json
{
  "gst_enabled": true,
  "tax_type_id": 1
}
```

### 20. Get Tax Settings
- **Method:** `GET`
- **Endpoint:** `/companies/:id/tax-settings`

### 21. Upsert Regional Settings
- **Method:** `PUT`
- **Endpoint:** `/companies/:id/regional-settings`
- **Request Body:**
```json
{
  "default_currency": "INR",
  "timezone": "Asia/Kolkata",
  "date_format": "DD/MM/YYYY"
}
```

### 22. Get Regional Settings
- **Method:** `GET`
- **Endpoint:** `/companies/:id/regional-settings`

---

## Invoices

### 1. Create Invoice
- **Method:** `POST`
- **Endpoint:** `/invoices`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Request Body:**
```json
{
  "customer_id": 1,
  "invoice_date": "2026-02-07",
  "due_date": "2026-03-07",
  "terms": "net_30",
  "subject": "Invoice for services",
  "salesperson_id": 1,
  "line_items": [
    {
      "item_id": "ITEM-001",
      "quantity": 2,
      "rate": 5000,
      "description": "Consulting services"
    }
  ],
  "shipping_charges": 500,
  "tax_id": 1,
  "adjustment": 100,
  "customer_notes": "Thank you for your business",
  "terms_and_conditions": "Net 30 payment terms"
}
```
- **Response:**
```json
{
  "success": true,
  "data": {
    "id": "INV-123456",
    "invoice_number": "INV-001",
    "customer_id": 1,
    "invoice_date": "2026-02-07",
    "due_date": "2026-03-07",
    "sub_total": 10000,
    "tax_amount": 1800,
    "total": 12400,
    "status": "draft",
    "created_at": "2026-02-07T10:00:00Z"
  },
  "message": "Invoice created successfully"
}
```

### 2. Get All Invoices
- **Method:** `GET`
- **Endpoint:** `/invoices?limit=10&offset=0`
- **Authentication:** Bearer Token (Required)

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

### 6. Update Invoice Status
- **Method:** `PATCH`
- **Endpoint:** `/invoices/:id/status`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Request Body:**
```json
{
  "status": "sent"
}
```

### 7. Get Invoices by Status
- **Method:** `GET`
- **Endpoint:** `/invoices/status/:status`
- **Authentication:** Bearer Token (Required)
- **Status Values:** draft, sent, partial, paid, overdue, void

### 8. Get Invoices by Customer
- **Method:** `GET`
- **Endpoint:** `/customers/:customerId/invoices`
- **Authentication:** Bearer Token (Required)

### 9. Get Payments by Invoice
- **Method:** `GET`
- **Endpoint:** `/invoices/:invoiceId/payments`
- **Authentication:** Bearer Token (Required)

---

## Bills

### 1. Create Bill
- **Method:** `POST`
- **Endpoint:** `/bills`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Request Body:**
```json
{
  "vendor_id": 1,
  "billing_address": "123 Vendor Street",
  "order_number": "PO-12345",
  "bill_date": "2026-02-07",
  "due_date": "2026-03-07",
  "payment_terms": "net_30",
  "subject": "Bill for materials",
  "line_items": [
    {
      "item_id": "ITEM-001",
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
- **Response:**
```json
{
  "success": true,
  "data": {
    "id": "BILL-123456",
    "bill_number": "BILL-1707298800",
    "vendor_id": 1,
    "bill_date": "2026-02-07",
    "due_date": "2026-03-07",
    "sub_total": 15000,
    "discount": 1000,
    "tax_amount": 2520,
    "total": 17020,
    "status": "draft",
    "created_at": "2026-02-07T10:00:00Z"
  },
  "message": "Bill created successfully"
}
```

### 2. Get All Bills
- **Method:** `GET`
- **Endpoint:** `/bills?limit=10&offset=0`
- **Authentication:** Bearer Token (Required)

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

### 6. Update Bill Status
- **Method:** `PATCH`
- **Endpoint:** `/bills/:id/status`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Request Body:**
```json
{
  "status": "sent"
}
```

### 7. Get Bills by Vendor
- **Method:** `GET`
- **Endpoint:** `/bills/vendor/:vendorId`
- **Authentication:** Bearer Token (Required)

### 8. Get Bills by Status
- **Method:** `GET`
- **Endpoint:** `/bills/status/:status`
- **Authentication:** Bearer Token (Required)
- **Status Values:** draft, sent, partial, paid, overdue, void

---

## Payments

### 1. Create Payment
- **Method:** `POST`
- **Endpoint:** `/payments`
- **Authentication:** Bearer Token + Admin Role (Required)
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
    "payment_method": "bank_transfer",
    "status": "completed",
    "created_at": "2026-02-07T10:00:00Z"
  },
  "message": "Payment created successfully"
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

## Tax Configuration

### 1. Create Tax
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
- **Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "CGST 9%",
    "tax_type": "CGST",
    "rate": 9,
    "created_at": "2026-02-07T10:00:00Z"
  }
}
```

### 2. Get All Taxes
- **Method:** `GET`
- **Endpoint:** `/taxes?limit=10&offset=0`
- **Authentication:** Bearer Token (Required)

### 3. Get Tax by ID
- **Method:** `GET`
- **Endpoint:** `/taxes/:id`
- **Authentication:** Bearer Token (Required)

### 4. Update Tax
- **Method:** `PUT`
- **Endpoint:** `/taxes/:id`
- **Authentication:** Bearer Token + Admin Role (Required)

### 5. Delete Tax
- **Method:** `DELETE`
- **Endpoint:** `/taxes/:id`
- **Authentication:** Bearer Token + SuperAdmin Role (Required)

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

## Purchase Orders

### 1. Create Purchase Order
- **Method:** `POST`
- **Endpoint:** `/purchase-orders`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Request Body:**
```json
{
  "vendor_id": 1,
  "reference_no": "REF-001",
  "po_date": "2026-02-07",
  "expected_delivery_date": "2026-02-14",
  "payment_terms": "net_30",
  "line_items": [
    {
      "item_id": "ITEM-001",
      "quantity": 100,
      "rate": 150
    }
  ],
  "shipping_charges": 500,
  "tax_id": 1,
  "adjustment": 100,
  "vendor_notes": "Notes for vendor"
}
```

### 2. Get All Purchase Orders
- **Method:** `GET`
- **Endpoint:** `/purchase-orders?limit=10&offset=0`
- **Authentication:** Bearer Token (Required)

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

### 6. Update Purchase Order Status
- **Method:** `PATCH`
- **Endpoint:** `/purchase-orders/:id/status`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Status Values:** draft, sent, partially_received, received, cancelled

### 7. Get Purchase Orders by Vendor
- **Method:** `GET`
- **Endpoint:** `/purchase-orders/vendor/:vendorId`
- **Authentication:** Bearer Token (Required)

### 8. Get Purchase Orders by Customer
- **Method:** `GET`
- **Endpoint:** `/purchase-orders/customer/:customerId`
- **Authentication:** Bearer Token (Required)

### 9. Get Purchase Orders by Status
- **Method:** `GET`
- **Endpoint:** `/purchase-orders/status/:status`
- **Authentication:** Bearer Token (Required)

---

## Sales Orders

### 1. Create Sales Order
- **Method:** `POST`
- **Endpoint:** `/sales-orders`
- **Authentication:** Bearer Token + Admin Role (Required)
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
      "item_id": "ITEM-001",
      "quantity": 10,
      "rate": 5000
    }
  ],
  "shipping_charges": 500,
  "tax_id": 1,
  "adjustment": 100,
  "customer_notes": "Handle with care"
}
```

### 2. Get All Sales Orders
- **Method:** `GET`
- **Endpoint:** `/sales-orders?limit=10&offset=0`
- **Authentication:** Bearer Token (Required)

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

### 6. Update Sales Order Status
- **Method:** `PATCH`
- **Endpoint:** `/sales-orders/:id/status`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Status Values:** draft, sent, confirmed, partial_shipped, shipped, delivered, cancelled

### 7. Get Sales Orders by Customer
- **Method:** `GET`
- **Endpoint:** `/sales-orders/customer/:customerId`
- **Authentication:** Bearer Token (Required)

### 8. Get Sales Orders by Status
- **Method:** `GET`
- **Endpoint:** `/sales-orders/status/:status`
- **Authentication:** Bearer Token (Required)

---

## Packages

### 1. Create Package
- **Method:** `POST`
- **Endpoint:** `/packages`
- **Authentication:** Bearer Token + Admin Role (Required)
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

### 6. Update Package Status
- **Method:** `PATCH`
- **Endpoint:** `/packages/:id/status`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Status Values:** created, packed, shipped, delivered, cancelled

### 7. Get Packages by Customer
- **Method:** `GET`
- **Endpoint:** `/packages/customer/:customer_id`
- **Authentication:** Bearer Token (Required)

### 8. Get Packages by Sales Order
- **Method:** `GET`
- **Endpoint:** `/packages/sales-order/:sales_order_id`
- **Authentication:** Bearer Token (Required)

### 9. Get Packages by Status
- **Method:** `GET`
- **Endpoint:** `/packages/status/:status`
- **Authentication:** Bearer Token (Required)

---

## Shipments

### 1. Create Shipment
- **Method:** `POST`
- **Endpoint:** `/shipments`
- **Authentication:** Bearer Token + Admin Role (Required)
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

### 6. Update Shipment Status
- **Method:** `PATCH`
- **Endpoint:** `/shipments/:id/status`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Status Values:** created, shipped, in_transit, delivered, cancelled

### 7. Get Shipments by Customer
- **Method:** `GET`
- **Endpoint:** `/shipments/customer/:customer_id`
- **Authentication:** Bearer Token (Required)

### 8. Get Shipments by Package
- **Method:** `GET`
- **Endpoint:** `/shipments/package/:package_id`
- **Authentication:** Bearer Token (Required)

### 9. Get Shipments by Sales Order
- **Method:** `GET`
- **Endpoint:** `/shipments/sales-order/:sales_order_id`
- **Authentication:** Bearer Token (Required)

### 10. Get Shipments by Status
- **Method:** `GET`
- **Endpoint:** `/shipments/status/:status`
- **Authentication:** Bearer Token (Required)

---

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

## Forward Auth Routes

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
  },
  "message": "Support ticket created successfully"
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

## Common Response Formats

### Success Response
```json
{
  "success": true,
  "data": { /* response data */ },
  "message": "Operation successful"
}
```

### Error Response
```json
{
  "success": false,
  "error": "Error message",
  "status_code": 400
}
```

### Paginated Response
```json
{
  "success": true,
  "data": [ /* array of items */ ],
  "total": 100,
  "limit": 10,
  "offset": 0
}
```

---

## Authentication Header

All protected endpoints require the following header:
```
Authorization: Bearer <access_token>
```

---

## Status Codes

- `200 OK` - Successful GET request
- `201 Created` - Successful POST request
- `204 No Content` - Successful DELETE request
- `400 Bad Request` - Invalid request body or parameters
- `401 Unauthorized` - Missing or invalid authentication token
- `403 Forbidden` - Insufficient permissions
- `404 Not Found` - Resource not found
- `409 Conflict` - Duplicate resource
- `500 Internal Server Error` - Server error

---

## Pagination

Most GET endpoints support pagination with query parameters:
- `limit` - Number of items per page (default: 10, max: 100)
- `offset` - Number of items to skip (default: 0)

Example: `/invoices?limit=20&offset=40`

---

## Swagger Documentation

For interactive API documentation, visit:
```
GET /docs/
```

This provides a Swagger UI where you can test all endpoints directly.

---

**Last Updated:** February 7, 2026
