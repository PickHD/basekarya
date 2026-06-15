CREATE TABLE onboarding_templates (
    id          BIGINT AUTO_INCREMENT PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    department  VARCHAR(50) NOT NULL,
    company_id  BIGINT NOT NULL DEFAULT 1,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at  TIMESTAMP NULL,
    INDEX idx_onboarding_templates_company_id (company_id),
    INDEX idx_onboarding_templates_deleted_at (deleted_at)
);

CREATE TABLE onboarding_template_items (
    id              BIGINT AUTO_INCREMENT PRIMARY KEY,
    template_id     BIGINT NOT NULL,
    task_name       VARCHAR(255) NOT NULL,
    description     TEXT NULL,
    sort_order      INT DEFAULT 0,
    company_id      BIGINT NOT NULL DEFAULT 1,
    INDEX idx_onboarding_template_items_company_id (company_id),
    FOREIGN KEY (template_id) REFERENCES onboarding_templates(id) ON DELETE CASCADE
);

ALTER TABLE onboarding_tasks ADD COLUMN template_item_id BIGINT NULL;
ALTER TABLE onboarding_tasks ADD COLUMN department VARCHAR(50) NOT NULL DEFAULT '';
ALTER TABLE onboarding_tasks ADD FOREIGN KEY (template_item_id) REFERENCES onboarding_template_items(id) ON DELETE SET NULL;
