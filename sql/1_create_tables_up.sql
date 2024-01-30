CREATE TABLE IF NOT EXISTS "users" (
  "user_id" VARCHAR(255),
  "password" VARCHAR(255),
  "full_name" VARCHAR(255),
  "class" INTEGER,
  "email" VARCHAR(255),
  "role" VARCHAR(255),
  PRIMARY KEY ("user_id")
);

CREATE TABLE IF NOT EXISTS "attendances" (
  "username" VARCHAR(255),
  "attendance_id" VARCHAR(255),
  "punch_in_date" TIMESTAMP WITH TIME ZONE,
  "punch_out_date" TIMESTAMP WITH TIME ZONE,
  PRIMARY KEY ("attendance_id")
);

ALTER TABLE IF EXISTS "attendances" ADD CONSTRAINT "fk_user_attendance" FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");
