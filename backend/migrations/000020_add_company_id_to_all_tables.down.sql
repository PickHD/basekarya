-- +migrate Down

-- 26. finance_transactions
ALTER TABLE finance_transactions DROP INDEX idx_finance_transactions_company_id;
ALTER TABLE finance_transactions DROP COLUMN company_id;

-- 25. finance_categories
ALTER TABLE finance_categories DROP INDEX idx_finance_categories_company_id;
ALTER TABLE finance_categories DROP COLUMN company_id;

-- 24. onboarding_tasks
ALTER TABLE onboarding_tasks DROP INDEX idx_onboarding_tasks_company_id;
ALTER TABLE onboarding_tasks DROP COLUMN company_id;

-- 23. onboarding_workflows
ALTER TABLE onboarding_workflows DROP INDEX idx_onboarding_workflows_company_id;
ALTER TABLE onboarding_workflows DROP COLUMN company_id;

-- 22. onboarding_template_items
ALTER TABLE onboarding_template_items DROP INDEX idx_onboarding_template_items_company_id;
ALTER TABLE onboarding_template_items DROP COLUMN company_id;

-- 21. onboarding_templates
ALTER TABLE onboarding_templates DROP INDEX idx_onboarding_templates_company_id;
ALTER TABLE onboarding_templates DROP COLUMN company_id;

-- 20. applicant_stage_histories
ALTER TABLE applicant_stage_histories DROP INDEX idx_applicant_stage_histories_company_id;
ALTER TABLE applicant_stage_histories DROP COLUMN company_id;

-- 19. applicants
ALTER TABLE applicants DROP INDEX idx_applicants_company_id;
ALTER TABLE applicants DROP COLUMN company_id;

-- 18. job_requisitions
ALTER TABLE job_requisitions DROP INDEX idx_job_requisitions_company_id;
ALTER TABLE job_requisitions DROP COLUMN company_id;

-- 17. contracts
ALTER TABLE contracts DROP INDEX idx_contracts_company_id;
ALTER TABLE contracts DROP COLUMN company_id;

-- 16. notifications
ALTER TABLE notifications DROP INDEX idx_notifications_company_id;
ALTER TABLE notifications DROP COLUMN company_id;

-- 15. payroll_details
ALTER TABLE payroll_details DROP INDEX idx_payroll_details_company_id;
ALTER TABLE payroll_details DROP COLUMN company_id;

-- 14. payrolls
ALTER TABLE payrolls DROP INDEX idx_payrolls_company_id;
ALTER TABLE payrolls DROP COLUMN company_id;

-- 13. reimbursements
ALTER TABLE reimbursements DROP INDEX idx_reimbursements_company_id;
ALTER TABLE reimbursements DROP COLUMN company_id;

-- 12. loans
ALTER TABLE loans DROP INDEX idx_loans_company_id;
ALTER TABLE loans DROP COLUMN company_id;

-- 11. overtimes
ALTER TABLE overtimes DROP INDEX idx_overtimes_company_id;
ALTER TABLE overtimes DROP COLUMN company_id;

-- 10. leave_requests
ALTER TABLE leave_requests DROP INDEX idx_leave_requests_company_id;
ALTER TABLE leave_requests DROP COLUMN company_id;

-- 9. leave_balances
ALTER TABLE leave_balances DROP INDEX idx_leave_balances_company_id;
ALTER TABLE leave_balances DROP COLUMN company_id;

-- 8. attendances
ALTER TABLE attendances DROP INDEX idx_attendances_company_id;
ALTER TABLE attendances DROP COLUMN company_id;

-- 7. ref_leave_types
ALTER TABLE ref_leave_types DROP INDEX idx_ref_leave_types_company_id;
ALTER TABLE ref_leave_types DROP COLUMN company_id;

-- 6. ref_shifts
ALTER TABLE ref_shifts DROP INDEX idx_ref_shifts_company_id;
ALTER TABLE ref_shifts DROP COLUMN company_id;

-- 5. ref_departments
ALTER TABLE ref_departments DROP INDEX idx_ref_departments_company_id;
ALTER TABLE ref_departments DROP COLUMN company_id;

-- 4. role_permissions
ALTER TABLE role_permissions DROP INDEX idx_role_permissions_company_id;
ALTER TABLE role_permissions DROP COLUMN company_id;

-- 3. roles — restore original unique index
ALTER TABLE roles DROP INDEX idx_roles_company_id;
ALTER TABLE roles DROP INDEX idx_roles_name_company_id;
ALTER TABLE roles ADD UNIQUE INDEX name (name);
ALTER TABLE roles DROP COLUMN company_id;

-- 2. employees
ALTER TABLE employees DROP INDEX idx_employees_company_id;
ALTER TABLE employees DROP COLUMN company_id;

-- 1. users — restore original unique index
ALTER TABLE users DROP INDEX idx_users_company_id;
ALTER TABLE users DROP INDEX idx_users_username_company_id;
ALTER TABLE users ADD UNIQUE INDEX username (username);
ALTER TABLE users DROP COLUMN company_id;
ALTER TABLE users DROP COLUMN is_platform_admin;
