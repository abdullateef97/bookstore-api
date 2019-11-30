package helpers

import "log"

// LogFatal logs errors and stops code execution
func LogFatal(err error) {
	if(err != nil){
		log.Fatalln(err)
	}
}