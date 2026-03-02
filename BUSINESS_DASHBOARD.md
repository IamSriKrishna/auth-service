# Business Dashboard API Documentation

## Overview

The Dashboard API provides comprehensive business analytics and metrics for tracking key business entities:
- **Customers** - Customer counts and activity
- **Vendors** - Vendor counts and activity
- **Items** - Item inventory, stock levels, and item groups
- **Shipments** - Shipment tracking and delivery status
- **Invoices** - Invoice metrics and payment status
- **Sales Orders** - Sales order tracking
- **Purchase Orders** - Purchase order tracking
- **Packages** - Package tracking and delivery

## API Endpoints

### 1. Get Main Dashboard
```
GET /dashboard
```
Returns all dashboard metrics in one comprehensive response.

**Response:**
```json
{
  "customer_metrics": {
    "total": 450,
    "active": 380,
    "inactive": 70,
    "created_today": 5
  },
  "vendor_metrics": {
    "total": 120,
    "active": 95,
    "inactive": 25,
    "created_today": 1
  },
  "item_metrics": {
    "total": 2500,
    "total_stock": 50000,
    "low_stock_items": 45,
    "item_groups": 80,
    "created_today": 12,
    "out_of_stock_items": 8
  },
  "shipment_metrics": {
    "total": 1200,
    "shipped": 950,
    "pending": 150,
    "in_transit": 80,
    "delivered": 920,
    "cancelled_shipped": 0,
    "average_delivery_time_days": 3.5
  },
  "invoice_metrics": {
    "total": 800,
    "total_amount": 150000.50,
    "outstanding_amount": 25000.00,
    "paid_count": 700,
    "pending_count": 85,
    "overdue_count": 15
  },
  "sales_order_metrics": {
    "total": 600,
    "total_amount": 120000,
    "completed_count": 500,
    "pending_count": 80,
    "cancelled_count": 20,
    "created_today": 8
  },
  "purchase_order_metrics": {
    "total": 400,
    "total_amount": 95000,
    "completed_count": 350,
    "pending_count": 40,
    "cancelled_count": 10,
    "created_today": 3
  },
  "package_metrics": {
    "total": 1500,
    "shipped_count": 1200,
    "pending_count": 200,
    "in_transit_count": 80,
    "delivered_count": 1180,
    "created_today": 25
  },
  "last_updated_at": "2024-03-02T16:00:00Z",
  "generated_at": "2024-03-02T16:05:00Z"
}
```

### 2. Get Activity Summary
```
GET /dashboard/activity
```
Returns today's activity including new items created, shipments, orders, etc.

**Response:**
```json
{
  "created_customers_today": 5,
  "created_vendors_today": 1,
  "created_items_today": 12,
  "created_sales_orders_today": 8,
  "created_purchase_orders_today": 3,
  "shipped_today": 45,
  "delivered_today": 38
}
```

### 3. Get Stock Information
```
GET /dashboard/stock
```
Returns detailed stock information across all items.

**Response:**
```json
{
  "data": [
    {
      "item_id": "item_001",
      "item_name": "Widget A",
      "current_quantity": 150,
      "available_quantity": 120,
      "reserved_quantity": 30,
      "in_transit_quantity": 0,
      "status": "in_stock"
    },
    {
      "item_id": "item_002",
      "item_name": "Widget B",
      "current_quantity": 45,
      "available_quantity": 25,
      "reserved_quantity": 20,
      "in_transit_quantity": 0,
      "status": "low_stock"
    }
  ],
  "total_items": 2500,
  "in_stock_count": 2447,
  "low_stock_count": 45,
  "out_of_stock_count": 8,
  "total_quantity": 50000
}
```

### 4. Get Shipment Tracking
```
GET /dashboard/shipment/{shipment_id}/tracking?limit=10
```
Returns tracking details for a specific shipment.

**Parameters:**
- `shipment_id` (path): The ID of the shipment
- `limit` (optional): Number of tracking records (default: 10)

**Response:**
```json
{
  "data": [
    {
      "id": "track_123",
      "shipment_id": "ship_456",
      "status": "delivered",
      "location": "Customer Address",
      "latitude": 40.7128,
      "longitude": -74.0060,
      "notes": "Delivered successfully",
      "timestamp": "2024-03-02T15:30:00Z"
    },
    {
      "id": "track_122",
      "shipment_id": "ship_456",
      "status": "in_transit",
      "location": "Local Warehouse",
      "latitude": 40.7505,
      "longitude": -73.9972,
      "notes": "On the way to customer",
      "timestamp": "2024-03-02T10:15:00Z"
    }
  ],
  "total": 2
}
```

### 5. Add Shipment Tracking
```
POST /dashboard/shipment/{shipment_id}/tracking
```
Add a new tracking record for a shipment.

**Request Body:**
```json
{
  "status": "delivered",
  "location": "Customer Address",
  "latitude": 40.7128,
  "longitude": -74.0060,
  "notes": "Delivered successfully"
}
```

**Response:**
```json
{
  "message": "Tracking record added successfully"
}
```

### 6. Get Entity Trends
```
GET /dashboard/trends/{entity_type}?days=30
```
Returns historical trend data for a specific entity.

**Parameters:**
- `entity_type` (path): Entity type (customer, vendor, item, etc.)
- `days` (optional): Number of days to look back (default: 30)

**Response:**
```json
{
  "entity_type": "customer",
  "data": [
    {
      "date": "2024-02-01",
      "count": 400,
      "active_count": 320,
      "created_today": 5
    },
    {
      "date": "2024-02-02",
      "count": 410,
      "active_count": 330,
      "created_today": 10
    }
  ]
}
```

### 7. Refresh Dashboard Metrics
```
POST /dashboard/refresh
```
Manually trigger a refresh of all dashboard metrics.

**Response:**
```json
{
  "message": "Dashboard metrics refreshed successfully"
}
```

## Usage Examples

### Get Complete Dashboard Overview
```bash
curl -X GET "http://localhost:3000/dashboard" \
  -H "Content-Type: application/json"
```

### Get Today's Activity
```bash
curl -X GET "http://localhost:3000/dashboard/activity" \
  -H "Content-Type: application/json"
```

### Get Stock Status
```bash
curl -X GET "http://localhost:3000/dashboard/stock" \
  -H "Content-Type: application/json"
```

### Track Shipment
```bash
curl -X GET "http://localhost:3000/dashboard/shipment/ship_12345/tracking" \
  -H "Content-Type: application/json"
```

### Add Shipment Update
```bash
curl -X POST "http://localhost:3000/dashboard/shipment/ship_12345/tracking" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "in_transit",
    "location": "Regional Distribution Center",
    "latitude": 40.7505,
    "longitude": -73.9972,
    "notes": "Shipment reached distribution center"
  }'
```

### Get Customer Trends
```bash
curl -X GET "http://localhost:3000/dashboard/trends/customer?days=30" \
  -H "Content-Type: application/json"
```

### Get Item Trends
```bash
curl -X GET "http://localhost:3000/dashboard/trends/item?days=30" \
  -H "Content-Type: application/json"
```

## Metrics Explained

### Customer Metrics
- **Total**: Total number of customers
- **Active**: Customers currently active (have recent activity)
- **Inactive**: Customers not active
- **Created Today**: Customers created in the last 24 hours

### Vendor Metrics
- **Total**: Total number of vendors
- **Active**: Vendors currently active
- **Inactive**: Vendors not active
- **Created Today**: Vendors created in the last 24 hours

### Item Metrics
- **Total**: Total number of items in catalog
- **Total Stock**: Total quantity in inventory across all items (from inventory_balance)
- **Low Stock Items**: Items below reorder point (threshold: 100 units)
- **Item Groups**: Number of item groups/categories
- **Created Today**: Items created in the last 24 hours
- **Out of Stock Items**: Items with zero quantity

### Stock Endpoint Data
The `/dashboard/stock` endpoint provides detailed stock information including:
- **Current Quantity**: Total quantity available
- **Available Quantity**: Quantity available for sale (current - reserved - in transit)
- **Reserved Quantity**: Quantity reserved in sales orders
- **In Transit Quantity**: Quantity currently being transported
- **Status**: One of `in_stock`, `low_stock`, or `out_of_stock`

### Shipment Metrics
- **Total**: Total number of shipments
- **Shipped**: Shipments that have been shipped
- **Pending**: Shipments awaiting shipment
- **In Transit**: Shipments on the way to customer
- **Delivered**: Shipments successfully delivered
- **Average Delivery Time**: Average days from shipment to delivery

### Invoice Metrics
- **Total**: Total number of invoices
- **Total Amount**: Sum of all invoice amounts
- **Outstanding Amount**: Amount pending payment (unpaid invoices)
- **Paid Count**: Number of paid invoices
- **Pending Count**: Number of pending/unpaid invoices
- **Overdue Count**: Invoices past due date

### Sales Order Metrics
- **Total**: Total sales orders
- **Total Amount**: Sum of all sales order amounts
- **Completed Count**: Completed orders
- **Pending Count**: Orders still in progress
- **Cancelled Count**: Cancelled orders
- **Created Today**: Orders created in last 24 hours

### Purchase Order Metrics
- **Total**: Total purchase orders
- **Total Amount**: Sum of all purchase order amounts
- **Completed Count**: Completed orders
- **Pending Count**: Orders in progress
- **Cancelled Count**: Cancelled orders
- **Created Today**: Orders created in last 24 hours

### Package Metrics
- **Total**: Total packages
- **Shipped Count**: Packages shipped
- **Pending Count**: Packages awaiting shipment
- **In Transit Count**: Packages on the way
- **Delivered Count**: Packages successfully delivered
- **Created Today**: Packages created in last 24 hours

## Database Tables

### dashboard_metrics
Stores aggregated business metrics for quick retrieval. Updated periodically or on-demand.

### shipment_tracking
Stores historical tracking data for each shipment. Enables viewing full shipment journey.

### entity_count_history
Stores daily counts for trend analysis. Helps track growth patterns and historical data.

## Performance Optimization

1. **Metrics Caching**: Dashboard metrics are cached and only updated when refreshed or on interval
2. **Historical Data**: Trends are calculated from entity_count_history for faster queries
3. **Indexes**: Database tables are indexed on frequently queried columns
4. **Async Updates**: Metrics can be updated asynchronously during off-peak hours

## Integration Notes

The dashboard automatically queries the following tables for metrics:
- `customers` - for customer metrics
- `vendors` - for vendor metrics
- `items` - for item metrics
- `inventory_balance` - for stock information
- `shipments` - for shipment metrics
- `invoices` - for invoice metrics
- `sales_orders` - for sales order metrics
- `purchase_orders` - for purchase order metrics
- `packages` - for package metrics

No modifications are required to existing entity handlers. The dashboard queries read-only data and calculates metrics on-demand.

## Future Enhancements

Potential improvements:
- Real-time metrics with WebSocket updates
- Customizable date range selection
- Comparison reports (month-over-month, year-over-year)
- Export functionality (PDF, Excel)
- Alerts for threshold breaches
- Predictive analytics
- User role-based access control
- Dashboard customization per user
