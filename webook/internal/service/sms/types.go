package sms

import "context"

type Service interface {
	Send(ctx context.Context, tplId string, args []string, number ...string) error
}
