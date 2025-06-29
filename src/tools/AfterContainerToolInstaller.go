package tools

import (
	"ContainDB/src/Docker"
	"fmt"
)

// AfterContainerToolInstaller provides post-installation setup for database management tools.
//
// For MySQL/MariaDB, it offers to install or reinstall phpMyAdmin. If phpMyAdmin is
// already running, it asks if the user wants to reinstall it.
//
// For MongoDB, it offers to install MongoDB Compass as a GUI management tool.
//
// For PostgreSQL, it offers to install PgAdmin as a GUI management tool.
//
// Parameters:
//   - database: A string identifying the database type ("mysql", "mariadb", "mongodb", or "postgresql")
//
// The function doesn't return any values but initiates the installation of
// the respective management tool based on user consent.
func AfterContainerToolInstaller(database string) {
	switch database {
	case "mysql", "mariadb":
		// Check if phpMyAdmin is already running
		if Docker.IsContainerRunning("phpmyadmin", true) {
			fmt.Println("phpMyAdmin is already running.")
			consentPhpMyAdmin := Docker.AskYesNo("Do you want to reinstall phpMyAdmin for this database?")
			if consentPhpMyAdmin {
				StartPHPMyAdmin()
			} else {
				fmt.Println("You can reinstall phpMyAdmin later using the 'phpmyadmin' option.")
			}
		} else {
			consentPhpMyAdmin := Docker.AskYesNo("Do you want to install phpMyAdmin for this database?")
			if consentPhpMyAdmin {
				StartPHPMyAdmin()
			} else {
				fmt.Println("You can install phpMyAdmin later using the 'phpmyadmin' option.")
			}
		}
	case "mongodb":
		consentCompass := Docker.AskYesNo("Do you want to install MongoDB Compass?")
		if consentCompass {
			DownloadMongoDBCompass()
		} else {
			fmt.Println("You can install MongoDB Compass later using the 'mongodb compass' option.")
		}
	case "postgresql":
		pgAdminConsent := Docker.AskYesNo("Do you want to install PgAdmin? (yes/no)")
		if pgAdminConsent {
			StartPgAdmin()
		}
	}
}