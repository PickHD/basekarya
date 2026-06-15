ALTER TABLE onboarding_tasks DROP FOREIGN KEY onboarding_tasks_ibfk_2;
ALTER TABLE onboarding_tasks DROP COLUMN template_item_id;
ALTER TABLE onboarding_tasks DROP COLUMN department;
DROP TABLE IF EXISTS onboarding_template_items;
DROP TABLE IF EXISTS onboarding_templates;
