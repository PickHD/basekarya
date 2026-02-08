import { useEffect, useState, useRef } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
  FormDescription,
} from "@/components/ui/form";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Loader2, Paperclip, X, CalendarIcon } from "lucide-react";
import { useApplyLeave, useLeaveTypes } from "../hooks/useLeave";
import { differenceInDays, parseISO } from "date-fns";
import { toast } from "sonner";

const leaveApplySchema = z
  .object({
    leave_type_id: z.string().min(1, "leave type id required"),
    start_date: z.string().min(1, "start date required"),
    end_date: z.string().min(1, "end date required"),
    reason: z.string().min(5, "reason required"),
    attachment_base64: z.string().optional(),
  })
  .refine(
    (data) => {
      if (!data.start_date || !data.end_date) return true;
      return new Date(data.end_date) >= new Date(data.start_date);
    },
    {
      message: "Tanggal selesai tidak boleh sebelum tanggal mulai",
      path: ["end_date"],
    },
  );

type LeaveApplyFormValues = z.infer<typeof leaveApplySchema>;

interface Props {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

export function LeaveApplyDialog({ open, onOpenChange }: Props) {
  const { data: leaveTypes } = useLeaveTypes();
  const { mutate: apply, isPending } = useApplyLeave();

  const [fileName, setFileName] = useState<string | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const form = useForm<LeaveApplyFormValues>({
    resolver: zodResolver(leaveApplySchema),
    defaultValues: {
      leave_type_id: "",
      start_date: "",
      end_date: "",
      reason: "",
      attachment_base64: "",
    },
  });

  useEffect(() => {
    if (open) {
      form.reset({
        leave_type_id: "",
        start_date: "",
        end_date: "",
        reason: "",
        attachment_base64: "",
      });
      setFileName(null);
      if (fileInputRef.current) fileInputRef.current.value = "";
    }
  }, [open, form]);

  const startDate = form.watch("start_date");
  const endDate = form.watch("end_date");

  const calculateDays = () => {
    if (!startDate || !endDate) return 0;
    const start = parseISO(startDate);
    const end = parseISO(endDate);
    const diff = differenceInDays(end, start);
    return diff >= 0 ? diff + 1 : 0;
  };

  const totalDays = calculateDays();

  const convertToBase64 = (file: File): Promise<string> => {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();
      reader.readAsDataURL(file);
      reader.onload = () => resolve(reader.result as string);
      reader.onerror = (error) => reject(error);
    });
  };

  const handleFileChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    if (file.size > 2 * 1024 * 1024) {
      toast.error("Ukuran file maksimal 2MB");
      return;
    }
    // Validasi Tipe
    if (!["image/jpeg", "image/png"].includes(file.type)) {
      toast.error("Format harus JPG, PNG");
      return;
    }

    try {
      const base64 = await convertToBase64(file);
      form.setValue("attachment_base64", base64, { shouldValidate: true });
      setFileName(file.name);
    } catch (err) {
      console.error("Gagal convert file", err);
      toast.error("Gagal memproses file");
    }
  };

  const removeFile = () => {
    form.setValue("attachment_base64", "", { shouldValidate: true });
    setFileName(null);
    if (fileInputRef.current) fileInputRef.current.value = "";
  };

  const onSubmit = (data: LeaveApplyFormValues) => {
    if (calculateDays() <= 0) {
      toast.error("Durasi cuti tidak valid");
      return;
    }
    const rawBase64 = data.attachment_base64?.split(",")[1];

    apply(
      {
        leave_type_id: Number(data.leave_type_id),
        start_date: data.start_date,
        end_date: data.end_date,
        reason: data.reason,
        attachment_base64: rawBase64,
      },
      {
        onSuccess: () => {
          onOpenChange(false);
        },
      },
    );
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>Form Pengajuan Cuti</DialogTitle>
        </DialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
            <FormField
              control={form.control}
              name="leave_type_id"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Jenis Cuti</FormLabel>
                  <Select
                    onValueChange={field.onChange}
                    defaultValue={field.value}
                  >
                    <FormControl>
                      <SelectTrigger>
                        <SelectValue placeholder="Pilih jenis cuti..." />
                      </SelectTrigger>
                    </FormControl>
                    <SelectContent>
                      {leaveTypes?.map((type) => (
                        <SelectItem key={type.id} value={String(type.id)}>
                          {type.name} (Kuota: {type.default_quota})
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                  <FormMessage />
                </FormItem>
              )}
            />

            <div className="grid grid-cols-2 gap-4">
              <FormField
                control={form.control}
                name="start_date"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Tanggal Mulai</FormLabel>
                    <FormControl>
                      <Input type="date" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="end_date"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Tanggal Selesai</FormLabel>
                    <FormControl>
                      <Input type="date" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>

            {totalDays > 0 && (
              <div className="text-sm text-blue-600 bg-blue-50 p-2 rounded flex items-center gap-2 animate-in fade-in">
                <CalendarIcon className="h-4 w-4" />
                Total durasi: <strong>{totalDays} hari</strong>
              </div>
            )}

            <FormField
              control={form.control}
              name="attachment_base64"
              render={() => (
                <FormItem>
                  <FormLabel>Lampiran (Opsional)</FormLabel>
                  <FormDescription>
                    Surat dokter atau bukti pendukung.
                  </FormDescription>
                  <FormControl>
                    <div>
                      <input
                        type="file"
                        ref={fileInputRef}
                        className="hidden"
                        accept="image/jpeg,image/png"
                        onChange={handleFileChange}
                      />

                      {!fileName ? (
                        <div
                          className="border-2 border-dashed rounded-md p-6 flex flex-col items-center justify-center cursor-pointer hover:bg-slate-50 transition-colors text-slate-500"
                          onClick={() => fileInputRef.current?.click()}
                        >
                          <Paperclip className="h-8 w-8 mb-2 text-slate-400" />
                          <span className="text-sm font-medium">
                            Klik untuk upload file
                          </span>
                          <span className="text-xs text-slate-400 mt-1">
                            Max 2MB (JPG/PNG/PDF)
                          </span>
                        </div>
                      ) : (
                        <div className="flex items-center justify-between p-3 bg-blue-50 border border-blue-100 rounded-md">
                          <div className="flex items-center gap-2 overflow-hidden">
                            <Paperclip className="h-4 w-4 text-blue-600 flex-shrink-0" />
                            <span className="text-sm text-blue-900 truncate max-w-[250px]">
                              {fileName}
                            </span>
                          </div>
                          <Button
                            type="button"
                            variant="ghost"
                            size="sm"
                            className="h-6 w-6 p-0 hover:bg-blue-100 rounded-full"
                            onClick={removeFile}
                          >
                            <X className="h-3 w-3 text-blue-600" />
                          </Button>
                        </div>
                      )}
                    </div>
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="reason"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Alasan / Keterangan</FormLabel>
                  <FormControl>
                    <Textarea
                      placeholder="Contoh: Acara keluarga di luar kota..."
                      className="resize-none"
                      rows={3}
                      {...field}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <DialogFooter>
              <Button
                type="button"
                variant="outline"
                onClick={() => onOpenChange(false)}
              >
                Batal
              </Button>
              <Button type="submit" disabled={isPending || totalDays <= 0}>
                {isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                {isPending ? "Mengirim..." : "Ajukan Permintaan"}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}
