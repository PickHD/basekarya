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
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Button } from "@/components/ui/button";
import { Loader2, Plus, Trash2, GripVertical } from "lucide-react";
import { useCreateTemplate, useUpdateTemplate } from "@/features/onboarding/hooks/useOnboarding";
import type { OnboardingTemplate } from "@/features/onboarding/types";

const itemSchema = z.object({
  task_name: z.string().min(2, "Task name required"),
  description: z.string().optional().default(""),
  sort_order: z.coerce.number().default(0),
});

const schema = z.object({
  name: z.string().min(2, "Name required"),
  department: z.string().min(1, "Department required"),
  items: z.array(itemSchema).min(1, "At least one task required"),
});

type FormValues = z.infer<typeof schema>;

interface Props {
  open: boolean;
  onOpenChange: (v: boolean) => void;
  editingTemplate?: OnboardingTemplate | null;
}

export function TemplateFormDialog({ open, onOpenChange, editingTemplate }: Props) {
  const isEdit = !!editingTemplate;
  const { mutateAsync: create, isPending: isCreating } = useCreateTemplate();
  const { mutateAsync: update, isPending: isUpdating } = useUpdateTemplate();
  const isPending = isCreating || isUpdating;

  const form = useForm<FormValues, unknown, FormValues>({
    resolver: zodResolver(schema) as any,
    defaultValues: {
      name: "",
      department: "",
      items: [{ task_name: "", description: "", sort_order: 1 }],
    },
  });

  const { fields, append, remove } = useFieldArray({
    control: form.control,
    name: "items",
  });

  useEffect(() => {
    if (open) {
      if (editingTemplate) {
        form.reset({
          name: editingTemplate.name,
          department: editingTemplate.department,
          items: editingTemplate.items.map((i) => ({
            task_name: i.task_name,
            description: i.description,
            sort_order: i.sort_order,
          })),
        });
      } else {
        form.reset({
          name: "",
          department: "",
          items: [{ task_name: "", description: "", sort_order: 1 }],
        });
      }
    }
  }, [open, editingTemplate, form]);

  const onSubmit = async (values: FormValues) => {
    const payload = {
      name: values.name,
      department: values.department,
      items: values.items.map((item, i) => ({
        task_name: item.task_name,
        description: item.description ?? "",
        sort_order: item.sort_order ?? i + 1,
      })),
    };
    if (isEdit && editingTemplate) {
      await update({ id: editingTemplate.id, payload });
    } else {
      await create(payload);
    }
    onOpenChange(false);
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>{isEdit ? "Edit Template" : "Create Onboarding Template"}</DialogTitle>
        </DialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <FormField
                control={form.control}
                name="name"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Template Name</FormLabel>
                    <FormControl>
                      <Input placeholder="e.g. IT Setup" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="department"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Department</FormLabel>
                    <FormControl>
                      <Input placeholder="e.g. IT, HR" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>

            {/* Task Items */}
            <div className="space-y-3">
              <div className="flex items-center justify-between">
                <h4 className="text-sm font-semibold text-slate-900">Tasks</h4>
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
                      name={`items.${index}.task_name`}
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
                      name={`items.${index}.description`}
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
              <Button type="submit" disabled={isPending} className="bg-blue-600 hover:bg-blue-700">
                {isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                {isEdit ? "Update Template" : "Create Template"}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}
