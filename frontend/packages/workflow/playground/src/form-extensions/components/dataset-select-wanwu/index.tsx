/*
 * Copyright 2025 coze-dev Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import React, { useMemo, useState, useEffect, useRef } from 'react';

import {
  StandardNodeType,
  useWorkflowNode,
  useNodeTestId,
} from '@coze-workflow/base';
import { FilterKnowledgeType } from '@coze-data/knowledge-modal-base';
import { useKnowledgeListModal } from '@coze-data/knowledge-modal-adapter';
import { useBizWorkflowKnowledgeIDEFullScreenModal } from '@coze-data/knowledge-ide-adapter';
import { I18n } from '@coze-arch/i18n';
import { ConfigProvider } from '@coze-arch/bot-semi';
import { useGlobalState, useDataSetInfos } from '@/hooks';

import { LibrarySelect } from '../library-select-wanwu';
import { MetadataFilterModal, DEFAULT_METADATA } from './metadata-filter-modal-wanwu'

interface ValueProps {
  dataset_id?: string;
  name?: string;
}

export const DatasetSelect = ({
  value: _value,
  onChange,
  readonly,
}: {
  value: ValueProps[];
  onChange: (value: ValueProps[]) => void;
  readonly: boolean;
}) => {
  const { getNodeSetterId } = useNodeTestId();
  const addButtonTestID = getNodeSetterId('dataset-select-add');
  const libraryCardTestID = getNodeSetterId('dataset-select-card');

  // When initializing, it will be worn to the default value: [null]
  const value:any = useMemo(() => _value?.filter?.(d => !!d) || [], [_value]);

  const { dataSets, cacheDataSetInfo } = useDataSetInfos({ ids: value.map(item => item.dataset_id) });
  const res = useWorkflowNode();
  const { type } = res;
  const isDatasetWrite = type === StandardNodeType.DatasetWrite;
  const knowledgeType = isDatasetWrite
    ? FilterKnowledgeType.TEXT
    : FilterKnowledgeType.ALL;
  const hiddenAddBtn = !!(isDatasetWrite && value?.length);

  const [curDateSetID, setCurDateSetID] = useState<string>();
  // metadata filter modal
  const [visible, setVisible] = useState(false);
  // metadata total data
  const [currentMetaData, setCurrentMetaData] = useState<any>(DEFAULT_METADATA);
  // knowledge id
  const [knowledgeId, setKnowledgeId] = useState<string>('');

  const handleClose = () => {
    setVisible(false)
  }

  const {
    node: knowledgePreviewModal,
  } = useBizWorkflowKnowledgeIDEFullScreenModal({
    biz: 'workflow',
  });

  const { projectId } = useGlobalState();
  const newWindowRef = useRef<WindowProxy | null>();

  const { node, open, close } = useKnowledgeListModal({
    datasetList: value.map(d => ({
      ...d,
    })),
    onClickKnowledgeDetail:(datasetID: string) => {
      setCurDateSetID(datasetID);
    },
    onDatasetListChange: list => {
      cacheDataSetInfo(list);
      onChange(list.map((item:any) => ({
        dataset_id: item.dataset_id,
        name: item.name,
        knowledgeId: item.knowledgeId,
        graphSwitch: item.graphSwitch,
        external: item.external,
        ragName: item.ragName,
        metaDataFilterParams: item.metaDataFilterParams
      })) as object[]);
    },
    defaultType: knowledgeType,
    // Pass undefined to display the full knowledge base
    knowledgeTypeConfigList: isDatasetWrite
      ? [FilterKnowledgeType.TEXT]
      : undefined,
    projectID: projectId,
    beforeCreate: shouldUpload => {
      if (shouldUpload && !projectId) {
        newWindowRef.current = window.open();
      }
    },
  });

  const showEditMetaData = async (id) => {
    // saved dataset data, get current dataset metadata object
    const currentValueItem:any = value.find(item => item.dataset_id === id) || {}
    const { metaDataFilterParams } = currentValueItem
    // origin dataset data, get current dataset knowledgeId
    const currentDataset:any = dataSets.find(item => item.dataset_id === id) || {}

    setCurDateSetID(id);
    setKnowledgeId(currentDataset.knowledgeId);
    setCurrentMetaData(metaDataFilterParams || DEFAULT_METADATA);
    setVisible(true);
  }

  useEffect(() => {
    // The knowledge base writing node can only select one knowledge base and then close it, limiting multiple selection
    if (isDatasetWrite && value?.length) {
      close();
    }
  }, [value, close, isDatasetWrite]);

  return (
    <div className="relative">
      {/* Semi is mounted on the body by default, but it is overridden and mounted on the current node in workflow. You need to manually overwrite it to the body here, otherwise the knowledge pop-up window will not open. */}
      <ConfigProvider getPopupContainer={() => document.body}>
        {knowledgePreviewModal}
        {node}
      </ConfigProvider>
      <MetadataFilterModal
        visible={visible}
        handleClose={handleClose}
        defaultValue={currentMetaData}
        knowledgeId={knowledgeId}
        onSubmit={(metaDataFilterParams) => {
          const newValue = value.map(item => (
            item.dataset_id === curDateSetID ? {...item, metaDataFilterParams} : {...item}
          ))
          onChange(newValue)
          handleClose()
        }}
      />
      <LibrarySelect
        libraries={dataSets?.map(
          ({
            dataset_id = '',
            name,
            description,
            external,
            orgName,
            share,
            category,
            icon_url,
          }) => ({
            id: dataset_id,
            name,
            description,
            external,
            iconUrl: icon_url,
            extraInfo: { orgName, share, category, external },

            // Invalid Knowledge Base is disabled
            isInvalid: true,
          }),
        )}
        readonly={readonly}
        onEditLibrary={id => {
          if (id) {
            showEditMetaData(id)
          }
        }}
        onDeleteLibrary={id => {
          onChange(value.filter(_v => _v?.dataset_id !== id));
        }}
        onAddLibrary={() => (open())}
        onClickLibrary={id => {
          if (id) {
            setCurDateSetID(id);
          }
        }}
        emptyText={
          isDatasetWrite
            ? I18n.t('kl_write_003')
            : I18n.t('workflow_knowledge_node_empty')
        }
        hideAddButton={hiddenAddBtn}
        addButtonTestID={addButtonTestID}
        libraryCardTestID={libraryCardTestID}
      />
    </div>
  );
};
