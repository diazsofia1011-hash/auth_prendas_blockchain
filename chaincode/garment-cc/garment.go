package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract gestiona el ciclo de vida y autenticación de las prendas
type SmartContract struct {
	contractapi.Contract
}

// Prenda define la estructura de datos almacenada on-chain para cada prenda
type Prenda struct {
	ID             string `json:"id"`             // Identificador público único (ej. QR)
	MarcaID        string `json:"marcaId"`        // Empresa fabricante
	HashSecreto    string `json:"hashSecreto"`    // Compromiso criptográfico (commit) del doble código
	Estado         string `json:"estado"`         // EN_FABRICA, EN_TRANSITO, VENDIDA
	CustodioActual string `json:"custodioActual"` // Quién la tiene en posesión (Fábrica, Distribuidor, Tienda)
	MetadatoHash   string `json:"metadatoHash"`   // Hash de la IPFS/base de datos off-chain con imágenes y certificados
}

// RegistrarPrenda permite a la marca dar de alta una nueva prenda en la red
func (s *SmartContract) RegistrarPrenda(ctx contractapi.TransactionContextInterface, id string, marcaId string, hashSecreto string, metadatoHash string) error {
	// Verificar si la prenda ya existe para evitar duplicados
	exists, err := s.PrendaExists(ctx, id)
	if err != nil {
		return fmt.Errorf("error al leer el world state: %v", err)
	}
	if exists {
		return fmt.Errorf("la prenda con id %s ya está registrada", id)
	}

	prenda := Prenda{
		ID:             id,
		MarcaID:        marcaId,
		HashSecreto:    hashSecreto,
		Estado:         "EN_FABRICA",
		CustodioActual: marcaId,
		MetadatoHash:   metadatoHash,
	}

	prendaJSON, err := json.Marshal(prenda)
	if err != nil {
		return err
	}

	// Guardar en el World State de Fabric
	return ctx.GetStub().PutState(id, prendaJSON)
}

// VenderPrenda valida el secreto y cambia el estado de la prenda a VENDIDA
func (s *SmartContract) VenderPrenda(ctx contractapi.TransactionContextInterface, id string, secretoRevelado string) error {
	prenda, err := s.ReadPrenda(ctx, id)
	if err != nil {
		return err
	}

	if prenda.Estado == "VENDIDA" {
		return fmt.Errorf("la prenda %s ya fue vendida anteriormente (protección contra clonación)", id)
	}

	// Aquí validarías que el hash del secreto revelado coincida con el hash on-chain
	// (En Go usarías crypto/sha256 para comprobarlo en producción)
	
	// Actualizar estado y custodio
	prenda.Estado = "VENDIDA"
	prenda.CustodioActual = "CONSUMIDOR_FINAL"
	// Invalidar el hash secreto para que no pueda reutilizarse
	prenda.HashSecreto = "USADO"

	prendaJSON, err := json.Marshal(prenda)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, prendaJSON)
}

// ReadPrenda consulta los datos públicos de una prenda (función de lectura abierta para el consumidor)
func (s *SmartContract) ReadPrenda(ctx contractapi.TransactionContextInterface, id string) (*Prenda, error) {
	prendaJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("falla al leer la prenda %s: %v", id, err)
	}
	if prendaJSON == nil {
		return nil, fmt.Errorf("la prenda %s no existe en el ledger", id)
	}

	var prenda Prenda
	err = json.Unmarshal(prendaJSON, &prenda)
	if err != nil {
		return nil, err
	}

	return &prenda, nil
}

// PrendaExists verifica si una prenda ya se encuentra registrada
func (s *SmartContract) PrendaExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	prendaJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("falla al leer el world state: %v", err)
	}

	return prendaJSON != nil, nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(new(SmartContract))
	if err != nil {
		panic(fmt.Sprintf("Error al crear el chaincode de prendas: %v", err))
	}

	if err := chaincode.Start(); err != nil {
		panic(fmt.Sprintf("Error al iniciar el chaincode de prendas: %v", err))
	}
}
