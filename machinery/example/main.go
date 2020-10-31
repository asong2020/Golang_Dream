package main

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/tasks"
)

func main()  {

	cnf,err := config.NewFromYaml("./config.yml",false)
	if err != nil{
		log.Println("config failed",err)
		return
	}

	server,err := machinery.NewServer(cnf)
	if err != nil{
		log.Println("start server failed",err)
		return
	}

	// 注册任务
	err = server.RegisterTask("sum",Sum)
	if err != nil{
		log.Println("reg task failed",err)
		return
	}
	err = server.RegisterTask("call",CallBack)
	if err != nil{
		log.Println("reg task failed",err)
		return
	}


	worker := server.NewWorker("asong", 1)
	go func() {
		err = worker.Launch()
		if err != nil {
			log.Println("start worker error",err)
			return
		}
	}()

	//task signature
	signature1 := &tasks.Signature{
		Name: "sum",
		Args: []tasks.Arg{
			{
				Type:  "[]int64",
				Value: []int64{1,2,3,4,5,6,7,8,9,10},
			},
		},
		RetryTimeout: 100,
		RetryCount: 3,
	}

	signature2 := &tasks.Signature{
		Name: "sum",
		Args: []tasks.Arg{
			{
				Type:  "[]int64",
				Value: []int64{1,2,3,4,5,6,7,8,9,10},
			},
		},
		RetryTimeout: 100,
		RetryCount: 3,
	}

	signature3 := &tasks.Signature{
		Name: "sum",
		Args: []tasks.Arg{
			{
				Type:  "[]int64",
				Value: []int64{1,2,3,4,5,6,7,8,9,10},
			},
		},
		RetryTimeout: 100,
		RetryCount: 3,
	}



	//// group
	//group,err :=tasks.NewGroup(signature1,signature2,signature3)
	//if err != nil{
	//	log.Println("add group failed",err)
	//	return
	//}
	//
	//asyncResults, err :=server.SendGroupWithContext(context.Background(),group,0)
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
	//for _, asyncResult := range asyncResults{
	//	results,err := asyncResult.Get(1)
	//	if err != nil{
	//		log.Println(err)
	//		continue
	//	}
	//	log.Printf(
	//		"%v  %v  %v\n",
	//		asyncResult.Signature.Args[0].Value,
	//		tasks.HumanReadableResults(results),
	//	)
	//}
	callback := &tasks.Signature{
		Name: "call",
	}



	//group, err := tasks.NewGroup(signature1, signature2, signature3)
	//if err != nil {
	//
	//	log.Printf("Error creating group: %s", err.Error())
	//	return
	//}
	//
	//chord, err := tasks.NewChord(group, callback)
	//if err != nil {
	//	log.Printf("Error creating chord: %s", err)
	//	return
	//}
	//
	//chordAsyncResult, err := server.SendChordWithContext(context.Background(), chord, 0)
	//if err != nil {
	//	log.Printf("Could not send chord: %s", err.Error())
	//	return
	//}
	//
	//results, err := chordAsyncResult.Get(time.Duration(time.Millisecond * 5))
	//if err != nil {
	//	log.Printf("Getting chord result failed with error: %s", err.Error())
	//	return
	//}
	//log.Printf("%v\n", tasks.HumanReadableResults(results))



	//chain
	chain,err := tasks.NewChain(signature1,signature2,signature3,callback)
	if err != nil {

		log.Printf("Error creating group: %s", err.Error())
		return
	}
	chainAsyncResult, err := server.SendChainWithContext(context.Background(), chain)
	if err != nil {
		log.Printf("Could not send chain: %s", err.Error())
		return
	}

	results, err := chainAsyncResult.Get(time.Duration(time.Millisecond * 5))
	if err != nil {
		log.Printf("Getting chain result failed with error: %s", err.Error())
	}
	log.Printf(" %v\n", tasks.HumanReadableResults(results))



	////
	//eta := time.Now().UTC().Add(time.Second * 20)
	//signature.ETA = &eta

	//asyncResult, err := server.SendTaskWithContext(context.Background())
	//if err != nil {
	//	log.Println(err)
	//}
	//
	//
	//res, err := asyncResult.Get(1* time.Second)
	//if err != nil {
	//	log.Println(err)
	//}
	//log.Printf("get res is %v\n", tasks.HumanReadableResults(res))
}


func Sum(args []int64) (int64, error) {
	sum := int64(0)
	for _, arg := range args {
		sum += arg
	}

	//return sum, tasks.NewErrRetryTaskLater("我说他错了", 4 * time.Second)
	return sum,errors.New("我说他错了")
	//return sum,nil
}

// Multiply ...
func CallBack(args ...int64) (int64, error) {
	sum := int64(1)
	for _, arg := range args {
		sum *= arg
	}
	return sum, nil
}