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

/* eslint-disable @typescript-eslint/no-explicit-any -- history code */
/**
 * Adapt coze graph 2.0 test form
 */

import React, { useRef } from 'react';

import cls from 'classnames';
import { TestFormFieldName } from '@coze-workflow/test-run-next';
import {
  FormItemSchemaType,
  FormPanelLayout,
  TestsetManageProvider,
} from '@coze-workflow/test-run';
import { userStoreService } from '@coze-studio/user-store';
import { I18n } from '@coze-arch/i18n';
import { Toast } from '@coze-arch/coze-design';
import { type WorkflowNodeEntity } from '@/test-run-kit';
import { useGlobalState } from '@/hooks';

import { stringifyValue } from '../utils/stringify-value';
import { TestFormV3Wanwu } from '../test-form-v3';
import { JsonEditorSemi } from '../test-form-materials/json-editor';
import { InputForm } from '../input-form';
import { useGetStartNode } from '../hooks/use-get-start-node';
import { TestsetBotProjectSelect } from '../chat-flow-test-form-panel/testset-bot-project-select';
import { TestFormSheetHeaderWanwu } from './header-wanwu';
import { TestFormSheetFooterV2Wanwu } from './footer-wanwu';

import styles from './index-wanwu.module.less';

interface TestWorkflowFormPanelPropsWanwu {
  node: WorkflowNodeEntity;
  onSubmit?: any;
}

const TestFormSheetV2Wanwu: React.FC<TestWorkflowFormPanelPropsWanwu> = ({ node, onSubmit }) => {
  const formApiRef = useRef<any>(null);
  const globalState = useGlobalState();
  const userInfo = userStoreService.useUserInfo();
  const { getNode } = useGetStartNode();

  const testRunFlowV3 = async () => {
    if (!formApiRef.current) {
      return;
    }
    const {
      empty,
      validate: formValidate,
      values,
    } = await formApiRef.current.submit();
    if (!formValidate) {
      Toast.error(I18n.t('workflow_testrun_form_vailate_error_toast'));
      return;
    }
    let inputData: object | undefined;
    let input: Record<string, string> | undefined;
    if (!empty) {
      inputData = values?.[TestFormFieldName.Node]?.[TestFormFieldName.Input];
      input = stringifyValue(inputData as object);
    }

    onSubmit?.({input, workflow_id: globalState.workflowId})
  };

  const handleSubmit = async () => {
    await testRunFlowV3();
  };

  return (
    <div className={styles['test-form-v2']}>
      <TestsetManageProvider
        spaceId={globalState.spaceId}
        workflowId={globalState.workflowId}
        userId={userInfo?.user_id_str}
        nodeId={getNode()?.id}
        projectId={globalState.projectId}
        formRenders={{
          [FormItemSchemaType.BOT]: TestsetBotProjectSelect as any,
          [FormItemSchemaType.LIST]: JsonEditorSemi as any,
          [FormItemSchemaType.OBJECT]: JsonEditorSemi as any,
        }}
      >
        <TestFormSheetHeaderWanwu />

        <div className={styles['test-form-content']}>
          <div
            className={cls('w-full h-full')}
          >
            <TestFormV3Wanwu node={node} onMounted={v => (formApiRef.current = v)} />
          </div>
        </div>
        <TestFormSheetFooterV2Wanwu onClick={handleSubmit} />
        <InputForm />
      </TestsetManageProvider>
    </div>
  );
};

const StartTestFormSheetWanwu: React.FC<TestWorkflowFormPanelPropsWanwu> = props => (
  <div className={styles['test-form-v2-wrapper-wanwu']}>
    <FormPanelLayout>
      <TestFormSheetV2Wanwu {...props} />
    </FormPanelLayout>
  </div>
);

export { TestFormSheetV2Wanwu, StartTestFormSheetWanwu, type TestWorkflowFormPanelPropsWanwu };
