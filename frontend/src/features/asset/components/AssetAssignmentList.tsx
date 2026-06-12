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
  Plus,
  Filter,
  Eye,
  ClipboardList,
} from "lucide-react";
import { AssetAssignmentStatusBadge } from "./StatusBadge";
import { PaginationControls } from "@/components/shared/PaginationControls";
import { format, isValid } from "date-fns";
import { usePermissions } from "@/hooks/usePermissions";
import { PERMISSIONS } from "@/config/permissions";
import { useAssetAssignments } from "@/features/asset/hooks/useAsset";
import { AssetAssignmentDetailDialog } from "./AssetAssignmentDetailDialog";
import { AssetAssignmentCreateDialog } from "./AssetAssignmentCreateDialog";

export const AssetAssignmentList = () => {
  const { hasPermission } = usePermissions();
  const [page, setPage] = useState(1);
  const [statusFilter, setStatusFilter] = useState("");

  const { data, isLoading } = useAssetAssignments({
    page,
    limit: 10,
    status: statusFilter,
  });

  const [selectedId, setSelectedId] = useState<number | null>(null);
  const [isDetailOpen, setIsDetailOpen] = useState(false);
  const [isCreateOpen, setIsCreateOpen] = useState(false);

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

  return (
    <div className="space-y-6">
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <div>
          <h2 className="text-2xl sm:text-3xl font-bold tracking-tight">Permintaan Aset</h2>
          <p className="text-sm sm:text-base text-slate-500">
            Monitor permintaan dan pengembalian aset karyawan.
          </p>
        </div>
        <div className="flex gap-2 w-full sm:w-auto">
          {hasPermission(PERMISSIONS.CREATE_ASSET) && (
            <Button
              onClick={() => setIsCreateOpen(true)}
              className="bg-blue-600 hover:bg-blue-700 w-full sm:w-auto"
            >
              <Plus className="mr-2 h-4 w-4" /> Ajukan Permintaan
            </Button>
          )}
        </div>
      </div>

      <Card>
        <CardHeader>
          <div className="flex flex-col md:flex-row justify-between md:items-center gap-4">
            <CardTitle className="flex items-center gap-2 text-lg">
              <ClipboardList className="h-5 w-5" /> Riwayat Permintaan
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
                  <option value="ACTIVE">Active</option>
                  <option value="RETURNED">Returned</option>
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
                        <h4 className="font-semibold line-clamp-1">{item.employee_name || "Karyawan"}</h4>
                        <p className="text-xs text-slate-500">{item.employee_nik}</p>
                      </div>
                      <AssetAssignmentStatusBadge status={item.status} />
                    </div>

                    <div className="space-y-1">
                      <p className="text-sm"><span className="text-slate-500">Aset:</span> {item.asset_name}</p>
                      {item.expected_return_date && (
                        <p className="text-sm text-slate-500">
                          Estimasi Kembali: {formatDateSafe(item.expected_return_date, "dd MMM yyyy")}
                        </p>
                      )}
                    </div>

                    <div className="pt-2 border-t">
                      <Button
                        variant="outline"
                        size="sm"
                        className="w-full"
                        onClick={() => handleViewDetail(item.id)}
                      >
                        <Eye className="mr-2 h-4 w-4" /> Detail
                      </Button>
                    </div>
                  </div>
                ))}
              </div>

              <div className="hidden md:block rounded-md border">
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>Karyawan</TableHead>
                      <TableHead>Aset</TableHead>
                      <TableHead>Tujuan</TableHead>
                      <TableHead>Estimasi Kembali</TableHead>
                      <TableHead>Status</TableHead>
                      <TableHead className="text-right">Aksi</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {data?.data.map((item) => (
                      <TableRow key={item.id}>
                        <TableCell className="font-medium">
                          {item.employee_name || "Karyawan"}
                          <div className="text-xs text-slate-400 font-normal">
                            NIK: {item.employee_nik || "-"}
                          </div>
                        </TableCell>
                        <TableCell>{item.asset_name || "-"}</TableCell>
                        <TableCell className="text-sm max-w-[200px] truncate">
                          {item.purpose || "-"}
                        </TableCell>
                        <TableCell className="text-sm">
                          {item.expected_return_date
                            ? formatDateSafe(item.expected_return_date, "dd MMM yyyy")
                            : "-"}
                        </TableCell>
                        <TableCell>
                          <AssetAssignmentStatusBadge status={item.status} />
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
                    {data?.data.length === 0 && (
                      <TableRow>
                        <TableCell colSpan={6} className="text-center py-8 text-slate-500">
                          Tidak ada permintaan aset ditemukan.
                        </TableCell>
                      </TableRow>
                    )}
                  </TableBody>
                </Table>
              </div>

              {data?.meta && (
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
              )}
            </>
          )}
        </CardContent>
      </Card>

      <AssetAssignmentCreateDialog
        open={isCreateOpen}
        onOpenChange={setIsCreateOpen}
      />
      <AssetAssignmentDetailDialog
        open={isDetailOpen}
        onOpenChange={setIsDetailOpen}
        assignmentId={selectedId}
      />
    </div>
  );
};
