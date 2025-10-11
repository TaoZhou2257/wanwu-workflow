import { I18n } from '@coze-arch/i18n';
import { NodeConfigForm } from '@/node-registries/common/components';

import { OutputsField } from '../common/fields';
import { ViewVariableType } from '@coze-workflow/base';
import { FileParamsField } from './components/file-params-field';

export const FormRender = () => (
  <NodeConfigForm>
    <FileParamsField
      inputFieldName="inputs.inputParameters.FileUrlList"
      testId="/inputs/inputParameters/FileUrlList"
      tooltip={I18n.t(
        'workflow_detail_multi_file_parse_input_tooltip',
        {},
        '输入文档url',
      )}
      paramName={'FileUrlList'}
      paramType={ViewVariableType.ArrayString}
      inputType={ViewVariableType.ArrayString}
    />

    <OutputsField
      title={I18n.t('workflow_detail_node_output')}
      tooltip={I18n.t('node_http_response_data')}
      id="multiFileParseWanwu-node-outputs"
      name="outputs"
      topLevelReadonly={true}
      customReadonly
    />
  </NodeConfigForm>
);
