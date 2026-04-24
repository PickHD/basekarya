import { useSortable } from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Button } from "@/components/ui/button";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { MoreHorizontal, Eye, GripVertical, ArrowRight } from "lucide-react";
import { format } from "date-fns";
import type { Applicant, ApplicantStage } from "../types";
import { cn } from "@/lib/utils";

const NEXT_STAGES: Partial<Record<ApplicantStage, ApplicantStage>> = {
  SCREENING: "INTERVIEW",
  INTERVIEW: "OFFERING",
  OFFERING: "HIRED",
};

interface Props {
  applicant: Applicant;
  isDragging?: boolean;
  onClick: () => void;
  canUpdate: boolean;
  onStageChange: (stage: ApplicantStage) => void;
}

export function ApplicantCard({ applicant, isDragging, onClick, canUpdate, onStageChange }: Props) {
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging: isSortableDragging,
  } = useSortable({ id: applicant.id, disabled: !canUpdate });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
  };

  const initials = applicant.full_name
    .split(" ")
    .slice(0, 2)
    .map((n) => n[0])
    .join("")
    .toUpperCase();

  const nextStage = NEXT_STAGES[applicant.stage];

  return (
    <div
      ref={setNodeRef}
      style={style}
      {...(canUpdate ? attributes : {})}
      {...(canUpdate ? listeners : {})}
      className={cn(
        "bg-white rounded-lg border border-slate-200 shadow-sm p-3 space-y-2",
        "hover:shadow-md hover:border-slate-300 transition-all",
        canUpdate ? "cursor-grab active:cursor-grabbing" : "",
        (isDragging || isSortableDragging) && "opacity-50 shadow-lg rotate-1 scale-105 z-50 relative"
      )}
    >
      <div className="flex items-start gap-2">
        {/* Drag handle visual */}
        {canUpdate && (
          <div className="mt-0.5 text-slate-300 flex-shrink-0">
            <GripVertical className="h-4 w-4 pointer-events-none" />
          </div>
        )}

        {/* Avatar */}
        <Avatar className="h-7 w-7 flex-shrink-0">
          <AvatarFallback className="text-xs bg-blue-100 text-blue-700">
            {initials}
          </AvatarFallback>
        </Avatar>

        {/* Name + email */}
        <div className="flex-1 min-w-0" onClick={onClick}>
          <p className="text-sm font-medium leading-tight truncate">{applicant.full_name}</p>
          <p className="text-xs text-slate-500 truncate">{applicant.email}</p>
        </div>

        {/* Actions dropdown */}
        <DropdownMenu>
          <DropdownMenuTrigger asChild onClick={(e) => e.stopPropagation()}>
            <Button variant="ghost" className="h-6 w-6 p-0 flex-shrink-0">
              <MoreHorizontal className="h-3 w-3" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end" className="w-44">
            <DropdownMenuItem onClick={onClick}>
              <Eye className="mr-2 h-4 w-4" />
              View Detail
            </DropdownMenuItem>
            {canUpdate && nextStage && (
              <>
                <DropdownMenuSeparator />
                <DropdownMenuItem
                  onClick={(e) => {
                    e.stopPropagation();
                    onStageChange(nextStage);
                  }}
                >
                  <ArrowRight className="mr-2 h-4 w-4" />
                  Move to {nextStage.charAt(0) + nextStage.slice(1).toLowerCase()}
                </DropdownMenuItem>
              </>
            )}
            {canUpdate && applicant.stage !== "REJECTED" && (
              <DropdownMenuItem
                className="text-red-600"
                onClick={(e) => {
                  e.stopPropagation();
                  onStageChange("REJECTED");
                }}
              >
                <ArrowRight className="mr-2 h-4 w-4" />
                Reject
              </DropdownMenuItem>
            )}
          </DropdownMenuContent>
        </DropdownMenu>
      </div>

      {/* Date */}
      <p className="text-xs text-slate-400 pl-9">
        {format(new Date(applicant.created_at), "dd MMM yyyy")}
      </p>
    </div>
  );
}
