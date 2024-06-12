package interfazArchivos

type EjecucionArchivos interface {
	// Agrega la informacion del archivo
	AgregarArchivo(string)

	// Lista todos los IPs en el rango determinado
	VerVisitantes(string, string)

	//// Se lista los N recursos más visitador
	VerMasVisitados(string)
}
