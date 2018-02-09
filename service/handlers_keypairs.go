// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2016-2017 Canonical Ltd
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

package service

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/CanonicalLtd/serial-vault/datastore"
	"github.com/gorilla/mux"
	"github.com/snapcore/snapd/asserts"
)

// KeypairWithPrivateKey is the JSON version of a keypair, including the base64 armored, signing-key
type KeypairWithPrivateKey struct {
	ID          int    `json:"id"`
	AuthorityID string `json:"authority-id"`
	PrivateKey  string `json:"private-key"`
	KeyName     string `json:"key-name"`
}

// KeypairStatusResponse is the JSON response from the API status of keypair generation
type KeypairStatusResponse struct {
	Success      bool                      `json:"success"`
	ErrorCode    string                    `json:"error_code"`
	ErrorSubcode string                    `json:"error_subcode"`
	ErrorMessage string                    `json:"message"`
	Status       []datastore.KeypairStatus `json:"status"`
}

// KeypairListHandler fetches the available keypairs for display from the database.
// Only viewable reference data is stored in the database, not the restricted private key.
func KeypairListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	authUser, err := checkIsAdminAndGetUserFromJWT(w, r)
	if err != nil {
		formatKeypairsResponse(false, "error-auth", "", "", nil, w)
		return
	}

	vars := mux.Vars(r)
	accountCode := vars["account_code"]

	userKeypairs, err := datastore.Environ.DB.ListAllowedKeypairs(authUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		formatKeypairsResponse(false, "error-fetch-keypairs", "", err.Error(), nil, w)
		return
	}

	// Filter the list down to the keypairs for the selected account
	keypairs := []datastore.Keypair{}
	for _, k := range userKeypairs {
		if k.AuthorityID == accountCode {
			keypairs = append(keypairs, k)
		}
	}

	// Return successful JSON response with the list of models
	w.WriteHeader(http.StatusOK)
	formatKeypairsResponse(true, "", "", "", keypairs, w)
}

// KeypairCreateHandler is the API method to create a new keypair that can be used
// for signing serial assertions. The keypairs are stored in the signing database
// and the authority-id/key-id is stored in the models database. Models can then be
// linked to one of the existing signing-keys.
func KeypairCreateHandler(w http.ResponseWriter, r *http.Request) {

	keypairWithKey, ok := verifyKeypair(w, r)
	if !ok {
		return
	}

	// Store the signing-key in the keypair store using the asserts module
	privateKey, sealedPrivateKey, err := datastore.Environ.KeypairDB.ImportSigningKey(keypairWithKey.AuthorityID, keypairWithKey.PrivateKey)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		formatBooleanResponse(false, "error-keypair-store", "", err.Error(), w)
		return
	}

	// Store the signing-key in the database
	keypair := datastore.Keypair{
		AuthorityID: keypairWithKey.AuthorityID,
		KeyID:       privateKey.PublicKey().ID(),
		SealedKey:   sealedPrivateKey,
	}
	errorCode, err := datastore.Environ.DB.PutKeypair(keypair)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		formatBooleanResponse(false, errorCode, "", err.Error(), w)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	formatBooleanResponse(true, "", "", "", w)
}

// KeypairDisableHandler disables an existing keypair, which will mean that any
// linked Models will not be able to be signed. The asserts module does not allow
// a keypair to be deleted, so the keypair will just be disabled in the local database.
func KeypairDisableHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	authUser, err := checkIsAdminAndGetUserFromJWT(w, r)
	if err != nil {
		formatBooleanResponse(false, "error-auth", "", "", w)
		return
	}

	// Get the keypair primary key
	vars := mux.Vars(r)
	keypairID, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errorMessage := fmt.Sprintf("%v", vars["id"])
		formatBooleanResponse(false, "error-invalid-key", "", errorMessage, w)
		return
	}

	// Update the keypair in the local database
	err = datastore.Environ.DB.UpdateAllowedKeypairActive(keypairID, false, authUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		formatBooleanResponse(false, "error-keypair-update", "", err.Error(), w)
		return
	}

	formatBooleanResponse(true, "", "", "", w)
}

// KeypairEnableHandler enables an existing keypair, which will mean that any
// linked Models will be able to be signed. The asserts module does not allow
// a keypair to be deleted, so the keypair will just be enabled in the local database.
func KeypairEnableHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	authUser, err := checkIsAdminAndGetUserFromJWT(w, r)
	if err != nil {
		formatBooleanResponse(false, "error-auth", "", "", w)
		return
	}

	// Get the keypair primary key
	vars := mux.Vars(r)
	keypairID, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		errorMessage := fmt.Sprintf("%v", vars["id"])
		formatBooleanResponse(false, "error-invalid-key", "", errorMessage, w)
		return
	}

	// Update the keypair in the local database
	err = datastore.Environ.DB.UpdateAllowedKeypairActive(keypairID, true, authUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		formatBooleanResponse(false, "error-keypair-update", "", err.Error(), w)
		return
	}
	formatBooleanResponse(true, "", "", "", w)
}

// KeypairAssertionHandler updates the account key assertion on a keypair
func KeypairAssertionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	authUser, err := checkIsAdminAndGetUserFromJWT(w, r)
	if err != nil {
		formatBooleanResponse(false, "error-auth", "", "", w)
		return
	}

	// Check that we have a message body
	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		formatBooleanResponse(false, "error-nil-data", "", "Uninitialized POST data", w)
		return
	}
	defer r.Body.Close()

	assertionRequest := AssertionRequest{}
	err = json.NewDecoder(r.Body).Decode(&assertionRequest)
	switch {
	// Check we have some data
	case err == io.EOF:
		w.WriteHeader(http.StatusBadRequest)
		formatBooleanResponse(false, "error-assertion-data", "", "No assertion data supplied", w)
		return
		// Check for parsing errors
	case err != nil:
		w.WriteHeader(http.StatusBadRequest)
		formatBooleanResponse(false, "error-assertion-json", "", err.Error(), w)
		return
	}

	// Check that a keypair ID has been provided
	if assertionRequest.ID == 0 {
		logMessage("KEYPAIR", "invalid-keypair", "ID of keypair not provided")
		w.WriteHeader(http.StatusBadRequest)
		formatBooleanResponse(false, "invalid-keypair", "", "ID of keypair not provided", w)
		return
	}

	// Decode the file
	decodedAssertion, err := base64.StdEncoding.DecodeString(assertionRequest.Assertion)
	if err != nil {
		logMessage("KEYPAIR", "invalid-assertion", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		formatBooleanResponse(false, "decode-assertion", "", err.Error(), w)
		return
	}

	// Validate the assertion in the request
	assertion, err := asserts.Decode(decodedAssertion)
	if err != nil {
		logMessage("KEYPAIR", "invalid-assertion", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		formatBooleanResponse(false, "decode-assertion", "", err.Error(), w)
		return
	}

	// Check that we have an account key assertion
	if assertion.Type().Name != asserts.AccountKeyType.Name {
		w.WriteHeader(http.StatusBadRequest)
		formatBooleanResponse(false, "invalid-assertion", "", fmt.Sprintf("An assertion of type '%s' is required", asserts.AccountKeyType.Name), w)
		return
	}

	keypair := datastore.Keypair{
		ID:          assertionRequest.ID,
		AuthorityID: assertion.HeaderString("account-id"),
		KeyID:       assertion.HeaderString("public-key-sha3-384"),
		Assertion:   string(decodedAssertion),
	}

	errorCode, err := datastore.Environ.DB.UpdateKeypairAssertion(keypair, authUser)
	if err != nil {
		logMessage("KEYPAIR", "invalid-assertion", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		formatBooleanResponse(false, errorCode, "", err.Error(), w)
		return
	}

	// Return the success response
	formatBooleanResponse(true, "", "", "", w)
}

// KeypairGenerateHandler is the API method to generate a new keypair that can be used
// for signing serial (or model) assertions. The keypairs are stored in the signing database
// and the authority-id/key-id is stored in the models database. Models can then be
// linked to one of the existing signing-keys.
func KeypairGenerateHandler(w http.ResponseWriter, r *http.Request) {

	keypair, ok := verifyKeypair(w, r)
	if !ok {
		return
	}

	go datastore.GenerateKeypair(keypair.AuthorityID, "", keypair.KeyName)

	// Return the URL to watch for the response
	statusURL := fmt.Sprintf("/v1/keypairs/status/%s/%s", keypair.AuthorityID, keypair.KeyName)
	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Location", statusURL)
	formatBooleanResponse(true, "", "", statusURL, w)
}

// KeypairStatusHandler returns the creation status of a keypair
func KeypairStatusHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	authUser, err := checkIsAdminAndGetUserFromJWT(w, r)
	if err != nil {
		formatBooleanResponse(false, "error-auth", "", "", w)
		return
	}

	// Check that the user has permissions to this authority-id
	if !datastore.Environ.DB.CheckUserInAccount(authUser.Username, vars["authorityID"]) {
		w.WriteHeader(http.StatusBadRequest)
		formatBooleanResponse(false, "error-auth", "", "Your user does not have permissions for the Signing Authority", w)
		return
	}

	ks, err := datastore.Environ.DB.GetKeypairStatus(vars["authorityID"], vars["keyName"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		formatBooleanResponse(false, "error-keypair-json", "", "Cannot find the status of the keypair", w)
		return
	}

	formatBooleanResponse(true, "", "", ks.Status, w)
}

// KeypairStatusProgressHandler returns the status of keypairs that are being generated
func KeypairStatusProgressHandler(w http.ResponseWriter, r *http.Request) {

	authUser, err := checkIsAdminAndGetUserFromJWT(w, r)
	if err != nil {
		formatBooleanResponse(false, "error-auth", "", "", w)
		return
	}

	ks, err := datastore.Environ.DB.ListAllowedKeypairStatus(authUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		formatBooleanResponse(false, "error-keypair-json", "", "Cannot find the status of the keypairs", w)
		return
	}

	formatKeypairStatusResponse(true, "", "", "", ks, w)
}

func verifyKeypair(w http.ResponseWriter, r *http.Request) (KeypairWithPrivateKey, bool) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	keypairWithKey := KeypairWithPrivateKey{}

	authUser, err := checkIsAdminAndGetUserFromJWT(w, r)
	if err != nil {
		formatBooleanResponse(false, "error-auth", "", "", w)
		return keypairWithKey, false
	}

	// Check that we have a message body
	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		formatBooleanResponse(false, "error-nil-data", "", "Uninitialized POST data", w)
		return keypairWithKey, false
	}
	defer r.Body.Close()

	// Decode the JSON body
	err = json.NewDecoder(r.Body).Decode(&keypairWithKey)
	switch {
	// Check we have some data
	case err == io.EOF:
		w.WriteHeader(http.StatusBadRequest)
		formatBooleanResponse(false, "error-keypair-data", "", "No keypair data supplied", w)
		return keypairWithKey, false
		// Check for parsing errors
	case err != nil:
		w.WriteHeader(http.StatusBadRequest)
		formatBooleanResponse(false, "error-keypair-json", "", err.Error(), w)
		return keypairWithKey, false
	}

	// Validate the authority-id
	keypairWithKey.AuthorityID = strings.TrimSpace(keypairWithKey.AuthorityID)
	if len(keypairWithKey.AuthorityID) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		formatBooleanResponse(false, "error-keypair-json", "", "The authority-id is mandatory", w)
		return keypairWithKey, false
	}

	// Check that the user has permissions to this authority-id
	if !datastore.Environ.DB.CheckUserInAccount(authUser.Username, keypairWithKey.AuthorityID) {
		w.WriteHeader(http.StatusBadRequest)
		formatBooleanResponse(false, "error-auth", "", "Your user does not have permissions for the Signing Authority", w)
		return keypairWithKey, false
	}

	return keypairWithKey, true

}
