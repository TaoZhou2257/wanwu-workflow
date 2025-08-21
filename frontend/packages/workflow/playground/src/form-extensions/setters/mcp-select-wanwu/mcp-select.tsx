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

import { useEffect, useState } from 'react'
import { I18n } from '@coze-arch/i18n';
import { Modal } from '@coze-arch/coze-design';
import { SelectMcpModal } from './components';
import { useWorkflowNode } from '@coze-workflow/base';

import { useGlobalState } from '@/hooks';
import { LibrarySelect } from '@/form-extensions/components/library-select';

import { useModal } from './use-modal';
import { type McpSelectValue } from './types';

interface McpSelectProps {
  value?: McpSelectValue;
  onChange?: (newValue: McpSelectValue) => void;
  readonly?: boolean;
  addButtonTestID?: string;
  libraryCardTestID?: string;
}

export const McpSelect = ({
  value,
  onChange,
  readonly = false,
  addButtonTestID,
  libraryCardTestID,
}: McpSelectProps) => {
  const { spaceId, projectId, getProjectApi, playgroundProps } =
    useGlobalState();
  const {
    isVisible: isSelectDatabaseModalVisible,
    openModal: openSelectDatabaseModal,
    closeModal: closeSelectMcpModal,
  } = useModal();
  const { data } = useWorkflowNode();
  const mcpList = data?.inputs?.mcpInfoList;
  const [libraries, setLibraries] = useState<any>([])

  function changeMcp(item: any) {
    onChange?.([item]);
  }

  function clearMcp() {
    onChange?.([]);
    setLibraries([])
  }

  function handleSelectMcpModalAdd(item: any) {
    changeMcp(item);
    setLibraries([item])
    closeSelectMcpModal();
  }

  function handleLibrarySelectDelete() {
    Modal.confirm({
      title: I18n.t(
        'workflow_mcp_delete_confirm_modal_title',
        {},
        '确认移除该条数据？',
      ),
      content: I18n.t(
        'workflow_mcp_delete_confirm_modal_content',
        {},
        '移除后，该节点配置的相关内容均会被删除且无法恢复',
      ),
      onOk: () => {
        clearMcp();
      },
      okText: I18n.t('workflow_confirm_modal_ok', {}, '确定'),
      cancelText: I18n.t('workflow_confirm_modal_cancel', {}, '取消'),
    });
  }

  useEffect(() => {
    setLibraries(mcpList)
  }, []);

  return (
    <>
      <LibrarySelect
        readonly={readonly}
        libraries={libraries}
        onAddLibrary={openSelectDatabaseModal}
        onDeleteLibrary={handleLibrarySelectDelete}
        emptyText={I18n.t('workflow_mcp_node_database_empty')}
        hideAddButton={value && value?.length > 0}
        addButtonTestID={addButtonTestID}
        libraryCardTestID={libraryCardTestID}
      />
      <SelectMcpModal
        spaceId={spaceId}
        visible={isSelectDatabaseModalVisible}
        onClose={closeSelectMcpModal}
        onAddMcp={handleSelectMcpModalAdd}
        enterFrom="workflow"
        projectID={projectId}
      />
    </>
  );
};
