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


var intlData = {
    zh: {},

    en: {
      "account": "Account",
      "account-api-key": "Management API key for the account",
      "account-description": "The authority-id for the account",
      "account-keys": "Account Key Assertions",
      "accounts": "Accounts",
      "activate": "Activate",
      "active": "Active",
      "add": "Add",
      "add-new-model": "Add a new model",
      "add-new-signing-key": "Import a signing key",
      "add-new-user": "Add a new user",
      "api-key": "API Key",
      "api-key-description": "API Key to sign a serial assertion request (min. 10 characters). Will be generated if blank or invalid",
      "assertion": "Assertion",
      "assertion-settings": "Configure the model assertion headers",
      "assertion-status": "Assertion Status",
      "authority-id-description": "The authority for signing models",
      "authority-id": "Signing Authority",
      "brand": "Brand",
      "brand-description": "The name of the device brand",
      "cancel": "Cancel",
      "close": "Close",
      "complete": "Complete",
      "confirm-log-delete": "Remove this log?",
      "confirm-model-delete": "Remove this model?",
      "confirm-store-delete": "Remove this sub-store model?",
      "confirm-user-delete": "Remove this user?",
      "create-assertion": "Error creating the assertion",
      "create-system-user": "Create System-User",
      "date": "Date",
      "deactivate": "Deactivate",
      "delete-log": "Delete log",
      "delete-model": "Delete model",
      "delete-user": "Delete user",
      "description": "The Serial Vault is a web service that generates cryptographically-signed serial assertions.",
      "download": "Download",
      "edit-model": "Edit Model",
      "edit-user": "Edit User",
      "email": "Email",
      "email-description": "Email for the Store",
      "error-adding-key": "Error adding a public key",
      "error-auth": "Unauthorized action",
      "error-created-model": "Cannot find the created model",
      "error-creating-model": "Error creating the model",
      "error-creating-store": "Error creating the sub-store model",
      "error-creating-user": "Error creating the user",
      "error-decode-json": "Error decoding JSON",
      "error-decode-key": "Error decoding the base64 Signing Key",
      "error-deleting-key": "Error deleting a public key",
      "error-fetch-models": "Error fetching the models",
      "error-fetch-users": "Error fetching the users",
      "error-format-assertions": "Error formatting the assertions",
      "error-get-model": "Cannot find the model",
      "error-get-non-user-accounts": "Cannot get user not related accounts",
      "error-get-user": "Cannot find the user",
      "error-invalid-key": "The Signing-key is invalid",
      "error-invalid-model": "Invalid model ID",
      "error-invalid-user": "Invalid user ID",
      "error-key-data": "No data supplied for the public key",
      "error-key-exists": 'The ssh public key already exists',
      "error-key-not-found": "The ssh public key cannot be found",
      "error-keypair-store": "Error storing the signing key",
      "error-model-data": "No model data supplied",
      "error-model-exists": "A device with the same Brand and Model already exists",
      "error-model-not-active": "The model is linked with an inactive signing-key",
      "error-model-not-found": "Cannot find model with the matching brand and model",
      "error-nil-data": "Uninitialized POST data",
      "error-no-permissions": "You do not have permissions to access this page",
      "error-read-private-key": "Error reading the private key",
      "error-sign-empty": "No data supplied for signing",
      "error-signing-assertions": "Error signing the assertions",
      "error-updating-model": "Error updating the model",
      "error-updating-store": "Error updating sub-store model",
      "error-updating-user": "Error updating the user",
      "error-user-data": "No user data supplied",
      "error-validate-key": "The public key must be entered",
      "error-validate-model": "The Brand and Model must be supplied",
      "error-validate-new-model": "The Brand, Model and Signing-Keys must be supplied",
      "error-validate-signingkey": "The Serial Assertion Key must be selected",
      "error-validate-userkey": "The System-User Assertion Key must be selected",
      "find-serialnumber": "find serial number",
      "fingerprint": "Fingerprint",
      "generate": "Generate",
      "generate-signing-key": "Generate Signing Key",
      "home": "Home",
      "inactive": "Inactive",
      "key-id": "Key ID",
      "key-name-description": "Unique name for the key in the store",
      "key-name": "Key Name",
      "login": "Login",
      "logout": "Logout",
      "makes": "Makes",
      "model-description": "The name of the device model",
      "model": "Model",
      "modelname": "Model Name",
      "modelname-description": "The name of the pivoted model",
      "models_available": "The following models are available",
      "models": "Models",
      "more": "More",
      "name": "Full name",
      "new-account-assertion": "Add or Replace Account Assertion",
      "new-account-key-assertion": "Add or Replace Account Key Assertion",
      "new-model": "New Model",
      "new-public-key-description": "Paste the public key of the machine that needs access to the Serial Vault",
      "new-public-key": "New Public Key",
      "new-signing-key-description": "Paste the signing-key or upload the file",
      "new-signing-key": "Import Signing Key",
      "new-substore-device": "Create a new sub-store mapping for a device",
      "new-user": "New User",
      "no": "No",
      "no-assertion-key": "No account key assertion found",
      "no-assertion": "No account assertion found",
      "no-assertions": "No assertions found",
      "no-signing-keys-found": "No signing keys found",
      "not-used-signing": "Not used for signing system-user assertions",
      "otp": "OTP",
      "otp-description": "One-time password for SSO",
      "password": "Password",
      "password-description": "Password for the Store",
      "private-key-description": "The signing-key that will be used to sign the device identity",
      "private-key-model": "Model Assertion Key",
      "private-key-model-short": "Assertion",
      "private-key": "Serial Assertion Key",
      "private-key-short": "Serial Key",
      "private-key-user-description": "The signing-key that will be used to sign the system-user assertion",
      "private-key-user-short": "System-User Key",
      "private-key-user": "System-User Assertion Key",
      "public-key-confirm": "Confirm deletion of the public key",
      "public-key-description": "Add a new public key",
      "public-key": "Public Key",
      "public-keys-authorized": "The following keys are authorized",
      "public-keys": "Public Keys",
      "register-signing-key": "Register Signing Key with the Store",
      "remove": "Remove",
      "reseller": "Reseller",
      "reseller-features": "Enable Reseller Features",
      "revision": "Revision",
      "role": "Role",
      "save": "Save",
      "select-accounts": "Select below the accounts this user belongs to:",
      "serial-number-description": "Serial Number of the device",
      "serial-number": "Serial Number",
      "signing-key": "Signing Key",
      "signing-keys": "Signing Keys",
      "signinglog-description": "Log of the serial numbers and device-key fingerprints that have been used",
      "signinglog": "Signing Log",
      "store-account-assertion": "Account Assertion from Store",
      "store-credentials": "Store Credentials",
      "substore": "Sub-Store",
      "substore-description": "Store ID of the Sub-Store",
      "substore-model": "Sub-Store Model",
      "substores": "Sub-Stores",
      "systemuser": "System-User",
      "title": "Serial Vault",
      "upload-account-assertion": "Upload Account Assertion",
      "user-accounts": "User Accounts",
      "user-email": "The email address of the user",
      "user-key": "User Key",
      "user-name": "The full name of the user",
      "username": "Username",
      "user-not-found": "You do not have permissions to access this site",
      "users_available": "The following users are available",
      "users": "Users",
      "user": "User",
      "user-username": "The nickname of the user",
      "version": "Version",
      "yes": "Yes",
    }
};

export default intlData;
