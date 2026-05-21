CREATE TABLE IF NOT EXISTS finance_categories (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  type ENUM('INCOME', 'EXPENSE') NOT NULL,
  description TEXT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

  INDEX idx_finance_categories_type (type)
);

INSERT INTO finance_categories (name, type, description) VALUES
  ('Gaji & Tunjangan', 'EXPENSE', 'Pembayaran gaji dan tunjangan karyawan'),
  ('Sewa & Utilitas', 'EXPENSE', 'Biaya sewa gedung, listrik, air, dan internet'),
  ('Perlengkapan Kantor', 'EXPENSE', 'Pembelian perlengkapan dan ATK'),
  ('Transportasi', 'EXPENSE', 'Biaya transportasi dan perjalanan dinas'),
  ('Pemasaran', 'EXPENSE', 'Biaya iklan dan promosi'),
  ('Pemeliharaan', 'EXPENSE', 'Biaya perbaikan dan pemeliharaan aset'),
  ('Lain-lain Pengeluaran', 'EXPENSE', 'Pengeluaran lain-lain'),
  ('Penjualan Produk', 'INCOME', 'Pendapatan dari penjualan produk'),
  ('Jasa Layanan', 'INCOME', 'Pendapatan dari jasa layanan'),
  ('Investasi', 'INCOME', 'Pendapatan dari investasi'),
  ('Lain-lain Pemasukan', 'INCOME', 'Pemasukan lain-lain');

CREATE TABLE IF NOT EXISTS finance_transactions (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

  finance_category_id BIGINT NOT NULL,
  created_by BIGINT NOT NULL,
  approved_by BIGINT NULL,

  type ENUM('INCOME', 'EXPENSE') NOT NULL,
  amount DECIMAL(15,2) NOT NULL,
  description TEXT NULL,
  transaction_date DATE NOT NULL,
  reference_number VARCHAR(100) NULL,

  status ENUM('PENDING', 'APPROVED', 'REJECTED') NOT NULL DEFAULT 'PENDING',
  rejection_reason TEXT NULL,

  INDEX idx_finance_transactions_type (type),
  INDEX idx_finance_transactions_status (status),
  INDEX idx_finance_transactions_date (transaction_date),
  INDEX idx_finance_transactions_category (finance_category_id),
  INDEX idx_finance_transactions_created_by (created_by),

  CONSTRAINT fk_finance_transactions_category
    FOREIGN KEY (finance_category_id) REFERENCES finance_categories(id)
    ON DELETE RESTRICT ON UPDATE CASCADE,
  CONSTRAINT fk_finance_transactions_created_by
    FOREIGN KEY (created_by) REFERENCES users(id)
    ON DELETE RESTRICT ON UPDATE CASCADE,
  CONSTRAINT fk_finance_transactions_approver
    FOREIGN KEY (approved_by) REFERENCES users(id)
    ON DELETE SET NULL ON UPDATE CASCADE
);
