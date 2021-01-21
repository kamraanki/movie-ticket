package main

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/op/go-logging"
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
