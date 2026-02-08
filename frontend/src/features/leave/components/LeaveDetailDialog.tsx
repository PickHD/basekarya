import { useState } from "react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from "@/components/ui/dialog";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Loader2, Calendar, FileText, Download } from "lucide-react";
import { useProfile } from "@/features/user/hooks/useProfile";
import {
  useLeaveDetail,
  useLeaveAction,
} from "@/features/leave/hooks/useLeave";
import { format } from "date-fns";
import { LeaveStatusBadge } from "@/features/leave/components/LeaveStatusBadge";

interface LeaveDetailDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  leaveId: number | null;
}

export function LeaveDetailDialog({
  open,
  onOpenChange,
  leaveId,
}: LeaveDetailDialogProps) {
  const { data: userProfile } = useProfile();
  const { data, isLoading } = useLeaveDetail(leaveId?.toString() || "");

  const { mutate: actionMutate, isPending } = useLeaveAction();

  const [actionType, setActionType] = useState<"APPROVE" | "REJECT" | null>(
    null,
  );
  const [rejectionReason, setRejectionReason] = useState("");
  const [isConfirmOpen, setIsConfirmOpen] = useState(false);

  const isSuperAdmin = userProfile?.role === "SUPERADMIN";
  const isPendingStatus = data?.status === "PENDING";

  const handleOpenChangeWrapper = (isOpen: boolean) => {
    onOpenChange(isOpen);

    if (!isOpen) {
      setTimeout(() => {
        setRejectionReason("");
        setActionType(null);
        setIsConfirmOpen(false);
      }, 300);
    }
  };

  const handleInitiateAction = (type: "APPROVE" | "REJECT") => {
    setActionType(type);
    setRejectionReason("");
    setIsConfirmOpen(true);
  };

  const handleConfirmAction = () => {
    if (!data || !actionType) return;

    if (actionType === "REJECT" && !rejectionReason.trim()) {
      return;
    }

    actionMutate(
      { id: data.id, action: actionType, rejection_reason: rejectionReason },
      {
        onSuccess: () => {
          setIsConfirmOpen(false);
          onOpenChange(false);
        },
      },
    );
  };

  const formatDate = (dateString: string) => {
    if (!dateString) return "-";
    return format(new Date(dateString), "dd MMMM yyyy");
  };

  return (
    <>
      <Dialog open={open} onOpenChange={handleOpenChangeWrapper}>
        <DialogContent className="sm:max-w-2xl max-h-[90vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle className="flex justify-between items-center pr-8">
              <span>Detail Permintaan Cuti</span>
              {data && <LeaveStatusBadge status={data.status} />}
            </DialogTitle>
            <DialogDescription>ID Request: #{leaveId}</DialogDescription>
          </DialogHeader>

          {isLoading ? (
            <div className="flex justify-center py-10">
              <Loader2 className="h-8 w-8 animate-spin text-blue-600" />
            </div>
          ) : data ? (
            <div className="space-y-6 py-4">
              <div className="bg-slate-50 p-4 rounded-lg border flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
                <div>
                  <h3 className="font-bold text-lg text-slate-900">
                    {data.leave_type?.name || "Tipe Cuti"}
                  </h3>
                  <div className="flex items-center gap-2 mt-1">
                    <p className="text-sm text-slate-500 font-medium">
                      {data.employee_name || "Karyawan"}
                    </p>
                    <span className="text-slate-300">|</span>
                    <p className="text-xs text-slate-400">
                      {data.employee_nik}
                    </p>
                  </div>
                </div>
                <div className="text-left sm:text-right bg-white sm:bg-transparent p-2 sm:p-0 rounded border sm:border-0 w-full sm:w-auto">
                  <p className="text-xs text-slate-500 mb-1">Total Durasi</p>
                  <p className="text-xl font-bold text-blue-600">
                    {data.total_days} Hari
                  </p>
                </div>
              </div>

              <div className="grid md:grid-cols-2 gap-6">
                <div className="space-y-5">
                  <div>
                    <span className="text-sm font-medium text-slate-500 flex items-center gap-2 mb-2">
                      <Calendar className="h-4 w-4" /> Periode Cuti
                    </span>
                    <div className="p-3 bg-white border rounded-md shadow-sm">
                      <p className="text-sm font-semibold text-slate-800">
                        {formatDate(data.start_date)}
                      </p>
                      <p className="text-xs text-slate-400 my-1 text-center font-medium">
                        s/d
                      </p>
                      <p className="text-sm font-semibold text-slate-800">
                        {formatDate(data.end_date)}
                      </p>
                    </div>
                  </div>

                  <div>
                    <span className="text-sm font-medium text-slate-500 flex items-center gap-2 mb-2">
                      <FileText className="h-4 w-4" /> Alasan / Keterangan
                    </span>
                    <div className="bg-slate-50 p-3 rounded border text-sm min-h-[80px] text-slate-700 leading-relaxed italic">
                      "{data.reason || "-"}"
                    </div>
                  </div>

                  {data.rejection_reason && (
                    <div className="bg-red-50 p-3 rounded border border-red-200 animate-in slide-in-from-bottom-2">
                      <span className="text-sm font-bold text-red-700 block mb-1">
                        Alasan Penolakan:
                      </span>
                      <p className="text-sm text-red-600">
                        {data.rejection_reason}
                      </p>
                    </div>
                  )}
                </div>

                <div className="space-y-2">
                  <span className="text-sm font-medium text-slate-500 mb-2 block">
                    Lampiran Dokumen
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
                          <span className="text-xs">
                            Preview tidak tersedia
                          </span>
                        </div>
                      )}

                      <a
                        href={data.attachment_url}
                        target="_blank"
                        rel="noreferrer"
                        className="absolute inset-0 bg-black/40 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity backdrop-blur-[1px]"
                      >
                        <Button variant="secondary" size="sm">
                          <Download className="mr-2 h-4 w-4" /> Buka / Download
                        </Button>
                      </a>
                    </div>
                  ) : (
                    <div className="border-2 border-dashed rounded-lg h-32 flex flex-col items-center justify-center text-slate-400 bg-slate-50">
                      <FileText className="h-8 w-8 mb-2 opacity-50" />
                      <span className="text-xs">Tidak ada lampiran</span>
                    </div>
                  )}
                </div>
              </div>
            </div>
          ) : (
            <div className="py-10 text-center text-slate-500">
              Data not found.
            </div>
          )}

          {isSuperAdmin && isPendingStatus && (
            <div className="flex flex-col-reverse sm:flex-row gap-2 justify-end pt-4 mt-4 border-t">
              <Button
                variant="destructive"
                onClick={() => handleInitiateAction("REJECT")}
                disabled={isPending}
                className="w-full sm:w-auto"
              >
                Tolak Permintaan
              </Button>
              <Button
                className="bg-green-600 hover:bg-green-700 w-full sm:w-auto"
                onClick={() => handleInitiateAction("APPROVE")}
                disabled={isPending}
              >
                {isPending ? (
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                ) : null}
                Setujui Permintaan
              </Button>
            </div>
          )}
        </DialogContent>
      </Dialog>

      <AlertDialog open={isConfirmOpen} onOpenChange={setIsConfirmOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>
              {actionType === "APPROVE"
                ? "Setujui Permintaan Cuti?"
                : "Tolak Permintaan Cuti?"}
            </AlertDialogTitle>
            <AlertDialogDescription>
              {actionType === "APPROVE"
                ? "Apakah Anda yakin ingin menyetujui cuti ini? Saldo cuti karyawan akan otomatis terpotong dan absensi akan diperbarui."
                : "Harap berikan alasan penolakan yang jelas agar karyawan mengerti."}
            </AlertDialogDescription>
          </AlertDialogHeader>

          {actionType === "REJECT" && (
            <div className="py-2 space-y-2">
              <Label htmlFor="reason" className="text-sm font-medium">
                Alasan Penolakan <span className="text-red-500">*</span>
              </Label>
              <Textarea
                id="reason"
                placeholder="Contoh: Kuota cuti tahunan sudah habis / Jadwal terlalu padat"
                value={rejectionReason}
                onChange={(e) => setRejectionReason(e.target.value)}
                className="resize-none focus-visible:ring-red-500"
              />
            </div>
          )}

          <AlertDialogFooter>
            <AlertDialogCancel disabled={isPending}>Batal</AlertDialogCancel>
            <AlertDialogAction
              onClick={(e) => {
                e.preventDefault();
                handleConfirmAction();
              }}
              disabled={
                isPending ||
                (actionType === "REJECT" && !rejectionReason.trim())
              }
              className={
                actionType === "REJECT"
                  ? "bg-red-600 hover:bg-red-700"
                  : "bg-green-600 hover:bg-green-700"
              }
            >
              {isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              {actionType === "APPROVE" ? "Ya, Setujui" : "Tolak Permintaan"}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </>
  );
}
