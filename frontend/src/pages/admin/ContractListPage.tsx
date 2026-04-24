import { useState } from "react";
import { Plus, Download, Search, Loader2, FileSpreadsheet } from "lucide-react";
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
import { ContractList } from "@/features/contract/components/ContractList";
import { ContractFormDialog } from "@/features/contract/components/ContractFormDialog";
import { ContractDetailDialog } from "@/features/contract/components/ContractDetailDialog";
import { useContracts, useExportContract } from "@/features/contract/hooks/useContract";
import type { Contract } from "@/features/contract/types";
import { usePermissions } from "@/hooks/usePermissions";
import { PERMISSIONS } from "@/config/permissions";
import { useDebounce } from "@/hooks/useDebounce";
import { PaginationControls } from "@/components/shared/PaginationControls";

export default function ContractListPage() {
  const { hasPermission } = usePermissions();

  const [page, setPage] = useState(1);
  const [search, setSearch] = useState("");
  const debouncedSearch = useDebounce(search, 500);
  const [typeFilter, setTypeFilter] = useState("ALL");
  const [expiringFilter, setExpiringFilter] = useState("ALL");

  const [isFormOpen, setIsFormOpen] = useState(false);
  const [isDetailOpen, setIsDetailOpen] = useState(false);
  const [selectedContract, setSelectedContract] = useState<Contract | null>(null);
  const [editingEmployeeId, setEditingEmployeeId] = useState<number | null>(null);

  const { data, isLoading } = useContracts({
    page,
    limit: 10,
    search: debouncedSearch,
    contract_type: typeFilter !== "ALL" ? typeFilter : undefined,
    expiring_within_days: expiringFilter === "30" ? 30 : undefined,
  });

  const { mutateAsync: exportContractExcel, isPending: isExporting } = useExportContract();

  const handleExport = () => {
    exportContractExcel({ expiring_within_days: expiringFilter === "ALL" ? "" : expiringFilter, contract_type: typeFilter === "ALL" ? "" : typeFilter, search: debouncedSearch });
  };

  const handleCreate = () => {
    setEditingEmployeeId(null);
    setIsFormOpen(true);
  };

  const handleEdit = (contract: Contract) => {
    setEditingEmployeeId(contract.employee_id);
    setIsFormOpen(true);
  };

  const handleFormClose = (open: boolean) => {
    setIsFormOpen(open);
    if (!open) setEditingEmployeeId(null);
  };

  const handleView = (contract: Contract) => {
    setSelectedContract(contract);
    setIsDetailOpen(true);
  };

  return (
    <div className="space-y-6">
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <div>
          <h2 className="text-3xl font-bold tracking-tight">Contract Management</h2>
          <p className="text-slate-500 mt-1">
            Manage all registered contracts.
          </p>
        </div>
        <div className="flex gap-2 w-full sm:w-auto">
          {hasPermission(PERMISSIONS.EXPORT_CONTRACT) && (
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
          {hasPermission(PERMISSIONS.CREATE_CONTRACT) && (
            <Button onClick={handleCreate} className="bg-blue-600 hover:bg-blue-700">
              <Plus className="mr-2 h-4 w-4" /> Add Contract
            </Button>
          )}
        </div>
      </div>

      <Card>
        <CardHeader>
          <div className="flex flex-col sm:flex-row justify-between items-center gap-4">
            <CardTitle className="flex items-center gap-2">
              <FileSpreadsheet className="h-5 w-5" /> Contract List
            </CardTitle>
            <div className="flex flex-col sm:flex-row gap-4 w-full sm:w-auto">
              <div className="relative w-full sm:w-64">
                <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-slate-400" />
                <Input
                  placeholder="Search name / NIK..."
                  className="pl-9"
                  value={search}
                  onChange={(e) => setSearch(e.target.value)}
                />
              </div>

              <Select value={typeFilter} onValueChange={setTypeFilter}>
                <SelectTrigger className="w-full sm:w-40">
                  <SelectValue placeholder="Contract Type" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="ALL">All Types</SelectItem>
                  <SelectItem value="PKWT">PKWT</SelectItem>
                  <SelectItem value="PKWTT">PKWTT</SelectItem>
                </SelectContent>
              </Select>

              <Select value={expiringFilter} onValueChange={setExpiringFilter}>
                <SelectTrigger className="w-full sm:w-48">
                  <SelectValue placeholder="Status" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="ALL">All Status</SelectItem>
                  <SelectItem value="30">Expiring in 30 Days</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>
        </CardHeader>
        <CardContent>
          <ContractList
            data={data?.data || []}
            isLoading={isLoading}
            onView={handleView}
            onEdit={handleEdit}
          />
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
        </CardContent>
      </Card>

      <ContractFormDialog
        open={isFormOpen}
        onOpenChange={handleFormClose}
        initialEmployeeId={editingEmployeeId}
      />

      <ContractDetailDialog
        open={isDetailOpen}
        onOpenChange={setIsDetailOpen}
        contract={selectedContract}
      />
    </div>
  );
}
