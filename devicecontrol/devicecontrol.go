package devicecontrol

import (
	"encoding/json"
	"fmt"
	"github.com/tarent/gomulocity/generic"
	"log"
	"net/http"
	"net/url"
)

func NewDeviceControlApi(client *generic.Client) DeviceControl {
	return &deviceControl{client, "/devicecontrol/operations", "/devicecontrol/bulkoperations"}
}

type DeviceControl interface {
	GetOperation(operationID string) (*Operation, *generic.Error)
	CreateOperation(operation *NewOperation) (*Operation, *generic.Error)
	UpdateOperation(operationID string, operation *UpdateOperation) (string, *generic.Error)
	GetOperationCollection(query OperationQuery, pageSize int) (*OperationCollection, *generic.Error)
	DeleteOperationCollection(query OperationQuery) *generic.Error
	CreateBulkOperation(bulkOperation *NewBulkOperation) (*BulkOperation, *generic.Error)
	GetCollectionOfBulkOperation(query OperationQuery, pageSize int) (*BulkOperationCollection, *generic.Error)
	UpdateBulkOperation(bulkOperationID string, operation *UpdateBulkOperation) (*BulkOperation, *generic.Error)
	GetBulkOperation(bulkOperationID string) (*BulkOperation, *generic.Error)
	DeleteBulkOperation(bulkOperationID string) *generic.Error
	NextPage(c *OperationCollection) (*OperationCollection, *generic.Error)
	PreviousPage(c *OperationCollection) (*OperationCollection, *generic.Error)
	FindOperationCollection(query OperationQuery, pageSize int) (*OperationCollection, *generic.Error)
	FindBulkOperationCollection(query OperationQuery, pageSize int) (*BulkOperationCollection, *generic.Error)
}

type deviceControl struct {
	client                 *generic.Client
	basePathOperations     string
	basePathBulkOperations string
}

func (d *deviceControl) GetOperation(operationID string) (*Operation, *generic.Error) {
	if len(operationID) == 0 {
		return nil, generic.ClientError("Getting operation without an id is not allowed", "GetOperation")
	}

	body, status, err := d.client.Get(fmt.Sprintf("%v/%v", d.basePathOperations, url.QueryEscape(operationID)), generic.AcceptHeader(OPERATION_ACCEPT_HEADER))
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while getting an operation: %s", err.Error()), "GetOperation")
	}

	if status != http.StatusOK {
		return nil, generic.CreateErrorFromResponse(body, status)
	}
	return parseOperationResponse(body)
}

func (d *deviceControl) CreateOperation(operation *NewOperation) (*Operation, *generic.Error) {
	bytes, err := generic.JsonFromObject(operation)
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while marhalling the operation: %s", err.Error()), "CreateOperation")
	}
	body, status, err := d.client.Post(d.basePathOperations, []byte(bytes), generic.AcceptAndContentTypeHeader(OPERATION_ACCEPT_HEADER, OPERATION_CONTENT_TYPE))
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while posting a new operation: %s", err), "CreateOperation")
	}

	if status != http.StatusCreated {
		return nil, generic.CreateErrorFromResponse(body, status)
	}
	return parseOperationResponse(body)
}

func (d *deviceControl) UpdateOperation(operationID string, operation *UpdateOperation) (string, *generic.Error) {
	if len(operationID) == 0 {
		return "", generic.ClientError("Updating operation without an id is not allowed", "UpdateOperation")
	}

	bytes, err := json.Marshal(operation)
	if err != nil {
		return "", generic.ClientError(fmt.Sprintf("Error while marshalling the operation: %s", err.Error()), "UpdateOperation")
	}

	body, status, err := d.client.Put(fmt.Sprintf("%v/%v", d.basePathOperations, url.QueryEscape(operationID)), bytes, generic.EmptyHeader())
	if err != nil {
		return "", generic.ClientError(fmt.Sprintf("Error while updating operation. Given operationID %v, %s", operationID, err), "UpdateOperation")
	}

	if status != http.StatusOK {
		return "", generic.CreateErrorFromResponse(body, status)
	}

	responseStatus := struct {
		Status string `json:"status"`
	}{}

	if err := json.Unmarshal(body, &responseStatus); err != nil {
		return "", generic.ClientError(fmt.Sprintf("Error while unmarshalling update status: %s", err), "UpdateOperation")
	}
	return responseStatus.Status, nil
}

func (d *deviceControl) GetOperationCollection(query OperationQuery, pageSize int) (*OperationCollection, *generic.Error) {
	return d.FindOperationCollection(query, pageSize)
}

func (d *deviceControl) DeleteOperationCollection(query OperationQuery) *generic.Error {
	operationQuery := &url.Values{}
	query.QueryParams(operationQuery)

	if len(*operationQuery) == 0 {
		return generic.ClientError("No filter set", "DeleteOperationCollection")
	}

	body, status, err := d.client.Delete(fmt.Sprintf("%v?%v", d.basePathOperations, operationQuery.Encode()), generic.EmptyHeader())
	if err != nil {
		return generic.ClientError(fmt.Sprintf("Error while deleting operation collection"), "DeleteOperationCollection")
	}

	if status != http.StatusNoContent {
		return generic.CreateErrorFromResponse(body, status)
	}
	return nil
}

func (d *deviceControl) CreateBulkOperation(bulkOperation *NewBulkOperation) (*BulkOperation, *generic.Error) {
	json, err := generic.JsonFromObject(bulkOperation)
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("failed to marshal bulkOperation: %s", err), "CreateBulkOperation")
	}

	body, status, err := d.client.Post(d.basePathBulkOperations, []byte(json), generic.AcceptAndContentTypeHeader(BULK_OPERATION_ACCEPT_HEADER, BULK_OPERATION_CONTENT_TYPE))
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while posting new bulk operation: %s", err), "CreateBulkOperation")
	}

	if status != http.StatusCreated {
		return nil, generic.CreateErrorFromResponse(body, status)
	}
	return parseBulkOperationResponse(body)
}

func (d *deviceControl) GetBulkOperation(bulkOperationID string) (*BulkOperation, *generic.Error) {
	if len(bulkOperationID) == 0 {
		return nil, generic.ClientError("Getting bulk operation without a bulkOperationID is not allowed", "GetBulkOperation")
	}

	body, status, err := d.client.Get(fmt.Sprintf("%v/%v", d.basePathBulkOperations, bulkOperationID), generic.AcceptHeader(BULK_OPERATION_ACCEPT_HEADER))
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while getting bulkOperation by ID: %v, %s", bulkOperationID, err), "GetBulkOperation")
	}

	if status != http.StatusOK {
		return nil, generic.CreateErrorFromResponse(body, status)
	}
	return parseBulkOperationResponse(body)
}

func (d *deviceControl) DeleteBulkOperation(bulkOperationID string) *generic.Error {
	if len(bulkOperationID) == 0 {
		return generic.ClientError("Deleting bulk operation without a bulkOperationID is not allowed", "DeleteBulkOperation")
	}

	body, status, err := d.client.Delete(fmt.Sprintf("%v/%v", d.basePathBulkOperations, bulkOperationID), generic.EmptyHeader())
	if err != nil {
		return generic.ClientError(fmt.Sprintf("Error while deleting bulkOperation by ID: %v, %s", bulkOperationID, err), "DeleteBulkOperation")
	}

	if status != http.StatusNoContent {
		return generic.CreateErrorFromResponse(body, status)
	}
	return nil
}

func (d *deviceControl) GetCollectionOfBulkOperation(query OperationQuery, pageSize int) (*BulkOperationCollection, *generic.Error) {
	return d.FindBulkOperationCollection(query, pageSize)
}

func (d *deviceControl) UpdateBulkOperation(bulkOperationID string, operation *UpdateBulkOperation) (*BulkOperation, *generic.Error) {
	if len(bulkOperationID) == 0 {
		return nil, generic.ClientError("Updating bulkOperation without a bulkOperationID is not allowed", "UpdateBulkOperation")
	}

	bytes, err := json.Marshal(operation)
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while marshalling update model: %s", err), "UpdateBulkOperation")
	}

	body, status, err := d.client.Put(fmt.Sprintf("%v/%v", d.basePathBulkOperations, bulkOperationID), bytes, generic.ContentTypeHeader(BULK_OPERATION_CONTENT_TYPE))
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while updating bulkOperation: %s", err), "UpdateBulkOperation")
	}

	if status != http.StatusOK {
		return nil, generic.CreateErrorFromResponse(body, status)
	}

	bulkOperation := &BulkOperation{}
	if err := generic.ObjectFromJson(body, bulkOperation); err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while marshalling response: %s", err), "UpdateBulkOperation")
	}
	return bulkOperation, nil
}

func (d *deviceControl) NextPage(c *OperationCollection) (*OperationCollection, *generic.Error) {
	return d.getPage(c.Next)
}

func (d *deviceControl) PreviousPage(c *OperationCollection) (*OperationCollection, *generic.Error) {
	return d.getPage(c.Prev)
}

func (d *deviceControl) getPage(reference string) (*OperationCollection, *generic.Error) {
	if reference == "" {
		log.Print("No page reference given. Returning nil.")
		return nil, nil
	}

	nextUrl, err := url.Parse(reference)
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Unparsable URL given for page reference: '%s'", reference), "GetPage")
	}

	collection, genErr := d.getCommon(fmt.Sprintf("%s?%s", nextUrl.Path, nextUrl.RawQuery))
	if genErr != nil {
		return nil, genErr
	}

	if len(collection.Operations) == 0 {
		log.Print("Returned collection is empty. Returning nil.")
		return nil, nil
	}

	return collection, nil
}

func (d *deviceControl) getCommon(path string) (*OperationCollection, *generic.Error) {
	body, status, err := d.client.Get(path, generic.EmptyHeader())

	if status != http.StatusOK {
		return nil, generic.CreateErrorFromResponse(body, status)
	}

	var result OperationCollection
	if len(body) > 0 {
		err = generic.ObjectFromJson(body, &result)
		if err != nil {
			return nil, generic.ClientError(fmt.Sprintf("Error while parsing response JSON: %s", err.Error()), "GetCollection")
		}
	} else {
		return nil, generic.ClientError("Response body was empty", "GetCollection")
	}

	return &result, nil
}

func (d *deviceControl) getCommonBulkCollection(path string) (*BulkOperationCollection, *generic.Error) {
	body, status, err := d.client.Get(path, generic.EmptyHeader())

	if status != http.StatusOK {
		return nil, generic.CreateErrorFromResponse(body, status)
	}

	var result BulkOperationCollection
	if len(body) > 0 {
		err = generic.ObjectFromJson(body, &result)
		if err != nil {
			return nil, generic.ClientError(fmt.Sprintf("Error while parsing response JSON: %s", err.Error()), "GetBulkCollection")
		}
	} else {
		return nil, generic.ClientError("Response body was empty", "GetCollection")
	}

	return &result, nil
}

func (d *deviceControl) FindOperationCollection(query OperationQuery, pageSize int) (*OperationCollection, *generic.Error) {
	queryParams := &url.Values{}
	query.QueryParams(queryParams)

	if len(*queryParams) == 0 {
		return nil, generic.ClientError("No filter set", "FindOperationCollection")
	}
	var err error
	if pageSize != 0 {
		err = generic.PageSizeParameter(pageSize, queryParams)
	}
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while building pageSize parameter to fetch measurements: %s", err.Error()), "FindMeasurements")
	}
	return d.getCommon(fmt.Sprintf("%s?%s", d.basePathOperations, queryParams.Encode()))
}

func (d *deviceControl) FindBulkOperationCollection(query OperationQuery, pageSize int) (*BulkOperationCollection, *generic.Error) {
	queryParams := &url.Values{}
	query.QueryParams(queryParams)

	if len(*queryParams) == 0 {
		return nil, generic.ClientError("No filter set", "FindBulkOperationCollection")
	}

	err := generic.PageSizeParameter(pageSize, queryParams)
	if err != nil {
		return nil, generic.ClientError(fmt.Sprintf("Error while building pageSize parameter to fetch measurements: %s", err.Error()), "FindMeasurements")
	}
	return d.getCommonBulkCollection(fmt.Sprintf("%s?%s", d.basePathBulkOperations, queryParams.Encode()))
}

func parseOperationResponse(body []byte) (*Operation, *generic.Error) {
	var operation Operation
	if len(body) > 0 {
		err := generic.ObjectFromJson(body, &operation)
		if err != nil {
			return nil, generic.ClientError(fmt.Sprintf("Error while parsing response JSON: %s", err.Error()), "ResponseParser")
		}
	} else {
		return nil, generic.ClientError("Response body was empty", "GetOperation")
	}
	return &operation, nil
}

func parseBulkOperationResponse(body []byte) (*BulkOperation, *generic.Error) {
	var operation BulkOperation
	if len(body) > 0 {
		err := generic.ObjectFromJson(body, &operation)
		if err != nil {
			return nil, generic.ClientError(fmt.Sprintf("Error while parsing response JSON: %s", err.Error()), "ResponseParser")
		}
	} else {
		return nil, generic.ClientError("Response body was empty", "GetBulkOperation")
	}
	return &operation, nil
}
