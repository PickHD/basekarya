import { format } from "date-fns";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Loader2, FileDown, Clock, User } from "lucide-react";
import { StageBadge } from "./StageBadge";
import { useApplicantDetail } from "../hooks/useApplicant";

interface Props {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  applicantId: number | null;
}

export function ApplicantDetailDialog({ open, onOpenChange, applicantId }: Props) {
  const { data: applicant, isLoading } = useApplicantDetail(applicantId);

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-lg max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>Applicant Detail</DialogTitle>
        </DialogHeader>

        {isLoading ? (
          <div className="flex justify-center py-10">
            <Loader2 className="h-5 w-5 animate-spin text-slate-400" />
          </div>
        ) : !applicant ? (
          <p className="text-center text-slate-500 py-6">No data found</p>
        ) : (
          <div className="space-y-5 text-sm">
            {/* Header */}
            <div className="flex items-start justify-between gap-3">
              <div>
                <p className="text-lg font-semibold">{applicant.full_name}</p>
                <p className="text-slate-500">{applicant.email}</p>
                {applicant.phone_number && (
                  <p className="text-slate-500">{applicant.phone_number}</p>
                )}
              </div>
              <StageBadge stage={applicant.stage} />
            </div>

            {/* Resume */}
            {applicant.resume_url && (
              <a href={applicant.resume_url} target="_blank" rel="noreferrer">
                <Button variant="outline" size="sm" className="gap-2">
                  <FileDown className="h-4 w-4" />
                  Download Resume
                </Button>
              </a>
            )}

            {/* Notes */}
            {applicant.notes && (
              <div>
                <p className="text-xs text-muted-foreground mb-1">Notes</p>
                <p className="bg-slate-50 border rounded p-3 whitespace-pre-wrap">
                  {applicant.notes}
                </p>
              </div>
            )}

            {/* Rejection reason */}
            {applicant.rejection_reason && (
              <div>
                <p className="text-xs text-red-500 mb-1">Rejection Reason</p>
                <p className="bg-red-50 border border-red-200 text-red-700 rounded p-3">
                  {applicant.rejection_reason}
                </p>
              </div>
            )}

            {/* Stage history timeline */}
            {applicant.stage_histories?.length > 0 && (
              <div>
                <p className="text-xs font-semibold uppercase text-muted-foreground mb-2 tracking-wider">
                  Stage History
                </p>
                <div className="space-y-2">
                  {applicant.stage_histories.map((h: any) => (
                    <div
                      key={h.id}
                      className="flex items-start gap-3 text-xs border-l-2 border-blue-200 pl-3 py-1"
                    >
                      <div className="flex-1">
                        <p className="font-medium">
                          {h.from_stage ? `${h.from_stage} → ` : ""}{h.to_stage}
                        </p>
                        {h.notes && (
                          <p className="text-slate-500 mt-0.5">{h.notes}</p>
                        )}
                      </div>
                      <div className="text-right text-slate-400 flex-shrink-0 space-y-0.5">
                        <div className="flex items-center gap-1 justify-end">
                          <User className="h-3 w-3" />
                          <span>{h.changed_by_name || "-"}</span>
                        </div>
                        <div className="flex items-center gap-1 justify-end">
                          <Clock className="h-3 w-3" />
                          <span>{format(new Date(h.created_at), "dd MMM HH:mm")}</span>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            )}
          </div>
        )}
      </DialogContent>
    </Dialog>
  );
}
