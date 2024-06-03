import "./App.css";
import { Box, CircularProgress, Modal, Paper } from "@mui/material";
import { AuctionGallery } from "./components/auction-gallery";
import React from "react";
import { States } from "./useObserverAuction";

export type Auction = {
  id: string;
  image: string;
  name: string;
  year: string;
  price: string;
  end_date: string;
  description: string[];
  starting_price: string;
};

interface Props {
  auction: Auction;
  onClick?: (auction: Auction) => void;
}

type View = (p: Props) => React.ReactNode


const App = ({ view, observedAuctionsState }: { view: View, observedAuctionsState: States }) => {

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
        <Box sx={{ display: "flex", marginY: "2rem", flexWrap: "wrap" }}>
          {observedAuctionsState.auctions.map((a) => {
            return (
              view({ auction: a, onClick: observedAuctionsState.onOpen(a) })
            );
          })}
        </Box>
      );
    }

    case "MODAL_OPEN": {
      return (
        <>
          <Box sx={{ display: "flex", marginY: "2rem", flexWrap: "wrap" }}>
            {observedAuctionsState.auctions.map((a) => {
              return (view({ auction: a }))
            })}
          </Box>
          <Modal
            open={true}
            onClose={observedAuctionsState.onClose}
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
