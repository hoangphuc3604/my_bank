CREATE TABLE "users" (
  "username" varchar PRIMARY KEY,
  "fullname" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "password" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "accounts" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");
ALTER TABLE "accounts" ADD CONSTRAINT "accounts_owner_currency_key" UNIQUE ("owner", "currency");