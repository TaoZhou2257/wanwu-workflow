import { I18n } from '@coze-arch/i18n';

import { NodeConfigForm } from '@/node-registries/common/components';

import { OutputsField } from '../common/fields';
import { INPUT_PATH } from './constants';
import { GuiAgentParamsField } from "./components/gui-agent-params-field";
import React from "react";

export const FormRender = () => (
  <NodeConfigForm>
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
);
