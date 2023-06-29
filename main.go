package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/go-sql-driver/mysql"
)

type DocumentType struct {
	ID     int    `json:"id"`
	Nombre string `json:"nombre"`
	Estado int    `json:"estado"`
}

type BodyData interface {
    
}

type ResponseBody struct {
    Message string `json:"message"`
    Data BodyData `json:"data"` 
}


func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	
    httpVerb := request.HTTPMethod
    fmt.Println("The HTTP verb is:", httpVerb)
    fmt.Println("The HTTP verb is:", request)

    documentTypes, err := queryDocumentTypes()
	if err != nil {
		log.Println("Error querying document types:", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Internal Server Error",
		}, nil
	}

    fmt.Println(documentTypes)

    searchDocumentTypes, err := searchDocumentTypes("ex")
	if err != nil {
		log.Println("Error querying document types:", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Internal Server Error",
		}, nil
	}

    //create response

	responseBody := ResponseBody{
		Message: "Query executed successfully",
		Data:   searchDocumentTypes,
	}

	responseJSON, err := json.Marshal(responseBody)
	if err != nil {
		log.Println("Error marshaling response body:", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Internal Server Error",
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseJSON),
	}, nil
}

func makeDatabaseConnection() (*sql.DB, error) {
	// Configurar los detalles de la conexión a la base de datos
	username := "admin_desa"
	password := "ImSign2023*"
	host := "imsigndesarrollo.cn8sqbaxyaac.us-east-1.rds.amazonaws.com"
	port := "3306"

	// Crear la cadena de conexión
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/", username, password, host, port)

	// Abrir la conexión a la base de datos
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func queryDocumentTypes() (BodyData, error) {
	db, err := makeDatabaseConnection()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Realizar la consulta a la tabla "tipo_documento"
	rows, err := db.Query("SELECT id, nombre, estado FROM Imsign.tipo_documento")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	documentTypes := []DocumentType{}
	for rows.Next() {
		var docType DocumentType
		err := rows.Scan(&docType.ID, &docType.Nombre, &docType.Estado)
		if err != nil {
			return nil, err
		}
		documentTypes = append(documentTypes, docType)
	}

	return documentTypes, nil
}


func searchDocumentTypes (queryText string) (BodyData, error){
    db, err := makeDatabaseConnection()
	if err != nil {
		return nil, err
	}
	defer db.Close()

    //query search by name
    query := "SELECT id, nombre, estado FROM Imsign.tipo_documento WHERE nombre LIKE CONCAT('%', ?, '%')"
    rows, err := db.Query(query,queryText)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	documentTypes := []DocumentType{}
	for rows.Next() {
		var docType DocumentType
		err := rows.Scan(&docType.ID, &docType.Nombre, &docType.Estado)
		if err != nil {
			return nil, err
		}
		documentTypes = append(documentTypes, docType)
	}

	return documentTypes, nil
}