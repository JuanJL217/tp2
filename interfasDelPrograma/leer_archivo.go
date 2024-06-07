package interfazArchivos

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	TDADiccionario "tdas/diccionario"
	TDAPila "tdas/pila"
)

const (
	_ESPACIO_VACIO             string  = " "
	_SEPARADOR_DIGITOS_IP      string  = "."
	_PARAMETRO_ENTRADA_AGREGAR string  = "agregar_archivo"
	_ERROR                     string  = "Error en comando "
	_CANT_MINIMA_DOS           int     = 5
	_SEGUNDO_MAXIMO            float64 = 0.2
	_DATOS_INSUFICIENTES       float64 = 0.0
	_VALOR_INICIAL             int     = 0
	_VER_IP                    int     = 0
	_VER_ZONA_HORARIA          int     = 1
	_VER_URL                   int     = 3
	_MAx_URLS_VISITADO         int     = 10
	_LAYOUT_PARSE              string  = "2022-12-18T17:55:00-00:00"
	_LAYOUT_PARSE2             string  = "2006-01-02T15:04:05-07:00"
)

func compararIPs(IP1, IP2 string) int {
	arrayIP1 := strings.Split(IP1, _SEPARADOR_DIGITOS_IP)
	arrayIP2 := strings.Split(IP2, _SEPARADOR_DIGITOS_IP)
	for i := 0; i < 4; i++ {
		intIP1, _ := strconv.Atoi(arrayIP1[i])
		intIP2, _ := strconv.Atoi(arrayIP2[i])
		if intIP1 > intIP2 {
			return 1
		}
		if intIP2 > intIP1 {
			return -1
		}
	}
	return 0
}

type informacionSesionUsuario struct {
	tiempo TDADiccionario.Diccionario[string, *TDAPila.Pila[[]string]]
	urls   TDADiccionario.Diccionario[string, int]
}

type informacionGeneral struct {
	informacionIPs  TDADiccionario.DiccionarioOrdenado[string, *informacionSesionUsuario]
	informacionUrls TDADiccionario.Diccionario[string, int]
}

func crearNuevoUsuario() *informacionSesionUsuario {
	return &informacionSesionUsuario{
		TDADiccionario.CrearHash[string, *TDAPila.Pila[[]string]](),
		TDADiccionario.CrearHash[string, int](),
	}
}

func CrearInformacionArchivos() EjecucionArchivos {
	return &informacionGeneral{
		TDADiccionario.CrearABB[string, *informacionSesionUsuario](compararIPs),
		TDADiccionario.CrearHash[string, int](),
	}
}

// AgregarArchivo implements EjecucionArchivos.
func (info informacionGeneral) AgregarArchivo(ruta string) {
	archivo, err := os.Open(ruta)
	if err != nil {
		fmt.Fprintf(os.Stderr, _ERROR+_PARAMETRO_ENTRADA_AGREGAR) // importar del main
	}
	defer archivo.Close()

	scanner := bufio.NewScanner(archivo)
	for scanner.Scan() {
		lineaTexto := scanner.Text()
		lineaInformacion := strings.Split(lineaTexto, _ESPACIO_VACIO)
		guardarInformacion(lineaInformacion[_VER_IP], lineaInformacion[_VER_URL], lineaInformacion[_VER_ZONA_HORARIA], info)
	}
}

func guardarInformacion(IP, url, tiempo string, info informacionGeneral) {
	reloj := tranformarTiempo(tiempo)
	if !info.informacionIPs.Pertenece(IP) {
		usuario := crearNuevoUsuario()
		pilaTiempos := TDAPila.CrearPilaDinamica[[]string]()
		pilaTiempos.Apilar(reloj)
		usuario.tiempo.Guardar(url, &pilaTiempos)
		usuario.urls.Guardar(url, _VALOR_INICIAL)
		info.informacionIPs.Guardar(IP, usuario)
	} else {
		usuario := info.informacionIPs.Obtener(IP)
		pilaTiempos := usuario.tiempo.Obtener(url)
		(*pilaTiempos).Apilar(reloj)
		usuario.urls.Guardar(IP, usuario.urls.Obtener(IP)+1)
		detectarDoS(usuario, IP, url)
	}
	contabilizarURLs(url, info)

	// reloj := tranformarTiempo(tiempo)
	// if !info.informacionIPs.Pertenece(IP) {
	// 	nuevoUsuario := crearNuevoUsuario()
	// 	pilaTiempos := TDAPila.CrearPilaDinamica[[]string]()
	// 	pilaTiempos.Apilar(reloj)
	// 	nuevoUsuario.tiempo.Guardar(url, pilaTiempos)
	// 	nuevoUsuario.urls.Guardar(url, 0)
	// 	info.informacionIPs.Guardar(IP, nuevoUsuario)
	// } else {
	// 	usuarioExistente := info.informacionIPs.Obtener(IP)
	// 	pilaTiempos := usuarioExistente.tiempo.Obtener(url)
	// 	pilaTiempos.Apilar(reloj)
	// 	usuarioExistente.tiempo.Guardar(url, pilaTiempos)
	// 	usuarioExistente.urls.Guardar(url, usuarioExistente.urls.Obtener(url)+1)
	// 	info.informacionIPs.Guardar(IP, usuarioExistente)
	// }
	// contabilizarURLs(url, info)
}

func detectarDoS(usuario *informacionSesionUsuario, IP, url string) {
	if usuario.urls.Obtener(url) >= _CANT_MINIMA_DOS {
		segundos := calcularTiempo(usuario, url)
		if segundos <= _SEGUNDO_MAXIMO {
			fmt.Println("DoS: " + IP)
		}
	}
}

func calcularTiempo(usuario *informacionSesionUsuario, url string) float64 {
	// desapilamos y apilamos
	return _DATOS_INSUFICIENTES
}

func tranformarTiempo(tiempo string) []string {
	//Transformar lo que necesitamos[dia, hora, minutos, segundos] ?
	return []string{tiempo}
}

func (info informacionGeneral) VerVisitantes(desdeIP, hastaIP string) {
	iterador := info.informacionIPs.IteradorRango(&desdeIP, &hastaIP)
	for iterador.HaySiguiente() {
		ip, _ := iterador.VerActual()
		fmt.Printf("\t" + ip)
		iterador.Siguiente()
	}
	fmt.Printf("OK")
}

func (info informacionGeneral) VerMasVisitados(mostrarCantidadUrls string) {
	cantidadTotal, err := strconv.Atoi(mostrarCantidadUrls)
	if err != nil {
		fmt.Fprintf(os.Stderr, _ERROR+"ver_mas_visitados") // importar del main
	}
	// Tengo que ver si hago la ejecucion acá o cuando lea el archivo
	// Como ya tengo el diccionario en contabilizarURLs(), tengo que aquí ordenarnos con un counting sort
	// Lo que ordenaré, será la cantidad de visitados que tenga cada URL
	// Como se ordenan de menor a mayor, pongo en una pila cada elemento
	// Y hare un ciclo for hasta que la Pila esté vacia, o cuanto el contador haya llegado a 10
	// Desapilando lo ltimo que apilé
	pila := TDAPila.CrearPilaDinamica[string]() //String no va a ser, tiene que ser una estructura
	// EjecuionDeOrdenamiento
	ordenamientoUrlVisitados(pila, info)
	contador := 0
	for !pila.EstaVacia() && contador <= cantidadTotal {
		fmt.Printf("\t" + pila.Desapilar())
		contador++
	}
	fmt.Printf("OK")
}

func ordenamientoUrlVisitados(pila TDAPila.Pila[string], info informacionGeneral) {
	// Hacer cositas
}

func contabilizarURLs(urlVisitado string, info informacionGeneral) {
	if !info.informacionUrls.Pertenece(urlVisitado) {
		info.informacionUrls.Guardar(urlVisitado, _VALOR_INICIAL)
	} else {
		info.informacionUrls.Guardar(urlVisitado, info.informacionUrls.Obtener(urlVisitado)+1)
	}
}

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
