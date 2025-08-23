package ofd

// AdditionalOfdParams
type AdditionalOfdParams struct {
	AgentInfoType                   int      `json:"agent_info.type"`
	AgentInfoPayingOperation        string   `json:"agent_info.paying.operation,omitempty"`
	AgentInfoPayingPhones           []string `json:"agent_info.paying.phones,omitempty"`
	AgentInfoPaymentsOperatorPhones []string `json:"agent_info.paymentsOperator.phones,omitempty"`
	AgentInfoMTOperatorAddress      string   `json:"agent_info.MTOperator.address,omitempty"`
	AgentInfoMTOperatorInn          string   `json:"agent_info.MTOperator.inn,omitempty"`
	AgentInfoMTOperatorName         string   `json:"agent_info.MTOperator.name,omitempty"`
	AgentInfoMTOperatorPhones       []string `json:"agent_info.MTOperator.phones,omitempty"`
	SupplierInfoPhones              []string `json:"supplier_info.phones,omitempty"`
	Cashier                         string   `json:"cashier,omitempty"`
	AdditionalCheckProps            string   `json:"additional_check_props,omitempty"`
	AdditionalUserPropsName         string   `json:"additional_user_props.name,omitempty"`
	AdditionalUserPropsValue        string   `json:"additional_user_props.value,omitempty"`
	CashierInn                      string   `json:"cashier_inn,omitempty"`
	ClientAddress                   string   `json:"client.address,omitempty"`
	ClientBirthDate                 string   `json:"client.birth_date,omitempty"`
	ClientCitizenship               string   `json:"client.citizenship,omitempty"`
	ClientDocumentCode              string   `json:"client.document_code,omitempty"`
	ClientPassportNumber            string   `json:"client.passport_number,omitempty"`
	ClientMail                      string   `json:"client.email,omitempty"`
	ClientPhone                     string   `json:"client.phone,omitempty"`
	ClientInn                       string   `json:"client.inn,omitempty"`
	ClientName                      string   `json:"client.name,omitempty"`
	OperatingCheckPropsName         string   `json:"operatingCheckProps.name,omitempty"`
	OperatingCheckPropsTimestamp    string   `json:"operatingCheckProps.timestamp,omitempty"`
	OperatingCheckPropsValue        string   `json:"operatingCheckProps.value,omitempty"`
	SectoralCheckPropsDate          string   `json:"sectoralCheckProps.date,omitempty"`
	SectoralCheckPropsFederalId     string   `json:"sectoralCheckProps.federalId,omitempty"`
	SectoralCheckPropsNumber        string   `json:"sectoralCheckProps.number,omitempty"`
	SectoralCheckPropsValue         string   `json:"sectoralCheckProps.value,omitempty"`
}
