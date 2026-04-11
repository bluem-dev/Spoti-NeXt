# Calidad y Fuentes de Descarga — Documentación técnica

## Estado

Implementado. Funcional.

---

## Fuentes disponibles

| Fuente | Calidad máxima | Formato | Notas |
|--------|---------------|---------|-------|
| **Tidal** | HI_RES_LOSSLESS (24-bit/192kHz) | FLAC | Requiere cuenta o proxy afkarxyz |
| **Qobuz** | Studio Master (24-bit/192kHz) | FLAC | Requiere proxy afkarxyz |
| **Amazon** | Original (lossless) | FLAC | Requiere proxy afkarxyz |
| **Deezer** | — | — | Solo resolver de links, no descarga real |

---

## Configuración de calidad por fuente

**Settings → Quality**

### Tidal

| Opción | Descripción |
|--------|-------------|
| `LOSSLESS` | 16-bit/44.1kHz (CD quality) |
| `HI_RES_LOSSLESS` | Hasta 24-bit/192kHz (MQA decodificado o Hi-Res nativo) |

Con Atmos activado, se intenta primero `DOLBY_ATMOS`. Si no está disponible, cae automáticamente a `HI_RES_LOSSLESS` → `LOSSLESS`.

### Qobuz

| Opción | Valor interno | Descripción |
|--------|--------------|-------------|
| MP3 320 | `5` | Lossy |
| FLAC CD | `6` | 16-bit/44.1kHz |
| FLAC Hi-Res | `7` | 24-bit hasta 96kHz |
| FLAC Studio Master | `27` | 24-bit hasta 192kHz |

### Amazon

Calidad fija: `original` — se descarga el stream de mayor calidad disponible para el track.

---

## Modo Auto

**Settings → General → Downloader → Auto**

Spoti NeXt intenta las fuentes en el orden definido por **Auto Order**. El orden por defecto es `Tidal → Qobuz → Amazon`.

El orden se puede cambiar desde Settings. No todas las permutaciones incluyen Deezer — ver BACKLOG (`autoOrder-deezer`).

### Auto Quality

Cuando el downloader es `Auto`, `Auto Quality` define el umbral mínimo de bit depth aceptable:

| Valor | Descripción |
|-------|-------------|
| `16` | Acepta 16-bit y 24-bit |
| `24` | Intenta fuentes hasta conseguir 24-bit; si ninguna lo tiene, usa lo disponible |

---

## Quality Fallback

**Settings → Advanced → Allow Quality Fallback**

Cuando está activo (por defecto), si la calidad solicitada no está disponible, la fuente cae a la calidad inferior más próxima. Ejemplo: si se solicita `HI_RES_LOSSLESS` y el track solo tiene `LOSSLESS`, descarga en `LOSSLESS`.

Cuando está desactivado, si la calidad exacta no está disponible, la descarga falla.

**Relación con Modo Directo:** si el modo directo está activo, `allow_fallback` se fuerza a `false` independientemente de este ajuste — no hay fallback de ningún tipo.

---

## Archivos relevantes

| Archivo | Función |
|---------|---------|
| `backend/tidal.go` | Descarga Tidal, lógica de calidad y fallback Atmos |
| `backend/qobuz.go` | Descarga Qobuz, map de calidad numérica |
| `backend/amazon.go` | Descarga Amazon |
| `frontend/src/hooks/useDownload.ts` | Lógica de Auto Order, Auto Quality, fallback entre fuentes |
| `frontend/src/lib/settings.ts` | Tipos: `tidalQuality`, `qobuzQuality`, `autoOrder`, `autoQuality`, `allowFallback` |
| `frontend/src/components/SettingsPage.tsx` | UI de todos los controles de calidad |
