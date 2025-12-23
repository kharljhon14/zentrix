CREATE TABLE IF NOT EXISTS products(
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "quote_id" UUID NOT NULL,
    "title" VARCHAR(255) NOT NULL,
    "unit_price" INTEGER NOT NULL,
    "quantity" INTEGER NOT NULL,
    "discount" INTEGER DEFAULT 0,

    CONSTRAINT fk_qoute_id
        FOREIGN KEY ("quote_id") REFERENCES "quotes"(id)
);