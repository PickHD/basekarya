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
  Package,
  FileSpreadsheet,
  Settings,
} from "lucide-react";
import { AssetStatusBadge, AssetConditionBadge } from "./StatusBadge";
import { PaginationControls } from "@/components/shared/PaginationControls";
import { usePermissions } from "@/hooks/usePermissions";
import { PERMISSIONS } from "@/config/permissions";
import { useAssets, useExportAssets } from "@/features/asset/hooks/useAsset";
import { AssetDetailDialog } from "./AssetDetailDialog";
import { AssetCreateDialog } from "./AssetCreateDialog";
import { AssetCategoryDialog } from "./AssetCategoryDialog";

export const AssetList = () => {
  const { hasPermission } = usePermissions();
  const [page, setPage] = useState(1);
  const [statusFilter, setStatusFilter] = useState("");
  const [conditionFilter, setConditionFilter] = useState("");

  const { data, isLoading } = useAssets({
    page,
    limit: 10,
    status: statusFilter,
    condition: conditionFilter,
  });

  const [selectedId, setSelectedId] = useState<number | null>(null);
  const [isDetailOpen, setIsDetailOpen] = useState(false);
  const [isCreateOpen, setIsCreateOpen] = useState(false);
  const [isCategoryOpen, setIsCategoryOpen] = useState(false);

  const { mutate: exportAssets, isPending: isExporting } = useExportAssets();

  const handleExport = () => {
    exportAssets({
      status: statusFilter,
      condition: conditionFilter,
    });
  };

  const handleViewDetail = (id: number) => {
    setSelectedId(id);
    setIsDetailOpen(true);
  };

  return (
    <div className="space-y-6">
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <div>
          <h2 className="text-3xl font-bold tracking-tight">Asset Management</h2>
          <p className="text-slate-500">
            Kelola aset perusahaan, tracking status dan kondisi.
          </p>
        </div>
        <div className="flex gap-2 w-full sm:w-auto">
          {hasPermission(PERMISSIONS.EXPORT_ASSET) && (
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
          {hasPermission(PERMISSIONS.MANAGE_ASSET) && (
            <>
              <Button
                onClick={() => setIsCategoryOpen(true)}
                variant="outline"
                className="w-full md:w-auto"
              >
                <Settings className="mr-2 h-4 w-4" /> Kategori
              </Button>
              <Button
                onClick={() => setIsCreateOpen(true)}
                className="bg-blue-600 hover:bg-blue-700 w-full sm:w-auto"
              >
                <Plus className="mr-2 h-4 w-4" /> Tambah Aset
              </Button>
            </>
          )}
        </div>
      </div>

      <Card>
        <CardHeader>
          <div className="flex flex-col md:flex-row justify-between md:items-center gap-4">
            <CardTitle className="flex items-center gap-2 text-lg">
              <Package className="h-5 w-5" /> Daftar Aset
            </CardTitle>

            <div className="flex flex-wrap gap-2 w-full md:w-auto">
              <div className="relative w-full md:w-40">
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
                  <option value="AVAILABLE">Tersedia</option>
                  <option value="ASSIGNED">Digunakan</option>
                  <option value="MAINTENANCE">Perbaikan</option>
                  <option value="RETIRED">Dihapus</option>
                </select>
              </div>
              <div className="relative w-full md:w-40">
                <Filter className="absolute left-2 top-2.5 h-4 w-4 text-slate-500" />
                <select
                  className="h-10 w-full rounded-md border border-input bg-background pl-8 pr-3 text-sm ring-offset-background focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                  value={conditionFilter}
                  onChange={(e) => {
                    setConditionFilter(e.target.value);
                    setPage(1);
                  }}
                >
                  <option value="">Semua Kondisi</option>
                  <option value="GOOD">Baik</option>
                  <option value="FAIR">Cukup</option>
                  <option value="DAMAGED">Rusak</option>
                  <option value="LOST">Hilang</option>
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
              <div className="hidden md:block rounded-md border">
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>Nama Aset</TableHead>
                      <TableHead>Kategori</TableHead>
                      <TableHead>Serial Number</TableHead>
                      <TableHead>Status</TableHead>
                      <TableHead>Kondisi</TableHead>
                      <TableHead className="text-right">Aksi</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {data?.data.map((item) => (
                      <TableRow key={item.id}>
                        <TableCell className="font-medium">
                          {item.name}
                          {item.current_employee && (
                            <div className="text-xs text-blue-600">
                              Digunakan: {item.current_employee}
                            </div>
                          )}
                        </TableCell>
                        <TableCell>{item.category_name || "-"}</TableCell>
                        <TableCell className="text-slate-500 text-sm">
                          {item.serial_number || "-"}
                        </TableCell>
                        <TableCell>
                          <AssetStatusBadge status={item.status} />
                        </TableCell>
                        <TableCell>
                          <AssetConditionBadge condition={item.condition} />
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
                          Tidak ada aset ditemukan.
                        </TableCell>
                      </TableRow>
                    )}
                  </TableBody>
                </Table>
              </div>

              <div className="md:hidden space-y-3">
                {data?.data.map((item) => (
                  <Card key={item.id} className="p-4">
                    <div className="flex items-start justify-between">
                      <div className="space-y-2 flex-1">
                        <div className="flex justify-between items-start gap-2">
                          <div>
                            <h4 className="font-semibold line-clamp-1">{item.name}</h4>
                            <p className="text-xs text-slate-500">
                              {item.category_name} {item.serial_number && `| SN: ${item.serial_number}`}
                            </p>
                          </div>
                          <AssetStatusBadge status={item.status} />
                        </div>

                        <div className="flex items-center gap-2">
                          <AssetConditionBadge condition={item.condition} />
                        </div>

                        {item.current_employee && (
                          <p className="text-sm text-blue-600">Pengguna: {item.current_employee}</p>
                        )}
                      </div>
                      <Button
                        variant="ghost"
                        size="icon"
                        className="h-8 w-8"
                        onClick={() => handleViewDetail(item.id)}
                      >
                        <Eye className="h-4 w-4" />
                      </Button>
                    </div>
                  </Card>
                ))}
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

      <AssetCreateDialog open={isCreateOpen} onOpenChange={setIsCreateOpen} />
      <AssetDetailDialog open={isDetailOpen} onOpenChange={setIsDetailOpen} assetId={selectedId} />
      <AssetCategoryDialog open={isCategoryOpen} onOpenChange={setIsCategoryOpen} />
    </div>
  );
};
