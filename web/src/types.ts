export interface IProfile {
  id: number;
  label: string;
  character: ICharacter;
  weapon: IWeapon;
  enemy: IEnemy;
  artifact_levels: number;
  artifact_main_stats: { [key: string]: string };
  abilities: IAbility[];
}

export interface IEnemy {
  level: number;
  phy_resist: number; //don't need to use a mod for this since they only have resist
  ele_resist: number;
}

export interface ICharacter {
  level: number;
  base_atk: number;
  mods: IModifier[];
}

export interface IWeapon {
  base_atk: number;
  mods: IModifier[];
}

export interface IAbility {
  id: number;
  multiplier: number;
  is_vape_melt: boolean;
  vape_melt_multiplier: number;
  is_physical: boolean;
  mods: IModifier[];
}

export interface IArtifact {
  level: number;
  slot: string;
  target_main_stat: string;
  main_stat: IStat;
  substat: IStat[];
}

export interface IArtifactSets {
  set: { [key: string]: IArtifact };
  mods: IModifier[];
}

export interface IStat {
  type: string;
  value: number;
}

export interface IStatProb {
  type: string;
  weight: number;
}

export interface IModifier {
  type: string;
  label: string;
  helper?: string;
  list: {
    value: number;
    desc: string;
  }[];
}

//artifact slots
export const SLOT_FLOWER = "FLOWER";
export const SLOT_FEATHER = "FEATHER";
export const SLOT_SANDS = "SANDS";
export const SLOT_GOBLET = "GOBLET";
export const SLOT_CIRCLET = "CIRCLET";

//stat types
export const STAT_TYPE_DEFP = "DEF%";
export const STAT_TYPE_DEF = "DEF";
export const STAT_TYPE_HP = "HP";
export const STAT_TYPE_HPP = "HP%";
export const STAT_TYPE_ATK = "ATK";
export const STAT_TYPE_ATKP = "ATK%";
export const STAT_TYPE_ER = "ER";
export const STAT_TYPE_EM = "EM";
export const STAT_TYPE_CR = "CR";
export const STAT_TYPE_CD = "CD";
export const STAT_TYPE_HEAL = "Heal";
export const STAT_TYPE_ELEP = "Ele%";
export const STAT_TYPE_PHYP = "Phys%";

//mod types
export const MOD_TYPE_ATKP = "ATKP_MOD";
export const MOD_TYPE_ELEP = "ELEP_MOD";
export const MOD_TYPE_PHYP = "PHYP_MOD";
export const MOD_TYPE_CR = "CR_MOD";
export const MOD_TYPE_CD = "CD_MOD";
export const MOD_TYPE_DMGP = "DMGP_MOD";
export const MOD_TYPE_EM = "EM_MOD";
export const MOD_TYPE_REACTION = "REACTION_MOD";
export const MOD_TYPE_RESIST = "RESIST_MOD";
export const MOD_TYPE_DEF = "DEF_MOD";

export const MAIN_STAT_BY_LVL: {
  [key: string]: number[];
} = {
  HP: [
    717,
    920,
    1123,
    1326,
    1530,
    1733,
    1936,
    2139,
    2342,
    2545,
    2749,
    2952,
    3155,
    3358,
    3561,
    3764,
    3967,
    4171,
    4374,
    4577,
    4780,
  ],
  ATK: [
    47,
    60,
    73,
    86,
    100,
    113,
    126,
    139,
    152,
    166,
    179,
    192,
    205,
    219,
    232,
    245,
    258,
    272,
    285,
    298,
    311,
  ],
  "HP%": [
    0.07,
    0.09,
    0.11,
    0.129,
    0.149,
    0.169,
    0.189,
    0.209,
    0.228,
    0.248,
    0.268,
    0.288,
    0.308,
    0.328,
    0.347,
    0.367,
    0.387,
    0.407,
    0.427,
    0.446,
    0.466,
  ],
  "ATK%": [
    0.07,
    0.09,
    0.11,
    0.129,
    0.149,
    0.169,
    0.189,
    0.209,
    0.228,
    0.248,
    0.268,
    0.288,
    0.308,
    0.328,
    0.347,
    0.367,
    0.387,
    0.407,
    0.427,
    0.446,
    0.466,
  ],
  "DEF%": [
    0.087,
    0.112,
    0.137,
    0.162,
    0.186,
    0.211,
    0.236,
    0.261,
    0.286,
    0.31,
    0.335,
    0.36,
    0.385,
    0.409,
    0.434,
    0.459,
    0.484,
    0.508,
    0.533,
    0.558,
    0.583,
  ],
  "Phys%": [
    0.087,
    0.112,
    0.137,
    0.162,
    0.162,
    0.211,
    0.236,
    0.261,
    0.286,
    0.31,
    0.335,
    0.36,
    0.385,
    0.409,
    0.434,
    0.459,
    0.484,
    0.508,
    0.533,
    0.558,
    0.583,
  ],
  "Ele%": [
    0.07,
    0.09,
    0.11,
    0.129,
    0.149,
    0.169,
    0.189,
    0.209,
    0.228,
    0.248,
    0.268,
    0.288,
    0.308,
    0.328,
    0.347,
    0.367,
    0.387,
    0.407,
    0.427,
    0.446,
    0.466,
  ],
  EM: [
    0.078,
    0.1,
    0.122,
    0.144,
    0.166,
    0.188,
    0.21,
    0.232,
    0.254,
    0.276,
    0.298,
    0.32,
    0.342,
    0.364,
    0.386,
    0.408,
    0.43,
    0.452,
    0.474,
    0.496,
    0.518,
  ],
  ER: [
    0.047,
    0.06,
    0.074,
    0.087,
    0.1,
    0.114,
    0.127,
    0.14,
    0.154,
    0.167,
    0.18,
    0.193,
    0.207,
    0.22,
    0.233,
    0.247,
    0.26,
    0.273,
    0.287,
    0.3,
    0.311,
  ],
  CD: [
    0.093,
    0.119,
    0.146,
    0.172,
    0.199,
    0.225,
    0.255,
    0.278,
    0.305,
    0.331,
    0.358,
    0.384,
    0.411,
    0.437,
    0.463,
    0.49,
    0.516,
    0.543,
    0.569,
    0.596,
    0.622,
  ],
  CR: [
    0.054,
    0.069,
    0.084,
    0.1,
    0.115,
    0.13,
    0.145,
    0.161,
    0.176,
    0.191,
    0.206,
    0.222,
    0.237,
    0.252,
    0.267,
    0.283,
    0.298,
    0.313,
    0.328,
    0.344,
    0.359,
  ],
  Heal: [
    5.4,
    6.9,
    8.4,
    10,
    11.5,
    13,
    14.5,
    16.1,
    17.6,
    19.1,
    20.6,
    22.2,
    23.7,
    25.2,
    26.7,
    28.3,
    29.8,
    31.3,
    32.8,
    34.4,
    35.9,
  ],
};

export const MAIN_STAT_PROB_BY_SLOT: {
  [key: string]: IStatProb[];
} = {
  FLOWER: [{ type: STAT_TYPE_HP, weight: 1 }],
  FEATHER: [{ type: STAT_TYPE_ATK, weight: 1 }],
  SANDS: [
    { type: STAT_TYPE_DEFP, weight: 210 },
    { type: STAT_TYPE_HPP, weight: 210 },
    { type: STAT_TYPE_ATKP, weight: 204 },
    { type: STAT_TYPE_ER, weight: 81 },
    { type: STAT_TYPE_EM, weight: 81 },
  ],
  GOBLET: [
    { type: STAT_TYPE_DEFP, weight: 60 },
    { type: STAT_TYPE_HPP, weight: 35 },
    { type: STAT_TYPE_ATKP, weight: 49 },
    { type: STAT_TYPE_EM, weight: 6 },
    { type: STAT_TYPE_ELEP, weight: 75 },
    { type: STAT_TYPE_PHYP, weight: 12 },
  ],
  CIRCLET: [
    { type: STAT_TYPE_DEFP, weight: 52 },
    { type: STAT_TYPE_HPP, weight: 52 },
    { type: STAT_TYPE_ATKP, weight: 59 },
    { type: STAT_TYPE_EM, weight: 9 },
    { type: STAT_TYPE_CR, weight: 31 },
    { type: STAT_TYPE_CD, weight: 21 },
    { type: STAT_TYPE_HEAL, weight: 23 },
  ],
};

export const SUB_STAT_PROB: {
  [key: string]: {
    [key: string]: IStatProb[];
  };
} = {
  FLOWER: {
    HP: [
      { type: STAT_TYPE_DEF, weight: 108 },
      { type: STAT_TYPE_DEFP, weight: 86 },
      { type: STAT_TYPE_HP, weight: 0 },
      { type: STAT_TYPE_HPP, weight: 86 },
      { type: STAT_TYPE_ATK, weight: 112 },
      { type: STAT_TYPE_ATKP, weight: 71 },
      { type: STAT_TYPE_ER, weight: 78 },
      { type: STAT_TYPE_EM, weight: 71 },
      { type: STAT_TYPE_CR, weight: 71 },
      { type: STAT_TYPE_CD, weight: 63 },
    ],
  },
  FEATHER: {
    ATK: [
      { type: STAT_TYPE_DEF, weight: 105 },
      { type: STAT_TYPE_DEFP, weight: 72 },
      { type: STAT_TYPE_HP, weight: 101 },
      { type: STAT_TYPE_HPP, weight: 66 },
      { type: STAT_TYPE_ATK, weight: 0 },
      { type: STAT_TYPE_ATKP, weight: 78 },
      { type: STAT_TYPE_ER, weight: 80 },
      { type: STAT_TYPE_EM, weight: 68 },
      { type: STAT_TYPE_CR, weight: 50 },
      { type: STAT_TYPE_CD, weight: 76 },
    ],
  },
  SANDS: {
    "DEF%": [
      { type: STAT_TYPE_DEF, weight: 36 },
      { type: STAT_TYPE_DEFP, weight: 0 },
      { type: STAT_TYPE_HP, weight: 34 },
      { type: STAT_TYPE_HPP, weight: 18 },
      { type: STAT_TYPE_ATK, weight: 30 },
      { type: STAT_TYPE_ATKP, weight: 23 },
      { type: STAT_TYPE_ER, weight: 22 },
      { type: STAT_TYPE_EM, weight: 25 },
      { type: STAT_TYPE_CR, weight: 20 },
      { type: STAT_TYPE_CD, weight: 13 },
    ],
    "HP%": [
      { type: STAT_TYPE_DEF, weight: 31 },
      { type: STAT_TYPE_DEFP, weight: 21 },
      { type: STAT_TYPE_HP, weight: 34 },
      { type: STAT_TYPE_HPP, weight: 0 },
      { type: STAT_TYPE_ATK, weight: 33 },
      { type: STAT_TYPE_ATKP, weight: 24 },
      { type: STAT_TYPE_ER, weight: 25 },
      { type: STAT_TYPE_EM, weight: 25 },
      { type: STAT_TYPE_CR, weight: 20 },
      { type: STAT_TYPE_CD, weight: 13 },
    ],
    "ATK%": [
      { type: STAT_TYPE_DEF, weight: 31 },
      { type: STAT_TYPE_DEFP, weight: 27 },
      { type: STAT_TYPE_HP, weight: 35 },
      { type: STAT_TYPE_HPP, weight: 23 },
      { type: STAT_TYPE_ATK, weight: 30 },
      { type: STAT_TYPE_ATKP, weight: 0 },
      { type: STAT_TYPE_ER, weight: 17 },
      { type: STAT_TYPE_EM, weight: 19 },
      { type: STAT_TYPE_CR, weight: 19 },
      { type: STAT_TYPE_CD, weight: 15 },
    ],
    ER: [
      { type: STAT_TYPE_DEF, weight: 15 },
      { type: STAT_TYPE_DEFP, weight: 7 },
      { type: STAT_TYPE_HP, weight: 14 },
      { type: STAT_TYPE_HPP, weight: 9 },
      { type: STAT_TYPE_ATK, weight: 14 },
      { type: STAT_TYPE_ATKP, weight: 5 },
      { type: STAT_TYPE_ER, weight: 0 },
      { type: STAT_TYPE_EM, weight: 5 },
      { type: STAT_TYPE_CR, weight: 9 },
      { type: STAT_TYPE_CD, weight: 7 },
    ],
    EM: [
      { type: STAT_TYPE_DEF, weight: 6 },
      { type: STAT_TYPE_DEFP, weight: 8 },
      { type: STAT_TYPE_HP, weight: 13 },
      { type: STAT_TYPE_HPP, weight: 15 },
      { type: STAT_TYPE_ATK, weight: 17 },
      { type: STAT_TYPE_ATKP, weight: 7 },
      { type: STAT_TYPE_ER, weight: 6 },
      { type: STAT_TYPE_EM, weight: 0 },
      { type: STAT_TYPE_CR, weight: 9 },
      { type: STAT_TYPE_CD, weight: 7 },
    ],
  },
  GOBLET: {
    "DEF%": [
      { type: STAT_TYPE_DEF, weight: 28 },
      { type: STAT_TYPE_DEFP, weight: 0 },
      { type: STAT_TYPE_HP, weight: 32 },
      { type: STAT_TYPE_HPP, weight: 15 },
      { type: STAT_TYPE_ATK, weight: 26 },
      { type: STAT_TYPE_ATKP, weight: 16 },
      { type: STAT_TYPE_ER, weight: 19 },
      { type: STAT_TYPE_EM, weight: 15 },
      { type: STAT_TYPE_CR, weight: 18 },
      { type: STAT_TYPE_CD, weight: 17 },
    ],
    "HP%": [
      { type: STAT_TYPE_DEF, weight: 11 },
      { type: STAT_TYPE_DEFP, weight: 15 },
      { type: STAT_TYPE_HP, weight: 17 },
      { type: STAT_TYPE_HPP, weight: 0 },
      { type: STAT_TYPE_ATK, weight: 14 },
      { type: STAT_TYPE_ATKP, weight: 7 },
      { type: STAT_TYPE_ER, weight: 13 },
      { type: STAT_TYPE_EM, weight: 11 },
      { type: STAT_TYPE_CR, weight: 12 },
      { type: STAT_TYPE_CD, weight: 10 },
    ],
    "ATK%": [
      { type: STAT_TYPE_DEF, weight: 25 },
      { type: STAT_TYPE_DEFP, weight: 17 },
      { type: STAT_TYPE_HP, weight: 19 },
      { type: STAT_TYPE_HPP, weight: 14 },
      { type: STAT_TYPE_ATK, weight: 23 },
      { type: STAT_TYPE_ATKP, weight: 0 },
      { type: STAT_TYPE_ER, weight: 16 },
      { type: STAT_TYPE_EM, weight: 15 },
      { type: STAT_TYPE_CR, weight: 9 },
      { type: STAT_TYPE_CD, weight: 20 },
    ],
    EM: [
      { type: STAT_TYPE_DEF, weight: 1 },
      { type: STAT_TYPE_DEFP, weight: 1 },
      { type: STAT_TYPE_HP, weight: 1 },
      { type: STAT_TYPE_HPP, weight: 2 },
      { type: STAT_TYPE_ATK, weight: 3 },
      { type: STAT_TYPE_ATKP, weight: 1 },
      { type: STAT_TYPE_ER, weight: 5 },
      { type: STAT_TYPE_EM, weight: 0 },
      { type: STAT_TYPE_CR, weight: 2 },
      { type: STAT_TYPE_CD, weight: 2 },
    ],
    "Ele%": [
      { type: STAT_TYPE_DEF, weight: 33 },
      { type: STAT_TYPE_DEFP, weight: 17 },
      { type: STAT_TYPE_HP, weight: 30 },
      { type: STAT_TYPE_HPP, weight: 28 },
      { type: STAT_TYPE_ATK, weight: 30 },
      { type: STAT_TYPE_ATKP, weight: 21 },
      { type: STAT_TYPE_ER, weight: 22 },
      { type: STAT_TYPE_EM, weight: 24 },
      { type: STAT_TYPE_CR, weight: 24 },
      { type: STAT_TYPE_CD, weight: 18 },
    ],
    "Phys%": [
      { type: STAT_TYPE_DEF, weight: 7 },
      { type: STAT_TYPE_DEFP, weight: 2 },
      { type: STAT_TYPE_HP, weight: 2 },
      { type: STAT_TYPE_HPP, weight: 3 },
      { type: STAT_TYPE_ATK, weight: 7 },
      { type: STAT_TYPE_ATKP, weight: 2 },
      { type: STAT_TYPE_ER, weight: 6 },
      { type: STAT_TYPE_EM, weight: 3 },
      { type: STAT_TYPE_CR, weight: 3 },
      { type: STAT_TYPE_CD, weight: 3 },
    ],
  },
  CIRCLET: {
    "DEF%": [
      { type: STAT_TYPE_DEF, weight: 30 },
      { type: STAT_TYPE_DEFP, weight: 0 },
      { type: STAT_TYPE_HP, weight: 28 },
      { type: STAT_TYPE_HPP, weight: 12 },
      { type: STAT_TYPE_ATK, weight: 24 },
      { type: STAT_TYPE_ATKP, weight: 14 },
      { type: STAT_TYPE_ER, weight: 22 },
      { type: STAT_TYPE_EM, weight: 11 },
      { type: STAT_TYPE_CR, weight: 16 },
      { type: STAT_TYPE_CD, weight: 13 },
    ],
    "HP%": [
      { type: STAT_TYPE_DEF, weight: 25 },
      { type: STAT_TYPE_DEFP, weight: 17 },
      { type: STAT_TYPE_HP, weight: 19 },
      { type: STAT_TYPE_HPP, weight: 0 },
      { type: STAT_TYPE_ATK, weight: 17 },
      { type: STAT_TYPE_ATKP, weight: 20 },
      { type: STAT_TYPE_ER, weight: 18 },
      { type: STAT_TYPE_EM, weight: 13 },
      { type: STAT_TYPE_CR, weight: 14 },
      { type: STAT_TYPE_CD, weight: 15 },
    ],
    "ATK%": [
      { type: STAT_TYPE_DEF, weight: 28 },
      { type: STAT_TYPE_DEFP, weight: 15 },
      { type: STAT_TYPE_HP, weight: 27 },
      { type: STAT_TYPE_HPP, weight: 18 },
      { type: STAT_TYPE_ATK, weight: 29 },
      { type: STAT_TYPE_ATKP, weight: 0 },
      { type: STAT_TYPE_ER, weight: 18 },
      { type: STAT_TYPE_EM, weight: 18 },
      { type: STAT_TYPE_CR, weight: 19 },
      { type: STAT_TYPE_CD, weight: 18 },
    ],
    EM: [
      { type: STAT_TYPE_DEF, weight: 3 },
      { type: STAT_TYPE_DEFP, weight: 3 },
      { type: STAT_TYPE_HP, weight: 4 },
      { type: STAT_TYPE_HPP, weight: 4 },
      { type: STAT_TYPE_ATK, weight: 3 },
      { type: STAT_TYPE_ATKP, weight: 3 },
      { type: STAT_TYPE_ER, weight: 6 },
      { type: STAT_TYPE_EM, weight: 0 },
      { type: STAT_TYPE_CR, weight: 1 },
      { type: STAT_TYPE_CD, weight: 1 },
    ],
    CR: [
      { type: STAT_TYPE_DEF, weight: 12 },
      { type: STAT_TYPE_DEFP, weight: 16 },
      { type: STAT_TYPE_HP, weight: 12 },
      { type: STAT_TYPE_HPP, weight: 9 },
      { type: STAT_TYPE_ATK, weight: 11 },
      { type: STAT_TYPE_ATKP, weight: 13 },
      { type: STAT_TYPE_ER, weight: 10 },
      { type: STAT_TYPE_EM, weight: 7 },
      { type: STAT_TYPE_CR, weight: 0 },
      { type: STAT_TYPE_CD, weight: 10 },
    ],
    CD: [
      { type: STAT_TYPE_DEF, weight: 9 },
      { type: STAT_TYPE_DEFP, weight: 2 },
      { type: STAT_TYPE_HP, weight: 11 },
      { type: STAT_TYPE_HPP, weight: 6 },
      { type: STAT_TYPE_ATK, weight: 12 },
      { type: STAT_TYPE_ATKP, weight: 4 },
      { type: STAT_TYPE_ER, weight: 11 },
      { type: STAT_TYPE_EM, weight: 6 },
      { type: STAT_TYPE_CR, weight: 7 },
      { type: STAT_TYPE_CD, weight: 0 },
    ],
    Heal: [
      { type: STAT_TYPE_DEF, weight: 10 },
      { type: STAT_TYPE_DEFP, weight: 5 },
      { type: STAT_TYPE_HP, weight: 10 },
      { type: STAT_TYPE_HPP, weight: 7 },
      { type: STAT_TYPE_ATK, weight: 8 },
      { type: STAT_TYPE_ATKP, weight: 6 },
      { type: STAT_TYPE_ER, weight: 8 },
      { type: STAT_TYPE_EM, weight: 8 },
      { type: STAT_TYPE_CR, weight: 8 },
      { type: STAT_TYPE_CD, weight: 4 },
    ],
  },
};

export const SUB_STAT_TIER: {
  [key: string]: number[];
} = {
  HP: [209, 239, 269, 299],
  DEF: [16, 19, 21, 23],
  ATK: [14, 16, 18, 19],
  "HP%": [0.041, 0.047, 0.053, 0.058],
  "DEF%": [0.051, 0.058, 0.066, 0.073],
  "ATK%": [0.041, 0.047, 0.053, 0.058],
  EM: [16, 19, 21, 23],
  ER: [0.045, 0.052, 0.058, 0.065],
  CR: [0.027, 0.031, 0.035, 0.039],
  CD: [0.054, 0.062, 0.07, 0.078],
};

export const PROB_MAX_SUB = 0.2221030042918455;
