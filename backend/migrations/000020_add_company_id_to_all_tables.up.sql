-- 1. users (also add is_platform_admin)
ALTER TABLE users ADD COLUMN is_platform_admin TINYINT(1) NOT NULL DEFAULT 0;
ALTER TABLE users ADD COLUMN company_id BIGINT NOT NULL DEFAULT 1;
ALTER TABLE users ADD INDEX idx_users_company_id (company_id);
ALTER TABLE users DROP INDEX username;
ALTER TABLE users ADD UNIQUE INDEX idx_users_username_company_id (username, company_id);

-- 2. employees
ALTER TABLE employees ADD COLUMN company_id BIGINT NOT NULL DEFAULT 1;
ALTER TABLE employees ADD INDEX idx_employees_company_id (company_id);

-- 3. roles
ALTER TABLE roles ADD COLUMN company_id BIGINT NOT NULL DEFAULT 1;
ALTER TABLE roles ADD INDEX idx_roles_company_id (company_id);
ALTER TABLE roles DROP INDEX name;
ALTER TABLE roles ADD UNIQUE INDEX idx_roles_name_company_id (name, company_id);

-- 4. role_permissions
ALTER TABLE role_permissions ADD COLUMN company_id BIGINT NOT NULL DEFAULT 1;
ALTER TABLE role_permissions ADD INDEX idx_role_permissions_company_id (company_id);

-- 5. ref_departments
ALTER TABLE ref_departments ADD COLUMN company_id BIGINT NOT NULL DEFAULT 1;
ALTER TABLE ref_departments ADD INDEX idx_ref_departments_company_id (company_id);

-- 6. ref_shifts
ALTER TABLE ref_shifts ADD COLUMN company_id BIGINT NOT NULL DEFAULT 1;
ALTER TABLE ref_shifts ADD INDEX idx_ref_shifts_company_id (company_id);

-- 7. ref_leave_types
ALTER TABLE ref_leave_types ADD COLUMN company_id BIGINT NOT NULL DEFAULT 1;
ALTER TABLE ref_leave_types ADD INDEX idx_ref_leave_types_company_id (company_id);

-- 8. attendances
ALTER TABLE attendances ADD COLUMN company_id BIGINT NOT NULL DEFAULT 1;
ALTER TABLE attendances ADD INDEX idx_attendances_company_id (company_id);

-- 9. leave_balances
ALTER TABLE leave_balances ADD COLUMN company_id BIGINT NOT NULL DEFAULT 1;
ALTER TABLE leave_balances ADD INDEX idx_leave_balances_company_id (company_id);

-- 10. leave_requests
ALTER TABLE leave_requests ADD COLUMN company_id BIGINT NOT NULL DEFAULT 1;
ALTER TABLE leave_requests ADD INDEX idx_leave_requests_company_id (company_id);

-- 11. overtimes
ALTER TABLE overtimes ADD COLUMN company_id BIGINT NOT NULL DEFAULT 1;
ALTER TABLE overtimes ADD INDEX idx_overtimes_company_id (company_id);

-- 12. loans
ALTER TABLE loans ADD COLUMN company_id BIGINT NOT NULL DEFAULT 1;
ALTER TABLE loans ADD INDEX idx_loans_company_id (company_id);

-- 13. reimbursements
ALTER TABLE reimbursements ADD COLUMN company_id BIGINT NOT NULL DEFAULT 1;
ALTER TABLE reimbursements ADD INDEX idx_reimbursements_company_id (company_id);

-- 14. payrolls
ALTER TABLE payrolls ADD COLUMN company_id BIGINT NOT NULL DEFAULT 1;
ALTER TABLE payrolls ADD INDEX idx_payrolls_company_id (company_id);

-- 15. payroll_details
ALTER TABLE payroll_details ADD COLUMN company_id BIGINT NOT NULL DEFAULT 1;
ALTER TABLE payroll_details ADD INDEX idx_payroll_details_company_id (company_id);

-- 16. notifications
ALTER TABLE notifications ADD COLUMN company_id BIGINT NOT NULL DEFAULT 1;
ALTER TABLE notifications ADD INDEX idx_notifications_company_id (company_id);

-- 17. contracts
ALTER TABLE contracts ADD COLUMN company_id BIGINT NOT NULL DEFAULT 1;
ALTER TABLE contracts ADD INDEX idx_contracts_company_id (company_id);

-- 18. job_requisitions
ALTER TABLE job_requisitions ADD COLUMN company_id BIGINT NOT NULL DEFAULT 1;
ALTER TABLE job_requisitions ADD INDEX idx_job_requisitions_company_id (company_id);

-- 19. applicants
ALTER TABLE applicants ADD COLUMN company_id BIGINT NOT NULL DEFAULT 1;
ALTER TABLE applicants ADD INDEX idx_applicants_company_id (company_id);

-- 20. applicant_stage_histories
ALTER TABLE applicant_stage_histories ADD COLUMN company_id BIGINT NOT NULL DEFAULT 1;
ALTER TABLE applicant_stage_histories ADD INDEX idx_applicant_stage_histories_company_id (company_id);

-- 21. onboarding_templates
ALTER TABLE onboarding_templates ADD COLUMN company_id BIGINT NOT NULL DEFAULT 1;
ALTER TABLE onboarding_templates ADD INDEX idx_onboarding_templates_company_id (company_id);

-- 22. onboarding_template_items
ALTER TABLE onboarding_template_items ADD COLUMN company_id BIGINT NOT NULL DEFAULT 1;
ALTER TABLE onboarding_template_items ADD INDEX idx_onboarding_template_items_company_id (company_id);

-- 23. onboarding_workflows
ALTER TABLE onboarding_workflows ADD COLUMN company_id BIGINT NOT NULL DEFAULT 1;
ALTER TABLE onboarding_workflows ADD INDEX idx_onboarding_workflows_company_id (company_id);

-- 24. onboarding_tasks
ALTER TABLE onboarding_tasks ADD COLUMN company_id BIGINT NOT NULL DEFAULT 1;
ALTER TABLE onboarding_tasks ADD INDEX idx_onboarding_tasks_company_id (company_id);

-- 25. finance_categories
ALTER TABLE finance_categories ADD COLUMN company_id BIGINT NOT NULL DEFAULT 1;
ALTER TABLE finance_categories ADD INDEX idx_finance_categories_company_id (company_id);

-- 26. finance_transactions
ALTER TABLE finance_transactions ADD COLUMN company_id BIGINT NOT NULL DEFAULT 1;
ALTER TABLE finance_transactions ADD INDEX idx_finance_transactions_company_id (company_id);
