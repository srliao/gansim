import {
  IModifier,
  MOD_TYPE_ATKP,
  MOD_TYPE_CR,
  MOD_TYPE_DMGP,
  MOD_TYPE_ELEP,
  MOD_TYPE_EM,
  MOD_TYPE_REACTION,
  MOD_TYPE_RESIST,
  MOD_TYPE_DEF,
  MOD_TYPE_CD,
} from "types";

export function defaultMods(): IModifier[] {
  return [
    {
      type: MOD_TYPE_ATKP,
      label: "Atk %",
      list: [],
    },
    {
      type: MOD_TYPE_ELEP,
      label: "Ele Damage %",
      list: [],
    },
    {
      type: MOD_TYPE_CR,
      label: "Crit %",
      list: [],
    },
    {
      type: MOD_TYPE_CD,
      label: "Crit Damage %",
      list: [],
    },
    {
      type: MOD_TYPE_DMGP,
      label: "Damage %",
      list: [],
    },
    {
      type: MOD_TYPE_EM,
      label: "Flat EM",
      helper: "Enter flat EM boost such as Albedo burst or Sucrose talent",
      list: [],
    },
    {
      type: MOD_TYPE_REACTION,
      label: "Reaction Bonus %",
      helper: "Forgot what this is for",
      list: [],
    },
    {
      type: MOD_TYPE_RESIST,
      label: "Resist %",
      helper: "Enter negative % for resist reduction such as Ganyu C1",
      list: [],
    },
    {
      type: MOD_TYPE_DEF,
      label: "Def %",
      helper: "Enter negative % for defense reduction",
      list: [],
    },
  ];
}
