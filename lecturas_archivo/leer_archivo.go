package lecturas_archivo

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	TDA "tdas/cola"
)

const (
	SEPARADOR_ARCHIVO           = "	"
	SEPARADOR_ENTRADA           = " "
	PARAMETRO_ENTRADA_AGREGAR   = "agregar_archivo"
	PARAMETRO_VER_VISITANTES    = "ver_visitantes"
	PARAMETRO_VER_MAS_VISITADOS = "ver_mas_visitados"
	VER_IP                      = 0
	VER_ZONA_HORARIA            = 1
	VER_URL                     = 3
	PARAMETRO_ARCHIVO           = 1
	LAYOUT_PARSE                = "2022-12-18T17:55:00-00:00"
	LAYOUT_PARSE2               = "2006-01-02T15:04:05-07:00"
)

func LeerStdin() {
	lectura := bufio.NewScanner(os.Stdin)
	cola := TDA.CrearColaEnlazada[[]string]()
	for lectura.Scan() {
		parametros := strings.Split(lectura.Text(), SEPARADOR_ENTRADA)
		agregarArchivo(parametros[PARAMETRO_ARCHIVO], cola)
		verDoS(cola)
	}

}

func agregarArchivo(ruta string, cola TDA.Cola[[]string]) {
	archivo, _ := os.Open(ruta)
	linea := bufio.NewScanner(archivo)
	for linea.Scan() {
		cola.Encolar(strings.Split(linea.Text(), SEPARADOR_ARCHIVO))
	}
	archivo.Close()
}

type datos struct {
	//una structura para el caso de ataques de denegacion
	ip     string
	tiempo []string
}

func verDoS(cola TDA.Cola[[]string]) {
	//en esta funcion vamos a encontrar los ataques de denegacion
	//contador := 0
	colaAux := TDA.CrearColaEnlazada[datos]()
	for !cola.EstaVacia() {
		info := cola.Desencolar()
		if colaAux.EstaVacia() {
			colaAux.Encolar(datos{info[VER_IP], []string{info[VER_ZONA_HORARIA]}})
		} else {
			datoAux := colaAux.Desencolar()
			fmt.Println(datoAux)
			//t1, _ := time.Parse(time.RFC3339, datos[VER_ZONA_HORARIA])
		}

	}
}
