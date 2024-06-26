package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Rajan-226/client-go-practice/pkg/clients/kubernetes"
	"github.com/Rajan-226/client-go-practice/pkg/utils"
)

func main() {
	ctx := context.Background()
	if err := kubernetes.Init(); err != nil {
		log.Fatal(err.Error())
	}

	if err := utils.PrintPods(ctx); err != nil {
		log.Fatal(err.Error())
	}

	// if err := utils.EditDeploymentImageTag(ctx, "nginx", "1.21.6"); err != nil {
	// 	log.Fatal(err.Error())
	// }

	// if err := utils.ListResources(ctx, "ns-one", "pods"); err != nil {
	// 	log.Fatal(err.Error())
	// }

	err, stopCh := utils.DeploymentController(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}

	// This will block the program
	<-stopCh

	fmt.Println("Program stopped")
}
