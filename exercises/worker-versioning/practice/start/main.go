package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	pizza "temporal-versioning/exercises/worker-versioning/practice"

	"go.temporal.io/sdk/client"
)

func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	// TODO Part A: call c.UpdateWorkerBuildIdCompatibility() to inform the Task
	// Queue of your Build ID. You can also do this via the CLI if you are changing
	// a currently running workflow. An example of how to do it via the SDK is
	// below. Don't forget to change the BuildID to match your Worker.
	//
	// c.UpdateWorkerBuildIdCompatibility(context.Background(), &client.UpdateWorkerBuildIdCompatibilityOptions{
	// 	TaskQueue: pizza.TaskQueueName,
	// 	Operation: &client.BuildIDOpAddNewIDInNewDefaultSet{
	// 		BuildID: "revision-yymmdd",
	// 	},
	// })
	// **Note:** This code would usually only need to be run once. In a production
	// system you would not run this as part of your client, but more likely as part
	// of your build system on initial deployment, either via the SDK or the CLI.

	order := *createPizzaOrder()

	workflowID := fmt.Sprintf("pizza-workflow-order-%s", order.OrderNumber)

	options := client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: pizza.TaskQueueName,
	}

	we, err := c.ExecuteWorkflow(context.Background(), options, pizza.PizzaWorkflow, order)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}
	log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())

	var result pizza.OrderConfirmation
	err = we.Get(context.Background(), &result)
	if err != nil {
		log.Fatalln("Unable get workflow result", err)
	}

	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatalln("Unable to format order confirmation as JSON", err)
	}
	log.Printf("Workflow result: %s\n", string(data))
}

func createPizzaOrder() *pizza.PizzaOrder {
	customer := pizza.Customer{
		CustomerID: 12983,
		Name:       "María García",
		Email:      "maria1985@example.com",
		Phone:      "415-555-7418",
	}

	address := pizza.Address{
		Line1:      "701 Mission Street",
		Line2:      "Apartment 9C",
		City:       "San Francisco",
		State:      "CA",
		PostalCode: "94103",
	}

	p1 := pizza.Pizza{
		Description: "Large, with mushrooms and onions",
		Price:       1500,
	}

	p2 := pizza.Pizza{
		Description: "Small, with pepperoni",
		Price:       1200,
	}

	p3 := pizza.Pizza{
		Description: "Medium, with extra cheese",
		Price:       1300,
	}

	items := []pizza.Pizza{p1, p2, p3}

	order := pizza.PizzaOrder{
		OrderNumber: "Z1238",
		Customer:    customer,
		Items:       items,
		Address:     address,
		IsDelivery:  true,
	}

	return &order
}
