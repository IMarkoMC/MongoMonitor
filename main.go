package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/iMarkoMC/MongoMonitor/pkg/pushover"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var pc pushover.PushoverClient

var offlineNodes = make(map[string]time.Time)

func main() {
	_, err := os.Stat(".env")

	if err == nil {
		log.Info("Found a .env file. Loading it")

		err := godotenv.Load()

		if err != nil {
			log.Fatalf("An error occurred while loading the env. Error %v", err)
		}
	}

	log.SetLevel(log.InfoLevel)

	if os.Getenv("DEBUG") == "true" {
		log.SetLevel(log.DebugLevel)
	}

	pc = pushover.New(os.Getenv("PUSHOVER_TOKEN"), os.Getenv("PUSHOVER_USER"))

	//TODO: Implement a autoreconnect in case the entire cluster goes down and/or send a notification that the cluster is offline
	serverMonitor := &event.ServerMonitor{}
	options := options.Client().ApplyURI(os.Getenv("MONGO_URI")).SetServerMonitor(serverMonitor)
	client, err := mongo.Connect(context.TODO(), options)

	if err != nil {
		log.Fatalf("An error occurred while connecting to mongo %v", err)
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	log.Info("Connected to MongoDB")

	serverMonitor.ServerDescriptionChanged = handleDescriptionChange

	//* keep the main thread alive
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
}

func handleDescriptionChange(e *event.ServerDescriptionChangedEvent) {
	log.Debugf("Event: %s", e)

	//* If the server was RSSecondary or RSPrimary and the new status is Unknown the node went down
	//TODO: Check if RSMember can be used insterad of this RSSecondary or RSPrimary
	if (e.PreviousDescription.Kind == description.RSSecondary || e.PreviousDescription.Kind == description.RSPrimary) && e.NewDescription.Kind == description.Unknown {
		log.Warnf("The node %s went offline", e.Address)

		msg := fmt.Sprintf("Address: %s\nRole: %s", e.Address, e.PreviousDescription.Kind.String())

		if e.NewDescription.LastError != nil {
			msg += "\nError: " + e.NewDescription.LastError.Error()
		}

		pc.SendNotification("A MongoDB node went offline!", msg)

		offlineNodes[string(e.Address)] = time.Now()

		return
	}

	//* If the previous state was Unknown and now it's Primary or secondary it's back online
	//TODO: Check if RSMember can be used insterad of this RSSecondary or RSPrimary
	if e.PreviousDescription.Kind == description.Unknown && (e.NewDescription.Kind == description.RSSecondary || e.NewDescription.Kind == description.RSPrimary) {
		log.Infof("The node %s is back online!", e.Address)

		msg := fmt.Sprintf("Address: %s\nRole: %s", e.Address, e.NewDescription.Kind.String())

		if offlineNodes[string(e.Address)] != (time.Time{}) {
			dt := time.Since(offlineNodes[string(e.Address)]).Round(time.Second)
			msg += "\nDowntime: " + dt.String()

			//* Remove the node from the map
			delete(offlineNodes, string(e.Address))
		}

		msg += "\nRTT: " + e.NewDescription.AverageRTT.Round(time.Millisecond).String()
		pc.SendNotification("A MongoDB node is online again!", msg)
	}
}
