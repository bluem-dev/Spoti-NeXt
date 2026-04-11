# Spoti NeXt

"Spoti NeXt" está basada en la versión original creada por afkarxyz.

![Image](https://res.cloudinary.com/djljnirz2/image/upload/v1775858383/1111_yl9tj3.jpg)

***

## CARACTERÍSTICAS

| Función | Estado |
|------|--------|
| Fuente única (DD) | 🟩 Implementado |
| Resolución de enlaces | 🟩 Corregido y extendido |
| Metadata: fuente (DD) | 🟩 Implementado |
| Descarga: Atmos (TIDAL) | 🟩 Implementado |
| Idiomas: inglés, español | 🟩 Implementado |
| Desc. simultáneas | 🟥 Pendiente |


# INFORMACIÓN DE CORRECCIONES

- Resolución de enlaces:<br>
→ Solución a errores de obtención de datos y cruce fallidos.<br>
→ Extensión de URL y regiones: los enlaces se resuelven sin importar región.<br>
→ Relacionados a "SpotiFetch API" & "Fetch Data".<br>
- ISRC/TTL: corrección de caché, tiempo y duplicado.
- Loggin plano: ahora usa log.Print de stdlib.
- ...

# INFORMACIÓN DE CARACTERISTICAS

- Dolby Atmos:<br>
→ Requiere cuenta TIDAL vinculada.<br>
Cuando está disponible, accede al stream E-AC3 en formato M4A.<br>
Si el track/cuenta no tienen acceso a D. Atmos, se descarga en HI_RES_LOSSLESS como alternativa.<br>
La disponibilidad depende del plan de suscripción y del catálogo regional.<br>
Spoti NeXt no puede garantizar Atmos en todos los casos.<br>

# FUTURAS IMPLEMENTACIONES

- Selector de formatos: FLAC, M4A, OPUS, MP3.<br>
→ Modalidades: manual y automática (definida por calidad adecuada).<br>
- División de álbumes: DISC 1 X DISC 2 x DISC 3, etc. Órden original.
- ...

# FUTURAS CORRECCIONES

- Conservar portadas y metadata en todos los formatos con compatibilidad: ID3Tag.<br>
- Expansión de directorios: distinción EPs, Remixes, Compilations, Playlists, Deluxe, Editions...<br>
- ...

## DECLARACIÓN

El proyecto está enfocado con própositos educacionales, y de uso privado. Sin ánimos de lucro. Su autor/creador no fomenta la infracción sobre derechos de autor.

**Spoti NeXt** (SpotiFLAC) es una herramienta de terceros y no está afiliada, respaldada ni vinculada a Spotify, Tidal, Qobuz, Amazon Music ni a ningún otro 
servicio de streaming.

***

## CRÉDITOS DE APIs & APLICACIÓN

[MusicBrainz](https://musicbrainz.org) · [LRCLIB](https://lrclib.net) · [Songlink/Odesli](https://song.link) · [hifi-api](https://github.com/binimum/hifi-api) · [dabmusic.xyz](https://dabmusic.xyz) · [SpotiFLAC](https://github.com/afkarxyz/SpotiFLAC)

