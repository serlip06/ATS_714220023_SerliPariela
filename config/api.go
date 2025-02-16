package config

import "os"

var ApiWaButton string = os.Getenv("URLAPIWABUTTON")
var GHAccessToken string = os.Getenv("TOKEN_GITHUB")
