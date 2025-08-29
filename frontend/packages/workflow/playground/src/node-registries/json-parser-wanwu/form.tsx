import { I18n } from '@coze-arch/i18n';

import { NodeConfigForm } from '@/node-registries/common/components';
import { INPUT_PATH } from './constants';
import { InputsField } from './components/inputs';
import {
  Field,
} from '@flowgram-adapter/free-layout-editor';
import { Outputs } from '@/nodes-v2/components/outputs';
import { ViewVariableType } from '@coze-workflow/base';

export const FormRender = ({form}) => (
  <NodeConfigForm>
    <InputsField
      name={INPUT_PATH}
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      defaultValue={[{ name: 'input' } as any]}
      title={I18n.t('node_http_request_params')}
      tooltip={I18n.t('workflow_250429_03')}
      required={false}
      layout="horizontal"
    />

      <Field
        name={'outputs'}
        deps={['batchMode']}
        defaultValue={[{ name: 'output', type: ViewVariableType.String }]}
      >
        {({ field, fieldState }) => (
          <Outputs
            id={'json-parser-wanwu-node-output'}
            value={field.value}
            onChange={field.onChange}
            batchMode={form.getValueIn('batchMode')}
            withDescription
            isRootNameDisabled={true}
            allowAppendRootData={false}
            titleTooltip={I18n.t('node_http_response_data')}
            disabledTypes={[]}
            errors={fieldState?.errors}
          />
        )}
      </Field>
  </NodeConfigForm>
);
