package main

const (
	name = "movie-ticket"
	showKeyIndex = "TheatreId~ShowDate~ShowTime~MovieHallNo"
)

/*
Assumptions:

1. All master data like theatre name, show name etc.will be handeled in DB by middleware. Chaincode will work with Theatre id, show id etc.
2. All validations will be done at middleware level.
3. Show time should be in format 00:00 - 23:59
*/

type Theatre struct {	
	// Represents a theatre structure
	TheatreId				string			`json:"theatreId"`
	MovieHallNos			int				`json:"movieHallNos"`	// No's of movie hall available in theatre
	TicketsPerShow			int				`json:"ticketsPerShow"`	// No's of seat per movie hall i.e. max no's of tickets per show can be sold
	TicketWindowNos			int				`json:"ticketWindowNos"`// No's of ticket windows
	RecordType				int				`json:"RecordType"`		// 3 for theatre
}

type Show struct {
	TheatreId				string			`json:"theatreId"`
	MovieHallNo				int				`json:"movieHallNo"`
	ShowId					string			`json:"showId"`	
	ShowName				string			`json:"showName"`
	ShowDate				string			`json:"showDate"`
	ShowStartDate			string			`json:"showStartDate"`
	ShowEndDate				string			`json:"showEndDate"` 
	ShowTime				string			`json:"showTime"`
	RecordType				int				`json:"recordType"`		// 1 for show
}

type ShowSearchQuery struct {
	TheatreId				string			`json:"theatreId"`
	ShowId					string			`json:"showId"`	
	ShowName				string			`json:"showName"`
	ShowDate				string			`json:"showDate"`
	ShowTime				string			`json:"showTime"`
}

type SeatAvailabilityQuery struct {
	TheatreId				string			`json:"theatreId"`
	ShowId					string			`json:"showId"`	
	ShowDate				string			`json:"showDate"`
	ShowTime				string			`json:"showTime"`
	MovieHallNo				int				`json:"movieHallNo"`
}

type Cafeteria struct {
	SodaBottleQuantity		int				`json:"sodaBottleQuantity"`	// Soda bottle quantity available in cafeteria
}

type Ticket struct {
	TicketId				string			`json:"ticketId"`
	TheatreId				string			`json:"theatreId"`
	ShowId					string			`json:"showId"`
	ShowDate				string			`json:"showDate"`
	ShowTime				string			`json:"showTime"`
	MovieHallNo				int				`json:"movieHallNo"`
	NoOfSeats				int				`json:"noOfSeats"`
	LuckyNo					int				`json:"luckyNo"`
	RecordType				int				`json:"recordType"`		// 2 for ticket		
}

type SodaBottleReplacement struct {
	TicketId				string			`json:"ticketId"`
}

type ShowSearchResult struct {
	ShowList				[]Show			`json:"showList"`
}

	