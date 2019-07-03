import FabricCAServices, { TLSOptions } from "fabric-ca-client"
import { FileSystemWallet, X509WalletMixin } from "fabric-network"
import * as config from "../config/connection.json"
import * as path from "path"
import * as fs from "fs"

const appRootPath = require("app-root-path").path

async function main() {
  try {
    // Parse connection config
    const registrar: any =
      config["certificateAuthorities"]["ca-org1"]["registrar"]
    const enrollId: string = registrar[0]["enrollId"]
    const enrollSecret: string = registrar[0]["enrollSecret"]
    const url: string = config["certificateAuthorities"]["ca-org1"]["url"]
    const clearUrl = url.replace("https://", "")
    const caURL = `https://${enrollId}:${enrollSecret}@${clearUrl}`
    const mspid = config["certificateAuthorities"]["ca-org1"]["caName"]

    const pem = Buffer.from(readAllFiles("./pem")[0])
    const options: TLSOptions = { trustedRoots: pem, verify: true }
    const caServices = new FabricCAServices(caURL, options)

    const walletPath = path.join(appRootPath, "wallet")
    const wallet = new FileSystemWallet(walletPath)
    console.log(`Wallet path: ${walletPath}`)

    const adminExists = await wallet.exists(enrollId)
    if (adminExists) {
      console.log(
        `An identity for the admin user "admin" already exists in the wallet`,
      )
      return
    }

    // Enroll the admin user, and import the new identity into the wallet.
    const enrollment = await caServices.enroll({
      enrollmentID: enrollId,
      enrollmentSecret: enrollSecret,
    })
    console.log("enrollment:", enrollment)

    const identity = X509WalletMixin.createIdentity(
      mspid,
      enrollment.certificate,
      enrollment.key.toBytes(),
    )

    console.log("identity:", identity)
    await wallet.import(enrollId, identity)

    console.log(
      `Successfully enrolled admin user "admin" and imported it into the wallet`,
    )
  } catch (error) {
    console.error(`Failed to enroll admin user "admin": ${error}`)
    process.exit(1)
  }
}

function readAllFiles(dir: string) {
  const files = fs.readdirSync(dir)
  const certs: any = []
  files.forEach(fileName => {
    const filePath = path.join(dir, fileName)
    const data = fs.readFileSync(filePath)
    certs.push(data)
  })
  return certs
}

main()
