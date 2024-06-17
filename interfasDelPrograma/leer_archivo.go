package interfazArchivos

import (
	"bufio"
	"cmp"
	"fmt"
	"os"
	"strconv"
	"strings"
	TDADiccionario "tdas/diccionario"
	TDAHeap "tdas/heap"
	TDALista "tdas/lista"
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
			return 1
		}
		if intIP2 > intIP1 {
			return -1
		}
	}
	return 0
}

type informacionArchivos struct {
	AbbVisitantes TDADiccionario.DiccionarioOrdenado[string, bool]
	DiccURL       TDADiccionario.Diccionario[string, int]
	DiccIP        TDADiccionario.Diccionario[string, TDALista.Lista[string]]
}

type sitiosVistados struct {
	Url      string
	Cantidad int
}

func CrearAnalisisDeArchivos() EjecucionArchivos {
	info := new(informacionArchivos)
	info.AbbVisitantes = TDADiccionario.CrearABB[string, bool](CompararIPs)
	info.DiccURL = TDADiccionario.CrearHash[string, int]()
	info.DiccIP = TDADiccionario.CrearHash[string, TDALista.Lista[string]]()
	return info
}

func (info informacionArchivos) AgregarArchivo(ruta string) {
	DiccDos := TDADiccionario.CrearHash[string, string]()
	archivo, _ := os.Open(ruta)
	defer archivo.Close()
	scanner := bufio.NewScanner(archivo)
	for scanner.Scan() {
		lineaTexto := strings.Split(scanner.Text(), _TABULACION)
		info.almacenarUsuarios(lineaTexto, DiccDos)
		info.contarSitiosVistados(lineaTexto)
	}
	if DiccDos.Cantidad() >= 1 {
		arrDoS := ordenarDos(DiccDos)
		for _, IP := range arrDoS {
			fmt.Println("DoS: " + IP)
		}
	}
	fmt.Println(_ACEPTADO)
}

func ordenarDos(listadoDicc TDADiccionario.Diccionario[string, string]) []string {
	arrDoS := make([]string, 0)
	iter := listadoDicc.Iterador()
	for iter.HaySiguiente() {
		IPDoS, _ := iter.VerActual()
		arrDoS = append(arrDoS, IPDoS)
		iter.Siguiente()
	}
	arrDoS = countingSortIPs(arrDoS, 3)
	arrDoS = countingSortIPs(arrDoS, 2)
	arrDoS = countingSortIPs(arrDoS, 1)
	arrDoS = countingSortIPs(arrDoS, 0)
	return arrDoS
}

func convertirANumero(array string, posicion int) int {
	ip := strings.Split(array, ".")
	elemento := ip[posicion]
	ipInt, _ := strconv.Atoi(elemento)
	return ipInt
}

func countingSortIPs(arrayIPs []string, posicion int) []string {

	frecuencias := make([]int, 256)
	for _, ipString := range arrayIPs {
		numero := convertirANumero(ipString, posicion)
		frecuencias[numero]++
	}

	sumasSucesivas := make([]int, 256)
	for i := 1; i < 256; i++ {
		sumasSucesivas[i] = sumasSucesivas[i-1] + frecuencias[i-1]
	}

	resultado := make([]string, len(arrayIPs))
	for _, ipString := range arrayIPs {
		numero := convertirANumero(ipString, posicion)
		pos := sumasSucesivas[numero]
		resultado[pos] = ipString
		sumasSucesivas[numero]++
	}

	return resultado
}

func (info informacionArchivos) VerMasVisitados(dato string) {
	mininimo, _ := strconv.Atoi(dato)
	heap := TDAHeap.CrearHeap[sitiosVistados](compararCantidades)
	info.DiccURL.Iterar(func(clave string, valor int) bool {
		heap.Encolar(sitiosVistados{Url: clave, Cantidad: valor})
		return true
	})
	fmt.Println("Sitios mas visitados:")
	for i := 0; i < mininimo && !heap.EstaVacia(); i++ {
		url := heap.Desencolar()
		fmt.Println(_TABULACION, url.Url, "-", url.Cantidad)
	}
}

// Terminado
func (info informacionArchivos) VerVisitantes(inicio, fin string) {
	iterRango := info.AbbVisitantes.IteradorRango(&inicio, &fin)
	for iterRango.HaySiguiente() {
		clave, _ := iterRango.VerActual()
		fmt.Println(_TABULACION, clave)
		iterRango.Siguiente()
	}
}

func (info *informacionArchivos) almacenarUsuarios(infoIP []string, diccIPDoS TDADiccionario.Diccionario[string, string]) {
	if info.DiccIP.Pertenece(infoIP[_VER_IP]) {
		agregarTiempo(infoIP, info.DiccIP)
		detectarDoS(infoIP, info.DiccIP, diccIPDoS)
	} else {
		info.AbbVisitantes.Guardar(infoIP[_VER_IP], true)
		registroTiempo := TDALista.CrearListaEnlazada[string]()
		registroTiempo.InsertarUltimo(infoIP[_VER_ZONA_HORARIA])
		info.DiccIP.Guardar(infoIP[_VER_IP], registroTiempo)
	}
}

func agregarTiempo(IP []string, diccIP TDADiccionario.Diccionario[string, TDALista.Lista[string]]) {
	listaTiempos := diccIP.Obtener(IP[_VER_IP])
	if listaTiempos.Largo() == _CANT_MINIMA_DOS {
		listaTiempos.BorrarPrimero()
		listaTiempos.InsertarUltimo(IP[_VER_ZONA_HORARIA])
	} else {
		listaTiempos.InsertarUltimo(IP[_VER_ZONA_HORARIA])
	}
}

func detectarDoS(IP []string, tiempoIP TDADiccionario.Diccionario[string, TDALista.Lista[string]], diccIPDoS TDADiccionario.Diccionario[string, string]) {
	registroTiempo := tiempoIP.Obtener(IP[_VER_IP])
	if registroTiempo.Largo() == _CANT_MINIMA_DOS {
		primerTiempo, _ := time.Parse(_LAYOUT_PARSE, registroTiempo.VerPrimero())
		segundoTiempo, _ := time.Parse(_LAYOUT_PARSE, registroTiempo.VerUltimo())
		if segundoTiempo.Sub(primerTiempo) <= _SEGUNDO_MAXIMO*time.Second {
			diccIPDoS.Guardar(IP[_VER_IP], "Es DoS")
		}
	}
}

// Terminado
func (info *informacionArchivos) contarSitiosVistados(elementos []string) {
	if info.DiccURL.Pertenece(elementos[_VER_URL]) {
		cantidad := info.DiccURL.Obtener(elementos[_VER_URL])
		info.DiccURL.Guardar(elementos[_VER_URL], cantidad+1)
	} else {
		info.DiccURL.Guardar(elementos[_VER_URL], _VALOR_INICIAL)
	}
}

func compararCantidades(a, b sitiosVistados) int {
	return cmp.Compare(a.Cantidad, b.Cantidad)
}
