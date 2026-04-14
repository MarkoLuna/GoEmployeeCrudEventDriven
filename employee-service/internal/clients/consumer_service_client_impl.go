package clients

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/MarkoLuna/EmployeeService/internal/models"
	"github.com/MarkoLuna/GoEmployeeCrudEventDriven/common/utils"
)

var (
	employeeConsumerHost = utils.GetEnv("EMPLOYEE_CONSUMER_HOST", "http://localhost:8081")
)

type EmployeeConsumerServiceClientImpl struct {
	client      http.Client
	serviceHost string
	token       string
}

func NewEmployeeConsumerServiceClient(client http.Client) EmployeeConsumerServiceClient {
	return &EmployeeConsumerServiceClientImpl{
		client:      client,
		serviceHost: employeeConsumerHost,
		token:       "",
	}
}

func (es *EmployeeConsumerServiceClientImpl) Create(e models.Employee) (models.Employee, error) {

	var employee models.Employee
	jsonStr, err := json.Marshal(e)
	if err != nil {
		return employee, err
	}

	req, err := http.NewRequest("POST", es.serviceHost+"/api/employee/", bytes.NewBuffer(jsonStr))
	if err != nil {
		return employee, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+es.token)

	resp, err := es.client.Do(req)
	if err != nil {
		return employee, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		log.Printf("Create: request failed with status: %s", resp.Status)
		return employee, errors.New("error status from service")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return employee, err
	}
	if err = json.Unmarshal(body, &employee); err != nil {
		return employee, err
	}

	log.Printf("Response Body: %s\n", body)
	return employee, nil
}

func (es *EmployeeConsumerServiceClientImpl) FindAll() ([]models.Employee, error) {

	req, err := http.NewRequest("GET", es.serviceHost+"/api/employee/", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+es.token)

	resp, err := es.client.Do(req)
	if err != nil {
		log.Printf("FindAll error: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("FindAll: request failed with status: %s", resp.Status)
		return nil, errors.New("error status from service")
	}

	employeesResponse := make([]models.Employee, 0)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(body, &employeesResponse); err != nil {
		return nil, err
	}

	log.Printf("Response Body: %s\n", body)
	return employeesResponse, nil
}

func (es *EmployeeConsumerServiceClientImpl) FindById(ID string) (models.Employee, error) {

	var employee models.Employee
	req, err := http.NewRequest("GET", es.serviceHost+"/api/employee/"+ID, nil)
	if err != nil {
		return employee, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+es.token)

	resp, err := es.client.Do(req)
	if err != nil {
		return employee, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		// handled below
	case http.StatusNotFound:
		return employee, errors.New("employee not found")
	default:
		log.Printf("FindById: request failed with status: %s", resp.Status)
		return employee, errors.New("error status from service")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return employee, err
	}
	if err = json.Unmarshal(body, &employee); err != nil {
		return employee, err
	}

	log.Printf("Response Body: %s\n", body)
	return employee, nil
}

func (es *EmployeeConsumerServiceClientImpl) DeleteById(ID string) error {

	req, err := http.NewRequest("DELETE", es.serviceHost+"/api/employee/"+ID, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+es.token)

	resp, err := es.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK, http.StatusNoContent, http.StatusNotFound:
		return nil
	default:
		log.Printf("DeleteById: request failed with status: %s", resp.Status)
		return errors.New("error status from service")
	}
}

func (es *EmployeeConsumerServiceClientImpl) Update(e models.Employee) (models.Employee, error) {

	var employee models.Employee
	jsonStr, err := json.Marshal(e)
	if err != nil {
		return employee, err
	}

	req, err := http.NewRequest("PUT", es.serviceHost+"/api/employee/"+e.Id, bytes.NewBuffer(jsonStr))
	if err != nil {
		return employee, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+es.token)

	resp, err := es.client.Do(req)
	if err != nil {
		return employee, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		// handled below
	case http.StatusNotFound:
		return employee, errors.New("employee not found")
	default:
		log.Printf("Update: request failed with status: %s", resp.Status)
		return employee, errors.New("error status from service")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return employee, err
	}
	if err = json.Unmarshal(body, &employee); err != nil {
		return employee, err
	}

	log.Printf("Response Body: %s\n", body)
	return employee, nil
}


