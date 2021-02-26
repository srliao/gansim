import calc, { IDmgResult } from "utils/calc";
import {
  IArtifactSets,
  IProfile,
  SLOT_CIRCLET,
  SLOT_FEATHER,
  SLOT_FLOWER,
  SLOT_GOBLET,
  SLOT_SANDS,
  STAT_TYPE_ATK,
  STAT_TYPE_HP,
} from "types";
import randArtifact from "utils/artifacts";

export function DSimWorker(id: number, p: IProfile): IDmgResult {
  // console.log(`worker ${id} got work!`);
  //create a set
  let s: IArtifactSets = {
    set: {},
    mods: [],
  };

  s.set[SLOT_FLOWER] = randArtifact(
    SLOT_FLOWER,
    STAT_TYPE_HP,
    p.artifact_levels
  );
  s.set[SLOT_FEATHER] = randArtifact(
    SLOT_FEATHER,
    STAT_TYPE_ATK,
    p.artifact_levels
  );
  s.set[SLOT_SANDS] = randArtifact(
    SLOT_SANDS,
    p.artifact_main_stats[SLOT_SANDS],
    p.artifact_levels
  );
  s.set[SLOT_GOBLET] = randArtifact(
    SLOT_GOBLET,
    p.artifact_main_stats[SLOT_GOBLET],
    p.artifact_levels
  );
  s.set[SLOT_CIRCLET] = randArtifact(
    SLOT_CIRCLET,
    p.artifact_main_stats[SLOT_CIRCLET],
    p.artifact_levels
  );

  //calculate dmg

  let r = calc(p, s);

  let out: IDmgResult = {
    normal: 0,
    average: 0,
    crit: 0,
  };
  r.forEach((x) => {
    out.normal += x.normal;
    out.average += x.average;
    out.crit += x.crit;
  });

  return out;
}
