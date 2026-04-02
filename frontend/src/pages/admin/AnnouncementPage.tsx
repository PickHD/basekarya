import { Megaphone } from "lucide-react";
import { CreateAnnouncementForm } from "@/features/announcement/components/CreateAnnouncementForm";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";

export default function AnnouncementPage() {
  return (
    <div className="flex flex-col gap-6 p-6 md:p-8 w-full max-w-5xl mx-auto">
      <div>
        <h1 className="text-3xl font-bold tracking-tight text-slate-900 flex items-center gap-2">
          <div className="p-2 bg-blue-100 rounded-lg">
            <Megaphone className="h-6 w-6 text-blue-700" />
          </div>
          Announcement
        </h1>
        <p className="text-slate-500 mt-2">
          Create and send important announcement to all employees.
        </p>
      </div>

      <Separator className="bg-slate-200" />

      <Card className="border-slate-200 shadow-sm overflow-hidden">
        <CardHeader className="bg-slate-50/50 border-b border-slate-100">
          <CardTitle className="text-xl text-slate-800">Publish Announcement</CardTitle>
          <CardDescription className="text-slate-500">
            Fill in the form below clearly. Once published, all active employees will receive a notification.
          </CardDescription>
        </CardHeader>
        <CardContent className="pt-6">
          <CreateAnnouncementForm />
        </CardContent>
      </Card>
    </div>
  );
}
