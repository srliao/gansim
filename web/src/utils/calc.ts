import { log } from "console";
import {
  IArtifactSets,
  IProfile,
  MOD_TYPE_ATKP,
  MOD_TYPE_CD,
  MOD_TYPE_CR,
  MOD_TYPE_DEF,
  MOD_TYPE_DMGP,
  MOD_TYPE_ELEP,
  MOD_TYPE_EM,
  MOD_TYPE_PHYP,
  MOD_TYPE_REACTION,
  MOD_TYPE_RESIST,
  STAT_TYPE_ATK,
  STAT_TYPE_ATKP,
  STAT_TYPE_CD,
  STAT_TYPE_CR,
  STAT_TYPE_ELEP,
  STAT_TYPE_EM,
  STAT_TYPE_PHYP,
} from "types";

export interface IDmgResult {
  normal: number;
  average: number;
  crit: number;
}

export interface ISimResult {
  min: number;
  max: number;
  mean: number;
  sd: number;
}

export default function calc(p: IProfile, s: IArtifactSets): IDmgResult[] {
  let debug = false;

  let as: { [key: string]: number } = {};

  for (const [key, artifact] of Object.entries(s.set)) {
    if (debug) console.log(artifact);

    if (!(artifact.main_stat.type in as)) {
      as[artifact.main_stat.type] = 0;
    }
    as[artifact.main_stat.type] =
      as[artifact.main_stat.type] + artifact.main_stat.value;

    artifact.substat.forEach((sub) => {
      if (!(sub.type in as)) {
        as[sub.type] = 0;
      }
      as[sub.type] = as[sub.type] + sub.value;
    });
  }

  if (debug) console.log("artifact total stats", as);

  let baseAtk = p.character.base_atk + p.weapon.base_atk;

  if (debug) console.log("char + weapon base atk: ", baseAtk);

  let r: IDmgResult[] = [];

  //add character mods
  let cmod: { [key: string]: number } = {};

  p.character.mods.forEach((mod) => {
    // console.log(mod);
    let sum = mod.list.reduce((a, b) => a + b.value, 0);
    if (!(mod.type in cmod)) {
      cmod[mod.type] = 0;
    }
    cmod[mod.type] = cmod[mod.type] + sum;
  });

  if (debug) console.log("total character mods", cmod);

  //weapon mods not implemented yet

  p.abilities.forEach((abil) => {
    //map the mods
    let amod: { [key: string]: number } = {};
    abil.mods.forEach((mod) => {
      let sum = mod.list.reduce((a, b) => a + b.value, 0);
      if (!(mod.type in amod)) {
        amod[mod.type] = 0;
      }
      amod[mod.type] = amod[mod.type] + sum;
    });

    //artifact stats
    let atk = as[STAT_TYPE_ATK] ? as[STAT_TYPE_ATK] : 0;
    let atkp = as[STAT_TYPE_ATKP] ? as[STAT_TYPE_ATKP] : 0;
    let cr = as[STAT_TYPE_CR] ? as[STAT_TYPE_CR] : 0;
    let cd = as[STAT_TYPE_CD] ? as[STAT_TYPE_CD] : 0;
    let elep = as[STAT_TYPE_ELEP] ? as[STAT_TYPE_ELEP] : 0;
    let phyp = as[STAT_TYPE_PHYP] ? as[STAT_TYPE_PHYP] : 0;
    let em = as[STAT_TYPE_EM] ? as[STAT_TYPE_EM] : 0;

    //add atk  % mod
    atkp += cmod[MOD_TYPE_ATKP] ? cmod[MOD_TYPE_ATKP] : 0;
    atkp += amod[MOD_TYPE_ATKP] ? amod[MOD_TYPE_ATKP] : 0;

    let totalAtk = baseAtk;
    totalAtk = totalAtk * (1 + atkp) + atk;

    elep += cmod[MOD_TYPE_ELEP] ? cmod[MOD_TYPE_ELEP] : 0;
    elep += amod[MOD_TYPE_ELEP] ? amod[MOD_TYPE_ELEP] : 0;

    phyp += cmod[MOD_TYPE_PHYP] ? cmod[MOD_TYPE_PHYP] : 0;
    phyp += amod[MOD_TYPE_PHYP] ? amod[MOD_TYPE_PHYP] : 0;

    //dmg mods
    let dmg_bonus: number = 0;
    if (abil.is_physical) {
      dmg_bonus += phyp;
    } else {
      dmg_bonus += elep;
    }

    dmg_bonus += amod[MOD_TYPE_DMGP] ? amod[MOD_TYPE_DMGP] : 0;
    dmg_bonus += cmod[MOD_TYPE_DMGP] ? cmod[MOD_TYPE_DMGP] : 0;

    //crit mods
    cr += amod[MOD_TYPE_CR] ? amod[MOD_TYPE_CR] : 0;
    cr += cmod[MOD_TYPE_CR] ? cmod[MOD_TYPE_CR] : 0;

    cd += amod[MOD_TYPE_CD] ? amod[MOD_TYPE_CD] : 0;
    cd += cmod[MOD_TYPE_CD] ? cmod[MOD_TYPE_CD] : 0;

    //cap cr at 1, cap cr/cd min at 0
    if (cr > 1) {
      cr = 1;
    }
    if (cr < 0) {
      cr = 0;
    }
    if (cd < 0) {
      cd = 0;
    }

    //add em mod
    em += amod[MOD_TYPE_EM] ? amod[MOD_TYPE_EM] : 0;
    em += cmod[MOD_TYPE_EM] ? cmod[MOD_TYPE_EM] : 0;
    let def_shred = 0;
    def_shred += amod[MOD_TYPE_DEF] ? amod[MOD_TYPE_DEF] : 0;
    def_shred += cmod[MOD_TYPE_DEF] ? cmod[MOD_TYPE_DEF] : 0;

    let def_adj =
      (100 + p.character.level) /
      (100 + p.character.level + (100 + p.enemy.level) * (1 - def_shred));

    let res = 0;
    if (abil.is_physical) {
      res += p.enemy.phy_resist;
    } else {
      res += p.enemy.ele_resist;
    }
    res += amod[MOD_TYPE_RESIST] ? amod[MOD_TYPE_RESIST] : 0;
    res += cmod[MOD_TYPE_RESIST] ? cmod[MOD_TYPE_RESIST] : 0;

    let res_adj = 0;
    if (res < 0) {
      res_adj = 1 - res / 2;
    } else if (res < 0.75) {
      res_adj = 1 - res;
    } else {
      res_adj = 1 / (4 * res + 1);
    }

    let react_bonus = 0;
    react_bonus += amod[MOD_TYPE_REACTION] ? amod[MOD_TYPE_REACTION] : 0;
    react_bonus += cmod[MOD_TYPE_REACTION] ? cmod[MOD_TYPE_REACTION] : 0;
    let vmMult = 1.0;
    if (abil.is_vape_melt) {
      vmMult =
        abil.vape_melt_multiplier *
          (1 + 0.00189266831 * em * Math.exp(-0.000505 * em)) +
        react_bonus;
    }

    if (debug) {
      console.log("--------------------------\n");
      console.log("\tAdjusted stats -");
      console.log(" total atk", totalAtk);
      console.log(" atkp", atkp);
      console.log(" cc", cr);
      console.log(" cd", cd);
      console.log(" dmg", dmg_bonus);
      console.log(" em", em);
      console.log("\n");
      console.log("\tModifiers -");
      console.log(" talent", abil.multiplier);
      console.log(" def adj", def_adj);
      console.log(" res adj", res_adj);
      console.log(" vape/melt", vmMult);
      console.log("\n");
    }

    let normalDmg =
      totalAtk * (1 + dmg_bonus) * abil.multiplier * def_adj * res_adj * vmMult;
    let critDmg = normalDmg * (1 + cd);
    let avgDmg = normalDmg * (1 + cr * cd);

    if (debug) console.log("avg dmg", avgDmg);

    r.push({
      normal: normalDmg,
      average: avgDmg,
      crit: critDmg,
    });
  });

  return r;
}
