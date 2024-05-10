ALTER TABLE "transfers"
    DROP CONSTRAINT "transfers_to_account_id_fkey",
    DROP CONSTRAINT "transfers_from_account_id_fkey";

ALTER TABLE "entries"
    DROP CONSTRAINT "entries_account_id_fkey";

DROP TABLE IF EXISTS "transfers";
DROP TABLE IF EXISTS "entries";
DROP TABLE IF EXISTS "accounts";