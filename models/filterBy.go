package models

const (
	FILTER_BYSOURCE       = "source"
	FILTER_BYSTATUS       = "status"
	FILTER_BYAGENTID      = "agentId"
	FILTER_BYTYPE         = "type"
	FILTER_BYDATEFROM     = "dateFrom"
	FILTER_BYDATETO       = "dateTo"
	FILTER_BYFRAGMENTTYPE = "fragmentType"
	FILTER_BYDEVICEID     = "deviceId"
	FILTER_BYTEXT         = "text"
	FILTER_BYLISTOFIDs    = "ids"
	FILTER_BYUSER         = "user"
	FILTER_BYAPPLICATION  = "application"
)

func GetFilters() []string {
	return []string{FILTER_BYSOURCE, FILTER_BYDATEFROM, FILTER_BYDATETO, FILTER_BYFRAGMENTTYPE, FILTER_BYTYPE}
}
