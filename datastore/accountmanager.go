// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2017-2018 Canonical Ltd
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package datastore

import (
	"database/sql"
	"errors"
	"log"
)

const createAccountTableSQL = `
	CREATE TABLE IF NOT EXISTS account (
		id            serial primary key not null,
		authority_id  varchar(200) not null unique,
		assertion     text default '',
		resellerapi   bool default false,
		api_key       varchar(200) not null
	)
`

const createAccountSQL = "INSERT INTO account (authority_id, assertion, resellerapi, api_key) VALUES ($1,$2,$3,$4)"
const listAccountsSQL = "select id, authority_id, assertion, resellerapi, api_key from account order by authority_id"
const getAccountSQL = "select id, authority_id, assertion, resellerapi, api_key from account where authority_id=$1"

const getAccountByIDSQL = "select id, authority_id, assertion, resellerapi, api_key from account where id=$1"
const getUserAccountByIDSQL = `
	select a.id, authority_id, assertion, resellerapi, api_key 
	from account a
	inner join useraccountlink l on a.id = l.account_id
	inner join userinfo u on l.user_id = u.id
	where id=$1 and u.username=$2`

const updateAccountSQL = "UPDATE account SET authority_id=$2, assertion=$3, resellerapi=$4, api_key=$5 WHERE id=$1"
const updateUserAccountSQL = `
	UPDATE account a
	SET authority_id=$3, assertion=$4, resellerapi=$5, api_key=$6
	INNER JOIN useraccountlink l on a.id = l.account_id
	INNER JOIN userinfo u on l.user_id = u.id
	WHERE id=$1 AND u.username=$2
`
const upsertAccountSQL = `
	WITH upsert AS (
		update account set authority_id=$1, assertion=$2
		where authority_id=$1
		RETURNING *
	)
	insert into account (authority_id,assertion)
	select $1, $2
	where not exists (select * from upsert)
`

const listUserAccountsSQL = `
	select a.id, authority_id, assertion, resellerapi, api_key 
	from account a
	inner join useraccountlink l on a.id = l.account_id
	inner join userinfo u on l.user_id = u.id
	where u.username=$1
`

const listNotUserAccountsSQL = `
	select id, authority_id, assertion, resellerapi, api_key 
	from account
	where id not in (
		select a.id 
		from account a
		inner join useraccountlink l on a.id = l.account_id
		inner join userinfo u on l.user_id = u.id
		where u.username=$1
	)
`

// Add the reseller API field to indicate whether the reseller functions are available for an account
const alterAccountResellerAPI = "ALTER TABLE account ADD COLUMN resellerapi bool DEFAULT false"

// Add the API key field to the account table (nullable)
const alterAccountAPIKey = "ALTER TABLE account ADD COLUMN api_key varchar(200) DEFAULT ''"

// Make the API key not-nullable
const alterAccountAPIKeyNotNullable = `
	ALTER TABLE account
	ALTER COLUMN api_key SET NOT null,
	ALTER COLUMN api_key DROP DEFAULT
`

// Account holds the store account assertion in the local database
type Account struct {
	ID          int
	AuthorityID string
	Assertion   string
	ResellerAPI bool
	APIKey      string
}

// CreateAccountTable creates the database table for an account.
func (db *DB) CreateAccountTable() error {
	_, err := db.Exec(createAccountTableSQL)
	return err
}

// AlterAccountTable modifies the database table for an account.
func (db *DB) AlterAccountTable() error {
	db.Exec(alterAccountResellerAPI)

	err := db.addAccountAPIKeyField()
	return err
}

// addAccountAPIKeyField adds and defaults the API key field to the account table
func (db *DB) addAccountAPIKeyField() error {

	// Add the API key field to the account table
	_, err := db.Exec(alterAccountAPIKey)
	if err != nil {
		// Field already exists so skip
		return nil
	}

	// Default the API key for any records where it is empty
	accounts, err := db.listAllAccounts()
	if err != nil {
		return err
	}
	for _, acc := range accounts {
		if len(acc.APIKey) > 0 {
			continue
		}

		// Generate an random API key and update the record
		apiKey, err := generateAPIKey()
		if err != nil {
			log.Printf("Could not generate random string for the API key")
			return errors.New("Error generating random string for the API key")
		}

		// Update the API key on the model
		acc.APIKey = apiKey
		db.updateAccount(acc)
	}

	// Add the constraints to the API key field
	_, err = db.Exec(alterAccountAPIKeyNotNullable)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) listAllAccounts() ([]Account, error) {
	return db.listAccountsFilteredByUser(anyUserFilter)
}

func (db *DB) listAccountsFilteredByUser(username string) ([]Account, error) {

	var (
		rows *sql.Rows
		err  error
	)

	if len(username) == 0 {
		rows, err = db.Query(listAccountsSQL)
	} else {
		rows, err = db.Query(listUserAccountsSQL, username)
	}
	if err != nil {
		log.Printf("Error retrieving database accounts: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	return rowsToAccounts(rows)
}

// CreateAccount creates an account in the database
func (db *DB) CreateAccount(account Account) error {

	apiKey, err := buildValidOrDefaultAPIKey(account.APIKey)
	if err != nil {
		return errors.New("Error in generating a valid API key")
	}
	account.APIKey = apiKey

	_, err = db.Exec(createAccountSQL, account.AuthorityID, account.Assertion, account.ResellerAPI, account.APIKey)
	if err != nil {
		log.Printf("Error creating the database account: %v\n", err)
		return err
	}
	return nil
}

// GetAccount fetches a single account from the database by the authority ID
func (db *DB) GetAccount(authorityID string) (Account, error) {
	account := Account{}

	err := db.QueryRow(getAccountSQL, authorityID).Scan(&account.ID, &account.AuthorityID, &account.Assertion, &account.ResellerAPI, &account.APIKey)
	if err != nil {
		log.Printf("Error retrieving account: %v\n", err)
		return account, err
	}

	return account, nil
}

// getAccountByID fetches a single account from the database by the ID
func (db *DB) getAccountByID(accountID int) (Account, error) {
	account := Account{}

	err := db.QueryRow(getAccountByIDSQL, accountID).Scan(&account.ID, &account.AuthorityID, &account.Assertion, &account.ResellerAPI, &account.APIKey)
	if err != nil {
		log.Printf("Error retrieving account: %v\n", err)
		return account, err
	}

	return account, nil
}

// getUserAccountByID fetches a single account from the database by the ID
func (db *DB) getUserAccountByID(accountID int, username string) (Account, error) {
	account := Account{}

	err := db.QueryRow(getUserAccountByIDSQL, accountID, username).Scan(&account.ID, &account.AuthorityID, &account.Assertion, &account.ResellerAPI, &account.APIKey)
	if err != nil {
		log.Printf("Error retrieving account: %v\n", err)
		return account, err
	}

	return account, nil
}

// updateAccount updates an account in the database
func (db *DB) updateAccount(account Account) error {
	apiKey, err := buildValidOrDefaultAPIKey(account.APIKey)
	if err != nil {
		return errors.New("Error in generating a valid API key")
	}
	account.APIKey = apiKey

	_, err = db.Exec(updateAccountSQL, account.ID, account.AuthorityID, account.Assertion, account.ResellerAPI, account.APIKey)
	if err != nil {
		log.Printf("Error updating the database account: %v\n", err)
		return err
	}

	return nil
}

// updateUserAccount updates an account in the database
func (db *DB) updateUserAccount(account Account, username string) error {
	apiKey, err := buildValidOrDefaultAPIKey(account.APIKey)
	if err != nil {
		return errors.New("Error in generating a valid API key")
	}
	account.APIKey = apiKey

	_, err = db.Exec(updateUserAccountSQL, account.ID, username, account.AuthorityID, account.Assertion, account.ResellerAPI, account.APIKey)
	if err != nil {
		log.Printf("Error updating the database account: %v\n", err)
		return err
	}

	return nil
}

// putAccount stores an account in the database
func (db *DB) putAccount(account Account) (string, error) {
	_, err := db.Exec(upsertAccountSQL, account.AuthorityID, account.Assertion)
	if err != nil {
		log.Printf("Error updating the database account: %v\n", err)
		return "", err
	}

	return "", nil
}

// ListUserAccounts returns a list of Account objects related with certain user
func (db *DB) ListUserAccounts(username string) ([]Account, error) {
	rows, err := db.Query(listUserAccountsSQL, username)
	if err != nil {
		log.Printf("Error retrieving database accounts of certain user: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	return rowsToAccounts(rows)
}

// ListNotUserAccounts returns a list of Account objects that are not related with certain user
func (db *DB) ListNotUserAccounts(username string) ([]Account, error) {
	rows, err := db.Query(listNotUserAccountsSQL, username)
	if err != nil {
		log.Printf("Error retrieving database accounts not belonging to certain user: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	return rowsToAccounts(rows)
}

func rowsToAccounts(rows *sql.Rows) ([]Account, error) {
	accounts := []Account{}

	for rows.Next() {
		account := Account{}
		err := rows.Scan(&account.ID, &account.AuthorityID, &account.Assertion, &account.ResellerAPI, &account.APIKey)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

// BuildAccountsFromAuthorityIDs from a list of strings representing authority ids, build related datastore.Account objects
func BuildAccountsFromAuthorityIDs(authorityIDs []string) []Account {
	var accounts []Account
	for _, authorityID := range authorityIDs {
		accounts = append(accounts, BuildAccountFromAuthorityID(authorityID))
	}
	return accounts
}

// BuildAccountFromAuthorityID from a string representing authority id, build related datastore.Account object
func BuildAccountFromAuthorityID(authorityID string) Account {
	return Account{AuthorityID: authorityID}
}
