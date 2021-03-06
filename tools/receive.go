package tools

import (
	"gopkg.in/stomp.v1"
)

type Conn interface {
	SendWithReceipt(destination, contentType string, body []byte, userDefined *stomp.Header) error
	Send(destination, contentType string, body []byte, userDefined *stomp.Header) error
}

func (pc *StompConfig) ReceiveLoop(queueName string, useTransactions bool, process func(Conn, *stomp.Message) error) error {
	conn, err := pc.NewConnection()
	//defer conn.Disconnect()

	if err != nil {
		return err
	}

	sub, err := conn.Subscribe(queueName, stomp.AckClientIndividual)
	if err != nil {
		return err
	}

	for {
		msg, err := sub.Read()
		if err != nil {
			return err
		}
		if useTransactions {
			tx := conn.Begin()
			err = process(tx, msg)
			if err != nil {
				tx.Abort()
			} else {
				tx.Ack(msg)
				err = tx.Commit()
			}
		} else {
			err = process(conn, msg)
			if err != nil {
				err = conn.Nack(msg)
			} else {
				err = conn.Ack(msg)
			}
		}
		if err != nil {
			return err
		}
	}
	return nil
}
