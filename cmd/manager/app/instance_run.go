/*
This file is part of Cloud Native PostgreSQL.

Copyright (C) 2019-2020 2ndQuadrant Italia SRL. Exclusively licensed to 2ndQuadrant Limited.
*/

package app

import (
	"context"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/2ndquadrant/cloud-native-postgresql/internal/management/controller"
	"github.com/2ndquadrant/cloud-native-postgresql/pkg/management/log"
	"github.com/2ndquadrant/cloud-native-postgresql/pkg/management/postgres/webserver"
)

var (
	postgresCommand *exec.Cmd
	reconciler      *controller.InstanceReconciler
)

func runSubCommand() {
	var err error

	reconciler, err = controller.NewInstanceReconciler(&instance)
	if err != nil {
		log.Log.Error(err, "Error while starting reconciler")
		return
	}

	err = reconciler.VerifyPgDataCoherence(context.Background())
	if err != nil {
		log.Log.Error(err, "Error while checking Kubernetes cluster status")
		return
	}

	startWebServer()
	startReconciler()
	registerSignalHandler()

	// Print the content of PostgreSQL control data, for debugging and tracing
	pgControlData := exec.Command("pg_controldata")
	pgControlData.Env = os.Environ()
	pgControlData.Env = append(pgControlData.Env, "PGDATA="+instance.PgData)
	pgControlData.Stdout = os.Stdout
	pgControlData.Stderr = os.Stderr
	err = pgControlData.Run()

	if err != nil {
		log.Log.Error(err, "Error printing the control information of this PostgreSQL instance")
		return
	}

	postgresCommand, err = instance.Run()
	if err != nil {
		log.Log.Error(err, "Unable to start PostgreSQL up")
	}

	if err = postgresCommand.Wait(); err != nil {
		log.Log.Error(err, "PostgreSQL exited with errors")
	}
}

// startWebServer start the web server for handling probes given
// a certain PostgreSQL instance
func startWebServer() {
	go func() {
		err := webserver.ListenAndServe(&instance)
		if err != nil {
			log.Log.Error(err, "Error while starting the web server")
		}
	}()
}

// startReconciler start the reconciliation loop
func startReconciler() {
	go reconciler.Run()
}

// registerSignalHandler handles signals from k8s, notifying postgres as
// needed
func registerSignalHandler() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-signals
		log.Log.Info("Received termination signal", "signal", sig)

		log.Log.Info("Shutting down web server")
		err := webserver.Shutdown()
		if err != nil {
			log.Log.Error(err, "Error while shutting down the web server")
		} else {
			log.Log.Info("Web server shut down")
		}

		log.Log.Info("Shutting down controller")
		reconciler.Stop()

		if postgresCommand != nil {
			log.Log.Info("Shutting down PostgreSQL instance")
			err := postgresCommand.Process.Signal(syscall.SIGTERM)
			if err != nil {
				log.Log.Error(err, "Unable to send SIGTERM to PostgreSQL instance")
			}
		}
	}()
}