package clients

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/MarkoLuna/EmployeeService/pkg/models"
	"github.com/MarkoLuna/EmployeeService/pkg/utils"
)

var (
	employeeConsumerHost = utils.GetEnv("EMPLOYEE_CONSUMER_HOST", "http://localhost:8081")
)

type EmployeeConsumerServiceClientImpl struct {
	client      http.Client
	serviceHost string
}

func NewEmployeeConsumerServiceClient(client http.Client) EmployeeConsumerServiceClient {
	service := &EmployeeConsumerServiceClientImpl{client, employeeConsumerHost}

	return service
}

func (es EmployeeConsumerServiceClientImpl) Create(e models.Employee) (models.Employee, error) {
	jsonStr, err := json.Marshal(e)

	var employee models.Employee
	if err != nil {
		return employee, err
	}

	req, err := http.NewRequest("POST", es.serviceHost+"/api/employee/", bytes.NewBuffer(jsonStr))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Request failed with status: %s", resp.Status)
		return employee, errors.New("Error status from service")
	}

	employee = models.Employee{}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return employee, err
	}

	err = json.Unmarshal(body, &employee)

	if err != nil {
		log.Fatal(err)
		return employee, err
	}

	log.Printf("Response Body: %s\n", body)
	return employee, nil
}

func (es EmployeeConsumerServiceClientImpl) FindAll() ([]models.Employee, error) {

	req, err := http.NewRequest("GET", es.serviceHost+"/api/employee/", nil)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resp, err := es.client.Do(req)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Request failed with status: %s", resp.Status)
		return nil, errors.New("Error status from service")
	}

	employeesResponse := make([]models.Employee, 0)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	err = json.Unmarshal(body, &employeesResponse)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	log.Printf("Response Body: %s\n", body)
	return employeesResponse, nil

}

func (es EmployeeConsumerServiceClientImpl) FindById(ID string) (models.Employee, error) {

	var employee models.Employee
	req, err := http.NewRequest("GET", es.serviceHost+"/api/employee/"+ID, nil)
	if err != nil {
		log.Fatal(err)
		return employee, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resp, err := es.client.Do(req)
	if err != nil {
		log.Fatal(err)
		return employee, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Request failed with status: %s", resp.Status)

		if resp.StatusCode == http.StatusNotFound {
			return employee, errors.New("Employee Not Found")
		}

		return employee, errors.New("Error status from service")
	}

	employee = models.Employee{}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return employee, err
	}

	err = json.Unmarshal(body, &employee)

	if err != nil {
		log.Fatal(err)
		return employee, err
	}

	log.Printf("Response Body: %s\n", body)
	return employee, nil
}

func (es EmployeeConsumerServiceClientImpl) DeleteById(ID string) error {
	req, err := http.NewRequest("DELETE", es.serviceHost+"/api/employee/"+ID, nil)
	if err != nil {
		log.Fatal(err)
		return err
	}

	req.Header.Add("Accept", "application/json")

	resp, err := es.client.Do(req)
	if err != nil {
		log.Fatal(err)
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {

		if resp.StatusCode == http.StatusNotFound {
			return nil
		}

		log.Printf("Request failed with status: %s", resp.Status)
		return errors.New("Error status from service")
	}

	return nil
}

func (es EmployeeConsumerServiceClientImpl) Update(e models.Employee) (models.Employee, error) {
	jsonStr, err := json.Marshal(e)

	var employee models.Employee
	if err != nil {
		return employee, err
	}

	req, err := http.NewRequest("PUT", es.serviceHost+"/api/employee/"+e.Id, bytes.NewBuffer(jsonStr))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return employee, errors.New("Employee Not Found")
		}

		log.Printf("Request failed with status: %s", resp.Status)
		return employee, errors.New("Error status from service")
	}

	employee = models.Employee{}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return employee, err
	}

	err = json.Unmarshal(body, &employee)

	if err != nil {
		log.Fatal(err)
		return employee, err
	}

	log.Printf("Response Body: %s\n", body)
	return employee, nil
}
