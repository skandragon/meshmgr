-- Add LoRa configuration fields to meshes table
ALTER TABLE meshes ADD COLUMN lora_region TEXT DEFAULT 'US';
ALTER TABLE meshes ADD COLUMN modem_preset TEXT DEFAULT 'LONG_FAST';
ALTER TABLE meshes ADD COLUMN frequency_slot INTEGER DEFAULT 0;

-- Add CHECK constraints for valid values
ALTER TABLE meshes ADD CONSTRAINT check_lora_region
    CHECK (lora_region IN (
        'UNSET', 'US', 'EU_433', 'EU_868', 'CN', 'JP', 'ANZ', 'ANZ_433',
        'KR', 'TW', 'RU', 'IN', 'NZ_865', 'TH', 'UA_433', 'UA_868',
        'MY_433', 'MY_919', 'SG_923', 'KZ_433', 'KZ_863', 'BR_902',
        'PH_433', 'PH_868', 'PH_915', 'NP_865', 'LORA_24'
    ));

ALTER TABLE meshes ADD CONSTRAINT check_modem_preset
    CHECK (modem_preset IN (
        'SHORT_TURBO', 'SHORT_FAST', 'SHORT_SLOW',
        'MEDIUM_FAST', 'MEDIUM_SLOW',
        'LONG_FAST', 'LONG_MODERATE', 'LONG_SLOW',
        'VERY_LONG_SLOW'
    ));

ALTER TABLE meshes ADD CONSTRAINT check_frequency_slot
    CHECK (frequency_slot >= 0 AND frequency_slot <= 7);
