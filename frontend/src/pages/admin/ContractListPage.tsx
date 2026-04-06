import { useState } from "react";
import { Plus, Download, Search } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { ContractList } from "@/features/contract/components/ContractList";
import { ContractFormDialog } from "@/features/contract/components/ContractFormDialog";
import { ContractDetailDialog } from "@/features/contract/components/ContractDetailDialog";
import { useContracts } from "@/features/contract/hooks/useContract";
import type { Contract } from "@/features/contract/types";
import { usePermissions } from "@/hooks/usePermissions";
import { PERMISSIONS } from "@/config/permissions";
import { api } from "@/lib/axios";
import { toast } from "sonner";
import { useDebounce } from "@/hooks/useDebounce";

export default function ContractListPage() {
  const { hasPermission } = usePermissions();

  const [page] = useState(1);
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

  const handleExport = async () => {
    try {
      const resp = await api.get("/contracts/export", {
        responseType: "blob",
        params: {
          search: debouncedSearch,
          contract_type: typeFilter !== "ALL" ? typeFilter : undefined,
          expiring_within_days: expiringFilter === "30" ? 30 : undefined,
        },
      });
      const url = window.URL.createObjectURL(new Blob([resp.data]));
      const link = document.createElement("a");
      link.href = url;
      link.setAttribute("download", `contracts-${new Date().getTime()}.xlsx`);
      document.body.appendChild(link);
      link.click();
      link.remove();
    } catch (error) {
      toast.error("Failed to download contract data");
    }
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
          <h2 className="text-2xl font-bold tracking-tight text-slate-800">
            Contract Management
          </h2>
          <p className="text-sm text-slate-500 mt-1">
            Manage employee contract data PKWT and PKWTT
          </p>
        </div>
        <div className="flex gap-2 w-full sm:w-auto">
          {hasPermission(PERMISSIONS.EXPORT_CONTRACT) && (
            <Button variant="outline" onClick={handleExport} className="w-full sm:w-auto">
              <Download className="mr-2 h-4 w-4" /> Export
            </Button>
          )}
          {hasPermission(PERMISSIONS.CREATE_CONTRACT) && (
            <Button onClick={handleCreate} className="w-full sm:w-auto">
              <Plus className="mr-2 h-4 w-4" /> Add Contract
            </Button>
          )}
        </div>
      </div>

      <div className="flex flex-col sm:flex-row justify-between items-center gap-4 bg-white p-4 rounded-lg border shadow-sm">
        <div className="relative w-full sm:w-72">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-slate-400" />
          <Input
            placeholder="Search name / NIK..."
            className="pl-9 bg-slate-50"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
          />
        </div>

        <div className="flex flex-col sm:flex-row gap-4 w-full sm:w-auto">
          <Select value={typeFilter} onValueChange={setTypeFilter}>
            <SelectTrigger className="w-full sm:w-40 bg-slate-50">
              <SelectValue placeholder="Contract Type" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="ALL">All Types</SelectItem>
              <SelectItem value="PKWT">PKWT</SelectItem>
              <SelectItem value="PKWTT">PKWTT</SelectItem>
            </SelectContent>
          </Select>

          <Select value={expiringFilter} onValueChange={setExpiringFilter}>
            <SelectTrigger className="w-full sm:w-48 bg-slate-50">
              <SelectValue placeholder="Status Expired" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="ALL">All Status</SelectItem>
              <SelectItem value="30">Expires in 30 Days</SelectItem>
            </SelectContent>
          </Select>
        </div>
      </div>

      <ContractList
        data={data?.data || []}
        isLoading={isLoading}
        onView={handleView}
        onEdit={handleEdit}
      />

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
