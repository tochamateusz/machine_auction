import { Box, CircularProgress } from "@mui/material";
//@ts-ignore
import ImageGallery from "react-image-gallery";
import "react-image-gallery/styles/css/image-gallery.css";
import { Auction } from "../App";
import { useEffect, useState } from "react";
import axios from "axios";

const style = {
  position: "absolute" as const,
  top: "50%",
  left: "50%",
  transform: "translate(-50%, -50%)",
  width: 800,
  bgcolor: "background.paper",
  border: "2px solid #000",
  boxShadow: 24,
  p: 4,
};

interface Props {
  children?: React.ReactNode;
  auction: Auction;
  onClick?: (auction: Auction) => void;
}

type GalleryState =
  | { type: "INIT" }
  | { type: "LOADING_GALLERY" }
  | { type: "GALLERY_LOADED"; images: string[] };

export const AuctionGallery: React.FC<Props> = ({ auction }) => {
  const [galleryState, setGalleryState] = useState<GalleryState>({
    type: "INIT",
  });

  const getImages = async (auctionId: string) => {
    setGalleryState({ type: "LOADING_GALLERY" });
    try {
      const images = await axios.get<string[]>(
        `${import.meta.env.VITE_DOMAIN}/scrapper/images/${auctionId}`,
      );
      setGalleryState({ type: "GALLERY_LOADED", images: images?.data || [] });
    } catch (e) {}

    return;
  };

  useEffect(() => {
    switch (galleryState.type) {
      case "INIT": {
        getImages(auction.id);
        return;
      }
      default: {
      }
    }
  }, [galleryState]);

  switch (galleryState.type) {
    case "INIT": {
      return <></>;
    }
    case "LOADING_GALLERY": {
      return (
        <Box sx={style}>
          <CircularProgress />
        </Box>
      );
    }
    case "GALLERY_LOADED": {
      return (
        <Box sx={style}>
          <ImageGallery
            items={galleryState.images.map((fileName: string) => {
              return {
                original: `${import.meta.env.VITE_DOMAIN}/scrapped/${
                  auction.id
                }/${fileName}`,
              };
            })}
          />
        </Box>
      );
    }
  }
};
