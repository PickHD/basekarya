import { useEffect } from "react";
import { useForm, useFieldArray } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Button } from "@/components/ui/button";
import { Loader2, Plus, Trash2, GripVertical } from "lucide-react";
import { useUpdateWorkflowTasks } from "@/features/onboarding/hooks/useOnboarding";
import type { OnboardingTask } from "@/features/onboarding/types";

const itemSchema = z.object({
  task_name: z.string().min(2, "Task name required"),
  description: z.string().optional().default(""),
  sort_order: z.coerce.number().default(0),
});

const schema = z.object({
  tasks: z.array(itemSchema).min(1, "At least one task required"),
});

type FormValues = z.infer<typeof schema>;

interface Props {
  open: boolean;
  onOpenChange: (v: boolean) => void;
  workflowId: number | null;
  pendingTasks: OnboardingTask[];
}

export function ManageTasksDialog({ open, onOpenChange, workflowId, pendingTasks }: Props) {
  const { mutateAsync: updateTasks, isPending } = useUpdateWorkflowTasks();

  const form = useForm<FormValues, unknown, FormValues>({
    resolver: zodResolver(schema) as any,
    defaultValues: {
      tasks: [{ task_name: "", description: "", sort_order: 1 }],
    },
  });

  const { fields, append, remove } = useFieldArray({
    control: form.control,
    name: "tasks",
  });

  useEffect(() => {
    if (open) {
      if (pendingTasks.length > 0) {
        form.reset({
          tasks: pendingTasks.map((t) => ({
            task_name: t.task_name,
            description: t.description,
            sort_order: t.sort_order,
          })),
        });
      } else {
        form.reset({
          tasks: [{ task_name: "", description: "", sort_order: 1 }],
        });
      }
    }
  }, [open, pendingTasks, form]);

  const onSubmit = async (values: FormValues) => {
    if (!workflowId) return;
    await updateTasks({
      id: workflowId,
      payload: {
        tasks: values.tasks.map((item, i) => ({
          task_name: item.task_name,
          description: item.description ?? "",
          sort_order: item.sort_order ?? i + 1,
        })),
      },
    });
    onOpenChange(false);
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>Manage Tasks</DialogTitle>
        </DialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
            <div className="space-y-3">
              <div className="flex items-center justify-between">
                <p className="text-sm text-slate-500">
                  Completed tasks are preserved. Updating will replace all pending tasks.
                </p>
                <Button
                  type="button"
                  variant="outline"
                  size="sm"
                  className="text-xs h-7"
                  onClick={() => append({ task_name: "", description: "", sort_order: fields.length + 1 })}
                >
                  <Plus className="h-3 w-3 mr-1" /> Add Task
                </Button>
              </div>

              {fields.map((field, index) => (
                <div key={field.id} className="flex gap-2 items-start p-3 bg-slate-50 rounded-lg border border-slate-200">
                  <GripVertical className="h-4 w-4 text-slate-300 mt-2.5 flex-shrink-0" />
                  <div className="flex-1 space-y-2">
                    <FormField
                      control={form.control}
                      name={`tasks.${index}.task_name`}
                      render={({ field }) => (
                        <FormItem>
                          <FormControl>
                            <Input placeholder={`Task ${index + 1} name`} className="text-sm" {...field} />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                    <FormField
                      control={form.control}
                      name={`tasks.${index}.description`}
                      render={({ field }) => (
                        <FormItem>
                          <FormControl>
                            <Textarea
                              placeholder="Description (optional)"
                              rows={2}
                              className="text-sm resize-none"
                              {...field}
                            />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                  </div>
                  <Button
                    type="button"
                    variant="ghost"
                    size="sm"
                    className="h-8 w-8 p-0 text-red-400 hover:text-red-600 hover:bg-red-50 flex-shrink-0 mt-0.5"
                    onClick={() => remove(index)}
                    disabled={fields.length === 1}
                  >
                    <Trash2 className="h-3.5 w-3.5" />
                  </Button>
                </div>
              ))}
            </div>

            <DialogFooter>
              <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
                Cancel
              </Button>
              <Button type="submit" disabled={isPending}>
                {isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                Save Tasks
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}
