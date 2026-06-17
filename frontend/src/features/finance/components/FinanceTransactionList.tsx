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
  Calendar,
  Wallet,
  FileSpreadsheet,
} from "lucide-react";
import { StatusBadge } from "./StatusBadge";
import { PaginationControls } from "@/components/shared/PaginationControls";
import { format, isValid } from "date-fns";
import { usePermissions } from "@/hooks/usePermissions";
import { PERMISSIONS } from "@/config/permissions";
import {
  useFinanceTransactions,
  useExportFinanceTransactions,
} from "@/features/finance/hooks/useFinance";
import { FinanceTransactionDetailDialog } from "./FinanceTransactionDetailDialog";
import { FinanceTransactionCreateDialog } from "./FinanceTransactionCreateDialog";
import { FinanceCategoryManager } from "./FinanceCategoryManager";

export const FinanceTransactionList = () => {
  const { hasPermission } = usePermissions();
  const [statusFilter, setStatusFilter] = useState("");
  const [typeFilter, setTypeFilter] = useState("");

  const {
    data,
    isLoading,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage,
  } = useFinanceTransactions({
    limit: 10,
    status: statusFilter,
    type: typeFilter,
  });

  const allTransactions = data?.pages.flatMap((page) => page.data) ?? [];

  const [selectedId, setSelectedId] = useState<number | null>(null);
  const [isDetailOpen, setIsDetailOpen] = useState(false);
  const [isCreateOpen, setIsCreateOpen] = useState(false);
  const [isCategoryOpen, setIsCategoryOpen] = useState(false);

  const { mutate: exportData, isPending: isExporting } = useExportFinanceTransactions();

  const handleExport = () => {
    exportData({ status: statusFilter, type: typeFilter });
  };

  const handleViewDetail = (id: number) => {
    setSelectedId(id);
    setIsDetailOpen(true);
  };

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat("id-ID", {
      style: "currency",
      currency: "IDR",
      minimumFractionDigits: 0,
    }).format(amount);
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
          <h2 className="text-3xl font-bold tracking-tight">
            Finance
          </h2>
          <p className="text-slate-500">
            Kelola pencatatan pemasukan dan pengeluaran perusahaan.
          </p>
        </div>
        <div className="flex gap-2 w-full sm:w-auto flex-wrap">
          {hasPermission(PERMISSIONS.EXPORT_FINANCE) && (
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
          {hasPermission(PERMISSIONS.MANAGE_FINANCE_CATEGORY) && (
            <Button
              variant="outline"
              onClick={() => setIsCategoryOpen(true)}
              className="w-full sm:w-auto"
            >
              Kelola Kategori
            </Button>
          )}
          {hasPermission(PERMISSIONS.CREATE_FINANCE) && (
            <Button
              onClick={() => setIsCreateOpen(true)}
              className="bg-blue-600 hover:bg-blue-700 w-full sm:w-auto"
            >
              <Plus className="mr-2 h-4 w-4" /> Buat Transaksi
            </Button>
          )}
        </div>
      </div>

      <Card>
        <CardHeader>
          <div className="flex flex-col md:flex-row justify-between md:items-center gap-4">
            <CardTitle className="flex items-center gap-2 text-lg">
              <Wallet className="h-5 w-5" /> Riwayat Transaksi
            </CardTitle>

            <div className="flex gap-2 w-full md:w-auto">
              <div className="relative w-full md:w-48">
                <Filter className="absolute left-2 top-2.5 h-4 w-4 text-slate-500" />
                <select
                  className="h-10 w-full rounded-md border border-input bg-background pl-8 pr-3 text-sm ring-offset-background focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2"
                  value={typeFilter}
                  onChange={(e) => {
                    setTypeFilter(e.target.value);
                  }}
                >
                  <option value="">All Types</option>
                  <option value="INCOME">Pemasukan</option>
                  <option value="EXPENSE">Pengeluaran</option>
                </select>
              </div>
              <div className="relative w-full md:w-48">
                <Filter className="absolute left-2 top-2.5 h-4 w-4 text-slate-500" />
                <select
                  className="h-10 w-full rounded-md border border-input bg-background pl-8 pr-3 text-sm ring-offset-background focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2"
                  value={statusFilter}
                  onChange={(e) => {
                    setStatusFilter(e.target.value);
                  }}
                >
                  <option value="">All Status</option>
                  <option value="PENDING">Pending</option>
                  <option value="APPROVED">Approved</option>
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
              <div className="hidden md:block rounded-md border">
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>Dibuat Oleh</TableHead>
                      <TableHead>Kategori</TableHead>
                      <TableHead>Tipe</TableHead>
                      <TableHead>Jumlah</TableHead>
                      <TableHead>Tanggal</TableHead>
                      <TableHead>Status</TableHead>
                      <TableHead className="text-right">Aksi</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {allTransactions.map((item) => (
                      <TableRow key={item.id}>
                        <TableCell className="font-medium">
                          {item.creator_name || "-"}
                        </TableCell>
                        <TableCell>{item.category_name}</TableCell>
                        <TableCell>
                          <span className={`font-medium ${item.type === "INCOME" ? "text-green-600" : "text-red-600"}`}>
                            {item.type === "INCOME" ? "Pemasukan" : "Pengeluaran"}
                          </span>
                        </TableCell>
                        <TableCell className="font-bold">{formatCurrency(item.amount)}</TableCell>
                        <TableCell>{formatDateSafe(item.transaction_date, "dd MMM yyyy")}</TableCell>
                        <TableCell><StatusBadge status={item.status} /></TableCell>
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
                    {allTransactions.length === 0 && (
                      <TableRow>
                        <TableCell colSpan={7} className="text-center py-8 text-slate-500">
                          Tidak ada data transaksi.
                        </TableCell>
                      </TableRow>
                    )}
                  </TableBody>
                </Table>
              </div>

              <div className="md:hidden space-y-3">
                {allTransactions.map((item) => (
                  <Card key={item.id} className="p-4">
                    <div className="flex items-start justify-between">
                      <div className="space-y-2 flex-1">
                        <div className="flex justify-between items-start gap-2">
                          <div>
                            <h4 className="font-semibold line-clamp-1">
                              {item.category_name}
                            </h4>
                            <div className="flex items-center text-xs text-slate-500 mt-1">
                              <Calendar className="mr-1 h-3 w-3" />
                              {formatDateSafe(item.transaction_date, "dd MMM yyyy")}
                            </div>
                          </div>
                          <StatusBadge status={item.status} />
                        </div>

                        <div className="space-y-1 mt-2">
                          <div className="flex items-center justify-between text-sm">
                            <span className="text-slate-500">Dibuat:</span>
                            <span className="font-medium">{item.creator_name || "-"}</span>
                          </div>
                          <div className="flex items-center justify-between text-sm">
                            <span className="text-slate-500">Tipe:</span>
                            <span className={`font-medium ${item.type === "INCOME" ? "text-green-600" : "text-red-600"}`}>
                              {item.type === "INCOME" ? "Pemasukan" : "Pengeluaran"}
                            </span>
                          </div>
                          <div className="flex items-center justify-between text-sm">
                            <span className="text-slate-500">Jumlah:</span>
                            <span className="font-bold">{formatCurrency(item.amount)}</span>
                          </div>
                        </div>
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

              {hasNextPage && (
                <PaginationControls
                  meta={{ limit: 10, has_next: hasNextPage }}
                  onLoadMore={() => fetchNextPage()}
                  isLoading={isFetchingNextPage}
                />
              )}
            </>
          )}
        </CardContent>
      </Card>

      <FinanceTransactionCreateDialog
        open={isCreateOpen}
        onOpenChange={setIsCreateOpen}
      />

      <FinanceTransactionDetailDialog
        open={isDetailOpen}
        onOpenChange={setIsDetailOpen}
        transactionId={selectedId}
      />

      <FinanceCategoryManager
        open={isCategoryOpen}
        onOpenChange={setIsCategoryOpen}
      />
    </div>
  );
};
