# Modo directo

## Estado
> Implementado. Funcional.
___

## ¿Qué es?

El modo directo fuerza a Spoti NeXt a descargar **únicamente desde la fuente seleccionada**, sin intentar fuentes alternativas si la descarga falla. Es la contraparte del comportamiento por defecto, que usa fallback automático entre proveedores.

Está disponible solo cuando el downloader está configurado en una fuente específica (Tidal, Qobuz, Amazon o Deezer) — no aplica en modo `Auto`.

---

## ¿Dónde se activa?

**Settings → General → Downloader**

El toggle "Direct Mode" aparece habilitado únicamente si se seleccionó un downloader distinto de `Auto`. Si el downloader es `Auto`, el toggle se muestra deshabilitado con el texto `"Direct Mode disabled for Auto"`.

---

## Cómo funciona

Con modo directo **desactivado** (comportamiento por defecto):

```
Descarga por Tidal → falla
  → intenta Qobuz
  → intenta Amazon
  → marca como fallido si todos fallan
```

Con modo directo **activado**:

```
Descarga por Tidal → falla
  → marca como fallido inmediatamente
  → no intenta otras fuentes
```

Internamente, cuando `directMode` es `true`, el parámetro `allow_fallback` que se envía al backend se fuerza a `false`, independientemente del ajuste de "Quality Fallback".

---

## Cuándo usarlo

| Situación | Recomendación |
|-----------|---------------|
| Quiero garantizar calidad de Qobuz en cada descarga | Activar |
| Quiero identificar qué pistas no están en una fuente específica | Activar |
| Quiero máxima tasa de éxito sin importar la fuente | Desactivar |
| Descargo playlists mixtas | Desactivar |

---

## Relación con otras opciones

- **Quality Fallback** (`allowFallback`): controla el fallback *dentro* de una fuente (ej. HI_RES → LOSSLESS en Tidal). El modo directo sobrepasa este ajuste — si el modo directo está activo, no hay fallback de ningún tipo.
- **Auto Order**: irrelevante en modo directo. El orden de prioridad solo aplica cuando el downloader es `Auto`.

---

## Archivos relevantes

| Archivo | Función |
|---------|---------|
| `frontend/src/hooks/useDownload.ts` | Lógica de `isDirectMode`, fuerza `allow_fallback: false` |
| `frontend/src/components/SettingsPage.tsx` | Toggle en UI, deshabilitado si `downloader === "auto"` |
| `frontend/src/lib/settings.ts` | Campo `directMode: boolean` en el tipo `Settings` |
