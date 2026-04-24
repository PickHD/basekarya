CREATE TABLE contracts (
    id              BIGINT AUTO_INCREMENT PRIMARY KEY,
    employee_id     BIGINT NOT NULL UNIQUE,   

    contract_type   ENUM('PKWT', 'PKWTT') NOT NULL,
    contract_number VARCHAR(50) NULL,

    start_date      DATE NOT NULL,
    end_date        DATE NULL,          

    notes           TEXT NULL,
    attachment_url  VARCHAR(255) NULL,  

    alerted_at      TIMESTAMP NULL,     

    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMP NULL,

    FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE,
    INDEX idx_contract_expiry (contract_type, end_date)
);
