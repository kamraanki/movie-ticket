package main

import (
	"encoding/json"
	"errors"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-protos-go/ledger/queryresult"
	"strconv"
)

/**
	Function to create composite key
**/
func getCompositeKey(ctx contractapi.TransactionContextInterface, compositeKeyIndex string, args ...string) (string, error) {
	key, err := ctx.GetStub().CreateCompositeKey(compositeKeyIndex, args)
	if err != nil {
		return "", errors.New("Error in creating composite Key : " + err.Error())
	}
	return key, nil
}

/**
	Function to create rich query string
*/
func CreateShowSearchQuery(showSearchQuery *ShowSearchQuery) string {
	queryStr := "{\"selector\":{\"recordType\": 1"

	if showSearchQuery.TheatreId != "" {
		queryStr = queryStr + ",\"theatreId\":\"" + showSearchQuery.TheatreId + "\""
	}

	if showSearchQuery.ShowId != "" {
		queryStr = queryStr + ",\"showId\":\"" + showSearchQuery.ShowId + "\""
	}

	if showSearchQuery.ShowName != "" {
		queryStr = queryStr + ",\"showName\":\"" + showSearchQuery.ShowName + "\""
	}

	if showSearchQuery.ShowTime != "" {
		queryStr = queryStr + ",\"showTime\":\"" + showSearchQuery.ShowTime + "\""
	}

	if showSearchQuery.ShowDate != "" {
		queryStr = queryStr + ",\"showDate\":\"" + showSearchQuery.ShowDate + "\""
	}

	return queryStr + "}}"
}

/**
	Function to get no of available seats
*/
func GetSeatAvailability(ctx contractapi.TransactionContextInterface, theatreId, showId, showDate, showTime string, movieHallNo int) (int, error) {

	var totalTicketsAvailable int

	// Get Show
	key, _ := getCompositeKey(ctx, showKeyIndex, theatreId, showDate, showTime, strconv.Itoa(movieHallNo))
	if data, err := ctx.GetStub().GetState(key); err != nil {
		return 0, err
	} else if data == nil {
		return 0, errors.New("INVALID_SHOW_INFO")
	}

	// Get Theatre
	if data, err := ctx.GetStub().GetState(theatreId); err != nil {
		return 0, err
	} else if data == nil {
		return 0, errors.New("INVALID_THEATRE_ID")
	} else {
		theatre := new(Theatre)
		_ = json.Unmarshal([]byte(data), &theatre)
		totalTicketsAvailable = theatre.TicketsPerShow
	}

	// Get no. of sold tickets
	queryString := "{\"selector\":{\"theatreId\":\"" + theatreId + "\",\"showId\":\"" + showId + "\",\"showDate\":\"" + showDate + "\",\"showTime\":\"" + showTime + "\",\"movieHallNo\":" + strconv.Itoa(movieHallNo) + ",\"recordType\":2}}"
	// get all tickets and count no. of sold tickets
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return 0, err
	}

	soldTickets := 0
	var queryResult *queryresult.KV
	for resultsIterator.HasNext() {
		queryResult, err = resultsIterator.Next()
		if err != nil {
			return 0, err
		}
		var ticket Ticket
		_ = json.Unmarshal(queryResult.Value, &ticket)
		soldTickets += ticket.NoOfSeats
	}

	return totalTicketsAvailable - soldTickets, nil
}
