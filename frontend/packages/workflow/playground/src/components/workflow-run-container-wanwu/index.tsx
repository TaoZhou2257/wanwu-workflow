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

import { forwardRef, useMemo } from 'react';
import classnames from 'classnames';
import { QueryClientProvider } from '@tanstack/react-query';
import { workflowQueryClient } from '@coze-workflow/base';
import { WorkflowMode } from '@coze-workflow/base/api';
import { Spin } from '@coze-arch/bot-semi';
import { CustomError } from '@coze-arch/bot-error';
import { type BotSpace } from '@coze-arch/bot-api/developer_api';
import { RunContainerWanwu } from '../toolbar';
import {
  type WorkflowPlaygroundProps,
  type WorkflowPlaygroundRef,
} from '../../typing';
import { useGlobalState } from '../../hooks';
import { WorkflowFloatLayout } from './workflow-float-layout';

import styles from './index.module.less';

/**
 * Process Canvas
 */
const WorkflowRunContainerWanwu = forwardRef<
  WorkflowPlaygroundRef,
  WorkflowPlaygroundProps & {
    spaceList: BotSpace[];
  }
>((props, ref) => {
  const workflowState = useGlobalState();
  const { loading, loadingError, info } = workflowState;
  let playgroundContent;

  // Synchronize component properties to globalStatus
  useMemo(() => {
    const { spaceList, ...playgroundProps } = props;

    workflowState.updateConfig({
      playgroundProps,
      spaceList,
    });
  }, [props]);

  if (loading) {
    playgroundContent = (
      <Spin spinning={true} style={{ height: '100%', width: '100%' }} />
    );
  } else if (loadingError) {
    // Trigger exception, go to the top error boundary fallback
    throw new CustomError('normal_error', loadingError);
  } else {
    playgroundContent = (
      <QueryClientProvider client={workflowQueryClient}>
        <div className="flex flex-1 h-full">
          <div className="flex flex-1 flex-col">
            <div className={`${styles.workflowContent} clean-code`}>
              <WorkflowFloatLayout components={{}} isChatflow={ info.flow_mode === WorkflowMode.ChatFlow }>
                <RunContainerWanwu />
              </WorkflowFloatLayout>
            </div>
          </div>
        </div>
      </QueryClientProvider>
    );
  }

  return (
    <>
      <div
        className={classnames({
          [styles.workflowContainer]: true,
          [styles.workflowContainerOp]: IS_BOT_OP,
          [props.className || '']: props.className,
        })}
        style={props.style}
      >
        {playgroundContent}
      </div>
    </>
  );
});

export default WorkflowRunContainerWanwu;
