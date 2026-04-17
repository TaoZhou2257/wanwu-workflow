import {
  ValidateTrigger,
  type FormMetaV2,
  type Validate,
} from '@flowgram-adapter/free-layout-editor';

import { createValueExpressionInputValidate } from '@/node-registries/common/validators';
import {
  fireNodeTitleChange,
  provideNodeOutputVariablesEffect,
} from '@/node-registries/common/effects';

import { type FormData } from './types';
import { FormRender } from './form';
import { transformOnInit, transformOnSubmit } from './data-transformer';
import { I18n } from '@coze-arch/i18n';

export const SKILL_FORM_META: FormMetaV2<FormData> = {
  // 节点表单渲染
  render: () => <FormRender />,

  // 验证触发时机
  validateTrigger: ValidateTrigger.onChange,

  // 验证规则
  validate: {
    // 必填
    'inputs.inputParameters.0.input': createValueExpressionInputValidate({
      required: true,
    }),
    'inputs.llmParam': (({ value }) => (value ? undefined : I18n.t('model_placeholder' as any))) as Validate,
  },

  // 副作用管理
  effect: {
    nodeMeta: fireNodeTitleChange,
    outputs: provideNodeOutputVariablesEffect,
  },

  // 节点后端数据 -> 前端表单数据
  formatOnInit: transformOnInit,

  // 前端表单数据 -> 节点后端数据
  formatOnSubmit: transformOnSubmit,
};
