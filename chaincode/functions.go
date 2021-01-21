package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-protos-go/ledger/queryresult"
	"github.com/op/go-logging"
	"strconv"
	"time"
)

func (s *MovieTicket) Init(ctx contractapi.TransactionContextInterface) error {
	log := logging.MustGetLogger(name)
	log.Infof("Chaincode initialized successfully")
	return nil
}

/**
	Method to register a theatre
*/
func (s *MovieTicket) Register_theatre(ctx contractapi.TransactionContextInterface, theatreStr string) error {
	log := logging.MustGetLogger(name)
	theatre := new(Theatre)
	if err := json.Unmarshal([]byte(theatreStr), &theatre); err != nil {
		log.Errorf("Invalid json input: %s, Error: %s", theatreStr, err.Error())
		return fmt.Errorf("Invalid json input: %s, Error: %s", theatreStr, err.Error())
	}

	// Check whether theatre id already registered or not
	if data, err := ctx.GetStub().GetState(theatre.TheatreId); err != nil {
		log.Errorf("Failed to get state for theatre id: %s, Got error: %s", theatre.TheatreId, err.Error())
		return fmt.Errorf("Failed to get state for theatre id: %s, Got error: %s", theatre.TheatreId, err.Error())
	} else if data != nil {
		log.Errorf("Theatre with theatre id %s already registered", theatre.TheatreId)
		return fmt.Errorf("Theatre with theatre id %s already registered", theatre.TheatreId)
	}
	theatre.RecordType = 3
	// Theatre with the same theatre id is not registered.
	theatreAsBytes, _ := json.Marshal(theatre)

	if err := ctx.GetStub().PutState(theatre.TheatreId, theatreAsBytes); err != nil {
		log.Errorf("Failed to register theatre with theatre id: %s, Error: %s", theatre.TheatreId, err.Error())
		return fmt.Errorf("Failed to register theatre with theatre id: %s, Error: %s", theatre.TheatreId, err.Error())
	}

	// Register cafeteria
	cafeteria := new(Cafeteria)
	cafeteriaAsBytes, _ := json.Marshal(cafeteria)

	if err := ctx.GetStub().PutState("cafeteria_"+theatre.TheatreId, cafeteriaAsBytes); err != nil {
		log.Errorf("Failed to register cafeteria with theatre id: %s, Error: %s", theatre.TheatreId, err.Error())
		return fmt.Errorf("Failed to register cafeteria with theatre id: %s, Error: %s", theatre.TheatreId, err.Error())
	}

	log.Infof("Theatre with theatre id: %s registered successfully !!", theatre.TheatreId)
	return nil
}

/**
	Method to register a show
*/
func (s *MovieTicket) Register_show(ctx contractapi.TransactionContextInterface, showStr string) error {
	log := logging.MustGetLogger(name)
	var err error
	var data []byte
	show := new(Show)
	if err = json.Unmarshal([]byte(showStr), &show); err != nil {
		log.Errorf("Invalid json input: %s, Error: %s", showStr, err.Error())
		return fmt.Errorf("Invalid json input: %s, Error: %s", showStr, err.Error())
	}

	// Check whether provided theatre id and movie hall id is valid or not
	if data, err = ctx.GetStub().GetState(show.TheatreId); err != nil {
		log.Errorf("Failed to get state for theatre id: %s, Got error: %s", show.TheatreId, err.Error())
		return fmt.Errorf("Failed to get state for theatre id: %s, Got error: %s", show.TheatreId, err.Error())
	} else if data == nil {
		log.Errorf("Theatre with theatre id %s does not exist", show.TheatreId)
		return fmt.Errorf("Theatre with theatre id %s does not exist", show.TheatreId)
	} else {
		theatre := new(Theatre)
		_ = json.Unmarshal([]byte(data), &theatre)
		if show.MovieHallNo < 1 || show.MovieHallNo > theatre.MovieHallNos {
			log.Errorf("Movie hall no %d in theatre %s does not exist.", show.MovieHallNo, theatre.TheatreId)
			return fmt.Errorf("Movie hall no %d in theatre %s does not exist.", show.MovieHallNo, theatre.TheatreId)
		}
	}

	var showStartDate, showEndDate time.Time
	if showStartDate, err = time.Parse("2006-01-02", show.ShowStartDate); err != nil {
		// Start date parsing issue
		log.Errorf("Invalid show start date: %s, Error: %s", show.ShowStartDate, err.Error())
		return fmt.Errorf("Invalid show start date: %s, Error: %s", show.ShowStartDate, err.Error())
	}

	if showEndDate, err = time.Parse("2006-01-02", show.ShowEndDate); err != nil {
		// End date parsing issue
		log.Errorf("Invalid show end date: %s, Error: %s", show.ShowEndDate, err.Error())
		return fmt.Errorf("Invalid show end date: %s, Error: %s", show.ShowEndDate, err.Error())
	}

	if showEndDate.Before(showStartDate) {
		// End date < Start date
		log.Errorf("Invalid show start & end dates. Start date: %s, End date: %s", show.ShowStartDate, show.ShowEndDate)
		return fmt.Errorf("Invalid show start & end dates. Start date: %s, End date: %s", show.ShowStartDate, show.ShowEndDate)
	}

	for d := showStartDate; !d.After(showEndDate); d = d.AddDate(0, 0, 1) {
		// Register show for each day
		key, _ := getCompositeKey(ctx, showKeyIndex, show.TheatreId, d.Format("2006-01-02"), show.ShowTime, strconv.Itoa(show.MovieHallNo))
		// Check whether any existing show exist on same date and time
		if data, err = ctx.GetStub().GetState(key); err != nil {
			log.Errorf("Failed to get state for existing show, Got error: %s", err.Error())
			return fmt.Errorf("Failed to get state for existing show, Got error: %s", err.Error())
		} else if data != nil {
			log.Errorf("Show already exist on date: %s and time %s", d.Format("2006-01-02"), show.ShowTime)
			return fmt.Errorf("Show already exist on date: %s and time %s", d.Format("2006-01-02"), show.ShowTime)
		} else {
			// register the show
			show.ShowDate = d.Format("2006-01-02")
			show.RecordType = 1
			showAsBytes, _ := json.Marshal(show)
			if err := ctx.GetStub().PutState(key, showAsBytes); err != nil {
				log.Errorf("Failed to register show with show id: %s, Error: %s", show.ShowId, err.Error())
				return fmt.Errorf("Failed to register show with show id: %s, Error: %s", show.ShowId, err.Error())
			}
		}
	}

	log.Infof("Show with show id: %s registered successfully !!", show.ShowId)
	return nil
}

/**
	Method to add soda bottles to cafeteria's inventory
*/
func (s *MovieTicket) Add_cafeteria_inventory(ctx contractapi.TransactionContextInterface, theatreId string, sodaBottleQuantity int) error {
	log := logging.MustGetLogger(name)

	// Get the cafeteria record
	if data, err := ctx.GetStub().GetState("cafeteria_" + theatreId); err != nil {
		log.Errorf("Failed to get state for theatre cafeteria for theatre id: %s, Got error: %s", theatreId, err.Error())
		return fmt.Errorf("Failed to get state for theatre cafeteria for theatre id: %s, Got error: %s", theatreId, err.Error())
	} else if data == nil {
		log.Errorf("Theatre with theatre id %s does not exist", theatreId)
		return fmt.Errorf("Theatre with theatre id %s does not exist", theatreId)
	} else {
		// Add soda bottle quantity
		cafeteria := new(Cafeteria)
		_ = json.Unmarshal([]byte(data), &cafeteria)
		cafeteria.SodaBottleQuantity += sodaBottleQuantity
		cafeteriaAsBytes, _ := json.Marshal(cafeteria)
		if err := ctx.GetStub().PutState("cafeteria_"+theatreId, cafeteriaAsBytes); err != nil {
			log.Errorf("Failed to register cafeteria with theatre id: %s, Error: %s", theatreId, err.Error())
			return fmt.Errorf("Failed to register cafeteria with theatre id: %s, Error: %s", theatreId, err.Error())
		}

		log.Infof("Inventry added to cafeteria successfully for theatre id: %s", theatreId)
		return nil
	}
}

/**
	Method to get list of shows using rich query
*/
func (s *MovieTicket) Get_shows(ctx contractapi.TransactionContextInterface, showSearchQueryStr string) (*ShowSearchResult, error) {
	log := logging.MustGetLogger(name)
	showSearchQuery := new(ShowSearchQuery)
	if err := json.Unmarshal([]byte(showSearchQueryStr), &showSearchQuery); err != nil {
		log.Errorf("Invalid json input: %s, Error: %s", showSearchQueryStr, err.Error())
		return nil, fmt.Errorf("Invalid json input: %s, Error: %s", showSearchQueryStr, err.Error())
	}
	
	// Create rich query string
	queryString := CreateShowSearchQuery(showSearchQuery)
	log.Info("Querying chaincode with query string: %s", queryString)
	
	// Execute couchdb rich query to get list of all available shows
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		log.Errorf("Error while running rich query: %s, error: %s", queryString, err.Error())
		return nil, fmt.Errorf("Error while running rich query: %s, error: %s", queryString, err.Error())
	}

	var shows []Show
	var queryResult *queryresult.KV
	for resultsIterator.HasNext() {
		queryResult, err = resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("Got error: %s", err.Error())
		}
		var show Show
		_ = json.Unmarshal(queryResult.Value, &show)
		shows = append(shows, show)
	}

	showSearchResult := new(ShowSearchResult)
	if shows != nil {
		showSearchResult.ShowList = shows
	} else {
		showSearchResult.ShowList = []Show{}
	}

	return showSearchResult, nil

}

/**
	Method to get no. of available seats/ ticket
*/
func (s *MovieTicket) Get_seat_availability(ctx contractapi.TransactionContextInterface, seatAvailabilityQueryStr string) (int, error) {
	log := logging.MustGetLogger(name)
	query := new(SeatAvailabilityQuery)

	if err := json.Unmarshal([]byte(seatAvailabilityQueryStr), &query); err != nil {
		log.Errorf("Invalid json input: %s, Error: %s", seatAvailabilityQueryStr, err.Error())
		return 0, fmt.Errorf("Invalid json input: %s, Error: %s", seatAvailabilityQueryStr, err.Error())
	} else if query.TheatreId == "" || query.ShowId == "" || query.ShowDate == "" || query.ShowTime == "" || query.MovieHallNo < 1 {
		log.Errorf("Invalid json input: %s", seatAvailabilityQueryStr)
		return 0, fmt.Errorf("Invalid json input: %s", seatAvailabilityQueryStr)
	}
	
	// Get no. of available seats
	availableSeats, err := GetSeatAvailability(ctx, query.TheatreId, query.ShowId, query.ShowDate, query.ShowTime, query.MovieHallNo)
	if err != nil {
		log.Errorf("Got error: %s", err.Error())
		return availableSeats, fmt.Errorf("Error: %s", err.Error())
	}

	return availableSeats, nil

}

/**
	Method to book a seat/ ticket
*/
func (s *MovieTicket) Book_ticket(ctx contractapi.TransactionContextInterface, ticketStr string) error {
	log := logging.MustGetLogger(name)
	ticket := new(Ticket)

	// Validate ticket json
	if err := json.Unmarshal([]byte(ticketStr), &ticket); err != nil {
		log.Errorf("Invalid json input: %s, Error: %s", ticketStr, err.Error())
		return fmt.Errorf("Invalid json input: %s, Error: %s", ticketStr, err.Error())
	} else if ticket.TheatreId == "" || ticket.ShowId == "" || ticket.ShowDate == "" || ticket.ShowTime == "" || ticket.MovieHallNo < 1 || ticket.NoOfSeats < 1 || ticket.LuckyNo < 1 {
		log.Errorf("Invalid json input: %s", ticketStr)
		return fmt.Errorf("Invalid json input: %s", ticketStr)
	}

	// Check availableSteats should be >= requiredSeats
	if availableSeats, err := GetSeatAvailability(ctx, ticket.TheatreId, ticket.ShowId, ticket.ShowDate, ticket.ShowTime, ticket.MovieHallNo); err != nil {
		log.Errorf("Got error: %s", err.Error())
		return fmt.Errorf("Got error: %s", err.Error())
	} else if availableSeats < ticket.NoOfSeats {
		log.Errorf("Seats not available")
		return fmt.Errorf("SEATS_NOT_AVAILABLE")
	}

	// Register ticket
	ticket.RecordType = 2
	ticketAsBytes, _ := json.Marshal(ticket)

	if err := ctx.GetStub().PutState(ticket.TicketId, ticketAsBytes); err != nil {
		log.Errorf("Failed to register ticket with ticket id: %s, Error: %s", ticket.TicketId, err.Error())
		return fmt.Errorf("Failed to register ticket with ticket id: %s, Error: %s", ticket.TicketId, err.Error())
	}

	return nil
}

/**
	Method to replace water bottle with soda bottle
*/
func (s *MovieTicket) Replace_with_soda_bottle(ctx contractapi.TransactionContextInterface, ticketId string) (bool, error) {
	log := logging.MustGetLogger(name)

	ticket := new(Ticket)

	// Check whether ticket id is valid or not, if valid get the ticket
	if data, err := ctx.GetStub().GetState(ticketId); err != nil {
		log.Errorf("Failed to get state for ticket id: %s, Got error: %s", ticketId, err.Error())
		return false, fmt.Errorf("Failed to get state for ticket id: %s, Got error: %s", ticketId, err.Error())
	} else if data == nil {
		log.Errorf("Invalid ticket id %s", ticketId)
		return false, fmt.Errorf("Invalid ticket id %s", ticketId)
	} else if err = json.Unmarshal([]byte(data), &ticket); err != nil {
		log.Errorf("Failed to get state for ticket id: %s, Got error: %s", ticketId, err.Error())
		return false, fmt.Errorf("Failed to get state for ticket id: %s, Got error: %s", ticketId, err.Error())
	} else if ticket.RecordType != 2 {
		log.Errorf("Invalid ticket id %s", ticketId)
		return false, fmt.Errorf("Invalid ticket id %s", ticketId)
	} else if ticket.LuckyNo%2 != 0 {
		log.Errorf("Not elegible for soda bottle replacement. Ticket id: %s", ticketId)
		return false, fmt.Errorf("NOT_ELIGIBLE")
	}

	// Check soda bottle is not replaced already for this ticket id
	if data, err := ctx.GetStub().GetState("replace_" + ticketId); err != nil {
		log.Errorf("Failed to check whether bottle is already replace for ticket or not. Ticket Id: %s, Error: %s", ticketId, err.Error())
		return false, fmt.Errorf("Failed to check whether bottle is already replace for ticket or not. Ticket Id: %s, Error: %s", ticketId, err.Error())
	} else if data != nil {
		log.Errorf("Soda bottle is already replaced for ticket id: %s", ticketId)
		return false, fmt.Errorf("ALREADY_REPLACED")
	}

	// Check cafeteria inventory
	cafeteria := new(Cafeteria)
	if data, err := ctx.GetStub().GetState("cafeteria_" + ticket.TheatreId); err != nil {
		log.Errorf("Failed to get state for cafeteria for theatre id: %s, Got error: %s", ticket.TheatreId, err.Error())
		return false, fmt.Errorf("Failed to get state for cafeteria for theatre id: %s, Got error: %s", ticket.TheatreId, err.Error())
	} else if err = json.Unmarshal([]byte(data), &cafeteria); err != nil {
		log.Errorf("Failed to get state for cafeteria for theatre id: %s, Got error: %s", ticket.TheatreId, err.Error())
		return false, fmt.Errorf("Failed to get state for cafeteria for theatre id: %s, Got error: %s", ticket.TheatreId, err.Error())
	} else if cafeteria.SodaBottleQuantity < 1 {
		log.Errorf("Soda bottle is out of stock for theatre id: %s", ticket.TheatreId)
		return false, fmt.Errorf("OUT_OF_STOCK")
	}

	sodaBottleReplacement := new(SodaBottleReplacement)
	sodaBottleReplacement.TicketId = ticketId
	sodaBottleReplacementAsBytes, _ := json.Marshal(sodaBottleReplacement)
	if err := ctx.GetStub().PutState("replace_"+ticketId, sodaBottleReplacementAsBytes); err != nil {
		log.Errorf("Failed to write soda replacement record for ticket id: %s, Error: %s", ticketId, err.Error())
		return false, fmt.Errorf("Failed to write soda replacement record for ticket id: %s, Error: %s", ticketId, err.Error())
	}

	cafeteria.SodaBottleQuantity--
	cafeteriaAsBytes, _ := json.Marshal(cafeteria)
	if err := ctx.GetStub().PutState("cafeteria_"+ticket.TheatreId, cafeteriaAsBytes); err != nil {
		log.Errorf("Failed to update cafeteria with theatre id: %s, Error: %s", ticket.TheatreId, err.Error())
		return false, fmt.Errorf("Failed to update cafeteria with theatre id: %s, Error: %s", ticket.TheatreId, err.Error())
	}

	return true, nil
}
