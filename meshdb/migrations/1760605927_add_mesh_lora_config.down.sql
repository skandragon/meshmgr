-- Remove LoRa configuration fields from meshes table
ALTER TABLE meshes DROP CONSTRAINT IF EXISTS check_frequency_slot;
ALTER TABLE meshes DROP CONSTRAINT IF EXISTS check_modem_preset;
ALTER TABLE meshes DROP CONSTRAINT IF EXISTS check_lora_region;

ALTER TABLE meshes DROP COLUMN IF EXISTS frequency_slot;
ALTER TABLE meshes DROP COLUMN IF EXISTS modem_preset;
ALTER TABLE meshes DROP COLUMN IF EXISTS lora_region;
