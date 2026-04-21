import { useState } from "react";
import { format } from "date-fns";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Loader2, Mail, CheckCircle2, Circle, User } from "lucide-react";
import { useOnboardingWorkflowDetail, useCompleteTask } from "@/features/onboarding/hooks/useOnboarding";
import type { OnboardingTask } from "@/features/onboarding/types";

interface Props {
  open: boolean;
  onOpenChange: (v: boolean) => void;
  workflowId: number | null;
  canComplete: boolean;
}

export function OnboardingDetailDialog({ open, onOpenChange, workflowId, canComplete }: Props) {
  const { data: workflow, isLoading } = useOnboardingWorkflowDetail(workflowId);
  const { mutate: completeTask } = useCompleteTask();
  const [completing, setCompleting] = useState<number | null>(null);

  const handleComplete = (taskId: number) => {
    setCompleting(taskId);
    completeTask(
      { id: taskId, notes: "" },
      { onSettled: () => setCompleting(null) }
    );
  };

  const TaskList = ({ tasks }: { tasks: OnboardingTask[] }) => (
    <div className="space-y-2">
      {tasks.length === 0 && (
        <p className="text-sm text-slate-400 py-4 text-center">No tasks in this category.</p>
      )}
      {tasks.map((task) => (
        <div
          key={task.id}
          className={`flex items-start gap-3 p-3 rounded-lg border transition-colors ${
            task.is_completed
              ? "bg-emerald-50 border-emerald-200"
              : "bg-white border-slate-200"
          }`}
        >
          {canComplete && !task.is_completed ? (
            <Checkbox
              id={`task-${task.id}`}
              checked={task.is_completed}
              disabled={completing === task.id}
              onCheckedChange={() => handleComplete(task.id)}
              className="mt-0.5"
            />
          ) : (
            <div className="mt-0.5">
              {task.is_completed ? (
                <CheckCircle2 className="h-4 w-4 text-emerald-500 flex-shrink-0" />
              ) : (
                <Circle className="h-4 w-4 text-slate-300 flex-shrink-0" />
              )}
            </div>
          )}
          <div className="flex-1 min-w-0">
            <p className={`text-sm font-medium ${task.is_completed ? "line-through text-slate-400" : "text-slate-800"}`}>
              {task.task_name}
            </p>
            {task.description && (
              <p className="text-xs text-slate-500 mt-0.5">{task.description}</p>
            )}
            {task.is_completed && (
              <div className="flex items-center gap-1.5 mt-1">
                <User className="h-3 w-3 text-emerald-500" />
                <span className="text-xs text-emerald-600">
                  {task.completed_by || "Unknown"} ·{" "}
                  {task.completed_at ? format(new Date(task.completed_at), "dd MMM yyyy HH:mm") : ""}
                </span>
              </div>
            )}
            {task.notes && (
              <p className="text-xs text-slate-500 mt-1 italic">"{task.notes}"</p>
            )}
          </div>
          {completing === task.id && (
            <Loader2 className="h-4 w-4 animate-spin text-slate-400 flex-shrink-0 mt-0.5" />
          )}
        </div>
      ))}
    </div>
  );

  const totalTasks = workflow
    ? (workflow.it_tasks?.length ?? 0) + (workflow.hr_tasks?.length ?? 0) + (workflow.other_tasks?.length ?? 0)
    : 0;

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>Onboarding Workflow</DialogTitle>
        </DialogHeader>

        {isLoading || !workflow ? (
          <div className="flex justify-center py-12">
            <Loader2 className="h-6 w-6 animate-spin text-slate-400" />
          </div>
        ) : (
          <div className="space-y-5">
            {/* New hire info card */}
            <div className="bg-gradient-to-r from-blue-50 to-indigo-50 rounded-xl p-4 border border-blue-100">
              <div className="flex items-start justify-between gap-4">
                <div>
                  <p className="text-base font-semibold text-slate-800">{workflow.new_hire_name}</p>
                  <p className="text-sm text-slate-500 mt-0.5">{workflow.new_hire_email}</p>
                  <div className="flex flex-wrap gap-3 mt-2 text-xs text-slate-600">
                    {workflow.position && <span>🏷️ {workflow.position}</span>}
                    {workflow.department && <span>🏢 {workflow.department}</span>}
                    {workflow.start_date && (
                      <span>📅 {format(new Date(workflow.start_date), "dd MMM yyyy")}</span>
                    )}
                  </div>
                </div>
                <div className="flex flex-col items-end gap-2">
                  {workflow.status === "COMPLETED" ? (
                    <Badge className="bg-emerald-100 text-emerald-700">✓ Completed</Badge>
                  ) : (
                    <Badge className="bg-blue-100 text-blue-700">In Progress</Badge>
                  )}
                  {workflow.welcome_email_sent && (
                    <span className="flex items-center gap-1 text-xs text-emerald-600">
                      <Mail className="h-3 w-3" /> Welcome email sent
                    </span>
                  )}
                </div>
              </div>

              {/* Progress bar */}
              <div className="mt-3">
                <div className="flex justify-between text-xs text-slate-500 mb-1">
                  <span>Progress</span>
                  <span>{workflow.progress}% ({totalTasks} tasks)</span>
                </div>
                <div className="h-2 bg-white/60 rounded-full overflow-hidden">
                  <div
                    className={`h-full rounded-full transition-all duration-500 ${
                      workflow.progress === 100 ? "bg-emerald-500" : "bg-blue-500"
                    }`}
                    style={{ width: `${workflow.progress}%` }}
                  />
                </div>
              </div>
            </div>

            {/* Tasks by department */}
            <Tabs defaultValue="it">
              <TabsList className="grid w-full grid-cols-3">
                <TabsTrigger value="it">
                  IT Tasks
                  <span className="ml-1.5 text-xs bg-slate-100 text-slate-600 rounded-full px-1.5 py-0.5">
                    {workflow.it_tasks?.length ?? 0}
                  </span>
                </TabsTrigger>
                <TabsTrigger value="hr">
                  HR Tasks
                  <span className="ml-1.5 text-xs bg-slate-100 text-slate-600 rounded-full px-1.5 py-0.5">
                    {workflow.hr_tasks?.length ?? 0}
                  </span>
                </TabsTrigger>
                <TabsTrigger value="other">
                  Other
                  <span className="ml-1.5 text-xs bg-slate-100 text-slate-600 rounded-full px-1.5 py-0.5">
                    {workflow.other_tasks?.length ?? 0}
                  </span>
                </TabsTrigger>
              </TabsList>
              <TabsContent value="it" className="mt-4">
                <TaskList tasks={workflow.it_tasks ?? []} />
              </TabsContent>
              <TabsContent value="hr" className="mt-4">
                <TaskList tasks={workflow.hr_tasks ?? []} />
              </TabsContent>
              <TabsContent value="other" className="mt-4">
                <TaskList tasks={workflow.other_tasks ?? []} />
              </TabsContent>
            </Tabs>

            <div className="flex justify-end">
              <Button variant="outline" onClick={() => onOpenChange(false)}>Close</Button>
            </div>
          </div>
        )}
      </DialogContent>
    </Dialog>
  );
}
