import { useState } from "react";
import {
  DndContext,
  DragOverlay,
  PointerSensor,
  useSensor,
  useSensors,
  closestCorners,
  useDroppable,
} from "@dnd-kit/core";
import type { DragEndEvent, DragStartEvent } from "@dnd-kit/core";
import {
  SortableContext,
  verticalListSortingStrategy,
} from "@dnd-kit/sortable";
import { Loader2, Plus } from "lucide-react";
import { Button } from "@/components/ui/button";
import { ApplicantCard } from "./ApplicantCard";
import { ApplicantFormDialog } from "./ApplicantFormDialog";
import { ApplicantDetailDialog } from "./ApplicantDetailDialog";
import { useApplicants, useUpdateApplicantStage } from "../hooks/useApplicant";
import { usePermissions } from "@/hooks/usePermissions";
import { PERMISSIONS } from "@/config/permissions";
import type { Applicant, ApplicantStage, KanbanBoard } from "../types";

const STAGES: { key: ApplicantStage; label: string; color: string; headerColor: string }[] = [
  {
    key: "SCREENING",
    label: "Screening",
    color: "bg-slate-50 border-slate-200",
    headerColor: "bg-slate-100 text-slate-700",
  },
  {
    key: "INTERVIEW",
    label: "Interview",
    color: "bg-blue-50 border-blue-200",
    headerColor: "bg-blue-100 text-blue-700",
  },
  {
    key: "OFFERING",
    label: "Offering",
    color: "bg-amber-50 border-amber-200",
    headerColor: "bg-amber-100 text-amber-700",
  },
  {
    key: "HIRED",
    label: "Hired ✓",
    color: "bg-emerald-50 border-emerald-200",
    headerColor: "bg-emerald-100 text-emerald-700",
  },
  {
    key: "REJECTED",
    label: "Rejected",
    color: "bg-red-50 border-red-200",
    headerColor: "bg-red-100 text-red-700",
  },
];

interface Props {
  requisitionId: number;
}

export function ApplicantKanbanBoard({ requisitionId }: Props) {
  const { hasPermission } = usePermissions();
  const { data: board, isLoading } = useApplicants(requisitionId);
  const { mutate: updateStage } = useUpdateApplicantStage();

  const [activeApplicant, setActiveApplicant] = useState<Applicant | null>(null);
  const [isAddOpen, setIsAddOpen] = useState(false);
  const [selectedApplicant, setSelectedApplicant] = useState<Applicant | null>(null);

  const canCreate = hasPermission(PERMISSIONS.CREATE_APPLICANT);
  const canUpdate = hasPermission(PERMISSIONS.UPDATE_APPLICANT);

  const sensors = useSensors(
    useSensor(PointerSensor, { activationConstraint: { distance: 5 } })
  );

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-20">
        <Loader2 className="h-6 w-6 animate-spin text-slate-400" />
        <span className="ml-2 text-slate-500">Loading applicants...</span>
      </div>
    );
  }

  const kanban: KanbanBoard = board ?? {
    SCREENING: [],
    INTERVIEW: [],
    OFFERING: [],
    HIRED: [],
    REJECTED: [],
  };

  const findApplicantById = (id: number): Applicant | null => {
    for (const stage of STAGES) {
      const found = kanban[stage.key].find((a) => a.id === id);
      if (found) return found;
    }
    return null;
  };

  const handleDragStart = (event: DragStartEvent) => {
    const applicant = findApplicantById(Number(event.active.id));
    setActiveApplicant(applicant);
  };

  const handleDragEnd = (event: DragEndEvent) => {
    setActiveApplicant(null);
    const { active, over } = event;
    if (!over || !canUpdate) return;

    const applicantId = Number(active.id);
    let targetStage: ApplicantStage | undefined;

    // The `over.id` could be a stage key (droppable col) or another applicant id (sortable item)
    if (STAGES.some((s) => s.key === over.id)) {
      targetStage = over.id as ApplicantStage;
    } else {
      const hoveredApplicant = findApplicantById(Number(over.id));
      if (hoveredApplicant) targetStage = hoveredApplicant.stage;
    }

    if (targetStage) {
      const currentApplicant = findApplicantById(applicantId);
      if (currentApplicant && currentApplicant.stage !== targetStage) {
        updateStage({ id: applicantId, payload: { stage: targetStage } });
      }
    }
  };

  return (
    <div className="space-y-4">
      {/* Add applicant button */}
      {canCreate && (
        <div className="flex justify-end">
          <Button
            onClick={() => setIsAddOpen(true)}
            className="bg-blue-600 hover:bg-blue-700"
          >
            <Plus className="mr-2 h-4 w-4" />
            Add Applicant
          </Button>
        </div>
      )}

      {/* Kanban Board */}
      <DndContext
        sensors={sensors}
        collisionDetection={closestCorners}
        onDragStart={handleDragStart}
        onDragEnd={handleDragEnd}
      >
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-5 gap-4 overflow-x-auto">
          {STAGES.map((stage) => {
            const applicants = kanban[stage.key] ?? [];
            return (
              <KanbanColumn
                key={stage.key}
                stage={stage}
                applicants={applicants}
                canUpdate={canUpdate}
                onSelectApplicant={setSelectedApplicant}
                onUpdateStage={(id, stage) => updateStage({ id, payload: { stage } })}
              />
            );
          })}
        </div>

        {/* Drag overlay — shows the card being dragged */}
        <DragOverlay>
          {activeApplicant && (
            <ApplicantCard
              applicant={activeApplicant}
              isDragging
              onClick={() => {}}
              canUpdate={false}
              onStageChange={() => {}}
            />
          )}
        </DragOverlay>
      </DndContext>

      {/* Dialogs */}
      <ApplicantFormDialog
        open={isAddOpen}
        onOpenChange={setIsAddOpen}
        requisitionId={requisitionId}
      />
      <ApplicantDetailDialog
        open={!!selectedApplicant}
        onOpenChange={(open) => !open && setSelectedApplicant(null)}
        applicantId={selectedApplicant?.id ?? null}
      />
    </div>
  );
}

interface KanbanColumnProps {
  stage: typeof STAGES[0];
  applicants: Applicant[];
  canUpdate: boolean;
  onUpdateStage: (id: number, stage: ApplicantStage) => void;
  onSelectApplicant: (applicant: Applicant) => void;
}

function KanbanColumn({ stage, applicants, canUpdate, onUpdateStage, onSelectApplicant }: KanbanColumnProps) {
  const { setNodeRef } = useDroppable({ id: stage.key });

  return (
    <div
      ref={setNodeRef}
      className={`flex flex-col rounded-xl border ${stage.color} min-h-[400px]`}
    >
      {/* Column header */}
      <div className={`flex items-center justify-between px-3 py-2 rounded-t-xl ${stage.headerColor}`}>
        <span className="text-sm font-semibold">{stage.label}</span>
        <span className="text-xs font-medium bg-white/60 px-2 py-0.5 rounded-full">
          {applicants.length}
        </span>
      </div>

      {/* Cards */}
      <SortableContext
        id={stage.key}
        items={applicants.map((a) => a.id)}
        strategy={verticalListSortingStrategy}
      >
        <div className="flex flex-col gap-2 p-2 flex-1 overflow-y-auto min-h-[100px]">
          {applicants.length === 0 ? (
            <div className="flex items-center justify-center h-20 text-xs text-slate-400 border-2 border-dashed border-slate-200 rounded-lg m-1">
              Drop here
            </div>
          ) : (
            applicants.map((applicant) => (
              <ApplicantCard
                key={applicant.id}
                applicant={applicant}
                onClick={() => onSelectApplicant(applicant)}
                canUpdate={canUpdate}
                onStageChange={(newStage) =>
                  onUpdateStage(applicant.id, newStage)
                }
              />
            ))
          )}
        </div>
      </SortableContext>
    </div>
  );
}
