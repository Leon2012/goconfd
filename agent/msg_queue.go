package agent

import (
	"github.com/Leon2012/goconfd/libs/kv"
	"github.com/Shopify/sysv_mq"
)

const (
	defaultMsgQueueMaxSize = 1024
)

type MsgQueue struct {
	mq  *sysv_mq.MessageQueue
	ctx *Context
}

func NewMsgQueue(ctx *Context, path string, projId int) (*MsgQueue, error) {
	mq, err := sysv_mq.NewMessageQueue(
		&sysv_mq.QueueConfig{
			Path:   path,
			ProjId: projId,
			//Key:     0xDEADBEEF,
			MaxSize: defaultMsgQueueMaxSize,
			Mode:    sysv_mq.IPC_CREAT | 0600,
		},
	)
	if err != nil {
		return nil, err
	}
	return &MsgQueue{
		mq:  mq,
		ctx: ctx,
	}, nil
}

func (q *MsgQueue) Receive() {
	for {
		select {
		case <-q.ctx.Agent.exitChan:
			return
		default:
			message, mtype, err := q.mq.ReceiveString(0, 0)
			q.ctx.Agent.logf("INFO: receive mq message - %s, type - %d", message, mtype)
			if err != nil {
				q.ctx.Agent.logf("ERROR: receive mq message (%s) error (%s)", message, err)
				continue
			}
			_, err = q.handle([]byte(message), mtype)
			if err != nil {
				q.ctx.Agent.logf("ERROR: handle mq message (%s) error (%s)", message, err)
				continue
			}
		}
	}
}

func (q *MsgQueue) SendString(message string) error {
	return q.mq.SendString(message, 1, 0)
}

func (q *MsgQueue) Send(message []byte) error {
	return q.mq.SendBytes(message, 1, 0)
}

func (q *MsgQueue) Close() {
	q.mq.Destroy()
}

func (q *MsgQueue) handle(message []byte, mtype int) ([]byte, error) {
	key := string(message)
	k, err := LoadValueByKey(q.ctx, key)
	if err != nil {
		return nil, err
	}
	data, err := kv.JsonEncode(k)
	if err != nil {
		return nil, err
	} else {
		return data, nil
	}
}
