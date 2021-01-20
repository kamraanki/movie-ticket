package main
import (
	"github.com/op/go-logging" 
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type MovieTicket struct {
	contractapi.Contract
}


func main() {
	log := logging.MustGetLogger(name)
	var chaincode *contractapi.ContractChaincode
	var err error
	if chaincode, err = contractapi.NewChaincode(new(MovieTicket)); err != nil {
		log.Errorf("Error while creating chaincode: %s", err.Error())
		return
	}

	if err = chaincode.Start(); err != nil {
		log.Errorf("Error while starting chaincode: %s", err.Error())
	}
}
