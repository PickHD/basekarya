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
  Calendar,
  CreditCard,
} from "lucide-react";
import { StatusBadge } from "./StatusBadge";
import { PaginationControls } from "@/components/shared/PaginationControls";
import { useReimbursements } from "../hooks/useReimbursement";
import { format, isValid } from "date-fns";
import { ReimbursementDetailDialog } from "./ReimbursementDetailDialog";
import { ReimbursementFormDialog } from "./ReimbursementCreateDialog";
import { useProfile } from "@/features/user/hooks/useProfile";

export const ReimbursementList = () => {
  const { data: user } = useProfile();
  const [page, setPage] = useState(1);
  const [statusFilter, setStatusFilter] = useState("");

  const { data, isLoading } = useReimbursements({
    page: page,
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
          <h2 className="text-2xl sm:text-3xl font-bold tracking-tight">
            Reimbursements
          </h2>
          <p className="text-sm sm:text-base text-slate-500">
            Manage financial claims and approvals.
          </p>
        </div>
        {user?.role !== "SUPERADMIN" && (
          <Button
            onClick={() => setIsCreateOpen(true)}
            className="bg-blue-600 hover:bg-blue-700 w-full sm:w-auto"
          >
            <Plus className="mr-2 h-4 w-4" /> New Request
          </Button>
        )}
      </div>

      <Card>
        <CardHeader>
          <div className="flex flex-col md:flex-row justify-between md:items-center gap-4">
            <CardTitle className="flex items-center gap-2 text-lg">
              <Receipt className="h-5 w-5" /> Request List
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
              <div className="grid grid-cols-1 gap-4 md:hidden">
                {data?.data.map((item) => (
                  <div
                    key={item.id}
                    className="flex flex-col rounded-lg border bg-card p-4 shadow-sm space-y-3"
                  >
                    <div className="flex justify-between items-start gap-2">
                      <div>
                        <h4 className="font-semibold line-clamp-1">
                          {item.title}
                        </h4>
                        <div className="flex items-center text-xs text-slate-500 mt-1">
                          <Calendar className="mr-1 h-3 w-3" />
                          {formatDateSafe(item.date_of_expense, "dd MMM yyyy")}
                        </div>
                      </div>
                      <StatusBadge status={item.status} />
                    </div>

                    <div className="flex items-center text-slate-900 font-bold text-lg">
                      <CreditCard className="mr-2 h-4 w-4 text-slate-400" />
                      {formatCurrency(item.amount)}
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
                      <TableHead>Date</TableHead>
                      <TableHead>Title</TableHead>
                      <TableHead>Amount</TableHead>
                      <TableHead>Status</TableHead>
                      <TableHead className="text-right">Action</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {data?.data.map((item) => (
                      <TableRow key={item.id}>
                        <TableCell className="font-medium">
                          {formatDateSafe(item.date_of_expense, "dd MMM yyyy")}
                          <div className="text-xs text-slate-400 font-normal">
                            {formatDateSafe(item.date_of_expense, "EEEE")}
                          </div>
                        </TableCell>
                        <TableCell className="max-w-[200px] truncate">
                          {item.title}
                        </TableCell>
                        <TableCell className="font-bold">
                          {formatCurrency(item.amount)}
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
                    <p>No reimbursement requests found.</p>
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

      <ReimbursementFormDialog
        open={isCreateOpen}
        onOpenChange={setIsCreateOpen}
      />

      <ReimbursementDetailDialog
        open={isDetailOpen}
        onOpenChange={setIsDetailOpen}
        reimbursementId={selectedId}
      />
    </div>
  );
};
