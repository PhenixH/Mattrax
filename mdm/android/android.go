package android

import (
	"context"
	"fmt"
	"io/ioutil"

	pubsub2 "cloud.google.com/go/pubsub"
	"github.com/gorilla/mux"
	mattrax "github.com/mattrax/Mattrax/internal"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/androidmanagement/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"google.golang.org/api/pubsub/v1"
)

// Protocol contains the handlers for the Base protocol
type Protocol struct {
	srv          *mattrax.Server
	ams          *androidmanagement.Service
	amsProjectID string
	pubsubTopic  string
}

// TODO: Unused remove
func (p *Protocol) ID() string {
	return "android"
}

func (p *Protocol) Init(srv *mattrax.Server) (err error) {
	if srv.Args.GoogleServiceAccountPath == "" {
		return fmt.Errorf("error starting Android management service: no service account configured")
	}

	if p.ams, err = androidmanagement.NewService(context.Background(), option.WithCredentialsFile(srv.Args.GoogleServiceAccountPath)); err != nil {
		return fmt.Errorf("error creating android management API service: %w", err)
	}
	p.srv = srv

	credsFile, err := ioutil.ReadFile(srv.Args.GoogleServiceAccountPath)
	if err != nil {
		panic(err)
	}
	creds, err := google.CredentialsFromJSON(context.Background(), credsFile)
	if err != nil {
		panic(err)
	}
	p.amsProjectID = creds.ProjectID

	p.pubsubTopic = "projects/" + p.amsProjectID + "/topics/mattrax-alerts"
	pubsubService, err := pubsub.NewService(context.Background(), option.WithCredentialsFile(srv.Args.GoogleServiceAccountPath))
	if err != nil {
		panic(err)
	}

	if _, err := pubsubService.Projects.Topics.Get(p.pubsubTopic).Do(); err != nil && err.(*googleapi.Error).Code == 404 {
		log.Info().Msg("Creating pubsub Topic on Google Cloud Platform for management updates")
		if _, err := pubsubService.Projects.Topics.Create(p.pubsubTopic, &pubsub.Topic{}).Do(); err != nil {
			panic(err)
		}

		if _, err := pubsubService.Projects.Topics.SetIamPolicy(p.pubsubTopic, &pubsub.SetIamPolicyRequest{
			Policy: &pubsub.Policy{
				Bindings: []*pubsub.Binding{
					{
						Members: []string{"serviceAccount:android-cloud-policy@system.gserviceaccount.com"},
						Role:    "roles/pubsub.publisher",
					},
				},
			},
		}).Do(); err != nil {
			panic(err)
		}
	} else if err != nil {
		panic(err)
	}

	pubsub, err := pubsub2.NewClient(context.Background(), p.amsProjectID, option.WithCredentialsFile(srv.Args.GoogleServiceAccountPath))
	if err != nil {
		panic(err)
	}

	sub := pubsub.Subscription("mattrax-alerts")
	if exist, err := sub.Exists(context.Background()); err != nil {
		panic(err)
	} else if !exist {
		if _, err := pubsub.CreateSubscription(context.Background(), "mattrax-alerts", pubsub2.SubscriptionConfig{Topic: pubsub.Topic("mattrax-alerts")}); err != nil {
			panic(err)
		}
	}

	go func() {
		if err = sub.Receive(context.Background(), recieve(p)); err != nil {
			panic(err)
		}
	}()

	return nil
}

func (p *Protocol) Mount(r *mux.Router) error {
	return nil
}

type status struct {
	ManagementAPI bool `json:"management_api"`
}

func (p *Protocol) Status() (interface{}, error) {
	return status{
		ManagementAPI: true, // TODO
	}, nil
}
