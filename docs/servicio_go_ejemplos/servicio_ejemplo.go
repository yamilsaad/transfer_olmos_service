// Arquitectura Interna del Servicio Go
package main

import (
	"database/sql"
	"log"
	"os/exec"
	"time"
)

func main() {
	// 1. Inicializar conexión a BD Don Web (reutilizable)
	db, err := sql.Open("mysql",
		"user:pass@tcp(don-web-server:3306)/branet_gesdrims")
	if err != nil {
		log.Fatal("Error conectando a BD", err)
	}
	defer db.Close()

	// 2. Crear TICKER para polling cada 20 segundos
	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()

	log.Println("✅ Servicio iniciado - Polling cada 20 segundos")

	// 3. LOOP INFINITO (Se mantiene activo)
	for {
		select {
		case <-ticker.C:
			// Cada 20 segundos ejecuta esta función
			procesarFacturasPendientes(db)
		}
	}
}

func procesarFacturasPendientes(db *sql.DB) {
	log.Println("[CICLO] Buscando facturas nuevas...")

	// FASE 1: SELECCIÓN
	rows, err := db.Query(`
        SELECT ID, TIPCOM, PTV, NROCOM, FECCOM, HORCOM 
        FROM maecompr 
        WHERE CGOSUC = '02' 
          AND ANULAD = 0 
          AND ID NOT IN (SELECT idmaecompr FROM control_patio_olmos) 
        ORDER BY ID ASC 
        LIMIT 10
    `)
	if err != nil {
		log.Println("❌ Error en query:", err)
		return // Reintentar en 20 segundos
	}
	defer rows.Close()

	// FASE 2: CONSTRUCCIÓN
	for rows.Next() {
		var id, tipcom, ptv, nrocom, feccom, horcom string
		rows.Scan(&id, &tipcom, &ptv, &nrocom, &feccom, &horcom)

		// Procesar factura
		recordC := construirRecordC(id, tipcom, ptv, nrocom, feccom, horcom)
		recordD := construirRecordD(db, id)
		recordP := construirRecordP(db, id)

		// FASE 3: DISPARO OPERATIVO (Ejecutar CLI)
		exitCodeC := ejecutarTransfer(recordC)
		exitCodeD := ejecutarTransfer(recordD)
		exitCodeP := ejecutarTransfer(recordP)

		// FASE 4: CONFIRMACIÓN
		if exitCodeC == 0 && exitCodeD == 0 && exitCodeP == 0 {
			registrarEnControl(db, id, "OK")
			log.Printf("✅ Factura %s enviada correctamente\n", id)
		} else {
			registrarEnControl(db, id, "ERROR")
			log.Printf("❌ Factura %s falló (códigos: %d,%d,%d)\n",
				id, exitCodeC, exitCodeD, exitCodeP)
		}
	}
}

func ejecutarTransfer(record string) int {
	cmd := exec.Command("transfer.exe", record)
	err := cmd.Run()

	if exitErr, ok := err.(*exec.ExitError); ok {
		return exitErr.ExitCode()
	}
	return 0
}

/************************************************************************/
/************************************************************************/
/************************************************************************/
//Manejo de Errores & Resiliencia
func procesarFacturasPendientes(db *sql.DB) {
	// Si la BD se desconecta
	err := db.Ping()
	if err != nil {
		log.Println("⚠️  BD no responde, reintentando en próximo ciclo...")
		return // No falla, espera 20 seg y reintenta
	}

	// Si transfer.exe no existe
	cmd := exec.Command("transfer.exe", record)
	err = cmd.Run()
	if err != nil {
		log.Println("⚠️  transfer.exe error:", err)
		// Se registra como fallido, reintenta en próximo ciclo
		return
	}

	// Reconexión automática
	if rows.Err() != nil {
		log.Println("⚠️  Error leyendo filas, reconectando...")
		db.Close()
		db, _ = sql.Open("mysql", connString)
	}
}
