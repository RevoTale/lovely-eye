# End-to-end testing
This module handle only the end-to-end testing. 

No source code should be included as a dependency inside tests! Exceptions are:
- root production server
- root config
- some basic types like specific constants to keep the data integrity.
- generated operations should be imported as `operations "github.com/lovely-eye/server/e2e/generated" `