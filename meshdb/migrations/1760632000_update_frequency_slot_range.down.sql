-- Revert frequency slot constraint back to 0-319
ALTER TABLE meshes DROP CONSTRAINT IF EXISTS check_frequency_slot;
ALTER TABLE meshes ADD CONSTRAINT check_frequency_slot
    CHECK (frequency_slot >= 0 AND frequency_slot <= 319);
