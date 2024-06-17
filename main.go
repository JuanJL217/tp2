package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	programa "tp2/interfasDelPrograma"
)

const (
	_ESPACIO_VACIO               string = " "
	_SEPARADOR_CODIGO            string = "."
	_PARAMETRO_ENTRADA_AGREGAR   string = "agregar_archivo"
	_PARAMETRO_VER_VISITANTES    string = "ver_visitantes"
	_PARAMETRO_VER_MAS_VISITADOS string = "ver_mas_visitados"
	ERROR                        string = "Error en comando: "
	_PARAMETRO_FUNCION           int    = 0
	_PARAMETRO_ARCHIVO           int    = 0
	_IP_DESDE                    int    = 1
	_IP_HASTA                    int    = 2
	_PARAMETRO_CANTIDAD          int    = 1
	pruebas                             = "./pruebas_analog/test05.log"
)

func ProcesarArchivos() {
	scanner := bufio.NewScanner(os.Stdin)
	informacionGeneral := programa.CrearAnalisisDeArchivos()

	for scanner.Scan() {
		lineaTexto := scanner.Text()
		arrayEjecuciones := strings.Split(lineaTexto, _ESPACIO_VACIO)

		if len(arrayEjecuciones) == 1 {
			fmt.Fprintf(os.Stderr, ERROR+arrayEjecuciones[_PARAMETRO_FUNCION])
			return

		} else {

			if arrayEjecuciones[_PARAMETRO_FUNCION] == _PARAMETRO_ENTRADA_AGREGAR {
				informacionGeneral.AgregarArchivo(arrayEjecuciones[_PARAMETRO_ARCHIVO])

			} else if arrayEjecuciones[_PARAMETRO_FUNCION] == _PARAMETRO_VER_VISITANTES {
				informacionGeneral.VerVisitantes(arrayEjecuciones[_IP_DESDE], arrayEjecuciones[_IP_HASTA])

			} else if arrayEjecuciones[_PARAMETRO_FUNCION] == _PARAMETRO_VER_MAS_VISITADOS {
				informacionGeneral.VerMasVisitados(arrayEjecuciones[_PARAMETRO_CANTIDAD])

			}
		}
	}
}

func main() {
	ProcesarArchivos()
}
