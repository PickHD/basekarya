-- Add marital status and dependents to employees
ALTER TABLE employees
  ADD COLUMN marital_status ENUM('TK','K') NULL AFTER npwp,
  ADD COLUMN dependents_count TINYINT UNSIGNED DEFAULT 0 AFTER marital_status;

-- Add BPJS registration numbers to companies
ALTER TABLE companies
  ADD COLUMN bpjs_kesehatan_number VARCHAR(50) NULL AFTER tax_number,
  ADD COLUMN bpjs_ketenagakerjaan_number VARCHAR(50) NULL AFTER bpjs_kesehatan_number;

-- Extend payroll_details with grouping and employer-borne tracking
ALTER TABLE payroll_details
  ADD COLUMN code VARCHAR(30) NULL AFTER title,
  ADD COLUMN `group` VARCHAR(30) NULL AFTER code,
  ADD COLUMN is_employer_borne TINYINT(1) DEFAULT 0 AFTER `group`,
  ADD INDEX idx_payroll_details_group (`group`);

-- BPJS rate configuration
CREATE TABLE bpjs_rate_configs (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL,
  company_id BIGINT NULL,
  type ENUM('KESEHATAN','JHT','JKK','JKM','JP') NOT NULL,
  employee_rate DECIMAL(5,4) NOT NULL DEFAULT 0,
  employer_rate DECIMAL(5,4) NOT NULL DEFAULT 0,
  max_salary_cap DECIMAL(15,2) NULL,
  industry_risk_level ENUM('VERY_LOW','LOW','MEDIUM','HIGH','VERY_HIGH') NULL,
  is_active TINYINT(1) DEFAULT 1,
  effective_from DATE NOT NULL,
  effective_until DATE NULL,
  INDEX idx_bpjs_config_company (company_id),
  INDEX idx_bpjs_config_type (type),
  INDEX idx_bpjs_config_active (is_active),
  INDEX idx_bpjs_config_deleted (deleted_at),
  CONSTRAINT fk_bpjs_config_company FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- PPh 21 TER rate brackets
CREATE TABLE pph21_term_configs (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL,
  company_id BIGINT NULL,
  category CHAR(1) NOT NULL,
  bracket_number INT NOT NULL,
  min_monthly_salary DECIMAL(15,2) NOT NULL,
  rate DECIMAL(5,4) NOT NULL,
  effective_from DATE NOT NULL,
  effective_until DATE NULL,
  INDEX idx_term_category (category),
  INDEX idx_term_effective (effective_from),
  CONSTRAINT fk_term_config_company FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- PTKP annual thresholds
CREATE TABLE ptkp_configs (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL,
  code VARCHAR(10) NOT NULL,
  annual_amount DECIMAL(15,2) NOT NULL,
  effective_year INT NOT NULL,
  INDEX idx_ptkp_code (code),
  INDEX idx_ptkp_year (effective_year)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
