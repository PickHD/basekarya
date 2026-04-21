import { useState } from "react";
import { format } from "date-fns";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Eye, Plus, Search, Loader2 } from "lucide-react";
import { useOnboardingWorkflows } from "@/features/onboarding/hooks/useOnboarding";
import type { OnboardingWorkflowList } from "@/features/onboarding/types";
import { PaginationControls } from "@/components/shared/PaginationControls";
import { useDebounce } from "@/hooks/useDebounce";

interface Props {
  onView: (workflow: OnboardingWorkflowList) => void;
  onCreateNew: () => void;
  canCreate: boolean;
}

export function OnboardingWorkflowList({ onView, onCreateNew, canCreate }: Props) {
  const [search, setSearch] = useState("");
  const [status, setStatus] = useState("");
  const [page, setPage] = useState(1);
  const limit = 10;

  const debouncedSearch = useDebounce(search, 500);

  const { data, isLoading } = useOnboardingWorkflows({ search: debouncedSearch, status, page, limit });
  const workflows: OnboardingWorkflowList[] = data?.data ?? [];
  const meta = data?.meta;

  const statusBadge = (s: string) => {
    if (s === "COMPLETED") return <Badge className="bg-emerald-100 text-emerald-700 border-emerald-200">Completed</Badge>;
    return <Badge className="bg-blue-100 text-blue-700 border-blue-200">In Progress</Badge>;
  };

  const progressBar = (pct: number) => (
    <div className="flex items-center gap-2">
      <div className="flex-1 h-2 bg-slate-100 rounded-full overflow-hidden">
        <div
          className={`h-full rounded-full transition-all ${pct === 100 ? "bg-emerald-500" : "bg-blue-500"}`}
          style={{ width: `${pct}%` }}
        />
      </div>
      <span className="text-xs text-slate-500 w-8 text-right">{pct}%</span>
    </div>
  );

  return (
    <div className="space-y-4">
      {/* Toolbar */}
      <div className="flex flex-col sm:flex-row gap-3 items-start sm:items-center justify-between">
        <div className="flex gap-2 flex-1 max-w-md">
          <div className="relative flex-1">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-slate-400" />
            <Input
              className="pl-9 text-sm"
              placeholder="Search by name or email..."
              value={search}
              onChange={(e) => { setSearch(e.target.value); setPage(1); }}
            />
          </div>
          <Select value={status} onValueChange={(v) => { setStatus(v === "ALL" ? "" : v); setPage(1); }}>
            <SelectTrigger className="w-36 text-sm">
              <SelectValue placeholder="Status" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="ALL">All Status</SelectItem>
              <SelectItem value="IN_PROGRESS">In Progress</SelectItem>
              <SelectItem value="COMPLETED">Completed</SelectItem>
            </SelectContent>
          </Select>
        </div>
        {canCreate && (
          <Button onClick={onCreateNew} className="bg-blue-600 hover:bg-blue-700 text-sm">
            <Plus className="mr-2 h-4 w-4" />
            New Onboarding
          </Button>
        )}
      </div>

      {/* Table */}
      <div className="rounded-xl border border-slate-200 overflow-hidden">
        <Table>
          <TableHeader>
            <TableRow className="bg-slate-50">
              <TableHead className="text-xs font-semibold text-slate-600">New Hire</TableHead>
              <TableHead className="text-xs font-semibold text-slate-600">Position</TableHead>
              <TableHead className="text-xs font-semibold text-slate-600">Department</TableHead>
              <TableHead className="text-xs font-semibold text-slate-600">Start Date</TableHead>
              <TableHead className="text-xs font-semibold text-slate-600">Status</TableHead>
              <TableHead className="text-xs font-semibold text-slate-600">Progress</TableHead>
              <TableHead className="text-xs font-semibold text-slate-600 text-right">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {isLoading ? (
              <TableRow>
                <TableCell colSpan={7} className="text-center py-16">
                  <Loader2 className="h-5 w-5 animate-spin text-slate-400 mx-auto" />
                </TableCell>
              </TableRow>
            ) : workflows.length === 0 ? (
              <TableRow>
                <TableCell colSpan={7} className="text-center py-16 text-slate-400 text-sm">
                  No onboarding workflows found.
                </TableCell>
              </TableRow>
            ) : (
              workflows.map((w) => (
                <TableRow key={w.id} className="hover:bg-slate-50 transition-colors">
                  <TableCell>
                    <div>
                      <p className="text-sm font-medium text-slate-800">{w.new_hire_name}</p>
                      <p className="text-xs text-slate-500">{w.new_hire_email}</p>
                    </div>
                  </TableCell>
                  <TableCell className="text-sm text-slate-600">{w.position || "—"}</TableCell>
                  <TableCell className="text-sm text-slate-600">{w.department || "—"}</TableCell>
                  <TableCell className="text-sm text-slate-600">
                    {w.start_date ? format(new Date(w.start_date), "dd MMM yyyy") : "—"}
                  </TableCell>
                  <TableCell>{statusBadge(w.status)}</TableCell>
                  <TableCell className="min-w-[140px]">{progressBar(w.progress)}</TableCell>
                  <TableCell className="text-right">
                    <Button
                      variant="ghost"
                      size="sm"
                      className="h-8 text-xs"
                      onClick={() => onView(w)}
                    >
                      <Eye className="h-3.5 w-3.5 mr-1.5" />
                      View
                    </Button>
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </div>

      {meta && (
        <PaginationControls
          meta={meta}
          onPageChange={setPage}
        />
      )}
    </div>
  );
}
