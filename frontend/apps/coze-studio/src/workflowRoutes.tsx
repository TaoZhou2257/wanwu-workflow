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

/*
 * 只包含 workflow 的路由
 */

import { createBrowserRouter } from 'react-router-dom';
import { lazy } from 'react';
import { Layout } from './layout';
import { GlobalError } from '@coze-foundation/layout';


const LoginPage = lazy(() =>
    import('@coze-foundation/account-ui-adapter').then(res => ({
      default: res.LoginPage,
    })),
);
const WorkflowPage = lazy(() =>
  import('@coze-workflow/playground-adapter').then(res => ({
    default: res.WorkflowPage,
  })),
);

export const workflowRouter: ReturnType<typeof createBrowserRouter> =
  createBrowserRouter([
    {
      path: '/',
      // Component: Layout,
      // errorElement: <GlobalError />,
      children: [
        /*{
          path: 'sign',
          Component: LoginPage,
          // errorElement: <GlobalError />,
          loader: () => ({
            hasSider: false,
            requireAuth: false,
          }),
        },*/
        {
          path: 'workflow',
          Component: WorkflowPage,
          loader: () => ({
            hasSider: false,
            requireAuth: false,
          }),
        },
      ],
    },
  ]);
