import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog";
import { Loader2, Plus, Edit2, Trash2, ChevronDown, ChevronUp, LayoutTemplate } from "lucide-react";
import { useOnboardingTemplates, useDeleteTemplate } from "@/features/onboarding/hooks/useOnboarding";
import { TemplateFormDialog } from "@/features/onboarding/components/TemplateFormDialog";
import type { OnboardingTemplate } from "@/features/onboarding/types";

export function OnboardingTemplateManager() {
  const { data: templates = [], isLoading } = useOnboardingTemplates();
  const { mutate: deleteTemplate, isPending: isDeleting } = useDeleteTemplate();
  const [isCreateOpen, setIsCreateOpen] = useState(false);
  const [editingTemplate, setEditingTemplate] = useState<OnboardingTemplate | null>(null);
  const [expandedId, setExpandedId] = useState<number | null>(null);

  const deptColor = (dept: string) => {
    if (dept === "IT") return "bg-blue-100 text-blue-700 border-blue-200";
    if (dept === "HR") return "bg-emerald-100 text-emerald-700 border-emerald-200";
    return "bg-slate-100 text-slate-700 border-slate-200";
  };

  return (
    <div className="space-y-4">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          <LayoutTemplate className="h-5 w-5 text-blue-600" />
          <div>
            <h3 className="text-sm font-semibold text-slate-800">Onboarding Templates</h3>
            <p className="text-xs text-slate-500">Manage checklist templates applied to all new hires</p>
          </div>
        </div>
        <Button
          className="bg-blue-600 hover:bg-blue-700 text-sm h-9"
          onClick={() => setIsCreateOpen(true)}
        >
          <Plus className="mr-2 h-4 w-4" /> New Template
        </Button>
      </div>

      {/* List */}
      {isLoading ? (
        <div className="flex justify-center py-12">
          <Loader2 className="h-5 w-5 animate-spin text-slate-400" />
        </div>
      ) : templates.length === 0 ? (
        <div className="rounded-xl border-2 border-dashed border-slate-200 py-12 text-center">
          <LayoutTemplate className="h-8 w-8 text-slate-300 mx-auto mb-2" />
          <p className="text-sm text-slate-400">No templates yet. Create your first template.</p>
        </div>
      ) : (
        <div className="space-y-3">
          {(templates as OnboardingTemplate[]).map((tmpl) => (
            <div key={tmpl.id} className="rounded-xl border border-slate-200 overflow-hidden">
              {/* Template header row */}
              <div className="flex items-center gap-3 px-4 py-3 bg-white hover:bg-slate-50 transition-colors">
                <button
                  className="flex-1 flex items-center gap-3 text-left"
                  onClick={() => setExpandedId(expandedId === tmpl.id ? null : tmpl.id)}
                >
                  <div className="flex-1">
                    <p className="text-sm font-medium text-slate-800">{tmpl.name}</p>
                    <div className="flex items-center gap-2 mt-0.5">
                      <Badge variant="outline" className={`text-xs ${deptColor(tmpl.department)}`}>
                        {tmpl.department}
                      </Badge>
                      <span className="text-xs text-slate-400">{tmpl.items.length} tasks</span>
                    </div>
                  </div>
                  {expandedId === tmpl.id ? (
                    <ChevronUp className="h-4 w-4 text-slate-400 flex-shrink-0" />
                  ) : (
                    <ChevronDown className="h-4 w-4 text-slate-400 flex-shrink-0" />
                  )}
                </button>

                <div className="flex items-center gap-1">
                  <Button
                    variant="ghost"
                    size="sm"
                    className="h-8 w-8 p-0"
                    onClick={() => setEditingTemplate(tmpl)}
                  >
                    <Edit2 className="h-3.5 w-3.5 text-slate-500" />
                  </Button>
                  <AlertDialog>
                    <AlertDialogTrigger asChild>
                      <Button variant="ghost" size="sm" className="h-8 w-8 p-0">
                        <Trash2 className="h-3.5 w-3.5 text-red-400" />
                      </Button>
                    </AlertDialogTrigger>
                    <AlertDialogContent>
                      <AlertDialogHeader>
                        <AlertDialogTitle>Delete Template</AlertDialogTitle>
                        <AlertDialogDescription>
                          Are you sure you want to delete "{tmpl.name}"? This will not affect existing workflows.
                        </AlertDialogDescription>
                      </AlertDialogHeader>
                      <AlertDialogFooter>
                        <AlertDialogCancel>Cancel</AlertDialogCancel>
                        <AlertDialogAction
                          className="bg-red-600 hover:bg-red-700"
                          onClick={() => deleteTemplate(tmpl.id)}
                          disabled={isDeleting}
                        >
                          Delete
                        </AlertDialogAction>
                      </AlertDialogFooter>
                    </AlertDialogContent>
                  </AlertDialog>
                </div>
              </div>

              {/* Expanded task list */}
              {expandedId === tmpl.id && (
                <div className="bg-slate-50 border-t border-slate-200 px-4 py-3">
                  <div className="space-y-2">
                    {tmpl.items.map((item, idx) => (
                      <div key={item.id} className="flex items-start gap-2 text-sm">
                        <span className="text-xs text-slate-400 w-5 mt-0.5 font-mono">{idx + 1}.</span>
                        <div>
                          <p className="text-slate-700 font-medium">{item.task_name}</p>
                          {item.description && (
                            <p className="text-xs text-slate-500 mt-0.5">{item.description}</p>
                          )}
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              )}
            </div>
          ))}
        </div>
      )}

      {/* Dialogs */}
      <TemplateFormDialog
        open={isCreateOpen}
        onOpenChange={setIsCreateOpen}
      />
      <TemplateFormDialog
        open={!!editingTemplate}
        onOpenChange={(v) => !v && setEditingTemplate(null)}
        editingTemplate={editingTemplate}
      />
    </div>
  );
}
