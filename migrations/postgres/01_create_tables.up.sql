
CREATE TABLE "book" (
    "id" UUID PRIMARY KEY,
    "name" VARCHAR NOT NULL,
    "price" NUMERIC NOT NULL,
    "status" VARCHAR,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP
);


