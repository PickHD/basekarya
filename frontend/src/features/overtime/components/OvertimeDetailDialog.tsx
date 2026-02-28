import { useState } from "react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
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
import { Loader2 } from "lucide-react";
import { StatusBadge } from "./StatusBadge";
import { useProfile } from "@/features/user/hooks/useProfile";
import { useOvertime, useOvertimeAction } from "@/features/overtime/hooks/useOvertime";
import { format, isValid } from "date-fns";

interface OvertimeDetailDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  overtimeId: number | null;
}

export function OvertimeDetailDialog({
  open,
  onOpenChange,
  overtimeId,
}: OvertimeDetailDialogProps) {
  const { data: userProfile } = useProfile();
  const { data, isLoading } = useOvertime(overtimeId?.toString() || "");
  const { mutate: actionMutate, isPending } = useOvertimeAction();

  const [actionType, setActionType] = useState<"APPROVE" | "REJECT" | null>(null);
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
      }
    );
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
    
    if (hours > 0 && mins > 0) return `${hours} jam ${mins} menit`;
    if (hours > 0) return `${hours} jam`;
    return `${mins} menit`;
  };

  return (
    <>
      <Dialog open={open} onOpenChange={handleOpenChangeWrapper}>
        <DialogContent className="sm:max-w-2xl max-h-[90vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle className="flex justify-between items-center pr-8">
              <span>Detail Lembur</span>
              {data && <StatusBadge status={data.status} />}
            </DialogTitle>
            <DialogDescription>ID Request: #{overtimeId}</DialogDescription>
          </DialogHeader>

          {isLoading ? (
            <div className="flex justify-center py-10">
              <Loader2 className="h-8 w-8 animate-spin text-blue-600" />
            </div>
          ) : data ? (
            <div className="grid gap-6 py-4">
              <div className="bg-slate-50 p-4 rounded-lg border flex justify-between items-start">
                <div>
                  <h3 className="font-bold text-lg text-slate-900">
                    {data.employee_name || "Karyawan"}
                  </h3>
                  <p className="text-sm text-slate-500">NIK: {data.employee_nik || "-"}</p>
                </div>
                <div className="text-right">
                  <p className="text-xs text-slate-500 mb-1">Durasi</p>
                  <p className="text-xl font-bold text-blue-600">
                    {formatDuration(data.duration_minutes)}
                  </p>
                </div>
              </div>

              <div className="grid md:grid-cols-2 gap-6">
                <div className="space-y-4">
                  <div>
                    <span className="text-sm font-medium text-slate-500">Tanggal Lembur</span>
                    <p className="text-sm font-medium text-slate-900">{formatDateSafe(data.date, "EEEE, dd MMMM yyyy")}</p>
                  </div>
                  <div>
                    <span className="text-sm font-medium text-slate-500">Waktu</span>
                    <p className="text-sm font-medium text-slate-900">{data.start_time} - {data.end_time}</p>
                  </div>
                  <div>
                    <span className="text-sm font-medium text-slate-500">Pekerjaan / Alasan</span>
                    <p className="text-sm text-slate-800 whitespace-pre-wrap">{data.reason || "-"}</p>
                  </div>
                  {data.rejection_reason && (
                    <div className="bg-red-50 p-3 rounded border border-red-200">
                      <span className="text-sm font-bold text-red-700 block">
                        Alasan Penolakan:
                      </span>
                      <p className="text-sm text-red-600">{data.rejection_reason}</p>
                    </div>
                  )}
                </div>
              </div>
            </div>
          ) : (
            <div className="py-10 text-center text-slate-500">Data not found.</div>
          )}

          {isSuperAdmin && isPendingStatus && (
            <DialogFooter className="gap-2 sm:gap-0 mt-4">
              <Button
                variant="destructive"
                onClick={() => handleInitiateAction("REJECT")}
                disabled={isPending}
              >
                Tolak Permintaan
              </Button>
              <Button
                className="bg-green-600 hover:bg-green-700"
                onClick={() => handleInitiateAction("APPROVE")}
                disabled={isPending}
              >
                Setujui Permintaan
              </Button>
            </DialogFooter>
          )}
        </DialogContent>
      </Dialog>

      <AlertDialog open={isConfirmOpen} onOpenChange={setIsConfirmOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>
              {actionType === "APPROVE" ? "Setujui Lembur?" : "Tolak Lembur?"}
            </AlertDialogTitle>
            <AlertDialogDescription>
              {actionType === "APPROVE"
                ? "Apakah Anda yakin ingin menyetujui lembur ini? Jam lembur akan masuk dalam perhitungan tunjangan."
                : "Harap berikan alasan penolakan agar karyawan mengetahui sebabnya."}
            </AlertDialogDescription>
          </AlertDialogHeader>

          {actionType === "REJECT" && (
            <div className="py-2 space-y-2">
              <Label htmlFor="reason" className="text-sm font-medium">
                Alasan Penolakan <span className="text-red-500">*</span>
              </Label>
              <Textarea
                id="reason"
                placeholder="Contoh: Pekerjaan bisa diselesaikan besok hari"
                value={rejectionReason}
                onChange={(e: any) => setRejectionReason(e.target.value)}
                className="resize-none"
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
              disabled={isPending || (actionType === "REJECT" && !rejectionReason.trim())}
              className={
                actionType === "REJECT"
                  ? "bg-red-600 hover:bg-red-700"
                  : "bg-green-600 hover:bg-green-700"
              }
            >
              {isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              {actionType === "APPROVE" ? "Ya, Setujui" : "Tolak Pengajuan"}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </>
  );
}
