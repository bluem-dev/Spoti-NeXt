/**
 * wails-extra.ts
 *
 * Bindings manuales para funciones de Go que Wails omite del autogenerado
 * (App.d.ts / App.js). Este archivo NO es sobrescrito por "wails generate".
 *
 * Patrón idéntico al de wailsjs/go/main/App.js generado por Wails.
 */

import type { backend } from "../../wailsjs/go/models";

export interface QobuzSearchRequest {
    query: string;
    limit: number;
}

export function SearchQobuz(req: QobuzSearchRequest): Promise<Array<backend.SearchResult>> {
    return (window as any)["go"]["main"]["App"]["SearchQobuz"](req);
}
