import { IProfile } from "types";
import { IDmgResult, ISimResult } from "utils/calc";
import { createWorkerFactory } from "@shopify/react-web-worker";

const createWorker = createWorkerFactory(() => import("./damage_sim_worker"));

export function DSim(
  n: number,
  b: number,
  w: number,
  p: IProfile
): Promise<ISimResult> {
  console.log("damage sim started");

  const do_work = async (id: number, num: number) => {
    let r: IDmgResult[] = [];
    let count = 0;
    const worker = createWorker();

    //send out worker while count < n
    while (count < num) {
      count++;
      let result = await worker.DSimWorker(id, p);
      r.push(result);
    }

    return r;
  };

  //fire up 4 workers

  return Promise.all([
    do_work(1, n / 8),
    do_work(2, n / 8),
    do_work(3, n / 8),
    do_work(4, n / 8),
    do_work(5, n / 8),
    do_work(6, n / 8),
    do_work(7, n / 8),
    do_work(8, n / 8),
  ]).then((r) => {
    let results: IDmgResult[] = [];
    console.log("all 4 promise done", r);
    r.forEach((x) => {
      results.push(...x);
    });

    let out: ISimResult = {
      min: 1000000000000,
      max: -1000000000000,
      mean: 0,
      sd: 0,
    };
    let ss = 0;
    let sum = 0;
    results.forEach((x) => {
      sum += x.average;
      if (out.max < x.average) out.max = x.average;
      if (out.min > x.average) out.min = x.average;
    });

    out.mean = sum / n;
    ss = results.reduce(
      (a, b) => a + (b.average - out.mean) * (b.average - out.mean),
      0
    );
    out.sd = Math.sqrt(ss / n);
    return out;
  });
}
