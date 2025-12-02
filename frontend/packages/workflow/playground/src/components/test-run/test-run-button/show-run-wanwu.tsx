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

import React, { useEffect, useState } from 'react';
import { I18n } from '@coze-arch/i18n';
import { IconCozLoading } from '@coze-arch/coze-design/icons';
import { workflowApi } from '@coze-arch/bot-api';
import { type ButtonProps, Empty, Toast } from '@coze-arch/coze-design';
import { IllustrationNoContent } from "@douyinfe/semi-illustrations";
import { useTestRunFlowV2Wanwu } from '../hooks/use-test-run-flow-v2-wanwu';
import { ResultLog } from '../test-form-sheet-v2/result-log';
import styles from './show-run-wanwu.module.less'

type StartTestRunButtonProps = Pick<ButtonProps, 'size'>;

export const ShowRunWanwu: React.FC<StartTestRunButtonProps> = props => {
  const [resultData, setResultData] = useState<any>('')
  const [loading, setLoading] = useState<any>(false)
  const { testRunFlow, isChatFlow } = useTestRunFlowV2Wanwu({
    onSubmit: async (params) => {
      try {
        setLoading(true)
        const res:any = await workflowApi.WorkFlowRunWanwu(params);
        const { data } = res || {}
        const resData = {
          NodeType: "End",
          input: data,
          output: data,
          raw_output: data,
          extra: "{\"response_extra\":{\"terminal_plan\":2}}"
        }
        setResultData(resData)
      } catch (err) {
        const { statusText, data } = err?.response || {}
        Toast.error(data?.msg || statusText || 'Server Error');
      } finally {
        setLoading(false)
      }
    }
  });

  const runFlow = async () => {
    await testRunFlow();
  }

  useEffect(() => {
    runFlow()
  }, []);

  return isChatFlow ? null : (
    <div className={styles['show-run-wanwu-result']}>
      {loading ? (
        <div className={styles['result-loading']}>
          <IconCozLoading className="animate-spin coz-fg-dim mb-[4px] text-[32px]"/>
        </div>
      ) : (
        resultData ? (
          <ResultLog result={resultData} />
        ) : (
          <div className={styles['result-loading']}>
            <Empty
              image={
                <IllustrationNoContent style={{width: 112, height: 112}}/>
              }
              description={I18n.t('workflow_detail_title_testrun_desc')}
            />
          </div>
        )
      )}
    </div>
  );
};
