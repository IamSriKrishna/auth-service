# Complete Employee Requests & Responses

## Overview
This guide shows the **complete request/response** structure for creating, retrieving (by user), and updating employees.

---

## Prerequisites & Setup

### Required Data Before Creating Employees

Before you can create employees, ensure the following data exists in your database:

#### 1. User Must Exist
The authenticated user must exist in the `users` table with correct `company_id`.

**Check user details:**
```sql
SELECT id, email, company_id, user_type FROM users WHERE id = 17;
```

#### 2. Company Must Exist
The authenticated user's `company_id` (from JWT token) must reference an existing company in the `companies` table.

**Check existing companies:**
```sql
SELECT id, name, company_code FROM companies;
```

**Create a company if needed:**
```sql
INSERT INTO companies (id, name, company_code, status, created_at, updated_at) 
VALUES (1, 'Your Company Name', 'COMP001', 'active', NOW(), NOW());
```

#### 3. Authentication Token
All requests require a valid JWT token in the Authorization header with proper `company_id` and `user_id` claims.

**JWT Token Inspection:**
```bash
# Decode JWT (use jwt.io or similar)
# Look for "company_id" and "user_id" fields
# Example decoded payload:
{
  "company_id": 1,
  "email": "srik2@company.com",
  "user_id": 17,
  ...
}
```

---

## Common Errors & Solutions

### Error: Unauthorized Access to Employee

**Error Message:**
```
Unauthorized access to employee
```

**Cause:**
The employee was created by a different user. Only the user who created the employee can retrieve or update it.

**Solution:**
1. Verify the `user_id` in your JWT token matches the creator
2. Use the correct authentication token for the user who created the employee
3. Check if the employee exists: `SELECT id, user_id FROM employees WHERE id = 1;`

### Error: Employee Not Found

**Error Message:**
```
Employee not found
```

**Cause:**
The employee with the provided ID doesn't exist in the database.

**Solution:**
1. Verify the employee ID is correct
2. Check existing employees: `SELECT id, name, user_id FROM employees WHERE user_id = 17;`
3. Create the employee first before trying to retrieve it

---

## CREATE EMPLOYEE

### Create Employee Request

```http
POST /auth/employees
Content-Type: application/json
Authorization: Bearer {JWT_TOKEN}
```

### Request Body - Complete Employee

```json
{
  "name": "John Doe",
  "email": "john.doe@company.com",
  "number": "9876543210",
  "address": "123 Business Street, Tech Park, City, State 560001, India"
}
```

### Request Field Breakdown

> **Note**: The employee is created by the authenticated user. The `user_id` and `company_id` are automatically set from the JWT token - they should NOT be included in the request body.

#### Main Employee Details
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| name | string | Yes | Employee full name |
| email | string | No | Email address |
| number | string | Yes | Contact phone number |
| address | string | Yes | Residential/office address |

### Create Employee Response (Success)

```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "John Doe",
    "email": "john.doe@company.com",
    "number": "9876543210",
    "address": "123 Business Street, Tech Park, City, State 560001, India",
    "user_id": 17,
    "company_id": 1,
    "created_at": "2026-03-04T19:30:00.000+05:30",
    "updated_at": "2026-03-04T19:30:00.000+05:30"
  }
}
```

### cURL Command - Create Employee

```bash
curl -X POST http://localhost:3000/auth/employees \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "John Doe",
    "email": "john.doe@company.com",
    "number": "9876543210",
    "address": "123 Business Street, Tech Park, City, State 560001, India"
  }'
```

---

## GET EMPLOYEE BY ID

### Get Employee by ID Request

```http
GET /auth/employees/:id
Content-Type: application/json
Authorization: Bearer {JWT_TOKEN}
```

### Get Employee by ID Response (Success)

```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "John Doe",
    "email": "john.doe@company.com",
    "number": "9876543210",
    "address": "123 Business Street, Tech Park, City, State 560001, India",
    "user_id": 17,
    "company_id": 1,
    "created_at": "2026-03-04T19:30:00.000+05:30",
    "updated_at": "2026-03-04T19:30:00.000+05:30"
  }
}
```

### Get Employee by ID Response (Unauthorized)

```json
{
  "error": true,
  "message": "Unauthorized access to employee"
}
```

### Get Employee by ID Response (Not Found)

```json
{
  "error": true,
  "message": "Employee not found"
}
```

### cURL Command - Get Employee by ID

```bash
curl -X GET http://localhost:3000/auth/employees/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

## GET EMPLOYEES BY CREATOR USER (Paginated)

### Get All Employees Created by User Request

```http
GET /auth/employees?page=1&limit=10
Content-Type: application/json
Authorization: Bearer {JWT_TOKEN}
```

### Query Parameters
| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| page | integer | 1 | Page number for pagination |
| limit | integer | 10 | Number of records per page |

### Get All Employees Created by User Response (Success)

```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "John Doe",
      "email": "john.doe@company.com",
      "number": "9876543210",
      "address": "123 Business Street, Tech Park, City, State 560001, India",
      "user_id": 17,
      "company_id": 1,
      "created_at": "2026-03-04T19:30:00.000+05:30"
    },
    {
      "id": 2,
      "name": "Jane Smith",
      "email": "jane.smith@company.com",
      "number": "9876543211",
      "address": "456 Corporate Avenue, Tech Hub, City, State 560002, India",
      "user_id": 17,
      "company_id": 1,
      "created_at": "2026-03-04T19:35:15.000+05:30"
    }
  ],
  "meta": {
    "current_page": 1,
    "per_page": 10,
    "total": 2,
    "total_pages": 1
  }
}
```

### Get All Employees - Empty Result Response

```json
{
  "success": true,
  "data": [],
  "meta": {
    "current_page": 1,
    "per_page": 10,
    "total": 0,
    "total_pages": 0
  }
}
```

### cURL Command - Get All Employees by User

```bash
curl -X GET "http://localhost:3000/auth/employees?page=1&limit=10" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### cURL Command - Get Employees with Custom Pagination

```bash
curl -X GET "http://localhost:3000/auth/employees?page=2&limit=20" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

## UPDATE EMPLOYEE

### Update Employee Request

```http
PUT /auth/employees/:id
Content-Type: application/json
Authorization: Bearer {JWT_TOKEN}
```

### Request Body - Partial Update Example

```json
{
  "name": "John Doe Updated",
  "email": "john.doe.updated@company.com",
  "number": "9876543215",
  "address": "789 New Street, Tech District, City, State 560003, India"
}
```

### Request Body - Update Only Some Fields

```json
{
  "email": "john.doe.newemail@company.com",
  "number": "9876543215"
}
```

### Update Employee Response (Success)

```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "John Doe Updated",
    "email": "john.doe.updated@company.com",
    "number": "9876543215",
    "address": "789 New Street, Tech District, City, State 560003, India",
    "user_id": 17,
    "company_id": 1,
    "created_at": "2026-03-04T19:30:00.000+05:30",
    "updated_at": "2026-03-04T20:45:30.000+05:30"
  }
}
```

### Update Employee Response (Unauthorized)

```json
{
  "error": true,
  "message": "Unauthorized access to employee"
}
```

### Update Employee Response (Not Found)

```json
{
  "error": true,
  "message": "Employee not found"
}
```

### cURL Command - Update Employee (Full Update)

```bash
curl -X PUT http://localhost:3000/auth/employees/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "John Doe Updated",
    "email": "john.doe.updated@company.com",
    "number": "9876543215",
    "address": "789 New Street, Tech District, City, State 560003, India"
  }'
```

### cURL Command - Update Employee (Partial Update)

```bash
curl -X PUT http://localhost:3000/auth/employees/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "email": "john.doe.newemail@company.com",
    "number": "9876543215"
  }'
```

---

## DELETE EMPLOYEE

### Delete Employee Request

```http
DELETE /auth/employees/:id
Authorization: Bearer {JWT_TOKEN}
```

### Delete Employee Response (Success)

```json
{
  "success": true,
  "message": "Employee deleted successfully"
}
```

### Delete Employee Response (Unauthorized)

```json
{
  "error": true,
  "message": "Unauthorized access to employee"
}
```

### Delete Employee Response (Not Found)

```json
{
  "error": true,
  "message": "Employee not found"
}
```

### cURL Command - Delete Employee

```bash
curl -X DELETE http://localhost:3000/auth/employees/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

## Data Types Reference

| Type | Example | Usage |
|------|---------|-------|
| string | "John Doe" | Text fields (name, email, address) |
| integer | 1, 17 | IDs, page numbers |
| timestamp | "2026-03-04T19:30:00.000+05:30" | Creation and update times |

---

## Validation Rules

### Required Fields
- **name**: Must be provided, minimum 1 character
- **number**: Must be provided, valid phone number
- **address**: Must be provided, minimum 1 character

### Optional Fields
- **email**: Optional, but must be valid email format if provided

### Format Validations
- Email must be valid email format (e.g., user@domain.com)
- Phone number should be 10-15 digits
- Address can contain special characters and line breaks

---

## Employee Lifecycle

```
CREATE EMPLOYEE
       ↓
GET EMPLOYEE BY ID (verify access)
       ↓
UPDATE EMPLOYEE (modify details)
       ↓
DELETE EMPLOYEE (remove from system)
```

### Authorization Model
- Only the **user who created the employee** can retrieve or update that employee
- Other users in the same company cannot access employees created by different users
- This ensures data isolation between team members

---

## Examples by Employee Type

### Full-Time Employee
```json
{
  "name": "Rajesh Kumar",
  "email": "rajesh.kumar@company.com",
  "number": "9876543210",
  "address": "123 Residential Complex, Tech City, Bangalore, Karnataka 560001, India"
}
```

### Part-Time Employee
```json
{
  "name": "Priya Sharma",
  "email": "priya.sharma@company.com",
  "number": "9876543211",
  "address": "456 Apartment, MG Road, Bangalore, Karnataka 560002, India"
}
```

### Consultant/Contractor
```json
{
  "name": "Amit Patel",
  "email": "amit@consultingfirm.com",
  "number": "9876543212",
  "address": "789 Business Center, Whitefield, Bangalore, Karnataka 560066, India"
}
```

---

## Common Use Cases

### Create Multiple Employees

Create employees one by one using the Create Employee endpoint:

```bash
# Employee 1
curl -X POST http://localhost:3000/auth/employees \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"name": "Employee 1", "email": "emp1@company.com", "number": "9876543210", "address": "Address 1"}'

# Employee 2
curl -X POST http://localhost:3000/auth/employees \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"name": "Employee 2", "email": "emp2@company.com", "number": "9876543211", "address": "Address 2"}'
```

### List All Employees with Pagination

```bash
# Get first page
curl -X GET "http://localhost:3000/auth/employees?page=1&limit=10" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Get second page
curl -X GET "http://localhost:3000/auth/employees?page=2&limit=10" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Get 25 employees per page
curl -X GET "http://localhost:3000/auth/employees?page=1&limit=25" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Update Employee Email

```bash
curl -X PUT http://localhost:3000/auth/employees/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"email": "newemail@company.com"}'
```

### Update Employee Contact Number

```bash
curl -X PUT http://localhost:3000/auth/employees/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"number": "9999999999"}'
```

---

## Response Structure

All responses follow a consistent structure:

### Success Response
```json
{
  "success": true,
  "data": { ... employee data ... }
}
```

### Error Response
```json
{
  "error": true,
  "message": "Error description"
}
```

### Paginated Response
```json
{
  "success": true,
  "data": [ ... array of employees ... ],
  "meta": {
    "current_page": 1,
    "per_page": 10,
    "total": 25,
    "total_pages": 3
  }
}
```

---

## Important Notes

### User-Specific Access
- Each employee is associated with the user who created it
- Only that user can retrieve or modify the employee record
- This prevents unauthorized access across teams

### Employee IDs
- Unique within the system
- Cannot be changed after creation
- Used to identify and retrieve specific employees

### Timestamps
- `created_at`: Set automatically when employee is created
- `updated_at`: Updated automatically when employee is modified
- Format: ISO 8601 with timezone (e.g., 2026-03-04T19:30:00.000+05:30)

### Pagination
- Default limit: 10 records per page
- Maximum recommended limit: 100 records
- Useful for displaying employees in UI with load more functionality

---

Generated: March 4, 2026
Updated: Complete Employee Request/Response Documentation
