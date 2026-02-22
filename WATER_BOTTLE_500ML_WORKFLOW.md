# 500ML Water Bottle Business - Complete Workflow (Zoho Inventory Style)

## Overview
This document provides a complete step-by-step workflow for managing a 500ml water bottle manufacturing and distribution business using a Zoho Inventory-like system. The workflow includes purchasing raw materials, manufacturing finished products, and selling to customers.

---

## STEP 1: SET UP YOUR CONTACTS

### 1.1 Create Vendors (Suppliers)

#### Vendor 1: Plastic Bottle Manufacturer

**Request:**
```json
POST /vendors
Content-Type: application/json
Authorization: Bearer {jwt_token}

{
  "salutation": "Mr.",
  "first_name": "Rajesh",
  "last_name": "Kumar",
  "company_name": "Premier Plastic Industries Pvt Ltd",
  "display_name": "Premier Plastic",
  "email_address": "contact@premierplastic.com",
  "work_phone": "1141234567",
  "work_phone_code": "+91",
  "mobile": "9876543210",
  "mobile_code": "+91",
  "vendor_language": "English",
  "other_details": {
    "pan": "AABPN1234H",
    "is_msme_registered": true,
    "currency": "INR",
    "payment_terms": "Net 30",
    "tds": "yes",
    "enable_portal": true,
    "website_url": "https://premierplastic.com",
    "department": "Sales",
    "designation": "Sales Manager",
    "twitter": "@premierplastic",
    "skype_name": "premierplastic.sales",
    "facebook": "https://facebook.com/premierplastic"
  },
  "billing_address": {
    "attention": "Mr. Rajesh Kumar",
    "country_region": "India",
    "address_line1": "Plot No. 123, Industrial Area",
    "address_line2": "MIDC Zone",
    "city": "Nashik",
    "state": "Maharashtra",
    "pin_code": "422212",
    "phone": "1141234567",
    "phone_code": "+91",
    "fax_number": "1141234568"
  },
  "shipping_address": {
    "attention": "Mr. Rajesh Kumar",
    "country_region": "India",
    "address_line1": "Plot No. 123, Industrial Area",
    "address_line2": "MIDC Zone",
    "city": "Nashik",
    "state": "Maharashtra",
    "pin_code": "422212",
    "phone": "1141234567",
    "phone_code": "+91",
    "fax_number": "1141234568"
  },
  "contact_persons": [
    {
      "salutation": "Mr.",
      "first_name": "Rajesh",
      "last_name": "Kumar",
      "email_address": "rajesh@premierplastic.com",
      "work_phone": "1141234567",
      "work_phone_code": "+91",
      "mobile": "9876543210",
      "mobile_code": "+91"
    }
  ],
  "bank_details": [
    {
      "bank_id": 1,
      "account_holder_name": "Premier Plastic Industries Pvt Ltd",
      "account_number": "12345678901234",
      "reenter_account_number": "12345678901234",
      "ifsc_code": "SBIN0001122",
      "branch_name": "Nashik Main Branch"
    }
  ]
}
```

**Response (Success - 201 Created):**
```json
{
  "success": true,
  "message": "Vendor created successfully",
  "data": {
    "vendor_id": 1,
    "company_name": "Premier Plastic Industries Pvt Ltd",
    "display_name": "Premier Plastic",
    "email_address": "contact@premierplastic.com",
    "status": "active",
    "created_at": "2024-02-22T10:00:00Z"
  }
}
```

---

#### Vendor 2: Label Printing Company

**Request:**
```json
POST /vendors
Content-Type: application/json
Authorization: Bearer {jwt_token}

{
  "salutation": "Ms.",
  "first_name": "Priya",
  "last_name": "Sharma",
  "company_name": "PrintPak Solutions Pvt Ltd",
  "display_name": "PrintPak Solutions",
  "email_address": "sales@printpak.com",
  "work_phone": "2241567890",
  "work_phone_code": "+91",
  "mobile": "9123456789",
  "mobile_code": "+91",
  "vendor_language": "English",
  "other_details": {
    "pan": "AABPS5678K",
    "is_msme_registered": true,
    "currency": "INR",
    "payment_terms": "Net 45",
    "tds": "yes",
    "enable_portal": true,
    "website_url": "https://printpak.com",
    "department": "Sales",
    "designation": "Sales Executive",
    "twitter": "@printpaksolutions",
    "skype_name": "printpak.sales",
    "facebook": "https://facebook.com/printpaksolutions"
  },
  "billing_address": {
    "attention": "Ms. Priya Sharma",
    "country_region": "India",
    "address_line1": "Building A, Business Park",
    "address_line2": "Print Zone",
    "city": "Pune",
    "state": "Maharashtra",
    "pin_code": "411001",
    "phone": "2241567890",
    "phone_code": "+91",
    "fax_number": "2241567891"
  },
  "shipping_address": {
    "attention": "Warehouse Team",
    "country_region": "India",
    "address_line1": "Warehouse Complex, Industrial Estate",
    "address_line2": "Logistics Hub",
    "city": "Pune",
    "state": "Maharashtra",
    "pin_code": "411041",
    "phone": "2241567892",
    "phone_code": "+91",
    "fax_number": "2241567893"
  },
  "contact_persons": [
    {
      "salutation": "Ms.",
      "first_name": "Priya",
      "last_name": "Sharma",
      "email_address": "priya@printpak.com",
      "work_phone": "2241567890",
      "work_phone_code": "+91",
      "mobile": "9123456789",
      "mobile_code": "+91"
    },
    {
      "salutation": "Mr.",
      "first_name": "Arjun",
      "last_name": "Desai",
      "email_address": "arjun@printpak.com",
      "work_phone": "2241567894",
      "work_phone_code": "+91",
      "mobile": "9123456790",
      "mobile_code": "+91"
    }
  ],
  "bank_details": [
    {
      "bank_id": 2,
      "account_holder_name": "PrintPak Solutions Pvt Ltd",
      "account_number": "98765432109876",
      "reenter_account_number": "98765432109876",
      "ifsc_code": "HDFC0001234",
      "branch_name": "Pune HQ"
    }
  ]
}
```

**Response (Success - 201 Created):**
```json
{
  "success": true,
  "message": "Vendor created successfully",
  "data": {
    "vendor_id": 2,
    "company_name": "PrintPak Solutions Pvt Ltd",
    "display_name": "PrintPak Solutions",
    "email_address": "sales@printpak.com",
    "status": "active",
    "created_at": "2024-02-22T10:05:00Z"
  }
}
```

---

### 1.2 Create Customers (Buyers)

#### Customer 1: Retail Distributor Company

**Request:**
```json
POST /customers
Content-Type: application/json
Authorization: Bearer {jwt_token}

{
  "salutation": "Mr.",
  "first_name": "Arun",
  "last_name": "Verma",
  "company_name": "Fresh Waters Distribution Ltd",
  "display_name": "Fresh Waters",
  "email_address": "procurement@freshwaters.com",
  "work_phone": "1141234567",
  "work_phone_code": "+91",
  "mobile": "9988776655",
  "mobile_code": "+91",
  "customer_language": "English",
  "other_details": {
    "pan": "AAFWV1234D",
    "currency": "INR",
    "payment_terms": "Net 15",
    "enable_portal": true
  },
  "billing_address": {
    "attention": "Mr. Arun Verma",
    "country_region": "India",
    "address_line1": "123 Market Street",
    "address_line2": "Business District",
    "city": "Mumbai",
    "state": "Maharashtra",
    "pin_code": "400001",
    "phone": "1141234567",
    "phone_code": "+91",
    "fax_number": "1141234568"
  },
  "shipping_address": {
    "attention": "Warehouse Manager",
    "country_region": "India",
    "address_line1": "450 Industrial Road",
    "address_line2": "Logistics Zone",
    "city": "Thane",
    "state": "Maharashtra",
    "pin_code": "400604",
    "phone": "9988776655",
    "phone_code": "+91",
    "fax_number": "9988776656"
  },
  "contact_persons": [
    {
      "salutation": "Mr.",
      "first_name": "Arun",
      "last_name": "Verma",
      "email_address": "arun@freshwaters.com",
      "work_phone": "1141234567",
      "work_phone_code": "+91",
      "mobile": "9988776655",
      "mobile_code": "+91"
    },
    {
      "salutation": "Ms.",
      "first_name": "Swati",
      "last_name": "Patel",
      "email_address": "swati@freshwaters.com",
      "work_phone": "1141234569",
      "work_phone_code": "+91",
      "mobile": "9988776657",
      "mobile_code": "+91"
    }
  ]
}
```

**Response (Success - 201 Created):**
```json
{
  "success": true,
  "message": "Customer created successfully",
  "data": {
    "customer_id": 1,
    "company_name": "Fresh Waters Distribution Ltd",
    "display_name": "Fresh Waters",
    "email_address": "procurement@freshwaters.com",
    "status": "active",
    "created_at": "2024-02-22T10:10:00Z"
  }
}
```

---

## STEP 2: BUILD YOUR DIGITAL WAREHOUSE (ITEMS)

### 2.1 Create Raw Material Item 1: 500ml PET Bottle (With Variants)

**Request:**
```json
POST /items
Content-Type: application/json
Authorization: Bearer {jwt_token}

{
  "name": "500ml PET Bottle (Clear)",
  "type": "goods",
  "item_details": {
    "structure": "variants",
    "unit": "pcs",
    "sku": "BOTTLE-500ML-CLR",
    "upc": "8901234500001",
    "ean": "8901234500001",
    "mpn": "BTL500-CLR-001",
    "isbn": "",
    "description": "Clear PET plastic bottle with 500ml capacity, food-grade material, suitable for drinking water and beverages. Available in variants for different flavor lines: Plain, Lemon, and Mint.",
    "attribute_definitions": [
      {
        "key": "Flavor",
        "options": ["Plain", "Lemon", "Mint"]
      }
    ],
    "variants": [
      {
        "sku": "BOTTLE-500ML-PLAIN",
        "attribute_map": {
          "flavor": "Plain"
        },
        "selling_price": 5.00,
        "cost_price": 3.50,
        "stock_quantity": 0
      },
      {
        "sku": "BOTTLE-500ML-LEMON",
        "attribute_map": {
          "flavor": "Lemon"
        },
        "selling_price": 5.00,
        "cost_price": 3.50,
        "stock_quantity": 0
      },
      {
        "sku": "BOTTLE-500ML-MINT",
        "attribute_map": {
          "flavor": "Mint"
        },
        "selling_price": 5.00,
        "cost_price": 3.50,
        "stock_quantity": 0
      }
    ]
  },
  "sales_info": {
    "account": "Raw Materials - Bottles",
    "selling_price": 5.00,
    "currency": "INR",
    "description": "Sold as component to manufacturing partners"
  },
  "purchase_info": {
    "account": "Raw Materials Purchase",
    "cost_price": 3.50,
    "currency": "INR",
    "preferred_vendor_id": 1,
    "description": "Purchased from Premier Plastic Industries in bulk quantities"
  },
  "inventory": {
    "track_inventory": true,
    "inventory_account": "Raw Materials Inventory - Bottles",
    "inventory_valuation_method": "FIFO",
    "reorder_point": 500
  },
  "return_policy": {
    "returnable": false
  }
}
```

**Response (Success - 201 Created):**
```json
{
  "success": true,
  "message": "Item created successfully with variants",
  "data": {
    "item_id": "1",
    "name": "500ml PET Bottle (Clear)",
    "sku": "BOTTLE-500ML-CLR",
    "type": "goods",
    "structure": "variants",
    "purchase_price": 3.50,
    "selling_price": 5.00,
    "track_inventory": true,
    "variants": [
      {
        "variant_id": "1",
        "variant_sku": "BOTTLE-500ML-PLAIN",
        "attribute_map": {"flavor": "Plain"},
        "selling_price": 5.00,
        "cost_price": 3.50
      },
      {
        "variant_id": "2",
        "variant_sku": "BOTTLE-500ML-LEMON",
        "attribute_map": {"flavor": "Lemon"},
        "selling_price": 5.00,
        "cost_price": 3.50
      },
      {
        "variant_id": "3",
        "variant_sku": "BOTTLE-500ML-MINT",
        "attribute_map": {"flavor": "Mint"},
        "selling_price": 5.00,
        "cost_price": 3.50
      }
    ],
    "created_at": "2024-02-22T10:15:00Z"
  }
}
```

---

### 2.2 Create Raw Material Item 2: Flip Cap with Seal (With Variants)

**Request:**
```json
POST /items
Content-Type: application/json
Authorization: Bearer {jwt_token}

{
  "name": "Flip Cap with Seal (28mm)",
  "type": "goods",
  "item_details": {
    "structure": "variants",
    "unit": "pcs",
    "sku": "CAP-FLIP-28MM",
    "upc": "8901234500002",
    "ean": "8901234500002",
    "mpn": "CAP-FLIP-28-001",
    "isbn": "",
    "description": "Flip cap with tamper-proof seal, 28mm neck fitting, food-grade plastic, reusable design. Available in variants with color coding for different flavor variants.",
    "attribute_definitions": [
      {
        "key": "Flavor",
        "options": ["Plain", "Lemon", "Mint"]
      }
    ],
    "variants": [
      {
        "sku": "CAP-FLIP-PLAIN",
        "attribute_map": {
          "flavor": "Plain"
        },
        "selling_price": 2.00,
        "cost_price": 1.20,
        "stock_quantity": 0
      },
      {
        "sku": "CAP-FLIP-LEMON",
        "attribute_map": {
          "flavor": "Lemon"
        },
        "selling_price": 2.00,
        "cost_price": 1.20,
        "stock_quantity": 0
      },
      {
        "sku": "CAP-FLIP-MINT",
        "attribute_map": {
          "flavor": "Mint"
        },
        "selling_price": 2.00,
        "cost_price": 1.20,
        "stock_quantity": 0
      }
    ]
  },
  "sales_info": {
    "account": "Raw Materials - Caps",
    "selling_price": 2.00,
    "currency": "INR",
    "description": "Flip cap component for bottle assembly"
  },
  "purchase_info": {
    "account": "Raw Materials Purchase",
    "cost_price": 1.20,
    "currency": "INR",
    "preferred_vendor_id": 1,
    "description": "Cap supplier - Premier Plastic Industries"
  },
  "inventory": {
    "track_inventory": true,
    "inventory_account": "Raw Materials Inventory - Caps",
    "inventory_valuation_method": "FIFO",
    "reorder_point": 500
  },
  "return_policy": {
    "returnable": false
  }
}
```

**Response (Success - 201 Created):**
```json
{
  "success": true,
  "message": "Item created successfully with variants",
  "data": {
    "item_id": "2",
    "name": "Flip Cap with Seal (28mm)",
    "sku": "CAP-FLIP-28MM",
    "type": "goods",
    "structure": "variants",
    "purchase_price": 1.20,
    "selling_price": 2.00,
    "track_inventory": true,
    "variants": [
      {
        "variant_id": "1",
        "variant_sku": "CAP-FLIP-PLAIN",
        "attribute_map": {"flavor": "Plain"},
        "selling_price": 2.00,
        "cost_price": 1.20
      },
      {
        "variant_id": "2",
        "variant_sku": "CAP-FLIP-LEMON",
        "attribute_map": {"flavor": "Lemon"},
        "selling_price": 2.00,
        "cost_price": 1.20
      },
      {
        "variant_id": "3",
        "variant_sku": "CAP-FLIP-MINT",
        "attribute_map": {"flavor": "Mint"},
        "selling_price": 2.00,
        "cost_price": 1.20
      }
    ],
    "created_at": "2024-02-22T10:20:00Z"
  }
}
```

---

### 2.3 Create Raw Material Item 3: Water Bottle Label (With Variants)

**Request:**
```json
POST /items
Content-Type: application/json
Authorization: Bearer {jwt_token}

{
  "name": "Brand Label - 500ml Water Bottle",
  "type": "goods",
  "item_details": {
    "structure": "variants",
    "unit": "pcs",
    "sku": "LABEL-WATER-500ML",
    "upc": "8901234500003",
    "ean": "8901234500003",
    "mpn": "LBL-WATER-500-001",
    "isbn": "",
    "description": "Printed water bottle label with brand logo, product information, nutrition facts, barcode, and regulatory compliance information. Different variants for Plain, Lemon, and Mint flavors with unique flavor identifiers.",
    "attribute_definitions": [
      {
        "key": "Flavor",
        "options": ["Plain", "Lemon", "Mint"]
      }
    ],
    "variants": [
      {
        "sku": "LABEL-PLAIN",
        "attribute_map": {
          "flavor": "Plain"
        },
        "selling_price": 0.50,
        "cost_price": 0.25,
        "stock_quantity": 0
      },
      {
        "sku": "LABEL-LEMON",
        "attribute_map": {
          "flavor": "Lemon"
        },
        "selling_price": 0.50,
        "cost_price": 0.25,
        "stock_quantity": 0
      },
      {
        "sku": "LABEL-MINT",
        "attribute_map": {
          "flavor": "Mint"
        },
        "selling_price": 0.50,
        "cost_price": 0.25,
        "stock_quantity": 0
      }
    ]
  },
  "sales_info": {
    "account": "Raw Materials - Labels",
    "selling_price": 0.50,
    "currency": "INR",
    "description": "Label component for water bottle packaging"
  },
  "purchase_info": {
    "account": "Raw Materials Purchase",
    "cost_price": 0.25,
    "currency": "INR",
    "preferred_vendor_id": 2,
    "description": "Labels from PrintPak Solutions"
  },
  "inventory": {
    "track_inventory": true,
    "inventory_account": "Raw Materials Inventory - Labels",
    "inventory_valuation_method": "FIFO",
    "reorder_point": 1000
  },
  "return_policy": {
    "returnable": false
  }
}
```

**Response (Success - 201 Created):**
```json
{
  "success": true,
  "message": "Item created successfully with variants",
  "data": {
    "item_id": "3",
    "name": "Brand Label - 500ml Water Bottle",
    "sku": "LABEL-WATER-500ML",
    "type": "goods",
    "structure": "variants",
    "purchase_price": 0.25,
    "selling_price": 0.50,
    "track_inventory": true,
    "variants": [
      {
        "variant_id": "1",
        "variant_sku": "LABEL-PLAIN",
        "attribute_map": {"flavor": "Plain"},
        "selling_price": 0.50,
        "cost_price": 0.25
      },
      {
        "variant_id": "2",
        "variant_sku": "LABEL-LEMON",
        "attribute_map": {"flavor": "Lemon"},
        "selling_price": 0.50,
        "cost_price": 0.25
      },
      {
        "variant_id": "3",
        "variant_sku": "LABEL-MINT",
        "attribute_map": {"flavor": "Mint"},
        "selling_price": 0.50,
        "cost_price": 0.25
      }
    ],
    "created_at": "2024-02-22T10:25:00Z"
  }
}
```

---

### 2.4 Create Finished Product Item: 500ml Purified Water Bottle (With Variants)

**Request:**
```json
POST /items
Content-Type: application/json
Authorization: Bearer {jwt_token}

{
  "name": "500ml Purified Water Bottle",
  "type": "goods",
  "item_details": {
    "structure": "variants",
    "unit": "pcs",
    "sku": "WATER-500ML-PURIFIED",
    "upc": "8901234600001",
    "ean": "8901234600001",
    "mpn": "WATER-500-PURE-001",
    "isbn": "",
    "description": "500ml purified drinking water in clear PET bottle with tamper-proof flip cap and branded label. Multi-stage filtration process ensures safety and purity. Available in multiple flavors. Ready for retail distribution and institutional use.",
    "attribute_definitions": [
      {
        "key": "Flavor",
        "options": ["Plain", "Lemon", "Mint"]
      }
    ],
    "variants": [
      {
        "sku": "WATER-500ML-PLAIN",
        "attribute_map": {
          "flavor": "Plain"
        },
        "selling_price": 20.00,
        "cost_price": 8.95,
        "stock_quantity": 0
      },
      {
        "sku": "WATER-500ML-LEMON",
        "attribute_map": {
          "flavor": "Lemon"
        },
        "selling_price": 22.00,
        "cost_price": 9.95,
        "stock_quantity": 0
      },
      {
        "sku": "WATER-500ML-MINT",
        "attribute_map": {
          "flavor": "Mint"
        },
        "selling_price": 22.00,
        "cost_price": 9.95,
        "stock_quantity": 0
      }
    ]
  },
  "sales_info": {
    "account": "Sales - Finished Goods",
    "selling_price": 20.00,
    "currency": "INR",
    "description": "Retail price for 500ml purified water bottle in bulk orders - Available in Plain, Lemon, and Mint flavors"
  },
  "purchase_info": {
    "account": "Manufacturing Cost",
    "cost_price": 8.95,
    "currency": "INR",
    "description": "Total cost = Bottle(3.50) + Cap(1.20) + Label(0.25) + Purified Water(2.00) + Packaging/Labor(1.00) + Overheads(1.00). Flavored variants add ₹1.00 per unit."
  },
  "inventory": {
    "track_inventory": true,
    "inventory_account": "Finished Goods Inventory",
    "inventory_valuation_method": "FIFO",
    "reorder_point": 200
  },
  "return_policy": {
    "returnable": true
  }
}
```

**Response (Success - 201 Created):**
```json
{
  "success": true,
  "message": "Item created successfully with variants",
  "data": {
    "item_id": "4",
    "name": "500ml Purified Water Bottle",
    "sku": "WATER-500ML-PURIFIED",
    "type": "goods",
    "structure": "variants",
    "purchase_price": 8.95,
    "selling_price": 20.00,
    "track_inventory": true,
    "variants": [
      {
        "variant_id": "1",
        "variant_sku": "WATER-500ML-PLAIN",
        "attribute_map": {
          "flavor": "Plain"
        },
        "selling_price": 20.00,
        "cost_price": 8.95
      },
      {
        "variant_id": "2",
        "variant_sku": "WATER-500ML-LEMON",
        "attribute_map": {
          "flavor": "Lemon"
        },
        "selling_price": 22.00,
        "cost_price": 9.95
      },
      {
        "variant_id": "3",
        "variant_sku": "WATER-500ML-MINT",
        "attribute_map": {
          "flavor": "Mint"
        },
        "selling_price": 22.00,
        "cost_price": 9.95
      }
    ],
    "created_at": "2024-02-22T10:30:00Z"
  }
}
```

---

### 2.5 Set Opening Stock for Raw Materials

#### Set Opening Stock - Bottles

**Request:**
```json
PUT /items/1/opening-stock
Content-Type: application/json
Authorization: Bearer {jwt_token}

{
  "opening_stock": 1000,
  "opening_stock_rate_per_unit": 3.50
}
```

**Response (Success - 200 OK):**
```json
{
  "success": true,
  "message": "Opening stock updated successfully",
  "data": {
    "item_id": "1",
    "item_name": "500ml PET Bottle (Clear)",
    "opening_stock": 1000,
    "opening_stock_rate_per_unit": 3.50,
    "total_opening_stock_value": 3500.00,
    "unit": "pcs"
  }
}
```

---

#### Set Opening Stock - Caps

**Request:**
```json
PUT /items/2/opening-stock
Content-Type: application/json
Authorization: Bearer {jwt_token}

{
  "opening_stock": 1000,
  "opening_stock_rate_per_unit": 1.20
}
```

**Response (Success - 200 OK):**
```json
{
  "success": true,
  "message": "Opening stock updated successfully",
  "data": {
    "item_id": "2",
    "item_name": "Flip Cap with Seal (28mm)",
    "opening_stock": 1000,
    "opening_stock_rate_per_unit": 1.20,
    "total_opening_stock_value": 1200.00,
    "unit": "pcs"
  }
}
```

---

#### Set Opening Stock - Labels

**Request:**
```json
PUT /items/3/opening-stock
Content-Type: application/json
Authorization: Bearer {jwt_token}

{
  "opening_stock": 1500,
  "opening_stock_rate_per_unit": 0.25
}
```

**Response (Success - 200 OK):**
```json
{
  "success": true,
  "message": "Opening stock updated successfully",
  "data": {
    "item_id": "3",
    "item_name": "Brand Label - 500ml Water Bottle",
    "opening_stock": 1500,
    "opening_stock_rate_per_unit": 0.25,
    "total_opening_stock_value": 375.00,
    "unit": "pcs"
  }
}
```

---

#### 2.6 Set Opening Stock for Finished Product Variants

**Request:**
```json
PUT /items/4/variants/opening-stock
Content-Type: application/json
Authorization: Bearer {jwt_token}

{
  "variants": [
    {
      "variant_sku": "WATER-500ML-PLAIN",
      "opening_stock": 100,
      "opening_stock_rate_per_unit": 8.95
    },
    {
      "variant_sku": "WATER-500ML-LEMON",
      "opening_stock": 50,
      "opening_stock_rate_per_unit": 9.95
    },
    {
      "variant_sku": "WATER-500ML-MINT",
      "opening_stock": 50,
      "opening_stock_rate_per_unit": 9.95
    }
  ]
}
```

**Response (Success - 200 OK):**
```json
{
  "success": true,
  "message": "Variant opening stock updated successfully",
  "data": {
    "item_id": "4",
    "item_name": "500ml Purified Water Bottle",
    "variants_updated": 3,
    "variants": [
      {
        "variant_sku": "WATER-500ML-PLAIN",
        "opening_stock": 100,
        "opening_stock_rate_per_unit": 8.95,
        "total_value": 895.00
      },
      {
        "variant_sku": "WATER-500ML-LEMON",
        "opening_stock": 50,
        "opening_stock_rate_per_unit": 9.95,
        "total_value": 497.50
      },
      {
        "variant_sku": "WATER-500ML-MINT",
        "opening_stock": 50,
        "opening_stock_rate_per_unit": 9.95,
        "total_value": 497.50
      }
    ],
    "total_opening_stock_value": 1890.00
  }
}
```

---

## STEP 3: PURCHASING STOCK (INBOUND)

### 3.1 Additional: Create Item Group (Bundle)

**Request:**
```json
POST /item-groups
Content-Type: application/json
Authorization: Bearer {jwt_token}

{
  "name": "500ml Water Bottle Assembly Kit",
  "description": "Complete kit for assembling 500ml water bottles containing bottle, flip cap with seal, and branded label. Used in production to create finished water bottles.",
  "is_active": true,
  "components": [
    {
      "item_id": "1",
      "variant_sku": null,
      "quantity": 1,
      "variant_details": {}
    },
    {
      "item_id": "2",
      "variant_sku": null,
      "quantity": 1,
      "variant_details": {}
    },
    {
      "item_id": "3",
      "variant_sku": null,
      "quantity": 1,
      "variant_details": {}
    }
  ]
}
```

**Response (Success - 201 Created):**
```json
{
  "success": true,
  "message": "Item group created successfully",
  "data": {
    "item_group_id": "1",
    "name": "500ml Water Bottle Assembly Kit",
    "description": "Complete kit for assembling 500ml water bottles",
    "is_active": true,
    "components": [
      {
        "component_id": "1",
        "item_id": "1",
        "item_name": "500ml PET Bottle (Clear)",
        "variant_sku": null,
        "quantity": 1
      },
      {
        "component_id": "2",
        "item_id": "2",
        "item_name": "Flip Cap with Seal (28mm)",
        "variant_sku": null,
        "quantity": 1
      },
      {
        "component_id": "3",
        "item_id": "3",
        "item_name": "Brand Label - 500ml Water Bottle",
        "variant_sku": null,
        "quantity": 1
      }
    ],
    "total_components": 3,
    "created_at": "2024-02-22T10:35:00Z"
  }
}
```

---

### 3.2 Create Purchase Order 1 - Bottles from Vendor

**Request:**
```json
POST /purchase-orders
Content-Type: application/json
Authorization: Bearer {jwt_token}

{
  "vendor_id": 1,
  "delivery_address_type": "organization",
  "delivery_address_id": null,
  "organization_name": "Water Bottle Manufacturing Plant",
  "organization_address": "Industrial Area, Manufacturing Zone, Nashik, MH 422212",
  "reference_no": "PO-BOTTLES-2024-001",
  "date": "2024-02-22T10:00:00Z",
  "delivery_date": "2024-03-05T10:00:00Z",
  "payment_terms": "Net 30",
  "shipment_preference": "Road",
  "line_items": [
    {
      "item_id": "1",
      "variant_sku": "",
      "description": "Clear 500ml PET Bottles - Food grade",
      "account": "Raw Materials Purchase",
      "quantity": 5000,
      "rate": 3.50,
      "variant_details": {}
    }
  ],
  "discount": 500,
  "discount_type": "fixed",
  "tax_id": 1,
  "adjustment": 0,
  "notes": "Bulk order for water bottle manufacturing - High quality bottles needed. Payment via bank transfer. Contact: Rajesh Kumar +91-9876543210",
  "terms_and_conditions": "Standard payment terms apply. Quality inspection at delivery. 2 days for return of damaged goods.",
  "attachments": []
}
```

**Response (Success - 201 Created):**
```json
{
  "success": true,
  "message": "Purchase Order created successfully",
  "data": {
    "purchase_order_id": "1",
    "vendor_id": 1,
    "vendor_name": "Premier Plastic Industries Pvt Ltd",
    "reference_no": "PO-BOTTLES-2024-001",
    "po_date": "2024-02-22",
    "delivery_date": "2024-03-05",
    "status": "draft",
    "total_amount": 17000.00,
    "line_items_count": 1,
    "created_at": "2024-02-22T10:00:00Z"
  }
}
```

---

### 3.3 Create Purchase Order 2 - Caps from Vendor

**Request:**
```json
POST /purchase-orders
Content-Type: application/json
Authorization: Bearer {jwt_token}

{
  "vendor_id": 1,
  "delivery_address_type": "organization",
  "organization_name": "Water Bottle Manufacturing Plant",
  "organization_address": "Industrial Area, Manufacturing Zone, Nashik, MH 422212",
  "reference_no": "PO-CAPS-2024-001",
  "date": "2024-02-22T10:00:00Z",
  "delivery_date": "2024-03-07T10:00:00Z",
  "payment_terms": "Net 30",
  "shipment_preference": "Road",
  "line_items": [
    {
      "item_id": "2",
      "description": "Flip caps with tamper-proof seal - 28mm",
      "account": "Raw Materials Purchase",
      "quantity": 5000,
      "rate": 1.20,
      "variant_details": {}
    }
  ],
  "discount": 0,
  "discount_type": "fixed",
  "tax_id": 1,
  "adjustment": 0,
  "notes": "Caps for 500ml water bottle assembly - High quality seal required for food safety compliance",
  "terms_and_conditions": "Standard payment terms apply. Quality inspection required."
}
```

**Response (Success - 201 Created):**
```json
{
  "success": true,
  "message": "Purchase Order created successfully",
  "data": {
    "purchase_order_id": "2",
    "vendor_id": 1,
    "vendor_name": "Premier Plastic Industries Pvt Ltd",
    "reference_no": "PO-CAPS-2024-001",
    "po_date": "2024-02-22",
    "delivery_date": "2024-03-07",
    "status": "draft",
    "total_amount": 6000.00,
    "line_items_count": 1,
    "created_at": "2024-02-22T10:05:00Z"
  }
}
```

---

### 3.4 Create Purchase Order 3 - Labels from Vendor

**Request:**
```json
POST /purchase-orders
Content-Type: application/json
Authorization: Bearer {jwt_token}

{
  "vendor_id": 2,
  "delivery_address_type": "organization",
  "organization_name": "Water Bottle Manufacturing Plant",
  "organization_address": "Industrial Area, Manufacturing Zone, Nashik, MH 422212",
  "reference_no": "PO-LABELS-2024-001",
  "date": "2024-02-22T10:00:00Z",
  "delivery_date": "2024-03-10T10:00:00Z",
  "payment_terms": "Net 45",
  "shipment_preference": "Road",
  "line_items": [
    {
      "item_id": "3",
      "description": "Water bottle labels with brand logo and nutrition facts - 500ml size",
      "account": "Raw Materials Purchase",
      "quantity": 5000,
      "rate": 0.25,
      "variant_details": {}
    }
  ],
  "discount": 0,
  "discount_type": "fixed",
  "tax_id": 1,
  "adjustment": 0,
  "notes": "Print labels for 500ml water bottles - Full color print as per approved design. Include barcode and regulatory text.",
  "terms_and_conditions": "Extended payment terms - Net 45 days. Quality approval required before delivery."
}
```

**Response (Success - 201 Created):**
```json
{
  "success": true,
  "message": "Purchase Order created successfully",
  "data": {
    "purchase_order_id": "3",
    "vendor_id": 2,
    "vendor_name": "PrintPak Solutions Pvt Ltd",
    "reference_no": "PO-LABELS-2024-001",
    "po_date": "2024-02-22",
    "delivery_date": "2024-03-10",
    "status": "draft",
    "total_amount": 1250.00,
    "line_items_count": 1,
    "created_at": "2024-02-22T10:10:00Z"
  }
}
```

---

### 3.5 Record Purchase Receive (Stock Received)

**Request (Example: Bottles received):**
```json
POST /purchase-orders/1/receive
Content-Type: application/json
Authorization: Bearer {jwt_token}

{
  "receive_date": "2024-03-05T14:30:00Z",
  "received_items": [
    {
      "item_id": "1",
      "quantity_ordered": 5000,
      "quantity_received": 5000,
      "quantity_rejected": 0,
      "notes": "All 5000 bottles received in good condition. QC inspection completed."
    }
  ],
  "reference_number": "GRN-BOTTLES-2024-001",
  "notes": "Purchase order PO-BOTTLES-2024-001 received completely. All items inspected and approved."
}
```

**Response (Success - 200 OK):**
```json
{
  "success": true,
  "message": "Purchase receipt recorded successfully",
  "data": {
    "receipt_id": "GRN-1",
    "purchase_order_id": "1",
    "receipt_date": "2024-03-05",
    "items_received": 1,
    "total_quantity": 5000,
    "status": "received",
    "inventory_updated": true,
    "created_at": "2024-03-05T14:30:00Z"
  }
}
```

---

### 3.6 Record Bills (Accounts Payable)

#### Bill 1 - Invoice from Bottle Supplier

**Request:**
```json
POST /bills
Content-Type: application/json
Authorization: Bearer {jwt_token}

{
  "vendor_id": 1,
  "bill_number": "BILL-BOTTLES-2024-001",
  "purchase_order_id": "1",
  "billing_address": "Supplier Factory, Industrial Zone, Nashik",
  "order_number": "PO-BOTTLES-2024-001",
  "bill_date": "2024-03-05T10:00:00Z",
  "due_date": "2024-04-05T10:00:00Z",
  "payment_terms": "Net 30",
  "subject": "Invoice for 5000 Clear PET Bottles - 500ml Capacity",
  "line_items": [
    {
      "item_id": "1",
      "variant_sku": "",
      "description": "Clear 500ml PET Bottles - Food grade plastic",
      "account": "Raw Materials Purchase",
      "quantity": 5000,
      "rate": 3.50,
      "variant_details": {}
    }
  ],
  "discount": 500,
  "discount_type": "fixed",
  "tax_id": 1,
  "tax_rate": 5.0,
  "adjustment": 0,
  "notes": "Invoice from bottle manufacturer for PO-BOTTLES-2024-001. Received on 2024-03-05. Quality approved.",
  "attachments": []
}
```

**Response (Success - 201 Created):**
```json
{
  "success": true,
  "message": "Bill created successfully",
  "data": {
    "bill_id": "1",
    "vendor_id": 1,
    "vendor_name": "Premier Plastic Industries Pvt Ltd",
    "bill_number": "BILL-BOTTLES-2024-001",
    "bill_date": "2024-03-05",
    "due_date": "2024-04-05",
    "subtotal": 17500.00,
    "discount": 500.00,
    "tax": 850.00,
    "total_amount": 17850.00,
    "status": "draft",
    "created_at": "2024-03-05T10:00:00Z"
  }
}
```

---

#### Bill 2 - Invoice from Cap Supplier

**Request:**
```json
POST /bills
Content-Type: application/json
Authorization: Bearer {jwt_token}

{
  "vendor_id": 1,
  "bill_number": "BILL-CAPS-2024-001",
  "purchase_order_id": "2",
  "billing_address": "Supplier Factory, Industrial Zone, Nashik",
  "order_number": "PO-CAPS-2024-001",
  "bill_date": "2024-03-07T10:00:00Z",
  "due_date": "2024-04-07T10:00:00Z",
  "payment_terms": "Net 30",
  "subject": "Invoice for 5000 Flip Caps with Tamper-Proof Seal",
  "line_items": [
    {
      "item_id": "2",
      "description": "Flip caps with tamper-proof seal - 28mm neck fit",
      "account": "Raw Materials Purchase",
      "quantity": 5000,
      "rate": 1.20,
      "variant_details": {}
    }
  ],
  "discount": 0,
  "discount_type": "fixed",
  "tax_id": 1,
  "tax_rate": 5.0,
  "adjustment": 0,
  "notes": "Invoice from cap supplier. Received on 2024-03-07. All items approved."
}
```

**Response (Success - 201 Created):**
```json
{
  "success": true,
  "message": "Bill created successfully",
  "data": {
    "bill_id": "2",
    "vendor_id": 1,
    "vendor_name": "Premier Plastic Industries Pvt Ltd",
    "bill_number": "BILL-CAPS-2024-001",
    "bill_date": "2024-03-07",
    "due_date": "2024-04-07",
    "subtotal": 6000.00,
    "discount": 0.00,
    "tax": 300.00,
    "total_amount": 6300.00,
    "status": "draft",
    "created_at": "2024-03-07T10:00:00Z"
  }
}
```

---

#### Bill 3 - Invoice from Label Supplier

**Request:**
```json
POST /bills
Content-Type: application/json
Authorization: Bearer {jwt_token}

{
  "vendor_id": 2,
  "bill_number": "BILL-LABELS-2024-001",
  "purchase_order_id": "3",
  "billing_address": "PrintPak Factory, Business Park, Pune",
  "order_number": "PO-LABELS-2024-001",
  "bill_date": "2024-03-10T10:00:00Z",
  "due_date": "2024-04-25T10:00:00Z",
  "payment_terms": "Net 45",
  "subject": "Invoice for 5000 Water Bottle Labels with Brand Logo",
  "line_items": [
    {
      "item_id": "3",
      "description": "Water bottle labels with brand logo and nutrition facts - 500ml",
      "account": "Raw Materials Purchase",
      "quantity": 5000,
      "rate": 0.25,
      "variant_details": {}
    }
  ],
  "discount": 0,
  "discount_type": "fixed",
  "tax_id": 1,
  "tax_rate": 5.0,
  "adjustment": 0,
  "notes": "Invoice from label printer. Full color print as specified. Received and approved on 2024-03-10."
}
```

**Response (Success - 201 Created):**
```json
{
  "success": true,
  "message": "Bill created successfully",
  "data": {
    "bill_id": "3",
    "vendor_id": 2,
    "vendor_name": "PrintPak Solutions Pvt Ltd",
    "bill_number": "BILL-LABELS-2024-001",
    "bill_date": "2024-03-10",
    "due_date": "2024-04-25",
    "subtotal": 1250.00,
    "discount": 0.00,
    "tax": 62.50,
    "total_amount": 1312.50,
    "status": "draft",
    "created_at": "2024-03-10T10:00:00Z"
  }
}
```

---

## STEP 4: SELLING TO CUSTOMERS (OUTBOUND)

### 4.1 Create Production Order (Additional Step - Not in Standard Zoho)

**Request:**
```json
POST /production-orders
Content-Type: application/json
Authorization: Bearer {jwt_token}

{
  "item_group_id": "1",
  "quantity_to_manufacture": 1000,
  "planned_start_date": "2024-02-24",
  "planned_end_date": "2024-02-28",
  "notes": "Assembly batch for 1000 units of 500ml water bottles. Production workflow: Receive raw materials from warehouse → Fill empty bottles with purified water using multi-stage filtration system → Inspect water quality → Apply tamper-proof seal → Apply branded label → Final QC check → Pack in cartons (20 bottles per carton) → Store in finished goods warehouse"
}
```

**Response (Success - 201 Created):**
```json
{
  "success": true,
  "message": "Production Order created successfully",
  "data": {
    "production_order_id": "1",
    "item_group_id": "1",
    "quantity_to_manufacture": 1000,
    "planned_start_date": "2024-02-24",
    "planned_end_date": "2024-02-28",
    "status": "draft",
    "created_at": "2024-02-22T10:40:00Z"
  }
}
```

---

### 4.2 Update Production Order Status - In Progress

**Request:**
```json
PUT /production-orders/1
Content-Type: application/json
Authorization: Bearer {jwt_token}

{
  "status": "in_progress",
  "quantity_manufactured": 500,
  "actual_start_date": "2024-02-24",
  "notes": "Production started as scheduled - 500 units completed out of 1000. Manufactured items: 500 bottles filled with purified water, labeled, capped and QC checked. Remaining 500 units in progress. Expected completion: 2024-02-28"
}
```

**Response (Success - 200 OK):**
```json
{
  "success": true,
  "message": "Production Order updated successfully",
  "data": {
    "production_order_id": "1",
    "status": "in_progress",
    "quantity_manufactured": 500,
    "quantity_remaining": 500,
    "actual_start_date": "2024-02-24",
    "updated_at": "2024-02-24T08:00:00Z"
  }
}
```

---

### 4.3 Consume Raw Materials - Bottles

**Request:**
```json
POST /production-orders/1/consume-item
Content-Type: application/json
Authorization: Bearer {jwt_token}

{
  "production_order_item_id": 1,
  "quantity_consumed": 1000,
  "notes": "Consumed 1000 PET bottles from raw materials inventory for filling with purified water. Bottles issued from warehouse at 12:00 PM on 2024-02-24. All bottles inspected for defects before use."
}
```

**Response (Success - 200 OK):**
```json
{
  "success": true,
  "message": "Material consumption recorded successfully",
  "data": {
    "consumption_id": "MC-1",
    "production_order_id": "1",
    "item_id": "1",
    "item_name": "500ml PET Bottle (Clear)",
    "quantity_consumed": 1000,
    "unit_cost": 3.50,
    "total_cost": 3500.00,
    "inventory_updated": true,
    "updated_at": "2024-02-24T12:00:00Z"
  }
}
```

---

### 4.4 Consume Raw Materials - Caps

**Request:**
```json
POST /production-orders/1/consume-item
Content-Type: application/json
Authorization: Bearer {jwt_token}

{
  "production_order_item_id": 2,
  "quantity_consumed": 1000,
  "notes": "Consumed 1000 flip caps with tamper-proof seals for sealing water bottles. Caps issued from inventory at 14:00 PM on 2024-02-24. All caps checked for tamper-proof seal integrity."
}
```

**Response (Success - 200 OK):**
```json
{
  "success": true,
  "message": "Material consumption recorded successfully",
  "data": {
    "consumption_id": "MC-2",
    "production_order_id": "1",
    "item_id": "2",
    "item_name": "Flip Cap with Seal (28mm)",
    "quantity_consumed": 1000,
    "unit_cost": 1.20,
    "total_cost": 1200.00,
    "inventory_updated": true,
    "updated_at": "2024-02-24T14:00:00Z"
  }
}
```

---

### 4.5 Consume Raw Materials - Labels

**Request:**
```json
POST /production-orders/1/consume-item
Content-Type: application/json
Authorization: Bearer {jwt_token}

{
  "production_order_item_id": 3,
  "quantity_consumed": 1000,
  "notes": "Applied 1000 branded labels to water bottles. Labels applied at 15:30 PM on 2024-02-24. All labels properly aligned, no defects observed."
}
```

**Response (Success - 200 OK):**
```json
{
  "success": true,
  "message": "Material consumption recorded successfully",
  "data": {
    "consumption_id": "MC-3",
    "production_order_id": "1",
    "item_id": "3",
    "item_name": "Brand Label - 500ml Water Bottle",
    "quantity_consumed": 1000,
    "unit_cost": 0.25,
    "total_cost": 250.00,
    "inventory_updated": true,
    "updated_at": "2024-02-24T15:30:00Z"
  }
}
```

---

### 4.6 Create Sales Order (With Variants)

**Request:**
```json
POST /sales-orders
Content-Type: application/json
Authorization: Bearer {jwt_token}

{
  "customer_id": 1,
  "reference_no": "SO-WATER-2024-001",
  "sales_order_date": "2024-02-25T10:00:00Z",
  "expected_shipment_date": "2024-03-01T10:00:00Z",
  "payment_terms": "Net 15",
  "delivery_method": "Courier",
  "salesperson_id": 1,
  "line_items": [
    {
      "item_id": "4",
      "variant_sku": "WATER-500ML-PLAIN",
      "description": "500ml Purified Water Bottle - Plain Flavor",
      "quantity": 200,
      "rate": 20.00,
      "variant_details": {
        "flavor": "Plain"
      }
    },
    {
      "item_id": "4",
      "variant_sku": "WATER-500ML-LEMON",
      "description": "500ml Purified Water Bottle - Lemon Flavor",
      "quantity": 150,
      "rate": 22.00,
      "variant_details": {
        "flavor": "Lemon"
      }
    },
    {
      "item_id": "4",
      "variant_sku": "WATER-500ML-MINT",
      "description": "500ml Purified Water Bottle - Mint Flavor",
      "quantity": 150,
      "rate": 22.00,
      "variant_details": {
        "flavor": "Mint"
      }
    }
  ],
  "shipping_charges": 1000,
  "tax_id": 1,
  "tax_rate": 5.0,
  "adjustment": 0,
  "customer_notes": "Bulk order for retail distribution chain with mixed flavors. Deliver to warehouse address in Thane. Installation of delivery not required. Please arrange pallet delivery.",
  "terms_and_conditions": "Payment due within 15 days. Free delivery for orders above 500 units. Goods are sold as per invoice attached. Acceptance deemed complete upon delivery unless noted otherwise within 48 hours."
}
```

**Response (Success - 201 Created):**
```json
{
  "success": true,
  "message": "Sales Order created successfully",
  "data": {
    "sales_order_id": "1",
    "customer_id": 1,
    "customer_name": "Fresh Waters Distribution Ltd",
    "reference_no": "SO-WATER-2024-001",
    "so_date": "2024-02-25",
    "expected_shipment_date": "2024-03-01",
    "status": "draft",
    "line_items": [
      {
        "line_item_id": "1",
        "item_id": "4",
        "variant_sku": "WATER-500ML-PLAIN",
        "description": "500ml Purified Water Bottle - Plain Flavor",
        "quantity": 200,
        "rate": 20.00,
        "total": 4000.00
      },
      {
        "line_item_id": "2",
        "item_id": "4",
        "variant_sku": "WATER-500ML-LEMON",
        "description": "500ml Purified Water Bottle - Lemon Flavor",
        "quantity": 150,
        "rate": 22.00,
        "total": 3300.00
      },
      {
        "line_item_id": "3",
        "item_id": "4",
        "variant_sku": "WATER-500ML-MINT",
        "description": "500ml Purified Water Bottle - Mint Flavor",
        "quantity": 150,
        "rate": 22.00,
        "total": 3300.00
      }
    ],
    "subtotal": 10600.00,
    "tax": 530.00,
    "shipping": 1000.00,
    "total_amount": 12130.00,
    "line_items_count": 3,
    "created_at": "2024-02-25T10:00:00Z"
  }
}
```

---

### 4.7 Update Sales Order Status - Confirmed

**Request:**
```json
PATCH /sales-orders/1/status
Content-Type: application/json
Authorization: Bearer {jwt_token}

{
  "status": "confirmed"
}
```

**Response (Success - 200 OK):**
```json
{
  "success": true,
  "message": "Sales Order status updated successfully",
  "data": {
    "sales_order_id": "1",
    "reference_no": "SO-WATER-2024-001",
    "old_status": "draft",
    "new_status": "confirmed",
    "updated_at": "2024-02-25T14:00:00Z"
  }
}
```

---

### 4.8 Create Package (Packing Slip - With Variants)

**Request:**
```json
POST /packages
Content-Type: application/json
Authorization: Bearer {jwt_token}

{
  "sales_order_id": "1",
  "customer_id": 1,
  "package_date": "2024-02-29T10:00:00Z",
  "items": [
    {
      "sales_order_item_id": 1,
      "item_id": "4",
      "variant_sku": "WATER-500ML-PLAIN",
      "ordered_qty": 200,
      "packed_qty": 200,
      "variant_details": {
        "flavor": "Plain"
      }
    },
    {
      "sales_order_item_id": 2,
      "item_id": "4",
      "variant_sku": "WATER-500ML-LEMON",
      "ordered_qty": 150,
      "packed_qty": 150,
      "variant_details": {
        "flavor": "Lemon"
      }
    },
    {
      "sales_order_item_id": 3,
      "item_id": "4",
      "variant_sku": "WATER-500ML-MINT",
      "ordered_qty": 150,
      "packed_qty": 150,
      "variant_details": {
        "flavor": "Mint"
      }
    }
  ],
  "internal_notes": "Packaged 500 units (200 Plain + 150 Lemon + 150 Mint) of 500ml water bottles in 25 cartons (20 bottles per carton). Carton breakdown: 10 cartons Plain, 8 cartons Lemon, 8 cartons Mint. Each carton dimensions: 30cm x 25cm x 20cm. Total carton weight: 12kg each. QC inspection completed - all items approved. Ready for shipment. Packing date: 2024-02-29 at 10:00 AM. Packed by: Warehouse Team. Checked by: Quality Control Manager."
}
```

**Response (Success - 201 Created):**
```json
{
  "success": true,
  "message": "Package created successfully",
  "data": {
    "package_id": "1",
    "sales_order_id": "1",
    "customer_id": 1,
    "package_date": "2024-02-29",
    "status": "created",
    "items": [
      {
        "item_id": "4",
        "variant_sku": "WATER-500ML-PLAIN",
        "variant_details": {
          "flavor": "Plain"
        },
        "packed_qty": 200
      },
      {
        "item_id": "4",
        "variant_sku": "WATER-500ML-LEMON",
        "variant_details": {
          "flavor": "Lemon"
        },
        "packed_qty": 150
      },
      {
        "item_id": "4",
        "variant_sku": "WATER-500ML-MINT",
        "variant_details": {
          "flavor": "Mint"
        },
        "packed_qty": 150
      }
    ],
    "total_items": 3,
    "total_quantity": 500,
    "cartons_count": 25,
    "created_at": "2024-02-29T10:00:00Z"
  }
}
```

---

### 4.9 Update Package Status - Packed

**Request:**
```json
PATCH /packages/1/status
Content-Type: application/json
Authorization: Bearer {jwt_token}

{
  "status": "packed"
}
```

**Response (Success - 200 OK):**
```json
{
  "success": true,
  "message": "Package status updated successfully",
  "data": {
    "package_id": "1",
    "old_status": "created",
    "new_status": "packed",
    "updated_at": "2024-02-29T14:30:00Z"
  }
}
```

---

### 4.10 Create Shipment

**Request:**
```json
POST /shipments
Content-Type: application/json
Authorization: Bearer {jwt_token}

{
  "package_id": "1",
  "sales_order_id": "1",
  "customer_id": 1,
  "ship_date": "2024-03-01T08:00:00Z",
  "carrier": "Premium Logistics Limited",
  "tracking_no": "TRACK-WATER-500ML-001",
  "tracking_url": "https://tracking.premiumlogistics.com/TRACK-WATER-500ML-001",
  "shipping_charges": 1000,
  "notes": "500ml water bottles shipped in 25 cartons (500 units total). Pallet weight: 300kg. Estimated delivery date: 2024-03-05. Shipped via road transport. Delivery address: Fresh Waters Warehouse, Thane, MH. Contact: +91-9988776655. Insurance included."
}
```

**Response (Success - 201 Created):**
```json
{
  "success": true,
  "message": "Shipment created successfully",
  "data": {
    "shipment_id": "1",
    "package_id": "1",
    "sales_order_id": "1",
    "carrier": "Premium Logistics Limited",
    "tracking_no": "TRACK-WATER-500ML-001",
    "ship_date": "2024-03-01",
    "status": "created",
    "estimated_delivery": "2024-03-05",
    "created_at": "2024-03-01T08:00:00Z"
  }
}
```

---

### 4.11 Update Shipment Status - Shipped

**Request:**
```json
PATCH /shipments/1/status
Content-Type: application/json
Authorization: Bearer {jwt_token}

{
  "status": "shipped"
}
```

**Response (Success - 200 OK):**
```json
{
  "success": true,
  "message": "Shipment status updated successfully",
  "data": {
    "shipment_id": "1",
    "tracking_no": "TRACK-WATER-500ML-001",
    "old_status": "created",
    "new_status": "shipped",
    "updated_at": "2024-03-01T08:30:00Z"
  }
}
```

---

### 4.12 Create Invoice (Bill to Customer - With Variants)

**Request:**
```json
POST /invoices
Content-Type: application/json
Authorization: Bearer {jwt_token}

{
  "customer_id": 1,
  "order_number": "SO-WATER-2024-001",
  "invoice_date": "2024-03-01T09:00:00Z",
  "due_date": "2024-03-15T09:00:00Z",
  "invoice_number": "INV-WATER-2024-001",
  "terms": "Net 15",
  "salesperson_id": 1,
  "subject": "Invoice for 500ml Water Bottle Bulk Order with Multiple Flavors - Fresh Waters Distribution",
  "line_items": [
    {
      "item_id": "4",
      "variant_sku": "WATER-500ML-PLAIN",
      "description": "500ml Purified Water Bottle - Plain Flavor (200 units @ ₹20.00 each)",
      "quantity": 200,
      "rate": 20.00,
      "variant_details": {
        "flavor": "Plain"
      }
    },
    {
      "item_id": "4",
      "variant_sku": "WATER-500ML-LEMON",
      "description": "500ml Purified Water Bottle - Lemon Flavor (150 units @ ₹22.00 each)",
      "quantity": 150,
      "rate": 22.00,
      "variant_details": {
        "flavor": "Lemon"
      }
    },
    {
      "item_id": "4",
      "variant_sku": "WATER-500ML-MINT",
      "description": "500ml Purified Water Bottle - Mint Flavor (150 units @ ₹22.00 each)",
      "quantity": 150,
      "rate": 22.00,
      "variant_details": {
        "flavor": "Mint"
      }
    }
  ],
  "tax_type": "exclusive",
  "tax_id": 1,
  "tax_rate": 5.0,
  "shipping_charges": 1000,
  "adjustment": 0,
  "customer_notes": "Invoice for 500ml water bottles with mixed flavors delivered. Please process payment within 15 days. Goods received at Fresh Waters Warehouse, Thane on 2024-03-05. Contact: Arun Verma +91-9988776655. Flavor breakdown: 200 Plain, 150 Lemon, 150 Mint.",
  "terms_and_conditions": "Payment terms: Net 15 days from invoice date. 2% discount if paid within 7 days. Goods sold as per invoice attached. All claims must be made within 30 days of delivery. Return policy: Goods can be returned if damaged within 48 hours of delivery with original packaging.",
  "payment_received": false,
  "payment_splits": [
    {
      "payment_mode": "bank_transfer",
      "deposit_to": "Company Bank Account - SBIN0001234",
      "amount_received": 0
    }
  ],
  "email_recipients": [
    "procurement@freshwaters.com",
    "billing@freshwaters.com",
    "arun@freshwaters.com"
  ]
}
```

**Response (Success - 201 Created):**
```json
{
  "success": true,
  "message": "Invoice created successfully",
  "data": {
    "invoice_id": "1",
    "customer_id": 1,
    "customer_name": "Fresh Waters Distribution Ltd",
    "invoice_number": "INV-WATER-2024-001",
    "invoice_date": "2024-03-01",
    "due_date": "2024-03-15",
    "status": "draft",
    "line_items": [
      {
        "line_item_id": "1",
        "item_id": "4",
        "variant_sku": "WATER-500ML-PLAIN",
        "variant_details": {
          "flavor": "Plain"
        },
        "quantity": 200,
        "rate": 20.00,
        "total": 4000.00
      },
      {
        "line_item_id": "2",
        "item_id": "4",
        "variant_sku": "WATER-500ML-LEMON",
        "variant_details": {
          "flavor": "Lemon"
        },
        "quantity": 150,
        "rate": 22.00,
        "total": 3300.00
      },
      {
        "line_item_id": "3",
        "item_id": "4",
        "variant_sku": "WATER-500ML-MINT",
        "variant_details": {
          "flavor": "Mint"
        },
        "quantity": 150,
        "rate": 22.00,
        "total": 3300.00
      }
    ],
    "subtotal": 10600.00,
    "tax": 530.00,
    "shipping": 1000.00,
    "total_amount": 12130.00,
    "line_items_count": 3,
    "payment_status": "unpaid",
    "created_at": "2024-03-01T09:00:00Z"
  }
}
```

---

## STEP 5: FINALIZING THE LOOP

### 5.1 Update Invoice Status - Sent

**Request:**
```json
PATCH /invoices/1/status
Content-Type: application/json
Authorization: Bearer {jwt_token}

{
  "status": "sent"
}
```

**Response (Success - 200 OK):**
```json
{
  "success": true,
  "message": "Invoice status updated successfully",
  "data": {
    "invoice_id": "1",
    "invoice_number": "INV-WATER-2024-001",
    "old_status": "draft",
    "new_status": "sent",
    "sent_to": [
      "procurement@freshwaters.com",
      "billing@freshwaters.com"
    ],
    "updated_at": "2024-03-01T10:00:00Z"
  }
}
```

---

### 5.2 Record Payment Received from Customer

**Request:**
```json
POST /payments
Content-Type: application/json
Authorization: Bearer {jwt_token}

{
  "invoice_id": "1",
  "payment_date": "2024-03-08T15:00:00Z",
  "amount": 12130,
  "payment_mode": "bank_transfer",
  "reference": "BANK-TXN-2024-0308-FW001",
  "notes": "Payment received from Fresh Waters Distribution Ltd via bank transfer. Invoice INV-WATER-2024-001 for 500 units of 500ml water bottles with mixed flavors (200 Plain + 150 Lemon + 150 Mint). Amount: ₹12,130 (including tax and shipping). Bank Reference: BANK-TXN-2024-0308-FW001. Deposited to SBIN0001234."
}
```

**Response (Success - 201 Created):**
```json
{
  "success": true,
  "message": "Payment recorded successfully",
  "data": {
    "payment_id": "1",
    "invoice_id": "1",
    "customer_id": 1,
    "invoice_number": "INV-WATER-2024-001",
    "payment_date": "2024-03-08",
    "amount": 12130.00,
    "payment_mode": "bank_transfer",
    "reference": "BANK-TXN-2024-0308-FW001",
    "status": "completed",
    "created_at": "2024-03-08T15:00:00Z"
  }
}
```

---

### 5.3 Update Invoice Status - Paid

**Request:**
```json
PATCH /invoices/1/status
Content-Type: application/json
Authorization: Bearer {jwt_token}

{
  "status": "paid"
}
```

**Response (Success - 200 OK):**
```json
{
  "success": true,
  "message": "Invoice payment status updated successfully",
  "data": {
    "invoice_id": "1",
    "invoice_number": "INV-WATER-2024-001",
    "old_status": "sent",
    "new_status": "paid",
    "payment_id": "1",
    "paid_amount": 12130.00,
    "updated_at": "2024-03-08T15:30:00Z"
  }
}
```

---

### 5.4 Get Variants Opening Stock

**Request:**
```json
GET /items/4/variants/opening-stock
Authorization: Bearer {jwt_token}
```

**Response (Success - 200 OK):**
```json
{
  "success": true,
  "message": "Variant opening stock retrieved successfully",
  "data": {
    "item_id": "4",
    "item_name": "500ml Purified Water Bottle",
    "variants": [
      {
        "variant_sku": "WATER-500ML-PLAIN",
        "variant_details": {
          "flavor": "Plain"
        },
        "opening_stock": 100,
        "opening_stock_rate_per_unit": 8.95,
        "total_value": 895.00,
        "current_stock": 100
      },
      {
        "variant_sku": "WATER-500ML-LEMON",
        "variant_details": {
          "flavor": "Lemon"
        },
        "opening_stock": 50,
        "opening_stock_rate_per_unit": 9.95,
        "total_value": 497.50,
        "current_stock": 50
      },
      {
        "variant_sku": "WATER-500ML-MINT",
        "variant_details": {
          "flavor": "Mint"
        },
        "opening_stock": 50,
        "opening_stock_rate_per_unit": 9.95,
        "total_value": 497.50,
        "current_stock": 50
      }
    ],
    "total_opening_stock": 200,
    "total_opening_stock_value": 1890.00
  }
}
```

---

### 5.5 Record Payment to Vendor 1 - Bottles

**Request:**
```json
POST /payments
Content-Type: application/json
Authorization: Bearer {jwt_token}

{
  "bill_id": "1",
  "payment_date": "2024-03-25T11:00:00Z",
  "amount": 17850,
  "payment_mode": "bank_transfer",
  "reference": "VENDOR-PAY-2024-0325-PPI001",
  "notes": "Payment made to Premier Plastic Industries for bill BILL-BOTTLES-2024-001. Amount: ₹17,850 (Invoice subtotal ₹17,500 + Tax ₹350 - Discount ₹0). Bank reference: VENDOR-PAY-2024-0325-PPI001. Payment due date: 2024-04-05. Paid early."
}
```

**Response (Success - 201 Created):**
```json
{
  "success": true,
  "message": "Vendor payment recorded successfully",
  "data": {
    "payment_id": "2",
    "bill_id": "1",
    "vendor_id": 1,
    "vendor_name": "Premier Plastic Industries Pvt Ltd",
    "bill_number": "BILL-BOTTLES-2024-001",
    "payment_date": "2024-03-25",
    "amount": 17850.00,
    "payment_mode": "bank_transfer",
    "reference": "VENDOR-PAY-2024-0325-PPI001",
    "status": "completed",
    "created_at": "2024-03-25T11:00:00Z"
  }
}
```

---

### 5.5 Record Payment to Vendor 2 - Labels

**Request:**
```json
POST /payments
Content-Type: application/json
Authorization: Bearer {jwt_token}

{
  "bill_id": "3",
  "payment_date": "2024-04-15T14:00:00Z",
  "amount": 1312.50,
  "payment_mode": "bank_transfer",
  "reference": "VENDOR-PAY-2024-0415-PRINTPAK001",
  "notes": "Payment made to PrintPak Solutions for bill BILL-LABELS-2024-001. Amount: ₹1,312.50 (Invoice subtotal ₹1,250 + Tax ₹62.50). Bank reference: VENDOR-PAY-2024-0415-PRINTPAK001. Payment terms: Net 45 (Due: 2024-04-25). Early payment discount applied."
}
```

**Response (Success - 201 Created):**
```json
{
  "success": true,
  "message": "Vendor payment recorded successfully",
  "data": {
    "payment_id": "3",
    "bill_id": "3",
    "vendor_id": 2,
    "vendor_name": "PrintPak Solutions Pvt Ltd",
    "bill_number": "BILL-LABELS-2024-001",
    "payment_date": "2024-04-15",
    "amount": 1312.50,
    "payment_mode": "bank_transfer",
    "reference": "VENDOR-PAY-2024-0415-PRINTPAK001",
    "status": "completed",
    "created_at": "2024-04-15T14:00:00Z"
  }
}
```

---

### 5.6 Update Shipment Status - Delivered

**Request:**
```json
PATCH /shipments/1/status
Content-Type: application/json
Authorization: Bearer {jwt_token}

{
  "status": "delivered"
}
```

**Response (Success - 200 OK):**
```json
{
  "success": true,
  "message": "Shipment delivery status updated successfully",
  "data": {
    "shipment_id": "1",
    "tracking_no": "TRACK-WATER-500ML-001",
    "customer_name": "Fresh Waters Distribution Ltd",
    "old_status": "shipped",
    "new_status": "delivered",
    "delivered_date": "2024-03-05",
    "updated_at": "2024-03-05T17:00:00Z"
  }
}
```

---

## BUSINESS ANALYTICS & DASHBOARD

### 6.1 Get Inventory Summary

**Request:**
```json
GET /items/1/stock-summary
Authorization: Bearer {jwt_token}
```

**Response (Success - 200 OK):**
```json
{
  "success": true,
  "message": "Stock summary retrieved successfully",
  "data": {
    "item_id": "1",
    "item_name": "500ml PET Bottle (Clear)",
    "sku": "BOTTLE-500ML-CLR",
    "opening_stock": 1000,
    "opening_stock_value": 3500.00,
    "purchases": 5000,
    "purchases_value": 17500.00,
    "consumed_in_production": 1000,
    "current_stock": 5000,
    "reorder_point": 500,
    "status": "good",
    "valuation_method": "FIFO",
    "last_purchase_rate": 3.50,
    "current_value": 17500.00
  }
}
```

---

### 6.2 Dashboard - Business Summary

**Request:**
```json
GET /auth/admin/dashboard/stats?start_date=2024-02-01&end_date=2024-03-31
Authorization: Bearer {jwt_token}
```

**Response (Success - 200 OK):**
```json
{
  "success": true,
  "message": "Dashboard stats retrieved successfully",
  "data": {
    "period": {
      "start_date": "2024-02-01",
      "end_date": "2024-03-31"
    },
    "sales": {
      "total_invoices": 1,
      "total_sales_amount": 12130.00,
      "total_sales_orders": 1,
      "average_order_value": 12130.00,
      "payment_received": 12130.00,
      "payment_pending": 0.00,
      "items_breakdown": {
        "plain_variant": {
          "quantity": 200,
          "rate": 20.00,
          "total": 4000.00
        },
        "lemon_variant": {
          "quantity": 150,
          "rate": 22.00,
          "total": 3300.00
        },
        "mint_variant": {
          "quantity": 150,
          "rate": 22.00,
          "total": 3300.00
        }
      }
    },
    "purchases": {
      "total_bills": 3,
      "total_purchase_amount": 25462.50,
      "total_purchase_orders": 3,
      "average_purchase_value": 8487.50,
      "payment_made": 19162.50,
      "payment_pending": 6300.00
    },
    "inventory": {
      "raw_materials_value": 21075.00,
      "finished_goods_value": 0.00,
      "work_in_progress_value": 0.00,
      "total_inventory_value": 21075.00,
      "items_count": 4,
      "variants_count": 3
    },
    "production": {
      "production_orders": 1,
      "units_manufactured": 500,
      "units_in_progress": 500,
      "manufacturing_cost": 4475.00,
      "variants_produced": {
        "plain": 200,
        "lemon": 150,
        "mint": 150
      }
    },
    "profitability": {
      "total_revenue": 12130.00,
      "total_cost_of_goods": 4850.00,
      "total_operating_costs": 1000.00,
      "gross_profit": 7280.00,
      "gross_profit_margin": 60.00,
      "net_profit": 6280.00,
      "net_profit_margin": 51.77,
      "roi_percentage": 29.10
    }
  }
}
```

---

## COMPLETE WORKFLOW SUMMARY

### Business Flow Overview

```
┌─────────────────────────────────────────────────────────────────┐
│         500ML WATER BOTTLE MANUFACTURING & SALES CYCLE          │
└─────────────────────────────────────────────────────────────────┘

PHASE 1: CONTACTS SETUP
├─ Create Vendor 1: Premier Plastic Industries (Bottles & Caps)
├─ Create Vendor 2: PrintPak Solutions (Labels)
└─ Create Customer: Fresh Waters Distribution Ltd

PHASE 2: INVENTORY SETUP
├─ Create Raw Material Items:
│  ├─ 500ml PET Bottle (Qty: 1000, Cost: ₹3.50)
│  ├─ Flip Cap with Seal (Qty: 1000, Cost: ₹1.20)
│  └─ Water Bottle Label (Qty: 1500, Cost: ₹0.25)
├─ Create Finished Product Item (With Variants):
│  ├─ 500ml Water Bottle - Plain (Cost: ₹8.95, Price: ₹20.00)
│  ├─ 500ml Water Bottle - Lemon (Cost: ₹9.95, Price: ₹22.00)
│  └─ 500ml Water Bottle - Mint (Cost: ₹9.95, Price: ₹22.00)
├─ Set Opening Inventory Stock
└─ Set Variant Opening Stock (200 total: 100 Plain, 50 Lemon, 50 Mint)

PHASE 3: PURCHASING (INBOUND)
├─ Create Item Group: 500ml Water Bottle Assembly Kit
├─ Purchase Order 1: 5000 Bottles @ ₹3.50 = ₹17,500
├─ Purchase Order 2: 5000 Caps @ ₹1.20 = ₹6,000
├─ Purchase Order 3: 5000 Labels @ ₹0.25 = ₹1,250
├─ Record Goods Received Notes (GRN)
├─ Receive all items to inventory
└─ Record Bills from Vendors (Accounts Payable)

PHASE 4: PRODUCTION (VALUE ADDITION)
├─ Create Production Order: 1000 units
├─ Mark as In Progress
├─ Consume Raw Materials:
│  ├─ Issue 1000 Bottles from inventory
│  ├─ Issue 1000 Caps from inventory
│  └─ Issue 1000 Labels from inventory
└─ Update Inventory: -1000 Bottles, -1000 Caps, -1000 Labels

PHASE 5: SALES (OUTBOUND - WITH VARIANTS)
├─ Create Sales Order: 500 units mixed flavors
│  ├─ Plain variant: 200 units @ ₹20.00 = ₹4,000
│  ├─ Lemon variant: 150 units @ ₹22.00 = ₹3,300
│  └─ Mint variant: 150 units @ ₹22.00 = ₹3,300
├─ Confirm Sales Order
├─ Create Package: 25 cartons (20 bottles each) with variant tracking
├─ Mark Package as Packed
├─ Create Shipment with Tracking Number
├─ Update Shipment Status to Shipped
└─ Update Shipment Status to Delivered

PHASE 6: INVOICING & PAYMENT
├─ Create Invoice to Customer: ₹12,130 (with tax & shipping - all variants)
├─ Send Invoice via Email
├─ Record Payment Received: ₹12,130 (Bank Transfer)
├─ Get Variant Opening Stock Details
├─ Update Invoice Status to Paid
├─ Pay Vendor Bills:
│  ├─ Pay Bottle Supplier: ₹17,850
│  └─ Pay Label Supplier: ₹1,312.50
└─ Finalize Accounts

PHASE 7: REPORTING & ANALYSIS
├─ View Inventory Summary
├─ Check Stock Levels
├─ Generate Profit & Loss Report
├─ View Dashboard Stats
└─ Monitor Margins and Efficiency
```

---

### Financial Summary (With Variants)

| Metric | Amount |
|--------|--------|
| **Revenue (500 units mixed)** | ₹10,600 |
| **Tax (5%)** | ₹530 |
| **Shipping Cost** | ₹1,000 |
| **Total Invoice Value** | ₹12,130 |
| **Cost of Goods Sold (COGS)** | ₹4,850 |
| **Gross Profit** | ₹7,280 |
| **Gross Margin** | 60.00% |
| **Operating Costs** | ₹1,000 |
| **Net Profit** | ₹6,280 |
| **Net Margin** | 51.77% |
| **ROI** | 29.10% |

**Variant Wise Revenue Breakdown:**
- Plain (200 units @ ₹20.00): ₹4,000
- Lemon (150 units @ ₹22.00): ₹3,300
- Mint (150 units @ ₹22.00): ₹3,300
- **Total Revenue**: ₹10,600

**Variant Wise Cost Breakdown:**
- Plain (200 units @ ₹8.95): ₹1,790
- Lemon (150 units @ ₹9.95): ₹1,492.50
- Mint (150 units @ ₹9.95): ₹1,492.50
- **Total COGS**: ₹4,850

---

### Inventory Tracking

**Raw Materials:**
- 500ml Bottles: 1000 → 5000 (purchased) → 4000 (5000 - 1000 used)
- Flip Caps: 1000 → 5000 (purchased) → 4000 (5000 - 1000 used)
- Labels: 1500 → 5000 (purchased) → 4000 (5000 - 1000 used)

**Work in Progress:**
- 500 units being manufactured (Mixed flavors: 200 Plain, 150 Lemon, 150 Mint)

**Finished Goods (By Variant):**
- Plain: 100 (opening) → 200 (produced) → 200 (sold) = 0 remaining
- Lemon: 50 (opening) → 150 (produced) → 150 (sold) = 0 remaining
- Mint: 50 (opening) → 150 (produced) → 150 (sold) = 0 remaining
- **Total**: 200 → 500 → 500 = 0 remaining

---

### Key Performance Indicators (KPIs)

1. **Order Fulfillment Cycle**: 7 days (SO Date → Delivered)
2. **Payment Collection**: 8 days (Invoice Date → Payment Received)
3. **Payment Disbursement**: 23 days to 44 days (Bill Date → Payment Made)
4. **Cash Conversion Cycle**: Positive (Cash collected before full payment to suppliers)
5. **Inventory Turnover**: High (Raw materials consumed in 5 days)
6. **Gross Margin**: 60.00% (Excellent, with variant pricing)
7. **ROI**: ₹6,280 profit on ₹21,575 total investment = 29.10%
8. **Variant Performance**:
   - Plain: ₹11.05 profit per unit (63.25% margin)
   - Lemon: ₹12.05 profit per unit (54.77% margin)
   - Mint: ₹12.05 profit per unit (54.77% margin)
9. **Variant Mix Optimization**: Plain (40%), Lemon (30%), Mint (30%) = Balanced portfolio

---

### Notes & Observations

1. **Variant Strategy**: Multiple flavors increase revenue per unit (Lemon/Mint @ ₹22 vs Plain @ ₹20)
2. **Scalability**: This workflow can easily be repeated with any combination of variants
3. **Profitability**: 
   - Plain generates ₹11.05 profit per unit (₹20 sale - ₹8.95 cost)
   - Lemon/Mint generate ₹12.05 profit per unit (₹22 sale - ₹9.95 cost)
   - Mixed portfolio optimizes overall margin at 60%
4. **Working Capital**: Positive flow - cash collected before full payment to suppliers
5. **Quality Control**: Included at multiple stages with variant-specific tracking
6. **Traceability**: Full tracking from raw material to final delivery, by variant
7. **Inventory Management**: Variant-specific opening stock, consumption, and sales tracking
8. **Compliance**: All documentation maintained for audit and regulatory requirements
9. **Customer Preference**: Tracking variant-wise sales helps identify customer preferences
10. **Production Flexibility**: Item groups enable easy scaling with different variant combinations

---

This complete workflow demonstrates a real-world manufacturing and distribution scenario using all the core functions of an inventory management system like Zoho Inventory, enhanced with variant management for multi-flavored products.
