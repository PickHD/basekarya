import { useState } from "react";
import { OnboardingWorkflowList } from "@/features/onboarding/components/OnboardingWorkflowList";
import { OnboardingDetailDialog } from "@/features/onboarding/components/OnboardingDetailDialog";
import { CreateWorkflowDialog } from "@/features/onboarding/components/CreateWorkflowDialog";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
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
      <div className="flex justify-between items-center">
        <div>
          <h2 className="text-3xl font-bold tracking-tight">Onboarding</h2>
          <p className="text-slate-500">
            Track onboarding progress for new hires across all departments.
          </p>
        </div>
      </div>

      <Card>
        <CardHeader className="pb-3">
          <CardTitle className="text-lg font-semibold">Onboarding Workflows</CardTitle>
        </CardHeader>
        <CardContent>
          <OnboardingWorkflowList
            onView={(w) => setSelectedWorkflow(w)}
            onCreateNew={() => setIsCreateOpen(true)}
            canCreate={canCreate}
          />
        </CardContent>
      </Card>

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
