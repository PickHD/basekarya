CREATE TABLE subscription_plans (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(50) NOT NULL UNIQUE,
    max_employees INT NOT NULL DEFAULT 0,
    price_monthly DECIMAL(10,2) NOT NULL DEFAULT 0,
    features JSON,
    is_active TINYINT(1) NOT NULL DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

ALTER TABLE companies ADD COLUMN subscription_plan_id BIGINT DEFAULT NULL;
ALTER TABLE companies ADD COLUMN subscription_status VARCHAR(20) DEFAULT 'ACTIVE';
ALTER TABLE companies ADD COLUMN subscription_expires_at TIMESTAMP NULL;
ALTER TABLE companies ADD COLUMN owner_user_id BIGINT DEFAULT NULL;
ALTER TABLE companies ADD INDEX idx_companies_subscription_plan_id (subscription_plan_id);
