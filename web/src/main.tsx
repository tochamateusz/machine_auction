import React, { useState } from "react";
import ReactDOM from "react-dom/client";
import App, { Auction } from "./App.tsx";
import "./index.css";
import { AppBar, Slide, Toolbar } from "@mui/material";
import { NavBar } from "./components/nav-bar.tsx";
import { ListAuction } from "./components/auction-list-item.tsx";
import { States, useObservedAuction } from "./useObserverAuction.ts";
import { ListRowsAuction } from "./components/auction-list-rows.tsx";

interface HideOnScrollProps {
  children: React.ReactElement;
}

interface ViewProps {
  auction: Auction;
  onClick?: (auction: Auction) => void;
}

function HideOnScroll(props: HideOnScrollProps) {
  const { children } = props;

  return (
    <Slide color="primary" appear={false} direction="down" in={true}>
      {children}
    </Slide>
  );
}


type ViewState =
  | { type: "DEFAULT", component: (p: ViewProps) => React.ReactNode }
  | { type: "ROWS", component: (p: ViewProps) => React.ReactNode }

export const AppHOF = () => {
  const auctions = useObservedAuction()
  const [view, setView] = useState<ViewState>({
    type: "DEFAULT", component: ({ auction, onClick }) =>
      <ListAuction
        key={auction.id}
        auction={auction}
        onClick={onClick}
      />

  });


  const nav = (observedAuctionsState: States) => {
    switch (observedAuctionsState.type) {
      case "AUCTIONS_LOADED": {
        return (
          <NavBar
            auctions={observedAuctionsState.auctions}
            onToggleView={() => {
              if (view.type === "DEFAULT") {
                setView({
                  type: "ROWS", component: ({ auction, onClick }) =>
                    <ListRowsAuction
                      key={auction.id}
                      auction={auction}
                      onClick={onClick}
                    />
                })
              } else {
                setView({
                  type: "DEFAULT", component: ({ auction, onClick }) =>
                    <ListAuction
                      key={auction.id}
                      auction={auction}
                      onClick={onClick}
                    />
                })
              }

            }} />
        )
      }
      default: {
        return (
          <NavBar
            auctions={[]}
            onToggleView={() => {
              if (view.type === "DEFAULT") {
                setView({
                  type: "ROWS", component: ({ auction, onClick }) =>
                    <ListRowsAuction
                      key={auction.id}
                      auction={auction}
                      onClick={onClick}
                    />
                })
              } else {
                setView({
                  type: "DEFAULT", component: ({ auction, onClick }) =>
                    <ListAuction
                      key={auction.id}
                      auction={auction}
                      onClick={onClick}
                    />
                })
              }

            }} />
        )
      }
    }
  }

  return (
    <React.StrictMode>
      <HideOnScroll>
        <AppBar>
          <Toolbar sx={{ background: "#011F26" }}>
            {nav(auctions)}
          </Toolbar>
        </AppBar>
      </HideOnScroll>
      <App view={view.component} observedAuctionsState={auctions} />
      <Toolbar />
    </React.StrictMode>
  )
}

ReactDOM.createRoot(document.getElementById("root")!).render(<AppHOF />)
