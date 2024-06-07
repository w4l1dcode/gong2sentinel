package main

import (
	"context"
	"flag"
	"fmt"
	_ "fmt"
	"github.com/sirupsen/logrus"
	"gong2sentinel/config"
	"gong2sentinel/pkg/gong/auditing"
	"gong2sentinel/pkg/gong/calls"
	msSentinel "gong2sentinel/pkg/sentinel"
	_ "io/ioutil"
	"log"
	_ "net/http"
	"sync"
)

func main() {
	ctx := context.Background()

	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	confFile := flag.String("config", "gong2sentinel_config.yml", "The YAML configuration file.")
	flag.Parse()

	conf := config.Config{}
	if err := conf.Load(*confFile); err != nil {
		logger.WithError(err).WithField("config", *confFile).Fatal("failed to load configuration")
	}

	if err := conf.Validate(); err != nil {
		logger.WithError(err).WithField("config", *confFile).Fatal("invalid configuration")
	}

	logrusLevel, err := logrus.ParseLevel(conf.Log.Level)
	if err != nil {
		logger.WithError(err).Error("invalid log level provided")
		logrusLevel = logrus.InfoLevel
	}
	logger.WithField("level", logrusLevel.String()).Info("set log level")
	logger.SetLevel(logrusLevel)

	// ---

	errors := make(chan error)
	wg := &sync.WaitGroup{}

	var allGongAuditLogs []map[string]string
	var allGongUserAccessLogs []map[string]string

	wg.Add(1)
	go func() {
		logger.Info("retrieving gong audit logs")

		var audErr error
		allGongAuditLogs, audErr = auditing.GetAuditLogs(conf.Gong.AccessKey, conf.Gong.AccessSecret, conf.Gong.LookupHours)
		if audErr != nil {
			//log.Fatalf("failed to retrieve Gong Audit Logs: %v", audErr)
			errors <- fmt.Errorf("failed to retrieve Gong Audit Logs: %v", audErr)
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		logger.Info("retrieving gong user access logs")

		// Get call IDs
		callIds, err := calls.GetCallIDs(conf.Gong.AccessKey, conf.Gong.AccessSecret)
		if err != nil {
			log.Fatalf("failed to retrieve call IDs: %v", err)
		}

		// Get user access logs
		allGongUserAccessLogs, err = calls.GetUserAccess(conf.Gong.AccessKey, conf.Gong.AccessSecret, callIds)
		if err != nil {
			log.Fatalf("failed to retrieve Gong User Access Logs: %v", err)
		}

		wg.Done()
	}()

	// ---

	doneChan := make(chan struct{})
	go func() {
		wg.Wait()
		close(doneChan)
	}()

	logger.Info("waiting for log retrieval to finish")
	select {
	case err := <-errors:
		logger.WithError(err).Fatal("failed to retrieve logs")
	case <-doneChan:
		logger.Info("finished retrieving logs")
	}

	// ---

	sentinel, err := msSentinel.New(logger, msSentinel.Credentials{
		TenantID:       conf.Microsoft.TenantID,
		ClientID:       conf.Microsoft.AppID,
		ClientSecret:   conf.Microsoft.SecretKey,
		SubscriptionID: conf.Microsoft.SubscriptionID,
	})
	if err != nil {
		logger.WithError(err).Fatal("could not create MS Sentinel client")
	}

	// ---
	errors2 := make(chan error)
	wg2 := &sync.WaitGroup{}

	wg2.Add(1)
	go func() {

		logger.WithField("total", len(allGongAuditLogs)).Info("shipping off Gong audit logs to Sentinel")

		if err := sentinel.SendLogs(ctx, logger,
			conf.Microsoft.DataCollection.Endpoint,
			conf.Microsoft.DataCollection.RuleID,
			conf.Microsoft.DataCollection.StreamNameAuditing,
			allGongAuditLogs); err != nil {
			errors2 <- fmt.Errorf("could not ship audit logs to sentinel: %v", err)
		}

		logger.WithField("total", len(allGongAuditLogs)).Info("successfully sent Gong Auditing logs to sentinel")
		wg2.Done()
	}()

	// ---

	wg2.Add(1)
	go func() {
		logger.WithField("total", len(allGongUserAccessLogs)).Info("shipping off Gong user access logs to Sentinel")

		if err := sentinel.SendLogs(ctx, logger,
			conf.Microsoft.DataCollection.Endpoint,
			conf.Microsoft.DataCollection.RuleID,
			conf.Microsoft.DataCollection.StreamNameCallUserAccess,
			allGongUserAccessLogs); err != nil {
			errors2 <- fmt.Errorf("could not ship access logs to sentinel: %v", err)

		}

		logger.WithField("total", len(allGongUserAccessLogs)).Info("successfully sent Gong User Access logs to sentinel")
		wg2.Done()
	}()

	doneChan2 := make(chan struct{})
	go func() {
		wg2.Wait()
		close(doneChan2)
	}()

	logger.Info("waiting for log ingestion to finish")
	select {
	case err := <-errors2:
		logger.WithError(err).Fatal("failed to ingest logs")
	case <-doneChan2:
		logger.Info("finished ingesting logs")
	}

}
