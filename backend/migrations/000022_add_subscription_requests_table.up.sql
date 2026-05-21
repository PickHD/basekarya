CREATE TABLE IF NOT EXISTS subscription_requests (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    company_id BIGINT UNSIGNED NOT NULL,
    current_plan_id BIGINT UNSIGNED NOT NULL,
    requested_plan_id BIGINT UNSIGNED NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    requested_by BIGINT UNSIGNED,
    reviewed_by BIGINT UNSIGNED,
    reviewed_at DATETIME,
    notes TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE INDEX idx_subscription_requests_company_id ON subscription_requests(company_id);
CREATE INDEX idx_subscription_requests_status ON subscription_requests(status);
