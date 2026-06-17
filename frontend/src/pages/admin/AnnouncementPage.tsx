import { CreateAnnouncementForm } from "@/features/announcement/components/CreateAnnouncementForm";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

export default function AnnouncementPage() {
  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-3xl font-bold tracking-tight">Announcement</h2>
        <p className="text-slate-500">
          Create and send important announcement to all employees.
        </p>
      </div>

      <Card>
        <CardHeader className="pb-3">
          <CardTitle className="text-lg font-semibold">Publish Announcement</CardTitle>
          <p className="text-slate-500">
            Fill in the form below clearly. Once published, all active employees will receive a notification.
          </p>
        </CardHeader>
        <CardContent>
          <CreateAnnouncementForm />
        </CardContent>
      </Card>
    </div>
  );
}
