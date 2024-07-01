ALTER TABLE IF EXISTS "accounts" DROP CONSTRAINT uk_user_id_currency;
ALTER TABLE IF EXISTS "accounts" DROP CONSTRAINT fk_accounts_user_id;

ALTER TABLE IF EXISTS "accounts" DROP COLUMN "user_id";

DROP TABLE IF EXISTS  "users";
