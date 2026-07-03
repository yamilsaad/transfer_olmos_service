
╔════════════════════════════════════════════════════════════════════╗
║          SERVICIO GO EN PC DE CAJA #1 - LÍNEA DE TIEMPO            ║
╚════════════════════════════════════════════════════════════════════╝

T=0 seg:   ✅ Servicio inicia (se instala como Windows Service)
           └─ Abre conexión a BD Don Web
           └─ Inicia TICKER de 20 segundos

T=20 seg:  🔄 CICLO 1
           ├─ Lee 10 facturas nuevas de BD
           ├─ Ejecuta transfer.exe (3 comandos por factura)
           ├─ Registra resultado en control_patio_olmos
           └─ Espera 20 segundos

T=40 seg:  🔄 CICLO 2
           ├─ Lee 10 facturas nuevas (las anteriores están en control)
           ├─ Ejecuta transfer.exe
           ├─ Registra resultado
           └─ Espera 20 segundos

T=60 seg:  🔄 CICLO 3
           └─ ... (continúa indefinidamente)

...

T=∞:       ✅ Servicio sigue corriendo hasta:
           • Reinicio manual: net stop OlmosTransfer
           • Reinicio de Windows (se reinicia automático)
           • Error grave no recuperable (logging + reinicio automático)