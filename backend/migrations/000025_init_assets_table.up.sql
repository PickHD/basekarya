CREATE TABLE asset_categories (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  company_id BIGINT NOT NULL,
  name VARCHAR(255) NOT NULL,
  description TEXT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

  INDEX idx_asset_categories_company_id (company_id),
  CONSTRAINT fk_asset_categories_company
    FOREIGN KEY (company_id) REFERENCES companies(id)
    ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE assets (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  company_id BIGINT NOT NULL,
  asset_category_id BIGINT NOT NULL,
  name VARCHAR(255) NOT NULL,
  description TEXT NULL,
  serial_number VARCHAR(100) NULL,
  status ENUM('AVAILABLE','ASSIGNED','MAINTENANCE','RETIRED') NOT NULL DEFAULT 'AVAILABLE',
  `condition` ENUM('GOOD','FAIR','DAMAGED','LOST') NOT NULL DEFAULT 'GOOD',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

  INDEX idx_assets_company_id (company_id),
  INDEX idx_assets_category_id (asset_category_id),
  INDEX idx_assets_status (status),
  CONSTRAINT fk_assets_company
    FOREIGN KEY (company_id) REFERENCES companies(id)
    ON DELETE RESTRICT ON UPDATE CASCADE,
  CONSTRAINT fk_assets_asset_categories
    FOREIGN KEY (asset_category_id) REFERENCES asset_categories(id)
    ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE asset_assignments (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  company_id BIGINT NOT NULL,
  asset_id BIGINT NOT NULL,
  employee_id BIGINT NOT NULL,
  user_id BIGINT NOT NULL,
  approved_by BIGINT NULL,
  purpose TEXT NULL,
  expected_return_date DATE NULL,
  actual_return_date DATE NULL,
  notes TEXT NULL,
  status ENUM('PENDING','ACTIVE','RETURNED','REJECTED') NOT NULL DEFAULT 'PENDING',
  rejection_reason TEXT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

  INDEX idx_asset_assignments_company_id (company_id),
  INDEX idx_asset_assignments_asset_id (asset_id),
  INDEX idx_asset_assignments_employee_id (employee_id),
  INDEX idx_asset_assignments_status (status),

  CONSTRAINT fk_asset_assignments_company
    FOREIGN KEY (company_id) REFERENCES companies(id)
    ON DELETE RESTRICT ON UPDATE CASCADE,
  CONSTRAINT fk_asset_assignments_asset
    FOREIGN KEY (asset_id) REFERENCES assets(id)
    ON DELETE RESTRICT ON UPDATE CASCADE,
  CONSTRAINT fk_asset_assignments_employee
    FOREIGN KEY (employee_id) REFERENCES employees(id)
    ON DELETE RESTRICT ON UPDATE CASCADE,
  CONSTRAINT fk_asset_assignments_user
    FOREIGN KEY (user_id) REFERENCES users(id)
    ON DELETE RESTRICT ON UPDATE CASCADE,
  CONSTRAINT fk_asset_assignments_approver
    FOREIGN KEY (approved_by) REFERENCES users(id)
    ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
