-- Down migration script
-- Remove the foreign key constraint for used_referral_code_id in circle table
ALTER TABLE "circle"
DROP CONSTRAINT IF EXISTS "circle_used_referral_code_id_fkey";

-- Remove the used_referral_code_id column from circle table
ALTER TABLE "circle"
DROP COLUMN IF EXISTS "used_referral_code_id";

-- Remove the index on circle_id in referral table
DROP INDEX IF EXISTS "idx_referral_circle_id";

-- Remove the foreign key constraint for circle_id in referral table
ALTER TABLE "referral"
DROP CONSTRAINT IF EXISTS "referral_circle_id_fkey";

-- Drop the referral table
DROP TABLE IF EXISTS "referral";