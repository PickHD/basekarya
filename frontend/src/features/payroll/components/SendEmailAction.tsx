import { useState } from "react";
import { Mail, Loader2, Send } from "lucide-react";

import { Button } from "@/components/ui/button";
import {
  AlertDialog,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog";
import { useSendPayslipEmail } from "@/features/payroll/hooks/usePayroll";

interface SendEmailActionProps {
  payrollId: number;
  employeeName: string;
  period: string;
  isDraft?: boolean;
}

export function SendEmailAction({
  payrollId,
  employeeName,
  period,
  isDraft,
}: SendEmailActionProps) {
  const { mutate: sendEmail, isPending } = useSendPayslipEmail();
  const [open, setOpen] = useState(false);

  const handleConfirm = (e: React.MouseEvent) => {
    e.preventDefault();
    sendEmail(payrollId, {
      onSuccess: () => {
        setOpen(false);
      },
    });
  };

  return (
    <AlertDialog open={open} onOpenChange={setOpen}>
      <AlertDialogTrigger asChild>
        <Button
          variant="outline"
          size="sm"
          className="gap-2 text-slate-600 hover:text-blue-600 hover:bg-blue-50"
          disabled={isDraft}
          title={isDraft ? "Cannot send draft payslip" : "Send to email"}
        >
          <Mail className="w-4 h-4" />
          <span className="hidden sm:inline">Email PDF</span>
        </Button>
      </AlertDialogTrigger>

      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle className="flex items-center gap-2">
            <Send className="w-5 h-5 text-blue-600" />
            Kirim Slip Gaji via Email
          </AlertDialogTitle>
          <AlertDialogDescription className="leading-relaxed">
            Apakah Anda yakin ingin mengirim slip gaji periode{" "}
            <strong>{period}</strong> kepada <strong>{employeeName}</strong>?{" "}
            <br />
            <br />
            Tindakan ini akan membuat dokumen PDF dan mengirimkannya langsung ke
            alamat email mereka yang terdaftar.
          </AlertDialogDescription>
        </AlertDialogHeader>

        <AlertDialogFooter>
          <AlertDialogCancel disabled={isPending}>Cancel</AlertDialogCancel>
          <Button
            disabled={isPending}
            onClick={handleConfirm}
            className="bg-blue-600 hover:bg-blue-700 text-white min-w-[120px]"
          >
            {isPending ? (
              <>
                <Loader2 className="w-4 h-4 mr-2 animate-spin" />
                Sedang proses...
              </>
            ) : (
              "Ya, Kirim Email"
            )}
          </Button>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
}
