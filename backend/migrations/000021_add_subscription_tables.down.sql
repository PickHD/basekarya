-- +migrate Down
ALTER TABLE companies DROP INDEX idx_companies_subscription_plan_id;
ALTER TABLE companies DROP COLUMN owner_user_id;
ALTER TABLE companies DROP COLUMN subscription_expires_at;
ALTER TABLE companies DROP COLUMN subscription_status;
ALTER TABLE companies DROP COLUMN subscription_plan_id;
DROP TABLE IF EXISTS subscription_plans;
