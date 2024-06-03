import axios from "axios";
import { Auction } from "./App";
import { useEffect, useState } from "react";

export type States =
  | { type: "INIT" }
  | { type: "LOADING_AUCTIONS" }
  | { type: "AUCTIONS_LOADED"; auctions: Auction[], onOpen: (a: Auction) => () => void }
  | { type: "AUCTIONS_LOADING_ERROR"; error: any }
  | { type: "MODAL_OPEN"; auctions: Auction[]; selectedAuction: Auction, onClose: () => void };

export const useObservedAuction = () => {
  const [observedAuctionsState, setObservedAuctionsState] = useState<States>({
    type: "INIT",
  });


  const getAuctions = async () => {
    setObservedAuctionsState({ type: "LOADING_AUCTIONS" });
    try {
      const auctions = await axios.get<Auction[]>(
        `${import.meta.env.VITE_DOMAIN}/scrapper`,
      );

      const onClose = () => {
        setObservedAuctionsState({
          type: "AUCTIONS_LOADED",
          auctions: auctions.data,
          onOpen: onOpen,
        });
      }

      const onOpen = (a: Auction) => () => {
        setObservedAuctionsState({
          type: "MODAL_OPEN",
          auctions: auctions.data,
          selectedAuction: a,
          onClose: onClose,
        });
      }

      setObservedAuctionsState({
        type: "AUCTIONS_LOADED",
        auctions: auctions.data,
        onOpen: onOpen
      })

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

  return observedAuctionsState
}
