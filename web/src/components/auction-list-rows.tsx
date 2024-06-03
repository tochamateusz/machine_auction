import {
  Box,
  Card,
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


export const ListRowsAuction: React.FC<Props> = ({ auction }) => {

  return (
    <Card key={auction.id} sx={{ marginY: "1rem", padding: "1rem", marginX: "1rem" }}>
      <Grid container spacing={4} width="120rem">

        <Grid item xs={2}>
          <CardMedia
            component="img"
            image={`${import.meta.env.VITE_DOMAIN}/scrapped/${auction.id
              }/0.jpg`}
            alt="machine"
          />
        </Grid>
        <Grid item xs={1}>
          {auction.id}
        </Grid>
        <Grid item xs={2}>
          {auction.name}
        </Grid>
        <Grid item xs={1}>
          {auction.starting_price}
        </Grid>
        <Grid item xs={1}>
          {auction.price}
        </Grid>
        <Grid item xs={1}>
          {auction.year}
        </Grid>
        <Grid item xs={1}>
          {auction.end_date}
        </Grid>
        <Grid item xs={3}>
          <Box sx={{ marginX: "0.5rem" }}>
            {(auction.description || []).map((description, index) => {
              return (<Typography key={index} display={"block"} variant="body2" color="text.secondary">
                {description}
              </Typography>)
            })}
          </Box>
        </Grid>
      </Grid>
    </Card>
  )
}
