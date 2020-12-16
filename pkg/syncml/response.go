package syncml

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/mattrax/xml"
)

// Response is a SyncML response body. It has helpers to make generating responses easier
type Response struct {
	res Message
}

// Set creates a generic command on the response
func (r *Response) Set(command, cmdID, uri string, meta *Meta, data interface{}) {
	r.res.Body.Commands = append(r.res.Body.Commands, Command{
		XMLName: xml.Name{
			Local: command,
		},
		CmdID: cmdID,
		Body: []Command{
			{
				XMLName: xml.Name{
					Local: "Item",
				},
				Target: &LocURI{
					URI: uri,
				},
				Meta: meta,
				Data: fmt.Sprintf("%v", data),
			},
		},
	})
}

// Respond creates the final element and encodes the response
func (r Response) Respond(w http.ResponseWriter) error {
	r.res.Body.Final = "<Final />"
	w.Header().Set("Content-Type", "application/vnd.syncml.dm+xml")
	if err := xml.NewEncoder(w).Encode(r.res); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}
	return nil
}

// SetStatus changes the SyncML body's status
func (r *Response) SetStatus(status int) {
	r.res.Body.Commands[0].Data = fmt.Sprintf("%v", status)
}

// FinalStatus returns the SyncML body's status
func (r *Response) FinalStatus() int32 {
	n, err := strconv.Atoi(r.res.Body.Commands[0].Data)
	if err != nil {
		return -1
	}
	return int32(n)
}

// NewResponse creates a new SyncML Envelope for the response
func NewResponse(cmd Message) Response {
	return Response{
		res: Message{
			XmlnA: "syncml:metinf",
			Header: Header{
				VerDTD:         cmd.Header.VerDTD,
				VerProto:       cmd.Header.VerProto,
				SessionID:      cmd.Header.SessionID,
				MsgID:          cmd.Header.MsgID,
				TargetURI:      cmd.Header.SourceURI,
				SourceURI:      cmd.Header.TargetURI,
				MetaMaxMsgSize: MaxRequestBodySize,
			},
			Body: Body{
				Commands: []Command{
					{
						XMLName: xml.Name{
							Local: "Status",
						},
						CmdID:  "1",
						MsgRef: cmd.Header.MsgID,
						CmdRef: "0",
						Cmd:    "SyncHdr",
						Data:   fmt.Sprintf("%v", StatusOK),
					},
				},
			},
		},
	}
}

// NewBlankResponse creates a new empty SyncML Envelope with only a Status element
func NewBlankResponse(cmd Message, status int) Response {
	return Response{
		res: Message{
			XmlnA: "syncml:metinf",
			Header: Header{
				VerDTD:         cmd.Header.VerDTD,
				VerProto:       cmd.Header.VerProto,
				SessionID:      cmd.Header.SessionID,
				MsgID:          cmd.Header.MsgID,
				TargetURI:      cmd.Header.SourceURI,
				SourceURI:      cmd.Header.TargetURI,
				MetaMaxMsgSize: MaxRequestBodySize,
			},
			Body: Body{
				Commands: []Command{
					{
						XMLName: xml.Name{
							Local: "Status",
						},
						CmdID:  "1",
						MsgRef: cmd.Header.MsgID,
						CmdRef: "0",
						Cmd:    "SyncHdr",
						Data:   fmt.Sprintf("%v", status),
					},
				},
			},
		},
	}
}
