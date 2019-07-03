import express from "express"
import cors from "cors"
import { queryChaincode } from "./helper"
import bodyParser = require("body-parser")

async function bootstrap() {
  const app = express()
  app.use(cors())
  app.use(bodyParser.json())

  app.get("/trips", async function(req, res) {
    const response = await queryChaincode(
      "citopia-channel",
      "admin",
      "trip-contract-go",
      "findTrips",
      [],
    )
    const jsonResults: any[] = JSON.parse(response)
    const trips = jsonResults.map(result => {
      return {
        id: result.id,
        userId: result.userId,
        providerId: result.providerId,
        serviceId: result.serviceId,
        mapBitId: result.mapBitId,
        status: result.status,
        currentUserLatitude: result.currentUserLatitude,
        currentUserLongitude: result.currentUserLongitude,
        startTime: parseInt(result.startTime),
        type: result.type,
      }
    })
    res.json(trips)
  })

  app.listen(4000, () => {
    console.log("Server is running at http://localhost:4000")
  })
}

bootstrap().catch(error => console.error(error))
