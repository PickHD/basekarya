CREATE TABLE IF NOT EXISTS overtimes (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  created_at DATETIME NULL,
  updated_at DATETIME NULL,
  
  user_id BIGINT NOT NULL,
  employee_id BIGINT NOT NULL,
  approved_by BIGINT NULL,
  
  date DATE NOT NULL,
  start_time TIME NOT NULL,
  end_time TIME NOT NULL,
  duration_minutes INT NOT NULL,
  reason TEXT NULL,
  
  status ENUM('PENDING', 'APPROVED', 'REJECTED', 'PAID') NOT NULL DEFAULT 'PENDING',
  rejection_reason TEXT NULL,
  
  INDEX idx_overtimes_user_id (user_id),
  INDEX idx_overtimes_status (status),

  CONSTRAINT fk_overtimes_user
      FOREIGN KEY (user_id) REFERENCES users(id)
      ON DELETE RESTRICT ON UPDATE CASCADE,

  CONSTRAINT fk_overtimes_employee
      FOREIGN KEY (employee_id) REFERENCES employees(id)
      ON DELETE RESTRICT ON UPDATE CASCADE,
      
  CONSTRAINT fk_overtimes_approver
      FOREIGN KEY (approved_by) REFERENCES users(id)
      ON DELETE SET NULL ON UPDATE CASCADE
);
