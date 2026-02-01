import { useState } from "react";
import {
  Building2,
  Briefcase,
  User,
  Lock,
  Loader2,
  TriangleAlert,
} from "lucide-react";

import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Badge } from "@/components/ui/badge";

import { useProfile } from "@/features/user/hooks/useProfile";
import { GeneralForm } from "@/features/user/components/GeneralForm";
import { PasswordForm } from "@/features/user/components/PasswordForm";

export default function ProfilePage() {
  const { data: user, isLoading } = useProfile();

  const [manualTab, setManualTab] = useState("general");

  const activeTab = user?.must_change_password ? "security" : manualTab;

  if (isLoading) {
    return (
      <div className="flex justify-center p-10">
        <Loader2 className="animate-spin h-8 w-8 text-blue-600" />
      </div>
    );
  }

  if (!user) return <div>Failed to load profile.</div>;

  const isLocked = user.must_change_password;

  return (
    <div className="space-y-6 max-w-6xl mx-auto pb-10">
      <div className="flex flex-col gap-2">
        <h2 className="text-3xl font-bold tracking-tight text-slate-900">
          Account Settings
        </h2>
        <p className="text-slate-500">
          Manage your identity and security preferences.
        </p>
      </div>

      {isLocked && (
        <Alert
          variant="destructive"
          className="bg-red-50 border-red-200 text-red-800 animate-in fade-in slide-in-from-top-2"
        >
          <TriangleAlert className="h-4 w-4" />
          <AlertTitle>Security Action Required</AlertTitle>
          <AlertDescription>
            For your security, you must change your default password before
            accessing other features.
          </AlertDescription>
        </Alert>
      )}

      <div className="grid gap-6 md:grid-cols-12">
        <div className="md:col-span-4 space-y-6">
          <Card className="border-t-4 border-t-blue-600 shadow-sm">
            <CardHeader className="text-center">
              <div className="mx-auto w-32 h-32 mb-4 relative">
                <Avatar className="w-32 h-32 border-4 border-slate-50 shadow-md">
                  <AvatarImage
                    src={user.profile_picture_url}
                    className="object-cover"
                  />
                  <AvatarFallback className="text-4xl bg-slate-100 text-slate-400">
                    {user.full_name?.charAt(0)}
                  </AvatarFallback>
                </Avatar>
              </div>
              <CardTitle className="text-xl">{user.full_name}</CardTitle>
              <CardDescription className="font-mono text-blue-600">
                {user.nik}
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="flex justify-center mb-6">
                <Badge variant="secondary" className="px-3 py-1 uppercase">
                  {user.role}
                </Badge>
              </div>
              <div className="space-y-4 text-sm border-t pt-4">
                <div className="flex items-center gap-3">
                  <Building2 className="w-4 h-4 text-slate-400" />
                  <span className="font-medium">
                    {user.department_name || "-"}
                  </span>
                </div>
                <div className="flex items-center gap-3">
                  <Briefcase className="w-4 h-4 text-slate-400" />
                  <span className="font-medium">{user.shift_name || "-"}</span>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>

        <div className="md:col-span-8">
          <Tabs
            value={activeTab}
            onValueChange={setManualTab}
            className="w-full"
          >
            <TabsList className="grid w-full grid-cols-2 mb-4">
              <TabsTrigger
                value="general"
                className="gap-2"
                disabled={isLocked}
              >
                <User className="w-4 h-4" /> General Info
              </TabsTrigger>
              <TabsTrigger value="security" className="gap-2">
                <Lock className="w-4 h-4" /> Security
              </TabsTrigger>
            </TabsList>

            <TabsContent value="general">
              <Card>
                <CardHeader>
                  <CardTitle>Personal Information</CardTitle>
                  <CardDescription>
                    Update your photo and contact details.
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <GeneralForm user={user} />
                </CardContent>
              </Card>
            </TabsContent>

            <TabsContent value="security">
              <Card
                className={
                  isLocked ? "border-red-200 shadow-red-100 shadow-md" : ""
                }
              >
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    Change Password
                    {isLocked && (
                      <Badge variant="destructive" className="text-xs">
                        Required
                      </Badge>
                    )}
                  </CardTitle>
                  <CardDescription>
                    Ensure your account stays secure.
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <PasswordForm />
                </CardContent>
              </Card>
            </TabsContent>
          </Tabs>
        </div>
      </div>
    </div>
  );
}
