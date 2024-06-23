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

		if arrayTexto[_PARAMETRO_FUNCION] == _PARAMETRO_ENTRADA_AGREGAR && len(arrayTexto) == 2 {
			informacionGeneral.AgregarArchivo(arrayTexto[_PARAMETRO_ARCHIVO])

		} else if arrayTexto[_PARAMETRO_FUNCION] == _PARAMETRO_VER_VISITANTES && len(arrayTexto) == 3 {
			informacionGeneral.VerVisitantes(arrayTexto[_IP_DESDE], arrayTexto[_IP_HASTA])

		} else if arrayTexto[_PARAMETRO_FUNCION] == _PARAMETRO_VER_MAS_VISITADOS && len(arrayTexto) == 2 {
			informacionGeneral.VerMasVisitados(arrayTexto[_PARAMETRO_CANTIDAD])

		} else {
			fmt.Fprintf(os.Stderr, "%s%s\n", ERROR, arrayTexto[_PARAMETRO_FUNCION])
		}
	}
}

func main() {
	ProcesarArchivos()
}
