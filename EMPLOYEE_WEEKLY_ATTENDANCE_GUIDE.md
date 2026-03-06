# Employee Weekly Attendance Management Guide

This guide walks you through the complete process of creating an employee and managing their weekly attendance records.

## Prerequisites

- Running auth-service with database migrations applied
- Valid authentication token (Bearer token)
- Admin access for employee and attendance management
- Date format: `YYYY-MM-DD` (e.g., `2026-03-05`)

---

## Step 1: Create a Company

First, create a company that will manage employees.

**Endpoint:** `POST /auth/manage/company`

**Request:**
```bash
curl -X POST http://localhost:3000/auth/manage/company \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "ABC Corporation",
    "email": "info@abccorp.com",
    "phone": "+1234567890",
    "address": "123 Business Street, City, Country",
    "city": "City Name",
    "state": "State",
    "country": "Country",
    "postal_code": "12345"
  }'
```

**Response:**
```json
{
  "id": 1,
  "name": "ABC Corporation",
  "email": "info@abccorp.com",
  "phone": "+1234567890",
  "address": "123 Business Street, City, Country",
  "city": "City Name",
  "state": "State",
  "country": "Country",
  "postal_code": "12345"
}
```

Save the **company ID** (e.g., `1`) for the next steps.

---

## Step 2: Create an Employee

Create an employee record for the company.

**Endpoint:** `POST /auth/manage/employee`

**Request:**
```bash
curl -X POST http://localhost:3000/auth/manage/employee \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "company_id": 1,
    "user_id": 2,
    "first_name": "John",
    "last_name": "Doe",
    "email": "john.doe@abccorp.com",
    "phone": "+1111111111",
    "position": "Senior Developer",
    "department": "Engineering",
    "employee_id": "EMP001",
    "joining_date": "2025-01-15",
    "address": "456 Employee Lane, City, Country",
    "city": "City Name",
    "state": "State",
    "country": "Country",
    "postal_code": "12345"
  }'
```

**Response:**
```json
{
  "id": 5,
  "company_id": 1,
  "user_id": 2,
  "first_name": "John",
  "last_name": "Doe",
  "email": "john.doe@abccorp.com",
  "phone": "+1111111111",
  "position": "Senior Developer",
  "department": "Engineering",
  "employee_id": "EMP001",
  "joining_date": "2025-01-15"
}
```

Save the **employee ID** (e.g., `5`) for marking attendance.

---

## Step 3: Mark Attendance for a Single Day

Mark attendance for a specific day.

**Endpoint:** `POST /auth/manage/attendance`

**Attendance Statuses:**
- `on_time` - Employee arrived on time
- `late` - Employee arrived late
- `absent` - Employee was absent
- `holiday` - Official holiday (no work expected)
- `half_day` - Employee worked only half day
- `leave` - Employee took leave

**Request:**
```bash
curl -X POST http://localhost:3000/auth/manage/attendance \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "employee_id": 5,
    "company_id": 1,
    "date": "2026-03-03",
    "status": "on_time",
    "check_in_time": "2026-03-03T09:00:00Z",
    "check_out_time": "2026-03-03T18:00:00Z",
    "working_hours": 9.0,
    "reason": "Regular working day",
    "notes": "Good productive day"
  }'
```

**Response:**
```json
{
  "id": 1,
  "employee_id": 5,
  "company_id": 1,
  "date": "2026-03-03",
  "status": "on_time",
  "check_in_time": "2026-03-03T09:00:00Z",
  "check_out_time": "2026-03-03T18:00:00Z",
  "working_hours": 9.0,
  "reason": "Regular working day",
  "notes": "Good productive day"
}
```

---

## Step 4: Mark Attendance for Multiple Days (Full Week)

Use bulk attendance creation to mark attendance for an entire week in a single request.

**Endpoint:** `POST /auth/manage/attendance/bulk`

**Request (Full Week: Mar 1-7, 2026):**
```bash
curl -X POST http://localhost:3000/auth/manage/attendance/bulk \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "employee_id": 5,
    "company_id": 1,
    "attendances": [
      {
        "date": "2026-03-01",
        "status": "on_time",
        "check_in_time": "2026-03-01T09:00:00Z",
        "check_out_time": "2026-03-01T18:00:00Z",
        "working_hours": 9.0,
        "reason": "Regular working day"
      },
      {
        "date": "2026-03-02",
        "status": "on_time",
        "check_in_time": "2026-03-02T09:15:00Z",
        "check_out_time": "2026-03-02T18:00:00Z",
        "working_hours": 8.75,
        "reason": "Regular working day"
      },
      {
        "date": "2026-03-03",
        "status": "late",
        "check_in_time": "2026-03-03T09:30:00Z",
        "check_out_time": "2026-03-03T18:00:00Z",
        "working_hours": 8.5,
        "reason": "Traffic delay",
        "notes": "20 minutes late"
      },
      {
        "date": "2026-03-04",
        "status": "on_time",
        "check_in_time": "2026-03-04T09:00:00Z",
        "check_out_time": "2026-03-04T18:00:00Z",
        "working_hours": 9.0,
        "reason": "Regular working day"
      },
      {
        "date": "2026-03-05",
        "status": "half_day",
        "check_in_time": "2026-03-05T09:00:00Z",
        "check_out_time": "2026-03-05T13:00:00Z",
        "working_hours": 4.0,
        "reason": "Medical appointment",
        "notes": "Left early for doctor visit"
      },
      {
        "date": "2026-03-06",
        "status": "absent",
        "reason": "Sick leave"
      },
      {
        "date": "2026-03-07",
        "status": "holiday",
        "reason": "Weekend/Public Holiday"
      }
    ]
  }'
```

**Response:**
```json
{
  "total": 7,
  "success": 7,
  "failed": 0,
  "successful_dates": [
    "2026-03-01",
    "2026-03-02",
    "2026-03-03",
    "2026-03-04",
    "2026-03-05",
    "2026-03-06",
    "2026-03-07"
  ],
  "failed_records": []
}
```

---

## Step 5: View Individual Employee Calendar

Get a specific employee's attendance calendar view for a date range.

**Endpoint:** `GET /auth/manage/attendance/employee/{employee_id}/calendar`

**Query Parameters:**
- `from_date` - Start date (YYYY-MM-DD)
- `to_date` - End date (YYYY-MM-DD)

**Request:**
```bash
curl -X GET 'http://localhost:3000/auth/manage/attendance/employee/5/calendar?from_date=2026-03-01&to_date=2026-03-07' \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Response:**
```json
{
  "employee_id": 5,
  "employee_name": "John Doe",
  "company_id": 1,
  "date_range": [
    "2026-03-01",
    "2026-03-02",
    "2026-03-03",
    "2026-03-04",
    "2026-03-05",
    "2026-03-06",
    "2026-03-07"
  ],
  "daily_status": {
    "2026-03-01": {
      "id": 1,
      "status": "on_time",
      "check_in_time": "2026-03-01T09:00:00Z",
      "check_out_time": "2026-03-01T18:00:00Z",
      "working_hours": 9.0
    },
    "2026-03-02": {
      "id": 2,
      "status": "on_time",
      "check_in_time": "2026-03-02T09:15:00Z",
      "check_out_time": "2026-03-02T18:00:00Z",
      "working_hours": 8.75
    },
    "2026-03-03": {
      "id": 3,
      "status": "late",
      "check_in_time": "2026-03-03T09:30:00Z",
      "check_out_time": "2026-03-03T18:00:00Z",
      "working_hours": 8.5,
      "notes": "20 minutes late"
    },
    "2026-03-04": {
      "id": 4,
      "status": "on_time",
      "check_in_time": "2026-03-04T09:00:00Z",
      "check_out_time": "2026-03-04T18:00:00Z",
      "working_hours": 9.0
    },
    "2026-03-05": {
      "id": 5,
      "status": "half_day",
      "check_in_time": "2026-03-05T09:00:00Z",
      "check_out_time": "2026-03-05T13:00:00Z",
      "working_hours": 4.0,
      "notes": "Left early for doctor visit"
    },
    "2026-03-06": {
      "id": 6,
      "status": "absent",
      "reason": "Sick leave"
    },
    "2026-03-07": {
      "id": 7,
      "status": "holiday",
      "reason": "Weekend/Public Holiday"
    }
  },
  "statistics": {
    "total_days": 7,
    "days_present": 4,
    "days_absent": 1,
    "days_late": 1,
    "days_half_day": 1,
    "days_holiday": 1,
    "total_working_hours": 38.25
  }
}
```

---

## Step 6: View Company-Wide Weekly Attendance

Get all employees' attendance for a specific week. This endpoint returns all employees in the company with their attendance records for the selected date range, along with per-employee and company-wide statistics.

**Endpoint:** `GET /auth/manage/attendance/company/week-view`

**Query Parameters:**
- `from_date` - Week start date (YYYY-MM-DD)
- `to_date` - Week end date (YYYY-MM-DD)

**Request:**
```bash
curl -X GET 'http://localhost:3000/auth/manage/attendance/company/week-view?from_date=2026-03-01&to_date=2026-03-07' \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Response Fields Explained:**
- `success` - API call status
- `data.company_id` - Company identifier
- `data.start_date` / `data.end_date` - Week boundaries
- `data.date_range` - Array of all dates in the range for frontend calendar
- `data.employees[]` - Array of all company employees
  - `employee_id` - Unique employee identifier
  - `employee_name` - Full name
  - `email` - Employee email address
  - `employee_type` - Employment type (full-time, part-time, contract, etc.)
  - `daily_attendance` - Map of dates with attendance records
    - Each date contains: `date`, `status`, `reason`, `check_in_time`, `check_out_time`, `working_hours`, `notes`
  - `week_stats` - Count breakdown for the week (total, on_time, absent, late, holiday, half_day, leave)
- `data.company_stats` - Aggregated statistics across all employees for the week

**Response:**
```json
{
  "success": true,
  "message": "",
  "data": {
    "company_id": 1,
    "start_date": "2026-03-01",
    "end_date": "2026-03-07",
    "date_range": [
      "2026-03-01",
      "2026-03-02",
      "2026-03-03",
      "2026-03-04",
      "2026-03-05",
      "2026-03-06",
      "2026-03-07"
    ],
    "employees": [
      {
        "employee_id": 3,
        "employee_name": "John Smith",
        "email": "john.smith@company.com",
        "employee_type": "full-time",
        "daily_attendance": {
          "2026-03-02": {
            "date": "2026-03-02",
            "status": "on_time",
            "reason": "Regular working day",
            "check_in_time": "2026-03-02T09:00:00+05:30",
            "check_out_time": "2026-03-02T18:00:00+05:30",
            "working_hours": 9,
            "notes": "Good productive day"
          },
          "2026-03-03": {
            "date": "2026-03-03",
            "status": "late",
            "reason": "Traffic delay",
            "check_in_time": "2026-03-03T09:30:00+05:30",
            "check_out_time": "2026-03-03T18:00:00+05:30",
            "working_hours": 8.5,
            "notes": "20 minutes late"
          },
          "2026-03-05": {
            "date": "2026-03-05",
            "status": "half_day",
            "reason": "Medical appointment",
            "check_in_time": "2026-03-05T09:00:00+05:30",
            "check_out_time": "2026-03-05T13:00:00+05:30",
            "working_hours": 4,
            "notes": "Left early for doctor visit"
          }
        },
        "week_stats": {
          "total": 3,
          "on_time": 1,
          "absent": 0,
          "late": 1,
          "holiday": 0,
          "half_day": 1,
          "leave": 0
        }
      },
      {
        "employee_id": 4,
        "employee_name": "sri krishna",
        "email": "srik@company.com",
        "employee_type": "part-time",
        "daily_attendance": {},
        "week_stats": {
          "total": 0,
          "on_time": 0,
          "absent": 0,
          "late": 0,
          "holiday": 0,
          "half_day": 0,
          "leave": 0
        }
      },
      {
        "employee_id": 5,
        "employee_name": "Jane Doe",
        "email": "jane.doe@company.com",
        "employee_type": "full-time",
        "daily_attendance": {
          "2026-03-01": {
            "date": "2026-03-01",
            "status": "on_time",
            "reason": "Regular working day",
            "check_in_time": "2026-03-01T09:00:00+05:30",
            "check_out_time": "2026-03-01T18:00:00+05:30",
            "working_hours": 9,
            "notes": ""
          },
          "2026-03-04": {
            "date": "2026-03-04",
            "status": "absence",
            "reason": "Sick leave",
            "check_in_time": null,
            "check_out_time": null,
            "working_hours": 0,
            "notes": ""
          }
        },
        "week_stats": {
          "total": 2,
          "on_time": 1,
          "absent": 1,
          "late": 0,
          "holiday": 0,
          "half_day": 0,
          "leave": 0
        }
      }
    ],
    "company_stats": {
      "total_employees": 3,
      "total_present": 2,
      "total_absent": 1,
      "total_late": 1,
      "total_holiday": 0,
      "total_half_day": 1,
      "total_leave": 0
    }
  }
}
```

---

## Common Workflow Summary

### Complete Process:

1. **Create Company** → Get `company_id`

2. **Create Employee** → Get `employee_id`

3. **Mark Weekly Attendance** (Choose one method):
   - Single day: `POST /auth/manage/attendance`
   - Full week: `POST /auth/manage/attendance/bulk`

4. **View Results**:
   - Individual employee: `GET /auth/manage/attendance/employee/{employee_id}/calendar`
   - Company-wide: `GET /auth/manage/attendance/company/week-view`

## Frontend Implementation Notes

### Company Week View Data Structure

The company week-view endpoint provides data optimized for calendar rendering:

```javascript
// Use daily_attendance map for quick date lookups
const dateAttendance = response.data.employees[0].daily_attendance["2026-03-02"];
// Returns: { date, status, reason, check_in_time, check_out_time, working_hours, notes }

// week_stats provides instant counts for that employee's week
const employeeStats = response.data.employees[0].week_stats;
// Returns: { total: 3, on_time: 1, absent: 1, late: 1, holiday: 0, half_day: 0, leave: 0 }

// company_stats aggregates across all employees
const companyStats = response.data.company_stats;
// Returns: { total_employees: 3, total_present: 2, total_absent: 1, ... }
```

### Building a Calendar Grid

```javascript
// Use date_range to create calendar header
const dates = response.data.date_range; // ["2026-03-01", "2026-03-02", ...]

// For each employee, check daily_attendance
response.data.employees.forEach(employee => {
  dates.forEach(date => {
    const status = employee.daily_attendance[date]?.status;
    // Apply color based on status:
    // on_time -> green, absent -> red, late -> yellow, half_day -> orange, etc.
  });
});
```

---

Use these status values for calendar highlighting:

| Status | Color | Code | Meaning |
|--------|-------|------|---------|
| `on_time` | 🟢 Green | #22c55e | Present and on time |
| `late` | 🟡 Yellow | #eab308 | Present but late |
| `absent` | 🔴 Red | #ef4444 | Absent from work |
| `half_day` | 🟠 Orange | #f97316 | Present for portion of day |
| `leave` | 🔵 Blue | #3b82f6 | On approved leave |
| `holiday` | ⚫ Gray | #6b7280 | Non-working day/Holiday |

---

## Error Handling

### Common Errors and Solutions

**"Invalid request body"**
- Ensure dates are in `YYYY-MM-DD` format
- Check that all required fields are provided
- Verify JSON syntax

**"Employee not found"**
- Verify `employee_id` is valid
- Ensure employee belongs to the specified company

**"Duplicate attendance record"**
- Attendance for this date already exists
- Use `PUT /auth/manage/attendance/{id}` to update instead

**"Invalid date range"**
- Ensure `from_date` is before `to_date`
- Both dates must be in valid format

---

## Tips and Best Practices

1. **Bulk Operations**: Use bulk attendance creation for weekly/monthly imports to reduce API calls.

2. **Time Format**: Always use RFC 3339 format for check-in/out times (e.g., `2026-03-05T09:00:00Z`). Handle timezone offsets in responses (e.g., `+05:30`).

3. **Working Hours**: Calculate as decimal hours (e.g., 9.5 for 9 hours 30 minutes). For absent records, set to 0.

4. **Date Range Keys**: The `daily_attendance` object uses date strings as keys (e.g., `"2026-03-05"`). Check if key exists before accessing to handle employees with no records for specific dates.

5. **Weekly Statistics**: Use `week_stats` for quick employee-level summaries without server-side aggregation. Use `company_stats` for company-wide dashboards.

6. **Empty Attendance**: Employees with no attendance records show `daily_attendance: {}` and `week_stats.total: 0`. This is valid data, not an error.

7. **Employee Type Distinction**: Use `employee_type` field to filter or display different employee categories (full-time, part-time, contract, etc.).

8. **Batch Updates**: Update multiple records' statuses in one bulk request when possible for efficiency.

9. **Caching Strategy**: The `date_range` array is consistent across all employees. Cache it once per request to avoid repeated calculations.

10. **Error Handling**: Check `success` field in response. If false, refer to `message` field for error details.

---

## API Reference Summary

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/auth/manage/attendance` | Create single day attendance |
| `POST` | `/auth/manage/attendance/bulk` | Create multiple days attendance |
| `PUT` | `/auth/manage/attendance/{id}` | Update attendance record |
| `GET` | `/auth/manage/attendance/employee/{id}/calendar` | Get employee calendar view |
| `GET` | `/auth/manage/attendance/company/week-view` | Get company-wide week view |
| `DELETE` | `/auth/manage/attendance/{id}` | Delete attendance record |

---

## Testing with Current Date

Current date: **March 5, 2026** (Wednesday)

**Recommended test date range:** Mar 1-7, 2026 (Full week: Sunday to Saturday)

Use real employee IDs and company IDs from your database when testing. The example responses use these test IDs:
- Company ID: 1
- Employee IDs: 3, 4, 5

---

## Quick Start Example

1. Get company ID from your database:
```bash
SELECT id FROM companies LIMIT 1;
# Result: company_id = 1
```

2. Get employee IDs:
```bash
SELECT id, first_name, last_name FROM employees WHERE company_id = 1;
# Results: employee_id = 3,4,5
```

3. Mark attendance for this week (Mar 1-7):
```bash
curl -X POST http://localhost:3000/auth/manage/attendance/bulk \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "employee_id": 3,
    "company_id": 1,
    "attendances": [
      {"date": "2026-03-01", "status": "on_time", "check_in_time": "2026-03-01T09:00:00Z", "check_out_time": "2026-03-01T18:00:00Z", "working_hours": 9},
      {"date": "2026-03-02", "status": "late", "check_in_time": "2026-03-02T09:30:00Z", "check_out_time": "2026-03-02T18:00:00Z", "working_hours": 8.5},
      {"date": "2026-03-03", "status": "absent", "reason": "Sick leave"}
    ]
  }'
```

4. View company-wide week attendance:
```bash
curl -X GET 'http://localhost:3000/auth/manage/attendance/company/week-view?from_date=2026-03-01&to_date=2026-03-07' \
  -H "Authorization: Bearer YOUR_TOKEN"
```
