import { Badge } from "@/components/ui/badge";
import type { ApplicantStage } from "../types";

const stageConfig: Record<ApplicantStage, { label: string; className: string }> = {
  SCREENING: {
    label: "Screening",
    className: "bg-slate-100 text-slate-700 border-slate-200",
  },
  INTERVIEW: {
    label: "Interview",
    className: "bg-blue-50 text-blue-700 border-blue-200",
  },
  OFFERING: {
    label: "Offering",
    className: "bg-amber-50 text-amber-700 border-amber-200",
  },
  HIRED: {
    label: "Hired",
    className: "bg-emerald-50 text-emerald-700 border-emerald-200",
  },
  REJECTED: {
    label: "Rejected",
    className: "bg-red-50 text-red-700 border-red-200",
  },
};

export function StageBadge({ stage }: { stage: ApplicantStage }) {
  const config = stageConfig[stage] ?? stageConfig.SCREENING;
  return (
    <Badge variant="outline" className={config.className}>
      {config.label}
    </Badge>
  );
}
