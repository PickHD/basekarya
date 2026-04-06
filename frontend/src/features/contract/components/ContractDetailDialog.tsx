import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { format } from "date-fns";
import { ContractTypeBadge } from "./ContractTypeBadge";
import type { Contract } from "../types";
import { FileText, Download, Calendar, Loader2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { useContractDetail } from "../hooks/useContract";

interface Props {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  contract: Contract | null;
}

export function ContractDetailDialog({ open, onOpenChange, contract }: Props) {
  // Fetch the full detail (includes notes + attachment_url) since the list only returns summary fields
  const { data: fullDetail, isLoading } = useContractDetail(contract?.id ?? null);

  const data = fullDetail ?? contract;

  const formatDate = (dateString?: string) => {
    if (!dateString) return "-";
    return format(new Date(dateString), "dd MMMM yyyy");
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-2xl max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle className="flex justify-between items-center pr-8">
            <span>Contract Detail</span>
            {data && <ContractTypeBadge type={data.contract_type} />}
          </DialogTitle>
        </DialogHeader>

        {isLoading ? (
          <div className="flex justify-center items-center py-12">
            <Loader2 className="h-8 w-8 animate-spin text-slate-400" />
          </div>
        ) : !data ? (
          <div className="py-10 text-center text-slate-500">Data not found.</div>
        ) : (
          <div className="space-y-6 py-4">
            <div className="bg-slate-50 p-4 rounded-lg border flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
              <div>
                <h3 className="font-bold text-lg text-slate-900">
                  {data.employee_name || "-"}
                </h3>
                <p className="text-sm text-slate-500 font-medium">
                  NIK: {data.employee_nik || "-"}
                </p>
              </div>
              <div className="text-left sm:text-right bg-white sm:bg-transparent p-2 sm:p-0 rounded border sm:border-0 w-full sm:w-auto">
                <p className="text-xs text-slate-500 mb-1">Contract Number</p>
                <p className="font-bold text-slate-800">
                  {data.contract_number || "-"}
                </p>
              </div>
            </div>

            <div className="grid md:grid-cols-2 gap-6">
              <div className="space-y-5">
                <div>
                  <span className="text-sm font-medium text-slate-500 flex items-center gap-2 mb-2">
                    <Calendar className="h-4 w-4" /> Periode
                  </span>
                  <div className="p-3 bg-white border rounded-md shadow-sm">
                    <p className="text-sm font-semibold text-slate-800">
                      {formatDate(data.start_date)}
                    </p>
                    <p className="text-xs text-slate-400 my-1 text-center font-medium">s/d</p>
                    <p className="text-sm font-semibold text-slate-800">
                      {data.contract_type === "PKWTT"
                        ? "Tak Terbatas (Permanent)"
                        : formatDate(data.end_date)}
                    </p>
                  </div>
                </div>

                <div>
                  <span className="text-sm font-medium text-slate-500 flex items-center gap-2 mb-2">
                    <FileText className="h-4 w-4" /> Notes
                  </span>
                  <div className="bg-slate-50 p-3 rounded border text-sm min-h-[80px] text-slate-700 leading-relaxed italic">
                    {data.notes || "-"}
                  </div>
                </div>
              </div>

              <div className="space-y-2">
                <span className="text-sm font-medium text-slate-500 mb-2 block">
                  Attachment Document
                </span>

                {data.attachment_url ? (
                  <div className="border rounded-lg overflow-hidden bg-slate-100 relative group h-48 flex items-center justify-center">
                    {data.attachment_url.match(/\.(jpeg|jpg|png|gif)$/i) ? (
                      <img
                        src={data.attachment_url}
                        alt="Lampiran"
                        className="max-w-full max-h-full object-contain"
                      />
                    ) : (
                      <div className="flex flex-col items-center text-slate-400">
                        <FileText className="h-12 w-12 mb-2" />
                        <span className="text-xs">Preview not available</span>
                      </div>
                    )}

                    <a
                      href={data.attachment_url}
                      target="_blank"
                      rel="noreferrer"
                      className="absolute inset-0 bg-black/40 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity backdrop-blur-[1px]"
                    >
                      <Button variant="secondary" size="sm">
                        <Download className="mr-2 h-4 w-4" /> Buka
                      </Button>
                    </a>
                  </div>
                ) : (
                  <div className="border-2 border-dashed rounded-lg h-32 flex flex-col items-center justify-center text-slate-400 bg-slate-50">
                    <FileText className="h-8 w-8 mb-2 opacity-50" />
                    <span className="text-xs">No attachment</span>
                  </div>
                )}
              </div>
            </div>
          </div>
        )}
      </DialogContent>
    </Dialog>
  );
}
