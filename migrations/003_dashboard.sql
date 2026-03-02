-- Dashboard Metrics table
CREATE TABLE IF NOT EXISTS dashboard_metrics (
    id VARCHAR(255) PRIMARY KEY,
    total_customers INT DEFAULT 0,
    active_customers INT DEFAULT 0,
    total_vendors INT DEFAULT 0,
    active_vendors INT DEFAULT 0,
    total_items INT DEFAULT 0,
    total_item_groups INT DEFAULT 0,
    total_stock BIGINT DEFAULT 0,
    low_stock_items INT DEFAULT 0,
    total_shipments INT DEFAULT 0,
    shipped_count INT DEFAULT 0,
    pending_shipments INT DEFAULT 0,
    total_invoices INT DEFAULT 0,
    invoice_amount FLOAT DEFAULT 0,
    total_sales_orders INT DEFAULT 0,
    sales_order_amount FLOAT DEFAULT 0,
    total_purchase_orders INT DEFAULT 0,
    purchase_order_amount FLOAT DEFAULT 0,
    total_packages INT DEFAULT 0,
    packages_shipped INT DEFAULT 0,
    pending_packages INT DEFAULT 0,
    last_updated_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Shipment Tracking table
CREATE TABLE IF NOT EXISTS shipment_tracking (
    id VARCHAR(255) PRIMARY KEY,
    shipment_id VARCHAR(255) NOT NULL,
    status VARCHAR(100) NOT NULL,
    location VARCHAR(255),
    latitude FLOAT,
    longitude FLOAT,
    updated_by VARCHAR(255),
    notes TEXT,
    timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_shipment_id (shipment_id),
    INDEX idx_status (status),
    INDEX idx_timestamp (timestamp)
);

-- Entity Count History table for trends
CREATE TABLE IF NOT EXISTS entity_count_history (
    id VARCHAR(255) PRIMARY KEY,
    entity_type VARCHAR(100) NOT NULL,
    date DATE NOT NULL,
    count INT DEFAULT 0,
    active_count INT DEFAULT 0,
    created_today INT DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_entity_type (entity_type),
    INDEX idx_date (date)
);
