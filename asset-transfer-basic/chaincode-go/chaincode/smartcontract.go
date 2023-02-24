package chaincode

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

// Asset describes basic details of what makes up a simple asset
type Asset struct {
	DocType             string `json:"docType"`
	RUC             string `json:"RUC"`
	IdBanco          int `json:"idBanco"`
	MontoMax           float64 `json:"montoMax"`
	NumeroDeCuenta           string `json:"numeroDeCuenta"`
	RazonSocial          string `json:"razonSocial"`
	TotalContribuciones  float64 `json:"totalContibuciones"`
	TotalGastos  float64 `json:"totalGastos"`
}

// Estructura para recibir el historico de un Asset
type HistoryQueryResult struct{
	Record             *Asset `json:"record"`
	TxId            string `json:"txId"`
	Timestamp             time.Time `json:"timestamp"`
	IsDelete             bool `json:"isDelete"`
}


// InitLedger adds a base set of assets to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	assets := []Asset{
		{DocType: "ORG-POL", RUC: "11223344556", IdBanco: 1,MontoMax: 5000.00, NumeroDeCuenta: "123456789", RazonSocial: "PSC", TotalContribuciones: 0.00, TotalGastos: 0.00},
		{DocType: "ORG-POL", RUC: "22334455667", IdBanco: 2,MontoMax: 5000.00, NumeroDeCuenta: "234567891", RazonSocial: "PSP", TotalContribuciones: 0.00, TotalGastos: 0.00},
		{DocType: "ORG-POL", RUC: "33445566778", IdBanco: 3,MontoMax: 5000.00, NumeroDeCuenta: "345678912", RazonSocial: "AVANZA", TotalContribuciones: 0.00, TotalGastos: 0.00},
		{DocType: "ORG-POL", RUC: "44556677889", IdBanco: 3,MontoMax: 5000.00, NumeroDeCuenta: "456789123", RazonSocial: "CREO", TotalContribuciones: 0.00, TotalGastos: 0.00},
		{DocType: "ORG-POL", RUC: "55667788990", IdBanco: 2,MontoMax: 5000.00, NumeroDeCuenta: "567891234", RazonSocial: "SUMA", TotalContribuciones: 0.00, TotalGastos: 0.00},
		{DocType: "ORG-POL", RUC: "66778899001", IdBanco: 1,MontoMax: 5000.00, NumeroDeCuenta: "678912345", RazonSocial: "MOVER", TotalContribuciones: 0.00, TotalGastos: 0.00},
	}

	for _, asset := range assets {
		assetJSON, err := json.Marshal(asset)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(asset.RUC, assetJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}

// CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, docType string, ruc string, idBanco int,montoMax float64, numeroDeCuenta string, razonSocial string, totalContibuciones float64, totalGastos float64) error {
	exists, err := s.AssetExists(ctx, ruc)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the asset %s already exists", ruc)
	}

	asset := Asset{
		DocType:		docType,
		RUC:             ruc,
		IdBanco:          idBanco,
		MontoMax:	montoMax,
		NumeroDeCuenta:           numeroDeCuenta,
		RazonSocial:          razonSocial,
		TotalContribuciones:	totalContibuciones,
		TotalGastos:	totalGastos,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(ruc, assetJSON)
}

// ReadAsset returns the asset stored in the world state with given ruc.
func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, ruc string) (*Asset, error) {
	assetJSON, err := ctx.GetStub().GetState(ruc)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", ruc)
	}

	var asset Asset
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

// UpdateAsset updates an existing asset in the world state with provided parameters.
func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, docType string, ruc string, idBanco int,montoMax float64, numeroDeCuenta string, razonSocial string, totalContibuciones float64, totalGastos float64) error {
	exists, err := s.AssetExists(ctx, ruc)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", ruc)
	}

	// overwriting original asset with new asset
	asset := Asset{
		DocType:		docType,
		RUC:             ruc,
		IdBanco:          idBanco,
		MontoMax:	montoMax,
		NumeroDeCuenta:           numeroDeCuenta,
		RazonSocial:          razonSocial,
		TotalContribuciones:	totalContibuciones,
		TotalGastos:	totalGastos,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(ruc, assetJSON)
}

// DeleteAsset deletes an given asset from the world state.
func (s *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, ruc string) error {
	exists, err := s.AssetExists(ctx, ruc)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", ruc)
	}

	return ctx.GetStub().DelState(ruc)
}

// AssetExists returns true when asset with given ID exists in world state
func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, ruc string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(ruc)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}

// TransferAsset updates the owner field of asset with given id in world state.
func (s *SmartContract) TransferAssetAporte(ctx contractapi.TransactionContextInterface, ruc string, tContribucion float64) error {
	asset, err := s.ReadAsset(ctx, ruc)
	if err != nil {
		return err
	}
	var total = asset.TotalContribuciones + tContribucion

	if total > asset.MontoMax{
		return fmt.Errorf("Total de contribuciones supera el maximo permitido")
	}

	asset.TotalContribuciones =  total
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(ruc, assetJSON)
}

// TransferAsset updates the owner field of asset with given id in world state.
func (s *SmartContract) TransferAssetAportante(ctx contractapi.TransactionContextInterface, ruc string, tContribucion float64) error {
	asset, err := s.ReadAsset(ctx, ruc)
	if err != nil {
		return err
	}

	asset.TotalContribuciones = asset.TotalContribuciones + tContribucion 
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(ruc, assetJSON)
}


func (s *SmartContract) TransferAssetPago(ctx contractapi.TransactionContextInterface, ruc string, tGastos float64) error {
	asset, err := s.ReadAsset(ctx, ruc)
	if err != nil {
		return err
	}
	var total = asset.TotalGastos + tGastos

	if total  > asset.TotalContribuciones {
		return fmt.Errorf("Total de gastos supera el Total de contribuciones")
	}

	asset.TotalGastos  = total 
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(ruc, assetJSON)
}

func (s *SmartContract) TransferAssetProveedor(ctx contractapi.TransactionContextInterface, ruc string, tGastos float64) error {
	asset, err := s.ReadAsset(ctx, ruc)
	if err != nil {
		return err
	}

	asset.TotalGastos  = asset.TotalGastos + tGastos 
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(ruc, assetJSON)
}



// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*Asset, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*Asset
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset Asset
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil
}

//Get history
func (s *SmartContract) GetAssetHistory(ctx contractapi.TransactionContextInterface, ruc string) ([]HistoryQueryResult, error){
	log.Print("GetAssetHistory: ID %v", ruc)

	resultsIterator, err:= ctx.GetStub().GetHistoryForKey(ruc)
	if err != nil{
		return nil, err
	}
	defer resultsIterator.Close()

	var records []HistoryQueryResult
	for resultsIterator.HasNext(){
		response, err  := resultsIterator.Next()
		if err != nil{
			return nil, err
		}

		var asset Asset
		if len(response.Value) > 0 {
			err=json.Unmarshal(response.Value, &asset)
			if err != nil{
				return nil, err
			}
		}else{
			asset=Asset{
				RUC: ruc,
			}
		}
		timestamp, err := ptypes.Timestamp(response.Timestamp)
		if err != nil{
			return nil, err
		}

		record := HistoryQueryResult{
			TxId:	response.TxId,
			Timestamp:	timestamp,
			Record:		&asset,
			IsDelete:	response.IsDelete,
		}
		records = append(records, record)
	}
	return records, nil
}
