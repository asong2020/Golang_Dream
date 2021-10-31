package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"strings"
	"time"
)

const (
	KEY = "trace_id"
)

func NewRequestID() string {
	return strings.Replace(uuid.New().String(), "-", "", -1)
}

func NewContextWithTraceID() context.Context {
	ctx := context.WithValue(context.Background(), KEY,NewRequestID())
	return ctx
}

func NewContextWithTimeout(ctx context.Context) (context.Context,context.CancelFunc){
	ctx, cancel := context.WithTimeout(ctx, time.Second * 3)
	return ctx,cancel
}

func PrintLog(ctx context.Context, message string)  {
	fmt.Printf("%s|info|trace_id=%s|%s",time.Now().Format("2006-01-02 15:04:05") , GetContextValue(ctx, KEY), message)
}

func GetContextValue(ctx context.Context,k string)  string{
	v, ok := ctx.Value(k).(string)
	if !ok{
		return ""
	}
	return v
}

func ProcessEnter(ctx context.Context) {
	PrintLog(ctx, "Golang梦工厂")
}


func main()  {
	ctx, cancel := NewContextWithTimeout(NewContextWithTraceID())
	defer cancel()
	ProcessEnter(ctx)
}