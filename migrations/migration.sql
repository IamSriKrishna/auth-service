-- Auth Service Database Migration
-- This script creates all required tables for the auth service
USE auth_service;

-- Create roles table
CREATE TABLE
    IF NOT EXISTS roles (
        id INT AUTO_INCREMENT PRIMARY KEY,
        role_name VARCHAR(255) NOT NULL UNIQUE,
        permissions JSON,
        description TEXT,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        deleted_at TIMESTAMP NULL,
        is_active BOOLEAN DEFAULT TRUE,
        INDEX idx_roles_deleted_at (deleted_at)
    );

-- Create users table
CREATE TABLE
    IF NOT EXISTS users (
        id INT AUTO_INCREMENT PRIMARY KEY,
        email VARCHAR(255) UNIQUE,
        phone VARCHAR(255) UNIQUE,
        username VARCHAR(255) UNIQUE,
        password_hash VARCHAR(255),
        google_id VARCHAR(255) UNIQUE,
        apple_id VARCHAR(255) UNIQUE,
        user_type ENUM ('mobile_user', 'superadmin', 'admin', 'partner') NOT NULL,
        role_id INT NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        deleted_at TIMESTAMP NULL,
        email_verified BOOLEAN DEFAULT FALSE,
        phone_verified BOOLEAN DEFAULT FALSE,
        status ENUM ('active', 'inactive', 'pending') DEFAULT 'active',
        created_by INT,
        last_login_at TIMESTAMP NULL,
        password_changed_at TIMESTAMP NULL,
        INDEX idx_users_email (email),
        INDEX idx_users_phone (phone),
        INDEX idx_users_username (username),
        INDEX idx_users_google_id (google_id),
        INDEX idx_users_apple_id (apple_id),
        INDEX idx_users_role_id (role_id),
        INDEX idx_users_created_by (created_by),
        INDEX idx_users_deleted_at (deleted_at),
        FOREIGN KEY (role_id) REFERENCES roles (id),
        FOREIGN KEY (created_by) REFERENCES users (id)
    );

-- Create refresh_tokens table
CREATE TABLE
    IF NOT EXISTS refresh_tokens (
        id INT AUTO_INCREMENT PRIMARY KEY,
        token_id VARCHAR(255) NOT NULL UNIQUE,
        user_id INT NOT NULL,
        expires_at TIMESTAMP NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        is_revoked BOOLEAN DEFAULT FALSE,
        INDEX idx_refresh_tokens_token_id (token_id),
        INDEX idx_refresh_tokens_user_id (user_id),
        FOREIGN KEY (user_id) REFERENCES users (id)
    );

-- Create user_sessions table
CREATE TABLE
    IF NOT EXISTS user_sessions (
        id INT AUTO_INCREMENT PRIMARY KEY,
        user_id INT NOT NULL,
        session_id VARCHAR(255) NOT NULL UNIQUE,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        expires_at TIMESTAMP NOT NULL,
        device_info TEXT,
        ip_address VARCHAR(45),
        user_agent TEXT,
        is_active BOOLEAN DEFAULT TRUE,
        INDEX idx_user_sessions_user_id (user_id),
        INDEX idx_user_sessions_session_id (session_id),
        FOREIGN KEY (user_id) REFERENCES users (id)
    );

-- Create support table
CREATE TABLE
    IF NOT EXISTS support (
        id INT AUTO_INCREMENT PRIMARY KEY,
        name VARCHAR(100) NOT NULL,
        phone VARCHAR(15),
        email VARCHAR(100),
        issue_type VARCHAR(100),
        description TEXT,
        status VARCHAR(50) DEFAULT 'open',
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        CONSTRAINT chk_contact CHECK (
            (
                phone IS NOT NULL
                AND phone <> ''
            )
            OR (
                email IS NOT NULL
                AND email <> ''
            )
        )
    );

CREATE TABLE
    bottle_sizes (
        id INT AUTO_INCREMENT PRIMARY KEY,
        size_ml INT NOT NULL UNIQUE,
        size_label VARCHAR(50) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

CREATE TABLE
    bottles (
        id INT AUTO_INCREMENT PRIMARY KEY,
        size_id INT NOT NULL,
        bottle_type VARCHAR(50) NOT NULL,
        neck_size_mm INT NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        CONSTRAINT fk_bottles_size FOREIGN KEY (size_id) REFERENCES bottle_sizes (id)
    );

CREATE TABLE
    caps (
        id INT AUTO_INCREMENT PRIMARY KEY,
        neck_size_mm INT NOT NULL,
        cap_type VARCHAR(50) NOT NULL,
        color VARCHAR(30),
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

CREATE TABLE
    products (
        id INT AUTO_INCREMENT PRIMARY KEY,
        product_name VARCHAR(255) NOT NULL,
        bottle_id INT NOT NULL,
        cap_id INT NOT NULL,
        mrp DECIMAL(10, 2),
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        CONSTRAINT fk_products_bottle FOREIGN KEY (bottle_id) REFERENCES bottles (id),
        CONSTRAINT fk_products_cap FOREIGN KEY (cap_id) REFERENCES caps (id)
    );

CREATE TABLE
    IF NOT EXISTS business_types (
        id INT AUTO_INCREMENT PRIMARY KEY,
        type_name VARCHAR(100) NOT NULL UNIQUE,
        description TEXT,
        is_active BOOLEAN DEFAULT TRUE,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        INDEX idx_business_types_active (is_active)
    );

CREATE TABLE
    IF NOT EXISTS companies (
        id INT AUTO_INCREMENT PRIMARY KEY,
        company_name VARCHAR(255) NOT NULL,
        business_type_id INT NOT NULL,
        gst_number VARCHAR(15) UNIQUE,
        pan_number VARCHAR(10) UNIQUE,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        deleted_at TIMESTAMP NULL,
        created_by INT,
        updated_by INT,
        INDEX idx_companies_business_type (business_type_id),
        INDEX idx_companies_gst (gst_number),
        INDEX idx_companies_pan (pan_number),
        INDEX idx_companies_deleted_at (deleted_at),
        CONSTRAINT fk_companies_business_type FOREIGN KEY (business_type_id) REFERENCES business_types (id),
        CONSTRAINT fk_companies_created_by FOREIGN KEY (created_by) REFERENCES users (id),
        CONSTRAINT fk_companies_updated_by FOREIGN KEY (updated_by) REFERENCES users (id)
    );

CREATE TABLE
    IF NOT EXISTS countries (
        id INT AUTO_INCREMENT PRIMARY KEY,
        country_name VARCHAR(100) NOT NULL UNIQUE,
        country_code VARCHAR(3) NOT NULL UNIQUE,
        phone_code VARCHAR(10),
        is_active BOOLEAN DEFAULT TRUE,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        INDEX idx_countries_code (country_code),
        INDEX idx_countries_active (is_active)
    );

CREATE TABLE
    IF NOT EXISTS states (
        id INT AUTO_INCREMENT PRIMARY KEY,
        country_id INT NOT NULL,
        state_name VARCHAR(100) NOT NULL,
        state_code VARCHAR(10),
        is_active BOOLEAN DEFAULT TRUE,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        INDEX idx_states_country (country_id),
        INDEX idx_states_active (is_active),
        CONSTRAINT fk_states_country FOREIGN KEY (country_id) REFERENCES countries (id),
        UNIQUE KEY unique_state_per_country (country_id, state_name)
    );

CREATE TABLE
    IF NOT EXISTS company_contacts (
        id INT AUTO_INCREMENT PRIMARY KEY,
        company_id INT NOT NULL UNIQUE,
        mobile VARCHAR(15) NOT NULL,
        alternate_mobile VARCHAR(15),
        email VARCHAR(255) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        INDEX idx_company_contacts_company (company_id),
        INDEX idx_company_contacts_email (email),
        INDEX idx_company_contacts_mobile (mobile),
        CONSTRAINT fk_company_contacts_company FOREIGN KEY (company_id) REFERENCES companies (id) ON DELETE CASCADE
    );

CREATE TABLE
    IF NOT EXISTS company_addresses (
        id INT AUTO_INCREMENT PRIMARY KEY,
        company_id INT NOT NULL UNIQUE,
        address_line1 VARCHAR(255) NOT NULL,
        address_line2 VARCHAR(255),
        city VARCHAR(100) NOT NULL,
        state_id INT NOT NULL,
        country_id INT NOT NULL,
        pincode VARCHAR(10) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        INDEX idx_company_addresses_company (company_id),
        INDEX idx_company_addresses_state (state_id),
        INDEX idx_company_addresses_country (country_id),
        INDEX idx_company_addresses_pincode (pincode),
        CONSTRAINT fk_company_addresses_company FOREIGN KEY (company_id) REFERENCES companies (id) ON DELETE CASCADE,
        CONSTRAINT fk_company_addresses_state FOREIGN KEY (state_id) REFERENCES states (id),
        CONSTRAINT fk_company_addresses_country FOREIGN KEY (country_id) REFERENCES countries (id)
    );

CREATE TABLE
    IF NOT EXISTS company_bank_details (
        id INT AUTO_INCREMENT PRIMARY KEY,
        company_id INT NOT NULL,
        bank_name VARCHAR(255) NOT NULL,
        account_holder_name VARCHAR(255) NOT NULL,
        account_number VARCHAR(50) NOT NULL,
        ifsc_code VARCHAR(11) NOT NULL,
        branch_name VARCHAR(255),
        is_primary BOOLEAN DEFAULT FALSE,
        is_active BOOLEAN DEFAULT TRUE,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        deleted_at TIMESTAMP NULL,
        INDEX idx_company_bank_company (company_id),
        INDEX idx_company_bank_ifsc (ifsc_code),
        INDEX idx_company_bank_primary (company_id, is_primary),
        INDEX idx_company_bank_deleted_at (deleted_at),
        CONSTRAINT fk_company_bank_company FOREIGN KEY (company_id) REFERENCES companies (id) ON DELETE CASCADE,
        UNIQUE KEY unique_primary_bank (company_id, is_primary)
    );

CREATE TABLE
    IF NOT EXISTS company_upi_details (
        id INT AUTO_INCREMENT PRIMARY KEY,
        company_id INT NOT NULL UNIQUE,
        upi_id VARCHAR(255) NOT NULL,
        upi_qr_url TEXT,
        is_active BOOLEAN DEFAULT TRUE,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        INDEX idx_company_upi_company (company_id),
        INDEX idx_company_upi_active (is_active),
        CONSTRAINT fk_company_upi_company FOREIGN KEY (company_id) REFERENCES companies (id) ON DELETE CASCADE
    );

CREATE TABLE
    IF NOT EXISTS tax_types (
        id INT AUTO_INCREMENT PRIMARY KEY,
        tax_name VARCHAR(50) NOT NULL UNIQUE,
        tax_code VARCHAR(10) NOT NULL UNIQUE,
        description TEXT,
        is_active BOOLEAN DEFAULT TRUE,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        INDEX idx_tax_types_active (is_active)
    );

CREATE TABLE
    IF NOT EXISTS company_invoice_settings (
        id INT AUTO_INCREMENT PRIMARY KEY,
        company_id INT NOT NULL UNIQUE,
        invoice_prefix VARCHAR(10) NOT NULL DEFAULT 'INV',
        invoice_start_number INT NOT NULL DEFAULT 1,
        current_invoice_number INT NOT NULL DEFAULT 1,
        show_logo BOOLEAN DEFAULT TRUE,
        show_signature BOOLEAN DEFAULT FALSE,
        round_off_total BOOLEAN DEFAULT TRUE,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        INDEX idx_company_invoice_company (company_id),
        CONSTRAINT fk_company_invoice_company FOREIGN KEY (company_id) REFERENCES companies (id) ON DELETE CASCADE,
        UNIQUE KEY unique_invoice_prefix (company_id, invoice_prefix)
    );

CREATE TABLE
    IF NOT EXISTS company_tax_settings (
        id INT AUTO_INCREMENT PRIMARY KEY,
        company_id INT NOT NULL UNIQUE,
        gst_enabled BOOLEAN DEFAULT TRUE,
        tax_type_id INT NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        INDEX idx_company_tax_company (company_id),
        INDEX idx_company_tax_type (tax_type_id),
        CONSTRAINT fk_company_tax_company FOREIGN KEY (company_id) REFERENCES companies (id) ON DELETE CASCADE,
        CONSTRAINT fk_company_tax_type FOREIGN KEY (tax_type_id) REFERENCES tax_types (id)
    );

CREATE TABLE
    IF NOT EXISTS company_regional_settings (
        id INT AUTO_INCREMENT PRIMARY KEY,
        company_id INT NOT NULL UNIQUE,
        timezone VARCHAR(50) DEFAULT 'Asia/Kolkata',
        date_format VARCHAR(20) DEFAULT 'DD/MM/YYYY',
        time_format VARCHAR(20) DEFAULT '24h',
        currency_code VARCHAR(3) DEFAULT 'INR',
        currency_symbol VARCHAR(10) DEFAULT '₹',
        language_code VARCHAR(5) DEFAULT 'en',
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        INDEX idx_company_regional_company (company_id),
        CONSTRAINT fk_company_regional_company FOREIGN KEY (company_id) REFERENCES companies (id) ON DELETE CASCADE
    );

ALTER TABLE bottle_sizes
ADD COLUMN updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;

ALTER TABLE bottles
ADD COLUMN thread_type VARCHAR(10) NOT NULL AFTER neck_size_mm,
ADD COLUMN updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;

CREATE INDEX idx_bottles_size_id ON bottles (size_id);

CREATE INDEX idx_bottles_neck_thread ON bottles (neck_size_mm, thread_type);

ALTER TABLE caps
ADD COLUMN thread_type VARCHAR(10) NOT NULL AFTER neck_size_mm,
ADD COLUMN material VARCHAR(30) NOT NULL AFTER color,
ADD COLUMN updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;

CREATE INDEX idx_caps_neck_thread ON caps (neck_size_mm, thread_type);

ALTER TABLE products MODIFY COLUMN mrp DECIMAL(10, 2) NOT NULL,
ADD COLUMN updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;

CREATE INDEX idx_products_bottle_id ON products (bottle_id);

CREATE INDEX idx_products_cap_id ON products (cap_id);

-- Add stock column to bottles table
ALTER TABLE bottles
ADD COLUMN stock INT NOT NULL DEFAULT 0 AFTER thread_type;

-- Add stock column to caps table
ALTER TABLE caps
ADD COLUMN stock INT NOT NULL DEFAULT 0 AFTER material;

-- Add index for stock queries
CREATE INDEX idx_bottles_stock ON bottles (stock);

CREATE INDEX idx_caps_stock ON caps (stock);

-- Add updated_at column to support table
ALTER TABLE support
ADD COLUMN updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;

-- Add deleted_at column to support table
ALTER TABLE support
ADD COLUMN deleted_at TIMESTAMP NULL;

-- Insert default roles
INSERT IGNORE INTO roles (role_name, permissions, description, is_active)
VALUES
    (
        'mobile_user',
        '["read:profile", "update:profile"]',
        'Default role for mobile app users',
        TRUE
    ),
    (
        'admin',
        '["read:users", "create:users", "update:users", "delete:users", "read:roles"]',
        'Administrator role with user management permissions',
        TRUE
    ),
    (
        'superadmin',
        '["*"]',
        'Super administrator with all permissions',
        TRUE
    ),
    (
        'partner',
        '["read:partner_data", "update:partner_profile"]',
        'Partner role with limited access',
        TRUE
    );

INSERT IGNORE INTO business_types (type_name, description)
VALUES
    (
        'Sole Proprietorship',
        'Business owned and run by one person'
    ),
    (
        'Partnership',
        'Business owned by two or more people'
    ),
    (
        'Limited Liability Partnership (LLP)',
        'Hybrid of partnership and company'
    ),
    (
        'Private Limited Company',
        'Company limited by shares'
    ),
    (
        'Public Limited Company',
        'Company whose shares are traded publicly'
    ),
    (
        'One Person Company (OPC)',
        'Company with single member'
    );

INSERT IGNORE INTO tax_types (tax_name, tax_code, description)
VALUES
    ('GST', 'GST', 'Goods and Services Tax'),
    (
        'IGST',
        'IGST',
        'Integrated Goods and Services Tax'
    ),
    ('VAT', 'VAT', 'Value Added Tax'),
    ('None', 'NONE', 'No Tax Applied');

INSERT IGNORE INTO countries (country_name, country_code, phone_code)
VALUES
    ('India', 'IN', '+91'),
    ('United States', 'US', '+1'),
    ('United Kingdom', 'GB', '+44');

INSERT IGNORE INTO states (country_id, state_name, state_code)
VALUES
    (1, 'Tamil Nadu', 'TN'),
    (1, 'Karnataka', 'KA'),
    (1, 'Maharashtra', 'MH'),
    (1, 'Delhi', 'DL'),
    (1, 'Gujarat', 'GJ'),
    (1, 'Rajasthan', 'RJ'),
    (1, 'Uttar Pradesh', 'UP'),
    (1, 'West Bengal', 'WB'),
    (1, 'Kerala', 'KL'),
    (1, 'Andhra Pradesh', 'AP');

CREATE
OR REPLACE VIEW company_complete_profile AS
SELECT
    c.id,
    c.company_name,
    bt.type_name AS business_type,
    c.gst_number,
    c.pan_number,
    cc.mobile,
    cc.alternate_mobile,
    cc.email,
    ca.address_line1,
    ca.address_line2,
    ca.city,
    s.state_name,
    co.country_name,
    ca.pincode,
    cbd.bank_name,
    cbd.account_number,
    cbd.ifsc_code,
    cud.upi_id,
    cis.invoice_prefix,
    cis.current_invoice_number,
    cis.show_logo,
    cis.show_signature,
    cis.round_off_total,
    cts.gst_enabled,
    tt.tax_name AS tax_type,
    crs.timezone,
    crs.currency_code,
    crs.language_code,
    c.created_at,
    c.updated_at
FROM
    companies c
    LEFT JOIN business_types bt ON c.business_type_id = bt.id
    LEFT JOIN company_contacts cc ON c.id = cc.company_id
    LEFT JOIN company_addresses ca ON c.id = ca.company_id
    LEFT JOIN states s ON ca.state_id = s.id
    LEFT JOIN countries co ON ca.country_id = co.id
    LEFT JOIN company_bank_details cbd ON c.id = cbd.company_id
    AND cbd.is_primary = TRUE
    LEFT JOIN company_upi_details cud ON c.id = cud.company_id
    LEFT JOIN company_invoice_settings cis ON c.id = cis.company_id
    LEFT JOIN company_tax_settings cts ON c.id = cts.company_id
    LEFT JOIN tax_types tt ON cts.tax_type_id = tt.id
    LEFT JOIN company_regional_settings crs ON c.id = crs.company_id;

-- Show created tables
SHOW TABLES;

-- Show table structures
DESCRIBE roles;

DESCRIBE users;

DESCRIBE refresh_tokens;

DESCRIBE user_sessions;