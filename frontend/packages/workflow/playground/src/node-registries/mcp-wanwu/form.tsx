import React, { useState } from "react";
import { I18n } from '@coze-arch/i18n';

import { NodeConfigForm } from '@/node-registries/common/components';
import { INPUT_PATH } from './constants'

import { OutputsField } from '../common/fields';
import { McpParamsField, McpSelectField } from "./componets";

export const FormRender = () => {
  const [params, setParams] = useState<any>(null)

  const afterChange = (newValue: any) => {
    const { properties } = newValue?.[0]?.inputSchema || {}
    setParams(properties ? Object.keys(properties) : [])
  }

  return (
    <NodeConfigForm>
      <McpSelectField
        name={'inputs.mcpInfoList'}
        afterChange={afterChange}
      />
      <McpParamsField
        name={INPUT_PATH}
        params={params}
        // defaultValue={params}
      />
      <OutputsField
        title={I18n.t('workflow_detail_node_output')}
        tooltip={I18n.t('node_http_response_data')}
        id="mcpWanwu-node-outputs"
        name="outputs"
        topLevelReadonly={true}
        customReadonly
      />
    </NodeConfigForm>
  )
}
