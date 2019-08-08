import express from "express"
import cors from "cors"
import { queryChaincode } from "./helper"
import bodyParser = require("body-parser")

async function bootstrap() {
  const app = express()
  app.use(cors())
  app.use(bodyParser.json())

  app.get("/trips", async function(req, res) {
    const filter = {
      providerId: process.env.PROVIDER_ID,
      completed: 0,
      paid: 0,
    }

    const response = await queryChaincode(
      "citopia-channel",
      "admin",
      "trip-contract",
      "findTrips",
      [JSON.stringify(filter)],
    )
    const jsonResults: any[] = JSON.parse(response)
    const trips = jsonResults.map(result => {
      return {
        id: result.id,
        userId: result.userId,
        providerId: result.providerId,
        mapBitId: result.mapBitId,
        serviceId: result.serviceId,
        serviceType: result.serviceType,
        serviceVehicleType: result.serviceVehicleType,
        completed: result.completed,
        paid: result.paid,
        currentLat: result.currentLat,
        currentLng: result.currentLng,
        destinationLat: result.destinationLat,
        destinationLng: result.destinationLng,
        co2: result.co2,
        traffic: result.traffic,
        health: result.health,
        startTime: result.startTime,
        endTime: result.endTime
      }
    })
    res.json(trips)
  })

  app.listen(4000, () => {
    console.log("Server is running at http://localhost:4000")
  })
}

bootstrap().catch(error => console.error(error))
