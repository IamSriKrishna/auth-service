# Employee Attendance System - API Usage Guide

This guide demonstrates how to create an employee and mark their attendance for single and multiple days.

## Prerequisites

- Server is running on `http://localhost:3000`
- User must be authenticated and have admin privileges
- Company ID is available from authentication token

---

## Step 1: Create an Employee

### Endpoint
```
POST /auth/manage/employees
```

### Request Headers
```
Authorization: Bearer <your_jwt_token>
Content-Type: application/json
```

### Request Body
```json
{
  "name": "John Doe",
  "email": "john.doe@company.com",
  "number": "9876543210",
  "address": "123 Main Street, City",
  "employee_type": "full-time"
}
```

### cURL Example
```bash
curl -X POST http://localhost:3000/auth/manage/employees \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john.doe@company.com",
    "number": "9876543210",
    "address": "123 Main Street, City",
    "employee_type": "full-time"
  }'
```

### Response (Success)
```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "John Doe",
    "email": "john.doe@company.com",
    "number": "9876543210",
    "address": "123 Main Street, City",
    "employee_type": "full-time",
    "user_id": 5,
    "company_id": 2,
    "created_at": "2026-03-05T10:30:00Z",
    "updated_at": "2026-03-05T10:30:00Z"
  }
}
```

**Note:** Save the `employee_id` (value: `1`) for the next steps.

---

## Step 2: Mark Attendance for a Single Day

### Endpoint
```
POST /auth/manage/attendance
```

### Request Headers
```
Authorization: Bearer <your_jwt_token>
Content-Type: application/json
```

### Request Body - Scenario A: Employee Present (On Time)
```json
{
  "employee_id": 1,
  "date": "2026-03-05",
  "status": "on_time",
  "reason": "",
  "check_in_time": "2026-03-05T09:00:00Z",
  "check_out_time": "2026-03-05T18:00:00Z",
  "working_hours": 9.0,
  "notes": "Regular working day"
}
```

### Request Body - Scenario B: Employee Absent
```json
{
  "employee_id": 1,
  "date": "2026-03-05",
  "status": "absent",
  "reason": "Sick leave",
  "check_in_time": null,
  "check_out_time": null,
  "working_hours": 0,
  "notes": "Not feeling well"
}
```

### Request Body - Scenario C: Employee Late
```json
{
  "employee_id": 1,
  "date": "2026-03-05",
  "status": "late",
  "reason": "Traffic",
  "check_in_time": "2026-03-05T10:30:00Z",
  "check_out_time": "2026-03-05T18:00:00Z",
  "working_hours": 7.5,
  "notes": "Arrived 1.5 hours late"
}
```

### Request Body - Scenario D: Holiday
```json
{
  "employee_id": 1,
  "date": "2026-03-05",
  "status": "holiday",
  "reason": "National Holiday",
  "check_in_time": null,
  "check_out_time": null,
  "working_hours": 0,
  "notes": "Independence Day"
}
```

### Request Body - Scenario E: Half Day
```json
{
  "employee_id": 1,
  "date": "2026-03-05",
  "status": "half_day",
  "reason": "Personal work",
  "check_in_time": "2026-03-05T09:00:00Z",
  "check_out_time": "2026-03-05T13:00:00Z",
  "working_hours": 4.0,
  "notes": "Left early for personal appointment"
}
```

### Request Body - Scenario F: Leave
```json
{
  "employee_id": 1,
  "date": "2026-03-05",
  "status": "leave",
  "reason": "Approved leave",
  "check_in_time": null,
  "check_out_time": null,
  "working_hours": 0,
  "notes": "Casual leave approved by manager"
}
```

### cURL Example - On Time
```bash
curl -X POST http://localhost:3000/auth/manage/attendance \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "employee_id": 1,
    "date": "2026-03-05",
    "status": "on_time",
    "reason": "",
    "check_in_time": "2026-03-05T09:00:00Z",
    "check_out_time": "2026-03-05T18:00:00Z",
    "working_hours": 9.0,
    "notes": "Regular working day"
  }'
```

### Response (Success)
```json
{
  "success": true,
  "data": {
    "id": 1,
    "employee_id": 1,
    "company_id": 2,
    "date": "2026-03-05",
    "status": "on_time",
    "reason": "",
    "check_in_time": "2026-03-05T09:00:00Z",
    "check_out_time": "2026-03-05T18:00:00Z",
    "working_hours": 9.0,
    "notes": "Regular working day",
    "created_at": "2026-03-05T10:35:00Z",
    "updated_at": "2026-03-05T10:35:00Z"
  }
}
```

---

## Step 3: Mark Attendance for Multiple Days

Mark attendance for multiple consecutive days (e.g., a week).

### Day 1: March 3, 2026 - On Time
```bash
curl -X POST http://localhost:3000/auth/manage/attendance \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "employee_id": 1,
    "date": "2026-03-03",
    "status": "on_time",
    "reason": "",
    "check_in_time": "2026-03-03T09:00:00Z",
    "check_out_time": "2026-03-03T18:00:00Z",
    "working_hours": 9.0,
    "notes": "Regular working day"
  }'
```

### Day 2: March 4, 2026 - Late
```bash
curl -X POST http://localhost:3000/auth/manage/attendance \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "employee_id": 1,
    "date": "2026-03-04",
    "status": "late",
    "reason": "Traffic jam",
    "check_in_time": "2026-03-04T10:00:00Z",
    "check_out_time": "2026-03-04T18:00:00Z",
    "working_hours": 8.0,
    "notes": "Arrived 1 hour late due to traffic"
  }'
```

### Day 3: March 5, 2026 - On Time
```bash
curl -X POST http://localhost:3000/auth/manage/attendance \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "employee_id": 1,
    "date": "2026-03-05",
    "status": "on_time",
    "reason": "",
    "check_in_time": "2026-03-05T09:00:00Z",
    "check_out_time": "2026-03-05T18:00:00Z",
    "working_hours": 9.0,
    "notes": "Regular working day"
  }'
```

### Day 4: March 6, 2026 - Half Day
```bash
curl -X POST http://localhost:3000/auth/manage/attendance \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "employee_id": 1,
    "date": "2026-03-06",
    "status": "half_day",
    "reason": "Doctor appointment",
    "check_in_time": "2026-03-06T09:00:00Z",
    "check_out_time": "2026-03-06T13:30:00Z",
    "working_hours": 4.5,
    "notes": "Left for medical appointment"
  }'
```

### Day 5: March 7, 2026 - Absent
```bash
curl -X POST http://localhost:3000/auth/manage/attendance \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "employee_id": 1,
    "date": "2026-03-07",
    "status": "absent",
    "reason": "Sick leave",
    "check_in_time": null,
    "check_out_time": null,
    "working_hours": 0,
    "notes": "No show, reported sick via SMS"
  }'
```

---

## Step 4: Retrieve Attendance Records

### Get Attendance for a Specific Employee

```bash
curl -X GET "http://localhost:3000/auth/manage/attendance/employee/1?page=1&limit=10" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Get Attendance by Date Range

```bash
curl -X GET "http://localhost:3000/auth/manage/attendance/date-range?from_date=2026-03-03&to_date=2026-03-07&page=1&limit=10" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Get Attendance for Employee in Date Range

```bash
curl -X GET "http://localhost:3000/auth/manage/attendance/employee/1/date-range?from_date=2026-03-03&to_date=2026-03-07&page=1&limit=10" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Get All Attendance Records (Company)

```bash
curl -X GET "http://localhost:3000/auth/manage/attendance?page=1&limit=10" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Get Attendance Statistics

```bash
curl -X GET "http://localhost:3000/auth/manage/attendance/stats/report?from_date=2026-03-03&to_date=2026-03-07" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Response Example
```json
{
  "success": true,
  "data": {
    "total": 5,
    "on_time": 2,
    "absent": 1,
    "late": 1,
    "holiday": 0,
    "half_day": 1,
    "leave": 0
  }
}
```

---

## Step 5: Quick Check-In / Check-Out

### Check-In Employee (Today)

```bash
curl -X POST http://localhost:3000/auth/manage/attendance/checkin/1 \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json"
```

### Check-Out Employee (Today)

```bash
curl -X POST http://localhost:3000/auth/manage/attendance/checkout/1 \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json"
```

### Response Example (Check-In)
```json
{
  "success": true,
  "data": {
    "id": 5,
    "employee_id": 1,
    "company_id": 2,
    "date": "2026-03-05",
    "status": "on_time",
    "reason": "",
    "check_in_time": "2026-03-05T09:15:30Z",
    "check_out_time": null,
    "working_hours": 0,
    "notes": "",
    "created_at": "2026-03-05T09:15:30Z",
    "updated_at": "2026-03-05T09:15:30Z"
  }
}
```

---

## Step 6: Update Attendance Record

### Endpoint
```
PUT /auth/manage/attendance/{attendance_id}
```

### Example: Update Status to Late

```bash
curl -X PUT http://localhost:3000/auth/manage/attendance/1 \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "late",
    "reason": "Traffic",
    "notes": "Updated: Verified traffic report"
  }'
```

---

## Attendance Status Values

| Status | Description | Use Case |
|--------|-------------|----------|
| `on_time` | Present on time | Regular attendance |
| `absent` | Not present | No call, no show |
| `late` | Arrived late | Delayed arrival |
| `holiday` | Company holiday | Public/company holiday |
| `half_day` | Half working day | Early leave, medical |
| `leave` | Approved leave | Vacation, sick leave |

---

## Error Responses

### Employee Not Found
```json
{
  "error": true,
  "message": "employee not found"
}
```

### Attendance Already Exists for Date
```json
{
  "error": true,
  "message": "attendance already recorded for this date"
}
```

### No Check-In Found for Check-Out
```json
{
  "error": true,
  "message": "no check-in found for today"
}
```

### Unauthorized Access
```json
{
  "error": true,
  "message": "unauthorized access to attendance record"
}
```

---

## Complete Flow Example

### 1. Create Employee
```bash
# Response: employee_id = 1
POST /auth/manage/employees
```

### 2. Mark Attendance for 3 Days
```bash
# Day 1: On Time
POST /auth/manage/attendance { date: 2026-03-05, status: on_time, ... }

# Day 2: Late
POST /auth/manage/attendance { date: 2026-03-06, status: late, ... }

# Day 3: Absent
POST /auth/manage/attendance { date: 2026-03-07, status: absent, ... }
```

### 3. Retrieve Records
```bash
GET /auth/manage/attendance/employee/1?page=1&limit=10
```

### 4. View Statistics
```bash
GET /auth/manage/attendance/stats/report?from_date=2026-03-05&to_date=2026-03-07
```

---

## Testing with Postman

1. Create a new collection named "Employee Attendance"
2. Set up authorization:
   - Type: Bearer Token
   - Token: YOUR_JWT_TOKEN
3. Create requests for each endpoint
4. Use the examples above as request bodies

---

## Notes

- All dates should be in `YYYY-MM-DD` format
- All times should be in ISO 8601 format (e.g., `2026-03-05T09:00:00Z`)
- Working hours are calculated automatically in check-in/check-out flow
- One attendance record per employee per day is allowed
- Edit/update operations available for corrections
- Delete operations soft-delete records (preserves audit trail)
