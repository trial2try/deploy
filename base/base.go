/*
Copyright IBM Corp 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"errors"
	"fmt"
	//"strconv"
	"encoding/json"
	"time"
	//"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}


/*	CONSTANTS	*/

var BICYCLE_GROUP_INDEX = "_bicycleindex"													//name for the key/value that will store a list of all bicycle groups
var SMARTPHONE_GROUP_INDEX = "_smartphoneindex"
var ID_CARD_GROUPINDEX = "_idcardindex"
var BICYCLE = "Bicycle"
var SMARTPHONE = "Smartphone"
var IDCARD = "Idcard"
/*	STRUCTURES PERSISTING IN BLOCKCHAIN	*/

type Claim struct{
	Id 			string 		`json:"id"`
	RiskId 		string 		`json:"riskid"`
	Claimed 	int 		`json:"claimed"`
	Settled 	int 		`json:"settled"`												//the fieldtags are needed to keep case from bouncing around
	Timestamp 	int64 		`json:"timestamp"`												//utc timestamp of creation
	Type 		string 		`json:"type"`
}

type Risk struct{
	Id 			string 		`json:"id"`
	Value 		float64 		`json:"value"`										
	Premium 	float64		`json:"premium"`
	Model 		string 		`json:"model"`
	Type 		string 		`json:"type"`
	Status 		string 		`json:"status"`
	OwnerId 	string 		`json:"ownerid"`
	ClaimIds	[]string 	`json:"claimsids"`
}

type Member struct{
	UserId 		string 		`json:userid`
	Name 		string 		`json:"name"`
    Email 		string 		`json:"email"`
    Contact 	string 		`json:"contact"`
    HomeAddress	string 		`json:"address"`
    Dob			string 		`json:"dob"`
    RiskIds 	[]string 	`json:"riskids"`
    Tokens 		float64 		`json:tokens`
}

type Group struct{
	Name 		string 		`json:"name"`
	RiskIds		[]string 	`json:"riskids"`	 
	RiskType	string		`json:"riskType"`	
	Status		string 		`json:"status"`
	PoolBalance	float64		`json:"poolBalance"`
	InsurerId	string 		`json:"insurer"`
	CreatedDate int64 		`json:"timestamp"`
}

type Insurer struct{
	Id			string 		`json:"id"`
	Name 		string 		`json:"name"`
	Tokens 		float64		`json:tokens`

}

/*	RESPONSE STRUCTURES	*/
type GroupInfo struct{
	Name 		string 		`json:"name"`
	Count 		int			`json:"count"`
	RiskType	string		`json:"riskType"`	
	Status		string 		`json:"status"`
	PoolBalance	int 		`json:"poolBalance"`

}

/*	PREDEFINED FUNCTIONS	*/

// ============================================================================================================================
// Init - Initializes all key value pairs to inital values
// ============================================================================================================================
//func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
func (t *SimpleChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	// Write the state to the ledger
	err := stub.PutState("hello", []byte(args[0]))				//making a test var "hello", I find it handy to read/write to it right away to test the network
	if err != nil {
		return nil, err
	}
	
	var bicycleInit []string
	bicycleInit = append(bicycleInit, "BI001")
	bicycleInit = append(bicycleInit, "BI002")
	bicycleInit = append(bicycleInit, "BI003")
	bicycleInit = append(bicycleInit, "BI004")

	jsonAsBytes, _ := json.Marshal(bicycleInit)								//marshal an emtpy array of strings to clear the index,  now intialising with hard coded values
	err = stub.PutState(BICYCLE_GROUP_INDEX, jsonAsBytes)
	if err != nil {
		return nil, err
	}
	/*
	var trades AllTrades
	jsonAsBytes, _ = json.Marshal(trades)								//clear the open trade struct
	err = stub.PutState(openTradesStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}
	*/

	group1 := Group{}
	group2 := Group{}
	group3 := Group{}
	group4 := Group{}
	member1 := Member{}
	member2 := Member{}
	member3 := Member{}
	risk1 := Risk{}
	risk2 := Risk{}
	risk3 := Risk{}
	claim1 := Claim{}
	claim2 := Claim{}
	insurer := Insurer{}

	claim1.Id = "CID001"
	claim1.RiskId = "RID002"
	claim1.Claimed = 40
	claim1.Timestamp = makeTimestamp()
	claim1.Type = "Damage"
    
	claim2.Id = "CID002"
	claim2.RiskId = "RID003"
	claim2.Claimed = 200
	claim2.Settled = 200
	claim2.Timestamp = makeTimestamp()
	claim2.Type = "Theft"

	risk1.Id = "RID001"
	risk1.Value = 100
	risk1.Premium = 12
	risk1.Model = "XYZ"
	risk1.Type = BICYCLE
	risk1.Status = "Active"
	risk1.OwnerId = "UID002"

	risk2.Id = "RID002"
	risk2.Value = 250
	risk2.Premium = 25
	risk2.Model = "PQR"
	risk2.Type = BICYCLE
	risk2.Status = "Active"
	risk2.OwnerId = "UID003"
	risk2.ClaimIds = append(risk2.ClaimIds, claim1.Id)

	risk3.Id = "RID003"
	risk3.Value = 200
	risk3.Premium = 20
	risk3.Model = "ABC"
	risk3.Type = BICYCLE
	risk3.Status = "Active"
	risk3.OwnerId = "UID001"
	risk3.ClaimIds = append(risk3.ClaimIds, claim2.Id)

	member1.UserId = "UID001"
	member1.Name = "John Snow"
	member1.Email = "johnsnow.s@gmail.com"	
	member1.Contact = "+1-541-754-3010"
	member1.HomeAddress = "New York"
	member1.Dob = "20 Nov 93"
	member1.RiskIds = append(member1.RiskIds, risk3.Id)
	member1.Tokens = 1000

	member2.UserId = "UID002"
	member2.Name = "Thomas Buck"
	member2.Email = "thomas@gmail.com"	
	member2.Contact = "+1-541-754-5987"
	member2.HomeAddress = "New Jersey"
	member2.Dob = "20 Sept 94"
	member2.RiskIds = append(member2.RiskIds, risk1.Id)
	member2.Tokens = 1200

	member3.UserId = "UID003"
	member3.Name = "George Tent"
	member3.Email = "george@gmail.com"
	member3.Contact = "+1-541-754-7811"
	member3.HomeAddress = "Seattle"
	member3.Dob = "20 Sept 94"
	member3.RiskIds = append(member3.RiskIds, risk2.Id)
	member3.Tokens = 1450

	insurer.Id = "INS001"
	insurer.Name = "JKL Insurance"
	insurer.Tokens = 13840
	

	group1.Name = "BI001"
	group1.RiskIds = append(group1.RiskIds, risk1.Id)
	group1.RiskIds = append(group1.RiskIds, risk2.Id)
	group1.RiskIds = append(group1.RiskIds, risk3.Id)
	group1.RiskType = BICYCLE
	group1.Status = "Open"
	group1.PoolBalance = 39.5
	group1.InsurerId = "INS001"
	group1.CreatedDate =  makeTimestamp()

	group2.Name = "BI002"
	for j :=0; j < 5; j++ {
		group2.RiskIds = append(group2.RiskIds, "dummyRisk")		
	}
	group2.RiskType = BICYCLE
	group2.Status = "Open"
	group2.PoolBalance = 45.6
	group2.InsurerId = "INS005"
	group2.CreatedDate = makeTimestamp() + 60000000

	group3.Name = "BI003"
	for j :=0; j < 8; j++ {
		group3.RiskIds = append(group3.RiskIds, "dummyRisk")		
	}
	group3.RiskType = BICYCLE
	group3.Status = "Closed"
	group3.PoolBalance = 67.2
	group3.InsurerId = "INS003"
	group3.CreatedDate = makeTimestamp() + 120000000

	group4.Name = "BI004"
	for j :=0; j < 8; j++ {
		group4.RiskIds = append(group4.RiskIds, "dummyRisk")		
	}
	group4.RiskType = BICYCLE
	group4.Status = "Closed"
	group4.PoolBalance = 23.7
	group4.InsurerId = "INS002"
	group4.CreatedDate = makeTimestamp() + 90000000
	    

	/*Persisting Groups*/
	jsonAsBytes, _ = json.Marshal(group1)
	err = stub.PutState(group1.Name, jsonAsBytes)				
	if err != nil {
		return nil, err
	}
	jsonAsBytes, _ = json.Marshal(group2)
	err = stub.PutState(group2.Name, jsonAsBytes)				
	if err != nil {
		return nil, err
	}
	jsonAsBytes, _ = json.Marshal(group3)
	err = stub.PutState(group3.Name, jsonAsBytes)				
	if err != nil {
		return nil, err
	}
	jsonAsBytes, _ = json.Marshal(group4)
	err = stub.PutState(group4.Name, jsonAsBytes)				
	if err != nil {
		return nil, err
	}

	/*Persisting Risks*/
	jsonAsBytes, _ = json.Marshal(risk1)
	err = stub.PutState(risk1.Id, jsonAsBytes)				
	if err != nil {
		return nil, err
	}
	jsonAsBytes, _ = json.Marshal(risk2)
	err = stub.PutState(risk2.Id, jsonAsBytes)				
	if err != nil {
		return nil, err
	}
	jsonAsBytes, _ = json.Marshal(risk3)
	err = stub.PutState(risk3.Id, jsonAsBytes)				
	if err != nil {
		return nil, err
	}

	/*Persisting Members */
	jsonAsBytes, _ = json.Marshal(member1)
	err = stub.PutState(member1.UserId, jsonAsBytes)				
	if err != nil {
		return nil, err
	}
	jsonAsBytes, _ = json.Marshal(member2)
	err = stub.PutState(member2.UserId, jsonAsBytes)				
	if err != nil {
		return nil, err
	}
	jsonAsBytes, _ = json.Marshal(member3)
	err = stub.PutState(member3.UserId, jsonAsBytes)				
	if err != nil {
		return nil, err
	}

	/*Persisting Claims*/
	jsonAsBytes, _ = json.Marshal(claim2)
	err = stub.PutState(claim2.Id, jsonAsBytes)				
	if err != nil {
		return nil, err
	}
	jsonAsBytes, _ = json.Marshal(claim1)
	err = stub.PutState(claim1.Id, jsonAsBytes)				
	if err != nil {
		return nil, err
	}

	/*Persisting Insurer*/
	jsonAsBytes, _ = json.Marshal(insurer)
	err = stub.PutState(insurer.Id, jsonAsBytes)				
	if err != nil {
		return nil, err
	}
	return nil, nil
}
// ============================================================================================================================
// Invoke - The entry point to invoke a chaincode function	
// ============================================================================================================================
//func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {	
func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "write" {
		return t.write(stub, args)
	} else if function == "delete" {
		return t.Delete(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation")
}
// ============================================================================================================================
// Query - The entry point for queries to a chaincode
// ============================================================================================================================
//func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	} 
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query")
}
// ============================================================================================================================
// write - Invoke function to write key/value pair
// ============================================================================================================================
//func (t *SimpleChaincode) write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
func (t *SimpleChaincode) write(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var key, value string
	var err error
	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}

	key = args[0] //rename for funsies
	value = args[1]
	err = stub.PutState(key, []byte(value)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}
// ============================================================================================================================
// read - Query function to read key/value pair
// ============================================================================================================================
//func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
func (t *SimpleChaincode) read(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}
// ============================================================================================================================
// Delete - Invoke function to remove a key/value pair from state
// ============================================================================================================================
//func (t *SimpleChaincode) Delete(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {												
func (t *SimpleChaincode) Delete(stub *shim.ChaincodeStub, args []string) ([]byte, error) {	
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	
	name := args[0]
	err := stub.DelState(name)													//remove the key from chaincode state
	if err != nil {
		return nil, errors.New("Failed to delete state")
	}

	return nil, nil
}


/*	USER DEFINED FUNCTIONS	*/


/*	UTILITY FUNCTIONS	*/

// ============================================================================================================================
// Make Timestamp - create a timestamp in ms
// ============================================================================================================================
func makeTimestamp() int64 {
    return time.Now().UnixNano() / (int64(time.Millisecond)/int64(time.Nanosecond))
}

