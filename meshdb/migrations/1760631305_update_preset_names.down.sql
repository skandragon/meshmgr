-- Revert modem preset names from CamelCase back to UPPER_SNAKE_CASE
-- and restore original frequency slot constraint

-- Revert data to use old UPPER_SNAKE_CASE names
UPDATE meshes SET modem_preset = 'LONG_FAST' WHERE modem_preset = 'LongFast';
UPDATE meshes SET modem_preset = 'LONG_MODERATE' WHERE modem_preset = 'LongMod';
UPDATE meshes SET modem_preset = 'LONG_SLOW' WHERE modem_preset = 'LongSlow';
UPDATE meshes SET modem_preset = 'MEDIUM_FAST' WHERE modem_preset = 'MediumFast';
UPDATE meshes SET modem_preset = 'MEDIUM_SLOW' WHERE modem_preset = 'MediumSlow';
UPDATE meshes SET modem_preset = 'SHORT_FAST' WHERE modem_preset = 'ShortFast';
UPDATE meshes SET modem_preset = 'SHORT_SLOW' WHERE modem_preset = 'ShortSlow';
UPDATE meshes SET modem_preset = 'SHORT_TURBO' WHERE modem_preset = 'ShortTurbo';

-- Restore default value
ALTER TABLE meshes ALTER COLUMN modem_preset SET DEFAULT 'LONG_FAST';

-- Restore old constraint
ALTER TABLE meshes DROP CONSTRAINT check_modem_preset;
ALTER TABLE meshes ADD CONSTRAINT check_modem_preset
    CHECK (modem_preset IN (
        'SHORT_TURBO', 'SHORT_FAST', 'SHORT_SLOW',
        'MEDIUM_FAST', 'MEDIUM_SLOW',
        'LONG_FAST', 'LONG_MODERATE', 'LONG_SLOW',
        'VERY_LONG_SLOW'
    ));

-- Restore original frequency slot constraint
ALTER TABLE meshes DROP CONSTRAINT check_frequency_slot;
ALTER TABLE meshes ADD CONSTRAINT check_frequency_slot
    CHECK (frequency_slot >= 0 AND frequency_slot <= 7);
