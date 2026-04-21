-- Onboarding checklists (templates)
CREATE TABLE onboarding_templates (
    id          BIGINT AUTO_INCREMENT PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    department  VARCHAR(50) NOT NULL,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Template items
CREATE TABLE onboarding_template_items (
    id              BIGINT AUTO_INCREMENT PRIMARY KEY,
    template_id     BIGINT NOT NULL,
    task_name       VARCHAR(255) NOT NULL,
    description     TEXT NULL,
    sort_order      INT DEFAULT 0,

    FOREIGN KEY (template_id) REFERENCES onboarding_templates(id) ON DELETE CASCADE
);

-- Onboarding workflows (instances per new hire)
CREATE TABLE onboarding_workflows (
    id                  BIGINT AUTO_INCREMENT PRIMARY KEY,
    applicant_id        BIGINT NULL,
    employee_id         BIGINT NULL,
    new_hire_name       VARCHAR(100) NOT NULL,
    new_hire_email      VARCHAR(255) NOT NULL,
    position            VARCHAR(100) NULL,
    department          VARCHAR(100) NULL,
    start_date          DATE NULL,

    status              ENUM('IN_PROGRESS', 'COMPLETED') DEFAULT 'IN_PROGRESS',
    welcome_email_sent  BOOLEAN DEFAULT FALSE,

    created_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (applicant_id) REFERENCES applicants(id) ON DELETE SET NULL,
    FOREIGN KEY (employee_id)  REFERENCES employees(id)  ON DELETE SET NULL
);

-- Onboarding tasks (workflow items)
CREATE TABLE onboarding_tasks (
    id                      BIGINT AUTO_INCREMENT PRIMARY KEY,
    onboarding_workflow_id  BIGINT NOT NULL,
    template_item_id        BIGINT NULL,

    task_name               VARCHAR(255) NOT NULL,
    description             TEXT NULL,
    department              VARCHAR(50) NOT NULL,
    is_completed            BOOLEAN DEFAULT FALSE,
    completed_by            BIGINT  NULL,
    completed_at            TIMESTAMP NULL,
    notes                   TEXT NULL,
    sort_order              INT DEFAULT 0,

    FOREIGN KEY (onboarding_workflow_id) REFERENCES onboarding_workflows(id) ON DELETE CASCADE,
    FOREIGN KEY (template_item_id)       REFERENCES onboarding_template_items(id) ON DELETE SET NULL,
    FOREIGN KEY (completed_by)           REFERENCES users(id) ON DELETE SET NULL
);
