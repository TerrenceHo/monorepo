/*
Package stackerrors provides simple error handling, on top of the stdlib's
implementation of error handling.

Context can be added to an error: instead of just

 	if err != nil {
 			return nil
 	}

We can wrap the original error with context of when this error was received
before passing the error up the chain.
*/

package stackerrors
