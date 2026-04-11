# Cuenta Tidal

## Estado
> Implementado. Funcional bajo las condiciones descritas abajo.

___

## ¿Qué es?

La vinculación de una cuenta Tidal personal permite a Spoti NeXt obtener streams directamente con las credenciales del usuario, en lugar de depender de proxies de terceros. Es requisito para acceder a Dolby Atmos.

---

## ¿Cómo vincular?

**Settings → Advanced → Tidal Account**

El flujo es Device Authorization Flow:

1. Se inicia el proceso desde Settings — Spoti NeXt muestra un código de dispositivo.
2. El usuario abre `link.tidal.com` en un navegador y autoriza el dispositivo con su cuenta Tidal.
3. Spoti NeXt completa el login y guarda `access_token`, `refresh_token`, `countryCode` y `userID` en los settings locales.

Las credenciales se guardan localmente en el archivo de configuración de la aplicación. No se envían a ningún servidor externo.

---

## Qué habilita

| Función | Sin cuenta | Con cuenta |
|---------|-----------|-----------|
| Descarga desde Tidal (vía proxy afkarxyz) | ✅ | ✅ |
| Calidad HI_RES_LOSSLESS | ✅ (proxy) | ✅ (directo) |
| Dolby Atmos | ❌ | ✅ (requiere HiFi Plus con Atmos) |
| Resolver Tidal nativo por ISRC | ❌ | ✅ |

---

## Estructura de la sesión

La aplicación mantiene dos tipos de sesión con Tidal:

| Sesión | Client ID | Uso |
|--------|-----------|-----|
| **Principal** | `fX2JxdmntZWK0ixT` | Login, refresh, descarga estándar |
| **TV / Atmos** | `cgiF7TQuB97BUIu3` | Re-auth efímera para obtener token con scope Atmos |

La sesión TV es efímera — se obtiene por descarga cuando Atmos está activo, usando el `refresh_token` de la sesión principal. No reemplaza la sesión principal.

---

## Detección de sesión incompatible

Al iniciar la aplicación, si se detecta una sesión guardada con un `clientID` incompatible con el esquema actual de autenticación, Spoti NeXt la invalida automáticamente y solicita re-vinculación. Esto ocurrió durante el desarrollo cuando se actualizó el client ID del login.

---

## Notas de seguridad

- Los tokens se guardan en el archivo de configuración local (`config.json` en el directorio de la app).
- El `refresh_token` es de larga duración — si el archivo de configuración se comparte, la cuenta queda expuesta.
- Para desvincular: **Settings → Advanced → Tidal Account → Unlink**.

---

## Relación con Atmos

Ver `tidal_atmos.md` para el flujo completo de autenticación y descarga Atmos.

---

## Archivos relevantes

| Archivo | Función |
|---------|---------|
| `backend/tidal_account.go` | Login (Device Auth Flow), TV re-auth, manifest request |
| `backend/tidal.go` | Integración con el flujo de descarga |
| `frontend/src/components/SettingsPage.tsx` | UI de vinculación, detección de sesión incompatible |
| `frontend/src/lib/settings.ts` | Campos `tidalAccount*` en el tipo `Settings` |
