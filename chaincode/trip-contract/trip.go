package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
	"time"
)

type SmartContract struct {
}

type Trip struct {
	Id                 string  `json:"id"`
	UserId             string  `json:"userId"`
	ProviderId         string  `json:"providerId"`
	MapBitId           string  `json:"mapBitId"`
	ServiceId          string  `json:"serviceId"`
	ServiceType        string  `json:"serviceType"`
	ServiceVehicleType string  `json:"serviceVehicleType"`
	Completed          bool    `json:"completed"`
	Paid               bool    `json:"paid"`
	CurrentLat         string  `json:"currentLat"`
	CurrentLng         string  `json:"currentLng"`
	DestinationLat     string  `json:"destinationLat"`
	DestinationLng     string  `json:"destinationLng"`
	Co2                float64 `json:"co2"`
	Traffic            float64 `json:"traffic"`
	Health             float64 `json:"health"`
	StartTime          int     `json:"startTime"`
	EndTime            int     `json:"endTime"`
	CreateDate         int64   `json:"createDate"`
	Type               string  `json:"type"`
}

type TripFilter struct {
	UserId       string   `json:"userId"`
	ProviderId   string   `json:"providerId"`
	ServiceId    string   `json:"serviceId"`
	ServiceTypes []string `json:"serviceTypes"`
	Completed    int      `json:"completed"`
	Paid         int      `json:"paid"`
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
	fmt.Println("=============== INVOKE ===============")
	fmt.Println("Function:", function)
	fmt.Println("Args:", args)
	fmt.Println("Datetime:", time.Now().Format("01-02-2006 15:04:05"))
	fmt.Println("======================================")

	switch function {
	case "findTrip":
		return s.findTrip(stub, args)
	case "findTrips":
		return s.findTrips(stub, args)
	default:
		return shim.Error("Invalid Smart Contract function name.")
	}
}

func (s *SmartContract) findTrip(stub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	key, compositeErr := stub.CreateCompositeKey("record", []string{args[0]})
	if compositeErr != nil {
		return shim.Error(fmt.Sprintf("Could not create a composite key for %s: %s", args[0], compositeErr.Error()))
	}
	recordAsBytes, _ := stub.GetState(key)
	return shim.Success(recordAsBytes)
}

func (s *SmartContract) findTrips(stub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	// Parse TripFilter from JSON
	tripFilterAsJson := args[0]
	tripFilterAsBytes := []byte(tripFilterAsJson)
	tripFilter := TripFilter{}
	err := json.Unmarshal(tripFilterAsBytes, &tripFilter)
	if err != nil {
		return shim.Error("Could not unmarshal given TripFilter")
	}

	resultsIterator, err := stub.GetStateByPartialCompositeKey("record", []string{})
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

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

		if len(tripFilter.UserId) > 0 && trip.UserId != tripFilter.UserId {
			continue
		}
		if len(tripFilter.ProviderId) > 0 && trip.ProviderId != tripFilter.ProviderId {
			continue
		}
		if len(tripFilter.ServiceId) > 0 && trip.ServiceId != tripFilter.ServiceId {
			continue
		}
		if len(tripFilter.ServiceTypes) > 0 && indexOf(tripFilter.ServiceTypes, trip.ServiceType) == -1 {
			continue
		}

		completed := false
		if tripFilter.Completed == 1 {
			completed = true
		}
		if tripFilter.Completed != 0 && trip.Completed != completed {
			continue
		}

		paid := false
		if tripFilter.Paid == 1 {
			paid = true
		}
		if tripFilter.Paid != 0 && trip.Paid != paid {
			continue
		}

		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		// Trip is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return shim.Success(buffer.Bytes())
}

func indexOf(arr []string, element string) int {
	for k, v := range arr {
		if v == element {
			return k
		}
	}
	return -1
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
