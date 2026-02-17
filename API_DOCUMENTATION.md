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
11. [Banks](#banks)
12. [Companies](#companies)
13. [Invoices](#invoices)
14. [Bills](#bills)
15. [Payments](#payments)
16. [Tax Configuration](#tax-configuration)
17. [Salespersons](#salespersons)
18. [Purchase Orders](#purchase-orders)
19. [Sales Orders](#sales-orders)
20. [Packages](#packages)
21. [Shipments](#shipments)
22. [Helper/Lookup Routes](#helperlookup-routes)
23. [Forward Auth Routes](#forward-auth-routes)
24. [Support](#support)
25. [API Workflow & Dependencies](#api-workflow--dependencies)

---

## ⚠️ IMPORTANT: Key API Architecture Notes

**Before using this API, please understand these critical concepts:**

### 1. **Banks are Master Reference Data (Now Creatable via API)**
- Banks can NOW be created via `POST /banks` endpoint (requires SuperAdmin role)
- Banks are master reference data used across the system
- When setting up company or vendor bank details, you must reference an existing `bank_id`
- **Workflow:** Create bank first via `POST /banks` → Then reference bank_id when adding bank details to companies/vendors
- **Error Resolution:** If you get "Foreign key constraint fails for bank_id", ensure the bank exists before referencing it

### 2. **Nested Entities are Sent Inline, Not Via Separate Endpoints**
- **Bank Details, Addresses, Contact Persons**: These are included in the vendor/company/customer creation request, NOT created via separate endpoints
- **Example:** When creating a vendor, include `bank_details`, `contact_persons`, `billing_address`, `shipping_address` in the same request
- **NO separate endpoints exist** for:
  - `POST /vendors/:id/bank-details` ❌ (Use POST /vendors with nested data instead)
  - `POST /vendors/:id/addresses` ❌ (Use POST /vendors with nested data instead)
  - `POST /vendors/:id/contact-persons` ❌ (Use POST /vendors with nested data instead)

### 3. **Prerequisites Before Creating Entities**
- **Company Setup Requires:** Business Type (from helpers), optionally a pre-existing Bank
- **Vendor Setup Requires:** Optionally a pre-existing Bank
- **Invoices/Bills Require:** Company, Customer/Vendor, Items, Tax, Salesperson (optional)
- **Purchase Orders Require:** Company, Vendor, Items, Tax
- **Sales Orders Require:** Company, Customer, Items, Tax, Salesperson (optional)

### 4. **How to Handle Business Logic**
- **Create Entity First, Update Details Later:** If you run into issues, create the main entity first without optional nested data, then update it
- **Status Transitions**: Some entities have status changes (Draft → Sent → Confirmed). Use PATCH endpoints to update status
- **Inventory Management**: After creating items, optionally set opening stock via `PUT /items/:id/opening-stock`

---

## API Workflow & Dependencies

### Entity Creation Workflow

This section outlines the correct order to create entities based on their foreign key dependencies. Follow these steps to avoid constraint errors.

#### **Step 1: Authentication & User Setup**
Create user accounts and authentication tokens first.

```
1. Register User (Auth) → Get Access Token → Use for subsequent API calls
```

**Endpoints:**
- `POST /auth/register/email` - Create user account
- `POST /auth/login/email` - Get access token

---

#### **Step 2: Master Data Setup (Banks, Business Types, Countries)**
Create and verify master data entities needed for subsequent operations.

```
2a. Create Banks (if needed)
2b. Get Business Types (Lookup - no creation needed)
2c. Get Countries & States (Lookup - no creation needed)
```

**Endpoints:**
- `POST /banks` - Create new bank master records (requires SuperAdmin role)
- `GET /banks` - Get all available banks
- `GET /helpers/business-types` - Get available business types
- `GET /helpers/countries` - Get countries
- `GET /helpers/countries/:country_id/states` - Get states

**Important:** Banks are now part of the master data that you can create via API. If you need a bank that doesn't exist, create it first using `POST /banks`, then reference the bank_id when setting up bank details for companies or vendors (Step 3 and 4).

---

#### **Step 3: Company Setup**
Create company records (requires business_type_id from Step 2, and optionally a bank_id if linking bank details).

```
3. Create Company → Add Contact → Add Address → Add Bank Details → Add Tax Settings
```

**Order:**
1. `POST /companies/setup` - Create company
2. `PUT /companies/:id/contact` - Add contact info
3. `PUT /companies/:id/address` - Add address
4. `POST /companies/:id/bank-details` - Link bank account (uses bank_id from Step 2)
5. `PUT /companies/:id/tax-settings` - Configure taxes

**Why This Order?** Company contact and address are optional but should be set up before using the company for transactions.

---

#### **Step 4: Vendors & Customers Setup**
Create vendors and customers (no dependencies, but can optionally link banks).

```
4a. Create Vendor → Add Contact Persons → Add Bank Details
4b. Create Customer → Add Contact Persons → Add Addresses
```

**Vendors - Step by Step:**
```
1. POST /vendors - Create basic vendor info
   {
     "first_name": "John",
     "last_name": "Doe",
     "company_name": "Acme Supplies",
     "email_address": "john@acme.com",
     "work_phone": "1234567890",
     "mobile": "9876543210"
   }

2. Then optionally add bank details (requires bank_id from pre-populated banks database):
   - Bank records must already exist in the database
   - Reference the bank_id when calling POST /vendors/:id/bank-details
   - Do NOT try to create banks via API - they are read-only master data
```

**Customers - Step by Step:**
```
1. POST /customers - Create basic customer info
   {
     "first_name": "Amit",
     "last_name": "Singh",
     "company_name": "Singh Enterprises",
     "email_address": "amit@singh.com",
     "work_phone": "9876543210"
   }

2. Then add contact persons and addresses as needed
```

---

#### **Step 5: Manufacturers & Brands**
Create manufacturers and brands (no dependencies).

```
5a. Create Manufacturers
5b. Create Brands (can reference manufacturer_id from Step 5a)
```

**Endpoints:**
- `POST /manufacturers` - Create manufacturer
- `POST /brands` - Create brand

---

#### **Step 6: Items & Inventory**
Create items and variants (requires manufacturer_id from Step 5).

```
6. Create Item → Add Variants → Set Opening Stock
```

**Order:**
1. `POST /items` - Create item with manufacturer reference
2. Item variants are created with the item
3. `PUT /items/:id/opening-stock` - Set initial inventory

---

#### **Step 7: Tax & Salesperson Configuration**
Set up taxes and salespersons for orders and invoices.

```
7a. Create Tax Types
7b. Create Salespersons
```

**Endpoints:**
- `POST /taxes` - Create tax configuration
- `POST /salespersons` - Create sales team members

---

#### **Step 8: Purchase Orders (Water Company)**
Create purchase orders with bottle suppliers (requires vendors from Step 4a and items from Step 6).

**Water Supplies Workflow:**
```
8a. Create PO with bottle manufacturer
8b. Add items: 500ml Bottles, 20L Cooler Bottles, Caps, Labels
8c. Add tax (5% on plastic), shipping charges
8d. Submit PO for delivery
```

**Example Items in PO:**
- 5000 × 500ml PET Bottles @ ₹8 = ₹40,000
- 50 × Cap packs (100 each) @ ₹25 = ₹1,250
- 100 × 20L Bottles @ ₹80 = ₹8,000

**Order:**
1. `POST /purchase-orders` - Create PO with supplier
2. `POST /purchase-orders/:po_id/line-items` - Add water bottles and caps
3. `PUT /purchase-orders/:po_id` - Set tax (5%), shipping, discounts
4. `PATCH /purchase-orders/:po_id/status` - Submit to vendor
5. `POST /bills` - Create bill when received

---

#### **Step 9: Sales Orders, Packages & Shipments (Water Company)**
Create sales orders to sell water bottles to retailers/customers.

**Water Distribution Workflow:**
```
9a. Create SO with customer (retail store, office, distributor)
9b. Add items: 500ml Bottles, 20L Cooler Bottles
9c. Set delivery terms and shipping address
9d. Create Package → Shipment → Invoice
```

**Example Order Items:**
- 1000 × 500ml Bottles @ ₹15 = ₹15,000 (for retail shop)
- 50 × 20L Bottles @ ₹150 = ₹7,500 (for offices/coolers)
- Total: ₹22,500 (with 5% tax: ₹23,625)

**Order:**
1. `POST /sales-orders` - Create SO with customer
2. `POST /sales-orders/:so_id/line-items` - Add water bottles
3. `PUT /sales-orders/:so_id/reserve-inventory` - Lock warehouse stock
4. `POST /packages` - Create shipment package
5. `POST /shipments` - Generate shipping label, tracking
6. `POST /invoices` - Issue customer invoice
7. `POST /payments` - Record payment received

---

#### **Step 10: Item Groups & Inventory (Water Company)**
Create item groups for packaged water products.

**Water Company BOMs:**
```
10a. Item Group: "500ml Complete Bottle" = Bottle + Cap + Label
10b. Item Group: "20L Cooler Set" = Bottle + Cap + Cleaning Kit
10c. Track inventory for each component
```

**Endpoints:**
- `POST /api/item-groups` - Create BOM for packaged products
- `PUT /items/:id/opening-stock` - Set warehouse stock levels

**Endpoints:**
- `POST /api/item-groups` - Create BOM
- `POST /api/production-orders` - Create manufacturing order
- `GET /api/inventory/*` - Track inventory

---

### Dependency Graph Summary

```
┌─────────────────────────────────────────────────────────────────┐
│ STEP 1: AUTH                                                    │
├─────────────────────────────────────────────────────────────────┤
│ Register User → Login → Get Access Token                         │
└────────────────────────────┬────────────────────────────────────┘
                             │
┌────────────────────────────▼────────────────────────────────────┐
│ STEP 2: MASTER DATA (No Dependencies)                           │
├─────────────────────────────────────────────────────────────────┤
│ • Banks                                                           │
│ • Business Types                                                  │
│ • Countries & States                                              │
└────────┬───────────────────────────────────┬─────────────────────┘
         │                                   │
    ┌────▼──────┐                    ┌──────▼─────┐
    │ Bank ID   │                    │ Bus Type ID│
    └────┬──────┘                    └──────┬─────┘
         │                                   │
┌────────▼───────────────────────────────────▼─────────────────────┐
│ STEP 3: COMPANY (needs: Business Type)                           │
├─────────────────────────────────────────────────────────────────┤
│ Company → Contact → Address → Bank Details → Settings            │
└────────────┬──────────────────────────────────────────────────────┘
             │
    ┌────────▼────────┐
    │                 │
    │ Company ID      │
    │ Bank ID         │
    │ (referenced)    │
    │                 │
    └────────┬────────┘
             │
┌────────────▼────────────────────────────────────────────────────┐
│ STEP 4: VENDORS & CUSTOMERS (No Dependencies)                   │
├──────────────────────────────────────────────────────────────────┤
│ 4a. Vendor → Bank Details                                        │
│ 4b. Customer → Addresses & Contacts                              │
└────────┬───────────────────────────────────────────┬─────────────┘
         │                                           │
    ┌────▼──────────┐                        ┌──────▼────────┐
    │ Vendor ID     │                        │ Customer ID   │
    │ Bank ID (ref) │                        │ Address Types │
    └────┬──────────┘                        └──────┬────────┘
         │                                          │
         │      ┌──────────────────────────────────┘
         │      │
┌────────▼──────▼────────────────────────────────────────────────┐
│ STEP 5: MANUFACTURERS & BRANDS (No Dependencies)               │
├──────────────────────────────────────────────────────────────────┤
│ Manufacturer → Brand (can reference Manufacturer ID)             │
└────────┬──────────────────────────────────────────────────────┬──┘
         │                                                       │
    ┌────▼──────────┐                                    ┌──────▼─────┐
    │ Manufacturer  │                                    │ Brand ID   │
    │       ID      │                                    └────────────┘
    └────┬──────────┘
         │
┌────────▼──────────────────────────────────────────────────────┐
│ STEP 6: ITEMS (needs: Manufacturer)                           │
├──────────────────────────────────────────────────────────────────┤
│ Item → Variants → Opening Stock                                 │
│ (references: Manufacturer ID)                                   │
└────────┬──────────────────────────────────────────┬──────────┬──┘
         │                                          │          │
    ┌────▼──────────┐                       ┌──────▼────┐ ┌──▼────┐
    │ Item ID       │                       │ Variant ID│ │Tax ID │
    │ Variant ID    │                       └───────────┘ └───────┘
    └────┬──────────┘
         │
┌────────▼──────────────────────────────────────────────────────┐
│ STEP 7: TAXES & SALESPERSONS (No Dependencies)                │
├──────────────────────────────────────────────────────────────────┤
│ Taxes → Salespersons                                            │
└────────┬──────────────────────────────────────┬─────────────┬──┘
         │                                      │             │
    ┌────▼──────┐                         ┌────▼────────┐ ┌─▼────────┐
    │ Tax ID    │                         │ Salesman ID │ │ Company  │
    └───────────┘                         └─────────────┘ │    ID    │
                                                          └──────────┘
         │                                      │             │
         └──────────────────┬───────────────────┘             │
                            │                                 │
┌───────────────────────────▼──────────────────────────────────▼──┐
│ STEP 8: PURCHASE ORDERS & BILLS                                 │
├─────────────────────────────────────────────────────────────────┤
│ (needs: Vendor, Item, Company, Tax)                             │
│ PurchaseOrder → Bill → Payment                                  │
└───────────────────────────┬──────────────────────────────────────┘
                            │
                    ┌───────▼────────┐
                    │ PO ID, Bill ID │
                    └────────────────┘
                            │
┌───────────────────────────▼──────────────────────────────────────┐
│ STEP 9: SALES ORDERS, PACKAGES & SHIPMENTS                      │
├─────────────────────────────────────────────────────────────────┤
│ (needs: Customer, Item, Salesman, Company, Tax)                 │
│ SalesOrder → Package → Shipment → Invoice                       │
└───────────────────────────┬──────────────────────────────────────┘
                            │
                    ┌───────▼──────────┐
                    │ SO, PKG, SHIP ID │
                    └──────────────────┘
                            │
┌───────────────────────────▼──────────────────────────────────────┐
│ STEP 10: INVENTORY & ITEM GROUPS                                │
├─────────────────────────────────────────────────────────────────┤
│ (needs: Items)                                                  │
│ ItemGroup → ProductionOrder → InventoryTracking                 │
└─────────────────────────────────────────────────────────────────┘
```

---

### Common Error Resolution

**Error: "Foreign key constraint fails: bank_id"**
- **Cause:** Creating vendor/company with non-existent bank_id
- **Solution:** CREATE bank first using `POST /banks` (requires SuperAdmin role), then use that bank_id when calling `POST /companies/:id/bank-details` or include it in vendor/company creation request. 

**Error: "Foreign key constraint fails: vendor_id"**
- **Cause:** Trying to create bill/PO with non-existent vendor
- **Solution:** CREATE vendor first using `POST /vendors`, then use that vendor_id in bill/PO creation

**Error: "Foreign key constraint fails: customer_id"**
- **Cause:** Trying to create sales order/invoice with non-existent customer
- **Solution:** CREATE customer first using `POST /customers`, then use that customer_id in SO/invoice creation

**Error: "Foreign key constraint fails: item_id"**
- **Cause:** Trying to use non-existent item in PO/SO/Invoice line items
- **Solution:** CREATE item first using `POST /items`, then reference item_id in line items

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
  "email": "jane.smith@example.com",
  "password": "SecurePassword123!",
  "first_name": "Jane",
  "last_name": "Smith",
  "phone": "9876543210",
  "phone_code": "+91",
  "role": "admin",
  "status": "active",
  "department": "Operations",
  "designation": "Senior Admin",
  "profile_picture_url": "https://example.com/jane-smith.jpg"
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
  "email": "contact@abccorp.com",
  "password": "SecurePassword123!",
  "first_name": "John",
  "last_name": "Smith",
  "phone": "9876543210",
  "phone_code": "+91",
  "company_name": "ABC Corporation Ltd.",
  "display_name": "ABC Corp",
  "gstin": "18AABCT1234H1Z0",
  "pan": "AAACT1234H",
  "website_url": "https://www.abccorp.com",
  "registration_number": "NCT/67890/2020",
  "status": "active"
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

### 1. Create Vendor (Complete Field Reference)

#### Available Fields Documentation

**Method:** `POST`  
**Endpoint:** `/vendors`  
**Authentication:** Bearer Token + SuperAdmin Role (Required)  
**Description:** Create vendor with complete information. All fields are optional unless marked as required.

##### Core Vendor Fields

| Field Name | Type | Required | Constraints | Alternative Names | Example |
|---|---|---|---|---|---|
| `salutation` | string | No | - | - | "Mr.", "Ms.", "Dr." |
| `first_name` | string | No | - | - | "John" |
| `last_name` | string | No | - | - | "Doe" |
| `company_name` | string | No | - | - | "Acme Supplies" |
| `display_name` | string | No | - | - | "Acme Supplies Inc." |
| `email_address` | string | No | Valid email format | - | "john.doe@acmesupplies.com" |
| `work_phone` | string | No | - | - | "1234567890" |
| `work_phone_code` | string | No | - | - | "+91", "+1" |
| `mobile` | string | No | - | `phone` | "9876543210" |
| `mobile_code` | string | No | - | `phone_code` | "+91", "+1" |
| `vendor_language` | string | No | - | - | "English", "Hindi" |

##### Other Details (Vendor-Specific)

The `other_details` object contains vendor-specific information:

| Field Name | Type | Required | Constraints | Example |
|---|---|---|---|---|
| `pan` | string | No | Exactly 10 characters | "ABCDE1234F" |
| `is_msme_registered` | boolean | No | true/false | true |
| `currency` | string | No | - | "INR", "USD", "EUR" |
| `payment_terms` | string | No | - | "Net 30", "Net 45", "Net 60" |
| `tds` | string | No | - | "10%", "5%" |
| `enable_portal` | boolean | No | true/false | true |
| `website_url` | string | No | Valid URL format | "https://www.acmesupplies.com" |
| `department` | string | No | - | "Sales", "Operations" |
| `designation` | string | No | - | "Manager", "Director" |
| `twitter` | string | No | - | "@acmesupplies" |
| `skype_name` | string | No | - | "acme.supplies" |
| `facebook` | string | No | - | "acmesupplies" |

##### Billing Address (Complete Fields)

The `billing_address` object contains:

| Field Name | Type | Required | Constraints | Alternative Names | Example |
|---|---|---|---|---|---|
| `attention` | string | No | - | - | "Finance Department" |
| `address_line1` | string | No | - | `street` | "123 Main Street" |
| `address_line2` | string | No | - | - | "Suite 100" |
| `city` | string | No | - | - | "Mumbai" |
| `state` | string | No | - | - | "Maharashtra" |
| `country_region` | string | No | - | `country` | "India" |
| `pin_code` | string | No | Max 10 chars | `postal_code` | "400001" |
| `phone` | string | No | - | - | "1234567890" |
| `phone_code` | string | No | - | - | "+91" |
| `fax_number` | string | No | - | - | "02212345678" |

##### Shipping Address (Complete Fields)

The `shipping_address` object contains the same fields as billing_address:

| Field Name | Type | Required | Constraints | Alternative Names | Example |
|---|---|---|---|---|---|
| `attention` | string | No | - | - | "Warehouse Department" |
| `address_line1` | string | No | - | `street` | "456 Shipping Rd" |
| `address_line2` | string | No | - | - | "Building B" |
| `city` | string | No | - | - | "Mumbai" |
| `state` | string | No | - | - | "Maharashtra" |
| `country_region` | string | No | - | `country` | "India" |
| `pin_code` | string | No | Max 10 chars | `postal_code` | "400001" |
| `phone` | string | No | - | - | "1234567890" |
| `phone_code` | string | No | - | - | "+91" |
| `fax_number` | string | No | - | - | "02212345678" |

##### Contact Persons Array (Each Contact Has)

The `contact_persons` array can contain multiple contact objects:

| Field Name | Type | Required | Constraints | Alternative Names | Example |
|---|---|---|---|---|---|
| `salutation` | string | No | - | `title` | "Mr.", "Ms.", "Dr." |
| `first_name` | string | No | - | - | "Rajesh" |
| `last_name` | string | No | - | - | "Kumar" |
| `email_address` | string | No | Valid email format | `email` | "rajesh@acme.com" |
| `work_phone` | string | No | - | - | "1234567890" |
| `work_phone_code` | string | No | - | - | "+91" |
| `mobile` | string | No | - | `phone` | "9876543211" |
| `mobile_code` | string | No | - | `phone_code` | "+91" |

##### Bank Details Array (Each Bank Account Has)

The `bank_details` array can contain multiple bank account objects:

| Field Name | Type | Required | Constraints | Description |
|---|---|---|---|---|
| `bank_id` | number | **Yes** | Must exist in banks table | Reference to bank master |
| `account_holder_name` | string | No | - | Name on bank account |
| `account_number` | string | **Yes** | - | Bank account number |
| `reenter_account_number` | string | **Yes** | Must match account_number | Confirmation field |

---

#### Request Body Examples

**Complete Example (All Fields):**
```json
{
  "salutation": "Mr.",
  "first_name": "John",
  "last_name": "Doe",
  "company_name": "Acme Supplies",
  "display_name": "Acme Supplies",
  "email_address": "john.doe@acmesupplies.com",
  "work_phone": "1234567890",
  "work_phone_code": "+91",
  "mobile": "9876543210",
  "mobile_code": "+91",
  "vendor_language": "English",
  "other_details": {
    "pan": "ABCDE1234F",
    "is_msme_registered": true,
    "currency": "INR",
    "payment_terms": "Net 30",
    "tds": "10%",
    "enable_portal": true,
    "website_url": "https://www.acmesupplies.com",
    "department": "Sales",
    "designation": "Director",
    "twitter": "@acmesupplies",
    "skype_name": "acme.supplies",
    "facebook": "acmesupplies"
  },
  "billing_address": {
    "attention": "Finance Department",
    "street": "123 Main Street",
    "address_line2": "Suite 100",
    "city": "Mumbai",
    "state": "Maharashtra",
    "country": "India",
    "postal_code": "400001",
    "phone": "1234567890",
    "phone_code": "+91",
    "fax_number": "02212345678"
  },
  "shipping_address": {
    "attention": "Warehouse Department",
    "street": "456 Shipping Rd",
    "city": "Mumbai",
    "state": "Maharashtra",
    "country": "India",
    "postal_code": "400001",
    "phone": "1234567890",
    "phone_code": "+91"
  },
  "contact_persons": [
    {
      "title": "Mr.",
      "first_name": "Rajesh",
      "last_name": "Kumar",
      "email": "rajesh@acme.com",
      "phone": "9876543211",
      "phone_code": "+91"
    }
  ],
  "bank_details": [
    {
      "bank_id": 1,
      "account_holder_name": "Acme Supplies",
      "account_number": "1234567890123456",
      "reenter_account_number": "1234567890123456"
    }
  ]
}
```

**Minimal Example (Required Fields Only):**
```json
{
  "company_name": "Acme Supplies",
  "display_name": "Acme Supplies"
}
```

**Alternative Field Names Example** (These fields are automatically normalized to primary names):
```json
{
  "salutation": "Mr.",
  "first_name": "John",
  "display_name": "Acme Supplies",
  "billing_address": {
    "street": "123 Main Street",
    "country": "India",
    "postal_code": "400001"
  },
  "contact_persons": [
    {
      "title": "Mr.",
      "first_name": "Rajesh",
      "email": "rajesh@acme.com",
      "phone": "9876543211"
    }
  ]
}
```

**Note:** The API intelligently maps alternative field names to primary fields:
- `street` → `address_line1`
- `country` → `country_region`
- `postal_code` → `pin_code`
- `title` → `salutation` (in contact persons)
- `email` → `email_address` (in contact persons)
- `phone` → `mobile` (in contact persons)
- `phone_code` → `mobile_code` (in contact persons)
- **Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "first_name": "John",
    "last_name": "Doe",
    "display_name": "Acme Supplies",
    "company_name": "Acme Supplies",
    "salutation": "Mr.",
    "email_address": "john.doe@acmesupplies.com",
    "work_phone": "1234567890",
    "work_phone_code": "+91",
    "mobile": "9876543210",
    "mobile_code": "+91",
    "vendor_language": "English",
    "other_details": {
      "pan": "ABCDE1234F",
      "is_msme_registered": true,
      "currency": "INR",
      "payment_terms": "Net 30",
      "tds": "10%",
      "enable_portal": true,
      "website_url": "https://www.acmesupplies.com",
      "department": "Sales",
      "designation": "Director",
      "twitter": "@acmesupplies",
      "skype_name": "acme.supplies",
      "facebook": "acmesupplies"
    },
    "billing_address": {
      "id": 1,
      "attention": "Finance Department",
      "address_line1": "123 Main Street",
      "address_line2": "Suite 100",
      "city": "Mumbai",
      "state": "Maharashtra",
      "country_region": "India",
      "pin_code": "400001",
      "phone": "1234567890",
      "phone_code": "+91",
      "fax_number": "02212345678"
    },
    "shipping_address": {
      "id": 2,
      "attention": "Warehouse Department",
      "address_line1": "456 Shipping Rd",
      "city": "Mumbai",
      "state": "Maharashtra",
      "country_region": "India",
      "pin_code": "400001"
    },
    "contact_persons": [
      {
        "id": 1,
        "salutation": "Mr.",
        "first_name": "Rajesh",
        "last_name": "Kumar",
        "email_address": "rajesh@acme.com",
        "mobile": "9876543211",
        "mobile_code": "+91"
      }
    ],
    "bank_details": [
      {
        "id": 1,
        "bank_id": 1,
        "account_holder_name": "Acme Supplies",
        "account_number": "****3456"
      }
    ],
    "created_at": "2026-02-07T10:00:00Z"
  },
  "message": "Vendor created successfully"
}
```

---

### 2. Update Vendor

**Method:** `PUT`  
**Endpoint:** `/vendors/:id`  
**Authentication:** Bearer Token + SuperAdmin Role (Required)  
**Description:** Update any vendor field. All fields from the create request are available for update.

**Update Request Examples:**

**Update Core Information:**
```json
{
  "email_address": "newemail@acmesupplies.com",
  "display_name": "Acme Global Supplies",
  "vendor_language": "Hindi"
}
```

**Update Only Other Details:**
```json
{
  "other_details": {
    "pan": "BCDEF5678G",
    "currency": "USD",
    "payment_terms": "Net 45",
    "website_url": "https://www.acmeglobal.com"
  }
}
```

**Update Bank Details:**
```json
{
  "bank_details": [
    {
      "bank_id": 2,
      "account_holder_name": "Acme Supplies Pvt Ltd",
      "account_number": "9876543210987654",
      "reenter_account_number": "9876543210987654"
    }
  ]
}
```

**Update Addresses and Contacts:**
```json
{
  "billing_address": {
    "address_line1": "789 New Street",
    "city": "Bangalore",
    "state": "Karnataka"
  },
  "contact_persons": [
    {
      "salutation": "Ms.",
      "first_name": "Priya",
      "last_name": "Singh",
      "email_address": "priya@acme.com",
      "mobile": "9876543212"
    }
  ]
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "display_name": "Acme Global Supplies",
    "email_address": "newemail@acmesupplies.com",
    "vendor_language": "Hindi",
    "updated_at": "2026-02-07T11:00:00Z"
  },
  "message": "Vendor updated successfully"
}
```

---

### 3. Get All Vendors

**Method:** `GET`  
**Endpoint:** `/vendors?limit=10&offset=0`  
**Authentication:** None (Public)  
**Query Parameters:**
- `limit` (optional, default: 10) - Number of vendors per page
- `offset` (optional, default: 0) - Pagination offset

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "first_name": "John",
      "last_name": "Doe",
      "display_name": "Acme Supplies",
      "company_name": "Acme Supplies",
      "email_address": "john.doe@acmesupplies.com",
      "work_phone": "1234567890",
      "mobile": "9876543210",
      "created_at": "2026-02-07T10:00:00Z"
    }
  ],
  "total": 1,
  "limit": 10,
  "offset": 0
}
```

---

### 4. Get Vendor by ID

**Method:** `GET`  
**Endpoint:** `/vendors/:id`  
**Authentication:** None (Public)  

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "first_name": "John",
    "last_name": "Doe",
    "display_name": "Acme Supplies",
    "company_name": "Acme Supplies",
    "salutation": "Mr.",
    "email_address": "john.doe@acmesupplies.com",
    "work_phone": "1234567890",
    "work_phone_code": "+91",
    "mobile": "9876543210",
    "mobile_code": "+91",
    "vendor_language": "English",
    "other_details": {
      "pan": "ABCDE1234F",
      "is_msme_registered": true,
      "currency": "INR",
      "payment_terms": "Net 30",
      "tds": "10%",
      "enable_portal": true,
      "website_url": "https://www.acmesupplies.com"
    },
    "billing_address": {
      "id": 1,
      "address_line1": "123 Main Street",
      "address_line2": "Suite 100",
      "city": "Mumbai",
      "state": "Maharashtra",
      "country_region": "India",
      "pin_code": "400001"
    },
    "shipping_address": {
      "id": 2,
      "address_line1": "456 Shipping Rd",
      "city": "Mumbai",
      "state": "Maharashtra",
      "country_region": "India",
      "pin_code": "400001"
    },
    "contact_persons": [
      {
        "id": 1,
        "salutation": "Mr.",
        "first_name": "Rajesh",
        "last_name": "Kumar",
        "email_address": "rajesh@acme.com",
        "mobile": "9876543211"
      }
    ],
    "bank_details": [
      {
        "id": 1,
        "bank_id": 1,
        "account_holder_name": "Acme Supplies",
        "account_number": "****3456"
      }
    ],
    "created_at": "2026-02-07T10:00:00Z"
  }
}
```

---

### 5. Delete Vendor

**Method:** `DELETE`  
**Endpoint:** `/vendors/:id`  
**Authentication:** Bearer Token + SuperAdmin Role (Required)  

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Vendor deleted successfully"
}
```

**Error Response (409 Conflict - Foreign Key Constraint):**
```json
{
  "success": false,
  "message": "Cannot delete vendor. It is referenced by other records (e.g., Purchase Orders, Bills)",
  "error": "Foreign key constraint violation"
}
```

---

## Customers

### 1. Create Customer (Complete Field Reference)

**Method:** `POST`  
**Endpoint:** `/customers`  
**Authentication:** Bearer Token + SuperAdmin Role (Required)  
**Description:** Create customer with complete information. All fields are optional unless marked as required.

#### Customer Core Fields

| Field Name | Type | Required | Constraints | Alternative Names | Example |
|---|---|---|---|---|---|
| `salutation` | string | **Yes** | - | - | "Mr.", "Ms.", "Dr." |
| `first_name` | string | **Yes** | - | - | "Amit" |
| `last_name` | string | No | - | - | "Singh" |
| `company_name` | string | No | - | - | "Singh Enterprises Ltd." |
| `display_name` | string | **Yes** | - | - | "Singh Enterprises" |
| `email_address` | string | No | Valid email format | - | "amit.singh@singenterprises.com" |
| `work_phone` | string | No | - | - | "9876543210" |
| `work_phone_code` | string | No | - | - | "+91" |
| `mobile` | string | No | - | `phone` | "9875432109" |
| `mobile_code` | string | No | - | `phone_code` | "+91" |
| `customer_language` | string | No | - | - | "English", "Hindi" |

#### Customer Other Details

The `other_details` object contains customer-specific information:

| Field Name | Type | Required | Constraints | Example |
|---|---|---|---|---|
| `pan` | string | No | Exactly 10 characters | "ABCDS1234H" |
| `currency` | string | No | - | "INR", "USD" |
| `payment_terms` | string | No | - | "Net 45", "Net 30" |
| `enable_portal` | boolean | No | true/false | true |

#### Billing Address Fields

Same 13 fields as vendor (see Vendor section for address field details):
- attention, address_line1 (or street), address_line2, city, state, country_region (or country), pin_code (or postal_code), phone, phone_code, fax_number

#### Shipping Address Fields

Same 13 fields as vendor (see Vendor section for address field details)

#### Contact Persons Array

Same 12 fields as vendor (see Vendor section for contact person field details)

---

#### Request Body Examples

**Complete Customer Example:**
```json
{
  "salutation": "Mr.",
  "first_name": "Amit",
  "last_name": "Singh",
  "company_name": "Singh Enterprises Ltd.",
  "display_name": "Singh Enterprises",
  "email_address": "amit.singh@singenterprises.com",
  "work_phone": "9876543210",
  "work_phone_code": "+91",
  "mobile": "9875432109",
  "mobile_code": "+91",
  "customer_language": "English",
  "other_details": {
    "pan": "ABCDS1234H",
    "currency": "INR",
    "payment_terms": "Net 45",
    "enable_portal": true
  },
  "billing_address": {
    "attention": "Finance Department",
    "street": "789 Park Avenue",
    "address_line2": "Suite 500",
    "city": "New Delhi",
    "state": "Delhi",
    "country": "India",
    "postal_code": "110001",
    "phone": "01112345678",
    "phone_code": "+91"
  },
  "shipping_address": {
    "attention": "Receiving Department",
    "street": "789 Warehouse Road",
    "address_line2": "Building B",
    "city": "Noida",
    "state": "Uttar Pradesh",
    "country": "India",
    "postal_code": "201301",
    "phone": "01204321098",
    "phone_code": "+91"
  },
  "contact_persons": [
    {
      "title": "Ms.",
      "first_name": "Priya",
      "last_name": "Singh",
      "email": "priya.singh@singenterprises.com",
      "work_phone": "9876543215",
      "phone": "9875432115"
    }
  ]
}
```

**Minimal Customer Example (Required Fields Only):**
```json
{
  "salutation": "Mr.",
  "first_name": "Amit",
  "display_name": "Singh Enterprises"
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "first_name": "Amit",
    "last_name": "Singh",
    "display_name": "Singh Enterprises",
    "company_name": "Singh Enterprises Ltd.",
    "salutation": "Mr.",
    "email_address": "amit.singh@singenterprises.com",
    "work_phone": "9876543210",
    "work_phone_code": "+91",
    "mobile": "9875432109",
    "mobile_code": "+91",
    "customer_language": "English",
    "other_details": {
      "pan": "ABCDS1234H",
      "currency": "INR",
      "payment_terms": "Net 45",
      "enable_portal": true
    },
    "billing_address": {
      "id": 1,
      "attention": "Finance Department",
      "address_line1": "789 Park Avenue",
      "city": "New Delhi",
      "state": "Delhi",
      "country_region": "India",
      "pin_code": "110001"
    },
    "shipping_address": {
      "id": 2,
      "attention": "Receiving Department",
      "address_line1": "789 Warehouse Road",
      "city": "Noida",
      "state": "Uttar Pradesh",
      "country_region": "India",
      "pin_code": "201301"
    },
    "contact_persons": [
      {
        "id": 1,
        "salutation": "Ms.",
        "first_name": "Priya",
        "last_name": "Singh",
        "email_address": "priya.singh@singenterprises.com",
        "mobile": "9875432115"
      }
    ],
    "created_at": "2026-02-07T10:00:00Z"
  },
  "message": "Customer created successfully"
}
```

---

### 2. Update Customer

**Method:** `PUT`  
**Endpoint:** `/customers/:id`  
**Authentication:** Bearer Token + SuperAdmin Role (Required)  
**Description:** Update any customer field. All fields from the create request are available for update.

**Update Request Examples:**

**Update Basic Information:**
```json
{
  "email_address": "amit.singh.new@singenterprises.com",
  "display_name": "Singh Global Enterprises",
  "customer_language": "Hindi"
}
```

**Update Other Details:**
```json
{
  "other_details": {
    "pan": "BCDEF5678G",
    "currency": "USD",
    "payment_terms": "Net 60",
    "enable_portal": false
  }
}
```

**Update Addresses and Contacts:**
```json
{
  "billing_address": {
    "address_line1": "New Park Avenue",
    "city": "Mumbai",
    "state": "Maharashtra"
  },
  "contact_persons": [
    {
      "salutation": "Ms.",
      "first_name": "Priya",
      "last_name": "Singh",
      "email_address": "priya.new@singenterprises.com",
      "mobile": "9875432115"
    }
  ]
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "display_name": "Singh Global Enterprises",
    "email_address": "amit.singh.new@singenterprises.com",
    "updated_at": "2026-02-07T11:00:00Z"
  },
  "message": "Customer updated successfully"
}
```

---

### 3. Get All Customers

**Method:** `GET`  
**Endpoint:** `/customers?limit=10&offset=0`  
**Authentication:** None (Public)  
**Query Parameters:**
- `limit` (optional, default: 10) - Number of customers per page
- `offset` (optional, default: 0) - Pagination offset

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "first_name": "Amit",
      "last_name": "Singh",
      "display_name": "Singh Enterprises",
      "email_address": "amit.singh@singenterprises.com"
    }
  ],
  "total": 1,
  "limit": 10,
  "offset": 0
}
```

---

### 4. Get Customer by ID

**Method:** `GET`  
**Endpoint:** `/customers/:id`  
**Authentication:** None (Public)  

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "first_name": "Amit",
    "last_name": "Singh",
    "display_name": "Singh Enterprises",
    "salutation": "Mr.",
    "email_address": "amit.singh@singenterprises.com",
    "work_phone": "9876543210",
    "mobile": "9875432109",
    "other_details": {
      "pan": "ABCDS1234H",
      "currency": "INR",
      "payment_terms": "Net 45"
    },
    "billing_address": {
      "id": 1,
      "address_line1": "789 Park Avenue",
      "city": "New Delhi"
    },
    "shipping_address": {
      "id": 2,
      "address_line1": "789 Warehouse Road",
      "city": "Noida"
    },
    "contact_persons": [
      {
        "id": 1,
        "first_name": "Priya",
        "email_address": "priya@singenterprises.com"
      }
    ],
    "created_at": "2026-02-07T10:00:00Z"
  }
}
```

---

### 5. Delete Customer

**Method:** `DELETE`  
**Endpoint:** `/customers/:id`  
**Authentication:** Bearer Token + SuperAdmin Role (Required)  

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Customer deleted successfully"
}
```

**Error Response (409 Conflict - Foreign Key Constraint):**
```json
{
  "success": false,
  "message": "Cannot delete customer. It is referenced by other records (e.g., Sales Orders, Invoices)",
  "error": "Foreign key constraint violation"
}
```

---

### 6. Get Customer Invoices (if available)

**Method:** `GET`  
**Endpoint:** `/customers/:customerId/invoices`  
**Authentication:** Bearer Token (Required)  

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "invoice_number": "INV-2026-001",
      "customer_id": 1,
      "amount": 150000,
      "status": "paid",
      "created_at": "2026-02-07T10:00:00Z"
    }
  ]
}
```

---

---

## Items & Inventory

### Water Company Item Examples

For a water company business, typical items include:

**Example 1: 500ml Plastic Bottle**
```json
{
  "name": "500ml Drinking Water Bottle - PET",
  "type": "product",
  "brand_id": 1,
  "manufacturer_id": 1,
  "unit": "piece",
  "item_details": {
    "sku": "WTR-BOT-500ML",
    "barcode": "8904220025641",
    "hsn_code": "3923.30",
    "description": "500ml PET plastic drinking water bottle with tamper-proof cap",
    "item_type": "beverage_container",
    "weight_grams": 25,
    "dimensions": "Height: 18cm, Diameter: 7cm"
  },
  "sales_info": {
    "sales_account": "Water Sales Revenue",
    "sales_tax_id": 1,
    "sales_tax_rate": 5,
    "selling_price": 15,
    "mrp": 20
  },
  "purchase_info": {
    "purchase_account": "Cost of Water Bottles",
    "purchase_tax_id": 1,
    "purchase_tax_rate": 5,
    "purchase_price": 8
  },
  "variants": [
    {
      "variant_name": "Regular Cap",
      "sku_suffix": "-REG",
      "cost_price": 8,
      "selling_price": 15
    },
    {
      "variant_name": "Sports Cap",
      "sku_suffix": "-SPORT",
      "cost_price": 9,
      "selling_price": 18
    }
  ],
  "is_active": true,
  "track_inventory": true
}
```

**Example 2: 20L Water Bottle (Bulk)**
```json
{
  "name": "20 Litre Polycarbonate Water Bottle",
  "type": "product",
  "brand_id": 1,
  "unit": "piece",
  "item_details": {
    "sku": "WTR-BOT-20L",
    "hsn_code": "3923.30",
    "description": "20 litre polycarbonate water bottle for water coolers",
    "item_type": "bulk_container",
    "weight_grams": 800
  },
  "sales_info": {
    "sales_account": "Water Sales Revenue",
    "sales_tax_id": 1,
    "sales_tax_rate": 5,
    "selling_price": 150,
    "mrp": 200
  },
  "purchase_info": {
    "purchase_account": "Cost of Bottles",
    "purchase_tax_id": 1,
    "purchase_tax_rate": 5,
    "purchase_price": 80
  },
  "is_active": true,
  "track_inventory": true
}
```

**Example 3: Plastic Caps & Lids**
```json
{
  "name": "Tamper-Proof Plastic Cap - Set of 100",
  "type": "product",
  "manufacturer_id": 2,
  "unit": "set",
  "item_details": {
    "sku": "CAP-100-TAM",
    "barcode": "8904220025642",
    "hsn_code": "3923.40",
    "description": "100 piece set of tamper-proof bottle caps, diameter 20mm",
    "item_type": "packaging_accessory",
    "weight_grams": 150
  },
  "sales_info": {
    "sales_tax_rate": 5,
    "selling_price": 45
  },
  "purchase_info": {
    "purchase_tax_rate": 5,
    "purchase_price": 25
  },
  "is_active": true,
  "track_inventory": true
}
```

**Example 4: Water Labels & Stickers**
```json
{
  "name": "Custom Water Bottle Labels - 1000 pieces",
  "type": "product",
  "unit": "pack",
  "item_details": {
    "sku": "LBL-1000-WAT",
    "hsn_code": "4816.40",
    "description": "Custom printed water bottle labels, waterproof, 1000 pieces per pack",
    "item_type": "packaging_label",
    "weight_grams": 200
  },
  "sales_info": {
    "sales_tax_rate": 18,
    "selling_price": 500
  },
  "purchase_info": {
    "purchase_tax_rate": 0,
    "purchase_price": 300
  },
  "is_active": true,
  "track_inventory": true
}
```

### 1. Create Item
- **Method:** `POST`
- **Endpoint:** `/items`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Request Body:**
```json
{
  "name": "500ml Drinking Water Bottle - PET",
  "type": "product",
  "brand_id": 1,
  "manufacturer_id": 1,
  "unit": "piece",
  "item_details": {
    "sku": "WTR-BOT-500ML",
    "barcode": "8904220025641",
    "hsn_code": "3923.30",
    "description": "500ml PET plastic drinking water bottle with tamper-proof cap",
    "item_type": "beverage_container",
    "weight_grams": 25
  },
  "sales_info": {
    "sales_account": "Water Sales Revenue",
    "sales_tax_id": 1,
    "sales_tax_rate": 5,
    "selling_price": 15,
    "mrp": 20
  },
  "purchase_info": {
    "purchase_account": "Cost of Water Bottles",
    "purchase_tax_id": 1,
    "purchase_tax_rate": 5,
    "purchase_price": 8
  },
  "variants": [
    {
      "variant_name": "Regular Cap",
      "sku_suffix": "-REG",
      "cost_price": 8,
      "selling_price": 15
    },
    {
      "variant_name": "Sports Cap",
      "sku_suffix": "-SPORT",
      "cost_price": 9,
      "selling_price": 18
    }
  ],
  "is_active": true,
  "track_inventory": true
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

### Water Company BOM Examples

**Example 1:** "500ml Complete Water Bottle" = 1 × 500ml Bottle + 1 × Cap + 1 × Label

**Example 2:** "20L Water Bottle Set" = 1 × 20L Bottle + 1 × Large Cap + Cleaning supplies

### 1. Create Item Group
- **Method:** `POST`
- **Endpoint:** `/api/item-groups`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Request Body (500ml Packaged Water Bottle):**
```json
{
  "name": "500ml Packaged Water Bottle",
  "description": "Complete 500ml water bottle with cap and label - ready for sale",
  "is_active": true,
  "components": [
    {
      "item_id": "wtr-bot-500ml",
      "variant_id": 1,
      "quantity": 1,
      "variant_details": {
        "capacity": "500ml",
        "material": "PET"
      }
    },
    {
      "item_id": "cap-100-tam",
      "quantity": 1,
      "variant_details": {
        "size": "20mm",
        "type": "tamper-proof"
      }
    },
    {
      "item_id": "lbl-1000-wat",
      "quantity": 0.001,
      "variant_details": {
        "type": "water_label"
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
  "description": "Leading global computer and technology manufacturer",
  "country": "USA",
  "state_region": "Texas",
  "city": "Round Rock",
  "address_line1": "1 Dell Way",
  "address_line2": "Round Rock Office",
  "postal_code": "78682",
  "email": "contact@dell.com",
  "phone": "+1-512-338-4400",
  "website_url": "https://www.dell.com",
  "logo_url": "https://example.com/dell-logo.png",
  "is_active": true
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
  "description": "Premium quality computer and technology products",
  "logo_url": "https://example.com/dell-logo.png",
  "website_url": "https://www.dell.com",
  "manufacturer_id": 1,
  "is_active": true,
  "established_year": 1984,
  "country": "USA",
  "email": "brand@dell.com"
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

## Banks

Bank Master table for managing bank information centrally. This is the master reference data for banks used throughout the system.

### 1. Create Bank
- **Method:** `POST`
- **Endpoint:** `/banks`
- **Authentication:** Bearer Token + SuperAdmin Role (Required)
- **Description:** Create a new bank master record. Banks created here can then be linked to companies, vendors, or customers.
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
- **Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "bank_name": "HDFC Bank",
    "ifsc_code": "HDFC0001234",
    "branch_name": "Mumbai Main Branch",
    "branch_code": "HDFC123",
    "address": "123 Banking Street",
    "city": "Mumbai",
    "state": "Maharashtra",
    "postal_code": "400001",
    "country": "India",
    "is_active": true,
    "created_at": "2026-02-17T10:00:00Z",
    "updated_at": "2026-02-17T10:00:00Z"
  }
}
```

### 2. Get All Banks
- **Method:** `GET`
- **Endpoint:** `/banks?limit=10&offset=0`
- **Authentication:** None (Public)
- **Query Parameters:**
  - `limit` (optional): Number of records to return (default: 10)
  - `offset` (optional): Number of records to skip (default: 0)
- **Response:**
```json
{
  "success": true,
  "data": {
    "banks": [
      {
        "id": 1,
        "bank_name": "HDFC Bank",
        "ifsc_code": "HDFC0001234",
        "branch_name": "Mumbai Main Branch",
        "is_active": true,
        "created_at": "2026-02-17T10:00:00Z",
        "updated_at": "2026-02-17T10:00:00Z"
      },
      {
        "id": 2,
        "bank_name": "ICICI Bank",
        "ifsc_code": "ICIC0000001",
        "branch_name": "Delhi Main Branch",
        "is_active": true,
        "created_at": "2026-02-17T10:00:00Z",
        "updated_at": "2026-02-17T10:00:00Z"
      }
    ],
    "total": 2
  }
}
```

### 3. Get Bank by ID
- **Method:** `GET`
- **Endpoint:** `/banks/:id`
- **Authentication:** None (Public)
- **Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "bank_name": "HDFC Bank",
    "ifsc_code": "HDFC0001234",
    "branch_name": "Mumbai Main Branch",
    "branch_code": "HDFC123",
    "address": "123 Banking Street",
    "city": "Mumbai",
    "state": "Maharashtra",
    "postal_code": "400001",
    "country": "India",
    "is_active": true,
    "created_at": "2026-02-17T10:00:00Z",
    "updated_at": "2026-02-17T10:00:00Z"
  }
}
```

### 4. Update Bank
- **Method:** `PUT`
- **Endpoint:** `/banks/:id`
- **Authentication:** Bearer Token + SuperAdmin Role (Required)
- **Description:** Update an existing bank master record. All fields are optional.
- **Request Body:**
```json
{
  "bank_name": "HDFC Bank Limited",
  "branch_name": "Mumbai New Branch",
  "city": "Mumbai",
  "is_active": true
}
```
- **Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "bank_name": "HDFC Bank Limited",
    "ifsc_code": "HDFC0001234",
    "branch_name": "Mumbai New Branch",
    "branch_code": "HDFC123",
    "address": "123 Banking Street",
    "city": "Mumbai",
    "state": "Maharashtra",
    "postal_code": "400001",
    "country": "India",
    "is_active": true,
    "created_at": "2026-02-17T10:00:00Z",
    "updated_at": "2026-02-17T10:30:00Z"
  }
}
```

### 5. Delete Bank
- **Method:** `DELETE`
- **Endpoint:** `/banks/:id`
- **Authentication:** Bearer Token + SuperAdmin Role (Required)
- **Response:**
```json
{
  "success": true,
  "message": "Bank deleted successfully"
}
```

---

## Companies

### 1. Complete Company Setup (Step-by-Step)

#### Step 1a: Create Basic Company (No Dependencies)
- **Method:** `POST`
- **Endpoint:** `/companies`
- **Authentication:** Bearer Token (Required)
- **Description:** Start with basic company information. Bank details will be added after company creation.
- **Request Body:**
```json
{
  "company_name": "ABC Corporation Pvt. Ltd.",
  "display_name": "ABC Corp",
  "gstin": "18AABCT1234H1Z0",
  "pan": "AAACT1234H",
  "business_type_id": 1,
  "email": "info@abccorp.com",
  "phone": "9876543210",
  "phone_code": "+91"
}
```
- **Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "company_name": "ABC Corporation Pvt. Ltd.",
    "gstin": "18AABCT1234H1Z0",
    "created_at": "2026-02-07T10:00:00Z"
  },
  "message": "Company created successfully"
}
```

---

#### Step 1b: Create/Get Available Banks
- **Create Bank (if not exists):**
  - **Method:** `POST`
  - **Endpoint:** `/banks`
  - **Authentication:** Bearer Token + SuperAdmin Role (Required)
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
  - **Response:** Bank object with `id` (save this for Step 1e)

- **Or Get Existing Banks:**
  - **Method:** `GET`
  - **Endpoint:** `/banks`
  - **Authentication:** Not Required
  - **Description:** Retrieve list of all available banks in the system
  - **Response Example:**
  ```json
  {
    "success": true,
    "data": {
      "banks": [
        {
          "id": 1,
          "bank_name": "HDFC Bank",
          "ifsc_code": "HDFC0001234",
          "branch_name": "Mumbai Main Branch",
          "is_active": true
        },
        {
          "id": 2,
          "bank_name": "ICICI Bank",
          "ifsc_code": "ICIC0000001",
          "branch_name": "Delhi Main Branch",
          "is_active": true
        }
      ],
      "total": 2
    }
  }
  ```

---

#### Step 1c: Add Company Contact
- **Method:** `PUT`
- **Endpoint:** `/companies/:id/contact`
- **Authentication:** Bearer Token (Required)
- **Request Body:**
```json
{
  "contact_person": "Rajesh Kumar",
  "email": "rajesh@abccorp.com",
  "phone": "9876543210",
  "phone_code": "+91"
}
```
- **Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "company_id": 1,
    "contact_person": "Rajesh Kumar",
    "email": "rajesh@abccorp.com"
  },
  "message": "Contact added successfully"
}
```

---

#### Step 1d: Add Company Address
- **Method:** `PUT`
- **Endpoint:** `/companies/:id/address`
- **Authentication:** Bearer Token (Required)
- **Request Body:**
```json
{
  "street": "123 Business Street",
  "address_line2": "Floor 5, Building A",
  "city": "Mumbai",
  "state": "Maharashtra",
  "country": "India",
  "postal_code": "400001"
}
```
- **Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "company_id": 1,
    "street": "123 Business Street",
    "city": "Mumbai"
  },
  "message": "Address added successfully"
}
```

---

#### Step 1e: Add Bank Detail to Company
- **Method:** `POST`
- **Endpoint:** `/companies/:id/bank-details`
- **Authentication:** Bearer Token (Required)
- **Description:** Link a bank account to the company. The bank must already exist in the database (from Step 1b). Provide the bank_id that corresponds to your chosen bank.
- **Request Body:**
```json
{
  "bank_id": 1,
  "account_holder_name": "ABC Corporation",
  "account_number": "1234567890",
  "is_primary": true
}
```
- **Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "company_id": 1,
    "bank_id": 1,
    "account_holder_name": "ABC Corporation",
    "account_number": "1234567890",
    "is_primary": true,
    "is_active": true,
    "created_at": "2026-02-07T10:00:00Z"
  },
  "message": "Bank detail created successfully"
}
```

---

#### Step 1f: Add UPI Detail (Optional)
- **Method:** `PUT`
- **Endpoint:** `/companies/:id/upi-details`
- **Authentication:** Bearer Token (Required)
- **Request Body:**
```json
{
  "upi_id": "abc@hdfc"
}
```
- **Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "company_id": 1,
    "upi_id": "abc@hdfc"
  },
  "message": "UPI detail updated successfully"
}
```

---

#### Step 1g: Configure Tax Settings (Optional)
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
- **Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "company_id": 1,
    "gst_enabled": true,
    "tax_type_id": 1
  },
  "message": "Tax settings updated successfully"
}
```

---

#### Step 1h: Configure Invoice Settings (Optional)
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
- **Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "company_id": 1,
    "invoice_prefix": "INV",
    "invoice_start_number": 1
  },
  "message": "Invoice settings updated successfully"
}
```

---

#### Step 1i: Configure Regional Settings (Optional)
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
- **Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "company_id": 1,
    "default_currency": "INR",
    "timezone": "Asia/Kolkata"
  },
  "message": "Regional settings updated successfully"
}
```

---

### Complete Company Setup Workflow Summary

```
1. POST /companies ─── Create basic company
                  │
                  ├─→ 2. PUT /companies/:id/contact ─── Add contact info
                  │
                  ├─→ 3. PUT /companies/:id/address ─── Add address
                  │
                  ├─→ 4. POST /banks (if needed) ─── Create bank master data
                  │         │
                  │         └─→ 5. POST /companies/:id/bank-details ─── Link bank to company
                  │
                  ├─→ 6. PUT /companies/:id/upi-details ─── Add UPI (optional)
                  │
                  ├─→ 7. PUT /companies/:id/tax-settings ─── Configure taxes (optional)
                  │
                  ├─→ 8. PUT /companies/:id/invoice-settings ─── Configure invoicing (optional)
                  │
                  └─→ 9. PUT /companies/:id/regional-settings ─── Configure regional (optional)
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
- **Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "company_id": 1,
    "bank_id": 1,
    "account_holder_name": "ABC Corporation",
    "account_number": "1234567890",
    "is_primary": true,
    "is_active": true,
    "created_at": "2026-02-17T10:00:00Z",
    "updated_at": "2026-02-17T10:00:00Z"
  },
  "message": "Bank detail created successfully"
}
```

### 12. Get Bank Details
- **Method:** `GET`
- **Endpoint:** `/companies/:id/bank-details`
- **Authentication:** Bearer Token (Required)
- **Response:** Returns all bank details linked to the company

### 13. Update Bank Detail
- **Method:** `PUT`
- **Endpoint:** `/companies/bank-details/:id`
- **Authentication:** Bearer Token (Required)
- **Request Body:**
```json
{
  "bank_id": 1,
  "account_holder_name": "ABC Corporation",
  "account_number": "9876543210",
  "is_primary": true,
  "is_active": true
}
```

### 14. Delete Bank Detail
- **Method:** `DELETE`
- **Endpoint:** `/companies/bank-details/:id`
- **Authentication:** Bearer Token (Required)

### 15. Upsert UPI Detail
- **Method:** `PUT`
- **Endpoint:** `/companies/:id/upi-details`
- **Authentication:** Bearer Token (Required)
- **Request Body:**
```json
{
  "upi_id": "abc@hdfc"
}
```

### 16. Get UPI Detail
- **Method:** `GET`
- **Endpoint:** `/companies/:id/upi-details`
- **Authentication:** Bearer Token (Required)

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
- **Authentication:** Bearer Token (Required)

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
- **Authentication:** Bearer Token (Required)

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
- **Authentication:** Bearer Token (Required)

---

## Invoices

### 1. Create Invoice (Step-by-Step)

#### Prerequisites:
- Customer must exist (create via `POST /customers`)
- Item(s) must exist (create via `POST /items`)
- Salesperson must exist (create via `POST /salespersons`) - Optional
- Tax configuration must exist (create via `POST /taxes`)
- Company must exist (create via `POST /companies`)
- (Optional) Sales Order must exist (create via `POST /sales-orders`)

#### Step 1: Create Invoice
- **Method:** `POST`
- **Endpoint:** `/invoices`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Description:** Create basic invoice without line items first
- **Request Body:**
```json
{
  "customer_id": 1,
  "company_id": 1,
  "invoice_number": "INV-2026-001",
  "reference_number": "REF-001",
  "invoice_date": "2026-02-07",
  "due_date": "2026-03-07",
  "terms": "net_30",
  "subject": "Invoice for IT Services & Products",
  "salesperson_id": 1
}
```
- **Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "id": "INV-2026-001",
    "invoice_number": "INV-2026-001",
    "customer_id": 1,
    "invoice_date": "2026-02-07",
    "due_date": "2026-03-07",
    "total": 0,
    "status": "draft",
    "created_at": "2026-02-07T10:00:00Z"
  },
  "message": "Invoice created successfully"
}
```

#### Step 2: Add Line Items to Invoice
- **Method:** `POST`
- **Endpoint:** `/invoices/:invoice_id/line-items`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Request Body (Water Bottle Sales):**
```json
{
  "item_id": "wtr-bot-500ml",
  "quantity": 2000,
  "rate": 15,
  "description": "500ml Packaged Drinking Water Bottles (2000 units)",
  "tax_id": 1
}
```
- **Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "invoice_id": "INV-2026-001",
    "item_id": "wtr-bot-500ml",
    "quantity": 2000,
    "rate": 15,
    "amount": 30000,
    "tax_amount": 1500,
    "created_at": "2026-02-07T10:00:00Z"
  },
  "message": "Line item added successfully"
}
```

#### Step 3: Add Multiple Line Items (Repeat Step 2)
```json
{
  "item_id": "wtr-bot-20l",
  "quantity": 100,
  "rate": 150,
  "description": "20 Litre Water Cooler Bottles - Bulk Supply (100 units)",
  "tax_id": 1
}
```

**Alternative Line Item (Packaging Refund/Deposit):**
```json
{
  "item_id": "bottle-deposit",
  "quantity": 100,
  "rate": -50,
  "description": "Refund for returned 20L bottles (100 units @ ₹50 each)",
  "tax_id": 1
}
```

#### Step 4: Set Tax, Shipping & Discount (Optional)
- **Method:** `PUT`
- **Endpoint:** `/invoices/:invoice_id`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Request Body:**
```json
{
  "shipping_charges": 500,
  "discount": 500,
  "tax_id": 1,
  "adjustment": 100,
  "customer_notes": "Thank you for your business. Please remit payment by due date.",
  "terms_and_conditions": "Payment terms: Net 30 days from invoice date. Late payment interest: 2% per month.",
  "delivery_address": "123 Business Park, Mumbai",
  "is_inclusive_tax": false
}
```
- **Response:**
```json
{
  "success": true,
  "data": {
    "id": "INV-2026-001",
    "sub_total": 22500,
    "shipping_charges": 500,
    "discount": 500,
    "tax_amount": 4050,
    "adjustment": 100,
    "total": 26650,
    "status": "draft"
  },
  "message": "Invoice updated successfully"
}
```

#### Step 5: Submit Invoice
- **Method:** `PATCH`
- **Endpoint:** `/invoices/:invoice_id/status`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Request Body:**
```json
{
  "status": "sent"
}
```
- **Response:**
```json
{
  "success": true,
  "data": {
    "id": "INV-2026-001",
    "status": "sent",
    "updated_at": "2026-02-07T10:30:00Z"
  },
  "message": "Invoice status updated to sent"
}
```

#### Step 6: Record Payment (Optional)
- **Method:** `POST`
- **Endpoint:** `/payments`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Request Body:**
```json
{
  "invoice_id": "INV-2026-001",
  "amount": 26650,
  "payment_date": "2026-02-10",
  "payment_method": "bank_transfer",
  "reference_number": "TXN-12345",
  "bank_id": 1,
  "transaction_id": "NEFT-2026-00123",
  "notes": "Payment received for Invoice INV-2026-001"
}
```
- **Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "invoice_id": "INV-2026-001",
    "amount": 26650,
    "payment_date": "2026-02-10",
    "payment_method": "bank_transfer",
    "reference_number": "TXN-12345",
    "status": "completed",
    "created_at": "2026-02-10T09:00:00Z"
  },
  "message": "Payment created successfully"
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

### 1. Create Bill (Step-by-Step)

#### Prerequisites:
- Vendor must exist (create via `POST /vendors`)
- Item(s) must exist (create via `POST /items`)
- Tax configuration must exist (create via `POST /taxes`)
- Company must exist (create via `POST /companies`)
- (Optional) Purchase Order must exist (create via `POST /purchase-orders`)

#### Step 1: Create Bill
- **Method:** `POST`
- **Endpoint:** `/bills`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Description:** Create basic bill without line items first
- **Request Body:**
```json
{
  "vendor_id": 1,
  "company_id": 1,
  "bill_number": "BILL-2026-001",
  "reference_number": "REF-BILL-001",
  "billing_address": "123 Vendor Street, Bangalore",
  "po_number": "PO-12345",
  "bill_date": "2026-02-07",
  "due_date": "2026-03-07",
  "payment_terms": "net_30",
  "subject": "Bill for Raw Materials Supply"
}
```
- **Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "id": "BILL-2026-001",
    "bill_number": "BILL-2026-001",
    "vendor_id": 1,
    "bill_date": "2026-02-07",
    "due_date": "2026-03-07",
    "total": 0,
    "status": "draft",
    "created_at": "2026-02-07T10:00:00Z"
  },
  "message": "Bill created successfully"
}
```

#### Step 2: Add Line Items to Bill
- **Method:** `POST`
- **Endpoint:** `/bills/:bill_id/line-items`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Request Body:**
```json
{
  "item_id": "ITEM-001",
  "quantity": 100,
  "rate": 150,
  "description": "High-quality plastic raw materials",
  "account": "Cost of Materials",
  "tax_id": 1
}
```
- **Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "bill_id": "BILL-2026-001",
    "item_id": "ITEM-001",
    "quantity": 100,
    "rate": 150,
    "amount": 15000,
    "created_at": "2026-02-07T10:00:00Z"
  },
  "message": "Line item added successfully"
}
```

#### Step 3: Add Multiple Line Items (Repeat Step 2)
```json
{
  "item_id": "ITEM-002",
  "quantity": 50,
  "rate": 200,
  "description": "Metal components for assembly",
  "account": "Cost of Materials",
  "tax_id": 1
}
```

#### Step 4: Set Tax, Shipping & Discount (Optional)
- **Method:** `PUT`
- **Endpoint:** `/bills/:bill_id`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Request Body:**
```json
{
  "shipping_charges": 2000,
  "discount": 1000,
  "tax_id": 1,
  "adjustment": 500,
  "notes": "Thank you for your supply. Please note the delivery was on schedule."
}
```
- **Response:**
```json
{
  "success": true,
  "data": {
    "id": "BILL-2026-001",
    "sub_total": 25000,
    "shipping_charges": 2000,
    "discount": 1000,
    "tax_amount": 4500,
    "adjustment": 500,
    "total": 31000,
    "status": "draft"
  },
  "message": "Bill updated successfully"
}
```

#### Step 5: Submit Bill
- **Method:** `PATCH`
- **Endpoint:** `/bills/:bill_id/status`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Request Body:**
```json
{
  "status": "sent"
}
```
- **Response:**
```json
{
  "success": true,
  "data": {
    "id": "BILL-2026-001",
    "status": "sent",
    "updated_at": "2026-02-07T10:30:00Z"
  },
  "message": "Bill status updated to sent"
}
```

#### Step 6: Record Payment
- **Method:** `POST`
- **Endpoint:** `/payments`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Request Body:**
```json
{
  "bill_id": "BILL-2026-001",
  "amount": 31000,
  "payment_date": "2026-02-12",
  "payment_method": "bank_transfer",
  "reference_number": "TXN-12345",
  "bank_id": 1,
  "transaction_id": "NEFT-2026-00123",
  "notes": "Payment made for Bill BILL-2026-001"
}
```
- **Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "bill_id": "BILL-2026-001",
    "amount": 31000,
    "payment_date": "2026-02-12",
    "payment_method": "bank_transfer",
    "reference_number": "TXN-12345",
    "status": "completed",
    "created_at": "2026-02-12T09:00:00Z"
  },
  "message": "Payment created successfully"
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
  "bill_id": null,
  "amount": 5000,
  "payment_date": "2026-02-07",
  "payment_method": "bank_transfer",
  "reference_number": "TXN-12345",
  "bank_id": 1,
  "cheque_number": "CHQ-12345",
  "cheque_date": "2026-02-07",
  "transaction_id": "NEFT-2026-00123",
  "notes": "Payment received - Invoice INV-123456",
  "description": "Customer payment for February supplies"
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
  "rate": 9,
  "description": "Central Goods and Services Tax at 9% rate",
  "tax_code": "CGST-9",
  "is_compound": false,
  "is_active": true,
  "effective_from": "2026-01-01",
  "priority": 1
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
  "first_name": "Raj",
  "last_name": "Kumar",
  "salutation": "Mr.",
  "email": "raj.kumar@example.com",
  "phone": "9876543210",
  "phone_code": "+91",
  "mobile": "9876543210",
  "mobile_code": "+91",
  "department": "Sales",
  "designation": "Senior Sales Executive",
  "employee_id": "EMP-001",
  "target_revenue": 1000000,
  "target_currency": "INR",
  "commission_percentage": 5,
  "reporting_manager_id": 2,
  "is_active": true,
  "date_of_joining": "2020-06-15"
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

### 1. Create Purchase Order (Step-by-Step)

#### Prerequisites:
- Vendor must exist (create via `POST /vendors`)
- Item(s) must exist (create via `POST /items`)
- Tax configuration must exist (create via `POST /taxes`)
- Company must exist (create via `POST /companies`)

#### Step 1: Create Purchase Order
- **Method:** `POST`
- **Endpoint:** `/purchase-orders`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Description:** Create basic PO without line items first
- **Request Body:**
```json
{
  "vendor_id": 1,
  "company_id": 1,
  "po_number": "PO-2026-001",
  "reference_no": "REF-001",
  "po_date": "2026-02-07",
  "expected_delivery_date": "2026-02-14",
  "delivery_address": "123 Warehouse Street, Bangalore",
  "payment_terms": "net_30",
  "incoterms": "FOB",
  "shipping_via": "Ground Transport",
  "notes": "Please ensure timely delivery as per schedule."
}
```
- **Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "id": "PO-2026-001",
    "purchase_order_no": "PO-2026-001",
    "vendor_id": 1,
    "po_date": "2026-02-07",
    "expected_delivery_date": "2026-02-14",
    "total": 0,
    "status": "draft",
    "created_at": "2026-02-07T10:00:00Z"
  },
  "message": "Purchase Order created successfully"
}
```

#### Step 2: Add Line Items to Purchase Order
- **Method:** `POST`
- **Endpoint:** `/purchase-orders/:po_id/line-items`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Request Body (Water Bottles):**
```json
{
  "item_id": "wtr-bot-500ml",
  "quantity": 5000,
  "rate": 8,
  "description": "500ml PET drinking water bottles (5000 units)",
  "tax_id": 1,
  "warehouse_id": 1
}
```
- **Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "purchase_order_id": "PO-2026-001",
    "item_id": "wtr-bot-500ml",
    "quantity": 5000,
    "rate": 8,
    "amount": 40000,
    "tax_amount": 2000,
    "created_at": "2026-02-07T10:00:00Z"
  },
  "message": "Line item added successfully"
}
```

#### Step 3: Add Multiple Line Items (Repeat Step 2)
```json
{
  "item_id": "cap-100-tam",
  "quantity": 50,
  "rate": 25,
  "description": "Tamper-proof caps (50 packs × 100 pieces)",
  "tax_id": 1,
  "warehouse_id": 1
}
```

**Alternative Step 3 (Water Cooler Bottles):**
```json
{
  "item_id": "wtr-bot-20l",
  "quantity": 100,
  "rate": 80,
  "description": "20 litre polycarbonate water cooler bottles (100 units)",
  "tax_id": 1,
  "warehouse_id": 1
}
```

#### Step 4: Set Tax, Shipping & Adjustment (Optional)
- **Method:** `PUT`
- **Endpoint:** `/purchase-orders/:po_id`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Request Body:**
```json
{
  "shipping_charges": 500,
  "discount": 500,
  "tax_id": 1,
  "adjustment": 100
}
```
- **Response:**
```json
{
  "success": true,
  "data": {
    "id": "PO-2026-001",
    "sub_total": 25000,
    "shipping_charges": 500,
    "discount": 500,
    "tax_amount": 4500,
    "adjustment": 100,
    "total": 29600,
    "status": "draft"
  },
  "message": "Purchase Order updated successfully"
}
```

#### Step 5: Submit PO
- **Method:** `PATCH`
- **Endpoint:** `/purchase-orders/:po_id/status`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Request Body:**
```json
{
  "status": "sent"
}
```
- **Response:**
```json
{
  "success": true,
  "data": {
    "id": "PO-2026-001",
    "status": "sent",
    "updated_at": "2026-02-07T10:30:00Z"
  },
  "message": "Purchase Order status updated to sent"
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

### 1. Create Sales Order (Step-by-Step)

#### Prerequisites:
- Customer must exist (create via `POST /customers`)
- Item(s) must exist (create via `POST /items`)
- Salesperson must exist (create via `POST /salespersons`)
- Tax configuration must exist (create via `POST /taxes`)
- Company must exist (create via `POST /companies`)

#### Step 1: Create Sales Order
- **Method:** `POST`
- **Endpoint:** `/sales-orders`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Description:** Create basic SO without line items first
- **Request Body:**
```json
{
  "customer_id": 1,
  "company_id": 1,
  "so_number": "SO-2026-001",
  "reference_no": "CUST-REF-001",
  "sales_order_date": "2026-02-07",
  "expected_shipment_date": "2026-02-14",
  "delivery_address": "789 Customer Plaza, New Delhi",
  "payment_terms": "net_30",
  "delivery_method": "courier",
  "courier_company": "FedEx",
  "salesperson_id": 1
}
```
- **Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "id": "SO-2026-001",
    "sales_order_no": "SO-2026-001",
    "customer_id": 1,
    "sales_order_date": "2026-02-07",
    "expected_shipment_date": "2026-02-14",
    "total": 0,
    "status": "draft",
    "created_at": "2026-02-07T10:00:00Z"
  },
  "message": "Sales Order created successfully"
}
```

#### Step 2: Add Line Items to Sales Order
- **Method:** `POST`
- **Endpoint:** `/sales-orders/:so_id/line-items`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Request Body (Water Bottles):**
```json
{
  "item_id": "wtr-bot-500ml",
  "quantity": 1000,
  "rate": 15,
  "description": "500ml PET Drinking Water Bottles (1000 units)",
  "tax_id": 1,
  "warehouse_id": 1
}
```
- **Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "sales_order_id": "SO-2026-001",
    "item_id": "wtr-bot-500ml",
    "quantity": 1000,
    "rate": 15,
    "amount": 15000,
    "tax_amount": 750,
    "created_at": "2026-02-07T10:00:00Z"
  },
  "message": "Line item added successfully"
}
```

#### Step 3: Add Multiple Line Items (Repeat Step 2)
```json
{
  "item_id": "wtr-bot-20l",
  "quantity": 50,
  "rate": 150,
  "description": "20 Litre Water Cooler Bottles (50 units)",
  "tax_id": 1,
  "warehouse_id": 1
}
```

**Alternative Option (Delivery with Bottles):**
```json
{
  "item_id": "cap-100-tam",
  "quantity": 10,
  "rate": 45,
  "description": "Tamper-proof caps - 10 packs (1000 total)",
  "tax_id": 1,
  "warehouse_id": 1
}
```

#### Step 4: Set Tax, Shipping & Adjustment (Optional)
- **Method:** `PUT`
- **Endpoint:** `/sales-orders/:so_id`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Request Body:**
```json
{
  "shipping_charges": 500,
  "discount": 500,
  "tax_id": 1,
  "adjustment": 100,
  "customer_notes": "Handle with care. Fragile items. Please call before delivery."
}
```
- **Response:**
```json
{
  "success": true,
  "data": {
    "id": "SO-2026-001",
    "sub_total": 62500,
    "shipping_charges": 500,
    "discount": 500,
    "tax_amount": 11250,
    "adjustment": 100,
    "total": 73850,
    "status": "draft"
  },
  "message": "Sales Order updated successfully"
}
```

#### Step 5: Reserve Inventory (Optional)
- **Method:** `PUT`
- **Endpoint:** `/sales-orders/:so_id/reserve-inventory`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Response:**
```json
{
  "success": true,
  "data": {
    "id": "SO-2026-001",
    "inventory_reserved": true,
    "reserved_date": "2026-02-07T10:30:00Z"
  },
  "message": "Inventory reserved successfully"
}
```

#### Step 6: Confirm and Submit SO
- **Method:** `PATCH`
- **Endpoint:** `/sales-orders/:so_id/status`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Request Body:**
```json
{
  "status": "confirmed"
}
```
- **Response:**
```json
{
  "success": true,
  "data": {
    "id": "SO-2026-001",
    "status": "confirmed",
    "updated_at": "2026-02-07T10:45:00Z"
  },
  "message": "Sales Order status updated to confirmed"
}
```

#### Step 7: Create Package from SO
- **Method:** `POST`
- **Endpoint:** `/packages`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Request Body:**
```json
{
  "sales_order_id": "SO-2026-001",
  "customer_id": 1
}
```
- **Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "id": "PKG-2026-001",
    "package_slip_no": "PKG-2026-001",
    "sales_order_id": "SO-2026-001",
    "customer_id": 1,
    "status": "created",
    "created_at": "2026-02-07T11:00:00Z"
  },
  "message": "Package created successfully"
}
```

#### Step 8: Create Shipment from Package
- **Method:** `POST`
- **Endpoint:** `/shipments`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Request Body:**
```json
{
  "package_id": "PKG-2026-001",
  "sales_order_id": "SO-2026-001",
  "customer_id": 1,
  "shipment_date": "2026-02-07",
  "expected_delivery": "2026-02-14",
  "carrier": "FedEx",
  "carrier_service_type": "Express",
  "tracking_number": "FEDEX123456"
}
```
- **Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "id": "SHIP-2026-001",
    "shipment_no": "SHIP-2026-001",
    "package_id": "PKG-2026-001",
    "carrier": "FedEx",
    "tracking_number": "FEDEX123456",
    "status": "created",
    "created_at": "2026-02-07T11:15:00Z"
  },
  "message": "Shipment created successfully"
}
```

#### Step 9: Create Invoice from SO
- **Method:** `POST`
- **Endpoint:** `/invoices`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Request Body:**
```json
{
  "customer_id": 1,
  "company_id": 1,
  "sales_order_id": "SO-2026-001",
  "invoice_date": "2026-02-07",
  "due_date": "2026-03-07",
  "terms": "net_30",
  "subject": "Invoice for Order SO-2026-001"
}
```
- **Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "id": "INV-2026-001",
    "invoice_number": "INV-2026-001",
    "customer_id": 1,
    "sales_order_id": "SO-2026-001",
    "total": 73850,
    "status": "draft",
    "created_at": "2026-02-07T11:30:00Z"
  },
  "message": "Invoice created successfully"
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
  "package_number": "PKG-2026-001",
  "package_type": "box",
  "package_items": [
    {
      "item_id": "ITEM-001",
      "quantity": 10,
      "item_description": "Dell Inspiron 15 Laptops",
      "sku": "LAP-DELL-001"
    },
    {
      "item_id": "ITEM-002",
      "quantity": 5,
      "item_description": "Monitor Accessories",
      "sku": "ACC-MON-001"
    }
  ],
  "package_weight": 25,
  "weight_unit": "kg",
  "package_dimensions": "50x40x30",
  "dimension_unit": "cm",
  "volume": 60000,
  "packaging_material": "cardboard",
  "fragile": true,
  "contents_description": "Electronic equipment - Handle with care",
  "package_status": "created"
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

### Water Company Shipment Examples

### 1. Create Shipment
- **Method:** `POST`
- **Endpoint:** `/shipments`
- **Authentication:** Bearer Token + Admin Role (Required)
- **Request Body (500ml Water Bottles):**
```json
{
  "package_id": "PKG-2026-WATER-001",
  "sales_order_id": "SO-2026-001",
  "customer_id": 15,
  "shipment_number": "SHIP-2026-W001",
  "shipment_date": "2026-02-07",
  "expected_delivery": "2026-02-09",
  "actual_delivery": null,
  "carrier": "Blue Dart Express",
  "carrier_service_type": "Standard",
  "tracking_number": "BDE2026W001450",
  "shipping_address": "456 Market Street, Building B, Bangalore, Karnataka 560001",
  "shipping_address_full": {
    "attention": "Retail Manager",
    "address_line1": "456 Market Street",
    "address_line2": "Building B",
    "city": "Bangalore",
    "state": "Karnataka",
    "country": "India",
    "postal_code": "560001",
    "phone": "9123456789"
  },
  "weight": 30,
  "weight_unit": "kg",
  "shipment_cost": 800,
  "insurance_amount": 50000,
  "insured": true,
  "shipping_mode": "ground",
  "shipment_status": "created",
  "notes": "Water bottles - Handle with care. Check for damages upon receipt. Maintain cool storage."
}
```

**Alternative (20L Water Cooler Bottles - Refrigerated/Special Handling):**
```json
{
  "package_id": "PKG-2026-COOLER-001",
  "sales_order_id": "SO-2026-002",
  "customer_id": 18,
  "shipment_number": "SHIP-2026-W002",
  "shipment_date": "2026-02-07",
  "expected_delivery": "2026-02-08",
  "carrier": "Local Logistics - AquaDeliver",
  "carrier_service_type": "Express",
  "tracking_number": "AQA2026002500",
  "shipping_address": "Office Complex, 789 Business Park, Hyderabad, Telangana 500001",
  "shipping_address_full": {
    "attention": "Admin Department",
    "address_line1": "Office Complex, 789 Business Park",
    "city": "Hyderabad",
    "state": "Telangana",
    "country": "India",
    "postal_code": "500001",
    "phone": "9876543215"
  },
  "weight": 800,
  "weight_unit": "kg",
  "shipment_cost": 2500,
  "insurance_amount": 100000,
  "insured": true,
  "shipping_mode": "ground",
  "shipment_status": "created",
  "notes": "20L bulk water bottles - Refrigerated transport required. Handle upright only. Delivery between 10 AM - 4 PM preferred. Call 1 hour before arrival."
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
  "phone": "9876543210",
  "phone_code": "+91",
  "company_name": "ABC Corporation",
  "subject": "Technical Issue - Purchase Order Not Processing",
  "category": "technical",
  "description": "I'm facing an issue with the purchase order module. When I try to create a PO, the system shows an error preventing the submission.",
  "priority": "high",
  "attachments": [
    {
      "filename": "error_screenshot.png",
      "file_url": "https://example.com/uploads/error_screenshot.png"
    }
  ],
  "expected_resolution_date": "2026-02-09",
  "ticket_status": "open"
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
