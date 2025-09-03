import React from "react";
import { I18n } from '@coze-arch/i18n';

import { NodeConfigForm } from '@/node-registries/common/components';
import {
  Field,
  type FieldRenderProps,
} from '@flowgram-adapter/free-layout-editor';
import { type IModelValue } from '@/typing';
import { FormCard } from '@/form-extensions/components/form-card';
import { useReadonly } from '@/nodes-v2/hooks/use-readonly';

import { OutputsField } from '../common/fields';
import { INPUT_PATH } from './constants';
import { GuiAgentParamsField, GuiModelWanwu } from "./components";

export const FormRender = () => {
  const readonly = useReadonly();
  return(
    <NodeConfigForm>
      <Field name={'inputs.guiParams.modelId'}>
        {({ field }: FieldRenderProps<IModelValue | undefined>) => (
          <FormCard
            header={I18n.t('workflow_detail_llm_model')}
            tooltip={I18n.t('workflow_detail_llm_prompt_tooltip')}
          >
            <GuiModelWanwu
              {...field}
              readonly={readonly}
            />
          </FormCard>
        )}
      </Field>
      <GuiAgentParamsField name={INPUT_PATH} />

      <OutputsField
        title={I18n.t('workflow_detail_node_output')}
        tooltip={I18n.t('node_http_response_data')}
        id="guiAgentWanwu-node-outputs"
        name="outputs"
        topLevelReadonly={true}
        customReadonly
      />
    </NodeConfigForm>
  )
};
