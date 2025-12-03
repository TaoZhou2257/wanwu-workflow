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

/**
 * Layout floating layer covering the entire canvas, carrying some linked suspension components
 */
import React, { useEffect, useLayoutEffect, useRef } from 'react';
import { useSize } from 'ahooks';
import cs from 'classnames';
import { type Render } from '@/services/workflow-float-layout-service';
import { useFloatLayoutService } from '@/hooks/use-float-layout-service';
import { FloatPanel } from './float-panel';

import styles from './float-layout-wanwu.module.less';

export interface FloatLayoutPropsWanwu {
  components: Record<string, Render>;
  isChatflow?: boolean;
}

export const FloatLayoutWanwu: React.FC<
  React.PropsWithChildren<FloatLayoutPropsWanwu>
> = ({ components, children, isChatflow }) => {
  const ref = useRef<HTMLDivElement>(null);
  const floatLayoutService = useFloatLayoutService();
  const size = useSize(ref);

  useLayoutEffect(() => {
    // Only trigger once
    floatLayoutService.register(components);
  }, []);

  useEffect(() => {
    if (size) {
      floatLayoutService.setLayoutSize(size);
    }
  }, [size, floatLayoutService]);

  return (
    <div className={styles['float-layout-wanwu']} ref={ref}>
      <div className={cs(styles['left-panel'], isChatflow ? styles['panel-content'] : null)}>
        <FloatPanel panel={floatLayoutService.right}/>
      </div>
      <div className={styles['right-panel']}>
        <div className={styles['left-main-panel']}>{children}</div>
      </div>
    </div>
  );
};
