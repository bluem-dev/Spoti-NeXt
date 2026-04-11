import { Badge } from "@/components/ui/badge";
import { Tooltip, TooltipContent, TooltipTrigger, } from "@/components/ui/tooltip";
import { openExternal } from "@/lib/utils";
import { formatRelativeTime } from "@/lib/relative-time";
import { getSettings } from "@/lib/settings";
import { t } from "@/lib/i18n";
import { TidalIcon, QobuzIcon, AmazonIcon } from "./PlatformIcons";
interface HeaderProps {
    version: string;
    hasUpdate: boolean;
    releaseDate?: string | null;
}
const SOURCE_LABELS: Record<string, string> = {
    tidal: "Tidal",
    qobuz: "Qobuz",
    amazon: "Amazon Music",
};
const SOURCE_ICONS: Record<string, React.ReactNode> = {
    tidal: <TidalIcon className="w-3 h-3" />,
    qobuz: <QobuzIcon className="w-3 h-3" />,
    amazon: <AmazonIcon className="w-3 h-3" />,
};
export function Header({ version, hasUpdate, releaseDate }: HeaderProps) {
    const settings = getSettings();
    const isDirectMode = settings.directMode && settings.downloader !== "auto";
    const activeSource = isDirectMode ? settings.downloader : null;
    return (<div className="relative">
      <div className="text-center space-y-2">
        <div className="flex items-center justify-center gap-3">
          <img src="/icon.svg" alt="Spoti NeXt" className="w-12 h-12 cursor-pointer" onClick={() => window.location.reload()}/>
          <h1 className="text-4xl font-bold cursor-pointer" onClick={() => window.location.reload()}>
            Spoti NeXt
          </h1>
          <div className="relative">
            <Tooltip>
              <TooltipTrigger asChild>
                <Badge variant="default" asChild>
                  <button type="button" onClick={() => openExternal("https://github.com/afkarxyz/SpotiFLAC/releases")} className="cursor-pointer hover:opacity-80 transition-opacity">
                    v{version}
                  </button>
                </Badge>
              </TooltipTrigger>
              {hasUpdate && releaseDate && (<TooltipContent>
                  <p>{formatRelativeTime(releaseDate)}</p>
                </TooltipContent>)}
            </Tooltip>
            {hasUpdate && (<span className="absolute -top-1 -right-1 flex h-3 w-3">
                <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-green-400 opacity-75"></span>
                <span className="relative inline-flex rounded-full h-3 w-3 bg-green-500"></span>
              </span>)}
          </div>
          {isDirectMode && activeSource && (
            <Tooltip>
              <TooltipTrigger asChild>
                <Badge variant="outline" className="flex items-center gap-1.5 text-xs font-medium cursor-default select-none">
                  {SOURCE_ICONS[activeSource]}
                  {SOURCE_LABELS[activeSource]}
                </Badge>
              </TooltipTrigger>
              <TooltipContent>
                <p>{t("directModeActive")} — {t("usingDirectSource").replace("{source}", SOURCE_LABELS[activeSource] ?? activeSource)}</p>
              </TooltipContent>
            </Tooltip>
          )}
        </div>
        <p className="text-muted-foreground">
          {t("headerSubtitle")}
        </p>
      </div>
    </div>);
}
