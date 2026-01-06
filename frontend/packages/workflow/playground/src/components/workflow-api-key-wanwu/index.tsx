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

import React, { useEffect, useState } from "react";
import { I18n } from '@coze-arch/i18n';
import { Button } from '@coze-arch/coze-design';
import { WorkflowMode, workflowApi } from '@coze-workflow/base/api';
import { useGlobalState } from '../../hooks';

import css from './index.module.less';

export const WorkflowApiKeyWanwu = () => {
  const globalState = useGlobalState();
  const [apiRootUrl, setApiRootUrl] = useState<string>('');

  const getApiRootUrl = async () => {
    const res = await workflowApi.GetWorkflowApiRootUrlWanwu({
      appId: globalState.workflowId,
      appType: globalState.flowMode === WorkflowMode.ChatFlow ? 'chatflow' : 'workflow'
    })
    setApiRootUrl(res?.data || '')
  }

  useEffect(() => {
    getApiRootUrl()
  }, []);

  return (
    <div className={css['workflow-api-key']}>
      <div className={css['api-key-content']}>
        <div className={css['api-root-url-title']}>
          {I18n.t('workflow_api_root_url_wanwu')}
        </div>
        <div>
          {apiRootUrl}
        </div>
      </div>
      <Button
        size="default"
        color="highlight"
        onClick={() => {
          window.open('/aibase/openApiKey')
        }}
      >
        {I18n.t('workflow_api_key_wanwu')}
      </Button>
    </div>
  );
};
