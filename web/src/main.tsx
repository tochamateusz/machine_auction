import React from "react";
import ReactDOM from "react-dom/client";
import App from "./App.tsx";
import "./index.css";
import { AppBar, Slide, Toolbar } from "@mui/material";
import { NavBar } from "./components/nav-bar.tsx";

interface Props {
  window?: () => Window;
  children: React.ReactElement;
}

function HideOnScroll(props: Props) {
  const { children } = props;

  return (
    <Slide color="primary" appear={false} direction="down" in={true}>
      {children}
    </Slide>
  );
}

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <HideOnScroll>
      <AppBar>
        <Toolbar sx={{ background: "#011F26" }}>
          <NavBar />
        </Toolbar>
      </AppBar>
    </HideOnScroll>
    <Toolbar />
    <App />
  </React.StrictMode>,
);
