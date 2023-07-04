-- DROP existing data structures & tables to start from scratch
DROP TABLE IF EXISTS contact;

DROP SEQUENCE IF EXISTS contact_id_seq;

DROP TYPE IF EXISTS link_precedence_enum;

-- CREATE a new enum for link_precedence field
-- Createing an enum since there are only 2 possible values
CREATE TYPE link_precedence_enum AS ENUM ('primary', 'secondary');

-- CREATE new table `contact` with fields defined in schema
CREATE TABLE IF NOT EXISTS contact (
    id SERIAL PRIMARY KEY,
    phone_number TEXT,
    email TEXT,
    linked_id INT,
    link_precedence link_precedence_enum,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- CREATE a foreign key reference to ID column in table `contact`
ALTER TABLE
    contact
ADD
    FOREIGN KEY (linked_id) REFERENCES contact(id) ON
DELETE
SET
    NULL;