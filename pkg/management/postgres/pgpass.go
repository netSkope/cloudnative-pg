/*
This file is part of Cloud Native PostgreSQL.

Copyright (C) 2019-2020 2ndQuadrant Italia SRL. Exclusively licensed to 2ndQuadrant Limited.
*/

package postgres

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"gitlab.2ndquadrant.com/k8s/cloud-native-postgresql/pkg/fileutils"
)

// CreatePgPass create and install a `.pgpass` file to allow
// connecting to all PostgreSQL server with password in the
// password file
func CreatePgPass(pwfile string) error {
	password, err := fileutils.ReadFile(pwfile)
	if err != nil {
		return err
	}

	// This is needed to create a normal connection
	pgpass := fmt.Sprintf(
		"%v:%v:%v:%v:%v\n",
		"*",        // host
		5432,       // port
		"postgres", // database name
		"postgres", // user name
		password)   // password

	// And this works for a replica connection
	pgpass += fmt.Sprintf(
		"%v:%v:%v:%v:%v\n",
		"*",           // host
		5432,          // port
		"replication", // database name
		"postgres",    // user name
		password)      // password

	if err = InstallPgPass(pgpass); err != nil {
		return err
	}

	return nil
}

// InstallPgPass install a pgpass file with the given content in the
// user home directory
func InstallPgPass(pgpassContent string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	targetPgPass := path.Join(homeDir, ".pgpass")
	if err = ioutil.WriteFile(targetPgPass, []byte(pgpassContent), 0600); err != nil {
		return err
	}

	return nil
}

// InstallCustomConfigurationFile install in the PgData the file with
// the PostgreSQL settings which have been generated by the operator
func InstallCustomConfigurationFile(pgdata, filename string) error {
	return installPgDataFile(pgdata, filename, PostgresqlCustomConfigurationFile)
}

// InstallPgHBAFile install in the PgData the file containing
// the host-based access rules which is being managed by the operator
func InstallPgHBAFile(pgdata, filename string) error {
	return installPgDataFile(pgdata, filename, "pg_hba.conf")
}

// installPgDataFile install a file in PgData
func installPgDataFile(pgdata, source, destination string) error {
	targetFile := path.Join(pgdata, destination)
	return fileutils.CopyFile(source, targetFile)
}