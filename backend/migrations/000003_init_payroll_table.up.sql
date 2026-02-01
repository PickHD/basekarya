ALTER TABLE employees
  ADD COLUMN base_salary DECIMAL(15, 2) DEFAULT 0,
  ADD COLUMN bank_name VARCHAR(50) NULL,
  ADD COLUMN bank_account_number VARCHAR(50) NULL,
  ADD COLUMN bank_account_holder VARCHAR(100) NULL,
  ADD COLUMN npwp VARCHAR(30) NULL;

CREATE TABLE payrolls (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL,
  
  employee_id BIGINT NOT NULL,
  period_date DATE NOT NULL,
  
  base_salary DECIMAL(15, 2) DEFAULT 0,
  total_allowance DECIMAL(15, 2) DEFAULT 0,
  total_deduction DECIMAL(15, 2) DEFAULT 0,
  net_salary DECIMAL(15, 2) DEFAULT 0,
  
  status VARCHAR(20) DEFAULT 'DRAFT',
  notes TEXT NULL,
  
  INDEX idx_payrolls_period_date (period_date),
  INDEX idx_payrolls_employee_id (employee_id),
  INDEX idx_payrolls_deleted_at (deleted_at),

  CONSTRAINT fk_payrolls_employee
    FOREIGN KEY (employee_id) 
    REFERENCES employees(id) 
    ON DELETE CASCADE 
    ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE payroll_details (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  payroll_id BIGINT NOT NULL,
  
  title VARCHAR(150) NOT NULL,
  type VARCHAR(20) NOT NULL,
  amount DECIMAL(15, 2) NOT NULL DEFAULT 0,
  
  INDEX idx_payroll_details_payroll_id (payroll_id),

  CONSTRAINT fk_payroll_details_payroll
    FOREIGN KEY (payroll_id) 
    REFERENCES payrolls(id) 
    ON DELETE CASCADE 
    ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;