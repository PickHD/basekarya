import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { TerBracketsTable } from "@/features/pph21-config/components/TerBracketsTable";
import { PtkpTable } from "@/features/pph21-config/components/PtkpTable";

export default function Pph21ConfigPage() {
  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h2 className="text-3xl font-bold tracking-tight">PPh 21 Config</h2>
          <p className="text-slate-500">
            Manage TER rate brackets and PTKP thresholds for PPh 21 calculations.
          </p>
        </div>
      </div>

      <Card>
        <CardHeader className="pb-3">
          <CardTitle className="text-lg font-semibold">Tax Configuration</CardTitle>
        </CardHeader>
        <CardContent>
          <Tabs defaultValue="ter-brackets">
            <TabsList>
              <TabsTrigger value="ter-brackets">TER Brackets</TabsTrigger>
              <TabsTrigger value="ptkp">PTKP Thresholds</TabsTrigger>
            </TabsList>
            <TabsContent value="ter-brackets" className="mt-4">
              <TerBracketsTable />
            </TabsContent>
            <TabsContent value="ptkp" className="mt-4">
              <PtkpTable />
            </TabsContent>
          </Tabs>
        </CardContent>
      </Card>
    </div>
  );
}
