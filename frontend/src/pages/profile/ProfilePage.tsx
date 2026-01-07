import { Building2, Briefcase, User, Lock, Loader2 } from "lucide-react";

import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Badge } from "@/components/ui/badge";

import { useProfile } from "@/features/user/hooks/useProfile";
import { GeneralForm } from "@/features/user/components/GeneralForm";
import { PasswordForm } from "@/features/user/components/PasswordForm";

export default function ProfilePage() {
  const { data: user, isLoading } = useProfile();

  if (isLoading) {
    return (
      <div className="flex justify-center p-10">
        <Loader2 className="animate-spin" />
      </div>
    );
  }

  if (!user) return <div>Failed to load profile.</div>;

  return (
    <div className="space-y-6 max-w-6xl mx-auto">
      <div>
        <h2 className="text-3xl font-bold tracking-tight text-slate-900">
          Account Settings
        </h2>
        <p className="text-slate-500">
          Manage your identity and security preferences.
        </p>
      </div>

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
                <Badge variant="secondary" className="px-3 py-1">
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
          <Tabs defaultValue="general" className="w-full">
            <TabsList className="grid w-full grid-cols-2 mb-4">
              <TabsTrigger value="general" className="gap-2">
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
              <Card>
                <CardHeader>
                  <CardTitle>Change Password</CardTitle>
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
