/*
This file is part of Cloud Native PostgreSQL.

Copyright (C) 2019-2021 EnterpriseDB Corporation.
*/

package controller

import (
	"github.com/spf13/cobra"
)

// NewCmd create a new cobra command
func NewCmd() *cobra.Command {
	var metricsAddr string
	var enableLeaderElection bool
	var configMapName string

	cmd := cobra.Command{
		Use: "controller [flags]",
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunController(metricsAddr, configMapName, enableLeaderElection)
		},
	}

	cmd.Flags().StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	cmd.Flags().BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. "+
			"If enabled, this will ensure there is only one active controller manager.")
	cmd.Flags().StringVar(&configMapName, "config-map-name", "", "The name of the ConfigMap containing "+
		"the operator configuration")

	return &cmd
}