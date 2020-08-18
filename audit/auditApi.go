package audit

import (
	"encoding/json"
	"fmt"
	"github.com/tarent/gomulocity/generic"
	"log"
	"net/http"
	"net/url"
)

type AuditApi interface {
	GetAuditRecord(auditID string) (*AuditRecord, *generic.Error)
	GetAuditRecords(auditQuery *AuditQuery, pageSize int) (*AuditRecordCollection, *generic.Error)
	CreateAuditRecord(record *AuditRecord) (*AuditRecord, *generic.Error)
	NextPage(c *AuditRecordCollection) (*AuditRecordCollection, *generic.Error)
	PreviousPage(c *AuditRecordCollection) (*AuditRecordCollection, *generic.Error)
}

type auditApi struct {
	client   *generic.Client
	basePath string
}

func NewAuditApi(client *generic.Client) AuditApi {
	return &auditApi{
		client:   client,
		basePath: "/audit/auditRecords",
	}
}

func (a *auditApi) GetAuditRecord(auditID string) (*AuditRecord, *generic.Error) {
	if len(auditID) == 0 {
		return nil, generic.ClientError("Getting an audit record without recordID is not allowed", "GetAuditRecord")
	}

	body, status, err := a.client.Get(fmt.Sprintf("%v/%v", a.basePath, auditID), generic.AcceptHeader(AUDIT_RECORD_ACCEPT))
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while getting audit record: %s", err), "GetAuditRecord")
	}

	if status != http.StatusOK {
		return nil, generic.CreateErrorFromResponse(body, status)
	}

	record := &AuditRecord{}
	if err := json.Unmarshal(body, record); err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while unmarshalling response: %s", err), "GetAuditRecord")
	}
	return record, nil
}

func (a *auditApi) GetAuditRecords(auditQuery *AuditQuery, pageSize int) (*AuditRecordCollection, *generic.Error) {
	return a.find(auditQuery, pageSize)
}

func (a *auditApi) CreateAuditRecord(record *AuditRecord) (*AuditRecord, *generic.Error) {
	bytes, err := json.Marshal(record)
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while marshalling audit record: %s", err), "CreateAuditRecord")
	}

	body, status, err := a.client.Post(a.basePath, bytes, generic.ContentTypeHeader(AUDIT_CONTENT_TYPE))
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while creating a new audit record: %s", err), "CreateAuditRecord")
	}

	if status != http.StatusCreated {
		return nil, generic.CreateErrorFromResponse(body, status)
	}

	auditRecord := &AuditRecord{}
	if err := json.Unmarshal(body, auditRecord); err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while unmarshalling response body: %s", err), "CreateAuditRecord")
	}
	return auditRecord, nil
}

func (a *auditApi) find(auditQuery *AuditQuery, pageSize int) (*AuditRecordCollection, *generic.Error) {
	queryParamsValues := &url.Values{}
	err := auditQuery.QueryParams(queryParamsValues)
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while building query parameters to search for audit records: %s", err.Error()), "FindAuditRecords")
	}

	err = generic.PageSizeParameter(pageSize, queryParamsValues)
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while building pageSize parameter to fetch audit records: %s", err.Error()), "FindAuditRecords")
	}

	return a.getCommon(fmt.Sprintf("%s?%s", a.basePath, queryParamsValues.Encode()))
}

func (a *auditApi) NextPage(c *AuditRecordCollection) (*AuditRecordCollection, *generic.Error) {
	return a.getPage(c.Next)
}

func (a *auditApi) PreviousPage(c *AuditRecordCollection) (*AuditRecordCollection, *generic.Error) {
	return a.getPage(c.Prev)
}

func (a *auditApi) getPage(reference string) (*AuditRecordCollection, *generic.Error) {
	if reference == "" {
		log.Print("No page reference given. Returning nil.")
		return nil, nil
	}

	nextUrl, err := url.Parse(reference)
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Unparsable URL given for page reference: '%s'", reference), "GetPage")
	}

	collection, genErr := a.getCommon(fmt.Sprintf("%s?%s", nextUrl.Path, nextUrl.RawQuery))
	if genErr != nil {
		return nil, genErr
	}

	if len(collection.AuditRecords) == 0 {
		log.Print("Returned collection is empty. Returning nil.")
		return nil, nil
	}

	return collection, nil
}

func (a *auditApi) getCommon(path string) (*AuditRecordCollection, *generic.Error) {
	body, status, err := a.client.Get(path, generic.AcceptHeader(AUDIT_RECORD_ACCEPT))
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while getting audit records: %s", err.Error()), "GetAuditRecords")
	}

	if status != http.StatusOK {
		return nil, generic.CreateErrorFromResponse(body, status)
	}

	return parseAuditRecordCollectionResponse(body)
}

func parseAuditRecordCollectionResponse(body []byte) (*AuditRecordCollection, *generic.Error) {
	var result AuditRecordCollection
	if len(body) > 0 {
		err := generic.ObjectFromJson(body, &result)
		if err != nil {
			return nil, generic.ClientError(fmt.Sprintf("Error while parsing response JSON: %s", err.Error()), "CollectionResponseParser")
		}
	} else {
		return nil, generic.ClientError("Response body was empty", "CollectionResponseParser")
	}

	return &result, nil
}
