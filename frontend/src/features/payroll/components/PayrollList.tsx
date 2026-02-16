import { useState } from "react";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import {
  Loader2,
  Plus,
  Search,
  Download,
  FileText,
  Calculator,
  Eye,
} from "lucide-react";
import { format } from "date-fns";
import { PaginationControls } from "@/components/shared/PaginationControls";
import { PayrollGenerateDialog } from "@/features/payroll/components/PayrollGenerateDialog";
import {
  usePayrolls,
  downloadPayslip,
} from "@/features/payroll/hooks/usePayroll";
import { useDebounce } from "@/hooks/useDebounce";
import { PayrollDetailDialog } from "./PayrollDetailDialog";

export default function PayrollList() {
  const [page, setPage] = useState(1);
  const [month, setMonth] = useState<string>(String(new Date().getMonth() + 1));
  const [year, setYear] = useState<string>(String(new Date().getFullYear()));
  const [search, setSearch] = useState("");

  const [isGenerateOpen, setIsGenerateOpen] = useState(false);
  const [isDownloading, setIsDownloading] = useState<number | null>(null);

  const debouncedSearch = useDebounce(search, 500);

  const { data, isLoading, isError } = usePayrolls({
    page,
    limit: 10,
    month: Number(month),
    year: Number(year),
    search: debouncedSearch,
  });

  const [selectedPayrollId, setSelectedPayrollId] = useState<number | null>(
    null,
  );
  const [isDetailOpen, setIsDetailOpen] = useState(false);

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat("id-ID", {
      style: "currency",
      currency: "IDR",
      minimumFractionDigits: 0,
    }).format(amount);
  };

  const handleDownload = async (
    id: number,
    nik: string,
    periodDate: string,
  ) => {
    setIsDownloading(id);
    const dateStr = format(new Date(periodDate), "MMMyyyy");
    const filename = `Payslip-${nik}-${dateStr}.pdf`;

    await downloadPayslip(id, filename);
    setIsDownloading(null);
  };

  return (
    <div className="space-y-6">
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <div>
          <h2 className="text-3xl font-bold tracking-tight text-slate-900">
            Payroll
          </h2>
          <p className="text-slate-500">
            Manage employee salaries and generate payslips.
          </p>
        </div>
        <Button
          onClick={() => setIsGenerateOpen(true)}
          className="bg-blue-600 hover:bg-blue-700"
        >
          <Plus className="mr-2 h-4 w-4" /> Generate Payroll
        </Button>
      </div>

      <Card>
        <CardHeader>
          <div className="flex flex-col md:flex-row justify-between gap-4">
            <CardTitle className="flex items-center gap-2">
              <Calculator className="h-5 w-5 text-black-600" /> Payslip History
            </CardTitle>

            <div className="flex flex-wrap gap-2">
              <div className="relative w-full md:w-60">
                <Search className="absolute left-2 top-2.5 h-4 w-4 text-slate-400" />
                <Input
                  placeholder="Search employee..."
                  className="pl-8"
                  value={search}
                  onChange={(e) => {
                    setSearch(e.target.value);
                    setPage(1);
                  }}
                />
              </div>

              <Select
                value={month}
                onValueChange={(val) => {
                  setMonth(val);
                  setPage(1);
                }}
              >
                <SelectTrigger className="w-[130px]">
                  <SelectValue placeholder="Month" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="0">All Month</SelectItem>
                  {Array.from({ length: 12 }, (_, i) => (
                    <SelectItem key={i + 1} value={String(i + 1)}>
                      {new Date(0, i).toLocaleString("id-ID", {
                        month: "short",
                      })}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>

              <Select
                value={year}
                onValueChange={(val) => {
                  setYear(val);
                  setPage(1);
                }}
              >
                <SelectTrigger className="w-[100px]">
                  <SelectValue placeholder="Year" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="0">All Year</SelectItem>
                  <SelectItem value="2025">2025</SelectItem>
                  <SelectItem value="2026">2026</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>
        </CardHeader>

        <CardContent>
          <div className="rounded-md border">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Period</TableHead>
                  <TableHead>Employee</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead className="text-right">Net Salary</TableHead>
                  <TableHead className="text-right">Action</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {isLoading ? (
                  <TableRow>
                    <TableCell colSpan={5} className="h-24 text-center">
                      <div className="flex justify-center items-center gap-2 text-slate-500">
                        <Loader2 className="h-5 w-5 animate-spin" /> Loading
                        data...
                      </div>
                    </TableCell>
                  </TableRow>
                ) : isError ? (
                  <TableRow>
                    <TableCell
                      colSpan={5}
                      className="h-24 text-center text-red-500"
                    >
                      Failed to load data.
                    </TableCell>
                  </TableRow>
                ) : data?.data.length === 0 ? (
                  <TableRow>
                    <TableCell
                      colSpan={5}
                      className="h-32 text-center text-slate-500"
                    >
                      <div className="flex flex-col items-center justify-center gap-2">
                        <FileText className="h-8 w-8 text-slate-300" />
                        <p>No payroll data found for this period.</p>
                      </div>
                    </TableCell>
                  </TableRow>
                ) : (
                  data?.data.map((item) => (
                    <TableRow key={item.id}>
                      <TableCell className="font-medium text-slate-600">
                        {format(new Date(item.period_date), "MMMM yyyy")}
                      </TableCell>
                      <TableCell>
                        <div className="font-semibold text-slate-900">
                          {item.employee_name}
                        </div>
                        <div className="text-xs text-slate-500">
                          {item.employee_nik}
                        </div>
                      </TableCell>
                      <TableCell>
                        <Badge
                          variant={
                            item.status === "PAID" ? "default" : "secondary"
                          }
                          className={
                            item.status === "PAID"
                              ? "bg-green-100 text-green-700 hover:bg-green-100"
                              : "bg-yellow-100 text-yellow-700 hover:bg-yellow-100"
                          }
                        >
                          {item.status}
                        </Badge>
                      </TableCell>
                      <TableCell className="text-right font-bold text-slate-700">
                        {formatCurrency(item.net_salary)}
                      </TableCell>
                      <TableCell className="text-right">
                        <Button
                          variant="ghost"
                          size="sm"
                          className="h-8 w-8 p-0"
                          onClick={() => {
                            setSelectedPayrollId(item.id);
                            setIsDetailOpen(true);
                          }}
                        >
                          <Eye className="h-4 w-4 text-slate-500 hover:text-blue-600" />
                        </Button>
                        <Button
                          variant="ghost"
                          size="sm"
                          className="h-8 w-8 p-0"
                          onClick={() =>
                            handleDownload(
                              item.id,
                              item.employee_nik,
                              item.period_date,
                            )
                          }
                          disabled={isDownloading === item.id}
                        >
                          {isDownloading === item.id ? (
                            <Loader2 className="h-4 w-4 animate-spin text-blue-600" />
                          ) : (
                            <Download className="h-4 w-4 text-slate-500 hover:text-blue-600" />
                          )}
                          <span className="sr-only">Download</span>
                        </Button>
                      </TableCell>
                    </TableRow>
                  ))
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
        </CardContent>
      </Card>

      <PayrollDetailDialog
        open={isDetailOpen}
        onOpenChange={setIsDetailOpen}
        payrollId={selectedPayrollId}
      />

      <PayrollGenerateDialog
        open={isGenerateOpen}
        onOpenChange={setIsGenerateOpen}
      />
    </div>
  );
}
