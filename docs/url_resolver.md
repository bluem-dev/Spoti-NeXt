# Resolver de enlaces

## Estado
> Implementado. Funcional. Algunos gaps pendientes (ver BACKLOG).
___

## ¿Qué es?

El resolver es el sistema que, dado un link de Spotify, encuentra las URLs equivalentes en Tidal, Qobuz, Amazon Music y Deezer para poder descargar el audio.

El proceso comienza siempre en Spotify: se extrae el ISRC del track, y ese ISRC se usa para buscar el mismo track en las plataformas de descarga.

---

## Flujo de resolución

```
Link de Spotify
  → extrae Spotify Track ID
  → lookupSpotifyISRC()         ← Spotify web player (TOTP) + cache 90d
  → consulta resolver cache     ← BoltDB, key ISRC, TTL 30d
      → hit: devuelve URLs guardadas
      → miss: continúa

  → resolver(es) configurado(s):
      → Songstats               ← scraping HTML + JSON-LD
      → Deezer + song.link      ← api.deezer.com/isrc → api.song.link
      → Tidal ISRC nativo       ← listen.tidal.com/v1/tracks?isrc= (si cuenta vinculada)

  → guarda en resolver cache
  → devuelve resolvedTrackLinks { TidalURL, AmazonURL, DeezerURL, QobuzURL, ISRC }
```

---

## Configuración disponible

**Settings → Advanced → Link Resolver**

| Opción | Valor | Descripción |
|--------|-------|-------------|
| **Resolver** | `Songlink` (por defecto) | Usa Deezer API + song.link (Odesli) |
| **Resolver** | `Songstats` | Usa scraping de Songstats (más frágil, sin límites de rate) |
| **Allow Resolver Fallback** | toggle | Si el resolver primario falla, intenta el alternativo |

---

## Cache de resolución

Las URLs resueltas se guardan en BoltDB por ISRC con TTL de 30 días. Esto significa que la segunda vez que se descarga un track de la misma playlist o álbum, el resolver no hace ningún request externo.

El cache de ISRC (Spotify → ISRC) tiene TTL de 90 días y es independiente.

---

## Proveedores que se resuelven

| Proveedor | Fuente de la URL | Notas |
|-----------|-----------------|-------|
| **Tidal** | song.link / Songstats / Tidal ISRC nativo | Principal fuente de descarga |
| **Amazon** | song.link / Songstats | Región puede variar |
| **Deezer** | api.deezer.com/isrc | Usado como resolver, no como fuente de descarga |
| **Qobuz** | song.link (`"qobuz"` key) | Disponible si song.link lo devuelve |

---

## Prioridad dinámica de proveedores

Además del resolver de links, existe un sistema de prioridad dinámica entre proveedores de descarga. Tras cada descarga, el resultado (éxito o fallo) se registra en BoltDB. En descargas sucesivas, los proveedores con mejor historial se intentan primero.

Este comportamiento es automático y no tiene configuración en UI — funciona en segundo plano.

---

## Limitaciones

- **Songstats**: el scraping puede romperse si el sitio cambia su estructura HTML.
- **song.link (Odesli)**: tiene rate limiting. En playlists largas puede generar errores transitorios.
- **Qobuz**: la URL de Qobuz que devuelve song.link no siempre está disponible — depende del catálogo regional.
- **Tidal ISRC nativo**: implementado como resolver adicional, requiere cuenta Tidal vinculada. Pendiente de validación en producción.

---

## Archivos relevantes

| Archivo | Función |
|---------|---------|
| `backend/link_resolver.go` | Orquestador principal del flujo de resolución |
| `backend/songlink.go` | Cliente song.link + merge de respuesta |
| `backend/songstats.go` | Scraper de Songstats |
| `backend/tidal_isrc_resolver.go` | Resolver Tidal nativo por ISRC |
| `backend/resolver_cache.go` | Cache BoltDB de URLs resueltas (TTL 30d) |
| `backend/isrc_cache.go` | Cache BoltDB de ISRCs (TTL 90d) |
| `backend/isrc_finder.go` | Obtención del ISRC desde Spotify web player |
| `backend/provider_priority.go` | Sistema de prioridad dinámica entre proveedores |
| `backend/config.go` | `GetLinkResolverSetting()`, `GetLinkResolverAllowFallback()` |
