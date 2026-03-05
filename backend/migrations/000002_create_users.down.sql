DROP INDEX IF EXISTS idx_users_subscription_status;
DROP INDEX IF EXISTS idx_users_plan_id;
DROP INDEX IF EXISTS idx_users_oauth;
ALTER TABLE users DROP CONSTRAINT IF EXISTS check_subscription_status;
ALTER TABLE users DROP CONSTRAINT IF EXISTS check_role;
DROP TABLE IF EXISTS users;
