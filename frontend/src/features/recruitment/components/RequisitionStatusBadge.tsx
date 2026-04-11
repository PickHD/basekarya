import { Badge } from "@/components/ui/badge";
import type { RequisitionStatus } from "../types";

const statusConfig: Record<RequisitionStatus, { label: string; className: string }> = {
  DRAFT: {
    label: "Draft",
    className: "bg-slate-100 text-slate-700 border-slate-200",
  },
  PENDING: {
    label: "Pending",
    className: "bg-amber-50 text-amber-700 border-amber-200",
  },
  APPROVED: {
    label: "Approved",
    className: "bg-emerald-50 text-emerald-700 border-emerald-200",
  },
  REJECTED: {
    label: "Rejected",
    className: "bg-red-50 text-red-700 border-red-200",
  },
  CLOSED: {
    label: "Closed",
    className: "bg-gray-100 text-gray-600 border-gray-200",
  },
};

export function RequisitionStatusBadge({ status }: { status: RequisitionStatus }) {
  const config = statusConfig[status] ?? statusConfig.DRAFT;
  return (
    <Badge variant="outline" className={config.className}>
      {config.label}
    </Badge>
  );
}
