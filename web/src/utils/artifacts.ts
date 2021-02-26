import {
  IArtifact,
  IStatProb,
  MAIN_STAT_BY_LVL,
  MAIN_STAT_PROB_BY_SLOT,
  PROB_MAX_SUB,
  SUB_STAT_PROB,
  SUB_STAT_TIER,
} from "types";

export default function randArtifact(
  slot: string,
  main_stat: string,
  lvl: number
): IArtifact {
  let r: IArtifact = {
    level: lvl,
    slot: slot,
    target_main_stat: main_stat,
    main_stat: {
      type: main_stat,
      value: 0,
    },
    substat: [],
  };

  r.main_stat.value = MAIN_STAT_BY_LVL[main_stat][lvl];

  //how many sub stats
  let p = Math.random();
  let lines = 3;
  if (p <= PROB_MAX_SUB) {
    lines = 4;
  }

  let n = 4;
  if (lvl < 4 && lines < 4) {
    n = 3;
  }

  if (!(slot in SUB_STAT_PROB)) {
    console.log("invalid slot!");
    throw "invalid slot type";
  }
  if (!(main_stat in SUB_STAT_PROB[slot])) {
    console.log(
      "main stat probablity not found:",
      main_stat,
      MAIN_STAT_PROB_BY_SLOT[slot]
    );
    throw "main stat probablity not found";
  }
  //make a copy of the prob
  let prb: IStatProb[] = [];

  SUB_STAT_PROB[slot][main_stat].forEach((s) => {
    prb.push({
      type: s.type,
      weight: s.weight,
    });
  });

  //initial rolls
  for (let i = 0; i < n; i++) {
    let sumWeights: number = prb.reduce((a, b) => a + b.weight, 0);
    let found: number = -1;
    let pick = Math.random() * sumWeights;
    //loop through prb and find the pick
    for (let j = 0; j < prb.length && found === -1; j++) {
      if (pick < prb[j].weight) {
        found = j;
      }
      pick -= prb[j].weight;
    }
    if (found === -1) {
      throw "unexpected error - random stat not found";
    }

    let t = prb[found].type;

    //set weight to 0 for next ieration
    prb[found].weight = 0;

    let tier = Math.floor(Math.random() * 4);
    let val = SUB_STAT_TIER[t][tier];

    r.substat.push({
      type: t,
      value: val,
    });
  }

  //check how many upgrades to do
  let up = Math.floor(lvl / 4);

  if (lines === 3) {
    up--;
  }

  if (r.substat.length !== 4) {
    console.log("invalid artifact", r);
    throw "invalid artifact, less than 4 lines";
  }

  //do more rolls

  for (let i = 0; i < up; i++) {
    let pick = Math.floor(Math.random() * 4);
    let tier = Math.floor(Math.random() * 4);
    // console.log(SUB_STAT_TIER, pick, r.substat[pick]);
    r.substat[pick].value += SUB_STAT_TIER[r.substat[pick].type][tier];
  }

  return r;
}
