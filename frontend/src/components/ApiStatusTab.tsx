import { t } from "@/lib/i18n";
import { Button } from "@/components/ui/button";
import { RefreshCw, CheckCircle2, XCircle, Loader2 } from "lucide-react";
import { TidalIcon, QobuzIcon, AmazonIcon } from "./PlatformIcons";
import { useApiStatus } from "@/hooks/useApiStatus";

function AtmosBadge({ mode }: { mode: string }) {
    if (mode === "checking") {
        return (
            <span className="flex items-center gap-1 text-xs text-muted-foreground">
                <Loader2 className="h-3 w-3 animate-spin" />
                {t("atmosChecking")}
            </span>
        );
    }
    if (mode === "DOLBY_ATMOS") {
        return (
            <span className="inline-flex items-center rounded-full px-2 py-0.5 text-xs font-semibold bg-sky-500/15 text-sky-400 ring-1 ring-inset ring-sky-500/30">
                {t("atmosSupported")}
            </span>
        );
    }
    if (mode === "STEREO" || mode === "SONY_360RA") {
        return (
            <span className="inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium bg-muted text-muted-foreground ring-1 ring-inset ring-border">
                {t("atmosNotSupported")}
            </span>
        );
    }
    if (mode === "UNKNOWN") {
        return (
            <span className="text-xs text-muted-foreground">{t("atmosUnknown")}</span>
        );
    }
    // empty string — not yet probed (API offline or check not started)
    return null;
}

export function ApiStatusTab() {
    const { sources, statuses, atmosModes, isCheckingAll, refreshAll } = useApiStatus();
    return (
        <div className="space-y-6">
            <div className="flex items-center justify-end">
                <Button variant="outline" onClick={() => void refreshAll()} disabled={isCheckingAll} className="gap-2">
                    <RefreshCw className={`h-4 w-4 ${isCheckingAll ? "animate-spin" : ""}`} />
                    {t("refreshAll")}
                </Button>
            </div>

            <div className="grid grid-cols-4 gap-4">
                {sources.map((source) => {
                    const status = statuses[source.id] || "idle";
                    const atmosMode = atmosModes?.[source.id] ?? "";
                    const isTidal = source.type === "tidal";

                    return (
                        <div key={source.id} className="flex items-center justify-between p-4 border rounded-lg bg-card text-card-foreground shadow-sm gap-3">
                            <div className="flex items-center gap-3 min-w-0">
                                {source.type === "tidal"
                                    ? <TidalIcon className="w-5 h-5 shrink-0 text-muted-foreground" />
                                    : source.type === "amazon"
                                        ? <AmazonIcon className="w-5 h-5 shrink-0 text-muted-foreground" />
                                        : <QobuzIcon className="w-5 h-5 shrink-0 text-muted-foreground" />}
                                <div className="flex flex-col gap-1 min-w-0">
                                    <p className="font-medium leading-none">{source.name}</p>
                                    {isTidal && status === "online" && (
                                        <AtmosBadge mode={atmosMode} />
                                    )}
                                </div>
                            </div>

                            <div className="flex items-center shrink-0">
                                {status === "checking" && <Loader2 className="h-5 w-5 animate-spin text-muted-foreground" />}
                                {status === "online" && <CheckCircle2 className="h-5 w-5 text-emerald-500" />}
                                {status === "offline" && <XCircle className="h-5 w-5 text-destructive" />}
                                {status === "idle" && <div className="h-5 w-5 rounded-full bg-muted" />}
                            </div>
                        </div>
                    );
                })}
            </div>
        </div>
    );
}
