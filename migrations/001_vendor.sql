CREATE TABLE IF NOT EXISTS vendors (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,

    salutation VARCHAR(10),
    first_name VARCHAR(100),
    last_name VARCHAR(100),

    company_name VARCHAR(255),
    display_name VARCHAR(255) NOT NULL,

    email VARCHAR(255),
    work_phone VARCHAR(20),
    mobile VARCHAR(20),

    vendor_language VARCHAR(50) DEFAULT 'English',
    remarks TEXT,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS vendor_compliance (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    vendor_id BIGINT NOT NULL,

    gstin VARCHAR(15),
    pan VARCHAR(10),
    is_msme_registered BOOLEAN DEFAULT FALSE,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_vendor_compliance_vendor
        FOREIGN KEY (vendor_id) REFERENCES vendors(id)
        ON DELETE CASCADE,

    UNIQUE KEY uk_vendor_compliance_vendor (vendor_id)
);

CREATE TABLE IF NOT EXISTS vendor_financials (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    vendor_id BIGINT NOT NULL,

    currency VARCHAR(10) DEFAULT 'INR',
    payment_terms VARCHAR(50) DEFAULT 'Due on Receipt',
    tds_tax VARCHAR(100),

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_vendor_financials_vendor
        FOREIGN KEY (vendor_id) REFERENCES vendors(id)
        ON DELETE CASCADE,

    UNIQUE KEY uk_vendor_financials_vendor (vendor_id)
);

CREATE TABLE IF NOT EXISTS vendor_settings (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    vendor_id BIGINT NOT NULL,

    portal_enabled BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_vendor_settings_vendor
        FOREIGN KEY (vendor_id) REFERENCES vendors(id)
        ON DELETE CASCADE,

    UNIQUE KEY uk_vendor_settings_vendor (vendor_id)
);


