package order

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/google/uuid"
	"strconv"
)

type SqsAccess interface {
	Enqueue(or OrderRequest) error
	AckQueue(handle *string) error
	Dequeue() (or []OrderRequest,err error)
}

type sqsAccess struct {
	queueUrl string
	linger int64
	sqs sqsiface.SQSAPI
}

func (o *sqsAccess) Enqueue(or OrderRequest) error{
	var queueURL = o.queueUrl

	orders := make([]OrderRequest,0)
	for i := 0 ; i < or.NumberOfOrders ; i ++ {
		c := or
		c.Id = uuid.New().String()
		orders = append(orders,c)
	}

	split := split(orders,10)
	for _,bucket := range split {
		go func(b []OrderRequest){
			messages := make([]*sqs.SendMessageBatchRequestEntry,len(b))
			for i := 0 ; i < len(b) ; i ++ {
				arr,err := json.Marshal(b[i])
				if err != nil {
					fmt.Println("marshal : ",err)
					continue
				}
				var body = string(arr)
				copyIndex := strconv.Itoa(i)
				m := sqs.SendMessageBatchRequestEntry{
					MessageBody: &body,
					Id: &copyIndex,
				}
				messages[i] = &m
			}
			if _, err := o.sqs.SendMessageBatch(&sqs.SendMessageBatchInput{
				Entries: messages,
				QueueUrl:    &queueURL,
			}); err != nil {
				fmt.Println("sqs : ",err)
			}
		}(bucket)
	}

	return nil
}

func (o *sqsAccess) AckQueue(handle *string) error {
	var q = o.queueUrl
	_, err := o.sqs.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      &q,
		ReceiptHandle: handle,
	})
	if err != nil {
		return err
	}
	return nil
}

func (o *sqsAccess) Dequeue() (or []OrderRequest,err error){
	var q = o.queueUrl
	var l = o.linger
	result, err := o.sqs.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl: &q,
		MaxNumberOfMessages: aws.Int64(10),
		WaitTimeSeconds: &l,
	})
	if err != nil {
		return nil,err
	}
	var os = make([]OrderRequest,0,10)
	for _,ms := range result.Messages {
		var order OrderRequest
		if err := json.Unmarshal([]byte(*ms.Body),&order); err != nil {
			return nil,err
		}
		order.ReceiptHandle = ms.ReceiptHandle
		os = append(os, order)
	}
	return os,nil
}

func split(all []OrderRequest,bucketSize int) [][]OrderRequest {
	a := make([][]OrderRequest,0)
	for index,order := range all {
		bucket := index / bucketSize
		if len(a) < bucket+1 {
			a = append(a,make([]OrderRequest,0))
		}
		a[bucket] = append(a[bucket],order)
	}
	return a
}

func NewSqsAccess(sqs sqsiface.SQSAPI,queueUrl string,linger int64) SqsAccess {
	return &sqsAccess{queueUrl,linger,sqs}
}



