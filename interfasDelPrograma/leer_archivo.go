package interfazArchivos

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	TDADiccionario "tdas/diccionario"
)

const (
	ESPACIO_VACIO               string = " "
	SEPARADOR_CODIGO            string = "."
	PARAMETRO_ENTRADA_AGREGAR   string = "agregar_archivo"
	PARAMETRO_VER_VISITANTES    string = "ver_visitantes"
	PARAMETRO_VER_MAS_VISITADOS string = "ver_mas_visitados"
	// _255_AL_CUBO                int    = 16581375
	// _255_AL_CUADRADO            int    = 65025
	PARAMETRO_FUNCION int    = 0
	VER_IP            int    = 0
	VER_ZONA_HORARIA  int    = 1
	VER_METODO        int    = 2
	VER_URL           int    = 3
	PARAMETRO_ARCHIVO int    = 0
	LAYOUT_PARSE      string = "2022-12-18T17:55:00-00:00"
	LAYOUT_PARSE2     string = "2006-01-02T15:04:05-07:00"
)

func comparacionNumerica(a, b uint32) int {
	if a > b {
		return 1
	}
	if a < b {
		return -1
	}
	return 0
}

type informacionUsuario struct {
	IP     string
	tiempo string
}

type informacionInterfaz struct {
	informacionIP   *TDADiccionario.DiccionarioOrdenado[uint32, informacionUsuario]
	informacionUrls *TDADiccionario.Diccionario[string, int]
}

func CrearInformacion(abb TDADiccionario.DiccionarioOrdenado[uint32, informacionUsuario], dic TDADiccionario.Diccionario[string, int]) EjecucionArchivos {
	return informacionInterfaz{&abb, &dic}
}

// AgregarArchivo implements EjecucionArchivos.
func (info informacionInterfaz) AgregarArchivo(ruta string) {
	archivo, err := os.Open(ruta)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error en comando "+PARAMETRO_ENTRADA_AGREGAR)
	}
	defer archivo.Close()

	scanner := bufio.NewScanner(archivo)
	for scanner.Scan() {
		lineaTexto := scanner.Text()
		lineaInformacion := strings.Split(lineaTexto, ESPACIO_VACIO)
		valorNumericoIP := transformarIP(lineaInformacion[VER_IP])
		informacionIP := informacionUsuario{
			IP:     lineaInformacion[VER_IP],
			tiempo: lineaInformacion[VER_ZONA_HORARIA],
		}
		contabilizarURLs(lineaInformacion[VER_URL], *info.informacionUrls)
		(*info.informacionIP).Guardar(valorNumericoIP, informacionIP)
	}
}

// VerMasVisitados implements EjecucionArchivos.
func (info informacionInterfaz) VerMasVisitados(int) {
	panic("unimplemented")
}

// VerVisitantes implements EjecucionArchivos.
func (info informacionInterfaz) VerVisitantes(string, string) {
	panic("unimplemented")
}

func LecturaArchivos() {
	scanner := bufio.NewScanner(os.Stdin)
	arbolIP := TDADiccionario.CrearABB[uint32, informacionUsuario](comparacionNumerica)
	DiccionarioURLs := TDADiccionario.CrearHash[string, int]()
	informacionGeneral := CrearInformacion(arbolIP, DiccionarioURLs)
	for scanner.Scan() {
		lineaTexto := scanner.Text()
		arrayEjecuciones := strings.Split(lineaTexto, ESPACIO_VACIO)

		if len(arrayEjecuciones) == 1 {
			fmt.Fprintf(os.Stderr, "Error en comando: "+arrayEjecuciones[PARAMETRO_FUNCION])
			return

		} else {
			if arrayEjecuciones[PARAMETRO_FUNCION] == PARAMETRO_ENTRADA_AGREGAR {
				informacionGeneral.AgregarArchivo(arrayEjecuciones[PARAMETRO_ARCHIVO])
				//AgregarArchivo(arrayEjecuciones[PARAMETRO_ARCHIVO], arbolIP, DiccionarioURLs)

			} else if arrayEjecuciones[PARAMETRO_FUNCION] == PARAMETRO_VER_VISITANTES {
				return

			} else if arrayEjecuciones[PARAMETRO_FUNCION] == PARAMETRO_VER_MAS_VISITADOS {
				return

			}
		}
	}
}

// func AgregarArchivo(ruta string, arbolIP TDADiccionario.DiccionarioOrdenado[uint32, informacionUsuario], dicionarioURL TDADiccionario.Diccionario[string, int]) {
// 	archivo, err := os.Open(ruta)
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "Error en comando "+PARAMETRO_ENTRADA_AGREGAR)
// 	}
// 	defer archivo.Close()

// 	scanner := bufio.NewScanner(archivo)
// 	for scanner.Scan() {
// 		lineaTexto := scanner.Text()
// 		lineaInformacion := strings.Split(lineaTexto, ESPACIO_VACIO)
// 		valorNumericoIP := transformarIP(lineaInformacion[VER_IP])
// 		informacionIP := informacionUsuario{
// 			IP:     lineaInformacion[VER_IP],
// 			tiempo: lineaInformacion[VER_ZONA_HORARIA],
// 		}
// 		contabilizarURLs(lineaInformacion[VER_URL], dicionarioURL)
// 		arbolIP.Guardar(valorNumericoIP, informacionIP)
// 	}
// }

func contabilizarURLs(urlVisitado string, dicURL TDADiccionario.Diccionario[string, int]) {
	if dicURL.Pertenece(urlVisitado) {
		dicURL.Guardar(urlVisitado, dicURL.Obtener(urlVisitado)+1)
	} else {
		dicURL.Guardar(urlVisitado, 0)
	}
}

func transformarIP(ip string) uint32 {
	IPstring := strings.Split(ip, SEPARADOR_CODIGO)
	IPint := make([]int, 4)
	for i := 0; i < 4; i++ {
		valor, _ := strconv.Atoi(IPstring[i])
		IPint[i] = valor
	}
	return uint32(IPint[0]*16777216 + IPint[1]*65536 + IPint[2]*255 + IPint[3])
}

// type datos struct {
// 	//una structura para el caso de ataques de denegacion
// 	ip     string
// 	tiempo []string
// }

// func LeerStdin() {
// 	lectura := bufio.NewScanner(os.Stdin)
// 	cola := TDA.CrearColaEnlazada[[]string]()
// 	for lectura.Scan() {
// 		parametros := strings.Split(lectura.Text(), ESPACIO)
// 		agregarArchivo(parametros[PARAMETRO_ARCHIVO], cola)
// 		verDoS(cola)
// 	}

// }

// func agregarArchivo(ruta string, cola TDA.Cola[[]string]) {
// 	archivo, _ := os.Open(ruta)
// 	linea := bufio.NewScanner(archivo)
// 	for linea.Scan() {
// 		cola.Encolar(strings.Split(linea.Text(), ESPACIO))
// 	}
// 	archivo.Close()
// }

// func verDoS(cola TDA.Cola[[]string]) {
// 	//en esta funcion vamos a encontrar los ataques de denegacion
// 	//contador := 0
// 	colaAux := TDA.CrearColaEnlazada[datos]()
// 	for !cola.EstaVacia() {
// 		info := cola.Desencolar()
// 		if colaAux.EstaVacia() {
// 			colaAux.Encolar(datos{info[VER_IP], []string{info[VER_ZONA_HORARIA]}})
// 		} else {
// 			datoAux := colaAux.Desencolar()
// 			fmt.Println(datoAux)
// 			//t1, _ := time.Parse(time.RFC3339, datos[VER_ZONA_HORARIA])
// 		}

// 	}
// }
