import React, {useEffect, useState} from "react";
import { I18n } from '@coze-arch/i18n';

import { NodeConfigForm } from '@/node-registries/common/components';

import { OutputsField, InputsParametersField } from '../common/fields';

import { McpParamsField } from "./componets/mcp-params-field";

export const FormRender = () => {
  const [params, setParams] = useState<any>([{name: 'key1', input: { type: 'ref' } }, {name: 'key2', input: { type: 'ref' } }])
  useEffect(() => {
    setParams([{name: 'key1', input: { type: 'ref' } }, {name: 'key2', input: { type: 'ref' } }])
  }, []);

  return (
    <NodeConfigForm>
      {/*<InputsParametersField
        name={INPUT_PATH}
        title={I18n.t('node_http_request_params')}
        tooltip={I18n.t('node_http_request_params_desc')}
        defaultValue={[]}
      />*/}
      <McpParamsField
        name="inputParameters"
        defaultValue={params}
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
