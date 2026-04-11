# Dolby Atmos

## Estado
> Implementado. Funcional bajo las condiciones descritas abajo.

#

## ¿Cómo funciona?

Cuando el usuario activa el toggle Atmos y tiene una cuenta Tidal vinculada, Spoti NeXt:

1. Hace login con el client ID de Device Authorization Flow (`fX2JxdmntZWK0ixT`) — obtiene `access_token` y `refresh_token`.
2. Para cada descarga con Atmos activado, re-autentica usando el `refresh_token` con el TV client (`cgiF7TQuB97BUIu3`) — obtiene un nuevo `access_token` con scope `TIDAL_Android_TV_Atmos_HiRes`.
3. Hace el request al endpoint `playbackinfopostpaywall/v4` con el TV token como `Bearer` y los headers del cliente Android TV.
4. Si la API devuelve `audioMode=DOLBY_ATMOS`, descarga el stream E-AC3 y lo guarda como `.m4a`.
5. Si la API devuelve `audioMode=STEREO`, cae automáticamente a HI_RES_LOSSLESS.

El mecanismo de re-autenticación funciona correctamente — el token TV con scope Atmos se obtiene en cada descarga. El resultado depende de lo que el servidor de Tidal autorice para la cuenta.

---

## Limitaciones

### Limitación 1 — Suscripción

Tidal evalúa server-side si la cuenta tiene Atmos habilitado antes de devolver el stream. El token TV correcto es condición necesaria pero no suficiente. Si la cuenta no tiene **Tidal HiFi Plus** con Atmos activo, la API devuelve `audioMode=STEREO` independientemente del client ID, endpoint o headers usados.

**Verificación:** Si en la app oficial de Tidal no aparece el badge "Dolby Atmos" en ningún track, la cuenta no tiene acceso y Spoti NeXt tampoco podrá obtenerlo.

### Limitación 2 — Catálogo regional

Algunos streams Atmos están geo-restringidos. Un track con badge Atmos en una región puede no tenerlo en otra. Spoti NeXt usa el `countryCode` de la sesión de la cuenta — no hay forma de forzar otra región.

### Limitación 3 — Disponibilidad por track

No todos los tracks tienen versión Atmos. Atmos es un mix separado — no una conversión del estéreo. Si el track no tiene Atmos en el catálogo de Tidal, no hay nada que descargar.

---

## Investigación realizada

Durante el desarrollo se investigaron y descartaron las siguientes variables:

| Variable | Resultado |
|----------|-----------|
| `countryCode=US` hardcodeado | Sin efecto — no es restricción regional |
| `X-Tidal-Token` con client TV como header | Insuficiente — el token Bearer debe ser emitido por ese client |
| Client ID `4N3n6Q1x95LL5K7p` (documentado en fuentes antiguas) | Incorrecto — reemplazado por `cgiF7TQuB97BUIu3` (orpheusdl actual) |
| Client ID `LXujKdmnc6QtydvY` (RedSea FireTV) | No probado directamente — mismo mecanismo |
| Headers `User-Agent`, `Connection` del cliente Android TV | Correctos, no cambian el resultado |
| `Accept-Encoding: gzip` manual | Rompe el parsing — Go lo maneja automáticamente |
| Endpoints alternativos (v4/LOSSLESS, v1/HI_RES, sin audioquality) | Todos devuelven `audioMode=STEREO` |
| Scope del token TV | Confirmado como `TIDAL_Android_TV_Atmos_HiRes` via `/sessions` |

La barrera es la autorización server-side por suscripción. No existe workaround desde el cliente.

---

## Comportamiento esperado

| Condición | Resultado |
|-----------|-----------|
| Cuenta vinculada + HiFi Plus con Atmos + track con Atmos | Descarga E-AC3 `.m4a` |
| Cuenta vinculada + sin Atmos en suscripción | Fallback a HI_RES_LOSSLESS automático |
| Cuenta vinculada + track sin Atmos en catálogo | Fallback a HI_RES_LOSSLESS automático |
| Sin cuenta vinculada | Proxies de terceros (sin acceso a Atmos) |

---

## Archivos relevantes

| Archivo | Función |
|---------|---------|
| `backend/tidal_account.go` | Login, TV re-auth, manifest request |
| `backend/tidal.go` | Integración con el flujo de descarga, fallback |
| `backend/debug_emitter.go` | Logging al panel de debug del frontend |
| `frontend/src/components/SettingsPage.tsx` | Toggle Atmos en UI |
