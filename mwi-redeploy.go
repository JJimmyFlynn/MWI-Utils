package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/charmbracelet/huh"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func setupClient() *ecs.Client {
	awsConfig, err := config.LoadDefaultConfig(context.Background())

	if err != nil {
		log.Printf("unable to load AWS credentials: %s", err)
		return nil
	}

	return ecs.NewFromConfig(awsConfig)
}
func promptForCluster() string {
	var cluster string

	err := huh.NewSelect[string]().
		Title("Select Environment").
		Options(
			huh.NewOption("Integration", "integration"),
			huh.NewOption("Staging", "staging"),
			huh.NewOption("Production", "production"),
		).
		Value(&cluster).
		Run()

	if err != nil {
		log.Fatalln(err)
	}

	return getCluster(cluster)
}

func getCluster(environment string) string {
	clusterNameMap := map[string]string{
		"integration": "mwi-int-ecs",
		"staging":     "mwi-stg-ecs",
		"production":  "mwi-prd-ecs",
	}

	clusterName, ok := clusterNameMap[environment]

	if ok == false {
		log.Fatalln("Invalid environment. Allowed values are: integration, staging, production")
	}

	return clusterName
}

func cmdRedeployCluster(ctx *cli.Context) error {
	var cluster string
	client := setupClient()

	if ctx.String("environment") != "" {
		cluster = getCluster(ctx.String("environment"))
	}

	if ctx.String("environment") == "" {
		cluster = promptForCluster()
	}

	response, err := client.ListServices(context.Background(), &ecs.ListServicesInput{
		Cluster: aws.String(cluster),
	})

	if err != nil {
		log.Println("Unable to reach the ECS service. Are you using the proper AWS profile and credentials?")
		log.Fatalln(err)
	}

	for _, service := range response.ServiceArns {

		fmt.Printf("Initiating redeploy of %s service\n", service)

		_, err := client.UpdateService(context.Background(), &ecs.UpdateServiceInput{
			Service:            aws.String(service),
			Cluster:            aws.String(cluster),
			ForceNewDeployment: true,
		})

		if err != nil {
			log.Printf("unable to initiate new deploy of %s\n", service)
			log.Fatal(err)
		}

	}

	return nil
}

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "environment",
				Aliases: []string{"e"},
				Usage:   "The environment to redeploy the services in: integration, staging, production",
			},
		},
		Action: cmdRedeployCluster,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
