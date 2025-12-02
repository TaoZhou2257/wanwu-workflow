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

import React, { useMemo } from 'react';
import { LayoutPanelKey } from '@/constants';

import {
  StartTestFormSheetWanwu,
  type TestWorkflowFormPanelPropsWanwu,
} from '../test-run/test-form-sheet-v2/index-wanwu';
import {
  ChatFlowTestFormPanel,
  type ChatFlowTestFormPanelProps,
} from '../test-run/chat-flow-test-form-panel';
import { FloatLayoutWanwu, type FloatLayoutPropsWanwu } from '../float-layout';

export const WorkflowFloatLayout: React.FC<
  React.PropsWithChildren<FloatLayoutPropsWanwu>
> = ({ components, children }) => {
  const registry = useMemo(
    () => ({
      ...components,
      [LayoutPanelKey.TestFlowForm]: (p: TestWorkflowFormPanelPropsWanwu) => (
        <StartTestFormSheetWanwu {...p} />
      ),
      [LayoutPanelKey.TestChatFlowForm]: (p: ChatFlowTestFormPanelProps) => (
        <ChatFlowTestFormPanel {...p} />
      ),
    }),
    [components],
  );

  return <FloatLayoutWanwu components={registry}>{children}</FloatLayoutWanwu>;
};
