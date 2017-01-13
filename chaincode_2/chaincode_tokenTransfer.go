package main

import (
  "encoding/json"
  "errors"
  "fmt"
  "strconv"
  "github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

//Token Struct

type MyToken struct{
	
	TID     string   `json:"tid"`
	Tname   string   `json:"tname"`
   Tamount string   `json:"tamount"`
}  

//Account struct 

type Account struct{ 
  ID     string  `json:"id"`
  Prefix string  `json:"prefix"`
  Token MyToken  `json:"token"`
}

//Transaction struct

type Transaction struct{

	TransactionID string  `json:"transactionid"`
	FromAccountID string  `json:"fromaccountid"`
	ToAccountID   string  `json:"toaccountid"`
	TokenID       string  `json:"tokenid"`
	Quantity      string  `json:"quantity"`

}

// ============================================================================================================================
// Main
// ============================================================================================================================

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

//Initialize
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	 
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	} 

    var token MyToken

    token = MyToken{TID: args[0],Tname: args[1],Tamount: "0"}
   
    tokenBytes,err := json.Marshal(&token)

    if err != nil {

    	fmt.Println("error Initializing token" + token.Tname)
		return nil, errors.New("Error Initializing token " + token.Tname)
    }

     err1:= stub.PutState(args[0],tokenBytes)

    if err1 != nil {

    	fmt.Println("error creating token" + token.Tname)
		return nil, errors.New("Error creating token " + token.Tname)
    }

   err2:= stub.PutState("counter",[]byte("1"))

   if err2 != nil {

    	fmt.Println("error creating counter")
		return nil, errors.New("Error creating counter ")
    }

  return nil, nil
}

//Invoke
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	   fmt.Println("invoke is running " + function)

	// Handle different functions
	  if function == "init" {
        return t.Init(stub, "init", args)
    }

    if function == "createUser" {
        return t.createUser(stub,args)
    }

     
    fmt.Println("invoke did not find func: " + function)					//error
	return nil, errors.New("Received unknown function invocation: " + function)

}

//Query
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)
    
    if function == "getUser"{
		return t.getUser(stub,args)
	}

    fmt.Println("query did not find func: " + function)			//error
	return nil, errors.New("Received unknown function query: " + function)
}

// Create User
func (t *SimpleChaincode) createUser(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

   var uSuffix,uPrefix,uID string

if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	} 

    uID = args[0]
    uSuffix="000U"
    counterBytes,err:= stub.GetState("counter") 

    if err != nil {
    	fmt.Println("could not retrieve counter")			//error
	    return nil,errors.New("could not retrieve counter")
     }

    counter,errx:= strconv.Atoi(string(counterBytes))

       if errx != nil {
        return nil,errx
        }
    
    //fmt.Println("Counter = %d",counter)

    if counter < 10 {
    	uPrefix = strconv.Itoa(counter)+"0"+uSuffix
    } else{
    	uPrefix = strconv.Itoa(counter)+uSuffix
    }

    counter = counter + 1   
    err1:= stub.PutState("counter",[]byte(strconv.Itoa(counter))) //Set new value of counter back to stub
   
   if err1 != nil {
    	fmt.Println("could not update counter")			//error
	    return nil,errors.New("could not update counter")
     }

       var user = Account{ID:uID,Prefix:uPrefix}
       userBytes,err2:= json.Marshal(&user)

  if err2 != nil {
			fmt.Println("error creating account" + user.ID)
			return nil,errors.New("Error creating account " + user.ID)
		}

     err3 := stub.PutState(uPrefix+uID,userBytes)  //comit account to stub

     if err3 != nil {
			fmt.Println("error commiting account" + user.ID)
			return nil, errors.New("Error commiting account " + user.ID)
		}
   
   //fmt.Println("UID = %s",user.Prefix+user.ID)

   return []byte(uPrefix+uID),nil
}

//get an user from the stub and return the user details as json 

func (t *SimpleChaincode) getUser(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    var key, jsonResp string
    var user Account

    if len(args) != 1 {
        return nil, errors.New("Incorrect number of arguments. Expecting user id key to query")
    }

    key = args[0]
    valAsbytes,err:= stub.GetState(key)
    
    if err != nil {
        jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
        return nil, errors.New(jsonResp)
    }

    err1:= json.Unmarshal(valAsbytes,&user)

    if err1 != nil {
        jsonResp = "{\"Error\":\"Failed to get object for " + string(valAsbytes) + "\"}"
        return nil, errors.New(jsonResp)
    }

    return []byte(user.Prefix+'&'+user.ID), nil
}