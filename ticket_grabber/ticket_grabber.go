package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"strconv"
	"time"

	"code.google.com/p/goprotobuf/proto"
	"github.com/contester/printing3/tickets"
	"github.com/contester/printing3/tools"
	"github.com/jmoiron/sqlx"

	"github.com/contester/printing3/grabber"
	_ "github.com/go-sql-driver/mysql"
"gopkg.in/stomp.v1"
	"encoding/json"
	"github.com/contester/printing3/printserver"
)

func createSelectSubmitQuery(extraWhere string) string {
	return `select
Submits.Contest, Submits.Team, Submits.Touched, Submits.Task, Submits.Compiled,
(unix_timestamp(Submits.Arrived) - unix_timestamp(Contests.Start)) as Arrived,
Submits.Passed, Submits.Taken, Submits.ID, Areas.Printer, Contests.SchoolMode, Submits.TestingID,
Schools.Name, Teams.Num, Contests.Name as ContestName,
Problems.Name as ProblemName, inet_ntoa(Submits.Computer) as ComputerID, CompLocations.Name as ComputerName,
Areas.ID as AreaID, Areas.Name as AreaName
from Submits, Contests, Areas, CompLocations, Teams, Schools, Participants, Problems
where
Contests.ID = Submits.Contest and Submits.Finished and ` + extraWhere + `
Submits.Computer = CompLocations.ID and CompLocations.Location = Areas.ID and Participants.LocalID = Submits.Team
and Teams.ID = Participants.Team and Participants.Contest = Submits.Contest and Teams.School = Schools.ID and
Problems.Contest = Submits.Contest and Problems.ID = Submits.Task and Contests.PrintTickets = 1 and Areas.Printer is not NULL
order by Submits.Arrived asc
`
}

var (
	allSubmitsQuery     = createSelectSubmitQuery("((Submits.Printed is null) or (Submits.Printed < Submits.Touched)) and")
	relatedSubmitsQuery = createSelectSubmitQuery("Submits.Contest = ? and Submits.Task = ? and Submits.Team = ? and Submits.ID < ? and")
	submitById = createSelectSubmitQuery("Submits.ID = ? and")
)

type scannedSubmit struct {
	Contest, Team                                      int64
	Touched                                            time.Time
	Task                                               string
	Compiled, Arrived, Passed, Taken, ID               int64
	Printer                                            string
	SchoolMode, TestingID                              int64
	SchoolName                                         string
	TeamNum                                            sql.NullInt64
	ContestName, ProblemName, ComputerID, ComputerName string
	AreaID                                             int64
	AreaName                                           string
}

type submitTicket struct {
	Submit struct {
		Id int32 `json:"id"`
		   } `json:"submit"`
}

func scanSubmit(r grabber.RowOrRows) (result *scannedSubmit, err error) {
	var sub scannedSubmit
	if err = r.Scan(&sub.Contest, &sub.Team, &sub.Touched, &sub.Task, &sub.Compiled, &sub.Arrived, &sub.Passed,
		&sub.Taken, &sub.ID, &sub.Printer, &sub.SchoolMode, &sub.TestingID, &sub.SchoolName, &sub.TeamNum,
		&sub.ContestName, &sub.ProblemName, &sub.ComputerID, &sub.ComputerName, &sub.AreaID, &sub.AreaName); err == nil {
		result = &sub
	}
	return
}

func findRelatedSubmits(db *sqlx.DB, sub *scannedSubmit) ([]*scannedSubmit, error) {
	rows, err := db.Query(relatedSubmitsQuery, sub.Contest, sub.Task, sub.Team, sub.ID)
	if err != nil {
		log.Printf("Error scanning for related submits: %s", err)
		return nil, err
	}

	result := make([]*scannedSubmit, 0)

	for rows.Next() {
		s, _ := scanSubmit(rows)
		if s != nil {
			result = append(result, s)
		}
	}

	return result, nil
}

func createSubmit(db *sqlx.DB, sub *scannedSubmit, submitNo int) *tickets.Ticket_Submit {
	var result tickets.Ticket_Submit
	result.SubmitNumber = proto.Uint32(uint32(submitNo))
	if sub.Arrived > 0 {
		result.Arrived = proto.Uint64(uint64(sub.Arrived))
	}
	if result.Compiled = proto.Bool(sub.Compiled == 1); !result.GetCompiled() {
		return &result
	}
	if sub.SchoolMode != 0 {
		result.School = &tickets.Ticket_Submit_School{TestsTaken: proto.Uint32(uint32(sub.Taken)), TestsPassed: proto.Uint32(uint32(sub.Passed))}
	} else {
		var description string
		var test int64
		err := db.QueryRow("select ResultDesc.Description, Results.Test from Results, ResultDesc where "+
			"Results.UID = ? and ResultDesc.ID = Results.Result and not ResultDesc.Success order by Results.Test",
			sub.TestingID).Scan(&description, &test)
		switch {
		case err == sql.ErrNoRows:
			if sub.Passed != 0 && sub.Passed == sub.Taken {
				result.Acm = &tickets.Ticket_Submit_ACM{Result: proto.String("ACCEPTED")}
			}
		case err != nil:
			log.Fatal(err)
			return nil
		default:
			result.Acm = &tickets.Ticket_Submit_ACM{Result: &description, TestId: proto.Uint32(uint32(test))}
		}
	}
	return &result
}

func processSubmit(db *sqlx.DB, sender func(msg proto.Message) error, rows grabber.RowOrRows) error {
	sub, err := scanSubmit(rows)
	if err != nil {
		return err
	}
	related, err := findRelatedSubmits(db, sub)
	if err != nil {
		return err
	}

	result := tickets.Ticket{
		SubmitId: proto.Uint32(uint32(sub.ID)),
		Printer:  &sub.Printer,
		Computer: &tickets.Computer{Id: &sub.ComputerID, Name: &sub.ComputerName},
		Area:     &tickets.IdName{Id: proto.Uint32(uint32(sub.AreaID)), Name: &sub.AreaName},
		Contest:  &tickets.IdName{Id: proto.Uint32(uint32(sub.Contest)), Name: &sub.ContestName},
		Problem:  &tickets.Ticket_Problem{Id: &sub.Task, Name: &sub.ProblemName},
	}

	teamName := sub.SchoolName
	if sub.TeamNum.Valid && sub.TeamNum.Int64 > 0 {
		teamName = teamName + " #" + strconv.FormatInt(sub.TeamNum.Int64, 10)
	}
	result.Team = &tickets.IdName{Id: proto.Uint32(uint32(sub.Team)), Name: &teamName}
	result.JudgeTime = proto.Uint64(uint64(sub.Touched.UnixNano() / 1000))

	result.Submit = make([]*tickets.Ticket_Submit, 0)
	result.Submit = append(result.Submit, createSubmit(db, sub, len(related)+1))
	for count := len(related); count > 0; {
		count -= 1
		result.Submit = append(result.Submit, createSubmit(db, related[count], count))
	}

	if err = sender(&result); err != nil {
		return err
	}
	if _, err = db.Exec("Update Submits set Printed = Touched where ID = ?", sub.ID); err != nil {
		return err
	}
	fmt.Printf("Printed submit %d\n", sub.ID)
	return nil
}

type grserver struct {
	db *sqlx.DB
}

func (g grserver) processIncoming(conn *printserver.ServerConn, msg *stomp.Message) error {
	var ticket submitTicket
	if err := json.Unmarshal(msg.Body, &ticket); err != nil {
		log.Printf("Received malformed job: %s", err)
		return err
	}

	rows, err := g.db.Query(submitById, ticket.Submit.Id)
	if err != nil {
		log.Printf("Error looking up submit: %s", err)
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err = processSubmit(g.db, conn.Send, rows); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	flag.Parse()

	config, err := tools.ReadConfig()
	if err != nil {
		log.Fatal(err)
	}

	dbSpec, err := config.GetString("server", "db")
	if err != nil {
		log.Fatal(err)
	}

	var gs grserver
	gs.db, err = tools.CreateDb(dbSpec)
	if err != nil {
		log.Fatal(err)
	}

	pserver := printserver.Server{
		Source:      "/amq/queue/contester.tickets",
		Destination: "/amq/queue/ticket_pb",
	}

	pserver.StompConfig, err = tools.ParseStompFlagOrConfig("", config, "messaging")
	if err != nil {
		log.Fatalf("Opening stomp connection: %s", err)
		return
	}

	for {
		log.Println(pserver.Process(gs.processIncoming))
		time.Sleep(15 * time.Second)
	}
}
