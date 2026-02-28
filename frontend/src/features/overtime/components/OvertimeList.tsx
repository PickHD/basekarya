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
  Receipt,
  Plus,
  Filter,
  Eye,
  CalendarDays,
  Clock,
  FileSpreadsheet,
} from "lucide-react";
import { StatusBadge } from "./StatusBadge";
import { PaginationControls } from "@/components/shared/PaginationControls";
import { format, isValid } from "date-fns";
import { useProfile } from "@/features/user/hooks/useProfile";
import { useOvertimes, useExportOvertimes } from "@/features/overtime/hooks/useOvertime";
import { OvertimeDetailDialog } from "./OvertimeDetailDialog";
import { OvertimeFormDialog } from "./OvertimeCreateDialog";

export const OvertimeList = () => {
  const { data: user } = useProfile();
  const [page, setPage] = useState(1);
  const [statusFilter, setStatusFilter] = useState("");

  const { data, isLoading } = useOvertimes({
    page: page,
    limit: 10,
    status: statusFilter,
  });

  const [selectedId, setSelectedId] = useState<number | null>(null);
  const [isDetailOpen, setIsDetailOpen] = useState(false);
  const [isCreateOpen, setIsCreateOpen] = useState(false);

  const { mutate: exportOvertimes, isPending: isExporting } = useExportOvertimes();

  const handleExport = () => {
    exportOvertimes({ status: statusFilter });
  };

  const handleViewDetail = (id: number) => {
    setSelectedId(id);
    setIsDetailOpen(true);
  };

  const formatDateSafe = (dateStr: string, pattern: string) => {
    if (!dateStr) return "-";
    const date = new Date(dateStr);
    if (!isValid(date)) return "-";
    return format(date, pattern);
  };

  const formatDuration = (minutes: number) => {
    const hours = Math.floor(minutes / 60);
    const mins = minutes % 60;
    
    if (hours > 0 && mins > 0) return `${hours} jam ${mins} m`;
    if (hours > 0) return `${hours} jam`;
    return `${mins} menit`;
  };

  return (
    <div className="space-y-6">
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <div>
          <h2 className="text-2xl sm:text-3xl font-bold tracking-tight">
            Overtime
          </h2>
          <p className="text-sm sm:text-base text-slate-500">
            Kelola pengajuan lembur dan persetujuan.
          </p>
        </div>
        <div className="flex gap-2 w-full sm:w-auto">
          {user?.role === "SUPERADMIN" && (
            <Button
              onClick={handleExport}
              disabled={isExporting}
              className="bg-green-600 hover:bg-green-700 text-white w-full md:w-auto"
            >
              {isExporting ? (
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              ) : (
                <>
                  <FileSpreadsheet className="mr-2 h-4 w-4" /> Export Excel
                </>
              )}
            </Button>
          )}
          {user?.role !== "SUPERADMIN" && (
            <Button
              onClick={() => setIsCreateOpen(true)}
              className="bg-blue-600 hover:bg-blue-700 w-full sm:w-auto"
            >
              <Plus className="mr-2 h-4 w-4" /> Ajukan Lembur
            </Button>
          )}
        </div>
      </div>

      <Card>
        <CardHeader>
          <div className="flex flex-col md:flex-row justify-between md:items-center gap-4">
            <CardTitle className="flex items-center gap-2 text-lg">
              <Clock className="h-5 w-5" /> Daftar Lembur
            </CardTitle>

            <div className="flex gap-2 w-full md:w-auto">
              <div className="relative w-full md:w-48">
                <Filter className="absolute left-2 top-2.5 h-4 w-4 text-slate-500" />
                <select
                  className="h-10 w-full rounded-md border border-input bg-background pl-8 pr-3 text-sm ring-offset-background focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                  value={statusFilter}
                  onChange={(e) => {
                    setStatusFilter(e.target.value);
                    setPage(1);
                  }}
                >
                  <option value="">Semua Status</option>
                  <option value="PENDING">Pending</option>
                  <option value="APPROVED">Approved</option>
                  <option value="PAID">Paid</option>
                  <option value="REJECTED">Rejected</option>
                </select>
              </div>
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
              <div className="grid grid-cols-1 gap-4 md:hidden">
                {data?.data.map((item) => (
                  <div
                    key={item.id}
                    className="flex flex-col rounded-lg border bg-card p-4 shadow-sm space-y-3"
                  >
                    <div className="flex justify-between items-start gap-2">
                      <div>
                        <h4 className="font-semibold line-clamp-1">
                          {item.employee_name || 'Karyawan'}
                        </h4>
                        <div className="flex items-center text-xs text-slate-500 mt-1">
                          <CalendarDays className="mr-1 h-3 w-3" />
                          {formatDateSafe(item.date, "dd MMM yyyy")}
                        </div>
                      </div>
                      <StatusBadge status={item.status} />
                    </div>

                    <div className="space-y-1 mt-2">
                      <div className="flex items-center justify-between text-sm">
                        <span className="text-slate-500">Waktu:</span>
                        <span className="font-medium">{item.start_time} - {item.end_time}</span>
                      </div>
                      <div className="flex items-center justify-between text-sm">
                        <span className="text-slate-500">Durasi:</span>
                        <span className="font-bold text-blue-600">{formatDuration(item.duration_minutes)}</span>
                      </div>
                    </div>

                    <div className="pt-2 border-t">
                      <Button
                        variant="outline"
                        size="sm"
                        className="w-full"
                        onClick={() => handleViewDetail(item.id)}
                      >
                        <Eye className="mr-2 h-4 w-4" /> View Details
                      </Button>
                    </div>
                  </div>
                ))}
              </div>

              <div className="hidden md:block rounded-md border">
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>Employee</TableHead>
                      <TableHead>Tanggal</TableHead>
                      <TableHead>Waktu</TableHead>
                      <TableHead>Durasi</TableHead>
                      <TableHead>Status</TableHead>
                      <TableHead className="text-right">Action</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {data?.data.map((item) => (
                      <TableRow key={item.id}>
                        <TableCell className="font-medium">
                          {item.employee_name || "Karyawan"}
                          <div className="text-xs text-slate-400 font-normal">
                            Pengajuan: {formatDateSafe(item.created_at, "dd MMM yyyy")}
                          </div>
                        </TableCell>
                        <TableCell>
                          {formatDateSafe(item.date, "dd MMM yyyy")}
                        </TableCell>
                        <TableCell>
                          {item.start_time} - {item.end_time}
                        </TableCell>
                        <TableCell className="text-blue-600 font-medium">
                          {formatDuration(item.duration_minutes)}
                        </TableCell>
                        <TableCell>
                          <StatusBadge status={item.status} />
                        </TableCell>
                        <TableCell className="text-right">
                          <Button
                            variant="ghost"
                            size="icon"
                            className="hover:bg-slate-100"
                            onClick={() => handleViewDetail(item.id)}
                          >
                            <Eye className="h-4 w-4 text-slate-500" />
                          </Button>
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </div>

              {data?.data.length === 0 && (
                <div className="text-center py-10 text-slate-500 border rounded-md mt-4 md:mt-0">
                  <div className="flex flex-col items-center justify-center gap-2">
                    <Receipt className="h-10 w-10 text-slate-300" />
                    <p>No overtime requests found.</p>
                  </div>
                </div>
              )}

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

      <OvertimeFormDialog
        open={isCreateOpen}
        onOpenChange={setIsCreateOpen}
      />

      <OvertimeDetailDialog
        open={isDetailOpen}
        onOpenChange={setIsDetailOpen}
        overtimeId={selectedId}
      />
    </div>
  );
};
