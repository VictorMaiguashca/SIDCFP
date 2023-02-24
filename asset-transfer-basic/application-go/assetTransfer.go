/*
Copyright 2020 IBM All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"fmt"
	"strings"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

func main() {
	err := os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
	if err != nil {
		log.Fatalf("Error setting DISCOVERY_AS_LOCALHOST environemnt variable: %v", err)
	}

	wallet, err := gateway.NewFileSystemWallet("wallet")
	if err != nil {
		log.Fatalf("Failed to create wallet: %v", err)
	}

	if !wallet.Exists("appUser") {
		err = populateWallet(wallet)
		if err != nil {
			log.Fatalf("Failed to populate wallet contents: %v", err)
		}
	}

	ccpPath := filepath.Join(
		"..",
		"..",
		"test-network",
		"organizations",
		"peerOrganizations",
		"org1.example.com",
		"connection-org1.yaml",
	)

	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(wallet, "appUser"),
	)
	if err != nil {
		log.Fatalf("Failed to connect to gateway: %v", err)
	}
	defer gw.Close()

	network, err := gw.GetNetwork("mychannel")
	if err != nil {
		log.Fatalf("Failed to get network: %v", err)
	}

	contract := network.GetContract("basic")

	fmt.Println("------------- Inicio de la aplicacion -------------")

	fmt.Println("--> Ejecutando funcion InitLedger:")
	fmt.Println("	Creando el conjunto inicial de registros del ledger")
	result, err := contract.SubmitTransaction("InicializarLibro")
	if err != nil {
		log.Fatalf("Failed to Submit transaction: %v", err)
	}
	log.Println(string(result))
	
	var ruc, banco, mMax, cta, rSoc, tCon, tGas string

	for true {
		fmt.Println("Ingrese una opcion:")
		fmt.Println("1.- Registrar Aportante")
		fmt.Println("2.- Registrar Proveedor")
		fmt.Println("3.- Registrar Aporte")
		fmt.Println("4.- Registrar Pago")
		fmt.Println("5.- Consultar Registro")
		fmt.Println("6.- Eliminar Registro")
		fmt.Println("7.- Consultar Historial")
		fmt.Println("8.- Mostrar todos los Registros")
		fmt.Println("9.-Salir")
		var op int
		fmt.Scanf("%d", &op)
		switch op{
			case 1:
				log.Println("--> Registrar Aportante")

				ruc = getRuc()
				banco = getBanco()
				mMax = "0"
				cta = getCta()
				rSoc = getRSoc()
				tCon = "0"
				tGas = "0"

				result, err = contract.SubmitTransaction("CrearActivo", "APORTANTE", ruc, banco, mMax, cta, rSoc, tCon, tGas)
				if err != nil {
					log.Println("Failed to Submit transaction: %v", err)
				}
				fmt.Println("Registro De Aportante Correcto")

			case 2:
				log.Println("--> Registrar Proveedor")
				
				ruc = getRuc()
				banco = getBanco()
				mMax = "0"
				cta = getCta()
				rSoc = getRSoc()
				tCon = "0"
				tGas = "0"

				result, err = contract.SubmitTransaction("CrearActivo", "PROVEEDOR", ruc, banco, mMax, cta, rSoc, tCon, tGas)
				if err != nil {
					log.Println("Failed to Submit transaction: %v", err)
				}

				fmt.Println("Registro De Proveedor Correcto")

			case 3:
				log.Println("--> Registrar Aporte")
				
				fmt.Println("Aportante:")
				rSoc = getRuc()
				banco = "0"
				mMax = "0"
				fmt.Println("Organizacion Politica:")
				cta = getRuc()
				tCon = getTCont()
				tGas = "0"
				ruc = getIdConsulta()

				_, err = contract.SubmitTransaction("TransferirActivoAporte", cta, tCon)
				if err != nil {
					log.Println("Failed to Submit transaction: %v", err)
				}

				_, err = contract.SubmitTransaction("TransferirActivoAportante", rSoc, tCon)
				if err != nil {
					log.Println("Failed to Submit transaction: %v", err)
				}

				result, err = contract.SubmitTransaction("CrearActivo", "APORTE", ruc, banco, mMax, cta, rSoc, tCon, tGas)
				if err != nil {
					log.Println("Failed to Submit transaction: %v", err)
				}

				fmt.Println("Registro De Aporte Correcto")

			case 4:
				log.Println("--> Registrar Pago")
				
				fmt.Println("Organizacion Politica:")
				rSoc = getRuc()
				banco = "0"
				mMax = "0"
				fmt.Println("Proveedor:")
				cta = getRuc()
				tCon = "0"
				tGas = getTGas()
				ruc = getIdConsulta()

				_, err = contract.SubmitTransaction("TransferirActivoPago", rSoc, tGas)
				if err != nil {
					log.Println("Failed to Submit transaction: %v", err)
				}

				_, err = contract.SubmitTransaction("TransferirActivoProveedor", cta, tGas)
				if err != nil {
					log.Println("Failed to Submit transaction: %v", err)
				}

				result, err = contract.SubmitTransaction("CrearActivo", "PAGO", ruc, banco, mMax, cta, rSoc, tCon, tGas)
				if err != nil {
					log.Println("Failed to Submit transaction: %v", err)
				}

				fmt.Println("Registro De Pago Correcto")

			case 5:
				log.Println("--> Consultar Registro")

				ruc = getRuc()

				result, err = contract.EvaluateTransaction("ActivoExiste", ruc)
				if result! = true {
					log.Println("No existe el registro buscado")
				}
				else{
					fmt.Println("Registro Encontrado.")
					result, err = contract.EvaluateTransaction("ConsultarActivo", ruc)
					if err != nil {
						log.Println("Failed to evaluate transaction: %v\n", err)
					}
				showResults(string(result))
				}

			case 6:
				log.Println("--> Eliminar Registro")
				ruc = getRuc()

				result, err = contract.EvaluateTransaction("ActivoExiste", ruc)
				if result! = true {
					log.Println("No existe el registro buscado")
				}
				else{
					result, err = contract.SubmitTransaction("EliminarActivo", ruc)
					if err != nil {
						log.Println("Failed to Submit transaction: %v", err)
					}
					fmt.Println("Registro De Pago Correcto")
				}

			case 7:
				log.Println("--> Consultar Historial")
				
				ruc = getRuc()

				result, err = contract.SubmitTransaction("HistorialdeActivo", ruc)
				if err != nil {
					log.Println("Failed to Submit transaction: %v", err)
				}
				showResults(string(result))

			case 8:
				log.Println("--> Mostrar todos los registros:")
				result, err = contract.EvaluateTransaction("ConsultarTodoslosActivos")
				if err != nil {
					log.Println("Failed to evaluate transaction: %v", err)
				}
				showResults(string(result))

			case 9:
				log.Println("============ Aplicacion Finalizada ============")
				os.Exit(0)

			default:
				log.Println("Opcion ingresada no encontrada.")
		}
	}
}

func populateWallet(wallet *gateway.Wallet) error {
	log.Println("============ Populating wallet ============")
	credPath := filepath.Join(
		"..",
		"..",
		"test-network",
		"organizations",
		"peerOrganizations",
		"org1.example.com",
		"users",
		"User1@org1.example.com",
		"msp",
	)

	certPath := filepath.Join(credPath, "signcerts", "cert.pem")
	// read the certificate pem
	cert, err := ioutil.ReadFile(filepath.Clean(certPath))
	if err != nil {
		return err
	}

	keyDir := filepath.Join(credPath, "keystore")
	// there's a single file in this dir containing the private key
	files, err := ioutil.ReadDir(keyDir)
	if err != nil {
		return err
	}
	if len(files) != 1 {
		return fmt.Errorf("keystore folder should have contain one file")
	}
	keyPath := filepath.Join(keyDir, files[0].Name())
	key, err := ioutil.ReadFile(filepath.Clean(keyPath))
	if err != nil {
		return err
	}

	identity := gateway.NewX509Identity("Org1MSP", string(cert), string(key))

	return wallet.Put("appUser", identity)
}

func getRuc() (string){
	var ruc string
	fmt.Print("Ingrese Ruc: ")
	fmt.Scanf("%s", &ruc)
	return ruc 
}

func getBanco()(string){
	var banco string
	fmt.Print("Ingrese Banco: ")
	fmt.Scanf("%s", &banco)
	return banco 
}

func getMMax()(string){
	var monMax string
	fmt.Print("Ingrese Monto Maximo: ")
	fmt.Scanf("%s", &monMax)
	return monMax 
}

func getCta()(string){
	var numCta string
	fmt.Print("Ingrese Numero De Cuenta: ")
	fmt.Scanf("%s", &numCta)
	return numCta
}

func getRSoc()(string){
	var razSoc string
	fmt.Print("Ingrese Razon Social: ")
	fmt.Scanf("%s", &razSoc)
	return razSoc 
}

func getTCont()(string){
	var totContrib string
	fmt.Print("Ingrese Total Contribuciones: ")
	fmt.Scanf("%s", &totContrib)
	return totContrib 
}

func getTGas()(string){
	var totGastos string
	fmt.Print("Ingrese el Total Gastos: ")
	fmt.Scanf("%s", &totGastos)
	return totGastos
}

func getIdConsulta()(string){
	var consulta string
	fmt.Print("Ingrese id para consultar la transaccion: ")
	fmt.Scanf("%s", &consulta )
	return consulta 
}


func showResults(result string){
	result = strings.Replace( result , ",\"" , ",\t\"" , -1)	
	result = strings.Replace( result , "," , ",\n" , -1)				

	fmt.Println(result)
}