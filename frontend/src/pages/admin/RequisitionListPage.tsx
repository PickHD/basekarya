import { useState } from "react";
import { Plus, Search, Briefcase } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { PaginationControls } from "@/components/shared/PaginationControls";
import { RequisitionList } from "@/features/recruitment/components/RequisitionList";
import { RequisitionFormDialog } from "@/features/recruitment/components/RequisitionFormDialog";
import { RequisitionDetailDialog } from "@/features/recruitment/components/RequisitionDetailDialog";
import { useRequisitions } from "@/features/recruitment/hooks/useRequisition";
import { usePermissions } from "@/hooks/usePermissions";
import { PERMISSIONS } from "@/config/permissions";
import { useDebounce } from "@/hooks/useDebounce";
import type { JobRequisition } from "@/features/recruitment/types";

export default function RequisitionListPage() {
  const { hasPermission } = usePermissions();
  const [page, setPage] = useState(1);
  const [search, setSearch] = useState("");
  const debouncedSearch = useDebounce(search, 500);
  const [statusFilter, setStatusFilter] = useState("ALL");
  const [priorityFilter, setPriorityFilter] = useState("ALL");

  const [isFormOpen, setIsFormOpen] = useState(false);
  const [selectedRequisition, setSelectedRequisition] = useState<JobRequisition | null>(null);

  const { data, isLoading } = useRequisitions({
    page,
    limit: 10,
    search: debouncedSearch || undefined,
    status: statusFilter !== "ALL" ? statusFilter : undefined,
    priority: priorityFilter !== "ALL" ? priorityFilter : undefined,
  });

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <div>
          <h2 className="text-3xl font-bold tracking-tight">Job Requisitions</h2>
          <p className="text-slate-500 mt-1">Manage job opening requests and approvals.</p>
        </div>
        {hasPermission(PERMISSIONS.CREATE_REQUISITION) && (
          <Button onClick={() => setIsFormOpen(true)} className="bg-blue-600 hover:bg-blue-700">
            <Plus className="mr-2 h-4 w-4" />
            New Requisition
          </Button>
        )}
      </div>

      {/* Card */}
      <Card>
        <CardHeader>
          <div className="flex flex-col sm:flex-row justify-between items-center gap-4">
            <CardTitle className="flex items-center gap-2">
              <Briefcase className="h-5 w-5" />
              Requisition List
            </CardTitle>
            <div className="flex flex-col sm:flex-row gap-3 w-full sm:w-auto">
              {/* Search */}
              <div className="relative w-full sm:w-64">
                <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-slate-400" />
                <Input
                  placeholder="Search title..."
                  className="pl-9"
                  value={search}
                  onChange={(e) => setSearch(e.target.value)}
                />
              </div>

              {/* Status filter */}
              <Select value={statusFilter} onValueChange={setStatusFilter}>
                <SelectTrigger className="w-full sm:w-36">
                  <SelectValue placeholder="Status" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="ALL">All Status</SelectItem>
                  <SelectItem value="DRAFT">Draft</SelectItem>
                  <SelectItem value="PENDING">Pending</SelectItem>
                  <SelectItem value="APPROVED">Approved</SelectItem>
                  <SelectItem value="REJECTED">Rejected</SelectItem>
                  <SelectItem value="CLOSED">Closed</SelectItem>
                </SelectContent>
              </Select>

              {/* Priority filter */}
              <Select value={priorityFilter} onValueChange={setPriorityFilter}>
                <SelectTrigger className="w-full sm:w-36">
                  <SelectValue placeholder="Priority" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="ALL">All Priority</SelectItem>
                  <SelectItem value="LOW">Low</SelectItem>
                  <SelectItem value="MEDIUM">Medium</SelectItem>
                  <SelectItem value="HIGH">High</SelectItem>
                  <SelectItem value="URGENT">Urgent</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>
        </CardHeader>
        <CardContent>
          <RequisitionList
            data={data?.data || []}
            isLoading={isLoading}
            onView={setSelectedRequisition}
          />
          {data?.meta && (
            <PaginationControls
              meta={data.meta}
              onPageChange={setPage}
              isLoading={isLoading}
            />
          )}
        </CardContent>
      </Card>

      {/* Dialogs */}
      <RequisitionFormDialog open={isFormOpen} onOpenChange={setIsFormOpen} />
      <RequisitionDetailDialog
        open={!!selectedRequisition}
        onOpenChange={(open) => !open && setSelectedRequisition(null)}
        requisition={selectedRequisition}
      />
    </div>
  );
}
