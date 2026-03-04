# Employee Management API Documentation

## Overview
The Employee Management API allows admin users to create, read, update, and delete employees. Only users with **admin** role can manage employees, and each user can only access employees they created.

---

## Base URL
```
http://localhost:3000
```

## Authentication
All employee endpoints require:
- **Authorization Header**: `Bearer {JWT_TOKEN}`
- **User Role**: `admin` or `superadmin`

---

## Endpoints

### 1. Create Employee
**Create a new employee record**

#### Request
```http
POST /auth/manage/employees
Content-Type: application/json
Authorization: Bearer {JWT_TOKEN}
```

#### Request Body
```json
{
  "name": "John Smith",
  "email": "john.smith@company.com",
  "number": "9876543210",
  "address": "123 Main Street, New York, NY 10001"
}
```

#### Request Fields
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| name | string | Yes | Employee's full name |
| email | string | No | Employee's email address (must be valid email format) |
| number | string | Yes | Employee's phone number |
| address | string | Yes | Employee's address |

#### Success Response
```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "John Smith",
    "email": "john.smith@company.com",
    "number": "9876543210",
    "address": "123 Main Street, New York, NY 10001",
    "user_id": 9,
    "company_id": 3,
    "created_at": "2026-03-04T10:30:45.123+05:30",
    "updated_at": "2026-03-04T10:30:45.123+05:30"
  }
}
```

#### Response Fields
| Field | Type | Description |
|-------|------|-------------|
| success | boolean | Always true on success |
| data.id | integer | Unique employee ID |
| data.name | string | Employee's name |
| data.email | string | Employee's email |
| data.number | string | Employee's phone number |
| data.address | string | Employee's address |
| data.user_id | integer | ID of the user who created this employee |
| data.company_id | integer | Company ID of the creator |
| data.created_at | string | Timestamp when created (ISO 8601) |
| data.updated_at | string | Timestamp when last updated (ISO 8601) |

#### Error Response
```json
{
  "error": true,
  "message": "request cannot be nil"
}
```

---

### 2. Get All Employees
**Retrieve all employees created by the authenticated user with pagination**

#### Request
```http
GET /auth/manage/employees?page=1&limit=10
Content-Type: application/json
Authorization: Bearer {JWT_TOKEN}
```

#### Query Parameters
| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| page | integer | 1 | Page number for pagination |
| limit | integer | 10 | Number of results per page |

#### Success Response
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "John Smith",
      "email": "john.smith@company.com",
      "number": "9876543210",
      "address": "123 Main Street, New York, NY 10001",
      "user_id": 9,
      "company_id": 3,
      "created_at": "2026-03-04T10:30:45.123+05:30"
    },
    {
      "id": 2,
      "name": "Jane Doe",
      "email": "jane.doe@company.com",
      "number": "9876543211",
      "address": "456 Oak Avenue, Boston, MA 02101",
      "user_id": 9,
      "company_id": 3,
      "created_at": "2026-03-04T11:15:30.456+05:30"
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

---

### 3. Get Employee by ID
**Retrieve a specific employee by ID**

#### Request
```http
GET /auth/manage/employees/1
Content-Type: application/json
Authorization: Bearer {JWT_TOKEN}
```

#### Path Parameters
| Parameter | Type | Description |
|-----------|------|-------------|
| id | integer | Employee ID |

#### Success Response
```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "John Smith",
    "email": "john.smith@company.com",
    "number": "9876543210",
    "address": "123 Main Street, New York, NY 10001",
    "user_id": 9,
    "company_id": 3,
    "created_at": "2026-03-04T10:30:45.123+05:30",
    "updated_at": "2026-03-04T10:30:45.123+05:30"
  }
}
```

#### Error Response
```json
{
  "error": true,
  "message": "unauthorized access to employee"
}
```

---

### 4. Update Employee
**Update an existing employee's information**

#### Request
```http
PUT /auth/manage/employees/1
Content-Type: application/json
Authorization: Bearer {JWT_TOKEN}
```

#### Request Body (All fields optional)
```json
{
  "name": "John Michael Smith",
  "email": "john.m.smith@company.com",
  "number": "9876543210",
  "address": "789 Pine Lane, Chicago, IL 60601"
}
```

#### Request Fields
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| name | string | No | Employee's full name |
| email | string | No | Employee's email address |
| number | string | No | Employee's phone number |
| address | string | No | Employee's address |

#### Success Response
```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "John Michael Smith",
    "email": "john.m.smith@company.com",
    "number": "9876543210",
    "address": "789 Pine Lane, Chicago, IL 60601",
    "user_id": 9,
    "company_id": 3,
    "created_at": "2026-03-04T10:30:45.123+05:30",
    "updated_at": "2026-03-04T12:45:20.789+05:30"
  }
}
```

---

### 5. Delete Employee
**Delete an employee**

#### Request
```http
DELETE /auth/manage/employees/1
Content-Type: application/json
Authorization: Bearer {JWT_TOKEN}
```

#### Path Parameters
| Parameter | Type | Description |
|-----------|------|-------------|
| id | integer | Employee ID |

#### Success Response
```json
{
  "success": true,
  "message": "Employee deleted successfully"
}
```

#### Error Response
```json
{
  "error": true,
  "message": "employee not found"
}
```

---

## cURL Examples

### Create Employee
```bash
curl -X POST http://localhost:3000/auth/manage/employees \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "John Smith",
    "email": "john.smith@company.com",
    "number": "9876543210",
    "address": "123 Main Street, New York, NY 10001"
  }'
```

### Get All Employees
```bash
curl -X GET "http://localhost:3000/auth/manage/employees?page=1&limit=10" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Get Single Employee
```bash
curl -X GET http://localhost:3000/auth/manage/employees/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Update Employee
```bash
curl -X PUT http://localhost:3000/auth/manage/employees/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "John Michael Smith",
    "email": "john.m.smith@newcompany.com",
    "number": "9876543999"
  }'
```

### Delete Employee
```bash
curl -X DELETE http://localhost:3000/auth/manage/employees/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

## Postman Collection

### Create Employee Request
```json
{
  "name": "Create Employee",
  "request": {
    "method": "POST",
    "header": [
      {
        "key": "Content-Type",
        "value": "application/json"
      },
      {
        "key": "Authorization",
        "value": "Bearer {{jwt_token}}"
      }
    ],
    "body": {
      "mode": "raw",
      "raw": "{\n  \"name\": \"John Smith\",\n  \"email\": \"john.smith@company.com\",\n  \"number\": \"9876543210\",\n  \"address\": \"123 Main Street, New York, NY 10001\"\n}"
    },
    "url": {
      "raw": "{{base_url}}/auth/manage/employees",
      "host": ["{{base_url}}"],
      "path": ["auth", "manage", "employees"]
    }
  }
}
```

---

## Important Notes

### Authorization
- Only users with `admin` or `superadmin` role can access employee endpoints
- Users can only view/modify employees they created
- Attempting to access another user's employees will return `401 Unauthorized`

### Data Isolation
- When you create an employee, it's automatically linked to:
  - Your `user_id` (creator)
  - Your `company_id` (from JWT token)
- When you fetch employees, you only see employees you created in your company

### Validation
- Email must be in valid email format (if provided)
- Name and address are required for creation
- Phone number is required for creation

### Response Format
- All successful responses have `success: true`
- All error responses have `error: true`
- Timestamps are in ISO 8601 format with timezone

---

## Error Responses

### 400 Bad Request
```json
{
  "error": true,
  "message": "Invalid request body"
}
```

### 401 Unauthorized (Invalid Token)
```json
{
  "error": true,
  "message": "Invalid or expired token"
}
```

### 403 Forbidden (Insufficient Permissions)
```json
{
  "error": true,
  "message": "Admin access required"
}
```

### 404 Not Found
```json
{
  "error": true,
  "message": "employee not found"
}
```

---

## Data Model

```
Employee {
  id: uint (Primary Key)
  name: string (Required) - Max 255 characters
  email: string (Optional) - Valid email format
  number: string (Required) - Phone number
  address: string (Required) - Address text
  user_id: uint (Required) - Creator user ID
  company_id: uint (Required) - Creator's company ID
  created_at: timestamp
  updated_at: timestamp
  deleted_at: timestamp (soft delete)
}
```

---

## API Usage Flow

```
1. User logs in and gets JWT token
2. User makes API call to /auth/manage/employees with Bearer token
3. Auth Middleware validates token and extracts user_id and company_id
4. Admin Middleware checks if user has admin role
5. Handler receives request and calls service
6. Service validates data and calls repository
7. Repository performs database operation
8. Response is returned with employee data
```

---

## Testing Tips

1. **Save JWT Token**: After login, save the JWT token from the response
2. **Include in Headers**: Always include the token in `Authorization: Bearer {token}`
3. **Test Pagination**: Try different page and limit values
4. **Handle Errors**: Check error messages for debugging
5. **Verify Isolation**: Create employees with different users and verify you can't access them

---

Generated: March 4, 2026
