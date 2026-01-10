CREATE TABLE ref_departments (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE ref_shifts (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE users (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(20) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role ENUM('SUPERADMIN', 'EMPLOYEE') DEFAULT 'EMPLOYEE',
    must_change_password BOOLEAN DEFAULT TRUE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE employees (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE,
    department_id BIGINT NOT NULL,
    shift_id BIGINT NOT NULL,
    
    nik VARCHAR(20) NOT NULL UNIQUE,
    full_name VARCHAR(100) NOT NULL,
    phone_number VARCHAR(20) NULL,
    profile_picture_url VARCHAR(255) NULL,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (department_id) REFERENCES ref_departments(id),
    FOREIGN KEY (shift_id) REFERENCES ref_shifts(id)
);

CREATE TABLE attendances (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    employee_id BIGINT NOT NULL,
    shift_id BIGINT NOT NULL,
    
    date DATE NOT NULL,
    
    check_in_time DATETIME NOT NULL,
    check_in_lat DECIMAL(10, 8) NOT NULL,
    check_in_long DECIMAL(11, 8) NOT NULL,
    check_in_image_url VARCHAR(255) NOT NULL,
    check_in_address TEXT NOT NULL,
    
    check_out_time DATETIME NULL,
    check_out_lat DECIMAL(10, 8) NULL,
    check_out_long DECIMAL(11, 8) NULL,
    check_out_image_url VARCHAR(255) NULL,
    check_out_address TEXT NULL,
    
    status ENUM('PRESENT', 'LATE', 'EXCUSED', 'ABSENT') DEFAULT 'ABSENT',
    is_suspicious BOOLEAN DEFAULT FALSE,
    notes TEXT,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (employee_id) REFERENCES employees(id),
    FOREIGN KEY (shift_id) REFERENCES ref_shifts(id),
    
    UNIQUE KEY unique_attendance_per_day (employee_id, date),
    INDEX idx_date (date),
    INDEX idx_status (status)
);

