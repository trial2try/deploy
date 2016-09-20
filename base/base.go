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
	"strconv"
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

var BICYCLE = "bicycle"
var SMARTPHONE = "smartphone"
var IDCARD = "idcard"

var RISK_INDEX = "_risksindex"
var MEMBER_INDEX = "_memebersindex"
var CLAIM_INDEX = "_claimsindex"
var INSURER_INDEX = "_insurersindex"
var ADMIN_FEE ="_adminfee"

/*	COUNTER VARIABLES	*/

var riskCounter = 3
var memberCounter = 3
var insurerCounter = 1
var claimCounter = 2
var bicycleCounter = 4

/*	STRUCTURES PERSISTING IN BLOCKCHAIN	*/

type Claim struct{
	Id 			string 		`json:"id"`
	RiskId 		string 		`json:"riskid"`
	Claimed 	float64		`json:"claimed"`
	Settled 	float64		`json:"settled"`												//the fieldtags are needed to keep case from bouncing around
	Timestamp 	int64 		`json:"timestamp"`												//utc timestamp of creation
	Type 		string 		`json:"type"`
}

type Risk struct{
	Id 			string 		`json:"id"`
	Value 		float64 	`json:"value"`										
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
    Tokens 		float64 	`json:tokens`
}

type Group struct{
	Name 		 string 	`json:"name"`
	RiskIds		 []string 	`json:"riskids"`	 
	RiskType	 string		`json:"riskType"`	
	Status		 string 	`json:"status"`
	PoolBalance	 float64	`json:"poolBalance"`
	InsurerId	 string 	`json:"insurer"`
	CreatedDate  int64 	`json:"createddate"`
	EndDate 	 int64 	`json:"enddate"`
	GroupPremium float64	`json:"groupPremium"`
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

	err = stub.PutState(ADMIN_FEE, []byte(strconv.FormatFloat(0, 'E', -1, 64)));
		if err != nil {
			return nil, errors.New("Billing account cannot be created")
	}
	
	
	var bicycleInit []string
	bicycleInit = append(bicycleInit, "bi001")
	bicycleInit = append(bicycleInit, "bi002")
	bicycleInit = append(bicycleInit, "bi003")
	bicycleInit = append(bicycleInit, "bi004")

	jsonAsBytes, _ := json.Marshal(bicycleInit)								//marshal an emtpy array of strings to clear the index,  now intialising with hard coded values
	err = stub.PutState(BICYCLE_GROUP_INDEX, jsonAsBytes)
	if err != nil {
		return nil, err
	}

	var riskInit []string
	riskInit = append(riskInit, "rid001")
	riskInit = append(riskInit, "rid002")
	riskInit = append(riskInit, "rid003")

	jsonAsBytes, _ = json.Marshal(riskInit)								//marshal an emtpy array of strings to clear the index,  now intialising with hard coded values
	err = stub.PutState(RISK_INDEX, jsonAsBytes)
	if err != nil {
		return nil, err
	}

	var memberInit []string
	memberInit = append(memberInit, "uid001")
	memberInit = append(memberInit, "uid002")
	memberInit = append(memberInit, "uid003")

	jsonAsBytes, _ = json.Marshal(memberInit)								//marshal an emtpy array of strings to clear the index,  now intialising with hard coded values
	err = stub.PutState(MEMBER_INDEX, jsonAsBytes)
	if err != nil {
		return nil, err
	}

	var claimInit []string
	claimInit = append(claimInit, "cid001")
	claimInit = append(claimInit, "cid002")

	jsonAsBytes, _ = json.Marshal(claimInit)								//marshal an emtpy array of strings to clear the index,  now intialising with hard coded values
	err = stub.PutState(CLAIM_INDEX, jsonAsBytes)
	if err != nil {
		return nil, err
	}

	var insurerInit []string
	insurerInit = append(insurerInit, "ins001")

	jsonAsBytes, _ = json.Marshal(insurerInit)								//marshal an emtpy array of strings to clear the index,  now intialising with hard coded values
	err = stub.PutState(INSURER_INDEX, jsonAsBytes)
	if err != nil {
		return nil, err
	}


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

	claim1.Id = "cid001"
	claim1.RiskId = "rid002"
	claim1.Claimed = 40.0
	claim1.Timestamp = makeTimestamp()
	claim1.Type = "damage"
    
	claim2.Id = "cid002"
	claim2.RiskId = "rid003"
	claim2.Claimed = 200.0
	claim2.Settled = 200.0
	claim2.Timestamp = makeTimestamp()
	claim2.Type = "theft"

	risk1.Id = "rid001"
	risk1.Value = 100
	risk1.Premium = 12
	risk1.Model = "XYZ"
	risk1.Type = BICYCLE
	risk1.Status = "covered"
	risk1.OwnerId = "uid002"

	risk2.Id = "rid002"
	risk2.Value = 250
	risk2.Premium = 25
	risk2.Model = "PQR"
	risk2.Type = BICYCLE
	risk2.Status = "covered"
	risk2.OwnerId = "uid003"
	risk2.ClaimIds = append(risk2.ClaimIds, claim1.Id)

	risk3.Id = "rid003"
	risk3.Value = 200
	risk3.Premium = 20
	risk3.Model = "ABC"
	risk3.Type = BICYCLE
	risk3.Status = "covered"
	risk3.OwnerId = "uid001"
	risk3.ClaimIds = append(risk3.ClaimIds, claim2.Id)

	member1.UserId = "uid001"
	member1.Name = "John Snow"
	member1.Email = "johnsnow.s@gmail.com"	
	member1.Contact = "+1-541-754-3010"
	member1.HomeAddress = "New York"
	member1.Dob = "20 Nov 93"
	member1.RiskIds = append(member1.RiskIds, risk3.Id)
	member1.Tokens = 1000

	member2.UserId = "uid002"
	member2.Name = "Thomas Buck"
	member2.Email = "thomas@gmail.com"	
	member2.Contact = "+1-541-754-5987"
	member2.HomeAddress = "New Jersey"
	member2.Dob = "20 Sept 94"
	member2.RiskIds = append(member2.RiskIds, risk1.Id)
	member2.Tokens = 1200

	member3.UserId = "uid003"
	member3.Name = "George Tent"
	member3.Email = "george@gmail.com"
	member3.Contact = "+1-541-754-7811"
	member3.HomeAddress = "Seattle"
	member3.Dob = "20 Sept 94"
	member3.RiskIds = append(member3.RiskIds, risk2.Id)
	member3.Tokens = 1450

	insurer.Id = "ins001"
	insurer.Name = "JKL Insurance"
	insurer.Tokens = 13840
	

	group1.Name = "bi001"
	group1.RiskIds = append(group1.RiskIds, risk1.Id)
	group1.RiskIds = append(group1.RiskIds, risk2.Id)
	group1.RiskIds = append(group1.RiskIds, risk3.Id)
	group1.RiskType = BICYCLE
	group1.Status = "open"
	group1.PoolBalance = 39.5
	group1.InsurerId = "ins001"
	group1.CreatedDate =  makeTimestamp()
	group1.EndDate = 1506752393
	group1.GroupPremium = 25

	group2.Name = "bi002"
	for j :=0; j < 5; j++ {
		group2.RiskIds = append(group2.RiskIds, "dummyRisk")		
	}
	group2.RiskType = BICYCLE
	group2.Status = "open"
	group2.PoolBalance = 45.6
	group2.InsurerId = "ins005"
	group2.CreatedDate = makeTimestamp() - 60000000
	group2.EndDate = 1506752393 - 60000000
	group2.GroupPremium = 30

	group3.Name = "bi003"
	for j :=0; j < 8; j++ {
		group3.RiskIds = append(group3.RiskIds, "dummyRisk")		
	}
	group3.RiskType = BICYCLE
	group3.Status = "closed"
	group3.PoolBalance = 67.2
	group3.InsurerId = "ins003"
	group3.CreatedDate = makeTimestamp() - 120000000
	group3.EndDate = 1506752393 - 120000000
	group3.GroupPremium = 24;

	group4.Name = "bi004"
	for j :=0; j < 8; j++ {
		group4.RiskIds = append(group4.RiskIds, "dummyRisk")		
	}
	group4.RiskType = BICYCLE
	group4.Status = "closed"
	group4.PoolBalance = 23.7
	group4.InsurerId = "ins002"
	group4.CreatedDate = makeTimestamp() - 90000000
	group4.EndDate = 1506752393 - 90000000   
	group4.GroupPremium = 28;

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

	/*	Initialize counter variable	*/
	riskCounter = 3
	memberCounter = 3
	insurerCounter = 1
	claimCounter = 2
	bicycleCounter = 4

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
	} else if function == "create_risk" {
		return t.CreateRisk(stub, args)
	} else if function == "add_risk" {
		return t.AddRisk(stub, args)
	} else if function == "raise_claim" {
		return t.RaiseClaim(stub, args)
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
	} else if function == "getGroupInfo"{

	} else if function == "getGroupMembers"{

	} else if function == "getGroupRisks"{

	} else if function == ""{

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
	fmt.Println("Read running for ",args[0])
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
// ============================================================================================================================
/* 
CreateRisk - Invoke function to create a new risk write key/value pair
Inputs: 	args[0]		args[1]	args[2]		args[3]
			value 		model 	type 		owner
			"100" 		"XYZ" 	"bicycle" 	
*/
// ============================================================================================================================
//func (t *SimpleChaincode) CreateRisk(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
func (t *SimpleChaincode) CreateRisk(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	fmt.Println("running CreateRisk()")

	if len(args) != 4 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4. ")
	}

	value, err := strconv.ParseFloat(args[0], 64)
	if err != nil {
		return nil, errors.New("1st argument must be a numeric string")
	}

	risk := Risk{}
	risk.Id = makeRiskId()
	risk.Value = value
	//risk.Premium = 12 											// Should be given at the period of joining a group
	risk.Model = args[1]
	risk.Type = args[2]
	risk.Status = "uncovered" 										// Should be changed active when joining a group
	risk.OwnerId = args[3]

	jsonAsBytes, _ := json.Marshal(risk)
	err = stub.PutState(risk.Id, jsonAsBytes)				
	if err != nil {
		return nil, err
	}

	riskIndexAsBytes, err := stub.GetState(RISK_INDEX)
	if err != nil {
		return nil, errors.New("Failed to get risk index")
	}

	var riskIndex []string
	json.Unmarshal(riskIndexAsBytes, &riskIndex)								//un stringify it aka JSON.parse()
	
	//append
	riskIndex = append(riskIndex, risk.Id)										//add risk id to index list
	fmt.Println("! risk index: ", riskIndex)
	jsonAsBytes, _ = json.Marshal(riskIndex)
	err = stub.PutState(RISK_INDEX, jsonAsBytes)						//store risk id of risk

	//Get Member
	memberAsBytes, err := stub.GetState(risk.OwnerId)
	if err != nil {
		return nil, errors.New("Failed to get member")
	}
	member := Member{}
	json.Unmarshal(memberAsBytes, &member)

	member.RiskIds = append(member.RiskIds, risk.Id)
	jsonAsBytes, _ = json.Marshal(member)
	err = stub.PutState(member.UserId, jsonAsBytes)				
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// ============================================================================================================================
/* 
AddRisk - Invoke function to create a new risk write key/value pair
Inputs: 	args[0]		args[1]		
			riskid 		groupid 	 	
			"rid002"	"bi002"			
*/
// ============================================================================================================================
//func (t *SimpleChaincode) AddRisk(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
func (t *SimpleChaincode) AddRisk(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	fmt.Println("running AddRisk()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2.")
	}

	/*premium, err := strconv.ParseFloat(args[2], 64)
	if err != nil {
		return nil, errors.New("3rd argument must be a numeric string")
	}*/

	//Get Risk
	riskAsBytes, err := stub.GetState(args[0])
	if err != nil {
		return nil, errors.New("Failed to get risk")
	}
	risk := Risk{}
	json.Unmarshal(riskAsBytes, &risk)								//un stringify it aka JSON.parse()



	//Get Member
	memberAsBytes, err := stub.GetState(risk.OwnerId)
	if err != nil {
		return nil, errors.New("Failed to get member")
	}
	member := Member{}
	json.Unmarshal(memberAsBytes, &member)


	//Get Group
		groupAsBytes, err := stub.GetState(args[1])
		if err != nil {
			return nil, errors.New("Failed to get group")
		}
		group := Group{}
		json.Unmarshal(groupAsBytes, &group)

	if member.Tokens-group.GroupPremium > 0.0 {

		if group.RiskType != risk.Type{
			return nil, errors.New("Can't add "+risk.Type+" risk to group of "+group.RiskType)	
		}
	
		//Get Insurer
		insurerAsBytes, err := stub.GetState(group.InsurerId)
		if err != nil {
			return nil, errors.New("Failed to get group")
		}
		insurer := Insurer{}
		json.Unmarshal(insurerAsBytes, &insurer)
		//Get Admin Account Fee
		value, err := stub.GetState(ADMIN_FEE)
		if err != nil {
			jsonResp := "{\"Error\":\"Failed to get state for " + ADMIN_FEE + "\"}"
			return nil, errors.New(jsonResp)
		}
		previous_val, _ := strconv.ParseFloat(string(value), 64)

		var premium float64
		var timestamp int64

		timestamp = makeTimestamp();

		premium = group.GroupPremium * float64((findTimeStampDiff(group.EndDate,timestamp)/findTimeStampDiff(group.EndDate,group.CreatedDate)));

		var percentage []float64 
		percentage = getPremiumPercentages(len(group.RiskIds))

		risk.Premium = 	premium											// Premium calculated at the time of risk added to group
		risk.Status = "covered" 										// risk is insured 
		//append risk to group
		group.RiskIds = append(group.RiskIds, risk.Id)
		fmt.Println("! risk "+risk.Id+"added to group: ", group.Name)
		group.PoolBalance = group.PoolBalance + premium * percentage[0]						// premium * pool share %
		member.Tokens = member.Tokens-premium
		insurer.Tokens = insurer.Tokens + premium * percentage[1] 						// premium * insurer share %
		

		//Update Member
		jsonAsBytes, _ := json.Marshal(member)
		err = stub.PutState(member.UserId, jsonAsBytes)				
		if err != nil {
			return nil, err
		}
		//Update Risk
		jsonAsBytes, _ = json.Marshal(risk)
		err = stub.PutState(risk.Id, jsonAsBytes)				
		if err != nil {
			return nil, err
		}
		//Update Group
		jsonAsBytes, _ = json.Marshal(group)
		err = stub.PutState(group.Name, jsonAsBytes)
		if err != nil {
			return nil, err
		}
		//Update Insurer
		jsonAsBytes, _ = json.Marshal(insurer)
		err = stub.PutState(insurer.Id, jsonAsBytes)				
		if err != nil {
			return nil, err
		}
		//Update Admin Account
		err = stub.PutState(ADMIN_FEE, []byte(strconv.FormatFloat(previous_val + risk.Premium- premium * percentage[0]-premium * percentage[1], 'E', -1, 64)));
		if err != nil {
			return nil, errors.New("Billing account not updated")
		}

	}else{
		return nil,  errors.New("Not enough funds for"+ member.UserId)
	}

	return nil, nil
}

// ============================================================================================================================
/* 
RaiseClaim - Invoke function to raise a new claim write key/value pair
Inputs: 	args[0]		args[1]		args[2]
			riskId 		claimed 	type
			"ri004"		"15.5" 		"damage" 	
*/
// ============================================================================================================================
//func (t *SimpleChaincode) RaiseClaim(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
func (t *SimpleChaincode) RaiseClaim(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	fmt.Println("running RaiseClaim()")

	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3. ")
	}

	claimedAmount, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		return nil, errors.New("2nd argument must be a numeric string")
	}

	//Get Risk
	riskAsBytes, err := stub.GetState(args[0])
	if err != nil {
		return nil, errors.New("Failed to get risk")
	}
	risk := Risk{}
	json.Unmarshal(riskAsBytes, &risk)								//un stringify it aka JSON.parse()

	if claimedAmount > risk.Value {
		return nil, errors.New("Claimed value cannot be greater than risk value ")
	}

	claim := Claim{}
	claim.Id = makeClaimId()
	claim.RiskId = args[0]
	claim.Claimed = claimedAmount
	//claim.Settled =												//Will be updated when claim is settled
	claim.Timestamp = makeTimestamp()
	claim.Type = args[2]


	//Put claim in chain
	jsonAsBytes, _ := json.Marshal(claim)
	err = stub.PutState(claim.Id, jsonAsBytes)				
	if err != nil {
		return nil, err
	}

	//Get chainIndex
	claimIndexAsBytes, err := stub.GetState(CLAIM_INDEX)
	if err != nil {
		return nil, errors.New("Failed to get claim index")
	}

	var claimIndex []string
	json.Unmarshal(claimIndexAsBytes, &claimIndex)								//un stringify it aka JSON.parse()
	
	//append
	claimIndex = append(claimIndex, claim.Id)										//add chain id to index list
	fmt.Println("! claim index: ", claimIndex)
	jsonAsBytes, _ = json.Marshal(claimIndex)
	//Update chainIndex in chain
	err = stub.PutState(CLAIM_INDEX, jsonAsBytes)								//store chain id of risk

	risk.ClaimIds = append(risk.ClaimIds, claim.Id)
	
	jsonAsBytes, _ = json.Marshal(risk)
	err = stub.PutState(risk.Id, jsonAsBytes)				
	if err != nil {
		return nil, err
	}

	return nil, nil
}

/*	UTILITY FUNCTIONS	*/

// ============================================================================================================================
// Make Timestamp - create a timestamp in ms
// ============================================================================================================================
func makeTimestamp() int64 {
    return time.Now().UnixNano() / (int64(time.Millisecond)/int64(time.Nanosecond))
}

func findTimeStampDiff(ts1 int64,ts2 int64) int64{
	return (ts1-ts2)/3600/24
}

func stringInSlice(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}

func makeRiskId() string {
	riskCounter = riskCounter+1
	id :="rid00"+strconv.Itoa(riskCounter)
	return id
}

func makeMemberId() string {
	memberCounter = memberCounter+1
	id :="uid00"+strconv.Itoa(memberCounter)
	return id
}

func makeClaimId() string {
	claimCounter = claimCounter+1
	id :="cid00"+strconv.Itoa(claimCounter)
	return id
}

func makeInsurerId() string {
	insurerCounter = insurerCounter+1
	id :="ins00"+strconv.Itoa(insurerCounter)
	return id
}

func makeBicycleId() string {
	bicycleCounter = bicycleCounter+1
	id :="bi00"+strconv.Itoa(bicycleCounter)
	return id
}
/*
grp-size  pool  fees  insurer
   < 20       25%   5%    70% 
 20-39       30%   5%    65%
 40-59       35%   5%    60%
 60-79       40%   5%    55%
 80-99       45%   5%    50%
 > 100       50%   5%    45%
*/
func getPremiumPercentages(size int) []float64{
	var percentages []float64
	if size < 20 {
		percentages = append(percentages, 0.25)										//pool share
		percentages = append(percentages, 0.70)										//insurer share
	} else if size >= 20 && size <40{
		percentages = append(percentages, 0.30)
		percentages = append(percentages, 0.65)
	} else if size >= 40 && size <60{
		percentages = append(percentages, 0.35)
		percentages = append(percentages, 0.60)
	} else if size >= 60 && size <70{
		percentages = append(percentages, 0.40)
		percentages = append(percentages, 0.55)
	} else if size >= 80 && size <100{
		percentages = append(percentages, 0.45)
		percentages = append(percentages, 0.50)
	} else if size >= 100{
		percentages = append(percentages, 0.50)
		percentages = append(percentages, 0.45)
	}

	return percentages

} 