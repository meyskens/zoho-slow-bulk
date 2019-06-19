Zoho slow bulk update
=====================

When needing to bulk update a field in Zoho I found the workflow rules not to work correct. This is a small Go script that updates individual reccords with a 3 second interval so the workflow rules work. I have used this script only to trigger them for all records.

## How to use
You need to edit `main.go` and edit the data structure in `newData`, toghether with the change in `newEntry`, when being ran it will update an Account every 3 seconds and executes the edit workflow(s)

The API OAuth client ID and Secret need to be set as environment variables as `ZOHO_CLIENT_ID` and `ZOHO_CLIENT_SECRET`
