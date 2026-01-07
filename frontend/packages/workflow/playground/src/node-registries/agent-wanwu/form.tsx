import React, { useState } from 'react';
import { I18n } from '@coze-arch/i18n';

import { NodeConfigForm } from '@/node-registries/common/components';
import { ViewVariableType } from '@coze-workflow/base';
import { OutputsField } from '../common/fields';
import { INPUT_PATH } from './constants';


import {
  Field,
  type FieldRenderProps,
  useForm,
} from '@flowgram-adapter/free-layout-editor';
import { type IModelValue } from '@/typing';
import { ModelSelect } from '@/components/model-select';
import { FormCard } from '@/form-extensions/components/form-card';
import { useReadonly } from '@/nodes-v2/hooks/use-readonly';
import { SystemPrompt } from '@/nodes-v2/llm-wanwu/system-prompt';
import { ToolSelectField, ToolInputParamsField } from './components';

export const FormRender = () => {
  const readonly = useReadonly();
  const form = useForm();
  const [params, setParams] = useState<any>(null);

  const handleToolChange = (newValue: any) => {
    const { properties } = newValue?.[0]?.inputSchema || {};
    setParams(properties ? Object.keys(properties) : []);
  };

  return (
    <NodeConfigForm>
      {/* 模型选择 */}
      <Field name={'inputs.llmParam'}>
        {({ field }: FieldRenderProps<IModelValue | undefined>) => (
          <FormCard
            header={I18n.t('workflow_detail_llm_model')}
            tooltip={I18n.t('workflow_detail_llm_prompt_tooltip')}
            required={true}
          >
            <ModelSelect {...field} readonly={readonly} />
          </FormCard>
        )}
      </Field>

      {/* 输入 */}
      <ToolInputParamsField name={INPUT_PATH} />

      {/* 系统提示词 */}
      <Field name="inputs.systemPrompt" deps={[INPUT_PATH]} defaultValue={''}>
        {({ field }: FieldRenderProps<string>) => (
          <SystemPrompt
            {...field}
            placeholder={I18n.t('workflow_detail_llm_sys_prompt_content')}
            inputParameters={form.getValueIn(INPUT_PATH)}
          />
        )}
      </Field>

      {/* 已选择工具 */}
      <ToolSelectField
        name={'inputs.toolInfoList'}
        afterChange={handleToolChange}
      />
      
      {/*隐藏域*/}
      <Field name="inputs.datasetSetting">
        {({ field }: FieldRenderProps<any>) => (
          <input type="hidden" {...field} />
        )}
      </Field>
      
      {/* 输出 */}
      <OutputsField
        title={I18n.t('workflow_detail_node_output')}
        tooltip={I18n.t('node_http_response_data')}
        id="agent-node-outputs"
        name="outputs"
        topLevelReadonly={true}
        customReadonly
      />
    </NodeConfigForm>
  );
};
