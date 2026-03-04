# Company Complete Workflow - Request & Response Guide

This guide provides complete request and response examples for creating a company with all its dependencies. Follow the sequence below to properly set up a company.

---

## Prerequisites Flow

Before creating a company, you need to set up the following dependencies:

```
Country → State → Business Type
        ↓
      Bank
        ↓
    Tax Type
        ↓
    Company
```

---

## Step 1: Create Country

**Endpoint**: `POST /api/countries`

### Request
```json
{
  "country_name": "India",
  "country_code": "IN",
  "currency_code": "INR",
  "is_active": true
}
```

### Response (201 Created)
```json
{
  "id": 1,
  "country_name": "India",
  "country_code": "IN",
  "currency_code": "INR",
  "is_active": true,
  "created_at": "2026-03-03T10:00:00Z",
  "updated_at": "2026-03-03T10:00:00Z"
}
```

**Save**: `country_id = 1`

---

## Step 2: Create State

**Endpoint**: `POST /api/states`

### Request
```json
{
  "state_name": "Maharashtra",
  "state_code": "MH",
  "country_id": 1,
  "is_active": true
}
```

### Response (201 Created)
```json
{
  "id": 1,
  "state_name": "Maharashtra",
  "state_code": "MH",
  "country_id": 1,
  "is_active": true,
  "created_at": "2026-03-03T10:05:00Z",
  "updated_at": "2026-03-03T10:05:00Z"
}
```

**Save**: `state_id = 1`

---

## Step 3: Create Business Type

**Endpoint**: `POST /api/business-types`

### Request
```json
{
  "business_type_name": "Private Limited",
  "description": "Private limited company",
  "is_active": true
}
```

### Response (201 Created)
```json
{
  "id": 1,
  "business_type_name": "Private Limited",
  "description": "Private limited company",
  "is_active": true,
  "created_at": "2026-03-03T10:10:00Z",
  "updated_at": "2026-03-03T10:10:00Z"
}
```

**Save**: `business_type_id = 1`

---

## Step 4: Create Tax Type

**Endpoint**: `POST /api/tax-types`

### Request
```json
{
  "tax_type_name": "GST",
  "tax_description": "Goods and Services Tax",
  "tax_rate": 18.0,
  "is_active": true
}
```

### Response (201 Created)
```json
{
  "id": 1,
  "tax_type_name": "GST",
  "tax_description": "Goods and Services Tax",
  "tax_rate": 18.0,
  "is_active": true,
  "created_at": "2026-03-03T10:15:00Z",
  "updated_at": "2026-03-03T10:15:00Z"
}
```

**Save**: `tax_type_id = 1`

---

## Step 5: Create Bank

**Endpoint**: `POST /api/banks`

### Request
```json
{
  "bank_name": "State Bank of India",
  "bank_code": "SBIN0000001",
  "ifsc_code": "SBIN0001234",
  "swift_code": "SBININBB123",
  "country_id": 1,
  "is_active": true
}
```

### Response (201 Created)
```json
{
  "id": 1,
  "bank_name": "State Bank of India",
  "bank_code": "SBIN0000001",
  "ifsc_code": "SBIN0001234",
  "swift_code": "SBININBB123",
  "country_id": 1,
  "is_active": true,
  "created_at": "2026-03-03T10:20:00Z",
  "updated_at": "2026-03-03T10:20:00Z"
}
```

**Save**: `bank_id = 1`

---

## Step 6: Create Company (Simple)

Now that you have all the reference IDs, create the company.

**Endpoint**: `POST /api/companies`

### Request
```json
{
  "company_name": "Acme Corporation Private Limited",
  "business_type_id": 1,
  "gst_number": "12ABCDE1234F1Z5",
  "pan_number": "AAAPL1234C"
}
```

### Response (201 Created)
```json
{
  "id": 1,
  "company_name": "Acme Corporation Private Limited",
  "business_type_id": 1,
  "gst_number": "12ABCDE1234F1Z5",
  "pan_number": "AAAPL1234C",
  "is_active": true,
  "created_at": "2026-03-03T10:25:00Z",
  "updated_at": "2026-03-03T10:25:00Z"
}
```

**Save**: `company_id = 1`

---

## Step 7: Complete Company Setup (All-in-One)

Alternatively, create the company with all details in one request using the Complete Setup endpoint.

**Endpoint**: `POST /api/companies/complete-setup`

### Request (Full)
```json
{
  "company": {
    "company_name": "Acme Corporation Private Limited",
    "business_type_id": 1,
    "gst_number": "12ABCDE1234F1Z5",
    "pan_number": "AAAPL1234C"
  },
  "contact": {
    "mobile": "9876543210",
    "alternate_mobile": "9876543211",
    "email": "contact@acme.com"
  },
  "address": {
    "address_line1": "123 Business Street",
    "address_line2": "Suite 100, Tech Park",
    "city": "Mumbai",
    "state_id": 1,
    "country_id": 1,
    "pincode": "400001"
  },
  "bank_details": {
    "bank_id": 1,
    "account_holder_name": "Acme Corporation Private Limited",
    "account_number": "1234567890123456",
    "is_primary": true
  },
  "upi_details": {
    "upi_id": "acmecorp@upi",
    "upi_qr_url": "https://example.com/qr/acme-upi-qr.png"
  },
  "invoice_settings": {
    "invoice_prefix": "INV",
    "invoice_start_number": 1001,
    "show_logo": true,
    "show_signature": true,
    "round_off_total": false
  },
  "tax_settings": {
    "gst_enabled": true,
    "tax_type_id": 1
  },
  "regional_settings": {
    "timezone": "Asia/Kolkata",
    "date_format": "DD/MM/YYYY",
    "time_format": "HH:mm:ss",
    "currency_code": "INR",
    "currency_symbol": "₹",
    "language_code": "en-IN"
  }
}
```

### Response (201 Created)
```json
{
  "company": {
    "id": 1,
    "company_name": "Acme Corporation Private Limited",
    "business_type_id": 1,
    "gst_number": "12ABCDE1234F1Z5",
    "pan_number": "AAAPL1234C",
    "is_active": true,
    "created_at": "2026-03-03T10:30:00Z",
    "updated_at": "2026-03-03T10:30:00Z"
  },
  "contact": {
    "id": 1,
    "company_id": 1,
    "mobile": "9876543210",
    "alternate_mobile": "9876543211",
    "email": "contact@acme.com",
    "created_at": "2026-03-03T10:30:00Z",
    "updated_at": "2026-03-03T10:30:00Z"
  },
  "address": {
    "id": 1,
    "company_id": 1,
    "address_line1": "123 Business Street",
    "address_line2": "Suite 100, Tech Park",
    "city": "Mumbai",
    "state_id": 1,
    "country_id": 1,
    "pincode": "400001",
    "created_at": "2026-03-03T10:30:00Z",
    "updated_at": "2026-03-03T10:30:00Z"
  },
  "bank_details": {
    "id": 1,
    "company_id": 1,
    "bank_id": 1,
    "account_holder_name": "Acme Corporation Private Limited",
    "account_number": "1234567890123456",
    "is_primary": true,
    "is_active": true,
    "created_at": "2026-03-03T10:30:00Z",
    "updated_at": "2026-03-03T10:30:00Z"
  },
  "upi_details": {
    "id": 1,
    "company_id": 1,
    "upi_id": "acmecorp@upi",
    "upi_qr_url": "https://example.com/qr/acme-upi-qr.png",
    "created_at": "2026-03-03T10:30:00Z",
    "updated_at": "2026-03-03T10:30:00Z"
  },
  "invoice_settings": {
    "id": 1,
    "company_id": 1,
    "invoice_prefix": "INV",
    "invoice_start_number": 1001,
    "next_invoice_number": 1001,
    "show_logo": true,
    "show_signature": true,
    "round_off_total": false,
    "created_at": "2026-03-03T10:30:00Z",
    "updated_at": "2026-03-03T10:30:00Z"
  },
  "tax_settings": {
    "id": 1,
    "company_id": 1,
    "gst_enabled": true,
    "tax_type_id": 1,
    "created_at": "2026-03-03T10:30:00Z",
    "updated_at": "2026-03-03T10:30:00Z"
  },
  "regional_settings": {
    "id": 1,
    "company_id": 1,
    "timezone": "Asia/Kolkata",
    "date_format": "DD/MM/YYYY",
    "time_format": "HH:mm:ss",
    "currency_code": "INR",
    "currency_symbol": "₹",
    "language_code": "en-IN",
    "created_at": "2026-03-03T10:30:00Z",
    "updated_at": "2026-03-03T10:30:00Z"
  }
}
```

---

## Complete Setup with Minimal Fields

### Request (Minimal - Only Required Fields)
```json
{
  "company": {
    "company_name": "Acme Corporation",
    "business_type_id": 1
  },
  "contact": {
    "mobile": "9876543210",
    "email": "contact@acme.com"
  },
  "address": {
    "address_line1": "123 Business Street",
    "city": "Mumbai",
    "state_id": 1,
    "country_id": 1,
    "pincode": "400001"
  },
  "tax_settings": {
    "tax_type_id": 1
  },
  "regional_settings": {
    "timezone": "Asia/Kolkata",
    "date_format": "DD/MM/YYYY",
    "time_format": "HH:mm:ss",
    "currency_code": "INR",
    "currency_symbol": "₹",
    "language_code": "en-IN"
  }
}
```

### Response (201 Created)
```json
{
  "company": {
    "id": 1,
    "company_name": "Acme Corporation",
    "business_type_id": 1,
    "is_active": true,
    "created_at": "2026-03-03T10:30:00Z",
    "updated_at": "2026-03-03T10:30:00Z"
  },
  "contact": {
    "id": 1,
    "company_id": 1,
    "mobile": "9876543210",
    "email": "contact@acme.com",
    "created_at": "2026-03-03T10:30:00Z",
    "updated_at": "2026-03-03T10:30:00Z"
  },
  "address": {
    "id": 1,
    "company_id": 1,
    "address_line1": "123 Business Street",
    "city": "Mumbai",
    "state_id": 1,
    "country_id": 1,
    "pincode": "400001",
    "created_at": "2026-03-03T10:30:00Z",
    "updated_at": "2026-03-03T10:30:00Z"
  },
  "tax_settings": {
    "id": 1,
    "company_id": 1,
    "tax_type_id": 1,
    "created_at": "2026-03-03T10:30:00Z",
    "updated_at": "2026-03-03T10:30:00Z"
  },
  "regional_settings": {
    "id": 1,
    "company_id": 1,
    "timezone": "Asia/Kolkata",
    "date_format": "DD/MM/YYYY",
    "time_format": "HH:mm:ss",
    "currency_code": "INR",
    "currency_symbol": "₹",
    "language_code": "en-IN",
    "created_at": "2026-03-03T10:30:00Z",
    "updated_at": "2026-03-03T10:30:00Z"
  }
}
```

---

## cURL Examples

### 1. Create Country
```bash
curl -X POST http://localhost:3000/api/countries \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
  "country_name": "India",
  "country_code": "IN",
  "currency_code": "INR",
  "is_active": true
}'
```

### 2. Create State
```bash
curl -X POST http://localhost:3000/api/states \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
  "state_name": "Maharashtra",
  "state_code": "MH",
  "country_id": 1,
  "is_active": true
}'
```

### 3. Create Business Type
```bash
curl -X POST http://localhost:3000/api/business-types \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
  "business_type_name": "Private Limited",
  "description": "Private limited company",
  "is_active": true
}'
```

### 4. Create Tax Type
```bash
curl -X POST http://localhost:3000/api/tax-types \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
  "tax_type_name": "GST",
  "tax_description": "Goods and Services Tax",
  "tax_rate": 18.0,
  "is_active": true
}'
```

### 5. Create Bank
```bash
curl -X POST http://localhost:3000/api/banks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
  "bank_name": "State Bank of India",
  "bank_code": "SBIN0000001",
  "ifsc_code": "SBIN0001234",
  "swift_code": "SBININBB123",
  "country_id": 1,
  "is_active": true
}'
```

### 6. Create Company (Simple)
```bash
curl -X POST http://localhost:3000/api/companies \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
  "company_name": "Acme Corporation Private Limited",
  "business_type_id": 1,
  "gst_number": "12ABCDE1234F1Z5",
  "pan_number": "AAAPL1234C"
}'
```

### 7. Complete Company Setup
```bash
curl -X POST http://localhost:3000/api/companies/complete-setup \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
  "company": {
    "company_name": "Acme Corporation Private Limited",
    "business_type_id": 1,
    "gst_number": "12ABCDE1234F1Z5",
    "pan_number": "AAAPL1234C"
  },
  "contact": {
    "mobile": "9876543210",
    "alternate_mobile": "9876543211",
    "email": "contact@acme.com"
  },
  "address": {
    "address_line1": "123 Business Street",
    "address_line2": "Suite 100, Tech Park",
    "city": "Mumbai",
    "state_id": 1,
    "country_id": 1,
    "pincode": "400001"
  },
  "bank_details": {
    "bank_id": 1,
    "account_holder_name": "Acme Corporation Private Limited",
    "account_number": "1234567890123456",
    "is_primary": true
  },
  "upi_details": {
    "upi_id": "acmecorp@upi",
    "upi_qr_url": "https://example.com/qr/acme-upi-qr.png"
  },
  "invoice_settings": {
    "invoice_prefix": "INV",
    "invoice_start_number": 1001,
    "show_logo": true,
    "show_signature": true,
    "round_off_total": false
  },
  "tax_settings": {
    "gst_enabled": true,
    "tax_type_id": 1
  },
  "regional_settings": {
    "timezone": "Asia/Kolkata",
    "date_format": "DD/MM/YYYY",
    "time_format": "HH:mm:ss",
    "currency_code": "INR",
    "currency_symbol": "₹",
    "language_code": "en-IN"
  }
}'
```

---

## Validation Reference

### CreateCompanyInput Validation
```
company_name: required, min=1, max=255 characters
business_type_id: required, must be valid ID
gst_number: optional, must be 15 characters
pan_number: optional, must be 10 characters
```

### UpsertCompanyContactInput Validation
```
mobile: required, 10-15 characters
alternate_mobile: optional, 10-15 characters
email: required, valid email format, max=255 characters
```

### UpsertCompanyAddressInput Validation
```
address_line1: required, min=1, max=255 characters
address_line2: optional, max=255 characters
city: required, min=1, max=100 characters
state_id: required, must be valid ID
country_id: required, must be valid ID
pincode: required, 4-10 characters
```

### CreateBankDetailInput Validation
```
bank_id: required, must be valid ID
account_holder_name: required, min=1, max=255 characters
account_number: required, min=1, max=50 characters
is_primary: optional, boolean
```

### UpsertTaxSettingsInput Validation
```
gst_enabled: optional, boolean
tax_type_id: required, must be valid ID
```

### UpsertRegionalSettingsInput Validation
```
timezone: required, max=50 characters (e.g., Asia/Kolkata)
date_format: required, max=20 characters (e.g., DD/MM/YYYY)
time_format: required, max=20 characters (e.g., HH:mm:ss)
currency_code: required, exactly 3 characters (e.g., INR)
currency_symbol: required, max=10 characters (e.g., ₹)
language_code: required, max=5 characters (e.g., en-IN)
```

---

## Common Error Responses

### 400 Bad Request - Invalid Business Type ID
```json
{
  "error": "business_type_id must be a valid reference"
}
```

### 400 Bad Request - Invalid State ID
```json
{
  "error": "state_id must be a valid reference"
}
```

### 400 Bad Request - Invalid Country ID
```json
{
  "error": "country_id must be a valid reference"
}
```

### 400 Bad Request - Invalid Bank ID
```json
{
  "error": "bank_id must be a valid reference"
}
```

### 400 Bad Request - Invalid Tax Type ID
```json
{
  "error": "tax_type_id must be a valid reference"
}
```

### 400 Bad Request - Validation Error
```json
{
  "error": "company_name is required and must have between 1 and 255 characters"
}
```

### 400 Bad Request - Invalid Email
```json
{
  "error": "email must be a valid email address"
}
```

### 400 Bad Request - Invalid Mobile
```json
{
  "error": "mobile must be between 10 and 15 characters"
}
```

### 400 Bad Request - Invalid PIN Code
```json
{
  "error": "pincode must be between 4 and 10 characters"
}
```

---

## Setup Summary Table

| Step | Endpoint | Method | Status | Variable |
|------|----------|--------|--------|----------|
| 1 | `/api/countries` | POST | 201 | `country_id = 1` |
| 2 | `/api/states` | POST | 201 | `state_id = 1` |
| 3 | `/api/business-types` | POST | 201 | `business_type_id = 1` |
| 4 | `/api/tax-types` | POST | 201 | `tax_type_id = 1` |
| 5 | `/api/banks` | POST | 201 | `bank_id = 1` |
| 6 | `/api/companies` | POST | 201 | `company_id = 1` |
| 7 | `/api/companies/complete-setup` | POST | 201 | Full company setup |

---

## Tips

1. **Always create dependencies first**: Country → State → Business Type → Tax Type → Bank → Company
2. **Save IDs**: Store the returned IDs from each step to use in subsequent requests
3. **Validation**: Ensure all required fields are provided and match the validation rules
4. **Complete Setup**: Use the `/api/companies/complete-setup` endpoint for atomic transactions
5. **Minimal Setup**: Only required fields must be provided; optional fields can be added later
6. **Current User**: The system automatically associates the company with the authenticated user
7. **Tax Type**: At least one tax type must exist for GST functionality
8. **Regional Settings**: These are important for invoice generation and currency handling

---

## Support

For more information, refer to [COMPANY_CREATION.md](COMPANY_CREATION.md) for detailed API documentation.
