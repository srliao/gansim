import { Button, Checkbox } from "@blueprintjs/core";
import { RootState } from "app/store";
import React from "react";
import { useSelector } from "react-redux";
import { createWorkerFactory, useWorker } from "@shopify/react-web-worker";

const createSim = createWorkerFactory(() => import("workers/damage_sim"));

// Create new instance

function DamageSim() {
  const sim = useWorker(createSim);
  const { profiles } = useSelector((state: RootState) => {
    return {
      profiles: state.profile.profiles,
    };
  });

  const [checked, setChecked] = React.useState<Array<number>>([]);
  const [n, setN] = React.useState<number>(100000);
  const [bin, setBin] = React.useState<number>(100);

  const handleCheck = (id: number) => {
    return () => {
      let clone = [...checked];
      //if id already in checked, remove it
      let index = checked.findIndex((x) => x === id);
      if (index === -1) {
        clone.push(id);
      } else {
        clone.splice(index, 1);
      }
      setChecked(clone);
    };
  };

  const handleSim = async () => {
    console.log("starting sim");

    if (checked.length > 0) {
      let p = profiles.find((e) => e.id === checked[0]);
      if (p !== undefined) {
        sim.DSim(n, 10, 4, p).then((r) => {
          console.log(r);
        });
        console.log("this shouldn't block?");
      }
    }
  };

  let p = profiles.map((e, i) => {
    let isChecked = checked.findIndex((x) => x === e.id) !== -1;
    return (
      <Checkbox key={e.id} checked={isChecked} onChange={handleCheck(e.id)}>
        {e.label}
      </Checkbox>
    );
  });

  if (!window.Worker) {
    return <div>No web worker?? can't sim</div>;
  }
  return (
    <div className="row">
      <div className="col-xs-offset-3 col-xs-6">
        Select profiles to sim
        {p}
        <Button intent="primary" fill onClick={handleSim}>
          Sim
        </Button>
      </div>
    </div>
  );
}

export default DamageSim;
