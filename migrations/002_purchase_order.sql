-- Migration for Purchase Order tables
-- This migration creates the purchase_orders and purchase_order_line_items tables

-- Create purchase_orders table
CREATE TABLE IF NOT EXISTS purchase_orders (
    id varchar(255) NOT NULL PRIMARY KEY,
    purchase_order_number varchar(100) NOT NULL UNIQUE,
    vendor_id INT NOT NULL,
    delivery_address_type varchar(50) NOT NULL,
    delivery_address_id INT,
    organization_name varchar(255),
    organization_address TEXT,
    customer_id INT,
    reference_no varchar(100),
    po_date DATETIME NOT NULL,
    delivery_date DATETIME NOT NULL,
    payment_terms varchar(50) NOT NULL,
    shipment_preference varchar(255),
    sub_total DECIMAL(18, 2) NOT NULL DEFAULT 0,
    discount DECIMAL(18, 2) DEFAULT 0,
    discount_type varchar(50),
    tax_type varchar(10),
    tax_id INT,
    tax_amount DECIMAL(18, 2) DEFAULT 0,
    adjustment DECIMAL(18, 2) DEFAULT 0,
    total DECIMAL(18, 2) NOT NULL DEFAULT 0,
    notes TEXT,
    terms_and_conditions TEXT,
    status varchar(50) NOT NULL DEFAULT 'draft',
    attachments JSON,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    created_by varchar(255),
    updated_by varchar(255),
    INDEX idx_vendor_id (vendor_id),
    INDEX idx_customer_id (customer_id),
    INDEX idx_status (status),
    INDEX idx_tax_id (tax_id),
    CONSTRAINT fk_po_vendor FOREIGN KEY (vendor_id) REFERENCES vendors(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    CONSTRAINT fk_po_customer FOREIGN KEY (customer_id) REFERENCES customers(id) ON UPDATE CASCADE ON DELETE SET NULL,
    CONSTRAINT fk_po_tax FOREIGN KEY (tax_id) REFERENCES taxes(id) ON UPDATE CASCADE ON DELETE SET NULL
);

-- Create purchase_order_line_items table
CREATE TABLE IF NOT EXISTS purchase_order_line_items (
    id INT AUTO_INCREMENT PRIMARY KEY,
    purchase_order_id varchar(255) NOT NULL,
    item_id varchar(255) NOT NULL,
    variant_sku VARCHAR(255),
    account varchar(100),
    quantity DECIMAL(18, 4) NOT NULL,
    rate DECIMAL(18, 2) NOT NULL,
    amount DECIMAL(18, 2) NOT NULL,
    received_quantity DECIMAL(18, 4) DEFAULT 0,
    variant_details JSON,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_po_id (purchase_order_id),
    INDEX idx_item_id (item_id),
    INDEX idx_variant_sku (variant_sku),
    CONSTRAINT fk_poli_po FOREIGN KEY (purchase_order_id) REFERENCES purchase_orders(id) ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT fk_poli_item FOREIGN KEY (item_id) REFERENCES items(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    CONSTRAINT fk_poli_variant FOREIGN KEY (variant_sku) REFERENCES variants(sku) ON UPDATE CASCADE ON DELETE SET NULL
);
