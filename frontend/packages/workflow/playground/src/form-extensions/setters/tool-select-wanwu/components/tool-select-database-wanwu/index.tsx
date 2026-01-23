import { useState, type FC } from 'react';
import { useForm } from '@/form';
import { FilterKnowledgeType } from '@coze-data/knowledge-modal-base';
import {
  useKnowledgeListModalContent,
  useQAKnowledgeListModalContent,
} from '@coze-data/knowledge-modal-base';

import { useGlobalState } from '@/hooks';

import { DatabaseModalSider } from './tool-select-database-sider';

const getInitialDatasetList = () => {
  const form = useForm();
  const selected = form.getValueIn('inputs.toolInfoList')?.filter(item => item.kind === 'database') || [];
  return selected;
};
interface ToolSelectDatabaseProps {
  spaceId: string;
  projectID?: string;
  onRemoveTool: (item: any) => void;
  onAddDatabase?: (data: {
    dataset_id: string;
    name?: string;
    knowledgeId?: string;
    graphSwitch?: boolean;
    ragName?: string;
    desc?: string;
    metaDataFilterParams?: any;
    external?: any;
  }) => void;
}

export const ToolSelectDatabase: FC<ToolSelectDatabaseProps> = ({
  spaceId,
  projectID,
  onAddDatabase,
  onRemoveTool,
}) => {
  const [activeTab, setActiveTab] = useState<'database' | 'qaBase'>('database');
  const [datasetList, setDatasetList] = useState<any[]>(getInitialDatasetList);
  const [qaDatasetList, setQaDatasetList] = useState<any[]>([]);
  const { projectId } = useGlobalState();

  const { renderContent: renderKnowledgeList, renderSearch: renderKnowledgeSearch } = useKnowledgeListModalContent({
    datasetList,
    onDatasetListChange: list => {
      setDatasetList(list);
      const addedItems = list.filter(item => 
        !datasetList.some(d => d.dataset_id === item.dataset_id)
      );
      const removedItems = datasetList.filter(item => 
        !list.some(d => d.dataset_id === item.dataset_id)
      );
      if (removedItems.length > 0 && onRemoveTool) {
        onRemoveTool(removedItems[0]['knowledgeId']);
        return;
      }
      if(addedItems.length > 0 && onAddDatabase) {
        const lastItem = list[list.length - 1] as any;
        onAddDatabase({
          dataset_id: lastItem.dataset_id || '',
          name: lastItem.name,
          knowledgeId: lastItem.knowledgeId,
          graphSwitch: lastItem.graphSwitch,
          external: lastItem.external,
          ragName: lastItem.ragName,
          desc: lastItem.description || '',
          metaDataFilterParams: lastItem.metaDataFilterParams
        });
      }
    },
    defaultType: FilterKnowledgeType.ALL,
    projectID: projectId,
    hideHeader: true,
    showFilters: [],
  });

  const { renderContent: renderQAKnowledgeList, renderSearch: renderQASearch } = useQAKnowledgeListModalContent({
    datasetList: qaDatasetList,
    onDatasetListChange: list => {
      setQaDatasetList(list);
      if (list.length > 0 && onAddDatabase) {
        const lastItem = list[list.length - 1] as any;
        onAddDatabase({
          dataset_id: lastItem.dataset_id || '',
          name: lastItem.name,
          knowledgeId: lastItem.knowledgeId,
          graphSwitch: lastItem.graphSwitch,
          external: lastItem.external,
          ragName: lastItem.ragName,
          desc: lastItem.description || '',
          metaDataFilterParams: lastItem.metaDataFilterParams
        });
      }
    },
    defaultType: FilterKnowledgeType.ALL,
    projectID: projectId,
    hideHeader: true,
    showFilters: [],
  });

  return (
    <div className="w-full h-full flex flex-row min-h-0">
      <div className="w-[220px] mr-[12px] border-r border-[rgba(255,255,255,0.06)]">
        <DatabaseModalSider
          activeTab={activeTab}
          renderSearch={renderKnowledgeSearch}
          renderQaSearch={renderQASearch}
          onTabChange={key => setActiveTab(key as 'database' | 'qaBase')}
        />
      </div>
      <div className="flex-1 min-h-0">
        {activeTab === 'database' && renderKnowledgeList()}
        {activeTab === 'qaBase' && renderQAKnowledgeList()}
      </div>
    </div>
  );
};

