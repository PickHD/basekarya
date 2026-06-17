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
import { Card } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Eye, Plus, Search, Loader2, GraduationCap } from "lucide-react";
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
    if (s === "COMPLETED") return <Badge variant="secondary" className="text-xs">Completed</Badge>;
    return <Badge variant="default" className="text-xs">In Progress</Badge>;
  };

  return (
    <div className="space-y-4">
      <div className="flex flex-col sm:flex-row gap-3 items-start sm:items-center justify-between">
        <div className="flex gap-2 flex-1 max-w-md">
          <div className="relative flex-1">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
            <Input
              className="pl-9"
              placeholder="Search by name or email..."
              value={search}
              onChange={(e) => { setSearch(e.target.value); setPage(1); }}
            />
          </div>
          <Select value={status} onValueChange={(v) => { setStatus(v === "ALL" ? "" : v); setPage(1); }}>
            <SelectTrigger className="w-36">
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
          <Button onClick={onCreateNew} className="bg-blue-600 hover:bg-blue-700">
            <Plus className="mr-2 h-4 w-4" /> New Onboarding
          </Button>
        )}
      </div>

      {isLoading ? (
        <div className="flex items-center justify-center py-8">
          <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
        </div>
      ) : workflows.length === 0 ? (
        <div className="flex flex-col items-center justify-center py-8 text-muted-foreground">
          <GraduationCap className="h-12 w-12 mb-2" />
          <p>No onboarding workflows found</p>
          {canCreate && <p className="text-sm">Create a new onboarding workflow to get started.</p>}
        </div>
      ) : (
        <>
          <div className="hidden md:block">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>New Hire</TableHead>
                  <TableHead>Position</TableHead>
                  <TableHead>Department</TableHead>
                  <TableHead>Start Date</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead>Progress</TableHead>
                  <TableHead className="w-20 text-right">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {workflows.map((w) => (
                  <TableRow key={w.id}>
                    <TableCell>
                      <div>
                        <p className="font-medium">{w.new_hire_name}</p>
                        <p className="text-xs text-muted-foreground">{w.new_hire_email}</p>
                      </div>
                    </TableCell>
                    <TableCell>{w.position || "—"}</TableCell>
                    <TableCell>{w.department || "—"}</TableCell>
                    <TableCell>
                      {w.start_date ? format(new Date(w.start_date), "dd MMM yyyy") : "—"}
                    </TableCell>
                    <TableCell>{statusBadge(w.status)}</TableCell>
                    <TableCell className="min-w-[140px]">
                      <div className="flex items-center gap-2">
                        <div className="flex-1 h-2 bg-slate-100 rounded-full overflow-hidden">
                          <div
                            className={`h-full rounded-full transition-all ${w.progress === 100 ? "bg-emerald-500" : "bg-blue-500"}`}
                            style={{ width: `${w.progress}%` }}
                          />
                        </div>
                        <span className="text-xs text-muted-foreground w-8 text-right">{w.progress}%</span>
                      </div>
                    </TableCell>
                    <TableCell className="text-right">
                      <Button
                        variant="ghost"
                        size="icon"
                        onClick={() => onView(w)}
                        className="h-8 w-8"
                      >
                        <Eye className="h-4 w-4" />
                      </Button>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </div>

          <div className="md:hidden space-y-3">
            {workflows.map((w) => (
              <Card key={w.id} className="p-4">
                <div className="flex items-start justify-between">
                  <div className="space-y-1">
                    <p className="font-medium">{w.new_hire_name}</p>
                    <p className="text-sm text-muted-foreground">{w.new_hire_email}</p>
                    <div className="flex flex-wrap gap-2 text-sm text-muted-foreground">
                      {w.position && <span>{w.position}</span>}
                      {w.department && <span>· {w.department}</span>}
                    </div>
                    <div className="flex items-center gap-2 pt-1">
                      {statusBadge(w.status)}
                      {w.start_date && (
                        <span className="text-xs text-muted-foreground">
                          {format(new Date(w.start_date), "dd MMM yyyy")}
                        </span>
                      )}
                    </div>
                    <div className="flex items-center gap-2 pt-1">
                      <div className="flex-1 max-w-[120px] h-2 bg-slate-100 rounded-full overflow-hidden">
                        <div
                          className={`h-full rounded-full transition-all ${w.progress === 100 ? "bg-emerald-500" : "bg-blue-500"}`}
                          style={{ width: `${w.progress}%` }}
                        />
                      </div>
                      <span className="text-xs text-muted-foreground">{w.progress}%</span>
                    </div>
                  </div>
                  <Button
                    variant="ghost"
                    size="icon"
                    onClick={() => onView(w)}
                    className="h-8 w-8 flex-shrink-0"
                  >
                    <Eye className="h-4 w-4" />
                  </Button>
                </div>
              </Card>
            ))}
          </div>

          {meta && (
            <PaginationControls
              meta={meta}
              onPageChange={setPage}
            />
          )}
        </>
      )}
    </div>
  );
}
