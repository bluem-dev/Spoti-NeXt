import { t } from "@/lib/i18n";
import { Input } from "@/components/ui/input";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue, } from "@/components/ui/select";
import { Search, ArrowUpDown, XCircle } from "lucide-react";
interface SearchAndSortProps {
    searchQuery: string;
    sortBy: string;
    onSearchChange: (value: string) => void;
    onSortChange: (value: string) => void;
}
export function SearchAndSort({ searchQuery, sortBy, onSearchChange, onSortChange, }: SearchAndSortProps) {
    return (<div className="flex gap-2">
      <div className="relative flex-1">
        <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground"/>
        <Input placeholder={t("searchTracks")} value={searchQuery} onChange={(e) => onSearchChange(e.target.value)} className="pl-10 pr-8"/>
        {searchQuery && (<button type="button" className="absolute right-2 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground transition-colors cursor-pointer" onClick={() => onSearchChange("")}>
            <XCircle className="h-4 w-4"/>
          </button>)}
      </div>
      <Select value={sortBy} onValueChange={onSortChange}>
        <SelectTrigger className="w-[200px] gap-1.5">
          <ArrowUpDown className="h-4 w-4"/>
          <SelectValue placeholder={t("sortBy")}/>
        </SelectTrigger>
        <SelectContent>
          <SelectItem value="default">{t("filterDefault")}</SelectItem>
          <SelectItem value="title-asc">{t("sortTitleAZ")}</SelectItem>
          <SelectItem value="title-desc">{t("sortTitleZA")}</SelectItem>
          <SelectItem value="artist-asc">{t("sortArtistAZ")}</SelectItem>
          <SelectItem value="artist-desc">{t("sortArtistZA")}</SelectItem>
          <SelectItem value="duration-asc">{t("sortDurationShort")}</SelectItem>
          <SelectItem value="duration-desc">{t("sortDurationLong")}</SelectItem>
          <SelectItem value="plays-asc">{t("playsLow")}</SelectItem>
          <SelectItem value="plays-desc">{t("playsHigh")}</SelectItem>
          <SelectItem value="downloaded">{t("filterDownloaded")}</SelectItem>
          <SelectItem value="not-downloaded">{t("filterNotDownloaded")}</SelectItem>
          <SelectItem value="failed">{t("filterFailedDownloads")}</SelectItem>
        </SelectContent>
      </Select>
    </div>);
}
