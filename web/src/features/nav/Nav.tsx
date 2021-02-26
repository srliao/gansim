import React from "react";
import { useLocation } from "wouter";

import { Navbar, NavbarGroup, Button, Classes } from "@blueprintjs/core";

function Nav() {
  const [location, setLocation] = useLocation();
  const navigate = (url: string) => {
    return () => {
      setLocation(url);
    };
  };

  return (
    <Navbar fixedToTop={true}>
      <NavbarGroup>
        <Button
          className={Classes.MINIMAL}
          icon="home"
          text="Dash"
          onClick={navigate("/")}
        />
        <Button
          className={Classes.MINIMAL}
          icon="user"
          text="Profile"
          onClick={navigate("/profile")}
        />
        <Button
          className={Classes.MINIMAL}
          icon="selection"
          text="Damage"
          onClick={navigate("/damage")}
        />
        <Button
          className={Classes.MINIMAL}
          icon="star"
          text="Artifacts"
          onClick={navigate("/artifact")}
        />
      </NavbarGroup>
    </Navbar>
  );
}

export default Nav;
