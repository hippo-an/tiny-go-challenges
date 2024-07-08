CREATE TABLE "sessions" (
    "id"                    uuid            PRIMARY KEY,
    "user_id"               bigint          NOT NULL,
    "refersh_token"          varchar         NOT NULL,
    "user_agent"             varchar         NOT NULL,
    "client_ip"                 varchar         NOT NULL,
    "is_blocked"                 boolean         NOT NULL  DEFAULT false,
    "expires_at"   timestamptz     NOT NULL,
    "created_at"            timestamptz     NOT NULL DEFAULT (now())
);

ALTER TABLE IF EXISTS "sessions" ADD CONSTRAINT "fk_sessions_user_id" FOREIGN KEY ("user_id") REFERENCES "users" ("id");