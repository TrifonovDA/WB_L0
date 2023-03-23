package config

type BdCredentials struct {
	Host     string
	Port     string
	Database string
	Username string
	Password string
}

var BdCred = BdCredentials{
	Host:     "localhost",
	Port:     "5432",
	Database: "L0_task",
	Username: "d.triphonov",
	Password: "C4kABAtg",
}
var Http_creds = "0.0.0.0:60606"
