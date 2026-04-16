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

import { useMemo, useState, useEffect,type FC } from 'react';
import { I18n } from '@coze-arch/i18n';
import { TOOL_TAB } from '../../types';
import classNames from 'classnames';
import {
  UICompositionModal,
  UICompositionModalMain,
} from '@coze-arch/bot-semi';
import {
  useWorkflowModalParts,
  WorkflowModalFrom,
} from '@coze-workflow/components';
import { usePluginModalParts } from '@coze-agent-ide/bot-plugin-export/agentSkillPluginModal/hooks';
import { From } from '@coze-agent-ide/plugin-shared';

import { useGlobalState } from '@/hooks';
import { useSelectMcpModal } from '../../../mcp-select-wanwu/components/mcp-select-modal-wanwu';
import { type ToolSelectValue } from '../../types';
import { ToolSelectDatabase } from '../tool-select-database-wanwu';
import { SkillSelect } from '../../../skill-select-wanwu';
import { useForm } from '@/form';

interface ToolSelectProps {
  visible: boolean;
  onClose: () => void;
  onAddTool?: (item: any) => void;
  onRemoveTool: (item: any) => void;
  spaceId: string;
  enterFrom?: string;
  projectID?: string;
  value?: ToolSelectValue;
  onChange?: (newValue: ToolSelectValue) => void;
  readonly?: boolean;
  addButtonTestID?: string;
  libraryCardTestID?: string;
  width?: string | number;
  showDataset?: boolean | undefined;
  onlyShowSkill?: boolean | undefined;
}

type ToolTab = typeof TOOL_TAB[keyof typeof TOOL_TAB];

export const SelectToolModal: FC<ToolSelectProps> = ({
  visible,
  onClose,
  onAddTool,
  onRemoveTool,
  showDataset,
  onlyShowSkill,
  spaceId
}) => {
  const { projectId } = useGlobalState();
  const form = useForm();
  const [activeTab, setActiveTab] = useState<ToolTab>(TOOL_TAB.TOOL);
  // const pluginModalFrom = projectId ? From.ProjectWorkflow : From.WorkflowAddNode;
  const pluginModalFrom = From.ProjectIde;//隐藏工具添加中的计数
  const [pluginApiList, setPluginApiList] = useState<any[]>(
    form.getValueIn('inputs.toolInfoList')?.filter(item => item.kind === TOOL_TAB.TOOL) || []
  );
  const [workFlowList, setWorkFlowList] = useState<any[]>(
    form.getValueIn('inputs.toolInfoList')?.filter(item => item.kind === TOOL_TAB.WORKFLOW) || []
  );
  const [mcpList, setMcpList] = useState<any[]>(
    form.getValueIn('inputs.toolInfoList')?.filter(item => item.kind === TOOL_TAB.MCP) || []
  );
  const [skillList, setSkillList] = useState<any[]>(
    form.getValueIn('inputs.toolInfoList')?.filter(item => item.kind === TOOL_TAB.SKILL) || []
  );


  useEffect(() => {
    const toolInfoList = form.getValueIn('inputs.toolInfoList') || [];
    if (!toolInfoList) return;
    const newPluginList = toolInfoList.filter(item => item.kind === TOOL_TAB.TOOL);
    const newWorkflowList = toolInfoList.filter(item => item.kind === TOOL_TAB.WORKFLOW);
    const newMcpList = toolInfoList.filter(item => item.kind === TOOL_TAB.MCP);
    const newSkillList = toolInfoList.filter(item => item.kind === TOOL_TAB.SKILL);
    setPluginApiList(newPluginList);
    setWorkFlowList(newWorkflowList);
    setMcpList(newMcpList);
    setSkillList(newSkillList);
  }, [form.getValueIn('inputs.toolInfoList')]);

  useEffect(() => {
    setActiveTab(showDataset ? TOOL_TAB.DATABASE : onlyShowSkill ? TOOL_TAB.SKILL : TOOL_TAB.TOOL)
  }, [showDataset, onlyShowSkill]);

  const workflowAddList = useMemo(() => {
    return mcpList.map(item => item.name);
  }, [mcpList]);

  const skillAddList = useMemo(() => {
    return skillList.map(item => item.id);
  }, [skillList]);

  const pluginModalParts = usePluginModalParts({
    pluginApiList,
    onPluginApiListChange: setPluginApiList,
    openModeCallback: api => {
      const isRemoved = pluginApiList.some(
        item => item.api_id === api?.api_id
      );
      if(isRemoved) {
        onRemoveTool(api?.api_id);
        const newList = pluginApiList.filter(
          item => item.api_id !== api?.api_id
        );
        setPluginApiList(newList);
        return true;
      }
      onAddTool?.({ kind: TOOL_TAB.TOOL, data: api });
      return true;
    },
    from: pluginModalFrom,
    projectId,
  });
  const workflowModalParts = useWorkflowModalParts({
    from: WorkflowModalFrom.BotSkills,//隐藏计数
    projectId,
    workFlowList,
    onWorkFlowListChange: setWorkFlowList,
    onAdd: item => onAddTool?.({ kind: TOOL_TAB.WORKFLOW, data: item }) || null,
    onRemove: item => onRemoveTool(item.workflow_id),
  });

  const mcpSelectParts = useSelectMcpModal({
    visible: activeTab === TOOL_TAB.MCP,
    onClose: () => null,
    onAddMcp: (item: any) => {
      onAddTool?.({ kind: TOOL_TAB.MCP, data: item });
    },
    enterFrom: TOOL_TAB.WORKFLOW,
    spaceId,
    projectID: projectId,
    workflowAddList,
    onRemoveMcp: (id: any) => {
      onRemoveTool(id)
    },
  });

  const tabList = useMemo<
    Array<{ key: ToolTab; label: string }>
  >(
    () => showDataset ? [
      {
        key: TOOL_TAB.DATABASE,
        label: I18n.t(
          'Datasets' as any,
          {
            resource: I18n.t('resource_type_database' as any),
          },
        ),
      },
    ] : onlyShowSkill ? [
      {
        key: TOOL_TAB.SKILL,
        label: 'Skills',
      },
    ] : [
      {
        key: TOOL_TAB.TOOL,
        label: I18n.t('workflow_tool_tab' as any, {}, '工具'),
      },
      {
        key: TOOL_TAB.WORKFLOW,
        label: I18n.t(
          'Submit_workflow_list' as any,
          {
            resource: I18n.t('library_resource_type_workflow' as any),
          },
        ),
      },
      {
        key: TOOL_TAB.MCP,
        label: 'MCP',
      },
      {
        key: TOOL_TAB.SKILL,
        label: 'Skills',
      },
    ],
    [],
  );

  const renderHeader = () => (
    <div className="flex flex-row gap-[8px] items-center w-full">
      <div className="flex flex-row gap-[8px] items-center">
        {tabList.map(item => (
          <div
            key={item.key}
            className={classNames(
              'px-[8px] py-[4px] text-[14px] rounded cursor-pointer transition-colors whitespace-nowrap',
              activeTab === item.key
                ? 'text-[var(--Primary-coz-primary,#3370ff)] bg-[var(--Primary-coz-primary-light,#3370ff1a)]'
                : 'coz-fg-secondary hover:coz-fg-primary',
            )}
            onClick={() => setActiveTab(item.key)}
          >
            {item.label}
          </div>
        ))}
      </div>
    </div>
  );

  const renderContent = () => (
    <div className="w-full h-full flex flex-col">
      {!showDataset && !onlyShowSkill && activeTab === TOOL_TAB.TOOL && (
        // 工具选择
        <div className="flex-1 flex flex-row min-h-0 h-full overflow-hidden">
          <div className="w-[200px] mr-[12px] border-r border-[rgba(255,255,255,0.06)] h-full overflow-y-auto">
            {pluginModalParts.sider}
          </div>
          <div className="flex-1 flex flex-col min-h-0 h-full overflow-hidden">
            <div className="mb-[8px] flex-shrink-0">{pluginModalParts.filter}</div>
            <div className="flex-1 min-h-0 overflow-y-auto max-h-full">{pluginModalParts.content}</div>
          </div>
        </div>
      )}
      {!showDataset && !onlyShowSkill &&  activeTab === TOOL_TAB.WORKFLOW && (
        // 工作流选择
        <div className="flex-1 flex flex-row min-h-0">
          <div className="w-[220px] mr-[12px] border-r border-[rgba(255,255,255,0.06)]">
            {workflowModalParts.sider}
          </div>
          <div className="flex-1 flex flex-col min-h-0">
            <div className="mb-[8px]">{workflowModalParts.filter}</div>
            <div className="flex-1 min-h-0 overflow-y-auto">{workflowModalParts.content}</div>
          </div>
        </div>
      )}
      {showDataset && activeTab === TOOL_TAB.DATABASE && (
        <ToolSelectDatabase
          spaceId={spaceId}
          projectID={projectId}
          onRemoveTool={onRemoveTool}
          onAddDatabase={databaseData => {
            onAddTool?.({ kind: TOOL_TAB.DATABASE, data: databaseData });
          }}
        />
      )}

      {!showDataset && !onlyShowSkill && activeTab === TOOL_TAB.MCP && (
        <div className="flex-1 min-h-0">
          {mcpSelectParts.renderContent()}
        </div>
      )}

      {!showDataset && activeTab === TOOL_TAB.SKILL && (
        <div className="flex-1 min-h-0">
          <SkillSelect
            spaceId={spaceId}
            projectID={projectId}
            skillList={skillAddList}
            onAddSkill={(item) => {
              onAddTool?.({ kind: TOOL_TAB.SKILL, data: item });
            }}
            onRemoveSkill={(id) => {
              onRemoveTool(id);
            }}
          />
        </div>
      )}
    </div>
  );

  if (!visible) {
    return null;
  }

  return (
    <UICompositionModal
      closable
      visible={visible}
      onCancel={onClose}
      header={onlyShowSkill ? null : renderHeader()}
      sider={null}
      style={{ width: '1200px' }}
      content={
        <UICompositionModalMain className="relative px-[12px] gap-[16px]">
          {renderContent()}
        </UICompositionModalMain>
      }
    />
  );
};
