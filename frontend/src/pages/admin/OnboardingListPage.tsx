import { useState } from "react";
import { OnboardingWorkflowList } from "@/features/onboarding/components/OnboardingWorkflowList";
import { OnboardingDetailDialog } from "@/features/onboarding/components/OnboardingDetailDialog";
import { CreateWorkflowDialog } from "@/features/onboarding/components/CreateWorkflowDialog";
import { usePermissions } from "@/hooks/usePermissions";
import { PERMISSIONS } from "@/config/permissions";
import type { OnboardingWorkflowList as WorkflowListType } from "@/features/onboarding/types";

export default function OnboardingListPage() {
  const { hasPermission } = usePermissions();
  const canCreate = hasPermission(PERMISSIONS.MANAGE_ONBOARDING_TEMPLATE);
  const canComplete = hasPermission(PERMISSIONS.UPDATE_ONBOARDING_TASK);

  const [selectedWorkflow, setSelectedWorkflow] = useState<WorkflowListType | null>(null);
  const [isCreateOpen, setIsCreateOpen] = useState(false);

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-xl font-bold text-slate-800">Onboarding Workflows</h1>
        <p className="text-sm text-slate-500 mt-1">
          Track onboarding progress for new hires across all departments.
        </p>
      </div>

      <OnboardingWorkflowList
        onView={(w) => setSelectedWorkflow(w)}
        onCreateNew={() => setIsCreateOpen(true)}
        canCreate={canCreate}
      />

      <OnboardingDetailDialog
        open={!!selectedWorkflow}
        onOpenChange={(v) => !v && setSelectedWorkflow(null)}
        workflowId={selectedWorkflow?.id ?? null}
        canComplete={canComplete}
      />

      <CreateWorkflowDialog
        open={isCreateOpen}
        onOpenChange={setIsCreateOpen}
      />
    </div>
  );
}
