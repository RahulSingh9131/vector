-- Enable pgcrypto for gen_random_uuid()
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

---- create above / drop below ----

-- Disable pgcrypto
-- DROP EXTENSION IF EXISTS "pgcrypto";
