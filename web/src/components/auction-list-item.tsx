import {
  Box,
  Card,
  CardActionArea,
  CardContent,
  CardMedia,
  Grid,
  Typography,
} from "@mui/material";
import React from "react";
import { Auction } from "../App";

interface Props {
  children?: React.ReactNode;
  auction: Auction;
  onClick?: (auction: Auction) => void;
}

export const ListAuction: React.FC<Props> = ({ auction, onClick }) => {
  return (
    <Card sx={{ maxWidth: 600, width: 600, marginY: "1rem", marginX: "1rem" }}>
      <CardActionArea
        onClick={() => {
          onClick && onClick(auction);
        }}
      >
        <Grid container spacing={2}>
          <Grid item xs={5}>
            <CardMedia
              component="img"
              height="250"
              image={`${import.meta.env.VITE_DOMAIN}/scrapped/${
                auction.id
              }/0.jpg`}
              alt="machine"
            />
          </Grid>
          <Grid item xs={7}>
            <CardContent>
              <Typography gutterBottom variant="h5" component="div">
                {auction.name}
              </Typography>
              <Box sx={{ display: "flex", marginX: "0.5rem" }}>
                <Typography variant="body2" color="text.secondary">
                  Price:
                </Typography>
                <Typography
                  variant="body2"
                  color="text.primary"
                  sx={{ marginX: "1rem" }}
                >
                  {auction.price}
                </Typography>
              </Box>

              <Box sx={{ display: "flex", marginX: "0.5rem" }}>
                <Typography variant="body2" color="text.secondary">
                  Year:
                </Typography>
                <Typography
                  variant="body2"
                  color="text.primary"
                  sx={{ marginX: "1rem" }}
                >
                  {auction.year}
                </Typography>
              </Box>

              <Box sx={{ display: "flex", marginX: "0.5rem" }}>
                <Typography variant="body2" color="text.secondary">
                  Auction finished at:
                </Typography>
                <Typography
                  variant="body2"
                  color="text.primary"
                  sx={{ marginX: "1rem" }}
                >
                  {auction.end_date}
                </Typography>
              </Box>
              <Box sx={{ display: "flex", marginX: "0.5rem" }}>
                <Typography variant="body2" color="text.secondary">
                  <>
                  {(auction.description || []).map((description)=>{
                    return <div>{description}</div>
                  })}
                  </>
                </Typography>
              </Box>
            </CardContent>
          </Grid>
        </Grid>
      </CardActionArea>
    </Card>
  );
};
