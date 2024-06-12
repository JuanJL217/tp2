package interfazArchivos

import (
	"bufio"
	"cmp"
	"fmt"
	"os"
	"strconv"
	"strings"
	TDAHeap "tdas/cola_prioridad"
	TDADiccionario "tdas/diccionario"
	"time"
)

const (
	_TABULACION                = "\t"
	_SEPARADOR_DIGITOS_IP      = "."
	_PARAMETRO_ENTRADA_AGREGAR = "agregar_archivo"
	_ERROR                     = "Error en comando "
	_DENEGACION_SERVICIO       = "DoS"
	_ACEPTADO                  = "OK"
	_CANT_MINIMA_DOS           = 5
	_SEGUNDO_MAXIMO            = 2
	_VALOR_INICIAL             = 1
	_VER_IP                    = 0
	_VER_ZONA_HORARIA          = 1
	_VER_URL                   = 3
	_LAYOUT_PARSE              = "2006-01-02T15:04:05-07:00"
)

func CompararIPs(IP1, IP2 string) int {
	arrayIP1 := strings.Split(IP1, _SEPARADOR_DIGITOS_IP)
	arrayIP2 := strings.Split(IP2, _SEPARADOR_DIGITOS_IP)
	for i := 0; i < len(arrayIP1); i++ {
		intIP1, _ := strconv.Atoi(arrayIP1[i])
		intIP2, _ := strconv.Atoi(arrayIP2[i])
		if intIP1 > intIP2 {
			// si la ip1 es mayor
			return 1
		}
		if intIP2 > intIP1 {
			//si la ip2 es mayor
			return -1
		}
	}
	//caso que sean iguales ips
	return 0
}

type informacionArchivos struct {
	visitantes TDADiccionario.Diccionario[string, int] // casos de denegacion DoS y todos los Visitantes
	visitados  TDADiccionario.Diccionario[string, int] // ver los sitios (url) mas visitados
}

type sitiosVistados struct {
	nombre   string
	cantidad int
}

func CrearAnalisisLogs() EjecucionArchivos {
	info := new(informacionArchivos)
	info.visitantes = TDADiccionario.CrearHash[string, int]()
	info.visitados = TDADiccionario.CrearHash[string, int]()
	return info
}

func (info informacionArchivos) AgregarArchivo(ruta string) {
	info.procesarArchivo(ruta)    // O(n)
	arrDoS := make([][]string, 0) // todas las ip que son DoS
	agregarIPDoS(arrDoS, info.visitantes)

	// hacer countingSort y radixSort
}

func (info informacionArchivos) VerMasVisitados(dato string) {
	mininimo, _ := strconv.Atoi(dato)
	heap := TDAHeap.CrearHeap[sitiosVistados](compararNumeros)
	info.visitados.Iterar(func(clave string, valor int) bool {
		heap.Encolar(sitiosVistados{nombre: clave, cantidad: valor})
		return true
	})
	fmt.Println("Sitios mas visitados:")
	for i := 0; i < mininimo && !heap.EstaVacia(); i++ {
		url := heap.Desencolar()
		fmt.Println("\t", url.nombre, "-", url.cantidad)
	}
}

func (info informacionArchivos) VerVisitantes(inicio, fin string) {
	fmt.Println("Visitantes:")
	dicOrd := TDADiccionario.CrearABB[string, int](CompararIPs)
	info.visitantes.Iterar(func(clave string, dato int) bool { // O(n.logn)
		dicOrd.Guardar(clave, dato)
		return true
	})
	if dicOrd.Cantidad() != 0 {
		fmt.Println("Visitados:")
		for iter := dicOrd.IteradorRango(&inicio, &fin); iter.HaySiguiente(); iter.Siguiente() { // O(n)
			clave, _ := iter.VerActual()
			fmt.Println(_TABULACION, clave)
		}
	}
	fmt.Println(_ACEPTADO)
	//Complejidad = O(n.log.n + n) = O(n.log.k) donde n es la cantidad de elementos del abb , k el rango (?
}

func (info *informacionArchivos) procesarArchivo(ruta string) {
	datosAux := TDADiccionario.CrearHash[string, string]()
	archivo, _ := os.Open(ruta)
	defer archivo.Close()
	scanner := bufio.NewScanner(archivo)
	for scanner.Scan() {
		arrTexto := strings.Split(scanner.Text(), _TABULACION)
		info.almacenarUsuarios(datosAux, arrTexto)
		info.contarSitiosVistados(arrTexto)
	}
	//Complejidad = O(n) , donde n es la cantidad de lineas que tiene el archivo
}

func (info *informacionArchivos) almacenarUsuarios(dic TDADiccionario.Diccionario[string, string], elementos []string) {
	if dic.Pertenece(elementos[_VER_IP]) {
		detectarDoS(elementos, dic, info.visitantes)
	} else {
		dic.Guardar(elementos[_VER_IP], elementos[_VER_ZONA_HORARIA])
		info.visitantes.Guardar(elementos[_VER_IP], _VALOR_INICIAL)
	}
	// Complejidad = O(1) , porque es un hash , acceder a sus campos es O(1)
}
func (info *informacionArchivos) contarSitiosVistados(elementos []string) {
	if info.visitados.Pertenece(elementos[_VER_URL]) {
		cantidad := info.visitados.Obtener(elementos[_VER_URL])
		cantidad++
		info.visitados.Guardar(elementos[_VER_URL], cantidad)
	} else {
		info.visitados.Guardar(elementos[_VER_URL], _VALOR_INICIAL)
	}
}

func detectarDoS(elementos []string, dic TDADiccionario.Diccionario[string, string], visitantes TDADiccionario.Diccionario[string, int]) {
	t1, _ := time.Parse(_LAYOUT_PARSE, dic.Obtener(elementos[_VER_IP]))
	t2, _ := time.Parse(_LAYOUT_PARSE, elementos[_VER_ZONA_HORARIA])
	if t2.Sub(t1) < _SEGUNDO_MAXIMO {
		cantidad := visitantes.Obtener(elementos[_VER_IP])
		cantidad++
		visitantes.Guardar(elementos[_VER_IP], cantidad)
	}
	dic.Guardar(elementos[_VER_IP], elementos[_VER_ZONA_HORARIA])
}
func compararNumeros(a, b sitiosVistados) int {
	return cmp.Compare(a.cantidad, b.cantidad)
}

func agregarIPDoS(elementos [][]string, dic TDADiccionario.Diccionario[string, int]) {
	dic.Iterar(func(clave string, dato int) bool {
		if dato >= _CANT_MINIMA_DOS {
			elementos = append(elementos, strings.Split(clave, _SEPARADOR_DIGITOS_IP))
		}
		return true
	})
}

func CountingSort(elementos [][]string, digito int) {

}

//func (info *informacionArchivos) almacenarVisitantes(Ip string) {
//	if !info.visitantes.Pertenece(Ip) {
//		info.visitantes.Guardar(Ip, _ACEPTADO)
//	}
//}

//func (info *informacionUsuario) guardarInformacion(ip string, datos []string) {
//
//}

//type informacionSesionUsuario struct {
//	tiempo TDADiccionario.Diccionario[string, *TDAPila.Pila[[]string]]
//	urls   TDADiccionario.Diccionario[string, int]
//}
//
//type informacionGeneral struct {
//	informacionIPs  TDADiccionario.DiccionarioOrdenado[string, *informacionSesionUsuario]
//	informacionUrls TDADiccionario.Diccionario[string, int]
//}

//func crearNuevoUsuario() *informacionSesionUsuario {
//	return &informacionSesionUsuario{
//		TDADiccionario.CrearHash[string, *TDAPila.Pila[[]string]](),
//		TDADiccionario.CrearHash[string, int](),
//	}
//}

//func CrearInformacionArchivos() EjecucionArchivos {
//	return &informacionGeneral{
//		TDADiccionario.CrearABB[string, *informacionSesionUsuario](compararIPs),
//		TDADiccionario.CrearHash[string, int](),
//	}
//}

// AgregarArchivo implements EjecucionArchivos.
//func (info informacionGeneral) AgregarArchivo(ruta string) {
//	archivo, err := os.Open(ruta)
//	if err != nil {
//		fmt.Fprintf(os.Stderr, _ERROR+_PARAMETRO_ENTRADA_AGREGAR) // importar del main
//	}
//	defer archivo.Close()
//
//	scanner := bufio.NewScanner(archivo)
//	for scanner.Scan() {
//		lineaTexto := scanner.Text()
//		lineaInformacion := strings.Split(lineaTexto, _ESPACIO_VACIO)
//		guardarInformacion(lineaInformacion[_VER_IP], lineaInformacion[_VER_URL], lineaInformacion[_VER_ZONA_HORARIA], info)
//	}
//}

//func guardarInformacion(IP, url, tiempo string, info informacionGeneral) {
//	reloj := tranformarTiempo(tiempo)
//	if !info.informacionIPs.Pertenece(IP) {
//		usuario := crearNuevoUsuario()
//		pilaTiempos := TDAPila.CrearPilaDinamica[[]string]()
//		pilaTiempos.Apilar(reloj)
//		usuario.tiempo.Guardar(url, &pilaTiempos)
//		usuario.urls.Guardar(url, _VALOR_INICIAL)
//		info.informacionIPs.Guardar(IP, usuario)
//	} else {
//		usuario := info.informacionIPs.Obtener(IP)
//		pilaTiempos := usuario.tiempo.Obtener(url)
//		(*pilaTiempos).Apilar(reloj)
//		usuario.urls.Guardar(IP, usuario.urls.Obtener(IP)+1)
//		detectarDoS(usuario, IP, url)
//	}
//	contabilizarURLs(url, info)
//
//	// reloj := tranformarTiempo(tiempo)
//	// if !info.informacionIPs.Pertenece(IP) {
//	// 	nuevoUsuario := crearNuevoUsuario()
//	// 	pilaTiempos := TDAPila.CrearPilaDinamica[[]string]()
//	// 	pilaTiempos.Apilar(reloj)
//	// 	nuevoUsuario.tiempo.Guardar(url, pilaTiempos)
//	// 	nuevoUsuario.urls.Guardar(url, 0)
//	// 	info.informacionIPs.Guardar(IP, nuevoUsuario)
//	// } else {
//	// 	usuarioExistente := info.informacionIPs.Obtener(IP)
//	// 	pilaTiempos := usuarioExistente.tiempo.Obtener(url)
//	// 	pilaTiempos.Apilar(reloj)
//	// 	usuarioExistente.tiempo.Guardar(url, pilaTiempos)
//	// 	usuarioExistente.urls.Guardar(url, usuarioExistente.urls.Obtener(url)+1)
//	// 	info.informacionIPs.Guardar(IP, usuarioExistente)
//	// }
//	// contabilizarURLs(url, info)
//}

//func detectarDoS(usuario *informacionSesionUsuario, IP, url string) {
//	if usuario.urls.Obtener(url) >= _CANT_MINIMA_DOS {
//		segundos := calcularTiempo(usuario, url)
//		if segundos <= _SEGUNDO_MAXIMO {
//			fmt.Println("DoS: " + IP)
//		}
//	}
//}

//func calcularTiempo(usuario *informacionSesionUsuario, url string) float64 {
//	// desapilamos y apilamos
//	return _DATOS_INSUFICIENTES
//}

//func tranformarTiempo(tiempo string) []string {
//	//Transformar lo que necesitamos[dia, hora, minutos, segundos] ?
//	return []string{tiempo}
//}

//func (info informacionGeneral) VerVisitantes(desdeIP, hastaIP string) {
//	iterador := info.informacionIPs.IteradorRango(&desdeIP, &hastaIP)
//	for iterador.HaySiguiente() {
//		ip, _ := iterador.VerActual()
//		fmt.Printf("\t" + ip)
//		iterador.Siguiente()
//	}
//	fmt.Printf("OK")
//}

//func (info informacionGeneral) VerMasVisitados(mostrarCantidadUrls string) {
//	cantidadTotal, err := strconv.Atoi(mostrarCantidadUrls)
//	if err != nil {
//		fmt.Fprintf(os.Stderr, _ERROR+"ver_mas_visitados") // importar del main
//	}
// Tengo que ver si hago la ejecucion acá o cuando lea el archivo
// Como ya tengo el diccionario en contabilizarURLs(), tengo que aquí ordenarnos con un counting sort
// Lo que ordenaré, será la cantidad de visitados que tenga cada URL
// Como se ordenan de menor a mayor, pongo en una pila cada elemento
// Y hare un ciclo for hasta que la Pila esté vacia, o cuanto el contador haya llegado a 10
// Desapilando lo ltimo que apilé
//pila := TDAPila.CrearPilaDinamica[string]() //String no va a ser, tiene que ser una estructura
// EjecuionDeOrdenamiento
//	ordenamientoUrlVisitados(pila, info)
//	contador := 0
//	for !pila.EstaVacia() && contador <= cantidadTotal {
//		fmt.Printf("\t" + pila.Desapilar())
//		contador++
//	}
//	fmt.Printf("OK")
//}

//func ordenamientoUrlVisitados(pila TDAPila.Pila[string], info informacionGeneral) {
//	// Hacer cositas
//}

//func contabilizarURLs(urlVisitado string, info informacionGeneral) {
//	if !info.informacionUrls.Pertenece(urlVisitado) {
//		info.informacionUrls.Guardar(urlVisitado, _VALOR_INICIAL)
//	} else {
//		info.informacionUrls.Guardar(urlVisitado, info.informacionUrls.Obtener(urlVisitado)+1)
//	}
//}
