-- Update modem preset names from UPPER_SNAKE_CASE to CamelCase
-- and relax frequency slot constraint to support full range (0-319)

-- First, drop the old constraint to allow data updates
ALTER TABLE meshes DROP CONSTRAINT IF EXISTS check_modem_preset;

-- Then update existing data to use new CamelCase names
UPDATE meshes SET modem_preset = 'LongFast' WHERE modem_preset = 'LONG_FAST';
UPDATE meshes SET modem_preset = 'LongMod' WHERE modem_preset = 'LONG_MODERATE';
UPDATE meshes SET modem_preset = 'LongSlow' WHERE modem_preset = 'LONG_SLOW';
UPDATE meshes SET modem_preset = 'MediumFast' WHERE modem_preset = 'MEDIUM_FAST';
UPDATE meshes SET modem_preset = 'MediumSlow' WHERE modem_preset = 'MEDIUM_SLOW';
UPDATE meshes SET modem_preset = 'ShortFast' WHERE modem_preset = 'SHORT_FAST';
UPDATE meshes SET modem_preset = 'ShortSlow' WHERE modem_preset = 'SHORT_SLOW';
UPDATE meshes SET modem_preset = 'ShortTurbo' WHERE modem_preset = 'SHORT_TURBO';
-- Remove VERY_LONG_SLOW, default to LongSlow
UPDATE meshes SET modem_preset = 'LongSlow' WHERE modem_preset = 'VERY_LONG_SLOW';

-- Recreate the modem preset constraint with CamelCase names
ALTER TABLE meshes ADD CONSTRAINT check_modem_preset
    CHECK (modem_preset IN (
        'ShortTurbo', 'ShortFast', 'ShortSlow',
        'MediumFast', 'MediumSlow',
        'LongFast', 'LongMod', 'LongSlow'
    ));

-- Update default value to use CamelCase
ALTER TABLE meshes ALTER COLUMN modem_preset SET DEFAULT 'LongFast';

-- Relax frequency slot constraint to support full range (0-319 for CN/LongSlow)
ALTER TABLE meshes DROP CONSTRAINT IF EXISTS check_frequency_slot;
ALTER TABLE meshes ADD CONSTRAINT check_frequency_slot
    CHECK (frequency_slot >= 0 AND frequency_slot <= 319);
