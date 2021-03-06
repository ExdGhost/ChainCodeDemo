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
  CashBalance float64 `json:"CashBalance"`
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

     err3:= stub.PutState("Tcounter",[]byte("1"))

   if err3 != nil {

    	fmt.Println("error creating Transaction counter")
		return nil, errors.New("Error Transaction creating counter ")
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

      if function == "seedToken" {
        return t.seedToken(stub,args)
    }

    if function == "sendToken" {
        return t.sendToken(stub,args)
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

   var uSuffix,uPrefix,uID,tID string
   var initialCash float64
   var errf error
   var token,newToken MyToken

if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3")
	} 

    uID = args[0]
    tID = args[1]
    uSuffix="000U"

    initialCash,errf = strconv.ParseFloat(args[2],64)

     if errf != nil {
	    return nil,errors.New("could not convert cash")
     }
    
    counterBytes,err:= stub.GetState("counter") 

    if err != nil {
    	fmt.Println("could not retrieve counter")			//error
	    return nil,errors.New("could not retrieve counter")
     }
     
    tokenBytes,erry:= stub.GetState(tID)
    
    if erry != nil {
        return nil,errors.New("Could not find the specified Token")
        }

       erry=json.Unmarshal(tokenBytes,&token)
      
       if erry != nil {
        return nil,errors.New("Could not unmsarshal token")
        }

    counter,errx:= strconv.Atoi(string(counterBytes))

       if errx != nil {
        return nil,errx
        }

      //Initialize the new token

       newToken = MyToken{TID:token.TID,Tname:token.Tname,Tamount:token.Tamount}
    
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
       
       //Initialize the user
       var user = Account{ID:uID,Prefix:uPrefix,CashBalance:initialCash,Token:newToken}
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

    return []byte(user.Prefix+"&"+user.ID+"&"+strconv.FormatFloat(user.CashBalance,'f',-1,64)+"&"+user.Token.TID+"&"+user.Token.Tamount), nil
}

//seed user account with token
func (t *SimpleChaincode) seedToken(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
     
    var user Account
    var jsonResp string

      if len(args) != 2 {
        return nil, errors.New("Incorrect number of arguments. Expecting user id and token amount to query")
     }

     uID := args[0]
     tamount := args[1]

   userBytes,err := stub.GetState(uID)

    if err != nil {
        jsonResp = "{\"Error\":\"Failed to get state for " + uID + "\"}"
        return nil, errors.New(jsonResp)
    }
   
   err = json.Unmarshal(userBytes,&user)

   if err != nil {
        jsonResp = "{\"Error\":\"Failed to get object for " + string(userBytes) + "\"}"
        return nil, errors.New(jsonResp)
    }

   /* newAmount,err1:= strconv.ParseFloat(tamount,64)
    
     if err1 != nil {
        jsonResp = "Failed to convert string to float64"
        return nil, errors.New(jsonResp)
    }*/

   user.Token.Tamount=tamount

   userBytes,err = json.Marshal(&user)

   if err != nil {
			fmt.Println("error Marshalling account for ID : " + user.ID)
			return nil,errors.New("error Marshalling account for ID : " + user.ID)
		}

      err = stub.PutState(uID,userBytes)

      if err != nil {
			fmt.Println("error updating account for ID : " + user.ID)
			return nil,errors.New("error updating account for ID : " + user.ID)
		}

   return nil,nil
    
}

//Transfer token from one user to another user
func (t *SimpleChaincode) sendToken(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	  var sender,receiver Account
      var err error
      var newAmount,senderTokenAmount,receiverTokenAmount int

      if len(args) != 3 {
        return nil, errors.New("Incorrect number of arguments.Expecting user id and token amount to query")
     }

	
	 senderID := args[0]
	 receiverID := args[1]
	  
	  //Conver the token amount to float64 
	  tamount,err := strconv.Atoi(args[2])

	  if err!= nil{
	  	return nil,errors.New("Could not convert from string to float in arguement")
	  }
       
      //Retrieve Sender account
      senderBytes,err := stub.GetState(senderID)
      
      if err != nil{
	  	return nil,errors.New("Could not find sender account")
	  }
      
      //Retrieve receiver account
	  receiverBytes,err := stub.GetState(receiverID)

	  if err != nil{
	  	return nil,errors.New("Could not find receiver account")
	  }

      //Unmarshal sender
      err = json.Unmarshal(senderBytes,&sender)

       if err != nil{
	  	return nil,errors.New("Could not Unmarshal sender")
	  }
      
       //Unmarshal receiver
	   err = json.Unmarshal(receiverBytes,&receiver)

       if err != nil{
	  	return nil,errors.New("Could not Unmarshal receiver")
	  }
     
   receiverTokenAmount,err = strconv.Atoi(receiver.Token.Tamount)

    if err != nil {
    	return nil,errors.New("Could not convert to int from receiver token amount")
    }

   senderTokenAmount,err = strconv.Atoi(sender.Token.Tamount)

    if err != nil {
    	return nil,errors.New("Could not convert to int from sender token amount")
    }
  
  //Make transaction

    //Substract amount to be transferred from sender and set new value in sender
  newAmount = senderTokenAmount - tamount

  sender.Token.Tamount = strconv.Itoa(newAmount)

  //Add amount to be transferred to receiver and set new value in reveiver
  newAmount = receiverTokenAmount + tamount

  receiver.Token.Tamount = strconv.Itoa(newAmount)
  
  //Re-Marshal new sender object
  senderBytes,err = json.Marshal(&sender)
  
  if err != nil {
    	return nil,errors.New("Could not Re-Marshal sender object")
    }
 
  //Commit updated sender to stub
  err = stub.PutState(senderID,senderBytes)

  if err != nil {
    	return nil,errors.New("Could not commit updated sender object")
    }

    //Re-Marshal new receiver object
  receiverBytes,err = json.Marshal(&receiver)
  
  if err != nil {
    	return nil,errors.New("Could not Re-Marshal receiver object")
    }
 
  //Commit updated receiver to stub
   err = stub.PutState(receiverID,receiverBytes)

  if err != nil {
    	return nil,errors.New("Could not commit updated receiver object")
    }

 return nil,nil
}
