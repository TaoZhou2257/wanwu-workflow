import { I18n } from '@coze-arch/i18n';
import { NodeConfigForm } from '@/node-registries/common/components';

import { OutputsField } from '../common/fields';
import { ViewVariableType } from '@coze-workflow/base';
import { FileParamsField } from './components/file-params-field';

export const FormRender = () => (
  <NodeConfigForm>
    <FileParamsField
      inputFieldName="inputs.inputParameters.FileUrl"
      testId="/inputs/inputParameters/FileUrl"
      tooltip={I18n.t(
        'workflow_detail_file_parse_input_tooltip',
        {},
        '输入文档url',
      )}
      paramName={'FileUrl'}
      paramType={ViewVariableType.String}
      inputType={ViewVariableType.String}
    />

    <OutputsField
      title={I18n.t('workflow_detail_node_output')}
      tooltip={I18n.t('node_http_response_data')}
      id="fileParseWanwu-node-outputs"
      name="outputs"
      topLevelReadonly={true}
      customReadonly
    />
  </NodeConfigForm>
);
