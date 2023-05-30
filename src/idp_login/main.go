package main

import (
    "fmt"
    "log"
    "os"
	"os/user"
	"database/sql"
    "encoding/json"
	"bytes"

    "github.com/urfave/cli/v2"
	_ "github.com/mattn/go-sqlite3"
)

const DATABASE_PATH = "/etc/project_1.sqlite"

func isAdministrator() bool {
	usr, err := user.Current()
	if err != nil {
		fmt.Println("[!] Error: could not get current user")
		return false
	}

	group, err := user.LookupGroup("idpadmins")
	if err != nil {
		fmt.Println("[!] Error: could not get idpadmins group information")
		return false
	}

	userGroups, err := usr.GroupIds()
	if err != nil {
		fmt.Println("[!] Error: could not get user groups")
		return false
	}

	for _, userGroup := range userGroups {
		if userGroup == group.Gid {
			return true
		}
	}

	return false
}

func check_requirements() {
	fmt.Println("[*] Checking requirements...")
	if os.Geteuid() != 0 {
		fmt.Println("[!] Error: program must be run with the 4750 permissions")
		return
	}

	if _, err := os.Stat(DATABASE_PATH); os.IsNotExist(err) {
		fmt.Println("[-] Database does not exist, creating it...")

		db, err := sql.Open("sqlite3", DATABASE_PATH)
		if err != nil {
			fmt.Println("[!] Error: could not create database")
			return
		}
		defer db.Close()

		_, err = db.Exec("CREATE TABLE idps (name TEXT PRIMARY KEY, params TEXT)")
		if err != nil {
			fmt.Println("[!] Error: could not create idps table")
			return
		}

		_, err = db.Exec("CREATE TABLE attributes (username TEXT, idp TEXT, attributes TEXT, PRIMARY KEY (username, idp))")
		if err != nil {
			fmt.Println("[!] Error: could not create attributes table")
			return
		}
	}

	fmt.Println("[+] All requirements met")
}

func idp_exists(idp string) bool {
	db, err := sql.Open("sqlite3", DATABASE_PATH)
	if err != nil {
		fmt.Println("[!] Error: could not open database")
		return false
	}

	defer db.Close()

	rows, err := db.Query("SELECT name FROM idps WHERE name = ?", idp)
	if err != nil {
		fmt.Println("[!] Error: could not query database")
		return false
	}

	defer rows.Close()

	if rows.Next() {
		return true
	}

	return false
}

func manage_idp(operation string, idp string, params string) {
	if operation == "set" {
		if params == "" {
			fmt.Println("[!] Error: params cannot be empty for set operation")
			return
		}

		if idp_exists(idp) {
			fmt.Println("[!] Error: IdP already exists")
			return
		}

		db, err := sql.Open("sqlite3", DATABASE_PATH)
		if err != nil {
			fmt.Println("[!] Error: could not open database")
			return
		}

		defer db.Close()

		_, err = db.Exec("INSERT INTO idps (name, params) VALUES (?, ?)", idp, params)
		if err != nil {
			fmt.Println("[!] Error: could not insert IdP into database")
			return
		}

		fmt.Println("[+] IdP successfully added")
	} else if operation == "change" {
		if params == "" {
			fmt.Println("[!] Error: params cannot be empty for change operation")
			return
		}

		if !idp_exists(idp) {
			fmt.Println("[!] Error: IdP does not exist")
			return
		}

		db, err := sql.Open("sqlite3", DATABASE_PATH)
		if err != nil {
			fmt.Println("[!] Error: could not open database")
			return
		}

		defer db.Close()

		_, err = db.Exec("UPDATE idps SET params = ? WHERE name = ?", params, idp)
		if err != nil {
			fmt.Println("[!] Error: could not update IdP in database")
			return
		}

		fmt.Println("[+] IdP successfully updated")
	} else if operation == "delete" {
		if !idp_exists(idp) {
			fmt.Println("[!] Error: IdP does not exist")
			return
		}

		db, err := sql.Open("sqlite3", DATABASE_PATH)
		if err != nil {
			fmt.Println("[!] Error: could not open database")
			return
		}

		defer db.Close()

		_, err = db.Exec("DELETE FROM idps WHERE name = ?", idp)
		if err != nil {
			fmt.Println("[!] Error: could not delete IdP from database")
			return
		}

		fmt.Println("[+] IdP successfully deleted")
	}
}

func getCurrentUser() string {
	user, err := user.Current()
	if err != nil {
		fmt.Println("[!] Error: could not get current user")
		return ""
	}

	return user.Username
}

func list_available_idps() {
	db, err := sql.Open("sqlite3", DATABASE_PATH)
	if err != nil {
		fmt.Println("[!] Error: could not open database")
		return
	}

	defer db.Close()

	rows, err := db.Query("SELECT name FROM idps")
	if err != nil {
		fmt.Println("[!] Error: could not query database")
		return
	}

	defer rows.Close()

	fmt.Println("[+] Available IdPs:")
	for rows.Next() {
		var name string
		rows.Scan(&name)
		fmt.Println("    - " + name)
	}
}

func list_users() {
	db, err := sql.Open("sqlite3", DATABASE_PATH)
	if err != nil {
		fmt.Println("[!] Error: could not open database")
		return
	}

	defer db.Close()

	rows, err := db.Query("SELECT username FROM attributes")
	if err != nil {
		fmt.Println("[!] Error: could not query database")
		return
	}

	defer rows.Close()

	fmt.Println("[+] Users:")
	for rows.Next() {
		var username string
		rows.Scan(&username)
		fmt.Println("    - " + username)
	}
}

func list_idps(username string) {
	db, err := sql.Open("sqlite3", DATABASE_PATH)
	if err != nil {
		fmt.Println("[!] Error: could not open database")
		return
	}

	defer db.Close()

	rows, err := db.Query("SELECT idp, attributes FROM attributes WHERE username = ?", username)
	if err != nil {
		fmt.Println("[!] Error: could not query database")
		return
	}

	defer rows.Close()

	for rows.Next() {
		var idp string
		var attributes string
		var out bytes.Buffer

		rows.Scan(&idp, &attributes)

		err := json.Indent(&out, []byte(attributes), "", "\t")
		if err != nil {
			fmt.Println("[!] Error: could not indent JSON")
			fmt.Println("[+] IdP name: " + idp)
			fmt.Println("[+] IdP params:")
			fmt.Println(attributes)
			return
		}

		fmt.Println("[+] IdP name: " + idp)
		fmt.Println("[+] IdP attributes:")
		fmt.Println(out.String())
	}
}

func print_attributes(idp string) {
	if !idp_exists(idp) {
		fmt.Println("[!] Error: IdP does not exist")
		return
	}

	db, err := sql.Open("sqlite3", DATABASE_PATH)
	if err != nil {
		fmt.Println("[!] Error: could not open database")
		return
	}

	defer db.Close()

	rows, err := db.Query("SELECT params FROM idps WHERE name = ?", idp)
	if err != nil {
		fmt.Println("[!] Error: could not query database")
		return
	}

	defer rows.Close()

	for rows.Next() {
		var params string
		rows.Scan(&params)
		var out bytes.Buffer
		err := json.Indent(&out, []byte(params), "", "\t")
		if err != nil {
			fmt.Println("[!] Error: could not indent JSON")
			fmt.Println("[+] IdP params:")
			fmt.Println(params)
			return
		}

		fmt.Println("[+] IdP params:")
		fmt.Println(out.String())
	}
}

func manage_attributes(username string, operation string, idp string, attributes string) {
	if !idp_exists(idp) {
		fmt.Println("[!] Error: IdP does not exist")
		return
	}

	if attributes == "" && operation != "delete" {
		fmt.Println("[!] Error: attributes cannot be empty")
		return
	}

	if operation != "set" && operation != "change" && operation != "delete" {
		fmt.Println("[!] Error: invalid operation")
		return
	}

	if operation == "set" {
		db, err := sql.Open("sqlite3", DATABASE_PATH)
		if err != nil {
			fmt.Println("[!] Error: could not open database")
			return
		}

		defer db.Close()

		_, err = db.Exec("INSERT INTO attributes (username, idp, attributes) VALUES (?, ?, ?)", username, idp, attributes)
		if err != nil {
			fmt.Println("[!] Error: could not insert attributes into database")
			return
		}

		fmt.Println("[+] Attributes successfully added")
	} else if operation == "change" {
		if attributes == "" {
			fmt.Println("[!] Error: attributes cannot be empty for change operation")
			return
		}

		db, err := sql.Open("sqlite3", DATABASE_PATH)
		if err != nil {
			fmt.Println("[!] Error: could not open database")
			return
		}

		defer db.Close()

		_, err = db.Exec("UPDATE attributes SET attributes = ? WHERE username = ? AND idp = ?", attributes, username, idp)
		if err != nil {
			fmt.Println("[!] Error: could not update attributes in database")
			return
		}

		fmt.Println("[+] Attributes successfully updated")
	} else if operation == "delete" {
		db, err := sql.Open("sqlite3", DATABASE_PATH)
		if err != nil {
			fmt.Println("[!] Error: could not open database")
			return
		}

		defer db.Close()

		_, err = db.Exec("DELETE FROM attributes WHERE username = ? AND idp = ?", username, idp)
		if err != nil {
			fmt.Println("[!] Error: could not delete attributes from database")
			return
		}

		fmt.Println("[+] Attributes successfully deleted")
	}
}

func main() {
	check_requirements()

    app := &cli.App{
        Name:  "idp_login",
        Usage: "Manage Identity Providers (IdPs) and identity attributes for users",
        
		Commands: []*cli.Command{
			{
				Name:    	 "manage-idp",
				Usage:       "manage-idp [--operation set|change|delete] [--idp IDP_NAME] [--params PARAMS]",
				Description: "Set, change or delete operational parameters for a given IdP, only users belonging to the idpadmins group can perform this operation",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "operation",
						Aliases: []string{"o"},
						Usage:   "Operation to perform (set, change, delete)",
						Required: true,
					},
					&cli.StringFlag{
						Name:    "idp",
						Aliases: []string{"i"},
						Usage:   "IdP name",
						Required: true,
					},
					&cli.StringFlag{
						Name:    "params",
						Aliases: []string{"p"},
						Usage:   "IdP operational parameters for the IdP",
					},
				},
				Action: func(c *cli.Context) error {
					if !isAdministrator() {
						fmt.Println("[!] Error: current user is not an administrator")
						return nil
					}

					operation := c.String("operation")
					if operation != "set" && operation != "change" && operation != "delete" {
						fmt.Println("[!] Error: invalid operation")
					} else {
						idp := c.String("idp")
						
						params := c.String("params")

						manage_idp(operation, idp, params)

					}
					return nil
				},
			},
			{
				Name:    	 "manage-attributes",
				Usage:   	 "manage-attributes [--operation set|change|delete] [--idp IDP_NAME] [--attributes ATTRIBUTES]",
				Description: "Set, change or delete identity attributes for a given IdP, the changes are applied only to the current user",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "operation",
						Aliases: []string{"o"},
						Usage:   "Operation to perform (set, change, delete)",
						Required: true,
					},
					&cli.StringFlag{
						Name:    "idp",
						Aliases: []string{"i"},
						Usage:   "IdP name",
					},
					&cli.StringFlag{
						Name:    "attributes",
						Aliases: []string{"a"},
						Usage:   "Identity attributes for the IdP",
					},
				},
				Action: func(c *cli.Context) error {
					operation := c.String("operation")
					if operation != "set" && operation != "change" && operation != "delete" {
						fmt.Println("[!] Error: invalid operation")
					} else {
						username := getCurrentUser()

						idp := c.String("idp")
						if idp == "" {
							list_available_idps()
						} else {
							attributes := c.String("attributes")

							if attributes == "" && operation != "delete" {
								print_attributes(idp)
							} else {
								manage_attributes(username, operation, idp, attributes)
							}
						}
					}

					return nil
				},
			},
			{
				Name:    	 "list-users",
				Usage:   	 "list-users",
				Description: "List all users with registered IdPs, only users belonging to the idpadmins group can perform this operation",
				Action: func(c *cli.Context) error {
					if !isAdministrator() {
						fmt.Println("[!] Error: current user is not an administrator")
						return nil
					}

					list_users()

					return nil
				},
			},
			{
				Name:    	 "list-idps",
				Usage:   	 "list-idps",
				Description: "List all registered IdPs, only for the current user",
				Action: func(c *cli.Context) error {
					username := getCurrentUser()

					list_idps(username)

					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}