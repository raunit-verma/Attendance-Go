package util

// const One = "You are not authorized to do this operation."
// const Two = "Cannot decode payload."
// const Three = "New User data is missing."
// const Four = "Username taken or user already exist."

const (
	NotAuthorized_One              = "You are not authorized to do this operation." // One
	CannotDecodePayload_Two        = "Cannot decode payload."                       // Two
	UserDataMissing_Three          = "New User data is missing."                    // Three
	UsernameTaken_Four             = "Username taken or user already exist."        // Four
	RequestDataValidation_Five     = "Request data validation failed."              // Five
	UserNotFound_Six               = "User not found."                              // Six
	DBError_Seven                  = "Error in doing DB operation."                 // Seven
	Success_Eight                  = "Operation completed."                         // Eight
	OperationNotAllowed_Nine       = "Operation not allowed."                       // Nine
	InternalServererror        int = 500
)
