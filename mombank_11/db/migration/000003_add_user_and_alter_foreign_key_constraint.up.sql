CREATE TABLE "users" (
    "id"                    bigserial       PRIMARY KEY,
    "username"              varchar         NOT NULL UNIQUE,
    "hashed_password"       varchar         NOT NULL,
    "full_name"             varchar         NOT NULL,
    "email"                 varchar         NOT NULL UNIQUE,
    "password_changed_at"   timestamptz     NOT NULL DEFAULT '0001-01-01 00:00:00Z'
    "created_at"            timestamptz     NOT NULL DEFAULT (now())
);


ALTER TABLE IF EXISTS "accounts"
ADD COLUMN "user_id" bigint NOT NULL;

ALTER TABLE IF EXISTS "accounts" ADD CONSTRAINT "fk_accounts_user_id" FOREIGN KEY ("user_id") REFERENCES "users" ("id");

-- user only has one account for one currency
ALTER TABLE IF EXISTS "accounts" ADD CONSTRAINT "uk_user_id_currency" UNIQUE ("user_id", "currency");