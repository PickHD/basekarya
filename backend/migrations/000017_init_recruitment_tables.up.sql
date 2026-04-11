-- Job Requisitions
CREATE TABLE job_requisitions (
    id              BIGINT AUTO_INCREMENT PRIMARY KEY,
    requester_id    BIGINT NOT NULL,
    department_id   BIGINT NOT NULL,

    title           VARCHAR(255) NOT NULL,
    description     TEXT NULL,
    quantity        INT DEFAULT 1,
    employment_type VARCHAR(10) NOT NULL,   -- 'PKWT' or 'PKWTT'

    priority        VARCHAR(10) DEFAULT 'MEDIUM',  -- 'LOW','MEDIUM','HIGH','URGENT'
    status          VARCHAR(10) DEFAULT 'DRAFT',   -- 'DRAFT','PENDING','APPROVED','REJECTED','CLOSED'

    approved_by      BIGINT NULL,
    rejection_reason TEXT NULL,

    target_date     DATE NULL,

    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMP NULL,

    FOREIGN KEY (requester_id) REFERENCES users(id),
    FOREIGN KEY (department_id) REFERENCES ref_departments(id),
    FOREIGN KEY (approved_by) REFERENCES users(id),
    INDEX idx_requisition_status (status)
);

-- Applicants
CREATE TABLE applicants (
    id                  BIGINT AUTO_INCREMENT PRIMARY KEY,
    job_requisition_id  BIGINT NOT NULL,

    full_name           VARCHAR(100) NOT NULL,
    email               VARCHAR(255) NOT NULL,
    phone_number        VARCHAR(20) NULL,
    resume_url          VARCHAR(255) NULL,

    stage               VARCHAR(15) DEFAULT 'SCREENING',  -- 'SCREENING','INTERVIEW','OFFERING','HIRED','REJECTED'
    stage_order         INT DEFAULT 0,

    notes               TEXT NULL,
    rejection_reason    TEXT NULL,

    created_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at          TIMESTAMP NULL,

    FOREIGN KEY (job_requisition_id) REFERENCES job_requisitions(id) ON DELETE CASCADE,
    INDEX idx_applicant_stage (job_requisition_id, stage)
);

-- Applicant stage history (audit trail)
CREATE TABLE applicant_stage_histories (
    id              BIGINT AUTO_INCREMENT PRIMARY KEY,
    applicant_id    BIGINT NOT NULL,
    from_stage      VARCHAR(15) NULL,
    to_stage        VARCHAR(15) NOT NULL,
    changed_by      BIGINT NOT NULL,
    notes           TEXT NULL,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (applicant_id) REFERENCES applicants(id) ON DELETE CASCADE,
    FOREIGN KEY (changed_by) REFERENCES users(id)
);
