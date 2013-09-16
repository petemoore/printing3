// Code generated by protoc-gen-go.
// source: tickets.proto
// DO NOT EDIT!

package tickets

import proto "code.google.com/p/goprotobuf/proto"
import json "encoding/json"
import math "math"

// Reference proto, json, and math imports to suppress error if they are not otherwise used.
var _ = proto.Marshal
var _ = &json.SyntaxError{}
var _ = math.Inf

type Ticket struct {
	ContestId        *uint32          `protobuf:"varint,1,opt,name=contest_id" json:"contest_id,omitempty"`
	ContestName      *string          `protobuf:"bytes,2,opt,name=contest_name" json:"contest_name,omitempty"`
	SubmitId         *uint32          `protobuf:"varint,3,opt,name=submit_id" json:"submit_id,omitempty"`
	JudgeTime        *string          `protobuf:"bytes,4,opt,name=judge_time" json:"judge_time,omitempty"`
	TeamId           *uint32          `protobuf:"varint,5,opt,name=team_id" json:"team_id,omitempty"`
	TeamName         *string          `protobuf:"bytes,6,opt,name=team_name" json:"team_name,omitempty"`
	LocationId       *uint32          `protobuf:"varint,7,opt,name=location_id" json:"location_id,omitempty"`
	LocationName     *string          `protobuf:"bytes,8,opt,name=location_name" json:"location_name,omitempty"`
	ComputerId       *string          `protobuf:"bytes,9,opt,name=computer_id" json:"computer_id,omitempty"`
	ComputerName     *string          `protobuf:"bytes,10,opt,name=computer_name" json:"computer_name,omitempty"`
	ProblemId        *string          `protobuf:"bytes,11,opt,name=problem_id" json:"problem_id,omitempty"`
	ProblemName      *string          `protobuf:"bytes,12,opt,name=problem_name" json:"problem_name,omitempty"`
	Submit           []*Ticket_Submit `protobuf:"bytes,13,rep,name=submit" json:"submit,omitempty"`
	Printer          *string          `protobuf:"bytes,14,opt,name=printer" json:"printer,omitempty"`
	XXX_unrecognized []byte           `json:"-"`
}

func (m *Ticket) Reset()         { *m = Ticket{} }
func (m *Ticket) String() string { return proto.CompactTextString(m) }
func (*Ticket) ProtoMessage()    {}

func (m *Ticket) GetContestId() uint32 {
	if m != nil && m.ContestId != nil {
		return *m.ContestId
	}
	return 0
}

func (m *Ticket) GetContestName() string {
	if m != nil && m.ContestName != nil {
		return *m.ContestName
	}
	return ""
}

func (m *Ticket) GetSubmitId() uint32 {
	if m != nil && m.SubmitId != nil {
		return *m.SubmitId
	}
	return 0
}

func (m *Ticket) GetJudgeTime() string {
	if m != nil && m.JudgeTime != nil {
		return *m.JudgeTime
	}
	return ""
}

func (m *Ticket) GetTeamId() uint32 {
	if m != nil && m.TeamId != nil {
		return *m.TeamId
	}
	return 0
}

func (m *Ticket) GetTeamName() string {
	if m != nil && m.TeamName != nil {
		return *m.TeamName
	}
	return ""
}

func (m *Ticket) GetLocationId() uint32 {
	if m != nil && m.LocationId != nil {
		return *m.LocationId
	}
	return 0
}

func (m *Ticket) GetLocationName() string {
	if m != nil && m.LocationName != nil {
		return *m.LocationName
	}
	return ""
}

func (m *Ticket) GetComputerId() string {
	if m != nil && m.ComputerId != nil {
		return *m.ComputerId
	}
	return ""
}

func (m *Ticket) GetComputerName() string {
	if m != nil && m.ComputerName != nil {
		return *m.ComputerName
	}
	return ""
}

func (m *Ticket) GetProblemId() string {
	if m != nil && m.ProblemId != nil {
		return *m.ProblemId
	}
	return ""
}

func (m *Ticket) GetProblemName() string {
	if m != nil && m.ProblemName != nil {
		return *m.ProblemName
	}
	return ""
}

func (m *Ticket) GetSubmit() []*Ticket_Submit {
	if m != nil {
		return m.Submit
	}
	return nil
}

func (m *Ticket) GetPrinter() string {
	if m != nil && m.Printer != nil {
		return *m.Printer
	}
	return ""
}

type Ticket_Submit struct {
	SubmitNumber     *uint32               `protobuf:"varint,1,opt,name=submit_number" json:"submit_number,omitempty"`
	Arrived          *uint32               `protobuf:"varint,2,opt,name=arrived" json:"arrived,omitempty"`
	Compiled         *bool                 `protobuf:"varint,3,opt,name=compiled" json:"compiled,omitempty"`
	School           *Ticket_Submit_School `protobuf:"bytes,4,opt,name=school" json:"school,omitempty"`
	Acm              *Ticket_Submit_ACM    `protobuf:"bytes,5,opt,name=acm" json:"acm,omitempty"`
	XXX_unrecognized []byte                `json:"-"`
}

func (m *Ticket_Submit) Reset()         { *m = Ticket_Submit{} }
func (m *Ticket_Submit) String() string { return proto.CompactTextString(m) }
func (*Ticket_Submit) ProtoMessage()    {}

func (m *Ticket_Submit) GetSubmitNumber() uint32 {
	if m != nil && m.SubmitNumber != nil {
		return *m.SubmitNumber
	}
	return 0
}

func (m *Ticket_Submit) GetArrived() uint32 {
	if m != nil && m.Arrived != nil {
		return *m.Arrived
	}
	return 0
}

func (m *Ticket_Submit) GetCompiled() bool {
	if m != nil && m.Compiled != nil {
		return *m.Compiled
	}
	return false
}

func (m *Ticket_Submit) GetSchool() *Ticket_Submit_School {
	if m != nil {
		return m.School
	}
	return nil
}

func (m *Ticket_Submit) GetAcm() *Ticket_Submit_ACM {
	if m != nil {
		return m.Acm
	}
	return nil
}

type Ticket_Submit_School struct {
	TestsTaken       *uint32 `protobuf:"varint,1,opt,name=tests_taken" json:"tests_taken,omitempty"`
	TestsPassed      *uint32 `protobuf:"varint,2,opt,name=tests_passed" json:"tests_passed,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *Ticket_Submit_School) Reset()         { *m = Ticket_Submit_School{} }
func (m *Ticket_Submit_School) String() string { return proto.CompactTextString(m) }
func (*Ticket_Submit_School) ProtoMessage()    {}

func (m *Ticket_Submit_School) GetTestsTaken() uint32 {
	if m != nil && m.TestsTaken != nil {
		return *m.TestsTaken
	}
	return 0
}

func (m *Ticket_Submit_School) GetTestsPassed() uint32 {
	if m != nil && m.TestsPassed != nil {
		return *m.TestsPassed
	}
	return 0
}

type Ticket_Submit_ACM struct {
	Result           *string `protobuf:"bytes,1,opt,name=result" json:"result,omitempty"`
	TestId           *uint32 `protobuf:"varint,2,opt,name=test_id" json:"test_id,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *Ticket_Submit_ACM) Reset()         { *m = Ticket_Submit_ACM{} }
func (m *Ticket_Submit_ACM) String() string { return proto.CompactTextString(m) }
func (*Ticket_Submit_ACM) ProtoMessage()    {}

func (m *Ticket_Submit_ACM) GetResult() string {
	if m != nil && m.Result != nil {
		return *m.Result
	}
	return ""
}

func (m *Ticket_Submit_ACM) GetTestId() uint32 {
	if m != nil && m.TestId != nil {
		return *m.TestId
	}
	return 0
}

func init() {
}
