import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { AssetList } from "@/features/asset/components/AssetList";
import { AssetAssignmentList } from "@/features/asset/components/AssetAssignmentList";

const AssetsListPage = () => {
  return (
    <Tabs defaultValue="assets" className="space-y-6">
      <TabsList>
        <TabsTrigger value="assets">Daftar Aset</TabsTrigger>
        <TabsTrigger value="assignments">Permintaan Aset</TabsTrigger>
      </TabsList>
      <TabsContent value="assets">
        <AssetList />
      </TabsContent>
      <TabsContent value="assignments">
        <AssetAssignmentList />
      </TabsContent>
    </Tabs>
  );
};

export default AssetsListPage;
