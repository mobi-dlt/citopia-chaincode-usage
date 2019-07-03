package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

type SmartContract struct {
}

type Trip struct {
	Id                   string `json:"id"`
	UserId               string `json:"userId"`
	ProviderId           string `json:"providerId"`
	ServiceId            string `json:"serviceId"`
	MapBitId             string `json:"mapBitId"`
	Status               string `json:"status"`
	CurrentUserLatitude  string `json:"currentUserLatitude"`
	CurrentUserLongitude string `json:"currentUserLongitude"`
	StartTime            string `json:"startTime"`
	Type                 string `json:"type"`
}

/*
 * The Init method is called when the Smart Contract is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(stub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(stub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("Invoke function:", function)
	fmt.Println("Invoke args:", args)

	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "findTrip" {
		return s.findTrip(stub, args)
	} else if function == "findTrips" {
		return s.findTrips(stub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

/*
 * Find trips by parameters
 *  args[0] - user id
 *  args[1] - provider id
 *  args[2] - serviceId id
 *  args[3] - status - "initiated"|"waiting"|"in-progress"|"canceled"|"completed-by-provider"|"completed"
 *
 * Returns all trips if parameters are not specified
 */
func (s *SmartContract) findTrips(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	userId := args[0]
	providerId := args[1]
	serviceId := args[2]
	status := args[3]

	// load all trips
	resultsIterator, err := stub.GetStateByPartialCompositeKey("trip", []string{})
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	// iterate over all trips and check parameters matching
	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		trip := Trip{}
		marshErr := json.Unmarshal(queryResponse.Value, &trip)
		if marshErr != nil {
			fmt.Printf("Error: %s", marshErr)
		}

		if len(userId) > 0 && trip.UserId != userId {
			continue
		}
		if len(providerId) > 0 && trip.ProviderId != providerId {
			continue
		}
		if len(serviceId) > 0 && trip.ServiceId != serviceId {
			continue
		}
		if len(status) > 0 && trip.Status != status {
			continue
		}

		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return shim.Success(buffer.Bytes())
}

/*
 * Find trip by id
 *  args[0] - trip id
 */
func (s *SmartContract) findTrip(stub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	// build composite key from trip type and id
	key, compositeErr := stub.CreateCompositeKey("trip", []string{args[0]})
	if compositeErr != nil {
		return shim.Error(fmt.Sprintf("Could not create a composite key for %s: %s", args[0], compositeErr.Error()))
	}

	// load trip
	tripAsBytes, _ := stub.GetState(key)
	return shim.Success(tripAsBytes)
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
