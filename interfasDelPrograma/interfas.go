package interfazArchivos

type ejecucionArchivos interface {
	AgregarArchivo(string)

	VerVisitantes(string, string)

	VerMasVisitados(int)
}
