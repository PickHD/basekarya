import { useState } from "react";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Loader2,
  CalendarDays,
  Plus,
  Filter,
  Eye,
  User,
  Clock,
} from "lucide-react";
import { PaginationControls } from "@/components/shared/PaginationControls";
import { format, isValid, parseISO } from "date-fns";
import { useProfile } from "@/features/user/hooks/useProfile";
import { useLeaves } from "@/features/leave/hooks/useLeave";
import { LeaveStatusBadge } from "@/features/leave/components/LeaveStatusBadge";
import { LeaveApplyDialog } from "@/features/leave/components/LeaveApplyDialog";
import { LeaveDetailDialog } from "@/features/leave/components/LeaveDetailDialog";
import { LeaveTypeBadge } from "@/features/leave/components/LeaveTypeBadge";

export const LeaveList = () => {
  const { data: user } = useProfile();
  const [page, setPage] = useState(1);
  const [statusFilter, setStatusFilter] = useState("");

  const { data, isLoading } = useLeaves({
    page,
    limit: 10,
    status: statusFilter,
  });

  const [isCreateOpen, setIsCreateOpen] = useState(false);
  const [selectedId, setSelectedId] = useState<number | null>(null);

  const handleViewDetail = (id: number) => {
    setSelectedId(id);
    setIsDetailOpen(true);
  };

  const [isDetailOpen, setIsDetailOpen] = useState(false);

  const formatDate = (dateStr: string) => {
    if (!dateStr) return "-";
    const date = parseISO(dateStr);
    if (!isValid(date)) return "-";
    return format(date, "dd MMM yyyy");
  };

  return (
    <div className="space-y-6">
      {/* HEADER */}
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <div>
          <h2 className="text-2xl sm:text-3xl font-bold tracking-tight">
            Leave Requests
          </h2>
          <p className="text-sm sm:text-base text-slate-500">
            Monitor employee leave and absence permissions.
          </p>
        </div>
        {user?.role !== "SUPERADMIN" && (
          <Button
            onClick={() => setIsCreateOpen(true)}
            className="bg-blue-600 hover:bg-blue-700 w-full sm:w-auto"
          >
            <Plus className="mr-2 h-4 w-4" /> Apply Leave
          </Button>
        )}
      </div>

      <Card>
        <CardHeader>
          <div className="flex flex-col md:flex-row justify-between md:items-center gap-4">
            <CardTitle className="flex items-center gap-2 text-lg">
              <CalendarDays className="h-5 w-5" /> Request History
            </CardTitle>

            <div className="relative w-full md:w-48">
              <Filter className="absolute left-2 top-2.5 h-4 w-4 text-slate-500" />
              <select
                className="h-10 w-full rounded-md border border-input bg-background pl-8 pr-3 text-sm focus:ring-2 focus:ring-blue-500 outline-none"
                value={statusFilter}
                onChange={(e) => {
                  setStatusFilter(e.target.value);
                  setPage(1);
                }}
              >
                <option value="">All Status</option>
                <option value="PENDING">Pending</option>
                <option value="APPROVED">Approved</option>
                <option value="REJECTED">Rejected</option>
              </select>
            </div>
          </div>
        </CardHeader>

        <CardContent>
          {isLoading ? (
            <div className="flex justify-center py-10">
              <Loader2 className="animate-spin h-8 w-8 text-blue-600" />
            </div>
          ) : (
            <>
              {/* --- MOBILE CARD VIEW --- */}
              <div className="grid grid-cols-1 gap-4 md:hidden">
                {data?.data.map((item) => (
                  <div
                    key={item.id}
                    className="flex flex-col rounded-lg border bg-card p-4 shadow-sm space-y-3"
                  >
                    {/* Header: Type & Status */}
                    <div className="flex justify-between items-start">
                      <div>
                        <h4 className="font-bold text-slate-800">
                          {item.leave_type?.name}
                        </h4>
                        <span className="text-xs text-slate-500 flex items-center gap-1 mt-1">
                          <Clock className="h-3 w-3" /> {item.total_days} Days
                        </span>
                      </div>
                      <LeaveStatusBadge status={item.status} />
                    </div>

                    <div className="flex items-center gap-2 text-sm text-slate-700">
                      <User className="h-4 w-4 text-slate-400" />
                      <span className="font-medium">
                        {item.employee_name || "Me"}
                      </span>
                    </div>

                    <div className="bg-slate-50 p-3 rounded text-sm grid grid-cols-2 gap-2 text-center">
                      <div>
                        <div className="text-xs text-slate-500">Start</div>
                        <div className="font-medium">
                          {formatDate(item.start_date)}
                        </div>
                      </div>
                      <div>
                        <div className="text-xs text-slate-500">End</div>
                        <div className="font-medium">
                          {formatDate(item.end_date)}
                        </div>
                      </div>
                    </div>

                    <Button
                      variant="outline"
                      size="sm"
                      className="w-full"
                      onClick={() => handleViewDetail(item.id)}
                    >
                      <Eye className="mr-2 h-4 w-4" /> View Details
                    </Button>
                  </div>
                ))}
              </div>

              <div className="hidden md:block rounded-md border">
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>Requested On</TableHead>
                      <TableHead>Employee</TableHead>
                      <TableHead>Leave Type</TableHead>
                      <TableHead>Duration</TableHead>
                      <TableHead>Status</TableHead>
                      <TableHead className="text-right">Action</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {data?.data.map((item) => (
                      <TableRow key={item.id}>
                        <TableCell className="text-slate-500">
                          {formatDate(item.created_at)}
                        </TableCell>
                        <TableCell>
                          <div className="font-bold">
                            {item.employee_name || "-"}
                          </div>
                          <div className="text-xs text-slate-500">
                            {item.employee_nik}
                          </div>
                        </TableCell>
                        <TableCell>
                          <LeaveTypeBadge status={item.leave_type?.name} />
                        </TableCell>
                        <TableCell>
                          <div className="flex flex-col">
                            <span className="font-medium">
                              {formatDate(item.start_date)} -{" "}
                              {formatDate(item.end_date)}
                            </span>
                            <span className="text-xs text-slate-500">
                              {item.total_days} days total
                            </span>
                          </div>
                        </TableCell>
                        <TableCell>
                          <LeaveStatusBadge status={item.status} />
                        </TableCell>
                        <TableCell className="text-right">
                          <Button
                            variant="ghost"
                            size="icon"
                            onClick={() => handleViewDetail(item.id)}
                          >
                            <Eye className="h-4 w-4 text-slate-500" />
                          </Button>
                        </TableCell>
                      </TableRow>
                    ))}
                    {data?.data.length === 0 && (
                      <TableRow>
                        <TableCell
                          colSpan={6}
                          className="text-center py-8 text-slate-500"
                        >
                          No leave requests found.
                        </TableCell>
                      </TableRow>
                    )}
                  </TableBody>
                </Table>
              </div>

              {data?.meta && (
                <div className="mt-4">
                  <PaginationControls
                    meta={{
                      limit: 10,
                      page: data.meta.page,
                      total_page: data.meta.total_page,
                      total_data: data.meta.total_data,
                    }}
                    onPageChange={setPage}
                    isLoading={isLoading}
                  />
                </div>
              )}
            </>
          )}
        </CardContent>
      </Card>

      <LeaveApplyDialog open={isCreateOpen} onOpenChange={setIsCreateOpen} />

      <LeaveDetailDialog
        open={isDetailOpen}
        onOpenChange={setIsDetailOpen}
        leaveId={selectedId}
      />
    </div>
  );
};
