CREATE TABLE ref_leave_types (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(50) NOT NULL UNIQUE,
  default_quota INT DEFAULT 12,
  is_deducted BOOLEAN DEFAULT TRUE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE leave_balances (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  employee_id BIGINT NOT NULL,
  leave_type_id BIGINT NOT NULL,
  year INT NOT NULL,
  quota_total INT DEFAULT 0,
  quota_used INT DEFAULT 0,
  quota_left INT DEFAULT 0,
  
  FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE,
  FOREIGN KEY (leave_type_id) REFERENCES ref_leave_types(id),
  UNIQUE KEY idx_balance_user_year (employee_id, leave_type_id, year)
);

CREATE TABLE leave_requests (
  id BIGINT  AUTO_INCREMENT PRIMARY KEY,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL,
  
  user_id BIGINT  NOT NULL,
  employee_id BIGINT NOT NULL,
  leave_type_id BIGINT NOT NULL,
  
  start_date DATE NOT NULL,
  end_date DATE NOT NULL,
  total_days INT NOT NULL,
  
  reason TEXT,
  attachment_url VARCHAR(255),
  
  status VARCHAR(20) DEFAULT 'PENDING',
  approved_by BIGINT NULL,
  rejection_reason VARCHAR(255),
  
  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (employee_id) REFERENCES employees(id),
  FOREIGN KEY (leave_type_id) REFERENCES ref_leave_types(id)
);