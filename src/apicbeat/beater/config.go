package beater

type ApicConfig struct {
	Period *int64
	Addr *string
	User *string
	Passwd *string
	ForwardSet struct {
		Tennant_Health *bool
		Tennant_Health_Cur *bool
		Tennant_Endpoint *bool
		Tennant_Endpoint_DN *bool
		Fault_Info *bool
	}
}

type ConfigSettings struct {
	Input    *ApicConfig `config:"input"`
	Apicbeat *ApicConfig `config:"apicbeat"`
}
