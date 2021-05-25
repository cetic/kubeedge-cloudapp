package api

type Response struct {
	Arg 	string  `json:"arg"`
	Job 	string	`json:"job"`
	Status 	string	`json:"status"`
	Trigger string	`json:"trigger"`
}
