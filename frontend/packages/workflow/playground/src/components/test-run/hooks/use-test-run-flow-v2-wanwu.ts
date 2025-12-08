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

import { useCallback } from 'react';
import { WorkflowMode } from '@coze-workflow/base/api';
import { useService } from '@flowgram-adapter/free-layout-editor';
import { WorkflowRunService } from '@/services';
import { useFloatLayoutService } from '@/hooks';
import { LayoutPanelKey } from '@/constants';
import { useGetStartNode } from './use-get-start-node';

export const useTestRunFlowV2Wanwu = ({onSubmit}) => {
  const runService = useService<WorkflowRunService>(WorkflowRunService);
  const floatLayoutService = useFloatLayoutService();
  const { getNode } = useGetStartNode();

  const isChatFlow = () => runService.globalState.flowMode === WorkflowMode.ChatFlow

  const testRunFlow = useCallback(async () => {
    const node = getNode();
    if (!node) {
      return;
    }
    if (isChatFlow()) {
      floatLayoutService.open(LayoutPanelKey.TestChatFlowForm, 'right', {
        node,
      });
      return;
    }
    // The form can be opened first.
    floatLayoutService.open(LayoutPanelKey.TestFlowForm, 'right', { node, onSubmit });
  }, [runService, floatLayoutService]);

  return {
    testRunFlow,
    isChatFlow: isChatFlow()
  };
};
