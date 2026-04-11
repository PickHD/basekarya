import { useState } from "react";
import { format } from "date-fns";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { Loader2, ExternalLink } from "lucide-react";
import { RequisitionStatusBadge } from "./RequisitionStatusBadge";
import { PriorityBadge } from "./PriorityBadge";
import { useRequisitionAction, useSubmitRequisition, useCloseRequisition } from "../hooks/useRequisition";
import { usePermissions } from "@/hooks/usePermissions";
import { PERMISSIONS } from "@/config/permissions";
import type { JobRequisition } from "../types";
import { useNavigate } from "react-router-dom";

interface Props {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  requisition: JobRequisition | null;
}

export function RequisitionDetailDialog({ open, onOpenChange, requisition }: Props) {
  const { hasPermission } = usePermissions();
  const navigate = useNavigate();
  const [rejectionReason, setRejectionReason] = useState("");
  const [showRejectForm, setShowRejectForm] = useState(false);

  const { mutate: doAction, isPending: isActioning } = useRequisitionAction();
  const { mutate: submitRequisition, isPending: isSubmitting } = useSubmitRequisition();
  const { mutate: closeRequisition, isPending: isClosing } = useCloseRequisition();

  if (!requisition) return null;

  const canApprove = hasPermission(PERMISSIONS.APPROVAL_REQUISITION);
  const canCreate = hasPermission(PERMISSIONS.CREATE_REQUISITION);
  const canViewApplicants = hasPermission(PERMISSIONS.VIEW_APPLICANT);

  const handleApprove = () => {
    doAction(
      { id: requisition.id, payload: { action: "APPROVE" } },
      { onSuccess: () => onOpenChange(false) }
    );
  };

  const handleReject = () => {
    doAction(
      { id: requisition.id, payload: { action: "REJECT", rejection_reason: rejectionReason } },
      {
        onSuccess: () => {
          setRejectionReason("");
          setShowRejectForm(false);
          onOpenChange(false);
        },
      }
    );
  };

  const handleSubmit = () => {
    submitRequisition(requisition.id, { onSuccess: () => onOpenChange(false) });
  };

  const handleClose = () => {
    closeRequisition(requisition.id, { onSuccess: () => onOpenChange(false) });
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-xl max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle className="text-lg">{requisition.title}</DialogTitle>
        </DialogHeader>

        <div className="space-y-4 text-sm">
          {/* Badges row */}
          <div className="flex flex-wrap items-center gap-2">
            <RequisitionStatusBadge status={requisition.status} />
            <PriorityBadge priority={requisition.priority} />
            <span className="text-xs text-slate-500 bg-slate-100 px-2 py-0.5 rounded">
              {requisition.employment_type}
            </span>
          </div>

          {/* Info grid */}
          <div className="grid grid-cols-2 gap-x-6 gap-y-3 text-sm">
            <div>
              <p className="text-muted-foreground text-xs">Department</p>
              <p className="font-medium">{requisition.department_name}</p>
            </div>
            <div>
              <p className="text-muted-foreground text-xs">Headcount</p>
              <p className="font-medium">{requisition.quantity}</p>
            </div>
            <div>
              <p className="text-muted-foreground text-xs">Requester</p>
              <p className="font-medium">{requisition.requester_name || "-"}</p>
            </div>
            <div>
              <p className="text-muted-foreground text-xs">Target Fill Date</p>
              <p className="font-medium">
                {requisition.target_date
                  ? format(new Date(requisition.target_date), "dd MMM yyyy")
                  : "-"}
              </p>
            </div>
            <div>
              <p className="text-muted-foreground text-xs">Created</p>
              <p className="font-medium">
                {format(new Date(requisition.created_at), "dd MMM yyyy")}
              </p>
            </div>
            {requisition.approver_name && (
              <div>
                <p className="text-muted-foreground text-xs">Approved By</p>
                <p className="font-medium">{requisition.approver_name}</p>
              </div>
            )}
          </div>

          {/* Description */}
          {requisition.description && (
            <div>
              <p className="text-muted-foreground text-xs mb-1">Description</p>
              <p className="text-sm leading-relaxed whitespace-pre-wrap bg-slate-50 rounded p-3 border">
                {requisition.description}
              </p>
            </div>
          )}

          {/* Rejection reason */}
          {requisition.rejection_reason && (
            <div>
              <p className="text-muted-foreground text-xs mb-1 text-red-600">Rejection Reason</p>
              <p className="text-sm leading-relaxed bg-red-50 text-red-700 rounded p-3 border border-red-200">
                {requisition.rejection_reason}
              </p>
            </div>
          )}

          {/* Reject form inline */}
          {showRejectForm && (
            <div className="space-y-2">
              <p className="text-xs font-medium text-red-600">Rejection Reason *</p>
              <Textarea
                placeholder="Please provide a reason for rejection..."
                value={rejectionReason}
                onChange={(e) => setRejectionReason(e.target.value)}
                rows={3}
              />
            </div>
          )}
        </div>

        <DialogFooter className="flex flex-col sm:flex-row gap-2 pt-2">
          {/* View Applicants */}
          {canViewApplicants && requisition.status === "APPROVED" && (
            <Button
              variant="outline"
              className="w-full sm:w-auto border-blue-200 text-blue-700 hover:bg-blue-50"
              onClick={() => {
                onOpenChange(false);
                navigate(`/admin/requisitions/${requisition.id}/applicants`);
              }}
            >
              <ExternalLink className="mr-2 h-4 w-4" />
              View Applicants
            </Button>
          )}

          {/* Submit (draft only, by requester) */}
          {canCreate && requisition.status === "DRAFT" && (
            <Button
              onClick={handleSubmit}
              disabled={isSubmitting}
              className="bg-blue-600 hover:bg-blue-700 w-full sm:w-auto"
            >
              {isSubmitting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              Submit for Approval
            </Button>
          )}

          {/* Approve/Reject (pending only, by approver) */}
          {canApprove && requisition.status === "PENDING" && !showRejectForm && (
            <>
              <Button
                onClick={handleApprove}
                disabled={isActioning}
                className="bg-emerald-600 hover:bg-emerald-700 w-full sm:w-auto"
              >
                {isActioning && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                Approve
              </Button>
              <Button
                variant="outline"
                onClick={() => setShowRejectForm(true)}
                className="border-red-200 text-red-600 hover:bg-red-50 w-full sm:w-auto"
              >
                Reject
              </Button>
            </>
          )}

          {/* Confirm reject */}
          {showRejectForm && (
            <>
              <Button
                onClick={handleReject}
                disabled={isActioning || !rejectionReason.trim()}
                className="bg-red-600 hover:bg-red-700 w-full sm:w-auto"
              >
                {isActioning && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                Confirm Reject
              </Button>
              <Button
                variant="outline"
                onClick={() => { setShowRejectForm(false); setRejectionReason(""); }}
                className="w-full sm:w-auto"
              >
                Cancel
              </Button>
            </>
          )}

          {/* Close */}
          {canApprove && requisition.status === "APPROVED" && (
            <Button
              variant="outline"
              onClick={handleClose}
              disabled={isClosing}
              className="w-full sm:w-auto"
            >
              {isClosing && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              Close Requisition
            </Button>
          )}
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
