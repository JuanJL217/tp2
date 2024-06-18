package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	programa "tp2/interfasDelPrograma"
)

const (
	_ESPACIO_VACIO               = " "
	_SEPARADOR_CODIGO            = "."
	_PARAMETRO_ENTRADA_AGREGAR   = "agregar_archivo"
	_PARAMETRO_VER_VISITANTES    = "ver_visitantes"
	_PARAMETRO_VER_MAS_VISITADOS = "ver_mas_visitados"
	ERROR                        = "Error en comando "
	_PARAMETRO_FUNCION           = 0
	_PARAMETRO_ARCHIVO           = 1
	_IP_DESDE                    = 1
	_IP_HASTA                    = 2
	_PARAMETRO_CANTIDAD          = 1
)

func ProcesarArchivos() {
	scanner := bufio.NewScanner(os.Stdin)
	informacionGeneral := programa.CrearAnalisisDeArchivos()

	for scanner.Scan() {
		lineaTexto := scanner.Text()
		arrayTexto := strings.Split(lineaTexto, _ESPACIO_VACIO)

		if len(arrayTexto) == 1 {
			fmt.Fprintf(os.Stderr, ERROR+arrayTexto[_PARAMETRO_FUNCION])
		}

		if arrayTexto[_PARAMETRO_FUNCION] == _PARAMETRO_ENTRADA_AGREGAR {
			informacionGeneral.AgregarArchivo(arrayTexto[_PARAMETRO_ARCHIVO])

		} else if arrayTexto[_PARAMETRO_FUNCION] == _PARAMETRO_VER_VISITANTES {
			informacionGeneral.VerVisitantes(arrayTexto[_IP_DESDE], arrayTexto[_IP_HASTA])

		} else if arrayTexto[_PARAMETRO_FUNCION] == _PARAMETRO_VER_MAS_VISITADOS {
			informacionGeneral.VerMasVisitados(arrayTexto[_PARAMETRO_CANTIDAD])

		}
	}
}

func main() {
	ProcesarArchivos()
}
