# Panel de Debug

## Estado
> Implementado. Funcional.

___

## ¿Qué es?

El panel de debug es una vista interna de Spoti NeXt que muestra los logs del proceso de descarga en tiempo real. Está orientado a diagnosticar errores, verificar el comportamiento de los resolvers y confirmar el resultado de operaciones de backend que de otro modo serían invisibles para el usuario.

---

## ¿Dónde está?

**Sidebar → Debug** (ícono de terminal)

---

## Qué muestra

El panel agrega dos flujos de logs:

### 1 — Logs del frontend (logger interno)

Generados por el hook `useDownload` durante el ciclo de descarga. Incluyen:

- Fuente intentada y nombre del track
- Resultado de cada intento (éxito / fallo + mensaje de error)
- Modo activo (directo o con fallback)
- Progreso de la cola

Los niveles de log tienen colores:

| Nivel | Color |
|-------|-------|
| `info` | Azul |
| `success` | Verde |
| `warning` | Amarillo |
| `error` | Rojo |
| `debug` | Gris |

### 2 — Logs del backend (DebugLog / eventos)

Emitidos por el backend Go y reenviados al frontend vía eventos Wails (`debug:atmos`). Actualmente incluyen los logs del flujo Atmos de Tidal:

- `[atmos-manifest]` — detalles del request al endpoint de Tidal: trackID, quality, URL, status HTTP
- `[atmos-tv]` — intento de re-autenticación con el TV client: clientID, token length, resultado
- `[atmos-manifest]` — audioMode, audioQuality, bitDepth, sampleRate del resultado
- `[atmos-manifest]` — body completo en caso de error HTTP

---

## Relación con la cuenta de Tidal

La cuenta de Tidal vinculada afecta directamente lo que aparece en el panel cuando Atmos está activo.

Lo que el panel permite confirmar:

| Log | Qué indica |
|-----|-----------|
| `[atmos-tv] success, tokenLen=NNN` | Re-auth con TV client funcionó |
| `[atmos-manifest] audioMode=DOLBY_ATMOS` | La cuenta tiene acceso Atmos — descarga Atmos real |
| `[atmos-manifest] audioMode=STEREO` | La cuenta no tiene Atmos en la suscripción — fallback |
| `[atmos-manifest] audioMode=HI_RES_LOSSLESS` | Fallback a lossless estándar |
| `[atmos-tv] error: ...` | Fallo en la re-auth TV — el access token de sesión fue usado como fallback |

Ver `tidal_atmos.md` para el diagnóstico completo del flujo Atmos.

---

## Acciones disponibles

| Acción | Descripción |
|--------|-------------|
| **Clear** | Limpia todos los logs del panel |
| **Copy** | Copia todos los logs al portapapeles en formato texto plano |
| **Export Failed** | Exporta a archivo los tracks que fallaron en la sesión actual (habilitado solo si hay descargas fallidas) |

---

## Notas técnicas

- Los logs no persisten entre sesiones. Al cerrar la aplicación se pierden.
- El panel se desplaza automáticamente al log más reciente.
- `log.Printf` del backend Go escribe a stderr — ese output **no** aparece en el panel. Solo aparece lo emitido explícitamente con `DebugLog(...)`.

---

## Archivos relevantes

| Archivo | Función |
|---------|---------|
| `backend/debug_emitter.go` | `DebugLog()` global — reenvía mensajes al frontend vía Wails events |
| `app.go` | Registra el emitter en `startup()` con `SetDebugEmitter` |
| `backend/tidal_account.go` | Fuente de los logs `[atmos-*]` |
| `frontend/src/lib/logger.ts` | Logger del frontend (niveles info/success/warning/error/debug) |
| `frontend/src/components/DebugLoggerPage.tsx` | UI del panel, suscripción al logger, botones de acción |
| `frontend/src/hooks/useDownloadQueueData.ts` | Estado de la cola (habilita botón Export Failed) |
