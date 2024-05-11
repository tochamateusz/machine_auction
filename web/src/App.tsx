import { useEffect, useState } from "react";
import "./App.css";
import axios from "axios";
import { Box, CircularProgress, Modal, Paper } from "@mui/material";
import { ListAuction } from "./components/auction-list-item";
import { AuctionGallery } from "./components/auction-gallery";

export type Auction = {
  id: string;
  image: string;
  name: string;
  year: string;
  price: string;
  end_date: string;
};

type States =
  | { type: "INIT" }
  | { type: "LOADING_AUCTIONS" }
  | { type: "AUCTIONS_LOADED"; auctions: Auction[] }
  | { type: "AUCTIONS_LOADING_ERROR"; error: any }
  | { type: "MODAL_OPEN"; auctions: Auction[]; selectedAuction: Auction };

function App() {
  const [observedAuctionsState, setObservedAuctionsState] = useState<States>({
    type: "INIT",
  });

  const getAuctions = async () => {
    setObservedAuctionsState({ type: "LOADING_AUCTIONS" });
    try {
      const auctions = await axios.get<Auction[]>(
        `${import.meta.env.VITE_DOMAIN}/scrapper`,
      );
      setObservedAuctionsState({
        type: "AUCTIONS_LOADED",
        auctions: auctions.data,
      });
    } catch (e) {
      setObservedAuctionsState({
        type: "AUCTIONS_LOADING_ERROR",
        error: e,
      });
    }
  };

  useEffect(() => {
    console.debug(observedAuctionsState);
    switch (observedAuctionsState.type) {
      case "INIT": {
        getAuctions();
      }
    }
  }, [observedAuctionsState]);

  switch (observedAuctionsState.type) {
    case "INIT": {
      return <></>;
    }
    case "LOADING_AUCTIONS": {
      return (
        <Paper elevation={0} sx={{ maxWidth: 256 }}>
          <Box sx={{ display: "flex" }}>
            <CircularProgress />
          </Box>
        </Paper>
      );
    }
    case "AUCTIONS_LOADED": {
      return (
        <Box sx={{ display: "flex", flexWrap: "wrap" }}>
          {observedAuctionsState.auctions.map((a) => {
            return (
              <ListAuction
                key={a.id}
                auction={a}
                onClick={(auction) => {
                  setObservedAuctionsState({
                    type: "MODAL_OPEN",
                    auctions: observedAuctionsState.auctions,
                    selectedAuction: auction,
                  });
                }}
              />
            );
          })}
        </Box>
      );
    }

    case "MODAL_OPEN": {
      return (
        <>
          <Box sx={{ display: "flex", flexWrap: "wrap" }}>
            {observedAuctionsState.auctions.map((a) => {
              return <ListAuction auction={a} />;
            })}
          </Box>
          <Modal
            open={true}
            onClose={() => {
              setObservedAuctionsState({
                type: "AUCTIONS_LOADED",
                auctions: observedAuctionsState.auctions,
              });
            }}
          >
            <AuctionGallery auction={observedAuctionsState.selectedAuction} />
          </Modal>
        </>
      );
    }

    case "AUCTIONS_LOADING_ERROR": {
      return <pre>{JSON.stringify(observedAuctionsState.error, null, 2)}</pre>;
    }
    default:
      return <></>;
  }
}

export default App;
