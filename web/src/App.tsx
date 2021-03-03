import React from "react";
import { Route, Switch } from "wouter";
import Nav from "features/nav/Nav";
import Dash from "features/dash/Dash";
import Profile from "features/profile/Profile";
import DamageSim from "features/damage/Damage";
import { useDispatch } from "react-redux";
import { load } from "features/profile/profileSlice";

function App() {
  const dispatch = useDispatch();
  React.useEffect(() => {
    dispatch(load());
    // @ts-ignore
    const go = new Go(); // eslint-disable-line
    WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject).then((result) => {
      go.run(result.instance);
    });
  }, []);

  return (
    <div className="App">
      <Nav />
      <div
        style={{ marginTop: "55px", marginRight: "20px", marginLeft: "20px" }}
      >
        <Switch>
          <Route path="/" component={Dash} />
          <Route path="/profile" component={Profile} />
          <Route path="/damage" component={DamageSim} />
          <Route path="/artifact" component={NI} />
          <Route path="/entry/details/:id">{(params) => <NI />}</Route>
        </Switch>
      </div>
    </div>
  );
}

function NI() {
  return <div></div>;
}

export default App;
