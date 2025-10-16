# Meshtastic Region/Preset Frequency Slot Table

This table shows all combinations of LoRa regions and modem presets,
including the number of available frequency slots and the default slot
for each preset (calculated from hash of preset name).

## Region: ANZ

**Frequency Range:** 915.00 - 928.00 MHz (13.00 MHz span)
**Max Power:** 30 dBm
**Band Type:** Sub-GHz

| Preset | Bandwidth | Num Channels | Default Slot (0-idx/user) | Slot Range |
|--------|-----------|--------------|---------------------------|------------|
| LongFast    | 250.00 kHz |           52 | 19/20                     | 0-51       |
| LongSlow    | 125.00 kHz |          104 | 26/27                     | 0-103      |
| LongMod     | 125.00 kHz |          104 | 5/6                       | 0-103      |
| MediumSlow  | 250.00 kHz |           52 | 51/52                     | 0-51       |
| MediumFast  | 250.00 kHz |           52 | 44/45                     | 0-51       |
| ShortSlow   | 250.00 kHz |           52 | 22/23                     | 0-51       |
| ShortFast   | 250.00 kHz |           52 | 15/16                     | 0-51       |
| ShortTurbo  | 500.00 kHz |           26 | 23/24                     | 0-25       |

## Region: ANZ_433

**Frequency Range:** 433.05 - 434.79 MHz (1.74 MHz span)
**Max Power:** 14 dBm
**Band Type:** Sub-GHz

| Preset | Bandwidth | Num Channels | Default Slot (0-idx/user) | Slot Range |
|--------|-----------|--------------|---------------------------|------------|
| LongFast    | 250.00 kHz |            6 | 5/6                       | 0-5        |
| LongSlow    | 125.00 kHz |           13 | 0/1                       | 0-12       |
| LongMod     | 125.00 kHz |           13 | 5/6                       | 0-12       |
| MediumSlow  | 250.00 kHz |            6 | 3/4                       | 0-5        |
| MediumFast  | 250.00 kHz |            6 | 0/1                       | 0-5        |
| ShortSlow   | 250.00 kHz |            6 | 4/5                       | 0-5        |
| ShortFast   | 250.00 kHz |            6 | 1/2                       | 0-5        |
| ShortTurbo  | 500.00 kHz |            3 | 2/3                       | 0-2        |

## Region: BR_902

**Frequency Range:** 902.00 - 907.50 MHz (5.50 MHz span)
**Max Power:** 30 dBm
**Band Type:** Sub-GHz

| Preset | Bandwidth | Num Channels | Default Slot (0-idx/user) | Slot Range |
|--------|-----------|--------------|---------------------------|------------|
| LongFast    | 250.00 kHz |           22 | 7/8                       | 0-21       |
| LongSlow    | 125.00 kHz |           44 | 10/11                     | 0-43       |
| LongMod     | 125.00 kHz |           44 | 37/38                     | 0-43       |
| MediumSlow  | 250.00 kHz |           22 | 21/22                     | 0-21       |
| MediumFast  | 250.00 kHz |           22 | 18/19                     | 0-21       |
| ShortSlow   | 250.00 kHz |           22 | 12/13                     | 0-21       |
| ShortFast   | 250.00 kHz |           22 | 9/10                      | 0-21       |
| ShortTurbo  | 500.00 kHz |           11 | 7/8                       | 0-10       |

## Region: CN

**Frequency Range:** 470.00 - 510.00 MHz (40.00 MHz span)
**Max Power:** 19 dBm
**Band Type:** Sub-GHz

| Preset | Bandwidth | Num Channels | Default Slot (0-idx/user) | Slot Range |
|--------|-----------|--------------|---------------------------|------------|
| LongFast    | 250.00 kHz |          160 | 35/36                     | 0-159      |
| LongSlow    | 125.00 kHz |          320 | 186/187                   | 0-319      |
| LongMod     | 125.00 kHz |          320 | 53/54                     | 0-319      |
| MediumSlow  | 250.00 kHz |          160 | 139/140                   | 0-159      |
| MediumFast  | 250.00 kHz |          160 | 148/149                   | 0-159      |
| ShortSlow   | 250.00 kHz |          160 | 90/91                     | 0-159      |
| ShortFast   | 250.00 kHz |          160 | 99/100                    | 0-159      |
| ShortTurbo  | 500.00 kHz |           80 | 65/66                     | 0-79       |

## Region: EU_433

**Frequency Range:** 433.00 - 434.00 MHz (1.00 MHz span)
**Max Power:** 10 dBm
**Band Type:** Sub-GHz

| Preset | Bandwidth | Num Channels | Default Slot (0-idx/user) | Slot Range |
|--------|-----------|--------------|---------------------------|------------|
| LongFast    | 250.00 kHz |            4 | 3/4                       | 0-3        |
| LongSlow    | 125.00 kHz |            8 | 2/3                       | 0-7        |
| LongMod     | 125.00 kHz |            8 | 5/6                       | 0-7        |
| MediumSlow  | 250.00 kHz |            4 | 3/4                       | 0-3        |
| MediumFast  | 250.00 kHz |            4 | 0/1                       | 0-3        |
| ShortSlow   | 250.00 kHz |            4 | 2/3                       | 0-3        |
| ShortFast   | 250.00 kHz |            4 | 3/4                       | 0-3        |
| ShortTurbo  | 500.00 kHz |            2 | 1/2                       | 0-1        |

## Region: EU_868

**Frequency Range:** 869.40 - 869.65 MHz (0.25 MHz span)
**Max Power:** 27 dBm
**Band Type:** Sub-GHz

| Preset | Bandwidth | Num Channels | Default Slot (0-idx/user) | Slot Range |
|--------|-----------|--------------|---------------------------|------------|
| LongFast    | 250.00 kHz |            1 | 0/1                       | 0-0        |
| LongSlow    | 125.00 kHz |            2 | 0/1                       | 0-1        |
| LongMod     | 125.00 kHz |            2 | 1/2                       | 0-1        |
| MediumSlow  | 250.00 kHz |            1 | 0/1                       | 0-0        |
| MediumFast  | 250.00 kHz |            1 | 0/1                       | 0-0        |
| ShortSlow   | 250.00 kHz |            1 | 0/1                       | 0-0        |
| ShortFast   | 250.00 kHz |            1 | 0/1                       | 0-0        |
| ShortTurbo  | 500.00 kHz |            0 | N/A                       | N/A        |

## Region: IN

**Frequency Range:** 865.00 - 867.00 MHz (2.00 MHz span)
**Max Power:** 30 dBm
**Band Type:** Sub-GHz

| Preset | Bandwidth | Num Channels | Default Slot (0-idx/user) | Slot Range |
|--------|-----------|--------------|---------------------------|------------|
| LongFast    | 250.00 kHz |            8 | 3/4                       | 0-7        |
| LongSlow    | 125.00 kHz |           16 | 10/11                     | 0-15       |
| LongMod     | 125.00 kHz |           16 | 5/6                       | 0-15       |
| MediumSlow  | 250.00 kHz |            8 | 3/4                       | 0-7        |
| MediumFast  | 250.00 kHz |            8 | 4/5                       | 0-7        |
| ShortSlow   | 250.00 kHz |            8 | 2/3                       | 0-7        |
| ShortFast   | 250.00 kHz |            8 | 3/4                       | 0-7        |
| ShortTurbo  | 500.00 kHz |            4 | 1/2                       | 0-3        |

## Region: JP

**Frequency Range:** 920.50 - 923.50 MHz (3.00 MHz span)
**Max Power:** 13 dBm
**Band Type:** Sub-GHz

| Preset | Bandwidth | Num Channels | Default Slot (0-idx/user) | Slot Range |
|--------|-----------|--------------|---------------------------|------------|
| LongFast    | 250.00 kHz |           12 | 11/12                     | 0-11       |
| LongSlow    | 125.00 kHz |           24 | 2/3                       | 0-23       |
| LongMod     | 125.00 kHz |           24 | 13/14                     | 0-23       |
| MediumSlow  | 250.00 kHz |           12 | 3/4                       | 0-11       |
| MediumFast  | 250.00 kHz |           12 | 0/1                       | 0-11       |
| ShortSlow   | 250.00 kHz |           12 | 10/11                     | 0-11       |
| ShortFast   | 250.00 kHz |           12 | 7/8                       | 0-11       |
| ShortTurbo  | 500.00 kHz |            6 | 5/6                       | 0-5        |

## Region: KR

**Frequency Range:** 920.00 - 923.00 MHz (3.00 MHz span)
**Max Power:** 23 dBm
**Band Type:** Sub-GHz

| Preset | Bandwidth | Num Channels | Default Slot (0-idx/user) | Slot Range |
|--------|-----------|--------------|---------------------------|------------|
| LongFast    | 250.00 kHz |           12 | 11/12                     | 0-11       |
| LongSlow    | 125.00 kHz |           24 | 2/3                       | 0-23       |
| LongMod     | 125.00 kHz |           24 | 13/14                     | 0-23       |
| MediumSlow  | 250.00 kHz |           12 | 3/4                       | 0-11       |
| MediumFast  | 250.00 kHz |           12 | 0/1                       | 0-11       |
| ShortSlow   | 250.00 kHz |           12 | 10/11                     | 0-11       |
| ShortFast   | 250.00 kHz |           12 | 7/8                       | 0-11       |
| ShortTurbo  | 500.00 kHz |            6 | 5/6                       | 0-5        |

## Region: KZ_433

**Frequency Range:** 433.07 - 434.77 MHz (1.70 MHz span)
**Max Power:** 10 dBm
**Band Type:** Sub-GHz

| Preset | Bandwidth | Num Channels | Default Slot (0-idx/user) | Slot Range |
|--------|-----------|--------------|---------------------------|------------|
| LongFast    | 250.00 kHz |            6 | 5/6                       | 0-5        |
| LongSlow    | 125.00 kHz |           13 | 0/1                       | 0-12       |
| LongMod     | 125.00 kHz |           13 | 5/6                       | 0-12       |
| MediumSlow  | 250.00 kHz |            6 | 3/4                       | 0-5        |
| MediumFast  | 250.00 kHz |            6 | 0/1                       | 0-5        |
| ShortSlow   | 250.00 kHz |            6 | 4/5                       | 0-5        |
| ShortFast   | 250.00 kHz |            6 | 1/2                       | 0-5        |
| ShortTurbo  | 500.00 kHz |            3 | 2/3                       | 0-2        |

## Region: KZ_863

**Frequency Range:** 863.00 - 868.00 MHz (5.00 MHz span)
**Max Power:** 30 dBm
**Band Type:** Sub-GHz

| Preset | Bandwidth | Num Channels | Default Slot (0-idx/user) | Slot Range |
|--------|-----------|--------------|---------------------------|------------|
| LongFast    | 250.00 kHz |           20 | 15/16                     | 0-19       |
| LongSlow    | 125.00 kHz |           40 | 26/27                     | 0-39       |
| LongMod     | 125.00 kHz |           40 | 13/14                     | 0-39       |
| MediumSlow  | 250.00 kHz |           20 | 19/20                     | 0-19       |
| MediumFast  | 250.00 kHz |           20 | 8/9                       | 0-19       |
| ShortSlow   | 250.00 kHz |           20 | 10/11                     | 0-19       |
| ShortFast   | 250.00 kHz |           20 | 19/20                     | 0-19       |
| ShortTurbo  | 500.00 kHz |           10 | 5/6                       | 0-9        |

## Region: LORA_24

**Frequency Range:** 2400.00 - 2483.50 MHz (83.50 MHz span)
**Max Power:** 10 dBm
**Band Type:** 2.4 GHz (Wide LoRa)

| Preset | Bandwidth | Num Channels | Default Slot (0-idx/user) | Slot Range |
|--------|-----------|--------------|---------------------------|------------|
| LongFast    | 812.50 kHz |          102 | 5/6                       | 0-101      |
| LongSlow    | 406.25 kHz |          205 | 86/87                     | 0-204      |
| LongMod     | 406.25 kHz |          205 | 128/129                   | 0-204      |
| MediumSlow  | 812.50 kHz |          102 | 51/52                     | 0-101      |
| MediumFast  | 812.50 kHz |          102 | 12/13                     | 0-101      |
| ShortSlow   | 812.50 kHz |          102 | 94/95                     | 0-101      |
| ShortFast   | 812.50 kHz |          102 | 55/56                     | 0-101      |
| ShortTurbo  | 1625.00 kHz |           51 | 29/30                     | 0-50       |

## Region: MY_433

**Frequency Range:** 433.00 - 435.00 MHz (2.00 MHz span)
**Max Power:** 20 dBm
**Band Type:** Sub-GHz

| Preset | Bandwidth | Num Channels | Default Slot (0-idx/user) | Slot Range |
|--------|-----------|--------------|---------------------------|------------|
| LongFast    | 250.00 kHz |            8 | 3/4                       | 0-7        |
| LongSlow    | 125.00 kHz |           16 | 10/11                     | 0-15       |
| LongMod     | 125.00 kHz |           16 | 5/6                       | 0-15       |
| MediumSlow  | 250.00 kHz |            8 | 3/4                       | 0-7        |
| MediumFast  | 250.00 kHz |            8 | 4/5                       | 0-7        |
| ShortSlow   | 250.00 kHz |            8 | 2/3                       | 0-7        |
| ShortFast   | 250.00 kHz |            8 | 3/4                       | 0-7        |
| ShortTurbo  | 500.00 kHz |            4 | 1/2                       | 0-3        |

## Region: MY_919

**Frequency Range:** 919.00 - 924.00 MHz (5.00 MHz span)
**Max Power:** 27 dBm
**Band Type:** Sub-GHz

| Preset | Bandwidth | Num Channels | Default Slot (0-idx/user) | Slot Range |
|--------|-----------|--------------|---------------------------|------------|
| LongFast    | 250.00 kHz |           20 | 15/16                     | 0-19       |
| LongSlow    | 125.00 kHz |           40 | 26/27                     | 0-39       |
| LongMod     | 125.00 kHz |           40 | 13/14                     | 0-39       |
| MediumSlow  | 250.00 kHz |           20 | 19/20                     | 0-19       |
| MediumFast  | 250.00 kHz |           20 | 8/9                       | 0-19       |
| ShortSlow   | 250.00 kHz |           20 | 10/11                     | 0-19       |
| ShortFast   | 250.00 kHz |           20 | 19/20                     | 0-19       |
| ShortTurbo  | 500.00 kHz |           10 | 5/6                       | 0-9        |

## Region: NP_865

**Frequency Range:** 865.00 - 868.00 MHz (3.00 MHz span)
**Max Power:** 30 dBm
**Band Type:** Sub-GHz

| Preset | Bandwidth | Num Channels | Default Slot (0-idx/user) | Slot Range |
|--------|-----------|--------------|---------------------------|------------|
| LongFast    | 250.00 kHz |           12 | 11/12                     | 0-11       |
| LongSlow    | 125.00 kHz |           24 | 2/3                       | 0-23       |
| LongMod     | 125.00 kHz |           24 | 13/14                     | 0-23       |
| MediumSlow  | 250.00 kHz |           12 | 3/4                       | 0-11       |
| MediumFast  | 250.00 kHz |           12 | 0/1                       | 0-11       |
| ShortSlow   | 250.00 kHz |           12 | 10/11                     | 0-11       |
| ShortFast   | 250.00 kHz |           12 | 7/8                       | 0-11       |
| ShortTurbo  | 500.00 kHz |            6 | 5/6                       | 0-5        |

## Region: NZ_865

**Frequency Range:** 864.00 - 868.00 MHz (4.00 MHz span)
**Max Power:** 36 dBm
**Band Type:** Sub-GHz

| Preset | Bandwidth | Num Channels | Default Slot (0-idx/user) | Slot Range |
|--------|-----------|--------------|---------------------------|------------|
| LongFast    | 250.00 kHz |           16 | 3/4                       | 0-15       |
| LongSlow    | 125.00 kHz |           32 | 26/27                     | 0-31       |
| LongMod     | 125.00 kHz |           32 | 21/22                     | 0-31       |
| MediumSlow  | 250.00 kHz |           16 | 11/12                     | 0-15       |
| MediumFast  | 250.00 kHz |           16 | 4/5                       | 0-15       |
| ShortSlow   | 250.00 kHz |           16 | 10/11                     | 0-15       |
| ShortFast   | 250.00 kHz |           16 | 3/4                       | 0-15       |
| ShortTurbo  | 500.00 kHz |            8 | 1/2                       | 0-7        |

## Region: PH_433

**Frequency Range:** 433.00 - 434.70 MHz (1.70 MHz span)
**Max Power:** 10 dBm
**Band Type:** Sub-GHz

| Preset | Bandwidth | Num Channels | Default Slot (0-idx/user) | Slot Range |
|--------|-----------|--------------|---------------------------|------------|
| LongFast    | 250.00 kHz |            6 | 5/6                       | 0-5        |
| LongSlow    | 125.00 kHz |           13 | 0/1                       | 0-12       |
| LongMod     | 125.00 kHz |           13 | 5/6                       | 0-12       |
| MediumSlow  | 250.00 kHz |            6 | 3/4                       | 0-5        |
| MediumFast  | 250.00 kHz |            6 | 0/1                       | 0-5        |
| ShortSlow   | 250.00 kHz |            6 | 4/5                       | 0-5        |
| ShortFast   | 250.00 kHz |            6 | 1/2                       | 0-5        |
| ShortTurbo  | 500.00 kHz |            3 | 2/3                       | 0-2        |

## Region: PH_868

**Frequency Range:** 868.00 - 869.40 MHz (1.40 MHz span)
**Max Power:** 14 dBm
**Band Type:** Sub-GHz

| Preset | Bandwidth | Num Channels | Default Slot (0-idx/user) | Slot Range |
|--------|-----------|--------------|---------------------------|------------|
| LongFast    | 250.00 kHz |            5 | 0/1                       | 0-4        |
| LongSlow    | 125.00 kHz |           11 | 10/11                     | 0-10       |
| LongMod     | 125.00 kHz |           11 | 4/5                       | 0-10       |
| MediumSlow  | 250.00 kHz |            5 | 4/5                       | 0-4        |
| MediumFast  | 250.00 kHz |            5 | 3/4                       | 0-4        |
| ShortSlow   | 250.00 kHz |            5 | 0/1                       | 0-4        |
| ShortFast   | 250.00 kHz |            5 | 4/5                       | 0-4        |
| ShortTurbo  | 500.00 kHz |            2 | 1/2                       | 0-1        |

## Region: PH_915

**Frequency Range:** 915.00 - 918.00 MHz (3.00 MHz span)
**Max Power:** 24 dBm
**Band Type:** Sub-GHz

| Preset | Bandwidth | Num Channels | Default Slot (0-idx/user) | Slot Range |
|--------|-----------|--------------|---------------------------|------------|
| LongFast    | 250.00 kHz |           12 | 11/12                     | 0-11       |
| LongSlow    | 125.00 kHz |           24 | 2/3                       | 0-23       |
| LongMod     | 125.00 kHz |           24 | 13/14                     | 0-23       |
| MediumSlow  | 250.00 kHz |           12 | 3/4                       | 0-11       |
| MediumFast  | 250.00 kHz |           12 | 0/1                       | 0-11       |
| ShortSlow   | 250.00 kHz |           12 | 10/11                     | 0-11       |
| ShortFast   | 250.00 kHz |           12 | 7/8                       | 0-11       |
| ShortTurbo  | 500.00 kHz |            6 | 5/6                       | 0-5        |

## Region: RU

**Frequency Range:** 868.70 - 869.20 MHz (0.50 MHz span)
**Max Power:** 20 dBm
**Band Type:** Sub-GHz

| Preset | Bandwidth | Num Channels | Default Slot (0-idx/user) | Slot Range |
|--------|-----------|--------------|---------------------------|------------|
| LongFast    | 250.00 kHz |            2 | 1/2                       | 0-1        |
| LongSlow    | 125.00 kHz |            4 | 2/3                       | 0-3        |
| LongMod     | 125.00 kHz |            4 | 1/2                       | 0-3        |
| MediumSlow  | 250.00 kHz |            2 | 1/2                       | 0-1        |
| MediumFast  | 250.00 kHz |            2 | 0/1                       | 0-1        |
| ShortSlow   | 250.00 kHz |            2 | 0/1                       | 0-1        |
| ShortFast   | 250.00 kHz |            2 | 1/2                       | 0-1        |
| ShortTurbo  | 500.00 kHz |            1 | 0/1                       | 0-0        |

## Region: SG_923

**Frequency Range:** 917.00 - 925.00 MHz (8.00 MHz span)
**Max Power:** 20 dBm
**Band Type:** Sub-GHz

| Preset | Bandwidth | Num Channels | Default Slot (0-idx/user) | Slot Range |
|--------|-----------|--------------|---------------------------|------------|
| LongFast    | 250.00 kHz |           32 | 3/4                       | 0-31       |
| LongSlow    | 125.00 kHz |           64 | 58/59                     | 0-63       |
| LongMod     | 125.00 kHz |           64 | 53/54                     | 0-63       |
| MediumSlow  | 250.00 kHz |           32 | 11/12                     | 0-31       |
| MediumFast  | 250.00 kHz |           32 | 20/21                     | 0-31       |
| ShortSlow   | 250.00 kHz |           32 | 26/27                     | 0-31       |
| ShortFast   | 250.00 kHz |           32 | 3/4                       | 0-31       |
| ShortTurbo  | 500.00 kHz |           16 | 1/2                       | 0-15       |

## Region: TH

**Frequency Range:** 920.00 - 925.00 MHz (5.00 MHz span)
**Max Power:** 16 dBm
**Band Type:** Sub-GHz

| Preset | Bandwidth | Num Channels | Default Slot (0-idx/user) | Slot Range |
|--------|-----------|--------------|---------------------------|------------|
| LongFast    | 250.00 kHz |           20 | 15/16                     | 0-19       |
| LongSlow    | 125.00 kHz |           40 | 26/27                     | 0-39       |
| LongMod     | 125.00 kHz |           40 | 13/14                     | 0-39       |
| MediumSlow  | 250.00 kHz |           20 | 19/20                     | 0-19       |
| MediumFast  | 250.00 kHz |           20 | 8/9                       | 0-19       |
| ShortSlow   | 250.00 kHz |           20 | 10/11                     | 0-19       |
| ShortFast   | 250.00 kHz |           20 | 19/20                     | 0-19       |
| ShortTurbo  | 500.00 kHz |           10 | 5/6                       | 0-9        |

## Region: TW

**Frequency Range:** 920.00 - 925.00 MHz (5.00 MHz span)
**Max Power:** 27 dBm
**Band Type:** Sub-GHz

| Preset | Bandwidth | Num Channels | Default Slot (0-idx/user) | Slot Range |
|--------|-----------|--------------|---------------------------|------------|
| LongFast    | 250.00 kHz |           20 | 15/16                     | 0-19       |
| LongSlow    | 125.00 kHz |           40 | 26/27                     | 0-39       |
| LongMod     | 125.00 kHz |           40 | 13/14                     | 0-39       |
| MediumSlow  | 250.00 kHz |           20 | 19/20                     | 0-19       |
| MediumFast  | 250.00 kHz |           20 | 8/9                       | 0-19       |
| ShortSlow   | 250.00 kHz |           20 | 10/11                     | 0-19       |
| ShortFast   | 250.00 kHz |           20 | 19/20                     | 0-19       |
| ShortTurbo  | 500.00 kHz |           10 | 5/6                       | 0-9        |

## Region: UA_433

**Frequency Range:** 433.00 - 434.70 MHz (1.70 MHz span)
**Max Power:** 10 dBm
**Band Type:** Sub-GHz

| Preset | Bandwidth | Num Channels | Default Slot (0-idx/user) | Slot Range |
|--------|-----------|--------------|---------------------------|------------|
| LongFast    | 250.00 kHz |            6 | 5/6                       | 0-5        |
| LongSlow    | 125.00 kHz |           13 | 0/1                       | 0-12       |
| LongMod     | 125.00 kHz |           13 | 5/6                       | 0-12       |
| MediumSlow  | 250.00 kHz |            6 | 3/4                       | 0-5        |
| MediumFast  | 250.00 kHz |            6 | 0/1                       | 0-5        |
| ShortSlow   | 250.00 kHz |            6 | 4/5                       | 0-5        |
| ShortFast   | 250.00 kHz |            6 | 1/2                       | 0-5        |
| ShortTurbo  | 500.00 kHz |            3 | 2/3                       | 0-2        |

## Region: UA_868

**Frequency Range:** 868.00 - 868.60 MHz (0.60 MHz span)
**Max Power:** 14 dBm
**Band Type:** Sub-GHz

| Preset | Bandwidth | Num Channels | Default Slot (0-idx/user) | Slot Range |
|--------|-----------|--------------|---------------------------|------------|
| LongFast    | 250.00 kHz |            2 | 1/2                       | 0-1        |
| LongSlow    | 125.00 kHz |            4 | 2/3                       | 0-3        |
| LongMod     | 125.00 kHz |            4 | 1/2                       | 0-3        |
| MediumSlow  | 250.00 kHz |            2 | 1/2                       | 0-1        |
| MediumFast  | 250.00 kHz |            2 | 0/1                       | 0-1        |
| ShortSlow   | 250.00 kHz |            2 | 0/1                       | 0-1        |
| ShortFast   | 250.00 kHz |            2 | 1/2                       | 0-1        |
| ShortTurbo  | 500.00 kHz |            1 | 0/1                       | 0-0        |

## Region: US

**Frequency Range:** 902.00 - 928.00 MHz (26.00 MHz span)
**Max Power:** 30 dBm
**Band Type:** Sub-GHz

| Preset | Bandwidth | Num Channels | Default Slot (0-idx/user) | Slot Range |
|--------|-----------|--------------|---------------------------|------------|
| LongFast    | 250.00 kHz |          104 | 19/20                     | 0-103      |
| LongSlow    | 125.00 kHz |          208 | 26/27                     | 0-207      |
| LongMod     | 125.00 kHz |          208 | 5/6                       | 0-207      |
| MediumSlow  | 250.00 kHz |          104 | 51/52                     | 0-103      |
| MediumFast  | 250.00 kHz |          104 | 44/45                     | 0-103      |
| ShortSlow   | 250.00 kHz |          104 | 74/75                     | 0-103      |
| ShortFast   | 250.00 kHz |          104 | 67/68                     | 0-103      |
| ShortTurbo  | 500.00 kHz |           52 | 49/50                     | 0-51       |

## Technical Notes

1. **Slot Numbers:**
   - Internally: 0-indexed (0 to numChannels-1)
   - User-facing: 1-indexed (1 to numChannels)
   - Channel 0 in config = use hash-based default slot

2. **Default Slot Calculation:**
   - Hash = djb2_hash(preset_name)
   - Default Slot = Hash % numChannels

3. **Frequency Calculation:**
   - freq = freqStart + (bandwidth/2) + (slot_num × bandwidth)
   - This centers each channel within its allocated bandwidth

4. **Default PSK (Pre-Shared Key):**
   - All preset channels use: d4f1bb3a20290759f0bcffabcf4e6901
   - This allows devices on same preset to communicate without key exchange

5. **Hash Values by Preset:**

   **Frequency Slot Hash (djb2):**
   - LongFast   : 0x07c63403 ( 130429955)
   - LongMod    : 0xe8f69d35 (3908476213)
   - LongSlow   : 0x07cd833a ( 130908986)
   - MediumFast : 0x57163d94 (1461075348)
   - MediumSlow : 0x571d8ccb (1461554379)
   - ShortFast  : 0x23fb41a3 ( 603668899)
   - ShortSlow  : 0x240290da ( 604147930)
   - ShortTurbo : 0xa46bbe81 (2758524545)
