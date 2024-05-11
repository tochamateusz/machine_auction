import { Box, Button, Typography } from "@mui/material";
import { useEffect, useRef, useState } from "react";
import { assertState } from "../types/assetType";
import axios from "axios";

type States =
  | { type: "INIT" }
  | { type: "LOADING_BACKUP" }
  | { type: "LOADED"; scrappingEvents: any; auctions: any };

export const NavBar = () => {
  const [backupState, setBackup] = useState<States>({
    type: "INIT",
  });

  const backup = async () => {
    assertState(backupState, "INIT", "LOADED");
    setBackup({ type: "LOADING_BACKUP" });
    const auctionsRes = await fetch(
      `${import.meta.env.VITE_DOMAIN}/backup/auctions.json`,
    );
    const scrappingEventsRes = await fetch(
      `${import.meta.env.VITE_DOMAIN}/backup/scrapping-events.json`,
    );
    const auctions = await auctionsRes.blob();
    const scrappingEvents = await scrappingEventsRes.blob();

    const auctionsBlob = window.URL.createObjectURL(auctions);
    const scrappingEventsBlob = window.URL.createObjectURL(scrappingEvents);

    setBackup({
      type: "LOADED",
      auctions: auctionsBlob,
      scrappingEvents: scrappingEventsBlob,
    });
  };

  const scrap = async () => {
    try {
      await axios.post(`${import.meta.env.VITE_DOMAIN}/scrapper/start`, {});
    } catch (e) {}
  };

  const auctionsLink = useRef<HTMLAnchorElement>(null);
  const scrappingEventsLink = useRef<HTMLAnchorElement>(null);

  useEffect(() => {
    switch (backupState.type) {
      case "INIT": {
        return;
      }
      case "LOADING_BACKUP": {
        return;
      }
      case "LOADED": {
        auctionsLink.current && auctionsLink.current.click();
        scrappingEventsLink.current && scrappingEventsLink.current.click();
        return;
      }
    }
  }, [backupState]);
  return (
    <Box
      sx={{
        display: "flex",
        width: "300px",
        justifyContent: "space-between",
      }}
    >
      <Typography variant="h6">Menu</Typography>
      <Button
        variant="contained"
        sx={{
          background: "#025E73",
          ":hover": { background: "#A5A692" },
        }}
        onClick={scrap}
      >
        Scrap
      </Button>
      <Button
        variant="contained"
        onClick={backup}
        sx={{
          background: "#025E73",
          ":hover": { background: "#A5A692" },
        }}
      >
        Get backup
      </Button>

      {backupState.type == "LOADED" ? (
        <Box>
          <a
            download="auctions.json"
            href={backupState.auctions}
            ref={auctionsLink}
          ></a>
          <a
            download="scrappingEvents.json"
            href={backupState.scrappingEvents}
            ref={scrappingEventsLink}
          ></a>
        </Box>
      ) : null}
    </Box>
  );
};
