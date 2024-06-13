ALTER TABLE IF EXISTS "entries"
    ALTER COLUMN "account_id" DROP NOT NULL;

ALTER TABLE IF EXISTS "transfers"
    ALTER COLUMN "from_account_id" DROP NOT NULL;

ALTER TABLE IF EXISTS "transfers"
    ALTER COLUMN "to_account_id" DROP NOT NULL;

ALTER TABLE IF EXISTS "transfers"
    ALTER COLUMN "created_at" DROP NOT NULL;