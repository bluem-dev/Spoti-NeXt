# Metadata — Documentación técnica

## Estado

Implementado. Funcional.

---

## Fuentes de metadata

Spoti NeXt obtiene metadata de múltiples fuentes y las combina en el archivo descargado:

| Fuente | Qué aporta |
|--------|-----------|
| **Spotify** | Título, artista, álbum, fecha de lanzamiento, track number, disc number, ISRC, portada, copyright, publisher |
| **MusicBrainz** | Géneros (tags de la comunidad) |
| **Discogs** | Label, año, géneros y estilos adicionales |
| **iTunes / Apple Music** | Metadata complementaria (título, artista, fecha) — sin streaming |
| **VGMdb** | Metadata para soundtracks de videojuegos y anime |
| **Lyrics providers** | Letras sincronizadas o en texto plano |

Las fuentes de enriquecimiento (Discogs, iTunes, VGMdb) están implementadas pero su integración al flujo de descarga está pendiente — ver BACKLOG (`wire-metadata-enrichment`).

---

## Portada

**Settings → Advanced → Embed Max Quality Cover**

| Opción | Descripción |
|--------|-------------|
| Desactivado | Portada estándar de Spotify (640×640px) |
| Activado | Portada de máxima resolución disponible (hasta 3000×3000px) |

La portada se embebe directamente en el archivo FLAC/M4A como imagen.

---

## Letras (Lyrics)

**Settings → Advanced → Embed Lyrics**

Cuando está activo, Spoti NeXt intenta obtener las letras del track y embeberlas en los metadatos del archivo.

Las letras se obtienen de providers externos. Si no hay letras disponibles, el campo queda vacío — no falla la descarga.

---

## Géneros

**Settings → Advanced → Embed Genre**

Cuando está activo, se embebe el género en los metadatos.

**Settings → Advanced → Use Single Genre**

Cuando está activo, solo se embebe el primer género disponible. Cuando está desactivado, se embeben todos los géneros separados por el separador configurado.

---

## Separador de campos múltiples

**Settings → Advanced → Separator**

Aplica a artistas múltiples, géneros y otros campos de lista.

| Opción | Carácter |
|--------|---------|
| Comma | `, ` |
| Semicolon | `; ` |

---

## Artista único

**Settings → Advanced → Use First Artist Only**

Cuando está activo, si un track tiene múltiples artistas (`Artist A, Artist B`), solo se embebe el primero (`Artist A`).

---

## Track Number

**Settings → Files → Track Number**

Cuando está activo, el número de track se embebe en los metadatos del archivo. Independientemente de este ajuste, la plantilla de nombre de archivo puede incluir `{track}` para usar el número en el nombre.

---

## Archivos relevantes

| Archivo | Función |
|---------|---------|
| `backend/metadata.go` | Escritura de metadata en archivos (ffmpeg) |
| `backend/spotify_metadata.go` | Obtención de metadata desde Spotify |
| `backend/musicbrainz.go` | Géneros desde MusicBrainz |
| `backend/discogs.go` | Metadata desde Discogs |
| `backend/itunes.go` | Metadata desde iTunes Search API |
| `backend/vgmdb.go` | Metadata desde VGMdb |
| `backend/lyrics.go` | Obtención de letras |
| `backend/cover.go` | Descarga y procesamiento de portada |
| `app.go` | Ensamblado de metadata completa antes de escritura |
