CREATE TABLE "accounts"
(
    "id"         bigserial PRIMARY KEY,
    "owner"      varchar     NOT NULL,
    "balance"    bigint      NOT NULL,
    "currency"   varchar     NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

select *
from accounts;

CREATE TABLE "entries"
(
    "id"         bigserial PRIMARY KEY,
    "account_id" bigint,
    "amount"     bigint      NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

select *
from entries;

CREATE TABLE "transfers"
(
    "id"              bigserial PRIMARY KEY,
    "from_account_id" bigint,
    "to_account_id"   bigint,
    "amount"          bigint NOT NULL,
    "created_at"      timestamptz DEFAULT (now())
);

select *
from transfers;

CREATE INDEX ON "accounts" ("owner");

CREATE INDEX ON "entries" ("account_id");

CREATE INDEX ON "transfers" ("from_account_id");

CREATE INDEX ON "transfers" ("to_account_id");

CREATE INDEX ON "transfers" ("from_account_id", "to_account_id");

COMMENT ON COLUMN "entries"."amount" IS 'can be negative or positive';

COMMENT ON COLUMN "transfers"."amount" IS 'must be positive';

ALTER TABLE IF EXISTS "entries"
    ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

ALTER TABLE IF EXISTS "transfers"
    ADD FOREIGN KEY ("from_account_id") REFERENCES "accounts" ("id");

ALTER TABLE IF EXISTS "transfers"
    ADD FOREIGN KEY ("to_account_id") REFERENCES "accounts" ("id");
