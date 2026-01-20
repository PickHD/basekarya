CREATE TABLE IF NOT EXISTS reimbursements (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  created_at DATETIME NULL,
  updated_at DATETIME NULL,
  deleted_at DATETIME NULL,
  
  user_id BIGINT NOT NULL,
  approved_by BIGINT NULL COMMENT 'User ID of the superadmin who approved/rejected',
  
  title VARCHAR(255) NOT NULL,
  description TEXT NULL,
  amount DECIMAL(15, 2) NOT NULL COMMENT 'Using DECIMAL for precise financial calculation',
  date_of_expense DATE NOT NULL,
  proof_file_url VARCHAR(255) NOT NULL,
  
  status ENUM('PENDING', 'APPROVED', 'REJECTED') NOT NULL DEFAULT 'PENDING',
  rejection_reason TEXT NULL,
  
  INDEX idx_reimbursements_user_id (user_id),
  INDEX idx_reimbursements_status (status),
  INDEX idx_reimbursements_deleted_at (deleted_at),

  CONSTRAINT fk_reimbursements_user
      FOREIGN KEY (user_id) REFERENCES users(id)
      ON DELETE RESTRICT ON UPDATE CASCADE,
      
  CONSTRAINT fk_reimbursements_approver
      FOREIGN KEY (approved_by) REFERENCES users(id)
      ON DELETE SET NULL ON UPDATE CASCADE
);