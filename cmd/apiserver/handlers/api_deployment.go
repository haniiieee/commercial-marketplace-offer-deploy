package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/api"
	"github.com/microsoft/commercial-marketplace-offer-deploy/pkg/operations"
	"gorm.io/gorm"
)

const DeploymenIdParameterName = "deploymentId"

type InvokeOperationDeploymentHandler func(int, api.InvokeDeploymentOperation, *gorm.DB) (*api.InvokedOperation, error)

func GetDeployment(c echo.Context) error {
	return c.JSON(http.StatusOK, "")
}

func InvokeOperation(c echo.Context, db *gorm.DB) error {
	deploymentId, err := strconv.Atoi(c.Param(DeploymenIdParameterName))

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("%s in route was not an int", DeploymenIdParameterName))
	}

	var operation api.InvokeDeploymentOperation
	err = c.Bind(&operation)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	log.Printf("Operation deserialized \n %v", operation)

	operationHandler := CreateOperationHandler(operation)

	if operationHandler == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "There was op OperationHandler registered for the invoked operation")
	}
	res, err := operationHandler(deploymentId, operation, db)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, res)
}

func CreateOperationHandler(operation api.InvokeDeploymentOperation) InvokeOperationDeploymentHandler {
	operationType := operations.GetOperationFromString(*operation.Name)

	switch operationType {
	case operations.DryRunDeploymentOperation:
		return CreateDryRun
	case operations.StartDeploymentOperation:
		return StartDeployment
	default:
		return nil
	}
}

func ListDeployments(c echo.Context) error {
	return c.JSON(http.StatusOK, "")
}

func UpdateDeployment(c echo.Context) error {
	return c.JSON(http.StatusOK, "")
}
