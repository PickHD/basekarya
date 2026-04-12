import { useState } from "react";
import { format } from "date-fns";
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
import { Loader2, ExternalLink, Calendar, FileText, Users } from "lucide-react";
import { RequisitionStatusBadge } from "./RequisitionStatusBadge";
import { PriorityBadge } from "./PriorityBadge";
import {
  useRequisitionAction,
  useSubmitRequisition,
  useCloseRequisition,
  useRequisitionDetail,
} from "../hooks/useRequisition";
import { usePermissions } from "@/hooks/usePermissions";
import { PERMISSIONS } from "@/config/permissions";
import { useNavigate } from "react-router-dom";

interface RequisitionDetailDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  requisitionId: number | null;
}

export function RequisitionDetailDialog({
  open,
  onOpenChange,
  requisitionId,
}: RequisitionDetailDialogProps) {
  const { hasPermission } = usePermissions();
  const navigate = useNavigate();

  const { mutate: doAction, isPending: isActioning } = useRequisitionAction();
  const { mutate: submitRequisition, isPending: isSubmitting } =
    useSubmitRequisition();
  const { mutate: closeRequisition, isPending: isClosing } =
    useCloseRequisition();

  const { data: requisition, isLoading } = useRequisitionDetail(requisitionId);

  const [actionType, setActionType] = useState<
    "APPROVE" | "REJECT" | "CLOSE" | "SUBMIT" | null
  >(null);
  const [rejectionReason, setRejectionReason] = useState("");
  const [isConfirmOpen, setIsConfirmOpen] = useState(false);

  const canApprove = hasPermission(PERMISSIONS.APPROVAL_REQUISITION);
  const canCreate = hasPermission(PERMISSIONS.CREATE_REQUISITION);
  const canViewApplicants = hasPermission(PERMISSIONS.VIEW_APPLICANT);

  const isPendingGlobal = isActioning || isSubmitting || isClosing;

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

  const handleInitiateAction = (
    type: "APPROVE" | "REJECT" | "CLOSE" | "SUBMIT",
  ) => {
    setActionType(type);
    setRejectionReason("");
    setIsConfirmOpen(true);
  };

  const handleConfirmAction = () => {
    if (!requisition || !actionType) return;

    if (actionType === "REJECT") {
      if (!rejectionReason.trim()) return;
      doAction(
        {
          id: requisition.id,
          payload: { action: "REJECT", rejection_reason: rejectionReason },
        },
        {
          onSuccess: () => {
            setIsConfirmOpen(false);
            onOpenChange(false);
          },
        },
      );
    } else if (actionType === "APPROVE") {
      doAction(
        { id: requisition.id, payload: { action: "APPROVE" } },
        {
          onSuccess: () => {
            setIsConfirmOpen(false);
            onOpenChange(false);
          },
        },
      );
    } else if (actionType === "SUBMIT") {
      submitRequisition(requisition.id, {
        onSuccess: () => {
          setIsConfirmOpen(false);
          onOpenChange(false);
        },
      });
    } else if (actionType === "CLOSE") {
      closeRequisition(requisition.id, {
        onSuccess: () => {
          setIsConfirmOpen(false);
          onOpenChange(false);
        },
      });
    }
  };

  return (
    <>
      <Dialog open={open} onOpenChange={handleOpenChangeWrapper}>
        <DialogContent className="sm:max-w-2xl max-h-[90vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle className="flex justify-between items-center pr-8">
              <span>Detail Permintaan Lowongan</span>
              {requisition && <RequisitionStatusBadge status={requisition.status} />}
            </DialogTitle>
            <DialogDescription>
              ID Request: #{requisition?.id || requisitionId || "-"}
            </DialogDescription>
          </DialogHeader>

          {isLoading ? (
            <div className="flex justify-center py-10">
              <Loader2 className="h-8 w-8 animate-spin text-blue-600" />
            </div>
          ) : requisition ? (
            <div className="space-y-6 py-4">
              <div className="bg-slate-50 p-4 rounded-lg border flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
                <div>
                  <h3 className="font-bold text-lg text-slate-900">
                    {requisition.title}
                  </h3>
                  <div className="flex items-center gap-2 mt-2">
                    <p className="text-sm text-slate-500 font-medium">
                      {requisition.department_name}
                    </p>
                    <span className="text-slate-300">|</span>
                    <p className="text-xs text-slate-400">
                      {requisition.employment_type}
                    </p>
                    <span className="text-slate-300">|</span>
                    <PriorityBadge priority={requisition.priority} />
                  </div>
                </div>
                <div className="text-left sm:text-right bg-white sm:bg-transparent p-2 sm:p-0 rounded border sm:border-0 w-full sm:w-auto">
                  <p className="text-xs text-slate-500 mb-1">Target Pegawai</p>
                  <p className="text-xl font-bold text-blue-600">
                    {requisition.quantity} Orang
                  </p>
                </div>
              </div>

              <div className="grid md:grid-cols-2 gap-6">
                <div className="space-y-5">
                  <div>
                    <span className="text-sm font-medium text-slate-500 flex items-center gap-2 mb-2">
                      <Calendar className="h-4 w-4" /> Tanggal Penting
                    </span>
                    <div className="p-3 bg-white border rounded-md shadow-sm">
                      <div className="grid grid-cols-2 gap-2 text-sm">
                        <div>
                          <p className="text-xs text-slate-400 font-medium mb-1">
                            Target Fill Date
                          </p>
                          <p className="font-semibold text-slate-800">
                            {requisition.target_date
                              ? format(
                                  new Date(requisition.target_date),
                                  "dd MMM yyyy",
                                )
                              : "-"}
                          </p>
                        </div>
                        <div>
                          <p className="text-xs text-slate-400 font-medium mb-1">
                            Tanggal Dibuat
                          </p>
                          <p className="font-semibold text-slate-800">
                            {format(
                              new Date(requisition.created_at),
                              "dd MMM yyyy",
                            )}
                          </p>
                        </div>
                      </div>
                    </div>
                  </div>

                  <div>
                    <span className="text-sm font-medium text-slate-500 flex items-center gap-2 mb-2">
                      <Users className="h-4 w-4" /> Informasi Personil
                    </span>
                    <div className="p-3 bg-white border rounded-md shadow-sm">
                      <div className="grid grid-cols-2 gap-2 text-sm">
                        <div>
                          <p className="text-xs text-slate-400 font-medium mb-1">
                            Pemohon
                          </p>
                          <p className="font-semibold text-slate-800">
                            {requisition.requester_name || "-"}
                          </p>
                        </div>
                        <div>
                          <p className="text-xs text-slate-400 font-medium mb-1">
                            Disetujui Oleh
                          </p>
                          <p className="font-semibold text-slate-800">
                            {requisition.approver_name || "-"}
                          </p>
                        </div>
                      </div>
                    </div>
                  </div>

                  {requisition.rejection_reason && (
                    <div className="bg-red-50 p-3 rounded border border-red-200 animate-in slide-in-from-bottom-2">
                      <span className="text-sm font-bold text-red-700 block mb-1">
                        Alasan Penolakan:
                      </span>
                      <p className="text-sm text-red-600">
                        {requisition.rejection_reason}
                      </p>
                    </div>
                  )}
                </div>

                <div className="space-y-5">
                  <div>
                    <span className="text-sm font-medium text-slate-500 flex items-center gap-2 mb-2">
                      <FileText className="h-4 w-4" /> Deskripsi Pekerjaan
                    </span>
                    <div className="bg-slate-50 p-3 rounded border text-sm min-h-[220px] text-slate-700 leading-relaxed whitespace-pre-wrap">
                      {requisition.description || "Tidak ada deskripsi"}
                    </div>
                  </div>
                </div>
              </div>

              <div className="flex flex-col-reverse sm:flex-row gap-2 justify-end pt-4 mt-4 border-t">
                {canViewApplicants && requisition.status === "APPROVED" && (
                  <Button
                    variant="outline"
                    className="w-full sm:w-auto border-blue-200 text-blue-700 hover:bg-blue-50 sm:mr-auto"
                    onClick={() => {
                      onOpenChange(false);
                      navigate(
                        `/admin/requisitions/${requisition.id}/applicants`,
                      );
                    }}
                  >
                    <ExternalLink className="mr-2 h-4 w-4" />
                    View Applicants
                  </Button>
                )}

                {canApprove && requisition.status === "APPROVED" && (
                  <Button
                    variant="outline"
                    onClick={() => handleInitiateAction("CLOSE")}
                    disabled={isPendingGlobal}
                    className="w-full sm:w-auto"
                  >
                    Tutup Permintaan
                  </Button>
                )}

                {canCreate && requisition.status === "DRAFT" && (
                  <Button
                    className="bg-blue-600 hover:bg-blue-700 w-full sm:w-auto"
                    onClick={() => handleInitiateAction("SUBMIT")}
                    disabled={isPendingGlobal}
                  >
                    Submit for Approval
                  </Button>
                )}

                {canApprove && requisition.status === "PENDING" && (
                  <>
                    <Button
                      variant="destructive"
                      onClick={() => handleInitiateAction("REJECT")}
                      disabled={isPendingGlobal}
                      className="w-full sm:w-auto"
                    >
                      Tolak Permintaan
                    </Button>
                    <Button
                      className="bg-green-600 hover:bg-green-700 w-full sm:w-auto"
                      onClick={() => handleInitiateAction("APPROVE")}
                      disabled={isPendingGlobal}
                    >
                      Setujui Permintaan
                    </Button>
                  </>
                )}
              </div>
            </div>
          ) : (
            <div className="py-10 text-center text-slate-500">
              Data not found.
            </div>
          )}
        </DialogContent>
      </Dialog>

      <AlertDialog open={isConfirmOpen} onOpenChange={setIsConfirmOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>
              {actionType === "APPROVE"
                ? "Setujui Permintaan?"
                : actionType === "REJECT"
                ? "Tolak Permintaan?"
                : actionType === "SUBMIT"
                ? "Submit Permintaan?"
                : "Tutup Permintaan?"}
            </AlertDialogTitle>
            <AlertDialogDescription>
              {actionType === "APPROVE"
                ? "Apakah Anda yakin ingin menyetujui permintaan lowongan ini? Lowongan akan dipublikasikan dan dapat diproses lebih lanjut."
                : actionType === "REJECT"
                ? "Harap berikan alasan penolakan yang jelas agar pemohon mengerti."
                : actionType === "SUBMIT"
                ? "Apakah Anda yakin ingin mengajukan permintaan ini untuk persetujuan?"
                : "Apakah Anda yakin ingin menutup permintaan ini? Anda tidak dapat membuka kembali nantinya."}
            </AlertDialogDescription>
          </AlertDialogHeader>

          {actionType === "REJECT" && (
            <div className="py-2 space-y-2">
              <Label htmlFor="reason" className="text-sm font-medium">
                Alasan Penolakan <span className="text-red-500">*</span>
              </Label>
              <Textarea
                id="reason"
                placeholder="Contoh: Budget rekrutmen tidak tersedia tahun ini"
                value={rejectionReason}
                onChange={(e) => setRejectionReason(e.target.value)}
                className="resize-none focus-visible:ring-red-500"
              />
            </div>
          )}

          <AlertDialogFooter>
            <AlertDialogCancel disabled={isPendingGlobal}>
              Batal
            </AlertDialogCancel>
            <AlertDialogAction
              onClick={(e) => {
                e.preventDefault();
                handleConfirmAction();
              }}
              disabled={
                isPendingGlobal ||
                (actionType === "REJECT" && !rejectionReason.trim())
              }
              className={
                actionType === "REJECT"
                  ? "bg-red-600 hover:bg-red-700"
                  : actionType === "APPROVE"
                  ? "bg-green-600 hover:bg-green-700"
                  : "bg-blue-600 hover:bg-blue-700"
              }
            >
              {isPendingGlobal && (
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              )}
              {actionType === "APPROVE"
                ? "Ya, Setujui"
                : actionType === "REJECT"
                ? "Tolak Permintaan"
                : actionType === "SUBMIT"
                ? "Ya, Submit"
                : "Ya, Tutup"}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </>
  );
}
