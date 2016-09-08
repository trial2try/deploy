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

var GROUPINDEX = "_bigroupindex"													//name for the key/value that will store a list of all known groups
var UNCOVEREDGROUPINDEX = "_uncoveredgroupindex"								//name for the key/value that will store all un covered groups

type Claim struct{
	Claimed 	int 		`json:"claimed"`
	Settled 	int 		`json:"settled"`										//the fieldtags are needed to keep case from bouncing around
	Timestamp 	int64 		`json:"timestamp"`									//utc timestamp of creation
	Type 		string 		`json:"type"`
}

type Risk struct{
	Value 		int 		`json:"value"`										
	Premium 	int 		`json:"premium"`
	Model 		string 		`json:"model"`
	status 		string 		`json:"status"`
	Claims 		[]Claim 	`json:"claims"`
}

type Member struct{
	Name 		string 		`json:"name"`
    Email 		string 		`json:"email"`
    Contact 	string 		`json:"contact"`
    HomeAddress	string 		`json:"address"`
    Dob			string 		`json:"dob"`
    Risks 		[]Risk 		`json:"risks"`
    Wallet 		int 		`json:"wallet"`
}

type Group struct{
	Name 		string 		`json:"name"`
	Members		[]Member 	`json:"members"`	 
	RiskType	string		`json:"riskType"`	
	Status		int 		`json:"status"`
	PoolBalance	int 		`json:"poolBalance"`
	Insurer		string 		`json:"insurer"`
}


// Init resets all the things
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
	
	var empty []string
	empty[0] = "bi001"
	empty[1] = "bi002"
	empty[2] = "bi003"
	empty[3] = "bi004"
	empty[4] = "bi005"

	jsonAsBytes, _ := json.Marshal(empty)								//marshal an emtpy array of strings to clear the index,  now intialising with hard coded values
	err = stub.PutState(GROUPINDEX, jsonAsBytes)
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

	group := Group{}
	member1 := Member{}
	member2 := Member{}
	member3 := Member{}
	risk1 := Risk{}
	risk2 := Risk{}
	risk3 := Risk{}
	claim3 := Claim{}

	
	risk1.Value  = 100									
	risk1.Premium  = 10
	risk1.Model  = 	"XYZ"
	risk1.status  = "Active"	
	
	risk2.Value  = 	200								
	risk2.Premium  = 20
	risk2.Model  = 	"ABC"
	risk2.status  = "Passive"

	claim3.Claimed = 20
	claim3.Settled = 20
	claim3.Timestamp =  makeTimestamp()
	claim3.Type = "Damage"

	risk3.Value  = 250									
	risk3.Premium  = 25
	risk3.Model  = 	"PQR"
	risk3.status  = "Active"
	risk3.Claims = 	append(risk3.Claims, claim3)

	
	member1.Name = "John Snow"
    member1.Email  = "johnsnow.s@gmail.com"	
    member1.Contact  = "+1-541-754-3010"
    member1.HomeAddress	 = "New York"
    member1.Dob		 = "20 Nov 93"
    member1.Wallet = 1000
    member1.Risks = append(member1.Risks, risk1)	
    member1.Risks = append(member1.Risks, risk2)
    
    member2.Name = "Thomas Buck"
    member2.Email  = "thomas@gmail.com"	
    member2.Contact  = "+1-541-754-5987"
    member2.HomeAddress	 = "New Jersey"
    member2.Dob	 = 	"20 Sept 94"
    member2.Wallet = 1450
    member2.Risks = append(member2.Risks, risk2)

    member3.Name = "George Tent"
    member3.Email  = "george@tcs.com"	
    member3.Contact  = "+1-541-754-7811"
    member3.HomeAddress	 = "Seattle"
    member3.Dob	 = 	"20 Sept 94"
    member3.Wallet = 1450
    member3.Risks = append(member3.Risks, risk2)
    member3.Risks = append(member3.Risks, risk3)

    
    group.Members = append(group.Members, member1)
    group.Members = append(group.Members, member2)
    group.Members = append(group.Members, member3)
	group.RiskType = "Bicycle"
	group.Status = 1
	group.PoolBalance = 1000
	group.Insurer = "JKL Insurance Company"
	group.Name = "bi001"
    


	jsonAsBytes, _ = json.Marshal(group)

	err = stub.PutState(group.Name, jsonAsBytes)				
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Invoke isur entry point to invoke a chaincode function
//func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {							//		NEED TO ADD
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

// Query is our entry point for queries
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

// write - invoke function to write key/value pair
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

// read - query function to read key/value pair
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
// Delete - remove a key/value pair from state
// ============================================================================================================================
//func (t *SimpleChaincode) Delete(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {												//NEED TO ADD FEW MORE STATE DELETE
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
// ============================================================================================================================
// Make Timestamp - create a timestamp in ms
// ============================================================================================================================
func makeTimestamp() int64 {
    return time.Now().UnixNano() / (int64(time.Millisecond)/int64(time.Nanosecond))
}

